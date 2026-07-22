package main

import (
	crand "crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
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
	"chatgo/payment"
	"chatgo/pseo"
	"chatgo/secret"
	"chatgo/store"
	"chatgo/wa"
	"chatgo/telegram"

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
	db.EnsureDefaultRoles()

	engine, err = wa.New(filepath.Join(dataDir, "session.db"), db)
	if err != nil {
		log.Fatalf("wa engine: %v", err)
	}
	if err := engine.Start(); err != nil {
		log.Printf("wa start: %v (will retry via QR page)", err)
	}
	engine.StartLoops()

	pseo.Init(getEnv("APP_NAME", "ChatGo"), appURL(), "6281296052010")
	pseo.InitIndexNow()
	setupProEngine()

	mux := http.NewServeMux()
	noDirFS := func(dir string) http.Handler {
		fs := http.FileServer(http.Dir(dir))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r); return
			}
			fs.ServeHTTP(w, r)
		})
	}
	mux.Handle("/assets/", http.StripPrefix("/assets/", noDirFS("web/assets")))
	mux.Handle("/web/", http.StripPrefix("/web/", noDirFS("web")))
	mux.Handle("/screens/", http.StripPrefix("/screens/", noDirFS("public/marketing/screens")))

	// Sitemap, robots.txt, IndexNow key
	mux.HandleFunc("/sitemap.xml", pseo.HandleSitemap)
	mux.HandleFunc("/sitemaps/", pseo.HandleSitemapPage)
	mux.HandleFunc("/robots.txt", pseo.HandleRobots)
	mux.HandleFunc("/indexnow-key.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, pseo.GetIndexNowKey())
	})

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
	mux.HandleFunc("/settings", authMiddleware(requireAdmin(handleSettings)))
	mux.HandleFunc("/admin/users/impersonate", authMiddleware(handleImpersonate))
	mux.HandleFunc("/exit-impersonation", handleExitImpersonation)
	mux.HandleFunc("/contacts", p("contacts"))
	mux.HandleFunc("/contacts/groups", p("groups"))
	mux.HandleFunc("/contacts/unsub", p("unsub"))
	mux.HandleFunc("/contacts/add", authMiddleware(limitGuard("contact", func(r *http.Request) {
		uid := getUserID(r)
		id, _ := db.AddContact(uid, r.FormValue("name"), r.FormValue("phone"), joinVals(r, "groups"))
		var tagIDs []int64
		for _, s := range r.Form["tag_ids"] {
			if tid, err := strconv.ParseInt(s, 10, 64); err == nil { tagIDs = append(tagIDs, tid) }
		}
		if len(tagIDs) > 0 { db.SetContactTags(id, tagIDs) }
	}, "/contacts")))
	mux.HandleFunc("/contacts/delete", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			db.DeleteContact(getUserID(r), id)
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	}))
	mux.HandleFunc("/contacts/import", authMiddleware(handleContactImport))
	mux.HandleFunc("/contacts/export", authMiddleware(handleContactExport))
	mux.HandleFunc("/contacts/bulk-delete", authMiddleware(handleContactBulkDelete))
	mux.HandleFunc("/groups/add", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			uid := getUserID(r)
			db.AddGroup(uid, r.FormValue("name"))
		}
		http.Redirect(w, r, "/contacts/groups", http.StatusSeeOther)
	}))
	mux.HandleFunc("/groups/delete", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			db.DeleteGroup(getUserID(r), id)
		}
		http.Redirect(w, r, "/contacts/groups", http.StatusSeeOther)
	}))
	mux.HandleFunc("/unsub/add", authMiddleware(crudPost(func(r *http.Request) { db.AddUnsub(r.FormValue("phone")) }, "/contacts/unsub")))
	mux.HandleFunc("/unsub/delete", authMiddleware(crudDel(func(id int64) { db.DeleteUnsub(id) }, "/contacts/unsub")))
	mux.HandleFunc("/broadcast", authMiddleware(handleBroadcast))
	mux.HandleFunc("/broadcast/stop", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 { db.UpdateCampaignStatus(getUserID(r), id, "stopped") }
		http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
	}))
	mux.HandleFunc("/broadcast/pause", authMiddleware(handleCampaignPause))
	mux.HandleFunc("/broadcast/retry", authMiddleware(handleCampaignRetry))
	mux.HandleFunc("/broadcast/delete", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 { db.DeleteCampaign(getUserID(r), id) }
		http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
	}))
	mux.HandleFunc("/drips", authMiddleware(handleDrips))
	mux.HandleFunc("/drips/add", authMiddleware(handleDripAdd))
	mux.HandleFunc("/drips/step/add", authMiddleware(handleDripStepAdd))
	mux.HandleFunc("/drips/step/delete", authMiddleware(crudDel(func(id int64) { db.DeleteDripStep(id) }, "/drips")))
	mux.HandleFunc("/drips/delete", authMiddleware(crudDel(func(id int64) { db.DeleteDrip(id) }, "/drips")))
		mux.HandleFunc("/drips/toggle", authMiddleware(handleDripToggle))
	mux.HandleFunc("/tags", authMiddleware(handleTags))
	mux.HandleFunc("/tags/add", authMiddleware(crudPost(func(r *http.Request) { db.AddTag(r.FormValue("name"), r.FormValue("color")) }, "/tags")))
	mux.HandleFunc("/tags/delete", authMiddleware(crudDel(func(id int64) { db.DeleteTag(id) }, "/tags")))
	mux.HandleFunc("/contacts/tags", authMiddleware(handleContactTags))
	mux.HandleFunc("/canned", authMiddleware(handleCanned))
	mux.HandleFunc("/canned/add", authMiddleware(limitGuard("canned", func(r *http.Request) { db.AddCanned(r.FormValue("shortcut"), r.FormValue("name"), r.FormValue("message")) }, "/canned")))
	mux.HandleFunc("/canned/delete", authMiddleware(crudDel(func(id int64) { db.DeleteCanned(id) }, "/canned")))
	mux.HandleFunc("/inbox/assign", authMiddleware(handleInboxAssign))
	mux.HandleFunc("/inbox/close", authMiddleware(handleInboxClose))
	mux.HandleFunc("/tracker", authMiddleware(p("tracker")))
	mux.HandleFunc("/ab-tests", authMiddleware(p("abtests")))
	mux.HandleFunc("/ab-tests/add", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cid, _ := strconv.ParseInt(r.FormValue("campaign_id"), 10, 64)
			if cid > 0 {
				db.CreateABTest(cid, r.FormValue("variant_a"), r.FormValue("variant_b"))
			}
		}
		http.Redirect(w, r, "/ab-tests", http.StatusSeeOther)
	}))
	mux.HandleFunc("/track/", handleLinkTrack)
	mux.HandleFunc("/subscribe", authMiddleware(handleSubscribe))
	mux.HandleFunc("/subscribe/checkout", authMiddleware(handleCheckout))
	mux.HandleFunc("/payment/callback/", handlePaymentCallback)
	mux.HandleFunc("/store", authMiddleware(p("store")))
	mux.HandleFunc("/store/add", authMiddleware(crudPost(func(r *http.Request) { p, _ := strconv.ParseFloat(r.FormValue("price"), 64); s, _ := strconv.Atoi(r.FormValue("stock")); db.AddProduct(r.FormValue("name"), r.FormValue("desc"), p, r.FormValue("image_url"), r.FormValue("category"), s) }, "/store")))
	mux.HandleFunc("/store/delete", authMiddleware(crudDel(func(id int64) { db.DeleteProduct(id) }, "/store")))
	mux.HandleFunc("/store/category/add", authMiddleware(crudPost(func(r *http.Request) { db.AddCategory(r.FormValue("name")) }, "/store")))
	mux.HandleFunc("/store/category/delete", authMiddleware(crudDel(func(id int64) { db.DeleteCategory(id) }, "/store")))
	mux.HandleFunc("/store/orders", authMiddleware(p("orders")))
	mux.HandleFunc("/store/orders/update", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		status := r.FormValue("status")
		if id > 0 { db.UpdateOrderStatus(id, status) }
		// WA notif to customer
		phone := r.FormValue("phone")
		if phone != "" {
			uid := getUserID(r)
			s := engine.FirstSession(uid)
			if s != nil {
				msg := fmt.Sprintf("Order #%d: *%s*\nStatus: %s", id, r.FormValue("product"), status)
				engine.SendFrom(uid, s.Phone, phone, msg)
			}
		}
		http.Redirect(w, r, "/store/orders", http.StatusSeeOther)
	}))
	mux.HandleFunc("/forms", authMiddleware(p("forms")))
	mux.HandleFunc("/forms/add", authMiddleware(limitGuard("form", func(r *http.Request) { db.AddForm(r.FormValue("name"), r.FormValue("fields")) }, "/forms")))
	mux.HandleFunc("/forms/delete", authMiddleware(crudDel(func(id int64) { db.DeleteForm(id) }, "/forms")))
	mux.HandleFunc("/forms/submissions", authMiddleware(p("submissions")))
	mux.HandleFunc("/reminders", authMiddleware(p("reminders")))
	mux.HandleFunc("/reminders/add", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		amt, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		db.AddReminder(r.FormValue("phone"), r.FormValue("name"), amt, r.FormValue("due_date"), r.FormValue("message"))
		http.Redirect(w, r, "/reminders", http.StatusSeeOther)
	}))
	mux.HandleFunc("/analytics", authMiddleware(handleAnalytics))
	mux.HandleFunc("/blacklist", authMiddleware(p("blacklist")))
	mux.HandleFunc("/blacklist/add", authMiddleware(crudPost(func(r *http.Request) { db.AddBlacklist(r.FormValue("phone"), r.FormValue("reason")) }, "/blacklist")))
	mux.HandleFunc("/blacklist/remove", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			db.RemoveBlacklist(r.FormValue("phone"))
		}
		http.Redirect(w, r, "/blacklist", http.StatusSeeOther)
	}))
	mux.HandleFunc("/csat", authMiddleware(p("csat")))
	mux.HandleFunc("/groups/language", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		gid, _ := strconv.ParseInt(r.FormValue("group_id"), 10, 64); db.SetGroupLanguage(gid, r.FormValue("language"))
		http.Redirect(w, r, "/contacts/groups", http.StatusSeeOther)
	}))
	mux.HandleFunc("/validate", authMiddleware(handleValidate))
	mux.HandleFunc("/depts", authMiddleware(p("depts")))
	mux.HandleFunc("/depts/add", authMiddleware(crudPost(func(r *http.Request) { db.AddDept(r.FormValue("name"), joinVals(r, "agents")) }, "/depts")))
	mux.HandleFunc("/depts/delete", authMiddleware(crudDel(func(id int64) { db.DeleteDept(id) }, "/depts")))
	mux.HandleFunc("/inbox/note", authMiddleware(handleInboxNote))
	mux.HandleFunc("/inbox/transfer", authMiddleware(handleInboxTransfer))
	mux.HandleFunc("/recurring", authMiddleware(p("recurring")))
	mux.HandleFunc("/recurring/add", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			uid := getUserID(r)
			if db.CountUserRecurring(uid) >= db.GetUserRecurringLimit(uid) {
				http.Redirect(w, r, "/recurring?msg=Recurring+limit+reached.+Upgrade+your+plan.", http.StatusSeeOther)
				return
			}
			dow, _ := strconv.Atoi(r.FormValue("day_of_week")); hr, _ := strconv.Atoi(r.FormValue("hour"))
			db.AddRecurring(r.FormValue("name"), joinVals(r, "groups"), r.FormValue("message"), dow, hr)
		}
		http.Redirect(w, r, "/recurring", http.StatusSeeOther)
	}))
	mux.HandleFunc("/recurring/delete", authMiddleware(crudDel(func(id int64) { db.DeleteRecurring(id) }, "/recurring")))
	mux.HandleFunc("/recurring/toggle", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64); db.ToggleRecurring(id)
		http.Redirect(w, r, "/recurring", http.StatusSeeOther)
	}))
	mux.HandleFunc("/uploads", authMiddleware(p("uploads")))
	mux.HandleFunc("/widget.js", handleWidgetJS)
	mux.HandleFunc("/widget/chat", handleWidgetChat)
	mux.HandleFunc("/inbox/label", authMiddleware(handleInboxLabel))
	mux.HandleFunc("/inbox/filter/", authMiddleware(handleInboxFilter))
	mux.HandleFunc("/inbox/canned", authMiddleware(handleInboxCanned))
	mux.HandleFunc("/inbox/star", authMiddleware(handleInboxStar))
	mux.HandleFunc("/inbox/suggest", authMiddleware(handleInboxSuggest))
	mux.HandleFunc("/faq", authMiddleware(handleFAQ))
	mux.HandleFunc("/faq/add", authMiddleware(handleFAQAdd))
	mux.HandleFunc("/faq/delete", authMiddleware(handleFAQDelete))
	mux.HandleFunc("/faq/import", authMiddleware(handleFAQImport))
	mux.HandleFunc("/customers", authMiddleware(p("customers")))
	mux.HandleFunc("/customers/profile", authMiddleware(handleCustomerProfile))
	mux.HandleFunc("/calendar", authMiddleware(p("calendar")))
	mux.HandleFunc("/backup", authMiddleware(handleBackup))
	mux.HandleFunc("/translate", authMiddleware(handleTranslate))
	mux.HandleFunc("/macros", authMiddleware(p("macros")))
	mux.HandleFunc("/macros/add", authMiddleware(limitGuard("macro", func(r *http.Request) { db.AddMacro(r.FormValue("name"), r.FormValue("actions")) }, "/macros")))
	mux.HandleFunc("/macros/delete", authMiddleware(crudDel(func(id int64) { db.DeleteMacro(id) }, "/macros")))
	mux.HandleFunc("/macros/execute", authMiddleware(handleMacroExecute))
	mux.HandleFunc("/merge", authMiddleware(p("merge")))
	mux.HandleFunc("/merge/execute", authMiddleware(handleMergeExecute))
	mux.HandleFunc("/priority/set", authMiddleware(handleSetPriority))
	mux.HandleFunc("/audit", authMiddleware(p("audit")))
	mux.HandleFunc("/email-webhook", handleEmailWebhook)
	mux.HandleFunc("/translate-tool", authMiddleware(p("translatetool")))
	mux.HandleFunc("/widget-info", authMiddleware(p("widgetinfo")))
	mux.HandleFunc("/email-wa", authMiddleware(p("emailwa")))
	mux.HandleFunc("/meta/send", authMiddleware(handleMetaSend))
	mux.HandleFunc("/meta/campaigns", authMiddleware(p("meta_campaigns")))
	mux.HandleFunc("/meta/inbox", authMiddleware(p("meta_inbox")))
	mux.HandleFunc("/meta/logs", authMiddleware(p("meta_logs")))
	mux.HandleFunc("/meta/analytics", authMiddleware(p("meta_analytics")))
	mux.HandleFunc("/meta/webhook", authMiddleware(p("meta_webhook")))
		mux.HandleFunc("/scheduled", authMiddleware(handleScheduled))
	mux.HandleFunc("/scheduled/delete", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 { db.DeleteScheduled(getUserID(r), id) }
		http.Redirect(w, r, "/scheduled", http.StatusSeeOther)
	}))
	mux.HandleFunc("/templates", p("templates"))
	mux.HandleFunc("/templates/add", authMiddleware(limitGuard("template", func(r *http.Request) { db.AddTemplate(r.FormValue("name"), r.FormValue("content")) }, "/templates")))
	mux.HandleFunc("/templates/delete", authMiddleware(crudDel(func(id int64) { db.DeleteTemplate(id) }, "/templates")))
	mux.HandleFunc("/apikeys", p("apikeys"))
	mux.HandleFunc("/apikeys/add", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			uid := getUserID(r)
			if db.CountUserApiKeys(uid) >= db.GetUserKeyLimit(uid) {
				http.Redirect(w, r, "/apikeys?msg=API+Key+limit+reached.+Upgrade+your+plan.", http.StatusSeeOther)
				return
			}
			secret := randSecret()
			db.AddAPIKey(r.FormValue("name"), secret)
			http.Redirect(w, r, "/apikeys?msg=API+Key+dibuat:+("+template.URLQueryEscaper(secret)+")+simpan+sekarang,+tidak+akan+ditampilkan+lagi", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/apikeys", http.StatusSeeOther)
	}))
	mux.HandleFunc("/apikeys/delete", authMiddleware(crudDel(func(id int64) { db.DeleteAPIKey(id) }, "/apikeys")))
	mux.HandleFunc("/webhooks", p("webhooks"))
	mux.HandleFunc("/webhooks/add", authMiddleware(limitGuard("webhook", func(r *http.Request) { db.AddWebhook(r.FormValue("name"), r.FormValue("url"), r.FormValue("event")) }, "/webhooks")))
	mux.HandleFunc("/webhooks/delete", authMiddleware(crudDel(func(id int64) { db.DeleteWebhook(id) }, "/webhooks")))
	mux.HandleFunc("/logger", p("logger"))
	mux.HandleFunc("/logger/clear", authMiddleware(func(w http.ResponseWriter, r *http.Request) { db.ClearLog(); http.Redirect(w, r, "/logger", http.StatusSeeOther) }))
	registerAdminRoutes(mux)
	initProRoutes(mux)
	mux.HandleFunc("/lang/", handleLang)
	mux.HandleFunc("/qr.png", handleQRImage)
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/webhook/meta", handleMetaWebhook)
	mux.HandleFunc("/autoreply/add", authMiddleware(handleAutoReplyAdd))
	mux.HandleFunc("/autoreply/delete", authMiddleware(handleAutoReplyDelete))
	mux.HandleFunc("/autoreply/toggle", authMiddleware(handleAutoReplyToggle))
	mux.HandleFunc("/autoreply/edit", authMiddleware(handleAutoReplyEdit))
	mux.HandleFunc("/upgrade", authMiddleware(func(w http.ResponseWriter, r *http.Request) { render(w, r, "upgrade") }))
	mux.HandleFunc("/pro/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		render(w, r, "upgrade")
	}))

	// Contact edit
	mux.HandleFunc("/contact/edit", authMiddleware(handleContactEdit))

	// Meta admin API routes
	mux.HandleFunc("/admin/meta/flows", authMiddleware(requireAdmin(handleMetaFlows)))
	mux.HandleFunc("/admin/meta/catalog", authMiddleware(requireAdmin(handleMetaCatalog)))
	mux.HandleFunc("/admin/meta/calling", authMiddleware(requireAdmin(handleMetaCalling)))
	mux.HandleFunc("/admin/meta/profile", authMiddleware(requireAdmin(handleMetaProfile)))
	mux.HandleFunc("/admin/meta/qr", authMiddleware(requireAdmin(handleMetaQR)))
	mux.HandleFunc("/admin/meta/health", authMiddleware(requireAdmin(handleMetaHealth)))
	mux.HandleFunc("/admin/meta/templates", authMiddleware(requireAdmin(handleMetaTemplates)))
	mux.HandleFunc("/admin/meta/carousel", authMiddleware(requireAdmin(handleMetaCarousel)))
	mux.HandleFunc("/admin/meta/webhooks", authMiddleware(requireAdmin(handleMetaWebhooks)))
	mux.HandleFunc("/admin/meta/payment", authMiddleware(requireAdmin(handleMetaPayment)))
	mux.HandleFunc("/admin/meta/register", authMiddleware(requireAdmin(handleMetaRegister)))
	mux.HandleFunc("/admin/meta/insights", authMiddleware(requireAdmin(handleMetaInsights)))
	mux.HandleFunc("/sheets", authMiddleware(handleSheets))
	mux.HandleFunc("/n8n-webhook", handleN8NWebhook)
	mux.HandleFunc("/n8n-node-definition", handleN8NNodeDefinition)
	mux.HandleFunc("/n8n-templates", handleN8NTemplates)
	mux.HandleFunc("/ig-webhook", handleIGWebhook)
	mux.HandleFunc("/admin/instagram", authMiddleware(requireAdmin(handleIGInbox)))
	mux.HandleFunc("/agency", authMiddleware(requireAdmin(handleAgency)))
	mux.HandleFunc("/ai-settings", authMiddleware(handleAISettings))
	mux.HandleFunc("/buttons-builder", authMiddleware(handleButtonsBuilder))
	mux.HandleFunc("/warmer", authMiddleware(requireAdmin(handleWarmer)))
	mux.HandleFunc("/flow-search", authMiddleware(handleFlowSearch))
	mux.HandleFunc("/dark-mode", authMiddleware(handleDarkMode))
	mux.HandleFunc("/flow-logs", authMiddleware(handleFlowLogs))
	mux.HandleFunc("/telegram-webhook", handleTelegramWebhook)
	mux.HandleFunc("/fb-webhook", handleFBWebhook)
	mux.HandleFunc("/omni/inbox", authMiddleware(handleOmniInbox))
	mux.HandleFunc("/omni/analytics", authMiddleware(handleOmniAnalytics))
	mux.HandleFunc("/omni/handoff", authMiddleware(handleOmniHandoff))

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
	fmt.Printf("\n  %s running at http://%s\n\n", getEnv("APP_NAME", "ChatGo"), addr)
	var handler http.Handler = mux
	handler = csrfMiddleware(handler)
	log.Fatal(http.ListenAndServe(addr, handler))
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
	IsAdmin                      bool
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
	EditPrice         string
	EditSendLimit     int
	EditDeviceLimit   int
	EditWaAccountLimit int
	EditContactLimit  int
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
	RateMaxDaily      string
	RateRandomMin     string
	RateRandomMax     string
	AutoCloseHours    string
	AutoCloseMessage  string
	AgentSignature    string
	BizHoursReply     string
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
	Drips      []store.Drip
	Tags       []store.Tag
	Canned     []store.CannedResponse
	LClicks    []store.LinkClick
	ABTests        []store.ABTest
	PaymentGateways []store.PaymentGateway
	Txs            []store.PaymentTransaction
	Products       []store.Product
	Categories     []store.ProductCategory
	Orders         []store.Order
	Forms          []store.ChatForm
	Submissions    []store.FormSubmission
	Reminders      []store.PaymentReminder
	AgentMetrics   []store.AgentMetric
	Blacklist      []store.BlacklistEntry
	CSATAvg        float64
	CSATCount      int
	Depts          []store.Department
	Recurrings     []store.RecurringCampaign
	Notes          []store.ChatNote
	CalEvents      []store.CalendarEvent
	Files          []string
	Macros         []store.InboxMacro
	Duplicates     []map[string]interface{}
	AuditLogs      []store.AuditLog
	Scheduleds     []store.Scheduled
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
	IsImpersonating    bool
	UserPackage        string
	UserPackageServices string
	UserPackageExpire   string
	AiKeys        []store.AiKey
	AiPlugins     []store.AiPlugin
	AiTrainings   []store.AiTraining
	Devices       []store.Device
	Ussds         []store.Ussd
	Knowledges    []store.KnowledgeEntry
	FAQ           []map[string]string
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
	ars, _ := db.ListAutoReplies(uid)
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
		d.IsAdmin = db.HasPermission(uid, "manage_users")
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
	d.RateMaxDaily = db.GetSetting("rate_max_daily", "0")
	d.RateRandomMin = db.GetSetting("rate_random_min", "0")
	d.RateRandomMax = db.GetSetting("rate_random_max", "0")
	d.AutoCloseHours = db.GetSetting("auto_close_hours", "0")
	d.AutoCloseMessage = db.GetSetting("auto_close_message", "Chat ini ditutup otomatis. Silakan hubungi kami kembali jika perlu bantuan.")
	d.AgentSignature = db.GetSetting("agent_signature", "")
	d.BizHoursReply = db.GetSetting("biz_hours_reply", "Saat ini di luar jam operasional. Pesan Anda akan kami balas saat jam kerja.")
	d.ForceOwnKey = db.GetSetting("force_own_key", "0") == "1"
	d.Registrations = db.GetSetting("registrations", "1") == "1"
	d.AiTokenQuota = int64(db.GetUserAiQuota(uid))
	d.AiTokenUsed = db.GetAiTokenUsage(uid)
	d.UnreadCount = db.UnreadCount()
	d.IsImpersonating = r.Header.Get("X-Impersonating") == "1"
	// load user's active subscription and package
	d.UserPackage = ""
	d.UserPackageServices = ""
	d.UserPackageExpire = ""
	if uid > 0 {
		if sub, err := db.GetActiveSubscription(uid); err == nil {
			d.UserPackageExpire = sub.Expire
			pkgID, _ := strconv.ParseInt(sub.Pkg, 10, 64)
			var pkg *store.Package
			if pkgID > 0 {
				pkg, _ = db.GetPackage(pkgID)
			} else {
				pkg, _ = db.GetPackageByName(sub.Pkg)
			}
			if pkg != nil {
				d.UserPackage = pkg.Name
				d.UserPackageServices = pkg.Services
			}
		} else {
			d.UserPackage = "Free"
		}
	}

	// load entity lists per page (only what's needed)
	switch page {
	case "contacts":
		d.Contacts, _ = db.ListContacts(uid)
		d.Groups, _ = db.ListGroups(uid)
		d.Tags, _ = db.ListTags()
	case "groups":
		d.Groups, _ = db.ListGroups(uid)
	case "unsub":
		d.Unsubs, _ = db.ListUnsub()
	case "templates":
		d.Templates, _ = db.ListTemplates()
	case "apikeys":
		d.APIKeys, _ = db.ListAPIKeys()
	case "autoreply":
		d.AiKeys, _ = db.ListAiKeys(uid)
		d.Knowledges, _ = db.ListKnowledge()
		d.AiTrainings, _ = db.ListAiTrainings()
		d.FAQ, _ = db.ListFAQ(uid)
		d.AiPlugins, _ = db.ListAiPlugins()
	case "webhooks":
		d.Webhooks, _ = db.ListWebhooks()
	case "broadcast":
		d.Campaigns, _ = db.ListCampaigns(uid)
		d.Groups, _ = db.ListGroups(uid)
		d.MetaAccounts, _ = db.ListMetaAccounts()
		d.MetaTemplates, _ = db.ListMetaTemplates()
	case "drips":
		d.Drips, _ = db.ListDrips()
	case "tags":
		d.Tags, _ = db.ListTags()
	case "tracker":
		d.LClicks, _ = db.ListLinkClicks()
	case "abtests":
		d.ABTests, _ = db.ListABTests()
		d.Campaigns, _ = db.ListCampaigns(uid)
	case "subscribe":
		d.Packages, _ = db.ListPackages()
		d.PaymentGateways, _ = db.ListPaymentGateways()
	case "admin_paygateways":
		d.PaymentGateways, _ = db.ListPaymentGateways()
	case "admin_transactions_pay":
		d.Txs, _ = db.ListPayTransactions()
	case "store":
		d.Products, _ = db.ListProducts()
		d.Categories, _ = db.ListCategories()
	case "orders":
		d.Orders, _ = db.ListOrders()
	case "forms":
		d.Forms, _ = db.ListForms()
	case "submissions":
		d.Forms, _ = db.ListForms()
		fid, _ := strconv.ParseInt(r.URL.Query().Get("form_id"), 10, 64)
		if fid > 0 { d.Submissions, _ = db.ListSubmissions(fid) }
	case "reminders":
		d.Reminders, _ = db.ListReminders()
	case "analytics":
		d.AgentMetrics = db.AgentMetrics()
	case "csat":
		d.CSATAvg = db.CSATAverage(30)
		d.CSATCount = db.CSATCount()
	case "depts":
		d.Depts, _ = db.ListDepts()
		d.Users, _ = db.ListUsers()
	case "recurring":
		d.Recurrings, _ = db.ListRecurring()
		d.Groups, _ = db.ListGroups(uid)
	case "uploads":
		d.Files = db.ListUploads("public/uploads")
	case "customers":
		d.Contacts, _ = db.ListContacts(uid)
	case "calendar":
		d.CalEvents = db.GetCalendarEvents()
	case "macros":
		d.Macros, _ = db.ListMacros()
	case "merge":
		d.Duplicates = db.FindDuplicateContacts()
	case "audit":
		d.AuditLogs, _ = db.ListAuditLogs()
	case "canned":
		d.Canned, _ = db.ListCanned()
		d.Users, _ = db.ListUsers()
	case "scheduled":
		d.Scheduleds, _ = db.ListScheduled()
	case "sent":
		d.SentPage = pageFromQuery(r)
		d.Sent, _ = db.ListSentPaginated(d.SentPage, d.SentPerPage)
		d.SentTotal = db.CountSent()
		d.SentPages = pageNums(d.SentPage, (d.SentTotal+d.SentPerPage-1)/d.SentPerPage)
	case "inbox":
		d.PageNum = pageFromQuery(r)
		d.InboxConversations, _ = db.GroupInboxPaginated(uid, d.PageNum, 10)
		d.InboxTotal = db.CountInbox(uid)
		d.InboxPages = pageNums(d.PageNum, (d.InboxTotal+9)/10)
		d.Statuses, _ = db.ListStatuses()
		d.Canned, _ = db.ListCanned()
	case "inbox_chat":
		d.Phone = r.URL.Query().Get("phone")
		d.ChatMessages, _ = db.ChatHistory(d.Phone, 100)
		d.Templates, _ = db.ListTemplates()
		d.Canned, _ = db.ListCanned()
		d.Users, _ = db.ListUsers()
		d.Notes, _ = db.GetNotes(d.Phone)
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
		d.AiKeys, _ = db.ListAiKeys(uid)
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
	case "faq":
		d.FAQ, _ = db.ListFAQ(uid)
	case "docs":
		d.DocsSteps = allDocsSteps
	}

	// edit mode – pre-fill forms
	if eid, _ := strconv.ParseInt(r.URL.Query().Get("edit"), 10, 64); eid > 0 {
		d.EditID = eid
		switch page {
		case "contacts":
			if c, err := db.GetContact(uid, eid); err == nil {
				d.EditName = c.Name; d.EditPhone = c.Phone; d.EditGroups = c.Groups
			}
		case "templates":
			if t, err := db.GetTemplate(eid); err == nil {
				d.EditName = t.Name; d.EditContent = t.Content
			}
		case "autoreply":
			if a, err := db.GetAutoReply(uid, eid); err == nil {
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
		case "admin_packages":
			if p, err := db.GetPackage(eid); err == nil {
				d.EditName = p.Name; d.EditPrice = p.Price
				d.EditSendLimit = p.SendLimit; d.EditDeviceLimit = p.DeviceLimit
				d.EditWaAccountLimit = p.WaAccountLimit; d.EditContactLimit = p.ContactLimit
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
	case "drips":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Drip Campaign", T("nav_whatsapp"), "Drip Campaign", "la-tint"
	case "tags":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Tags", T("nav_contacts"), "Contact Tags", "la-tags"
	case "canned":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Canned Responses", T("nav_tools"), "Canned Responses", "la-comment-dots"
	case "store":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Store", "Products", "Product Catalog", "la-store"
	case "orders":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Orders", "Store", "Orders", "la-shopping-bag"
	case "forms":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Forms", "Tools", "Interactive Forms", "la-wpforms"
	case "submissions":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Submissions", "Forms", "Form Data", "la-database"
	case "reminders":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Reminders", "Tools", "Payment Reminders", "la-bell"
	case "analytics":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Analytics", "Reports", "Conversation Analytics", "la-chart-pie"
	case "blacklist":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Blacklist", "Safety", "Blocked Numbers", "la-ban"
	case "csat":
		d.Title, d.Pretitle, d.Heading, d.Icon = "CSAT", "Reports", "Customer Satisfaction", "la-star"
	case "depts":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Departments", "Admin", "Departments", "la-building"
	case "recurring":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Recurring", "Broadcast", "Auto-Repeat", "la-redo-alt"
	case "uploads":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Files", "Tools", "File Manager", "la-folder-open"
	case "customers":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Customers", "CRM", "Customer Directory", "la-users"
	case "calendar":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Calendar", "Schedule", "Campaign Calendar", "la-calendar"
	case "backup":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Backup", "System", "Database Backup", "la-database"
	case "macros":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Macros", "Tools", "Inbox Macros", "la-bolt"
	case "merge":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Merge", "Contacts", "Merge Duplicates", "la-code-branch"
	case "audit":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Audit", "System", "Audit Log", "la-history"
	case "translatetool":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Translate", "Tools", "Auto Translate", "la-language"
	case "widgetinfo":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Widget", "Tools", "Web Widget", "la-code"
	case "emailwa":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Email→WA", "Tools", "Email Gateway", "la-envelope"
	case "meta_send":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Send", "Meta", "Send via Cloud API", "la-paper-plane"
		d.MetaAccounts, _ = db.ListMetaAccounts()
	case "meta_campaigns":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Campaigns", "Meta", "Campaign via Cloud API", "la-bullhorn"
		d.Campaigns, _ = db.ListCampaigns(uid)
	case "meta_inbox":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Inbox", "Meta", "Meta Live Chat", "la-comments"
	case "meta_logs":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Logs", "Meta", "Webhook Activity", "la-clipboard-list"
		d.Logs, _ = db.ListLog(50)
	case "meta_analytics":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Stats", "Meta", "Cloud API Analytics", "la-chart-bar"
	case "meta_webhook":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Meta Webhook", "Meta", "Webhook Config", "la-link"
	case "tracker":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Link Tracker", T("nav_tools"), "Link Clicks", "la-link"
	case "abtests":
		d.Title, d.Pretitle, d.Heading, d.Icon = "A/B Tests", T("nav_whatsapp"), "A/B Testing", "la-balance-scale"
	case "subscribe":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Pilih Paket", "Subscription", "Pricing", "la-shopping-cart"
	case "admin_paygateways":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Payment Gateways", "Admin", "Pay Gateways", "la-credit-card"
	case "admin_transactions_pay":
		d.Title, d.Pretitle, d.Heading, d.Icon = "Payment Transactions", "Admin", "Transactions", "la-receipt"
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
		"add": func(a, b int) int { return a + b },
		"mult": func(a, b float64) float64 { return a * b },
		"split": func(s, sep string) []string { return strings.Split(s, sep) },
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
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 { return nil }
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, _ := values[i].(string)
				dict[key] = values[i+1]
			}
			return dict
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
		if isPSEOPath(r.URL.Path) {
			pseo.HandlePSEO(w, r)
			return
		}
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

func isPSEOPath(path string) bool {
	prefixes := []string{
		"/best-", "/alternatives-to-", "/compare/",
		"/whatsapp-marketing-untuk-", "/beli-aplikasi-",
		"/source-code-", "/aplikasi-whatsapp-",
		"/jual-aplikasi-", "/jual-source-code-",
		"/harga-source-code-", "/jasa-whatsapp-",
		"/cara-", "/chatbot-", "/panduan-",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
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
		uid := getUserID(r)
		if uid == 0 || engine.GetSessionUserID(id) == uid || engine.GetSessionUserID(id) == 0 {
			_ = engine.LogoutAccount(id)
		}
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
	uid := getUserID(r)
	if uid > 0 && db.CountSentByUser(uid) >= db.GetUserSendLimit(uid) {
		http.Redirect(w, r, "/send?msg="+template.URLQueryEscaper("Send limit reached. Upgrade your plan."), http.StatusSeeOther)
		return
	}
	if err := engine.SendFrom(uid, strings.TrimPrefix(r.FormValue("account_phone"), "+"), phone, msgtemplate.Render(message, msgtemplate.Vars{Phone: phone})); err != nil {
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
	uid := getUserID(r)
	if err := engine.SendMedia(uid, accountPhone, phone, mediaType, dest, caption); err != nil {
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
			uid := getUserID(r)
			if err := engine.SendFrom(uid, accountPhone, phone, message); err != nil {
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
	w.Header().Set("Access-Control-Allow-Origin", appURL())
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
	sig := db.GetSetting("agent_signature", "")
	if sig != "" { message += "\n\n" + sig }
	accountPhone := strings.TrimPrefix(r.FormValue("account_phone"), "+")
	if phone == "" || message == "" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":false,"error":"phone and message required"}`)
		return
	}
	uid := getUserID(r)
	if err := engine.SendFrom(uid, accountPhone, phone, message); err != nil {
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
	uid := getUserID(r)
	if uid != 0 && acc.UserID != 0 && acc.UserID != uid {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":false,"error":"forbidden"}`)
		return
	}
	mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
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
func saveAutoreplyMedia(r *http.Request) (mediaType, mediaURL string) {
	file, header, err := r.FormFile("media_file")
	if err != nil { return }
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	mediaDir := "public/uploads/"
	os.MkdirAll(mediaDir, 0755)
	fname := fmt.Sprintf("%s%d%s", mediaDir, time.Now().UnixNano(), ext)
	out, err := os.Create(fname)
	if err != nil { return }
	io.Copy(out, file)
	out.Close()
	mediaURL = "/" + fname
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp": mediaType = "image"
	case ".mp4", ".mov", ".avi", ".mkv": mediaType = "video"
	case ".mp3", ".ogg", ".wav", ".aac", ".m4a": mediaType = "audio"
	default: mediaType = "document"
	}
	return
}

func handleAutoReplyAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	r.ParseMultipartForm(10 << 20)
	keyword := r.FormValue("keyword")
	match := r.FormValue("match")
	reply := r.FormValue("reply")
	useAI := r.FormValue("use_ai") == "on" || r.FormValue("use_ai") == "1"
	aiKeyID, _ := strconv.ParseInt(r.FormValue("ai_key_id"), 10, 64)
	if match == "ai" {
		if faq := r.FormValue("faq"); faq != "" { reply = faq }
		if keyword == "" { keyword = match }
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
	if match == "" { match = "contains" }
	aid := joinVals(r, "account_ids"); uid := getUserID(r)
	if !engine.ValidateAccountIDs(uid, aid) {
		http.Redirect(w, r, "/autoreply?msg="+template.URLQueryEscaper("Nomor tidak valid"), http.StatusSeeOther)
		return
	}
	mediaType, mediaURL := saveAutoreplyMedia(r)
	_, _ = db.AddAutoReply(uid, keyword, match, reply, useAI, aiKeyID, aid, func()int64{t,_:=strconv.ParseInt(r.FormValue("training_id"),10,64);return t}(), mediaType, mediaURL)
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}

func handleAutoReplyDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		_ = db.DeleteAutoReply(getUserID(r), id)
	}
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}

func handleAutoReplyToggle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		_ = db.ToggleAutoReply(getUserID(r), id)
	}
	http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
}
func handleAutoReplyEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if r.Method == http.MethodPost {
		if id > 0 {
			r.ParseMultipartForm(10 << 20)
			keyword := r.FormValue("keyword")
			match := r.FormValue("match")
			reply := r.FormValue("reply")
			useAI := r.FormValue("use_ai") == "on" || r.FormValue("use_ai") == "1"
			aiKeyID, _ := strconv.ParseInt(r.FormValue("ai_key_id"), 10, 64)
			if match == "" { match = "contains" }
			mediaType, mediaURL := saveAutoreplyMedia(r)
			_ = db.UpdateAutoReply(getUserID(r), id, keyword, match, reply, useAI, aiKeyID, joinVals(r, "account_ids"), func()int64{t,_:=strconv.ParseInt(r.FormValue("training_id"),10,64);return t}(), mediaType, mediaURL)
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
			_ = db.UpdateContact(getUserID(r), id, r.FormValue("name"), r.FormValue("phone"), joinVals(r, "groups"))
			var tagIDs []int64
			for _, s := range r.Form["tag_ids"] {
				if tid, err := strconv.ParseInt(s, 10, 64); err == nil { tagIDs = append(tagIDs, tid) }
			}
			if len(tagIDs) > 0 { db.SetContactTags(id, tagIDs) }
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
	_ = db.SetSetting("rate_max_daily", r.FormValue("rate_max_daily"))
	_ = db.SetSetting("rate_random_min", r.FormValue("rate_random_min"))
	_ = db.SetSetting("rate_random_max", r.FormValue("rate_random_max"))
	_ = db.SetSetting("auto_close_hours", r.FormValue("auto_close_hours"))
	_ = db.SetSetting("auto_close_message", r.FormValue("auto_close_message"))
	_ = db.SetSetting("agent_signature", r.FormValue("agent_signature"))
	_ = db.SetSetting("biz_hours_reply", r.FormValue("biz_hours_reply"))
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
	redir := "/settings"
	if strings.Contains(r.Referer(), "/autoreply") { redir = "/autoreply" }
	http.Redirect(w, r, redir, http.StatusSeeOther)
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

func limitGuard(resource string, fn func(*http.Request), redirect string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return
		}
		uid := getUserID(r)
		var count, limit int
		switch resource {
		case "template":
			count, limit = db.CountUserTemplates(uid), db.GetUserTemplateLimit(uid)
		case "canned":
			count, limit = db.CountUserCanned(uid), db.GetUserCannedLimit(uid)
		case "drip":
			count, limit = db.CountUserDrips(uid), db.GetUserDripLimit(uid)
		case "scheduled":
			count, limit = db.CountUserScheduled(uid), db.GetUserScheduledLimit(uid)
		case "webhook":
			count, limit = db.CountUserWebhooks(uid), db.GetUserWebhookLimit(uid)
		case "recurring":
			count, limit = db.CountUserRecurring(uid), db.GetUserRecurringLimit(uid)
		case "form":
			count, limit = db.CountUserForms(uid), db.GetUserFormLimit(uid)
		case "macro":
			count, limit = db.CountUserMacros(uid), db.GetUserMacroLimit(uid)
		case "ai_key":
			count, limit = db.CountUserAiKeys(uid), db.GetUserAiKeyLimit(uid)
		case "knowledge":
			count, limit = db.CountUserKnowledge(uid), db.GetUserKnowledgeLimit(uid)
		case "meta":
			count, limit = db.CountMetaByUser(uid), db.GetUserMetaLimit(uid)
		case "contact":
			count, limit = db.CountUserContacts(uid), db.GetUserContactLimit(uid)
		default:
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return
		}
		if count >= limit {
			http.Redirect(w, r, redirect+"?msg=Limit+reached.+Upgrade+your+plan.", http.StatusSeeOther)
			return
		}
		fn(r)
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

func decryptOrPlain(s string) string {
	if s == "" { return s }
	if d, err := secret.Decrypt(s); err == nil && d != "" { return d }
	return s
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
	uid := getUserID(r)
	if !engine.ValidateAccountIDs(uid, accountID) {
		http.Redirect(w, r, "/broadcast?msg="+template.URLQueryEscaper("Nomor pengirim tidak valid"), http.StatusSeeOther)
		return
	}
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
		list, _ := db.ContactsByGroup(uid, strings.TrimSpace(gid))
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

	metaAccountID, _ := strconv.ParseInt(r.FormValue("meta_account_id"), 10, 64)
	metaTemplate := r.FormValue("meta_template")
	tags := joinVals(r, "tags")
	mediaType := ""; mediaURL := ""
	if r.FormValue("media_type") != "" {
		mediaType = r.FormValue("media_type")
		mediaURL = r.FormValue("media_url")
	}
	// handle file upload
	file, header, ferr := r.FormFile("media_file")
	if ferr == nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		mediaDir := "public/uploads/"
		os.MkdirAll(mediaDir, 0755)
		fname := fmt.Sprintf("%s%d%s", mediaDir, time.Now().UnixNano(), ext)
		out, _ := os.Create(fname)
		if out != nil {
			io.Copy(out, file)
			out.Close()
			mediaURL = "/" + fname
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif", ".webp": mediaType = "image"
			case ".mp4", ".mov", ".avi": mediaType = "video"
			case ".pdf", ".doc", ".docx", ".xls", ".xlsx": mediaType = "document"
			default: mediaType = "document"
			}
		}
	}
	interval, _ := strconv.Atoi(r.FormValue("interval"))
	if interval <= 0 { interval = 300 }
	_, _ = db.AddCampaign(uid, name, groups, normalizedNumbers, mediaType, mediaURL, message, len(seen), accountID, sendMode, interval, metaAccountID, metaTemplate, tags)
	http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
}

func handleCampaignPause(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		uid := getUserID(r)
		camps, _ := db.ListCampaigns(uid)
		for _, c := range camps {
			if c.ID == id {
				if c.Status == "paused" {
					db.UpdateCampaignStatus(uid, id, "running")
				} else if c.Status == "running" {
					db.UpdateCampaignStatus(uid, id, "paused")
				}
				break
			}
		}
	}
	http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
}

func handleCampaignRetry(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		uid := getUserID(r)
		camps, _ := db.ListCampaigns(uid)
		for _, c := range camps {
			if c.ID == id {
				// clone campaign as new pending
				db.AddCampaign(uid, c.Name+" (retry)", c.Groups, c.Numbers, c.MediaType, c.MediaURL, c.Message, c.Total, c.AccountIDs, c.SendMode, c.Interval, c.MetaAccountID, c.MetaTemplate, c.Tags)
				break
			}
		}
	}
	http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
}

// ---- Drip Campaigns ----

func handleDrips(w http.ResponseWriter, r *http.Request) {
	render(w, r, "drips")
}
func handleDripAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/drips", http.StatusSeeOther)
		return
	}
	uid := getUserID(r)
	if db.CountUserDrips(uid) >= db.GetUserDripLimit(uid) {
		http.Redirect(w, r, "/drips?msg=Drip+limit+reached.+Upgrade+your+plan.", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	if name == "" { name = "Drip Campaign" }
	db.AddDrip(uid, name)
	http.Redirect(w, r, "/drips", http.StatusSeeOther)
}
func handleDripStepAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/drips", http.StatusSeeOther)
		return
	}
	dripID, _ := strconv.ParseInt(r.FormValue("drip_id"), 10, 64)
	delay, _ := strconv.Atoi(r.FormValue("delay"))
	message := r.FormValue("message")
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))
	if dripID > 0 && message != "" {
		db.AddDripStep(dripID, delay, message, sortOrder)
	}
	http.Redirect(w, r, "/drips", http.StatusSeeOther)
}
func handleDripToggle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 {
		drips, _ := db.ListDrips()
		for _, d := range drips {
			if d.ID == id {
				if d.Status == "active" {
					db.UpdateDripStatus(id, "inactive")
				} else {
					db.UpdateDripStatus(id, "active")
				}
				break
			}
		}
	}
	http.Redirect(w, r, "/drips", http.StatusSeeOther)
}

func handleTags(w http.ResponseWriter, r *http.Request) {
	render(w, r, "tags")
}

func handleContactTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	contactID, _ := strconv.ParseInt(r.FormValue("contact_id"), 10, 64)
	if contactID > 0 {
		var tagIDs []int64
		for _, s := range r.Form["tag_ids"] {
			if id, err := strconv.ParseInt(s, 10, 64); err == nil {
				tagIDs = append(tagIDs, id)
			}
		}
		db.SetContactTags(contactID, tagIDs)
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func handleCanned(w http.ResponseWriter, r *http.Request) {
	render(w, r, "canned")
}

func handleInboxAssign(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		agentID, _ := strconv.ParseInt(r.FormValue("agent_id"), 10, 64)
		if agentID > 0 {
			db.AssignAgent(phone, agentID)
		}
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleInboxClose(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		db.CloseConversation(phone)
		// send CSAT survey
		uid := getUserID(r)
		if s := engine.FirstSession(uid); s != nil {
			msg := "Terima kasih! Bagaimana pengalaman Anda? Balas dengan rating 1-5 ⭐"
			engine.SendFrom(uid, s.Phone, phone, msg)
		}
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleLinkTrack(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/track/")
	db.LogLinkClick(token)
	var url string
	db.QueryRow(`SELECT url FROM link_clicks WHERE token=?`, token).Scan(&url)
	if url == "" { url = "/" }
	http.Redirect(w, r, url, http.StatusFound)
}

func handleAnalytics(w http.ResponseWriter, r *http.Request) {
	render(w, r, "analytics")
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/broadcast", http.StatusSeeOther)
		return
	}
	phones := strings.Split(r.FormValue("numbers"), "\n")
	var valid, invalid int
	for _, p := range phones {
		p = strings.TrimSpace(p)
		if p == "" { continue }
		if store.ValidFormat(strings.TrimPrefix(strings.TrimPrefix(p, "+"), "0")) && !db.IsBlacklisted(p) && !db.IsUnsub(p) {
			valid++
		} else {
			invalid++
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/broadcast?msg=Valid:+%d,+Invalid:+%d", valid, invalid), http.StatusSeeOther)
}

func handleInboxNote(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid := getUserID(r)
		db.AddNote(r.FormValue("phone"), uid, r.FormValue("note"))
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleInboxTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		agentID, _ := strconv.ParseInt(r.FormValue("agent_id"), 10, 64)
		dept := r.FormValue("dept")
		if agentID > 0 { db.AssignAgent(phone, agentID) }
		if dept != "" { db.AssignToDept(phone, dept) }
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleWidgetJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprintf(w, `(function(){var d=document;d.write('<div id="cwa" style="position:fixed;bottom:20px;right:20px;z-index:9999;font-family:sans-serif"><button id="cwab" onclick="var p=document.getElementById(\'cwac\');p.style.display=p.style.display==\'none\'?\'block\':\'none\'" style="width:56px;height:56px;border-radius:50%%;background:#25D366;border:none;color:#fff;font-size:28px;cursor:pointer;box-shadow:0 4px 12px rgba(0,0,0,.15)">💬</button><div id="cwac" style="display:none;width:350px;height:450px;background:#fff;border-radius:12px;box-shadow:0 4px 24px rgba(0,0,0,.15);margin-bottom:12px;overflow:hidden"><div style="background:#075E54;color:#fff;padding:16px;font-weight:700">Chat with us</div><div id="cwamsg" style="height:330px;overflow-y:auto;padding:12px;font-size:14px"></div><form id="cwaf" onsubmit="var m=document.getElementById(\'cwai\').value;if(!m)return false;var x=new XMLHttpRequest();x.open(\'POST\',\'/widget/chat\');x.setRequestHeader(\'Content-Type\',\'application/x-www-form-urlencoded\');x.send(\'message=\'+encodeURIComponent(m)+\'&phone=\'+encodeURIComponent(\'web_\'+Date.now()));document.getElementById(\'cwai\').value=\'\';document.getElementById(\'cwamsg\').innerHTML+=\'<div style=text-align:right;margin:4px 0><span style=background:#DCF8C6;padding:8px 12px;border-radius:8px;display:inline-block;max-width:80%%>\'+m+\'</span></div>\';return false" style="display:flex;border-top:1px solid #eee"><input id="cwai" placeholder="Message..." style="flex:1;border:none;padding:12px;font-size:14px;outline:none"><button style="background:#075E54;color:#fff;border:none;padding:12px 16px;cursor:pointer">Send</button></form></div></div>')})()`)
}

func handleWidgetChat(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")
	phone := r.FormValue("phone")
	if msg != "" && phone != "" {
		db.LogReceived(phone, "Web Visitor", msg, false, "", "", "widget")
		engine.Notify(phone)
		if s := engine.FirstSession(0); s != nil {
			engine.SendFrom(0, s.Phone, phone, "Terima kasih! Tim kami akan segera membalas.")
		}
	}
	w.WriteHeader(200)
}

func handleInboxLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db.SetConversationLabel(r.FormValue("phone"), r.FormValue("label"))
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleInboxFilter(w http.ResponseWriter, r *http.Request) {
	ft := strings.TrimPrefix(r.URL.Path, "/inbox/filter/")
	w.Header().Set("Content-Type", "application/json")
	items := db.InboxFiltered(getUserID(r), ft)
	if items == nil { fmt.Fprint(w, "[]"); return }
	json.NewEncoder(w).Encode(items)
}

func handleCustomerProfile(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" { http.Error(w, "phone required", 400); return }
	p := db.GetCustomerProfile(phone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func handleBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fname := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
		db.BackupDB("public/backups/" + fname)
		http.Redirect(w, r, "/backup?msg=Backup+created:+backup.sql", http.StatusSeeOther)
		return
	}
	render(w, r, "backup")
}

func handleTranslate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	text := r.FormValue("text")
	to := r.FormValue("to")
	if text == "" || to == "" { http.Error(w, "text and to required", 400); return }
	keys, _ := db.ListAiKeys(getUserID(r))
	if len(keys) == 0 { http.Error(w, "no ai key configured", 400); return }
	ak := keys[0]
	dec, _ := secret.Decrypt(ak.APIKey)
	if dec == "" { dec = ak.APIKey }
	prompt := fmt.Sprintf("Translate to %s, return only translation: %s", to, text)
	if reply, err := aiservice.Reply(dec, ak.Provider, ak.Model, ak.BaseURL, "", prompt, nil, nil); err == nil {
		fmt.Fprint(w, reply)
	} else {
		http.Error(w, err.Error(), 500)
	}
}

func handleMacroExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Redirect(w, r, "/inbox", http.StatusSeeOther); return }
	phone := r.FormValue("phone")
	actions := r.FormValue("actions")
	uid := getUserID(r)
	for _, a := range strings.Split(actions, ";") {
		parts := strings.SplitN(a, ":", 2)
		if len(parts) < 2 { continue }
		switch parts[0] {
		case "assign": aid, _ := strconv.ParseInt(parts[1], 10, 64); db.AssignAgent(phone, aid)
		case "tag": db.SetConversationLabel(phone, parts[1])
		case "reply":
			sig := db.GetSetting("agent_signature", "")
			msg := parts[1]; if sig != "" { msg += "\n\n" + sig }
			if s := engine.FirstSession(uid); s != nil { engine.SendFrom(uid, s.Phone, phone, msg) }
		case "close": db.CloseConversation(phone)
		}
	}
	db.LogAudit(uid, "macro", fmt.Sprintf("%s on %s", actions, phone), r.RemoteAddr)
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleMergeExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		kid, _ := strconv.ParseInt(r.FormValue("keep_id"), 10, 64)
		var mids []int64
		for _, s := range r.Form["merge_ids"] {
			if id, err := strconv.ParseInt(s, 10, 64); err == nil { mids = append(mids, id) }
		}
		db.MergeContacts(kid, mids)
	}
	http.Redirect(w, r, "/merge", http.StatusSeeOther)
}

func handleSetPriority(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		pri, _ := strconv.Atoi(r.FormValue("priority"))
		db.SetContactPriority(r.FormValue("phone"), pri)
	}
	http.Redirect(w, r, "/inbox", http.StatusSeeOther)
}

func handleEmailWebhook(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	subject := r.FormValue("subject")
	body := r.FormValue("text")
	if from == "" && subject == "" { w.WriteHeader(200); return }
	msg := fmt.Sprintf("Email dari %s:\n*%s*\n\n%s", from, subject, body)
	db.LogReceived(from, subject, msg, false, "", "", "email")
	engine.Notify(from)
	w.WriteHeader(200)
}

func handleMetaSend(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		accID, _ := strconv.ParseInt(r.FormValue("account_id"), 10, 64)
		acc, err := db.GetMetaAccount(accID)
		if err == nil {
			uid := getUserID(r)
			if uid == 0 || acc.UserID == 0 || acc.UserID == uid {
				mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
				mc.SendText(r.FormValue("phone"), r.FormValue("message"))
			}
			http.Redirect(w, r, "/meta/send?msg=Sent", http.StatusSeeOther)
			return
		}
	}
	render(w, r, "meta_send")
}

func handleInboxCanned(w http.ResponseWriter, r *http.Request) {
	canned, _ := db.ListCanned()
	var list []map[string]string
	for _, c := range canned {
		list = append(list, map[string]string{"name": c.Name, "message": c.Message})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleInboxSuggest(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" { w.WriteHeader(400); return }
	keys, _ := db.ListAiKeys(getUserID(r))
	if len(keys) == 0 { w.Header().Set("Content-Type", "application/json"); fmt.Fprint(w, `{"reply":""}`); return }
	ak := keys[0]
	dec, _ := secret.Decrypt(ak.APIKey)
	if dec == "" { dec = ak.APIKey }
	msgs, _ := db.ChatHistory(phone, 10)
	context := ""
	for i := len(msgs) - 1; i >= 0; i-- {
		prefix := "Customer"
		if msgs[i].Type == "sent" { prefix = "Agent" }
		context += prefix + ": " + msgs[i].Message + "\n"
	}
	prompt := "Based on this conversation, suggest a helpful short reply in Indonesian:\n" + context + "\nReply:"
	reply, err := aiservice.Reply(dec, ak.Provider, ak.Model, ak.BaseURL, "", prompt, nil, nil)
	w.Header().Set("Content-Type", "application/json")
	if err != nil || reply == "" { fmt.Fprint(w, `{"reply":""}`); return }
	reply = strings.Trim(reply, "\"' \n")
	b, _ := json.Marshal(map[string]string{"reply": reply})
	w.Write(b)
}

func handleInboxStar(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	typ := r.FormValue("type")
	if id > 0 && (typ == "sent" || typ == "received") { db.ToggleStar(typ, id) }
	w.WriteHeader(200)
}

// ---- Payment / Subscription ----

func handleFAQ(w http.ResponseWriter, r *http.Request) {
	render(w, r, "faq")
}

func handleFAQAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db.AddFAQ(getUserID(r), r.FormValue("question"), r.FormValue("answer"))
	}
	http.Redirect(w, r, "/faq", http.StatusSeeOther)
}

func handleFAQDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if id > 0 { db.DeleteFAQ(getUserID(r), id) }
	http.Redirect(w, r, "/faq", http.StatusSeeOther)
}

func handleFAQImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/faq", http.StatusSeeOther)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil { http.Redirect(w, r, "/faq?msg=File+required", http.StatusSeeOther); return }
	defer file.Close()
	reader := csv.NewReader(file)
	headers, _ := reader.Read()
	colQ, colA := -1, -1
	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch h {
		case "question", "pertanyaan": colQ = i
		case "answer", "jawaban": colA = i
		}
	}
	if colQ < 0 || colA < 0 { http.Redirect(w, r, "/faq?msg=CSV+must+have+question+and+answer+columns", http.StatusSeeOther); return }
	uid := getUserID(r)
	count := 0
	for {
		record, err := reader.Read()
		if err != nil { break }
		q := strings.TrimSpace(safeGet(record, colQ))
		a := strings.TrimSpace(safeGet(record, colA))
		if q != "" && a != "" { db.AddFAQ(uid, q, a); count++ }
	}
	http.Redirect(w, r, "/faq?msg=Imported+"+strconv.Itoa(count)+"+FAQs", http.StatusSeeOther)
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	render(w, r, "subscribe")
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/subscribe", http.StatusSeeOther)
		return
	}
	uid := getUserID(r)
	packageID, _ := strconv.ParseInt(r.FormValue("package_id"), 10, 64)
	gatewayID, _ := strconv.ParseInt(r.FormValue("gateway_id"), 10, 64)
	voucher := r.FormValue("voucher")

	if uid == 0 || packageID == 0 {
		http.Redirect(w, r, "/subscribe?msg=Invalid", http.StatusSeeOther)
		return
	}

	// Redeem free voucher (no payment needed)
	if voucher != "" {
		if _, err := db.RedeemVoucher(uid, voucher); err == nil {
			http.Redirect(w, r, "/?msg=Subscription+activated+via+voucher!", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/subscribe?msg=Invalid+voucher+code", http.StatusSeeOther)
		return
	}

	if gatewayID == 0 {
		http.Redirect(w, r, "/subscribe?msg=Select+a+payment+method", http.StatusSeeOther)
		return
	}
	pkg, err := db.GetPackage(packageID)
	if err != nil {
		http.Redirect(w, r, "/subscribe?msg=Package+not+found", http.StatusSeeOther)
		return
	}
	gw, err := db.GetPaymentGateway(gatewayID)
	if err != nil || gw.Status != "active" {
		http.Redirect(w, r, "/subscribe?msg=Gateway+not+available", http.StatusSeeOther)
		return
	}
	price, _ := strconv.ParseFloat(pkg.Price, 64)
	if price <= 0 { price = 99000 }
	currency := gw.Currency
	invoiceID := store.GenInvoiceID()

	cfg := payment.GatewayConfig{
		APIKey: decryptOrPlain(gw.APIKey), APISecret: decryptOrPlain(gw.APISecret),
		WebhookSecret: decryptOrPlain(gw.WebhookSecret), BaseURL: gw.BaseURL, Currency: currency,
	}
	pg, err := payment.New(gw.Provider, cfg)
	if err != nil {
		http.Error(w, "Gateway error: "+err.Error(), 500)
		return
	}
	user, _ := db.GetUserByID(uid)
	email := ""
	if user != nil { email = user.Email }
	callbackURL := appURL() + "/payment/callback/" + gw.Provider
	result, err := pg.CreateCharge(payment.ChargeParams{
		InvoiceID: invoiceID, Amount: price, Currency: currency,
		Description: pkg.Name + " Subscription",
		CustomerName: r.FormValue("name"), CustomerEmail: email,
		Items:       []payment.ChargeItem{{ID: pkg.Name, Name: pkg.Name, Price: price, Quantity: 1}},
		CallbackURL: callbackURL, ReturnURL: appURL() + "/subscribe",
	})
	if err != nil {
		http.Error(w, "Payment error: "+err.Error(), 500)
		return
	}
	db.CreateTransaction(uid, packageID, gatewayID, price, currency, invoiceID, result.RedirectURL, result.ExternalID)
	http.Redirect(w, r, result.RedirectURL, http.StatusSeeOther)
}

func handlePaymentCallback(w http.ResponseWriter, r *http.Request) {
	provider := strings.TrimPrefix(r.URL.Path, "/payment/callback/")
	body, _ := io.ReadAll(r.Body)
	gateways, _ := db.ListPaymentGateways()
	var gw *store.PaymentGateway
	for _, g := range gateways {
		if g.Provider == provider && g.Status == "active" { gw = &g; break }
	}
	if gw == nil { http.Error(w, "Unknown gateway", 400); return }
	cfg := payment.GatewayConfig{
		APIKey: decryptOrPlain(gw.APIKey), APISecret: decryptOrPlain(gw.APISecret),
		WebhookSecret: decryptOrPlain(gw.WebhookSecret), BaseURL: gw.BaseURL, Currency: gw.Currency,
	}
	pg, err := payment.New(provider, cfg)
	if err != nil { http.Error(w, err.Error(), 500); return }
	headers := map[string]string{}
	for k := range r.Header {
		headers[k] = r.Header.Get(k)
	}
	result, err := pg.VerifyCallback(body, headers)
	if err != nil {
		db.Log("payment", "verify_error", fmt.Sprintf("%s: %s", provider, err.Error()))
		http.Error(w, "Verification failed", 400)
		return
	}
	if result.Status == "paid" {
		trx, _ := db.GetTransactionByInvoice(result.InvoiceID)
		if trx != nil && trx.Status != "paid" {
			db.UpdateTransactionStatus(result.InvoiceID, "paid")
			db.ActivateSubscription(trx.UserID, trx.PackageID)
			db.Log("payment", "paid", fmt.Sprintf("%s paid via %s: %.0f %s", result.InvoiceID, provider, result.Amount, gw.Currency))
		}
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
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
	uid := getUserID(r)
	existingGroups, _ := db.ListGroups(uid)
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
					id, err := db.AddGroup(uid, gn)
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
		existing, _ := db.FindContactByPhone(uid, phone)
		if existing != nil {
			skipped++
			continue
		}
		if name == "" { name = phone }
		if _, err := db.AddContact(uid, name, phone, gidStr); err == nil {
			imported++
		}
	}
	msg := fmt.Sprintf("Imported+%d+contacts", imported)
	if skipped > 0 { msg += fmt.Sprintf(",+%d+skipped+(duplicate)", skipped) }
	http.Redirect(w, r, "/contacts?msg="+msg, http.StatusSeeOther)
}

func handleContactExport(w http.ResponseWriter, r *http.Request) {
	contacts, _ := db.ListContacts(getUserID(r))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=contacts.csv")
	w.Write([]byte("\xEF\xBB\xBF")) // BOM for Excel
	cw := csv.NewWriter(w)
	cw.Write([]string{"name", "phone", "groups"})
	for _, c := range contacts {
		cw.Write([]string{c.Name, c.Phone, c.Groups})
	}
	cw.Flush()
}

func handleContactBulkDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	count := 0
	uid := getUserID(r)
	for _, idStr := range r.Form["ids"] {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil { continue }
		db.DeleteContact(uid, id)
		count++
	}
	http.Redirect(w, r, fmt.Sprintf("/contacts?msg=Deleted+%d+contacts", count), http.StatusSeeOther)
}

// ---- Scheduled ----

func handleScheduled(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		render(w, r, "scheduled")
		return
	}
	uid := getUserID(r)
	if db.CountUserScheduled(uid) >= db.GetUserScheduledLimit(uid) {
		http.Redirect(w, r, "/scheduled?msg=Scheduled+limit+reached.+Upgrade+your+plan.", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	message := r.FormValue("message")
	sendAt := r.FormValue("send_at")
	repeat, _ := strconv.Atoi(r.FormValue("repeat"))
	if name == "" || phone == "" || message == "" || sendAt == "" {
		http.Redirect(w, r, "/scheduled", http.StatusSeeOther)
		return
	}
	sendAt = strings.Replace(sendAt, "T", " ", 1) + ":00"
	accountIDs := joinVals(r, "account_ids")
	if !engine.ValidateAccountIDs(uid, accountIDs) {
		http.Redirect(w, r, "/scheduled?msg="+template.URLQueryEscaper("Nomor pengirim tidak valid"), http.StatusSeeOther)
		return
	}
	_, _ = db.AddScheduled(uid, name, phone, message, sendAt, repeat, accountIDs)
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
				if db.GetAssignedAgent(m.From) == 0 {
					db.AssignNextRoundRobin(m.From)
				}
				db.MarkRead(m.From)
				engine.Notify(m.From)
				db.Log("meta", "received", fmt.Sprintf("%s -> %s: %s", m.From, acc.Name, text))

				// Drip: auto-enroll + STOP
				trimmed := strings.ToLower(strings.TrimSpace(text))
				if trimmed == "stop" || trimmed == "berhenti" || trimmed == "unsub" {
					db.UnenrollFromDrip(m.From)
				} else {
					drips, _ := db.ListDrips()
					for _, d := range drips {
						if d.Status == "active" {
							db.EnrollInDrip(d.ID, m.From, "")
						}
					}
				}

				// Auto-reply for Meta
				mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)

					// PRO: Flow builder check
				if wa.MetaFlowCallback != nil {
					if replies, matched := wa.MetaFlowCallback(acc.UserID, acc.PhoneNumberID, m.From, text, ""); matched {
						for _, reply := range replies {
							switch reply.Action {
							case "poll":
								if opts, ok := reply.ActionData["options"].([]string); ok {
									mc.SendPoll(m.From, reply.Text, opts)
								}
							case "template":
								tplName := ""
								if v, ok := reply.ActionData["template_name"].(string); ok { tplName = v }
								mc.SendTemplate(m.From, tplName, "id", nil)
							case "flow":
								fid := ""; ft := ""
								if v, ok := reply.ActionData["flow_id"].(string); ok { fid = v }
								if v, ok := reply.ActionData["flow_token"].(string); ok { ft = v }
								mc.SendFlow(m.From, fid, ft)
							case "location":
								lat := 0.0; lng := 0.0
								if v, ok := reply.ActionData["lat"].(float64); ok { lat = v }
								if v, ok := reply.ActionData["lng"].(float64); ok { lng = v }
								mc.SendLocation(m.From, reply.Text, "", lat, lng)
							case "contacts":
								mc.SendText(m.From, reply.Text)
							case "reaction":
								msgID := ""
								if v, ok := reply.ActionData["message_id"].(string); ok { msgID = v }
								mc.SendReaction(m.From, msgID, "❤️")
							case "reply":
								msgID := ""
								if v, ok := reply.ActionData["message_id"].(string); ok { msgID = v }
								mc.SendReply(m.From, msgID, reply.Text)
							case "carousel":
								mc.SendText(m.From, reply.Text)
							default:
								if reply.MediaURL != "" {
									if mediaID, err := mc.UploadMedia(reply.MediaURL); err == nil {
										mc.SendMediaByID(m.From, reply.MediaType, mediaID, reply.Text)
									} else {
										mc.SendMedia(m.From, reply.MediaType, reply.MediaURL, reply.Text)
									}
								} else if reply.Text != "" {
									mc.SendText(m.From, reply.Text)
								}
							}
						}
						continue
					}
				}

				// Spam detection
				if db.TrackSpam(m.From, fmt.Sprintf("%x", text[:min(len(text), 20)])) {
					db.AddBlacklist(m.From, "auto: spam detected (meta)")
					continue
				}
				// FAQ check
				if reply, found := db.FindFAQAnswer(acc.UserID, text); found {
					mc.SendText(m.From, msgtemplate.Render(reply, msgtemplate.Vars{Phone: m.From, Message: text}))
					db.LogSent(m.From, reply, "faq", "meta"); continue
				}
				// Store Bot
				if trimmed == "menu" || trimmed == "katalog" {
					if cats, _ := db.ListCategories(); len(cats) > 0 {
						var msg string
						for _, c := range cats { msg += "• " + c.Name + "\n" }
						mc.SendText(m.From, "📋 Kategori:\n"+msg+"\nKetik nama kategori."); continue
					}
				}
				// Anti-spam mute
				if aiservice.CheckSpam(m.From, text) || aiservice.CheckJailbreak(text) { continue }
				// Department auto-detect
				if depts, _ := db.ListDepts(); len(depts) > 0 {
					for _, d := range depts {
						if strings.Contains(strings.ToLower(text), strings.ToLower(d.Name)) {
							db.AssignToDept(m.From, d.Name); break
						}
					}
				}
				// Human handoff
				if db.GetSetting("handoff_enabled", "0") == "1" {
					keywords := strings.Split(strings.ToLower(db.GetSetting("handoff_keywords", "admin,telp,manusia,cs,operator")), ",")
					for _, kw := range keywords {
						if strings.TrimSpace(kw) != "" && strings.Contains(strings.ToLower(text), strings.TrimSpace(kw)) {
							handoffMsg := msgtemplate.Render(db.GetSetting("handoff_message", "Silakan hubungi admin kami."), msgtemplate.Vars{Phone: m.From, Message: text})
							mc.SendText(m.From, handoffMsg); continue
						}
					}
				}
				// Welcome message
				if db.GetSetting("welcome_enabled", "0") == "1" && db.MarkWelcomed(m.From) {
					if wmsg := db.GetSetting("welcome_message", ""); wmsg != "" {
						mc.SendText(m.From, msgtemplate.Render(wmsg, msgtemplate.Vars{Phone: m.From, Message: text}))
						db.LogSent(m.From, wmsg, "welcome", "meta")
					}
				}
				// Fallback message
				fallbackSent := false
				// Webhook dispatch
				db.Log("meta", "received", fmt.Sprintf("%s -> %s: %s", m.From, acc.Name, text))

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
						fallbackSent = true
					}
				}
				// Fallback message
				if !fallbackSent && db.GetSetting("fallback_enabled", "0") == "1" {
					if fmsg := db.GetSetting("fallback_message", ""); fmsg != "" {
						mc.SendText(m.From, msgtemplate.Render(fmsg, msgtemplate.Vars{Phone: m.From, Message: text}))
						db.LogSent(m.From, fmsg, "fallback", "meta")
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

func handleMetaFlows(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	var result []map[string]interface{}
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		flows, _ := mc.ListFlows()
		if flows != nil { result = append(result, flows...) }
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleMetaCatalog(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	if r.Method == http.MethodPost {
		for _, acc := range accounts {
			mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
			mc.SyncProduct(r.FormValue("name"), r.FormValue("description"), parseFloat(r.FormValue("price")), r.FormValue("image_url"), r.FormValue("website"))
		}
		http.Redirect(w, r, "/admin/meta/catalog", http.StatusSeeOther)
		return
	}
	var result []map[string]interface{}
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		products, _ := mc.ListProducts()
		if products != nil { result = append(result, products...) }
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleMetaCalling(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	voices, _ := mc.FetchElevenLabsVoices()
	status, _ := mc.GetHealthStatus()
	json.NewEncoder(w).Encode(map[string]interface{}{"status":status,"voices":voices})
}

func handleMetaProfile(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/admin/meta/health", http.StatusSeeOther) }

func handleMetaQR(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	msg := r.URL.Query().Get("message"); if msg == "" { msg = "Hello!" }
	qr, _ := mc.GenerateQRCode(msg)
	json.NewEncoder(w).Encode(qr)
}

func handleMetaHealth(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	var result []map[string]interface{}
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		health, _ := mc.GetHealthStatus()
		quality, _ := mc.GetQualityScore()
		phone, _ := mc.GetPhoneNumberStatus()
		result = append(result, map[string]interface{}{"phone_id":acc.PhoneNumberID,"name":acc.Name,"health":health,"quality":quality,"phone":phone})
	}
	json.NewEncoder(w).Encode(result)
}

func parseFloat(s string) float64 { f, _ := strconv.ParseFloat(s, 64); return f }

func handleMetaTemplates(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	if r.Method == http.MethodPost {
		name := r.FormValue("name"); language := r.FormValue("language"); category := r.FormValue("category")
		if name != "" { mc.CreateTemplate(name, language, category, nil) }
	}
	templates, _ := mc.FetchTemplates()
	json.NewEncoder(w).Encode(templates)
}

func handleMetaCarousel(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	if r.Method == http.MethodPost {
		var products []map[string]string
		json.NewDecoder(r.Body).Decode(&products)
		mc.SendCarousel(r.FormValue("to"), products)
	}
	json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}

func handleMetaWebhooks(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	if r.Method == http.MethodPost {
		fields := r.Form["fields"]
		if len(fields) > 0 { mc.SubscribeWebhook(fields) }
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"webhook_url": appURL() + "/webhook/meta",
		"fields":      []string{"messages", "message_template_status_update", "message_template_quality_update"},
	})
}

func handleMetaPayment(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	if r.Method == http.MethodPost {
		result, _ := mc.SendPaymentRequest(r.FormValue("to"), r.FormValue("token"), r.FormValue("amount"), r.FormValue("currency"))
		json.NewEncoder(w).Encode(map[string]string{"result": result})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status":"ready"})
}

func handleMetaRegister(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	if r.Method == http.MethodPost {
		if r.FormValue("action") == "register" { mc.RegisterPhoneNumber(r.FormValue("pin")) }
		if r.FormValue("action") == "deregister" { mc.DeregisterPhoneNumber() }
	}
	status, _ := mc.GetPhoneNumberStatus()
	json.NewEncoder(w).Encode(status)
}

func handleSheets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ready"})
}

func handleN8NWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(r.Body)
	var input map[string]interface{}
	json.Unmarshal(body, &input)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "input": input})
}

func handleN8NNodeDefinition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nodes := []map[string]interface{}{
		{"name": "chatgo.sendMessage", "displayName": "ChatGo — Send Message", "inputs": []string{"main"}, "outputs": []string{"main"}, "properties": []interface{}{map[string]interface{}{"displayName": "Phone", "name": "phone", "type": "string", "required": true}, map[string]interface{}{"displayName": "Message", "name": "message", "type": "string", "required": true}}},
		{"name": "chatgo.addContact", "displayName": "ChatGo — Add Contact", "inputs": []string{"main"}, "outputs": []string{"main"}, "properties": []interface{}{map[string]interface{}{"displayName": "Name", "name": "name", "type": "string", "required": true}, map[string]interface{}{"displayName": "Phone", "name": "phone", "type": "string", "required": true}}},
		{"name": "chatgo.trigger", "displayName": "ChatGo — WhatsApp Trigger", "inputs": []string{}, "outputs": []string{"main"}},
	}
	json.NewEncoder(w).Encode(nodes)
}

func handleN8NTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := []map[string]string{
		{"id": "wa-to-sheets", "name": "WhatsApp → Google Sheets"},
		{"id": "form-to-wa", "name": "Form → WhatsApp"},
		{"id": "broadcast-scheduler", "name": "Broadcast Scheduler"},
	}
	json.NewEncoder(w).Encode(t)
}

func handleMetaInsights(w http.ResponseWriter, r *http.Request) {
	accounts, _ := db.ListMetaAccounts()
	w.Header().Set("Content-Type", "application/json")
	if len(accounts) == 0 { json.NewEncoder(w).Encode(map[string]string{"status":"no_meta_account"}); return }
	mc := meta.New(accounts[0].PhoneNumberID, decryptOrPlain(accounts[0].AccessToken), accounts[0].VerifyToken)
	msgID := r.URL.Query().Get("message_id")
	var insights map[string]interface{}
	if msgID != "" { insights, _ = mc.GetMessageInsights(msgID) }
	rateLimit, _ := mc.GetRateLimit()
	quality, _ := mc.GetQualityScore()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"insights":   insights,
		"rate_limit": rateLimit,
		"quality":    quality,
	})
}

func handleIGWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mode := r.URL.Query().Get("hub.mode")
		challenge := r.URL.Query().Get("hub.challenge")
		if mode == "subscribe" { fmt.Fprint(w, challenge) }
		return
	}
	body, _ := io.ReadAll(r.Body)
	msgs, ok := meta.ParseIGWebhook(body)
	if !ok { return }
	accounts, _ := db.ListMetaAccounts()
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		for _, m := range msgs {
			db.LogReceivedForWA(acc.PhoneNumberID, m.From, "", m.Text.Body, false, "", "", "instagram")
			if wa.IGCallback != nil {
				if replies, matched := wa.IGCallback(acc.UserID, m.From, m.Text.Body, ""); matched {
					for _, reply := range replies {
						if reply.Text != "" { mc.SendIGMessage(m.From, reply.Text) }
					}
				}
			}
		}
	}
	fmt.Fprint(w, "ok")
}

func handleIGInbox(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	accounts, _ := db.ListMetaAccounts()
	var result []map[string]interface{}
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		convs, _ := mc.GetIGConversations()
		if convs != nil { result = append(result, convs...) }
	}
	json.NewEncoder(w).Encode(result)
}

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(r.Body)
	token := db.GetSetting("telegram_token", "")
	if token == "" { json.NewEncoder(w).Encode(map[string]string{"status":"no_token"}); return }
	bot := telegram.New(token)
	update, err := telegram.ParseUpdate(body)
	if err != nil || update == nil { json.NewEncoder(w).Encode(map[string]string{"status":"invalid"}); return }

	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	text := update.Message.Text
	name := update.Message.Chat.FirstName

	db.LogReceivedForWA("tg", chatID, name, text, false, "", "", "telegram")

	// Flow builder integration
	if wa.TGCallback != nil {
		if replies, matched := wa.TGCallback(0, chatID, text, name); matched {
			for _, reply := range replies {
				if reply.Text != "" { bot.SendMessage(update.Message.Chat.ID, reply.Text) }
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}












func handleOmniInbox(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	channel := r.URL.Query().Get("channel")
	var conversations []map[string]interface{}

	// WhatsApp conversations
	if channel == "" || channel == "wa" {
		recv, _ := db.ListReceivedPaginated(1, 50)
		for _, msg := range recv {
			conversations = append(conversations, map[string]interface{}{
				"id": msg.ID, "phone": msg.Phone, "message": msg.Message,
				"channel": "wa", "time": msg.Created, "name": msg.Name,
			})
		}
	}

	// Instagram conversations (from Meta accounts)
	if channel == "" || channel == "ig" {
		accounts, _ := db.ListMetaAccounts()
		for _, acc := range accounts {
			mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
			convs, _ := mc.GetIGConversations()
			for _, c := range convs {
				conversations = append(conversations, map[string]interface{}{
					"channel": "ig", "data": c,
				})
			}
		}
	}

	json.NewEncoder(w).Encode(conversations)
}

func handleOmniAnalytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	waSent := db.CountSent()
	waRecv := db.CountReceived()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"channels": map[string]interface{}{
			"whatsapp": map[string]int{"sent": waSent, "received": waRecv},
			"instagram": map[string]int{"sent": 0, "received": 0},
			"telegram": map[string]int{"sent": 0, "received": 0},
		},
		"total": map[string]int{
			"sent": waSent, "received": waRecv,
		},
	})
}

func handleOmniHandoff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fromChannel := r.FormValue("from")
	toChannel := r.FormValue("to")
	phone := r.FormValue("phone")
	message := r.FormValue("message")

	if toChannel == "wa" && phone != "" {
		engine.SendFrom(0, "", phone,
			fmt.Sprintf("🔀 Handoff dari %s:\n%s\n\nSilakan dibantu.", fromChannel, message))
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "handoff_initiated",
		"from":   fromChannel,
		"to":     toChannel,
	})
}

func handleAgency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, _ := db.ListUsers()
	var clients []map[string]interface{}
	for _, u := range users {
		if u.Role != "admin" {
			clients = append(clients, map[string]interface{}{
				"id": u.ID, "name": u.Name, "email": u.Email,
				"role": u.Role, "created": u.Created,
			})
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clients": clients,
		"total":   len(clients),
	})
}

func handleAISettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := getUserID(r)
	if r.Method == http.MethodPost {
		db.SetSetting("ai_provider_"+strconv.FormatInt(uid, 10), r.FormValue("provider"))
		db.SetSetting("ai_model_"+strconv.FormatInt(uid, 10), r.FormValue("model"))
		db.SetSetting("ai_temp_"+strconv.FormatInt(uid, 10), r.FormValue("temperature"))
	}
	json.NewEncoder(w).Encode(map[string]string{
		"provider": db.GetSetting("ai_provider_"+strconv.FormatInt(uid, 10), "openai"),
		"model":    db.GetSetting("ai_model_"+strconv.FormatInt(uid, 10), "gpt-4o"),
		"temp":     db.GetSetting("ai_temp_"+strconv.FormatInt(uid, 10), "0.7"),
	})
}

func handleButtonsBuilder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone"); title := r.FormValue("title"); footer := r.FormValue("footer")
		buttons := strings.Split(r.FormValue("buttons"), ",")
		engine.SendButtons(0, "", phone, title, footer, buttons)
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func handleWarmer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	accounts := engine.Accounts(0)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"accounts": len(accounts), "status": "ready",
	})
}

func handleFBWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		mode := r.URL.Query().Get("hub.mode"); challenge := r.URL.Query().Get("hub.challenge")
		if mode == "subscribe" { fmt.Fprint(w, challenge) }; return
	}
	body, _ := io.ReadAll(r.Body)
	msgs, ok := meta.ParseIGWebhook(body)
	if !ok { json.NewEncoder(w).Encode(map[string]string{"status":"no_messages"}); return }
	accounts, _ := db.ListMetaAccounts()
	for _, acc := range accounts {
		mc := meta.New(acc.PhoneNumberID, decryptOrPlain(acc.AccessToken), acc.VerifyToken)
		for _, m := range msgs {
			db.LogReceivedForWA(acc.PhoneNumberID, m.From, "", m.Text.Body, false, "", "", "facebook")
			if wa.FBCallback != nil {
				if replies, matched := wa.FBCallback(acc.UserID, m.From, m.Text.Body, ""); matched {
					for _, reply := range replies {
						if reply.Text != "" { mc.SendFBMessage(m.From, reply.Text) }
					}
				}
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}

func handleFlowSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	flows, _ := db.ListFlows(0)
	var result []interface{}
	for _, f := range flows {
		if q != "" && !strings.Contains(strings.ToLower(f.Name), strings.ToLower(q)) { continue }
		result = append(result, map[string]interface{}{"id": f.ID, "name": f.Name, "active": f.Active})
	}
	json.NewEncoder(w).Encode(result)
}

func handleDarkMode(w http.ResponseWriter, r *http.Request) {
	uid := getUserID(r)
	mode := db.GetSetting("dark_mode_"+strconv.FormatInt(uid, 10), "light")
	if r.Method == http.MethodPost {
		mode = r.FormValue("mode")
		db.SetSetting("dark_mode_"+strconv.FormatInt(uid, 10), mode)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mode": mode})
}

func handleFlowLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("flow_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	logs, _ := db.GetFlowExecutionLog(id, 50)
	json.NewEncoder(w).Encode(logs)
}