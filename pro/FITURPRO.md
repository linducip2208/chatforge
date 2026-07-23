# ChatGo PRO Features — $799

> **⚠️ NEVER PUSH REAL IMPLEMENTATIONS TO GITHUB**
> Semua file implementasi disimpan di folder `pro/` dengan build tag `//go:build pro`.
> Folder `pro/` di-gitignore (kecuali `pro/stub.go`).
> Distributed via `pro.zip` + `license.key` ke customer.

---

## 1. Message Status Tracking

**File:** `pro/message_status.go`

Tracking real-time status pengiriman pesan WhatsApp:
- `sent` → pesan terkirim ke server WhatsApp
- `delivered` → pesan sampai ke perangkat penerima
- `read` → pesan dibaca (centang biru)
- `failed` → pesan gagal

**API:**
```
GET  /pro/message-status/{messageID}   → return status + timestamp
GET  /pro/message-stats                → summary semua message stats
```

**Data:** whatsmeow sudah menyediakan ACK events. Tinggal capture + simpan di DB.

**Integrasi dengan existing:** Update `sent` table tambah kolom `delivered_at`, `read_at`.

---

## 2. Message Buttons

**File:** `pro/message_buttons.go`

Kirim pesan interaktif dengan tombol:
- **CTA Buttons** — "Beli Sekarang" / "Hubungi CS" → URL atau quick reply
- **List Message** — menu pilihan (max 10 items) — "Makanan", "Minuman", "Snack"
- **Quick Reply Buttons** — pilihan cepat (max 3) — "Ya", "Tidak", "Nanti"

**API:**
```
POST /pro/send/buttons       → kirim pesan dengan tombol
POST /pro/send/list          → kirim list message
POST /pro/send/reply-buttons → kirim quick reply buttons
```

**Data:** whatsmeow `SendMessage` + `waE2E.ButtonsMessage` / `waE2E.ListMessage`.

**Integrasi dengan existing:** Tambah opsi di form kirim pesan + auto-reply + broadcast.

---

## 3. Message Attributes

**File:** `pro/message_attributes.go`

Formatting pesan yang lebih kaya:
- **Header** — text, image, video, document di atas pesan
- **Footer** — teks kecil di bawah pesan
- **Body** — bold, italic, strikethrough, monospace
- **Template Variables** — {name}, {phone}, {date}, {amount}

**Integrasi dengan existing:** Auto-reply + template system.

---

## 4. Check WA Number

**File:** `pro/wa_number_check.go`

Validasi apakah nomor terdaftar di WhatsApp:
- Single check: `POST /pro/check-number` → `{ phone, exists, wa_name }`
- Bulk check: `POST /pro/check-numbers` → array of results
- Auto-check saat import kontak
- Auto-check sebelum broadcast

**Data:** whatsmeow `IsOnWhatsApp()` method.

---

## 5. Telegram Bot

**File:** `pro/telegram.go`

Integrasi Telegram Bot API:
- **Setup:** Bot token + webhook registration
- **Receive:** Webhook inbound messages (private + group)
- **Send:** Text, image, video, document, buttons (inline keyboard)
- **Auto-reply:** Keyword matching + AI (mirror WA auto-reply)
- **Broadcast:** Kirim massal ke subscriber Telegram

**Routes:**
```
POST /pro/telegram/webhook           → inbound messages
GET  /pro/telegram/setup             → setup page
POST /pro/telegram/send              → send message
POST /pro/telegram/broadcast         → broadcast
```

**DB Tables:** `telegram_bots`, `telegram_subscribers`, `telegram_messages`.

---

## 6. Instagram DM

**File:** `pro/instagram.go`

Integrasi Instagram Messaging via Meta Graph API:
- **Setup:** Facebook App + Instagram Business Account + Webhook
- **Receive:** DM inbound via Meta webhook
- **Send:** Text, image, story mentions, comment replies
- **Auto-reply:** keyword + AI
- **Auto-engage:** komentar → DM otomatis (social auto-engagement)

**Routes:**
```
GET  /pro/instagram/webhook          → Meta verification
POST /pro/instagram/webhook          → inbound messages + comments
GET  /pro/instagram/setup            → setup page
POST /pro/instagram/send             → send DM
```

**DB Tables:** `instagram_accounts`, `instagram_messages`, `instagram_comments`.

---

## 7. Facebook Messenger

**File:** `pro/facebook.go`

Integrasi Facebook Messenger via Meta Graph API:
- **Setup:** Facebook Page + Webhook
- **Receive:** Messenger inbound
- **Send:** Text, image, video, buttons, quick replies, templates
- **Auto-reply:** keyword + AI

**Routes:**
```
GET  /pro/facebook/webhook           → Meta verification
POST /pro/facebook/webhook           → inbound messages
GET  /pro/facebook/setup             → setup page
POST /pro/facebook/send              → send message
```

**DB Tables:** `facebook_pages`, `facebook_messages`.

---

## 8. Omnichannel Inbox

**File:** `pro/omnichannel_inbox.go`

Satu inbox untuk semua channel:
```
┌────────────────────────────────────────┐
│  OMNI INBOX                            │
│                                        │
│  Filter: [WA] [IG] [FB] [TG] [Web]    │
│                                        │
│  ┌──────────────────────────────────┐  │
│  │ Contact List (left panel)        │  │
│  │ 🟢 WA  +62812...  "Mau beli"     │  │
│  │ 🟣 IG  @customer  "Ready stok?"  │  │
│  │ 🔵 FB  User 123  "Jam buka?"    │  │
│  │ 🔷 TG  @buyer99   "Kirim mana?" │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │ Chat View (right panel)         │  │
│  │ Channel badge + message history  │  │
│  │ Reply input + channel selector   │  │
│  └──────────────────────────────────┘  │
└────────────────────────────────────────┘
```

**Routes:**
```
GET  /pro/omni/inbox               → omnichannel inbox page
GET  /pro/omni/events              → SSE multi-channel events
POST /pro/omni/send                → send via selected channel
```

**Integrasi dengan existing:** Extend `/inbox` system. Tambah `channel` column di `received` table.

---

## 9. Visual Flow Builder

**File:** `pro/flow_builder.go`, `pro/flow_builder.js`

Drag-and-drop visual editor untuk chatbot flow.

**Node types:**
| Node | Icon | Fungsi |
|------|------|--------|
| Trigger | 🚀 | Keyword / welcome / fallback |
| Send Message | 💬 | Text, image, buttons, list |
| Condition | 🔀 | If/else branching (multi-output) |
| Wait/Delay | ⏱️ | Jeda N detik |
| API Call | 🔌 | HTTP request ke endpoint |
| Set Variable | 📦 | Simpan data ke variable |
| Tag Contact | 🏷️ | Assign tag |
| Transfer Agent | 👤 | Lempar ke CS manusia |
| Close Chat | ✅ | Tutup percakapan |

**Routes:**
```
GET  /pro/flow-builder             → editor UI
POST /pro/flow-builder/save        → save flow JSON
GET  /pro/flow-builder/load/{id}   → load flow
DELETE /pro/flow-builder/{id}      → delete flow
GET  /pro/flow-builder/list        → list saved flows
```

**Format Flow JSON:**
```json
{
  "id": "flow_1",
  "name": "Welcome Flow",
  "trigger": { "type": "keyword", "value": "halo,hai,menu" },
  "nodes": [
    { "id": "n1", "type": "message", "data": { "text": "Halo! Mau pesan apa?" } },
    { "id": "n2", "type": "buttons", "data": { "items": ["Makanan", "Minuman"] } },
    { "id": "n3", "type": "condition", "data": { "field": "{{input}}", "branches": { "Makanan": "n4", "Minuman": "n5" } } }
  ],
  "edges": [
    { "from": "n1", "to": "n2" },
    { "from": "n2", "to": "n3" }
  ]
}
```

**UI:** React Flow + Tailwind CSS → build sebagai static SPA di `web/flow-builder/`.

---

## 10. n8n / Zapier / Make Integration

**File:** `pro/integrations.go`

Webhook-based integration untuk automation platforms:
- **Trigger:** New message, new contact, broadcast complete
- **Action:** Send message, add contact, update contact

**Routes:**
```
POST /pro/webhooks/n8n/{event}      → trigger untuk n8n
POST /pro/webhooks/zapier/{event}   → trigger untuk Zapier
POST /pro/webhooks/make/{event}     → trigger untuk Make.com
```

**Documentation:** Generate n8n node definition JSON.

---

## 11. Google Sheets Sync

**File:** `pro/google_sheets.go`

Auto-sync data ke Google Sheets:
- **Chat Log:** setiap pesan masuk/keluar → append row
- **Lead Capture:** kontak baru → append row
- **Campaign Report:** summary broadcast → new sheet
- **Config:** Sheet ID + Google Service Account credential

**Routes:**
```
GET  /pro/sheets/setup              → setup page
POST /pro/sheets/config             → save sheet config
POST /pro/sheets/test               → test sync
```

---

## 12. Agency Dashboard

**File:** `pro/agency_dashboard.go`

Multi-client management untuk agency/reseller:
- **Client CRUD:** nama, domain, package, status
- **Per-client:** WA accounts, contacts, broadcasts, usage stats
- **White-label per client:** logo, nama app, domain
- **Billing per client:** invoice, payment status
- **Impersonate:** login sebagai client (one-click)

**Routes:**
```
GET  /pro/agency                    → agency dashboard
POST /pro/agency/clients/add        → tambah client
POST /pro/agency/clients/{id}/edit  → edit client
POST /pro/agency/clients/{id}/delete → hapus client
GET  /pro/agency/clients/{id}       → detail client
POST /pro/agency/clients/{id}/impersonate → impersonate
```

**DB Tables:** `agency_clients`, `agency_client_usage`, `agency_invoices`.

---

## Development Order

| # | Fitur | Estimasi | Priority |
|---|-------|----------|----------|
| 1 | Message Status Tracking | 2h | 🔥 Quick win |
| 2 | Message Buttons + Attributes | 3h | 🔥 High impact |
| 3 | Check WA Number | 1h | 🔥 Quick win |
| 4 | Telegram Bot | 4h | 🚀 Test omnichannel |
| 5 | Instagram DM | 6h | Meta API setup |
| 6 | Facebook Messenger | 4h | (share setup dengan IG) |
| 7 | Omnichannel Inbox | 5h | Gabungin #4+#5+#6+#7 |
| 8 | Google Sheets Sync | 3h | |
| 9 | n8n/Zapier/Make | 2h | |
| 10 | Visual Flow Builder | 12h | 🔥 Paling kompleks |
| 11 | Agency Dashboard | 6h | |

**Total: ~48 jam development**

---

## 13. Flow Builder — Advanced Features

### 13.1 AI Node
Node yang pakai AI engine untuk generate reply atau decide next path.
```
[AI Reply] → pakai AI key + system prompt + context variables → generate reply
[AI Decide] → kirim chat history ke AI → AI pilih branch mana
```
**Files:** `pro/flow_ai_node.go`

### 13.2 Flow Templates
Library template siap pakai, import 1-click:
- "Restoran Ordering" — menu → pilih → konfirmasi → payment link
- "CSAT Survey" — close → rating → follow-up / complaint
- "Lead Qualification" — greet → qualify → transfer sales / nurturing
- "Appointment Booking" — tanya tanggal → konfirmasi → reminder
- "Abandoned Cart Recovery" — detect → reminder → discount → checkout
- "Customer Support Triage" — classify → FAQ / transfer agent
- "Event Registration" — greet → collect data → confirm → reminder
- "Product Recommendation" — tanya preferensi → rekomendasi → order

**Files:** `pro/flow_templates.go`, `web/flow_templates/`

### 13.3 Flow Test / Simulate
Test flow sebelum publish. Chat simulator di panel kanan.
```
┌─────────────────┬──────────────────────┐
│   Flow Editor   │   Test Chat          │
│                 │   Bot: Halo! Mau..   │
│                 │   You: Makanan       │
│                 │   Bot: Ini menu...   │
│                 │                      │
│                 │   [Reset] [Step-by-step] │
└─────────────────┴──────────────────────┘
```
**Files:** `pro/flow_simulator.go`

### 13.4 Flow Analytics
Track performa setiap flow:
- Trigger count (berapa kali flow jalan)
- Completion rate (berapa % user sampai end)
- Drop-off per node (node mana user berhenti)
- Avg time to complete
- Revenue generated (dari order yang dibuat)
```
Flow "Welcome": 1,234x | 89% done | avg 3.2 nodes | 45 orders
```
**Files:** `pro/flow_analytics.go`, `store/flow_analytics.go`

### 13.5 E-commerce Nodes
Node yang terintegrasi dengan ChatGo Store:
- **Product Carousel** — kirim katalog produk (nama, harga, gambar)
- **Create Order** — auto-buat order dari input user
- **Payment Link** — generate invoice + payment URL (Midtrans/Xendit)
- **Stock Check** — cek stok sebelum konfirmasi
- **Order Status** — kirim status order terbaru

**Files:** `pro/flow_ecommerce_nodes.go`

### 13.6 Smart CRM Nodes
Auto-enrich contact data dari flow:
- **Save to Contact** — simpan jawaban user ke field contact (nama, email, alamat)
- **Check Order History** — cek riwayat order, total spent
- **Check Tags** — cek apakah contact punya tag tertentu
- **Language Detect** — auto-detect bahasa → multi-language reply
- **Sentiment Check** — deteksi marah/kecewa → auto-transfer agent

**Files:** `pro/flow_crm_nodes.go`

### 13.7 Advanced Logic Nodes
- **Random** — pilih random dari beberapa path (A/B testing, variasi konten)
- **Loop** — ulang N kali atau sampai kondisi terpenuhi
- **Sub-flow Call** — panggil flow lain sebagai subroutine (reusable components)
- **Regex Match** — condition pakai regular expression
- **Numeric Compare** — `{{var.price}} > 100000`
- **Date Compare** — `{{var.due_date}} < today + 3 days`
- **List Contains** — `{{var.tags}} contains "VIP"`
- **Empty Check** — `{{input}} is empty`

**Files:** `pro/flow_advanced_nodes.go`

### 13.8 Flow Versioning
- **Draft/Publish** — flow dalam mode draft sebelum publish
- **Rollback** — kembali ke versi sebelumnya
- **Changelog** — auto-catat perubahan setiap save

**Files:** `pro/flow_versioning.go`

### 13.9 Flow Debugger
- **Step-by-step execution** — jalanin flow 1 node per 1 node
- **Breakpoints** — pause di node tertentu
- **Variable Inspector** — lihat nilai semua variable di setiap step
- **Edge Highlight** — highlight path yang diambil (hijau = diambil, merah = tidak)

**Files:** `pro/flow_debugger.go`

### 13.10 New Trigger Types
- **Webhook** — trigger flow dari HTTP request eksternal
- **Button Click** — user klik button di pesan sebelumnya (Quick Reply, CTA)
- **Cron/Schedule** — trigger flow di jam/tanggal tertentu
- **Inactivity** — user tidak reply selama N jam
- **Order Event** — order dibuat, dibayar, dikirim
- **Tag Change** — contact ditambahkan ke tag tertentu

**Files:** `pro/flow_triggers_advanced.go`

---

## Development Order (Combined)

| # | Fitur | Estimasi | Priority |
|---|-------|----------|----------|
| 1-12 | (Existing Pro features) | 48h | — |
| 13.5 | E-commerce Nodes | 3h | 🔥 |
| 13.6 | Smart CRM Nodes | 3h | 🔥 |
| 13.1 | AI Node | 1h | 🔥 |
| 13.2 | Flow Templates | 2h | 🔥 |
| 13.10 | New Trigger Types | 2h | 🔥 |
| 13.7 | Advanced Logic Nodes | 3h | |
| 13.3 | Flow Test/Simulate | 2h | |
| 13.4 | Flow Analytics | 3h | |
| 13.8 | Flow Versioning | 2h | |
| 13.9 | Flow Debugger | 2h | |

---

## 14. Check WA Number

Validate apakah nomor terdaftar di WhatsApp.

**Single Check:**
```
POST /pro/check-number
Body: {"phone": "6281296052010"}
Response: {"phone": "6281296052010", "exists": true, "wa_name": "Budi"}
```

**Bulk Check:**
```
POST /pro/check-numbers
Body: {"phones": ["62812...", "62813..."]}
Response: [{"phone": "...", "exists": bool}, ...]
```

**UI Integration:**
- Tombol "Validate" di contacts import → cek semua nomor sebelum import
- Broadcast → validasi nomor sebelum kirim
- API endpoint public → bisa dipakai sistem eksternal

**Technology:** whatsmeow `IsOnWhatsApp()` method.

**Files:** `pro/wa_number_check.go`

---

## 15. Flow Builder — WhatsApp-Specific Nodes

### 15.1 Media Gallery (Album)
Kirim multiple gambar/video sebagai album WA (maks 10).
```json
{"type": "gallery", "data": {"items": [
  {"url": "/uploads/img1.jpg", "caption": "Depan"},
  {"url": "/uploads/img2.jpg", "caption": "Belakang"}
]}}
```

### 15.2 Voice Note
Kirim pesan suara. Bisa pre-recorded file atau Text-to-Speech.
```json
{"type": "voice", "data": {"url": "/uploads/audio.ogg", "tts": "Halo, selamat datang!"}}
```

### 15.3 Document / PDF
Kirim file document (PDF, DOCX, XLSX) dengan caption.
```json
{"type": "document", "data": {"url": "/uploads/brosur.pdf", "caption": "Brosur kami"}}
```

### 15.4 Poll Node
Buat WhatsApp Poll interaktif. User vote, hasil masuk ke variable.
```json
{"type": "poll", "data": {
  "question": "Mau pesan apa?",
  "options": ["Makanan", "Minuman", "Snack"],
  "max_select": 1,
  "var_result": "poll_answer"
}}
```

### 15.5 Location Node
Request atau kirim lokasi ke user.
```json
{"type": "location_request", "data": {"button_text": "Share Lokasi"}}
{"type": "location_send", "data": {"lat": -6.2, "lng": 106.8, "name": "Kantor Kami"}}
```

### 15.6 Contact Card (vCard)
Kirim kartu kontak ke user — nama, nomor, alamat.
```json
{"type": "contact_card", "data": {"name": "CS Kami", "phone": "62812...", "org": "ChatGo"}}
```

### 15.7 Typing Indicator
Tampilkan "typing..." selama N detik sebelum kirim pesan — terasa lebih natural.
```json
{"type": "typing", "data": {"seconds": 2}}
```

### 15.8 Scheduled Delay
Tunggu sampai jam/tanggal spesifik (bukan sekian detik).
```json
{"type": "delay_until", "data": {"time": "09:00", "timezone": "Asia/Jakarta"}}
```

### 15.9 Group Mention
Mention specific person di group chat.
```json
{"type": "group_mention", "data": {"phone": "62812...", "text": "Tolong dibantu"}}
```

---

## 16. Flow Builder — Engine Enhancements

### 16.1 Inline Metrics
Setiap node tampilkan counter: X kali dieksekusi, Y kali success, Z kali error.
```
┌──────────┐
│ Message  │  234x ⬆️
│ "Halo!"  │  89%  success
└──────────┘
```

### 16.2 Fallback Path
Setiap node bisa punya "on error" / "on timeout" edge. Kalau API call gagal atau user ga reply, ada path alternatif.

### 16.3 Rate Limit Per Flow
Batas berapa kali flow tertentu jalan per user per jam/hari. Cegah spam.
```json
{"rate_limit": {"per_user": 3, "per_hour": true}}
```

### 16.4 Split / Merge Parallel
Jalankan beberapa path parallel (misal: kirim WA + update Google Sheet bersamaan), lalu merge kembali.

### 16.5 Interaction Tracking
Catat interaksi user di flow: berapa kali user reply, berapa lama reply, berapa kali restart flow.

---

## Development Order (Updated)

| # | Fitur | Estimasi |
|---|-------|----------|
| 14 | Check WA Number | 1h |
| 13.5 | E-commerce Nodes | 3h |
| 13.6 | Smart CRM Nodes | 3h |
| 13.1 | AI Node | 1h |
| 13.2 | Flow Templates | 2h |
| 13.10 | New Trigger Types | 2h |
| 15.1 | Media Gallery | 1h |
| 15.3 | Document Node | 0.5h |
| 15.4 | Poll Node | 2h |
| 15.5 | Location Node | 1h |
| 15.7 | Typing Indicator | 0.5h |
| 15.8 | Scheduled Delay | 1h |
| 13.7 | Advanced Logic Nodes | 3h |
| 13.3 | Flow Test/Simulate | 2h |
| 13.4 | Flow Analytics | 3h |
| 16.1 | Inline Metrics | 2h |
| 16.3 | Rate Limit Per Flow | 1h |

---

## 17. Math + String + Date Nodes

### 17.1 Formula / Math Node
Kalkulasi di flow. Operasi matematika dengan variable.
```json
{"type": "math", "data": {"formula": "{{var.price}} * {{var.qty}} + {{var.shipping}}", "var_result": "total"}}
```
Support: +, -, *, /, %, round, min, max.

### 17.2 String Template Node
Gabungin multiple values jadi satu text.
```json
{"type": "string_template", "data": {
  "template": "Halo {{name}}! Pesanan #{{order_id}} total Rp{{total}}. Estimasi {{eta}} hari.",
  "var_result": "message"
}}
```

### 17.3 Date / Time Node
Manipulasi tanggal di flow.
```json
{"type": "date_math", "data": {"operation": "add_days", "value": 3, "var_result": "due_date"}}
```
Support: add_days, add_hours, format_date, day_of_week, now.

---

## 18. Loop + Sub-flow + Random

### 18.1 Loop Node
Ulangi section flow N kali atau sampai kondisi.
```json
{"type": "loop", "data": {"count": 3, "condition": "{{var.status}} != 'done'"}}
```
Loop body: nodes between loop start and loop end edge.

### 18.2 Sub-flow Call
Panggil flow lain sebagai reusable component.
```json
{"type": "subflow", "data": {"flow_id": 5, "inputs": {"name": "{{var.name}}"}, "outputs": {"result": "var.result"}}}
```
Component library: bikin sekali, pakai di banyak flow.

### 18.3 Random / A/B Node
Pilih random dari beberapa path. Untuk A/B testing.
```json
{"type": "random", "data": {"weights": [50, 30, 20]}}
```
Weight-based distribution untuk variant testing.

---

## 19. Contact Lookup + Counter + DB Query

### 19.1 Contact Lookup Node
Cek data kontak dari database. Auto-enrich flow context.
```json
{"type": "contact_lookup", "data": {"phone": "{{phone}}", "fields": ["last_order", "total_spent", "tags"]}}
```
Fields tersimpan di: {{contact.last_order}}, {{contact.total_spent}}, etc.

### 19.2 Counter / Increment Node
Track berapa kali flow jalan per user. Cegah spam.
```json
{"type": "counter", "data": {"key": "flow_{{flow_id}}_per_user", "max": 3, "var_result": "count"}}
```
Kalau melebihi max → route ke overflow path.

### 19.3 Database Query Node
Execute SQL query dan simpan hasil ke variable.
```json
{"type": "db_query", "data": {"query": "SELECT COUNT(*) FROM orders WHERE phone='{{phone}}'", "var_result": "order_count"}}
```
Read-only. Untuk custom logic yang ga bisa di-cover built-in nodes.

---

## 20. Export / Import + Marketplace + Conditional

### 20.1 Flow Export / Import
Export flow as JSON file. Import di instance lain.
```
GET /pro/flow-builder/export?id=5 → download JSON
POST /pro/flow-builder/import → upload JSON file
```
Portable flows. Share antar team/instance.

### 20.2 Flow Conditional Enable
Flow hanya aktif jika kondisi tertentu terpenuhi.
```json
{"conditional_enable": {"field": "time", "operator": "between", "value": "08:00-17:00"}}
```
Support: time range, day of week, stock available, variable exists.

### 20.3 Flow Marketplace (Future)
Share/publish flow template ke community marketplace.
- Browse public flows
- Rate & review
- One-click install
- Version tracking

---

## 21. WhatsApp-Specific (Additional)

### 21.1 Read Receipt Node
Tunggu sampai pesan sebelumnya dibaca (centang biru), lalu lanjut flow.

### 21.2 Sticker / GIF Node
Kirim sticker atau GIF animasi sebagai bagian dari flow.
```json
{"type": "sticker", "data": {"url": "/uploads/sticker.webp"}}
```

### 21.3 Multiple Message Sequence
Kirim beberapa pesan dalam sequence dengan delay.
```json
{"type": "message_sequence", "data": {"messages": [
  {"text": "Halo!", "delay": 1},
  {"text": "Ada yang bisa dibantu?", "delay": 2}
]}}
```

---

## Development Order (Updated)

| # | Fitur | Estimasi |
|---|-------|----------|
| 17.1 | Math Node | 1h |
| 17.2 | String Template Node | 0.5h |
| 17.3 | Date Node | 1h |
| 17.1 | Loop Node | 2h |
| 18.2 | Sub-flow Call | 2h |
| 18.3 | Random Node | 1h |
| 19.1 | Contact Lookup | 2h |
| 19.2 | Counter Node | 1h |
| 20.1 | Flow Export/Import | 1h |
| 20.2 | Conditional Enable | 1h |
| 21.1 | Read Receipt Node | 2h |
| 21.3 | Message Sequence | 1h |

**Total tambahan: ~15.5 jam**
