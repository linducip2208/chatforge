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
	ID       int64
	Name     string
	Count    int
	Language string
	Created  string
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
type Tag struct {
	ID      int64
	Name    string
	Color   string
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
	ID              int64
	Name            string
	Groups          string
	Tags            string
	Numbers         string
	MediaType       string
	MediaURL        string
	AccountID       string
	AccountIDs      string
	MetaAccountID   int64
	MetaTemplate    string
	SendMode        string
	Message         string
	Total           int
	Sent            int
	Status          string
	Interval        int
	SentTo          string
	Created         string
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
		`CREATE TABLE IF NOT EXISTS tags (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, color VARCHAR(7) NOT NULL DEFAULT '#2c7be5', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS contact_tags (contact_id BIGINT NOT NULL, tag_id BIGINT NOT NULL, PRIMARY KEY (contact_id, tag_id), FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE, FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS canned_responses (id BIGINT AUTO_INCREMENT PRIMARY KEY, shortcut VARCHAR(32) NOT NULL DEFAULT '', name VARCHAR(255) NOT NULL, message TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS agent_assignments (phone VARCHAR(64) PRIMARY KEY, agent_id BIGINT NOT NULL DEFAULT 0, assigned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, status VARCHAR(20) NOT NULL DEFAULT 'open') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS link_clicks (id BIGINT AUTO_INCREMENT PRIMARY KEY, token VARCHAR(16) NOT NULL UNIQUE, url VARCHAR(2048) NOT NULL, campaign_id BIGINT NOT NULL DEFAULT 0, phone VARCHAR(64) NOT NULL DEFAULT '', clicked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS ab_tests (id BIGINT AUTO_INCREMENT PRIMARY KEY, campaign_id BIGINT NOT NULL, variant_a TEXT NOT NULL, variant_b TEXT NOT NULL, a_sent INT NOT NULL DEFAULT 0, b_sent INT NOT NULL DEFAULT 0, a_replied INT NOT NULL DEFAULT 0, b_replied INT NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN send_mode VARCHAR(20) NOT NULL DEFAULT 'round_robin' AFTER account_ids`)
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN numbers TEXT NOT NULL AFTER ` + "`groups`")
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN meta_account_id BIGINT NOT NULL DEFAULT 0 AFTER send_mode`)
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN meta_template VARCHAR(255) NOT NULL DEFAULT '' AFTER meta_account_id`)
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN media_type VARCHAR(20) NOT NULL DEFAULT '' AFTER numbers`)
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN media_url VARCHAR(1024) NOT NULL DEFAULT '' AFTER media_type`)
	d.sql.Exec(`ALTER TABLE campaigns ADD COLUMN tags VARCHAR(512) NOT NULL DEFAULT '' AFTER ` + "`groups`")
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
func (d *DB) FindContactByPhone(phone string) (*Contact, error) {
	var c Contact
	err := d.sql.QueryRow(`SELECT id, name, phone, `+"`groups`"+`, created_at FROM contacts WHERE phone=?`, phone).Scan(&c.ID, &c.Name, &c.Phone, &c.Groups, &c.Created)
	if err != nil { return nil, err }
	return &c, nil
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
	rows, err := d.sql.Query(`SELECT g.id, g.name, COUNT(c.id), IFNULL(g.language,''), g.created_at FROM contact_groups g LEFT JOIN contacts c ON FIND_IN_SET(g.id, REPLACE(c.`+"`groups`"+`, ' ', '')) GROUP BY g.id ORDER BY g.id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Group
	for rows.Next() {
		var g Group
		rows.Scan(&g.ID, &g.Name, &g.Count, &g.Language, &g.Created)
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
func (d *DB) AddCampaign(name, groups, numbers, mediaType, mediaURL, message string, total int, accountIDs, sendMode string, interval int, metaAccountID int64, metaTemplate, tags string) (int64, error) {
	res, err := d.sql.Exec("INSERT INTO campaigns (name, `groups`, tags, numbers, media_type, media_url, message, total, status, account_ids, send_mode, meta_account_id, meta_template, msg_interval) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?, ?, ?, ?, ?)", name, groups, tags, numbers, mediaType, mediaURL, message, total, accountIDs, sendMode, metaAccountID, metaTemplate, interval)
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
	rows, err := d.sql.Query("SELECT id, name, `groups`, IFNULL(tags,''), IFNULL(numbers,''), IFNULL(media_type,''), IFNULL(media_url,''), message, total, sent, status, IFNULL(account_id,''), IFNULL(account_ids,''), IFNULL(send_mode,'round_robin'), IFNULL(meta_account_id,0), IFNULL(meta_template,''), IFNULL(msg_interval,3), IFNULL(sent_to,''), created_at FROM campaigns ORDER BY id DESC")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Campaign
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Name, &c.Groups, &c.Tags, &c.Numbers, &c.MediaType, &c.MediaURL, &c.Message, &c.Total, &c.Sent, &c.Status, &c.AccountID, &c.AccountIDs, &c.SendMode, &c.MetaAccountID, &c.MetaTemplate, &c.Interval, &c.SentTo, &c.Created)
		out = append(out, c)
	}
	return out, nil
}
func (d *DB) PendingCampaigns() ([]Campaign, error) {
	rows, err := d.sql.Query("SELECT id, name, `groups`, IFNULL(tags,''), IFNULL(numbers,''), IFNULL(media_type,''), IFNULL(media_url,''), message, total, sent, status, IFNULL(account_id,''), IFNULL(account_ids,''), IFNULL(send_mode,'round_robin'), IFNULL(meta_account_id,0), IFNULL(meta_template,''), IFNULL(msg_interval,3), IFNULL(sent_to,''), created_at FROM campaigns WHERE status='running' OR status='pending'")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Campaign
	for rows.Next() {
		var c Campaign
		rows.Scan(&c.ID, &c.Name, &c.Groups, &c.Tags, &c.Numbers, &c.MediaType, &c.MediaURL, &c.Message, &c.Total, &c.Sent, &c.Status, &c.AccountID, &c.AccountIDs, &c.SendMode, &c.MetaAccountID, &c.MetaTemplate, &c.Interval, &c.SentTo, &c.Created)
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
	return d.ListLogPaginated(1, limit)
}
func (d *DB) ListLogPaginated(page, perPage int) ([]LogEntry, error) {
	offset := (page - 1) * perPage
	rows, err := d.sql.Query(`SELECT id, type, reason, content, created_at FROM logger ORDER BY id DESC LIMIT ? OFFSET ?`, perPage, offset)
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
func (d *DB) CountLog() int {
	var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM logger`).Scan(&n); return n
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

// ---- Tags ----
func (d *DB) AddTag(name, color string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO tags (name, color) VALUES (?, ?)`, name, color)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteTag(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM tags WHERE id=?`, id)
	return err
}
func (d *DB) ListTags() ([]Tag, error) {
	rows, err := d.sql.Query(`SELECT id, name, color, created_at FROM tags ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Tag
	for rows.Next() {
		var t Tag
		rows.Scan(&t.ID, &t.Name, &t.Color, &t.Created)
		out = append(out, t)
	}
	return out, nil
}
func (d *DB) SetContactTags(contactID int64, tagIDs []int64) error {
	d.sql.Exec(`DELETE FROM contact_tags WHERE contact_id=?`, contactID)
	for _, tid := range tagIDs {
		d.sql.Exec(`INSERT IGNORE INTO contact_tags (contact_id, tag_id) VALUES (?, ?)`, contactID, tid)
	}
	return nil
}
func (d *DB) ContactsByTag(tagID int64) ([]Contact, error) {
	rows, err := d.sql.Query(`SELECT c.id, c.name, c.phone, c.groups, c.created_at FROM contacts c JOIN contact_tags ct ON ct.contact_id=c.id WHERE ct.tag_id=? ORDER BY c.name`, tagID)
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
func (d *DB) GetContactTags(contactID int64) []int64 {
	rows, err := d.sql.Query(`SELECT tag_id FROM contact_tags WHERE contact_id=?`, contactID)
	if err != nil { return nil }
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var tid int64
		rows.Scan(&tid)
		out = append(out, tid)
	}
	return out
}
func (d *DB) CampaignReport(campaignID int64) (sent, replied int) {
	d.sql.QueryRow(`SELECT COUNT(*) FROM sent WHERE id IN (SELECT id FROM campaigns WHERE id=?)`, campaignID).Scan(&sent)
	// count replies: messages from same phone numbers within 24h of campaign start
	return sent, 0
}

// Rate limit
func (d *DB) TodaySentCount(phone string) int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM sent WHERE channel='whatsmeow' AND DATE(created_at)=CURDATE()`).Scan(&n)
	return n
}

// ---- Canned Responses ----
type CannedResponse struct {
	ID       int64
	Shortcut string
	Name     string
	Message  string
	Created  string
}

func (d *DB) AddCanned(shortcut, name, message string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO canned_responses (shortcut, name, message) VALUES (?, ?, ?)`, shortcut, name, message)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteCanned(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM canned_responses WHERE id=?`, id)
	return err
}
func (d *DB) ListCanned() ([]CannedResponse, error) {
	rows, err := d.sql.Query(`SELECT id, shortcut, name, message, created_at FROM canned_responses ORDER BY shortcut`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []CannedResponse
	for rows.Next() {
		var c CannedResponse
		rows.Scan(&c.ID, &c.Shortcut, &c.Name, &c.Message, &c.Created)
		out = append(out, c)
	}
	return out, nil
}

// ---- Agent Assignment ----
func (d *DB) AssignAgent(phone string, agentID int64) error {
	_, err := d.sql.Exec(`INSERT INTO agent_assignments (phone, agent_id, status) VALUES (?, ?, 'open') ON DUPLICATE KEY UPDATE agent_id=VALUES(agent_id)`, phone, agentID)
	return err
}
func (d *DB) CloseConversation(phone string) error {
	_, err := d.sql.Exec(`UPDATE agent_assignments SET status='closed' WHERE phone=?`, phone)
	return err
}
func (d *DB) GetAssignedAgent(phone string) int64 {
	var n int64
	d.sql.QueryRow(`SELECT agent_id FROM agent_assignments WHERE phone=?`, phone).Scan(&n)
	return n
}
func (d *DB) GetUnassignedCount() int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM agent_assignments WHERE agent_id=0 AND status='open'`).Scan(&n)
	return n
}
func (d *DB) AssignNextRoundRobin(phone string) int64 {
	users, _ := d.ListUsers()
	var ids []int64
	for _, u := range users {
		if u.Role == "admin" || u.Role == "agent" { ids = append(ids, u.ID) }
	}
	if len(ids) == 0 { return 0 }
	var assigned int64
	d.sql.QueryRow(`SELECT IFNULL(MAX(agent_id),0) FROM agent_assignments ORDER BY assigned_at DESC LIMIT 1`).Scan(&assigned)
	nextIdx := 0
	for i, id := range ids {
		if id > assigned { nextIdx = i; break }
	}
	if nextIdx >= len(ids) { nextIdx = 0 }
	agentID := ids[nextIdx]
	d.sql.Exec(`INSERT INTO agent_assignments (phone, agent_id, status) VALUES (?, ?, 'open') ON DUPLICATE KEY UPDATE agent_id=VALUES(agent_id)`, phone, agentID)
	return agentID
}

// ---- Link Tracking ----
func (d *DB) TrackLink(token, url string, campaignID int64, phone string) {
	d.sql.Exec(`INSERT INTO link_clicks (token, url, campaign_id, phone) VALUES (?, ?, ?, ?)`, token, url, campaignID, phone)
}
func (d *DB) LogLinkClick(token string) {
	d.sql.Exec(`UPDATE link_clicks SET clicked_at=NOW() WHERE token=? AND clicked_at=created_at`, token)
}
func (d *DB) LinkClicks(campaignID int64) int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM link_clicks WHERE campaign_id=? AND clicked_at > created_at`, campaignID).Scan(&n)
	return n
}
func (d *DB) ListLinkClicks() ([]LinkClick, error) {
	rows, err := d.sql.Query(`SELECT id, token, url, campaign_id, phone, clicked_at > created_at as is_clicked, created_at FROM link_clicks ORDER BY id DESC LIMIT 200`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []LinkClick
	for rows.Next() {
		var v LinkClick
		rows.Scan(&v.ID, &v.Token, &v.URL, &v.CampaignID, &v.Phone, &v.Clicked, &v.Created)
		out = append(out, v)
	}
	return out, nil
}

// ---- A/B Testing ----
type ABTest struct {
	ID         int64
	CampaignID int64
	VariantA   string
	VariantB   string
	ASent      int
	BSent      int
	AReplied   int
	BReplied   int
	Created    string
}

type LinkClick struct {
	ID         int64
	Token      string
	URL        string
	CampaignID int64
	Phone      string
	Clicked    bool
	Created    string
}

func (d *DB) CreateABTest(campaignID int64, variantA, variantB string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO ab_tests (campaign_id, variant_a, variant_b) VALUES (?, ?, ?)`, campaignID, variantA, variantB)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) GetABTest(campaignID int64) (*ABTest, error) {
	var ab ABTest
	err := d.sql.QueryRow(`SELECT id, campaign_id, variant_a, variant_b, a_sent, b_sent, a_replied, b_replied, created_at FROM ab_tests WHERE campaign_id=?`, campaignID).Scan(&ab.ID, &ab.CampaignID, &ab.VariantA, &ab.VariantB, &ab.ASent, &ab.BSent, &ab.AReplied, &ab.BReplied, &ab.Created)
	if err != nil { return nil, err }
	return &ab, nil
}
func (d *DB) ListABTests() ([]ABTest, error) {
	rows, err := d.sql.Query(`SELECT id, campaign_id, variant_a, variant_b, a_sent, b_sent, a_replied, b_replied, created_at FROM ab_tests ORDER BY id DESC LIMIT 50`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ABTest
	for rows.Next() {
		var ab ABTest
		rows.Scan(&ab.ID, &ab.CampaignID, &ab.VariantA, &ab.VariantB, &ab.ASent, &ab.BSent, &ab.AReplied, &ab.BReplied, &ab.Created)
		out = append(out, ab)
	}
	return out, nil
}
func (d *DB) IncABSent(campaignID int64, variant string) {
	if variant == "b" {
		d.sql.Exec(`UPDATE ab_tests SET b_sent=b_sent+1 WHERE campaign_id=?`, campaignID)
	} else {
		d.sql.Exec(`UPDATE ab_tests SET a_sent=a_sent+1 WHERE campaign_id=?`, campaignID)
	}
}
