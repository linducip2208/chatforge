package store

import (
	"fmt"
	"strings"
	"time"
)

type Contact struct {
	ID      int64
	Name    string
	Phone   string
	Groups  string
	Created string
}
type Group struct {
	ID      int64
	Name    string
	Count   int
	Created string
}
type Template struct {
	ID      int64
	Name    string
	Content string
	Created string
}
type Unsub struct {
	ID      int64
	Phone   string
	Created string
}
type APIKey struct {
	ID      int64
	Name    string
	Secret  string
	Created string
}
type Webhook struct {
	ID      int64
	Name    string
	URL     string
	Event   string
	Created string
}
type Campaign struct {
	ID         int64
	Name       string
	Groups     string
	AccountID  string
	AccountIDs string
	Message    string
	Total      int
	Sent       int
	Status     string
	Interval   int
	SentTo     string
	Created    string
}
type Scheduled struct {
	ID         int64
	Name       string
	Phone      string
	Message    string
	SendAt     string
	Repeat     int
	Status     string
	AccountIDs string
	Created    string
}
type LogEntry struct {
	ID      int64
	Type    string
	Reason  string
	Content string
	Created string
}

func (d *DB) migrateExtra() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS contacts (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, phone VARCHAR(64) NOT NULL, ` + "`groups`" + ` VARCHAR(512) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS contact_groups (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS unsubscribed (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS templates (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, content TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS api_keys (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, secret VARCHAR(128) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS webhooks (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, url VARCHAR(512) NOT NULL, event VARCHAR(32) NOT NULL DEFAULT 'all', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS campaigns (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, ` + "`groups`" + ` VARCHAR(255) NOT NULL DEFAULT '', message TEXT NOT NULL, total INT NOT NULL DEFAULT 0, sent INT NOT NULL DEFAULT 0, status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS scheduled (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, phone VARCHAR(64) NOT NULL, message TEXT NOT NULL, send_at DATETIME NOT NULL, repeat_min INT NOT NULL DEFAULT 0, status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS logger (id BIGINT AUTO_INCREMENT PRIMARY KEY, type VARCHAR(40) NOT NULL, reason VARCHAR(255) NOT NULL DEFAULT '', content TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Contacts ----
func (d *DB) AddContact(name, phone, groups string) (int64, error) {
	res, err := d.sql.Exec("INSERT INTO contacts (name, phone, `groups`) VALUES (?, ?, ?)", name, phone, groups)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateContact(id int64, name, phone, groups string) error {
	_, err := d.sql.Exec("UPDATE contacts SET name=?, phone=?, `groups`=? WHERE id=?", name, phone, groups, id)
	return err
}
func (d *DB) GetContact(id int64) (*Contact, error) {
	var c Contact
	err := d.sql.QueryRow(`SELECT id, name, phone, `+"`groups`"+`, created_at FROM contacts WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Groups, &c.Created)
	if err != nil { return nil, err }
	return &c, nil
}
func (d *DB) DeleteContact(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM contacts WHERE id=?`, id)
	return err
}
func (d *DB) ListContacts() ([]Contact, error) {
	rows, err := d.sql.Query("SELECT id, name, phone, `groups`, created_at FROM contacts ORDER BY id DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Contact
	for rows.Next() {
		var c Contact
		rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Groups, &c.Created)
		out = append(out, c)
	}
	return out, nil
}
func (d *DB) ContactsByGroup(gid string) ([]Contact, error) {
	rows, err := d.sql.Query("SELECT id, name, phone, `groups`, created_at FROM contacts WHERE `groups` LIKE ? ORDER BY name", "%"+gid+"%")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Contact
	for rows.Next() {
		var c Contact
		rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Groups, &c.Created)
		out = append(out, c)
	}
	return out, nil
}

// ---- Groups ----
func (d *DB) AddGroup(name string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO contact_groups (name) VALUES (?)`, name)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteGroup(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM contact_groups WHERE id=?`, id)
	return err
}
func (d *DB) ListGroups() ([]Group, error) {
	rows, err := d.sql.Query(`SELECT g.id, g.name, COUNT(c.id), g.created_at FROM contact_groups g LEFT JOIN contacts c ON FIND_IN_SET(g.id, REPLACE(c.`+"`groups`"+`, ' ', '')) GROUP BY g.id ORDER BY g.id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Group
	for rows.Next() {
		var g Group
		rows.Scan(&g.ID, &g.Name, &g.Count, &g.Created)
		out = append(out, g)
	}
	return out, nil
}

// ---- Unsub ----
func (d *DB) AddUnsub(phone string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO unsubscribed (phone) VALUES (?)`, phone)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteUnsub(id int64) error { return d.del("unsubscribed", id) }
func (d *DB) IsUnsub(phone string) bool {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM unsubscribed WHERE phone=?`, phone).Scan(&n)
	return n > 0
}
func (d *DB) ListUnsub() ([]Unsub, error) {
	rows, err := d.sql.Query(`SELECT id, phone, created_at FROM unsubscribed ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Unsub
	for rows.Next() {
		var u Unsub
		rows.Scan(&u.ID, &u.Phone, &u.Created)
		out = append(out, u)
	}
	return out, nil
}

// ---- Templates ----
func (d *DB) AddTemplate(name, content string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO templates (name, content) VALUES (?, ?)`, name, content)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateTemplate(id int64, name, content string) error {
	_, err := d.sql.Exec(`UPDATE templates SET name=?, content=? WHERE id=?`, name, content, id)
	return err
}
func (d *DB) GetTemplate(id int64) (*Template, error) {
	var t Template
	err := d.sql.QueryRow(`SELECT id, name, content, created_at FROM templates WHERE id=?`, id).Scan(&t.ID, &t.Name, &t.Content, &t.Created)
	if err != nil { return nil, err }
	return &t, nil
}
func (d *DB) DeleteTemplate(id int64) error { _, err := d.sql.Exec(`DELETE FROM templates WHERE id=?`, id); return err }
func (d *DB) ListTemplates() ([]Template, error) {
	rows, err := d.sql.Query(`SELECT id, name, content, created_at FROM templates ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Template
	for rows.Next() {
		var t Template
		rows.Scan(&t.ID, &t.Name, &t.Content, &t.Created)
		out = append(out, t)
	}
	return out, nil
}

// ---- API Keys ----
func (d *DB) AddAPIKey(name, secret string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO api_keys (name, secret) VALUES (?, ?)`, name, secret)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteAPIKey(id int64) error { return d.del("api_keys", id) }
func (d *DB) ValidAPIKey(secret string) bool {
	var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM api_keys WHERE secret=?`, secret).Scan(&n); return n > 0
}
func (d *DB) ListAPIKeys() ([]APIKey, error) {
	rows, err := d.sql.Query(`SELECT id, name, secret, created_at FROM api_keys ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []APIKey
	for rows.Next() {
		var a APIKey
		rows.Scan(&a.ID, &a.Name, &a.Secret, &a.Created)
		out = append(out, a)
	}
	return out, nil
}

// ---- Webhooks ----
func (d *DB) AddWebhook(name, url, event string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO webhooks (name, url, event) VALUES (?, ?, ?)`, name, url, event)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteWebhook(id int64) error { _, err := d.sql.Exec(`DELETE FROM webhooks WHERE id=?`, id); return err }
func (d *DB) ListWebhooks() ([]Webhook, error) {
	rows, err := d.sql.Query(`SELECT id, name, url, event, created_at FROM webhooks ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Webhook
	for rows.Next() {
		var web Webhook
		rows.Scan(&web.ID, &web.Name, &web.URL, &web.Event, &web.Created)
		out = append(out, web)
	}
	return out, nil
}
func (d *DB) WebhooksForEvent(event string) ([]Webhook, error) {
	rows, err := d.sql.Query(`SELECT id, name, url, event, created_at FROM webhooks WHERE event=? OR event='all'`, event)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Webhook
	for rows.Next() {
		var web Webhook
		rows.Scan(&web.ID, &web.Name, &web.URL, &web.Event, &web.Created)
		out = append(out, web)
	}
	return out, nil
}

// ---- Campaigns ----
func (d *DB) AddCampaign(name, groups, message string, total int, accountIDs string, interval int) (int64, error) {
	res, err := d.sql.Exec("INSERT INTO campaigns (name, `groups`, message, total, status, account_ids, msg_interval) VALUES (?, ?, ?, ?, 'pending', ?, ?)", name, groups, message, total, accountIDs, interval)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateCampaignStatus(id int64, status string) error {
	_, err := d.sql.Exec(`UPDATE campaigns SET status=? WHERE id=?`, status, id)
	return err
}
func (d *DB) IncCampaignSent(id int64) error {
	_, err := d.sql.Exec(`UPDATE campaigns SET sent=sent+1 WHERE id=?`, id)
	return err
}
func (d *DB) AppendCampaignSentTo(id int64, phone string) error {
	_, err := d.sql.Exec(`UPDATE campaigns SET sent_to=CONCAT(IFNULL(sent_to,''),?) WHERE id=?`, phone+",", id)
	return err
}
func (d *DB) DeleteCampaign(id int64) error { _, err := d.sql.Exec(`DELETE FROM campaigns WHERE id=?`, id); return err }
func (d *DB) ListCampaigns() ([]Campaign, error) {
	rows, err := d.sql.Query("SELECT id, name, `groups`, message, total, sent, status, IFNULL(account_id,''), IFNULL(account_ids,''), IFNULL(msg_interval,3), IFNULL(sent_to,''), created_at FROM campaigns ORDER BY id DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Campaign
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Name, &c.Groups, &c.Message, &c.Total, &c.Sent, &c.Status, &c.AccountID, &c.AccountIDs, &c.Interval, &c.SentTo, &c.Created)
		out = append(out, c)
	}
	return out, nil
}
func (d *DB) PendingCampaigns() ([]Campaign, error) {
	rows, err := d.sql.Query("SELECT id, name, `groups`, message, total, sent, status, IFNULL(account_id,''), IFNULL(account_ids,''), IFNULL(msg_interval,3), IFNULL(sent_to,''), created_at FROM campaigns WHERE status='running' OR status='pending'")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Campaign
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Name, &c.Groups, &c.Message, &c.Total, &c.Sent, &c.Status, &c.AccountID, &c.AccountIDs, &c.Interval, &c.SentTo, &c.Created)
		out = append(out, c)
	}
	return out, nil
}

// ---- Scheduled ----
func (d *DB) AddScheduled(name, phone, message, sendAt string, repeat int, accountIDs string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO scheduled (name, phone, message, send_at, repeat_min, status, account_ids) VALUES (?, ?, ?, ?, ?, 'pending', ?)`, name, phone, message, sendAt, repeat, accountIDs)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteScheduled(id int64) error { _, err := d.sql.Exec(`DELETE FROM scheduled WHERE id=?`, id); return err }
func (d *DB) MarkScheduledSent(id int64) error {
	_, err := d.sql.Exec(`UPDATE scheduled SET status='sent' WHERE id=?`, id)
	return err
}
func (d *DB) RescheduleAfter(id int64, minutes int) error {
	_, err := d.sql.Exec(`UPDATE scheduled SET send_at = DATE_ADD(send_at, INTERVAL ? MINUTE) WHERE id=?`, minutes, id)
	return err
}
func (d *DB) ListScheduled() ([]Scheduled, error) {
	rows, err := d.sql.Query(`SELECT id, name, phone, message, send_at, repeat_min, status, IFNULL(account_ids,''), created_at FROM scheduled ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Scheduled
	for rows.Next() {
		var s Scheduled
		rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Message, &s.SendAt, &s.Repeat, &s.Status, &s.AccountIDs, &s.Created)
		out = append(out, s)
	}
	return out, nil
}
func (d *DB) DueScheduled() ([]Scheduled, error) {
	rows, err := d.sql.Query(`SELECT id, name, phone, message, send_at, repeat_min, status, IFNULL(account_ids,''), created_at FROM scheduled WHERE status='pending' AND send_at <= NOW()`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Scheduled
	for rows.Next() {
		var s Scheduled
		rows.Scan(&s.ID, &s.Name, &s.Phone, &s.Message, &s.SendAt, &s.Repeat, &s.Status, &s.AccountIDs, &s.Created)
		out = append(out, s)
	}
	return out, nil
}

// ---- Logger ----
func (d *DB) Log(typ, reason, content string) {
	d.sql.Exec(`INSERT INTO logger (type, reason, content) VALUES (?,?,?)`, typ, reason, content)
}
func (d *DB) ListLog(limit int) ([]LogEntry, error) {
	rows, err := d.sql.Query(`SELECT id, type, reason, content, created_at FROM logger ORDER BY id DESC LIMIT ?`, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []LogEntry
	for rows.Next() {
		var l LogEntry
		rows.Scan(&l.ID, &l.Type, &l.Reason, &l.Content, &l.Created)
		out = append(out, l)
	}
	return out, nil
}
func (d *DB) ClearLog() error {
	_, err := d.sql.Exec(`DELETE FROM logger`)
	return err
}

func (d *DB) CountRunningCampaigns() (int, error) {
	var n int; err := d.sql.QueryRow(`SELECT COUNT(*) FROM campaigns WHERE status='running'`).Scan(&n); return n, err
}

func (d *DB) MessageChartData() (labels, sent, received string) {
	var lb, s, r []string
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("02 Jan")
		lb = append(lb, `"`+date+`"`)
		var sc, rc int
		d.sql.QueryRow(`SELECT COUNT(*) FROM sent WHERE DATE(created_at)=DATE_SUB(CURDATE(), INTERVAL ? DAY)`, i).Scan(&sc)
		d.sql.QueryRow(`SELECT COUNT(*) FROM received WHERE DATE(created_at)=DATE_SUB(CURDATE(), INTERVAL ? DAY)`, i).Scan(&rc)
		s = append(s, fmt.Sprintf("%d", sc))
		r = append(r, fmt.Sprintf("%d", rc))
	}
	return strings.Join(lb, ","), strings.Join(s, ","), strings.Join(r, ",")
}

func (d *DB) CountUsers() int {
	var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); return n
}
