package wa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"chatgo/msgtemplate"
	"chatgo/aiservice"
	"chatgo/secret"
	"chatgo/store"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	waCompanionReg "go.mau.fi/whatsmeow/proto/waCompanionReg"
	wmstore "go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	_ "modernc.org/sqlite" // pure-Go SQLite driver for whatsmeow session store
)

// session is one WhatsApp account (one device).
type session struct {
	id        string
	client    *whatsmeow.Client
	status    string
	qr        string
	Phone     string
	createdAt time.Time
	userID    int64
}

// Engine manages MANY whatsmeow accounts (multi-number).
type Engine struct {
	container *sqlstore.Container
	db        *store.DB
	log       waLog.Logger

	mu       sync.RWMutex
	sessions map[string]*session // key = session id
	newSeq   int

	notifyCh chan string
}

// New opens the whatsmeow session store.
func New(sessionPath string, appDB *store.DB) (*Engine, error) {
	logger := waLog.Stdout("chatgo", "INFO", true)
	wmstore.DeviceProps.Os = proto.String("Chrome (Linux)")
	wmstore.DeviceProps.PlatformType = waCompanionReg.DeviceProps_CHROME.Enum()

	ctx := context.Background()
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", sessionPath)
	container, err := sqlstore.New(ctx, "sqlite", dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("open session store: %w", err)
	}
	return &Engine{
		container: container,
		db:        appDB,
		log:       logger,
		sessions:  map[string]*session{},
		notifyCh:  make(chan string, 256),
	}, nil
}

// Start loads all previously-paired devices and connects them.
func (e *Engine) Start() error {
	ctx := context.Background()
	devices, err := e.container.GetAllDevices(ctx)
	if err != nil {
		return fmt.Errorf("get devices: %w", err)
	}
	for _, dev := range devices {
		e.connectDevice(dev)
	}
	return nil
}

// AccountLimit returns how many accounts are allowed (from the biggest package;
// falls back to a default). This drives the "based on plan" rule on /wa.
func (e *Engine) AccountLimit(userID int64) int {
	return e.db.GetUserPackageLimit(userID)
}

func (e *Engine) UserAccountLimit(userID int64) int {
	return e.db.GetUserPackageLimit(userID)
}

// connectDevice wraps a store device in a live session and connects it.
func (e *Engine) connectDevice(dev *wmstore.Device) *session {
	client := whatsmeow.NewClient(dev, e.log)
	id := "unknown"
	var phone string
	if dev.ID != nil {
		id = dev.ID.String()
		phone = dev.ID.User
	}
	s := &session{id: id, client: client, status: "connecting", createdAt: time.Now()}
	if dev.ID != nil {
		s.Phone = phone
		s.userID = e.db.GetSessionOwner(phone)
	}
	e.mu.Lock()
	e.sessions[id] = s
	e.mu.Unlock()

	client.AddEventHandler(func(evt interface{}) { e.handleEvent(s, evt) })
	if err := client.Connect(); err != nil {
		s.status = "disconnected"
	}
	return s
}

// AddAccount starts a brand-new pairing session (new number) and returns its id.
func (e *Engine) AddAccount(userID int64) (string, error) {
	if e.CountAccounts(userID) >= e.AccountLimit(userID) {
		return "", fmt.Errorf("account limit reached (%d)", e.AccountLimit(userID))
	}
	// WA Server package restriction
	if !e.userAllowedOnWAServer(userID) {
		return "", fmt.Errorf("your package is not allowed on any WA server (upgrade required)")
	}
	ctx := context.Background()
	dev := e.container.NewDevice()
	client := whatsmeow.NewClient(dev, e.log)

	e.mu.Lock()
	e.newSeq++
	id := fmt.Sprintf("new:%d", e.newSeq)
	s := &session{id: id, client: client, status: "qr", createdAt: time.Now(), userID: userID}
	e.sessions[id] = s
	e.mu.Unlock()

	client.AddEventHandler(func(evt interface{}) { e.handleEvent(s, evt) })

	qrChan, _ := client.GetQRChannel(ctx)
	if err := client.Connect(); err != nil {
		return "", err
	}
	go func() {
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				e.mu.Lock()
				s.qr = evt.Code
				s.status = "qr"
				e.mu.Unlock()
			case "success":
				e.mu.Lock()
				s.qr = ""
				s.status = "connected"
				if client.Store.ID != nil {
					newID := client.Store.ID.String()
					s.Phone = client.Store.ID.User
					s.userID = userID
					delete(e.sessions, id)
					s.id = newID
					e.sessions[newID] = s
					e.db.SaveSessionOwner(s.Phone, userID)
				}
				e.mu.Unlock()
			case "timeout":
				e.mu.Lock()
				s.status = "disconnected"
				e.mu.Unlock()
			}
		}
	}()
	// Wait for the first QR code (up to 15s)
	deadline := time.Now().Add(15 * time.Second)
	for {
		e.mu.RLock()
		qr := s.qr
		e.mu.RUnlock()
		if qr != "" || s.status == "connected" {
			break
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	return id, nil
}

// Accounts returns a snapshot of all sessions for the UI.
type AccountInfo struct {
	ID     string
	Phone  string
	Status string
}

func (e *Engine) Accounts(userID int64) []AccountInfo {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []AccountInfo
	now := time.Now()
	for id, s := range e.sessions {
		if s.userID != 0 && s.userID != userID {
			continue
		}
		if s.Phone == "" && s.status != "connected" && s.createdAt.Add(5*time.Minute).Before(now) {
			continue
		}
		out = append(out, AccountInfo{ID: id, Phone: s.Phone, Status: s.status})
	}
	return out
}

func (e *Engine) CountAccounts(userID int64) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	n := 0
	for _, s := range e.sessions {
		if s.userID != 0 && s.userID != userID { continue }
		n++
	}
	return n
}

func (e *Engine) CountConnected(userID int64) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	n := 0
	for _, s := range e.sessions {
		if s.userID != 0 && s.userID != userID { continue }
		if s.status == "connected" { n++ }
	}
	return n
}
func (e *Engine) CountDisconnected(userID int64) int {
	return e.CountAccounts(userID) - e.CountConnected(userID)
}

// QRFor returns the latest QR for a pairing session id.
func (e *Engine) QRFor(id string) string {
	// wait up to 8s for the QR to become available (pairing may still be connecting)
	deadline := time.Now().Add(8 * time.Second)
	for {
		e.mu.RLock()
		s := e.sessions[id]
		e.mu.RUnlock()
		if s != nil && s.qr != "" {
			return s.qr
		}
		if s != nil && s.status == "connected" {
			return "" // already connected, no QR needed
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return ""
}

// StatusFor returns status+phone for a session.
func (e *Engine) StatusFor(id string) (string, string) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if s := e.sessions[id]; s != nil {
		return s.status, s.Phone
	}
	return "offline", ""
}

// Overall status: "connected" if any account is connected.
func (e *Engine) Status() (string, string) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, s := range e.sessions {
		if s.status == "connected" {
			return "connected", s.Phone
		}
	}
	for _, s := range e.sessions {
		if s.status == "qr" {
			return "qr", ""
		}
	}
	return "disconnected", ""
}

func (e *Engine) NotifyChan() <-chan string {
	return e.notifyCh
}

func (e *Engine) Notify(phone string) {
	select {
	case e.notifyCh <- phone:
	default:
	}
}

// LogoutAccount logs out & removes a single account.
func (e *Engine) LogoutAccount(id string) error {
	e.mu.RLock()
	s := e.sessions[id]
	e.mu.RUnlock()
	if s == nil {
		return nil
	}
	ctx := context.Background()
	_ = s.client.Logout(ctx)
	s.client.Disconnect()
	e.mu.Lock()
	delete(e.sessions, id)
	e.mu.Unlock()
	return nil
}

// AccountBelongsToUser checks whether the given phone number belongs to the user.
func (e *Engine) AccountBelongsToUser(userID int64, phone string) bool {
	if userID == 0 { return true }
	phone = strings.TrimPrefix(phone, "+")
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, s := range e.sessions {
		if s.Phone == phone {
			if s.userID == 0 { return true }
			return s.userID == userID
		}
	}
	return false
}

// GetSessionUserID returns the userID of a session by its internal ID.
func (e *Engine) GetSessionUserID(id string) int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if s, ok := e.sessions[id]; ok {
		return s.userID
	}
	return 0
}

// ValidateAccountIDs checks that all account phone numbers in the comma-separated list belong to the user.
func (e *Engine) ValidateAccountIDs(userID int64, accountIDs string) bool {
	if userID == 0 || accountIDs == "" { return true }
	for _, phone := range strings.Split(accountIDs, ",") {
		phone = strings.TrimSpace(phone)
		if phone == "" { continue }
		if !e.AccountBelongsToUser(userID, phone) {
			return false
		}
	}
	return true
}

// firstConnected returns any connected session (for sending).
func (e *Engine) userAllowedOnWAServer(userID int64) bool {
	servers, _ := e.db.ListWaServers()
	hasRestriction := false
	for _, s := range servers {
		if s.Packages != "" {
			hasRestriction = true
			break
		}
	}
	if !hasRestriction { return true } // no restrictions configured
	pkgName := e.db.GetUserPackageName(userID)
	if pkgName == "" { return false }
	for _, s := range servers {
		if s.Packages == "" { continue }
		for _, p := range strings.Split(s.Packages, ",") {
			if strings.TrimSpace(p) == pkgName { return true }
		}
	}
	return false
}
func (e *Engine) firstConnected(userID int64) *session {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, s := range e.sessions {
		if s.status == "connected" {
			if userID != 0 && s.userID != 0 && s.userID != userID { continue }
			return s
		}
	}
	return nil
}

func (e *Engine) FirstSession(userID int64) *session {
	return e.firstConnected(userID)
}

func (e *Engine) handleEvent(s *session, rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Connected:
		e.mu.Lock()
		s.status = "connected"
		if s.client.Store.ID != nil {
			s.Phone = s.client.Store.ID.User
		}
		e.mu.Unlock()
		e.db.LogInstance(s.Phone, "connected")
	case *events.Disconnected:
		e.mu.Lock()
		if s.status != "qr" {
			s.status = "disconnected"
		}
		e.mu.Unlock()
		e.db.LogInstance(s.Phone, "disconnected")
	case *events.LoggedOut:
		e.mu.Lock()
		s.status = "disconnected"
		e.mu.Unlock()
		e.db.LogInstance(s.Phone, "logged_out")
	case *events.Message:
		e.onMessage(s, evt)
	}
}

// onMessage: log + auto reply. Matching logic UNCHANGED.
func (e *Engine) onMessage(s *session, evt *events.Message) {
	if evt.Info.IsFromMe {
		text := extractText(evt.Message)
		if text != "" {
			to := evt.Info.Chat.User
			name := evt.Info.PushName
			e.db.LogSentForWA(s.Phone, to, text, "sent", "phone_sync")
			e.db.Log("send", "sent", fmt.Sprintf("outgoing (phone) -> %s: %s", to, text))
			e.dispatchWebhooks("sent", to, name, text, "private")
			select {
			case e.notifyCh <- to:
			default:
			}
		}
		return
	}
	isStatus := evt.Info.Chat == waTypes.StatusBroadcastJID
	if isStatus {
		senderPhone := evt.Info.Sender.User
		name := evt.Info.PushName
		text := extractText(evt.Message)
		mediaURL := extractMedia(evt.Message)
		e.db.LogStatus(senderPhone, name, text, mediaURL)
		e.log.Infof("status from %s: %s", name, text)
		return
	}
	text := extractText(evt.Message)
	if text == "" {
		return
	}
	senderPhone := evt.Info.Sender.User
	if evt.Info.Sender.Server == waTypes.HiddenUserServer && !evt.Info.SenderAlt.IsEmpty() {
		senderPhone = evt.Info.SenderAlt.User
	}
	name := evt.Info.PushName

	if evt.Info.IsGroup {
		groupJID := evt.Info.Chat.User
		groupName := e.db.GetGroupName(groupJID)
		if groupName == "" {
			gjid, err := waTypes.ParseJID(groupJID)
			if err == nil {
				if ginfo, gerr := s.client.GetGroupInfo(context.Background(), gjid); gerr == nil {
					groupName = ginfo.GroupName.Name
					e.db.SaveGroupName(groupJID, groupName)
				}
			}
		}
		e.db.LogReceivedForWA(s.Phone, groupJID, groupName, text, true, senderPhone, name, "whatsmeow")
		e.db.Log("received", "group", fmt.Sprintf("[%s] %s → %s: %s", groupName, name, groupJID, text))
		select {
		case e.notifyCh <- groupJID:
		default:
		}
		e.dispatchWebhooks("received", senderPhone, name, text, "group")
		e.log.Infof("group msg: %s → %s: %s", name, groupName, text)
	} else {
		e.db.LogReceivedForWA(s.Phone, senderPhone, name, text, false, "", "", "whatsmeow"); e.db.Log("received", "private", fmt.Sprintf("%s (%s): %s", name, senderPhone, text))
		// spam detection
		if e.db.TrackSpam(senderPhone, fmt.Sprintf("%x", text[:minInt(len(text), 20)])) {
			e.db.AddBlacklist(senderPhone, "auto: spam detected")
			e.db.Log("safety", "blacklisted", fmt.Sprintf("%s auto-blocked for spam", senderPhone))
		}
		// auto-assign round robin for new conversations
		if e.db.GetAssignedAgent(senderPhone) == 0 {
			e.db.AssignNextRoundRobin(senderPhone)
		}
		// auto-detect department from keywords
		if depts, _ := e.db.ListDepts(); len(depts) > 0 {
			for _, d := range depts {
				if strings.Contains(strings.ToLower(text), strings.ToLower(d.Name)) {
					e.db.AssignToDept(senderPhone, d.Name)
					break
				}
			}
		}
		select {
		case e.notifyCh <- senderPhone:
		default:
		}
		e.dispatchWebhooks("received", senderPhone, name, text, "private")
	}

	// Drip: auto-enroll + STOP handler
	trimmed := strings.TrimSpace(strings.ToLower(text))
	if trimmed == "stop" || trimmed == "berhenti" || trimmed == "unsub" {
		e.db.UnenrollFromDrip(senderPhone)
		e.db.Log("drip", "stop", fmt.Sprintf("%s stopped all drips", senderPhone))
	} else {
		drips, _ := e.db.ListDrips()
		for _, d := range drips {
			if d.Status == "active" {
				e.db.EnrollInDrip(d.ID, senderPhone, name)
			}
		}
	}

	// Store Bot: menu, category, product, order flow
	to := evt.Info.Chat
	if !evt.Info.IsGroup && to.Server == waTypes.HiddenUserServer && !evt.Info.SenderAlt.IsEmpty() {
		to = evt.Info.SenderAlt
	}
	if trimmed == "menu" || trimmed == "katalog" || trimmed == "produk" {
		cats, _ := e.db.ListCategories()
		if len(cats) > 0 {
			reply := "*Katalog Produk*\n\nBalas nama kategori:\n"
			for _, c := range cats { reply += "• " + c.Name + "\n" }
			e.sendVia(s, to, reply)
		} else { e.sendVia(s, to, "Belum ada produk.") }
		return
	}
	// Check if matches a category -> show products in that category
	if cats, _ := e.db.ListCategories(); len(cats) > 0 {
		for _, cat := range cats {
			if strings.EqualFold(strings.TrimSpace(cat.Name), trimmed) {
				prods, _ := e.db.ProductsByCategory(cat.Name)
				if len(prods) > 0 {
					reply := fmt.Sprintf("*%s*\n\n", cat.Name)
					for _, p := range prods {
						reply += fmt.Sprintf("*%s*\n%s\nRp%.0f\n\n", p.Name, p.Description, p.Price)
					}
					reply += "Balas nama produk untuk order."
					e.sendVia(s, to, reply)
				}
				return
			}
		}
	}
	// Check if matches a product -> create order
	if prods, _ := e.db.ListProducts(); len(prods) > 0 {
		for _, p := range prods {
			if strings.Contains(strings.ToLower(p.Name), trimmed) || strings.Contains(trimmed, strings.ToLower(p.Name)) {
				total := p.Price
				orderID, _ := e.db.CreateOrder(senderPhone, name, p.ID, 1, total)
				reply := fmt.Sprintf("✅ *Order #%d*\nProduk: %s\nHarga: Rp%.0f\n\nKetik *BAYAR* untuk pembayaran.", orderID, p.Name, total)
				e.sendVia(s, to, reply)
				return
			}
		}
	}
	if trimmed == "bayar" || trimmed == "checkout" || trimmed == "payment" {
		e.sendVia(s, to, "Untuk pembayaran, silakan kunjungi: "+e.db.GetSetting("app_url", "/")+"/subscribe")
		return
	}

	if evt.Info.IsGroup && e.db.GetSetting("reply_in_group", "0") != "1" {
		return
	}

	tv := msgtemplate.Vars{Name: name, Phone: senderPhone, Message: text}

	// Anti-spam guard: mute repeat messages
	if aiservice.CheckSpam(senderPhone, text) {
		e.log.Infof("spam muted: %s", senderPhone)
		return
	}

	// Jailbreak guard: reject prompt injection
	if aiservice.CheckJailbreak(text) {
		e.log.Infof("jailbreak blocked: %s", senderPhone)
		return
	}

	// AI global: reply to EVERY message using the configured AI key
	if e.db.GetSetting("ai_all_enabled", "0") == "1" && !evt.Info.IsGroup {
		// Fallback-only mode: skip if auto-reply already would have matched
		if e.db.GetSetting("ai_fallback_only", "0") == "1" {
			if e.hasKeywordMatch(text, s.Phone) {
				goto skipAIAll
			}
		}
		// Business hours check
		if e.db.GetSetting("biz_hours_enabled", "0") == "1" && !e.inBusinessHours() {
			goto skipAIAll
		}
		if aiAllIDStr := e.db.GetSetting("ai_all_key_id", "0"); aiAllIDStr != "0" {
			if aiAllID, err := strconv.ParseInt(aiAllIDStr, 10, 64); err == nil && aiAllID > 0 {
				if aik, err := e.db.GetAiKey(aiAllID); err == nil {
					decKey, _ := secret.Decrypt(aik.APIKey)
					if decKey == "" { decKey = aik.APIKey }
					if aiReply, aiErr := aiservice.Reply(decKey, aik.Provider, aik.Model, aik.BaseURL, aik.SystemPrompt, text, e.knowledgeRows(), nil); aiErr == nil && aiReply != "" {
						if err := e.sendVia(s, to, aiReply); err != nil {
							e.db.LogSent(to.User, aiReply, "failed", "whatsmeow"); e.db.Log("autoreply", "failed", fmt.Sprintf("FAILED -> %s: %s", to.User, aiReply))
						} else {
							e.db.LogSent(to.User, aiReply, "ai_all", "whatsmeow"); e.db.Log("ai", "sent", fmt.Sprintf("AI reply -> %s: %s", to.User, aiReply))
						}
						return
					}
				}
			}
		}
	}
skipAIAll:

	if e.db.GetSetting("welcome_enabled", "0") == "1" && !evt.Info.IsGroup {
		if e.db.MarkWelcomed(senderPhone) {
			if wmsg := e.db.GetSetting("welcome_message", ""); wmsg != "" {
				rendered := msgtemplate.Render(wmsg, tv)
				if err := e.sendVia(s, to, rendered); err == nil {
					e.db.LogSent(to.User, rendered, "welcome", "whatsmeow"); e.db.Log("welcome", "sent", fmt.Sprintf("welcome -> %s: %s", to.User, rendered))
				}
			}
		}
	}

	// Basic auto reply — filter by account (comma-separated)
	rule, ok := e.db.FindReplyFullForAccount(text, s.Phone)
	if !ok {
		rule, ok = e.db.FindReplyFullForAccount(text, "") // fallback: rules with no account set
	}
	if ok {
		// build knowledge context (global DB + per-rule FAQ)
		krows := e.knowledgeRows()
		if rule.Match == "ai" && rule.Reply != "" {
			for _, line := range strings.Split(rule.Reply, "\n") {
				line = strings.TrimSpace(line)
				if line == "" { continue }
				parts := strings.SplitN(line, "|", 2)
				q := strings.TrimSpace(parts[0])
				a := ""
				if len(parts) > 1 { a = strings.TrimSpace(parts[1]) }
				if q != "" { krows = append(krows, aiservice.KnowledgeRow{Question: q, Answer: a}) }
			}
		}

		// Human handoff: keyword trigger → stop AI → send admin contact
		if e.db.GetSetting("handoff_enabled", "0") == "1" {
			keywords := strings.Split(strings.ToLower(e.db.GetSetting("handoff_keywords", "admin,telp,manusia,cs,operator")), ",")
			msgLower := strings.ToLower(text)
			for _, kw := range keywords {
				if strings.TrimSpace(kw) != "" && strings.Contains(msgLower, strings.TrimSpace(kw)) {
					handoffMsg := msgtemplate.Render(e.db.GetSetting("handoff_message", "Silakan hubungi admin kami."), tv)
					e.sendVia(s, to, handoffMsg)
					e.db.LogSent(to.User, handoffMsg, "handoff", "whatsmeow"); e.db.Log("handoff", "sent", fmt.Sprintf("handoff -> %s: %s", to.User, handoffMsg))
					return
				}
			}
		}

		rendered := msgtemplate.Render(rule.Reply, tv)

		// AI Mode with store context (RAG)
		if rule.UseAI && rule.AiKeyID > 0 {
			// force own key check
			if e.db.GetSetting("force_own_key", "0") == "1" {
				if keys, _ := e.db.ListAiKeys(0); len(keys) == 0 { goto afterAI }
			}
			// training campaign override
			sysPrompt := ""
			aiKeyID := rule.AiKeyID
			if rule.TrainingID > 0 {
				for _, t := range e.db.MustListAiTrainings() {
					if t.ID == rule.TrainingID {
						sysPrompt = t.SystemPrompt
						if t.AiKeyID > 0 { aiKeyID = t.AiKeyID }
						break
					}
				}
			}
			if aik, err := e.db.GetAiKey(aiKeyID); err == nil {
				decKey, _ := secret.Decrypt(aik.APIKey)
				if decKey == "" { decKey = aik.APIKey }
				// memory window: fetch recent N messages
				memWindow, _ := strconv.Atoi(e.db.GetSetting("ai_memory_window", "5"))
				history := e.getChatHistory(senderPhone, memWindow)
				// reasoning level → temperature
				temp := 0.7
				switch e.db.GetSetting("ai_reasoning_level", "medium") {
				case "low": temp = 0.2
				case "high": temp = 0.9
				}
				cooldown := e.fallbackCooldownByPhone(senderPhone)
				// RAG: inject store catalog + customer profile into AI context
				prods, _ := e.db.ListProducts()
				if len(prods) > 0 {
					sysPrompt += "\nStore Products:\n"
					for i, p := range prods {
						if i >= 10 { break }
						sysPrompt += fmt.Sprintf("- %s: Rp%.0f (Stok:%d)\n", p.Name, p.Price, p.Stock)
					}
				}
				if cp := e.db.GetCustomerProfile(senderPhone); cp != nil && cp.TotalOrders > 0 {
					sysPrompt += fmt.Sprintf("\nCustomer Profile: %s, %d orders, total Rp%.0f\n", cp.Name, cp.TotalOrders, cp.TotalSpent)
				}
				if aiReply, aiErr := aiservice.ReplyWithContext(decKey, aik.Provider, aik.Model, aik.BaseURL, sysPrompt + aik.SystemPrompt, text, krows, cooldown, history, temp); aiErr == nil && aiReply != "" {
					rendered = aiReply
					e.db.RecordAiUsage(0, len(text)/4+len(aiReply)/4, aik.Provider, aik.Model)
				}
			}
		}

		if err := e.sendVia(s, to, rendered); err != nil {
			e.db.LogSent(to.User, rendered, "failed", "whatsmeow"); e.db.Log("autoreply", "failed", fmt.Sprintf("FAILED -> %s: %s", to.User, rendered))
			return
		}
		e.db.LogSent(to.User, rendered, "autoreply", "whatsmeow"); e.db.Log("autoreply", "sent", fmt.Sprintf("auto-reply -> %s: %s", to.User, rendered))
		return
	}
	// Fallback: no rule matched
	if e.db.GetSetting("biz_hours_enabled", "0") == "1" {
		reply := e.db.GetSetting("biz_hours_reply", "Saat ini di luar jam operasional.")
		e.sendVia(s, to, reply)
		e.db.LogSent(to.User, reply, "autoreply", "whatsmeow")
		return
	}
afterAI:

	if e.db.GetSetting("fallback_enabled", "0") == "1" && !evt.Info.IsGroup {
		if fmsg := e.db.GetSetting("fallback_message", ""); fmsg != "" {
			rendered := msgtemplate.Render(fmsg, tv)
			if err := e.sendVia(s, to, rendered); err == nil {
				e.db.LogSent(to.User, rendered, "fallback", "whatsmeow"); e.db.Log("fallback", "sent", fmt.Sprintf("fallback -> %s: %s", to.User, rendered))
			}
		}
	}
}

func (e *Engine) Send(userID int64, phone, message string) error {
	return e.SendFrom(userID, "", phone, message)
}

func (e *Engine) SendFrom(userID int64, accountPhone, phone, message string) error {
	s := e.findSession(userID, accountPhone)
	if s == nil {
		return fmt.Errorf("not connected")
	}
	var jid waTypes.JID
	if strings.HasSuffix(phone, "@g.us") || strings.HasSuffix(phone, "@s.whatsapp.net") {
		var err error
		jid, err = waTypes.ParseJID(phone)
		if err != nil {
			return fmt.Errorf("invalid jid: %w", err)
		}
	} else {
		digits := onlyDigits(phone)
		if digits == "" {
			return fmt.Errorf("invalid phone")
		}
		jid = waTypes.NewJID(digits, waTypes.DefaultUserServer)
	}
	if err := e.sendVia(s, jid, message); err != nil {
		e.db.LogSentForWA(s.Phone, phone, message, "failed", "whatsmeow"); e.db.Log("send", "failed", fmt.Sprintf("FAILED -> %s: %s", phone, message))
		return err
	}
	e.db.LogSentForWA(s.Phone, phone, message, "sent", "whatsmeow"); e.db.Log("send", "sent", fmt.Sprintf("outgoing -> %s: %s", phone, message))
	e.dispatchWebhooks("sent", phone, "", message, "private")
	return nil
}

func (e *Engine) findSession(userID int64, accountPhone string) *session {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, s := range e.sessions {
		if s.status != "connected" { continue }
		if userID != 0 && s.userID != 0 && s.userID != userID { continue }
		if accountPhone != "" {
			if s.Phone == accountPhone { return s }
		} else {
			return s
		}
	}
	return nil
}

type SenderSelector struct {
	Phones []string
	Index  int
	Mode   string // "round_robin" or "random"
}

func (sel *SenderSelector) Next(e *Engine) *session {
	if len(sel.Phones) == 0 { return e.firstConnected(0) }
	var phone string
	switch sel.Mode {
	case "round_robin":
		phone = sel.Phones[sel.Index%len(sel.Phones)]
		sel.Index++
	default: // random
		phone = sel.Phones[time.Now().UnixNano()%int64(len(sel.Phones))]
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, s := range e.sessions {
		if s.status == "connected" && s.Phone == phone {
			return s
		}
	}
	return e.firstConnected(0)
}

func (e *Engine) SendMedia(userID int64, accountPhone, phone, mediaType, filePath, caption string) error {
	s := e.findSession(userID, accountPhone)
	if s == nil { return fmt.Errorf("not connected") }
	digits := onlyDigits(phone)
	if digits == "" { return fmt.Errorf("invalid phone") }
	data, err := os.ReadFile(filePath)
	if err != nil { return err }
	jid := waTypes.NewJID(digits, waTypes.DefaultUserServer)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var msg waProto.Message
	thumb := data
	if len(thumb) > 50000 { thumb = thumb[:50000] }
	switch mediaType {
	case "image":
		msg.ImageMessage = &waProto.ImageMessage{
			Caption:        proto.String(caption),
			JPEGThumbnail: thumb,
		}
		resp, err := s.client.Upload(ctx, data, whatsmeow.MediaImage)
		if err != nil { return err }
		msg.ImageMessage.URL = proto.String(resp.URL)
		msg.ImageMessage.DirectPath = proto.String(resp.DirectPath)
		msg.ImageMessage.MediaKey = resp.MediaKey
		msg.ImageMessage.FileEncSHA256 = resp.FileEncSHA256
		msg.ImageMessage.FileSHA256 = resp.FileSHA256
		msg.ImageMessage.Mimetype = proto.String("image/jpeg")
	case "video":
		msg.VideoMessage = &waProto.VideoMessage{
			Caption:        proto.String(caption),
			JPEGThumbnail: thumb,
		}
		resp, err := s.client.Upload(ctx, data, whatsmeow.MediaVideo)
		if err != nil { return err }
		msg.VideoMessage.URL = proto.String(resp.URL)
		msg.VideoMessage.DirectPath = proto.String(resp.DirectPath)
		msg.VideoMessage.MediaKey = resp.MediaKey
		msg.VideoMessage.FileEncSHA256 = resp.FileEncSHA256
		msg.VideoMessage.FileSHA256 = resp.FileSHA256
		msg.VideoMessage.Mimetype = proto.String("video/mp4")
	case "document":
		msg.DocumentMessage = &waProto.DocumentMessage{
			Caption: proto.String(caption),
		}
		resp, err := s.client.Upload(ctx, data, whatsmeow.MediaDocument)
		if err != nil { return err }
		msg.DocumentMessage.URL = proto.String(resp.URL)
		msg.DocumentMessage.DirectPath = proto.String(resp.DirectPath)
		msg.DocumentMessage.MediaKey = resp.MediaKey
		msg.DocumentMessage.FileEncSHA256 = resp.FileEncSHA256
		msg.DocumentMessage.FileSHA256 = resp.FileSHA256
		msg.DocumentMessage.Mimetype = proto.String("application/octet-stream")
		msg.DocumentMessage.FileName = proto.String(filepath.Base(filePath))
	}
	_, err = s.client.SendMessage(ctx, jid, &msg)
	if err != nil {
		e.db.LogSentForWA(s.Phone, digits, caption+" [media]", "failed", "whatsmeow")
		return err
	}
	e.db.LogSentForWA(s.Phone, digits, caption+" [media]", "sent", "whatsmeow")
	e.dispatchWebhooks("sent", digits, "", caption+" [media]", "private")
	return nil
}

func (e *Engine) sendVia(s *session, to waTypes.JID, message string) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("no session/client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := s.client.SendMessage(ctx, to, &waProto.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		e.log.Errorf("sendVia failed to %s: %v", to.User, err)
	}
	return err
}

func extractText(m *waProto.Message) string {
	if m == nil {
		return ""
	}
	if m.GetConversation() != "" {
		return m.GetConversation()
	}
	if ext := m.GetExtendedTextMessage(); ext != nil {
		return ext.GetText()
	}
	if img := m.GetImageMessage(); img != nil {
		return img.GetCaption()
	}
	if vid := m.GetVideoMessage(); vid != nil {
		return vid.GetCaption()
	}
	return ""
}

func extractMedia(m *waProto.Message) string {
	if m == nil { return "" }
	if img := m.GetImageMessage(); img != nil {
		if img.GetURL() != "" { return img.GetURL() }
		return "image"
	}
	if vid := m.GetVideoMessage(); vid != nil {
		if vid.GetURL() != "" { return vid.GetURL() }
		return "video"
	}
	return ""
}

func (e *Engine) getChatHistory(phone string, max int) []string {
	var out []string
	msgs, _ := e.db.ListReceivedPaginated(1, max)
	for _, m := range msgs {
		if m.Phone == phone { out = append(out, m.Message) }
	}
	return out
}

func (e *Engine) hasKeywordMatch(text, phone string) bool {
	_, ok := e.db.FindReplyFullForAccount(text, phone)
	if !ok {
		_, ok = e.db.FindReplyFullForAccount(text, "")
	}
	return ok
}

func (e *Engine) inBusinessHours() bool {
	now := time.Now()
	day := now.Format("Monday")
	offDays := strings.ToLower(e.db.GetSetting("biz_hours_off_days", "Saturday,Sunday"))
	for _, d := range strings.Split(offDays, ",") {
		if strings.ToLower(strings.TrimSpace(d)) == strings.ToLower(day) { return false }
	}
	start := e.db.GetSetting("biz_hours_start", "08:00")
	end := e.db.GetSetting("biz_hours_end", "17:00")
	current := now.Format("15:04")
	return current >= start && current <= end
}

func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (e *Engine) knowledgeRows() []aiservice.KnowledgeRow {
	rows := e.db.ActiveKnowledgeRows()
	out := make([]aiservice.KnowledgeRow, len(rows))
	for i, r := range rows {
		out[i] = aiservice.KnowledgeRow{
			Question: r.Question,
			Answer:   r.Answer,
			Category: r.Category,
		}
	}
	return out
}

var fallbackCooldownsMu sync.Mutex
var fallbackCooldowns = map[string]*aiservice.Cooldown{}

func (e *Engine) fallbackCooldownByPhone(phone string) *aiservice.Cooldown {
	if e.db.GetSetting("fallback_enabled", "0") != "1" {
		return nil
	}
	fallbackCooldownsMu.Lock()
	defer fallbackCooldownsMu.Unlock()
	c, ok := fallbackCooldowns[phone]
	if !ok {
		c = &aiservice.Cooldown{Window: 10 * time.Minute, Max: 3}
		fallbackCooldowns[phone] = c
	}
	return c
}

func (e *Engine) dispatchWebhooks(event, phone, name, message, msgType string) {
	go func() {
		hooks, err := e.db.WebhooksForEvent(event)
		if err != nil || len(hooks) == 0 {
			return
		}
		payload, _ := json.Marshal(map[string]string{
			"event": event, "phone": phone, "name": name,
			"message": message, "type": msgType,
		})
		for _, h := range hooks {
			go func(url string) {
				http.Post(url, "application/json", bytes.NewReader(payload))
			}(h.URL)
		}
	}()
}
