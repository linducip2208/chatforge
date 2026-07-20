package pseo

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PageData struct {
	Title       string
	Description string
	Canonical   string
	OGImage     string
	SiteName    string
	AppURL      string
	WaNumber    string
	Year        int
	Content     template.HTML
	Breadcrumbs []Breadcrumb
	JSONLD      template.JS
	IsHome      bool
	IsList      bool
	IsCompare   bool
	Items       []ListItem
	CompareA    string
	CompareB    string
	CurrentYear int
}

type Breadcrumb struct {
	Name string
	URL  string
}

type ListItem struct {
	Title       string
	URL         string
	Description string
}

var (
	AppName string
	AppURL  string
)

func Init(appName, appURL, waNumber string) {
	AppName = appName
	AppURL = strings.TrimRight(appURL, "/")
}

func safeSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func humanize(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Split(s, " ")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func renderPSEO(w http.ResponseWriter, r *http.Request, d PageData) {
	d.Year = time.Now().Year()
	d.CurrentYear = d.Year
	d.SiteName = AppName
	d.WaNumber = "6281296052010"

	if d.Canonical == "" {
		d.Canonical = AppURL + r.URL.Path
	}
	if d.OGImage == "" {
		d.OGImage = AppURL + "/assets/theme/default-favicon.png"
	}

	d.JSONLD = template.JS(generateJSONLD(d))

	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}

	tmpl := template.Must(template.New("pseo").Funcs(funcMap).Parse(pseoTemplate))
	if err := tmpl.Execute(w, d); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func generateJSONLD(d PageData) string {
	if d.IsHome {
		return fmt.Sprintf(`{"@context":"https://schema.org","@type":"WebSite","name":"%s","url":"%s","description":"%s"}`,
			template.JSEscapeString(d.SiteName), template.JSEscapeString(d.Canonical), template.JSEscapeString(d.Description))
	}
	if d.IsCompare {
		return fmt.Sprintf(`{"@context":"https://schema.org","@type":"WebPage","name":"%s","url":"%s","description":"%s","mainEntity":{"@type":"FAQPage","mainEntity":[{"@type":"Question","name":"Apa perbedaan %s dan %s?","acceptedAnswer":{"@type":"Answer","text":"%s adalah solusi self-hosted one-time purchase, sedangkan %s biasanya SaaS berlangganan. %s menawarkan semua fitur: broadcast, auto-reply AI, multi-account, chatbot — tanpa biaya bulanan."}}]}}`,
			template.JSEscapeString(d.Title), template.JSEscapeString(d.Canonical), template.JSEscapeString(d.Description),
			template.JSEscapeString(d.CompareA), template.JSEscapeString(d.CompareB),
			template.JSEscapeString(d.CompareA), template.JSEscapeString(d.CompareB), template.JSEscapeString(AppName))
	}
	if d.IsList {
		items := ""
		for i, item := range d.Items {
			if i > 0 {
				items += ","
			}
			items += fmt.Sprintf(`{"@type":"ListItem","position":%d,"url":"%s","name":"%s"}`,
				i+1, template.JSEscapeString(item.URL), template.JSEscapeString(item.Title))
		}
		return fmt.Sprintf(`{"@context":"https://schema.org","@type":"ItemList","name":"%s","url":"%s","description":"%s","itemListElement":[%s]}`,
			template.JSEscapeString(d.Title), template.JSEscapeString(d.Canonical), template.JSEscapeString(d.Description), items)
	}
	return fmt.Sprintf(`{"@context":"https://schema.org","@type":"WebPage","name":"%s","url":"%s","description":"%s"}`,
		template.JSEscapeString(d.Title), template.JSEscapeString(d.Canonical), template.JSEscapeString(d.Description))
}

func ctaBanner() template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="background:linear-gradient(135deg,#4F46E5,#7C3AED);border-radius:16px;padding:32px;color:#fff;margin:32px 0;text-align:center">
<h2 style="font-size:1.8rem;font-weight:800;margin:0 0 8px">Dapatkan Source Code Aplikasi</h2>
<p style="font-size:1.1rem;opacity:.9;margin:0 0 20px">Self-hosted, bayar sekali, lifetime update. Bebas biaya bulanan selamanya.</p>
<div style="display:flex;gap:12px;justify-content:center;flex-wrap:wrap">
<a href="https://wa.me/%s" style="background:#25D366;color:#fff;padding:14px 32px;border-radius:12px;text-decoration:none;font-weight:700;font-size:1.05rem;display:inline-flex;align-items:center;gap:8px" target="_blank" rel="noopener">💬 WhatsApp %s</a>
<a href="/docs" style="background:rgba(255,255,255,.2);color:#fff;padding:14px 32px;border-radius:12px;text-decoration:none;font-weight:700;font-size:1.05rem">📖 Lihat Dokumentasi</a>
</div></div>`, AppName, "6281296052010", "6281296052010"))
}

func listCTABanner(title string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="background:#f8f9fc;border:2px solid #4F46E5;border-radius:16px;padding:24px 28px;margin:24px 0;display:flex;align-items:center;gap:20px;flex-wrap:wrap">
<div style="flex:1;min-width:200px">
<h3 style="font-size:1.2rem;font-weight:700;color:#152e4d;margin:0 0 4px">Cari %s?</h3>
<p style="color:#5e6b7e;margin:0;font-size:.9rem">ChatGo adalah solusi self-hosted lengkap. Broadcast, auto-reply AI, multi-account, chatbot — tanpa biaya bulanan.</p>
</div>
<a href="https://wa.me/6281296052010" target="_blank" rel="noopener" style="background:#25D366;color:#fff;padding:12px 24px;border-radius:10px;text-decoration:none;font-weight:700;white-space:nowrap">💬 Tanya via WA</a>
</div>`, title))
}

// HandlePSEO is the main catch-all handler
func HandlePSEO(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check for compound patterns first (contain "untuk-" or "-di-")
	if strings.Contains(path, "-untuk-") || strings.Contains(path, "-di-") {
		p := parseCompoundSlug(path)
		if p.IsCompound {
			renderCompoundPage(w, r, p)
			return
		}
	}

	slug := strings.TrimPrefix(path, "/")

	switch {
	case strings.HasPrefix(slug, "best-"):
		handleBest(w, r)
	case strings.HasPrefix(slug, "alternatives-to-"):
		handleAlternative(w, r)
	case strings.HasPrefix(slug, "compare/"):
		handleCompare(w, r)
	case strings.HasPrefix(slug, "whatsapp-marketing-untuk-"):
		handleIndustry(w, r)
	case strings.HasPrefix(slug, "beli-aplikasi-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "source-code-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "aplikasi-whatsapp-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "jual-aplikasi-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "jual-source-code-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "harga-source-code-"):
		handleSourceCode(w, r)
	case strings.HasPrefix(slug, "jasa-whatsapp-"):
		handleCitySourceCode(w, r)
	case strings.HasPrefix(slug, "cara-"):
		handleCara(w, r)
	case strings.HasPrefix(slug, "chatbot-"):
		handleChatbot(w, r)
	case strings.HasPrefix(slug, "panduan-"):
		handlePanduan(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleBest(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/best-")
	// extract year if present
	year := time.Now().Year()
	baseSlug := slug
	for _, y := range []int{2024, 2025, 2026, 2027, 2028} {
		suffix := fmt.Sprintf("-%d", y)
		if strings.HasSuffix(slug, suffix) {
			baseSlug = strings.TrimSuffix(slug, suffix)
			year = y
			break
		}
	}

	var bp *struct{ Slug, Title, Keyword string }
	for i := range BestPages {
		if BestPages[i].Slug == baseSlug {
			bp = &BestPages[i]
			break
		}
	}
	if bp == nil {
		// Dynamic best page
		bp = &struct{ Slug, Title, Keyword string }{slug, humanize(slug), humanize(slug)}
	}

	title := fmt.Sprintf("%d Rekomendasi %s %d", len(Competitors)+1, bp.Title, year)
	desc := fmt.Sprintf("Daftar %s rekomendasi %d. Bandingkan fitur, harga, dan kelebihan masing-masing. %s adalah solusi self-hosted tanpa biaya bulanan.", bp.Title, year, len(Competitors)+1, AppName)

	items := make([]ListItem, 0, len(Competitors)+1)
	items = append(items, ListItem{Title: AppName + " (Self-Hosted)", URL: AppURL, Description: "Solusi self-hosted lengkap: broadcast, auto-reply AI, multi-account, chatbot. Bayar sekali, lifetime."})
	for _, c := range Competitors {
		items = append(items, ListItem{Title: humanize(c), URL: "/alternatives-to-" + c, Description: fmt.Sprintf("Lihat alternatif %s — bandingkan fitur dengan %s.", humanize(c), AppName)})
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      true,
		Items:       items,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {bp.Title, ""}},
		Content:     listCTABanner(bp.Title) + bestContent(title, bp.Title, year, bp.Keyword),
	})
}

func handleAlternative(w http.ResponseWriter, r *http.Request) {
	tool := strings.TrimPrefix(r.URL.Path, "/alternatives-to-")
	toolName := humanize(tool)

	title := fmt.Sprintf("5 Alternatif %s — Pilih Solusi WhatsApp Marketing Terbaik", toolName)
	desc := fmt.Sprintf("Bandingkan %s dengan %s dan alternatif lainnya. Mana yang lebih murah, lengkap fitur, dan cocok untuk bisnis Anda?", toolName, AppName)

	items := make([]ListItem, 0, 6)
	items = append(items, ListItem{Title: AppName + " (Rekomendasi Utama)", URL: AppURL, Description: "Self-hosted, bayar sekali, semua fitur unlimited. Broadcast, auto-reply AI, chatbot, multi-account."})
	for i, c := range Competitors {
		if strings.EqualFold(c, tool) || i >= 5 {
			continue
		}
		items = append(items, ListItem{Title: humanize(c), URL: "/compare/" + safeSlug(AppName) + "-vs-" + c, Description: fmt.Sprintf("Bandingkan %s vs %s secara detail.", AppName, humanize(c))})
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      true,
		Items:       items,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Alternatif Tools", ""}},
		Content:     listCTABanner(fmt.Sprintf("Alternatif %s", toolName)) + alternativeContent(toolName),
	})
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/compare/")
	parts := strings.Split(slug, "-vs-")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	a, b := humanize(parts[0]), humanize(parts[1])

	title := fmt.Sprintf("%s vs %s: Perbandingan Fitur, Harga, Kelebihan 2026", a, b)
	desc := fmt.Sprintf("Bandingkan %s vs %s secara detail: fitur broadcast, auto-reply AI, chatbot, multi-account, harga, dan kelebihan. Mana yang lebih cocok?", a, b)

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsCompare:   true,
		CompareA:    a,
		CompareB:    b,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Compare", ""}},
		Content:     ctaBanner() + compareContent(a, b),
	})
}

func handleIndustry(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/whatsapp-marketing-untuk-")
	var ind *struct{ Name, Slug, Icon string }
	for i := range Industries {
		if Industries[i].Slug == slug {
			ind = &Industries[i]
			break
		}
	}
	if ind == nil {
		ind = &struct{ Name, Slug, Icon string }{humanize(slug), slug, "📱"}
	}

	title := fmt.Sprintf("WhatsApp Marketing untuk %s %s: Solusi Lengkap", ind.Icon, ind.Name)
	desc := fmt.Sprintf("Optimalkan WhatsApp marketing untuk bisnis %s Anda. Auto-reply, broadcast promo, chatbot CS 24/7, manajemen kontak — dalam satu platform.", ind.Name)

	features := []ListItem{
		{"Broadcast Promo Massal", "/best-whatsapp-broadcast-tools", "Kirim promo ke ribuan customer sekaligus dengan satu klik."},
		{"Auto Reply Cerdas", "/best-whatsapp-auto-reply", "Balas otomatis pertanyaan customer 24/7 dengan AI."},
		{"Chatbot AI", "/best-whatsapp-chatbot", "Chatbot pintar yang bisa jawab FAQ, proses order, booking."},
		{"Manajemen Kontak", "/best-whatsapp-marketing-tools", "Import kontak, grouping, tagging, segmentasi."},
		{"Laporan & Analitik", "/best-whatsapp-marketing-tools", "Dashboard pengiriman, delivery rate, response rate."},
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      true,
		Items:       features,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"WhatsApp Marketing untuk " + ind.Name, ""}},
		Content:     listCTABanner("WhatsApp Marketing untuk " + ind.Name) + industryContent(ind.Name),
	})
}

func handleSourceCode(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	title := humanize(path)
	desc := fmt.Sprintf("%s — WhatsApp Gateway, Multi-Account, Broadcast, Auto-Reply AI, Chatbot. Self-hosted, bayar sekali, lifetime update. Source code lengkap + dokumentasi.", title)

	// Find matching source code page
	for _, sc := range SourceCodePages {
		if sc.Slug == path {
			title = sc.Title
			break
		}
	}

	features := []ListItem{
		{"WhatsApp Multi-Account", "/best-whatsapp-gateway", "Kelola puluhan nomor WA dalam satu dashboard."},
		{"Broadcast & Blast", "/best-whatsapp-broadcast-tools", "Kirim pesan massal aman anti-banned."},
		{"Auto Reply + AI", "/best-whatsapp-auto-reply", "Auto-reply cerdas dengan AI OpenAI/DeepSeek."},
		{"REST API", "/best-whatsapp-api", "Integrasi dengan sistem lain via API."},
		{"Chatbot 24/7", "/best-whatsapp-chatbot", "Layanan customer service otomatis."},
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      true,
		Items:       features,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Source Code", ""}},
		Content:     ctaBanner() + sourceCodeContent(title),
	})
}

func handleCitySourceCode(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/jasa-whatsapp-")
	slug = strings.TrimSuffix(slug, "/")
	parts := strings.Split(slug, "-")

	var cityName string
	for _, c := range Cities {
		if strings.Contains(slug, c.Slug) {
			cityName = c.Name
			break
		}
	}
	if cityName == "" && len(parts) > 0 {
		cityName = humanize(slug)
	}

	title := fmt.Sprintf("Jasa WhatsApp Marketing %s — Self-Hosted Source Code", cityName)
	desc := fmt.Sprintf("Cari jasa WhatsApp marketing di %s? Lebih hemat dengan self-hosted source code. Broadcast, auto-reply, chatbot — tanpa biaya bulanan. Support & training termasuk.", cityName)

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Jasa WhatsApp " + cityName, ""}},
		Content:     ctaBanner() + cityContent(cityName),
	})
}

func handleCara(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/cara-")

	var cp *struct{ Slug, Title, Desc string }
	for i := range CaraPages {
		if CaraPages[i].Slug == slug {
			cp = &CaraPages[i]
			break
		}
	}
	if cp == nil {
		cp = &struct{ Slug, Title, Desc string }{slug, humanize(slug), "Panduan langkah demi langkah."}
	}

	title := cp.Title
	desc := cp.Desc

	steps := []ListItem{
		{"1. Setup Akun WhatsApp", "", "Hubungkan nomor WA via scan QR di dashboard. Multi-account support."},
		{"2. Tambah Kontak", "", "Import kontak via CSV atau tambah manual. Grouping & tagging."},
		{"3. Konfigurasi Auto Reply", "", "Setup keyword + AI auto-reply. Pilih nomor WA pengirim."},
		{"4. Mulai Broadcast", "", "Buat campaign, pilih grup, kirim massal dengan interval aman."},
		{"5. Pantau Laporan", "", "Dashboard real-time: sent, delivered, read, replied."},
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      true,
		Items:       steps,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Panduan", ""}},
		Content:     listCTABanner(title) + caraContent(cp.Title, cp.Desc),
	})
}

func handleChatbot(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/chatbot-")

	var cp *struct{ Slug, Title, Desc string }
	for i := range ChatbotPages {
		if ChatbotPages[i].Slug == slug {
			cp = &ChatbotPages[i]
			break
		}
	}
	if cp == nil {
		cp = &struct{ Slug, Title, Desc string }{slug, humanize(slug), "Chatbot WhatsApp pintar dengan AI."}
	}

	title := cp.Title
	desc := cp.Desc

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Chatbot WhatsApp", ""}},
		Content:     ctaBanner() + chatbotContent(title, desc),
	})
}

func handlePanduan(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/panduan-")

	var pp *struct{ Slug, Title, Desc string }
	for i := range PanduanPages {
		if PanduanPages[i].Slug == slug {
			pp = &PanduanPages[i]
			break
		}
	}
	if pp == nil {
		pp = &struct{ Slug, Title, Desc string }{slug, humanize(slug), "Panduan lengkap WhatsApp marketing."}
	}

	title := pp.Title
	desc := pp.Desc

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		Breadcrumbs: []Breadcrumb{{"Home", "/"}, {"Panduan", ""}},
		Content:     listCTABanner(title) + panduanContent(title, desc),
	})
}

func HandleSitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>`+"\n")
	fmt.Fprint(w, `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`+"\n")

	total := CountPSEOURLs()
	perPage := 45000
	pages := (total + perPage - 1) / perPage
	for i := 1; i <= pages; i++ {
		fmt.Fprintf(w, "  <sitemap><loc>%s/sitemaps/%d</loc></sitemap>\n", AppURL, i)
	}

	fmt.Fprint(w, "</sitemapindex>")
}

func HandleSitemapPage(w http.ResponseWriter, r *http.Request) {
	numStr := strings.TrimPrefix(r.URL.Path, "/sitemaps/")
	num, err := strconv.Atoi(numStr)
	if err != nil || num < 1 {
		http.NotFound(w, r)
		return
	}

	perPage := 45000
	startIdx := (num - 1) * perPage

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	fmt.Fprint(w, `<?xml version="1.0" encoding="UTF-8"?>`+"\n")
	fmt.Fprint(w, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`+"\n")

	count := 0
	add := func(loc, prio string) {
		if count >= startIdx && count < startIdx+perPage {
			fmt.Fprintf(w, "  <url><loc>%s</loc><priority>%s</priority></url>\n", AppURL+loc, prio)
		}
		count++
	}

	// 1. Core pages
	add("/", "1.0"); add("/docs", "0.9"); add("/login", "0.3"); add("/register", "0.3")
	for _, bp := range BestPages {
		add("/best-"+bp.Slug, "0.8")
		add(fmt.Sprintf("/best-%s-%d", bp.Slug, time.Now().Year()), "0.8")
	}
	for _, c := range Competitors {
		add("/alternatives-to-"+c, "0.7")
	}
	for _, feat := range Features {
		add("/best-whatsapp-"+feat.Slug, "0.6")
	}
	for _, sc := range SourceCodePages {
		add("/"+sc.Slug, "0.8")
	}
	for _, cp := range CaraPages {
		add("/cara-"+cp.Slug, "0.7")
	}
	for _, cp := range ChatbotPages {
		add("/chatbot-"+cp.Slug, "0.7")
	}
	for _, pp := range PanduanPages {
		add("/panduan-"+pp.Slug, "0.7")
	}

	// 2. Industry × City × 2
	for _, ind := range Industries {
		for _, city := range Cities {
			add(fmt.Sprintf("/whatsapp-marketing-untuk-%s-di-%s", ind.Slug, city.Slug), "0.7")
			add(fmt.Sprintf("/jasa-whatsapp-marketing-untuk-%s-di-%s", ind.Slug, city.Slug), "0.6")
		}
	}

	// 3. Compare: all competitor pairs + industry variants
	for _, comp := range Competitors {
		add("/compare/"+safeSlug(AppName)+"-vs-"+comp, "0.6")
		for _, ind := range Industries {
			add(fmt.Sprintf("/compare/%s-vs-%s-untuk-%s", safeSlug(AppName), comp, ind.Slug), "0.6")
		}
	}

	// 4. Compare × Industry × City (3-way)
	for _, comp := range Competitors {
		for _, ind := range Industries {
			for _, city := range Cities {
				add(fmt.Sprintf("/compare/%s-vs-%s-untuk-%s-di-%s", safeSlug(AppName), comp, ind.Slug, city.Slug), "0.5")
			}
		}
	}

	// 5. Alternatives × Industry × City (3-way)
	for _, c := range Competitors {
		for _, ind := range Industries {
			for _, city := range Cities {
				add(fmt.Sprintf("/alternatives-to-%s-untuk-%s-di-%s", c, ind.Slug, city.Slug), "0.5")
			}
		}
	}

	// 6. Source code × Industry × City
	scPrefixes := []string{"beli-aplikasi-whatsapp-marketing", "beli-aplikasi-whatsapp-broadcast", "beli-aplikasi-whatsapp-chatbot", "beli-aplikasi-whatsapp-auto-reply", "source-code-whatsapp-marketing", "source-code-whatsapp-blast", "source-code-chatbot-whatsapp", "source-code-whatsapp-gateway", "aplikasi-whatsapp-marketing-self-hosted", "aplikasi-whatsapp-gateway-multi-account", "jual-aplikasi-whatsapp-marketing", "jual-source-code-whatsapp"}
	for _, sc := range scPrefixes {
		for _, ind := range Industries {
			for _, city := range Cities {
				add(fmt.Sprintf("/%s-untuk-%s-di-%s", sc, ind.Slug, city.Slug), "0.6")
			}
		}
	}

	// 7. Cara × Industry × City × Feature (4-way for max volume)
	for _, cp := range CaraPages {
		for _, feat := range Features {
			for _, ind := range Industries {
				for _, city := range Cities {
					add(fmt.Sprintf("/cara-%s-%s-untuk-%s-di-%s", cp.Slug, feat.Slug, ind.Slug, city.Slug), "0.4")
				}
			}
		}
	}

	// 8. Best × Industry × City × Year
	for _, bp := range BestPages {
		for _, ind := range Industries {
			for _, city := range Cities {
				for y := 2024; y <= 2026; y++ {
					add(fmt.Sprintf("/best-%s-untuk-%s-di-%s-%d", bp.Slug, ind.Slug, city.Slug, y), "0.5")
				}
			}
		}
	}

	// 9. City pages
	for _, city := range Cities {
		add("/jasa-whatsapp-marketing-"+city.Slug, "0.6")
		add("/aplikasi-whatsapp-"+city.Slug, "0.6")
	}

	// 10. Chatbot × Industry × City
	for _, cb := range ChatbotPages {
		for _, ind := range Industries {
			for _, city := range Cities {
				add(fmt.Sprintf("/chatbot-%s-untuk-%s-di-%s", cb.Slug, ind.Slug, city.Slug), "0.5")
			}
		}
	}

	// 11. Panduan × Industry × City
	for _, pp := range PanduanPages {
		for _, ind := range Industries {
			for _, city := range Cities {
				add(fmt.Sprintf("/panduan-%s-untuk-%s-di-%s", pp.Slug, ind.Slug, city.Slug), "0.5")
			}
		}
	}

	// 12. Feature pages
	for _, feat := range Features {
		add("/best-whatsapp-"+feat.Slug, "0.6")
		for _, ind := range Industries {
			add(fmt.Sprintf("/best-whatsapp-%s-untuk-%s", feat.Slug, ind.Slug), "0.6")
		}
	}

	fmt.Fprint(w, "</urlset>")
}

func HandleRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, `User-agent: *
Allow: /$
Allow: /docs
Allow: /marketing/
Allow: /screens/
Allow: /best-
Allow: /alternatives-to-
Allow: /compare/
Allow: /whatsapp-marketing-untuk-
Allow: /beli-aplikasi-
Allow: /source-code-
Allow: /aplikasi-whatsapp-
Allow: /jual-aplikasi-
Allow: /jual-source-code-
Allow: /harga-source-code-
Allow: /jasa-whatsapp-
Allow: /cara-
Allow: /chatbot-
Allow: /panduan-
Allow: /sitemap
Allow: /sitemaps/
Allow: /indexnow-key.txt
Allow: /login
Allow: /register
Disallow: /admin
Disallow: /api
Disallow: /__pair
Disallow: /webhooks
Disallow: /wa
Disallow: /send
Disallow: /sent
Disallow: /received
Disallow: /inbox
Disallow: /contacts
Disallow: /broadcast
Disallow: /autoreply
Disallow: /settings

Sitemap: /sitemap.xml
`)
}

func CountPSEOURLs() int {
	c := 0

	// Core
	c += 4 + len(BestPages)*2 + len(Competitors) + len(Features) + len(SourceCodePages) + len(CaraPages) + len(ChatbotPages) + len(PanduanPages)

	// Industry × City × 2
	c += len(Industries) * len(Cities) * 2

	// Compare pairs + industry
	c += len(Competitors) + len(Competitors)*len(Industries)

	// Compare × Industry × City
	c += len(Competitors) * len(Industries) * len(Cities)

	// Alternatives × Industry × City
	c += len(Competitors) * len(Industries) * len(Cities)

	// Source code × Industry × City (12 prefixes)
	c += 12 * len(Industries) * len(Cities)

	// Cara × Feature × Industry × City (4-way)
	c += len(CaraPages) * len(Features) * len(Industries) * len(Cities)

	// Best × Industry × City × Year (3 years)
	c += len(BestPages) * len(Industries) * len(Cities) * 3

	// City pages
	c += len(Cities) * 2

	// Chatbot × Industry × City
	c += len(ChatbotPages) * len(Industries) * len(Cities)

	// Panduan × Industry × City
	c += len(PanduanPages) * len(Industries) * len(Cities)

	// Feature × Industry
	c += len(Features) + len(Features)*len(Industries)

	return c
}
