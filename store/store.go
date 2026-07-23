package store

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	sql *sql.DB
	mu  sync.Mutex
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.sql.QueryRow(query, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.sql.Query(query, args...)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.sql.Exec(query, args...)
}

type AutoReply struct {
	ID         int64
	Keyword    string
	Match      string
	Reply      string
	IsActive   bool
	UseAI      bool
	AiKeyID    int64
	AccountID  string
	TrainingID int64
	UserID     int64
	Created    string
	MediaType  string
	MediaURL   string
}

type SentMessage struct {
	ID      int64
	Phone   string
	Message string
	Status  string
	Channel string
	Created string
}

type ReceivedMessage struct {
	ID          int64
	Phone       string
	Name        string
	Message     string
	IsGroup     bool
	SenderPhone string
	SenderName  string
	IsRead      bool
	Channel     string
	Created     string
}

type InboxConversation struct {
	Phone    string
	Name     string
	LastMsg  string
	LastTime string
	Unread   int
	IsGroup  bool
	Channel  string
}

type ChatMessage struct {
	Type       string
	ID         int64
	Phone      string
	Name       string
	Message    string
	Created    string
	IsRead     bool
	SenderName string
	IsGroup    bool
	Channel    string
}

func Open(dsn string) (*DB, error) {
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping: %w", err)
	}
	sqldb.SetMaxOpenConns(10)
	if _, err := sqldb.Exec("SET time_zone = '+07:00'"); err != nil {
		return nil, fmt.Errorf("set timezone: %w", err)
	}
	db := &DB{sql: sqldb}
	if err := db.migrate(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d *DB) columnExists(table, column string) bool {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=? AND COLUMN_NAME=?`, table, column).Scan(&n)
	return n > 0
}

func (d *DB) safeAddColumn(table, column, def string) {
	if !d.columnExists(table, column) {
		d.sql.Exec(`ALTER TABLE ` + table + ` ADD COLUMN ` + column + ` ` + def)
	}
}

func (d *DB) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS autoreplies (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			keyword VARCHAR(255) NOT NULL,
			match_type VARCHAR(20) NOT NULL DEFAULT 'contains',
			reply TEXT NOT NULL,
			is_active TINYINT NOT NULL DEFAULT 1,
			use_ai TINYINT NOT NULL DEFAULT 0,
			ai_key_id BIGINT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS sent (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			phone VARCHAR(64) NOT NULL,
			message TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'sent',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS received (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			phone VARCHAR(64) NOT NULL,
			name VARCHAR(255) NOT NULL DEFAULT '',
			message TEXT NOT NULL,
			is_group TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS settings (
			name VARCHAR(64) PRIMARY KEY,
			value TEXT NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS welcomed (
			phone VARCHAR(64) PRIMARY KEY,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	if err := d.migrateExtra(); err != nil {
		return err
	}
	if err := d.migrateAdmin(); err != nil {
		return err
	}
	if _, err := d.sql.Exec(`ALTER TABLE received ADD COLUMN is_read TINYINT NOT NULL DEFAULT 0`); err != nil {
		if !strings.Contains(err.Error(), "Duplicate") && !strings.Contains(err.Error(), "1060") {
			return fmt.Errorf("add is_read: %w", err)
		}
	}
	d.safeAddColumn("received", "sender_phone", "VARCHAR(64) NOT NULL DEFAULT ''")
	d.safeAddColumn("received", "sender_name", "VARCHAR(255) NOT NULL DEFAULT ''")
	d.safeAddColumn("received", "channel", "VARCHAR(20) NOT NULL DEFAULT 'whatsmeow'")
	d.safeAddColumn("sent", "channel", "VARCHAR(20) NOT NULL DEFAULT 'whatsmeow'")
	d.safeAddColumn("packages", "meta_limit", "INT NOT NULL DEFAULT 0")
	d.safeAddColumn("meta_accounts", "user_id", "BIGINT NOT NULL DEFAULT 0")
	d.safeAddColumn("meta_accounts", "parent_id", "BIGINT NOT NULL DEFAULT 0")
	d.safeAddColumn("autoreplies", "user_id", "BIGINT NOT NULL DEFAULT 0")
	d.safeAddColumn("autoreplies", "media_type", "VARCHAR(20) NOT NULL DEFAULT ''")
	d.safeAddColumn("autoreplies", "media_url", "VARCHAR(1024) NOT NULL DEFAULT ''")
	d.safeAddColumn("subscriptions", "status", "VARCHAR(20) NOT NULL DEFAULT 'active'")
	d.safeAddColumn("subscriptions", "user_id", "BIGINT NOT NULL DEFAULT 0")
	d.safeAddColumn("subscriptions", "package_id", "BIGINT NOT NULL DEFAULT 0")
	d.safeAddColumn("sent", "starred", "TINYINT NOT NULL DEFAULT 0")
	d.safeAddColumn("received", "starred", "TINYINT NOT NULL DEFAULT 0")
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS faq (id BIGINT AUTO_INCREMENT PRIMARY KEY, question VARCHAR(500) NOT NULL, answer TEXT NOT NULL, user_id BIGINT NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	d.sql.Exec(`CREATE TABLE IF NOT EXISTS wa_session_owners (phone VARCHAR(64) PRIMARY KEY, user_id BIGINT NOT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	d.safeAddColumn("received", "wa_phone", "VARCHAR(64) NOT NULL DEFAULT ''")
	d.safeAddColumn("sent", "wa_phone", "VARCHAR(64) NOT NULL DEFAULT ''")
	d.safeAddColumn("sent", "msg_id", "VARCHAR(128) NOT NULL DEFAULT ''")
	d.migrateInstanceLog()
	d.migrateStatuses()
	d.migrateMeta()
	d.migrateDrip()
	d.migratePayment()
	d.migrateStore()
	d.migrateSafety()
	d.migratePlus()
	d.migrateFinal()
	d.migrateFlow()
	d.migrateFlowStats()
	d.migrateFlowVersions()
	d.migrateFlowMarket()
	return nil
}

// AutoReply CRUD
func (d *DB) AddAutoReply(userID int64, keyword, match, reply string, useAI bool, aiKeyID int64, accountID string, trainingID int64, mediaType, mediaURL string) (int64, error) {
	ai := 0; if useAI { ai = 1 }
	res, err := d.sql.Exec(`INSERT INTO autoreplies (user_id,keyword,match_type,reply,is_active,use_ai,ai_key_id,account_id,training_id,media_type,media_url) VALUES (?,?,?,?,1,?,?,?,?,?,?)`, userID, keyword, match, reply, ai, aiKeyID, accountID, trainingID, mediaType, mediaURL)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteAutoReply(userID int64, id int64) error {
	if userID == 0 {
		_, err := d.sql.Exec(`DELETE FROM autoreplies WHERE id=?`, id)
		return err
	}
	_, err := d.sql.Exec(`DELETE FROM autoreplies WHERE id=? AND user_id=?`, id, userID)
	return err
}
func (d *DB) ToggleAutoReply(userID int64, id int64) error {
	if userID == 0 {
		_, err := d.sql.Exec(`UPDATE autoreplies SET is_active=1-is_active WHERE id=?`, id)
		return err
	}
	_, err := d.sql.Exec(`UPDATE autoreplies SET is_active=1-is_active WHERE id=? AND user_id=?`, id, userID)
	return err
}
func (d *DB) GetAutoReply(userID int64, id int64) (*AutoReply, error) {
	var a AutoReply; var active int
	var err error
	if userID == 0 {
		err = d.sql.QueryRow(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),IFNULL(user_id,0),IFNULL(media_type,''),IFNULL(media_url,''),created_at FROM autoreplies WHERE id=?`, id).Scan(&a.ID, &a.Keyword, &a.Match, &a.Reply, &active, &a.UseAI, &a.AiKeyID, &a.AccountID, &a.TrainingID, &a.UserID, &a.MediaType, &a.MediaURL, &a.Created)
	} else {
		err = d.sql.QueryRow(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),IFNULL(user_id,0),IFNULL(media_type,''),IFNULL(media_url,''),created_at FROM autoreplies WHERE id=? AND user_id=?`, id, userID).Scan(&a.ID, &a.Keyword, &a.Match, &a.Reply, &active, &a.UseAI, &a.AiKeyID, &a.AccountID, &a.TrainingID, &a.UserID, &a.MediaType, &a.MediaURL, &a.Created)
	}
	a.IsActive = active == 1
	if err != nil { return nil, err }
	return &a, nil
}
func (d *DB) UpdateAutoReply(userID int64, id int64, keyword, match, reply string, useAI bool, aiKeyID int64, accountID string, trainingID int64, mediaType, mediaURL string) error {
	use := 0
	if useAI { use = 1 }
	if userID == 0 {
		_, err := d.sql.Exec(`UPDATE autoreplies SET keyword=?, match_type=?, reply=?, use_ai=?, ai_key_id=?, account_id=?, training_id=?, media_type=?, media_url=? WHERE id=?`, keyword, match, reply, use, aiKeyID, accountID, trainingID, mediaType, mediaURL, id)
		return err
	}
	_, err := d.sql.Exec(`UPDATE autoreplies SET keyword=?, match_type=?, reply=?, use_ai=?, ai_key_id=?, account_id=?, training_id=?, media_type=?, media_url=? WHERE id=? AND user_id=?`, keyword, match, reply, use, aiKeyID, accountID, trainingID, mediaType, mediaURL, id, userID)
	return err
}
func (d *DB) ListAutoReplies(userID int64) ([]AutoReply, error) {
	var rows *sql.Rows
	var err error
	if userID == 0 {
		rows, err = d.sql.Query(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),IFNULL(user_id,0),IFNULL(media_type,''),IFNULL(media_url,''),created_at FROM autoreplies ORDER BY id DESC`)
	} else {
		rows, err = d.sql.Query(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),IFNULL(user_id,0),IFNULL(media_type,''),IFNULL(media_url,''),created_at FROM autoreplies WHERE user_id=? ORDER BY id DESC`, userID)
	}
	if err != nil { return nil, err }
	defer rows.Close()
	var out []AutoReply
	for rows.Next() {
		var a AutoReply; var active int
		if err := rows.Scan(&a.ID,&a.Keyword,&a.Match,&a.Reply,&active,&a.UseAI,&a.AiKeyID,&a.AccountID,&a.TrainingID,&a.UserID,&a.MediaType,&a.MediaURL,&a.Created); err != nil { return nil, err }
		a.IsActive = active==1
		out = append(out, a)
	}
	return out, nil
}
func (d *DB) FindReply(incoming string) (string, bool) {
	r, ok := d.FindReplyFull(0, incoming)
	return r.Reply, ok
}
func (d *DB) FindReplyFullForAccount(userID int64, incoming string, accountPhone string) (AutoReply, bool) {
	rules, _ := d.ListAutoReplies(userID)
	msg := strings.ToLower(strings.TrimSpace(incoming))
	for _, r := range rules {
		if !r.IsActive { continue }
		// account filter: empty = all, otherwise check if phone is in comma-separated list
		if accountPhone != "" && r.AccountID != "" {
			phone := "+" + accountPhone
			found := false
			for _, a := range strings.Split(r.AccountID, ",") {
				if strings.TrimSpace(a) == phone {
					found = true
					break
				}
			}
			if !found { continue }
		}
		kw := strings.ToLower(strings.TrimSpace(r.Keyword))
		// support comma-separated keywords
		keywords := strings.Split(kw, ",")
		matched := false
		for _, k := range keywords {
			k = strings.TrimSpace(k)
			if k == "" { continue }
			switch r.Match {
			case "exact":
				if msg == k { matched = true }
			case "starts_with":
				if strings.HasPrefix(msg, k) { matched = true }
			case "ai":
				matched = true
			default:
				if strings.Contains(msg, k) { matched = true }
			}
			if matched { break }
		}
		if matched { return r, true }
	}
	return AutoReply{}, false
}
func (d *DB) FindReplyFull(userID int64, incoming string) (AutoReply, bool) {
	rules, _ := d.ListAutoReplies(userID)
	msg := strings.ToLower(strings.TrimSpace(incoming))
	for _, r := range rules {
		if !r.IsActive { continue }
		kw := strings.ToLower(strings.TrimSpace(r.Keyword))
		if kw == "" { continue }
		keywords := strings.Split(kw, ",")
		matched := false
		for _, k := range keywords {
			k = strings.TrimSpace(k)
			if k == "" { continue }
			switch r.Match {
			case "exact": matched = msg == k
			case "starts_with": matched = strings.HasPrefix(msg, k)
			case "ai": matched = true
			default: matched = strings.Contains(msg, k)
			}
			if matched { break }
		}
		if matched { return r, true }
	}
	return AutoReply{}, false
}

// Sent log
func (d *DB) LogSent(phone, message, status, channel string) {
	if channel == "" { channel = "whatsmeow" }
	d.sql.Exec(`INSERT INTO sent (phone,message,status,channel) VALUES (?,?,?,?)`, phone, message, status, channel)
}
func (d *DB) LogSentForWA(waPhone, phone, message, status, channel string) {
	if channel == "" { channel = "whatsmeow" }
	d.sql.Exec(`INSERT INTO sent (wa_phone,phone,message,status,channel) VALUES (?,?,?,?,?)`, waPhone, phone, message, status, channel)
}
func (d *DB) UpdateSentStatus(messageID, newStatus string) {
	d.sql.Exec(`UPDATE sent SET status=? WHERE msg_id=?`, newStatus, messageID)
}
func (d *DB) GetSentStatus(messageID string) string {
	var s string
	d.sql.QueryRow(`SELECT status FROM sent WHERE msg_id=? LIMIT 1`, messageID).Scan(&s)
	return s
}
func (d *DB) ListSent(limit int) ([]SentMessage, error) { return d.ListSentPaginated(0, 1, limit) }
func (d *DB) ListSentPaginated(userID int64, page, perPage int) ([]SentMessage, error) {
	offset := (page-1)*perPage
	rows, err := d.sql.Query(`SELECT id,phone,message,status,channel,created_at FROM sent WHERE wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) ORDER BY id DESC LIMIT ? OFFSET ?`, userID, perPage, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []SentMessage
	for rows.Next() { var m SentMessage; rows.Scan(&m.ID,&m.Phone,&m.Message,&m.Status,&m.Channel,&m.Created); out = append(out, m) }
	return out, nil
}
func (d *DB) CountSent(userID int64) int { var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM sent WHERE wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?)`, userID).Scan(&n); return n }

// Received log
func (d *DB) LogReceived(phone, name, message string, isGroup bool, senderPhone, senderName, channel string) {
	g := 0; if isGroup { g = 1 }
	if channel == "" { channel = "whatsmeow" }
	d.sql.Exec(`INSERT INTO received (phone,name,message,is_group,sender_phone,sender_name,channel) VALUES (?,?,?,?,?,?,?)`, phone, name, message, g, senderPhone, senderName, channel)
}
func (d *DB) LogReceivedForWA(waPhone, phone, name, message string, isGroup bool, senderPhone, senderName, channel string) {
	g := 0; if isGroup { g = 1 }
	if channel == "" { channel = "whatsmeow" }
	d.sql.Exec(`INSERT INTO received (wa_phone,phone,name,message,is_group,sender_phone,sender_name,channel) VALUES (?,?,?,?,?,?,?,?)`, waPhone, phone, name, message, g, senderPhone, senderName, channel)
}
func (d *DB) ListReceived(limit int) ([]ReceivedMessage, error) { return d.ListReceivedPaginated(0, 1, limit) }
func (d *DB) ListReceivedPaginated(userID int64, page, perPage int) ([]ReceivedMessage, error) {
	offset := (page-1)*perPage
	rows, err := d.sql.Query(`SELECT id,phone,name,message,is_group,is_read,sender_phone,sender_name,channel,created_at FROM received WHERE wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) ORDER BY id DESC LIMIT ? OFFSET ?`, userID, perPage, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ReceivedMessage
	for rows.Next() { var m ReceivedMessage; var g, r int; rows.Scan(&m.ID,&m.Phone,&m.Name,&m.Message,&g,&r,&m.SenderPhone,&m.SenderName,&m.Channel,&m.Created); m.IsGroup = g==1; m.IsRead = r==1; out = append(out, m) }
	return out, nil
}
func (d *DB) CountReceived(userID int64) int { var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM received WHERE wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?)`, userID).Scan(&n); return n }

func (d *DB) GroupInboxPaginated(userID int64, page, perPage int) ([]InboxConversation, error) {
	offset := (page - 1) * perPage
	var rows *sql.Rows
	var err error
	baseSent := `SELECT phone, '' as name, 0 as is_group, 1 as is_read, channel, created_at FROM sent`
	if userID == 0 {
		rows, err = d.sql.Query(`SELECT t.phone, COALESCE(g.name, MAX(t.name)) as name, MAX(t.is_group) as is_group, COUNT(CASE WHEN t.is_read=0 THEN 1 END) as unread, MAX(t.channel) as channel, MAX(t.created_at) as last_time FROM (SELECT phone, name, is_group, is_read, channel, created_at FROM received UNION ALL `+baseSent+`) t LEFT JOIN wa_groups g ON t.phone = g.jid GROUP BY t.phone ORDER BY last_time DESC LIMIT ? OFFSET ?`, perPage, offset)
	} else {
		baseRecv := `SELECT r.phone, r.name, r.is_group, r.is_read, r.channel, r.created_at FROM received r INNER JOIN wa_session_owners o ON r.wa_phone = o.phone AND o.user_id = ?`
		baseSentU := `SELECT s.phone, '' as name, 0 as is_group, 1 as is_read, s.channel, s.created_at FROM sent s LEFT JOIN wa_session_owners o ON s.wa_phone = o.phone AND o.user_id = ? WHERE s.wa_phone = '' OR o.phone IS NOT NULL`
		rows, err = d.sql.Query(`SELECT t.phone, COALESCE(g.name, MAX(t.name)) as name, MAX(t.is_group) as is_group, COUNT(CASE WHEN t.is_read=0 THEN 1 END) as unread, MAX(t.channel) as channel, MAX(t.created_at) as last_time FROM (`+baseRecv+` UNION ALL `+baseSentU+`) t LEFT JOIN wa_groups g ON t.phone = g.jid GROUP BY t.phone ORDER BY last_time DESC LIMIT ? OFFSET ?`, userID, userID, perPage, offset)
	}
	if err != nil { return nil, err }
	defer rows.Close()
	var out []InboxConversation
	for rows.Next() {
		var c InboxConversation; var g int; var lastTime string
		if err := rows.Scan(&c.Phone, &c.Name, &g, &c.Unread, &c.Channel, &lastTime); err != nil { return nil, err }
		c.IsGroup = g == 1
		c.LastTime = lastTime
		// Jika tidak ada nama, coba cari di contacts
		if c.Name == "" && !c.IsGroup {
			cleanPhone := extractCleanPhone(c.Phone)
			var contactName string
			d.sql.QueryRow(`SELECT name FROM contacts WHERE phone=? OR phone=?`, c.Phone, cleanPhone).Scan(&contactName)
			if contactName != "" { c.Name = contactName }
		}
		d.sql.QueryRow(`SELECT message, created_at FROM received WHERE phone=? UNION ALL SELECT message, created_at FROM sent WHERE phone=? ORDER BY created_at DESC LIMIT 1`, c.Phone, c.Phone).Scan(&c.LastMsg, &c.LastTime)
		out = append(out, c)
	}
	return out, nil
}

func (d *DB) CountInbox(userID int64) int {
	var n int
	if userID == 0 {
		d.sql.QueryRow(`SELECT COUNT(DISTINCT phone) FROM (SELECT phone FROM received UNION ALL SELECT phone FROM sent) t`).Scan(&n)
	} else {
		d.sql.QueryRow(`SELECT COUNT(DISTINCT t.phone) FROM (SELECT r.phone FROM received r INNER JOIN wa_session_owners o ON r.wa_phone = o.phone AND o.user_id = ? UNION ALL SELECT s.phone FROM sent s INNER JOIN wa_session_owners o ON s.wa_phone = o.phone AND o.user_id = ?) t`, userID, userID).Scan(&n)
	}
	return n
}

func (d *DB) GroupInbox(userID int64) ([]InboxConversation, error) {
	return d.GroupInboxPaginated(userID, 1, 100)
}

func (d *DB) ChatHistory(userID int64, phone string, limit int) ([]ChatMessage, error) {
	rows, err := d.sql.Query(`SELECT 'received' as type, id, phone, name, message, created_at, is_read, sender_name, is_group, channel FROM received WHERE phone=? AND wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) UNION ALL SELECT 'sent' as type, id, phone, '' as name, message, created_at, 1 as is_read, '' as sender_name, 0 as is_group, channel FROM sent WHERE phone=? AND wa_phone IN (SELECT phone FROM wa_session_owners WHERE user_id=?) ORDER BY created_at ASC LIMIT ?`, phone, userID, phone, userID, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ChatMessage
	for rows.Next() {
		var m ChatMessage; var r, g int
		if err := rows.Scan(&m.Type, &m.ID, &m.Phone, &m.Name, &m.Message, &m.Created, &r, &m.SenderName, &g, &m.Channel); err != nil { return nil, err }
		m.IsRead = r == 1
		m.IsGroup = g == 1
		out = append(out, m)
	}
	return out, nil
}

func (d *DB) MarkRead(phone string) error {
	_, err := d.sql.Exec(`UPDATE received SET is_read=1 WHERE phone=? AND is_read=0`, phone)
	return err
}

func (d *DB) UnreadCount() int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM received WHERE is_read=0`).Scan(&n)
	return n
}

func (d *DB) SearchInbox(query string) ([]InboxConversation, error) {
	q := "%" + query + "%"
	rows, err := d.sql.Query(`SELECT r.phone, COALESCE(g.name, MAX(r.name)) as name, MAX(r.is_group) as is_group, COUNT(CASE WHEN r.is_read=0 THEN 1 END) as unread, MAX(r.channel) as channel FROM received r LEFT JOIN wa_groups g ON r.phone=g.jid WHERE r.phone LIKE ? OR r.name LIKE ? OR r.message LIKE ? OR g.name LIKE ? OR r.sender_name LIKE ? GROUP BY r.phone ORDER BY MAX(r.id) DESC LIMIT 50`, q, q, q, q, q)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []InboxConversation
	for rows.Next() {
		var c InboxConversation; var g int
		if err := rows.Scan(&c.Phone, &c.Name, &g, &c.Unread, &c.Channel); err != nil { return nil, err }
		c.IsGroup = g == 1
		d.sql.QueryRow(`SELECT message, created_at FROM received WHERE phone=? ORDER BY id DESC LIMIT 1`, c.Phone).Scan(&c.LastMsg, &c.LastTime)
		out = append(out, c)
	}
	return out, nil
}

func (d *DB) GetGroupName(jid string) string {
	var name string
	d.sql.QueryRow(`SELECT name FROM wa_groups WHERE jid=?`, jid).Scan(&name)
	return name
}

func (d *DB) SaveGroupName(jid, name string) {
	d.sql.Exec(`INSERT INTO wa_groups (jid, name, updated_at) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE name=VALUES(name), updated_at=NOW()`, jid, name)
}

// Settings
func (d *DB) GetSetting(name, def string) string {
	var v string
	if err := d.sql.QueryRow(`SELECT value FROM settings WHERE name=?`, name).Scan(&v); err != nil { return def }
	return v
}
func (d *DB) SetSetting(name, value string) error {
	_, err := d.sql.Exec(`INSERT INTO settings (name,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=VALUES(value)`, name, value)
	return err
}

func (d *DB) SaveSessionOwner(phone string, userID int64) {
	d.sql.Exec(`INSERT INTO wa_session_owners (phone, user_id) VALUES (?,?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id)`, phone, userID)
}

func (d *DB) GetSessionOwner(phone string) int64 {
	var uid int64
	d.sql.QueryRow(`SELECT user_id FROM wa_session_owners WHERE phone=?`, phone).Scan(&uid)
	return uid
}

func (d *DB) SaveSession(token string, userID int64) {
	d.sql.Exec(`INSERT INTO sessions (token, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE user_id=VALUES(user_id), created_at=NOW()`, token, userID)
}
func (d *DB) GetSession(token string) (int64, bool) {
	var uid int64
	err := d.sql.QueryRow(`SELECT user_id FROM sessions WHERE token=?`, token).Scan(&uid)
	if err != nil { return 0, false }
	return uid, true
}
func (d *DB) DeleteSession(token string) {
	d.sql.Exec(`DELETE FROM sessions WHERE token=?`, token)
}

func (d *DB) ToggleStar(typ string, id int64) {
	if typ == "sent" {
		d.sql.Exec(`UPDATE sent SET starred=1-starred WHERE id=?`, id)
	} else {
		d.sql.Exec(`UPDATE received SET starred=1-starred WHERE id=?`, id)
	}
}

func extractCleanPhone(jid string) string {
	if len(jid) <= 16 { return jid }
	for cut := 14; cut >= 10; cut-- {
		if cut <= len(jid) {
			return jid[:cut]
		}
	}
	return jid
}

func (d *DB) AddFAQ(userID int64, question, answer string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO faq (user_id, question, answer) VALUES (?,?,?)`, userID, question, answer)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteFAQ(userID int64, id int64) error {
	if userID == 0 { _, err := d.sql.Exec(`DELETE FROM faq WHERE id=?`, id); return err }
	_, err := d.sql.Exec(`DELETE FROM faq WHERE id=? AND user_id=?`, id, userID)
	return err
}
func (d *DB) ListFAQ(userID int64) ([]map[string]string, error) {
	var rows *sql.Rows; var err error
	if userID == 0 { rows, err = d.sql.Query(`SELECT id, question, answer FROM faq ORDER BY id DESC`)
	} else { rows, err = d.sql.Query(`SELECT id, question, answer FROM faq WHERE user_id=? ORDER BY id DESC`, userID) }
	if err != nil { return nil, err }
	defer rows.Close()
	var out []map[string]string
	for rows.Next() { var id int64; var q, a string; rows.Scan(&id, &q, &a); out = append(out, map[string]string{"id": fmt.Sprintf("%d", id), "question": q, "answer": a}) }
	return out, nil
}
func (d *DB) FindFAQAnswer(userID int64, incoming string) (string, bool) {
	faqs, _ := d.ListFAQ(userID)
	msg := strings.ToLower(strings.TrimSpace(incoming))
	for _, f := range faqs {
		kw := strings.ToLower(strings.TrimSpace(f["question"]))
		if kw == "" { continue }
		for _, k := range strings.Split(kw, ",") { if strings.Contains(msg, strings.TrimSpace(k)) { return f["answer"], true } }
	}
	return "", false
}

// Welcome tracking
func (d *DB) MarkWelcomed(phone string) bool {
	var lastWelcomed string
	err := d.sql.QueryRow(`SELECT last_welcomed FROM welcomed WHERE phone=?`, phone).Scan(&lastWelcomed)
	if err == nil {
		if t, err := time.Parse("2006-01-02 15:04:05", lastWelcomed); err == nil {
			if time.Since(t) < 24*time.Hour {
				return false
			}
		}
		d.sql.Exec(`UPDATE welcomed SET last_welcomed=NOW() WHERE phone=?`, phone)
		return true
	}
	d.sql.Exec(`INSERT INTO welcomed (phone, last_welcomed) VALUES (?, NOW())`, phone)
	return true
}
