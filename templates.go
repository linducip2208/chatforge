package main

const templates = `
{{define "layout"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
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
  .badge-soft-primary{background:#4F46E5!important;color:#fff!important;font-weight:600}
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
   /* ===== BS5 -> BS4 utility compatibility ===== */
   .me-1{margin-right:.25rem!important}.me-2{margin-right:.5rem!important}.me-3{margin-right:1rem!important}.me-4{margin-right:1.5rem!important}
   .ms-1{margin-left:.25rem!important}.ms-2{margin-left:.5rem!important}.ms-3{margin-left:1rem!important}.ms-4{margin-left:1.5rem!important}
   .ms-auto{margin-left:auto!important}.me-auto{margin-right:auto!important}
   .pe-1{padding-right:.25rem!important}.pe-2{padding-right:.5rem!important}.pe-3{padding-right:1rem!important}
   .ps-1{padding-left:.25rem!important}.ps-2{padding-left:.5rem!important}.ps-3{padding-left:1rem!important}
   .gap-1{gap:.25rem!important}.gap-2{gap:.5rem!important}.gap-3{gap:1rem!important}
   .min-w-0{min-width:0!important}
   .fw-bold{font-weight:700!important}.fw-semibold{font-weight:600!important}
   .text-end{text-align:right!important}.text-start{text-align:left!important}
   .form-select{display:block;width:100%;height:calc(1.5em + .75rem + 2px);padding:.375rem 1.75rem .375rem .75rem;font-size:.9375rem;line-height:1.5;color:#12263f;vertical-align:middle;background:#fff url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 4 5'%3e%3cpath fill='%2395aac9' d='M2 0L0 2h4zm0 5L0 3h4z'/%3e%3c/svg%3e") right .75rem center/8px 10px no-repeat;border:1px solid #d2ddec;border-radius:.375rem;-webkit-appearance:none;appearance:none}
   .form-select:focus{border-color:#2c7be5;outline:0;box-shadow:0 0 0 3px rgba(44,123,229,.12)}
   .form-select-sm{height:calc(1.5em + .5rem + 2px);padding-top:.25rem;padding-bottom:.25rem;padding-left:.5rem;font-size:.8125rem}
   .bg-primary.bg-opacity-10{background-color:rgba(44,123,229,.12)!important}
   .bg-success.bg-opacity-10{background-color:rgba(0,217,126,.12)!important}
   .bg-warning.bg-opacity-10{background-color:rgba(246,195,67,.15)!important}
   .bg-danger.bg-opacity-10{background-color:rgba(230,55,87,.12)!important}
   .bg-info.bg-opacity-10{background-color:rgba(57,160,237,.12)!important}
   .bg-secondary.bg-opacity-10{background-color:rgba(110,120,140,.12)!important}
   /* ===== Responsive base ===== */
   img,canvas,svg,video{max-width:100%;height:auto}
   .table-responsive{-webkit-overflow-scrolling:touch}
   .card-header{flex-wrap:wrap;gap:6px;height:auto;min-height:60px}
   /* Tablet <=1024px */
   @media(max-width:1024px){
     .container-fluid{padding-left:16px;padding-right:16px}
     .header-title{font-size:1.35rem}
     .card-body{padding:1.15rem}
     .table th,.table td{padding:.6rem .75rem}
   }
   /* Mobile <768px — Dashkit collapses sidebar into top hamburger bar */
   @media(max-width:767.98px){
     .main-content{padding-top:4px}
     .header{margin-bottom:14px}
     .header-body{padding-top:16px;padding-bottom:16px}
     .header-title{font-size:1.15rem}
     .header-pretitle{font-size:.64rem}
     .container-fluid{padding-left:12px;padding-right:12px}
     .card{border-radius:10px;margin-bottom:14px}
     .card-header{padding:10px 14px}
     .card-header-title{font-size:.95rem}
     .card-body{padding:1rem}
     h1,.h1{font-size:1.35rem}h2,.h2{font-size:1.2rem}h3,.h3{font-size:1.05rem}
     .table th,.table td{padding:8px 10px;font-size:12.5px}
     .btn{min-height:38px}
     .btn-sm,.btn-group-sm>.btn{min-height:32px}
     .form-control,.form-select,select.custom-select,textarea{font-size:16px!important}
     .form-control-sm{font-size:14px!important}
     .modal-dialog{margin:10px}
     .pagination{gap:3px}
     .pagination a,.pagination span{padding:5px 10px;font-size:12px}
     .nav-tabs .nav-link{padding:8px 10px;font-size:13px}
     .alert{padding:10px 14px;font-size:13.5px}
     #qrimg{width:220px;height:220px}
     .navbar-vertical .nav-link{padding:10px 20px}
     .dropdown-menu{font-size:14px}
   }
   /* Very small <=400px */
   @media(max-width:400px){
     .header-title{font-size:1.02rem}
     .container-fluid{padding-left:10px;padding-right:10px}
     .table th,.table td{font-size:12px;padding:7px 8px}
     #qrimg{width:180px;height:180px}
     .btn{font-size:13px}
   }
   /* Touch targets (WCAG 2.5.5) */
   @media(pointer:coarse){
     .btn{min-height:42px}
     .navbar-vertical .nav-link,.dropdown-item{min-height:42px;display:flex;align-items:center}
     .page-link{min-width:38px;min-height:38px;display:inline-flex;align-items:center;justify-content:center}
   }
   /* Accessibility: reduced motion */
   @media(prefers-reduced-motion:reduce){
     *,*::before,*::after{animation-duration:.01ms!important;animation-iteration-count:1!important;transition-duration:.01ms!important;scroll-behavior:auto!important}
   }
   /* Print */
   @media print{
     #sidebar,#topbar,.pagination,.btn{display:none!important}
     .main-content{margin-left:0!important}
   }
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
          <a href="#" class="dropdown-toggle text-muted" role="button" data-toggle="dropdown" style="text-decoration:none">
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
          <a href="#" class="dropdown-toggle text-muted" role="button" data-toggle="dropdown" style="text-decoration:none">
            <i class="la la-user-circle la-lg"></i>
          </a>
          <div class="dropdown-menu dropdown-menu-end">
            <a class="dropdown-item" href="/settings"><i class="la la-cog me-2"></i> {{T "nav_settings"}}</a>
            <div class="dropdown-divider"></div>
            {{if .IsImpersonating}}
            <a class="dropdown-item text-warning" href="/exit-impersonation"><i class="la la-times-circle me-2"></i> {{T "nav_exit_impersonation"}}</a>
            {{else}}
            <a class="dropdown-item text-danger" href="/logout"><i class="la la-sign-out me-2"></i> {{T "auth_logout"}}</a>
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
  document.getElementById(">{{T "ar_faq_tab"}}Group").style.display=v==="ai"?"block":"none";
  var kw=document.querySelector('input[name="keyword"]');
  if(kw) kw.required=v!=="ai";
}
// tab switcher
document.querySelectorAll('.nav-tabs .nav-link').forEach(function(t){t.addEventListener('click',function(e){e.preventDefault();document.querySelectorAll('.nav-tabs .nav-link').forEach(function(x){x.classList.remove('active')});document.querySelectorAll('.tab-pane').forEach(function(x){x.classList.remove('show','active')});this.classList.add('active');var el=document.querySelector(this.getAttribute('href'));if(el){el.classList.add('show','active')}})});
// truncate messages: show first 20 chars
document.querySelectorAll(".msg-full").forEach(function(el){
  var text=el.textContent.trim();
  if(text.length>20){ el.setAttribute("data-full",text); el.textContent=text.substring(0,20)+"..."; el.style.cursor="pointer"; el.title='{{T "tooltip_click_expand"}}';
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
<style>
.nav-group-header{cursor:pointer;display:flex;align-items:center;justify-content:space-between;padding:6px 12px;margin:2px 8px;border-radius:8px;color:#8895b7;font-size:.75rem;font-weight:700;text-transform:uppercase;letter-spacing:.5px;transition:all .15s;user-select:none}
.nav-group-header:hover{color:#fff;background:rgba(255,255,255,.04)}
.nav-group-header i.la-chevron-down{font-size:11px;transition:transform .2s}
.nav-group-header.collapsed i.la-chevron-down{transform:rotate(-90deg)}
.nav-group-body{overflow:hidden;transition:max-height .3s ease}
.nav-group-body.collapsed{max-height:0!important}
.nav-sub{padding-left:16px}
</style>
<nav class="navbar navbar-vertical fixed-left navbar-expand-md navbar-dark navbar-vibrant" id="sidebar">
  <div class="container-fluid">
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#sidebarCollapse"><span class="navbar-toggler-icon"></span></button>
    <a class="navbar-brand" href="/"><img src="{{.AppLogo}}" class="navbar-brand-img mx-auto" alt="{{.AppName}}" onerror="this.outerHTML='<span style=&quot;color:#fff;font-weight:800;font-size:20px&quot;>{{.AppName}}</span>'"></a>
    <div class="collapse navbar-collapse" id="sidebarCollapse">

<div class="d-md-none px-3 pt-3 pb-2" style="border-bottom:1px solid rgba(255,255,255,.08);margin-bottom:8px">
  <div class="mb-2 text-center">
    {{if eq .Status "connected"}}
      <span class="badge badge-soft-success"><span class="status-dot" style="background:#00d97e"></span> {{T "status_connected"}} +{{.Phone}}</span>
    {{else if eq .Status "qr"}}
      <span class="badge badge-soft-warning"><span class="status-dot" style="background:#f6c343"></span> {{T "status_scanqr"}}</span>
    {{else}}
      <span class="badge badge-soft-danger"><span class="status-dot" style="background:#e63757"></span> {{T "status_disconnected"}}</span>
    {{end}}
  </div>
  <div class="d-flex mb-2" style="gap:8px">
    <a class="btn btn-sm btn-primary flex-fill" href="/wa"><i class="la la-whatsapp me-1"></i> {{T "nav_whatsapp"}}</a>
    <a class="btn btn-sm btn-primary flex-fill" href="/send"><i class="la la-paper-plane me-1"></i> {{T "nav_send"}}</a>
  </div>
  <div class="d-flex align-items-center justify-content-between" style="gap:8px;flex-wrap:wrap">
    <span>{{range .Languages}}<a href="/lang/{{.Code}}" class="me-2" title="{{.Name}}"><span class="flag-icon flag-icon-{{.Flag}} lang-flag"></span></a>{{end}}</span>
    <span>
      <a href="/settings" class="text-muted me-3" style="font-size:13px;text-decoration:none"><i class="la la-cog"></i> {{T "nav_settings"}}</a>
      {{if .IsImpersonating}}
      <a href="/exit-impersonation" class="text-warning" style="font-size:13px;text-decoration:none"><i class="la la-times-circle"></i> {{T "nav_exit_impersonation"}}</a>
      {{else}}
      <a href="/logout" class="text-danger" style="font-size:13px;text-decoration:none"><i class="la la-sign-out"></i> {{T "auth_logout"}}</a>
      {{end}}
    </span>
  </div>
</div>

{{template "sgroup" dict "id" "dashboard" "icon" "la-chart-bar" "label" (T "sg_dashboard") "open" true}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "home"}}active{{end}}" href="/"><i class="la la-chart-bar la-lg"></i> {{T "nav_dashboard"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "inbox" "icon" "la-comments" "label" (T "sg_inbox") "open" true}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "inbox"}}active{{end}}" href="/inbox"><i class="la la-inbox la-lg"></i> {{T "nav_inbox"}}{{if gt .UnreadCount 0}} <span class="badge badge-pill badge-danger ml-1">{{.UnreadCount}}</span>{{end}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "contacts" "icon" "la-address-book" "label" (T "sg_contacts")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "contacts"}}active{{end}}" href="/contacts"><i class="la la-address-book la-lg"></i> {{T "nav_contacts_all"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "groups"}}active{{end}}" href="/contacts/groups"><i class="la la-list la-lg"></i> {{T "nav_contacts_groups"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "tags"}}active{{end}}" href="/tags"><i class="la la-tags la-lg"></i> {{T "nav_tags"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "unsub"}}active{{end}}" href="/contacts/unsub"><i class="la la-unlink la-lg"></i> {{T "nav_contacts_unsub"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "merge"}}active{{end}}" href="/merge"><i class="la la-code-branch la-lg"></i> {{T "nav_merge"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "broadcast" "icon" "la-bullhorn" "label" (T "sg_broadcast")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "broadcast"}}active{{end}}" href="/broadcast"><i class="la la-paper-plane la-lg"></i> {{T "nav_broadcast"}}<span class="badge ms-1" style="background:#4F46E5;color:#fff;font-size:8px;padding:2px 5px;border-radius:3px">META</span></a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "drips"}}active{{end}}" href="/drips"><i class="la la-tint la-lg"></i> {{T "nav_drips"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "recurring"}}active{{end}}" href="/recurring"><i class="la la-redo-alt la-lg"></i> {{T "nav_recurring"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "scheduled"}}active{{end}}" href="/scheduled"><i class="la la-clock la-lg"></i> {{T "nav_scheduled"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "abtests"}}active{{end}}" href="/ab-tests"><i class="la la-balance-scale la-lg"></i> {{T "nav_abtest"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "channels" "icon" "la-plug" "label" (T "sg_channels")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "wa"}}active{{end}}" href="/wa"><i class="la la-whatsapp la-lg"></i> {{T "nav_wa_unofficial"}}</a></li>
  <li class="nav-item" style="padding-left:16px"><a class="nav-link {{if eq .Active "wa"}}active{{end}}" href="/wa"><i class="la la-qrcode la-lg"></i> {{T "nav_account_qr"}}</a></li>
  <li class="nav-item" style="padding-left:16px"><a class="nav-link {{if eq .Active "hosts_whatsapp"}}active{{end}}" href="/hosts/whatsapp"><i class="la la-server la-lg"></i> {{T "nav_devices"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_meta"}}active{{end}}" href="/admin/meta"><i class="la la-cloud la-lg"></i> {{T "nav_meta_api"}}<span class="badge ms-1" style="background:#4F46E5;color:#fff;font-size:8px;padding:2px 5px;border-radius:3px">META</span></a></li>
  <li class="nav-item" style="padding-left:16px"><a class="nav-link {{if eq .Active "admin_metatemplates"}}active{{end}}" href="/admin/metatemplates"><i class="la la-file-alt la-lg"></i> {{T "nav_meta_templates"}}</a></li>
  <li class="nav-item" style="padding-left:16px"><a class="nav-link {{if eq .Active "meta_webhook"}}active{{end}}" href="/meta/webhook"><i class="la la-link la-lg"></i> {{T "nav_meta_webhook"}}</a></li>
  <li class="nav-item" style="padding-left:16px"><a class="nav-link {{if eq .Active "admin_meta"}}active{{end}}" href="/admin/meta"><i class="la la-cog la-lg"></i> {{T "nav_meta_config"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "emailwa"}}active{{end}}" href="/email-wa"><i class="la la-envelope la-lg"></i> {{T "nav_email_wa"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "logs" "icon" "la-envelope" "label" (T "sg_logs")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "send"}}active{{end}}" href="/send"><i class="la la-share la-lg"></i> {{T "nav_send"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "sent"}}active{{end}}" href="/sent"><i class="la la-arrow-up la-lg"></i> {{T "nav_sent"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "received"}}active{{end}}" href="/received"><i class="la la-arrow-down la-lg"></i> {{T "nav_received"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "meta_logs"}}active{{end}}" href="/meta/logs"><i class="la la-clipboard-list la-lg"></i> {{T "nav_meta_logs"}}<span class="badge ms-1" style="background:#4F46E5;color:#fff;font-size:8px;padding:2px 5px;border-radius:3px">META</span></a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "automation" "icon" "la-robot" "label" (T "sg_automation")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "autoreply"}}active{{end}}" href="/autoreply"><i class="la la-reply la-lg"></i> {{T "nav_autoreply"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "flowbuilder"}}active{{end}}" href="/pro/flow-builder"><i class="la la-project-diagram la-lg"></i> Flow Builder</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "canned"}}active{{end}}" href="/canned"><i class="la la-comment-dots la-lg"></i> {{T "nav_canned"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "macros"}}active{{end}}" href="/macros"><i class="la la-bolt la-lg"></i> {{T "nav_macros"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "forms"}}active{{end}}" href="/forms"><i class="la la-wpforms la-lg"></i> {{T "nav_forms"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "reminders"}}active{{end}}" href="/reminders"><i class="la la-bell la-lg"></i> {{T "nav_reminders"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "translate"}}active{{end}}" href="/translate-tool"><i class="la la-language la-lg"></i> {{T "nav_translate"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "widget"}}active{{end}}" href="/widget-info"><i class="la la-code la-lg"></i> {{T "nav_widget"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "commerce" "icon" "la-store" "label" (T "sg_commerce")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "store"}}active{{end}}" href="/store"><i class="la la-store la-lg"></i> {{T "nav_store"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "orders"}}active{{end}}" href="/store/orders"><i class="la la-shopping-bag la-lg"></i> {{T "nav_orders"}}</a></li>
{{template "egroup"}}

{{if .IsAdmin}}
<hr class="navbar-divider my-2">
<div style="color:#6e84a3;font-size:11px;text-align:center;padding:4px;letter-spacing:1px">{{T "nav_admin_divider"}}</div>

{{template "sgroup" dict "id" "admin_overview" "icon" "la-shield-alt" "label" (T "sg_overview")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin"}}active{{end}}" href="/admin"><i class="la la-chart-bar la-lg"></i> {{T "nav_admin_overview"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "admin_business" "icon" "la-building" "label" (T "sg_business")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_users"}}active{{end}}" href="/admin/users"><i class="la la-users la-lg"></i> {{T "adm_users"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_roles"}}active{{end}}" href="/admin/roles"><i class="la la-user-shield la-lg"></i> {{T "adm_roles"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_packages"}}active{{end}}" href="/admin/packages"><i class="la la-box la-lg"></i> {{T "adm_packages"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_subscriptions"}}active{{end}}" href="/admin/subscriptions"><i class="la la-star la-lg"></i> {{T "adm_subscriptions"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "admin_finance" "icon" "la-money-bill" "label" (T "sg_finance")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_vouchers"}}active{{end}}" href="/admin/vouchers"><i class="la la-ticket-alt la-lg"></i> {{T "adm_vouchers"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_transactions"}}active{{end}}" href="/admin/transactions"><i class="la la-receipt la-lg"></i> {{T "adm_transactions"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_paygateways"}}active{{end}}" href="/admin/gateways-pay"><i class="la la-credit-card la-lg"></i> {{T "nav_paygateways"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_transactions_pay"}}active{{end}}" href="/admin/transactions-pay"><i class="la la-file-invoice la-lg"></i> {{T "nav_paylogs"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_payouts"}}active{{end}}" href="/admin/payouts"><i class="la la-hand-holding-usd la-lg"></i> {{T "adm_payouts"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "admin_system" "icon" "la-server" "label" (T "sg_system")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "backup"}}active{{end}}" href="/backup"><i class="la la-database la-lg"></i> {{T "nav_backup"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "audit"}}active{{end}}" href="/audit"><i class="la la-history la-lg"></i> {{T "nav_audit"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "logger"}}active{{end}}" href="/logger"><i class="la la-clipboard-list la-lg"></i> {{T "nav_logger"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "admin_content" "icon" "la-file" "label" (T "sg_content")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_pages"}}active{{end}}" href="/admin/pages"><i class="la la-copy la-lg"></i> {{T "adm_pages"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_marketing"}}active{{end}}" href="/admin/marketing"><i class="la la-bullhorn la-lg"></i> {{T "adm_marketing"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_languages"}}active{{end}}" href="/admin/languages"><i class="la la-language la-lg"></i> {{T "adm_languages"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "admin_infra" "icon" "la-network-wired" "label" (T "sg_infra")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_waservers"}}active{{end}}" href="/admin/waservers"><i class="la la-server la-lg"></i> {{T "adm_waservers"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_gateways"}}active{{end}}" href="/admin/gateways"><i class="la la-code la-lg"></i> {{T "adm_gateways"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_shorteners"}}active{{end}}" href="/admin/shorteners"><i class="la la-link la-lg"></i> {{T "adm_shorteners"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "admin_plugins"}}active{{end}}" href="/admin/plugins"><i class="la la-puzzle-piece la-lg"></i> {{T "adm_plugins"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "reports" "icon" "la-chart-pie" "label" (T "sg_reports")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "analytics"}}active{{end}}" href="/analytics"><i class="la la-chart-pie la-lg"></i> {{T "nav_analytics"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "meta_analytics"}}active{{end}}" href="/meta/analytics"><i class="la la-chart-bar la-lg"></i> {{T "nav_meta_stats"}}<span class="badge ms-1" style="background:#4F46E5;color:#fff;font-size:8px;padding:2px 5px;border-radius:3px">META</span></a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "csat"}}active{{end}}" href="/csat"><i class="la la-star la-lg"></i> {{T "nav_csat"}}</a></li>
{{template "egroup"}}

{{template "sgroup" dict "id" "settings" "icon" "la-cog" "label" (T "sg_settings")}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "settings"}}active{{end}}" href="/settings"><i class="la la-cog la-lg"></i> {{T "nav_general"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "templates"}}active{{end}}" href="/templates"><i class="la la-file-alt la-lg"></i> {{T "nav_templates"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "tracker"}}active{{end}}" href="/tracker"><i class="la la-link la-lg"></i> {{T "nav_links"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "blacklist"}}active{{end}}" href="/blacklist"><i class="la la-ban la-lg"></i> {{T "nav_blacklist"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "uploads"}}active{{end}}" href="/uploads"><i class="la la-folder-open la-lg"></i> {{T "nav_files"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "calendar"}}active{{end}}" href="/calendar"><i class="la la-calendar la-lg"></i> {{T "nav_calendar"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "customers"}}active{{end}}" href="/customers"><i class="la la-users la-lg"></i> {{T "nav_customers"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "apikeys"}}active{{end}}" href="/apikeys"><i class="la la-key la-lg"></i> {{T "nav_apikeys"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "webhooks"}}active{{end}}" href="/webhooks"><i class="la la-code-branch la-lg"></i> {{T "nav_webhooks"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "logger"}}active{{end}}" href="/logger"><i class="la la-clipboard-list la-lg"></i> {{T "nav_logger"}}</a></li>
{{template "egroup"}}
{{end}}

<ul class="navbar-nav mt-2">
  {{if .UserPackage}}
  <li class="nav-item px-4 mb-2">
    <div style="background:rgba(79,70,229,.12);border-radius:8px;padding:10px 12px">
      <div style="font-size:.7rem;color:#8895b7;text-transform:uppercase;letter-spacing:.06em">Paket Aktif</div>
      <div style="font-weight:700;color:#fff;font-size:.85rem">{{.UserPackage}}</div>
      {{if .UserPackageExpire}}<div style="font-size:.7rem;color:#5a6780;margin-top:2px">Sampai {{.UserPackageExpire}}</div>{{end}}
    </div>
  </li>
  {{end}}
  <li class="nav-item"><a class="nav-link {{if eq .Active "upgrade"}}active{{end}}" href="/upgrade"><i class="la la-rocket la-lg"></i> {{T "nav_upgrade"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "subscribe"}}active{{end}}" href="/subscribe"><i class="la la-shopping-cart la-lg"></i> {{T "nav_subscribe"}}</a></li>
  <li class="nav-item"><a class="nav-link {{if eq .Active "docs"}}active{{end}}" href="/docs"><i class="la la-book la-lg"></i> {{T "nav_docs"}}</a></li>
</ul>
    </div>
  </div>
</nav>
<script>
(function(){
var groups=document.querySelectorAll('.nav-group-header');
groups.forEach(function(h){
var body=h.nextElementSibling;
var height=body.scrollHeight+'px';
body.style.maxHeight=height;
h.addEventListener('click',function(){
h.classList.toggle('collapsed');
body.classList.toggle('collapsed');
if(body.classList.contains('collapsed'))body.style.maxHeight='0';
else body.style.maxHeight=body.scrollHeight+'px';
});
var sub=body.querySelectorAll('.nav-link');
sub.forEach(function(a){
if(a.classList.contains('active')){h.classList.remove('collapsed');body.classList.remove('collapsed');body.style.maxHeight=body.scrollHeight+'px'}
});
});
})();
</script>
{{end}}

{{define "sgroup"}}
<div class="nav-group-header{{if not .open}} collapsed{{end}}">
<span><i class="la {{.icon}} la-lg me-1"></i>{{.label}}</span>
<i class="la la-chevron-down"></i>
</div>
<ul class="navbar-nav nav-group-body{{if not .open}} collapsed{{end}}" style="max-height:{{if .open}}800{{else}}0{{end}}px;padding-left:4px">
{{end}}

{{define "egroup"}}
</ul>
{{end}}

{{define "landing"}}<!DOCTYPE html>
<html lang="{{.LangCode}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.AppName}} — {{T "landing_title"}}</title>
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
@media(max-width:768px){
.hero{padding:120px 20px 56px}
.hero h1{font-size:2rem}
.hero p{font-size:1rem}
.feature-grid{grid-template-columns:1fr}
.navbar .container{padding:0 16px}
.nav-links{gap:10px}
.nav-links a{font-size:13px}
.nav-links .btn-login{padding:7px 14px}
.features{padding:32px 16px 56px}
.features h2{font-size:1.5rem}
.demo-section{padding:32px 16px 56px}
.cta-banner{padding:44px 20px}
.cta-banner h2{font-size:1.5rem}
}
@media(max-width:400px){
.hero h1{font-size:1.6rem}
.nav-links a:not(.btn-login){display:none}
.nav-links .lang-switch{display:inline-block}
}
@media(pointer:coarse){.nav-links .btn-login,.hero .cta-group a,.cta-banner a{min-height:44px;display:inline-flex;align-items:center;justify-content:center}}
@media(prefers-reduced-motion:reduce){*,*::before,*::after{animation-duration:.01ms!important;transition-duration:.01ms!important;scroll-behavior:auto!important}}
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
<h1>{{T "landing_hero_title"}}</h1>
<p>{{T "landing_hero_desc"}}</p>
<div class="cta-group">
<a href="/register" class="btn-primary">{{T "landing_cta_free"}}</a>
<a href="/docs" class="btn-outline">{{T "landing_cta_docs"}}</a>
</div>
<div style="max-width:400px;margin:32px auto 0;background:#fff;border-radius:14px;padding:24px;box-shadow:0 4px 24px rgba(0,0,0,.08)">
<form method="post" action="/login/post">
<div style="margin-bottom:12px"><input type="email" name="email" class="form-control" placeholder="{{T "auth_email"}}" value="{{.AppEmail}}" style="border-radius:8px;padding:10px 14px;border:1px solid #ddd;width:100%;font-size:14px"></div>
<div style="margin-bottom:12px"><input type="password" name="password" class="form-control" placeholder="{{T "auth_password"}}" style="border-radius:8px;padding:10px 14px;border:1px solid #ddd;width:100%;font-size:14px"></div>
<button type="submit" style="width:100%;padding:10px;background:#4F46E5;color:#fff;border:none;border-radius:8px;font-weight:600;font-size:14px;cursor:pointer">{{T "auth_signin"}}</button>
</form>
<div style="text-align:center;margin-top:12px;font-size:12px;color:#999">{{T "landing_demo_label"}} <code style="background:#f0f0f0;padding:2px 6px;border-radius:4px">{{.AppEmail}}</code> / <code style="background:#f0f0f0;padding:2px 6px;border-radius:4px">password</code></div>
</div>
</section>

<section class="features">
<h2>{{T "landing_features_title"}}</h2>
<p class="subtitle">{{T "landing_features_subtitle"}}</p>
<div class="feature-grid">
<div class="feature-card"><i class="la la-comments"></i><h4>{{T "inbox_title"}}</h4><p>{{if eq .LangCode "id"}}Inbox real-time dengan SSE, reply langsung, group chat, filter private/group.{{else}}Real-time inbox with SSE, direct reply, group chat, private/group filter.{{end}}</p></div>
<div class="feature-card"><i class="la la-robot"></i><h4>{{T "landing_feat_ai_title"}}</h4><p>{{T "landing_feat_ai_desc"}}</p></div>
<div class="feature-card"><i class="la la-bullhorn"></i><h4>{{T "landing_feat_broadcast_title"}}</h4><p>{{T "landing_feat_broadcast_desc"}}</p></div>
<div class="feature-card"><i class="la la-whatsapp"></i><h4>{{T "landing_feat_multi_title"}}</h4><p>{{T "landing_feat_multi_desc"}}</p></div>
<div class="feature-card"><i class="la la-cloud"></i><h4>{{T "landing_feat_meta_title"}}</h4><p>{{T "landing_feat_meta_desc"}}</p></div>
<div class="feature-card"><i class="la la-clock"></i><h4>{{T "landing_feat_schedule_title"}}</h4><p>{{T "landing_feat_schedule_desc"}}</p></div>
<div class="feature-card"><i class="la la-paint-brush"></i><h4>{{T "landing_feat_whitelabel_title"}}</h4><p>{{T "landing_feat_whitelabel_desc"}}</p></div>
<div class="feature-card"><i class="la la-chart-bar"></i><h4>{{T "landing_feat_analytics_title"}}</h4><p>{{T "landing_feat_analytics_desc"}}</p></div>
</div>
</section>

<section class="demo-section">
<div class="container">
<h2>{{T "landing_demo_title"}}</h2>
<div class="demo-box" style="max-width:480px;margin:0 auto">
<div class="demo-row"><strong>{{T "landing_demo_admin_label"}}</strong> {{.AppEmail}} / password</div>
</div>
<div style="text-align:center;margin-top:24px">
<a href="/login" style="display:inline-block;padding:12px 28px;border-radius:10px;font-weight:600;text-decoration:none;background:#4F46E5;color:#fff">{{T "landing_go_dashboard"}}</a>
</div>
</div>
</section>

<section class="cta-banner">
<h2>{{T "landing_cta_headline"}}</h2>
<p>{{T "landing_cta_sub"}}</p>
<a href="/register">{{T "landing_cta_signup"}}</a>
</section>

<footer class="footer">&copy; 2026 {{.AppName}}. {{T "landing_footer"}}</footer>

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
@media(max-width:768px){.auth-left{display:none}.auth-right{min-width:100%;padding:28px 20px}.auth-card h2{font-size:1.5rem}.auth-card input{font-size:16px}}
@media(pointer:coarse){.auth-card .btn-submit{min-height:46px}}
@media(prefers-reduced-motion:reduce){*,*::before,*::after{animation-duration:.01ms!important;transition-duration:.01ms!important}}
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
<div class="form-group"><label>{{T "auth_password"}}</label><input type="password" name="password" placeholder="••••••••" required></div>
<button type="submit" class="btn-submit">{{T "auth_signin"}}</button>
</form>
<div class="auth-divider"><span>{{T "auth_or"}}</span></div>
<div class="demo-box">
<div class="demo-title">{{T "auth_demo"}}</div>
<div class="demo-row"><strong>{{T "landing_demo_admin_label"}}</strong> {{.AppEmail}} / password</div>
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
<div class="form-group"><label>Kode Voucher (Opsional)</label><input name="voucher" placeholder="Masukkan kode voucher"></div>
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
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_in"}}</h4><a href="/received" class="btn btn-sm btn-white">{{T "btn_all"}}</a></div>
<div class="table-responsive"><table class="table table-sm table-nowrap card-table mb-0"><thead><tr><th>{{T "col_from"}}</th><th>{{T "col_message"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
{{range .Received}}<tr><td><strong>{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}</strong></td><td class="text-muted small" style="max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{.Message}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div>
</div>
<div class="col-12 col-lg-6">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dash_recent_out"}}</h4><a href="/sent" class="btn btn-sm btn-white">{{T "btn_all"}}</a></div>
<div class="table-responsive"><table class="table table-sm table-nowrap card-table mb-0"><thead><tr><th>{{T "col_to"}}</th><th>{{T "col_message"}}</th><th>{{T "col_status"}}</th></tr></thead><tbody>
{{range .Sent}}<tr><td><strong>+{{.Phone}}</strong></td><td class="text-muted small" style="max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{.Message}}</td><td><span class="badge badge-soft-success">{{.Status}}</span></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div>
</div>
</div>
<script>
new Chart(document.getElementById('msgChart'),{type:'line',data:{labels:[{{.ChartLabels}}],datasets:[{label:'{{T "chart_sent"}}',data:[{{.ChartSent}}],borderColor:'#4F46E5',backgroundColor:'rgba(79,70,229,.1)',fill:true,tension:.3,pointRadius:2,pointHoverRadius:5},{label:'{{T "chart_received"}}',data:[{{.ChartReceived}}],borderColor:'#10B981',backgroundColor:'rgba(16,185,129,.1)',fill:true,tension:.3,pointRadius:2,pointHoverRadius:5}]},options:{responsive:true,interaction:{intersect:false,mode:'index'},plugins:{legend:{position:'bottom'}},scales:{y:{beginAtZero:true,grid:{color:'rgba(0,0,0,.05)'}},x:{grid:{display:false}}}}})
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
              {{if eq .Status "connected"}}<span class="badge badge-soft-success"><i class="la la-check-circle me-1"></i>{{T "wa_connected"}}</span>
              {{else if eq .Status "qr"}}<span class="badge badge-soft-warning"><i class="la la-qrcode me-1"></i>{{T "wa_scanqr"}}</span>
              {{else}}<span class="badge badge-soft-danger"><i class="la la-times-circle me-1"></i>{{T "wa_disconnected"}}</span>{{end}}
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
            <div class="form-group"><label>{{T "send_from"}}</label><select name="account_phone" class="form-control" required {{if not .HasConnected}}disabled{{end}}><option value="">{{T "send_select_sender"}}</option>{{range .ConnectedAccounts}}{{if eq .Status "connected"}}<option value="+{{.Phone}}">+{{.Phone}}</option>{{end}}{{end}}</select>{{if not .HasConnected}}<small class="form-text text-muted">{{T "send_connect_hint"}}</small>{{end}}</div>
            <div class="form-group"><label>{{T "send_to"}}</label><input name="phone" class="form-control" placeholder="628123456789" value="{{.SendTo}}" required><small class="form-text text-muted">{{T "send_to_hint"}}</small></div>
            <div class="form-group"><label>{{T "send_message"}}</label><textarea name="message" class="form-control" rows="4" placeholder="{{T "send_message_ph"}}" required></textarea></div>
             <button class="btn btn-primary lift" {{if ne .Status "connected"}}disabled{{end}}><i class="la la-paper-plane me-1"></i> {{T "send_btn"}}</button>
             {{if ne .Status "connected"}}<span class="text-muted ms-2">{{T "send_connect_first"}}</span>{{end}}
           </form>
         </div>
       </div>
       <div class="card mt-3">
         <div class="card-header"><h4 class="card-header-title"><i class="la la-paperclip me-1"></i> Kirim Media / Attachment</h4></div>
         <div class="card-body">
           <form method="post" action="/send/media" enctype="multipart/form-data">
             <div class="form-group"><label>{{T "send_from"}}</label><select name="account_phone" class="form-control" required {{if not .HasConnected}}disabled{{end}}><option value="">{{T "send_select_sender"}}</option>{{range .ConnectedAccounts}}{{if eq .Status "connected"}}<option value="+{{.Phone}}">+{{.Phone}}</option>{{end}}{{end}}</select></div>
             <div class="form-group"><label>{{T "send_to"}}</label><input name="phone" class="form-control" placeholder="628123456789" value="{{.SendTo}}" required></div>
             <div class="form-group"><label>{{T "col_type"}}</label><select name="media_type" class="form-control"><option value="image">🖼️ Image</option><option value="video">🎬 Video</option><option value="document">📄 Document</option><option value="audio">🎵 Audio</option></select></div>
             <div class="form-group"><label>{{T "col_file"}}</label><input type="file" name="media_file" class="form-control" required></div>
             <div class="form-group"><label>{{T "send_message"}} <small class="text-muted">(caption)</small></label><input name="caption" class="form-control" placeholder="Teks caption (opsional)"></div>
             <button class="btn btn-warning lift" {{if ne .Status "connected"}}disabled{{end}}><i class="la la-cloud-upload me-1"></i> Kirim Media</button>
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
<style>
.autoreply-tabs{display:flex;flex-wrap:wrap;border-bottom:2px solid #e0e0e0;margin-bottom:16px}
.autoreply-tabs .at{background:none;border:none;padding:10px 16px;font-size:13px;font-weight:600;color:#6e788c;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-2px;white-space:nowrap}
.autoreply-tabs .at:hover{color:#152e4d}
.autoreply-tabs .at.active{color:#2c7be5;border-bottom-color:#2c7be5}
.at-panel{display:none}
.at-panel.active{display:block}
</style>
<div class="autoreply-tabs">
<button class="at active" onclick="var p=document.querySelectorAll('.at-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('at-welcome').classList.add('active');var b=this.parentElement.querySelectorAll('.at');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-hand-sparkles me-1"></i>{{T "at_tab_welcome_fallback"}}</button>
<button class="at" onclick="var p=document.querySelectorAll('.at-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('at-rules').classList.add('active');var b=this.parentElement.querySelectorAll('.at');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-reply me-1"></i>{{T "nav_autoreply"}}</button>
<button class="at" onclick="var p=document.querySelectorAll('.at-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('at-faq').classList.add('active');var b=this.parentElement.querySelectorAll('.at');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-question-circle me-1"></i>{{T "ar_faq_tab"}}</button>
<button class="at" onclick="var p=document.querySelectorAll('.at-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('at-kb').classList.add('active');var b=this.parentElement.querySelectorAll('.at');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-book me-1"></i>{{T "nav_knowledge"}}</button>
<button class="at" onclick="var p=document.querySelectorAll('.at-panel');for(var i=0;i<p.length;i++)p[i].classList.remove('active');document.getElementById('at-ai').classList.add('active');var b=this.parentElement.querySelectorAll('.at');for(var i=0;i<b.length;i++)b[i].classList.remove('active');this.classList.add('active')"><i class="la la-robot me-1"></i>{{T "at_tab_ai_setup"}}</button>
</div>

<div class="at-panel active" id="at-welcome">
<form method="post" action="/settings">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "set_welcome_title"}}</h4><div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="welcome_enabled" {{if .WelcomeEnabled}}checked{{end}}></div></div><div class="card-body">
<label>{{T "set_welcome_msg"}}</label>
<textarea name="welcome_message" class="form-control" rows="2" placeholder="{{T "set_welcome_ph"}}">{{.WelcomeMessage}}</textarea>
<small class="form-text text-muted">{{T "set_vars_hint"}}</small>
</div></div>
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "set_fallback_title"}}</h4><div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="fallback_enabled" {{if .FallbackEnabled}}checked{{end}}></div></div><div class="card-body">
<label>{{T "set_fallback_msg"}}</label>
<textarea name="fallback_message" class="form-control" rows="2" placeholder="{{T "set_fallback_ph"}}">{{.FallbackMessage}}</textarea>
<small class="form-text text-muted">{{T "set_fallback_hint"}}</small>
</div></div>
<div class="text-end"><button class="btn btn-primary"><i class="la la-save me-1"></i> {{T "btn_save"}}</button></div>
</form>
</div>

<div class="at-panel" id="at-rules">
<div class="row">
<div class="col-12 col-lg-5">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ar_add_title"}}</h4></div>
<div class="card-body"><form method="post" action="/autoreply/add" enctype="multipart/form-data">
<div class="form-group"><label>{{T "ar_matchtype"}}</label><select name="match" class="form-control" onchange="onMatchTypeChange(this.value)"><option value="contains">{{T "ar_contains"}}</option><option value="exact">{{T "ar_exact"}}</option><option value="starts_with">{{T "ar_starts"}}</option><option value="ai">{{T "ar_ai_type"}}</option></select></div>
<div id="keywordGroup"><div class="form-group"><label>{{T "ar_keyword"}}</label><input name="keyword" class="form-control" placeholder="halo, hi, menu"></div></div>
<div id="faqGroup" style="display:none"><div class="form-group"><label>{{T "ar_faq"}}</label><textarea name="faq" class="form-control" rows="5" placeholder="Apa produk?|Software WA marketing"></textarea></div></div>
<div class="form-group"><label>{{T "ar_reply"}}</label><textarea name="reply" class="form-control" rows="3" placeholder="{{T "ar_reply_ph"}}"></textarea></div>
<div class="bg-light border rounded p-3 mb-3">
<div class="form-check"><input class="form-check-input" type="checkbox" name="use_ai" value="1" id="useAiCheck"><label class="form-check-label" for="useAiCheck">{{T "ar_use_ai"}}</label></div>
<div class="form-group mt-2" id="aiKeyGroup" style="display:none"><label>{{T "ar_ai_key"}}</label><select name="ai_key_id" class="form-control">{{range .AiKeys}}<option value="{{.ID}}">{{.Name}} ({{.Provider}})</option>{{end}}</select></div>
</div>
<div class="form-group"><label>{{T "ar_wa_number"}}</label><div class="border rounded p-2" style="max-height:120px;overflow-y:auto">{{range .ConnectedAccounts}}{{if .Phone}}<div class="form-check form-check-inline small"><input class="form-check-input" type="checkbox" name="account_ids" value="+{{.Phone}}" id="a_{{.Phone}}"><label for="a_{{.Phone}}">+{{.Phone}}</label></div>{{end}}{{end}}</div></div>
<div class="form-group"><label><i class="la la-paperclip me-1"></i> Media / File <small class="text-muted">(opsional)</small></label><input type="file" name="media_file" class="form-control" accept="image/*,video/*,audio/*,.pdf,.doc,.docx,.xls,.xlsx"></div>
<div class="mb-2"><label class="field-label">{{T "ar_training_camp"}}</label><select name="training_id" class="form-control form-control-sm"><option value="0">{{T "ar_default_global"}}</option>{{range .AiTrainings}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
<button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
</form></div></div></div>
<div class="col-12 col-lg-7">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "ar_list_title"}}</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "ar_keyword"}}</th><th>{{T "ar_reply_label"}}</th><th>{{T "col_status"}}</th><th></th></tr></thead><tbody>
{{range .AutoReplies}}<tr><td>{{.ID}}</td><td><strong>{{.Keyword}}</strong></td><td>{{if .MediaURL}}<span class="badge bg-info bg-opacity-10 text-info me-1 small" title="{{.MediaType}}: {{.MediaURL}}">📎</span>{{end}}{{if .UseAI}}<span class="badge bg-warning bg-opacity-10 text-warning me-1 small">AI</span>{{end}}{{.Reply}}</td><td>{{if .IsActive}}<span class="badge bg-success bg-opacity-10 text-success small">{{T "ar_on"}}</span>{{else}}<span class="badge bg-danger bg-opacity-10 text-danger small">{{T "ar_off"}}</span>{{end}}</td><td><a class="btn btn-sm btn-white px-2" href="/autoreply?edit={{.ID}}">{{T "ar_edit"}}</a><form method="post" action="/autoreply/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white px-2">{{if .IsActive}}{{T "ar_off"}}{{else}}{{T "ar_on"}}{{end}}</button></form><form method="post" action="/autoreply/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm text-danger px-2">&times;</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted py-3 text-center">{{T "ar_empty"}}</td></tr>{{end}}
</tbody></table></div></div></div></div>
</div>

<div class="at-panel" id="at-faq">
<div class="row">
<div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "faq_add_title"}}</h4></div><div class="card-body">
<form method="post" action="/faq/add"><div class="form-group"><label>{{T "col_question"}}</label><input name="question" class="form-control" placeholder="{{T "faq_question_ph"}}" required></div><div class="form-group"><label>{{T "col_answer"}}</label><textarea name="answer" class="form-control" rows="3" placeholder="{{T "faq_answer_ph"}}" required></textarea></div><button class="btn btn-primary"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></form>
<hr><h6>{{T "kb_import"}}</h6><form method="post" action="/faq/import" enctype="multipart/form-data"><div class="form-group"><label>{{T "faq_import_csv"}}</label><input type="file" name="file" class="form-control" accept=".csv" required></div><button class="btn btn-outline-primary btn-sm"><i class="la la-upload me-1"></i> {{T "btn_import"}}</button></form>
</div></div></div>
<div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "faq_list_title"}} <small class="text-muted">— {{T "faq_subtitle"}}</small></h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_question"}}</th><th>{{T "col_answer"}}</th><th></th></tr></thead><tbody>
{{range .FAQ}}<tr><td>{{.id}}</td><td>{{.question}}</td><td>{{.answer}}</td><td><form method="post" action="/faq/delete" style="display:inline"><input type="hidden" name="id" value="{{.id}}"><button class="btn btn-sm btn-danger">&times;</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">{{T "faq_empty"}}</td></tr>{{end}}</tbody></table></div></div></div></div>
</div>

<div class="at-panel" id="at-kb">
<div class="row">
<div class="col-12 col-lg-5">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "kb_add"}}</h4></div><div class="card-body">
<form method="post" action="/knowledge/add">
<div class="form-group"><label>{{T "kb_title"}}</label><input name="title" class="form-control" placeholder="{{T "ar_faq_tab"}}" required></div>
<div class="form-group"><label>{{T "kb_question"}}</label><input name="question" class="form-control" placeholder="{{T "kb_question_dot"}}..." required></div>
<div class="form-group"><label>{{T "kb_answer"}}</label><textarea name="answer" class="form-control" rows="3" placeholder="{{T "kb_answer_dot"}}..." required></textarea></div>
<div class="form-group"><label>{{T "kb_category"}}</label><input name="category" class="form-control" placeholder="{{T "kb_placeholder_category"}}"></div>
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
<label class="text-muted small mb-2 d-block">📄 {{T "kb_upload_pdf"}}</label>
<div class="input-group"><input type="text" name="title" class="form-control" placeholder="{{T "kb_placeholder_title"}}"><input type="file" name="file" class="form-control" accept=".pdf" required><button class="btn btn-white">{{T "kb_upload"}}</button></div>
<small class="form-text text-muted">{{T "kb_pdf_hint"}}</small>
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
</div>

<div class="at-panel" id="at-ai">
<div class="row">
<div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "aik_add"}}</h4></div><div class="card-body">
<form method="post" action="/ai/keys/add">
<div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
<div class="form-group"><label>{{T "col_provider"}}</label><select name="provider" class="form-control"><option value="openai">OpenAI</option><option value="geminiai">Gemini</option><option value="claudeai">Claude</option><option value="deepseekai">DeepSeek</option></select></div>
<div class="form-group"><label>{{T "ar_ai_model"}}</label><input name="model" class="form-control" placeholder="gpt-4o"></div>
<div class="form-group"><label>{{T "aik_api_key"}}</label><input name="apikey" class="form-control" required></div>
<div class="form-group"><label>{{T "ar_system_prompt"}}</label><textarea name="system_prompt" class="form-control" rows="3"></textarea></div>
<button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
</form></div></div></div>
<div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_keys"}}</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>Provider</th><th>Model</th><th>{{T "col_action"}}</th></tr></thead><tbody>
{{range .AiKeys}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><span class="badge badge-soft-secondary">{{.Provider}}</span></td><td>{{.Model}}</td><td><form method="post" action="/ai/keys/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div></div>
</div>
<div class="row mt-3">
<div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "aip_add"}}</h4></div><div class="card-body">
<form method="post" action="/ai/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><div class="form-group"><label>{{T "ar_ai_endpoint"}}</label><input name="endpoint" class="form-control" placeholder="https://..."></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></form></div></div></div>
<div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_plugins"}}</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "ar_ai_endpoint"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
{{range .AiPlugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Endpoint}}</td><td><form method="post" action="/ai/plugins/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
</tbody></table></div></div></div>
</div>
</div>
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
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "set_tab_branding"}}</h4></div>
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
<div class="card-header"><h4 class="card-header-title">{{T "set_group_title"}}</h4>
<div class="form-check form-switch"><input class="form-check-input" type="checkbox" name="reply_in_group" {{if .ReplyInGroup}}checked{{end}}></div>
</div>
<div class="card-body"><small class="form-text text-muted">{{T "set_group_hint"}}</small></div>
</div>

<div class="card mt-3">
<div class="card-header"><h4 class="card-header-title"><i class="la la-clock me-1"></i> {{T "set_autoclose_title"}}</h4></div>
<div class="card-body">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>{{T "set_autoclose_hours"}}</label><input type="number" name="auto_close_hours" class="form-control" value="{{.AutoCloseHours}}"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "set_autoclose_followup"}}</label><input name="auto_close_message" class="form-control" placeholder="{{T "set_autoclose_ph"}}" value="{{.AutoCloseMessage}}"></div></div>
</div>
</div></div>
<div class="card-body">
<div class="row">
<div class="col-md-4"><div class="form-group"><label>{{T "set_rate_max_daily"}}</label><input type="number" name="rate_max_daily" class="form-control" value="{{.RateMaxDaily}}"></div></div>
<div class="col-md-4"><div class="form-group"><label>{{T "set_rate_min"}}</label><input type="number" name="rate_random_min" class="form-control" value="{{.RateRandomMin}}"></div></div>
<div class="col-md-4"><div class="form-group"><label>{{T "set_rate_max"}}</label><input type="number" name="rate_random_max" class="form-control" value="{{.RateRandomMax}}"></div></div>
</div>
<small class="form-text text-muted">{{T "set_rate_hint"}}</small>
</div></div>
</div>
</div>

<div class="st-panel" id="st-system">
<div class="card"><div class="card-header"><h4 class="card-header-title">{{T "set_system_title"}}</h4></div>
<div class="card-body">
<div class="form-group"><label>{{T "set_registrations"}}</label>
<select name="registrations" class="form-control"><option value="1" {{if .Registrations}}selected{{end}}>{{T "set_enabled"}}</option><option value="0" {{if not .Registrations}}selected{{end}}>{{T "set_disabled"}}</option></select></div>
<div class="form-group"><label>{{T "set_listen_addr"}}</label><input class="form-control" value="0.0.0.0:8080" disabled><small class="form-text text-muted">{{T "set_addr_hint"}}</small></div>
<div class="form-group"><label>{{T "set_mysql_conn"}}</label><input class="form-control" value="***" disabled><small class="form-text text-muted">{{T "set_mysql_hint"}}</small></div>
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
           <div class="form-group"><label>{{T "nav_tags"}}</label><div class="border rounded p-2" style="max-height:100px;overflow-y:auto">{{range .Tags}}<div class="form-check form-check-inline small"><input class="form-check-input" type="checkbox" name="tag_ids" value="{{.ID}}" id="addtag_{{.ID}}"><label for="addtag_{{.ID}}">{{.Name}}</label></div>{{else}}<small class="text-muted">Belum ada tag</small>{{end}}</div></div>
           <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
      <div class="card mt-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-upload me-1"></i> {{T "ct_import_csv"}}</h4></div>
        <div class="card-body"><form method="post" action="/contacts/import" enctype="multipart/form-data">
          <div class="form-group"><label>{{T "ct_upload_csv"}}</label><input type="file" name="file" class="form-control" accept=".csv" required></div>
          <small class="form-text text-muted mb-2 d-block">{{T "ct_import_hint"}}</small>
          <button class="btn btn-white lift"><i class="la la-cloud-upload me-1"></i> {{T "btn_import"}}</button>
        </form></div>
      </div>
    </div>
    {{if .EditID}}
    <div class="col-12 col-lg-4">
      <div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> {{T "ar_edit"}}</h4></div>
        <div class="card-body"><form method="post" action="/contacts/edit">
          <input type="hidden" name="id" value="{{.EditID}}">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div>
          <div class="form-group"><label>{{T "col_from"}}</label><input name="phone" class="form-control" value="{{.EditPhone}}" required></div>
           <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
           <div class="form-group"><label>{{T "nav_tags"}}</label><div class="border rounded p-2" style="max-height:100px;overflow-y:auto">{{range .Tags}}<div class="form-check form-check-inline small"><input class="form-check-input" type="checkbox" name="tag_ids" value="{{.ID}}" id="tag_{{.ID}}"><label for="tag_{{.ID}}">{{.Name}}</label></div>{{else}}<small class="text-muted">Belum ada tag</small>{{end}}</div></div>
           <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "set_save"}}</button> <a href="/contacts" class="btn btn-white ms-2">{{T "ar_cancel"}}</a>
        </form></div>
      </div>
    </div>
    {{end}}
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header d-flex justify-content-between align-items-center"><h4 class="card-header-title mb-0">{{T "nav_contacts_saved"}}</h4><div><a href="/contacts/export" class="btn btn-sm btn-white me-1"><i class="la la-download me-1"></i> {{T "ct_export_csv"}}</a><button class="btn btn-sm btn-danger" onclick="bulkDeleteContacts()"><i class="la la-trash me-1"></i> {{T "ar_delete"}}</button></div></div>
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
      if (ids.length === 0) { alert('{{T "ct_alert_select"}}'); return; }
      if (!confirm('{{T "ct_confirm_bulk_delete"}}'.replace('%d', ids.length))) return;
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
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "tag_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/tags/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="VIP" required></div>
          <div class="form-group"><label>{{T "tag_color"}}</label><input type="color" name="color" class="form-control form-control-color" value="#2c7be5" style="height:40px;padding:4px"></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "tag_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "tag_color"}}</th><th>{{T "col_name"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
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
           {{if .Tags}}<div class="form-group"><label>{{T "bc_tags"}} <small class="text-muted">— {{T "bc_tags_hint"}}</small></label><select name="tags" class="form-control" multiple>{{range .Tags}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select><small class="form-text text-muted">{{T "bc_tags_desc"}}</small></div>{{end}}
           <div class="form-group"><label>{{T "bc_direct_numbers"}} <small class="text-muted">— {{T "bc_direct_numbers_hint"}}</small></label><textarea name="numbers" class="form-control" rows="4" placeholder="628123456789&#10;628987654321&#10;..."></textarea><small class="form-text text-muted">{{T "bc_direct_numbers_desc"}} <form method="post" action="/validate" style="display:inline" target="_blank"><input type="hidden" name="numbers" value="" id="validateInput"><button type="button" class="btn btn-sm btn-outline-warning" onclick="document.getElementById('validateInput').value=document.querySelector('textarea[name=numbers]').value;this.form.submit()">{{T "bc_validate_btn"}}</button></form></small></div>
           <div class="form-group"><label><i class="la la-image me-1"></i> {{T "bc_media"}} <small class="text-muted">— {{T "bc_media_hint"}}</small></label><input type="file" name="media_file" class="form-control" accept="image/*,video/*,.pdf,.doc,.docx,.xls,.xlsx"><small class="form-text text-muted">{{T "bc_media_desc"}}</small></div>
           <div class="form-group"><label>{{T "bc_account"}}</label><div class="border rounded p-2" style="max-height:160px;overflow-y:auto">{{range .ConnectedAccounts}}{{if .Phone}}<div class="form-check"><input class="form-check-input" type="checkbox" name="account_ids" value="+{{.Phone}}" id="bc_{{.Phone}}"><label class="form-check-label small" for="bc_{{.Phone}}">+{{.Phone}}</label></div>{{end}}{{end}}{{if not .HasConnected}}<small class="text-muted">{{T "bc_no_connected"}}</small>{{end}}</div><small class="form-text text-muted">{{T "bc_account_hint"}}</small></div>
           <div class="form-group"><label>{{T "bc_send_mode"}}</label><div class="border rounded p-2"><div class="form-check"><input class="form-check-input" type="radio" name="send_mode" value="round_robin" id="mode_rr" checked><label class="form-check-label" for="mode_rr"><strong>{{T "bc_mode_rr"}}</strong> <small class="text-muted">— kirim bergantian merata ke tiap nomor</small></label></div><div class="form-check mt-1"><input class="form-check-input" type="radio" name="send_mode" value="random" id="mode_rand"><label class="form-check-label" for="mode_rand"><strong>{{T "bc_mode_random"}}</strong> <small class="text-muted">— kirim acak ke nomor manapun</small></label></div></div></div>
           <div class="form-group"><label>Interval (detik) <small class="text-muted">jeda antar pesan</small></label><input name="interval" type="number" class="form-control" value="300" min="30" placeholder="300-400"></div>
           {{if .MetaAccounts}}
           <div class="form-group"><label><i class="la la-cloud me-1"></i> {{T "bc_meta_api"}} <small class="text-muted">— kirim lewat Cloud API</small></label><select name="meta_account_id" class="form-control" onchange="toggleMetaTemplate(this)"><option value="0">{{T "bc_meta_none"}}</option>{{range .MetaAccounts}}<option value="{{.ID}}">{{.Name}} ({{.PhoneNumberID}})</option>{{end}}</select></div>
            <div class="form-group" id="metaTemplateGroup" style="display:none"><label>{{T "bc_meta_template"}} <small class="text-muted">— opsional</small></label><select name="meta_template" class="form-control"><option value="">{{T "bc_meta_plain"}}</option>{{range .MetaTemplates}}<option value="{{.Name}}">[Meta] {{.Name}} ({{.Language}})</option>{{end}}</select><small class="form-text text-muted">Jika dipilih, template akan dipakai. Variabel dari pesan akan masuk ke parameter.</small></div>
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
           {{range .Campaigns}}<tr><td>{{.ID}}</td><td>{{.Name}}{{if .MetaAccountID}} <span class="badge badge-soft-primary" style="font-size:9px">Meta</span>{{end}}</td><td>{{.Sent}}/{{.Total}}</td><td>{{if eq .Status "running"}}<span class="badge badge-soft-primary">{{T "bc_status_running"}}</span>{{else if eq .Status "paused"}}<span class="badge badge-soft-warning">{{T "bc_status_paused"}}</span>{{else if eq .Status "done"}}<span class="badge badge-soft-success">{{T "bc_status_done"}}</span>{{else}}<span class="badge badge-soft-secondary">{{.Status}}</span>{{end}}</td><td class="text-nowrap">
             {{if eq .Status "running"}}<form method="post" action="/broadcast/pause" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-warning" title=">{{T "bc_pause"}}"><i class="la la-pause"></i></button></form>{{end}}
             {{if eq .Status "paused"}}<form method="post" action="/broadcast/pause" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-success" title=">{{T "bc_resume"}}"><i class="la la-play"></i></button></form>{{end}}
             {{if eq .Status "done"}}<form method="post" action="/broadcast/retry" style="display:inline" onsubmit="return confirm('{{T "bc_retry_confirm"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-info" title=">{{T "bc_retry"}}"><i class="la la-redo"></i></button></form>{{end}}
             {{if eq .Status "stopped"}}<form method="post" action="/broadcast/retry" style="display:inline" onsubmit="return confirm('{{T "bc_retry_confirm"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-info" title=">{{T "bc_retry"}}"><i class="la la-redo"></i></button></form>{{end}}
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
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "drip_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/drips/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="{{T "drip_name_ph"}}" required></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      {{range .Drips}}
      <div class="card mb-3">
        <div class="card-header d-flex justify-content-between align-items-center">
          <div><h4 class="card-header-title mb-0">{{.Name}}</h4><small class="text-muted">{{len .Steps}} {{T "drip_steps"}} &middot; {{if eq .Status "active"}}<span class="text-success">{{T "drip_active"}}</span>{{else}}<span class="text-muted">{{T "drip_inactive"}}</span>{{end}}</small></div>
          <div>
            <form method="post" action="/drips/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}">{{if eq .Status "active"}}<button class="btn btn-sm btn-warning">{{T "drip_pause"}}</button>{{else}}<button class="btn btn-sm btn-success">{{T "drip_resume"}}</button>{{end}}</form>
            <form method="post" action="/drips/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form>
          </div>
        </div>
        <div class="table-responsive"><table class="table table-sm card-table mb-0"><thead><tr><th>#</th><th>{{T "col_delay"}}</th><th>{{T "col_message"}}</th><th></th></tr></thead><tbody>
          {{range $i, $s := .Steps}}<tr><td>{{add $i 1}}</td><td>{{if eq $i 0}}{{T "drip_instant"}}{{else}}{{$s.DelayMinutes}} {{T "drip_min"}}{{end}}</td><td>{{$s.Message}}</td><td><form method="post" action="/drips/step/delete" style="display:inline"><input type="hidden" name="id" value="{{$s.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-times"></i></button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">{{T "drip_no_steps"}}</td></tr>{{end}}
        </tbody></table></div>
        <div class="card-body border-top"><form method="post" action="/drips/step/add" class="row g-2">
          <input type="hidden" name="drip_id" value="{{.ID}}">
          <input type="hidden" name="sort_order" value="{{len .Steps}}">
          <div class="col-md-2"><input type="number" name="delay" class="form-control form-control-sm" placeholder="{{T "drip_delay_ph"}}" value="0"></div>
          <div class="col-md-7"><input type="text" name="message" class="form-control form-control-sm" placeholder="{{T "drip_msg_ph"}}" required></div>
          <div class="col-md-3"><button class="btn btn-sm btn-primary w-100"><i class="la la-plus"></i> {{T "drip_add_step"}}</button></div>
        </form></div>
      </div>
      {{else}}
      <div class="card"><div class="card-body text-center text-muted py-5">{{T "drip_empty"}}</div></div>
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
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "canned_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/canned/add">
          <div class="form-group"><label>{{T "canned_shortcut"}} <small class="text-muted">{{T "canned_shortcut_hint"}}</small></label><input name="shortcut" class="form-control" placeholder="/salam"></div>
          <div class="form-group"><label>{{T "canned_title"}}</label><input name="name" class="form-control" placeholder="Salam Pembuka" required></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "canned_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>{{T "col_shortcut"}}</th><th>{{T "col_name"}}</th><th>{{T "col_message"}}</th><th></th></tr></thead><tbody>
          {{range .Canned}}<tr><td><code>{{.Shortcut}}</code></td><td>{{.Name}}</td><td style="max-width:250px">{{.Message}}</td><td><form method="post" action="/canned/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "tracker"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-link me-1"></i> {{T "tracker_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_url"}}</th><th>{{T "col_campaign"}}</th><th>{{T "col_phone"}}</th><th>{{T "col_clicked"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .LClicks}}<tr><td>{{.ID}}</td><td style="max-width:300px;word-break:break-all"><a href="/track/{{.Token}}" target="_blank">{{.URL}}</a></td><td>{{if .CampaignID}}#{{.CampaignID}}{{else}}-{{end}}</td><td>{{.Phone}}</td><td>{{if .Clicked}}<span class="badge badge-soft-success">{{T "tracker_yes"}}</span>{{else}}<span class="badge badge-soft-secondary">{{T "tracker_no"}}</span>{{end}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">{{T "tracker_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "abtests"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title"><i class="la la-balance-scale me-1"></i> {{T "ab_title"}}</h4><button class="btn btn-sm btn-primary" data-toggle="modal" data-target="#abModal"><i class="la la-plus"></i> {{T "ab_add"}}</button></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "ab_col_campaign"}}</th><th>{{T "ab_col_var_a"}}</th><th>{{T "ab_col_var_b"}}</th><th>{{T "ab_col_a_sent"}}</th><th>{{T "ab_col_b_sent"}}</th></tr></thead><tbody>
      {{range .ABTests}}<tr><td>{{.ID}}</td><td>#{{.CampaignID}}</td><td style="max-width:200px">{{.VariantA}}</td><td style="max-width:200px">{{.VariantB}}</td><td>{{.ASent}}</td><td>{{.BSent}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">{{T "ab_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
  <div class="modal fade" id="abModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><form method="post" action="/ab-tests/add">
    <div class="modal-header"><h5>{{T "ab_modal_title"}}</h5><button type="button" class="close" data-dismiss="modal"><span>&times;</span></button></div>
    <div class="modal-body">
      <div class="form-group"><label>{{T "ab_campaign_id"}}</label><input type="number" name="campaign_id" class="form-control" required></div>
      <div class="form-group"><label>{{T "ab_var_a"}}</label><textarea name="variant_a" class="form-control" rows="2" required></textarea></div>
      <div class="form-group"><label>{{T "ab_var_b"}}</label><textarea name="variant_b" class="form-control" rows="2" required></textarea></div>
    </div>
    <div class="modal-footer"><button class="btn btn-primary">{{T "ab_create"}}</button></div>
  </form></div></div></div>
{{end}}

{{if eq .Page "store"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card mb-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "store_add_product"}}</h4></div>
        <div class="card-body"><form method="post" action="/store/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>{{T "store_desc"}}</label><textarea name="desc" class="form-control" rows="2"></textarea></div>
          <div class="form-group"><label>{{T "store_price"}}</label><input name="price" type="number" step="0.01" class="form-control" required></div>
          <div class="form-group"><label>{{T "store_image"}}</label><input name="image_url" class="form-control"></div>
          <div class="form-group"><label>{{T "store_category"}}</label><select name="category" class="form-control">{{range .Categories}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div>
          <div class="form-group"><label>{{T "store_stock"}}</label><input name="stock" type="number" class="form-control" value="0"></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "store_categories"}}</h4></div>
        <div class="card-body"><form method="post" action="/store/category/add"><div class="input-group"><input name="name" class="form-control" placeholder="{{T "store_cat_name_ph"}}"><button class="btn btn-primary">{{T "ar_add_btn"}}</button></div></form>
          <div class="mt-2">{{range .Categories}}<span class="badge badge-soft-primary me-1 mb-1">{{.Name}} <form method="post" action="/store/category/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button style="background:none;border:none;color:inherit;cursor:pointer;font-size:10px">&times;</button></form></span>{{end}}</div>
        </div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "store_products"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "store_col_image"}}</th><th>{{T "col_name"}}</th><th>{{T "store_col_price"}}</th><th>{{T "store_col_category"}}</th><th>{{T "store_col_stock"}}</th><th></th></tr></thead><tbody>
          {{range .Products}}<tr><td>{{.ID}}</td><td>{{if .ImageURL}}<img src="{{.ImageURL}}" style="width:40px;height:40px;object-fit:cover;border-radius:6px">{{else}}-{{end}}</td><td>{{.Name}}</td><td>{{.Price}}</td><td>{{.Category}}</td><td>{{.Stock}}</td><td><form method="post" action="/store/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="7" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "orders"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "orders_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_phone"}}</th><th>{{T "col_name"}}</th><th>{{T "col_product"}}</th><th>{{T "col_qty"}}</th><th>{{T "col_total"}}</th><th>{{T "col_status"}}</th><th></th></tr></thead><tbody>
      {{range .Orders}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Name}}</td><td>#{{.ProductID}}</td><td>{{.Quantity}}</td><td>{{.Total}}</td><td><span class="badge badge-soft-{{if eq .Status "new"}}warning{{else if eq .Status "paid"}}success{{else}}secondary{{end}}">{{.Status}}</span></td><td>
        <form method="post" action="/store/orders/update" style="display:inline" class="d-flex gap-1"><input type="hidden" name="id" value="{{.ID}}"><select name="status" class="form-select form-select-sm" style="width:auto" onchange="this.form.submit()"><option value="new" selected>{{T "order_new"}}</option><option value="confirmed">{{T "order_confirmed"}}</option><option value="paid">{{T "order_paid"}}</option><option value="shipped">{{T "order_shipped"}}</option><option value="cancelled">{{T "order_cancelled"}}</option></select></form>
      </td></tr>{{else}}<tr><td colspan="8" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "forms"}}
  <div class="row">
    <div class="col-12 col-lg-4">
       <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "form_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/forms/add">
          <div class="form-group"><label>{{T "form_name"}}</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>{{T "form_fields"}}</label><textarea name="fields" class="form-control" rows="6" placeholder='[{"label":"Nama","type":"text"},{"label":"Email","type":"text"},{"label":"Rating","type":"number"}]'></textarea><small class="form-text text-muted">{{T "form_fields_hint"}}</small></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> {{T "form_create"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "form_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_fields"}}</th><th></th></tr></thead><tbody>
          {{range .Forms}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><code>{{.Fields}}</code></td><td><a href="/forms/submissions?form_id={{.ID}}" class="btn btn-sm btn-white">{{T "form_data_btn"}}</a> <form method="post" action="/forms/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "submissions"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">{{T "form_submissions"}}</h4>
    <select class="form-select form-select-sm" style="width:auto" onchange="window.location='?form_id='+this.value"><option value="">{{T "form_select"}}</option>{{range .Forms}}<option value="{{.ID}}" {{if eq .ID $.QueryFormID}}selected{{end}}>{{.Name}}</option>{{end}}</select>
  </div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_phone"}}</th><th>{{T "col_data"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Submissions}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td><code>{{.Data}}</code></td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "reminders"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "rem_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/reminders/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control"></div>
          <div class="form-group"><label>{{T "col_phone"}}</label><input name="phone" class="form-control" required></div>
          <div class="form-group"><label>{{T "col_amount"}}</label><input name="amount" type="number" step="0.01" class="form-control" required></div>
          <div class="form-group"><label>{{T "col_due"}}</label><input type="date" name="due_date" class="form-control" required></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="2">{{T "rem_default_msg"}}</textarea></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "rem_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_phone"}}</th><th>{{T "col_name"}}</th><th>{{T "col_amount"}}</th><th>{{T "col_due"}}</th><th>Status</th></tr></thead><tbody>
          {{range .Reminders}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Name}}</td><td>{{.Amount}}</td><td>{{.DueDate}}</td><td>{{.Status}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "analytics"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-chart-pie me-1"></i> {{T "analytics_agent_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>{{T "col_agent"}}</th><th>{{T "col_chats"}}</th><th>{{T "col_replies"}}</th><th>{{T "col_avg_resp"}}</th></tr></thead><tbody>
      {{range .AgentMetrics}}<tr><td>{{.AgentName}}</td><td>{{.Chats}}</td><td>{{.Replied}}</td><td>{{printf "%.0f" .AvgTime}}</td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center py-4">{{T "analytics_no_data"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
  <div class="card mt-3"><div class="card-header"><h4 class="card-header-title"><i class="la la-star me-1"></i> {{T "analytics_csat_title"}}</h4></div>
    <div class="card-body text-center">
      <div class="display-3 fw-bold text-warning">{{printf "%.1f" .CSATAvg}} ⭐</div>
      <p class="text-muted">{{T "analytics_csat_from" .CSATCount}}</p>
    </div>
  </div>
{{end}}

{{if eq .Page "depts"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "dept_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/depts/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="{{T "dept_name_ph"}}" required></div>
          <div class="form-group"><label>{{T "dept_agents"}}</label><select name="agents" class="form-control" multiple>{{range .Users}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div>
          <button class="btn btn-primary"><i class="la la-plus"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "dept_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "dept_agents"}}</th><th></th></tr></thead><tbody>
          {{range .Depts}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Agents}}</td><td><form method="post" action="/depts/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "recurring"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "rec_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/recurring/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
          <div class="form-group"><label>{{T "nav_contacts_groups"}}</label><select name="groups" class="form-control" multiple>{{range .Groups}}<option value="{{.ID}}">{{.Name}} ({{.Count}})</option>{{end}}</select></div>
          <div class="form-group"><label>{{T "col_message"}}</label><textarea name="message" class="form-control" rows="3" required></textarea></div>
          <div class="row">
            <div class="col-6"><label>{{T "rec_day"}}</label><select name="day_of_week" class="form-control"><option value="0">{{T "rec_daily"}}</option><option value="1">{{T "rec_mon"}}</option><option value="2">{{T "rec_tue"}}</option><option value="3">{{T "rec_wed"}}</option><option value="4">{{T "rec_thu"}}</option><option value="5">{{T "rec_fri"}}</option><option value="6">{{T "rec_sat"}}</option><option value="7">{{T "rec_sun"}}</option></select></div>
            <div class="col-6"><label>{{T "rec_hour"}}</label><input type="number" name="hour" class="form-control" value="9" min="0" max="23"></div>
          </div>
          <button class="btn btn-primary mt-2"><i class="la la-plus"></i> {{T "rec_create"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "rec_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "nav_contacts_groups"}}</th><th>{{T "col_schedule"}}</th><th>{{T "col_status"}}</th><th></th></tr></thead><tbody>
          {{range .Recurrings}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Groups}}</td><td>{{if eq .DayOfWeek 0}}{{T "rec_daily"}}{{else}}Day {{.DayOfWeek}}{{end}} @ {{.Hour}}:00</td><td>{{.Status}}</td><td>
            <form method="post" action="/recurring/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{if eq .Status "active"}}{{T "rec_pause"}}{{else}}{{T "rec_activate"}}{{end}}</button></form>
            <form method="post" action="/recurring/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form>
          </td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "uploads"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-folder-open me-1"></i> {{T "uploads_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_file"}}</th><th>{{T "col_url"}}</th></tr></thead><tbody>
      {{range $i, $f := .Files}}<tr><td>{{add $i 1}}</td><td>{{$f}}</td><td><code>/public/uploads/{{$f}}</code></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center py-4">{{T "uploads_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "blacklist"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i>{{T "bl_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/blacklist/add">
          <div class="form-group"><label>{{T "col_phone"}}</label><input name="phone" class="form-control" placeholder="628xxx" required></div>
          <div class="form-group"><label>{{T "col_reason"}}</label><input name="reason" class="form-control" placeholder="{{T "bl_reason_ph"}}"></div>
          <button class="btn btn-danger"><i class="la la-ban me-1"></i>{{T "bl_block"}}</button>
        </form></div>
      </div>
      <div class="card mt-3"><div class="card-header"><h4 class="card-header-title">{{T "bl_validate"}}</h4></div>
        <div class="card-body"><form method="post" action="/validate">
          <textarea name="numbers" class="form-control" rows="6" placeholder="628xxx&#10;628xxx" required></textarea>
          <button class="btn btn-primary mt-2"><i class="la la-check me-1"></i>{{T "bl_validate_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "bl_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_phone"}}</th><th>{{T "col_reason"}}</th><th>{{T "col_time"}}</th><th></th></tr></thead><tbody>
          {{range .Blacklist}}<tr><td>{{.ID}}</td><td>{{.Phone}}</td><td>{{.Reason}}</td><td class="text-muted small">{{.Created}}</td><td><form method="post" action="/blacklist/remove"><input type="hidden" name="phone" value="{{.Phone}}"><button class="btn btn-sm btn-success">{{T "bl_unblock"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center py-4">{{T "bl_empty"}}</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "csat"}}
  <div class="row">
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-warning">{{printf "%.1f" .CSATAvg}}</h1><p>{{T "csat_avg_rating"}}</p></div></div>
    </div>
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-primary">{{.CSATCount}}</h1><p>{{T "csat_total_resp"}}</p></div></div>
    </div>
    <div class="col-12 col-md-4">
      <div class="card text-center"><div class="card-body"><h1 class="display-3 fw-bold text-success">{{printf "%.0f" (mult .CSATAvg 20)}}%</h1><p>{{T "csat_score"}}</p></div></div>
    </div>
  </div>
{{end}}

{{if eq .Page "customers"}}
  <div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">{{T "cust_title"}}</h4><input type="text" id="custSearch" class="form-control form-control-sm" placeholder="{{T "cust_search_ph"}}" style="width:250px" oninput="var q=this.value.toLowerCase();document.querySelectorAll('.cust-row').forEach(r=>r.style.display=r.textContent.toLowerCase().includes(q)?'':'none')"></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>{{T "col_phone"}}</th><th>{{T "col_name"}}</th><th>{{T "col_orders"}}</th><th>{{T "col_last_active"}}</th><th></th></tr></thead><tbody>
      {{range .Contacts}}<tr class="cust-row"><td>{{.Phone}}</td><td>{{.Name}}</td><td>-</td><td>-</td><td><a href="/inbox/chat?phone={{.Phone}}" class="btn btn-sm btn-primary">{{T "cust_chat_btn"}}</a></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "calendar"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-calendar me-1"></i>{{T "cal_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>{{T "col_date"}}</th><th>{{T "col_title"}}</th><th>{{T "col_type"}}</th></tr></thead><tbody>
      {{range .CalEvents}}<tr><td>{{.Date}}</td><td>{{.Title}}</td><td><span class="badge badge-soft-{{if eq .Type "Campaign"}}primary{{else if eq .Type "Recurring"}}success{{else}}warning{{end}}">{{if eq .Type "Campaign"}}{{T "cal_campaign"}}{{else if eq .Type "Recurring"}}{{T "cal_recurring"}}{{else}}{{.Type}}{{end}}</span></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center py-4">{{T "cal_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "backup"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-database me-1"></i>{{T "backup_title"}}</h4></div>
    <div class="card-body text-center py-5">
      <form method="post" action="/backup">
        <button class="btn btn-primary btn-lg"><i class="la la-download me-1"></i>{{T "backup_btn"}}</button>
      </form>
      <p class="text-muted mt-2">{{T "backup_hint"}}</p>
    </div>
  </div>
{{end}}

{{if eq .Page "macros"}}
  <div class="row">
    <div class="col-12 col-lg-4">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i>{{T "macro_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/macros/add">
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="{{T "macro_name_ph"}}" required></div>
          <div class="form-group"><label>{{T "macro_actions"}}</label><textarea name="actions" class="form-control" rows="4" placeholder="assign:1;tag:resolved;reply:Terima kasih!;close" required></textarea><small class="form-text text-muted">{{T "macro_hint"}}</small></div>
          <button class="btn btn-primary"><i class="la la-plus"></i>{{T "macro_create"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-8">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "macro_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "macro_actions"}}</th><th></th></tr></thead><tbody>
          {{range .Macros}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><code>{{.Actions}}</code></td><td><form method="post" action="/macros/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "merge"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-code-branch me-1"></i>{{T "merge_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>{{T "col_phone"}}</th><th>{{T "col_count"}}</th><th>{{T "col_names"}}</th><th>{{T "col_ids"}}</th><th></th></tr></thead><tbody>
      {{range .Duplicates}}<tr><td>{{index . "phone"}}</td><td>{{index . "cnt"}}</td><td>{{index . "names"}}</td><td>{{index . "ids"}}</td><td><form method="post" action="/merge/execute"><input type="hidden" name="keep_id" value="{{index (split (index . "ids") ",") 0}}">{{range $i, $id := split (index . "ids") ","}}{{if gt $i 0}}<input type="hidden" name="merge_ids" value="{{$id}}">{{end}}{{end}}<button class="btn btn-sm btn-warning">{{T "merge_btn"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center py-4">{{T "merge_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "audit"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-history me-1"></i>{{T "audit_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_user"}}</th><th>{{T "col_action"}}</th><th>{{T "col_detail"}}</th><th>{{T "col_ip"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .AuditLogs}}<tr><td>{{.ID}}</td><td>#{{.UserID}}</td><td>{{.Action}}</td><td>{{.Detail}}</td><td>{{.IP}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center py-4">{{T "audit_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "translatetool"}}
  <div class="card mb-4"><div class="card-header"><h4 class="card-header-title"><i class="la la-language me-1"></i> {{T "tr_title"}}</h4></div>
    <div class="card-body">
      <div class="alert alert-warning mb-3"><i class="la la-info-circle me-1"></i> <strong>{{T "tr_ai_required"}}</strong> &mdash; {{T "tr_ai_required_desc"}} <a href="/autoreply" class="alert-link">{{T "tr_ai_configure"}}</a></div>
      <div class="row">
        <div class="col-md-5"><div class="form-group"><label>{{T "tr_source"}}</label><textarea id="srcText" class="form-control" rows="5" placeholder="{{T "tr_source_ph"}}"></textarea></div></div>
        <div class="col-md-2 d-flex align-items-end pb-3"><button onclick="doTranslate()" id="trBtn" class="btn btn-primary w-100"><i class="la la-sync me-1"></i> {{T "tr_btn"}}</button></div>
        <div class="col-md-5"><div class="form-group"><label>{{T "tr_result"}} <select id="langTo" class="form-select form-select-sm" style="width:auto;display:inline"><option value="id">ID</option><option value="en">EN</option><option value="es">ES</option><option value="fr">FR</option><option value="de">DE</option><option value="zh">ZH</option><option value="ja">JA</option><option value="ko">KO</option><option value="ar">AR</option></select></label><textarea id="resText" class="form-control" rows="5" readonly></textarea><div id="trError" class="text-danger small mt-1 d-none"></div></div>
      </div>
    </div>
  </div>
  <script>
  function doTranslate(){
    var t=document.getElementById('srcText').value.trim();
    if(!t) return;
    var to=document.getElementById('langTo').value;
    var btn=document.getElementById('trBtn');
    var err=document.getElementById('trError');
    err.classList.add('d-none');
    btn.disabled=true;
    btn.innerHTML='<span class="spinner-border spinner-border-sm me-1"></span> {{T "tr_btn"}}';
    fetch('/translate',{method:'POST',headers:{'Content-Type':'application/x-www-form-urlencoded'},body:'text='+encodeURIComponent(t)+'&to='+to})
    .then(function(r){if(!r.ok)throw new Error('Server error');return r.text()})
    .then(function(r){
      if(r===t){err.textContent='{{T "tr_no_ai"}}';err.classList.remove('d-none')}
      document.getElementById('resText').value=r;
    })
    .catch(function(e){err.textContent=e.message;err.classList.remove('d-none')})
    .finally(function(){btn.disabled=false;btn.innerHTML='<i class="la la-sync me-1"></i> {{T "tr_btn"}}'});
  }
  </script>
{{end}}

{{if eq .Page "widgetinfo"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-code me-1"></i> {{T "widget_title"}}</h4></div>
    <div class="card-body">
      <p>{{T "widget_desc"}}</p>
      <pre class="bg-light p-3 rounded"><code>&lt;script src="{{.AppURL}}/widget.js"&gt;&lt;/script&gt;</code></pre>
      <p class="text-muted small">{{T "widget_pos_desc"}}</p>
      <p class="text-muted small">{{T "widget_webhook_desc"}} <code>{{.AppURL}}/email-webhook</code></p>
    </div>
  </div>
{{end}}

{{if eq .Page "emailwa"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-envelope me-1"></i> {{T "emailwa_title"}}</h4></div>
    <div class="card-body">
      <p>{{T "emailwa_desc"}}</p>
      <h5 class="mt-3">{{T "emailwa_setup"}}</h5>
      <pre class="bg-light p-3 rounded"><code>POST {{.AppURL}}/email-webhook
Content-Type: application/x-www-form-urlencoded

from=sender@email.com&subject=Judul Email&text=Isi email</code></pre>
      <p class="text-muted small">{{T "emailwa_integration"}}</p>
    </div>
  </div>
{{end}}

{{if eq .Page "meta_send"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-paper-plane me-1"></i> {{T "meta_send_title"}}</h4></div>
    <div class="card-body"><form method="post" action="/meta/send">
      <div class="row">
        <div class="col-md-4"><select name="account_id" class="form-control" required><option value="">{{T "meta_select_account"}}</option>{{range .MetaAccounts}}<option value="{{.ID}}">{{.Name}} ({{.PhoneNumberID}})</option>{{end}}</select></div>
        <div class="col-md-4"><input name="phone" class="form-control" placeholder="628123456789" required></div>
        <div class="col-md-4"><button class="btn btn-primary w-100"><i class="la la-paper-plane"></i> {{T "send_btn"}}</button></div>
      </div>
      <textarea name="message" class="form-control mt-2" rows="4" placeholder="{{T "meta_msg_ph"}}" required></textarea>
    </form></div>
  </div>
{{end}}

{{if eq .Page "meta_campaigns"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-bullhorn me-1"></i> {{T "meta_camp_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Groups</th><th>Progress</th><th>Status</th></tr></thead><tbody>
      {{range .Campaigns}}{{if .MetaAccountID}}<tr><td>{{.ID}}</td><td>{{.Name}}<span class="badge ms-1" style="background:#4F46E5;color:#fff;font-size:9px">META</span></td><td>{{.Groups}}</td><td>{{.Sent}}/{{.Total}}</td><td>{{.Status}}</td></tr>{{end}}{{else}}<tr><td colspan="5" class="text-muted text-center py-4">{{T "meta_camp_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "meta_inbox"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-comments me-1"></i> {{T "meta_inbox_title"}}</h4></div>
    <div class="card-body text-center py-5">
      <p class="text-muted">{{T "meta_inbox_desc"}}</p>
      <p>Webhook URL: <code>{{.AppURL}}/webhook/meta</code></p>
      <p class="small text-muted">{{T "meta_inbox_url_hint"}}</p>
      <a href="/inbox" class="btn btn-primary">{{T "meta_inbox_btn"}}</a>
    </div>
  </div>
{{end}}

{{if eq .Page "meta_logs"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-clipboard-list me-1"></i> {{T "meta_logs_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_type"}}</th><th>{{T "col_reason"}}</th><th>{{T "col_content"}}</th></tr></thead><tbody>
      {{range .Logs}}{{if eq .Type "meta"}}<tr><td>{{.ID}}</td><td>{{.Type}}</td><td>{{.Reason}}</td><td>{{.Content}}</td></tr>{{end}}{{else}}<tr><td colspan="4" class="text-muted text-center">{{T "meta_logs_empty"}}</td></tr>{{end}}
    </tbody></table></div>
  </div>
{{end}}

{{if eq .Page "meta_analytics"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-chart-bar me-1"></i> {{T "meta_stats_title"}}</h4></div>
    <div class="card-body">
      <div class="row text-center">
        <div class="col-4"><div class="display-4 fw-bold text-primary">{{.CountSent}}</div><small class="text-muted">{{T "meta_stats_sent"}}</small></div>
        <div class="col-4"><div class="display-4 fw-bold text-success">{{.CountReceived}}</div><small class="text-muted">{{T "meta_stats_recv"}}</small></div>
        <div class="col-4"><div class="display-4 fw-bold text-info">{{len .MetaAccounts}}</div><small class="text-muted">{{T "meta_stats_accts"}}</small></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "meta_webhook"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-link me-1"></i> {{T "meta_webhook_title"}}</h4></div>
    <div class="card-body">
      <p><strong>Webhook URL:</strong> <code>{{.AppURL}}/webhook/meta</code></p>
      <p class="text-muted">{{T "meta_webhook_desc"}}</p>
      <ol class="small">
        <li>{{T "meta_webhook_step1"}} <a href="https://developers.facebook.com" target="_blank">developers.facebook.com</a></li>
        <li>{{T "meta_webhook_step2"}}</li>
        <li>{{T "meta_webhook_step3"}}</li>
        <li>{{T "meta_webhook_step4"}}</li>
        <li>Subscribe ke field <code>messages</code></li>
      </ol>
      <hr>
      <h5>{{T "meta_accts_registered"}}</h5>
      <table class="table table-sm"><thead><tr><th>{{T "col_name"}}</th><th>{{T "meta_phone_col"}}</th><th>{{T "meta_verify_token"}}</th></tr></thead><tbody>
        {{range .MetaAccounts}}<tr><td>{{.Name}}</td><td><code>{{.PhoneNumberID}}</code></td><td><code>{{.VerifyToken}}</code></td></tr>{{else}}<tr><td colspan="3" class="text-muted">{{T "meta_accts_empty"}}</td></tr>{{end}}
      </tbody></table>
    </div>
  </div>
{{end}}

{{if eq .Page "upgrade"}}
<div class="row justify-content-center">
  <div class="col-12 text-center mb-4">
    <h2 class="fw-bold">Pilih Paket ChatGo</h2>
    <p class="text-muted">Mulai dari gratis sampai enterprise. Upgrade kapan saja.</p>
  </div>

  <!-- FREE -->
  <div class="col-12 col-md-6 col-lg-4 mb-4">
    <div class="card border shadow-sm h-100" style="border-radius:16px">
      <div class="card-body text-center p-4">
        <span class="badge bg-success mb-2">Open Source</span>
        <h3 class="fw-bold">Free</h3>
        <div class="display-4 fw-bold text-success my-3">Rp 0</div>
        <p class="text-muted small">Selamanya gratis. Full source code.</p>
        <hr>
        <ul class="text-start small mb-3" style="list-style:none;padding:0;line-height:2.2">
        <li>✅ Multi-Account WA</li>
        <li>✅ Auto Reply (keyword + spintax)</li>
        <li>✅ AI Auto Reply (BYOK)</li>
        <li>✅ Broadcast + Drip + Recurring</li>
        <li>✅ Inbox real-time + Agent</li>
        <li>✅ Multi-User SaaS + RBAC</li>
        <li>✅ Payment Gateway</li>
        <li>✅ REST API + Webhook</li>
        <li>✅ E-commerce + White-label</li>
        <li>✅ PSEO 1.3M+ halaman</li>
        <li>✅ Single binary Go</li>
        </ul>
        <div class="text-muted small mb-3">
          <strong>⚠️ Build sendiri</strong><br>
          git clone → go build → jalan.<br>
          Support: komunitas GitHub.
        </div>
        <a href="https://github.com/linducip2208/chatforge" target="_blank" class="btn btn-outline-success w-100 mb-2">
          <i class="la la-github me-1"></i> GitHub
        </a>
        <a href="/register" class="btn btn-outline-primary w-100">
          <i class="la la-user-plus me-1"></i> Coba Demo
        </a>
      </div>
    </div>
  </div>

  <!-- STANDARD -->
  <div class="col-12 col-md-6 col-lg-4 mb-4">
    <div class="card border-primary shadow h-100" style="border-radius:16px;border-width:2px">
      <div class="card-body text-center p-4">
        <span class="badge bg-primary mb-2">Populer</span>
        <h3 class="fw-bold">Standard</h3>
        <div class="display-4 fw-bold text-primary my-3">$66</div>
        <p class="text-muted small">One-time purchase. Lifetime.</p>
        <hr>
        <ul class="text-start small mb-3" style="list-style:none;padding:0;line-height:2.2">
        <li>✅ Semua fitur Free</li>
        <li>📦 Binary siap pakai (.exe)</li>
        <li>🛠️ Bantuan instalasi + setup</li>
        <li>📞 Support WA 6 bulan</li>
        <li>🔄 Update binary gratis</li>
        <li>📖 Dokumentasi lengkap</li>
        <li>🎯 Jaminan berfungsi</li>
        </ul>
        <div class="text-muted small mb-3">
          <strong>Cocok untuk:</strong> UMKM, bisnis kecil-menengah, startup.
        </div>
        <a href="https://wa.me/6281296052010?text=Halo%20saya%20mau%20beli%20ChatGo%20Standard%20($66)" target="_blank" class="btn btn-primary w-100">
          <i class="la la-whatsapp me-1"></i> Beli via WhatsApp
        </a>
      </div>
    </div>
  </div>

  <!-- PRO -->
  <div class="col-12 col-md-6 col-lg-4 mb-4">
    <div class="card border-dark shadow-sm h-100" style="border-radius:16px;background:linear-gradient(135deg,#0f1f33,#152e4d)">
      <div class="card-body text-center p-4" style="color:#fff">
        <span class="badge bg-warning text-dark mb-2">Coming Soon</span>
        <h3 class="fw-bold">Pro</h3>
        <div class="display-4 fw-bold text-warning my-3">$799</div>
        <p class="small" style="opacity:.8">One-time. Early adopter price.</p>
        <hr style="border-color:rgba(255,255,255,.15)">
        <ul class="text-start small mb-3" style="list-style:none;padding:0;line-height:2.2">
        <li>✅ Semua fitur Standard</li>
        <li>📸 Instagram DM Automation</li>
        <li>📘 Facebook Messenger</li>
        <li>✈️ Telegram Bot</li>
        <li>🌐 Omnichannel Inbox</li>
        <li>🧩 Visual Flow Builder</li>
        <li>🔗 n8n / Zapier / Make</li>
        <li>📊 Google Sheets Sync</li>
        <li>🏢 Agency Dashboard</li>
        <li>⭐ Priority Support</li>
        </ul>
        <div class="small mb-3" style="opacity:.7">
          <strong>Target:</strong> Agency, SaaS Provider, Enterprise.
        </div>
        <a href="https://wa.me/6281296052010?text=Halo%20saya%20tertarik%20ChatGo%20Pro%20($799)" target="_blank" class="btn btn-warning w-100">
          <i class="la la-whatsapp me-1"></i> Pre-order via WA
        </a>
      </div>
    </div>
  </div>

  <!-- FAQ -->
  <div class="col-12 col-lg-8 mt-5">
    <h4 class="fw-bold mb-3">Pertanyaan Umum</h4>
    <div class="accordion" id="faq">
      <div class="card mb-2"><div class="card-header p-3" style="cursor:pointer" data-toggle="collapse" data-target="#q1"><strong>Apa bedanya Free dan Standard?</strong></div>
      <div id="q1" class="collapse"><div class="card-body">Fitur <strong>sama persis</strong>. Bedanya: Free Anda clone dari GitHub dan build sendiri. Standard kami yang build, kirim binary siap pakai, dan beri support 6 bulan. Anda bayar untuk kenyamanan — tidak perlu install Go, setup environment, atau debugging masalah build.</div></div></div>
      <div class="card mb-2"><div class="card-header p-3" style="cursor:pointer" data-toggle="collapse" data-target="#q2"><strong>Apakah Standard dapat source code?</strong></div>
      <div id="q2" class="collapse"><div class="card-body">Tidak. Standard hanya binary (.exe). Kalau mau source code, gunakan versi Free dari GitHub — gratis, full source.</div></div></div>
      <div class="card mb-2"><div class="card-header p-3" style="cursor:pointer" data-toggle="collapse" data-target="#q3"><strong>Kapan Pro tersedia?</strong></div>
      <div id="q3" class="collapse"><div class="card-body">Pro sedang dalam pengembangan. Fitur omnichannel (Instagram, Facebook, Telegram) + Visual Flow Builder. Pre-order sekarang untuk harga early adopter.</div></div></div>
      <div class="card mb-2"><div class="card-header p-3" style="cursor:pointer" data-toggle="collapse" data-target="#q4"><strong>Apakah bisa upgrade nanti?</strong></div>
      <div id="q4" class="collapse"><div class="card-body">Bisa. Free → Standard ($66): beli Standard, dapat binary. Standard → Pro ($799): beli Pro, dapat binary Pro (update di atas Standard). Tidak perlu instalasi ulang — cukup replace .exe.</div></div></div>
      <div class="card mb-2"><div class="card-header p-3" style="cursor:pointer" data-toggle="collapse" data-target="#q5"><strong>Kenapa Standard murah tapi Pro mahal?</strong></div>
      <div id="q5" class="collapse"><div class="card-body">Standard adalah binary pre-built dari source code yang sama dengan Free — Anda bayar convenience. Pro adalah produk berbeda: omnichannel (Instagram, Facebook, Telegram), Visual Flow Builder, agency dashboard — fitur yang tidak ada di Free/Standard dan dikembangkan terpisah.</div></div></div>
    </div>
  </div>
</div>
{{end}}

{{if eq .Page "subscribe"}}
  <div class="row">
    <div class="col-12"><h2 class="mb-4">{{T "sub_choose"}}</h2></div>
    {{range .Packages}}
    {{$pkgID := .ID}}
    <div class="col-12 col-md-6 col-lg-4 mb-4">
      <div class="card border-0 shadow-sm h-100" style="border-radius:14px">
        <div class="card-body text-center p-4">
          <h4 class="fw-bold">{{.Name}}</h4>
          <div class="display-4 fw-bold text-primary my-3">{{.Price}}</div>
          <p class="text-muted small">
            <strong>{{T "sub_limits_label"}}</strong> Send:{{.SendLimit}} WA:{{.WaAccountLimit}} Dev:{{.DeviceLimit}} Contact:{{.ContactLimit}}<br>
            <strong>{{T "sub_features_label"}}</strong> {{.Services}}
          </p>
          <div class="my-3">
            <form method="post" action="/subscribe/checkout" class="mb-2">
              <input type="hidden" name="package_id" value="{{$pkgID}}">
              <input type="hidden" name="gateway_id" value="free">
              <div class="input-group"><input type="text" name="voucher" class="form-control form-control-sm" placeholder="Kode Voucher"><button class="btn btn-sm btn-outline-success">Aktifkan</button></div>
            </form>
          </div>
          {{range $.PaymentGateways}}
          {{if eq .Status "active"}}
          <form method="post" action="/subscribe/checkout" class="mt-2">
            <input type="hidden" name="package_id" value="{{$pkgID}}">
            <input type="hidden" name="gateway_id" value="{{.ID}}">
            <button class="btn btn-primary w-100"><i class="la la-credit-card me-1"></i> {{T "sub_pay_via"}} {{.Name}}</button>
          </form>
          {{end}}
          {{end}}
          {{if not $.PaymentGateways}}<p class="text-muted small mt-2">{{T "sub_no_gateway"}}</p>{{end}}
        </div>
      </div>
    </div>
    {{else}}
    <div class="col-12 text-center py-5 text-muted">{{T "sub_no_packages"}}</div>
    {{end}}
  </div>
{{end}}

{{if eq .Page "admin_paygateways"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-plus me-1"></i> {{T "paygw_add"}}</h4></div>
        <div class="card-body"><form method="post" action="/admin/gateways-pay/add">
          <div class="form-group"><label>{{T "col_provider"}}</label><select name="provider" class="form-control"><option value="midtrans">Midtrans (ID)</option><option value="xendit">Xendit (ID)</option><option value="paypal">PayPal (Intl)</option><option value="stripe">Stripe (Intl)</option></select></div>
          <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="My Gateway"></div>
          <div class="form-group"><label>{{T "paygw_apikey"}}</label><input name="api_key" class="form-control"></div>
          <div class="form-group"><label>{{T "paygw_apisecret"}}</label><input name="api_secret" class="form-control"></div>
          <div class="form-group"><label>{{T "paygw_webhooksecret"}}</label><input name="webhook_secret" class="form-control"></div>
          <div class="form-group"><label>{{T "paygw_currency"}}</label><select name="currency" class="form-control"><option value="IDR">IDR</option><option value="USD">USD</option><option value="EUR">EUR</option><option value="SGD">SGD</option></select></div>
          <div class="form-group"><label>{{T "paygw_baseurl"}}</label><input name="base_url" class="form-control" placeholder="Kosongkan untuk default"></div>
          <button class="btn btn-primary"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form></div>
      </div>
    </div>
    <div class="col-12 col-lg-7">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "paygw_list"}}</h4></div>
        <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>Name</th><th>Provider</th><th>Currency</th><th>Status</th><th></th></tr></thead><tbody>
          {{range .PaymentGateways}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Provider}}</td><td>{{.Currency}}</td><td>{{if eq .Status "active"}}<span class="badge badge-soft-success">{{T "paygw_active"}}</span>{{else}}<span class="badge badge-soft-secondary">{{T "paygw_inactive"}}</span>{{end}}</td><td>
            <form method="post" action="/admin/gateways-pay/toggle" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-white">{{if eq .Status "active"}}{{T "paygw_disable"}}{{else}}{{T "paygw_enable"}}{{end}}</button></form>
            <form method="post" action="/admin/gateways-pay/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form>
          </td></tr>{{else}}<tr><td colspan="6" class="text-muted text-center">-</td></tr>{{end}}
        </tbody></table></div>
      </div>
    </div>
  </div>
{{end}}

{{if eq .Page "admin_transactions_pay"}}
  <div class="card"><div class="card-header"><h4 class="card-header-title"><i class="la la-receipt me-1"></i> {{T "paytx_title"}}</h4></div>
    <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_invoice"}}</th><th>{{T "col_user"}}</th><th>{{T "col_package"}}</th><th>{{T "col_amount"}}</th><th>{{T "col_status"}}</th><th>{{T "col_time"}}</th></tr></thead><tbody>
      {{range .Txs}}<tr><td>{{.ID}}</td><td><code>{{.InvoiceID}}</code></td><td>#{{.UserID}}</td><td>#{{.PackageID}}</td><td>{{.Amount}} {{.Currency}}</td><td>{{if eq .Status "paid"}}<span class="badge badge-soft-success">{{T "paytx_paid"}}</span>{{else if eq .Status "failed"}}<span class="badge badge-soft-danger">{{T "paytx_failed"}}</span>{{else}}<span class="badge badge-soft-warning">{{.Status}}</span>{{end}}</td><td class="text-muted small">{{.Created}}</td></tr>{{else}}<tr><td colspan="7" class="text-muted text-center py-4">{{T "paytx_empty"}}</td></tr>{{end}}
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
      <form method="post" action="/ussd/add"><div class="form-group"><label>{{T "ussd_code"}}</label><input name="code" class="form-control" placeholder="*123#" required></div><button class="btn btn-primary lift"><i class="la la-satellite-dish me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">{{T "btn_cancel"}}</a>{{end}}</form></div></div></div>
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
        <div class="form-group"><label>{{T "aik_api_key"}}</label><input name="apikey" class="form-control" required></div>
        <div class="form-group"><label>Prompt</label><textarea name="system_prompt" class="form-control" rows="3"></textarea></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "nav_ai_keys"}}</h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_provider"}}</th><th>{{T "ar_ai_model"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .AiKeys}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><span class="badge badge-soft-secondary">{{.Provider}}</span></td><td>{{.Model}}</td><td><form method="post" action="/ai/keys/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "ai_plugins"}}
  <div class="row">
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "aip_add"}}</h4></div><div class="card-body">
      <form method="post" action="/ai/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><div class="form-group"><label>Endpoint</label><input name="endpoint" class="form-control" placeholder="https://..."></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">{{T "btn_cancel"}}</a>{{end}}</form></div></div></div>
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
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "admin_stat_users"}}</h6><span class="h2 mb-0">{{.TotalUsers}}</span></div>
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
<div class="col"><h6 class="text-uppercase text-muted mb-2 small">{{T "admin_stat_campaigns"}}</h6><span class="h2 mb-0">{{.RunningCampaigns}}</span></div>
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
<div class="col-12"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "admin_overview_title"}}</h4></div>
<div class="card-body"><canvas id="adminChart" height="80"></canvas></div></div></div>
</div>
<div class="row mt-3">
<div class="col-6 col-xl-3"><a href="/admin/users" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-users la-2x text-primary mb-2 d-block"></i><strong>{{T "admin_card_users"}}</strong><br><small class="text-muted">{{T "admin_card_users_desc"}}</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/packages" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-box la-2x text-success mb-2 d-block"></i><strong>{{T "admin_card_packages"}}</strong><br><small class="text-muted">{{T "admin_card_packages_desc"}}</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/waservers" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-server la-2x text-warning mb-2 d-block"></i><strong>{{T "admin_card_servers"}}</strong><br><small class="text-muted">{{T "admin_card_servers_desc"}}</small></div></a></div>
<div class="col-6 col-xl-3"><a href="/admin/subscriptions" class="card text-decoration-none"><div class="card-body text-center py-4"><i class="la la-star la-2x text-danger mb-2 d-block"></i><strong>{{T "admin_card_subs"}}</strong><br><small class="text-muted">{{T "admin_card_subs_desc"}}</small></div></a></div>
</div>
<script>
new Chart(document.getElementById('adminChart'),{type:'line',data:{labels:[{{.ChartLabels}}],datasets:[{label:'{{T "chart_sent"}}',data:[{{.ChartSent}}],borderColor:'#4F46E5',backgroundColor:'rgba(79,70,229,.1)',fill:true,tension:.3,pointRadius:2},{label:'{{T "chart_received"}}',data:[{{.ChartReceived}}],borderColor:'#10B981',backgroundColor:'rgba(16,185,129,.1)',fill:true,tension:.3,pointRadius:2}]},options:{responsive:true,plugins:{legend:{position:'bottom'}},scales:{y:{beginAtZero:true}}}})
</script>
{{end}}

{{if eq .Page "admin_users"}}
<div class="row">
<div class="col-12"><div class="card"><div class="card-header d-flex justify-content-between"><h4 class="card-header-title">{{T "adm_users"}}</h4><button class="btn btn-primary btn-sm lift" onclick="document.getElementById('addUserForm').style.display=document.getElementById('addUserForm').style.display==='none'?'block':'none'"><i class="la la-plus me-1"></i> {{T "usr_add"}}</button></div>
<div id="addUserForm" style="display:none;border-bottom:1px solid #eee;padding:16px"><form method="post" action="/admin/users/add">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "auth_email"}}</label><input name="email" type="email" class="form-control" required></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "auth_password"}}</label><input name="password" type="password" class="form-control"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "usr_role"}}</label><select name="role" class="form-control">{{range .Roles}}<option value="{{.Name}}">{{.Name}}</option>{{end}}</select></div></div>
<div class="col-12"><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></div>
</div></form></div>
<div class="table-responsive"><table class="table table-sm card-table mb-0"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "auth_email"}}</th><th>{{T "usr_role"}}</th><th>{{T "col_registered"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
{{range .Users}}<tr>
<td>{{.ID}}</td>
<td>{{.Name}}</td>
<td>{{.Email}}</td>
<td><span class="badge {{if eq .Role "admin"}}badge-soft-primary{{else}}badge-soft-secondary{{end}}">{{.Role}}</span></td>
<td class="text-muted small">{{.Created}}</td>
<td class="text-nowrap">
 <a class="btn btn-sm btn-white" href="/admin/users?edit={{.ID}}"><i class="la la-edit"></i></a>
 <a class="btn btn-sm btn-warning" href="/admin/users/impersonate?id={{.ID}}" title="{{T "usr_impersonate"}}"><i class="la la-user-circle"></i></a>
 <form method="post" action="/admin/users/reset-password" style="display:inline" onsubmit="return confirm('Reset password untuk {{.Name}}?')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-info" title="Reset Password"><i class="la la-key"></i></button></form>
 <form method="post" action="/admin/users/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form>
</td>
</tr>{{else}}<tr><td colspan="7" class="text-muted text-center py-4">-</td></tr>{{end}}
</tbody></table></div></div></div>
{{if .EditID}}
<div class="col-12 mt-3"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> {{T "usr_edit_title"}} #{{.EditID}}</h4></div>
<div class="card-body"><form method="post" action="/admin/users/edit">
<input type="hidden" name="id" value="{{.EditID}}">
<div class="row">
<div class="col-md-6"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "auth_email"}}</label><input name="email" type="email" class="form-control" value="{{.EditPhone}}"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "usr_password_keep"}}</label><input name="password" type="password" class="form-control" placeholder="••••••"></div></div>
<div class="col-md-6"><div class="form-group"><label>{{T "usr_role"}}</label><select name="role" class="form-control">{{range .Roles}}<option value="{{.Name}}" {{if eq .Name $.EditRole}}selected{{end}}>{{.Name}}</option>{{end}}</select></div></div>
</div>
<button class="btn btn-warning lift"><i class="la la-save me-1"></i> {{T "usr_update"}}</button> <a href="/admin/users" class="btn btn-white ms-2">{{T "btn_cancel"}}</a>
</form></div></div></div>
{{end}}
</div>
{{end}}

{{if eq .Page "admin_roles"}}
  <div class="row">
    {{if .EditID}}
    <div class="col-12 col-lg-4"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h4 class="card-header-title"><i class="la la-edit me-1"></i> {{T "role_edit_title"}} #{{.EditID}}</h4></div><div class="card-body"><form method="post" action="/admin/roles/edit"><input type="hidden" name="id" value="{{.EditID}}"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div><div class="form-group"><label>{{T "role_perms"}}</label><select name="permissions" class="form-control" multiple size="18" style="overflow-y:auto;min-height:360px"><option value="manage_users">{{T "perm_manage_users"}}</option><option value="manage_roles">{{T "perm_manage_roles"}}</option><option value="manage_packages">{{T "perm_manage_packages"}}</option><option value="manage_vouchers">{{T "perm_manage_vouchers"}}</option><option value="manage_subscriptions">{{T "perm_manage_subscriptions"}}</option><option value="manage_transactions">{{T "perm_manage_transactions"}}</option><option value="manage_payouts">{{T "perm_manage_payouts"}}</option><option value="manage_pages">{{T "perm_manage_pages"}}</option><option value="manage_marketing">{{T "perm_manage_marketing"}}</option><option value="manage_languages">{{T "perm_manage_languages"}}</option><option value="manage_waservers">{{T "perm_manage_waservers"}}</option><option value="manage_gateways">{{T "perm_manage_gateways"}}</option><option value="manage_shorteners">{{T "perm_manage_shorteners"}}</option><option value="manage_plugins">{{T "perm_manage_plugins"}}</option><option value="manage_meta">{{T "perm_manage_meta"}}</option><option value="manage_metatemplates">{{T "perm_manage_metatemplates"}}</option><option value="wa_send">{{T "perm_wa_send"}}</option><option value="wa_broadcast">{{T "perm_wa_broadcast"}}</option><option value="wa_scheduled">{{T "perm_wa_scheduled"}}</option><option value="wa_sent">{{T "perm_wa_sent"}}</option><option value="wa_received">{{T "perm_wa_received"}}</option><option value="wa_inbox">{{T "perm_wa_inbox"}}</option><option value="wa_status">{{T "perm_wa_status"}}</option><option value="wa_autoreply">{{T "perm_wa_autoreply"}}</option><option value="wa_ai_keys">{{T "perm_wa_ai_keys"}}</option><option value="wa_ai_plugins">{{T "perm_wa_ai_plugins"}}</option><option value="wa_knowledge">{{T "perm_wa_knowledge"}}</option><option value="wa_contacts">{{T "perm_wa_contacts"}}</option><option value="wa_groups">{{T "perm_wa_groups"}}</option><option value="wa_unsub">{{T "perm_wa_unsub"}}</option><option value="wa_templates">{{T "perm_wa_templates"}}</option><option value="wa_apikeys">{{T "perm_wa_apikeys"}}</option><option value="wa_webhooks">{{T "perm_wa_webhooks"}}</option><option value="wa_logger">{{T "perm_wa_logger"}}</option><option value="wa_settings">{{T "perm_wa_settings"}}</option><option value="wa_docs">{{T "perm_wa_docs"}}</option><option value="wa_hosts">{{T "perm_wa_hosts"}}</option><option value="wa_ussd">{{T "perm_wa_ussd"}}</option><option value="wa_impersonate">{{T "perm_wa_impersonate"}}</option><option value="wa_drips">{{T "perm_wa_drips"}}</option><option value="wa_tags">{{T "perm_wa_tags"}}</option><option value="wa_canned">{{T "perm_wa_canned"}}</option><option value="wa_recurring">{{T "perm_wa_recurring"}}</option><option value="wa_store">{{T "perm_wa_store"}}</option><option value="wa_orders">{{T "perm_wa_orders"}}</option><option value="wa_forms">{{T "perm_wa_forms"}}</option><option value="wa_reminders">{{T "perm_wa_reminders"}}</option><option value="wa_analytics">{{T "perm_wa_analytics"}}</option><option value="wa_blacklist">{{T "perm_wa_blacklist"}}</option><option value="wa_csat">{{T "perm_wa_csat"}}</option><option value="wa_depts">{{T "perm_wa_depts"}}</option><option value="wa_customers">{{T "perm_wa_customers"}}</option><option value="wa_calendar">{{T "perm_wa_calendar"}}</option><option value="wa_macros">{{T "perm_wa_macros"}}</option><option value="wa_files">{{T "perm_wa_files"}}</option><option value="wa_merge">{{T "perm_wa_merge"}}</option><option value="wa_translate">{{T "perm_wa_translate"}}</option><option value="wa_audit">{{T "perm_wa_audit"}}</option><option value="wa_backup">{{T "perm_wa_backup"}}</option><option value="wa_subscribe">{{T "perm_wa_subscribe"}}</option><option value="manage_paygateways">{{T "perm_manage_paygateways"}}</option><option value="manage_paytx">{{T "perm_manage_paytx"}}</option></select></div><button class="btn btn-warning lift"><i class="la la-save me-1"></i> Update</button> <a href="/admin/roles" class="btn btn-white ms-2">{{T "ar_cancel"}}</a></form><script>document.addEventListener('DOMContentLoaded',function(){var s=document.querySelector('form[action=\"/admin/roles/edit\"] select[name=\"permissions\"]');if(s){var v='{{.EditContent}}';v.split(',').forEach(function(p){var o=s.querySelector('option[value=\"'+p.replace(/^\\s+|\\s+$/g,'')+'\"]');if(o)o.selected=true})}})</script></div></div></div>
    {{else}}
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "role_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/roles/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div>
      <div class="form-group"><label>{{T "role_perms"}}</label><select name="permissions" class="form-control" multiple size="18" style="overflow-y:auto;min-height:360px"><option value="manage_users">{{T "perm_manage_users"}}</option><option value="manage_roles">{{T "perm_manage_roles"}}</option><option value="manage_packages">{{T "perm_manage_packages"}}</option><option value="manage_vouchers">{{T "perm_manage_vouchers"}}</option><option value="manage_subscriptions">{{T "perm_manage_subscriptions"}}</option><option value="manage_transactions">{{T "perm_manage_transactions"}}</option><option value="manage_payouts">{{T "perm_manage_payouts"}}</option><option value="manage_pages">{{T "perm_manage_pages"}}</option><option value="manage_marketing">{{T "perm_manage_marketing"}}</option><option value="manage_languages">{{T "perm_manage_languages"}}</option><option value="manage_waservers">{{T "perm_manage_waservers"}}</option><option value="manage_gateways">{{T "perm_manage_gateways"}}</option><option value="manage_shorteners">{{T "perm_manage_shorteners"}}</option><option value="manage_plugins">{{T "perm_manage_plugins"}}</option><option value="manage_meta">{{T "perm_manage_meta"}}</option><option value="manage_metatemplates">{{T "perm_manage_metatemplates"}}</option><option value="wa_send">{{T "perm_wa_send"}}</option><option value="wa_broadcast">{{T "perm_wa_broadcast"}}</option><option value="wa_scheduled">{{T "perm_wa_scheduled"}}</option><option value="wa_sent">{{T "perm_wa_sent"}}</option><option value="wa_received">{{T "perm_wa_received"}}</option><option value="wa_inbox">{{T "perm_wa_inbox"}}</option><option value="wa_status">{{T "perm_wa_status"}}</option><option value="wa_autoreply">{{T "perm_wa_autoreply"}}</option><option value="wa_ai_keys">{{T "perm_wa_ai_keys"}}</option><option value="wa_ai_plugins">{{T "perm_wa_ai_plugins"}}</option><option value="wa_knowledge">{{T "perm_wa_knowledge"}}</option><option value="wa_contacts">{{T "perm_wa_contacts"}}</option><option value="wa_groups">{{T "perm_wa_groups"}}</option><option value="wa_unsub">{{T "perm_wa_unsub"}}</option><option value="wa_templates">{{T "perm_wa_templates"}}</option><option value="wa_apikeys">{{T "perm_wa_apikeys"}}</option><option value="wa_webhooks">{{T "perm_wa_webhooks"}}</option><option value="wa_logger">{{T "perm_wa_logger"}}</option><option value="wa_settings">{{T "perm_wa_settings"}}</option><option value="wa_docs">{{T "perm_wa_docs"}}</option><option value="wa_hosts">{{T "perm_wa_hosts"}}</option><option value="wa_ussd">{{T "perm_wa_ussd"}}</option><option value="wa_impersonate">{{T "perm_wa_impersonate"}}</option><option value="wa_drips">{{T "perm_wa_drips"}}</option><option value="wa_tags">{{T "perm_wa_tags"}}</option><option value="wa_canned">{{T "perm_wa_canned"}}</option><option value="wa_recurring">{{T "perm_wa_recurring"}}</option><option value="wa_store">{{T "perm_wa_store"}}</option><option value="wa_orders">{{T "perm_wa_orders"}}</option><option value="wa_forms">{{T "perm_wa_forms"}}</option><option value="wa_reminders">{{T "perm_wa_reminders"}}</option><option value="wa_analytics">{{T "perm_wa_analytics"}}</option><option value="wa_blacklist">{{T "perm_wa_blacklist"}}</option><option value="wa_csat">{{T "perm_wa_csat"}}</option><option value="wa_depts">{{T "perm_wa_depts"}}</option><option value="wa_customers">{{T "perm_wa_customers"}}</option><option value="wa_calendar">{{T "perm_wa_calendar"}}</option><option value="wa_macros">{{T "perm_wa_macros"}}</option><option value="wa_files">{{T "perm_wa_files"}}</option><option value="wa_merge">{{T "perm_wa_merge"}}</option><option value="wa_translate">{{T "perm_wa_translate"}}</option><option value="wa_audit">{{T "perm_wa_audit"}}</option><option value="wa_backup">{{T "perm_wa_backup"}}</option><option value="wa_subscribe">{{T "perm_wa_subscribe"}}</option><option value="manage_paygateways">{{T "perm_manage_paygateways"}}</option><option value="manage_paytx">{{T "perm_manage_paytx"}}</option></select></div>
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
        <div class="form-group col-3"><label>Durasi (hari)</label><input name="duration" type="number" class="form-control" value="30"></div></div>
        <div class="form-group"><label>{{T "pkg_services"}}</label><select name="services" class="form-control" multiple size="15"><option value="whatsapp">{{T "svc_whatsapp"}}</option><option value="meta">{{T "svc_meta"}}</option><option value="broadcast">{{T "svc_broadcast"}}</option><option value="drips">{{T "svc_drips"}}</option><option value="recurring">{{T "svc_recurring"}}</option><option value="ai">{{T "svc_ai"}}</option><option value="inbox">{{T "svc_inbox"}}</option><option value="contacts">{{T "svc_contacts"}}</option><option value="tags">{{T "svc_tags"}}</option><option value="groups">{{T "svc_groups"}}</option><option value="merge">{{T "svc_merge"}}</option><option value="store">{{T "svc_store"}}</option><option value="payment">{{T "svc_payment"}}</option><option value="forms">{{T "svc_forms"}}</option><option value="reminders">{{T "svc_reminders"}}</option><option value="analytics">{{T "svc_analytics"}}</option><option value="csat">{{T "svc_csat"}}</option><option value="api">{{T "svc_api"}}</option><option value="webhooks">{{T "svc_webhooks"}}</option><option value="templates">{{T "svc_templates"}}</option><option value="canned">{{T "svc_canned"}}</option><option value="macros">{{T "svc_macros"}}</option><option value="translate">{{T "svc_translate"}}</option><option value="widget">{{T "svc_widget"}}</option><option value="email">{{T "svc_email"}}</option><option value="blacklist">{{T "svc_blacklist"}}</option><option value="files">{{T "svc_files"}}</option><option value="calendar">{{T "svc_calendar"}}</option><option value="knowledge">{{T "svc_knowledge"}}</option><option value="audit">{{T "svc_audit"}}</option><option value="backup">{{T "svc_backup"}}</option></select></div>
        <div class="form-group"><label>{{T "pkg_footermark"}}</label><select name="footermark" class="form-control"><option value="2">{{T "pkg_footermark_off"}}</option><option value="1">{{T "pkg_footermark_on"}}</option></select></div>
        <hr><h6 class="text-uppercase text-muted small">{{T "pkg_limits"}}</h6>
        <div class="form-row">
        <div class="form-group col-4"><label>{{T "pkg_limit_send"}}</label><input name="send_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_receive"}}</label><input name="receive_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_ussd"}}</label><input name="ussd_limit" type="number" class="form-control" value="0"></div></div>
        <div class="form-row">
        <div class="form-group col-4"><label>{{T "pkg_limit_device"}}</label><input name="device_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_wa_send"}}</label><input name="wa_send_limit" type="number" class="form-control" value="100"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_wa_receive"}}</label><input name="wa_receive_limit" type="number" class="form-control" value="100"></div></div>
        <div class="form-row">
        <div class="form-group col-3"><label>{{T "pkg_limit_wa"}}</label><input name="wa_account_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_contact"}}</label><input name="contact_limit" type="number" class="form-control" value="50"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_scheduled"}}</label><input name="scheduled_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-3"><label>{{T "aik_api_key"}}</label><input name="key_limit" type="number" class="form-control" value="5"></div></div>
        <div class="form-row">
        <div class="form-group col-4"><label>{{T "pkg_limit_webhook"}}</label><input name="webhook_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_action"}}</label><input name="action_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-4"><label>{{T "pkg_limit_meta"}}</label><input name="meta_limit" type="number" class="form-control" value="0"></div></div>
        <div class="form-row">
        <div class="form-group col-3"><label>{{T "pkg_limit_drips"}}</label><input name="drip_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_recurring"}}</label><input name="recurring_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_forms"}}</label><input name="form_limit" type="number" class="form-control" value="1"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_template"}}</label><input name="template_limit" type="number" class="form-control" value="5"></div></div>
        <div class="form-row">
        <div class="form-group col-3"><label>{{T "pkg_limit_canned"}}</label><input name="canned_limit" type="number" class="form-control" value="10"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_macros"}}</label><input name="macro_limit" type="number" class="form-control" value="5"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_ai_key"}}</label><input name="ai_key_limit" type="number" class="form-control" value="3"></div>
        <div class="form-group col-3"><label>{{T "pkg_limit_knowledge"}}</label><input name="knowledge_limit" type="number" class="form-control" value="10"></div></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_packages"}}</h4></div>
      {{if .EditID}}<div class="border-bottom p-3 bg-light"><form method="post" action="/admin/packages/edit"><input type="hidden" name="id" value="{{.EditID}}"><h6>Edit #{{.EditID}}: {{.EditName}}</h6>
        <div class="form-row"><div class="form-group col-8"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div><div class="form-group col-4"><label>{{T "pkg_price"}}</label><input name="price" class="form-control" value="{{.EditPrice}}"></div></div>
        <hr><h6 class="text-uppercase text-muted small">{{T "pkg_limits"}}</h6>
        <div class="form-row"><div class="form-group col-3"><label>{{T "pkg_limit_send"}}</label><input name="send_limit" type="number" class="form-control" value="{{.EditSendLimit}}"></div><div class="form-group col-3"><label>{{T "pkg_limit_device"}}</label><input name="device_limit" type="number" class="form-control" value="{{.EditDeviceLimit}}"></div><div class="form-group col-3"><label>{{T "pkg_limit_wa"}}</label><input name="wa_account_limit" type="number" class="form-control" value="{{.EditWaAccountLimit}}"></div><div class="form-group col-3"><label>{{T "pkg_limit_contact"}}</label><input name="contact_limit" type="number" class="form-control" value="{{.EditContactLimit}}"></div></div>
        <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{T "pkg_update"}}</button> <a href="/admin/packages" class="btn btn-white btn-sm ms-2">{{T "pkg_cancel"}}</a></form></div>{{end}}
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "pkg_price"}}</th><th>{{T "pkg_limits"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Packages}}<tr><td>{{.ID}}</td><td><strong>{{.Name}}</strong><br><small class="text-muted">{{.Services}}</small></td><td>{{.Price}}</td><td><small>Dev:{{.DeviceLimit}} WA:{{.WaAccountLimit}} Meta:{{.MetaLimit}} AI:{{.AiKeyLimit}}</small></td><td class="text-nowrap"><a class="btn btn-sm btn-white" href="/admin/packages?edit={{.ID}}"><i class="la la-edit"></i></a> <form method="post" action="/admin/packages/delete" style="display:inline" onsubmit="return confirm('{{T "ar_confirm_delete"}}')"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger"><i class="la la-trash"></i></button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
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
      <form method="post" action="/admin/subscriptions/add"><div class="form-group"><label>{{T "col_user"}}</label><select name="user_id" class="form-control" required>{{range .Users}}<option value="{{.ID}}">{{.Name}} ({{.Email}})</option>{{else}}<option value="">{{T "sub_no_users"}}</option>{{end}}</select></div><div class="form-group"><label>{{T "adm_packages"}}</label><select name="package_id" class="form-control">{{range .Packages}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></div><div class="form-group"><label>{{T "sub_expire"}}</label><input name="expire" type="date" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_subscriptions"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_user"}}</th><th>{{T "adm_packages"}}</th><th>{{T "sub_expire"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
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
      <form method="post" action="/admin/marketing/add"><div class="form-group"><label>{{T "pg_title"}}</label><input name="title" class="form-control" required></div><div class="form-group"><label>{{T "col_message"}}</label><textarea name="content" class="form-control" rows="4"></textarea></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">{{T "btn_cancel"}}</a>{{end}}</form></div></div></div>
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
    <div class="col-12 col-lg-5"><div class="card border-warning"><div class="card-header bg-warning bg-opacity-10"><h5 class="card-header-title mb-0">{{T "was_edit_title"}} #{{.EditID}}</h5></div><div class="card-body">
      <form method="post" action="/admin/waservers/edit"><input type="hidden" name="id" value="{{.EditID}}">
    {{else}}
    <div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "was_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/waservers/add">
    {{end}}
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" value="{{.EditName}}" required></div>
        <div class="form-row"><div class="form-group col-6"><label>{{T "was_accounts"}}</label><input name="accounts" type="number" class="form-control" value="{{if .EditID}}{{.EditContent}}{{else}}100{{end}}"></div>
        <div class="form-group col-6"><label>{{T "adm_packages"}}</label><select name="packages" class="form-control" multiple>{{range .Packages}}<option value="{{.Name}}" {{if and $.EditID (contains $.EditGroups .Name)}}selected{{end}}>{{.Name}}</option>{{end}}</select></div></div>
        <div class="form-row"><div class="form-group col-8"><label>{{T "was_url"}}</label><input name="url" class="form-control" placeholder="http://127.0.0.1" value="{{.EditContent}}"></div>
        <div class="form-group col-4"><label>{{T "was_port"}}</label><input name="port" class="form-control" placeholder="8080" value="{{.EditPhone}}"></div></div>
        <div class="form-group"><label>{{T "was_secret"}}</label><input name="secret" class="form-control" value="{{.EditKeyword}}"></div>
        <button class="btn btn-primary lift"><i class="la la-save me-1"></i> {{if .EditID}}{{T "was_save"}}{{else}}{{T "ar_add_btn"}}{{end}}</button>
        {{if .EditID}}<a href="/admin/waservers" class="btn btn-white ms-2">{{T "btn_cancel"}}</a>{{end}}
      </form></div></div></div>
    <div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_waservers"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>URL</th><th>{{T "was_accounts"}}</th><th>{{T "col_packages"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
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
      <form method="post" action="/admin/shorteners/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">{{T "btn_cancel"}}</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_shorteners"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Shorteners}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td><form method="post" action="/admin/shorteners/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="3" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_plugins"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "plg_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/plugins/add"><div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" required></div><div class="form-group"><label>{{T "plg_dir"}}</label><input name="dir" class="form-control"></div><button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>{{if .EditID}}<a href="/admin/waservers" class="btn btn-white btn-sm ms-2">{{T "btn_cancel"}}</a>{{end}}</form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "adm_plugins"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "plg_dir"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .Plugins}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Dir}}</td><td><form method="post" action="/admin/plugins/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
{{end}}

{{if eq .Page "admin_meta"}}
  <div class="row">
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "meta_add_title"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/meta/add">
        <div class="form-group"><label>{{T "col_name"}}</label><input name="name" class="form-control" placeholder="My Business" required></div>
        <div class="form-group"><label>{{T "meta_phone_id"}}</label><input name="phone_number_id" class="form-control" placeholder="123456789..." required></div>
        <div class="form-group"><label>{{T "meta_access_token"}}</label><input name="access_token" class="form-control" placeholder="EAA..." required></div>
        <div class="form-group"><label>{{T "meta_app_id"}}</label><input name="app_id" class="form-control" placeholder="123456..."></div>
        <div class="form-group"><label>{{T "meta_app_secret"}}</label><input name="app_secret" class="form-control" placeholder="abc123..."></div>
        <div class="form-group"><label>{{T "meta_verify_token"}}</label><input name="verify_token" class="form-control" placeholder="chatgo_webhook_123"></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "meta_list_title"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "meta_phone_id"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .MetaAccounts}}<tr><td>{{.ID}}</td><td>{{.Name}} <span class="badge badge-soft-primary" style="font-size:9px">Meta</span></td><td>{{.PhoneNumberID}}</td><td><form method="post" action="/admin/meta/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">-</td></tr>{{end}}
      </tbody></table></div></div></div>
  </div>
  <div class="card mt-3"><div class="card-header"><h4 class="card-header-title">{{T "meta_webhook_title"}}</h4></div>
  <div class="card-body">
    <p class="small text-muted">{{T "meta_webhook_desc"}}</p>
    <code id="webhookUrl" style="word-break:break-all">{{.AppURL}}/webhook/meta</code>
    <p class="small text-muted mt-2">{{T "meta_webhook_verify"}}</p>
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
    <div class="col-12 col-lg-4"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "meta_tpl_add"}}</h4></div><div class="card-body">
      <form method="post" action="/admin/metatemplates/add">
        <div class="form-group"><label>{{T "meta_tpl_name"}}</label><input name="name" class="form-control" placeholder="hello_world" required></div>
        <div class="form-group"><label>{{T "meta_tpl_lang"}}</label><select name="language" class="form-control"><option value="id">{{T "meta_tpl_lang_id"}}</option><option value="en">{{T "meta_tpl_lang_en"}}</option><option value="en_US">{{T "meta_tpl_lang_en_us"}}</option></select></div>
        <div class="form-group"><label>{{T "meta_tpl_cat"}}</label><select name="category" class="form-control"><option value="marketing">{{T "meta_tpl_cat_marketing"}}</option><option value="utility">{{T "meta_tpl_cat_utility"}}</option><option value="authentication">{{T "meta_tpl_cat_auth"}}</option></select></div>
        <div class="form-group"><label>{{T "meta_tpl_comp"}}</label><textarea name="components" class="form-control" rows="4" placeholder='[{"type":"body","text":"Halo {{1}}, pesanan {{2}} sudah diproses"}]'></textarea></div>
        <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
      </form></div></div></div>
    <div class="col-12 col-lg-8"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "meta_tpl_list"}}</h4></div>
      <div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_name"}}</th><th>{{T "meta_tpl_lang"}}</th><th>{{T "meta_tpl_cat"}}</th><th>{{T "col_action"}}</th></tr></thead><tbody>
        {{range .MetaTemplates}}<tr><td>{{.ID}}</td><td>{{.Name}}</td><td>{{.Language}}</td><td>{{.Category}}</td><td><form method="post" action="/admin/metatemplates/delete" style="display:inline"><input type="hidden" name="id" value="{{.ID}}"><button class="btn btn-sm btn-danger">{{T "ar_delete"}}</button></form></td></tr>{{else}}<tr><td colspan="5" class="text-muted text-center">-</td></tr>{{end}}
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
.inbox-split{display:flex;height:calc(100vh - 140px);min-height:500px}
.panel-left{width:380px;min-width:340px;border-right:1px solid #e5e7eb;display:flex;flex-direction:column;overflow:hidden}
.panel-right{flex:1;display:flex;flex-direction:column;background:#efeae2;overflow:hidden}
.chat-empty{display:flex;align-items:center;justify-content:center;height:100%;color:#8696a0;font-size:16px;flex-direction:column;gap:8px}
.chat-empty i{font-size:64px;opacity:.3}
.chat-header{display:flex;align-items:center;padding:10px 16px;background:#f0f2f5;border-bottom:1px solid #e0e0e0;gap:12px;min-height:60px}
.chat-area{flex:1;overflow-y:auto;padding:16px 40px;background-color:#efeae2;background-image:url('data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22200%22 height=%22200%22><rect width=%22200%22 height=%22200%22 fill=%22%23efeae2%22/><circle cx=%22100%22 cy=%22100%22 r=%2260%22 fill=%22%23e5ddd5%22 opacity=%220.5%22/></svg>')}
.chat-footer{border-top:1px solid #e0e0e0;padding:8px 16px;background:#f0f2f5}
.chat-bubble{max-width:70%;padding:6px 10px 6px 10px;border-radius:8px;word-wrap:break-word;box-shadow:0 1px 0.5px rgba(0,0,0,.13);font-size:14.2px;line-height:1.4;position:relative;margin-bottom:3px}
.chat-bubble.received{background:#fff;align-self:flex-start;border-top-left-radius:0}
.chat-bubble.sent{background:#d9fdd3;align-self:flex-end;border-top-right-radius:0}
.chat-bubble.sent.phone-sync{background:#dbeafe;border:1px solid #93c5fd}
.chat-bubble.sent.phone-sync::after{content:'📱';position:absolute;top:-8px;right:4px;font-size:10px}
.chat-sender{font-size:12.5px;font-weight:600;color:#10B981;margin-bottom:2px}
.chat-time{font-size:11px;color:#667781;float:right;margin-left:8px;margin-top:2px}
.chat-input-wrap{display:flex;align-items:center;gap:8px}
.chat-input-wrap textarea{flex:1;resize:none;border-radius:8px;padding:9px 12px;min-height:42px;max-height:100px;border:1px solid #e0e0e0;font-size:14px;outline:none}
.chat-input-wrap textarea:focus{border-color:#00a884}
.chat-input-wrap button{width:42px;height:42px;border-radius:50%;border:none;background:#00a884;color:#fff;cursor:pointer;flex-shrink:0;display:flex;align-items:center;justify-content:center}
.conv-item{cursor:pointer;transition:background .15s;border-bottom:1px solid #f0f0f0;padding:12px 16px;display:flex;gap:12px;align-items:center}
.conv-item:hover{background:#f0f2f5}
.conv-item.active{background:#e8f3ff}
.conv-item.unread{background:#eef2ff}
.conv-item .conv-name{font-size:15px;font-weight:500;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:200px}
.conv-item .conv-msg{font-size:13px;color:#667781;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:220px}
.conv-item .conv-time{font-size:11px;color:#8696a0;white-space:nowrap}
@media(max-width:767.98px){
.inbox-split{height:calc(100vh - 120px);min-height:400px}
.panel-left{width:100%;min-width:0}
.panel-right{display:none}
.panel-right.open{display:flex;position:fixed;top:0;left:0;right:0;bottom:0;z-index:1050;background:#efeae2}
.chat-area{padding:12px 16px}
.conv-item{padding:10px 12px}
.conv-item .conv-name{max-width:52vw}
.conv-item .conv-msg{max-width:58vw}
.chat-bubble{max-width:85%}
.chat-input-wrap textarea{font-size:16px}
}
</style>
<div class="inbox-split" id="inboxSplit">
<div class="panel-left">
<div style="padding:12px 16px;background:#f0f2f5;border-bottom:1px solid #e0e0e0;display:flex;justify-content:space-between;align-items:center">
<strong>{{T "inbox_chat_title"}}{{if gt .UnreadCount 0}} <span class="badge badge-danger">{{.UnreadCount}}</span>{{end}}</strong>
<div style="width:180px"><input type="text" id="inboxSearch" class="form-control form-control-sm" placeholder="{{T "inbox_search"}}"></div>
</div>
<div style="overflow-y:auto;flex:1" id="convList">
{{range .InboxConversations}}
<div class="conv-item{{if gt .Unread 0}} unread{{end}}" data-phone="{{.Phone}}" data-name="{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}" data-group="{{.IsGroup}}" data-channel="{{.Channel}}" onclick="openChat('{{.Phone}}','{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}','{{.IsGroup}}','{{.Channel}}')">
<div class="avatar {{if .IsGroup}}group{{else}}person{{end}}" style="width:48px;height:48px;font-size:16px">{{if .IsGroup}}G{{else}}{{slice .Name 0 1}}{{if not .Name}}+{{end}}{{end}}</div>
<div class="flex-grow-1 min-w-0">
<div class="d-flex justify-content-between"><span class="conv-name">{{if .Name}}{{.Name}}{{else}}+{{.Phone}}{{end}}</span><span class="conv-time">{{.LastTime}}</span></div>
<div class="d-flex align-items-center gap-2"><span class="conv-msg">{{.LastMsg}}</span>{{if gt .Unread 0}}<span class="badge badge-pill" style="background:#25d366;font-size:10px;min-width:20px">{{.Unread}}</span>{{end}}</div>
</div></div>{{else}}<div class="text-center text-muted py-4">{{T "inbox_empty"}}</div>{{end}}
</div></div>
<div class="panel-right" id="chatPanel">
<div class="chat-empty" id="chatEmpty"><i class="la la-comments"></i><div>{{T "inbox_select_hint"}}</div></div>
<div id="chatView" style="display:none;flex-direction:column;height:100%">
<div class="chat-header" id="chatHeader">
<button class="btn btn-sm d-md-none" onclick="closeChat()" style="border:none;background:none;font-size:20px">&larr;</button>
<div class="avatar person" style="width:40px;height:40px;font-size:14px" id="chatAvatar">+</div>
<div class="flex-grow-1"><strong id="chatTitle">-</strong><div><small class="text-muted" id="chatSubtitle"></small></div></div>
<select id="chatAccountPhone" class="form-select form-select-sm" style="width:auto">{{range .ConnectedAccounts}}{{if eq .Status "connected"}}<option value="+{{.Phone}}">+{{.Phone}}</option>{{end}}{{end}}</select>
<button class="btn btn-sm btn-outline-secondary" onclick="exportChat()" title="Export CSV"><i class="la la-download"></i></button>
</div>
<div class="chat-area" id="chatMessages"></div>
<div class="chat-footer">
<div class="chat-input-wrap">
<textarea id="chatInput" rows="1" placeholder="{{T "inbox_chat_ph"}}" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();sendChatMsg()}"></textarea>
<button onclick="sendChatMsg()"><i class="la la-send" style="font-size:18px"></i></button>
</div></div></div></div></div>
<script>
var chatPhone='',chatName='',chatIsGroup='';
function openChat(phone,name,isGroup,channel){
document.querySelectorAll('.conv-item').forEach(function(e){e.classList.remove('active')});
var el=document.querySelector('.conv-item[data-phone="'+phone+'"]');if(el)el.classList.add('active');
chatPhone=phone;chatName=name;chatIsGroup=isGroup==='true';
document.getElementById('chatEmpty').style.display='none';
document.getElementById('chatView').style.display='flex';
document.getElementById('chatTitle').textContent=name;
document.getElementById('chatSubtitle').textContent=(isGroup==='true'?'Group':'')+(channel==='meta'?' Meta':'');
document.getElementById('chatAvatar').textContent=name.charAt(0)||'+';
if(window.innerWidth<768)document.getElementById('chatPanel').classList.add('open');
loadMessages();
}
function closeChat(){document.getElementById('chatPanel').classList.remove('open');}
function loadMessages(){
fetch('/inbox/messages?phone='+encodeURIComponent(chatPhone)).then(function(r){return r.json()}).then(function(msgs){
var box=document.getElementById('chatMessages');
if(!msgs||!msgs.length){box.innerHTML='<div class="text-center text-muted py-4">{{T "inbox_no_messages"}}</div>';scrollChat();return}
var html='';
for(var i=0;i<msgs.length;i++){
var m=msgs[i],side=m.type==='sent'?'flex-end':'flex-start';
var bubble='chat-bubble '+m.type;
if(m.type==='sent'&&m.channel==='phone_sync')bubble+=' phone-sync';
html+='<div class="d-flex w-100 mb-1" style="justify-content:'+side+'"><div class="'+bubble+'">';
if(m.type==='received'&&m.sender_name)html+='<div class="chat-sender">'+m.sender_name+'</div>';
html+='<div>'+m.message+'<span class="chat-time">'+m.created+'</span></div></div></div>';
}
box.innerHTML=html;scrollChat();
});
}
function sendChatMsg(){
var inp=document.getElementById('chatInput'),msg=inp.value.trim();
if(!msg||!chatPhone)return;
var f=new FormData();f.append('phone',chatPhone);f.append('message',msg);
var acp=document.getElementById('chatAccountPhone');if(acp)f.append('account_phone',acp.value);
fetch('/inbox/send',{method:'POST',body:f}).then(function(r){return r.json()}).then(function(d){
if(d.ok){inp.value='';loadMessages();inp.focus();}
});
}
function scrollChat(){var box=document.getElementById('chatMessages');if(box)box.scrollTop=box.scrollHeight;}
function exportChat(){
fetch('/inbox/messages?phone='+encodeURIComponent(chatPhone)).then(function(r){return r.json()}).then(function(msgs){
var csv='type,message,time\n';msgs.forEach(function(m){csv+='"'+m.type+'","'+m.message.replace(/"/g,'""')+'","'+m.created+'"\n'});
var blob=new Blob(['\uFEFF'+csv],{type:'text/csv'});
var a=document.createElement('a');a.href=URL.createObjectURL(blob);a.download=chatPhone+'.csv';a.click();
});
}
var evtSource=new EventSource('/inbox/events');
evtSource.onmessage=function(e){var d=JSON.parse(e.data);if(d.phone===chatPhone)loadMessages();};
document.getElementById('inboxSearch').addEventListener('input',function(){
var q=this.value.toLowerCase();
document.querySelectorAll('.conv-item').forEach(function(el){el.style.display=el.textContent.toLowerCase().includes(q)?'':'none'});
});
scrollChat();
</script>
{{end}}
{{if eq .Page "inbox_chat"}}
<style>
.chat-area{height:55vh;overflow-y:auto;padding:16px;background:#efeae2;background-image:url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><rect width="200" height="200" fill="%23efeae2"/><circle cx="100" cy="100" r="60" fill="%23e5ddd5" opacity="0.5"/></svg>')}
.chat-bubble{max-width:75%;padding:8px 12px;border-radius:8px;word-wrap:break-word;box-shadow:0 1px 1px rgba(0,0,0,.08);font-size:14px;line-height:1.4;position:relative}
.chat-bubble.received{background:#fff;align-self:flex-start;border-top-left-radius:0}
.chat-bubble.sent{background:#d9fdd3;align-self:flex-end;border-top-right-radius:0}
.chat-bubble.sent.phone-sync{background:#dbeafe;border:1px solid #93c5fd;border-top-right-radius:0}
.chat-bubble.sent.phone-sync::after{content:'📱';position:absolute;top:-8px;right:4px;font-size:10px}
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
<div><strong>{{if .ChatName}}{{.ChatName}}{{else}}+{{.Phone}}{{end}}</strong>{{if .IsGroup}}<small class="text-success ms-1">{{T "type_group"}}</small>{{end}}{{if .Channel}}<small class="badge badge-soft-primary ms-1">{{.Channel}}</small>{{end}}</div>
</h6>
<div class="d-flex gap-2 align-items-center">
<select id="chatChannel" class="form-select form-select-sm" style="width:auto;font-size:12px" onchange="onChannelChange()">
<option value="whatsmeow">{{T "inbox_wa_channel"}}</option>
{{if .MetaAccounts}}<option value="meta">{{T "inbox_meta_channel"}}</option>{{end}}
</select>
<select id="chatAccountPhone" class="form-select form-select-sm" style="width:auto;display:inline"> {{range .ConnectedAccounts}}{{if eq .Status "connected"}}<option value="+{{.Phone}}">+{{.Phone}}</option>{{end}}{{end}}</select>
<select id="chatMetaAccount" class="form-select form-select-sm" style="width:auto;display:none"> {{range .MetaAccounts}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select>
{{if .Users}}<form method="post" action="/inbox/assign" style="display:inline" class="me-1"><input type="hidden" name="phone" value="{{.Phone}}"><select name="agent_id" class="form-select form-select-sm" style="width:auto;font-size:11px" onchange="this.form.submit()"><option value="0">{{T "inbox_unassigned"}}</option>{{range .Users}}<option value="{{.ID}}">{{.Name}}</option>{{end}}</select></form><form method="post" action="/inbox/close" style="display:inline"><input type="hidden" name="phone" value="{{.Phone}}"><button class="btn btn-sm btn-outline-danger" style="font-size:11px;padding:2px 8px">{{T "inbox_close_btn"}}</button></form>{{end}}</div>
</div>
<div class="card-body p-0">
<div class="chat-area" id="chatMessages">
{{range .ChatMessages}}
<div class="d-flex w-100 mb-1" style="{{if eq .Type "sent"}}justify-content:flex-end{{else}}justify-content:flex-start{{end}}">
<div class="chat-bubble {{.Type}}{{if eq .Channel "phone_sync"}} phone-sync{{end}}">
{{if and (eq .Type "received") .SenderName}}<div class="chat-sender">{{.SenderName}}</div>{{end}}
<div>{{.Message}}<span class="chat-time">{{.Created}}</span></div>
</div>
</div>
{{else}}
<div class="text-center text-muted py-4">{{T "inbox_chat_empty"}}</div>
{{end}}
</div>
</div>
<div class="card-footer bg-white border-top" style="padding:8px 16px">
{{if .Notes}}<div class="mb-2" style="max-height:100px;overflow-y:auto">{{range .Notes}}<div class="small text-muted mb-1"><i class="la la-sticky-note me-1"></i> {{.Note}} <span class="text-muted" style="font-size:10px">{{.Created}}</span></div>{{end}}</div>{{end}}
<div class="d-flex gap-1 mb-1">
  <form method="post" action="/inbox/note" class="d-flex gap-1 flex-grow-1"><input type="hidden" name="phone" value="{{.Phone}}"><input name="note" class="form-control form-control-sm" placeholder="{{T "inbox_note_ph"}}" style="font-size:12px"><button class="btn btn-sm btn-outline-secondary" style="font-size:11px"><i class="la la-sticky-note"></i></button></form>
</div>
<form id="chatForm" onsubmit="return sendChat(event)">
<div class="chat-input-group">
<textarea id="chatInput" name="message" class="form-control" placeholder="{{T "inbox_chat_ph"}}" rows="1" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();sendChat(event)}"></textarea>
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
if(!msgs||!msgs.length){chatBox.innerHTML='<div class="text-center text-muted py-4">{{T "inbox_no_messages"}}</div>';return}
var html='';
for(var i=0;i<msgs.length;i++){
var m=msgs[i];
var side=m.type==='sent'?'flex-end':'flex-start';
var bubbleClass='chat-bubble '+m.type;
if(m.type==='sent'&&m.channel==='phone_sync')bubbleClass+=' phone-sync';
html+='<div class="d-flex w-100 mb-1" style="justify-content:'+side+'"><div class="'+bubbleClass+'">';
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
.docs-section h3{font-weight:700;margin-bottom:16px;padding-bottom:8px;border-bottom:2px solid #eee}
.docs-section h4{font-weight:600;margin:16px 0 8px}
.docs-step{background:#f8f9fc;border-left:3px solid #4F46E5;padding:12px 16px;margin:8px 0;border-radius:0 8px 8px 0}
.docs-step strong{color:#4F46E5}
pre code{font-size:13px}
</style>
<div class="row">
<div class="col-12 col-lg-3"><div class="docs-nav">
<a href="#quick">{{T "docs_nav_quick"}}</a>
<a href="#wa">{{T "docs_nav_wa"}}</a>
<a href="#contacts">{{T "docs_nav_contacts"}}</a>
<a href="#broadcast">{{T "docs_nav_broadcast"}}</a>
<a href="#drip">{{T "docs_nav_drip"}}</a>
<a href="#ai">{{T "docs_nav_ai"}}</a>
<a href="#inbox">{{T "docs_nav_inbox"}}</a>
<a href="#store">{{T "docs_nav_store"}}</a>
<a href="#payment">{{T "docs_nav_payment"}}</a>
<a href="#team">{{T "docs_nav_team"}}</a>
<a href="#analytics">{{T "docs_nav_analytics"}}</a>
<a href="#tools">{{T "docs_nav_tools"}}</a>
<a href="#api">{{T "docs_nav_api"}}</a>
</div></div>
<div class="col-12 col-lg-9">

<div class="docs-section" id="quick">
<h3>{{T "docs_quick_title"}}</h3>
<p>{{T "landing_demo_label"}} <code>{{.AppEmail}}</code> / <code>password</code></p>
<div class="docs-step"><strong>1. Hubungkan WA:</strong> Buka WhatsApp di HP → Perangkat Tertaut → Tautkan Perangkat. Di ChatGo: buka <code>/wa</code> → Tambah Akun → Scan QR.</div>
<div class="docs-step"><strong>2. Tambah Kontak:</strong> <code>/contacts</code> → import CSV atau tambah manual.</div>
<div class="docs-step"><strong>3. Kirim Broadcast:</strong> <code>/broadcast</code> → pilih grup → tulis pesan → kirim massal.</div>
<div class="docs-step"><strong>4. Setup AI:</strong> Tambah API key di <code>/ai/keys</code> → buat rule di <code>/autoreply</code>.</div>
<div class="docs-step"><strong>5. Live Chat:</strong> <code>/inbox</code> — semua pesan masuk real-time.</div>
</div>

<div class="docs-section" id="wa">
<h3>Hubungkan WhatsApp</h3>
<table class="table table-sm">
<tr><td width="180"><strong>WhatsApp Web</strong> <span class="badge badge-soft-success">WA Web</span></td><td>Scan QR di <code>/wa</code> → perangkat tertaut. Multi-akun, unlimited nomor (sesuai paket).</td></tr>
<tr><td><strong>Meta Cloud API</strong> <span class="badge" style="background:#4F46E5;color:#fff;font-size:10px">META</span></td><td><code>/admin/meta</code> → input Phone ID + Token dari Facebook. Support template, webhook, tidak perlu HP online.</td></tr>
<tr><td><strong>Web Widget</strong></td><td>Embed <code>&lt;script src="/widget.js"&gt;&lt;/script&gt;</code> di website. Chat muncul di pojok kanan bawah.</td></tr>
<tr><td><strong>Email to WA</strong></td><td>POST ke <code>/email-webhook</code> — forward email ke WA inbox.</td></tr>
</table>
</div>

<div class="docs-section" id="contacts">
<h3>Kontak & Groups</h3>
<div class="docs-step"><strong>1. Tambah Manual:</strong> <code>/contacts</code> → nama + nomor + group</div>
<div class="docs-step"><strong>2. CSV Import:</strong> Upload CSV (name, phone, groups). Group auto-create.</div>
<div class="docs-step"><strong>3. CSV Export:</strong> Tombol Export di Contacts</div>
<div class="docs-step"><strong>4. Bulk Delete:</strong> Centang → Delete</div>
<div class="docs-step"><strong>5. Groups:</strong> <code>/contacts/groups</code> — atur grup + language ID/EN per grup</div>
<div class="docs-step"><strong>6. Tags:</strong> <code>/tags</code> — label warna (VIP, Leads). Filter di broadcast.</div>
<div class="docs-step"><strong>7. Merge Duplikat:</strong> <code>/merge</code> — auto-detect & gabung</div>
<div class="docs-step"><strong>8. Number Validator:</strong> Tombol Validate di broadcast form</div>
</table>
</div>

<div class="docs-section" id="broadcast">
<h3>Broadcast</h3>
<div class="docs-step"><strong>1.</strong> <code>/broadcast</code> → nama campaign + pilih grup target + nomor langsung</div>
<div class="docs-step"><strong>2.</strong> Upload media (gambar/video/dokumen) — dikirim bersama pesan</div>
<div class="docs-step"><strong>3.</strong> Pilih akun WA + mode (Round Robin / Random)</div>
<div class="docs-step"><strong>4.</strong> Interval + Rate Limiter di Settings</div>
<div class="docs-step"><strong>5.</strong> Pause / Resume / Stop / Retry campaign</div>
<div class="docs-step"><strong>6.</strong> Link Tracking auto + A/B Testing</div>
<div class="docs-step"><strong>7.</strong> Recurring Campaign: <code>/recurring</code> — jadwal daily/weekly</div>
<div class="docs-step"><strong>8.</strong> Campaign Calendar: <code>/calendar</code></div>
</div>

<div class="docs-section" id="drip">
<h3>Drip Campaign — Follow-up Otomatis</h3>
<div class="docs-step"><strong>1. Buat Drip:</strong> <code>/drips</code> → nama + add steps (message + delay menit)</div>
<div class="docs-step"><strong>2. Auto-Enroll:</strong> Setiap chat masuk → otomatis masuk semua drip aktif</div>
<div class="docs-step"><strong>3. STOP:</strong> User reply "STOP" atau "berhenti" → unenroll</div>
<div class="docs-step"><strong>4. Contoh Step:</strong> 0: "Halo Kak!" → +120min: "Ada yang bisa dibantu?" → +1day: "Promo spesial!"</div>
</div>

<div class="docs-section" id="ai">
<h3>AI Auto Reply</h3>
<div class="docs-step"><strong>1. AI Key:</strong> <code>/ai/keys</code> → tambah OpenAI/Gemini/Claude/DeepSeek</div>
<div class="docs-step"><strong>2. Knowledge Base:</strong> <code>/knowledge</code> → Q&A manual atau import CSV/PDF</div>
<div class="docs-step"><strong>3. Buat Rule:</strong> <code>/autoreply</code> → Match Type (Contains/Exact/AI) → centang AI</div>
<div class="docs-step"><strong>4. AI Global:</strong> Settings → AI untuk Semua Pesan</div>
<div class="docs-step"><strong>5. Store Agent:</strong> AI auto-inject katalog produk + profil customer</div>
<div class="docs-step"><strong>6. Human Handoff:</strong> Keyword trigger → stop AI → kontak admin</div>
<div class="docs-step"><strong>7. Working Hours:</strong> Jam kerja + pesan luar jam</div>
<div class="docs-step"><strong>8. Variables:</strong> {name} {phone} {message}. Spintax: {Halo|Hai}</div>
</div>

<div class="docs-section" id="inbox">
<h3>Live Chat Inbox</h3>
<div class="docs-step"><strong>1.</strong> <code>/inbox</code> — WhatsApp-style. Tab Chat + Status</div>
<div class="docs-step"><strong>2.</strong> Filter: Semua / Private / Group / Unread</div>
<div class="docs-step"><strong>3.</strong> Channel: WA Web atau Meta (dropdown di chat header)</div>
<div class="docs-step"><strong>4.</strong> Canned Responses: <code>/canned</code> — tombol shortcut di bawah input</div>
<div class="docs-step"><strong>5.</strong> Templates: <code>/templates</code> — klik untuk insert</div>
</div>

<div class="docs-section" id="store">
<h3>WA Store Bot</h3>
<div class="docs-step"><strong>1. Setup:</strong> <code>/store</code> → tambah kategori + produk (nama, harga, gambar)</div>
<div class="docs-step"><strong>2. Flow Customer:</strong> Chat "menu" → lihat kategori → pilih produk → auto-order</div>
<div class="docs-step"><strong>3. Pembayaran:</strong> Chat "BAYAR" → diarahkan ke halaman pembayaran</div>
<div class="docs-step"><strong>4. Order:</strong> <code>/store/orders</code> → update status. WA notif otomatis.</div>
</div>

<div class="docs-section" id="payment">
<h3>Payment Gateway</h3>
<div class="docs-step"><strong>1.</strong> <code>/admin/gateways-pay</code> → tambah Midtrans/PayPal/Stripe/Xendit</div>
<div class="docs-step"><strong>2.</strong> Buat Package: <code>/admin/packages</code> → harga + limit + services (fitur)</div>
<div class="docs-step"><strong>3.</strong> User subscribe: <code>/subscribe</code> → pilih package → gateway → bayar</div>
<div class="docs-step"><strong>4.</strong> Auto-activate: callback → verify → subscription aktif</div>
<div class="docs-step"><strong>5.</strong> Reminder: <code>/reminders</code> — jadwal tagihan + auto-WA</div>
</div>

<div class="docs-section" id="team">
<h3>Team & Support</h3>
<div class="docs-step"><strong>Agent Assignment:</strong> Auto round-robin. Assign manual via dropdown.</div>
<div class="docs-step"><strong>Departments:</strong> <code>/depts</code> — Sales/Support/Billing. Auto-detect keyword.</div>
<div class="docs-step"><strong>Chat Transfer + Notes + Labels + Close + CSAT Survey</strong></div>
<div class="docs-step"><strong>Macros:</strong> <code>/macros</code> — one-click assign+tag+reply+close</div>
<div class="docs-step"><strong>Auto-Close + VIP Priority + Agent Signature</strong></div>
</div>

<div class="docs-section" id="analytics">
<h3>Analytics</h3>
<div class="docs-step"><strong>Dashboard:</strong> Chart 7 hari sent vs received</div>
<div class="docs-step"><strong>Agent:</strong> <code>/analytics</code> — chats, replies, response time</div>
<div class="docs-step"><strong>CSAT:</strong> <code>/csat</code> & <code>/analytics</code> — rating survey</div>
<div class="docs-step"><strong>Link Tracker:</strong> <code>/tracker</code> — URL clicks</div>
<div class="docs-step"><strong>Audit:</strong> <code>/audit</code> — user activity log</div>
</div>

<div class="docs-section" id="tools">
<h3>Tools</h3>
<div class="docs-step"><strong>Translate:</strong> <code>/translate-tool</code> — AI translate ke 9 bahasa</div>
<div class="docs-step"><strong>Forms:</strong> <code>/forms</code> — interactive form builder</div>
<div class="docs-step"><strong>Files:</strong> <code>/uploads</code> — browse uploaded media</div>
<div class="docs-step"><strong>Backup:</strong> <code>/backup</code> — one-click DB backup</div>
<div class="docs-step"><strong>Safety:</strong> <code>/blacklist</code> — auto spam detection. Settings — Rate Limiter</div>
</div>

<div class="docs-section" id="api">
<h3>API Reference</h3>
<pre class="bg-light p-3 rounded"><code>POST /api/send
Header: X-API-Key: YOUR_KEY
Body: {"phone":"628xx","message":"text"}

GET /api/contacts
Header: X-API-Key: YOUR_KEY</code></pre>
<p class="small text-muted">API Keys: <code>/apikeys</code>. Webhooks: <code>/webhooks</code>.</p>
</div>

</div></div>
<script>
(function(){
var n=document.querySelectorAll('.docs-nav a');
var s=document.querySelectorAll('.docs-section');
addEventListener('scroll',function(){
var st=scrollY+90;
s.forEach(function(sec,i){
var t=sec.offsetTop,h=sec.offsetHeight;
n.forEach(function(a){a.classList.remove('active')});
if(st>=t && st<t+h && n[i+1])n[i+1].classList.add('active');
});
});
n.forEach(function(a){a.addEventListener('click',function(e){e.preventDefault();var el=document.getElementById(this.getAttribute('href').slice(1));if(el)el.scrollIntoView({behavior:'smooth',block:'start'})})});
})();
</script>
{{end}}


{{if eq .Page "knowledge"}}
  <div class="row">
    <div class="col-12 col-lg-5">
      <div class="card"><div class="card-header"><h4 class="card-header-title">{{T "kb_add"}}</h4></div><div class="card-body">
        <form method="post" action="/knowledge/add">
          <div class="form-group"><label>{{T "kb_title"}}</label><input name="title" class="form-control" placeholder="{{T "ar_faq_tab"}}" required></div>
          <div class="form-group"><label>{{T "kb_question"}}</label><input name="question" class="form-control" placeholder="{{T "kb_question_dot"}}..." required></div>
          <div class="form-group"><label>{{T "kb_answer"}}</label><textarea name="answer" class="form-control" rows="3" placeholder="{{T "kb_answer_dot"}}..." required></textarea></div>
          <div class="form-group"><label>{{T "kb_category"}}</label><input name="category" class="form-control" placeholder="{{T "kb_placeholder_category"}}"></div>
          <button class="btn btn-primary lift"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/import" enctype="multipart/form-data">
          <label class="text-muted small mb-2 d-block">{{T "kb_import"}}</label>
<div class="input-group"><input type="text" name="title" class="form-control" placeholder="{{T "kb_placeholder_title"}}"><input type="file" name="file" class="form-control" accept=".csv,.txt" required><button class="btn btn-white">{{T "kb_upload"}}</button></div>
          <small class="form-text text-muted">{{T "kb_csv_hint"}} <a href="/web/sample-knowledge.csv" target="_blank">{{T "kb_sample"}}</a></small>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/url">
          <label class="text-muted small mb-2 d-block">{{T "kb_url"}}</label>
<div class="input-group"><input type="text" name="title" class="form-control" placeholder="{{T "kb_placeholder_title"}}"><input type="url" name="url" class="form-control" placeholder="https://..." required><button class="btn btn-white">{{T "kb_train"}}</button></div>
          <small class="form-text text-muted">{{T "kb_url_hint"}}</small>
        </form>
        <hr class="my-3">
        <form method="post" action="/knowledge/pdf" enctype="multipart/form-data">
          <label class="text-muted small mb-2 d-block">📄 {{T "kb_upload_pdf"}}</label>
          <div class="input-group"><input type="text" name="title" class="form-control" placeholder="{{T "kb_placeholder_title"}}"><input type="file" name="file" class="form-control" accept=".pdf" required><button class="btn btn-white">{{T "kb_upload"}}</button></div>
          <small class="form-text text-muted">{{T "kb_pdf_hint"}}</small>
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
{{end}}
{{if eq .Page "faq"}}
<div class="row">
<div class="col-12 col-lg-5"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "faq_add_title"}}</h4></div><div class="card-body">
<form method="post" action="/faq/add"><div class="form-group"><label>{{T "col_question"}}</label><input name="question" class="form-control" placeholder="{{T "faq_question_ph"}}" required></div><div class="form-group"><label>{{T "col_answer"}}</label><textarea name="answer" class="form-control" rows="3" placeholder="{{T "faq_answer_ph"}}" required></textarea></div><button class="btn btn-primary"><i class="la la-plus me-1"></i> {{T "ar_add_btn"}}</button></form>
<hr><h6>{{T "kb_import"}}</h6><form method="post" action="/faq/import" enctype="multipart/form-data"><div class="form-group"><label>{{T "faq_import_csv"}}</label><input type="file" name="file" class="form-control" accept=".csv" required></div><button class="btn btn-outline-primary btn-sm"><i class="la la-upload me-1"></i> {{T "btn_import"}}</button></form>
</div></div></div>
<div class="col-12 col-lg-7"><div class="card"><div class="card-header"><h4 class="card-header-title">{{T "faq_list_title"}} <small class="text-muted">— {{T "faq_subtitle"}}</small></h4></div>
<div class="table-responsive"><table class="table table-sm card-table"><thead><tr><th>#</th><th>{{T "col_question"}}</th><th>{{T "col_answer"}}</th><th></th></tr></thead><tbody>
{{range .FAQ}}<tr><td>{{.id}}</td><td>{{.question}}</td><td>{{.answer}}</td><td><form method="post" action="/faq/delete" style="display:inline"><input type="hidden" name="id" value="{{.id}}"><button class="btn btn-sm btn-danger">&times;</button></form></td></tr>{{else}}<tr><td colspan="4" class="text-muted text-center">{{T "faq_empty"}}</td></tr>{{end}}</tbody></table></div></div></div></div>
{{end}}{{end}}`

