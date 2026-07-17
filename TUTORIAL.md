# 🏗️ Chatforge — Complete Tutorial

> **Tutorial Lengkap dari Awal Sampai Mahir**  
> 🇬🇧 English | [🇮🇩 Bahasa Indonesia](#-chatforge--tutorial-lengkap)

---

## 📋 Table of Contents

1. [Server Setup & Installation](#1-server-setup--installation)
2. [Login & First Access](#2-login--first-access)
3. [Connect WhatsApp](#3-connect-whatsapp)
4. [Send Your First Message](#4-send-your-first-message)
5. [Manage Contacts](#5-manage-contacts)
6. [Auto Reply Setup](#6-auto-reply-setup)
7. [AI Auto Reply](#7-ai-auto-reply)
8. [Broadcast Campaign](#8-broadcast-campaign)
9. [Drip Campaign](#9-drip-campaign)
10. [Scheduled Messages](#10-scheduled-messages)
11. [Live Chat Inbox](#11-live-chat-inbox)
12. [Message Templates](#12-message-templates)
13. [Recurring Campaigns](#13-recurring-campaigns)
14. [A/B Testing](#14-ab-testing)
15. [Canned Responses](#15-canned-responses)
16. [Macros & Workflows](#16-macros--workflows)
17. [Contact Tags & Groups](#17-contact-tags--groups)
18. [Blacklist & Spam Protection](#18-blacklist--spam-protection)
19. [CSAT Surveys](#19-csat-surveys)
20. [Store & Orders](#20-store--orders)
21. [Forms & Reminders](#21-forms--reminders)
22. [Web Widget](#22-web-widget)
23. [Email → WA Gateway](#23-email--wa-gateway)
24. [Link Tracker](#24-link-tracker)
25. [API Keys & Webhooks](#25-api-keys--webhooks)
26. [File Manager](#26-file-manager)
27. [WhatsApp Cloud API (Meta)](#27-whatsapp-cloud-api-meta)
28. [Admin: User Management](#28-admin-user-management)
29. [Admin: Roles & Permissions](#29-admin-roles--permissions)
30. [Admin: Packages & Limits](#30-admin-packages--limits)
31. [Admin: Subscriptions](#31-admin-subscriptions)
32. [Admin: Payment Gateways](#32-admin-payment-gateways)
33. [Admin: System & Backup](#33-admin-system--backup)
34. [Settings & Configuration](#34-settings--configuration)
35. [REST API Usage](#35-rest-api-usage)
36. [Deploy to Production](#36-deploy-to-production)

---

## 1. Server Setup & Installation

### Prerequisites
- **Go 1.21+** installed
- **MySQL 5.7+ or 8.0+** running
- **Git** installed

### Step 1.1 — Clone Repository

```bash
git clone https://github.com/linducip2208/chatforge.git
cd chatforge
```

![Clone Repository](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/01-login.png)

### Step 1.2 — Create Database

```bash
mysql -u root -p
```

```sql
CREATE DATABASE chatgo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### Step 1.3 — Configure Environment

Create `.env` file in project root:

```env
CHATGO_MYSQL=root:yourpassword@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4&parseTime=true
CHATGO_ENC_KEY=your-32-byte-random-key-here-change-me
APP_URL=http://localhost:8080
APP_NAME=ChatGo
APP_EMAIL=admin@yourdomain.com
```

> **Generate AES Key:** `openssl rand -base64 32`

### Step 1.4 — Build & Run

```bash
# Build
go build -o chatgo.exe .

# Run
./chatgo.exe
```

Output: `ChatGo running at http://0.0.0.0:8080`

### Step 1.5 — Run Database Migration

The app auto-migrates on startup. Alternatively:

```bash
mysql -u root -p chatgo < migrate.sql
```

✅ **Done!** Open `http://127.0.0.1:8080`

---

## 2. Login & First Access

### Default Credentials

| Role | Email | Password |
|------|-------|----------|
| Admin | `admin@chatgo.test` | `password` |

![Login Page](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/01-login.png)

### Dashboard Overview

After login, you'll see the Dashboard with:

- **Chart**: 7-day sent vs received messages
- **Stats**: Total users, active WA accounts, running campaigns
- **Quick access**: Connected WhatsApp numbers

![Dashboard](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/02-dashboard.png)

### Sidebar Navigation

The sidebar is organized into collapsible groups:

- **Dashboard** — overview & stats
- **Inbox** — live chat conversations
- **Contacts** — contact management
- **Broadcast & Campaigns** — mass messaging tools
- **Channels** — WhatsApp, Meta API, Email
- **Message Logs** — outgoing & incoming history
- **Automation** — auto reply, macros, forms
- **Commerce** — store & orders
- **Settings** — general configuration
- **Admin** (Admin only) — users, roles, packages, system

---

## 3. Connect WhatsApp

### Step 3.1 — Open Devices Page

Navigate to **Channels → WhatsApp (Unofficial) → Account & QR**

### Step 3.2 — Add Account

Click **Add Account** button. A QR code appears.

![WA QR](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/03-wa-qr.png)

### Step 3.3 — Scan with Phone

1. Open WhatsApp on your phone
2. Go to **Settings → Linked Devices**
3. Tap **Link a Device**
4. Scan the QR code on screen
5. Wait for "Connected" status

### Step 3.4 — Verify Connection

The status changes to **Connected** with your phone number displayed:

```
+628123456789  Connected
```

### Add Multiple Numbers

Repeat steps to add more WhatsApp numbers. Each number shows in the Accounts list.

**Notes:**
- Free tier: 1 WhatsApp account
- Growth tier: 5 WhatsApp accounts
- Enterprise: Unlimited

---

## 4. Send Your First Message

### Step 4.1 — Open Send Page

Navigate to **Message Logs → Send Message**

### Step 4.2 — Fill Form

![Send Message](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/04-send-message.png)

| Field | Description |
|-------|-------------|
| **Kirim Dari** | Select your WhatsApp number |
| **Nomor WA** | Recipient phone (e.g., `628123456789`) |
| **Pesan** | Your message text |

### Step 4.3 — Send

Click **Send** button. Success message appears: `Send OK`

### Send Media

To send images, videos, or documents:

1. Click the media attachment button
2. Select file type (Image/Video/Document)
3. Upload file
4. Add optional caption
5. Click Send

### Phone Sync

Messages sent directly from your phone also appear in the system with a **blue bubble** indicator 📱.

---

## 5. Manage Contacts

### Step 5.1 — Open Contacts

Navigate to **Contacts → All Contacts**

### Step 5.2 — Add Contact

Click the add form:

1. Enter **Name**
2. Enter **Phone Number** (e.g., `628123456789`)
3. Select **Groups** (optional)
4. Click **Add**

### Step 5.3 — Import CSV

![Contacts](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/09-contacts.png)

1. Click **Import CSV**
2. Upload CSV file with columns: `name, phone, groups`
3. System auto-creates new groups from CSV
4. Duplicates auto-skipped

**CSV Format:**
```csv
name,phone,groups
John Doe,628123456789,Customers
Jane Smith,628987654321,VIP
```

### Step 5.4 — Create Groups

Navigate to **Contacts → Groups**

1. Enter group name
2. Click **Add**

Groups help organize contacts for broadcast targeting.

### Step 5.5 — Tags

Navigate to **Contacts → Tags**

Create tags (e.g., "VIP", "New Lead", "Follow Up") and assign them to contacts.

### Step 5.6 — Merge Duplicates

Navigate to **Contacts → Merge Duplicates**

Find and merge contacts with the same phone number.

### Step 5.7 — Unsubscribed

View contacts who sent "STOP" or "Unsub" — they won't receive broadcasts.

---

## 6. Auto Reply Setup

### Step 6.1 — Open Auto Reply

Navigate to **Automation → Auto Reply**

### Step 6.2 — Create Rule

![Auto Reply](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/08-autoreply-rules.png)

| Setting | Description |
|---------|-------------|
| **Match Type** | Contains / Exact / Starts With / AI |
| **Keyword** | Trigger word (comma-separated for multiple) |
| **Reply** | Response text. Supports spintax: `{Halo|Hai}` |
| **AI** | Check to use AI for response |
| **AI Key** | Select AI provider key |
| **Nomor WA** | Select which numbers this rule applies to |

### Step 6.3 — Examples

**Basic Auto Reply:**
- Match: `contains` → Keyword: `menu` → Reply: `Berikut menu kami: ...`

**Spintax Reply:**
- Reply: `{Halo|Hai|Selamat datang}! Ada yang bisa dibantu?`

**Multiple Keywords:**
- Keyword: `harga, price, biaya` → matches any of these

### Step 6.4 — Toggle & Delete

- Click **ON/OFF** to enable/disable a rule
- Click **×** to delete

**Pro tip:** Rules are matched in order. Put specific rules before general ones.

---

## 7. AI Auto Reply

### Step 7.1 — Setup AI Key

Navigate to **Settings** (or Admin → AI for admins)

1. Click **AI Keys** tab
2. Click **Add**
3. Fill in:

| Field | Example |
|-------|---------|
| Name | `My OpenAI Key` |
| Provider | `openai` |
| Model | `gpt-4o` |
| API Key | `sk-...` (your key) |
| Base URL | `https://api.openai.com/v1` (or custom endpoint) |
| System Prompt | `You are a helpful customer support agent for ChatGo.` |

**Supported Providers:**
- OpenAI (`gpt-4o`, `gpt-4o-mini`, `gpt-3.5-turbo`)
- DeepSeek (`deepseek-chat`, `deepseek-reasoner`)
- Gemini (`gemini-2.0-flash`)
- Claude (`claude-3-5-sonnet`)
- Any OpenAI-compatible API (Ollama, vLLM, Groq, Mistral, etc.)

### Step 7.2 — Create AI Rule

![AI Auto Reply](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/07-autoreply-ai.png)

In Auto Reply, create a rule:
- Match Type: **AI**
- Check **Use AI**
- Select your AI Key
- Optional: add FAQ for grounding

### Step 7.3 — AI Settings

In **Settings → General:**

| Setting | Description |
|---------|-------------|
| **AI All Enabled** | AI responds to ALL incoming messages |
| **AI Fallback Only** | AI runs only when no keyword matches |
| **Memory Window** | Number of previous messages sent as context |
| **Delay Seconds** | Wait before AI response (typing effect) |
| **Reasoning Level** | Low / Medium / High |
| **Business Hours** | Restrict AI to working hours only |
| **Handoff Keywords** | Words that stop AI (e.g., `admin, operator`) |
| **Force Own Key** | Reseller requires sub-users to use their own key |

---

## 8. Broadcast Campaign

### Step 8.1 — Open Broadcast

Navigate to **Broadcast & Campaigns → Broadcast**

### Step 8.2 — Create Campaign

![Broadcast](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/05-broadcast.png)

| Field | Description |
|-------|-------------|
| **Name** | Campaign identifier |
| **Message** | Broadcast text with spintax |
| **Groups** | Select target contact groups |
| **Numbers** | Paste phone numbers (one per line) |
| **Account** | Select sender WA numbers |
| **Send Mode** | Round Robin or Random |
| **Interval** | Seconds between each message (default: 300) |
| **Tags** | Filter contacts by tag |
| **Media** | Optional image/video/document |

### Step 8.3 — Monitor Campaign

Campaign statuses:
- **Pending** — waiting to start
- **Running** — currently sending
- **Paused** — temporarily stopped
- **Done** — completed
- **Stopped** — manually stopped

Actions:
- **Pause/Resume** — toggle
- **Retry** — clone as new campaign
- **Stop** — cancel permanently

### Step 8.4 — Broadcast Tips

- Use Round Robin to distribute across multiple WA numbers
- Set interval 300-600 seconds to avoid rate limiting
- Preview message count before sending
- Test with small group first

---

## 9. Drip Campaign

### Step 9.1 — Open Drip Campaign

Navigate to **Broadcast & Campaigns → Drip Campaign**

### Step 9.2 — Create Drip

1. Click **Add Drip**
2. Enter campaign name
3. Campaign auto-starts in "active" status

### Step 9.3 — Add Steps

For each drip campaign, add steps:

| Field | Description |
|-------|-------------|
| **Drip** | Select campaign |
| **Delay (min)** | Minutes to wait before sending |
| **Message** | Text to send at this step |
| **Sort Order** | Sequence number |

**Example Sequence:**
1. Day 0: `Halo! Terima kasih sudah menghubungi kami 🙏`
2. Day 1: `Ada yang bisa kami bantu?`
3. Day 3: `Jangan lupa cek katalog produk kami di...`
4. Day 7: `Butuh bantuan? Tim kami siap 24/7`

### Step 9.4 — Auto Enrollment

Users automatically enroll in active drips when they send a message. To stop receiving drips, user sends `STOP` or `BERHENTI`.

---

## 10. Scheduled Messages

### Step 10.1 — Open Scheduled

Navigate to **Broadcast & Campaigns → Scheduled**

### Step 10.2 — Schedule Message

![Scheduled](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/06-scheduled.png)

| Field | Description |
|-------|-------------|
| **Name** | Label |
| **Phone** | Recipient number |
| **Message** | Text |
| **Send At** | Date & time (e.g., `2026-07-20T15:00`) |
| **Repeat** | Minutes to repeat (0 = once) |
| **Account** | Sender WA number |

**Example:**
- Daily reminder at 9 AM: Repeat = `1440` (24 hours × 60 min)

---

## 11. Live Chat Inbox

### Step 11.1 — Open Inbox

Navigate to **Inbox → All Conversations**

![Inbox](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/12-inbox.png)

### Step 11.2 — Chat Features

| Feature | Description |
|---------|-------------|
| **Chat View** | WhatsApp-style bubbles (green = sent, white = received, blue = phone-synced) |
| **Send Reply** | Type and send directly |
| **Assign Agent** | Route to team member |
| **Notes** | Internal notes (not visible to contact) |
| **Labels** | Tag conversations |
| **Close** | End conversation + send CSAT |
| **Transfer** | Move to another department |
| **Canned** | Quick reply shortcuts |
| **Macros** | Multi-action workflows |

### Step 11.3 — Real-time Updates

Inbox auto-updates via SSE (Server-Sent Events). New messages appear within seconds without page refresh.

### Step 11.4 — Meta Inbox

Messages from WhatsApp Cloud API also appear in the same Inbox view, unified with WhatsApp personal messages.

---

## 12. Message Templates

### Step 12.1 — Open Templates

Navigate to **Settings → Templates**

### Step 12.2 — Create Template

| Field | Description |
|-------|-------------|
| **Name** | Template identifier |
| **Content** | Message with variables: `{name}`, `{phone}`, `{date}` |

**Example Template:**
```
Halo {name},

Terima kasih sudah menghubungi ChatGo.
Pesanan Anda #{order_id} sedang diproses.
Estimasi pengiriman: {date}

Salam,
Tim ChatGo
```

### Step 12.3 — Use Templates

In Send Message, Broadcast, or Auto Reply, use templates with variables auto-filled from contact data.

---

## 13. Recurring Campaigns

Navigate to **Broadcast & Campaigns → Recurring**

Create automated campaigns that repeat on schedule:

| Setting | Description |
|---------|-------------|
| **Name** | Campaign name |
| **Message** | Broadcast text |
| **Groups** | Target groups |
| **Schedule** | Repeat interval (daily/weekly/monthly) |

---

## 14. A/B Testing

Navigate to **Broadcast & Campaigns → A/B Test**

1. Select a campaign
2. Enter **Variant A** and **Variant B** messages
3. System sends both variants to split audience
4. Compare **reply rates** to determine winner

---

## 15. Canned Responses

Navigate to **Automation → Canned Responses**

Create quick-reply shortcuts for inbox chat:

| Field | Description |
|-------|-------------|
| **Shortcut** | Quick key (e.g., `greet`) |
| **Name** | Display name |
| **Message** | Full response text |

In chat, type the shortcut to instantly insert the full response.

---

## 16. Macros & Workflows

Navigate to **Automation → Macros**

Create one-click workflows that execute multiple actions:

**Example Macro:**
```
assign:agent_id:3;tag:Resolved;reply:Terima kasih!;close;
```

Available actions: `assign`, `tag`, `reply`, `close`

---

## 17. Contact Tags & Groups

### Tags

Navigate to **Contacts → Tags**

1. Create tags (e.g., `VIP`, `New Lead`)
2. Assign tags in contact detail page
3. Filter broadcasts by tag

### Groups

Navigate to **Contacts → Groups**

1. Create group names
2. Assign contacts to groups
3. Use groups for broadcast targeting

---

## 18. Blacklist & Spam Protection

Navigate to **Settings → Blacklist**

- **Manual Block**: Add phone number + reason
- **Auto Detection**: System tracks spam patterns and auto-blocks
- **Remove**: Unblock numbers anytime

Blacklisted numbers won't receive broadcasts or trigger auto-replies.

---

## 19. CSAT Surveys

Navigate to **Reports → CSAT**

After closing a conversation, system sends:
```
Terima kasih! Bagaimana pengalaman Anda? Balas dengan rating 1-5 ⭐
```

View CSAT dashboard with average rating and total responses.

---

## 20. Store & Orders

### Products

Navigate to **Commerce → Store**

1. Add product categories
2. Add products with name, description, price
3. Products appear in WhatsApp chat when user types `menu`

### Orders

Navigate to **Commerce → Orders**

Track WhatsApp orders:
- New orders auto-created when customers order via chat
- Update status: Pending → Processing → Shipped → Delivered
- Send status notification to customer via WA

---

## 21. Forms & Reminders

### Forms

Navigate to **Automation → Forms**

Create interactive forms that customers can fill via WhatsApp chat conversation.

### Reminders

Navigate to **Automation → Reminders**

Schedule payment reminders:
| Field | Description |
|-------|-------------|
| **Phone** | Customer number |
| **Name** | Customer name |
| **Amount** | Payment amount |
| **Due Date** | Deadline |
| **Message** | Custom reminder text |

---

## 22. Web Widget

Navigate to **Automation → Widget**

Embed a WhatsApp chat button on any website:

```html
<script src="http://your-server.com/widget.js"></script>
```

The widget shows a floating WhatsApp button with inline chat.

---

## 23. Email → WA Gateway

Navigate to **Channels → Email → WA**

Configure email-to-WhatsApp forwarding. Emails sent to your configured address are automatically forwarded as WhatsApp messages.

---

## 24. Link Tracker

Navigate to **Settings → Links**

Create shortened URLs with click tracking:

1. Enter destination URL
2. Get short link (e.g., `http://your-server.com/track/abc123`)
3. Track clicks per campaign
4. View analytics dashboard

---

## 25. API Keys & Webhooks

### API Keys

Navigate to **Settings → API Keys**

1. Click **Add** with a name
2. **Copy the key immediately** — it won't be shown again
3. Use in API requests via `X-API-Key` header or `?apikey=` query param

### Webhooks

Navigate to **Settings → Webhooks**

1. Add webhook URL
2. Select events: `sent`, `received`, or `all`
3. System POSTs JSON on each matching event

**Webhook Payload:**
```json
{
  "event": "received",
  "phone": "628123456789",
  "name": "John Doe",
  "message": "Hello",
  "timestamp": "2026-07-17 15:30:00"
}
```

---

## 26. File Manager

Navigate to **Settings → Files**

Upload, browse, and manage media files:
- Images (JPG, PNG, GIF, WebP)
- Videos (MP4, MOV)
- Documents (PDF, DOC, XLS)

Files can be used in broadcasts and auto replies.

---

## 27. WhatsApp Cloud API (Meta)

### Step 27.1 — Setup Meta Account

Navigate to **Admin → Meta API** (or Channels → WhatsApp Cloud API)

1. Click **Add**
2. Enter:
   - **Name**: Account label
   - **Phone Number ID**: From Meta Developer Console
   - **Access Token**: Permanent token from Meta
   - **App ID**: From Meta App Dashboard
   - **App Secret**: From Meta App Dashboard
   - **Verify Token**: Custom string for webhook verification

### Step 27.2 — Configure Webhook

Navigate to **Channels → WhatsApp Cloud API → Webhook**

1. Copy the webhook URL
2. In Meta Developer Console, set Callback URL to `https://your-server.com/meta/webhook`
3. Enter your verify token
4. Subscribe to `messages` webhook field

### Step 27.3 — Message Templates

Navigate to **Channels → WhatsApp Cloud API → Templates**

Create and manage Meta-approved message templates for outbound messaging.

---

## 28. Admin: User Management

Navigate to **Admin → Business → Users**

![Admin Users](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/14-admin-users.png)

### Add User

| Field | Description |
|-------|-------------|
| **Name** | Full name |
| **Email** | Login email |
| **Password** | Login password |
| **Role** | Select from roles list |

### Impersonate

Click the impersonate icon to log in as any user. Click "Exit Impersonation" in topbar to return.

### Edit User

Click edit icon to modify name, email, role, or password.

---

## 29. Admin: Roles & Permissions

Navigate to **Admin → Business → Roles**

### Create Role

1. Enter **Name** (e.g., "Manager", "Agent")
2. Select **Permissions** (multi-select):
   - `manage_users` — access admin panel
   - `wa_send` — send messages
   - `wa_broadcast` — create campaigns
   - `wa_inbox` — access live chat
   - `wa_autoreply` — manage auto replies
   - ... and 50+ more permissions

### Assign Role

In Users page, select the role when creating/editing users. Permissions are checked on every request.

---

## 30. Admin: Packages & Limits

Navigate to **Admin → Business → Packages**

![Admin Packages](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/15-admin-packages.png)

### Create Package

| Field | Description |
|-------|-------------|
| **Name** | Package name |
| **Price** | Price amount |
| **Services** | Multi-select features included |
| **Limits** | Numerical limits per feature |

### Available Limits

| Limit | Controls |
|-------|----------|
| Send | Max sent messages |
| Device | Max WA accounts |
| WA Account | Max WA numbers |
| Contact | Max contacts |
| Meta | Max Meta accounts |
| Drips | Max drip campaigns |
| Recurring | Max recurring campaigns |
| Forms | Max forms |
| Templates | Max message templates |
| Canned | Max canned responses |
| Macros | Max macros |
| AI Key | Max AI provider keys |
| Knowledge | Max knowledge base entries |

---

## 31. Admin: Subscriptions

Navigate to **Admin → Business → Subscriptions**

### Assign Subscription

1. Select **User** (by email)
2. Select **Package**
3. Set **Expiry Date**
4. Click **Add**

System auto-enforces package limits based on active subscription.

---

## 32. Admin: Payment Gateways

Navigate to **Admin → Finance → Payment Gateways**

### Configure Gateway

| Field | Description |
|-------|-------------|
| **Name** | Gateway label |
| **Provider** | `midtrans`, `xendit`, `tripay`, `duitku` |
| **API Key** | Provider API key (encrypted at rest) |
| **API Secret** | Provider secret (encrypted at rest) |
| **Webhook Secret** | Callback verification key |
| **Base URL** | API endpoint URL |
| **Currency** | `IDR` (Indonesian Rupiah) |

### User Self-Service

Users can subscribe via the **Upgrade** page. Payment flow:
1. Select package
2. Choose gateway
3. Redirect to payment page
4. Auto-activate on payment confirmation

---

## 33. Admin: System & Backup

### Backup

Navigate to **Admin → System → Backup**

Download database backup as SQL file.

### Audit Log

Navigate to **Admin → System → Audit**

Track all admin actions: user creation, role changes, package updates, etc.

### WA Servers

Navigate to **Admin → Infrastructure → WA Servers**

Manage external WA server configurations with package restrictions.

---

## 34. Settings & Configuration

Navigate to **Settings → General**

![Settings](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/11-settings.png)

### Key Settings

| Setting | Description |
|---------|-------------|
| **Welcome Message** | Auto-send on first contact |
| **Fallback Message** | Send when no rule matches |
| **Reply in Group** | Enable/disable group chat replies |
| **Agent Signature** | Footer added to all agent messages |
| **Rate Limit** | Max daily messages + random delay |
| **Auto Close** | Auto-close inactive conversations |
| **App Name/Logo** | Branding customization |
| **Registrations** | Enable/disable public signup |

---

## 35. REST API Usage

### Authentication

All API requests require `X-API-Key` header or `?apikey=` query parameter.

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/send` | Send message |
| GET | `/api/status` | WA connection status |
| GET | `/api/messages` | Message history |
| GET | `/api/contacts` | Contact list |
| GET | `/api/devices` | Connected devices |

### Send Message

```bash
curl -X POST http://localhost:8080/api/send \
  -H "X-API-Key: YOUR_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "628123456789",
    "message": "Hello from API!",
    "account_phone": "+628123456789"
  }'
```

### Response

```json
{"status": "ok", "message": "sent"}
```

### Get Status

```bash
curl http://localhost:8080/api/status -H "X-API-Key: YOUR_KEY"
```

```json
{
  "status": "connected",
  "phone": "628123456789",
  "accounts": [
    {"id": "628123456789@s.whatsapp.net", "phone": "628123456789", "status": "connected"}
  ]
}
```

---

## 36. Deploy to Production

### Step 36.1 — Build Linux Binary

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o chatgo_linux .
```

### Step 36.2 — Upload to Server

```bash
scp chatgo_linux .env user@your-server:/opt/chatgo/
```

### Step 36.3 — Systemd Service

```bash
sudo cp chatgo.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable chatgo
sudo systemctl start chatgo
```

### Step 36.4 — Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name chatgo.yourdomain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Step 36.5 — SSL with Certbot

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d chatgo.yourdomain.com
```

### Step 36.6 — Verify

Open `https://chatgo.yourdomain.com` and login.

---

## 🎉 Congratulations!

You've completed the full ChatGo setup from installation to production deployment. For help:

- **WhatsApp Support**: [+62 812-9605-2010](https://wa.me/6281296052010)
- **GitHub Issues**: [github.com/linducip2208/chatforge/issues](https://github.com/linducip2208/chatforge/issues)

---

# 🏗️ Chatforge — Tutorial Lengkap

> 🇬🇧 [English](#-chatforge--complete-tutorial) | **🇮🇩 Bahasa Indonesia**

---

## 📋 Daftar Isi

1. [Setup Server & Instalasi](#1-setup-server--instalasi)
2. [Login & Akses Pertama](#2-login--akses-pertama)
3. [Hubungkan WhatsApp](#3-hubungkan-whatsapp)
4. [Kirim Pesan Pertama](#4-kirim-pesan-pertama)
5. [Kelola Kontak](#5-kelola-kontak)
6. [Setup Auto Reply](#6-setup-auto-reply)
7. [AI Auto Reply](#7-ai-auto-reply)
8. [Broadcast Campaign](#8-broadcast-campaign)
9. [Drip Campaign](#9-drip-campaign)
10. [Jadwal Pesan](#10-jadwal-pesan)
11. [Live Chat Inbox](#11-live-chat-inbox)
12. [Template Pesan](#12-template-pesan)
13. [Recurring Campaigns](#13-recurring-campaigns)
14. [A/B Testing](#14-ab-testing)
15. [Canned Responses](#15-canned-responses)
16. [Macros & Workflows](#16-macros--workflows)
17. [Tag & Grup Kontak](#17-tag--grup-kontak)
18. [Blacklist & Anti-Spam](#18-blacklist--anti-spam)
19. [CSAT Survey](#19-csat-survey)
20. [Toko & Pesanan](#20-toko--pesanan)
21. [Form & Pengingat](#21-form--pengingat)
22. [Web Widget](#22-web-widget)
23. [Email → WA Gateway](#23-email--wa-gateway)
24. [Link Tracker](#24-link-tracker)
25. [API Keys & Webhooks](#25-api-keys--webhooks)
26. [File Manager](#26-file-manager)
27. [WhatsApp Cloud API (Meta)](#27-whatsapp-cloud-api-meta)
28. [Admin: Manajemen User](#28-admin-manajemen-user)
29. [Admin: Role & Permission](#29-admin-role--permission)
30. [Admin: Paket & Limit](#30-admin-paket--limit)
31. [Admin: Subscription](#31-admin-subscription)
32. [Admin: Payment Gateway](#32-admin-payment-gateway)
33. [Admin: Sistem & Backup](#33-admin-sistem--backup)
34. [Pengaturan & Konfigurasi](#34-pengaturan--konfigurasi)
35. [REST API](#35-rest-api)
36. [Deploy ke Production](#36-deploy-ke-production)

---

## 1. Setup Server & Instalasi

### Persiapan
- **Go 1.21+** terinstall
- **MySQL 5.7+ / 8.0+** berjalan
- **Git** terinstall

### Langkah 1.1 — Clone Repository

```bash
git clone https://github.com/linducip2208/chatforge.git
cd chatforge
```

![Clone Repo](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/01-login.png)

### Langkah 1.2 — Buat Database

```bash
mysql -u root -p
```

```sql
CREATE DATABASE chatgo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### Langkah 1.3 — Konfigurasi Environment

Buat file `.env` di root project:

```env
CHATGO_MYSQL=root:passwordanda@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4&parseTime=true
CHATGO_ENC_KEY=kunci-aes-32-byte-acak-anda-ganti-ini
APP_URL=http://localhost:8080
APP_NAME=ChatGo
APP_EMAIL=admin@domainanda.com
```

> **Generate AES Key:** `openssl rand -base64 32`

### Langkah 1.4 — Build & Jalankan

```bash
# Build
go build -o chatgo.exe .

# Jalankan
./chatgo.exe
```

Output: `ChatGo running at http://0.0.0.0:8080`

### Langkah 1.5 — Jalankan Migrasi Database

Aplikasi auto-migrasi saat startup. Atau manual:

```bash
mysql -u root -p chatgo < migrate.sql
```

✅ **Selesai!** Buka `http://127.0.0.1:8080`

---

## 2. Login & Akses Pertama

### Kredensial Default

| Role | Email | Password |
|------|-------|----------|
| Admin | `admin@chatgo.test` | `password` |

![Halaman Login](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/01-login.png)

### Tampilan Dashboard

Setelah login, Anda akan melihat Dashboard dengan:

- **Chart**: Grafik pesan terkirim vs diterima 7 hari
- **Stats**: Total user, akun WA aktif, campaign berjalan
- **Quick access**: Nomor WhatsApp yang terkoneksi

![Dashboard](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/02-dashboard.png)

### Navigasi Sidebar

Sidebar diorganisir dalam grup collapsible:

- **Dashboard** — overview & statistik
- **Inbox** — percakapan live chat
- **Contacts** — manajemen kontak
- **Broadcast & Campaigns** — alat kirim massal
- **Channels** — WhatsApp, Meta API, Email
- **Message Logs** — riwayat keluar & masuk
- **Automation** — auto reply, macros, form
- **Commerce** — toko & pesanan
- **Settings** — konfigurasi umum
- **Admin** (khusus Admin) — user, role, paket, sistem

---

## 3. Hubungkan WhatsApp

### Langkah 3.1 — Buka Halaman Devices

Navigasi ke **Channels → WhatsApp (Unofficial) → Account & QR**

### Langkah 3.2 — Tambah Akun

Klik tombol **Add Account**. QR code akan muncul.

![WA QR](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/03-wa-qr.png)

### Langkah 3.3 — Scan dengan HP

1. Buka WhatsApp di HP Anda
2. Masuk ke **Setelan → Perangkat Tertaut**
3. Ketuk **Tautkan Perangkat**
4. Scan QR code di layar
5. Tunggu status "Connected"

### Langkah 3.4 — Verifikasi Koneksi

Status berubah jadi **Connected** dengan nomor HP:

```
+628123456789  Connected
```

### Tambah Banyak Nomor

Ulangi langkah untuk menambah nomor WhatsApp lain. Setiap nomor muncul di daftar Akun.

**Catatan:**
- Paket Free: 1 akun WhatsApp
- Paket Growth: 5 akun WhatsApp
- Paket Enterprise: Unlimited

---

## 4. Kirim Pesan Pertama

### Langkah 4.1 — Buka Halaman Kirim

Navigasi ke **Message Logs → Send Message**

### Langkah 4.2 — Isi Form

![Kirim Pesan](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/04-send-message.png)

| Field | Deskripsi |
|-------|-----------|
| **Kirim Dari** | Pilih nomor WhatsApp Anda |
| **Nomor WA** | Nomor tujuan (contoh: `628123456789`) |
| **Pesan** | Teks pesan Anda |

### Langkah 4.3 — Kirim

Klik tombol **Send**. Pesan sukses muncul: `Send OK`

### Kirim Media

Untuk mengirim gambar, video, atau dokumen:

1. Klik tombol attachment
2. Pilih tipe file (Gambar/Video/Dokumen)
3. Upload file
4. Tambah caption (opsional)
5. Klik Send

### Phone Sync

Pesan yang dikirim langsung dari HP Anda juga muncul di sistem dengan indikator **bubble biru** 📱.

---

## 5. Kelola Kontak

### Langkah 5.1 — Buka Kontak

Navigasi ke **Contacts → All Contacts**

### Langkah 5.2 — Tambah Kontak

Klik form tambah:

1. Masukkan **Nama**
2. Masukkan **Nomor HP** (contoh: `628123456789`)
3. Pilih **Grup** (opsional)
4. Klik **Add**

### Langkah 5.3 — Import CSV

![Kontak](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/09-contacts.png)

1. Klik **Import CSV**
2. Upload file CSV dengan kolom: `name, phone, groups`
3. Sistem otomatis buat grup baru dari CSV
4. Duplikat otomatis dilewati

**Format CSV:**
```csv
name,phone,groups
Budi Santoso,628123456789,Pelanggan
Siti Nurhaliza,628987654321,VIP
```

### Langkah 5.4 — Buat Grup

Navigasi ke **Contacts → Groups**

1. Masukkan nama grup
2. Klik **Add**

### Langkah 5.5 — Tag

Navigasi ke **Contacts → Tags**

Buat tag (contoh: "VIP", "Leads Baru", "Follow Up") dan pasang ke kontak.

### Langkah 5.6 — Gabung Duplikat

Navigasi ke **Contacts → Merge Duplicates**

Temukan dan gabungkan kontak dengan nomor yang sama.

---

## 6. Setup Auto Reply

### Langkah 6.1 — Buka Auto Reply

Navigasi ke **Automation → Auto Reply**

### Langkah 6.2 — Buat Rule

![Auto Reply](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/08-autoreply-rules.png)

| Setting | Deskripsi |
|---------|-----------|
| **Match Type** | Contains / Exact / Starts With / AI |
| **Keyword** | Kata pemicu (pisahkan koma untuk multi) |
| **Reply** | Teks balasan. Support spintax: `{Halo|Hai}` |
| **AI** | Centang untuk pakai AI |
| **AI Key** | Pilih provider AI |
| **Nomor WA** | Pilih nomor mana rule ini berlaku |

### Langkah 6.3 — Contoh

**Auto Reply Dasar:**
- Match: `contains` → Keyword: `menu` → Reply: `Berikut menu kami: ...`

**Spintax Reply:**
- Reply: `{Halo|Hai|Selamat datang}! Ada yang bisa dibantu?`

**Multi Keyword:**
- Keyword: `harga, price, biaya` → cocok salah satu

### Langkah 6.4 — Toggle & Hapus

- Klik **ON/OFF** untuk mengaktifkan/nonaktifkan rule
- Klik **×** untuk menghapus

**Pro tip:** Rule dicocokkan berurutan. Taruh rule spesifik sebelum rule umum.

---

## 7. AI Auto Reply

### Langkah 7.1 — Setup AI Key

Navigasi ke **Settings** (atau Admin → AI untuk admin)

1. Klik tab **AI Keys**
2. Klik **Add**
3. Isi:

| Field | Contoh |
|-------|--------|
| Name | `OpenAI Key Saya` |
| Provider | `openai` |
| Model | `gpt-4o` |
| API Key | `sk-...` (kunci Anda) |
| Base URL | `https://api.openai.com/v1` |
| System Prompt | `Anda adalah customer service ChatGo yang ramah.` |

**Provider yang Didukung:**
- OpenAI (`gpt-4o`, `gpt-4o-mini`)
- DeepSeek (`deepseek-chat`)
- Gemini (`gemini-2.0-flash`)
- Claude (`claude-3-5-sonnet`)
- Semua API OpenAI-compatible (Ollama, Groq, Mistral, dll.)

### Langkah 7.2 — Buat AI Rule

![AI Auto Reply](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/07-autoreply-ai.png)

Di Auto Reply, buat rule:
- Match Type: **AI**
- Centang **Use AI**
- Pilih AI Key Anda

### Langkah 7.3 — Pengaturan AI

Di **Settings → General:**

| Setting | Deskripsi |
|---------|-----------|
| **AI All Enabled** | AI balas SEMUA pesan masuk |
| **AI Fallback Only** | AI hanya jalan saat tidak ada keyword match |
| **Memory Window** | Jumlah pesan sebelumnya sebagai konteks |
| **Delay Seconds** | Jeda sebelum AI membalas |
| **Reasoning Level** | Low / Medium / High |
| **Business Hours** | Batasi AI hanya di jam kerja |
| **Handoff Keywords** | Kata pemicu stop AI (contoh: `admin, operator`) |

---

## 8. Broadcast Campaign

### Langkah 8.1 — Buka Broadcast

Navigasi ke **Broadcast & Campaigns → Broadcast**

### Langkah 8.2 — Buat Campaign

![Broadcast](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/05-broadcast.png)

| Field | Deskripsi |
|-------|-----------|
| **Name** | Nama campaign |
| **Message** | Teks broadcast dengan spintax |
| **Groups** | Pilih grup kontak target |
| **Numbers** | Tempel nomor HP (satu per baris) |
| **Account** | Pilih nomor WA pengirim |
| **Send Mode** | Round Robin atau Random |
| **Interval** | Detik antar pesan (default: 300) |
| **Tags** | Filter kontak by tag |
| **Media** | Gambar/video/dokumen opsional |

### Langkah 8.3 — Monitor Campaign

Status campaign:
- **Pending** — menunggu mulai
- **Running** — sedang mengirim
- **Paused** — dijeda sementara
- **Done** — selesai
- **Stopped** — dihentikan manual

Aksi:
- **Pause/Resume** — jeda/lanjutkan
- **Retry** — clone jadi campaign baru
- **Stop** — batalkan permanen

---

## 9. Drip Campaign

### Langkah 9.1 — Buka Drip Campaign

Navigasi ke **Broadcast & Campaigns → Drip Campaign**

### Langkah 9.2 — Buat Drip

1. Klik **Add Drip**
2. Masukkan nama campaign
3. Campaign otomatis aktif

### Langkah 9.3 — Tambah Step

| Field | Deskripsi |
|-------|-----------|
| **Drip** | Pilih campaign |
| **Delay (menit)** | Menit tunggu sebelum kirim |
| **Message** | Teks yang dikirim |
| **Sort Order** | Urutan step |

**Contoh Urutan:**
1. Hari 0: `Halo! Terima kasih sudah menghubungi kami 🙏`
2. Hari 1: `Ada yang bisa kami bantu?`
3. Hari 3: `Jangan lupa cek katalog produk kami di...`
4. Hari 7: `Butuh bantuan? Tim kami siap 24/7`

### Langkah 9.4 — Auto Enrollment

User otomatis masuk ke drip aktif saat kirim pesan. Untuk berhenti, user kirim `STOP` atau `BERHENTI`.

---

## 10. Jadwal Pesan

### Langkah 10.1 — Buka Scheduled

Navigasi ke **Broadcast & Campaigns → Scheduled**

### Langkah 10.2 — Jadwalkan Pesan

![Scheduled](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/06-scheduled.png)

| Field | Deskripsi |
|-------|-----------|
| **Name** | Label |
| **Phone** | Nomor tujuan |
| **Message** | Teks |
| **Send At** | Tanggal & jam (`2026-07-20T15:00`) |
| **Repeat** | Menit untuk ulang (0 = sekali) |
| **Account** | Nomor WA pengirim |

---

## 11. Live Chat Inbox

### Langkah 11.1 — Buka Inbox

Navigasi ke **Inbox → All Conversations**

![Inbox](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/12-inbox.png)

### Langkah 11.2 — Fitur Chat

| Fitur | Deskripsi |
|-------|-----------|
| **Chat View** | Bubble ala WhatsApp (hijau=terkirim, putih=diterima, biru=sync HP) |
| **Kirim Balasan** | Ketik dan kirim langsung |
| **Assign Agent** | Arahkan ke anggota tim |
| **Catatan** | Catatan internal (tidak terlihat kontak) |
| **Label** | Tag percakapan |
| **Close** | Akhiri percakapan + kirim CSAT |
| **Transfer** | Pindah ke departemen lain |
| **Canned** | Shortcut balasan cepat |
| **Macros** | Workflow multi-aksi |

### Langkah 11.3 — Real-time Update

Inbox auto-update via SSE. Pesan baru muncul dalam hitungan detik tanpa refresh.

---

## 12 - 27. Fitur Lanjutan

> Untuk langkah 12-27, lihat [English version](#12-message-templates) di atas. Prosedurnya sama persis — cukup ganti bahasa.

### Ringkasan Fitur Lanjutan:

| # | Fitur | Navigasi |
|---|-------|----------|
| 12 | Template Pesan | Settings → Templates |
| 13 | Recurring | Broadcast → Recurring |
| 14 | A/B Test | Broadcast → A/B Test |
| 15 | Canned Responses | Automation → Canned |
| 16 | Macros | Automation → Macros |
| 17 | Tag & Grup | Contacts → Tags / Groups |
| 18 | Blacklist | Settings → Blacklist |
| 19 | CSAT | Reports → CSAT |
| 20 | Toko & Pesanan | Commerce → Store / Orders |
| 21 | Form & Reminder | Automation → Forms / Reminders |
| 22 | Widget | Automation → Widget |
| 23 | Email→WA | Channels → Email→WA |
| 24 | Link Tracker | Settings → Links |
| 25 | API Keys | Settings → API Keys |
| 26 | File Manager | Settings → Files |
| 27 | Meta API | Channels → WhatsApp Cloud API |

---

## 28. Admin: Manajemen User

Navigasi ke **Admin → Business → Users**

![Admin Users](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/14-admin-users.png)

### Tambah User

| Field | Deskripsi |
|-------|-----------|
| **Name** | Nama lengkap |
| **Email** | Email login |
| **Password** | Password login |
| **Role** | Pilih dari daftar role |

### Impersonate

Klik ikon impersonate untuk login sebagai user lain. Klik "Exit Impersonation" di topbar untuk kembali.

---

## 29. Admin: Role & Permission

Navigasi ke **Admin → Business → Roles**

### Buat Role

1. Masukkan **Nama** (contoh: "Manager", "Agent")
2. Pilih **Permissions** (multi-select)

### Assign Role

Di halaman Users, pilih role saat membuat/edit user.

---

## 30. Admin: Paket & Limit

Navigasi ke **Admin → Business → Packages**

![Admin Packages](https://raw.githubusercontent.com/linducip2208/chatforge/main/public/marketing/screens/15-admin-packages.png)

### Buat Paket

| Field | Deskripsi |
|-------|-----------|
| **Name** | Nama paket |
| **Price** | Harga |
| **Services** | Fitur yang termasuk (multi-select) |
| **Limits** | Batasan numerik per fitur |

---

## 31-36. Admin Lanjutan & Deploy

> Untuk langkah 31-36, lihat [English version](#31-admin-subscriptions) di atas.

### Ringkasan:

| # | Topik | Navigasi |
|---|-------|----------|
| 31 | Subscription | Admin → Business → Subscriptions |
| 32 | Payment Gateway | Admin → Finance → Payment Gateways |
| 33 | Backup & Audit | Admin → System |
| 34 | Settings | Settings → General |
| 35 | REST API | Dokumentasi API endpoint |
| 36 | Deploy Production | Build Linux, systemd, Nginx, SSL |

---

## 🎉 Selamat!

Anda telah menyelesaikan tutorial lengkap ChatGo. Untuk bantuan:

- **WhatsApp Support**: [+62 812-9605-2010](https://wa.me/6281296052010)
- **GitHub Issues**: [github.com/linducip2208/chatforge/issues](https://github.com/linducip2208/chatforge/issues)
