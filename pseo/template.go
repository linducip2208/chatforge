package pseo

const pseoTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}} | {{.SiteName}}</title>
<meta name="description" content="{{.Description}}">
<meta name="robots" content="index, follow">
<link rel="canonical" href="{{.Canonical}}">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta property="og:type" content="website">
<meta property="og:site_name" content="{{.SiteName}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<link rel="icon" href="/assets/theme/default-favicon.png">
<style>
*,::before,::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:'Inter',-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;color:#152e4d;background:#fff;line-height:1.6}
a{color:#4F46E5;text-decoration:none}
a:hover{text-decoration:underline}
.container{max-width:1100px;margin:0 auto;padding:0 20px}
/* Header */
.site-header{background:#fff;border-bottom:1px solid #e5e7eb;position:sticky;top:0;z-index:100;backdrop-filter:blur(12px)}
.header-inner{display:flex;align-items:center;justify-content:space-between;padding:14px 0;flex-wrap:wrap;gap:12px}
.site-logo{display:flex;align-items:center;gap:10px;font-weight:800;font-size:1.3rem;color:#152e4d;text-decoration:none}
.site-logo span{color:#4F46E5}
.header-cta{display:flex;gap:10px;align-items:center;flex-wrap:wrap}
.header-cta a{padding:10px 20px;border-radius:10px;font-weight:600;font-size:.9rem;text-decoration:none}
.btn-demo{background:#4F46E5;color:#fff}
.btn-demo:hover{background:#4338CA;text-decoration:none}
.btn-wa{background:#25D366;color:#fff}
.btn-wa:hover{background:#1da851;text-decoration:none}
/* Breadcrumb */
.breadcrumb{padding:16px 0;color:#6b7280;font-size:.85rem}
.breadcrumb a{color:#6b7280}
.breadcrumb span{color:#152e4d;font-weight:500}
/* Hero */
.pseo-hero{background:linear-gradient(135deg,#0f1f33,#152e4d 50%,#1a3a5c);color:#fff;padding:60px 0;text-align:center}
.pseo-hero h1{font-size:2.2rem;font-weight:800;line-height:1.2;margin-bottom:12px}
.pseo-hero p{font-size:1.1rem;opacity:.85;max-width:700px;margin:0 auto;line-height:1.7}
/* Content */
.pseo-content{padding:20px 0 60px}
.pseo-content h2{font-size:1.4rem;font-weight:700;color:#152e4d;margin:28px 0 14px}
.pseo-content h3{font-size:1.15rem;font-weight:600;color:#152e4d;margin:20px 0 10px}
.pseo-content p{color:#4b5563;line-height:1.8;margin-bottom:16px}
.pseo-content ul,.pseo-content ol{padding-left:24px;margin-bottom:16px}
.pseo-content li{color:#4b5563;line-height:1.8}
/* Feature Grid */
.feat-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(250px,1fr));gap:16px;margin:24px 0}
.feat-card{background:#f9fafb;border:1px solid #e5e7eb;border-radius:12px;padding:20px}
.feat-card h3{font-size:1rem;font-weight:700;color:#4F46E5;margin:0 0 6px}
.feat-card p{color:#6b7280;margin:0;font-size:.9rem;line-height:1.6}
/* List Items */
.list-items{margin:20px 0}
.list-item{display:flex;align-items:flex-start;gap:12px;padding:14px 16px;background:#f9fafb;border-radius:10px;margin-bottom:10px;transition:all .2s}
.list-item:hover{background:#f0efff}
.list-item-num{background:#4F46E5;color:#fff;font-weight:700;font-size:.85rem;min-width:28px;height:28px;border-radius:50%;display:flex;align-items:center;justify-content:center;flex-shrink:0}
.list-item-text{flex:1}
.list-item-text strong{display:block;color:#152e4d;margin-bottom:2px}
.list-item-text span{color:#6b7280;font-size:.9rem}
/* Table */
.compare-table{overflow-x:auto;margin:20px 0;border-radius:12px;border:1px solid #e5e7eb}
.compare-table table{width:100%;border-collapse:collapse;font-size:.9rem}
.compare-table th{background:#f9fafb;padding:12px 16px;text-align:left;font-weight:700;color:#152e4d;border-bottom:2px solid #e5e7eb}
.compare-table td{padding:10px 16px;border-bottom:1px solid #f3f4f6;color:#4b5563}
.compare-table .highlight{background:#f5f3ff;font-weight:600}
/* CTA */
.cta-section{background:linear-gradient(135deg,#4F46E5,#7C3AED);color:#fff;border-radius:16px;padding:40px;text-align:center;margin:40px 0}
.cta-section h2{font-size:1.6rem;font-weight:800;color:#fff;margin:0 0 8px}
.cta-section p{font-size:1.05rem;opacity:.9;color:#fff;margin:0 0 24px}
.cta-buttons{display:flex;gap:12px;justify-content:center;flex-wrap:wrap}
.cta-btn{display:inline-flex;align-items:center;gap:8px;padding:14px 28px;border-radius:12px;font-weight:700;font-size:1rem;text-decoration:none}
.cta-btn-primary{background:#25D366;color:#fff}
.cta-btn-primary:hover{background:#1da851;text-decoration:none}
.cta-btn-secondary{background:rgba(255,255,255,.2);color:#fff}
.cta-btn-secondary:hover{background:rgba(255,255,255,.3);text-decoration:none}
/* Footer */
.site-footer{background:#0f1f33;color:rgba(255,255,255,.6);padding:40px 0;margin-top:40px}
.footer-inner{display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:16px}
.footer-inner a{color:rgba(255,255,255,.8)}
.footer-inner a:hover{color:#fff}
/* Responsive */
@media(max-width:768px){
  .pseo-hero h1{font-size:1.5rem}
  .pseo-hero p{font-size:1rem}
  .pseo-hero{padding:40px 0}
  .feat-grid{grid-template-columns:1fr}
  .cta-section{padding:28px 20px}
  .cta-section h2{font-size:1.3rem}
  .header-inner{flex-direction:column;gap:8px}
}
@media(max-width:480px){
  .pseo-hero h1{font-size:1.3rem}
  .cta-buttons{flex-direction:column}
  .cta-btn{justify-content:center}
}
/* FAQ style for compare pages */
.faq-item{background:#f9fafb;border-radius:10px;padding:18px;margin-bottom:10px}
.faq-item strong{display:block;color:#152e4d;margin-bottom:6px}
.faq-item p{color:#4b5563;margin:0;font-size:.9rem;line-height:1.7}
</style>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.bunny.net/css?family=inter:400,500,600,700,800" rel="stylesheet">
</head>
<body>

<header class="site-header">
<div class="container header-inner">
<a href="/" class="site-logo">
<svg width="32" height="32" viewBox="0 0 32 32"><rect width="32" height="32" rx="8" fill="#4F46E5"/><text x="16" y="22" text-anchor="middle" fill="#fff" font-size="18" font-weight="800">C</text></svg>
{{.SiteName}}<span>.</span>
</a>
<div class="header-cta">
<a href="/docs" class="btn-demo">📖 Dokumentasi</a>
<a href="https://wa.me/{{.WaNumber}}" class="btn-wa" target="_blank" rel="noopener">💬 WhatsApp</a>
</div>
</div>
</header>

<div class="container">
{{if .Breadcrumbs}}
<nav class="breadcrumb" aria-label="Breadcrumb">
{{range $i, $b := .Breadcrumbs}}
{{if $i}} › {{end}}
{{if $b.URL}}<a href="{{$b.URL}}">{{$b.Name}}</a>{{else}}<span>{{$b.Name}}</span>{{end}}
{{end}}
</nav>
{{end}}
</div>

<section class="pseo-hero">
<div class="container">
<h1>{{.Title}}</h1>
<p>{{.Description}}</p>
</div>
</section>

<div class="container pseo-content">
{{if .IsList}}
<div class="list-items">
{{range $i, $item := .Items}}
<div class="list-item">
<div class="list-item-num">{{inc $i}}</div>
<div class="list-item-text">
{{if $item.URL}}<a href="{{$item.URL}}"><strong>{{$item.Title}}</strong></a>{{else}}<strong>{{$item.Title}}</strong>{{end}}
{{if $item.Description}}<span>{{$item.Description}}</span>{{end}}
</div>
</div>
{{end}}
</div>
{{end}}

{{.Content}}

<div class="cta-section">
<h2>Punya Pertanyaan?</h2>
<p>Konsultasikan kebutuhan WhatsApp marketing Anda dengan kami. Gratis!</p>
<div class="cta-buttons">
<a href="https://wa.me/{{.WaNumber}}?text=Halo%20saya%20tertarik%20dengan%20{{.SiteName}}%20untuk%20WhatsApp%20Marketing" target="_blank" rel="noopener" class="cta-btn cta-btn-primary">💬 Chat via WhatsApp</a>
<a href="/docs" class="cta-btn cta-btn-secondary">📖 Lihat Dokumentasi</a>
</div>
</div>

</div>

<footer class="site-footer">
<div class="container footer-inner">
<div>© {{.Year}} {{.SiteName}}. Solusi WhatsApp Marketing Self-Hosted Indonesia.</div>
<div><a href="/docs">Dokumentasi</a> · <a href="/best-whatsapp-marketing-tools">Tools Terbaik</a> · <a href="https://wa.me/{{.WaNumber}}" target="_blank" rel="noopener">WhatsApp {{.WaNumber}}</a></div>
</div>
</footer>

</body>
</html>`
