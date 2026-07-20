package pseo

import (
	"fmt"
	"html/template"
	"strings"
)

func bestContent(title, pageTitle string, year int, keyword string) template.HTML {
	items := make([]string, 0, len(Competitors)+1)
	items = append(items, fmt.Sprintf(`<li><strong>%s</strong> — Pilihan #1. Self-hosted, bayar sekali, lifetime update. Broadcast massal, auto-reply AI (OpenAI/DeepSeek/Gemini), multi-account WhatsApp, chatbot, REST API, webhook. Tanpa biaya bulanan.</li>`, AppName))
	for i, c := range Competitors {
		items = append(items, fmt.Sprintf(`<li><strong>%s</strong> — Alternatif populer untuk %s. Bandingkan <a href="/compare/%s-vs-%s">%s vs %s</a>.</li>`, humanize(c), keyword, safeSlug(AppName), c, AppName, humanize(c)))
		if i >= 8 {
			items = append(items, fmt.Sprintf(`<li>...dan %d tools lainnya.</li>`, len(Competitors)-9))
			break
		}
	}

	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Mengapa %s Jadi Pilihan Utama?</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Di tahun %d, kebutuhan WhatsApp marketing semakin tinggi. Bisnis butuh tools yang <strong>handal, murah, dan fleksibel</strong>. %s hadir sebagai solusi <strong>self-hosted</strong> — Anda punya kendali penuh atas data, tidak ada biaya bulanan, dan semua fitur tersedia tanpa batasan.</p>

<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">💡 Self-Hosted</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Install di server sendiri. Data aman, kendali penuh.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">⚡ Multi-Account</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Jalankan puluhan nomor WA sekaligus.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">🤖 AI Auto-Reply</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">OpenAI, DeepSeek, Gemini — BYOK.</p></div>
<div style="background:#fdf2f8;border:1px solid #fbcfe8;border-radius:12px;padding:20px"><strong style="color:#831843">💰 One-Time Purchase</strong><p style="color:#be185d;margin:8px 0 0;font-size:.9rem">Bayar sekali, pakai selamanya.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Daftar Lengkap %s</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px">%s</ol>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Fitur yang Wajib Ada di Tools %s</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>✅ Broadcast massal dengan interval configurable</li>
<li>✅ Multi-account WhatsApp (kelola banyak nomor)</li>
<li>✅ Auto-reply berbasis keyword + AI</li>
<li>✅ REST API + Webhook untuk integrasi</li>
<li>✅ Dashboard real-time & laporan analitik</li>
<li>✅ Manajemen kontak (CSV import, tag, group)</li>
<li>✅ Drip campaign multi-step</li>
<li>✅ Harga transparan, tanpa hidden fee</li>
</ul>
</div>`, AppName, year, AppName, keyword, strings.Join(items, "\n"), keyword))
}

func alternativeContent(toolName string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Mengapa Mencari Alternatif %s?</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Banyak pengguna %s mencari alternatif karena <strong>harga berlangganan yang mahal</strong>, <strong>fitur terbatas</strong> di paket murah, atau <strong>batasan jumlah pesan</strong> per bulan. %s menjawab semua masalah ini dengan solusi <strong>self-hosted one-time purchase</strong>.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Perbandingan %s vs %s</h2>
<div style="overflow-x:auto;margin-bottom:32px">
<table style="width:100%%;border-collapse:collapse;font-size:.9rem">
<thead><tr style="background:#f8f9fc"><th style="padding:12px;text-align:left;border-bottom:2px solid #e0e4e9">Fitur</th><th style="padding:12px;text-align:center;border-bottom:2px solid #e0e4e9">%s</th><th style="padding:12px;text-align:center;border-bottom:2px solid #4F46E5;color:#4F46E5">%s</th></tr></thead>
<tbody>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Model Harga</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Berlangganan / bulan</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">One-time purchase</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Multi-Account WA</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Unlimited</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Broadcast</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">✅</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ + Round-Robin</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">AI Auto-Reply</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">❌ / Terbatas</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ BYOK (OpenAI/DeepSeek/Gemini)</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Chatbot</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Full AI + FAQ Search</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">REST API</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Paket mahal</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Included</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Self-Hosted</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">❌</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ 100%%</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Source Code</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">❌</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Termasuk</td></tr>
</tbody></table></div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kesimpulan</h2>
<p style="color:#5e6b7e;line-height:1.8">Jika Anda bosan dengan biaya berlangganan %s yang terus naik, %s adalah <strong>alternatif terbaik</strong>. Dengan one-time purchase, Anda dapat semua fitur tanpa batasan. <strong>Cocok untuk agency, UMKM, dan enterprise</strong> yang butuh solusi WhatsApp marketing handal tanpa biaya bulanan.</p>
</div>`, toolName, toolName, AppName, AppName, toolName, toolName, AppName, toolName, AppName))
}

func compareContent(a, b string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s vs %s: Perbandingan Head-to-Head</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Membandingkan <strong>%s</strong> dengan <strong>%s</strong> untuk kebutuhan WhatsApp marketing bisnis Anda. Kami evaluasi dari segi <strong>fitur, harga, kemudahan penggunaan, dan skalabilitas</strong>.</p>

<div style="overflow-x:auto;margin-bottom:32px">
<table style="width:100%%;border-collapse:collapse;font-size:.9rem">
<thead><tr style="background:#f8f9fc"><th style="padding:12px;text-align:left;border-bottom:2px solid #e0e4e9">Kategori</th><th style="padding:12px;text-align:center;border-bottom:2px solid #4F46E5;color:#4F46E5">%s</th><th style="padding:12px;text-align:center;border-bottom:2px solid #e0e4e9">%s</th></tr></thead>
<tbody>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Model Bisnis</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Self-Hosted · One-Time</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Cloud SaaS · Subscription</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Harga</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Sekali bayar</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Per bulan / tahun</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Multi WA Account</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Unlimited</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Tergantung paket</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Broadcast</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Round-Robin</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">✅ Standard</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>AI Auto-Reply</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ BYOK Multi-Model</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Tergantung paket</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Drip Campaign</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Multi-Step</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>REST API</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Included</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Paket Enterprise</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Meta Cloud API</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Built-in</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">✅</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Source Code</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Full Access</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">❌ Proprietary</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee"><strong>Kustomisasi</strong></td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ Unlimited</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas / Tidak bisa</td></tr>
</tbody></table></div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Verdict: Mana yang Lebih Baik?</h2>
<div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-bottom:24px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px">
<h3 style="font-size:1.1rem;font-weight:700;color:#166534;margin:0 0 8px">✅ Kelebihan %s</h3>
<ul style="color:#15803d;margin:0;padding-left:20px;font-size:.9rem;line-height:1.8">
<li>One-time purchase, tanpa biaya bulanan</li>
<li>Full source code — kustomisasi bebas</li>
<li>Unlimited WA account & broadcast</li>
<li>AI BYOK — pakai API key sendiri</li>
<li>Self-hosted — data 100%% milik Anda</li>
</ul></div>
<div style="background:#fff7ed;border:1px solid #fed7aa;border-radius:12px;padding:20px">
<h3 style="font-size:1.1rem;font-weight:700;color:#9a3412;margin:0 0 8px">⚠️ Kekurangan %s</h3>
<ul style="color:#c2410c;margin:0;padding-left:20px;font-size:.9rem;line-height:1.8">
<li>Butuh server/VPS sendiri</li>
<li>Setup awal perlu teknisi</li>
<li>Tidak ada free tier cloud</li>
</ul></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">FAQ: %s vs %s</h2>
<div style="margin-bottom:32px">
<div style="background:#f8f9fc;border-radius:12px;padding:20px;margin-bottom:12px">
<strong style="color:#152e4d">Q: Apakah %s lebih murah dari %s?</strong>
<p style="color:#5e6b7e;margin:8px 0 0">Ya. %s one-time purchase — setelah 3-6 bulan, Anda sudah balik modal dibanding biaya berlangganan %s. Jangka panjang jauh lebih hemat.</p>
</div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px;margin-bottom:12px">
<strong style="color:#152e4d">Q: Apakah perlu teknisi untuk setup?</strong>
<p style="color:#5e6b7e;margin:8px 0 0">Setup awal sederhana: upload binary ke VPS, setting .env, jalankan. Dokumentasi lengkap tersedia di /docs. Bisa juga minta bantuan instalasi via WhatsApp.</p>
</div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px">
<strong style="color:#152e4d">Q: Apakah fiturnya selengkap %s?</strong>
<p style="color:#5e6b7e;margin:8px 0 0">Bahkan lebih lengkap: multi-account, AI auto-reply, drip campaign, chatbot, REST API, Meta Cloud API — semua included tanpa batasan.</p>
</div>
</div>
</div>`, a, b, a, b, a, b, a, b, a, b, AppName, AppName, b, AppName, b))
}

func industryContent(name string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Mengapa Bisnis %s Butuh WhatsApp Marketing?</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Di Indonesia, <strong>90%%+ pengguna smartphone</strong> menggunakan WhatsApp setiap hari. Untuk bisnis %s, WhatsApp bukan sekadar chat — ini adalah <strong>channel marketing, sales, dan customer service</strong> dalam satu platform.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Fitur yang Cocok untuk Bisnis %s</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📢 Broadcast Promo</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Kirim penawaran spesial ke semua pelanggan dalam satu klik.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🤖 CS 24/7</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Chatbot AI jawab pertanyaan pelanggan kapan saja.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📅 Booking & Reminder</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Jadwal + reminder otomatis via WhatsApp.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📊 Analitik</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Pantau performa: sent, delivered, read, replied.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Tips WhatsApp Marketing untuk %s</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px">
<li><strong>Segmentasi kontak</strong> — kelompokkan pelanggan berdasarkan minat, lokasi, atau riwayat pembelian</li>
<li><strong>Gunakan template personal</strong> — variabel {name} membuat pesan terasa personal</li>
<li><strong>Atur jadwal broadcast</strong> — kirim di jam yang tepat (pagi 8-10, siang 12-14, sore 16-18)</li>
<li><strong>Manfaatkan drip campaign</strong> — rangkaian pesan bertahap untuk nurture leads</li>
<li><strong>Aktifkan auto-reply AI</strong> — jawab FAQ otomatis, hemat waktu tim CS</li>
<li><strong>Pantau analytics</strong> — evaluasi performa campaign secara berkala</li>
</ol>
</div>`, name, name, name, name))
}

func sourceCodeContent(title string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Anda bisa <strong>memiliki aplikasi WhatsApp marketing sendiri</strong> tanpa tergantung pihak ketiga. Source code lengkap, self-hosted di server Anda, <strong>bayar sekali — pakai selamanya</strong>.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Apa yang Anda Dapatkan?</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">📦 Source Code Lengkap</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Full source code Go — bisa dikustomisasi sesuai kebutuhan.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">📖 Dokumentasi</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Dokumentasi instalasi & penggunaan lengkap.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">🔄 Lifetime Update</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">Update gratis seumur hidup via GitHub.</p></div>
<div style="background:#fdf2f8;border:1px solid #fbcfe8;border-radius:12px;padding:20px"><strong style="color:#831843">🎓 Support & Training</strong><p style="color:#be185d;margin:8px 0 0;font-size:.9rem">Bantuan instalasi + training penggunaan.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Fitur Unggulan</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px;margin-bottom:24px">
<li>✅ WhatsApp Multi-Account (unlimited nomor)</li>
<li>✅ Broadcast & Blast Massal (round-robin, anti-banned)</li>
<li>✅ Auto Reply AI (OpenAI / DeepSeek / Gemini / Claude)</li>
<li>✅ Chatbot AI + FAQ + Training Campaign</li>
<li>✅ Drip Campaign Multi-Step</li>
<li>✅ Scheduled & Recurring Messages</li>
<li>✅ Live Chat Inbox Real-Time (SSE)</li>
<li>✅ REST API + Webhook</li>
<li>✅ Meta Cloud API (WhatsApp Business API)</li>
<li>✅ Multi-User SaaS (admin, roles, packages, subscriptions)</li>
<li>✅ Payment Gateway (Midtrans, Xendit, PayPal, Stripe)</li>
<li>✅ Dashboard Analytics + Chart</li>
<li>✅ Single Binary — No Dependency</li>
</ul>
</div>`, title))
}

func cityContent(cityName string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">WhatsApp Marketing di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Bisnis di <strong>%s</strong> butuh solusi WhatsApp marketing yang handal dan terjangkau. Dengan source code self-hosted, Anda bisa <strong>menjalankan sendiri</strong> platform WhatsApp marketing tanpa ketergantungan pihak ketiga.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kenapa Pilih Self-Hosted untuk Bisnis di %s?</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">💡 Kendali Penuh</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Server sendiri, data sendiri, aturan sendiri.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">⚡ Hemat Biaya</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Bayar sekali, tanpa biaya bulanan forever.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">🎯 Support Lokal</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">Training & support via WhatsApp dalam bahasa Indonesia.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Cocok untuk Bisnis di %s</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px;margin-bottom:24px">
<li>🏪 Toko retail / grosir — broadcast promo ke reseller</li>
<li>🍽️ Restoran / cafe — auto-reply menu & reservasi</li>
<li>🏥 Klinik / apotek — reminder janji temu & resep</li>
<li>🏠 Agen properti — kirim listing terbaru ke buyer</li>
<li>👗 Fashion / butik — katalog & order via WhatsApp</li>
<li>📦 Distributor — blast harga update ke agen</li>
</ul>
</div>`, cityName, cityName, cityName, cityName))
}

func caraContent(title, desc string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">%s</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Langkah Demi Langkah</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px;margin-bottom:24px">
<li><strong>Persiapan Server</strong> — Siapkan VPS/Linux server. %s adalah single binary Go — cukup upload file dan jalankan.</li>
<li><strong>Setup Database MySQL</strong> — Buat database, setting koneksi di file .env. Auto-migrasi ~40 tabel.</li>
<li><strong>Jalankan Aplikasi</strong> — ./chatgo atau chatgo.exe. Akses via browser di port 8080.</li>
<li><strong>Hubungkan WhatsApp</strong> — Scan QR code via WhatsApp mobile (Linked Devices). Multi-account support.</li>
<li><strong>Import Kontak</strong> — Upload CSV atau tambah manual. Grouping & tagging untuk segmentasi.</li>
<li><strong>Setup Auto Reply</strong> — Buat rule keyword atau aktifkan AI mode (OpenAI/DeepSeek/Gemini).</li>
<li><strong>Mulai Broadcast</strong> — Pilih grup target, tulis pesan, kirim massal dengan interval aman.</li>
<li><strong>Pantau Dashboard</strong> — Real-time stats: sent, delivered, read, replied.</li>
</ol>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Tips & Trik</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>Gunakan <strong>multi-account</strong> untuk menghindari limit WhatsApp (250 broadcast/hari/akun)</li>
<li>Aktifkan <strong>rate limiter</strong> di settings untuk interval aman antar pesan</li>
<li>Manfaatkan <strong>template + spintax</strong> untuk variasi pesan otomatis</li>
<li>Setup <strong>webhook</strong> untuk integrasi real-time dengan sistem lain</li>
<li>Gunakan <strong>AI fallback mode</strong> — hanya jawab pakai AI jika tidak ada keyword cocok</li>
</ul>
</div>`, title, desc, AppName))
}

func chatbotContent(title, desc string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">%s Dengan %s, Anda bisa memiliki <strong>chatbot WhatsApp AI</strong> sendiri yang berjalan 24/7, menjawab pertanyaan customer, memproses order, dan memberikan rekomendasi — semuanya otomatis.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kemampuan Chatbot AI</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🧠 AI Multi-Model</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">OpenAI, DeepSeek, Gemini, Claude — pilih sesuai budget.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📚 Knowledge Base</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Upload FAQ, CSV, PDF, URL — AI search sebelum jawab.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🔄 Memory Window</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Ingat konteks percakapan N chat terakhir.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">👨‍💼 Human Handoff</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Keyword trigger → AI berhenti → transfer ke agent manusia.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🕐 Jam Kerja</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">AI hanya aktif di jam kerja. Di luar jam — balasan otomatis.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🛡️ Anti-Spam</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Deteksi spam, mute 30 menit, anti-jailbreak regex.</p></div>
</div>
</div>`, title, desc, AppName))
}

func panduanContent(title, desc string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">%s</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Isi Panduan</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px;margin-bottom:24px">
<li><strong>Pengenalan WhatsApp Marketing</strong> — Kenapa WhatsApp adalah channel marketing #1 di Indonesia</li>
<li><strong>Persiapan Infrastruktur</strong> — Pilih server/VPS, setup domain, install MySQL</li>
<li><strong>Installasi %s</strong> — Single binary, setting .env, auto-migration database</li>
<li><strong>Koneksi WhatsApp</strong> — Scan QR, multi-account, Meta Cloud API</li>
<li><strong>Manajemen Kontak</strong> — Import CSV, segmentasi grup, tagging, merge duplikat</li>
<li><strong>Auto Reply & AI</strong> — Keyword rules, AI config, training campaign, knowledge base</li>
<li><strong>Broadcast Campaign</strong> — Setup campaign, pilih target, round-robin, interval</li>
<li><strong>Analitik & Optimasi</strong> — Dashboard chart, campaign analytics, agent performance</li>
</ol>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Tools yang Dibutuhkan</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>VPS/Linux server (minimal 1GB RAM) atau Windows Server</li>
<li>MySQL 5.7+ / MariaDB 10.3+</li>
<li>Domain (opsional, untuk production)</li>
<li>Nomor WhatsApp (bisa nomor biasa, tidak perlu API Business)</li>
<li>API Key OpenAI/DeepSeek (opsional, untuk fitur AI)</li>
</ul>
</div>`, title, desc, AppName))
}
