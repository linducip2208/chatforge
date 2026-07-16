package main

const templates = `
{{define "layout"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
<title>{{.Title}} &middot; ChatGo</title>
<link rel="icon" href="/assets/theme/default-favicon.png">
<link rel="stylesheet" href="/assets/_assets/css/libs/line-awesome.min.css">
<link rel="stylesheet" href="/assets/_assets/css/libs/flag-icon.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/fonts/feather/feather.css">
<link rel="stylesheet" href="/assets/dashboard/css/libs/bootstrap.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/style.min.css">
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
<style>
  .navbar-vibrant{background:#0B1220}
  .navbar-vertical .navbar-heading{color:#5a6780;font-size:0.68rem;font-weight:700;letter-spacing:0.08em;padding:12px 24px 6px;text-transform:uppercase}
  .navbar-vertical .nav-link{color:#8895b7;font-size:0.85rem;padding:8px 24px;margin:1px 8px;border-radius:8px;transition:all .15s}
  .navbar-vertical .nav-link:hover{background:rgba(255,255,255,.04);color:#fff}
  .navbar-vertical .nav-link.active{background:rgba(44,123,229,.15);color:#2c7be5;font-weight:600}
  .navbar-vertical .nav-link i{font-size:1.1rem;width:22px;text-align:center;opacity:.7}
  .navbar-vertical .nav-link.active i{opacity:1}
  .navbar-divider{border-color:rgba(255,255,255,.06)}
  #qrimg{width:260px;height:260px;background:#fff;border-radius:12px;padding:8px}
  .status-dot{height:9px;width:9px;border-radius:50%;display:inline-block;margin-right:5px}
  .badge-soft-success{background:rgba(0,217,126,.1);color:#00d97e}
  .badge-soft-danger{background:rgba(230,55,87,.1);color:#e63757}
  .badge-soft-warning{background:rgba(246,195,67,.15);color:#f6c343}
  .badge-soft-secondary{background:rgba(110,120,140,.12);color:#6e788c}
  .lang-flag{width:20px;height:15px;border-radius:2px;margin-right:6px;object-fit:cover}
  .msg-trunc{cursor:pointer;color:#2c7be5}
  .msg-trunc:hover{text-decoration:underline}
  .pagination{display:flex;gap:4px;padding:12px 0;flex-wrap:wrap}
  .pagination a,.pagination span{padding:6px 12px;border-radius:6px;border:1px solid #ddd;text-decoration:none;color:#152e4d;font-size:13px}
  .pagination .active{background:#2c7be5;color:#fff;border-color:#2c7be5}
   .auth-page{min-height:100vh;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#152e4d 0%,#1a3a5c 50%,#0f1f33 100%)}
   .auth-card{width:100%;max-width:420px;border-radius:14px;box-shadow:0 20px 60px rgba(0,0,0,.3)}
   .auth-split{min-height:100vh;display:flex;flex-wrap:wrap}
   .auth-left{flex:1;min-width:300px;background:linear-gradient(160deg,#0f1f33,#152e4d 40%,#1a3a5c);display:flex;flex-direction:column;justify-content:space-between;padding:48px 40px;position:relative;overflow:hidden}
   .auth-left::before{content:'';position:absolute;top:-80px;right:-80px;width:300px;height:300px;border-radius:50%;background:rgba(44,123,229,.08)}
   .auth-left::after{content:'';position:absolute;bottom:-60px;left:-60px;width:250px;height:250px;border-radius:50%;background:rgba(44,123,229,.05)}
   .auth-left-content{position:relative;z-index:1}
   .auth-right{flex:1;min-width:320px;background:#fff;display:flex;align-items:center;justify-content:center;padding:40px}
   .auth-logo{display:inline-flex;align-items:center;gap:10px;text-decoration:none}
   .auth-logo-text{font-size:28px;font-weight:800;color:#fff;letter-spacing:-.5px}
   .auth-logo-text span{color:#2c7be5}
   .auth-hero h2{font-size:2rem;font-weight:700;color:#fff;line-height:1.2;margin-bottom:12px}
   .auth-hero p{color:rgba(255,255,255,.65);font-size:0.95rem;line-height:1.6;margin-bottom:32px;max-width:360px}
   .auth-features{display:flex;flex-direction:column;gap:12px;max-width:320px}
   .auth-feat{display:flex;align-items:center;gap:12px;padding:12px 16px;background:rgba(255,255,255,.06);border-radius:10px;backdrop-filter:blur(8px)}
   .auth-feat i{font-size:20px;color:#2c7be5;width:24px;text-align:center}
   .auth-feat span{color:rgba(255,255,255,.8);font-size:0.85rem;font-weight:500}
   .auth-copy{color:rgba(255,255,255,.3);font-size:0.75rem}
   .auth-form-wrap{width:100%;max-width:400px}
   .auth-form-wrap h3{font-size:1.5rem;font-weight:700;color:#152e4d;margin-bottom:4px}
   .auth-form-wrap .sub{color:#6e788c;font-size:0.85rem;margin-bottom:24px}
   .auth-form-wrap .sub a{color:#2c7be5;font-weight:600;text-decoration:none}
   .auth-form-wrap .form-control{border-radius:10px;padding:10px 14px;border:1.5px solid #e0e4e9;font-size:0.9rem}
   .auth-form-wrap .form-control:focus{border-color:#2c7be5;box-shadow:0 0 0 3px rgba(44,123,229,.12)}
   .auth-form-wrap .btn{width:100%;border-radius:10px;padding:12px;font-weight:600;font-size:0.95rem;background:linear-gradient(135deg,#2c7be5,#1a5bbf);border:none;color:#fff}
   .auth-form-wrap .btn:hover{transform:translateY(-1px);box-shadow:0 8px 24px rgba(44,123,229,.3)}
   .auth-divider{display:flex;align-items:center;margin:20px 0;color:#aab0b9;font-size:0.8rem}
   .auth-divider::before,.auth-divider::after{content:'';flex:1;height:1px;background:#e0e4e9}
   .auth-divider span{padding:0 12px}
   .demo-box{background:#f7f8fa;border:1px solid #e0e4e9;border-radius:10px;padding:14px 16px}
   .demo-box .demo-title{font-weight:700;color:#152e4d;font-size:0.8rem;margin-bottom:8px}
   .demo-box .demo-row{font-size:0.75rem;color:#5e6b7e;padding:3px 0;font-family:'Courier New',monospace}
   .demo-box .demo-row strong{color:#152e4d}
   @media(max-width:767px){
     .auth-left{min-width:100%;padding:32px 24px;min-height:auto}
     .auth-left::before,.auth-left::after{display:none}
     .auth-hero h2{font-size:1.4rem}
     .auth-hero p{font-size:.85rem;margin-bottom:16px}
     .auth-features{display:none}
     .auth-right{min-width:100%;padding:24px}
     .auth-form-wrap{max-width:100%}
   }
</style>
</head>
<body>

{{template "sidebar" .}}
<div class="main-content">
  <nav class="navbar navbar-expand-md navbar-light d-none d-md-flex" id="topbar">
    <div class="container-fluid">
      <div class="me-4">
        <a class="btn btn-md btn-primary mb-1 lift" href="/wa"><i class="la la-whatsapp la-lg me-1"></i> {{T "nav_whatsapp"}}</a>
        <a class="btn btn-md btn-primary mb-1 lift" href="/send"><i class="la la-paper-plane la-lg me-1"></i> {{T "nav_send"}}</a>
      </div>
      <div class="navbar-user d-flex align-items-center">
        <div class="dropdown me-3">
          <a href="#" class="dropdown-toggle text-muted" role="button" data-bs-toggle="dropdown" style="text-decoration:none">
            <span class="flag-icon flag-icon-{{.LangFlag}} lang-flag"></span>{{.LangName}}
          </a>
          <div class="dropdown-menu dropdown-menu-end">
            {{range .Languages}}
            <a class="dropdown-item" href="/lang/{{.Code}}"><span class="flag-icon flag-icon-{{.Flag}} lang-flag"></span>{{.Name}}</a>
            {{end}}
          </div>
        </div>
        {{if eq .Status "connected"}}
          <span class="badge badge-soft-success"><span class="status-dot" style="background:#00d97e"></span> {{T "status_connected"}} +{{.Phone}}</span>
        {{else if eq .Status "qr"}}
          <span class="badge badge-soft-warning"><span class="status-dot" style="background:#f6c343"></span> {{T "status_scanqr"}}</span>
        {{else}}
          <span class="badge badge-soft-danger"><span class="status-dot" style="background:#e63757"></span> {{T "status_disconnected"}}</span>
        {{end}}
      </div>
    </div>
  </nav>

  <div class="header">
    <div class="container-fluid">
      <div class="header-body">
        <div class="row align-items-end">
          <div class="col">
            <h6 class="header-pretitle">{{.Pretitle}}</h6>
            <h1 class="header-title"><i class="la {{.Icon}} la-lg me-1"></i> {{.Heading}}</h1>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="container-fluid">
    {{if .Flash}}<div class="alert alert-success" role="alert">{{.Flash}}</div>{{end}}
    {{template "content" .}}
  </div>
</div>

<script src="/assets/_assets/js/libs/jquery.min.js"></script>
<script src="/assets/dashboard/assets/js/libs/bootstrap.min.js"></script>
<script>
var lastStatus="{{.Status}}";
setInterval(async()=>{try{const r=await fetch("/status");const s=await r.json();
 if(s.status!==lastStatus){location.reload();return;}
 if(s.status!=="connected"&&s.qr){var i=document.getElementById("qrimg");if(i)i.src="/qr.png?t="+Date.now();}
}catch(e){}},3000);
// AI toggle: mirror wabot
var cb=document.getElementById("useAiCheck");
if(cb){cb.addEventListener("change",function(){
  var on=this.checked;
  document.getElementById("aiKeyGroup").style.display=on?"block":"none";
  var f=this.closest("form");
  if(f.querySelector('input[name="keyword"]')) f.querySelector('input[name="keyword"]').required=!on;
  if(f.querySelector('textarea[name="reply"]')) f.querySelector('textarea[name="reply"]').required=!on;
  if(f.querySelector('select[name="ai_key_id"]')) f.querySelector('select[name="ai_key_id"]').required=on;
})}
// match type: AI Reply → hide keyword, show {{T "ar_faq_tab"}}
function onMatchTypeChange(v){
  document.getElementById("keywordGroup").style.display=v==="ai"?"none":"block";
  document.getElementById("{{T "ar_faq_tab"}}Group").style.display=v==="ai"?"block":"none";
  var kw=document.querySelector('input[name="keyword"]');
  if(kw) kw.required=v!=="ai";
}
// tab switcher
document.querySelectorAll('.nav-tabs .nav-link').forEach(function(t){t.addEventListener('click',function(e){e.preventDefault();document.querySelectorAll('.nav-tabs .nav-link').forEach(function(x){x.classList.remove('active')});document.querySelectorAll('.tab-pane').forEach(function(x){x.classList.remove('show','active')});this.classList.add('active');var el=document.querySelector(this.getAttribute('href'));if(el){el.classList.add('show','active')}})});
// truncate messages: show first 20 chars
document.querySelectorAll(".msg-full").forEach(function(el){
  var text=el.textContent.trim();
  if(text.length>20){ el.setAttribute("data-full",text); el.textContent=text.substring(0,20)+"..."; el.style.cursor="pointer"; el.title="Klik untuk lihat selengkapnya";
    el.addEventListener("click",function(){
      if(this.getAttribute("data-full")===this.textContent){ this.textContent=text.substring(0,20)+"..."; return; }
      this.textContent=this.getAttribute("data-full");
    });
  }
});
</script>
</body>
</html>{{end}}

{{define "sidebar"}}
<nav class="navbar navbar-vertical fixed-left navbar-expand-md navbar-dark navbar-vibrant" id="sidebar">
  <div class="container-fluid">
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#sidebarCollapse">
      <span class="navbar-toggler-icon"></span>
    </button>
    <a class="navbar-brand" href="/">
      <img src="/assets/theme/default-logo-light.png" class="navbar-brand-img mx-auto" alt="ChatGo" onerror="this.outerHTML='<span style=&quot;color:#fff;font-weight:800;font-size:20px&quot;>chat<span style=&quot;color:#2c7be5&quot;>go</span></span>'">
    </a>
    <div class="collapse navbar-collapse" id="sidebarCollapse">
      <h6 class="navbar-heading">{{T "nav_overview"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "home"}}active{{end}}" href="/"><i class="la la-chart-bar la-lg"></i> {{T "nav_dashboard"}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_whatsapp"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "wa"}}active{{end}}" href="/wa"><i class="la la-whatsapp la-lg"></i> {{T "nav_account_qr"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "send"}}active{{end}}" href="/send"><i class="la la-paper-plane la-lg"></i> {{T "nav_send"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "broadcast"}}active{{end}}" href="/broadcast"><i class="la la-bullhorn la-lg"></i> {{T "nav_broadcast"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "scheduled"}}active{{end}}" href="/scheduled"><i class="la la-clock la-lg"></i> {{T "nav_scheduled"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "inbox"}}active{{end}}" href="/inbox"><i class="la la-comments la-lg"></i> Inbox</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "sent"}}active{{end}}" href="/sent"><i class="la la-telegram la-lg"></i> {{T "nav_sent"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "received"}}active{{end}}" href="/received"><i class="la la-comment la-lg"></i> {{T "nav_received"}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      {{if eq .Role "admin"}}
      <h6 class="navbar-heading">{{T "nav_hosts"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "hosts_whatsapp"}}active{{end}}" href="/hosts/whatsapp"><i class="la la-whatsapp la-lg"></i> {{T "nav_hosts_whatsapp"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "hosts_android"}}active{{end}}" href="/hosts/android"><i class="la la-mobile la-lg"></i> {{T "nav_hosts_android"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "ussd"}}active{{end}}" href="/ussd"><i class="la la-satellite-dish la-lg"></i> {{T "nav_ussd"}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      {{end}}
      <h6 class="navbar-heading">{{T "nav_contacts"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "contacts"}}active{{end}}" href="/contacts"><i class="la la-address-book la-lg"></i> {{T "nav_contacts_saved"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "groups"}}active{{end}}" href="/contacts/groups"><i class="la la-list la-lg"></i> {{T "nav_contacts_groups"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "unsub"}}active{{end}}" href="/contacts/unsub"><i class="la la-unlink la-lg"></i> {{T "nav_contacts_unsub"}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_tools"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "autoreply"}}active{{end}}" href="/autoreply"><i class="la la-robot la-lg"></i> {{T "nav_autoreply"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "templates"}}active{{end}}" href="/templates"><i class="la la-file-alt la-lg"></i> {{T "nav_templates"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "apikeys"}}active{{end}}" href="/apikeys"><i class="la la-key la-lg"></i> {{T "nav_apikeys"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "webhooks"}}active{{end}}" href="/webhooks"><i class="la la-code-branch la-lg"></i> {{T "nav_webhooks"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "logger"}}active{{end}}" href="/logger"><i class="la la-clipboard-list la-lg"></i> {{T "nav_logger"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "settings"}}active{{end}}" href="/settings"><i class="la la-cog la-lg"></i> {{T "nav_settings"}}</a></li>
      </ul>
      {{if eq .Role "admin"}}
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_admin"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin"}}active{{end}}" href="/admin"><i class="la la-chart-bar la-lg"></i> {{T "nav_overview"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_users"}}active{{end}}" href="/admin/users"><i class="la la-users la-lg"></i> {{T "adm_users"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_roles"}}active{{end}}" href="/admin/roles"><i class="la la-user-shield la-lg"></i> {{T "adm_roles"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_packages"}}active{{end}}" href="/admin/packages"><i class="la la-box la-lg"></i> {{T "adm_packages"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_vouchers"}}active{{end}}" href="/admin/vouchers"><i class="la la-ticket-alt la-lg"></i> {{T "adm_vouchers"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_subscriptions"}}active{{end}}" href="/admin/subscriptions"><i class="la la-star la-lg"></i> {{T "adm_subscriptions"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_transactions"}}active{{end}}" href="/admin/transactions"><i class="la la-money-bill la-lg"></i> {{T "adm_transactions"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_payouts"}}active{{end}}" href="/admin/payouts"><i class="la la-hand-holding-usd la-lg"></i> {{T "adm_payouts"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_pages"}}active{{end}}" href="/admin/pages"><i class="la la-file la-lg"></i> {{T "adm_pages"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_marketing"}}active{{end}}" href="/admin/marketing"><i class="la la-bullhorn la-lg"></i> {{T "adm_marketing"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_languages"}}active{{end}}" href="/admin/languages"><i class="la la-language la-lg"></i> {{T "adm_languages"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_waservers"}}active{{end}}" href="/admin/waservers"><i class="la la-server la-lg"></i> {{T "adm_waservers"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_gateways"}}active{{end}}" href="/admin/gateways"><i class="la la-code la-lg"></i> {{T "adm_gateways"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_shorteners"}}active{{end}}" href="/admin/shorteners"><i class="la la-link la-lg"></i> {{T "adm_shorteners"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_plugins"}}active{{end}}" href="/admin/plugins"><i class="la la-puzzle-piece la-lg"></i> {{T "adm_plugins"}}</a></li>
      </ul>
      {{end}}
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_docs"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "docs"}}active{{end}}" href="/docs"><i class="la la-book la-lg"></i> {{T "nav_docs"}}</a></li>
      </ul>
    </div>
  </div>
</nav>{{end}}

{{define "home"}}{{template "layout" .}}{{end}}
{{define "content"}}
{{if eq .Page "home"}}
  <div class="row">
    <div class="col-12 col-lg-6 col-xl-3">
      <div class="card"><div class="card-body"><div class="row align-items-center">
        <div class="col"><h6 class="text-uppercase text-muted mb-2">{{T "dash_status"}}</h6><span class="h2 mb-0">{{if eq .Status "connected"}}{{T "dash_active"}}{{else}}{{T "dash_inactive"}}{{end}}</span></div>
        <div class="col-auto"><span class="h2 la la-whatsapp la-lg text-muted mb-0"></span></div>
      </div></div></div>
    </div>
    <div class="col-12 col-lg-6 col-xl-3">
      <div class="card"><div class="card-body"><div class="row align-items-center">
        <div class="col"><h6 class="text-uppercase text-muted mb-2">{{T "dash_connected_number"}}</h6><span class="h2 mb-0">{{if .Phone}}+{{.Phone}}{{else}}-{{end}}</span></div>
        <div class="col-auto"><span class="h2 la la-mobile la-lg text-muted mb-0"></span></div>
      </div></div></div>
    </div>
    <div class="col-12 col-lg-6 col-xl-3">
      <div class="card"><div class="card-body"><div class="row align-items-center">
        <div class="col"><h6 class="text-uppercase text-muted mb-2">{{T "dash_total_out"}}</h6><span class="h2 mb-0">{{.CountSent}}</span></div>
        <div class="col-auto"><span class="h2 la la-telegram la-lg text-muted mb-0"></span></div>
      </div></div></div>
    </div>
    <div class="col-12 col-lg-6 col-xl-3">
      <div class="card"><div class="card-body"><div class="row align-items-center">
        <div class="col"><h6 class="text-uppercase text-muted mb-2">Total Users</h6><span class="h2 mb-0">{{.TotalUsers}}</span></div>
        <div class="col-auto"><span class="h2 la la-users la-lg text-muted mb-0"></span></div>
      </div></div></div>
    </div>
  </div>
  <div class="row"><div class="col-6 col-xl-3"><div class="card"><div class="card-body"><div class="row align-items-center"><div class="col"><h6 class="text-uppercase text-muted mb-2">Active WA</h6><span class="h2 mb-0">{{.ActiveAccounts}}</span></div><div class="col-auto"><span class="h2 la la-whatsapp la-lg text-success mb-0"></span></div></div></div></div></div>
  <div class="col-6 col-xl-3"><div class="card"><div class="card-body"><div class="row align-items-center"><div class="col"><h6 class="text-uppercase text-muted mb-2">Campaigns</h6><span class="h2 mb-0">{{.RunningCampaigns}}</span></div><div class="col-auto"><span class="h2 la la-bullhorn la-lg text-warning mb-0"></span></div></div></div></div></div>
  <div class="col-12 col-xl-6"><div class="card"><div class="card-body p-2"><div class="d-flex gap-2 flex-wrap align-items-center">{{range .ActiveAccountList}}{{if .Phone}}<a href="/send?to=+{{.Phone}}" class="badge bg-success bg-opacity-10 text-success text-decoration-none small py-2 px-3">+{{.Phone}} ✉️</a>{{end}}{{else}}<span class="text-muted small">No active accounts</span>{{end}}</div></div></div></div></div>
  <div class="row"><div class="col-12"><div class="card"><div class="card-header"><h4 class="card-header-title">Message Activity (7 days)</h4></div><div class="card-body"><canvas id="msgChart" height="80"></canvas></div></div></div></div>
  <script>new Chart(document.getElementById('msgChart'),{type:'bar',data:{labels:[{{.ChartLabels}}],datasets:[{label:'Sent',data:[{{.ChartSent}}],backgroundColor:'#4F46E5',borderRadius:4},{label:'Received',data:[{{.ChartReceived}}],backgroundColor:'#10B981',borderRadius:4}]},options:{responsive:true,scales:{y:{beginAtZero:true}}}})</script>
  <div class="row">
    <div class="col-12 col-xl-6">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_in"}}</h4><a href="/received" class="btn btn-sm btn-white">{{T "btn_all"}}</a></div>
      <div class="table-responsive"><table class="table table-sm table-nowrap card-table"><thead><tr><th>{{T "col_from"}}</th><th>{{T "col_name"}}</th><th>{{T "col_message"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
        {{range .Received}}<tr><td>{{.Phone}}</td><td>{{.Name}}</td><td>{{.Message}}</td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div>
    </div>
    <div class="col-12 col-xl-6">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_out"}}</h4><a href="/sent" class="btn btn-sm btn-white">{{T "btn_all"}}</a></div>
      <div class="table-responsive"><table class="table table-sm table-nowrap card-table"><thead><tr><th>{{T "col_to"}}</th><th>{{T "col_message"}}</th><th>{{T "col_status"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
        {{range .Sent}}<tr><td>{{.Phone}}</td><td>{{.Message}}</td><td><span class="badge badge-soft-success">{{.Status}}</span></td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div>
    </div>
  </div>
{{end}}

{{if eq .Page "wa"}}
  <div class="row">
    <div class="col-12 col-lg-6">
      <div class="card">
        <div class="card-header"><h4 class="card-header-title"><i class="la la-whatsapp text-success me-1"></i> {{T "nav_account_qr"}}</h4>
          <form method="post" action="/wa/add" style="display:inline"><button class="btn btn-sm btn-primary lift" {{if ge (len .Accounts) .AccountLimit}}disabled{{end}}><i class="la la-plus me-1"></i> {{T "wa_add_account"}}</button></form>
          <span class="text-muted small ms-2">{{len .Accounts}} / {{.AccountLimit}}</span>
        </div>
        <div class="card-body">
          {{if .ScanAccount}}
            <div class="text-center">
              <img id="qrimg" src="/qr.png?id={{.ScanAccount}}" alt="QR Code" onerror="this.style.display='none'">
              <p class="text-muted mt-2">{{T "wa_new_qr"}}</p>
            </div>
          {{end}}
          {{range .Accounts}}
          <div class="d-flex align-items-center justify-content-between border rounded p-3 mb-2">
            <div>{{if .Phone}}<strong>+{{.Phone}}</strong>{{else}}<span class="text-muted">{{T "wa_pairing"}}</span>{{end}}<br>
              {{if eq .Status "connected"}}<span class="badge badge-soft-success"><i class="la la-check-circle me-1"></i>Connected</span>
              {{else if eq .Status "qr"}}<span class="badge badge-soft-warning"><i class="la la-qrcode me-1"></i>Scan QR</span>
              {{else}}<span class="badge badge-soft-danger"><i class="la la-times-circle me-1"></i>Disconnected</span>{{end}}
            </div>
            <form method="post" action="/wa/logout" onsubmit="return confirm('Logout nomor ini?')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger lift"><i class="la la-sign-out"></i></button></form>
          </div>
          {{else}}
          <div class="text-center py-4"><span class="h1 la la-whatsapp text-muted"></span><p class="text-muted mt-3">{{T "wa_no_accounts"}}</p></div>
          {{end}}
        </div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "send"}}
  <div class="row justify-content-center">
    <div class="col-12 col-lg-8">
      <div class="card">
        <div class="card-header"><h4 class="card-header-title">{{T "send_title"}}</h4></div>
        <div class="card-body">
          <form method="post" action="/send"><input type="hidden" name="is_text" value="1">
            <div class="form-group"><label>{{T "send_to"}}</label><input name="phone" class="form-control" placeholder="628123456789" value="{{.SendTo}}" required><small class="form-text text-muted">{{T "send_to_hint"}}</small></div>
            <div class="form-group"><label>{{T "send_message"}}</label><textarea name="message" class="form-control" rows="4" placeholder="{{T "send_message_ph"}}" required></textarea></div>
            <button class="btn btn-primary lift" {{if ne .Status "connected"}}disabled{{end}}><i class="la la-paper-plane me-1"></i> {{T "send_btn"}}</button>
            {{if ne .Status "connected"}}<span class="text-muted ms-2">{{T "send_connect_first"}}</span>{{end}}
          </form>
        </div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "sent"}}
  <div class="card">
    <div class="card-header"><h4 class="card-header-title">{{T "sent_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm table-nowrap card-table"><thead><tr><th>{{T "col_no"}}</th><th>{{T "col_to"}}</th><th>{{T "col_message"}}</th><th>{{T "col_status"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Sent}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td><span class="msg-full">{{.Message}}</span></td><td><span class="badge badge-soft-success">{{.Status}}</span></td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">{{T "sent_empty"}}</td></tr>{{end}}
    </tbody></table></div>
    <div class="pagination px-3 pb-2">{{range .SentPages}}<a href="/sent?page={{.}}" class="{{if eq . $.SentPage}}active{{end}}">{{.}}</a>{{end}}</div>
  </div>
{{end}}

{{if eq .Page "received"}}
  <div class="card">
    <div class="card-header"><h4 class="card-header-title">{{T "received_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm table-nowrap card-table"><thead><tr><th>{{T "col_no"}}</th><th>{{T "col_from"}}</th><th>{{T "col_name"}}</th><th>{{T "col_message"}}</th><th>{{T "col_type"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Received}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Name}}</td><td><span class="msg-full">{{.Message}}</span></td><td>{{if .IsGroup}}<span class="badge badge-soft-warning">{{T "type_group"}}</span>{{else}}<span class="badge badge-soft-success">{{T "type_private"}}</span>{{end}}</td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">{{T "received_empty"}}</td></tr>{{end}}
    </tbody></table></div>
    <div class="pagination px-3 pb-2">{{range .ReceivedPages}}<a href="/received?page={{.}}" class="{{if eq . $.ReceivedPage}}active{{end}}">{{.}}</a>{{end}}</div>
  </div>
{{end}}

{{if eq .Page "autoreply"}}
<div class="row">
<div class="col-12 col-lg-5">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ar_add_title"}}</h4></div>
<div class="card-body"><form method="post" action="/autoreply/add">
<div class="form-group"><label>{{T "ar_matchtype"}}</label><select name="match" class="form-control" onchange="onMatchTypeChange(this.value)"><option value="contains">{{T "ar_contains"}}</option><option value="exact">{{T "ar_exact"}}</option><option value="starts_with">{{T "ar_starts"}}</option><option value="ai">{{T "ar_ai_type"}}</option></select></div>
<div id="keywordGroup"><div class="form-group"><label>{{T "ar_keyword"}}</label><input name="keyword" class="form-control" placeholder="halo, hi, menu"></div></div>
<div id="faqGroup" style="display:none"><div class="form-group"><label>{{T "ar_faq"}}</label><textarea name="faq" class="form-control" rows="5" placeholder="Apa produk?|Software WA marketing"></textarea></div></div>
<div class="form-group"><label>{{T "ar_reply"}}</label><textarea name="reply" class="form-control" rows="3" placeholder="{{T "ar_reply_ph"}}"></textarea></div>
<div class="bg-light border rounded p-3 mb-3">
<div class="form-check"><input class="form-check-input" type="checkbox" name="use_ai" value="1" id="useAiCheck"><label class="form-check-label" for="useAiCheck">{{T "ar_use_ai"}}</label></div>
<div class="form-group mt-2" id="aiKeyGroup" style="display:none"><label>{{T "ar_ai_key"}}</label><select name="ai_key_id" class="form-control">{{range .AiKeys}}<option value="{{.ID}}">{{.Name}} ({{.Provider}})</option>{{end}}</select></div>
</div>
<div class="form-group"><label>Nomor WA</label><div class="border rounded p-2" style="max-height:120px;overflow-y:auto">{{range .ConnectedAccounts}}{{if .Phone}}<div class="form-check form-check-inline small"><input class="form-check-input" type="checkbox" name="account_ids" value="+{{.Phone}}" id="a_{{.Phone}}"><label for="a_{{.Phone}}">+{{.Phone}}</label></div>{{end}}{{end}}</div></div>
<div class="mb-2"><label class="field-label">Training Campaign</label><select name="training_id" class="form-control form-control-sm"><option value="0">-- Default --</option>{{range .AiTrainings}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
<button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
</form></div></div></div>
<div class="col-12 col-lg-7">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ar_list_title"}}</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Keyword</th><th>Reply</th><th>Status</th><th></th></tr></thead><tbody>
{{range .AutoReplies}}<tr><td>{{.ID}}</td><td><strong>{{.Keyword}}</strong></td><td>{{if .UseAI}}<span class="badge bg-warning bg-opacity-10 text-warning me-1 small">AI</span>{{end}}{{.Reply}}</td><td>{{if .IsActive}}<span class="badge bg-success bg-opacity-10 text-success small">ON</span>{{else}}<span class="badge bg-danger bg-opacity-10 text-danger small">OFF</span>{{end}}</td><td><a class="btn btn-sm btn-white px-2" href="/autoreply?edit={{.ID}}">??</a><form method="post" action="/autoreply/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white px-2">{{if .IsActive}}OFF{{else}}ON{{end}}</button></form><form method="post" action="/autoreply/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm text-danger px-2">&times;</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted py-3 text-center">Belum ada rule.</td></tr>{{end}}
</tbody></table></div></div></div></div>
{{end}}
{{if eq .Page "settings"}}
  <div class="row justify-content-center">
    <div class="col-12 col-lg-8">
      <form method="post" action="/settings">
        <div class="card">
          <div class="card-header"><h4 class="card-header-title">{{T "set_welcome_title"}}</h4>
            <div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="welcome_enabled" {{if .WelcomeEnabled}}checked{{end}}></div>
          </div>
          <div class="card-body">
            <label>{{T "set_welcome_msg"}}</label>
            <textarea name="welcome_message" class="form-control" rows="3" placeholder="{{T "set_welcome_ph"}}">{{.WelcomeMessage}}</textarea>
            <small class="form-text text-muted">{{T "set_vars_hint"}}</small>
          </div>
        </div>

        <div class="card">
          <div class="card-header"><h4 class="card-header-title">{{T "set_fallback_title"}}</h4>
            <div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="fallback_enabled" {{if .FallbackEnabled}}checked{{end}}></div>
          </div>
          <div class="card-body">
            <label>{{T "set_fallback_msg"}}</label>
            <textarea name="fallback_message" class="form-control" rows="3" placeholder="{{T "set_fallback_ph"}}">{{.FallbackMessage}}</textarea>
            <small class="form-text text-muted">{{T "set_fallback_hint"}}</small>
          </div>
        </div>

        <div class="card">
          <div class="card-header"><h4 class="card-header-title">{{T "set_group_title"}}</h4>
            <div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="reply_in_group" {{if .ReplyInGroup}}checked{{end}}></div>
          </div>
          <div class="card-body"><small class="form-text text-muted">{{T "set_group_hint"}}</small></div>
        </div>

        <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "set_save"}}</button>
      </form>
    </div>
  </div>
{{end}}

{{if eq .Page "contacts"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ct_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/contacts/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    {{if .EditID}}
    <div class="col-12 col-lg-4">
      <div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> Edit</h4></div>
        <div class="card-body"><form method="post" action="/contacts/edit">
          <input type="hidden" name="id" value="{{.EditID}}">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" value="{{.EditName}}" required></div>
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" value="{{.EditPhone}}" required></div>
          <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "set_save"}}</button> <a href="/contacts" class="btn btn-white ms-2">{{T "ar_cancel"}}</a>
        </form></div>
      </div>
    </div>
    {{end}}
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_contacts_saved"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_from"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Contacts}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Phone}}</td><td>
            <a class="btn btn-sm btn-white" href="/contacts?edit={{.ID}}"><i class="la la-edit"></i></a>
            <form method="post" action="/contacts/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "groups"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "grp_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/groups/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_contacts_groups"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "grp_members"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Groups}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Count}}</td><td><form method="post" action="/groups/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "unsub"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "unsub_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/unsub/add">
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_contacts_unsub"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_from"}}</th><th>{{T "col_time"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Unsubs}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td class="text-muted">{{.Created}}</td><td><form method="post" action="/unsub/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "broadcast"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "bc_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/broadcast">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <div class="form-group"><label>{{T "bc_groups"}}</label><select name="groups" class="form-control" multiple required>{{range .Groups}}<option value="{{.ID}}">{{.Name}} ({{.Count}})</option>{{end}}</select></div>
          <div class="form-group"><label>{{T "bc_account"}}</label><div class="border rounded p-2" style="max-height:160px;overflow-y:auto">{{range .ConnectedAccounts}}{{if .Phone}}<div class="form-check"><input class="form-check-input" type="checkbox" name="account_ids" value="+{{.Phone}}" id="bc_{{.Phone}}"><label class="form-check-label small" for="bc_{{.Phone}}">+{{.Phone}}</label></div>{{end}}{{end}}{{if not .HasConnected}}<small class="text-muted">Belum ada nomor terkoneksi</small>{{end}}</div><small class="form-text text-muted">Biarkan kosong = random semua nomor. Checklist = hanya nomor itu.</small><div class="form-group"><label>Interval (detik) <small class="text-muted">jeda antar pesan</small></label><input name="interval" type="number" class="form-control" value="300" min="30" placeholder="300-400"></div></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea><small class="form-text text-muted">{{T "set_vars_hint"}}</small></div>
          <button class="btn btn-primary lift" {{if not .HasConnected}}disabled{{end}}><i class="la la-bullhorn me-1"></i> {{T "bc_start"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_broadcast"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "bc_progress"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Campaigns}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><a href="/broadcast/detail?id={{.ID}}" title="Lihat detail nomor terkirim">{{.Sent}}/{{.Total}}</a></td><td><span class="badge badge-soft-warning">{{.Status}}</span></td><td>
            <form method="post" action="/broadcast/stop" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{T "bc_stop"}}</button></form>
            <form method="post" action="/broadcast/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form>
          </td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "scheduled"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "sch_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/scheduled">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <div class="form-group"><label>{{T "col_to"}}</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <div class="form-group"><label>{{T "sch_time"}}</label><input type="datetime-local" name="send_at" class="form-control" required></div>
          <div class="form-group"><label>{{T "sch_repeat"}}</label><input type="number" name="repeat" class="form-control" value="0"><small class="form-text text-muted">{{T "sch_repeat_hint"}}</small></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea></div>
          <button class="btn btn-primary lift"><i class="la la-clock me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_scheduled"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_to"}}</th><th>{{T "sch_time"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Scheduleds}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Phone}}</td><td class="text-muted">{{.SendAt}}</td><td><span class="badge badge-soft-warning">{{.Status}}</span></td><td><form method="post" action="/scheduled/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "templates"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "tpl_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/templates/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="content" class="form-control" rows="4" required></textarea><small class="form-text text-muted">{{T "set_vars_hint"}}</small></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    {{if .EditID}}
    <div class="col-12 col-lg-5">
      <div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> Edit</h4></div>
        <div class="card-body"><form method="post" action="/templates/edit">
          <input type="hidden" name="id" value="{{.EditID}}">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" value="{{.EditName}}" required></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="content" class="form-control" rows="4" required>{{.EditContent}}</textarea></div>
          <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "set_save"}}</button> <a href="/templates" class="btn btn-white ms-2">{{T "ar_cancel"}}</a>
        </form></div>
      </div>
    </div>
    {{end}}
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_templates"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_message"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Templates}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Content}}</td><td>
            <a class="btn btn-sm btn-white" href="/templates?edit={{.ID}}"><i class="la la-edit"></i></a>
            <form method="post" action="/templates/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "apikeys"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "key_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/apikeys/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "key_generate"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_apikeys"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Secret</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .APIKeys}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><code>{{.Secret}}</code></td><td><form method="post" action="/apikeys/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "webhooks"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "wh_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/webhooks/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
          <div class="form-group"><label>URL</label><input name="url" class="form-control" placeholder="https://..." required></div>
          <div class="form-group"><label>{{T "wh_event"}}</label><select name="event" class="form-control"><option value="all">{{T "wh_all"}}</option><option value="received">{{T "nav_received"}}</option><option value="sent">{{T "nav_sent"}}</option></select></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_webhooks"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>URL</th><th>{{T "wh_event"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Webhooks}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.URL}}</td><td><span class="badge badge-soft-secondary">{{.Event}}</span></td><td><form method="post" action="/webhooks/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "logger"}}
  <div class="card">
    <div class="card-header"><h4 class="card-header-title">{{T "nav_logger"}}</h4>
      <form method="post" action="/logger/clear"><button class="btn btn-sm btn-white">{{T "log_clear"}}</button></form>
    </div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "log_type"}}</th><th>{{T "log_reason"}}</th><th>{{T "col_message"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Logs}}<tr><td>{{.ID}}</td><td><span class="badge badge-soft-secondary">{{.Type}}</span></td><td>{{.Reason}}</td><td>{{.Content}}</td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "hosts_android"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dev_add"}}</h4></div><div class="card-body">
      <form method="post" action="/devices/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-group"><label>Device ID</label><input name="did" class="form-control"></div>
        <div class="form-group"><label>{{T "dev_manuf"}}</label><input name="manufacturer" class="form-control"></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_hosts_android"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Device ID</th><th>{{T "dev_manuf"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Devices}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.DID}}</td><td>{{.Manufacturer}}</td><td><form method="post" action="/devices/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "login"}}
<h3>Masuk</h3>
<p class="sub">Belum punya akun? <a href="/register">Daftar gratis</a></p>
<form method="post" action="/login/post">
  <div class="form-group mb-3"><label class="small fw-bold text-muted">Email</label><input type="email" name="email" class="form-control" placeholder="admin@chatgo.test" required></div>
  <div class="form-group mb-3"><label class="small fw-bold text-muted">Password</label><input type="password" name="password" class="form-control" placeholder="••••••••" required></div>
  <button class="btn" type="submit">Masuk</button>
</form>
<div class="auth-divider"><span>atau</span></div>
<div class="demo-box">
  <div class="demo-title">🧪 Demo Login</div>
  <div class="demo-row"><strong>Admin:</strong> admin@chatgo.test / password</div>
  <div class="demo-row"><strong>User:</strong> saas_005357@test.com / password</div>
</div>
{{end}}

{{if eq .Page "register"}}
<h3>Daftar</h3>
<p class="sub">Sudah punya akun? <a href="/login">Masuk</a></p>
<form method="post" action="/register/post">
  <div class="form-group mb-3"><label class="small fw-bold text-muted">{{T "ar_nama"}}</label><input name class="form-control" placeholder="Nama Anda" required></div>
  <div class="form-group mb-3"><label class="small fw-bold text-muted">Email</label><input type="email" name="email" class="form-control" placeholder="email@domain.com" required></div>
  <div class="form-group mb-3"><label class="small fw-bold text-muted">Password</label><input type="password" name="password" class="form-control" placeholder="Min. 6 karakter" required></div>
  <button class="btn" type="submit">Daftar</button>
</form>
{{end}}

{{if eq .Page "hosts_whatsapp"}}
  <div class="row justify-content-center"><div class="col-12 col-lg-8"><div class="card">
    <div class="card-header"><h4 class="card-header-title">{{T "nav_hosts_whatsapp"}}</h4><a href="/wa" class="btn btn-sm btn-primary">{{T "nav_account_qr"}}</a></div>
    <div class="card-body text-center">
      {{if eq .Status "connected"}}<span class="h1 la la-check-circle text-success"></span><h3>{{T "wa_connected_title"}}</h3><p class="text-muted">+{{.Phone}}</p>
      {{else}}<p class="text-muted">{{T "hosts_wa_hint"}}</p><a href="/wa" class="btn btn-primary lift"><i class="la la-whatsapp me-1"></i> {{T "nav_account_qr"}}</a>{{end}}
    </div>
  </div></div></div>
{{end}}

{{if eq .Page "ussd"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ussd_add"}}</h4></div><div class="card-body">
      <form method="post" action="/ussd/add"><div class="form-group"><label>{{T "ussd_code"}}</label><input name="code" class="form-control" placeholder="*123#" required></div><button class="btn btn-primary lift"><i class="la la-satellite-dish me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ussd"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "ussd_code"}}</th><th>{{T "ussd_response"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Ussds}}<tr><td>{{.ID}}</td><td>{{.Code}}</td><td>{{.Response}}</td><td><span class="badge badge-soft-warning">{{.Status}}</span></td><td><form method="post" action="/ussd/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "ai_keys"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "aik_add"}}</h4></div><div class="card-body">
      <form method="post" action="/ai/keys/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-group"><label>Provider</label><select name="provider" class="form-control"><option value="openai">OpenAI</option><option value="geminiai">Gemini</option><option value="claudeai">Claude</option><option value="deepseekai">DeepSeek</option></select></div>
        <div class="form-group"><label>Model</label><input name="model" class="form-control" placeholder="gpt-4o"></div>
        <div class="form-group"><label>API Key</label><input name="apikey" class="form-control" required></div>
        <div class="form-group"><label>Prompt</label><textarea name="system_prompt" class="form-control" rows="3"></textarea></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_keys"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Provider</th><th>Model</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .AiKeys}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><span class="badge badge-soft-secondary">{{.Provider}}</span></td><td>{{.Model}}</td><td><form method="post" action="/ai/keys/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "ai_plugins"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "aip_add"}}</h4></div><div class="card-body">
      <form method="post" action="/ai/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div><div class="form-group"><label>Endpoint</label><input name="endpoint" class="form-control" placeholder="https://..."></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_plugins"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Endpoint</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .AiPlugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Endpoint}}</td><td><form method="post" action="/ai/plugins/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin"}}
  <div class="row">
    <div class="col-6 col-lg-3"><div class="card"><div class="card-body"><h6 class="text-uppercase text-muted mb-2">{{T "adm_users"}}</h6><a href="/admin/users" class="h2 mb-0 d-block">{{T "btn_all"}} <i class="la la-users"></i></a></div></div></div>
    <div class="col-6 col-lg-3"><div class="card"><div class="card-body"><h6 class="text-uppercase text-muted mb-2">{{T "adm_packages"}}</h6><a href="/admin/packages" class="h2 mb-0 d-block">{{T "btn_all"}} <i class="la la-box"></i></a></div></div></div>
    <div class="col-6 col-lg-3"><div class="card"><div class="card-body"><h6 class="text-uppercase text-muted mb-2">{{T "adm_waservers"}}</h6><a href="/admin/waservers" class="h2 mb-0 d-block">{{T "btn_all"}} <i class="la la-server"></i></a></div></div></div>
    <div class="col-6 col-lg-3"><div class="card"><div class="card-body"><h6 class="text-uppercase text-muted mb-2">{{T "adm_transactions"}}</h6><a href="/admin/transactions" class="h2 mb-0 d-block">{{T "btn_all"}} <i class="la la-money-bill"></i></a></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_users"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "usr_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/users/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-group"><label>Email</label><input name="email" type="email" class="form-control" required></div>
        <div class="form-group"><label>Password</label><input name="password" type="password" class="form-control"></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "usr_role"}}</label><select name="role" class="form-control">{{range .Roles}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div>
        <div class="form-group col-6"><label>{{T "usr_country"}}</label><select name="country" class="form-control"><option value="ID">Indonesia</option><option value="US">United States</option></select></div></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "usr_lang"}}</label><select name="language" class="form-control">{{range .LanguagesAdm}}<option value="{{.ISO}}">{{.Name}}</option>{{end}}</select></div>
        <div class="form-group col-6"><label>{{T "usr_theme"}}</label><select name="theme_color" class="form-control"><option value="light">Light</option><option value="dark">Dark</option></select></div></div>
        <div class="form-row"><div class="form-group col-4"><label>{{T "usr_clock"}}</label><select name="clock_format" class="form-control"><option value="g:i A">12h</option><option value="H:i">24h</option></select></div>
        <div class="form-group col-4"><label>{{T "usr_date"}}</label><select name="date_format" class="form-control"><option value="n/j/Y">MM/DD/YYYY</option><option value="j/n/Y">DD/MM/YYYY</option></select></div>
        <div class="form-group col-4"><label>{{T "usr_sep"}}</label><select name="date_separator" class="form-control"><option value="/">/</option><option value="-">-</option></select></div></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "usr_timezone"}}</label><input name="timezone" class="form-control" value="asia/jakarta"></div>
        <div class="form-group col-6"><label>{{T "usr_credits"}}</label><input name="credits" type="number" class="form-control" value="0"></div></div>
        <div class="form-row"><div class="form-group col-4"><label>{{T "usr_alertsound"}}</label><select name="alertsound" class="form-control"><option value="1">On</option><option value="2">Off</option></select></div>
        <div class="form-group col-4"><label>{{T "usr_partner"}}</label><select name="partner" class="form-control"><option value="1">Yes</option><option value="2">No</option></select></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_users"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Email</th><th>{{T "usr_role"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Users}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Email}}</td><td><span class="badge badge-soft-secondary">{{.Role}}</span></td><td><form method="post" action="/admin/users/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_roles"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "role_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/roles/add"><div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
      <div class="form-group"><label>{{T "role_perms"}}</label><select name="permissions" class="form-control" multiple><option value="manage_users">Users</option><option value="manage_packages">Packages</option><option value="manage_waservers">WA Servers</option><option value="manage_plugins">Plugins</option></select></div>
      <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_roles"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "role_perms"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Roles}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Permissions}}</td><td><form method="post" action="/admin/roles/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_packages"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "pkg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/packages/add">
        <div class="form-row"><div class="form-group col-6"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-group col-3"><label>{{T "pkg_price"}}</label><input name="price" class="form-control" value="0"></div>
        <div class="form-group col-3"><label>{{T "pkg_hidden"}}</label><select name="hidden" class="form-control"><option value="1">Hidden</option><option value="2">Visible</option></select></div></div>
        <div class="form-group"><label>{{T "pkg_services"}}</label><select name="services" class="form-control" multiple><option value="whatsapp">WhatsApp</option><option value="api">API</option><option value="webhooks">Webhooks</option><option value="templates">Templates</option><option value="ai">AI</option></select></div>
        <div class="form-group"><label>{{T "pkg_footermark"}}</label><select name="footermark" class="form-control"><option value="2">Off</option><option value="1">On</option></select></div>
        <hr><h6 class="text-uppercase text-muted small">Limits</h6>
        <div class="form-row">
        <div class="form-group col-4"><label>Send</label><input name="send_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>Receive</label><input name="receive_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>USSD</label><input name="ussd_limit" type="number" class="form-control" value="0"></div></div>
        <div class="form-row">
        <div class="form-group col-4"><label>Device</label><input name="device_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-4"><label>WA Send</label><input name="wa_send_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>WA Receive</label><input name="wa_receive_limit" type="number" class="form-control" value="100"></div></div>
        <div class="form-row">
        <div class="form-group col-3"><label>WA Acc</label><input name="wa_account_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-3"><label>Contact</label><input name="contact_limit" type="number" class="form-control" value="50"></div>
        <div class="form-group col-3"><label>Scheduled</label><input name="scheduled_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-3"><label>API Key</label><input name="key_limit" type="number" class="form-control" value="5"></div></div>
        <div class="form-row">
        <div class="form-group col-4"><label>Webhook</label><input name="webhook_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-4"><label>Action</label><input name="action_limit" type="number" class="form-control" value="5"></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_packages"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "pkg_price"}}</th><th>Limits</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Packages}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Price}}</td><td><small>S:{{.SendLimit}} D:{{.DeviceLimit}}</small></td><td><form method="post" action="/admin/packages/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_vouchers"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "vch_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/vouchers/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "vch_count"}}</label><input name="count" type="number" class="form-control" value="1"></div>
        <div class="form-group col-6"><label>{{T "vch_duration"}}</label><input name="duration" type="number" class="form-control" value="30"></div></div>
        <div class="form-group"><label>{{T "adm_packages"}}</label><select name="pkg" class="form-control">{{range .Packages}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "vch_generate"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_vouchers"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "vch_code"}}</th><th>{{T "adm_packages"}}</th><th>{{T "vch_duration"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Vouchers}}<tr><td>{{.ID}}</td><td><code>{{.Code}}</code></td><td>{{.Pkg}}</td><td>{{.Duration}}d</td><td><form method="post" action="/admin/vouchers/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_subscriptions"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "sub_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/subscriptions/add"><div class="form-group"><label>User</label><input name="user" class="form-control" required></div><div class="form-group"><label>{{T "adm_packages"}}</label><select name="pkg" class="form-control">{{range .Packages}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div><div class="form-group"><label>{{T "sub_expire"}}</label><input name="expire" type="date" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_subscriptions"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>User</th><th>{{T "adm_packages"}}</th><th>{{T "sub_expire"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Subscriptions}}<tr><td>{{.ID}}</td><td>{{.User}}</td><td>{{.Pkg}}</td><td>{{.Expire}}</td><td><form method="post" action="/admin/subscriptions/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_transactions"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_transactions"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>User</th><th>{{T "trx_amount"}}</th><th>Provider</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Transactions}}<tr><td>{{.ID}}</td><td>{{.User}}</td><td>{{.Amount}}</td><td>{{.Provider}}</td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "admin_payouts"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_payouts"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>User</th><th>{{T "trx_amount"}}</th><th>{{T "pay_address"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
      {{range .Payouts}}<tr><td>{{.ID}}</td><td>{{.User}}</td><td>{{.Amount}}</td><td>{{.Address}}</td><td><span class="badge badge-soft-warning">{{.Status}}</span></td><td><form method="post" action="/admin/payouts/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "admin_pages"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "pg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/pages/add">
        <div class="form-row"><div class="form-group col-8"><label>{{T "pg_title"}}</label><input name="title" class="form-control" required></div>
        <div class="form-group col-4"><label>Slug</label><input name="slug" class="form-control" placeholder="about-us"></div></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "pg_roles"}}</label><select name="roles" class="form-control" multiple>{{range .Roles}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
        <div class="form-group col-6"><label>{{T "pg_logged"}}</label><select name="logged" class="form-control"><option value="1">All</option><option value="2">Logged In Only</option></select></div></div>
        <div class="form-group"><label>{{T "col_message"}}</label><textarea name="content" class="form-control" rows="6"></textarea></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_pages"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "pg_title"}}</th><th>Slug</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Pages}}<tr><td>{{.ID}}</td><td>{{.Title}}</td><td>{{.Slug}}</td><td><form method="post" action="/admin/pages/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_marketing"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "mkt_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/marketing/add"><div class="form-group"><label>{{T "pg_title"}}</label><input name="title" class="form-control" required></div><div class="form-group"><label>{{T "col_message"}}</label><textarea name="content" class="form-control" rows="4"></textarea></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_marketing"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "pg_title"}}</th><th>{{T "col_time"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Marketings}}<tr><td>{{.ID}}</td><td>{{.Title}}</td><td class="text-muted">{{.Created}}</td><td><form method="post" action="/admin/marketing/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_languages"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "lng_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/languages/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-row"><div class="form-group col-6"><label>ISO</label><input name="iso" class="form-control" value="us" maxlength="2"></div>
        <div class="form-group col-6"><label>RTL</label><select name="rtl" class="form-control"><option value="2">LTR</option><option value="1">RTL</option></select></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_languages"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>ISO</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .LanguagesAdm}}<tr><td>{{.ID}}</td><td><span class="flag-icon flag-icon-{{.ISO}} lang-flag"></span>{{.Name}}</td><td>{{.ISO}}</td><td><form method="post" action="/admin/languages/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_waservers"}}
  <div class="row mb-4">
    <div class="col-12 col-lg-3"><div class="card"><div class="card-body text-center"><span class="h2 la la-server text-success"></span><h4>{{T "was_active"}}</h4><span class="h2">{{if eq .Status "connected"}}1{{else}}0{{end}}</span></div></div></div>
    <div class="col-12 col-lg-3"><div class="card"><div class="card-body text-center"><span class="h2 la la-whatsapp text-success"></span><h4>{{T "was_connected"}}</h4><span class="h2">{{.ConnectedCount}}</span></div></div></div>
    <div class="col-12 col-lg-3"><div class="card"><div class="card-body text-center"><span class="h2 la la-whatsapp text-danger"></span><h4>{{T "was_disconnected"}}</h4><span class="h2">{{.DisconnectedCount}}</span></div></div></div>
    <div class="col-12 col-lg-3"><div class="card"><div class="card-body text-center"><span class="h2 la la-list text-muted"></span><h4>{{T "was_total"}}</h4><span class="h2">{{.AccountLimit}}</span></div></div></div>
  </div>
  <div class="row">
    {{if .EditID}}
    <div class="col-12 col-lg-5"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h5 class="card-header-title mb-0">Edit WA Server #{{.EditID}}</h5></div><div class="card-body">
      <form method="post" action="/admin/waservers/edit"><input type="hidden" name="id" value="{{.EditID}}">
    {{else}}
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "was_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/waservers/add">
    {{end}}
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "was_accounts"}}</label><input name="accounts" type="number" class="form-control" value="{{if .EditID}}{{.EditContent}}{{else}}100{{end}}"></div>
        <div class="form-group col-6"><label>{{T "adm_packages"}}</label><select name="packages" class="form-control" multiple>{{range .Packages}}<option value="{{.Name}}" {{if and $.EditID (contains $.EditGroups .Name)}}selected{{end}}>{{.Name}}</option>{{end}}</select></div></div>
        <div class="form-row"><div class="form-group col-8"><label>URL</label><input name="url" class="form-control" placeholder="http://127.0.0.1" value="{{.EditContent}}"></div>
        <div class="form-group col-4"><label>Port</label><input name="port" class="form-control" placeholder="8080" value="{{.EditPhone}}"></div></div>
        <div class="form-group"><label>Secret</label><input name="secret" class="form-control" value="{{.EditKeyword}}"></div>
        <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{if .EditID}}Save{{else}}{{T "ar_add_btn"}}{{end}}</button>
        {{if .EditID}}<a href="/admin/waservers" class="btn btn-white ms-2">Batal</a>{{end}}
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_waservers"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>URL</th><th>{{T "was_accounts"}}</th><th>Packages</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .WaServers}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.URL}}:{{.Port}}</td><td>{{.Accounts}}</td><td><span class="badge bg-info bg-opacity-10 text-info small">{{.Packages}}</span></td><td><a class="btn btn-sm btn-white px-2" href="/admin/waservers?edit={{.ID}}">✏️</a><form method="post" action="/admin/waservers/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_gateways"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "gw_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/gateways/add">
        <div class="form-row"><div class="form-group col-8"><label>{{T "col_name"}}</label><input name class="form-control" required></div>
        <div class="form-group col-4"><label>{{T "gw_controller"}}</label><input name="controller" class="form-control" placeholder="gateway.php"></div></div>
        <div class="form-row"><div class="form-group col-4"><label>{{T "gw_callback"}}</label><select name="callback" class="form-control"><option value="1">On</option><option value="2">Off</option></select></div>
        <div class="form-group col-8"><label>{{T "gw_callback_id"}}</label><input name="callback_id" class="form-control"></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_gateways"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "gw_callback"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Gateways}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><span class="badge badge-soft-secondary">{{.Callback}}</span></td><td><form method="post" action="/admin/gateways/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_shorteners"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "sh_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/shorteners/add"><div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_shorteners"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Shorteners}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><form method="post" action="/admin/shorteners/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_plugins"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "plg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name class="form-control" required></div><div class="form-group"><label>{{T "plg_dir"}}</label><input name="dir" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_plugins"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "plg_dir"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Plugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Dir}}</td><td><form method="post" action="/admin/plugins/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "inbox"}}
<div class="card"><div class="card-header"><h4 class="card-header-title">Percakapan Terbaru</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Nomor</th><th>{{T "ar_nama"}}</th><th>Pesan Terakhir</th><th>Waktu</th></tr></thead><tbody>
{{range .Received}}<tr><td><a href="/inbox/chat?phone={{.Phone}}">{{.Phone}}</a></td><td>{{.Name}}</td><td class="msg-full">{{.Message}}</td><td class="text-muted">{{.Created}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">Belum ada percakapan</td></tr>{{end}}
</tbody></table></div></div>
{{end}}
{{if eq .Page "inbox_chat"}}
<div class="card"><div class="card-header"><h4 class="card-header-title"><a href="/inbox" class="text-decoration-none">&larr; Kembali</a> &nbsp; Chat dengan {{.Phone}}</h4></div>
<div style="max-height:500px;overflow-y:auto;padding:16px">
{{range .Received}}{{if eq .Phone $.Phone}}<div class="d-flex mb-3"><div class="bg-light rounded p-2" style="max-width:75%"><small class="text-muted">{{.Name}} &middot; {{.Created}}</small><br>{{.Message}}</div></div>{{end}}{{end}}
{{range .Sent}}{{if eq .Phone $.Phone}}<div class="d-flex mb-3 justify-content-end"><div class="bg-primary bg-opacity-10 rounded p-2" style="max-width:75%"><small class="text-muted">{{.Created}}</small><br>{{.Message}}</div></div>{{end}}{{end}}
</div></div>
{{end}}
{{if eq .Page "docs"}}
  <div class="card"><div class="card-body">
    <h3>{{T "docs_title"}}</h3>
    <p class="text-muted">{{T "docs_intro"}}</p>
    <div class="card mt-3"><div class="card-header"><h4>📋 {{T "docs_demo"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Role</th><th>Email</th><th>Keterangan</th></tr></thead>
        <tbody><tr><td>Admin</td><td>admin@chatgo.test</td><td>Akses penuh semua menu</td></tr></tbody></table></div>
    </div>
    <h4 class="mt-4">📖 {{T "docs_tutorial"}}</h4>
    <div class="row mt-3">{{range $i, $step := .DocsSteps}}
      <div class="col-12 col-lg-6"><div class="card mb-2"><div class="card-body"><h5>{{$step.Num}}. {{$step.Title}}</h5><p class="text-muted small">{{$step.Desc}}</p></div></div></div>{{end}}
    </div>
    <h4 class="mt-4">🤖 AI Auto Reply <a id="ai"></a></h4>
    <div class="card mt-2"><div class="card-body">
      <h5>Cara Setup AI Auto Reply</h5>
      <ol class="small text-muted mb-3">
        <li><strong>{{T "ar_ai_keys_list"}}</strong> — tambah API key (OpenAI/Gemini/Claude/DeepSeek), isi provider + model + key</li>
        <li><strong>Knowledge Base</strong> (opsional) — upload {{T "ar_faq_tab"}} via CSV, atau input manual Q&A. AI akan pakai ini sebagai referensi {{T "kb_answer_dot"}}</li>
        <li><strong>Auto Reply</strong> — buat rule baru:
          <ul>
            <li>Match Type: <code>AI</code> → balas semua pesan tanpa keyword</li>
            <li>Atau match type lain + <strong>centang "{{T "ar_use_ai"}}"</strong> → AI dipakai sebagai fallback</li>
          </ul>
        </li>
        <li>{{T "ar_faq_tab"}} di textarea Auto Reply (format: <code>{{T "kb_question_dot"}}|{{T "kb_answer_dot"}}</code> per baris) akan digabung dengan Knowledge Base</li>
      </ol>
      <h5>BYOK (Bring Your Own Key)</h5>
      <p class="small text-muted">Kamu input API key sendiri. ChatGo tidak menyediakan key bawaan. Key di-encrypt sebelum disimpan.</p>
      <h5>Variabel Spintax</h5>
      <p class="small text-muted">Semua balasan support: <code>{name}</code> <code>{phone}</code> <code>{message}</code> <code>{time}</code> <code>{date}</code>. Spintax: <code>{Halo|Hai|Hi}</code> — random setiap kirim.</p>
    </div></div>
    <h4 class="mt-4">🖼️ Screenshots</h4>
    <div class="row g-2 mt-2">
    <div class="col-4"><a href="/screens/01-login.png" target="_blank"><img src="/screens/01-login.png" class="img-fluid rounded border" alt="Login"></a><small class="text-muted d-block text-center">Login</small></div>
    <div class="col-4"><a href="/screens/02-dashboard.png" target="_blank"><img src="/screens/02-dashboard.png" class="img-fluid rounded border" alt="Dashboard"></a><small class="text-muted d-block text-center">Dashboard</small></div>
    <div class="col-4"><a href="/screens/03-wa-qr.png" target="_blank"><img src="/screens/03-wa-qr.png" class="img-fluid rounded border" alt="WA QR"></a><small class="text-muted d-block text-center">WA Account</small></div>
    <div class="col-4"><a href="/screens/04-send-message.png" target="_blank"><img src="/screens/04-send-message.png" class="img-fluid rounded border" alt="Send"></a><small class="text-muted d-block text-center">Kirim Pesan</small></div>
    <div class="col-4"><a href="/screens/05-broadcast.png" target="_blank"><img src="/screens/05-broadcast.png" class="img-fluid rounded border" alt="Broadcast"></a><small class="text-muted d-block text-center">Broadcast</small></div>
    <div class="col-4"><a href="/screens/07-autoreply-ai.png" target="_blank"><img src="/screens/07-autoreply-ai.png" class="img-fluid rounded border" alt="Auto Reply"></a><small class="text-muted d-block text-center">Auto Reply</small></div>
    <div class="col-4"><a href="/screens/09-contacts.png" target="_blank"><img src="/screens/09-contacts.png" class="img-fluid rounded border" alt="Contacts"></a><small class="text-muted d-block text-center">Kontak</small></div>
    <div class="col-4"><a href="/screens/12-inbox.png" target="_blank"><img src="/screens/12-inbox.png" class="img-fluid rounded border" alt="Inbox"></a><small class="text-muted d-block text-center">Inbox</small></div>
    <div class="col-4"><a href="/screens/14-admin-users.png" target="_blank"><img src="/screens/14-admin-users.png" class="img-fluid rounded border" alt="Admin Users"></a><small class="text-muted d-block text-center">Admin Users</small></div>
    </div>
    <h4 class="mt-4">🔌 API Reference</h4>
    <pre class="bg-light p-3 rounded"><code>POST /api/send  -H "X-API-Key: &lt;key&gt;"  -d '{"phone":"628xx","message":"text"}'</code></pre>
  </div></div>
{{end}}



{{if eq .Page "knowledge"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "kb_add"}}</h4></div><div class="card-body">
        <form method="post" action="/knowledge/add">
          <div class="form-group"><label>{{T "kb_title"}}</label><input name="title" class="form-control" placeholder="{{T "ar_faq_tab"}} Produk" required></div>
          <div class="form-group"><label>{{T "kb_question"}}</label><input name="question" class="form-control" placeholder="{{T "kb_question_dot"}}..." required></div>
          <div class="form-group"><label>{{T "kb_answer"}}</label><textarea name="answer" class="form-control" rows="3" placeholder="{{T "kb_answer_dot"}}..." required></textarea></div>
          <div class="form-group"><label>{{T "kb_category"}}</label><input name="category" class="form-control" placeholder="produk, harga"></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/import" enctype="multipart/form-data">
          <label class="text-muted small mb-2 d-block">{{T "kb_import"}}</label>
          <div class="input-group"><input type="text" name="title" class="form-control" placeholder="Judul (opsional)"><input type="file" name="file" class="form-control" accept=".csv,.txt" required><button class="btn btn-white">{{T "kb_upload"}}</button></div>
          <small class="form-text text-muted">{{T "kb_csv_hint"}} <a href="/web/sample-knowledge.csv" target="_blank">{{T "kb_sample"}}</a></small>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/url">
          <label class="text-muted small mb-2 d-block">{{T "kb_url"}}</label>
          <div class="input-group"><input type="text" name="title" class="form-control" placeholder="Judul (opsional)"><input type="url" name="url" class="form-control" placeholder="https://..." required><button class="btn btn-white">{{T "kb_train"}}</button></div>
          <small class="form-text text-muted">{{T "kb_url_hint"}}</small>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/pdf" enctype="multipart/form-data">
          <label class="text-muted small mb-2 d-block">📄 Upload PDF</label>
          <div class="input-group"><input type="text" name="title" class="form-control" placeholder="Judul (opsional)"><input type="file" name="file" class="form-control" accept=".pdf" required><button class="btn btn-white">Upload</button></div>
          <small class="form-text text-muted">Upload PDF (brosur, daftar harga). Teks akan diekstrak otomatis.</small>
        </form>
      </div></div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-robot text-warning me-1"></i> {{T "nav_knowledge"}}</h4></div>
        <p class="text-muted px-4 pt-2">{{T "kb_info"}}</p>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "kb_title"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Knowledges}}<tr><td>{{.ID}}</td><td>{{.Title}}</td>
            <td>{{if .Active}}<span class="badge badge-soft-success">{{T "ar_active"}}</span>{{else}}<span class="badge badge-soft-danger">{{T "ar_off"}}</span>{{end}}</td>
            <td>
              <form method="post" action="/knowledge/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{if .Active}}{{T "ar_off"}}{{else}}{{T "ar_on"}}{{end}}</button></form>
              <form method="post" action="/knowledge/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form>
            </td>
          </tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}{{end}}`
