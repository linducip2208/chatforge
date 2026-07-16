package main

import (
	crand "crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"chatgo/aiservice"
	"chatgo/i18n"
	"chatgo/meta"
	"chatgo/msgtemplate"
	"chatgo/secret"
	"chatgo/store"
	"chatgo/wa"

	qrcode "github.com/skip2/go-qrcode"
)

var (
	db     *store.DB
	engine *wa.Engine
)

func main() {
	// load .env file if present
	if _, err := os.Stat(".env"); err == nil {
		raw, _ := os.ReadFile(".env")
		for _, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) == 2 {
				k := strings.TrimSpace(kv[0])
				v := strings.TrimSpace(kv[1])
				if os.Getenv(k) == "" {
					os.Setenv(k, v)
				}
			}
		}
	}
	dataDir := "data"
	_ = os.MkdirAll(dataDir, 0o755)

	if err := i18n.Load("lang"); err != nil {
		log.Fatalf("load lang: %v", err)
	}

	dsn := os.Getenv("CHATGO_MYSQL")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4"
	}
	var err error
	db, err = store.Open(dsn)
	if err != nil {
		log.Fatalf("open mysql db: %v", err)
	}

	engine, err = wa.New(filepath.Join(dataDir, "session.db"), db)
	if err != nil {
		log.Fatalf("wa engine: %v", err)
	}
	if err := engine.Start(); err != nil {
		log.Printf("wa start: %v (will retry via QR page)", err)
	}
	engine.StartLoops()

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	mux.Handle("/screens/", http.StripPrefix("/screens/", http.FileServer(http.Dir("public/marketing/screens"))))
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { render(w, r, "login") })
	mux.HandleFunc("/login/post", loginUser)
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { render(w, r, "register") })
	mux.HandleFunc("/register/post", registerUser)
	mux.HandleFunc("/logout", logoutUser)

	// REST API (no auth — uses X-API-Key)
	mux.HandleFunc("/api/send", handleAPISend)
	mux.HandleFunc("/api/status", handleAPIStatus)
	mux.HandleFunc("/api/messages", handleAPIMessages)
	mux.HandleFunc("/api/contacts", handleAPIContacts)
	mux.HandleFunc("/api/devices", handleAPIDevices)

	// All page routes wrapped with auth middleware
	mux.HandleFunc("/wa", p("wa"))
	mux.HandleFunc("/wa/add", authMiddleware(handleWaAdd))
	mux.HandleFunc("/wa/logout", authMiddleware(handleWaLogout))
	mux.HandleFunc("/send", authMiddleware(handleSend))
	mux.HandleFunc("/send/media", authMiddleware(handleSendMedia))
	mux.HandleFunc("/sent", p("sent"))
	mux.HandleFunc("/received", p("received"))
	mux.HandleFunc("/inbox", p("inbox"))
	mux.HandleFunc("/inbox/chat", authMiddleware(handleInboxChat))
	mux.HandleFunc("/inbox/events", authMiddleware(handleInboxEvents))
	mux.HandleFunc("/inbox/send", authMiddleware(handleInboxSend))
	mux.HandleFunc("/inbox/send-meta", authMiddleware(handleInboxSendMeta))
	mux.HandleFunc("/inbox/messages", authMiddleware(handleInboxMessages))
	mux.HandleFunc("/inbox/unread-count", authMiddleware(handleInboxUnreadCount))
	mux.HandleFunc("/inbox/mark-read", authMiddleware(handleInboxMarkRead))
	mux.HandleFunc("/inbox/search", authMiddleware(handleInboxSearch))
	mux.HandleFunc("/autoreply", p("autoreply"))
	mux.HandleFunc("/settings", authMiddleware(handleSettings))
	mux.HandleFunc("/admin/users/impersonate", authMiddleware(handleImpersonate))
	mux.HandleFunc("/exit-impersonation", handleExitImpersonation)
	mux.HandleFunc("/contacts", p("contacts"))
	mux.HandleFunc("/contacts/groups", p("groups"))
	mux.HandleFunc("/contacts/unsub", p("unsub"))
	mux.HandleFunc("/contacts/add", authMiddleware(crudPost(func(r *http.Request) { db.AddContact(r.FormValue("name"), r.FormValue("phone"), joinVals(r, "groups")) }, "/contacts")))
	mux.HandleFunc("/contacts/delete", authMiddleware(crudDel(func(id int64) { db.DeleteContact(id) }, "/contacts")))
	mux.HandleFunc("/contacts/import", authMiddleware(handleContactImport))
	mux.HandleFunc("/groups/add", authMiddleware(crudPost(func(r *http.Request) { db.AddGroup(r.FormValue("name")) }, "/contacts/groups")))
	mux.HandleFunc("/groups/delete", authMiddleware(crudDel(func(id int64) { db.DeleteGroup(id) }, "/contacts/groups")))
	mux.HandleFunc("/unsub/add", authMiddleware(crudPost(func(r *http.Request) { db.AddUnsub(r.FormValue("phone")) }, "/contacts/unsub")))
	mux.HandleFunc("/unsub/delete", authMiddleware(crudDel(func(id int64) { db.DeleteUnsub(id) }, "/contacts/unsub")))
	mux.HandleFunc("/broadcast", authMiddleware(handleBroadcast))
	mux.HandleFunc("/broadcast/stop", authMiddleware(crudDel(func(id int64) { db.UpdateCampaignStatus(id, "stopped") }, "/broadcast")))
	mux.HandleFunc("/broadcast/delete", authMiddleware(crudDel(func(id int64) { db.DeleteCampaign(id) }, "/broadcast")))
		mux.HandleFunc("/scheduled", authMiddleware(handleScheduled))
	mux.HandleFunc("/scheduled/delete", authMiddleware(crudDel(func(id int64) { db.DeleteScheduled(id) }, "/scheduled")))
	mux.HandleFunc("/templates", p("templates"))
	mux.HandleFunc("/templates/add", authMiddleware(crudPost(func(r *http.Request) { db.AddTemplate(r.FormValue("name"), r.FormValue("content")) }, "/templates")))
	mux.HandleFunc("/templates/delete", authMiddleware(crudDel(func(id int64) { db.DeleteTemplate(id) }, "/templates")))
	mux.HandleFunc("/apikeys", p("apikeys"))
	mux.HandleFunc("/apikeys/add", authMiddleware(crudPost(func(r *http.Request) { db.AddAPIKey(r.FormValue("name"), randSecret()) }, "/apikeys")))
	mux.HandleFunc("/apikeys/delete", authMiddleware(crudDel(func(id int64) { db.DeleteAPIKey(id) }, "/apikeys")))
	mux.HandleFunc("/webhooks", p("webhooks"))
	mux.HandleFunc("/webhooks/add", authMiddleware(crudPost(func(r *http.Request) { db.AddWebhook(r.FormValue("name"), r.FormValue("url"), r.FormValue("event")) }, "/webhooks")))
	mux.HandleFunc("/webhooks/delete", authMiddleware(crudDel(func(id int64) { db.DeleteWebhook(id) }, "/webhooks")))
	mux.HandleFunc("/logger", p("logger"))
	mux.HandleFunc("/logger/clear", authMiddleware(func(w http.ResponseWriter, r *http.Request) { db.ClearLog(); http.Redirect(w, r, "/logger", http.StatusSeeOther) }))
	registerAdminRoutes(mux)
	mux.HandleFunc("/lang/", handleLang)
	mux.HandleFunc("/qr.png", handleQRImage)
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/webhook/meta", handleMetaWebhook)
	mux.HandleFunc("/autoreply/add", authMiddleware(handleAutoReplyAdd))
	mux.HandleFunc("/autoreply/delete", authMiddleware(handleAutoReplyDelete))
	mux.HandleFunc("/autoreply/toggle", authMiddleware(handleAutoReplyToggle))
	mux.HandleFunc("/autoreply/edit", authMiddleware(handleAutoReplyEdit))

	// Contact edit
	mux.HandleFunc("/contacts/edit", authMiddleware(handleContactEdit))

	// Template edit
	mux.HandleFunc("/templates/edit", authMiddleware(handleTemplateEdit))

	// Android Hosts
	mux.HandleFunc("/hosts/android", p("hosts_android"))
	mux.HandleFunc("/devices/add", authMiddleware(crudPost(func(r *http.Request) { db.AddDevice(r.FormValue("name"), r.FormValue("did"), r.FormValue("manufacturer")) }, "/hosts/android")))
	mux.HandleFunc("/devices/delete", authMiddleware(crudDel(func(id int64) { db.DeleteDevice(id) }, "/hosts/android")))

	// USSD
	mux.HandleFunc("/ussd", p("ussd"))
	mux.HandleFunc("/ussd/add", authMiddleware(crudPost(func(r *http.Request) { db.AddUssd(r.FormValue("code")) }, "/ussd")))
	mux.HandleFunc("/ussd/delete", authMiddleware(crudDel(func(id int64) { db.DeleteUssd(id) }, "/ussd")))

	addr := "127.0.0.1:8080"
	if v := os.Getenv("CHATGO_ADDR"); v != "" {
		addr = v
	}
	fmt.Printf("\n  %s running at http://%s\n\n", getEnv("APP_NAME", "chatgo"), addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

type langInfo struct {
	Code, Name, Flag string
}

type pageData struct {
	Title, Pretitle, Heading, Icon string
	Active, Page                   string
	Status, Phone                  string
	Flash                          string
	AutoReplies                    []store.AutoReply
	Sent                           []store.SentMessage
	Received                       []store.ReceivedMessage
	CountSent, CountReceived       int
	// i18n
	LangCode, LangName, LangFlag string
	Languages                    []langInfo
	WaConnectedDesc              template.HTML
	WaScanHint                   template.HTML
	// multi-account
	Accounts          []wa.AccountInfo
	ConnectedAccounts []wa.AccountInfo
	HasConnected      bool
	ScanAccount       string
	AccountLimit      int
	ConnectedCount    int
	DisconnectedCount int
	// role
	Role            string
	EditID          int64
	EditName        string
	EditPhone       string
	EditContent     string
	EditGroups      string
	EditKeyword     string
	EditMatch       string
	EditReply       string
	EditUseAI       bool
	EditAiKeyID     int64
	EditAccountID   string
	EditTrainingID  int64
	EditRole        string
	WelcomeEnabled  bool
	WelcomeMessage  string
	FallbackEnabled bool
	FallbackMessage string
	ReplyInGroup    bool
	AiAllEnabled  bool
	AiAllKeyID    int64
	HandoffEnabled   bool
	HandoffMessage   string
	HandoffKeywords  string
	AiFallbackOnly   bool
	AiMemoryWindow   int
	AiDelaySeconds   int
	AiReasoningLevel string
	BizHoursEnabled  bool
	BizHoursStart    string
	BizHoursEnd      string
	BizHoursOffDays  string
	ForceOwnKey      bool
	Registrations    bool
	AiTokenQuota     int64
	AiTokenUsed      int64
	ChartLabels       template.JS
	ChartSent         template.JS
	ChartReceived     template.JS
	TotalUsers        int
	TotalLogins       int
	LoginsToday       int
	ActiveAccounts    int
	RunningCampaigns  int
	ActiveAccountList []wa.AccountInfo
	SendTo            string
	// full-menu entities
	Contacts   []store.Contact
	Groups     []store.Group
	Unsubs     []store.Unsub
	Templates  []store.Template
	APIKeys    []store.APIKey
	Webhooks   []store.Webhook
	Campaigns  []store.Campaign
	Scheduleds []store.Scheduled
	Logs       []store.LogEntry
	// admin/ai/devices
	Users         []store.User
	Roles         []store.Role
	Packages      []store.Package
	Vouchers      []store.Voucher
	Subscriptions []store.Subscription
	Transactions  []store.Transaction
	Payouts       []store.Payout
	Pages         []store.Page
	Marketings    []store.Marketing
	LanguagesAdm  []store.Language
	WaServers     []store.WaServer
	Gateways      []store.Gateway
	Shorteners    []store.Shortener
	Plugins       []store.Plugin
	MetaAccounts  []store.MetaAccount
	MetaTemplates []store.MetaTemplate
	IsImpersonating bool
	AiKeys        []store.AiKey
	AiPlugins     []store.AiPlugin
	AiTrainings   []store.AiTraining
	Devices       []store.Device
	Ussds         []store.Ussd
	Knowledges    []store.KnowledgeEntry
	DocsSteps     []DocsStep
	InboxConversations []store.InboxConversation
	ChatMessages  []store.ChatMessage
	UnreadCount   int
	IsGroup       bool
	ChatName      string
	Channel       string
	AppName       string
	AppLogo       string
	AppEmail      string
	AppURL        string
	Statuses      []store.WAStatus
	// pagination
	SentPage       int
	SentPerPage    int
	SentTotal      int
	SentPages      []int
	ReceivedPage   int
	ReceivedPerPage int
	ReceivedTotal  int
	ReceivedPages  []int
	PageNum        int
	InboxTotal     int
	InboxPages     []int
	LogTotal       int
	LogPages       []int
}

type DocsStep struct {
	Num   int
	Title string
	Desc  string
}

var allDocsSteps = []DocsStep{
	{1, "Login & Register", "Buka halaman login, masukkan email & password. Admin default: admin@chatgo.test / password. User baru bisa register."},
	{2, "Setup WA Account", "Buka tab WA > Akun & QR. Klik Tambah Akun. Scan QR code dengan WhatsApp (Linked Devices). Tunggu status Connected."},
	{3, "Tambah Kontak", "Menu Kontak > Tersimpan. Tambah kontak: nama + nomor WA + grup. Support multiple grup."},
	{4, "Buat Grup Kontak", "Menu Kontak > Grup. Buat grup untuk broadcast tertarget (VIP, Reseller, dll)."},
	{5, "Auto Reply - Basic", "Buka Auto Reply > tab Rules. Tambah rule: Match Type (Contains/Exact), Keyword, Reply Text. Support spintax {Halo|Hai}."},
	{6, "Auto Reply - AI Mode", "Tab AI Config: tambah AI Key (OpenAI/DeepSeek/Gemini). Tab Rules: pilih Match Type AI, centang Use AI, pilih key. AI auto balas."},
	{7, "Auto Reply - Multi WA", "Pilih nomor WA mana yang reply (checkbox WA Account). Kosongkan = semua nomor. Rule jalan sesuai nomor penerima."},
	{8, "FAQ / Knowledge Base", "Tab FAQ: tambah FAQ manual atau upload CSV/PDF/URL. AI akan search FAQ sebelum jawab (function calling)."},
	{9, "Training Campaign", "Tab Training: buat campaign dengan System Prompt berbeda (CS Produk, CS Teknis). Assign ke rule via dropdown."},
	{10, "AI Settings", "Tab AI Config: AI Global (balas semua chat), Fallback Only, Memory Window, Reasoning Level, Jam Kerja, Force Own Key."},
	{11, "Human Handoff", "Di AI Config > Controls: enable Handoff. Keyword trigger (admin, operator) → AI stop → kirim kontak admin."},
	{12, "Kirim Pesan", "Menu Kirim Pesan. Input nomor + pesan. Pilih nomor WA pengirim. Support text + gambar/video/dokumen."},
	{13, "Broadcast / Campaign", "Menu Broadcast. Pilih grup kontak, pilih nomor WA, set interval (default 300 detik). Round-robin otomatis."},
	{14, "Pesan Terjadwal", "Menu Terjadwal. Set nama, nomor, pesan, jadwal (datetime), repeat (menit). Pilih nomor WA pengirim."},
	{15, "Template Pesan", "Menu Template. Simpan template pesan yang sering dipakai. Support variabel {name} {phone}."},
	{16, "Welcome Message", "Menu Pengaturan > Pesan Sambutan. Dikirim ke kontak baru (24 jam cooldown)."},
	{17, "Fallback Message", "Menu Pengaturan > Balasan Default. Dikirim saat tidak ada keyword cocok (max 3x/10 menit)."},
	{18, "API Keys & Webhooks", "Menu API Keys: generate key. REST API: POST /api/send. Menu Webhooks: notifikasi real-time ke URL."},
	{19, "Multi-User & Paket", "Menu Admin: kelola user, role, paket (limit WA/quota), voucher, subscription. Dashboard SaaS."},
	{20, "Deployment", "Edit .env untuk MySQL & listen address. Single binary: chatgo.exe. Reverse proxy Nginx. Satu file, zero dependency."},
}

// current language from cookie (fallback to default)
func currentLang(r *http.Request) string {
	if c, err := r.Cookie("chatgo_lang"); err == nil && i18n.Has(c.Value) {
		return c.Value
	}
	return i18n.Default()
}

func getUserID(r *http.Request) int64 {
	idStr := r.Header.Get("X-User-ID")
	if idStr == "" {
		return 0
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" { return v }
	return def
}

func appURL() string {
	u := getEnv("APP_URL", "http://127.0.0.1:8080")
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}
	return strings.TrimRight(u, "/")
}

func render(w http.ResponseWriter, r *http.Request, page string) {
	lang := currentLang(r)
	T := i18n.Translator(lang)

	status, phone := engine.Status()
	uid := getUserID(r)
	ars, _ := db.ListAutoReplies()
	sent, _ := db.ListSentPaginated(1, 20)
	received, _ := db.ListReceivedPaginated(1, 20)

	langs := make([]langInfo, 0)
	for _, l := range i18n.List() {
		langs = append(langs, langInfo{Code: l.Code, Name: l.Name, Flag: l.Flag})
	}
	cur := i18n.Get(lang)

	d := pageData{
		Page: page, Active: page, Status: status, Phone: phone,
		Flash:       r.URL.Query().Get("msg"),
		AutoReplies: ars, Sent: sent, Received: received,
		CountSent: db.CountSent(), CountReceived: db.CountReceived(),
		LangCode: cur.Code, LangName: cur.Name, LangFlag: cur.Flag,
		Languages: langs,
		SentPerPage: 10, ReceivedPerPage: 10,
		SentPage: 1, ReceivedPage: 1,
		AppName: db.GetSetting("app_name", getEnv("APP_NAME", "ChatGo")),
		AppLogo: db.GetSetting("app_logo", getEnv("APP_LOGO", "/assets/theme/default-logo-light.png")),
		AppEmail: db.GetSetting("app_email", getEnv("APP_EMAIL", "admin@chatgo.test")),
		AppURL: appURL(),
	}
	// pre-rendered translated strings with HTML/format
	d.WaConnectedDesc = template.HTML(fmt.Sprintf(T("wa_connected_desc"), template.HTMLEscapeString(phone)))
	d.WaScanHint = template.HTML(T("wa_scan_hint"))
	if uid > 0 {
		if u, err := db.GetUserByID(uid); err == nil {
			d.Role = u.Role
		}
	}
	d.Accounts = engine.Accounts(uid)
	d.AccountLimit = engine.AccountLimit(uid)
	for _, a := range d.Accounts {
		if a.Status == "connected" { d.ConnectedAccounts = append(d.ConnectedAccounts, a) }
	}
	d.HasConnected = len(d.ConnectedAccounts) > 0
	d.SendTo = r.URL.Query().Get("to")
	d.ConnectedCount = engine.CountConnected(uid)
	d.DisconnectedCount = engine.CountDisconnected(uid)
	if scan := r.URL.Query().Get("scan"); scan != "" {
		d.ScanAccount = scan
		log.Printf("DEBUG ScanAccount set to: %s", scan)
	}
	// force if wa page and there's a new: session
	if page == "wa" && d.ScanAccount == "" {
		for _, a := range d.Accounts {
			if strings.HasPrefix(a.ID, "new:") {
				d.ScanAccount = a.ID
				break
			}
		}
	}

	// settings values (for the settings page)
	d.WelcomeEnabled = db.GetSetting("welcome_enabled", "0") == "1"
	d.WelcomeMessage = db.GetSetting("welcome_message", "")
	d.FallbackEnabled = db.GetSetting("fallback_enabled", "0") == "1"
	d.FallbackMessage = db.GetSetting("fallback_message", "")
	d.ReplyInGroup = db.GetSetting("reply_in_group", "0") == "1"
	d.AiAllEnabled = db.GetSetting("ai_all_enabled", "0") == "1"
	d.AiAllKeyID, _ = strconv.ParseInt(db.GetSetting("ai_all_key_id", "0"), 10, 64)
	d.HandoffEnabled = db.GetSetting("handoff_enabled", "0") == "1"
	d.HandoffMessage = db.GetSetting("handoff_message", "Silakan hubungi admin kami di nomor ini.")
	d.HandoffKeywords = db.GetSetting("handoff_keywords", "admin,telp,manusia,cs,operator")
	d.AiFallbackOnly = db.GetSetting("ai_fallback_only", "0") == "1"
	d.AiMemoryWindow, _ = strconv.Atoi(db.GetSetting("ai_memory_window", "5"))
	d.AiDelaySeconds, _ = strconv.Atoi(db.GetSetting("ai_delay_seconds", "0"))
	d.AiReasoningLevel = db.GetSetting("ai_reasoning_level", "medium")
	d.BizHoursEnabled = db.GetSetting("biz_hours_enabled", "0") == "1"
	d.BizHoursStart = db.GetSetting("biz_hours_start", "08:00")
	d.BizHoursEnd = db.GetSetting("biz_hours_end", "17:00")
	d.BizHoursOffDays = db.GetSetting("biz_hours_off_days", "Saturday,Sunday")
	d.ForceOwnKey = db.GetSetting("force_own_key", "0") == "1"
	d.Registrations = db.GetSetting("registrations", "1") == "1"
	d.AiTokenQuota = int64(db.GetUserAiQuota(uid))
	d.AiTokenUsed = db.GetAiTokenUsage(uid)
	d.UnreadCount = db.UnreadCount()
	d.IsImpersonating = r.Header.Get("X-Impersonating") == "1"

	// load entity lists per page (only what's needed)
	switch page {
	case "contacts":
		d.Contacts, _ = db.ListContacts()
		d.Groups, _ = db.ListGroups()
	case "groups":
		d.Groups, _ = db.ListGroups()
	case "unsub":
		d.Unsubs, _ = db.ListUnsub()
	case "templates":
		d.Templates, _ = db.ListTemplates()
	case "apikeys":
		d.APIKeys, _ = db.ListAPIKeys()
	case "autoreply":
		d.AiKeys, _ = db.ListAiKeys()
		d.Knowledges, _ = db.ListKnowledge()
		d.AiTrainings, _ = db.ListAiTrainings()
	case "webhooks":
		d.Webhooks, _ = db.ListWebhooks()
	case "broadcast":
		d.Campaigns, _ = db.ListCampaigns()
		d.Groups, _ = db.ListGroups()
	case "scheduled":
		d.Scheduleds, _ = db.ListScheduled()
	case "sent":
		d.SentPage = pageFromQuery(r)
		d.Sent, _ = db.ListSentPaginated(d.SentPage, d.SentPerPage)
		d.SentTotal = db.CountSent()
		d.SentPages = pageNums(d.SentPage, (d.SentTotal+d.SentPerPage-1)/d.SentPerPage)
	case "inbox":
		d.PageNum = pageFromQuery(r)
		d.InboxConversations, _ = db.GroupInboxPaginated(d.PageNum, 10)
		d.InboxTotal = db.CountInbox()
		d.InboxPages = pageNums(d.PageNum, (d.InboxTotal+9)/10)
		d.Statuses, _ = db.ListStatuses()
	case "inbox_chat":
		d.Phone = r.URL.Query().Get("phone")
		d.ChatMessages, _ = db.ChatHistory(d.Phone, 100)
		d.Templates, _ = db.ListTemplates()
		d.MetaAccounts, _ = db.ListMetaAccounts()
		if msgs, _ := db.ChatHistory(d.Phone, 1); len(msgs) > 0 {
			d.IsGroup = msgs[0].IsGroup
			d.Channel = msgs[0].Channel
			if d.IsGroup {
				if nm := db.GetGroupName(d.Phone); nm != "" {
					d.ChatName = nm
				}
			}
		}
	case "received":
		d.ReceivedPage = pageFromQuery(r)
		d.Received, _ = db.ListReceivedPaginated(d.ReceivedPage, d.ReceivedPerPage)
		d.ReceivedTotal = db.CountReceived()
		d.ReceivedPages = pageNums(d.ReceivedPage, (d.ReceivedTotal+d.ReceivedPerPage-1)/d.ReceivedPerPage)
			case "logger":
		d.PageNum = pageFromQuery(r)
		d.Logs, _ = db.ListLogPaginated(d.PageNum, 10)
		d.LogTotal = db.CountLog()
		d.LogPages = pageNums(d.PageNum, (d.LogTotal+9)/10)
	case "hosts_android":
		d.Devices, _ = db.ListDevices()
	case "ussd":
		d.Ussds, _ = db.ListUssd()
	case "ai_keys":
		d.AiKeys, _ = db.ListAiKeys()
	case "ai_plugins":
		d.AiPlugins, _ = db.ListAiPlugins()
	case "admin_users":
		d.Users, _ = db.ListUsers()
		d.Roles, _ = db.ListRoles()
	case "admin_roles":
		d.Roles, _ = db.ListRoles()
	case "admin_packages":
		d.Packages, _ = db.ListPackages()
	case "admin_vouchers":
		d.Vouchers, _ = db.ListVouchers()
		d.Packages, _ = db.ListPackages()
	case "admin_subscriptions":
		d.Subscriptions, _ = db.ListSubscriptions()
		d.Packages, _ = db.ListPackages()
		d.Users, _ = db.ListUsers()
	case "admin_transactions":
		d.Transactions, _ = db.ListTransactions()
	case "admin_payouts":
		d.Payouts, _ = db.ListPayouts()
	case "admin_pages":
		d.Roles, _ = db.ListRoles()
		d.Pages, _ = db.ListPages()
	case "admin_marketing":
		d.Marketings, _ = db.ListMarketing()
	case "admin_languages":
		d.LanguagesAdm, _ = db.ListLanguagesAdmin()
	case "admin_waservers":
		d.WaServers, _ = db.ListWaServers()
		d.Packages, _ = db.ListPackages()
	case "admin_gateways":
		d.Gateways, _ = db.ListGateways()
	case "admin_shorteners":
		d.Shorteners, _ = db.ListShorteners()
	case "admin_plugins":
		d.Plugins, _ = db.ListPlugins()
	case "admin_meta":
		d.MetaAccounts, _ = db.ListMetaAccounts()
	case "admin_metatemplates":
		d.MetaTemplates, _ = db.ListMetaTemplates()
	case "knowledge":
		d.Knowledges, _ = db.ListKnowledge()
	case "docs":
		d.DocsSteps = allDocsSteps
	}

	// edit mode – pre-fill forms
	if eid, _ := strconv.ParseInt(r.URL.Query().Get("edit"), 10, 64); eid > 0 {
		d.EditID = eid
		switch page {
		case "contacts":
			if c, err := db.GetContact(eid); err == nil {
				d.EditName = c.Name; d.EditPhone = c.Phone; d.EditGroups = c.Groups
			}
		case "templates":
			if t, err := db.GetTemplate(eid); err == nil {
				d.EditName = t.Name; d.EditContent = t.Content
			}
		case "autoreply":
			if a, err := db.GetAutoReply(eid); err == nil {
				d.EditKeyword = a.Keyword; d.EditMatch = a.Match; d.EditReply = a.Reply; d.EditAccountID = a.AccountID
				d.EditUseAI = a.UseAI; d.EditAiKeyID = a.AiKeyID; d.EditTrainingID = a.TrainingID
			}
		case "admin_waservers":
			if w, err := db.GetWaServer(eid); err == nil {
				d.EditName = w.Name; d.EditContent = w.URL
				d.EditPhone = w.Port; d.EditKeyword = w.Secret
				d.EditGroups = w.Packages
			}
		case "admin_users":
			if u, err := db.GetUserByID(eid); err == nil {
				d.EditName = u.Name; d.EditPhone = u.Email; d.EditRole = u.Role
			}
		case "admin_roles":
			if r, err := db.GetRole(eid); err == nil {
				d.EditName = r.Name; d.EditContent = r.Permissions
			}
		}
	}

	switch page {
	case "home":
		l, s, r := db.MessageChartData()
		d.ChartLabels = template.JS(l)
		d.ChartSent = template.JS(s)
		d.ChartReceived = template.JS(r)
		d.TotalUsers = db.CountUsers()
		d.ActiveAccounts = engine.CountConnected(uid)
		d.ActiveAccountList = d.ConnectedAccounts
		d.RunningCampaigns, _ = db.CountRunningCampaigns()
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_dashboard"), T("nav_overview"), T("nav_dashboard"), "la-chart-bar"
		if len(d.Sent) > 8 {
			d.Sent = d.Sent[:8]
		}
		if len(d.Received) > 8 {
			d.Received = d.Received[:8]
		}
	case "wa":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_whatsapp"), T("nav_whatsapp"), T("nav_account_qr"), "la-whatsapp"
	case "send":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_send"), T("nav_whatsapp"), T("nav_send"), "la-paper-plane"
	case "sent":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_sent"), T("nav_whatsapp"), T("nav_sent"), "la-telegram"
	case "received":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_received"), T("nav_whatsapp"), T("nav_received"), "la-comment"
	case "inbox":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Live Chat", "Messaging", "Percakapan", "la-comments"
	case "inbox_chat":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Live Chat", "Messaging", "Percakapan", "la-comment"
	case "settings":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_settings"), T("nav_tools"), T("nav_settings"), "la-cog"
	case "autoreply":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_autoreply"), T("nav_tools"), T("nav_autoreply"), "la-robot"
	case "login":
		d.Title = "Login"
	case "register":
		d.Title = "Register"
	case "landing":
		d.Title = d.AppName
	case "contacts":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_contacts_saved"), T("nav_contacts"), T("nav_contacts_saved"), "la-address-book"
	case "groups":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_contacts_groups"), T("nav_contacts"), T("nav_contacts_groups"), "la-list"
	case "unsub":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_contacts_unsub"), T("nav_contacts"), T("nav_contacts_unsub"), "la-unlink"
	case "broadcast":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_broadcast"), T("nav_whatsapp"), T("nav_broadcast"), "la-bullhorn"
	case "scheduled":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_scheduled"), T("nav_whatsapp"), T("nav_scheduled"), "la-clock"
	case "templates":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_templates"), T("nav_tools"), T("nav_templates"), "la-file-alt"
	case "apikeys":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_apikeys"), T("nav_tools"), T("nav_apikeys"), "la-key"
	case "webhooks":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_webhooks"), T("nav_tools"), T("nav_webhooks"), "la-code-branch"
	case "logger":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_logger"), T("nav_tools"), T("nav_logger"), "la-clipboard-list"
	case "hosts_android":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_hosts_android"), T("nav_hosts"), T("nav_hosts_android"), "la-mobile"
	case "hosts_whatsapp":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_hosts_whatsapp"), T("nav_hosts"), T("nav_hosts_whatsapp"), "la-whatsapp"
	case "ussd":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_ussd"), T("nav_android"), T("nav_ussd"), "la-satellite-dish"
	case "ai_keys":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_ai_keys"), T("nav_ai"), T("nav_ai_keys"), "la-key"
	case "ai_plugins":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_ai_plugins"), T("nav_ai"), T("nav_ai_plugins"), "la-plug"
	case "knowledge":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_knowledge"), T("nav_ai"), T("nav_knowledge"), "la-book"
	case "docs":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_docs"), T("nav_docs"), T("nav_docs"), "la-book"
	case "admin":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("nav_admin"), T("nav_admin"), T("nav_admin"), "la-shield-alt"
	case "admin_users":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_users"), T("nav_admin"), T("adm_users"), "la-users"
	case "admin_roles":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_roles"), T("nav_admin"), T("adm_roles"), "la-user-shield"
	case "admin_packages":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_packages"), T("nav_admin"), T("adm_packages"), "la-box"
	case "admin_vouchers":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_vouchers"), T("nav_admin"), T("adm_vouchers"), "la-ticket-alt"
	case "admin_subscriptions":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_subscriptions"), T("nav_admin"), T("adm_subscriptions"), "la-star"
	case "admin_transactions":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_transactions"), T("nav_admin"), T("adm_transactions"), "la-money-bill"
	case "admin_payouts":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_payouts"), T("nav_admin"), T("adm_payouts"), "la-hand-holding-usd"
	case "admin_pages":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_pages"), T("nav_admin"), T("adm_pages"), "la-file"
	case "admin_marketing":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_marketing"), T("nav_admin"), T("adm_marketing"), "la-bullhorn"
	case "admin_languages":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_languages"), T("nav_admin"), T("adm_languages"), "la-language"
	case "admin_waservers":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_waservers"), T("nav_admin"), T("adm_waservers"), "la-server"
	case "admin_gateways":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_gateways"), T("nav_admin"), T("adm_gateways"), "la-code"
	case "admin_shorteners":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_shorteners"), T("nav_admin"), T("adm_shorteners"), "la-link"
	case "admin_plugins":
		d.Title, d.Pretitle, d.Heading, d.Icon = T("adm_plugins"), T("nav_admin"), T("adm_plugins"), "la-puzzle-piece"
	case "admin_meta":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Accounts", "Admin", "Meta Cloud API", "la-cloud"
	case "admin_metatemplates":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Templates", "Admin", "Message Templates", "la-file-alt"
	}

	// parse templates with a language-bound T function
	tpl := template.Must(template.New("").Funcs(template.FuncMap{
		"T":        T,
		"contains": strings.Contains,
		"slice": func(s string, start, end int) string {
			if s == "" || start >= len(s) { return "" }
			if end > len(s) { end = len(s) }
			return s[start:end]
		},
		"js": func(s string) template.JS { return template.JS(s) },
		"permBadges": func(perms string) template.HTML {
			if perms == "" { return "-" }
			parts := strings.Split(perms, ",")
			var buf strings.Builder
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" { continue }
				c := "secondary"
				if strings.HasPrefix(p, "manage_") { c = "primary" }
				if strings.HasPrefix(p, "wa_") { c = "success" }
				buf.WriteString(fmt.Sprintf(`<span class="badge badge-soft-%s me-1" style="line-height:1.4;margin-bottom:2px">%s</span>`, c, p))
			}
			return template.HTML(buf.String())
		},
	}).Parse(templates))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	switch page {
	case "landing":
		if err := tpl.ExecuteTemplate(w, "landing", d); err != nil {
			http.Error(w, err.Error(), 500)
		}
	case "login", "register":
		if err := tpl.ExecuteTemplate(w, "authpage", d); err != nil {
			http.Error(w, err.Error(), 500)
		}
	default:
		if err := tpl.ExecuteTemplate(w, "home", d); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if c, err := r.Cookie("chatgo_sess"); err == nil {
		if uid, ok := db.GetSession(c.Value); ok {
			r.Header.Set("X-User-ID", strconv.FormatInt(uid, 10))
		}
	}
	if _, err := r.Cookie("chatgo_orig"); err == nil {
		r.Header.Set("X-Impersonating", "1")
	}
	uid := getUserID(r)
	if uid == 0 {
		render(w, r, "landing")
		return
	}
	render(w, r, "home")
}

func pageHandler(page string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			return
		}
		render(w, r, page)
	}
}

// /lang/{code} -> set language cookie
func handleLang(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/lang/")
	if i18n.Has(code) {
		http.SetCookie(w, &http.Cookie{Name: "chatgo_lang", Value: code, Path: "/", MaxAge: 31536000})
	}
	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

func handleQRImage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var code string
	if id != "" {
		code = engine.QRFor(id)
	}
	if code == "" {
		http.Error(w, "no qr", 404)
		return
	}
	png, err := qrcode.Encode(code, qrcode.Medium, 300)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(png)
}

// add a new WhatsApp account (multi-number)
func handleWaAdd(w http.ResponseWriter, r *http.Request) {
	id, err := engine.AddAccount(getUserID(r))
	if err != nil {
		http.Redirect(w, r, "/wa?msg="+template.URLQueryEscaper(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/wa?scan="+id, http.StatusSeeOther)
}

// logout a specific account
func handleWaLogout(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id != "" {
		_ = engine.LogoutAccount(id)
	}
	http.Redirect(w, r, "/wa", http.StatusSeeOther)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	status, phone := engine.Status()
	uid := getUserID(r)
	hasQR := false
	for _, a := range engine.Accounts(uid) {
		if a.Status == "qr" {
			hasQR = true
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":%q,"phone":%q,"qr":%v}`, status, phone, hasQR)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, r, "send")
		return
	}
	T := i18n.Translator(currentLang(r))
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	if phone == "" || message == "" {
		http.Redirect(w, r, "/send?msg="+template.URLQueryEscaper(T("send_connect_first")), http.StatusSeeOther)
		return
	}
	if err := engine.SendFrom(strings.TrimPrefix(r.FormValue("account_phone"), "+"), phone, msgtemplate.Render(message, msgtemplate.Vars{Phone: phone})); err != nil {
		http.Redirect(w, r, "/send?msg="+template.URLQueryEscaper(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/send?msg="+template.URLQueryEscaper(T("send_btn")+" OK"), http.StatusSeeOther)
}

func handleSendMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/send", http.StatusSeeOther)
		return
	}
	phone := r.FormValue("phone")
	mediaType := r.FormValue("media_type")
	caption := r.FormValue("caption")
	accountPhone := strings.TrimPrefix(r.FormValue("account_phone"), "+")
	if phone == "" {
		http.Redirect(w, r, "/send?msg=Phone+required", http.StatusSeeOther)
		return
	}
	file, header, err := r.FormFile("media_file")
	if err != nil {
		http.Redirect(w, r, "/send?msg=File+required", http.StatusSeeOther)
		return
	}
	defer file.Close()
	os.MkdirAll("data/media", 0o755)
	dest := filepath.Join("data/media", strconv.FormatInt(time.Now().UnixNano(), 36)+"_"+filepath.Base(header.Filename))
	f, err := os.Create(dest)
	if err != nil { http.Redirect(w, r, "/send?msg=Save+error", http.StatusSeeOther); return }
	io.Copy(f, file)
	f.Close()
	if err := engine.SendMedia(accountPhone, phone, mediaType, dest, caption); err != nil {
		http.Redirect(w, r, "/send?msg="+template.URLQueryEscaper(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/send?msg=Media+sent+OK", http.StatusSeeOther)
}
func handleInboxChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		message := r.FormValue("message")
		accountPhone := strings.TrimPrefix(r.FormValue("account_phone"), "+")
		if phone != "" && message != "" {
			if err := engine.SendFrom(accountPhone, phone, msgtemplate.Render(message, msgtemplate.Vars{Phone: phone})); err != nil {
				http.Redirect(w, r, "/inbox/chat?phone="+phone+"&msg="+template.URLQueryEscaper(err.Error()), http.StatusSeeOther)
				return
			}
		}
		http.Redirect(w, r, "/inbox/chat?phone="+phone+"&msg=OK", http.StatusSeeOther)
		return
	}
	render(w, r, "inbox_chat")
}

func handleInboxEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ch := engine.NotifyChan()
	ctx := r.Context()
	for {
		select {
		case phone := <-ch:
			fmt.Fprintf(w, "data: {\"phone\":%q}\n\n", phone)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func handleInboxSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	accountPhone := strings.TrimPrefix(r.FormValue("account_phone"), "+")
	if phone == "" || message == "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":false,"error":"phone and message required"}`)
		return
	}
	if err := engine.SendFrom(accountPhone, phone, msgtemplate.Render(message, msgtemplate.Vars{Phone: phone})); err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":%q}`, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"ok":true}`)
}

func handleInboxMessages(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "phone required", 400)
		return
	}
	msgs, _ := db.ChatHistory(phone, 100)
	db.MarkRead(phone)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "[")
	for i, m := range msgs {
		if i > 0 { fmt.Fprint(w, ",") }
		fmt.Fprintf(w, `{"type":%q,"id":%d,"phone":%q,"name":%q,"message":%q,"created":%q,"sender_name":%q,"is_group":%v,"channel":%q}`,
			m.Type, m.ID, m.Phone, m.Name, m.Message, m.Created, m.SenderName, m.IsGroup, m.Channel)
	}
	fmt.Fprint(w, "]")
}

func handleInboxUnreadCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"unread":%d}`, db.UnreadCount())
}

func handleInboxMarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	phone := r.FormValue("phone")
	if phone != "" {
		db.MarkRead(phone)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"ok":true}`)
}

func handleInboxSendMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	accountID, _ := strconv.ParseInt(r.FormValue("account_id"), 10, 64)
	if phone == "" || message == "" || accountID == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":false,"error":"phone, message, and account_id required"}`)
		return
	}
	acc, err := db.GetMetaAccount(accountID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":"meta account not found"}`)
		return
	}
	mc := meta.New(acc.PhoneNumberID, acc.AccessToken, acc.VerifyToken)
	_, err = mc.SendText(phone, message)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ok":false,"error":%q}`, err.Error())
		return
	}
	db.LogSent(phone, message, "sent", "meta")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"ok":true}`)
}

func handleInboxSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results, _ := db.SearchInbox(q)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "[")
	for i, c := range results {
		if i > 0 { fmt.Fprint(w, ",") }
		fmt.Fprintf(w, `{"phone":%q,"name":%q,"last_msg":%q,"last_time":%q,"unread":%d,"is_group":%v}`,
			c.Phone, c.Name, c.LastMsg, c.LastTime, c.Unread, c.IsGroup)
	}
	fmt.Fprint(w, "]")
}
func handleAutoReplyAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	keyword := r.FormValue("keyword")
	match := r.FormValue("match")
	reply := r.FormValue("reply")
	useAI := r.FormValue("use_ai") == "on" || r.FormValue("use_ai") == "1"
	aiKeyID, _ := strconv.ParseInt(r.FormValue("ai_key_id"), 10, 64)
	if match == "ai" {
		if faq := r.FormValue("faq"); faq != "" {
			reply = faq
		}
		if keyword == "" {
			keyword = match // "ai"
		}
		if reply == "" && !useAI {
			http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
			return
		}
	} else {
		if keyword == "" || (!useAI && reply == "") {
			http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
			return
		}
	}
	if match == "" {
		match = "contains"
	}
	_, _ = db.AddAutoReply(keyword, match, reply, useAI, aiKeyID, joinVals(r, "account_ids"), func()int64{t,_:=strconv.ParseInt(r.FormValue("training_id"),10,64);return t}())
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}

func handleAutoReplyDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		_ = db.DeleteAutoReply(id)
	}
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}

func handleAutoReplyToggle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		_ = db.ToggleAutoReply(id)
	}
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}
func handleAutoReplyEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if r.Method == http.MethodPost {
		if id > 0 {
			keyword := r.FormValue("keyword")
			match := r.FormValue("match")
			reply := r.FormValue("reply")
			useAI := r.FormValue("use_ai") == "on" || r.FormValue("use_ai") == "1"
			aiKeyID, _ := strconv.ParseInt(r.FormValue("ai_key_id"), 10, 64)
			if match == "" { match = "contains" }
			_ = db.UpdateAutoReply(id, keyword, match, reply, useAI, aiKeyID, joinVals(r, "account_ids"), func()int64{t,_:=strconv.ParseInt(r.FormValue("training_id"),10,64);return t}())
		}
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	render(w, r, "autoreply")
}
func handleContactEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if r.Method == http.MethodPost {
		if id > 0 {
			_ = db.UpdateContact(id, r.FormValue("name"), r.FormValue("phone"), joinVals(r, "groups"))
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	render(w, r, "contacts")
}
func handleTemplateEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if r.Method == http.MethodPost {
		if id > 0 {
			_ = db.UpdateTemplate(id, r.FormValue("name"), r.FormValue("content"))
		}
		http.Redirect(w, r, "/templates", http.StatusSeeOther)
		return
	}
	render(w, r, "templates")
}


func handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, r, "settings")
		return
	}
	setBool := func(name, field string) {
		if r.FormValue(field) == "on" {
			_ = db.SetSetting(name, "1")
		} else {
			_ = db.SetSetting(name, "0")
		}
	}
	setBool("welcome_enabled", "welcome_enabled")
	_ = db.SetSetting("welcome_message", r.FormValue("welcome_message"))
	setBool("fallback_enabled", "fallback_enabled")
	_ = db.SetSetting("fallback_message", r.FormValue("fallback_message"))
	setBool("reply_in_group", "reply_in_group")
	setBool("ai_all_enabled", "ai_all_enabled")
	_ = db.SetSetting("ai_all_key_id", r.FormValue("ai_all_key_id"))
	setBool("handoff_enabled", "handoff_enabled")
	_ = db.SetSetting("handoff_message", r.FormValue("handoff_message"))
	_ = db.SetSetting("handoff_keywords", r.FormValue("handoff_keywords"))
	setBool("ai_fallback_only", "ai_fallback_only")
	_ = db.SetSetting("ai_memory_window", r.FormValue("ai_memory_window"))
	_ = db.SetSetting("ai_delay_seconds", r.FormValue("ai_delay_seconds"))
	_ = db.SetSetting("ai_reasoning_level", r.FormValue("ai_reasoning_level"))
	setBool("biz_hours_enabled", "biz_hours_enabled")
	_ = db.SetSetting("biz_hours_start", r.FormValue("biz_hours_start"))
	_ = db.SetSetting("biz_hours_end", r.FormValue("biz_hours_end"))
	_ = db.SetSetting("biz_hours_off_days", r.FormValue("biz_hours_off_days"))
	setBool("registrations", "registrations")
	_ = db.SetSetting("app_name", r.FormValue("app_name"))
	_ = db.SetSetting("app_email", r.FormValue("app_email"))

	os.MkdirAll("web/assets/theme", 0o755)
	if file, header, err := r.FormFile("logo_file"); err == nil {
		defer file.Close()
		dest := filepath.Join("web/assets/theme", "logo-"+strconv.FormatInt(time.Now().Unix(), 10)+filepath.Ext(header.Filename))
		if f, e := os.Create(dest); e == nil {
			io.Copy(f, file)
			f.Close()
			db.SetSetting("app_logo", "/"+strings.ReplaceAll(dest, "\\", "/"))
		}
	}
	_ = db.SetSetting("app_name", r.FormValue("app_name"))
	setBool("force_own_key", "force_own_key")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

// ---- generic CRUD helpers ----

func crudPost(fn func(*http.Request), redirect string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fn(r)
		}
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}

func crudDel(fn func(int64), redirect string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			fn(id)
		}
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}

func joinVals(r *http.Request, field string) string {
	_ = r.ParseForm()
	vals := r.Form[field]
	return strings.Join(vals, ",")
}

func randSecret() string {
	b := make([]byte, 24)
	_, _ = crand.Read(b)
	return hex.EncodeToString(b)
}

// ---- Broadcast (campaign) ----

func handleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, r, "broadcast")
		return
	}
	name := r.FormValue("name")
	message := r.FormValue("message")
	groups := joinVals(r, "groups")
	accountID := joinVals(r, "account_ids")
	sendMode := r.FormValue("send_mode")
	if sendMode != "round_robin" && sendMode != "random" { sendMode = "round_robin" }
	numbers := strings.TrimSpace(r.FormValue("numbers"))
	if name == "" || message == "" || (groups == "" && numbers == "") {
		http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
		return
	}
	// count unique recipients from groups + direct numbers
	seen := map[string]bool{}
	for _, gid := range strings.Split(groups, ",") {
		list, _ := db.ContactsByGroup(strings.TrimSpace(gid))
		for _, c := range list {
			if c.Phone != "" { seen[c.Phone] = true }
		}
	}
	if numbers != "" {
		for _, n := range strings.Split(numbers, "\n") {
			n = strings.TrimSpace(n)
			// strip +, -, spaces from phone
			n = strings.Map(func(r rune) rune {
				if r == '+' || r == '-' || r == ' ' { return -1 }
				return r
			}, n)
			if n != "" { seen[n] = true }
		}
	}
	// normalize numbers to comma-separated (strip duplicates from groups)
	var numList []string
	for n := range seen {
		numList = append(numList, n)
	}
	normalizedNumbers := strings.Join(numList, ",")

	interval, _ := strconv.Atoi(r.FormValue("interval"))
	if interval <= 0 { interval = 300 }
	_, _ = db.AddCampaign(name, groups, normalizedNumbers, message, len(seen), accountID, sendMode, interval)
	http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
}

// ---- Contact CSV Import ----
func handleContactImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Redirect(w, r, "/contacts?msg=File+required", http.StatusSeeOther)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		http.Redirect(w, r, "/contacts?msg=Invalid+CSV", http.StatusSeeOther)
		return
	}
	colName := -1; colPhone := -1; colGroups := -1
	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch h {
		case "name", "nama": colName = i
		case "phone", "no", "nomor", "telepon": colPhone = i
		case "groups", "group", "grup": colGroups = i
		}
	}
	if colPhone < 0 {
		http.Redirect(w, r, "/contacts?msg=CSV+must+have+phone+column", http.StatusSeeOther)
		return
	}
	// Preload all existing groups into a name→id map
	existingGroups, _ := db.ListGroups()
	gnameToID := map[string]int64{}
	for _, g := range existingGroups { gnameToID[strings.ToLower(strings.TrimSpace(g.Name))] = g.ID }

	imported := 0; skipped := 0
	for {
		record, err := reader.Read()
		if err == io.EOF { break }
		if err != nil { continue }
		phone := strings.TrimSpace(safeGet(record, colPhone))
		name := strings.TrimSpace(safeGet(record, colName))
		if phone == "" { continue }
		// Resolve groups from the CSV — auto-create missing groups
		var gids []string
		groupStr := strings.TrimSpace(safeGet(record, colGroups))
		if groupStr != "" {
			for _, gn := range strings.Split(groupStr, ",") {
				gn = strings.TrimSpace(gn)
				if gn == "" { continue }
				key := strings.ToLower(gn)
				gid, ok := gnameToID[key]
				if !ok {
					id, err := db.AddGroup(gn)
					if err == nil {
						gnameToID[key] = id
						gid = id
					}
				}
				if gid > 0 { gids = append(gids, strconv.FormatInt(gid, 10)) }
			}
		}
		gidStr := strings.Join(gids, ",")
		// Deduplicate by phone
		existing, _ := db.FindContactByPhone(phone)
		if existing != nil {
			skipped++
			continue
		}
		if name == "" { name = phone }
		if _, err := db.AddContact(name, phone, gidStr); err == nil {
			imported++
		}
	}
	msg := fmt.Sprintf("Imported+%d+contacts", imported)
	if skipped > 0 { msg += fmt.Sprintf(",+%d+skipped+(duplicate)", skipped) }
	http.Redirect(w, r, "/contacts?msg="+msg, http.StatusSeeOther)
}

// ---- Scheduled ----

func handleScheduled(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, r, "scheduled")
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	sendAt := r.FormValue("send_at") // "2006-01-02T15:04"
	repeat, _ := strconv.Atoi(r.FormValue("repeat"))
	if name == "" || phone == "" || message == "" || sendAt == "" {
		http.Redirect(w, r, "/scheduled", http.StatusSeeOther)
		return
	}
	sendAt = strings.Replace(sendAt, "T", " ", 1) + ":00"
	accountIDs := joinVals(r, "account_ids")
	_, _ = db.AddScheduled(name, phone, message, sendAt, repeat, accountIDs)
	http.Redirect(w, r, "/scheduled", http.StatusSeeOther)
}

func p(page string) http.HandlerFunc {
	return authMiddleware(pageHandler(page))
}

func pageFromQuery(r *http.Request) int {
	p, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if p < 1 { p = 1 }
	return p
}
func pageNums(current, total int) []int {
	var out []int
	start := current - 2
	if start < 1 { start = 1 }
	end := start + 4
	if end > total { end = total; start = end - 4; if start < 1 { start = 1 } }
	for i := start; i <= end; i++ { out = append(out, i) }
	return out
}

func handleMetaWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mode := r.URL.Query().Get("hub.mode")
		challenge := r.URL.Query().Get("hub.challenge")
		vt := r.URL.Query().Get("hub.verify_token")
		accounts, _ := db.ListMetaAccounts()
		for _, acc := range accounts {
			if acc.VerifyToken == vt && mode == "subscribe" {
				fmt.Fprint(w, challenge)
				return
			}
		}
		http.Error(w, "verification failed", 403)
		return
	}
	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		msgs, phoneNumberID, ok := meta.ParseWebhook(body, "")
		if !ok {
			return
		}
		accs, _ := db.ListMetaAccounts()
		for _, acc := range accs {
			if acc.PhoneNumberID != phoneNumberID {
				continue
			}
			for _, m := range msgs {
				text := m.Text.Body
				if m.Interactive != nil && m.Interactive.ButtonReply != nil {
					text = m.Interactive.ButtonReply.Title
				}
				if text == "" {
					continue
				}
				db.LogReceived(m.From, "", text, false, "", "", "meta")
				db.MarkRead(m.From)
				engine.Notify(m.From)
				db.Log("meta", "received", fmt.Sprintf("%s -> %s: %s", m.From, acc.Name, text))

				// Auto-reply for Meta
				mc := meta.New(acc.PhoneNumberID, acc.AccessToken, acc.VerifyToken)
				if ar, found := db.FindReplyFullForAccount(text, ""); found && ar.IsActive {
					reply := msgtemplate.Render(ar.Reply, msgtemplate.Vars{Phone: m.From, Name: "", Message: text})
					if ar.UseAI && ar.AiKeyID > 0 {
						if aik, err := db.GetAiKey(ar.AiKeyID); err == nil {
							decKey, _ := secret.Decrypt(aik.APIKey)
							if decKey == "" { decKey = aik.APIKey }
							if aiReply, aiErr := aiservice.Reply(decKey, aik.Provider, aik.Model, aik.BaseURL, aik.SystemPrompt, text, nil, nil); aiErr == nil && aiReply != "" {
								reply = aiReply
							}
						}
					}
					if reply != "" {
						mc.SendText(m.From, reply)
						db.LogSent(m.From, reply, "autoreply", "meta")
						db.Log("meta", "autoreply", fmt.Sprintf("-> %s: %s", m.From, reply))
					}
				}
			}
			break
		}
		fmt.Fprint(w, "ok")
		return
	}
	http.Error(w, "method not allowed", 405)
}











