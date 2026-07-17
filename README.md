# 🏗️ Chatforge — Open Source WhatsApp AI Gateway

**🇬🇧 English** | [🇮🇩 Bahasa Indonesia](#-chatforge--open-source-whatsapp-ai-gateway-1)

Single-binary WhatsApp marketing platform with AI auto-reply, multi-account, broadcast, and SaaS management. Zero external dependency. Pure Go.

---

## 🆘 Installation Help

Need help with setup, installation, or custom development?

**📱 WhatsApp: [+62 812-9605-2010](https://wa.me/6281296052010)**

---

## ✨ Features

### 🤖 AI Auto Reply
- **BYOK** — Bring Your Own Key (OpenAI, Gemini, Claude, DeepSeek)
- **Function calling** — AI auto-searches your FAQ before answering
- **Training campaigns** — per-rule system prompt + AI key override
- **Fallback mode** — AI only activates when no keyword matches
- **Memory window** — N-chat conversation context for AI
- **Human handoff** — keyword triggers to stop AI & route to admin
- **Business hours** — timezone-aware gating
- **Reasoning level** — Low / Medium / High (temperature control)
- **Token tracking** — per-user token usage counter

### 📱 WhatsApp
- **Multi-account** — manage multiple WA numbers in one panel
- **QR pairing** — scan & connect via WhatsApp Linked Devices
- **Broadcast** — send to contact groups, configurable interval, round-robin
- **Scheduled messages** — time-based, recurring per minute
- **Auto Reply** — keyword matching: Contains, Exact, Starts With, AI
- **Media messages** — send images, videos, documents
- **Webhook dispatch** — real-time notifications on send/receive
- **Inbox** — conversation view per contact

### 👥 SaaS Platform
- **Multi-user** — admin + user roles with menu filtering
- **Role-based permissions** — dynamic role system with per-feature access control
- **Packages & subscriptions** — limit devices, contacts, drips, templates, AI keys, and more per plan
- **Multi-tenant isolation** — user_id on all tables, inbox filtered per user, session ownership
- **WA Server enforcement** — restrict accounts by package
- **Dashboard analytics** — 7-day charts, user stats, active accounts
- **Session persistence** — MySQL-backed, survives server restarts
- **Force own key** — reseller can require sub-users to use their own API key
- **Security** — bcrypt passwords, SHA-256 API keys, AES-256-GCM encrypted secrets, HttpOnly cookies

### 🔧 Developer
- **REST API** — `POST /api/send`, `GET /api/status`, `/api/messages`, `/api/contacts`, `/api/devices`
- **Single binary** — `go build` = one `.exe`, deploy anywhere
- **No Node.js** — uses [whatsmeow](https://github.com/tulir/whatsmeow) (pure Go)
- **Multi-language** — English + Indonesian via `lang/*.json`
- **Spintax** — `{Hello|Hi|Hey}` random per message
- **.env config** — MySQL DSN, AES encryption key, app URL
- **Auto-migration** — column-existence check before ALTER TABLE, safe restarts

### 📋 Additional Features
- **Drip Campaigns** — multi-step automated message sequences
- **Recurring Campaigns** — auto-repeat broadcast on schedule
- **A/B Testing** — split-test message variants
- **Canned Responses** — quick reply shortcuts
- **Contact Tags & Groups** — organize & segment contacts
- **CSAT Surveys** — post-chat satisfaction rating
- **Store & Orders** — product catalog with WhatsApp ordering
- **Forms & Reminders** — interactive forms, payment reminders
- **Web Widget** — embeddable chat widget
- **Email → WA Gateway** — forward emails to WhatsApp
- **Link Tracker** — URL shortener with click analytics
- **File Manager** — upload & share media
- **Blacklist** — block spam numbers
- **Macros** — one-click multi-action workflows
- **Auto Translate** — on-the-fly message translation via AI

---

## 🚀 Quick Start

```bash
# Clone
git clone https://github.com/linducip2208/chatforge.git
cd chatforge

# Build (single binary)
go build -o chatforge.exe .

# Database
mysql -u root -e "CREATE DATABASE chatforge CHARACTER SET utf8mb4"

# Run
./chatforge.exe
```

Open **http://127.0.0.1:8080**  
Login: `admin@chatgo.test` / `password`

---

## 📦 Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go |
| Database | MySQL |
| WhatsApp | [whatsmeow](https://github.com/tulir/whatsmeow) |
| UI | Bootstrap 5 |
| Charts | Chart.js |
| Screenshots | Playwright |
| Session | MySQL (survives restart) |

---

## 📸 Screenshots

| | | |
|---|---|---|
| ![Login](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/01-login.png) | ![Dashboard](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/02-dashboard.png) | ![Auto Reply](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/07-autoreply-ai.png) |
| **Login** | **Dashboard** | **Auto Reply AI** |
| ![Broadcast](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/05-broadcast.png) | ![Contacts](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/09-contacts.png) | ![Admin](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/14-admin-users.png) |
| **Broadcast** | **Contacts** | **Admin Users** |

---

## 🏗️ Architecture

```
chatforge (single Go binary)
├── wa/           — WhatsApp engine (whatsmeow sessions, send, auto-reply, loops)
├── store/        — MySQL database layer (CRUD, migrations, multi-tenant isolation)
│   ├── store.go          — core (autoreplies, inbox, sent/received, sessions, settings)
│   ├── store_admin.go    — users, roles, packages, AI keys, devices, permissions
│   ├── store_extra.go    — contacts, groups, campaigns, scheduled, templates, tags
│   ├── store_safety.go   — blacklist, CSAT, spam detection
│   ├── store_drip.go     — drip campaigns, steps, enrollments
│   ├── store_meta.go     — WhatsApp Cloud API accounts, templates
│   ├── store_payment.go  — payment gateways, transactions
│   ├── store_plus.go     — departments, recurring, notes, labels
│   ├── store_final.go    — audit logs, inbox macros
│   └── store_knowledge.go— knowledge base, AI trainings
├── aiservice/    — AI provider adapter (OpenAI-compatible, BYOK)
├── meta/         — WhatsApp Cloud API client (Meta Graph API v22)
├── msgtemplate/  — Spintax + variable render engine
├── i18n/         — Multi-language JSON loader
├── secret/       — AES-256-GCM encryption for API keys & secrets
├── payment/      — Payment gateway adapter (Midtrans, Xendit, Tripay, Duitku)
├── lang/         — id.json + en.json translation files
└── web/          — Static assets (CSS, JS, Font Awesome, Line Awesome, images)
```

### Security Architecture
- **Passwords** — bcrypt with auto-upgrade from legacy SHA-256
- **API Keys** — SHA-256 hashed, shown once at creation
- **Secrets** — AES-256-GCM encrypted (env key via `CHATGO_ENC_KEY`)
- **Sessions** — HttpOnly + Secure + SameSite cookies
- **Multi-tenant** — `user_id` on 15+ tables, ownership-validated CRUD, inbox filtered by WA session owner

---

## 🔌 Live Chat Widget

Embed Chatforge as a **WhatsApp live chat widget** on any website:

```html
<a href="https://wa.me/6281296052010?text=Hello%20I%20have%20a%20question" 
   target="_blank" 
   style="position:fixed;bottom:20px;right:20px;width:60px;height:60px;background:#25D366;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 12px rgba(0,0,0,.2);z-index:9999">
  <svg width="32" height="32" viewBox="0 0 24 24" fill="white"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/></svg>
</a>
```

---

## 📖 Usage Guide

### 1. Connect WhatsApp
Go to **Account & QR** → click Add Account → scan QR code with WhatsApp (Linked Devices).

### 2. Send Message
Go to **Send Message** → select your WA number, enter recipient phone, type message, click Send.

### 3. Auto Reply
Go to **Auto Reply** → add rule:
- **Match Type**: Contains / Exact / Starts With / AI
- **Keyword**: trigger word
- **Reply**: response text (supports spintax `{Hi|Hello}`)
- **AI**: check "Use AI" and select an AI Key
- **Account**: select which WA number the rule applies to

### 4. Broadcast
Go to **Broadcast** → enter campaign name, message, select groups or paste numbers, check sender numbers, click Send.

### 5. Contacts
Go to **Contacts** → add contacts manually or import CSV. Create groups, assign tags, merge duplicates.

### 6. Drip Campaign
Go to **Drip Campaign** → create campaign, add steps with delay, activate. Users auto-enroll when they message you.

### 7. AI Setup
Go to **AI Keys** (Admin) → add provider (OpenAI, DeepSeek, Gemini, etc.), paste API key, select model. Then use in Auto Reply rules.

### 8. Packages & Limits
Go to **Packages** (Admin) → create plans with limits (devices, contacts, drips, templates, AI keys). Assign subscriptions to users.

### 9. API Usage
```bash
# Send message
curl -X POST http://localhost:8080/api/send \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"phone":"628123456789","message":"Hello World"}'

# Check status
curl http://localhost:8080/api/status -H "X-API-Key: YOUR_API_KEY"

# List contacts
curl http://localhost:8080/api/contacts -H "X-API-Key: YOUR_API_KEY"
```

### 10. Environment Variables
```env
# .env file
CHATGO_MYSQL=root:password@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4
CHATGO_ENC_KEY=your-32-byte-aes-key-here
APP_URL=https://your-domain.com
APP_NAME=ChatGo
```

### 11. Deploy to Linux
```bash
# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o chatgo_linux .

# Upload
scp chatgo_linux .env user@server:/opt/chatgo/

# Systemd service
sudo cp chatgo.service /etc/systemd/system/
sudo systemctl enable chatgo --now
```

---

## ⚠️ Disclaimer

This is an unofficial WhatsApp client using the [whatsmeow](https://github.com/tulir/whatsmeow) library. Use at your own risk. Not affiliated with or endorsed by Meta Platforms, Inc. or WhatsApp LLC.

**META** and **WhatsApp** are registered trademarks of Meta Platforms, Inc.

---

---

# 🏗️ Chatforge — Open Source WhatsApp AI Gateway

[🇬🇧 English](#-chatforge--open-source-whatsapp-ai-gateway) | **🇮🇩 Bahasa Indonesia**

Platform WhatsApp marketing single-binary dengan AI auto-reply, multi-akun, broadcast, dan manajemen SaaS. Tanpa dependency eksternal. Pure Go.

## 🆘 Bantuan Instalasi

Butuh bantuan setup, instalasi, atau custom development?

**📱 WhatsApp: [+62 812-9605-2010](https://wa.me/6281296052010)**

## ✨ Fitur

### 🤖 AI Auto Reply
- **BYOK** — Bawa API Key Sendiri (OpenAI, Gemini, Claude, DeepSeek)
- **Function calling** — AI otomatis cari FAQ sebelum menjawab
- **Training campaigns** — system prompt + AI key per rule
- **Fallback mode** — AI hanya jalan saat tidak ada keyword match
- **Memory window** — kirim N chat terakhir ke AI sebagai konteks
- **Human handoff** — keyword trigger `admin|operator` → stop AI → kirim kontak admin
- **Business hours** — batasi AI hanya di jam kerja
- **Reasoning level** — Low / Medium / High (temperature control)
- **Token tracking** — hitung pemakaian token per user

### 📱 WhatsApp
- **Multi-account** — kelola banyak nomor WA dalam satu panel
- **QR pairing** — scan & connect via WhatsApp Linked Devices
- **Broadcast** — kirim massal ke grup kontak, interval configurable, round-robin
- **Scheduled messages** — jadwal kirim pesan, recurring per menit
- **Auto Reply** — keyword matching: Contains, Exact, Starts With, AI
- **Media messages** — kirim gambar, video, dokumen
- **Webhook dispatch** — notifikasi real-time saat pesan masuk/keluar
- **Inbox** — lihat percakapan per kontak

### 👥 SaaS Platform
- **Multi-user** — admin + user roles dengan menu filtering
- **Role-based permissions** — sistem role dinamis dengan kontrol akses per fitur
- **Packages & subscriptions** — batasi device, kontak, drips, template, AI key, dan lainnya per paket
- **Multi-tenant isolation** — user_id di semua tabel, inbox terfilter per user, kepemilikan session
- **WA Server enforcement** — user hanya bisa pakai server yg diizinkan paketnya
- **Dashboard analytics** — chart 7 hari, stats user, nomor WA aktif
- **Session persistence** — MySQL-backed, survive restart server
- **Force own key** — reseller bisa paksa sub-user pakai API key sendiri
- **Security** — bcrypt passwords, SHA-256 API keys, AES-256-GCM encrypted secrets, HttpOnly cookies

### 🔧 Developer
- **REST API** — `POST /api/send`, `GET /api/status`, `/api/messages`, `/api/contacts`, `/api/devices`
- **Single binary** — `go build` = 1 file `.exe`, deploy ke mana aja
- **No Node.js** — pakai [whatsmeow](https://github.com/tulir/whatsmeow) (pure Go)
- **Multi-language** — Indonesia + English via `lang/*.json`
- **Spintax** — `{Halo|Hai|Hi}` random tiap kirim
- **Auto-migration** — cek kolom sebelum ALTER TABLE, restart aman

### 📋 Fitur Tambahan
- **Drip Campaigns** — rangkaian pesan otomatis multi-step
- **Recurring Campaigns** — broadcast ulang otomatis sesuai jadwal
- **A/B Testing** — uji coba varian pesan
- **Canned Responses** — shortcut balasan cepat
- **Contact Tags & Groups** — atur & segmentasi kontak
- **CSAT Surveys** — rating kepuasan setelah chat
- **Store & Orders** — katalog produk dengan order via WhatsApp
- **Forms & Reminders** — form interaktif, pengingat pembayaran
- **Web Widget** — widget chat embed
- **Email → WA Gateway** — teruskan email ke WhatsApp
- **Link Tracker** — pemendek URL dengan analytics klik
- **File Manager** — upload & bagikan media
- **Blacklist** — blokir nomor spam
- **Macros** — workflow multi-aksi satu klik
- **Auto Translate** — terjemahan pesan otomatis via AI

## 📖 Panduan Penggunaan

### 1. Hubungkan WhatsApp
Buka **Account & QR** → klik Add Account → scan QR code dengan WhatsApp (Linked Devices).

### 2. Kirim Pesan
Buka **Send Message** → pilih nomor WA pengirim, masukkan nomor tujuan, ketik pesan, klik Send.

### 3. Auto Reply
Buka **Auto Reply** → tambah rule:
- **Match Type**: Contains / Exact / Starts With / AI
- **Keyword**: kata pemicu
- **Reply**: teks balasan (support spintax `{Halo|Hai}`)
- **AI**: centang "Use AI" lalu pilih AI Key
- **Account**: pilih nomor WA mana yang menerapkan rule ini

### 4. Broadcast
Buka **Broadcast** → isi nama campaign, pesan, pilih grup atau tempel nomor, centang nomor pengirim, klik Send.

### 5. Kontak
Buka **Contacts** → tambah kontak manual atau import CSV. Buat grup, beri tag, gabung duplikat.

### 6. Drip Campaign
Buka **Drip Campaign** → buat campaign, tambah step dengan delay, aktifkan. User otomatis masuk saat kirim pesan.

### 7. Setup AI
Buka **AI Keys** (Admin) → tambah provider (OpenAI, DeepSeek, Gemini, dll.), tempel API key, pilih model. Lalu pakai di rule Auto Reply.

### 8. Paket & Limit
Buka **Packages** (Admin) → buat paket dengan limit (device, kontak, drips, template, AI key). Assign subscription ke user.

### 9. Penggunaan API
```bash
# Kirim pesan
curl -X POST http://localhost:8080/api/send \
  -H "X-API-Key: API_KEY_ANDA" \
  -H "Content-Type: application/json" \
  -d '{"phone":"628123456789","message":"Halo Dunia"}'

# Cek status
curl http://localhost:8080/api/status -H "X-API-Key: API_KEY_ANDA"

# Lihat kontak
curl http://localhost:8080/api/contacts -H "X-API-Key: API_KEY_ANDA"
```

### 10. Environment Variables
```env
# File .env
CHATGO_MYSQL=root:password@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4
CHATGO_ENC_KEY=kunci-aes-32-byte-anda
APP_URL=https://domain-anda.com
APP_NAME=ChatGo
```

### 11. Deploy ke Linux
```bash
# Build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o chatgo_linux .

# Upload
scp chatgo_linux .env user@server:/opt/chatgo/

# Systemd service
sudo cp chatgo.service /etc/systemd/system/
sudo systemctl enable chatgo --now
```

## ⚠️ Disclaimer

Ini adalah WhatsApp client tidak resmi menggunakan library [whatsmeow](https://github.com/tulir/whatsmeow). Gunakan dengan risiko sendiri. Tidak berafiliasi dengan Meta Platforms, Inc. atau WhatsApp LLC.

**META** dan **WhatsApp** adalah merek dagang terdaftar Meta Platforms, Inc.
