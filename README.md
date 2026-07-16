# 🏗️ Chatforge — Open Source WhatsApp AI Gateway

Single-binary WhatsApp marketing platform with AI auto-reply, multi-account, broadcast, and SaaS management. Zero external dependency. Pure Go.

---

## 🆘 Bantuan Instalasi

Butuh bantuan setup, instalasi, atau custom development?

**📱 WhatsApp: [+62 812-9605-2010](https://wa.me/6281296052010)**

---

## ✨ Features

### 🤖 AI Auto Reply
- **BYOK** — Bring Your Own Key (OpenAI, Gemini, Claude, DeepSeek)
- **Function calling** — AI auto-searches your FAQ before answering
- **Training campaigns** — per-rule system prompt + AI key override
- **Fallback mode** — AI hanya jalan saat tidak ada keyword match
- **Memory window** — kirim N chat terakhir ke AI sebagai konteks
- **Human handoff** — keyword trigger `admin|operator` → stop AI → kirim kontak admin
- **Business hours** — batasi AI hanya di jam kerja
- **Reasoning level** — Low / Medium / High (temperature control)
- **Token tracking** — hitung pemakaian token per user

### 📱 WhatsApp
- **Multi-account** — kelola banyak nomor WA dalam satu panel
- **QR pairing** — scan & connect via WhatsApp Linked Devices
- **Broadcast** — kirim massal ke grup kontak, interval configurable, round-robin multi-nomor
- **Scheduled messages** — jadwal kirim pesan, recurring per menit
- **Auto Reply** — keyword matching: Contains, Exact, Starts With, AI
- **Media messages** — kirim gambar, video, dokumen
- **Webhook dispatch** — notifikasi real-time saat pesan masuk/keluar
- **Inbox** — lihat percakapan per kontak

### 👥 SaaS Platform
- **Multi-user** — admin + user roles dengan menu filtering
- **Packages & subscriptions** — batasi device, token, fitur per paket
- **WA Server enforcement** — user hanya bisa pakai server yg diizinkan paketnya
- **Dashboard analytics** — chart 7 hari, stats user, nomor WA aktif
- **Session persistence** — MySQL-backed, survive restart server
- **Force own key** — reseller bisa paksa sub-user pakai API key sendiri

### 🔧 Developer
- **REST API** — `POST /api/send`, `GET /api/status`, `/api/messages`, `/api/contacts`, `/api/devices`
- **Single binary** — `go build` = 1 file `.exe`, deploy ke mana aja
- **No Node.js** — pakai [whatsmeow](https://github.com/tulir/whatsmeow) (pure Go library)
- **Multi-language** — Indonesia + English via `lang/*.json`
- **Spintax** — `{Halo|Hai|Hi}` random tiap kirim
- **.env config** — MySQL DSN + listen address

---

## 🚀 Quick Start

```bash
# Clone
git clone https://github.com/linducipta/chatforge.git
cd chatforge

# Build (single binary)
go build -o chatforge.exe .

# Database
mysql -u root -e "CREATE DATABASE chatforge CHARACTER SET utf8mb4"

# Run
./chatforge.exe
```

Buka **http://127.0.0.1:8080**  
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
| ![Login](screens/01-login.png) | ![Dashboard](screens/02-dashboard.png) | ![Auto Reply](screens/07-autoreply-ai.png) |
| **Login** | **Dashboard** | **Auto Reply AI** |
| ![Broadcast](screens/05-broadcast.png) | ![Contacts](screens/09-contacts.png) | ![Admin](screens/14-admin-users.png) |
| **Broadcast** | **Contacts** | **Admin Users** |

---

## 🏗️ Architecture

```
chatforge.exe (single binary)
├── wa/         — WhatsApp engine (whatsmeow)
├── store/      — MySQL database layer
├── aiservice/  — AI provider adapter (OpenAI-compatible)
├── msgtemplate/— Spintax + variable engine
├── i18n/       — Multi-language JSON loader
├── secret/     — AES encryption for API keys
├── lang/       — id.json + en.json
└── web/        — Static assets (CSS, JS, images)
```

---

## 🔌 Live Chat Widget

Chatforge bisa di-embed ke website sebagai **live chat WhatsApp widget**. Contoh:

```html
<!-- Floating WhatsApp Button -->
<a href="https://wa.me/6281296052010?text=Halo%20saya%20mau%20tanya" 
   target="_blank" 
   style="position:fixed;bottom:20px;right:20px;width:60px;height:60px;background:#25D366;border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 12px rgba(0,0,0,.2);z-index:9999">
  <svg width="32" height="32" viewBox="0 0 24 24" fill="white"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/></svg>
</a>
```

---

## ⚠️ Disclaimer

This is an unofficial WhatsApp client using the [whatsmeow](https://github.com/tulir/whatsmeow) library. Use at your own risk. Not affiliated with or endorsed by Meta Platforms, Inc. or WhatsApp LLC.

**META** and **WhatsApp** are registered trademarks of Meta Platforms, Inc.
