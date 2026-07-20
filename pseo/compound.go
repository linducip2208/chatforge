package pseo

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type parsedSlug struct {
	Type      string // best, alternative, compare, cara, etc.
	Keyword   string
	Industry  string
	City      string
	Year      int
	Feature   string
	ToolA     string
	ToolB     string
	IsCompound bool
}

func parseCompoundSlug(path string) parsedSlug {
	p := parsedSlug{}
	slug := strings.TrimPrefix(path, "/")
	slug = strings.TrimSuffix(slug, "/")

	// Determine type and base slug
	switch {
	case strings.HasPrefix(slug, "best-"):
		p.Type = "best"
		p.Keyword = strings.TrimPrefix(slug, "best-")
	case strings.HasPrefix(slug, "alternatives-to-"):
		p.Type = "alternative"
		p.Keyword = strings.TrimPrefix(slug, "alternatives-to-")
	case strings.HasPrefix(slug, "compare/"):
		p.Type = "compare"
		p.Keyword = strings.TrimPrefix(slug, "compare/")
		parts := strings.Split(p.Keyword, "-vs-")
		if len(parts) == 2 {
			p.ToolA = parts[0]
			// Remove compound suffix from toolB
			p.ToolB = parts[1]
		}
	case strings.HasPrefix(slug, "whatsapp-marketing-untuk-"):
		p.Type = "industry"
		p.Keyword = strings.TrimPrefix(slug, "whatsapp-marketing-untuk-")
	case strings.HasPrefix(slug, "beli-aplikasi-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "beli-aplikasi-")
	case strings.HasPrefix(slug, "source-code-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "source-code-")
	case strings.HasPrefix(slug, "aplikasi-whatsapp-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "aplikasi-whatsapp-")
	case strings.HasPrefix(slug, "jual-aplikasi-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "jual-aplikasi-")
	case strings.HasPrefix(slug, "jual-source-code-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "jual-source-code-")
	case strings.HasPrefix(slug, "harga-source-code-"):
		p.Type = "sourcecode"
		p.Keyword = strings.TrimPrefix(slug, "harga-source-code-")
	case strings.HasPrefix(slug, "jasa-whatsapp-"):
		p.Type = "cityservice"
		p.Keyword = strings.TrimPrefix(slug, "jasa-whatsapp-")
	case strings.HasPrefix(slug, "cara-"):
		p.Type = "cara"
		p.Keyword = strings.TrimPrefix(slug, "cara-")
	case strings.HasPrefix(slug, "chatbot-"):
		p.Type = "chatbot"
		p.Keyword = strings.TrimPrefix(slug, "chatbot-")
	case strings.HasPrefix(slug, "panduan-"):
		p.Type = "panduan"
		p.Keyword = strings.TrimPrefix(slug, "panduan-")
	default:
		return p
	}

	// Extract year from keyword suffix
	for y := 2024; y <= 2030; y++ {
		ys := strconv.Itoa(y)
		if strings.HasSuffix(p.Keyword, "-"+ys) {
			p.Keyword = strings.TrimSuffix(p.Keyword, "-"+ys)
			p.Year = y
			break
		}
	}

	// Extract industry and city from keyword (compound patterns)
	// Pattern: {keyword}-untuk-{industri}-di-{city}
	if idx := strings.Index(p.Keyword, "-untuk-"); idx >= 0 {
		p.IsCompound = true
		rest := p.Keyword[idx+7:] // after "-untuk-"
		p.Keyword = p.Keyword[:idx]

		// Check if rest contains "-di-" for city
		if diIdx := strings.Index(rest, "-di-"); diIdx >= 0 {
			indSlug := rest[:diIdx]
			citySlug := rest[diIdx+4:] // after "-di-"

			// Match industry
			for _, ind := range Industries {
				if ind.Slug == indSlug {
					p.Industry = ind.Name
					break
				}
			}
			if p.Industry == "" {
				p.Industry = humanize(indSlug)
			}

			// Match city
			for _, c := range Cities {
				if c.Slug == citySlug {
					p.City = c.Name
					break
				}
			}
			if p.City == "" {
				p.City = humanize(citySlug)
			}
		} else {
			// Only industry, no city
			for _, ind := range Industries {
				if ind.Slug == rest {
					p.Industry = ind.Name
					break
				}
			}
			if p.Industry == "" {
				// Check if it's a city
				for _, c := range Cities {
					if c.Slug == rest {
						p.City = c.Name
						break
					}
				}
			}
			if p.Industry == "" && p.City == "" {
				p.Industry = humanize(rest)
			}
		}
	}

	// For compare pages, toolB might have compound suffix
	if p.Type == "compare" && strings.Contains(p.ToolB, "-untuk-") {
		p.IsCompound = true
		if idx := strings.Index(p.ToolB, "-untuk-"); idx >= 0 {
			rest := p.ToolB[idx+7:]
			p.ToolB = p.ToolB[:idx]
			for _, ind := range Industries {
				if ind.Slug == rest {
					p.Industry = ind.Name
					break
				}
			}
			if p.Industry == "" {
				p.Industry = humanize(rest)
			}
		}
	}

	// For alternatives, check compound
	if p.Type == "alternative" && strings.Contains(p.Keyword, "-untuk-") {
		p.IsCompound = true
		// Already handled by the industry extraction above (keyword prefix is the tool name)
		// The keyword here is actually {tool}-untuk-{industry}
		// But wait, the alternatives prefix already removed "alternatives-to-"
		// So keyword = "wablas-untuk-restoran" -> tool = "wablas", industry = "restoran"
		if idx := strings.Index(p.Keyword, "-untuk-"); idx >= 0 {
			p.Industry = humanize(p.Keyword[idx+7:])
			// Re-match with actual industry data
			for _, ind := range Industries {
				if ind.Slug == p.Keyword[idx+7:] {
					p.Industry = ind.Name
					break
				}
			}
			p.Keyword = p.Keyword[:idx]
		}
	}

	return p
}

func renderCompoundPage(w http.ResponseWriter, r *http.Request, p parsedSlug) {
	var title, desc string
	var items []ListItem
	var content template.HTML
	isList := false

	switch p.Type {
	case "best":
		title = fmt.Sprintf("%d Rekomendasi %s untuk %s di %s %d", len(Competitors)+1, humanize(p.Keyword), p.Industry, p.City, p.Year)
		if p.Year == 0 {
			p.Year = time.Now().Year()
		}
		desc = fmt.Sprintf("Daftar %s untuk bisnis %s di %s tahun %d. Bandingkan fitur, harga, dan kelebihan. %s self-hosted tanpa biaya bulanan.", humanize(p.Keyword), p.Industry, p.City, p.Year, AppName)
		content = compoundBestContent(p)
		isList = true
		items = generateFeatureItems()

	case "alternative":
		title = fmt.Sprintf("5 Alternatif %s untuk %s di %s", humanize(p.Keyword), p.Industry, p.City)
		desc = fmt.Sprintf("Cari alternatif %s untuk bisnis %s di %s? Bandingkan fitur, harga, dan pilih yang paling cocok.", humanize(p.Keyword), p.Industry, p.City)
		content = compoundAlternativeContent(p)
		isList = true
		items = generateAlternativeItems(p)

	case "compare":
		title = fmt.Sprintf("%s vs %s untuk %s: Perbandingan Lengkap", humanize(p.ToolA), humanize(p.ToolB), p.Industry)
		desc = fmt.Sprintf("Bandingkan %s vs %s untuk kebutuhan %s. Fitur, harga, kelebihan — mana yang lebih cocok untuk bisnis Anda?", humanize(p.ToolA), humanize(p.ToolB), p.Industry)
		content = compoundCompareContent(p)

	case "industry":
		title = fmt.Sprintf("WhatsApp Marketing untuk %s di %s: Solusi Lengkap", p.Industry, p.City)
		desc = fmt.Sprintf("Optimalkan WhatsApp marketing untuk bisnis %s di %s. Auto-reply, broadcast promo, chatbot CS 24/7 — dalam satu platform.", p.Industry, p.City)
		content = compoundIndustryContent(p)

	case "sourcecode":
		title = fmt.Sprintf("%s untuk %s di %s — Source Code Self-Hosted", humanize(p.Keyword), p.Industry, p.City)
		desc = fmt.Sprintf("Beli source code WhatsApp %s untuk bisnis %s di %s. Self-hosted, bayar sekali, lifetime update. Termasuk dokumentasi & support.", humanize(p.Keyword), p.Industry, p.City)
		content = compoundSourceCodeContent(p)
		isList = true
		items = generateFeatureItems()

	case "cara":
		title = fmt.Sprintf("Cara %s untuk %s di %s", humanize(p.Keyword), p.Industry, p.City)
		desc = fmt.Sprintf("Panduan lengkap cara %s untuk bisnis %s di %s. Step-by-step dengan screenshot.", humanize(p.Keyword), p.Industry, p.City)
		content = compoundCaraContent(p)

	case "cityservice":
		if p.Industry != "" {
			title = fmt.Sprintf("Jasa %s untuk %s di %s", humanize(p.Keyword), p.Industry, p.City)
			desc = fmt.Sprintf("Butuh jasa %s untuk bisnis %s di %s? Solusi self-hosted lebih hemat. Source code + training termasuk.", humanize(p.Keyword), p.Industry, p.City)
		} else {
			title = fmt.Sprintf("Jasa %s di %s", humanize(p.Keyword), p.City)
			desc = fmt.Sprintf("Cari jasa %s di %s? Lebih hemat dengan self-hosted source code. Tanpa biaya bulanan.", humanize(p.Keyword), p.City)
		}
		content = compoundCityContent(p)
	}

	breadcrumbs := []Breadcrumb{{"Home", "/"}}
	if p.Type != "" {
		breadcrumbs = append(breadcrumbs, Breadcrumb{humanize(p.Type), ""})
	}

	renderPSEO(w, r, PageData{
		Title:       title,
		Description: desc,
		IsList:      isList,
		Items:       items,
		Breadcrumbs: breadcrumbs,
		Content:     ctaBanner() + content,
	})
}

func generateFeatureItems() []ListItem {
	return []ListItem{
		{"Broadcast & Blast Massal", "/best-whatsapp-broadcast-tools", "Kirim pesan ke ribuan kontak sekaligus dengan interval aman."},
		{"Auto Reply AI", "/best-whatsapp-auto-reply", "Balas otomatis 24/7 dengan AI OpenAI/DeepSeek/Gemini."},
		{"Chatbot Pintar", "/best-whatsapp-chatbot", "Layanan CS otomatis dengan memory & FAQ search."},
		{"Multi-Account WhatsApp", "/best-whatsapp-gateway", "Kelola puluhan nomor WA dalam satu dashboard."},
		{"REST API + Webhook", "/best-whatsapp-api", "Integrasi dengan CRM, e-commerce, sistem lain."},
		{"Drip Campaign", "/best-whatsapp-marketing-tools", "Rangkaian pesan otomatis bertahap nurture leads."},
	}
}

func generateAlternativeItems(p parsedSlug) []ListItem {
	items := []ListItem{
		{Title: AppName + " (Rekomendasi Utama)", URL: "/", Description: "Self-hosted, bayar sekali, semua fitur unlimited."},
	}
	count := 0
	for _, c := range Competitors {
		if strings.EqualFold(c, p.Keyword) {
			continue
		}
		count++
		if count > 5 {
			break
		}
		if p.Industry != "" {
			items = append(items, ListItem{
				Title:       humanize(c),
				URL:         "/compare/" + safeSlug(AppName) + "-vs-" + c + "-untuk-" + safeSlug(p.Industry),
				Description: fmt.Sprintf("Bandingkan %s vs %s untuk %s.", AppName, humanize(c), p.Industry),
			})
		} else {
			items = append(items, ListItem{
				Title:       humanize(c),
				URL:         "/compare/" + safeSlug(AppName) + "-vs-" + c,
				Description: fmt.Sprintf("Bandingkan %s vs %s secara detail.", AppName, humanize(c)),
			})
		}
	}
	return items
}
