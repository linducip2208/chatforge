package pseo

import (
	"fmt"
	"html/template"
)

func compoundBestContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s Terbaik untuk %s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Bisnis <strong>%s</strong> di <strong>%s</strong> butuh tools %s yang handal. Di tahun %d, persaingan makin ketat — Anda perlu solusi yang <strong>cepat, murah, dan bisa dikustomisasi</strong>. %s adalah jawabannya: self-hosted, bayar sekali, semua fitur tanpa batasan.</p>

<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">💡 Self-Hosted</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Install di server sendiri. Data aman 100%% milik Anda.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">🎯 Kustomisasi</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Full source code — modifikasi sesuai kebutuhan %s.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">💰 Hemat</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">One-time purchase. Setelah 3 bulan sudah balik modal.</p></div>
<div style="background:#fdf2f8;border:1px solid #fbcfe8;border-radius:12px;padding:20px"><strong style="color:#831843">🏢 Lokal</strong><p style="color:#be185d;margin:8px 0 0;font-size:.9rem">Support bahasa Indonesia. Training untuk tim di %s.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Mengapa %s Cocok untuk Bisnis %s?</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Setiap industri punya kebutuhan berbeda. Untuk <strong>%s</strong>, WhatsApp marketing harus bisa: broadcast promo rutin, auto-reply pertanyaan customer, segmentasi kontak berdasarkan lokasi & perilaku, dan integrasi dengan sistem existing. %s menyediakan semua itu — dalam satu platform.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Tips %s untuk Bisnis %s di %s</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>✅ Gunakan <strong>multi-account WA</strong> — satu nomor untuk CS, satu untuk marketing</li>
<li>✅ Setup <strong>auto-reply AI</strong> — jawab FAQ otomatis, hemat biaya CS</li>
<li>✅ Broadcast promo <strong>per segment</strong> — jangan blast semua kontak sekaligus</li>
<li>✅ <strong>Pantau analytics</strong> — evaluasi performa campaign mingguan</li>
<li>✅ Manfaatkan <strong>drip campaign</strong> — nurture leads dari inquiry sampai closing</li>
</ul>
</div>`, humanize(p.Keyword), p.Industry, p.City, p.Industry, p.City, humanize(p.Keyword), p.Year, AppName, p.Industry, p.City, AppName, p.Industry, p.Industry, AppName, humanize(p.Keyword), p.Industry, p.City))
}

func compoundAlternativeContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Alternatif %s untuk %s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Banyak bisnis <strong>%s</strong> di <strong>%s</strong> mencari alternatif <strong>%s</strong> karena alasan harga, fitur terbatas, atau tidak fleksibel. %s hadir sebagai solusi <strong>self-hosted</strong> yang memberikan kendali penuh tanpa biaya bulanan.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Perbandingan untuk Kebutuhan %s</h2>
<div style="overflow-x:auto;margin-bottom:32px">
<table style="width:100%%;border-collapse:collapse;font-size:.9rem">
<thead><tr style="background:#f8f9fc"><th style="padding:12px;text-align:left;border-bottom:2px solid #e0e4e9">Fitur</th><th style="padding:12px;text-align:center;border-bottom:2px solid #e0e4e9">%s</th><th style="padding:12px;text-align:center;border-bottom:2px solid #4F46E5;color:#4F46E5">%s</th></tr></thead>
<tbody>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Model Harga</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Berlangganan</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">One-Time Purchase</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Kustomisasi</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Full Source Code</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Multi-Account</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Terbatas</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">Unlimited</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">AI Auto-Reply</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Tidak termasuk</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">✅ BYOK Multi-Model</td></tr>
<tr><td style="padding:10px;border-bottom:1px solid #eee">Support Lokal</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee">Bahasa Inggris</td><td style="padding:10px;text-align:center;border-bottom:1px solid #eee;background:#f5f3ff;font-weight:600">🇮🇩 Bahasa Indonesia</td></tr>
</tbody></table></div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kenapa Bisnis %s di %s Pilih %s?</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>💡 <strong>Kendali penuh:</strong> install di server sendiri, atur sendiri</li>
<li>💰 <strong>Hemat jangka panjang:</strong> tidak ada biaya bulanan</li>
<li>🔧 <strong>Kustomisasi bebas:</strong> full source code, modifikasi sesuka hati</li>
<li>📞 <strong>Support lokal:</strong> training & bantuan dalam bahasa Indonesia</li>
<li>🚀 <strong>Cocok untuk %s:</strong> fitur broadcast, chatbot, auto-reply siap pakai</li>
</ul>
</div>`, humanize(p.Keyword), p.Industry, p.City, p.Industry, p.City, humanize(p.Keyword), AppName, p.Industry, humanize(p.Keyword), AppName, p.Industry, p.City, AppName, p.Industry))
}

func compoundCompareContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s vs %s untuk %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Mana yang lebih cocok untuk bisnis <strong>%s</strong>: <strong>%s</strong> atau <strong>%s</strong>? Kami bandingkan berdasarkan fitur yang relevan untuk industri %s.</p>

<div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-bottom:24px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px">
<h3 style="font-size:1.1rem;font-weight:700;color:#166534;margin:0 0 8px">✅ %s: Kelebihan</h3>
<ul style="color:#15803d;margin:0;padding-left:20px;font-size:.9rem;line-height:1.8">
<li>One-time purchase, tanpa biaya bulanan</li>
<li>Full source code — kustomisasi bebas</li>
<li>Support bahasa Indonesia</li>
<li>Training & setup bantuan</li>
</ul></div>
<div style="background:#fff7ed;border:1px solid #fed7aa;border-radius:12px;padding:20px">
<h3 style="font-size:1.1rem;font-weight:700;color:#9a3412;margin:0 0 8px">%s: Perbandingan</h3>
<ul style="color:#c2410c;margin:0;padding-left:20px;font-size:.9rem;line-height:1.8">
<li>Model langganan — biaya terus menerus</li>
<li>Tidak bisa kustomisasi</li>
<li>Data di server pihak ketiga</li>
<li>Fitur terbatas di paket murah</li>
</ul></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kesimpulan untuk %s</h2>
<p style="color:#5e6b7e;line-height:1.8">Untuk kebutuhan <strong>%s</strong>, %s adalah pilihan yang lebih baik karena: <strong>self-hosted</strong> (data aman), <strong>one-time purchase</strong> (hemat jangka panjang), <strong>full source code</strong> (kustomisasi bebas), dan <strong>support lokal</strong>. Cocok untuk bisnis yang serius dengan WhatsApp marketing tanpa mau terikat biaya bulanan.</p>
</div>`, humanize(p.ToolA), humanize(p.ToolB), p.Industry, p.Industry, humanize(p.ToolA), humanize(p.ToolB), p.Industry, AppName, humanize(p.ToolB), p.Industry, p.Industry, AppName))
}

func compoundIndustryContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">WhatsApp Marketing untuk %s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Bisnis <strong>%s</strong> di <strong>%s</strong> punya tantangan unik: kompetitor banyak, customer makin demanding, dan biaya operasional terus naik. WhatsApp marketing dengan %s bisa jadi <strong>game-changer</strong> — efisien, personal, dan terukur.</p>

<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📢 Broadcast Promo</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Kirim promo, diskon, produk baru ke semua customer %s di %s.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">🤖 CS 24/7</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Chatbot AI jawab pertanyaan kapan saja — bahkan di luar jam kerja.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📅 Booking System</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Terima reservasi, janji temu, order — semua via WhatsApp.</p></div>
<div style="background:#f8f9fc;border-radius:12px;padding:20px"><strong style="color:#4F46E5">📊 Laporan</strong><p style="color:#5e6b7e;margin:8px 0 0;font-size:.9rem">Pantau performa campaign real-time. Data-driven decision.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Strategi WhatsApp Marketing untuk %s</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px">
<li><strong>Bangun database kontak</strong> — kumpulkan nomor WA customer, segmentasi per minat/lokasi</li>
<li><strong>Setup auto-reply cerdas</strong> — FAQ, jam operasional, menu produk — semua auto</li>
<li><strong>Broadcast terjadwal</strong> — kirim promo di jam prime (pagi, siang, sore)</li>
<li><strong>Drip campaign loyalitas</strong> — rangkaian pesan untuk customer baru & repeat order</li>
<li><strong>Integrasi sistem</strong> — hubungkan dengan POS, CRM, atau website via API</li>
</ol>
</div>`, p.Industry, p.City, p.Industry, p.City, AppName, p.Industry, p.City, p.Industry))
}

func compoundSourceCodeContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s untuk %s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Punya <strong>aplikasi WhatsApp marketing sendiri</strong> untuk bisnis %s di %s. Source code lengkap, self-hosted, <strong>bayar sekali — pakai selamanya</strong>. Cocok untuk agency, software house, atau bisnis yang ingin mandiri secara teknologi.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Mengapa Bisnis %s di %s Butuh Self-Hosted?</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">📦 Source Code</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Full Go source code — bisa dipelajari & dimodifikasi.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">📖 Dokumentasi</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Instalasi & penggunaan lengkap. Bahasa Indonesia.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">🔄 Lifetime Update</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">Update gratis seumur hidup. Fitur baru terus nambah.</p></div>
<div style="background:#fdf2f8;border:1px solid #fbcfe8;border-radius:12px;padding:20px"><strong style="color:#831843">🎓 Support</strong><p style="color:#be185d;margin:8px 0 0;font-size:.9rem">Bantuan instalasi & training untuk tim di %s.</p></div>
</div>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Fitur untuk %s</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>✅ Multi-Account WhatsApp (unlimited nomor)</li>
<li>✅ Broadcast & Blast Massal (anti-banned, round-robin)</li>
<li>✅ Auto Reply AI (OpenAI / DeepSeek / Gemini)</li>
<li>✅ Chatbot + FAQ Search + Training Campaign</li>
<li>✅ Live Chat Inbox Real-Time (SSE)</li>
<li>✅ REST API + Webhook + Meta Cloud API</li>
<li>✅ Multi-User SaaS (admin, roles, packages)</li>
<li>✅ Dashboard Analytics + Chart</li>
</ul>
</div>`, humanize(p.Keyword), p.Industry, p.City, p.Industry, p.City, p.Industry, p.City, p.City, p.Industry))
}

func compoundCaraContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Cara %s untuk %s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Panduan step-by-step untuk bisnis <strong>%s</strong> di <strong>%s</strong>. Dengan %s, Anda bisa setup sendiri — tanpa teknisi mahal.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Langkah Demi Langkah</h2>
<ol style="color:#5e6b7e;line-height:2;padding-left:20px;margin-bottom:24px">
<li><strong>Persiapan Server</strong> — Siapkan VPS/Linux murah (mulai Rp 100rb/bulan). %s single binary — cukup upload 1 file.</li>
<li><strong>Setup Database</strong> — MySQL/MariaDB, auto-migrasi ~40 tabel. Tinggal setting .env.</li>
<li><strong>Hubungkan WhatsApp</strong> — Scan QR di HP (Linked Devices). Bisa banyak nomor sekaligus.</li>
<li><strong>Import Data %s</strong> — Upload CSV kontak customer, produk, atau data lain.</li>
<li><strong>Konfigurasi Auto Reply</strong> — Setup keyword rules khusus untuk bisnis %s di %s.</li>
<li><strong>Mulai Broadcast</strong> — Kirim pesan ke customer %s dengan interval aman.</li>
<li><strong>Pantau & Optimasi</strong> — Cek dashboard analytics, evaluasi, perbaiki.</li>
</ol>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Tips Khusus untuk %s</h2>
<ul style="color:#5e6b7e;line-height:2;padding-left:20px">
<li>Gunakan <strong>template + spintax</strong> untuk variasi pesan otomatis</li>
<li>Setup <strong>jam kerja</strong> — auto-reply beda saat jam kerja & di luar jam</li>
<li>Aktifkan <strong>CSAT survey</strong> — ukur kepuasan customer otomatis</li>
<li>Manfaatkan <strong>drip campaign</strong> — nurture leads sampai closing</li>
</ul>
</div>`, humanize(p.Keyword), p.Industry, p.City, p.Industry, p.City, AppName, AppName, p.Industry, p.Industry, p.City, p.Industry, p.Industry))
}

func compoundCityContent(p parsedSlug) template.HTML {
	return template.HTML(fmt.Sprintf(`<div style="max-width:800px;margin:0 auto">
<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">%s di %s</h2>
<p style="color:#5e6b7e;line-height:1.8;margin-bottom:24px">Bisnis di <strong>%s</strong> berkembang pesat. Untuk tetap kompetitif, Anda butuh tools WhatsApp marketing yang <strong>handal, terjangkau, dan fleksibel</strong>. %s adalah solusi self-hosted — install di server sendiri, kendali penuh, tanpa biaya bulanan.</p>

<h2 style="font-size:1.5rem;font-weight:700;color:#152e4d;margin:32px 0 16px">Kenapa Bisnis %s Pilih Self-Hosted?</h2>
<div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:32px">
<div style="background:#f0fdf4;border:1px solid #bbf7d0;border-radius:12px;padding:20px"><strong style="color:#166534">💡 Kendali Penuh</strong><p style="color:#15803d;margin:8px 0 0;font-size:.9rem">Server sendiri, aturan sendiri. Data tidak dikirim ke pihak ketiga.</p></div>
<div style="background:#eff6ff;border:1px solid #bfdbfe;border-radius:12px;padding:20px"><strong style="color:#1e40af">💰 Hemat</strong><p style="color:#2563eb;margin:8px 0 0;font-size:.9rem">Bayar sekali. Dibanding SaaS, balik modal dalam 3-6 bulan.</p></div>
<div style="background:#fefce8;border:1px solid #fef08a;border-radius:12px;padding:20px"><strong style="color:#854d0e">🎯 Custom</strong><p style="color:#a16207;margin:8px 0 0;font-size:.9rem">Full source code — tambah fitur sesuai kebutuhan bisnis %s.</p></div>
</div>
</div>`, humanize(p.Keyword), p.City, p.City, AppName, p.City, p.City))
}
