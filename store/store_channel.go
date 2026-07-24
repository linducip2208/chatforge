package store

import "database/sql"

type ChannelKey struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Platform  string `json:"platform"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	Token     string `json:"token"`
	TokenSecret string `json:"token_secret"`
	WebhookURL string `json:"webhook_url,omitempty"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (d *DB) migrateChannelKeys() error {
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS channel_keys (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		user_id BIGINT NOT NULL DEFAULT 0,
		platform VARCHAR(32) NOT NULL,
		name VARCHAR(128) NOT NULL DEFAULT '',
		api_key VARCHAR(512) NOT NULL DEFAULT '',
		api_secret VARCHAR(512) NOT NULL DEFAULT '',
		token VARCHAR(1024) NOT NULL DEFAULT '',
		token_secret VARCHAR(1024) NOT NULL DEFAULT '',
		webhook_url VARCHAR(512) NOT NULL DEFAULT '',
		is_active TINYINT NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_ck_platform (platform),
		INDEX idx_ck_user (user_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return nil
}

func (d *DB) AddChannelKey(uid int64, platform, name, apiKey, apiSecret, token, tokenSecret string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO channel_keys (user_id,platform,name,api_key,api_secret,token,token_secret,is_active) VALUES (?,?,?,?,?,?,?,1)`,
		uid, platform, name, apiKey, apiSecret, token, tokenSecret)
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (d *DB) ListChannelKeys(uid int64, platform string) ([]ChannelKey, error) {
	var rows *sql.Rows
	var err error
	if platform == "" {
		rows, err = d.sql.Query(`SELECT id,user_id,platform,name,api_key,api_secret,token,token_secret,webhook_url,is_active,created_at FROM channel_keys WHERE user_id=? ORDER BY platform,id`, uid)
	} else {
		rows, err = d.sql.Query(`SELECT id,user_id,platform,name,api_key,api_secret,token,token_secret,webhook_url,is_active,created_at FROM channel_keys WHERE user_id=? AND platform=? ORDER BY id`, uid, platform)
	}
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ChannelKey
	for rows.Next() {
		var ck ChannelKey; var active int
		rows.Scan(&ck.ID, &ck.UserID, &ck.Platform, &ck.Name, &ck.APIKey, &ck.APISecret, &ck.Token, &ck.TokenSecret, &ck.WebhookURL, &active, &ck.CreatedAt)
		ck.IsActive = active == 1
		out = append(out, ck)
	}
	return out, nil
}

func (d *DB) DeleteChannelKey(id, uid int64) error {
	_, err := d.sql.Exec(`DELETE FROM channel_keys WHERE id=? AND user_id=?`, id, uid)
	return err
}

func (d *DB) ToggleChannelKey(id, uid int64) error {
	_, err := d.sql.Exec(`UPDATE channel_keys SET is_active=1-is_active WHERE id=? AND user_id=?`, id, uid)
	return err
}
