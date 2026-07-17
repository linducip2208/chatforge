package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"regexp"
	"chatgo/secret"
	"net/http"
	"strconv"
	"strings"
)

// registerAdminRoutes wires every remaining Zender menu (SMS/Hosts/Android/AI/Admin/Docs).
// Auto-reply core & WhatsApp engine untouched.
func registerAdminRoutes(mux *http.ServeMux) {
	// helper: wrap pageHandler with auth
	ap := func(page string) http.HandlerFunc {
		return authMiddleware(pageHandler(page))
	}
	acp := func(fn func(*http.Request), redirect string) http.HandlerFunc {
		return authMiddleware(crudPost(fn, redirect))
	}
	acd := func(fn func(int64), redirect string) http.HandlerFunc {
		return authMiddleware(crudDel(fn, redirect))
	}
	a := func(fn http.HandlerFunc) http.HandlerFunc {
		return authMiddleware(fn)
	}

	// Hosts
	mux.HandleFunc("/hosts/whatsapp", ap("hosts_whatsapp"))

	// AI
	mux.HandleFunc("/ai/keys", ap("ai_keys"))
	mux.HandleFunc("/ai/keys/add", acp(func(r *http.Request) {
		enc, _ := secret.Encrypt(r.FormValue("apikey"))
		uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
		db.AddAiKey(uid, r.FormValue("name"), r.FormValue("provider"), r.FormValue("model"), enc, r.FormValue("base_url"), r.FormValue("system_prompt"))
	}, "/autoreply"))
	mux.HandleFunc("/ai/keys/delete", a(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
			db.DeleteAiKey(uid, id)
		}
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
	}))
	// Knowledge base
	mux.HandleFunc("/knowledge", ap("knowledge"))
	mux.HandleFunc("/knowledge/add", acp(func(r *http.Request) {
		r.ParseForm()
		title := r.FormValue("title")
		question := r.FormValue("question")
		answer := r.FormValue("answer")
		category := r.FormValue("category")
		rows, _ := json.Marshal([]map[string]string{{"question": question, "answer": answer, "category": category}})
		db.AddKnowledge(title, string(rows))
	}, "/autoreply"))
	mux.HandleFunc("/knowledge/delete", acd(func(id int64) { db.DeleteKnowledge(id) }, "/autoreply"))
	mux.HandleFunc("/knowledge/toggle", a(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 { db.ToggleKnowledge(id) }
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
	}))
	// CSV import
	mux.HandleFunc("/knowledge/import", a(handleKnowledgeImport))
	// URL training
	mux.HandleFunc("/knowledge/url", a(handleKnowledgeURL))
	// PDF upload
	mux.HandleFunc("/knowledge/pdf", a(handleKnowledgePDF))
	// AI Training Campaigns
	mux.HandleFunc("/ai/training/add", acp(func(r *http.Request) {
		aiKeyID, _ := strconv.ParseInt(r.FormValue("ai_key_id"), 10, 64)
		db.AddAiTraining(r.FormValue("name"), r.FormValue("system_prompt"), aiKeyID)
	}, "/autoreply"))
	mux.HandleFunc("/ai/training/delete", acd(func(id int64) { db.DeleteAiTraining(id) }, "/autoreply"))
	mux.HandleFunc("/ai/plugins", ap("ai_plugins"))
	mux.HandleFunc("/ai/plugins/add", acp(func(r *http.Request) { db.AddAiPlugin(r.FormValue("name"), r.FormValue("endpoint")) }, "/autoreply"))
	mux.HandleFunc("/ai/plugins/delete", acd(func(id int64) { db.DeleteAiPlugin(id) }, "/autoreply"))

	// Admin
	mux.HandleFunc("/admin", ap("admin"))
	mux.HandleFunc("/admin/users", ap("admin_users"))
	mux.HandleFunc("/admin/users/add", acp(func(r *http.Request) {
		id, err := db.AddUser(r.FormValue("name"), r.FormValue("email"), r.FormValue("role"), r.FormValue("country"))
		if err == nil && r.FormValue("password") != "" {
			hash, _ := hashPassword(r.FormValue("password"))
			db.SetUserPassword(id, hash)
		}
	}, "/admin/users"))
	mux.HandleFunc("/admin/users/delete", acd(func(id int64) { db.DeleteUser(id) }, "/admin/users"))
	mux.HandleFunc("/admin/users/edit", a(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			db.UpdateUser(id, r.FormValue("name"), r.FormValue("email"), r.FormValue("role"))
			if pw := r.FormValue("password"); pw != "" {
				hash, _ := hashPassword(pw)
				db.SetUserPassword(id, hash)
			}
		}
		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
	}))
	mux.HandleFunc("/admin/roles", ap("admin_roles"))
	mux.HandleFunc("/admin/roles/add", acp(func(r *http.Request) { db.AddRole(r.FormValue("name"), joinVals(r, "permissions")) }, "/admin/roles"))
	mux.HandleFunc("/admin/roles/delete", acd(func(id int64) { db.DeleteRole(id) }, "/admin/roles"))
	mux.HandleFunc("/admin/roles/edit", a(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			db.UpdateRole(id, r.FormValue("name"), joinVals(r, "permissions"))
		}
		http.Redirect(w, r, "/admin/roles", http.StatusSeeOther)
	}))
	mux.HandleFunc("/admin/packages", ap("admin_packages"))
	mux.HandleFunc("/admin/packages/add", acp(func(r *http.Request) {
		s, _ := strconv.Atoi(r.FormValue("send_limit"))
		rc, _ := strconv.Atoi(r.FormValue("receive_limit"))
		dv, _ := strconv.Atoi(r.FormValue("device_limit"))
		us, _ := strconv.Atoi(r.FormValue("ussd_limit"))
		ws, _ := strconv.Atoi(r.FormValue("wa_send_limit"))
		wr, _ := strconv.Atoi(r.FormValue("wa_receive_limit"))
		wa, _ := strconv.Atoi(r.FormValue("wa_account_limit"))
		co, _ := strconv.Atoi(r.FormValue("contact_limit"))
		sc, _ := strconv.Atoi(r.FormValue("scheduled_limit"))
		kl, _ := strconv.Atoi(r.FormValue("key_limit"))
		wl, _ := strconv.Atoi(r.FormValue("webhook_limit"))
		al, _ := strconv.Atoi(r.FormValue("action_limit"))
		sv := joinVals(r, "services")
		hd, _ := strconv.Atoi(r.FormValue("hidden"))
		fm, _ := strconv.Atoi(r.FormValue("footermark"))
		ml, _ := strconv.Atoi(r.FormValue("meta_limit"))
		dl, _ := strconv.Atoi(r.FormValue("drip_limit"))
		rl, _ := strconv.Atoi(r.FormValue("recurring_limit"))
		fl, _ := strconv.Atoi(r.FormValue("form_limit"))
		tl, _ := strconv.Atoi(r.FormValue("template_limit"))
		cl, _ := strconv.Atoi(r.FormValue("canned_limit"))
		mcl, _ := strconv.Atoi(r.FormValue("macro_limit"))
		akl, _ := strconv.Atoi(r.FormValue("ai_key_limit"))
		knl, _ := strconv.Atoi(r.FormValue("knowledge_limit"))
		db.AddPackage(r.FormValue("name"), r.FormValue("price"), s, rc, dv, us, ws, wr, wa, co, sc, kl, wl, al, ml, dl, rl, fl, tl, cl, mcl, akl, knl, sv, hd, fm)
	}, "/admin/packages"))
	mux.HandleFunc("/admin/packages/edit", a(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
			return
		}
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			s, _ := strconv.Atoi(r.FormValue("send_limit"))
			rc, _ := strconv.Atoi(r.FormValue("receive_limit"))
			dv, _ := strconv.Atoi(r.FormValue("device_limit"))
			us, _ := strconv.Atoi(r.FormValue("ussd_limit"))
			ws, _ := strconv.Atoi(r.FormValue("wa_send_limit"))
			wr, _ := strconv.Atoi(r.FormValue("wa_receive_limit"))
			wa, _ := strconv.Atoi(r.FormValue("wa_account_limit"))
			co, _ := strconv.Atoi(r.FormValue("contact_limit"))
			sc, _ := strconv.Atoi(r.FormValue("scheduled_limit"))
			kl, _ := strconv.Atoi(r.FormValue("key_limit"))
			wl, _ := strconv.Atoi(r.FormValue("webhook_limit"))
			al, _ := strconv.Atoi(r.FormValue("action_limit"))
			ml, _ := strconv.Atoi(r.FormValue("meta_limit"))
			dl, _ := strconv.Atoi(r.FormValue("drip_limit"))
			rl, _ := strconv.Atoi(r.FormValue("recurring_limit"))
			fl, _ := strconv.Atoi(r.FormValue("form_limit"))
			tl, _ := strconv.Atoi(r.FormValue("template_limit"))
			cl, _ := strconv.Atoi(r.FormValue("canned_limit"))
			mcl, _ := strconv.Atoi(r.FormValue("macro_limit"))
			akl, _ := strconv.Atoi(r.FormValue("ai_key_limit"))
			knl, _ := strconv.Atoi(r.FormValue("knowledge_limit"))
			sv := joinVals(r, "services")
			hd, _ := strconv.Atoi(r.FormValue("hidden"))
			fm, _ := strconv.Atoi(r.FormValue("footermark"))
			db.UpdatePackage(id, r.FormValue("name"), r.FormValue("price"), s, rc, dv, us, ws, wr, wa, co, sc, kl, wl, al, ml, dl, rl, fl, tl, cl, mcl, akl, knl, sv, hd, fm)
		}
		http.Redirect(w, r, "/admin/packages", http.StatusSeeOther)
	}))
	mux.HandleFunc("/admin/packages/delete", acd(func(id int64) { db.DeletePackage(id) }, "/admin/packages"))
	mux.HandleFunc("/admin/vouchers", ap("admin_vouchers"))
	mux.HandleFunc("/admin/vouchers/add", acp(func(r *http.Request) {
		dur, _ := strconv.Atoi(r.FormValue("duration"))
		count, _ := strconv.Atoi(r.FormValue("count"))
		if count <= 0 { count = 1 }
		for i := 0; i < count; i++ {
			db.AddVoucher(randSecret()[:10], r.FormValue("pkg"), dur)
		}
	}, "/admin/vouchers"))
	mux.HandleFunc("/admin/vouchers/delete", acd(func(id int64) { db.DeleteVoucher(id) }, "/admin/vouchers"))
	mux.HandleFunc("/admin/subscriptions", ap("admin_subscriptions"))
	mux.HandleFunc("/admin/subscriptions/add", acp(func(r *http.Request) { db.AddSubscription(r.FormValue("user"), r.FormValue("pkg"), r.FormValue("expire")) }, "/admin/subscriptions"))
	mux.HandleFunc("/admin/subscriptions/delete", acd(func(id int64) { db.DeleteSubscription(id) }, "/admin/subscriptions"))
	mux.HandleFunc("/admin/transactions", ap("admin_transactions"))
	mux.HandleFunc("/admin/payouts", ap("admin_payouts"))
	mux.HandleFunc("/admin/payouts/delete", acd(func(id int64) { db.DeletePayout(id) }, "/admin/payouts"))
	mux.HandleFunc("/admin/pages", ap("admin_pages"))
	mux.HandleFunc("/admin/pages/add", acp(func(r *http.Request) { db.AddPage(r.FormValue("title"), r.FormValue("slug"), r.FormValue("content")) }, "/admin/pages"))
	mux.HandleFunc("/admin/pages/delete", acd(func(id int64) { db.DeletePage(id) }, "/admin/pages"))
	mux.HandleFunc("/admin/marketing", ap("admin_marketing"))
	mux.HandleFunc("/admin/marketing/add", acp(func(r *http.Request) { db.AddMarketing(r.FormValue("title"), r.FormValue("content")) }, "/admin/marketing"))
	mux.HandleFunc("/admin/marketing/delete", acd(func(id int64) { db.DeleteMarketing(id) }, "/admin/marketing"))
	mux.HandleFunc("/admin/languages", ap("admin_languages"))
	mux.HandleFunc("/admin/languages/add", acp(func(r *http.Request) { db.AddLanguageAdmin(r.FormValue("name"), r.FormValue("iso")) }, "/admin/languages"))
	mux.HandleFunc("/admin/languages/delete", acd(func(id int64) { db.DeleteLanguageAdmin(id) }, "/admin/languages"))
	mux.HandleFunc("/admin/waservers", ap("admin_waservers"))
	mux.HandleFunc("/admin/waservers/add", acp(func(r *http.Request) {
		acc, _ := strconv.Atoi(r.FormValue("accounts"))
		pkgs := joinVals(r, "packages")
		db.AddWaServer(r.FormValue("name"), r.FormValue("url"), r.FormValue("port"), r.FormValue("secret"), acc, pkgs)
	}, "/admin/waservers"))
	mux.HandleFunc("/admin/waservers/delete", acd(func(id int64) { db.DeleteWaServer(id) }, "/admin/waservers"))
	mux.HandleFunc("/admin/waservers/edit", acp(func(r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		acc, _ := strconv.Atoi(r.FormValue("accounts"))
		pkgs := joinVals(r, "packages")
		db.UpdateWaServer(id, r.FormValue("name"), r.FormValue("url"), r.FormValue("port"), r.FormValue("secret"), acc, pkgs)
	}, "/admin/waservers"))
	mux.HandleFunc("/admin/gateways", ap("admin_gateways"))
	mux.HandleFunc("/admin/gateways/add", acp(func(r *http.Request) { db.AddGateway(r.FormValue("name")) }, "/admin/gateways"))
	mux.HandleFunc("/admin/gateways/delete", acd(func(id int64) { db.DeleteGateway(id) }, "/admin/gateways"))
	mux.HandleFunc("/admin/shorteners", ap("admin_shorteners"))
	mux.HandleFunc("/admin/shorteners/add", acp(func(r *http.Request) { db.AddShortener(r.FormValue("name")) }, "/admin/shorteners"))
	mux.HandleFunc("/admin/shorteners/delete", acd(func(id int64) { db.DeleteShortener(id) }, "/admin/shorteners"))
	mux.HandleFunc("/admin/plugins", ap("admin_plugins"))
	mux.HandleFunc("/admin/plugins/add", acp(func(r *http.Request) { db.AddPlugin(r.FormValue("name"), r.FormValue("dir")) }, "/admin/plugins"))
	mux.HandleFunc("/admin/plugins/delete", acd(func(id int64) { db.DeletePlugin(id) }, "/admin/plugins"))

	mux.HandleFunc("/admin/meta", ap("admin_meta"))
	mux.HandleFunc("/admin/meta/add", a(func(w http.ResponseWriter, r *http.Request) {
		uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
		if db.CountMetaByUser(uid) >= db.GetUserMetaLimit(uid) {
			http.Redirect(w, r, "/admin/meta?msg=Meta+limit+reached", http.StatusSeeOther)
			return
		}
		encToken, _ := secret.Encrypt(r.FormValue("access_token"))
		encSecret, _ := secret.Encrypt(r.FormValue("app_secret"))
		db.AddMetaAccount(r.FormValue("name"), r.FormValue("phone_number_id"), encToken, r.FormValue("app_id"), encSecret, r.FormValue("verify_token"), uid, 0)
		http.Redirect(w, r, "/admin/meta", http.StatusSeeOther)
	}))
	mux.HandleFunc("/admin/meta/delete", a(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 {
			uid, _ := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
			acc, err := db.GetMetaAccount(id)
			if err == nil && (acc.UserID == 0 || acc.UserID == uid || uid == 0) {
				db.DeleteMetaAccount(id)
			}
		}
		http.Redirect(w, r, "/admin/meta", http.StatusSeeOther)
	}))

	mux.HandleFunc("/admin/metatemplates", ap("admin_metatemplates"))
	mux.HandleFunc("/admin/metatemplates/add", acp(func(r *http.Request) {
		db.AddMetaTemplate(r.FormValue("name"), r.FormValue("language"), r.FormValue("category"), r.FormValue("components"), "active")
	}, "/admin/metatemplates"))
	mux.HandleFunc("/admin/metatemplates/delete", acd(func(id int64) { db.DeleteMetaTemplate(id) }, "/admin/metatemplates"))

	// Payment Gateways
	mux.HandleFunc("/admin/gateways-pay", ap("admin_paygateways"))
	mux.HandleFunc("/admin/gateways-pay/add", a(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			encKey, _ := secret.Encrypt(r.FormValue("api_key"))
			encSecret, _ := secret.Encrypt(r.FormValue("api_secret"))
			encWebhook, _ := secret.Encrypt(r.FormValue("webhook_secret"))
			db.AddPaymentGateway(r.FormValue("name"), r.FormValue("provider"), encKey, encSecret, encWebhook, r.FormValue("base_url"), r.FormValue("currency"), r.FormValue("config"))
			http.Redirect(w, r, "/admin/gateways-pay", http.StatusSeeOther)
		}
	}))
	mux.HandleFunc("/admin/gateways-pay/delete", acd(func(id int64) { db.DeletePaymentGateway(id) }, "/admin/gateways-pay"))
	mux.HandleFunc("/admin/gateways-pay/toggle", a(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		if id > 0 { db.TogglePaymentGateway(id) }
		http.Redirect(w, r, "/admin/gateways-pay", http.StatusSeeOther)
	}))
	mux.HandleFunc("/admin/transactions-pay", ap("admin_transactions_pay"))

	// Docs
	mux.HandleFunc("/docs", ap("docs"))
}

func handleKnowledgeImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		title = "CSV Import"
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Redirect(w, r, "/autoreply?msg=File+required", http.StatusSeeOther)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		http.Redirect(w, r, "/autoreply?msg=Invalid+CSV", http.StatusSeeOther)
		return
	}
	colQ := -1; colA := -1; colC := -1
	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch h {
		case "question", "pertanyaan": colQ = i
		case "answer", "jawaban": colA = i
		case "category", "kategori": colC = i
		}
	}
	if colQ < 0 || colA < 0 {
		http.Redirect(w, r, "/autoreply?msg=CSV+must+have+question+and+answer+columns", http.StatusSeeOther)
		return
	}
	var rows []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF { break }
		if err != nil { continue }
		q := safeGet(record, colQ)
		a := safeGet(record, colA)
		if q == "" || a == "" { continue }
		c := ""
		if colC >= 0 { c = safeGet(record, colC) }
		rows = append(rows, map[string]string{"question": q, "answer": a, "category": c})
	}
	if len(rows) == 0 {
		http.Redirect(w, r, "/autoreply?msg=No+valid+rows", http.StatusSeeOther)
		return
	}
	content, _ := json.Marshal(map[string]interface{}{"rows": rows})
	db.AddKnowledge(title, string(content))
	http.Redirect(w, r, "/autoreply?msg=Imported+"+strconv.Itoa(len(rows))+"+rows", http.StatusSeeOther)
}

func safeGet(record []string, i int) string {
	if i >= 0 && i < len(record) { return strings.TrimSpace(record[i]) }
	return ""
}

func handleKnowledgeURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	urlStr := strings.TrimSpace(r.FormValue("url"))
	if urlStr == "" {
		http.Redirect(w, r, "/autoreply?msg=URL+required", http.StatusSeeOther)
		return
	}
	resp, err := http.DefaultClient.Get(urlStr)
	if err != nil {
		http.Redirect(w, r, "/autoreply?msg=URL+fetch+failed", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 100*1024))
	text := string(body)
	for _, re := range []string{`<script[^>]*>[\s\S]*?</script>`, `<style[^>]*>[\s\S]*?</style>`, `<[^>]+>`, `&[a-z]+;`} {
		text = regexp.MustCompile(re).ReplaceAllString(text, " ")
	}
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	if len(text) > 3000 {
		text = text[:3000]
	}
	if len(text) < 20 {
		http.Redirect(w, r, "/autoreply?msg=URL+content+too+short", http.StatusSeeOther)
		return
	}
	if title == "" {
		title = "URL: " + urlStr
		if len(title) > 100 { title = title[:100] }
	}
	content, _ := json.Marshal(map[string]interface{}{"rows": []map[string]string{{"content": text}}})
	db.AddKnowledge(title, string(content))
	http.Redirect(w, r, "/autoreply?msg=URL+trained:+"+title, http.StatusSeeOther)
}

func handleKnowledgePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/autoreply", http.StatusSeeOther)
		return
	}
	title := r.FormValue("title")
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Redirect(w, r, "/autoreply?msg=File+required", http.StatusSeeOther)
		return
	}
	defer file.Close()
	raw, _ := io.ReadAll(io.LimitReader(file, 5*1024*1024))
	if len(raw) < 100 {
		http.Redirect(w, r, "/autoreply?msg=PDF+too+small", http.StatusSeeOther)
		return
	}
	text := extractPDFText(string(raw))
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	if len(text) > 5000 { text = text[:5000] }
	if len(text) < 20 {
		http.Redirect(w, r, "/autoreply?msg=PDF+text+not+found+(scanned?)", http.StatusSeeOther)
		return
	}
	if title == "" { title = "PDF Import" }
	content, _ := json.Marshal(map[string]interface{}{"rows": []map[string]string{{"content": text}}})
	db.AddKnowledge(title, string(content))
	http.Redirect(w, r, "/autoreply?msg=PDF+ok+("+strconv.Itoa(len(text))+"+chars)", http.StatusSeeOther)
}

func extractPDFText(raw string) string {
	var out strings.Builder
	re := regexp.MustCompile(`(?s)BT(.*?)ET`)
	textRe := regexp.MustCompile(`\(([^)]*)\)`)
	for _, match := range re.FindAllStringSubmatch(raw, -1) {
		if len(match) > 1 {
			for _, t := range textRe.FindAllStringSubmatch(match[1], -1) {
				if len(t) > 1 { out.WriteString(t[1] + " ") }
			}
		}
	}
	return out.String()
}
