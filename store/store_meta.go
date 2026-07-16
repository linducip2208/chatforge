package store

type MetaAccount struct {
	ID            int64
	Name          string
	PhoneNumberID string
	AccessToken   string
	AppID         string
	AppSecret     string
	VerifyToken   string
	UserID        int64
	ParentID      int64
	Created       string
}

type MetaTemplate struct {
	ID        int64
	Name      string
	Language  string
	Category  string
	Components string
	Status    string
	Created   string
}

func (d *DB) migrateMeta() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS meta_accounts (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			phone_number_id VARCHAR(64) NOT NULL DEFAULT '',
			access_token TEXT NOT NULL,
			app_id VARCHAR(64) NOT NULL DEFAULT '',
			app_secret VARCHAR(128) NOT NULL DEFAULT '',
			verify_token VARCHAR(64) NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS meta_templates (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			language VARCHAR(10) NOT NULL DEFAULT 'id',
			category VARCHAR(50) NOT NULL DEFAULT 'marketing',
			components TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) AddMetaAccount(name, phoneNumberID, accessToken, appID, appSecret, verifyToken string, userID, parentID int64) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO meta_accounts (name, phone_number_id, access_token, app_id, app_secret, verify_token, user_id, parent_id) VALUES (?,?,?,?,?,?,?,?)`, name, phoneNumberID, accessToken, appID, appSecret, verifyToken, userID, parentID)
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (d *DB) DeleteMetaAccount(id int64) error { _, err := d.sql.Exec(`DELETE FROM meta_accounts WHERE id=?`, id); return err }

func (d *DB) ListMetaAccounts() ([]MetaAccount, error) {
	rows, err := d.sql.Query(`SELECT id, name, phone_number_id, access_token, app_id, app_secret, verify_token, IFNULL(user_id,0), IFNULL(parent_id,0), created_at FROM meta_accounts ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []MetaAccount
	for rows.Next() {
		var m MetaAccount
		rows.Scan(&m.ID, &m.Name, &m.PhoneNumberID, &m.AccessToken, &m.AppID, &m.AppSecret, &m.VerifyToken, &m.UserID, &m.ParentID, &m.Created)
		out = append(out, m)
	}
	return out, nil
}

func (d *DB) GetMetaAccount(id int64) (*MetaAccount, error) {
	var m MetaAccount
	err := d.sql.QueryRow(`SELECT id, name, phone_number_id, access_token, app_id, app_secret, verify_token, IFNULL(user_id,0), IFNULL(parent_id,0), created_at FROM meta_accounts WHERE id=?`, id).Scan(&m.ID, &m.Name, &m.PhoneNumberID, &m.AccessToken, &m.AppID, &m.AppSecret, &m.VerifyToken, &m.UserID, &m.ParentID, &m.Created)
	if err != nil { return nil, err }
	return &m, nil
}

func (d *DB) AddMetaTemplate(name, language, category, components, status string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO meta_templates (name, language, category, components, status) VALUES (?,?,?,?,?)`, name, language, category, components, status)
	if err != nil { return 0, err }
	return res.LastInsertId()
}

func (d *DB) DeleteMetaTemplate(id int64) error { _, err := d.sql.Exec(`DELETE FROM meta_templates WHERE id=?`, id); return err }

func (d *DB) GetUserMetaLimit(userID int64) int {
	u, err := d.GetUserByID(userID)
	if err != nil { return 0 }
	var pkgName string
	err = d.sql.QueryRow(`SELECT pkg FROM subscriptions WHERE user=? ORDER BY id DESC LIMIT 1`, u.Email).Scan(&pkgName)
	if err != nil { return 0 }
	var limit int
	err = d.sql.QueryRow(`SELECT meta_limit FROM packages WHERE name=? LIMIT 1`, pkgName).Scan(&limit)
	if err != nil { return 0 }
	return limit
}

func (d *DB) CountMetaByUser(userID int64) int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM meta_accounts WHERE user_id=?`, userID).Scan(&n)
	return n
}

func (d *DB) ListMetaTemplates() ([]MetaTemplate, error) {
	rows, err := d.sql.Query(`SELECT id, name, language, category, components, status, created_at FROM meta_templates ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []MetaTemplate
	for rows.Next() {
		var m MetaTemplate
		rows.Scan(&m.ID, &m.Name, &m.Language, &m.Category, &m.Components, &m.Status, &m.Created)
		out = append(out, m)
	}
	return out, nil
}
