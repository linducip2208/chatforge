package wa

import (
	"strings"
	"time"

	"chatgo/msgtemplate"
	"chatgo/store"

	waTypes "go.mau.fi/whatsmeow/types"
)

// StartLoops runs background processors for campaigns and scheduled messages.
// This is an ADDITIVE feature layer — it does not touch auto-reply matching.
func (e *Engine) StartLoops() {
	go e.campaignLoop()
	go e.scheduledLoop()
	go e.heartbeatLoop()
}

// campaignLoop drains running campaigns, sending to each contact in the target groups.
func (e *Engine) campaignLoop() {
	for {
		time.Sleep(5 * time.Second)
		status, _ := e.Status()
		if status != "connected" {
			continue
		}
		camps, err := e.db.PendingCampaigns()
		if err != nil {
			continue
		}
		for _, c := range camps {
			if c.Status == "pending" {
				_ = e.db.UpdateCampaignStatus(c.ID, "running")
			}
			e.runCampaign(c)
		}
	}
}

func (e *Engine) runCampaign(c store.Campaign) {
	// Build sender selector from comma-separated account_ids
	var sel *SenderSelector
	sendMode := c.SendMode
	if sendMode == "" { sendMode = "round_robin" }
	if c.AccountIDs != "" {
		phones := strings.Split(c.AccountIDs, ",")
		var clean []string
		for _, p := range phones {
			p = strings.TrimSpace(p)
			if p != "" { clean = append(clean, strings.TrimPrefix(p, "+")) }
		}
		if len(clean) > 0 {
			sel = &SenderSelector{Phones: clean, Mode: sendMode}
		}
	}
	if sel == nil { sel = &SenderSelector{Mode: sendMode} }
	// gather unique recipients from groups + direct numbers
	seen := map[string]bool{}
	var targets []store.Contact
	for _, gid := range strings.Split(c.Groups, ",") {
		gid = strings.TrimSpace(gid)
		if gid == "" { continue }
		list, _ := e.db.ContactsByGroup(gid)
		for _, ct := range list {
			if !seen[ct.Phone] && ct.Phone != "" {
				seen[ct.Phone] = true
				targets = append(targets, ct)
			}
		}
	}
	// direct numbers (comma-separated)
	if c.Numbers != "" {
		for _, n := range strings.Split(c.Numbers, ",") {
			n = strings.TrimSpace(n)
			if n == "" || seen[n] { continue }
			seen[n] = true
			targets = append(targets, store.Contact{Phone: n})
		}
	}
	for _, ct := range targets {
		st, _ := e.Status()
		if st != "connected" {
			return
		}
		if e.campaignStopped(c.ID) {
			return
		}
		if e.db.IsUnsub(ct.Phone) {
			continue
		}
		sendSession := sel.Next(e)
		if sendSession == nil { continue }
		msg := msgtemplate.Render(c.Message, msgtemplate.Vars{Name: ct.Name, Phone: ct.Phone, Message: ""})
		digits := onlyDigits(ct.Phone)
		if digits != "" {
			jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
			if err := e.sendVia(sendSession, jid, msg); err == nil {
				_ = e.db.IncCampaignSent(c.ID)
			}
		}
		interval := c.Interval
		if interval <= 0 { interval = 300 }
		time.Sleep(time.Duration(interval) * time.Second)
	}
	_ = e.db.UpdateCampaignStatus(c.ID, "done")
	e.db.Log("campaign", "done", "Campaign #"+itoa(c.ID)+" ("+c.Name+") finished")
}

func (e *Engine) campaignStopped(id int64) bool {
	for _, c := range mustCampaigns(e) {
		if c.ID == id {
			return c.Status == "stopped"
		}
	}
	return false
}

func mustCampaigns(e *Engine) []store.Campaign {
	c, _ := e.db.ListCampaigns()
	return c
}

// scheduledLoop sends due scheduled messages.
func (e *Engine) scheduledLoop() {
	for {
		time.Sleep(20 * time.Second)
		status, _ := e.Status()
		if status != "connected" {
			continue
		}
		due, err := e.db.DueScheduled()
		if err != nil {
			continue
		}
		for _, s := range due {
			msg := msgtemplate.Render(s.Message, msgtemplate.Vars{Phone: s.Phone, Message: ""})
			// Use specific account if set, otherwise use Send (first connected)
			if s.AccountIDs != "" {
				phones := strings.Split(s.AccountIDs, ",")
				var p string
				for _, ph := range phones {
					if ph = strings.TrimSpace(strings.TrimPrefix(ph, "+")); ph != "" { p = ph; break }
				}
				_ = e.SendFrom(p, s.Phone, msg)
			} else {
				_ = e.Send(s.Phone, msg)
			}
			if s.Repeat > 0 {
				_ = e.db.RescheduleAfter(s.ID, s.Repeat) // keep pending, move next time
			} else {
				_ = e.db.MarkScheduledSent(s.ID)
			}
			e.db.Log("scheduled", "sent", "Scheduled #"+itoa(s.ID)+" -> "+s.Phone)
		}
	}
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func minInt(a, b int) int { if a < b { return a }; return b }

// heartbeatLoop checks connected sessions and attempts reconnect on stale ones.
func (e *Engine) heartbeatLoop() {
	retryCounts := map[string]int{}
	for {
		time.Sleep(60 * time.Second)
		e.mu.RLock()
		snap := make(map[string]*session, len(e.sessions))
		for k, v := range e.sessions {
			snap[k] = v
		}
		e.mu.RUnlock()
		for _, s := range snap {
			if s.status == "connected" && s.client != nil {
				if !s.client.IsConnected() {
					e.log.Warnf("heartbeat: %s lost, reconnecting (backoff)...", s.id)
					e.db.LogInstance(s.phone, "heartbeat_fail")
					retry := retryCounts[s.id]
					retryCounts[s.id] = retry + 1
					// exponential backoff 5s → 10s → 20s → ... max 5min
					backoffSec := 5 << minInt(retry, 6) // 5,10,20,40,80,160,~320
					if backoffSec > 300 { backoffSec = 300 }
					if retry > 10 {
						e.db.LogInstance(s.phone, "max_retries")
						continue
					}
					time.Sleep(time.Duration(backoffSec) * time.Second)
					_ = s.client.Connect()
				}
			}
		}
	}
}
