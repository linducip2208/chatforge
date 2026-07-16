# ChatGo Production Deployment

## Struktur Folder
```
produksi/
├── chatgo.exe        # Windows binary
├── chatgo_linux      # Linux binary
├── .env              # Konfigurasi (edit sebelum jalan)
├── web/              # Web assets (CSS, JS, images, theme)
├── lang/             # Language files (id.json, en.json)
└── public/           # Marketing screenshots
```

## Quick Deploy (3 langkah)

### 1. Setup Database
Database auto-create saat pertama jalan — cukup bikin database kosong:
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS chatgo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
```

### 2. Seed Default Data (opsional)
```bash
mysql -u root -p chatgo < seed.sql
```
> Pakai `INSERT IGNORE` — aman di-run ulang, tidak overwrite data existing.

### 2. Edit .env
```bash
nano .env
```
Ganti `CHATGO_MYSQL`, `APP_URL`, `APP_NAME`, `APP_EMAIL` sesuai server.

### 3. Jalankan
```bash
# Linux
chmod +x chatgo_linux
./chatgo_linux

# Systemd (recommended)
sudo nano /etc/systemd/system/chatgo.service
```

```ini
[Unit]
Description=ChatGo Server
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/chatgo
ExecStart=/var/www/chatgo/chatgo_linux
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable chatgo
sudo systemctl start chatgo
```

## Nginx Reverse Proxy
```nginx
server {
    listen 80;
    server_name mydomain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 86400s;
    }
}
```

## Meta Webhook Setup
Buka Admin > Meta API, copy webhook URL: `https://mydomain.com/webhook/meta`

## Default Login
- **Admin**: admin@chatgo.test / password

## Update ke Versi Baru
```bash
sudo systemctl stop chatgo
scp chatgo_linux user@server:/var/www/chatgo/
scp -r web/* user@server:/var/www/chatgo/web/
scp -r lang/* user@server:/var/www/chatgo/lang/
sudo systemctl start chatgo
```
> JANGAN replace: `.env`, `data/` folder
