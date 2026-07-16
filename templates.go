package main

const templates = `
{{define "layout"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
<title>{{.Title}} &middot; {{.AppName}}</title>
<link rel="icon" href="/assets/theme/default-favicon.png">
<link rel="stylesheet" href="/assets/_assets/css/libs/line-awesome.min.css">
<link rel="stylesheet" href="/assets/_assets/css/libs/flag-icon.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/fonts/feather/feather.css">
<link rel="stylesheet" href="/assets/dashboard/css/libs/bootstrap.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/style.min.css">
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
<style>
  .navbar-vertical{overflow-y:auto}
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
  .badge-soft-primary{background:rgba(44,123,229,.12);color:#2c7be5}
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
  <nav class="navbar navbar-expand-md navbar-light d-none d-md-flex pe-3" id="topbar">
    <div class="container-fluid">
      <div class="me-4">
        <a class="btn btn-md btn-primary mb-1 lift" href="/wa"><i class="la la-whatsapp la-lg me-1"></i> {{T "nav_whatsapp"}}</a>
        <a class="btn btn-md btn-primary mb-1 lift" href="/send"><i class="la la-paper-plane la-lg me-1"></i> {{T "nav_send"}}</a>
      </div>
      <div class="ms-auto"></div>
      <div class="navbar-user d-flex align-items-center me-2">
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
        <div class="dropdown ms-3">
          <a href="#" class="dropdown-toggle text-muted" role="button" data-bs-toggle="dropdown" style="text-decoration:none">
            <i class="la la-user-circle la-lg"></i>
          </a>
          <div class="dropdown-menu dropdown-menu-end">
            <a class="dropdown-item" href="/settings"><i class="la la-cog me-2"></i> {{T "nav_settings"}}</a>
            <div class="dropdown-divider"></div>
            {{if .IsImpersonating}}
            <a class="dropdown-item text-warning" href="/exit-impersonation"><i class="la la-times-circle me-2"></i> Exit Impersonation</a>
            {{else}}
            <a class="dropdown-item text-danger" href="/logout"><i class="la la-sign-out me-2"></i> Logout</a>
            {{end}}
          </div>
        </div>
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
<script>
window.addEventListener('DOMContentLoaded',function(){
setTimeout(function(){
var s=document.getElementById('sidebar');
var a=s&&s.querySelector('.nav-link.active');
if(a){var t=a.offsetTop-s.offsetHeight/2;if(t>0)s.scrollTop=t}
},100);
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
      <img src="{{.AppLogo}}" class="navbar-brand-img mx-auto" alt="{{.AppName}}" onerror="this.outerHTML='<span style=&quot;color:#fff;font-weight:800;font-size:20px&quot;>{{.AppName}}</span>'">
    </a>
    <div class="collapse navbar-collapse" id="sidebarCollapse">
      <h6 class="navbar-heading">{{T "nav_overview"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "home"}}active{{end}}" href="/"><i class="la la-chart-bar la-lg"></i> {{T "nav_dashboard"}} <span class="badge badge-soft-primary" style="font-size:8px">META</span></a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "inbox"}}active{{end}}" href="/inbox"><i class="la la-comments la-lg"></i> Live Chat{{if gt .UnreadCount 0}} <span class="badge badge-pill badge-danger ml-1 inbox-badge">{{.UnreadCount}}</span>{{end}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      {{if eq .Role "admin"}}
      <h6 class="navbar-heading">{{T "nav_whatsapp"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "wa"}}active{{end}}" href="/wa"><i class="la la-whatsapp la-lg"></i> {{T "nav_account_qr"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "send"}}active{{end}}" href="/send"><i class="la la-paper-plane la-lg"></i> {{T "nav_send"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "broadcast"}}active{{end}}" href="/broadcast"><i class="la la-bullhorn la-lg"></i> {{T "nav_broadcast"}} <span class="badge badge-soft-primary" style="font-size:8px">META</span></a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "drips"}}active{{end}}" href="/drips"><i class="la la-tint la-lg"></i> Drip Campaign</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "scheduled"}}active{{end}}" href="/scheduled"><i class="la la-clock la-lg"></i> {{T "nav_scheduled"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "sent"}}active{{end}}" href="/sent"><i class="la la-telegram la-lg"></i> {{T "nav_sent"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "received"}}active{{end}}" href="/received"><i class="la la-comment la-lg"></i> {{T "nav_received"}}</a></li>
      </ul>
      {{else}}
      <h6 class="navbar-heading">{{T "nav_whatsapp"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "wa"}}active{{end}}" href="/wa"><i class="la la-whatsapp la-lg"></i> {{T "nav_account_qr"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "send"}}active{{end}}" href="/send"><i class="la la-paper-plane la-lg"></i> {{T "nav_send"}}</a></li>
      </ul>
      {{end}}
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
      {{if eq .Role "admin"}}
      <h6 class="navbar-heading">{{T "nav_contacts"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "contacts"}}active{{end}}" href="/contacts"><i class="la la-address-book la-lg"></i> {{T "nav_contacts_saved"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "groups"}}active{{end}}" href="/contacts/groups"><i class="la la-list la-lg"></i> {{T "nav_contacts_groups"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "tags"}}active{{end}}" href="/tags"><i class="la la-tags la-lg"></i> Tags</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "unsub"}}active{{end}}" href="/contacts/unsub"><i class="la la-unlink la-lg"></i> {{T "nav_contacts_unsub"}}</a></li>
      </ul>
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_tools"}}</h6>
      <ul class="nav nav-sm flex-column">
        <li class="nav-item"><a class="nav-link {{if eq .Active "settings"}}active{{end}}" href="/settings"><i class="la la-cog la-lg"></i> {{T "nav_settings"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "autoreply"}}active{{end}}" href="/autoreply"><i class="la la-robot la-lg"></i> {{T "nav_autoreply"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "templates"}}active{{end}}" href="/templates"><i class="la la-file-alt la-lg"></i> {{T "nav_templates"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "canned"}}active{{end}}" href="/canned"><i class="la la-comment-dots la-lg"></i> Canned</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "tracker"}}active{{end}}" href="/tracker"><i class="la la-link la-lg"></i> Links</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "abtests"}}active{{end}}" href="/ab-tests"><i class="la la-balance-scale la-lg"></i> A/B Test</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "store"}}active{{end}}" href="/store"><i class="la la-store la-lg"></i> Store</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "orders"}}active{{end}}" href="/store/orders"><i class="la la-shopping-bag la-lg"></i> Orders</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "forms"}}active{{end}}" href="/forms"><i class="la la-wpforms la-lg"></i> Forms</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "reminders"}}active{{end}}" href="/reminders"><i class="la la-bell la-lg"></i> Reminders</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "analytics"}}active{{end}}" href="/analytics"><i class="la la-chart-pie la-lg"></i> Analytics</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "blacklist"}}active{{end}}" href="/blacklist"><i class="la la-ban la-lg"></i> Blacklist</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "csat"}}active{{end}}" href="/csat"><i class="la la-star la-lg"></i> CSAT</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "recurring"}}active{{end}}" href="/recurring"><i class="la la-redo-alt la-lg"></i> Recurring</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "uploads"}}active{{end}}" href="/uploads"><i class="la la-folder-open la-lg"></i> Files</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "calendar"}}active{{end}}" href="/calendar"><i class="la la-calendar la-lg"></i> Calendar</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "customers"}}active{{end}}" href="/customers"><i class="la la-users la-lg"></i> Customers</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "merge"}}active{{end}}" href="/merge"><i class="la la-code-branch la-lg"></i> Merge</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "macros"}}active{{end}}" href="/macros"><i class="la la-bolt la-lg"></i> Macros</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "translate"}}active{{end}}" href="/translate-tool"><i class="la la-language la-lg"></i> Translate</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "widget"}}active{{end}}" href="/widget-info"><i class="la la-code la-lg"></i> Widget</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "emailwa"}}active{{end}}" href="/email-wa"><i class="la la-envelope la-lg"></i> Email→WA</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "apikeys"}}active{{end}}" href="/apikeys"><i class="la la-key la-lg"></i> {{T "nav_apikeys"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "webhooks"}}active{{end}}" href="/webhooks"><i class="la la-code-branch la-lg"></i> {{T "nav_webhooks"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "logger"}}active{{end}}" href="/logger"><i class="la la-clipboard-list la-lg"></i> {{T "nav_logger"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "settings"}}active{{end}}" href="/settings"><i class="la la-cog la-lg"></i> {{T "nav_settings"}}</a></li>
      </ul>
      {{else}}
      <h6 class="navbar-heading">{{T "nav_tools"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "templates"}}active{{end}}" href="/templates"><i class="la la-file-alt la-lg"></i> {{T "nav_templates"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "apikeys"}}active{{end}}" href="/apikeys"><i class="la la-key la-lg"></i> {{T "nav_apikeys"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "settings"}}active{{end}}" href="/settings"><i class="la la-cog la-lg"></i> {{T "nav_settings"}}</a></li>
      </ul>
      {{end}}
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
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_paygateways"}}active{{end}}" href="/admin/gateways-pay"><i class="la la-credit-card la-lg"></i> Pay Gateways</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_transactions_pay"}}active{{end}}" href="/admin/transactions-pay"><i class="la la-receipt la-lg"></i> Pay Logs</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "backup"}}active{{end}}" href="/backup"><i class="la la-database la-lg"></i> Backup</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "audit"}}active{{end}}" href="/audit"><i class="la la-history la-lg"></i> Audit</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_payouts"}}active{{end}}" href="/admin/payouts"><i class="la la-hand-holding-usd la-lg"></i> {{T "adm_payouts"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_pages"}}active{{end}}" href="/admin/pages"><i class="la la-file la-lg"></i> {{T "adm_pages"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_marketing"}}active{{end}}" href="/admin/marketing"><i class="la la-bullhorn la-lg"></i> {{T "adm_marketing"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_languages"}}active{{end}}" href="/admin/languages"><i class="la la-language la-lg"></i> {{T "adm_languages"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_waservers"}}active{{end}}" href="/admin/waservers"><i class="la la-server la-lg"></i> {{T "adm_waservers"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_gateways"}}active{{end}}" href="/admin/gateways"><i class="la la-code la-lg"></i> {{T "adm_gateways"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_shorteners"}}active{{end}}" href="/admin/shorteners"><i class="la la-link la-lg"></i> {{T "adm_shorteners"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_plugins"}}active{{end}}" href="/admin/plugins"><i class="la la-puzzle-piece la-lg"></i> {{T "adm_plugins"}}</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_meta"}}active{{end}}" href="/admin/meta"><i class="la la-cloud la-lg"></i> Meta API <span class="badge badge-soft-primary ms-1" style="font-size:9px">Meta</span></a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "admin_metatemplates"}}active{{end}}" href="/admin/metatemplates"><i class="la la-file-alt la-lg"></i> Templates <span class="badge badge-soft-primary ms-1" style="font-size:9px">Meta</span></a></li>
      </ul>
      {{end}}
      <hr class="navbar-divider my-3">
      <h6 class="navbar-heading">{{T "nav_docs"}}</h6>
      <ul class="navbar-nav">
        <li class="nav-item"><a class="nav-link {{if eq .Active "subscribe"}}active{{end}}" href="/subscribe"><i class="la la-shopping-cart la-lg"></i> Upgrade</a></li>
        <li class="nav-item"><a class="nav-link {{if eq .Active "docs"}}active{{end}}" href="/docs"><i class="la la-book la-lg"></i> {{T "nav_docs"}}</a></li>
      </ul>
    </div>
  </div>
</nav>{{end}}

{{define "landing"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.AppName}} — WhatsApp Marketing Platform</title>
<link rel="stylesheet" href="/assets/_assets/css/libs/line-awesome.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/libs/bootstrap.min.css">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;color:#1a1a2e;line-height:1.6}
.navbar{padding:16px 0;position:fixed;top:0;width:100%;z-index:100;background:rgba(255,255,255,.92);backdrop-filter:blur(12px);border-bottom:1px solid rgba(0,0,0,.06)}
.navbar .container{max-width:1140px;margin:0 auto;padding:0 24px;display:flex;justify-content:space-between;align-items:center}
.navbar-brand{font-size:22px;font-weight:800;color:#1a1a2e;text-decoration:none}
.navbar-brand span{color:#4F46E5}
.nav-links{display:flex;gap:16px;align-items:center}
.nav-links a{text-decoration:none;color:#555;font-weight:500;font-size:14px}
.nav-links .btn-login{padding:8px 20px;background:#4F46E5;color:#fff;border-radius:8px;font-weight:600}
.lang-switch{position:relative;display:inline-block}
.lang-switch select{appearance:none;padding:6px 28px 6px 10px;border:1px solid #ddd;border-radius:6px;font-size:13px;background:#fff;cursor:pointer;color:#555}
.lang-switch .flag-icon{position:absolute;right:8px;top:50%;transform:translateY(-50%);pointer-events:none;font-size:11px;color:#999}
.hero{padding:140px 24px 80px;text-align:center;max-width:900px;margin:0 auto}
.hero h1{font-size:3rem;font-weight:800;line-height:1.2;margin-bottom:16px;background:linear-gradient(135deg,#4F46E5,#7C3AED);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.hero p{font-size:1.15rem;color:#666;max-width:600px;margin:0 auto 32px}
.hero .cta-group{display:flex;gap:12px;justify-content:center;flex-wrap:wrap}
.hero .cta-group a{padding:12px 28px;border-radius:10px;font-weight:600;font-size:15px;text-decoration:none;transition:all .2s}
.hero .btn-primary{background:#4F46E5;color:#fff;box-shadow:0 4px 14px rgba(79,70,229,.3)}
.hero .btn-primary:hover{background:#4338CA;transform:translateY(-1px)}
.hero .btn-outline{border:2px solid #ddd;color:#444;background:#fff}
.hero .btn-outline:hover{border-color:#4F46E5;color:#4F46E5}
.features{padding:40px 24px 80px;max-width:1140px;margin:0 auto}
.features h2{text-align:center;font-size:2rem;font-weight:700;margin-bottom:12px}
.features .subtitle{text-align:center;color:#666;margin-bottom:48px}
.feature-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:24px}
.feature-card{padding:28px;border-radius:14px;border:1px solid #eee;transition:all .2s;background:#fff}
.feature-card:hover{box-shadow:0 8px 30px rgba(0,0,0,.08);transform:translateY(-2px)}
.feature-card i{font-size:2rem;color:#4F46E5;margin-bottom:12px}
.feature-card h4{font-size:1.1rem;font-weight:700;margin-bottom:6px}
.feature-card p{color:#666;font-size:.9rem}
.demo-section{padding:40px 24px 80px;background:#f8f9fc}
.demo-section .container{max-width:800px;margin:0 auto}
.demo-section h2{text-align:center;font-size:1.8rem;font-weight:700;margin-bottom:32px}
.demo-box{background:#f0f4ff;border:1px solid #d0d7f0;border-radius:12px;padding:20px;font-family:monospace;font-size:14px}
.demo-row{padding:6px 0}
.demo-row strong{color:#4F46E5}
.cta-banner{padding:60px 24px;text-align:center;background:linear-gradient(135deg,#4F46E5,#7C3AED);color:#fff}
.cta-banner h2{font-size:2rem;font-weight:700;margin-bottom:12px}
.cta-banner p{margin-bottom:24px;opacity:.9}
.cta-banner a{padding:14px 32px;background:#fff;color:#4F46E5;border-radius:10px;font-weight:700;text-decoration:none;display:inline-block;font-size:15px}
.footer{padding:24px;text-align:center;color:#888;font-size:13px;border-top:1px solid #eee}
@media(max-width:768px){.hero h1{font-size:2rem}.feature-grid{grid-template-columns:1fr}}
</style>
</head>
<body>

<nav class="navbar">
<div class="container">
<a href="/" class="navbar-brand">{{.AppName}}</a>
<div class="nav-links">
<div class="lang-switch"><select onchange="window.location=this.value">{{range .Languages}}<option value="/lang/{{.Code}}" {{if eq .Code $.LangCode}}selected{{end}}>{{.Flag}} {{.Name}}</option>{{end}}</select><span class="flag-icon">&#9660;</span></div>
<a href="/login">{{if eq .LangCode "id"}}Masuk{{else}}Sign In{{end}}</a>
<a href="/register" class="btn-login">{{if eq .LangCode "id"}}Daftar Gratis{{else}}Sign Up Free{{end}}</a>
</div>
</div>
</nav>

<section class="hero">
<h1>{{if eq .LangCode "id"}}WhatsApp Marketing Jadi Mudah{{else}}WhatsApp Marketing Made Easy{{end}}</h1>
<p>{{if eq .LangCode "id"}}{{.AppName}} adalah platform all-in-one untuk kirim broadcast, auto-reply AI, kelola multi-akun WhatsApp, dan live chat real-time — semua dalam satu dashboard.{{else}}{{.AppName}} is an all-in-one platform for sending broadcasts, AI auto-reply, managing multiple WhatsApp accounts, and real-time live chat — all in one dashboard.{{end}}</p>
<div class="cta-group">
<a href="/register" class="btn-primary">{{if eq .LangCode "id"}}Coba Gratis{{else}}Try Free{{end}}</a>
<a href="/docs" class="btn-outline">{{if eq .LangCode "id"}}Lihat Dokumentasi{{else}}View Documentation{{end}}</a>
</div>
<div style="max-width:400px;margin:32px auto 0;background:#fff;border-radius:14px;padding:24px;box-shadow:0 4px 24px rgba(0,0,0,.08)">
<form method="post" action="/login/post">
<div style="margin-bottom:12px"><input type="email" name="email" class="form-control" placeholder="Email" value="{{.AppEmail}}" style="border-radius:8px;padding:10px 14px;border:1px solid #ddd;width:100%;font-size:14px"></div>
<div style="margin-bottom:12px"><input type="password" name="password" class="form-control" placeholder="Password" value="password" style="border-radius:8px;padding:10px 14px;border:1px solid #ddd;width:100%;font-size:14px"></div>
<button type="submit" style="width:100%;padding:10px;background:#4F46E5;color:#fff;border:none;border-radius:8px;font-weight:600;font-size:14px;cursor:pointer">{{if eq .LangCode "id"}}Masuk{{else}}Sign In{{end}}</button>
</form>
<div style="text-align:center;margin-top:12px;font-size:12px;color:#999">Demo: <code style="background:#f0f0f0;padding:2px 6px;border-radius:4px">{{.AppEmail}}</code> / <code style="background:#f0f0f0;padding:2px 6px;border-radius:4px">password</code></div>
</div>
</section>

<section class="features">
<h2>{{if eq .LangCode "id"}}Fitur Lengkap{{else}}Complete Features{{end}}</h2>
<p class="subtitle">{{if eq .LangCode "id"}}Semua yang kamu butuhkan untuk WhatsApp marketing{{else}}Everything you need for WhatsApp marketing{{end}}</p>
<div class="feature-grid">
<div class="feature-card"><i class="la la-comments"></i><h4>{{T "inbox_title"}}</h4><p>{{if eq .LangCode "id"}}Inbox real-time dengan SSE, reply langsung, group chat, filter private/group.{{else}}Real-time inbox with SSE, direct reply, group chat, private/group filter.{{end}}</p></div>
<div class="feature-card"><i class="la la-robot"></i><h4>AI Auto Reply</h4><p>{{if eq .LangCode "id"}}Balas otomatis pakai AI (OpenAI/Gemini/Claude/DeepSeek) + knowledge base.{{else}}Auto reply with AI (OpenAI/Gemini/Claude/DeepSeek) + knowledge base.{{end}}</p></div>
<div class="feature-card"><i class="la la-bullhorn"></i><h4>Broadcast</h4><p>{{if eq .LangCode "id"}}Kirim pesan massal ke grup kontak, round-robin multi-akun WA.{{else}}Send bulk messages to contact groups, round-robin multi-WA accounts.{{end}}</p></div>
<div class="feature-card"><i class="la la-whatsapp"></i><h4>{{if eq .LangCode "id"}}Multi Akun{{else}}Multi Account{{end}}</h4><p>{{if eq .LangCode "id"}}Kelola banyak nomor WhatsApp sekaligus, scan QR pairing.{{else}}Manage multiple WhatsApp numbers at once, QR scan pairing.{{end}}</p></div>
<div class="feature-card"><i class="la la-cloud"></i><h4>Meta Cloud API</h4><p>{{if eq .LangCode "id"}}Integrasi resmi WhatsApp Business API + template pesan.{{else}}Official WhatsApp Business API integration + message templates.{{end}}</p></div>
<div class="feature-card"><i class="la la-clock"></i><h4>{{if eq .LangCode "id"}}Pesan Terjadwal{{else}}Scheduled Messages{{end}}</h4><p>{{if eq .LangCode "id"}}Jadwalkan pesan, repeat otomatis, pilih nomor pengirim.{{else}}Schedule messages, auto repeat, select sender number.{{end}}</p></div>
<div class="feature-card"><i class="la la-paint-brush"></i><h4>Whitelabel</h4><p>{{if eq .LangCode "id"}}Ganti logo, nama, email — satu binary, banyak domain.{{else}}Replace logo, name, email — one binary, many domains.{{end}}</p></div>
<div class="feature-card"><i class="la la-chart-bar"></i><h4>Dashboard Analytics</h4><p>{{if eq .LangCode "id"}}Chart aktivitas, statistik pesan, status koneksi real-time.{{else}}Activity charts, message statistics, real-time connection status.{{end}}</p></div>
</div>
</section>

<section class="demo-section">
<div class="container">
<h2>{{if eq .LangCode "id"}}Akun Demo{{else}}Demo Account{{end}}</h2>
<div class="demo-box" style="max-width:480px;margin:0 auto">
<div class="demo-row"><strong>Admin:</strong> {{.AppEmail}} / password</div>
</div>
<div style="text-align:center;margin-top:24px">
<a href="/login" style="display:inline-block;padding:12px 28px;border-radius:10px;font-weight:600;text-decoration:none;background:#4F46E5;color:#fff">{{if eq .LangCode "id"}}Masuk ke Dashboard{{else}}Go to Dashboard{{end}}</a>
</div>
</div>
</section>

<section class="cta-banner">
<h2>{{if eq .LangCode "id"}}Siap Tingkatkan WhatsApp Marketing Kamu?{{else}}Ready to Level Up Your WhatsApp Marketing?{{end}}</h2>
<p>{{if eq .LangCode "id"}}Daftar sekarang — gratis. Tanpa kartu kredit.{{else}}Sign up now — free. No credit card.{{end}}</p>
<a href="/register">{{if eq .LangCode "id"}}Daftar Gratis{{else}}Sign Up Free{{end}}</a>
</section>

<footer class="footer">&copy; 2026 {{.AppName}}. Powered by ChatGo.</footer>

</body>
</html>{{end}}

{{define "authpage"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}} &middot; {{.AppName}}</title>
<link rel="stylesheet" href="/assets/_assets/css/libs/line-awesome.min.css">
<link rel="stylesheet" href="/assets/_assets/css/libs/flag-icon.min.css">
<link rel="stylesheet" href="/assets/dashboard/css/libs/bootstrap.min.css">
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;min-height:100vh;display:flex;flex-direction:column}
.auth-wrap{flex:1;display:flex}
.auth-left{flex:1;background:linear-gradient(160deg,#0f1f33,#152e4d 40%,#1a3a5c);display:flex;flex-direction:column;justify-content:center;padding:60px 48px;position:relative;overflow:hidden;min-width:360px}
.auth-left::before{content:'';position:absolute;top:-100px;right:-100px;width:400px;height:400px;border-radius:50%;background:rgba(79,70,229,.08)}
.auth-left::after{content:'';position:absolute;bottom:-80px;left:-80px;width:300px;height:300px;border-radius:50%;background:rgba(79,70,229,.05)}
.auth-left h1{font-size:2.2rem;font-weight:800;color:#fff;line-height:1.2;margin-bottom:12px;position:relative;z-index:1}
.auth-left p{color:rgba(255,255,255,.6);font-size:1rem;position:relative;z-index:1;max-width:400px}
.auth-left .features{position:relative;z-index:1;margin-top:40px;display:flex;flex-direction:column;gap:16px}
.auth-left .feat{display:flex;align-items:center;gap:12px;color:rgba(255,255,255,.7)}
.auth-left .feat i{color:#4F46E5;font-size:1.3rem}
.auth-right{flex:1;display:flex;align-items:center;justify-content:center;padding:40px;background:#fff;min-width:360px}
.auth-card{width:100%;max-width:400px}
.auth-card h2{font-size:1.8rem;font-weight:700;margin-bottom:4px}
.auth-card .sub{color:#888;margin-bottom:24px;font-size:14px}
.auth-card .sub a{color:#4F46E5;font-weight:600;text-decoration:none}
.auth-card .form-group{margin-bottom:14px}
.auth-card label{font-size:12px;font-weight:600;text-transform:uppercase;color:#666;letter-spacing:.5px;margin-bottom:4px;display:block}
.auth-card input{border-radius:10px;padding:11px 14px;border:1.5px solid #e0e0e0;width:100%;font-size:14px;transition:border .15s;background:#fff;color:#1a1a2e;display:block}
.auth-card input:focus{outline:none;border-color:#4F46E5;box-shadow:0 0 0 3px rgba(79,70,229,.1)}
.auth-card .btn-submit{width:100%;padding:12px;background:#4F46E5;color:#fff;border:none;border-radius:10px;font-weight:600;font-size:15px;cursor:pointer;box-shadow:0 4px 14px rgba(79,70,229,.3);transition:all .15s}
.auth-card .btn-submit:hover{background:#4338CA;transform:translateY(-1px)}
.auth-divider{display:flex;align-items:center;gap:12px;margin:20px 0;color:#aaa;font-size:13px}
.auth-divider::before,.auth-divider::after{content:'';flex:1;height:1px;background:#e0e0e0}
.demo-box{background:#f4f6ff;border:1px solid #d0d7f0;border-radius:10px;padding:16px;font-family:monospace;font-size:13px}
.demo-box .demo-title{font-weight:700;color:#4F46E5;margin-bottom:8px;font-size:14px}
.demo-row{padding:3px 0;color:#555}
.demo-row strong{color:#333}
.alert-danger{background:rgba(230,55,87,.1);color:#e63757;padding:10px 14px;border-radius:8px;margin-bottom:14px;font-size:13px}
@media(max-width:768px){.auth-left{display:none}.auth-right{min-width:100%}}
</style>
</head>
<body>
<div class="auth-wrap">
<div class="auth-left">
<h1>{{.AppName}}</h1>
<p>WhatsApp Marketing Platform</p>
<div class="features">
<div class="feat"><i class="la la-comments"></i> Live Chat real-time SSE</div>
<div class="feat"><i class="la la-robot"></i> AI Auto Reply</div>
<div class="feat"><i class="la la-bullhorn"></i> Broadcast massal</div>
</div>
</div>
<div class="auth-right">{{if eq .Page "login"}}{{template "login" .}}{{else}}{{template "register" .}}{{end}}</div>
</div>
</body>
</html>{{end}}

{{define "login"}}
{{$d := .}}
<div class="auth-card">
<div style="display:flex;justify-content:flex-end;margin-bottom:12px">
{{range .Languages}}<a href="/lang/{{.Code}}" class="text-decoration-none mx-1"><span class="flag-icon flag-icon-{{.Flag}}" style="width:22px;height:16px;border-radius:2px" title="{{.Name}}"></span></a>{{end}}
</div>
<h2>{{T "auth_login"}}</h2>
<p class="sub">{{T "auth_no_account"}} <a href="/register">{{T "auth_signup"}}</a></p>
{{if .Flash}}<div class="alert-danger">{{.Flash}}</div>{{end}}
<form method="post" action="/login/post">
<div class="form-group"><label>{{T "auth_email"}}</label><input type="email" name="email" placeholder="{{.AppEmail}}" value="{{.AppEmail}}" required></div>
<div class="form-group"><label>{{T "auth_password"}}</label><input type="password" name="password" placeholder="••••••••" value="password" required></div>
<button type="submit" class="btn-submit">{{T "auth_signin"}}</button>
</form>
<div class="auth-divider"><span>{{T "auth_or"}}</span></div>
<div class="demo-box">
<div class="demo-title">{{T "auth_demo"}}</div>
<div class="demo-row"><strong>Admin:</strong> {{.AppEmail}} / password</div>
</div>
</div>
{{end}}

{{define "register"}}
<div class="auth-card">
<div style="display:flex;justify-content:flex-end;margin-bottom:12px">
{{range .Languages}}<a href="/lang/{{.Code}}" class="text-decoration-none mx-1"><span class="flag-icon flag-icon-{{.Flag}}" style="width:22px;height:16px;border-radius:2px" title="{{.Name}}"></span></a>{{end}}
</div>
<h2>{{T "auth_register"}}</h2>
<p class="sub">{{T "auth_has_account"}} <a href="/login">{{T "auth_signin"}}</a></p>
{{if .Flash}}<div class="alert-danger">{{.Flash}}</div>{{end}}
<form method="post" action="/register/post">
<div class="form-group"><label>{{T "auth_name"}}</label><input name="name" placeholder="Nama Anda" required></div>
<div class="form-group"><label>{{T "auth_email"}}</label><input type="email" name="email" placeholder="email@domain.com" required></div>
<div class="form-group"><label>{{T "auth_password"}}</label><input type="password" name="password" placeholder="••••••••" required></div>
<button type="submit" class="btn-submit">{{T "auth_register"}}</button>
</form>
</div>
{{end}}

{{define "home"}}{{template "layout" .}}{{end}}
{{define "content"}}
{{if eq .Page "home"}}
<div class="row">
<div class="col-12 col-sm-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_total_sent"}}</h6><span class="h2 mb-0">{{.CountSent}}</span></div>
<div class="col-auto"><span class="h2 la la-telegram la-lg text-primary mb-0"></span></div>
</div></div></div>
</div>
<div class="col-12 col-sm-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_total_recv"}}</h6><span class="h2 mb-0">{{.CountReceived}}</span></div>
<div class="col-auto"><span class="h2 la la-comment la-lg text-success mb-0"></span></div>
</div></div></div>
</div>
<div class="col-12 col-sm-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_active_wa"}}</h6><span class="h2 mb-0">{{.ActiveAccounts}}</span></div>
<div class="col-auto"><span class="h2 la la-whatsapp la-lg text-success mb-0"></span></div>
</div></div></div>
</div>
<div class="col-12 col-sm-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_unread"}}</h6><span class="h2 mb-0">{{.UnreadCount}}</span></div>
<div class="col-auto"><span class="h2 la la-envelope la-lg text-warning mb-0"></span></div>
</div></div></div>
</div>
</div>
<div class="row">
<div class="col-12 col-lg-8">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_chart_title"}}</h4><small class="text-muted">{{T "dash_chart_sub"}}</small></div>
<div class="card-body"><canvas id="msgChart" height="100"></canvas></div></div>
</div>
<div class="col-12 col-lg-4">
<div class="card mb-3"><div class="card-header"><h4 class="card-header-title">{{T "dash_wa_status"}}</h4></div>
<div class="card-body">
{{if .ConnectedAccounts}}
{{range .ConnectedAccounts}}<div class="d-flex align-items-center justify-content-between mb-2"><span><span class="status-dot" style="background:#00d97e"></span> +{{.Phone}}</span><a href="/send?to=+{{.Phone}}" class="badge bg-primary bg-opacity-10 text-primary text-decoration-none small py-1 px-2">Kirim</a></div>{{end}}
{{else}}<span class="text-muted small">Tidak ada WA terkoneksi. <a href="/wa">{{T "dash_connect"}}</a></span>{{end}}
</div></div>
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_campaigns"}}</h4></div>
<div class="card-body">
<div class="d-flex justify-content-between mb-1"><span>{{T "dash_running"}}</span><span class="badge badge-soft-warning">{{.RunningCampaigns}}</span></div>
<a href="/broadcast" class="btn btn-sm btn-white w-100">{{T "dash_manage_bc"}}</a>
</div></div>
</div>
</div>
<div class="row mt-3">
<div class="col-12 col-lg-6">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_in"}}</h4><a href="/received" class="btn btn-sm btn-white">All</a></div>
<div class="table-responsive"><table class="table table-sm table-nowrap card-table mb-0"><thead><tr><th>From</th><th>Message</th><th>Time</th></tr></thead><tbody>
{{range .Received}}<tr><td><strong>{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}</strong></td><td class="text-muted small" style="max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{.Message}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div>
</div>
<div class="col-12 col-lg-6">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_out"}}</h4><a href="/sent" class="btn btn-sm btn-white">All</a></div>
<div class="table-responsive"><table class="table table-sm table-nowrap card-table mb-0"><thead><tr><th>To</th><th>Message</th><th>Status</th></tr></thead><tbody>
{{range .Sent}}<tr><td><strong>+{{.Phone}}</strong></td><td class="text-muted small" style="max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{.Message}}</td><td><span class="badge badge-soft-success">{{.Status}}</span></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div>
</div>
</div>
<script>
new Chart(document.getElementById('msgChart'),{type:'line',data:{labels:[{{.ChartLabels}}],datasets:[{label:'Sent',data:[{{.ChartSent}}],borderColor:'#4F46E5',backgroundColor:'rgba(79,70,229,.1)',fill:true,tension:.3,pointRadius:2,pointHoverRadius:5},{label:'Received',data:[{{.ChartReceived}}],borderColor:'#10B981',backgroundColor:'rgba(16,185,129,.1)',fill:true,tension:.3,pointRadius:2,pointHoverRadius:5}]},options:{responsive:true,interaction:{intersect:false,mode:'index'},plugins:{legend:{position:'bottom'}},scales:{y:{beginAtZero:true,grid:{color:'rgba(0,0,0,.05)'}},x:{grid:{display:false}}}}})
</script>
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
            <form method="post" action="/wa/logout" onsubmit="return confirm('{{T "wa_logout_confirm"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger lift"><i class="la la-sign-out"></i></button></form>
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
<style>
.setting-tabs{display:flex;border-bottom:2px solid #e0e0e0;margin-bottom:16px}
.setting-tabs .st{background:none;border:none;padding:10px 20px;font-size:14px;font-weight:600;color:#6e788c;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-2px}
.setting-tabs .st:hover{color:#152e4d}
.setting-tabs .st.active{color:#2c7be5;border-bottom-color:#2c7be5}
.st-panel{display:none}
.st-panel.active{display:block}
</style>
<div class="setting-tabs">
<button class="st active" onclick="var p=document.querySelectorAll('.st-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('st-branding').classList.add('active');var b=this.parentElement.querySelectorAll('.st');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-paint-brush me-1"></i>{{T "set_tab_branding"}}</button>
<button class="st" onclick="var p=document.querySelectorAll('.st-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('st-messaging').classList.add('active');var b=this.parentElement.querySelectorAll('.st');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-comment me-1"></i>{{T "set_tab_messaging"}}</button>
<button class="st" onclick="var p=document.querySelectorAll('.st-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('st-system').classList.add('active');var b=this.parentElement.querySelectorAll('.st');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-cog me-1"></i>{{T "set_tab_system"}}</button>
</div>

<form method="post" action="/settings" enctype="multipart/form-data">
<div class="st-panel active" id="st-branding">
<div class="card"><div class="card-header"><h4 class="card-header-title">Branding</h4></div>
<div class="card-body">
<div class="form-group"><label>{{T "set_app_name"}}</label><input name="app_name" class="form-control" value="{{.AppName}}"></div>
<div class="form-group"><label>{{T "set_logo_upload"}}</label>
<div class="d-flex gap-2 align-items-center"><input type="file" name="logo_file" class="form-control" accept="image/*" style="flex:1"><img src="{{.AppLogo}}" onerror="this.style.display='none'" style="height:38px;border-radius:6px;border:1px solid #eee"></div>
<small class="form-text text-muted">{{T "set_logo_hint"}}: <code>{{.AppLogo}}</code></small>
</div>
<div class="form-group"><label>{{T "set_admin_email"}}</label><input name="app_email" class="form-control" value="{{.AppEmail}}"></div>
<div class="form-group"><label>{{T "set_domain"}}</label><input class="form-control" value="{{.AppURL}}" disabled><small class="form-text text-muted">{{T "set_domain_hint"}}</small></div>
</div></div>
</div>

<div class="st-panel" id="st-messaging">
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

<div class="card mt-3">
<div class="card-header"><h4 class="card-header-title"><i class="la la-clock me-1"></i> Auto-Close Idle Chats</h4></div>
<div class="card-body">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>Close After (hours, 0=disabled)</label><input type="number" name="auto_close_hours" class="form-control" value="{{.AutoCloseHours}}"></div></div>
<div class="col-md-6"><div class="form-group"><label>Follow-up Message</label><input name="auto_close_message" class="form-control" placeholder="Chat ini ditutup otomatis." value="{{.AutoCloseMessage}}"></div></div>
</div>
</div></div>
<div class="card-body">
<div class="row">
<div class="col-md-4"><div class="form-group"><label>Max Per Day (0=unlimited)</label><input type="number" name="rate_max_daily" class="form-control" value="{{.RateMaxDaily}}"></div></div>
<div class="col-md-4"><div class="form-group"><label>Random Min (detik)</label><input type="number" name="rate_random_min" class="form-control" value="{{.RateRandomMin}}"></div></div>
<div class="col-md-4"><div class="form-group"><label>Random Max (detik)</label><input type="number" name="rate_random_max" class="form-control" value="{{.RateRandomMax}}"></div></div>
</div>
<small class="form-text text-muted">Set 0 untuk unlimited. Random delay akan menggantikan interval tetap, dipilih acak antara min-max.</small>
</div></div>
</div>
</div>

<div class="st-panel" id="st-system">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "set_system_title"}}</h4></div>
<div class="card-body">
<div class="form-group"><label>{{T "set_registrations"}}</label>
<select name="registrations" class="form-control"><option value="1" {{if .Registrations}}selected{{end}}>{{T "set_enabled"}}</option><option value="0" {{if not .Registrations}}selected{{end}}>{{T "set_disabled"}}</option></select></div>
<div class="form-group"><label>{{T "set_listen_addr"}}</label><input class="form-control" value="0.0.0.0:8080" disabled><small class="form-text text-muted">Edit <code>CHATGO_ADDR</code> di <code>.env</code></small></div>
<div class="form-group"><label>MySQL Connection</label><input class="form-control" value="***" disabled><small class="form-text text-muted">Edit <code>CHATGO_MYSQL</code> di <code>.env</code></small></div>
</div></div>
</div>

<button class="btn btn-primary lift mt-3"><i class="la la-save me-1"></i> {{T "set_save"}}</button>
</form>
{{end}}

{{if eq .Page "contacts"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ct_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/contacts/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
      <div class="card mt-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-upload me-1"></i> Import CSV</h4></div>
        <div class="card-body"><form method="post" action="/contacts/import" enctype="multipart/form-data">
          <div class="form-group"><label>Upload CSV</label><input type="file" name="file" class="form-control" accept=".csv" required></div>
          <small class="form-text text-muted mb-2 d-block">Kolom: <code>name</code>, <code>phone</code>, <code>groups</code> (nama grup, koma). Group otomatis dibuat jika belum ada.</small>
          <button class="btn btn-white lift"><i class="la la-cloud-upload me-1"></i> Import</button>
        </form></div>
      </div>
    </div>
    {{if .EditID}}
    <div class="col-12 col-lg-4">
      <div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> Edit</h4></div>
        <div class="card-body"><form method="post" action="/contacts/edit">
          <input type="hidden" name="id" value="{{.EditID}}">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div>
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" value="{{.EditPhone}}" required></div>
          <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "set_save"}}</button> <a href="/contacts" class="btn btn-white ms-2">{{T "ar_cancel"}}</a>
        </form></div>
      </div>
    </div>
    {{end}}
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header d-flex justify-content-between align-items-center"><h4 class="card-header-title mb-0">{{T "nav_contacts_saved"}}</h4><div><a href="/contacts/export" class="btn btn-sm btn-white me-1"><i class="la la-download me-1"></i> Export CSV</a><button class="btn btn-sm btn-danger" onclick="bulkDeleteContacts()"><i class="la la-trash me-1"></i> Delete</button></div></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th style="width:30px"><input type="checkbox" onchange="toggleAll(this)"></th><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_from"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Contacts}}<tr><td><input type="checkbox" name="cid" value="{{.ID}}" class="contact-check"></td><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Phone}}</td><td>
            <a class="btn btn-sm btn-white" href="/contacts?edit={{.ID}}"><i class="la la-edit"></i></a>
            <form method="post" action="/contacts/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
  <script>
    function toggleAll(el) { document.querySelectorAll('.contact-check').forEach(c => c.checked = el.checked) }
    function bulkDeleteContacts() {
      var ids = []; document.querySelectorAll('.contact-check:checked').forEach(c => ids.push(c.value));
      if (ids.length === 0) { alert('Pilih kontak dulu'); return; }
      if (!confirm('Hapus ' + ids.length + ' kontak?')) return;
      var f = document.createElement('form'); f.method = 'POST'; f.action = '/contacts/bulk-delete';
      ids.forEach(function(id) { var i = document.createElement('input'); i.type = 'hidden'; i.name = 'ids'; i.value = id; f.appendChild(i) });
      document.body.appendChild(f); f.submit();
    }
  </script>
{{end}}

{{if eq .Page "groups"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "grp_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/groups/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_contacts_groups"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "grp_members"}}</th><th>Lang</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Groups}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Count}}</td><td>
            <form method="post" action="/groups/language" style="display:inline"><input type="hidden" name="group_id" value="{{.ID}}"><select name="language" class="form-select form-select-sm" style="width:80px;font-size:11px" onchange="this.form.submit()"><option value="" {{if not .Language}}selected{{end}}>-</option><option value="id" {{if eq .Language "id"}}selected{{end}}>ID</option><option value="en" {{if eq .Language "en"}}selected{{end}}>EN</option></select></form>
          </td><td><form method="post" action="/groups/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "tags"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Tambah Tag</h4></div>
        <div class="card-body"><form method="post" action="/tags/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="VIP" required></div>
          <div class="form-group"><label>Warna</label><input type="color" name="color" class="form-control form-control-color" value="#2c7be5" style="height:40px;padding:4px"></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Tags</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Warna</th><th>{{T "col_name"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
          {{range .Tags}}<tr><td>{{.ID}}</td><td><span style="display:inline-block;width:20px;height:20px;border-radius:4px;background:{{.Color}}"></span></td><td>{{.Name}}</td><td><form method="post" action="/tags/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
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
        <div class="card-body"><form method="post" action="/broadcast" enctype="multipart/form-data">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
           <div class="form-group"><label>{{T "bc_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}} ({{.Count}})</option>{{end}}</select></div>
           {{if .Tags}}<div class="form-group"><label>Tags <small class="text-muted">— filter by tag</small></label><select name="tags" class="form-control" multiple>{{range .Tags}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select><small class="form-text text-muted">Pilih tag untuk kirim hanya ke kontak dengan tag tsb.</small></div>{{end}}
           <div class="form-group"><label>Nomor Langsung <small class="text-muted">— satu per baris</small></label><textarea name="numbers" class="form-control" rows="4" placeholder="628123456789&#10;628987654321&#10;..."></textarea><small class="form-text text-muted">Tempel nomor langsung (tanpa grup). <form method="post" action="/validate" style="display:inline" target="_blank"><input type="hidden" name="numbers" value="" id="validateInput"><button type="button" class="btn btn-sm btn-outline-warning" onclick="document.getElementById('validateInput').value=document.querySelector('textarea[name=numbers]').value;this.form.submit()">Validate</button></form></small></div>
           <div class="form-group"><label><i class="la la-image me-1"></i> Media <small class="text-muted">— opsional</small></label><input type="file" name="media_file" class="form-control" accept="image/*,video/*,.pdf,.doc,.docx,.xls,.xlsx"><small class="form-text text-muted">Upload gambar/video/dokumen untuk dikirim bersama pesan.</small></div>
           <div class="form-group"><label>{{T "bc_account"}}</label><div class="border rounded p-2" style="max-height:160px;overflow-y:auto">{{range .ConnectedAccounts}}{{if .Phone}}<div class="form-check"><input class="form-check-input" type="checkbox" name="account_ids" value="+{{.Phone}}" id="bc_{{.Phone}}"><label class="form-check-label small" for="bc_{{.Phone}}">+{{.Phone}}</label></div>{{end}}{{end}}{{if not .HasConnected}}<small class="text-muted">Belum ada nomor terkoneksi</small>{{end}}</div><small class="form-text text-muted">Biarkan kosong = semua nomor terhubung. Checklist = hanya nomor itu.</small></div>
           <div class="form-group"><label>Mode Pengiriman</label><div class="border rounded p-2"><div class="form-check"><input class="form-check-input" type="radio" name="send_mode" value="round_robin" id="mode_rr" checked><label class="form-check-label" for="mode_rr"><strong>Round Robin</strong> <small class="text-muted">— kirim bergantian merata ke tiap nomor</small></label></div><div class="form-check mt-1"><input class="form-check-input" type="radio" name="send_mode" value="random" id="mode_rand"><label class="form-check-label" for="mode_rand"><strong>Random</strong> <small class="text-muted">— kirim acak ke nomor manapun</small></label></div></div></div>
           <div class="form-group"><label>Interval (detik) <small class="text-muted">jeda antar pesan</small></label><input name="interval" type="number" class="form-control" value="300" min="30" placeholder="300-400"></div>
           {{if .MetaAccounts}}
           <div class="form-group"><label><i class="la la-cloud me-1"></i> Meta API <small class="text-muted">— kirim lewat Cloud API</small></label><select name="meta_account_id" class="form-control" onchange="toggleMetaTemplate(this)"><option value="0">-- Tidak pakai Meta --</option>{{range .MetaAccounts}}<option value="{{.ID}}">{{.Name}} ({{.PhoneNumberID}})</option>{{end}}</select></div>
            <div class="form-group" id="metaTemplateGroup" style="display:none"><label>Meta Template <small class="text-muted">— opsional</small></label><select name="meta_template" class="form-control"><option value="">-- Plain text --</option>{{range .MetaTemplates}}<option value="{{.Name}}">[Meta] {{.Name}} ({{.Language}})</option>{{end}}</select><small class="form-text text-muted">Jika dipilih, template akan dipakai. Variabel dari pesan akan masuk ke parameter.</small></div>
           <script>function toggleMetaTemplate(el){document.getElementById('metaTemplateGroup').style.display=el.value!=='0'?'block':'none'}</script>
           {{end}}
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea><small class="form-text text-muted">{{T "set_vars_hint"}}</small></div>
          <button class="btn btn-primary lift" {{if not .HasConnected}}disabled{{end}}><i class="la la-bullhorn me-1"></i> {{T "bc_start"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_broadcast"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "bc_progress"}}</th><th>{{T "col_status"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
           {{range .Campaigns}}<tr><td>{{.ID}}</td><td>{{.Name}}{{if .MetaAccountID}} <span class="badge badge-soft-primary" style="font-size:9px">Meta</span>{{end}}</td><td><a href="/broadcast/detail?id={{.ID}}" title="Lihat detail nomor terkirim">{{.Sent}}/{{.Total}}</a></td><td>{{if eq .Status "running"}}<span class="badge badge-soft-primary">running</span>{{else if eq .Status "paused"}}<span class="badge badge-soft-warning">paused</span>{{else if eq .Status "done"}}<span class="badge badge-soft-success">done</span>{{else}}<span class="badge badge-soft-secondary">{{.Status}}</span>{{end}}</td><td class="text-nowrap">
             {{if eq .Status "running"}}<form method="post" action="/broadcast/pause" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-warning" title="Pause"><i class="la la-pause"></i></button></form>{{end}}
             {{if eq .Status "paused"}}<form method="post" action="/broadcast/pause" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-success" title="Resume"><i class="la la-play"></i></button></form>{{end}}
             {{if eq .Status "done"}}<form method="post" action="/broadcast/retry" style="display:inline" onsubmit="return confirm('Jalankan ulang campaign ini?')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-info" title="Retry"><i class="la la-redo"></i></button></form>{{end}}
             {{if eq .Status "stopped"}}<form method="post" action="/broadcast/retry" style="display:inline" onsubmit="return confirm('Jalankan ulang campaign ini?')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-info" title="Retry"><i class="la la-redo"></i></button></form>{{end}}
             <form method="post" action="/broadcast/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form>
           </td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "drips"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> New Drip</h4></div>
        <div class="card-body"><form method="post" action="/drips/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="Welcome Series" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      {{range .Drips}}
      <div class="card mb-3">
        <div class="card-header d-flex justify-content-between align-items-center">
          <div><h4 class="card-header-title mb-0">{{.Name}}</h4><small class="text-muted">{{len .Steps}} steps &middot; {{if eq .Status "active"}}<span class="text-success">Active</span>{{else}}<span class="text-muted">Inactive</span>{{end}}</small></div>
          <div>
            <form method="post" action="/drips/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}">{{if eq .Status "active"}}<button class="btn btn-sm btn-warning">Pause</button>{{else}}<button class="btn btn-sm btn-success">Resume</button>{{end}}</form>
            <form method="post" action="/drips/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form>
          </div>
        </div>
        <div class="table-responsive"><table class="table table-sm card-table mb-0"><thead><tr><th>#</th><th>Delay</th><th>{{T "col_message"}}</th><th></th></tr></thead><tbody>
          {{range $i, $s := .Steps}}<tr><td>{{add $i 1}}</td><td>{{if eq $i 0}}Instant{{else}}{{$s.DelayMinutes}} min{{end}}</td><td>{{$s.Message}}</td><td><form method="post" action="/drips/step/delete" style="display:inline"><input type="hidden" name="id" value="{{$s.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-times"></i></button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">No steps yet</td></tr>{{end}}
        </tbody></table></div>
        <div class="card-body border-top"><form method="post" action="/drips/step/add" class="row g-2">
          <input type="hidden" name="drip_id" value="{{.ID}}">
          <input type="hidden" name="sort_order" value="{{len .Steps}}">
          <div class="col-md-2"><input type="number" name="delay" class="form-control form-control-sm" placeholder="Min" value="0"></div>
          <div class="col-md-7"><input type="text" name="message" class="form-control form-control-sm" placeholder="Pesan..." required></div>
          <div class="col-md-3"><button class="btn btn-sm btn-primary w-100"><i class="la la-plus"></i> Add Step</button></div>
        </form></div>
      </div>
      {{else}}
      <div class="card"><div class="card-body text-center text-muted py-5">Belum ada drip campaign. Buat yang pertama!</div></div>
      {{end}}
    </div>
  </div>
{{end}}

{{if eq .Page "scheduled"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "sch_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/scheduled">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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

{{if eq .Page "canned"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Canned Response</h4></div>
        <div class="card-body"><form method="post" action="/canned/add">
          <div class="form-group"><label>Shortcut <small class="text-muted">— ketik di inbox</small></label><input name="shortcut" class="form-control" placeholder="/salam"></div>
          <div class="form-group"><label>Judul</label><input name="name" class="form-control" placeholder="Salam Pembuka" required></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Canned Responses</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Shortcut</th><th>Nama</th><th>{{T "col_message"}}</th><th></th></tr></thead><tbody>
          {{range .Canned}}<tr><td><code>{{.Shortcut}}</code></td><td>{{.Name}}</td><td style="max-width:250px">{{.Message}}</td><td><form method="post" action="/canned/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "tracker"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-link me-1"></i> Link Clicks</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>URL</th><th>Campaign</th><th>Phone</th><th>Clicked</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .LClicks}}<tr><td>{{.ID}}</td><td style="max-width:300px;word-break:break-all"><a href="/track/{{.Token}}" target="_blank">{{.URL}}</a></td><td>{{if .CampaignID}}#{{.CampaignID}}{{else}}-{{end}}</td><td>{{.Phone}}</td><td>{{if .Clicked}}<span class="badge badge-soft-success">Yes</span>{{else}}<span class="badge badge-soft-secondary">No</span>{{end}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">Belum ada link yang di-track.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "abtests"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title"><i class="la la-balance-scale me-1"></i> A/B Test Results</h4><button class="btn btn-sm btn-primary" data-bs-toggle="modal" data-bs-target="#abModal"><i class="la la-plus"></i> New Test</button></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Campaign</th><th>Variant A</th><th>Variant B</th><th>A Sent</th><th>B Sent</th></tr></thead><tbody>
      {{range .ABTests}}<tr><td>{{.ID}}</td><td>#{{.CampaignID}}</td><td style="max-width:200px">{{.VariantA}}</td><td style="max-width:200px">{{.VariantB}}</td><td>{{.ASent}}</td><td>{{.BSent}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">Belum ada A/B test.</td></tr>{{end}}
    </tbody></table></div>
  </div>
  <div class="modal fade" id="abModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><form method="post" action="/ab-tests/add">
    <div class="modal-header"><h5>New A/B Test</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>
    <div class="modal-body">
      <div class="form-group"><label>Campaign ID</label><input type="number" name="campaign_id" class="form-control" required></div>
      <div class="form-group"><label>Variant A</label><textarea name="variant_a" class="form-control" rows="2" required></textarea></div>
      <div class="form-group"><label>Variant B</label><textarea name="variant_b" class="form-control" rows="2" required></textarea></div>
    </div>
    <div class="modal-footer"><button class="btn btn-primary">Create</button></div>
  </form></div></div></div>
{{end}}

{{if eq .Page "store"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card mb-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Product</h4></div>
        <div class="card-body"><form method="post" action="/store/add">
          <div class="form-group"><label>Nama</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>Deskripsi</label><textarea name="desc" class="form-control" rows="2"></textarea></div>
          <div class="form-group"><label>Harga</label><input name="price" type="number" step="0.01" class="form-control" required></div>
          <div class="form-group"><label>Gambar URL</label><input name="image_url" class="form-control"></div>
          <div class="form-group"><label>Kategori</label><select name="category" class="form-control">{{range .Categories}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div>
          <div class="form-group"><label>Stok</label><input name="stock" type="number" class="form-control" value="0"></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> Add</button>
        </form></div>
      </div>
      <div class="card"><div class="card-header"><h4 class="card-header-title">Kategori</h4></div>
        <div class="card-body"><form method="post" action="/store/category/add"><div class="input-group"><input name="name" class="form-control" placeholder="Nama kategori"><button class="btn btn-primary">Add</button></div></form>
          <div class="mt-2">{{range .Categories}}<span class="badge badge-soft-primary me-1 mb-1">{{.Name}} <form method="post" action="/store/category/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button style="background:none;border:none;color:inherit;cursor:pointer;font-size:10px">&times;</button></form></span>{{end}}</div>
        </div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Products</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Image</th><th>Nama</th><th>Price</th><th>Category</th><th>Stock</th><th></th></tr></thead><tbody>
          {{range .Products}}<tr><td>{{.ID}}</td><td>{{if .ImageURL}}<img src="{{.ImageURL}}" style="width:40px;height:40px;object-fit:cover;border-radius:6px">{{else}}-{{end}}</td><td>{{.Name}}</td><td>{{.Price}}</td><td>{{.Category}}</td><td>{{.Stock}}</td><td><form method="post" action="/store/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form></td></tr>{{else}}<tr><td colspan="7" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "orders"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title">Orders</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Phone</th><th>Name</th><th>Product</th><th>Qty</th><th>Total</th><th>Status</th><th></th></tr></thead><tbody>
      {{range .Orders}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Name}}</td><td>#{{.ProductID}}</td><td>{{.Quantity}}</td><td>{{.Total}}</td><td><span class="badge badge-soft-{{if eq .Status "new"}}warning{{else if eq .Status "paid"}}success{{else}}secondary{{end}}">{{.Status}}</span></td><td>
        <form method="post" action="/store/orders/update" style="display:inline" class="d-flex gap-1"><input type="hidden" name="id" value="{{.ID}}"><select name="status" class="form-select form-select-sm" style="width:auto" onchange="this.form.submit()"><option value="new" selected>New</option><option value="confirmed">Confirmed</option><option value="paid">Paid</option><option value="shipped">Shipped</option><option value="cancelled">Cancelled</option></select></form>
      </td></tr>{{else}}<tr><td colspan="8" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "forms"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Create Form</h4></div>
        <div class="card-body"><form method="post" action="/forms/add">
          <div class="form-group"><label>Nama Form</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>Fields (JSON)</label><textarea name="fields" class="form-control" rows="6" placeholder='[{"label":"Nama","type":"text"},{"label":"Email","type":"text"},{"label":"Rating","type":"number"}]'></textarea><small class="form-text text-muted">Array of {label, type}. Type: text, number, email, textarea.</small></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> Create</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Forms</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Fields</th><th></th></tr></thead><tbody>
          {{range .Forms}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><code>{{.Fields}}</code></td><td><a href="/forms/submissions?form_id={{.ID}}" class="btn btn-sm btn-white">Data</a> <form method="post" action="/forms/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "submissions"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">Form Submissions</h4>
    <select class="form-select form-select-sm" style="width:auto" onchange="window.location='?form_id='+this.value"><option value="">Pilih Form</option>{{range .Forms}}<option value="{{.ID}}" {{if eq .ID $.QueryFormID}}selected{{end}}>{{.Name}}</option>{{end}}</select>
  </div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Phone</th><th>Data</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Submissions}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td><code>{{.Data}}</code></td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "reminders"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Reminder</h4></div>
        <div class="card-body"><form method="post" action="/reminders/add">
          <div class="form-group"><label>Nama</label><input name="name" class="form-control"></div>
          <div class="form-group"><label>Phone</label><input name="phone" class="form-control" required></div>
          <div class="form-group"><label>Amount</label><input name="amount" type="number" step="0.01" class="form-control" required></div>
          <div class="form-group"><label>Due Date</label><input type="date" name="due_date" class="form-control" required></div>
          <div class="form-group"><label>Pesan</label><textarea name="message" class="form-control" rows="2">Pengingat: tagihan sebesar {amount} jatuh tempo {date}.</textarea></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> Add</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Payment Reminders</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Phone</th><th>Name</th><th>Amount</th><th>Due</th><th>Status</th></tr></thead><tbody>
          {{range .Reminders}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Name}}</td><td>{{.Amount}}</td><td>{{.DueDate}}</td><td>{{.Status}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "analytics"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-chart-pie me-1"></i> Agent Performance (30 Hari)</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Agent</th><th>Chats</th><th>Replies</th><th>Avg Response (s)</th></tr></thead><tbody>
      {{range .AgentMetrics}}<tr><td>{{.AgentName}}</td><td>{{.Chats}}</td><td>{{.Replied}}</td><td>{{printf "%.0f" .AvgTime}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center py-4">No data yet.</td></tr>{{end}}
    </tbody></table></div>
  </div>
  <div class="card mt-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-star me-1"></i> CSAT Score</h4></div>
    <div class="card-body text-center">
      <div class="display-3 fw-bold text-warning">{{printf "%.1f" .CSATAvg}} ⭐</div>
      <p class="text-muted">dari {{.CSATCount}} penilaian (30 hari)</p>
    </div>
  </div>
{{end}}

{{if eq .Page "depts"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Department</h4></div>
        <div class="card-body"><form method="post" action="/depts/add">
          <div class="form-group"><label>Name</label><input name="name" class="form-control" placeholder="Sales" required></div>
          <div class="form-group"><label>Agents</label><select name="agents" class="form-control" multiple>{{range .Users}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> Add</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Departments</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Agents</th><th></th></tr></thead><tbody>
          {{range .Depts}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Agents}}</td><td><form method="post" action="/depts/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "recurring"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> New Recurring</h4></div>
        <div class="card-body"><form method="post" action="/recurring/add">
          <div class="form-group"><label>Name</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>Groups</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}} ({{.Count}})</option>{{end}}</select></div>
          <div class="form-group"><label>Message</label><textarea name="message" class="form-control" rows="3" required></textarea></div>
          <div class="row">
            <div class="col-6"><label>Day</label><select name="day_of_week" class="form-control"><option value="0">Daily</option><option value="1">Monday</option><option value="2">Tuesday</option><option value="3">Wednesday</option><option value="4">Thursday</option><option value="5">Friday</option><option value="6">Saturday</option><option value="7">Sunday</option></select></div>
            <div class="col-6"><label>Hour (0-23)</label><input type="number" name="hour" class="form-control" value="9" min="0" max="23"></div>
          </div>
          <button class="btn btn-primary mt-2"><i class="la la-plus"></i> Create</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Recurring Campaigns</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Groups</th><th>Schedule</th><th>Status</th><th></th></tr></thead><tbody>
          {{range .Recurrings}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Groups}}</td><td>{{if eq .DayOfWeek 0}}Daily{{else}}Day {{.DayOfWeek}}{{end}} @ {{.Hour}}:00</td><td>{{.Status}}</td><td>
            <form method="post" action="/recurring/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{if eq .Status "active"}}Pause{{else}}Activate{{end}}</button></form>
            <form method="post" action="/recurring/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form>
          </td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "uploads"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-folder-open me-1"></i> Uploaded Files</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>File</th><th>URL</th></tr></thead><tbody>
      {{range $i, $f := .Files}}<tr><td>{{add $i 1}}</td><td>{{$f}}</td><td><code>/public/uploads/{{$f}}</code></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center py-4">No files uploaded yet.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "blacklist"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add to Blacklist</h4></div>
        <div class="card-body"><form method="post" action="/blacklist/add">
          <div class="form-group"><label>Phone</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <div class="form-group"><label>Reason</label><input name="reason" class="form-control" placeholder="spam / abuse"></div>
          <button class="btn btn-danger"><i class="la la-ban me-1"></i> Block</button>
        </form></div>
      </div>
      <div class="card mt-3"><div class="card-header"><h4 class="card-header-title">Validate Numbers</h4></div>
        <div class="card-body"><form method="post" action="/validate">
          <textarea name="numbers" class="form-control" rows="6" placeholder="628xxx&#10;628xxx" required></textarea>
          <button class="btn btn-primary mt-2"><i class="la la-check me-1"></i> Validate</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Blocked Numbers</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Phone</th><th>Reason</th><th>{{T "col_time"}}</th><th></th></tr></thead><tbody>
          {{range .Blacklist}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Reason}}</td><td class="text-muted small">{{.Created}}</td><td><form method="post" action="/blacklist/remove"><input type="hidden" name="phone" value="{{.Phone}}"><button class="btn btn-sm btn-success">Unblock</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center py-4">Blacklist kosong.</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "csat"}}
  <div class="row">
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-warning">{{printf "%.1f" .CSATAvg}}</h1><p>Average Rating (30d)</p></div></div>
    </div>
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-primary">{{.CSATCount}}</h1><p>Total Responses</p></div></div>
    </div>
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-success">{{printf "%.0f" (mult .CSATAvg 20)}}%</h1><p>Satisfaction Score</p></div></div>
    </div>
  </div>
{{end}}

{{if eq .Page "customers"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">Customer Directory</h4><input type="text" id="custSearch" class="form-control form-control-sm" placeholder="Search phone..." style="width:250px" oninput="var q=this.value.toLowerCase();document.querySelectorAll('.cust-row').forEach(r=>r.style.display=r.textContent.toLowerCase().includes(q)?'':'none')"></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Phone</th><th>Name</th><th>Orders</th><th>Last Active</th><th></th></tr></thead><tbody>
      {{range .Contacts}}<tr class="cust-row"><td>{{.Phone}}</td><td>{{.Name}}</td><td>-</td><td>-</td><td><a href="/inbox/chat?phone={{.Phone}}" class="btn btn-sm btn-primary">Chat</a></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "calendar"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-calendar me-1"></i> Campaign Calendar</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Date</th><th>Title</th><th>Type</th></tr></thead><tbody>
      {{range .CalEvents}}<tr><td>{{.Date}}</td><td>{{.Title}}</td><td><span class="badge badge-soft-{{if eq .Type "Campaign"}}primary{{else if eq .Type "Recurring"}}success{{else}}warning{{end}}">{{.Type}}</span></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center py-4">No scheduled events.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "backup"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-database me-1"></i> Database Backup</h4></div>
    <div class="card-body text-center py-5">
      <form method="post" action="/backup">
        <button class="btn btn-primary btn-lg"><i class="la la-download me-1"></i> Backup Database Now</button>
      </form>
      <p class="text-muted mt-2">Backup akan disimpan ke <code>public/backups/</code></p>
    </div>
  </div>
{{end}}

{{if eq .Page "macros"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Macro</h4></div>
        <div class="card-body"><form method="post" action="/macros/add">
          <div class="form-group"><label>Name</label><input name="name" class="form-control" placeholder="Quick Resolve" required></div>
          <div class="form-group"><label>Actions</label><textarea name="actions" class="form-control" rows="4" placeholder="assign:1;tag:resolved;reply:Terima kasih!;close" required></textarea><small class="form-text text-muted">Format: action:value;action:value. Actions: assign, tag, reply, close</small></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> Create</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Macros</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Actions</th><th></th></tr></thead><tbody>
          {{range .Macros}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><code>{{.Actions}}</code></td><td><form method="post" action="/macros/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "merge"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-code-branch me-1"></i> Duplicate Contacts</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>Phone</th><th>Count</th><th>Names</th><th>IDs</th><th></th></tr></thead><tbody>
      {{range .Duplicates}}<tr><td>{{index . "phone"}}</td><td>{{index . "cnt"}}</td><td>{{index . "names"}}</td><td>{{index . "ids"}}</td><td><form method="post" action="/merge/execute"><input type="hidden" name="keep_id" value="{{index (split (index . "ids") ",") 0}}">{{range $i, $id := split (index . "ids") ","}}{{if gt $i 0}}<input type="hidden" name="merge_ids" value="{{$id}}">{{end}}{{end}}<button class="btn btn-sm btn-warning">Merge</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center py-4">No duplicates found.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "audit"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-history me-1"></i> Audit Log</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>User</th><th>Action</th><th>Detail</th><th>IP</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .AuditLogs}}<tr><td>{{.ID}}</td><td>#{{.UserID}}</td><td>{{.Action}}</td><td>{{.Detail}}</td><td>{{.IP}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">No audit entries yet.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "translatetool"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-language me-1"></i> Auto Translate</h4></div>
    <div class="card-body">
      <div class="row">
        <div class="col-md-5"><div class="form-group"><label>Source Text</label><textarea id="srcText" class="form-control" rows="5" placeholder="Masukkan teks..."></textarea></div></div>
        <div class="col-md-2 d-flex align-items-end pb-3"><button onclick="doTranslate()" class="btn btn-primary w-100"><i class="la la-sync me-1"></i> Translate</button></div>
        <div class="col-md-5"><div class="form-group"><label>Result <select id="langTo" class="form-select form-select-sm" style="width:auto;display:inline"><option value="id">ID</option><option value="en">EN</option><option value="es">ES</option><option value="fr">FR</option><option value="de">DE</option><option value="zh">ZH</option><option value="ja">JA</option><option value="ko">KO</option><option value="ar">AR</option></select></label><textarea id="resText" class="form-control" rows="5" readonly></textarea></div>
      </div>
    </div>
  </div>
  <script>
  function doTranslate(){
    var t=document.getElementById('srcText').value;
    var to=document.getElementById('langTo').value;
    fetch('/translate',{method:'POST',headers:{'Content-Type':'application/x-www-form-urlencoded'},body:'text='+encodeURIComponent(t)+'&to='+to})
    .then(r=>r.text()).then(r=>document.getElementById('resText').value=r);
  }
  </script>
{{end}}

{{if eq .Page "widgetinfo"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-code me-1"></i> Web Chat Widget</h4></div>
    <div class="card-body">
      <p>Sisipkan script ini di HTML website kamu untuk menampilkan chat widget.</p>
      <pre class="bg-light p-3 rounded"><code>&lt;script src="{{.AppURL}}/widget.js"&gt;&lt;/script&gt;</code></pre>
      <p class="text-muted small">Widget akan muncul di pojok kanan bawah website. Pengunjung bisa chat langsung ke WhatsApp kamu.</p>
      <p class="text-muted small"><strong>Webhook Email:</strong> POST ke <code>{{.AppURL}}/email-webhook</code> dengan field <code>from</code>, <code>subject</code>, <code>text</code> untuk forward email ke WA inbox.</p>
    </div>
  </div>
{{end}}

{{if eq .Page "emailwa"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-envelope me-1"></i> Email → WhatsApp Gateway</h4></div>
    <div class="card-body">
      <p>Forward email ke WhatsApp inbox via webhook.</p>
      <h5 class="mt-3">Setup</h5>
      <pre class="bg-light p-3 rounded"><code>POST {{.AppURL}}/email-webhook
Content-Type: application/x-www-form-urlencoded

from=sender@email.com&subject=Judul Email&text=Isi email</code></pre>
      <p class="text-muted small">Bisa diintegrasikan dengan Zapier, Make, n8n, atau custom webhook dari email provider.</p>
    </div>
  </div>
{{end}}

{{if eq .Page "subscribe"}}
  <div class="row">
    <div class="col-12"><h2 class="mb-4">Pilih Paket Langganan</h2></div>
    {{range .Packages}}
    {{$pkgID := .ID}}
    <div class="col-12 col-md-6 col-lg-4 mb-4">
      <div class="card border-0 shadow-sm h-100" style="border-radius:14px">
        <div class="card-body text-center p-4">
          <h4 class="fw-bold">{{.Name}}</h4>
          <div class="display-4 fw-bold text-primary my-3">{{.Price}}</div>
          <p class="text-muted small">
            Send: {{.SendLimit}} | Device: {{.DeviceLimit}} | WA: {{.WaAccountLimit}}<br>
            Contact: {{.ContactLimit}} | AI: {{if .KeyLimit}}{{.KeyLimit}}{{else}}0{{end}}
          </p>
          {{range $.PaymentGateways}}
          {{if eq .Status "active"}}
          <form method="post" action="/subscribe/checkout" class="mt-2">
            <input type="hidden" name="package_id" value="{{$pkgID}}">
            <input type="hidden" name="gateway_id" value="{{.ID}}">
            <button class="btn btn-primary w-100"><i class="la la-credit-card me-1"></i> Bayar via {{.Name}}</button>
          </form>
          {{end}}
          {{end}}
          {{if not $.PaymentGateways}}<p class="text-muted small mt-2">Belum ada gateway pembayaran.</p>{{end}}
        </div>
      </div>
    </div>
    {{else}}
    <div class="col-12 text-center py-5 text-muted">Belum ada paket tersedia.</div>
    {{end}}
  </div>
{{end}}

{{if eq .Page "admin_paygateways"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> Add Gateway</h4></div>
        <div class="card-body"><form method="post" action="/admin/gateways-pay/add">
          <div class="form-group"><label>Provider</label><select name="provider" class="form-control"><option value="midtrans">Midtrans (ID)</option><option value="xendit">Xendit (ID)</option><option value="paypal">PayPal (Intl)</option><option value="stripe">Stripe (Intl)</option></select></div>
          <div class="form-group"><label>Nama</label><input name="name" class="form-control" placeholder="My Gateway"></div>
          <div class="form-group"><label>API Key</label><input name="api_key" class="form-control"></div>
          <div class="form-group"><label>API Secret</label><input name="api_secret" class="form-control"></div>
          <div class="form-group"><label>Webhook Secret</label><input name="webhook_secret" class="form-control"></div>
          <div class="form-group"><label>Currency</label><select name="currency" class="form-control"><option value="IDR">IDR</option><option value="USD">USD</option><option value="EUR">EUR</option><option value="SGD">SGD</option></select></div>
          <div class="form-group"><label>Base URL (optional)</label><input name="base_url" class="form-control" placeholder="Kosongkan untuk default"></div>
          <button class="btn btn-primary"><i class="la la-plus me-1"></i> Add</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">Payment Gateways</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Provider</th><th>Currency</th><th>Status</th><th></th></tr></thead><tbody>
          {{range .PaymentGateways}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Provider}}</td><td>{{.Currency}}</td><td>{{if eq .Status "active"}}<span class="badge badge-soft-success">Active</span>{{else}}<span class="badge badge-soft-secondary">Inactive</span>{{end}}</td><td>
            <form method="post" action="/admin/gateways-pay/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{if eq .Status "active"}}Disable{{else}}Enable{{end}}</button></form>
            <form method="post" action="/admin/gateways-pay/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Del</button></form>
          </td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "admin_transactions_pay"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-receipt me-1"></i> Payment Transactions</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Invoice</th><th>User</th><th>Package</th><th>Amount</th><th>Status</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Txs}}<tr><td>{{.ID}}</td><td><code>{{.InvoiceID}}</code></td><td>#{{.UserID}}</td><td>#{{.PackageID}}</td><td>{{.Amount}} {{.Currency}}</td><td>{{if eq .Status "paid"}}<span class="badge badge-soft-success">Paid</span>{{else if eq .Status "failed"}}<span class="badge badge-soft-danger">Failed</span>{{else}}<span class="badge badge-soft-warning">{{.Status}}</span>{{end}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="7" class="text-muted text-center py-4">No transactions yet.</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "templates"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "tpl_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/templates/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div>
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
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
    {{if .LogPages}}<div class="card-footer"><div class="pagination">{{range .LogPages}}{{if eq . $.PageNum}}<span class="active">{{.}}</span>{{else}}<a href="?page={{.}}">{{.}}</a>{{end}}{{end}}</div></div>{{end}}
  </div>
{{end}}

{{if eq .Page "hosts_android"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dev_add"}}</h4></div><div class="card-body">
      <form method="post" action="/devices/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
      <form method="post" action="/ai/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><div class="form-group"><label>Endpoint</label><input name="endpoint" class="form-control" placeholder="https://..."></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_plugins"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Endpoint</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .AiPlugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Endpoint}}</td><td><form method="post" action="/ai/plugins/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin"}}
<div class="row">
<div class="col-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">Total Users</h6><span class="h2 mb-0">{{.TotalUsers}}</span></div>
<div class="col-auto"><span class="h2 la la-users la-lg text-primary mb-0"></span></div>
</div></div></div>
</div>
<div class="col-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_active_wa"}}</h6><span class="h2 mb-0">{{.ActiveAccounts}}</span></div>
<div class="col-auto"><span class="h2 la la-whatsapp la-lg text-success mb-0"></span></div>
</div></div></div>
</div>
<div class="col-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">Campaigns</h6><span class="h2 mb-0">{{.RunningCampaigns}}</span></div>
<div class="col-auto"><span class="h2 la la-bullhorn la-lg text-warning mb-0"></span></div>
</div></div></div>
</div>
<div class="col-6 col-xl-3">
<div class="card"><div class="card-body"><div class="row align-items-center">
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "dash_total_sent"}}</h6><span class="h2 mb-0">{{.CountSent}}</span></div>
<div class="col-auto"><span class="h2 la la-telegram la-lg text-info mb-0"></span></div>
</div></div></div>
</div>
</div>
<div class="row mt-3">
<div class="col-12"><div class="card"><div class="card-header"><h4 class="card-header-title">System Overview</h4></div>
<div class="card-body"><canvas id="adminChart" height="80"></canvas></div></div></div>
</div>
<div class="row mt-3">
<div class="col-6 col-xl-3"><a href="/admin/users" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-users la-2x text-primary mb-2 d-block"></i><strong>Users</strong><br><small class="text-muted">Manage Users</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/packages" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-box la-2x text-success mb-2 d-block"></i><strong>Packages</strong><br><small class="text-muted">Manage Packages</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/waservers" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-server la-2x text-warning mb-2 d-block"></i><strong>WA Servers</strong><br><small class="text-muted">Manage Servers</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/subscriptions" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-star la-2x text-danger mb-2 d-block"></i><strong>Subscriptions</strong><br><small class="text-muted">Manage Subs</small></div></a></div>
</div>
<script>
new Chart(document.getElementById('adminChart'),{type:'line',data:{labels:[{{.ChartLabels}}],datasets:[{label:'Sent',data:[{{.ChartSent}}],borderColor:'#4F46E5',backgroundColor:'rgba(79,70,229,.1)',fill:true,tension:.3,pointRadius:2},{label:'Received',data:[{{.ChartReceived}}],borderColor:'#10B981',backgroundColor:'rgba(16,185,129,.1)',fill:true,tension:.3,pointRadius:2}]},options:{responsive:true,plugins:{legend:{position:'bottom'}},scales:{y:{beginAtZero:true}}}})
</script>
{{end}}

{{if eq .Page "admin_users"}}
<div class="row">
<div class="col-12"><div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">{{T "adm_users"}}</h4><button class="btn btn-primary btn-sm lift" onclick="document.getElementById('addUserForm').style.display=document.getElementById('addUserForm').style.display==='none'?'block':'none'"><i class="la la-plus me-1"></i> {{T "usr_add"}}</button></div>
<div id="addUserForm" style="display:none;border-bottom:1px solid #eee;padding:16px"><form method="post" action="/admin/users/add">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div></div>
<div class="col-md-6"><div class="form-group"><label>Email</label><input name="email" type="email" class="form-control" required></div></div>
<div class="col-md-6"><div class="form-group"><label>Password</label><input name="password" type="password" class="form-control"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "usr_role"}}</label><select name="role" class="form-control">{{range .Roles}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div></div>
<div class="col-12"><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></div>
</div></form></div>
<div class="table-responsive"><table class="table table-sm card-table mb-0"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Email</th><th>{{T "usr_role"}}</th><th>Registered</th><th>{{T "col_action"}}</th></tr></thead><tbody>
{{range .Users}}<tr>
<td>{{.ID}}</td>
<td>{{.Name}}</td>
<td>{{.Email}}</td>
<td><span class="badge {{if eq .Role "admin"}}badge-soft-primary{{else}}badge-soft-secondary{{end}}">{{.Role}}</span></td>
<td class="text-muted small">{{.Created}}</td>
<td class="text-nowrap">
<a class="btn btn-sm btn-white" href="/admin/users?edit={{.ID}}"><i class="la la-edit"></i></a>
<a class="btn btn-sm btn-warning" href="/admin/users/impersonate?id={{.ID}}" title="Impersonate"><i class="la la-user-circle"></i></a>
<form method="post" action="/admin/users/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form>
</td>
</tr>{{else}}<tr><td colspan="7" class="text-muted text-center py-4">-</td></tr>{{end}}
</tbody></table></div></div></div>
{{if .EditID}}
<div class="col-12 mt-3"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> Edit User #{{.EditID}}</h4></div>
<div class="card-body"><form method="post" action="/admin/users/edit">
<input type="hidden" name="id" value="{{.EditID}}">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div></div>
<div class="col-md-6"><div class="form-group"><label>Email</label><input name="email" type="email" class="form-control" value="{{.EditPhone}}"></div></div>
<div class="col-md-6"><div class="form-group"><label>Password (biarkan kosong)</label><input name="password" type="password" class="form-control" placeholder="••••••"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "usr_role"}}</label><select name="role" class="form-control">{{range .Roles}}<option value="{{.Name}}" {{if eq .Name $.EditRole}}selected{{end}}>{{.Name}}</option>{{end}}</select></div></div>
</div>
<button class="btn btn-warning lift"><i class="la la-save me-1"></i> Update</button> <a href="/admin/users" class="btn btn-white ms-2">{{T "btn_cancel"}}</a>
</form></div></div></div>
{{end}}
</div>
{{end}}

{{if eq .Page "admin_roles"}}
  <div class="row">
    {{if .EditID}}
    <div class="col-12 col-lg-4"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> Edit Role #{{.EditID}}</h4></div><div class="card-body"><form method="post" action="/admin/roles/edit"><input type="hidden" name="id" value="{{.EditID}}"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div><div class="form-group"><label>{{T "role_perms"}}</label><select name="permissions" class="form-control" multiple size="18" style="overflow-y:auto;min-height:360px"><option value="manage_users">Users</option><option value="manage_roles">Roles</option><option value="manage_packages">Packages</option><option value="manage_vouchers">Vouchers</option><option value="manage_subscriptions">Subscriptions</option><option value="manage_transactions">Transactions</option><option value="manage_payouts">Payouts</option><option value="manage_pages">Pages</option><option value="manage_marketing">Marketing</option><option value="manage_languages">Languages</option><option value="manage_waservers">WA Servers</option><option value="manage_gateways">Gateways</option><option value="manage_shorteners">Shorteners</option><option value="manage_plugins">Plugins</option><option value="manage_meta">Meta API</option><option value="manage_metatemplates">Meta Templates</option><option value="wa_send">Send Message</option><option value="wa_broadcast">Broadcast</option><option value="wa_scheduled">Scheduled</option><option value="wa_sent">Sent Messages</option><option value="wa_received">Received Messages</option><option value="wa_inbox">Live Chat</option><option value="wa_status">WA Status</option><option value="wa_autoreply">Auto Reply</option><option value="wa_ai_keys">AI Keys</option><option value="wa_ai_plugins">AI Plugins</option><option value="wa_knowledge">Knowledge Base</option><option value="wa_contacts">Contacts</option><option value="wa_groups">Contact Groups</option><option value="wa_unsub">Unsubscribed</option><option value="wa_templates">Templates</option><option value="wa_apikeys">API Keys</option><option value="wa_webhooks">Webhooks</option><option value="wa_logger">Logger</option><option value="wa_settings">Settings</option><option value="wa_docs">Documentation</option><option value="wa_hosts">Hosts</option><option value="wa_ussd">USSD</option><option value="wa_impersonate">Impersonate</option><option value="wa_drips">Drip Campaign</option><option value="wa_tags">Tags</option><option value="wa_canned">Canned Responses</option><option value="wa_recurring">Recurring Campaigns</option><option value="wa_store">Store Products</option><option value="wa_orders">Orders</option><option value="wa_forms">Forms</option><option value="wa_reminders">Payment Reminders</option><option value="wa_analytics">Analytics</option><option value="wa_blacklist">Blacklist</option><option value="wa_csat">CSAT Survey</option><option value="wa_depts">Departments</option><option value="wa_customers">Customer Directory</option><option value="wa_calendar">Calendar</option><option value="wa_macros">Inbox Macros</option><option value="wa_files">File Manager</option><option value="wa_merge">Contact Merge</option><option value="wa_translate">Auto Translate</option><option value="wa_audit">Audit Log</option><option value="wa_backup">Backup</option><option value="wa_subscribe">Pricing/Subscribe</option><option value="manage_paygateways">Payment Gateways</option><option value="manage_paytx">Payment Transactions</option></select></div><button class="btn btn-warning lift"><i class="la la-save me-1"></i> Update</button> <a href="/admin/roles" class="btn btn-white ms-2">{{T "ar_cancel"}}</a></form><script>document.addEventListener('DOMContentLoaded',function(){var s=document.querySelector('form[action=\"/admin/roles/edit\"] select[name=\"permissions\"]');if(s){var v='{{.EditContent}}';v.split(',').forEach(function(p){var o=s.querySelector('option[value=\"'+p.replace(/^\\s+|\\s+$/g,'')+'\"]');if(o)o.selected=true})}})</script></div></div></div>
    {{else}}
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "role_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/roles/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
      <div class="form-group"><label>{{T "role_perms"}}</label><select name="permissions" class="form-control" multiple size="18" style="overflow-y:auto;min-height:360px"><option value="manage_users">Users</option><option value="manage_roles">Roles</option><option value="manage_packages">Packages</option><option value="manage_vouchers">Vouchers</option><option value="manage_subscriptions">Subscriptions</option><option value="manage_transactions">Transactions</option><option value="manage_payouts">Payouts</option><option value="manage_pages">Pages</option><option value="manage_marketing">Marketing</option><option value="manage_languages">Languages</option><option value="manage_waservers">WA Servers</option><option value="manage_gateways">Gateways</option><option value="manage_shorteners">Shorteners</option><option value="manage_plugins">Plugins</option><option value="manage_meta">Meta API</option><option value="manage_metatemplates">Meta Templates</option><option value="wa_send">Send Message</option><option value="wa_broadcast">Broadcast</option><option value="wa_scheduled">Scheduled</option><option value="wa_sent">Sent Messages</option><option value="wa_received">Received Messages</option><option value="wa_inbox">Live Chat</option><option value="wa_status">WA Status</option><option value="wa_autoreply">Auto Reply</option><option value="wa_ai_keys">AI Keys</option><option value="wa_ai_plugins">AI Plugins</option><option value="wa_knowledge">Knowledge Base</option><option value="wa_contacts">Contacts</option><option value="wa_groups">Contact Groups</option><option value="wa_unsub">Unsubscribed</option><option value="wa_templates">Templates</option><option value="wa_apikeys">API Keys</option><option value="wa_webhooks">Webhooks</option><option value="wa_logger">Logger</option><option value="wa_settings">Settings</option><option value="wa_docs">Documentation</option><option value="wa_hosts">Hosts</option><option value="wa_ussd">USSD</option><option value="wa_impersonate">Impersonate</option><option value="wa_drips">Drip Campaign</option><option value="wa_tags">Tags</option><option value="wa_canned">Canned Responses</option><option value="wa_recurring">Recurring Campaigns</option><option value="wa_store">Store Products</option><option value="wa_orders">Orders</option><option value="wa_forms">Forms</option><option value="wa_reminders">Payment Reminders</option><option value="wa_analytics">Analytics</option><option value="wa_blacklist">Blacklist</option><option value="wa_csat">CSAT Survey</option><option value="wa_depts">Departments</option><option value="wa_customers">Customer Directory</option><option value="wa_calendar">Calendar</option><option value="wa_macros">Inbox Macros</option><option value="wa_files">File Manager</option><option value="wa_merge">Contact Merge</option><option value="wa_translate">Auto Translate</option><option value="wa_audit">Audit Log</option><option value="wa_backup">Backup</option><option value="wa_subscribe">Pricing/Subscribe</option><option value="manage_paygateways">Payment Gateways</option><option value="manage_paytx">Payment Transactions</option></select></div>
      <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></form></div></div></div>
    {{end}}
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_roles"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "role_perms"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Roles}}<tr class="align-middle"><td>{{.ID}}</td><td><strong>{{.Name}}</strong></td><td style="max-width:480px">{{permBadges .Permissions}}</td><td class="text-nowrap"><a class="btn btn-sm btn-white" href="/admin/roles?edit={{.ID}}"><i class="la la-edit"></i></a> <form method="post" action="/admin/roles/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center py-4">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_packages"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "pkg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/packages/add">
        <div class="form-row"><div class="form-group col-6"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
        <div class="form-group col-4"><label>Action</label><input name="action_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-4"><label>Meta</label><input name="meta_limit" type="number" class="form-control" value="0"></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_packages"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "pkg_price"}}</th><th>Limits</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Packages}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Price}}</td><td><small>S:{{.SendLimit}} D:{{.DeviceLimit}} M:{{.MetaLimit}}</small></td><td><form method="post" action="/admin/packages/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_vouchers"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "vch_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/vouchers/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
      <form method="post" action="/admin/subscriptions/add"><div class="form-group"><label>User</label><select name="user" class="form-control" required>{{range .Users}}<option value="{{.Email}}">{{.Name}} ({{.Email}})</option>{{else}}<option value="">No users</option>{{end}}</select></div><div class="form-group"><label>{{T "adm_packages"}}</label><select name="pkg" class="form-control">{{range .Packages}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div><div class="form-group"><label>{{T "sub_expire"}}</label><input name="expire" type="date" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
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
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
        <div class="form-row"><div class="form-group col-8"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
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
      <form method="post" action="/admin/shorteners/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_shorteners"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Shorteners}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><form method="post" action="/admin/shorteners/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_plugins"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "plg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><div class="form-group"><label>{{T "plg_dir"}}</label><input name="dir" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">Batal</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_plugins"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "plg_dir"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Plugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Dir}}</td><td><form method="post" action="/admin/plugins/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_meta"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">Tambah Meta Account</h4></div><div class="card-body">
      <form method="post" action="/admin/meta/add">
        <div class="form-group"><label>Nama</label><input name="name" class="form-control" placeholder="My Business" required></div>
        <div class="form-group"><label>Phone Number ID</label><input name="phone_number_id" class="form-control" placeholder="123456789..." required></div>
        <div class="form-group"><label>Access Token</label><input name="access_token" class="form-control" placeholder="EAA..." required></div>
        <div class="form-group"><label>App ID</label><input name="app_id" class="form-control" placeholder="123456..."></div>
        <div class="form-group"><label>App Secret</label><input name="app_secret" class="form-control" placeholder="abc123..."></div>
        <div class="form-group"><label>Verify Token</label><input name="verify_token" class="form-control" placeholder="chatgo_webhook_123"></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> Tambah</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">Meta Accounts</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Nama</th><th>Phone ID</th><th>Action</th></tr></thead><tbody>
        {{range .MetaAccounts}}<tr><td>{{.ID}}</td><td>{{.Name}} <span class="badge badge-soft-primary" style="font-size:9px">Meta</span></td><td>{{.PhoneNumberID}}</td><td><form method="post" action="/admin/meta/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Delete</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
  <div class="card mt-3"><div class="card-header"><h4 class="card-header-title">Webhook URL</h4></div>
  <div class="card-body">
    <p class="small text-muted">Copy URL ini ke Facebook Developer Console &gt; WhatsApp &gt; Configuration &gt; Webhook:</p>
    <code id="webhookUrl" style="word-break:break-all">{{.AppURL}}/webhook/meta</code>
    <p class="small text-muted mt-2">Verify Token: sesuai yang diisi di form atas.</p>
</div></div>
<script>
(function(){
var navs=document.querySelectorAll('.docs-nav a');
var sections=document.querySelectorAll('.docs-section');
function onScroll(){
var scroll=window.scrollY+90;
sections.forEach(function(sec, i){
var top=sec.offsetTop;
var h=sec.offsetHeight;
navs.forEach(function(a){a.classList.remove('active')});
if(scroll >= top && scroll < top+h && navs[i]) navs[i].classList.add('active');
});
}
window.addEventListener('scroll',onScroll);
navs.forEach(function(a){a.addEventListener('click',function(e){e.preventDefault();var id=this.getAttribute('href').slice(1);var el=document.getElementById(id);if(el)el.scrollIntoView({behavior:'smooth',block:'start'})})});
onScroll();
})();
</script>
{{end}}

{{if eq .Page "admin_metatemplates"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">Tambah Template</h4></div><div class="card-body">
      <form method="post" action="/admin/metatemplates/add">
        <div class="form-group"><label>Nama Template</label><input name="name" class="form-control" placeholder="hello_world" required></div>
        <div class="form-group"><label>Language</label><select name="language" class="form-control"><option value="id">Indonesia</option><option value="en">English</option><option value="en_US">English (US)</option></select></div>
        <div class="form-group"><label>Category</label><select name="category" class="form-control"><option value="marketing">Marketing</option><option value="utility">Utility</option><option value="authentication">Authentication</option></select></div>
        <div class="form-group"><label>Components (JSON)</label><textarea name="components" class="form-control" rows="4" placeholder='[{"type":"body","text":"Halo {{1}}, pesanan {{2}} sudah diproses"}]'></textarea></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> Tambah</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">Templates</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Nama</th><th>Lang</th><th>Category</th><th>Action</th></tr></thead><tbody>
        {{range .MetaTemplates}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Language}}</td><td>{{.Category}}</td><td><form method="post" action="/admin/metatemplates/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">Delete</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "inbox"}}
<style>
.inbox-conv{cursor:pointer;transition:background .15s;border-left:3px solid transparent}
.inbox-conv:hover{background:#f8f9fc;border-left-color:#2c7be5}
.inbox-conv.unread{font-weight:600;background:#eef2ff}
.inbox-conv .last-msg{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:260px}
.inbox-search{position:relative}
.inbox-search i{position:absolute;left:12px;top:50%;transform:translateY(-50%);color:#95aac9}
.inbox-search input{padding-left:36px}
.avatar{width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:700;font-size:16px;flex-shrink:0}
.avatar.group{background:linear-gradient(135deg,#10B981,#059669);font-size:14px}
.avatar.person{background:linear-gradient(135deg,#4F46E5,#6366F1)}
.main-tabs{display:flex;border-bottom:2px solid #e0e0e0;margin-bottom:0}
.main-tabs .tab-item{padding:10px 24px;font-size:14px;font-weight:600;cursor:pointer;color:#6e788c;border:none;border-bottom:2px solid transparent;margin-bottom:-2px;transition:all .15s;background:none}
.main-tabs .tab-item:hover{color:#152e4d}
.main-tabs .tab-item.active{color:#2c7be5;border-bottom-color:#2c7be5}
.sub-tabs{display:flex;gap:2px;padding:10px 16px}
.sub-tabs .sub-btn{padding:4px 14px;border:1px solid #e0e0e0;background:#fff;border-radius:6px;font-size:12px;font-weight:500;cursor:pointer;transition:all .15s}
.sub-tabs .sub-btn:hover{background:#f5f5f5}
.sub-tabs .sub-btn.active{background:#2c7be5;color:#fff;border-color:#2c7be5}
.tab-panel{display:none}
.tab-panel.active{display:block}
.wa-status-card{border:1px solid #e0e0e0;border-radius:10px;padding:16px 20px;margin:12px 16px}
.wa-status-card .acc-row{display:flex;align-items:center;justify-content:space-between;padding:10px 0;border-bottom:1px solid #f0f0f0}
.wa-status-card .acc-row:last-child{border-bottom:0}
.acc-status{display:flex;align-items:center;gap:8px;font-size:13px;font-weight:500}
.acc-status .dot{width:9px;height:9px;border-radius:50%}
.acc-status .dot.green{background:#00d97e}
.acc-status .dot.yellow{background:#f6c343}
.acc-status .dot.red{background:#e63757}
.acc-phone{font-family:monospace;font-size:13px;color:#6e788c}
</style>

<div class="card">
<div class="card-header pb-0">
<div class="d-flex justify-content-between align-items-center">
<h4 class="card-header-title mb-0">Live Chat{{if gt .UnreadCount 0}} <span class="badge badge-danger">{{.UnreadCount}} baru</span>{{end}}</h4>
<div class="inbox-search" style="width:220px"><i class="la la-search"></i><input type="text" id="inboxSearch" class="form-control form-control-sm" placeholder="Cari..."></div>
</div>
<div class="main-tabs mt-2">
<button class="tab-item active" onclick={{js "var p=document.getElementById('chat-panel'),s=document.getElementById('status-panel');p.style.display='block';s.style.display='none';var btns=this.parentElement.querySelectorAll('button');for(var i=0;i<btns.length;i++)btns[i].classList.remove('active');this.classList.add('active');return false"}}>Chat</button>
<button class="tab-item" onclick={{js "var p=document.getElementById('chat-panel'),s=document.getElementById('status-panel');p.style.display='none';s.style.display='block';var btns=this.parentElement.querySelectorAll('button');for(var i=0;i<btns.length;i++)btns[i].classList.remove('active');this.classList.add('active');return false"}}>Status</button>
</div>
</div>

<div class="tab-panel active" id="chat-panel">
<div class="sub-tabs">
<button class="sub-btn active" onclick={{js "var d=document.querySelectorAll('#inboxList .inbox-conv');for(var i=0;i<d.length;i++)d[i].style.display='';var s=this.parentElement.querySelectorAll('button');for(var i=0;i<s.length;i++)s[i].classList.remove('active');this.classList.add('active')"}}>Semua</button>
<button class="sub-btn" onclick={{js "var d=document.querySelectorAll('#inboxList .inbox-conv');for(var i=0;i<d.length;i++){var g=d[i].getAttribute('data-group');d[i].style.display=g==='private'?'':'none'};var s=this.parentElement.querySelectorAll('button');for(var i=0;i<s.length;i++)s[i].classList.remove('active');this.classList.add('active')"}}>Private</button>
<button class="sub-btn" onclick={{js "var d=document.querySelectorAll('#inboxList .inbox-conv');for(var i=0;i<d.length;i++){var g=d[i].getAttribute('data-group');d[i].style.display=g==='group'?'':'none'};var s=this.parentElement.querySelectorAll('button');for(var i=0;i<s.length;i++)s[i].classList.remove('active');this.classList.add('active')"}}>Group</button>
<button class="sub-btn" onclick="window.location='/inbox?filter=unread'">Unread</button>
</div>
<div class="list-group list-group-flush" id="inboxList">
{{range .InboxConversations}}
<a href="/inbox/chat?phone={{.Phone}}" class="list-group-item list-group-item-action inbox-conv{{if gt .Unread 0}} unread{{end}}" data-group="{{if .IsGroup}}group{{else}}private{{end}}">
<div class="d-flex align-items-center gap-3">
<div class="avatar {{if .IsGroup}}group{{else}}person{{end}}">{{if .IsGroup}}G{{else}}{{slice .Name 0 1}}{{if not .Name}}+{{end}}{{end}}</div>
<div class="flex-grow-1 min-w-0">
<div class="d-flex justify-content-between align-items-start">
<div class="d-flex align-items-center gap-2">
<strong>{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}</strong>
{{if .IsGroup}}<span class="badge badge-soft-success" style="font-size:10px">Group</span>{{end}}
{{if eq .Channel "meta"}}<span class="badge badge-soft-primary" style="font-size:10px">Meta</span>{{end}}
{{if gt .Unread 0}}<span class="badge badge-pill badge-danger" style="font-size:10px">{{.Unread}}</span>{{end}}
</div>
<small class="text-muted">{{.LastTime}}</small>
</div>
<div class="last-msg text-muted small mt-1">{{.LastMsg}}</div>
</div>
</div>
</a>
{{else}}
<div class="list-group-item text-center text-muted py-4">Belum ada percakapan</div>
{{end}}
</div>
{{if .InboxPages}}<div class="card-footer"><div class="pagination">{{range .InboxPages}}{{if eq . $.PageNum}}<span class="active">{{.}}</span>{{else}}<a href="?page={{.}}">{{.}}</a>{{end}}{{end}}</div></div>{{end}}
</div>

<div class="tab-panel" id="status-panel">
<div class="p-3">
<h6 class="mb-3">Status</h6>
<div class="row g-3">
{{range .Statuses}}
<div class="col-6 col-md-4 col-lg-3">
<div class="card border" style="border-radius:12px;overflow:hidden">
<div class="d-flex align-items-center gap-2 p-3">
<div class="avatar person" style="width:40px;height:40px;font-size:14px">{{slice .Name 0 1}}{{if not .Name}}+{{end}}</div>
<div class="min-w-0">
<div class="fw-bold small">{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}</div>
<div class="text-muted" style="font-size:11px">{{.Created}}</div>
</div>
</div>
{{if .MediaURL}}<div style="height:120px;background:#f0f0f0;display:flex;align-items:center;justify-content:center;color:#aaa;font-size:12px">Media</div>{{end}}
{{if .Caption}}<div class="p-2 small">{{.Caption}}</div>{{end}}
</div>
</div>
{{else}}
<div class="col-12 text-center text-muted py-4">Belum ada status. Status muncul saat kontak posting story.</div>
{{end}}
</div>
</div>
</div>
</div>

<script>
var srch=document.getElementById('inboxSearch');if(srch)srch.addEventListener('input',function(){
var q=this.value.toLowerCase();
document.querySelectorAll('#inboxList .inbox-conv').forEach(function(el){
if(el.style.display==='none')return;
el.style.display=el.textContent.toLowerCase().includes(q)?'':'none';
});
});
setInterval(function(){
fetch('/inbox/unread-count').then(function(r){return r.json()}).then(function(d){
var b=document.querySelector('.inbox-badge');
if(d.unread>0){if(b)b.textContent=d.unread;else{var n=document.querySelector('.nav-link[href="/inbox"]');if(n)n.innerHTML+=' <span class="badge badge-pill badge-danger ml-1 inbox-badge">'+d.unread+'</span>'}}
else{if(b)b.remove()}
});
},5000);
</script>
{{end}}
{{if eq .Page "inbox_chat"}}
<style>
.chat-area{height:55vh;overflow-y:auto;padding:16px;background:#efeae2;background-image:url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><rect width="200" height="200" fill="%23efeae2"/><circle cx="100" cy="100" r="60" fill="%23e5ddd5" opacity="0.5"/></svg>')}
.chat-bubble{max-width:75%;padding:8px 12px;border-radius:8px;word-wrap:break-word;box-shadow:0 1px 1px rgba(0,0,0,.08);font-size:14px;line-height:1.4;position:relative}
.chat-bubble.received{background:#fff;align-self:flex-start;border-top-left-radius:0}
.chat-bubble.sent{background:#d9fdd3;align-self:flex-end;border-top-right-radius:0}
.chat-meta{font-size:10.5px;color:#667781;margin-bottom:1px}
.chat-sender{font-size:12px;font-weight:600;color:#10B981;margin-bottom:1px}
.chat-time{font-size:10px;color:#99aab5;float:right;margin-left:8px;margin-top:2px}
.chat-input-group{position:relative}
.chat-input-group textarea{resize:none;border-radius:8px;padding:10px 48px 10px 14px;min-height:44px;max-height:120px;border:1px solid #e0e0e0;font-size:14px}
.chat-input-group textarea:focus{outline:none;border-color:#10B981;box-shadow:0 0 0 3px rgba(16,185,129,.1)}
.chat-input-group button{position:absolute;right:6px;bottom:6px;border-radius:50%;width:36px;height:36px;padding:0;display:flex;align-items:center;justify-content:center;background:#10B981;border:none;color:#fff;cursor:pointer}
.chat-input-group button:hover{background:#059669}
.chat-input-group button:disabled{background:#ccc;cursor:default}
</style>
<div class="card border-0 shadow-sm">
<div class="card-header d-flex justify-content-between align-items-center bg-white border-bottom" style="padding:10px 16px">
<h6 class="mb-0 d-flex align-items-center gap-2">
<a href="/inbox" class="text-decoration-none text-muted">&larr;</a>
<div class="avatar {{if .IsGroup}}group{{else}}person{{end}}" style="width:36px;height:36px;font-size:13px">{{if .ChatName}}{{slice .ChatName 0 1}}{{else}}+{{end}}</div>
<div><strong>{{if .ChatName}}{{.ChatName}}{{else}}+{{.Phone}}{{end}}</strong>{{if .IsGroup}}<small class="text-success ms-1">Group</small>{{end}}{{if .Channel}}<small class="badge badge-soft-primary ms-1">{{.Channel}}</small>{{end}}</div>
</h6>
<div class="d-flex gap-2 align-items-center">
<select id="chatChannel" class="form-select form-select-sm" style="width:auto;font-size:12px" onchange="onChannelChange()">
<option value="whatsmeow">WA</option>
{{if .MetaAccounts}}<option value="meta">Meta</option>{{end}}
</select>
<select id="chatAccountPhone" class="form-select form-select-sm" style="width:auto;display:inline"> {{range .ConnectedAccounts}}{{if eq .Status "connected"}}<option value="+{{.Phone}}">+{{.Phone}}</option>{{end}}{{end}}</select>
<select id="chatMetaAccount" class="form-select form-select-sm" style="width:auto;display:none"> {{range .MetaAccounts}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select>
{{if .Users}}<form method="post" action="/inbox/assign" style="display:inline" class="me-1"><input type="hidden" name="phone" value="{{.Phone}}"><select name="agent_id" class="form-select form-select-sm" style="width:auto;font-size:11px" onchange="this.form.submit()"><option value="0">Unassigned</option>{{range .Users}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></form><form method="post" action="/inbox/close" style="display:inline"><input type="hidden" name="phone" value="{{.Phone}}"><button class="btn btn-sm btn-outline-danger" style="font-size:11px;padding:2px 8px">Close</button></form>{{end}}</div>
</div>
<div class="card-body p-0">
<div class="chat-area" id="chatMessages">
{{range .ChatMessages}}
<div class="d-flex w-100 mb-1" style="{{if eq .Type "sent"}}justify-content:flex-end{{else}}justify-content:flex-start{{end}}">
<div class="chat-bubble {{.Type}}">
{{if and (eq .Type "received") .SenderName}}<div class="chat-sender">{{.SenderName}}</div>{{end}}
<div>{{.Message}}<span class="chat-time">{{.Created}}</span></div>
</div>
</div>
{{else}}
<div class="text-center text-muted py-4">Belum ada pesan. Kirim pesan pertama!</div>
{{end}}
</div>
</div>
<div class="card-footer bg-white border-top" style="padding:8px 16px">
{{if .Notes}}<div class="mb-2" style="max-height:100px;overflow-y:auto">{{range .Notes}}<div class="small text-muted mb-1"><i class="la la-sticky-note me-1"></i> {{.Note}} <span class="text-muted" style="font-size:10px">{{.Created}}</span></div>{{end}}</div>{{end}}
<div class="d-flex gap-1 mb-1">
  <form method="post" action="/inbox/note" class="d-flex gap-1 flex-grow-1"><input type="hidden" name="phone" value="{{.Phone}}"><input name="note" class="form-control form-control-sm" placeholder="Tambah catatan internal..." style="font-size:12px"><button class="btn btn-sm btn-outline-secondary" style="font-size:11px"><i class="la la-sticky-note"></i></button></form>
</div>
<form id="chatForm" onsubmit="return sendChat(event)">
<div class="chat-input-group">
<textarea id="chatInput" name="message" class="form-control" placeholder="Ketik pesan..." rows="1" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();sendChat(event)}"></textarea>
<input type="hidden" name="phone" value="{{.Phone}}">
<button type="submit" id="sendBtn"><i class="la la-paper-plane"></i></button>
</div>
</form>
        <div class="d-flex gap-1 mt-2 flex-wrap">
          {{range .Templates}}<button class="btn btn-sm btn-outline-secondary template-btn" data-content="{{.Content}}" title="{{.Name}}" style="font-size:11px;padding:2px 8px">{{.Name}}</button>{{end}}
          {{range .Canned}}<button class="btn btn-sm btn-outline-info canned-btn" data-content="{{.Message}}" title="{{.Name}}" style="font-size:11px;padding:2px 8px">{{if .Shortcut}}/{{.Shortcut}}{{else}}{{.Name}}{{end}}</button>{{end}}
        </div>
</div>
</div>
<script>
var chatPhone="{{.Phone}}";
var chatBox=document.getElementById('chatMessages');
function scrollToBottom(){chatBox.scrollTop=chatBox.scrollHeight}
function onChannelChange(){
var ch=document.getElementById('chatChannel').value;
document.getElementById('chatAccountPhone').style.display=ch==='whatsmeow'?'':'none';
document.getElementById('chatMetaAccount').style.display=ch==='meta'?'':'none';
}

function sendChat(e){
e.preventDefault();
var msg=document.getElementById('chatInput').value.trim();
if(!msg)return false;
var ch=document.getElementById('chatChannel').value;
var f=new FormData();
f.append('phone',chatPhone);
f.append('message',msg);
var url='/inbox/send';
if(ch==='meta'){
var mid=document.getElementById('chatMetaAccount').value;
f.append('account_id',mid);
url='/inbox/send-meta';
}else{
var acp=document.getElementById('chatAccountPhone');
if(acp)f.append('account_phone',acp.value);
}
document.getElementById('chatInput').value='';
document.getElementById('sendBtn').disabled=true;
fetch(url,{method:'POST',body:f}).then(function(r){return r.json()}).then(function(d){
document.getElementById('sendBtn').disabled=false;
if(d.ok)loadMessages();
});
return false;
}

function loadMessages(){
fetch('/inbox/messages?phone='+encodeURIComponent(chatPhone)).then(function(r){return r.json()}).then(function(msgs){
if(!msgs||!msgs.length){chatBox.innerHTML='<div class="text-center text-muted py-4">Belum ada pesan</div>';return}
var html='';
for(var i=0;i<msgs.length;i++){
var m=msgs[i];
var side=m.type==='sent'?'flex-end':'flex-start';
html+='<div class="d-flex w-100 mb-1" style="justify-content:'+side+'"><div class="chat-bubble '+m.type+'">';
if(m.type==='received'&&m.sender_name)html+='<div class="chat-sender">'+m.sender_name+'</div>';
html+='<div>'+m.message+'<span class="chat-time">'+m.created+'</span></div></div></div>';
}
chatBox.innerHTML=html;
scrollToBottom();
});
}

scrollToBottom();

var evtSource=new EventSource('/inbox/events');
evtSource.onmessage=function(e){
var d=JSON.parse(e.data);
if(d.phone===chatPhone)loadMessages();
fetch('/inbox/unread-count').then(function(r){return r.json()}).then(function(d){
var b=document.querySelector('.inbox-badge');
if(d.unread>0){if(b)b.textContent=d.unread;else{var n=document.querySelector('.nav-link[href="/inbox"]');if(n)n.innerHTML+=' <span class="badge badge-pill badge-danger ml-1 inbox-badge">'+d.unread+'</span>'}}
else{if(b)b.remove()}
});
};

document.querySelectorAll('.template-btn').forEach(function(btn){
btn.addEventListener('click',function(){
document.getElementById('chatInput').value=this.dataset.content;
document.getElementById('chatInput').focus();
});
});
document.querySelectorAll('.canned-btn').forEach(function(btn){
btn.addEventListener('click',function(){
document.getElementById('chatInput').value=this.dataset.content;
document.getElementById('chatInput').focus();
});
});
</script>
{{end}}
{{if eq .Page "docs"}}
<style>
html{scroll-behavior:smooth}
.docs-nav{position:sticky;top:80px}
.docs-nav a{display:block;padding:4px 12px;color:#555;text-decoration:none;font-size:13px;border-left:2px solid transparent;transition:.15s}
.docs-nav a:hover,.docs-nav a.active{border-left-color:#4F46E5;color:#4F46E5;background:rgba(79,70,229,.05)}
.docs-section{margin-bottom:32px}
.docs-section h3{font-weight:700;margin-bottom:4px;padding-bottom:8px;border-bottom:2px solid #eee}
.badge-soft-info{background:rgba(44,123,229,.1);color:#2c7be5}
</style>
<div class="row">
<div class="col-12 col-lg-3">
  <div class="docs-nav">
    <a href="#start">Quick Start</a>
    <a href="#wa">WhatsApp Setup</a>
    <a href="#broadcast">Broadcast & Campaign</a>
    <a href="#drip">Drip Campaign</a>
    <a href="#store">Store & Products</a>
    <a href="#payment">Payment Gateway</a>
    <a href="#inbox">Live Chat Inbox</a>
    <a href="#ai">AI Auto Reply</a>
    <a href="#team">Team & Support</a>
    <a href="#contacts">Contacts & CRM</a>
    <a href="#safety">Safety & Blacklist</a>
    <a href="#analytics">Analytics</a>
    <a href="#tools">Tools & System</a>
    <a href="#widget">Web Widget</a>
    <a href="#api">API Reference</a>
  </div>
</div>
<div class="col-12 col-lg-9">

<div class="docs-section" id="start"><h3>Quick Start</h3>
<p class="text-muted">Langkah pertama menggunakan {{.AppName}}</p>
<p><strong>Demo Login:</strong> <code>{{.AppEmail}}</code> / <code>password</code></p>
<p class="small text-muted mb-3"><span class="badge badge-soft-success" style="font-size:10px">WA Web</span> WhatsApp Web (whatsmeow) &nbsp; <span class="badge badge-soft-primary" style="font-size:10px">Meta</span> WhatsApp Cloud API &nbsp; <span class="badge badge-soft-info" style="font-size:10px">Both</span> Tersedia di kedua channel</p>
<div class="row g-2">
<div class="col-12 col-md-6"><div class="card"><div class="card-body"><strong>1. Hubungkan WA</strong><p class="small text-muted mb-0">Buka Akun & QR → Scan QR dengan WhatsApp → Connected</p></div></div></div>
<div class="col-12 col-md-6"><div class="card"><div class="card-body"><strong>2. Tambah Kontak</strong><p class="small text-muted mb-0">Import CSV atau tambah manual di Contacts</p></div></div></div>
<div class="col-12 col-md-6"><div class="card"><div class="card-body"><strong>3. Buat Broadcast</strong><p class="small text-muted mb-0">Pilih grup kontak → tulis pesan → kirim massal</p></div></div></div>
<div class="col-12 col-md-6"><div class="card"><div class="card-body"><strong>4. Setup AI</strong><p class="small text-muted mb-0">Tambah AI Key → buat Auto Reply rule → AI balas otomatis</p></div></div></div>
</div></div>

<div class="docs-section" id="wa"><h3>WhatsApp Setup</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>WA Web Connect</td><td>/wa</td><td>Klik Tambah Akun → Scan QR <span class="badge badge-soft-success" style="font-size:9px">WA Web</span></td></tr>
<tr><td>Multi-Akun WA</td><td>/wa</td><td>Tambah beberapa nomor <span class="badge badge-soft-success" style="font-size:9px">WA Web</span></td></tr>
<tr><td>Meta Cloud API</td><td>/admin/meta</td><td>Phone ID + Access Token <span class="badge badge-soft-primary" style="font-size:9px">Meta</span></td></tr>
<tr><td>Meta Templates</td><td>/admin/metatemplates</td><td>Template WA Business <span class="badge badge-soft-primary" style="font-size:9px">Meta</span></td></tr>
<tr><td>Kirim Pesan</td><td>/send</td><td>Kirim pesan ke satu nomor <span class="badge badge-soft-success" style="font-size:9px">WA Web</span></td></tr>
</tbody></table></div>

<div class="docs-section" id="broadcast"><h3>Broadcast & Campaign</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Keterangan</th></tr></thead>
<tbody>
<tr><td>Broadcast Massal</td><td>/broadcast</td><td>Pilih grup + nomor + tag filter <span class="badge badge-soft-info" style="font-size:9px">Both</span></td></tr>
<tr><td>Media Broadcast</td><td>/broadcast</td><td>Upload gambar/video/dokumen <span class="badge badge-soft-info" style="font-size:9px">Both</span></td></tr>
<tr><td>Round Robin / Random</td><td>/broadcast</td><td>Radio select di form broadcast</td></tr>
<tr><td>Pause / Resume</td><td>/broadcast</td><td>Tombol ⏸/▶ di campaign list</td></tr>
<tr><td>Retry Campaign</td><td>/broadcast</td><td>Tombol ↩ untuk jalankan ulang</td></tr>
<tr><td>Auto-Resume</td><td>Background</td><td>Restart server → lanjut dari nomor terakhir</td></tr>
<tr><td>Link Tracking</td><td>/broadcast</td><td>URL otomatis di-track via /track/:token</td></tr>
<tr><td>A/B Testing</td><td>/ab-tests</td><td>Buat varian A/B, lihat hasil</td></tr>
<tr><td>Recurring Campaign</td><td>/recurring</td><td>Jadwal broadcast otomatis daily/weekly</td></tr>
<tr><td>Scheduled Message</td><td>/scheduled</td><td>Jadwalkan pesan ke nomor spesifik</td></tr>
<tr><td>Anti-Ban Rate Limiter</td><td>/settings</td><td>Max/day + random delay interval</td></tr>
</tbody></table></div>

<div class="docs-section" id="drip"><h3>Drip Campaign</h3>
<p class="text-muted">Multi-step follow-up otomatis saat customer kirim pesan.</p>
<table class="table table-sm">
<tr><td><strong>Setup</strong></td><td>/drips → buat drip → tambah steps (delay + message)</td></tr>
<tr><td><strong>Auto-Enroll</strong></td><td>Setiap pesan masuk → otomatis masuk semua drip aktif</td></tr>
<tr><td><strong>STOP</strong></td><td>Customer reply "STOP" → unenroll dari semua drip</td></tr>
<tr><td><strong>Pause</strong></td><td>Toggle Active/Inactive di halaman /drips</td></tr>
</table></div>

<div class="docs-section" id="store"><h3>WA Store Bot</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>Product Catalog</td><td>/store</td><td>Tambah produk: nama, harga, gambar, kategori, stok</td></tr>
<tr><td>Kategori</td><td>/store</td><td>Tambah/hapus kategori produk</td></tr>
<tr><td>Order Management</td><td>/store/orders</td><td>Lihat & update status order (new→paid→shipped)</td></tr>
<tr><td>WA Menu Bot</td><td>Auto</td><td>Chat "menu" → daftar kategori → pilih produk → order</td></tr>
<tr><td>WA Order Notif</td><td>Auto</td><td>Status order berubah → WA otomatis ke customer</td></tr>
</tbody></table></div>

<div class="docs-section" id="payment"><h3>Payment Gateway</h3>
<table class="table table-sm"><thead><tr><th>Gateway</th><th>Region</th><th>Setup</th></tr></thead>
<tbody>
<tr><td>Midtrans</td><td>Indonesia</td><td>/admin/gateways-pay → pilih midtrans → input Server Key</td></tr>
<tr><td>Xendit</td><td>Indonesia</td><td>/admin/gateways-pay → pilih xendit → input API Key</td></tr>
<tr><td>PayPal</td><td>International</td><td>/admin/gateways-pay → pilih paypal → input Access Token</td></tr>
<tr><td>Stripe</td><td>International</td><td>/admin/gateways-pay → pilih stripe → input Secret Key</td></tr>
</tbody></table>
<p class="text-muted"><strong>Flow:</strong> User pilih paket di /subscribe → checkout via gateway → callback → subscription auto-active</p></div>

<div class="docs-section" id="inbox"><h3>Live Chat Inbox</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>Real-time Inbox</td><td>/inbox — SSE auto-refresh, filter Private/Group/Unread</td></tr>
<tr><td>Send/Receive</td><td>Klik conversation → kirim pesan langsung</td></tr>
<tr><td>Media Preview</td><td>Gambar/video tampil inline</td></tr>
<tr><td>Status Tab</td><td>Lihat story/status dari kontak</td></tr>
</tbody></table></div>

<div class="docs-section" id="ai"><h3>AI Auto Reply</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>AI Keys</td><td>/ai/keys → tambah key (OpenAI/Gemini/Claude/DeepSeek)</td></tr>
<tr><td>Auto Reply Rules</td><td>/autoreply → match type (contains/exact/starts/AI)</td></tr>
<tr><td>Knowledge Base</td><td>/knowledge → tambah Q&A, import CSV, upload PDF/URL</td></tr>
<tr><td>AI Global Mode</td><td>/settings → centang "AI untuk Semua Pesan"</td></tr>
<tr><td>AI Store Agent (RAG)</td><td>AI auto-inject product catalog + customer profile ke context</td></tr>
<tr><td>Human Handoff</td><td>/settings → keyword trigger untuk stop AI → kirim kontak admin</td></tr>
<tr><td>Memory Window</td><td>/settings → berapa pesan terakhir dikirim ke AI</td></tr>
<tr><td>Working Hours</td><td>/settings → jam kerja + pesan luar jam + off days</td></tr>
</tbody></table></div>

<div class="docs-section" id="team"><h3>Team & Support</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>Agent Assignment</td><td>Inbox Chat</td><td>Dropdown assign agent, auto round-robin</td></tr>
<tr><td>Departments</td><td>/depts</td><td>Buat Sales/Support/Billing + assign agents</td></tr>
<tr><td>Chat Transfer</td><td>Inbox Chat</td><td>Dropdown transfer ke agent lain</td></tr>
<tr><td>Conversation Notes</td><td>Inbox Chat</td><td>Catat note internal, tampil di bawah chat</td></tr>
<tr><td>Close Conversation</td><td>Inbox Chat</td><td>Tombol Close → kirim CSAT survey</td></tr>
<tr><td>Canned Responses</td><td>/canned</td><td>Shortcut balasan cepat, muncul di inbox</td></tr>
<tr><td>Agent Signature</td><td>/settings</td><td>Auto-append di akhir pesan agent</td></tr>
<tr><td>Chat Labels</td><td>Inbox</td><td>Tag conversation (urgent/follow-up/resolved)</td></tr>
<tr><td>Inbox Macros</td><td>/macros</td><td>One-click multi-action: assign+tag+reply+close</td></tr>
<tr><td>Auto-Close Idle</td><td>/settings</td><td>Tutup otomatis setelah X jam + follow-up</td></tr>
<tr><td>CSAT Survey</td><td>Auto</td><td>Auto-kirim rating request saat close</td></tr>
<tr><td>VIP Priority</td><td>Inbox</td><td>Set priority=1, antrian lebih cepat</td></tr>
</tbody></table></div>

<div class="docs-section" id="contacts"><h3>Contacts & CRM</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>Contact List</td><td>/contacts</td><td>Tambah/edit/hapus kontak manual</td></tr>
<tr><td>CSV Import</td><td>/contacts</td><td>Upload CSV (name,phone,groups)</td></tr>
<tr><td>CSV Export</td><td>/contacts</td><td>Download semua kontak ke CSV</td></tr>
<tr><td>Bulk Delete</td><td>/contacts</td><td>Centang → Delete Selected</td></tr>
<tr><td>Contact Groups</td><td>/contacts/groups</td><td>Atur grup + language per grup</td></tr>
<tr><td>Contact Tags</td><td>/tags</td><td>Tag custom (VIP, Leads) + warna</td></tr>
<tr><td>Multi-Lang Groups</td><td>/contacts/groups</td><td>Set bahasa ID/EN per grup</td></tr>
<tr><td>Unsubscribe</td><td>/contacts/unsub</td><td>Kelola nomor yang opt-out</td></tr>
<tr><td>Customer Directory</td><td>/customers</td><td>Search & chat langsung</td></tr>
<tr><td>Contact Merge</td><td>/merge</td><td>Deteksi & gabung duplikat</td></tr>
<tr><td>Number Validator</td><td>/broadcast</td><td>Tombol Validate cek format+blacklist</td></tr>
</tbody></table></div>

<div class="docs-section" id="safety"><h3>Safety & Blacklist</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Cara</th></tr></thead>
<tbody>
<tr><td>Smart Blacklist</td><td>/blacklist</td><td>Auto-detect spam → auto-block. Manual add/remove</td></tr>
<tr><td>Spam Detection</td><td>Auto</td><td>8+ identik dalam 10 menit → auto-blacklist</td></tr>
<tr><td>Rate Limiter</td><td>/settings</td><td>Max/day + random delay anti-ban</td></tr>
</tbody></table></div>

<div class="docs-section" id="analytics"><h3>Analytics & Reports</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Keterangan</th></tr></thead>
<tbody>
<tr><td>Dashboard Chart</td><td>/home</td><td>Grafik 7 hari sent vs received</td></tr>
<tr><td>Agent Performance</td><td>/analytics</td><td>Chats, replies, avg response time</td></tr>
<tr><td>CSAT Score</td><td>/analytics</td><td>Average rating dari survey</td></tr>
<tr><td>Link Tracker</td><td>/tracker</td><td>Lihat semua link + klik status</td></tr>
<tr><td>Campaign Calendar</td><td>/calendar</td><td>Visual timeline campaign + recurring</td></tr>
<tr><td>Audit Log</td><td>/audit</td><td>Trace siapa ngapain kapan</td></tr>
<tr><td>Payment Reminders</td><td>/reminders</td><td>Auto-WA tagihan + due tracking</td></tr>
</tbody></table></div>

<div class="docs-section" id="tools"><h3>Tools & System</h3>
<table class="table table-sm"><thead><tr><th>Fitur</th><th>Lokasi</th><th>Keterangan</th></tr></thead>
<tbody>
<tr><td>Templates</td><td>/templates</td><td>Template pesan reusable</td></tr>
<tr><td>API Keys</td><td>/apikeys</td><td>Generate key untuk API access</td></tr>
<tr><td>Webhooks</td><td>/webhooks</td><td>Kirim event ke URL eksternal</td></tr>
<tr><td>Language</td><td>Navbar</td><td>Switch ID ↔ EN, semua konten bilingual</td></tr>
<tr><td>Whitelabel</td><td>/settings</td><td>Ganti logo, nama, email</td></tr>
<tr><td>Database Backup</td><td>/backup</td><td>One-click backup DB</td></tr>
<tr><td>File Manager</td><td>/uploads</td><td>Browse uploaded media</td></tr>
<tr><td>Auto-Translate</td><td>POST /translate</td><td>AI translate text ke bahasa target</td></tr>
<tr><td>Email→WA Gateway</td><td>POST /email-webhook</td><td>Forward email ke WA inbox</td></tr>
</tbody></table></div>

<div class="docs-section" id="widget"><h3>Web Widget</h3>
<p class="text-muted">Embeddable chat widget untuk website.</p>
<pre class="bg-light p-3 rounded"><code>&lt;script src="{{.AppURL}}/widget.js"&gt;&lt;/script&gt;</code></pre>
<p class="small text-muted">Tambah script ini di HTML website kamu. Tombol chat otomatis muncul di pojok kanan bawah.</p></div>

<div class="docs-section" id="api"><h3>API Reference</h3>
<pre class="bg-light p-3 rounded"><code># Send message
POST /api/send
Header: X-API-Key: &lt;your-api-key&gt;
Body: {"phone":"628123456789","message":"Hello"}

# List contacts
GET /api/contacts
Header: X-API-Key: &lt;your-api-key&gt;

# Send from specific account
POST /api/send
Body: {"phone":"628xx","message":"text","account_phone":"+628xx"}</code></pre></div>

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

