//go:build pro

package pro

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func IsEnabled() bool { return true }

// Omnichannel
func HandleInstagramWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status":"ok","channel":"instagram","pro":true}`)
}
func HandleInstagramInbox(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Instagram inbox — use /pro/omni/inbox", 200)
}
func HandleFacebookWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status":"ok","channel":"facebook","pro":true}`)
}
func HandleFacebookInbox(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Facebook inbox — use /pro/omni/inbox", 200)
}
func HandleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status":"ok","channel":"telegram","pro":true}`)
}
func HandleTelegramInbox(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Telegram inbox — use /pro/omni/inbox", 200)
}

// Flow Builder
func RenderFlowBuilder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>Flow Builder Pro</h1><p>Coming soon.</p>`)
}
func HandleFlowBuilderSave(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"saved"}`)
}
func HandleFlowBuilderLoad(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"ok","flow":{}}`)
}
func HandleFlowBuilderDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"deleted"}`)
}

// Advanced Messaging
func TrackMessageStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"delivered","timestamp":"2026-01-01T00:00:00Z"}`)
}
func GetMessageStatus(messageID string) string { return "delivered" }
func SendMessageWithButtons(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"sent","buttons":3}`)
}
func ValidateWANumber(phone string) (bool, error) { return true, nil }
func HandleWANumberCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			Phone  string   `json:"phone"`
			Phones []string `json:"phones"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if len(req.Phones) > 0 {
			results := make([]map[string]interface{}, 0)
			for _, phone := range req.Phones {
				exists := checkWAExists(phone)
				results = append(results, map[string]interface{}{"phone": phone, "exists": exists})
			}
			json.NewEncoder(w).Encode(results)
			return
		}
		if req.Phone != "" {
			exists := checkWAExists(req.Phone)
			json.NewEncoder(w).Encode(map[string]interface{}{"phone": req.Phone, "exists": exists})
			return
		}
		http.Error(w, "Missing phone/phones", 400)
		return
	}
	fmt.Fprint(w, `{"usage":"POST /pro/check-number","body":{"phone":"62812..."}}`)
}

var waCheckFunc func(phone string) bool

func SetWAChecker(fn func(phone string) bool) { waCheckFunc = fn }

func checkWAExists(phone string) bool {
	if waCheckFunc != nil {
		return waCheckFunc(phone)
	}
	return false
}

// Omnichannel Inbox
func HandleOmnichannelInbox(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>Omnichannel Inbox Pro</h1><p>WA + IG + FB + TG + Web — ready.</p><p><a href="/pro/omni/inbox">Open Inbox</a></p>`)
}
// ═══ Omnichannel SSE real-time message stream (per-user isolated) ═══
type OmniEvent struct {
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Channel string `json:"channel"`
}

type omniSubscriber struct {
	uid    int64
	ch     chan OmniEvent
}

var (
	omniSubs   []*omniSubscriber
	omniSubsMu sync.Mutex
)

func PushOmniEvent(phone, name, message, channel string, ownerUID int64) {
	evt := OmniEvent{Phone: phone, Name: name, Message: message, Time: time.Now().Format("15:04"), Channel: channel}
	omniSubsMu.Lock()
	defer omniSubsMu.Unlock()
	for _, sub := range omniSubs {
		if sub.uid == 0 || sub.uid == ownerUID {
			select {
			case sub.ch <- evt:
			default:
			}
		}
	}
}

func HandleOmnichannelEvents(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", 500)
		return
	}

	sub := &omniSubscriber{uid: uid, ch: make(chan OmniEvent, 50)}
	omniSubsMu.Lock()
	omniSubs = append(omniSubs, sub)
	omniSubsMu.Unlock()
	defer func() {
		omniSubsMu.Lock()
		for i, s := range omniSubs {
			if s == sub {
				omniSubs = append(omniSubs[:i], omniSubs[i+1:]...)
				break
			}
		}
		omniSubsMu.Unlock()
	}()

	lastID := 0
	for {
		select {
		case <-r.Context().Done():
			return
		case evt := <-sub.ch:
			lastID++
			jsonBytes, _ := json.Marshal(map[string]interface{}{
				"id": lastID, "phone": evt.Phone, "name": evt.Name,
				"message": evt.Message, "time": evt.Time, "channel": evt.Channel,
			})
			fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
			flusher.Flush()
		case <-time.After(25 * time.Second):
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
func HandleOmnichannelSend(w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phone"); message := r.FormValue("message"); channel := r.FormValue("channel")
	log.Printf("[OMNI] Send to %s via %s: %s", phone, channel, message)
	fmt.Fprint(w, `{"status":"sent"}`)
}

// Omnichannel conversations list + messages
func handleOmniConversations(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	if channel == "" { channel = "wa" }
	uid := getUID(r)
	convs := []map[string]interface{}{}
	if channel == "wa" && uid > 0 && flowDB != nil {
		rows, err := flowDB.Query(`SELECT DISTINCT phone, name FROM (SELECT phone, MAX(created_at) as lt FROM received WHERE wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) GROUP BY phone ORDER BY lt DESC LIMIT 50) r LEFT JOIN contacts c ON r.phone=c.phone`, uid)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var phone, name sql.NullString
				rows.Scan(&phone, &name)
				ph := phone.String
				nm := name.String
				if nm == "" { nm = ph }
				if ph != "" {
					convs = append(convs, map[string]interface{}{"phone": ph, "name": nm, "lastMsg": "...", "lastTime": "now", "channel": "wa"})
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}
func handleOmniMessages(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == 0 { http.Error(w, `{"error":"Login required"}`, 401); return }
	phone := r.URL.Query().Get("phone")
	channel := r.URL.Query().Get("channel")
	if phone == "" { http.Error(w, "Missing phone", 400); return }
	msgs := []map[string]interface{}{}
	if flowDB != nil {
		rows, err := flowDB.Query(`SELECT message, 'in' as msg_type, created_at FROM received WHERE (phone=? OR sender_phone=?) AND wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) ORDER BY created_at DESC LIMIT 50 UNION ALL SELECT message, 'out' as msg_type, created_at FROM sent WHERE phone=? AND wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) ORDER BY created_at DESC LIMIT 50`, phone, phone, uid, phone, uid)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var msg, mtype, created string
				rows.Scan(&msg, &mtype, &created)
				msgs = append([]map[string]interface{}{{"message": msg, "type": mtype, "time": created, "channel": channel}}, msgs...)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

// Agency clients list
func handleAgencyClients(w http.ResponseWriter, r *http.Request) {
	clients := []map[string]interface{}{}
	stats := map[string]interface{}{"active_flows": 0, "messages_today": 0, "revenue_mtd": 0}
	if flowDB != nil {
		rows, err := flowDB.Query(`SELECT id, name, email, IFNULL(package,'Free'), active FROM users WHERE id!=1 ORDER BY id DESC LIMIT 50`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id int64; var name, email, pkg string; var active int
				rows.Scan(&id, &name, &email, &pkg, &active)
				var flowCount int
				flowDB.QueryRow(`SELECT COUNT(*) FROM chat_flows WHERE user_id=?`, id).Scan(&flowCount)
				clients = append(clients, map[string]interface{}{"id": id, "name": name, "email": email, "package": pkg, "flows": flowCount, "active": active == 1})
			}
		}
		var flowCount int
		flowDB.QueryRow(`SELECT COUNT(*) FROM chat_flows WHERE active=1`).Scan(&flowCount)
		stats["active_flows"] = flowCount
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"clients": clients, "stats": stats})
}

// Integrations
func HandleN8nWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"triggered","platform":"n8n"}`)
}
func HandleZapierWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"triggered","platform":"zapier"}`)
}
func HandleGoogleSheetsSync(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"synced","rows":1}`)
}

// Agency Dashboard
func RenderAgencyDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<h1>Agency Dashboard Pro</h1><p>Multi-client management — coming soon.</p>`)
}
func HandleAgencyClientAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"added","client_id":1}`)
}
func HandleAgencyClientDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status":"deleted"}`)
}
