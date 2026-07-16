# ChatGo Production Deployment

## Struktur Folder
```
produksi/
├── chatgo.exe        # Windows binary
├── chatgo_linux      # Linux binary
├── .env              # Konfigurasi (edit)
├── seed.sql          # Default data (roles, packages)
├── web/              # Web assets
├── lang/             # Language files
└── public/           # Screenshots
```

## Deploy

### 1. Database
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS chatgo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
```

### 2. Seed (opsional)
```bash
mysql -u root -p chatgo < seed.sql
```

### 3. Edit .env
```bash
nano .env
```

### 4. Run
```bash
chmod +x chatgo_linux && ./chatgo_linux
```

## Systemd
```ini
[Unit]
Description=ChatGo
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

## Nginx
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
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 86400s;
    }
}
```

## Update
```bash
sudo systemctl stop chatgo
scp chatgo_linux user@server:/var/www/chatgo/
scp -r web/* user@server:/var/www/chatgo/web/
scp -r lang/* user@server:/var/www/chatgo/lang/
sudo systemctl start chatgo
```
> JANGAN replace `.env`, `data/`

## Login
admin@chatgo.test / password
