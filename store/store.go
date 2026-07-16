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

type AutoReply struct {
	ID        int64
	Keyword   string
	Match     string
	Reply     string
	IsActive  bool
	UseAI     bool
	AiKeyID    int64
	AccountID  string
	TrainingID int64
	Created    string
}

type SentMessage struct {
	ID      int64
	Phone   string
	Message string
	Status  string
	Created string
}

type ReceivedMessage struct {
	ID      int64
	Phone   string
	Name    string
	Message string
	IsGroup bool
	Created string
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
	return nil
}

// AutoReply CRUD
func (d *DB) AddAutoReply(keyword, match, reply string, useAI bool, aiKeyID int64, accountID string, trainingID int64) (int64, error) {
	ai := 0; if useAI { ai = 1 }
	res, err := d.sql.Exec(`INSERT INTO autoreplies (keyword,match_type,reply,is_active,use_ai,ai_key_id,account_id,training_id) VALUES (?,?,?,1,?,?,?,?)`, keyword, match, reply, ai, aiKeyID, accountID, trainingID)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteAutoReply(id int64) error { _, err := d.sql.Exec(`DELETE FROM autoreplies WHERE id=?`, id); return err }
func (d *DB) ToggleAutoReply(id int64) error { _, err := d.sql.Exec(`UPDATE autoreplies SET is_active=1-is_active WHERE id=?`, id); return err }
func (d *DB) GetAutoReply(id int64) (*AutoReply, error) {
	var a AutoReply; var active int
	err := d.sql.QueryRow(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),created_at FROM autoreplies WHERE id=?`, id).Scan(&a.ID, &a.Keyword, &a.Match, &a.Reply, &active, &a.UseAI, &a.AiKeyID, &a.AccountID, &a.TrainingID, &a.Created)
	a.IsActive = active == 1
	if err != nil { return nil, err }
	return &a, nil
}
func (d *DB) UpdateAutoReply(id int64, keyword, match, reply string, useAI bool, aiKeyID int64, accountID string, trainingID int64) error {
	use := 0
	if useAI { use = 1 }
	_, err := d.sql.Exec(`UPDATE autoreplies SET keyword=?, match_type=?, reply=?, use_ai=?, ai_key_id=?, account_id=?, training_id=? WHERE id=?`, keyword, match, reply, use, aiKeyID, accountID, trainingID, id)
	return err
}
func (d *DB) ListAutoReplies() ([]AutoReply, error) {
	rows, err := d.sql.Query(`SELECT id,keyword,match_type,reply,is_active,use_ai,ai_key_id,IFNULL(account_id,''),IFNULL(training_id,0),created_at FROM autoreplies ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []AutoReply
	for rows.Next() {
		var a AutoReply; var active int
		if err := rows.Scan(&a.ID,&a.Keyword,&a.Match,&a.Reply,&active,&a.UseAI,&a.AiKeyID,&a.AccountID,&a.TrainingID,&a.Created); err != nil { return nil, err }
		a.IsActive = active==1
		out = append(out, a)
	}
	return out, nil
}
func (d *DB) FindReply(incoming string) (string, bool) {
	r, ok := d.FindReplyFull(incoming)
	return r.Reply, ok
}
func (d *DB) FindReplyFullForAccount(incoming string, accountPhone string) (AutoReply, bool) {
	rules, _ := d.ListAutoReplies()
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
func (d *DB) FindReplyFull(incoming string) (AutoReply, bool) {
	rules, _ := d.ListAutoReplies()
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
func (d *DB) LogSent(phone, message, status string) { d.sql.Exec(`INSERT INTO sent (phone,message,status) VALUES (?,?,?)`, phone, message, status) }
func (d *DB) ListSent(limit int) ([]SentMessage, error) { return d.ListSentPaginated(1, limit) }
func (d *DB) ListSentPaginated(page, perPage int) ([]SentMessage, error) {
	offset := (page-1)*perPage
	rows, err := d.sql.Query(`SELECT id,phone,message,status,created_at FROM sent ORDER BY id DESC LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []SentMessage
	for rows.Next() { var m SentMessage; rows.Scan(&m.ID,&m.Phone,&m.Message,&m.Status,&m.Created); out = append(out, m) }
	return out, nil
}
func (d *DB) CountSent() int { var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM sent`).Scan(&n); return n }

// Received log
func (d *DB) LogReceived(phone, name, message string, isGroup bool) {
	g := 0; if isGroup { g = 1 }
	d.sql.Exec(`INSERT INTO received (phone,name,message,is_group) VALUES (?,?,?,?)`, phone, name, message, g)
}
func (d *DB) ListReceived(limit int) ([]ReceivedMessage, error) { return d.ListReceivedPaginated(1, limit) }
func (d *DB) ListReceivedPaginated(page, perPage int) ([]ReceivedMessage, error) {
	offset := (page-1)*perPage
	rows, err := d.sql.Query(`SELECT id,phone,name,message,is_group,created_at FROM received ORDER BY id DESC LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ReceivedMessage
	for rows.Next() { var m ReceivedMessage; var g int; rows.Scan(&m.ID,&m.Phone,&m.Name,&m.Message,&g,&m.Created); m.IsGroup = g==1; out = append(out, m) }
	return out, nil
}
func (d *DB) CountReceived() int { var n int; d.sql.QueryRow(`SELECT COUNT(*) FROM received`).Scan(&n); return n }

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
