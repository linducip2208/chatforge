package wa

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"chatgo/meta"
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
	go e.dripLoop()
	go e.reminderLoop()
	go e.recurringLoop()
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
	// check if paused before starting
	status := c.Status
	if status == "paused" { return }
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

	// resume: skip already-sent numbers
	alreadySent := map[string]bool{}
	if c.SentTo != "" {
		for _, p := range strings.Split(c.SentTo, ",") {
			p = strings.TrimSpace(p)
			if p != "" { alreadySent[p] = true }
		}
	}

	// gather unique recipients from groups + direct numbers
	seen := map[string]bool{}
	var targets []store.Contact
	for _, gid := range strings.Split(c.Groups, ",") {
		gid = strings.TrimSpace(gid)
		if gid == "" { continue }
		list, _ := e.db.ContactsByGroup(gid)
		for _, ct := range list {
			if !seen[ct.Phone] && ct.Phone != "" && !alreadySent[ct.Phone] {
				seen[ct.Phone] = true
				targets = append(targets, ct)
			}
		}
	}
	// direct numbers (comma-separated)
	if c.Numbers != "" {
		for _, n := range strings.Split(c.Numbers, ",") {
			n = strings.TrimSpace(n)
			if n == "" || seen[n] || alreadySent[n] { continue }
			seen[n] = true
			targets = append(targets, store.Contact{Phone: n})
		}
	}
	// tag filter: if tags specified, intersect with contacts from those tags
	if c.Tags != "" {
		tagSet := map[string]bool{}
		for _, tid := range strings.Split(c.Tags, ",") {
			tid = strings.TrimSpace(tid)
			if tid == "" { continue }
			tagID, _ := strconv.ParseInt(tid, 10, 64)
			list, _ := e.db.ContactsByTag(tagID)
			for _, ct := range list {
				tagSet[ct.Phone] = true
			}
		}
		// filter targets: only keep those in tagSet
		var filtered []store.Contact
		for _, ct := range targets {
			if tagSet[ct.Phone] {
				filtered = append(filtered, ct)
			}
		}
		targets = filtered
	}
	// determine Meta client (if selected)
	var metaClient *meta.Client
	if c.MetaAccountID > 0 {
		acc, err := e.db.GetMetaAccount(c.MetaAccountID)
		if err == nil {
			metaClient = meta.New(acc.PhoneNumberID, acc.AccessToken, acc.VerifyToken)
		}
	}
	// rate limit config
	maxDaily, _ := strconv.Atoi(e.db.GetSetting("rate_max_daily", "0"))
	rndMin, _ := strconv.Atoi(e.db.GetSetting("rate_random_min", "0"))
	rndMax, _ := strconv.Atoi(e.db.GetSetting("rate_random_max", "0"))
	if rndMax <= rndMin { rndMax = rndMin + 5 }

	for _, ct := range targets {
		if e.campaignTerminated(c.ID) { return }
		if e.db.IsUnsub(ct.Phone) || e.db.IsBlacklisted(ct.Phone) { continue }
		// validate number format
		if !store.ValidFormat(ct.Phone) { continue }
		// rate limit
		if maxDaily > 0 {
			todaySent := e.db.TodaySentCount("")
			if todaySent >= maxDaily {
				e.db.Log("campaign", "rate_limit", "Daily limit reached: "+itoa(int64(maxDaily)))
				return
			}
		}
		var sendErr error
		msg := msgtemplate.Render(c.Message, msgtemplate.Vars{Name: ct.Name, Phone: ct.Phone, Message: ""})
		// A/B test: check if there's a test for this campaign, alternate variants
		if ab, abErr := e.db.GetABTest(c.ID); abErr == nil {
			if ab.ASent+ab.BSent < ab.ASent*2+10 {
				var v string
				if ab.ASent <= ab.BSent { msg = ab.VariantA; v = "a" } else { msg = ab.VariantB; v = "b" }
				msg = msgtemplate.Render(msg, msgtemplate.Vars{Name: ct.Name, Phone: ct.Phone, Message: ""})
				e.db.IncABSent(c.ID, v)
			}
		}
		// link tracking: replace URLs with tracking links
		msg = replaceTrackURLs(e, msg, c.ID, ct.Phone)
		if metaClient != nil {
			// Meta API sending
			if c.MediaURL != "" && c.MediaType == "image" {
				_, sendErr = metaClient.SendImage(onlyDigits(ct.Phone), c.MediaURL, msg)
			} else if c.MediaURL != "" && c.MediaType == "document" {
				_, sendErr = metaClient.SendDocument(onlyDigits(ct.Phone), c.MediaURL, "document", msg)
			} else if c.MetaTemplate != "" {
				lang := "id"
				_, sendErr = metaClient.SendTemplate(onlyDigits(ct.Phone), c.MetaTemplate, lang, []string{msg})
			} else {
				_, sendErr = metaClient.SendText(onlyDigits(ct.Phone), msg)
			}
			if sendErr == nil { _ = e.db.IncCampaignSent(c.ID) }
			e.db.AppendCampaignSentTo(c.ID, ct.Phone)
		} else {
			st, _ := e.Status()
			if st != "connected" { return }
			sendSession := sel.Next(e)
			if sendSession == nil { continue }
			digits := onlyDigits(ct.Phone)
			if digits != "" {
				if c.MediaURL != "" && c.MediaType == "image" {
					sendErr = e.SendMedia(sendSession.Phone, ct.Phone, "image", c.MediaURL, msg)
				} else if c.MediaURL != "" && c.MediaType == "document" {
					sendErr = e.SendMedia(sendSession.Phone, ct.Phone, "document", c.MediaURL, msg)
				} else {
					jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
					sendErr = e.sendVia(sendSession, jid, msg)
				}
				if sendErr == nil { _ = e.db.IncCampaignSent(c.ID) }
			}
			_ = e.db.AppendCampaignSentTo(c.ID, ct.Phone)
		}
		interval := c.Interval
		if interval <= 0 { interval = 300 }
		// random delay
		if rndMax > 0 {
			randDelay := rndMin + int(time.Now().UnixNano())%(rndMax-rndMin)
			interval = randDelay
			if interval < 30 { interval = 30 }
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	_ = e.db.UpdateCampaignStatus(c.ID, "done")
	e.db.Log("campaign", "done", "Campaign #"+itoa(c.ID)+" ("+c.Name+") finished")
}

func (e *Engine) campaignTerminated(id int64) bool {
	for _, c := range mustCampaigns(e) {
		if c.ID == id {
			return c.Status == "stopped" || c.Status == "paused"
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
					e.db.LogInstance(s.Phone, "heartbeat_fail")
					retry := retryCounts[s.id]
					retryCounts[s.id] = retry + 1
					// exponential backoff 5s → 10s → 20s → ... max 5min
					backoffSec := 5 << minInt(retry, 6) // 5,10,20,40,80,160,~320
					if backoffSec > 300 { backoffSec = 300 }
					if retry > 10 {
						e.db.LogInstance(s.Phone, "max_retries")
						continue
					}
					time.Sleep(time.Duration(backoffSec) * time.Second)
					_ = s.client.Connect()
				}
			}
		}
	}
}

// dripLoop sends drip campaign messages on schedule.
func (e *Engine) dripLoop() {
	for {
		time.Sleep(60 * time.Second)
		enrollments, err := e.db.DueDripEnrollments()
		if err != nil {
			continue
		}
		for _, en := range enrollments {
			drip, err := e.db.GetDrip(en.DripID)
			if err != nil || len(drip.Steps) == 0 {
				continue
			}
			if en.CurrentStep >= len(drip.Steps) {
				e.db.AdvanceDripStep(en.ID, 0) // mark completed
				continue
			}
			step := drip.Steps[en.CurrentStep]
			// send via first connected WA session
			e.mu.RLock()
			var s *session
			for _, ss := range e.sessions {
				if ss.status == "connected" { s = ss; break }
			}
			e.mu.RUnlock()
			if s == nil { continue }
			msg := msgtemplate.Render(step.Message, msgtemplate.Vars{Name: en.Name, Phone: en.Phone, Message: ""})
			digits := onlyDigits(en.Phone)
			if digits != "" {
				jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
				if err := e.sendVia(s, jid, msg); err == nil {
					e.db.Log("drip", "sent", "Drip #"+itoa(en.DripID)+" step "+itoa(int64(en.CurrentStep+1))+" -> "+en.Phone)
				}
			}
			nextDelay := 0
			if en.CurrentStep+1 < len(drip.Steps) {
				nextDelay = drip.Steps[en.CurrentStep+1].DelayMinutes
			}
			e.db.AdvanceDripStep(en.ID, nextDelay)
			time.Sleep(2 * time.Second)
		}
	}
}

// reminderLoop checks for due payment reminders and sends them.
func (e *Engine) reminderLoop() {
	for {
		time.Sleep(4 * time.Hour)
		reminders, err := e.db.DueReminders()
		if err != nil { continue }
		e.mu.RLock()
		var s *session
		for _, ss := range e.sessions {
			if ss.status == "connected" { s = ss; break }
		}
		e.mu.RUnlock()
		if s == nil { continue }
		for _, r := range reminders {
			msg := strings.ReplaceAll(r.Message, "{amount}", fmt.Sprintf("%.0f", r.Amount))
			msg = strings.ReplaceAll(msg, "{date}", r.DueDate)
			digits := onlyDigits(r.Phone)
			if digits != "" {
				jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
				if err := e.sendVia(s, jid, msg); err == nil {
					e.db.MarkReminderSent(r.ID)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
}

// replaceTrackURLs finds URLs in message, creates tracking links, and replaces them.
var urlRe = regexp.MustCompile(`https?://[^\s"]+`)

func replaceTrackURLs(e *Engine, msg string, campaignID int64, phone string) string {
	matches := urlRe.FindAllString(msg, -1)
	for _, url := range matches {
		tok := make([]byte, 6)
		rand.Read(tok)
		token := hex.EncodeToString(tok)
		e.db.TrackLink(token, url, campaignID, phone)
		appURL := e.db.GetSetting("app_url", "http://localhost:8080")
		msg = strings.Replace(msg, url, appURL+"/track/"+token, 1)
	}
	return msg
}

func (e *Engine) recurringLoop() {
	for {
		time.Sleep(30 * time.Minute)
		camps, err := e.db.DueRecurring()
		if err != nil { continue }
		e.mu.RLock()
		var s *session
		for _, ss := range e.sessions {
			if ss.status == "connected" { s = ss; break }
		}
		e.mu.RUnlock()
		if s == nil { continue }
		for _, c := range camps {
			seen := map[string]bool{}
			for _, gid := range strings.Split(c.Groups, ",") {
				gid = strings.TrimSpace(gid)
				if gid == "" { continue }
				list, _ := e.db.ContactsByGroup(gid)
				for _, ct := range list {
					if !seen[ct.Phone] && ct.Phone != "" { seen[ct.Phone] = true }
				}
			}
			for phone := range seen {
				msg := msgtemplate.Render(c.Message, msgtemplate.Vars{Phone: phone, Message: ""})
				digits := onlyDigits(phone)
				if digits != "" {
					jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
					e.sendVia(s, jid, msg)
					time.Sleep(3 * time.Second)
				}
			}
			e.db.MarkRecurringRun(c.ID)
			e.db.Log("recurring", "run", fmt.Sprintf("%s sent to %d contacts", c.Name, len(seen)))
		}
	}
}
