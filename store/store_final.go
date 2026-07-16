package store

import (
	"fmt"
	"time"
)

type InboxMacro struct {
	ID        int64
	Name      string
	Actions   string // JSON: [{type:"assign",value:"agent_id"},{type:"tag",value:"label"},{type:"reply",value:"message"},{type:"close"}]
	Created   string
}

type WebhookRetry struct {
	ID         int64
	WebhookID  int64
	URL        string
	Event      string
	Payload    string
	Attempts   int
	MaxAttempts int
	NextRetry  string
	Status     string
	Created    string
}

type AuditLog struct {
	ID        int64
	UserID    int64
	Action    string
	Detail    string
	IP        string
	Created   string
}

func (d *DB) migrateFinal() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS inbox_macros (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, actions TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS webhook_retries (id BIGINT AUTO_INCREMENT PRIMARY KEY, webhook_id BIGINT NOT NULL DEFAULT 0, url VARCHAR(1024) NOT NULL, event VARCHAR(50) NOT NULL, payload TEXT NOT NULL, attempts INT NOT NULL DEFAULT 0, max_attempts INT NOT NULL DEFAULT 5, next_retry DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS audit_logs (id BIGINT AUTO_INCREMENT PRIMARY KEY, user_id BIGINT NOT NULL DEFAULT 0, action VARCHAR(50) NOT NULL, detail TEXT NOT NULL, ip VARCHAR(45) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`ALTER TABLE agent_assignments ADD COLUMN IF NOT EXISTS priority INT NOT NULL DEFAULT 0`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Inbox Macros ----
func (d *DB) AddMacro(name, actions string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO inbox_macros (name, actions) VALUES (?,?)`, name, actions)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteMacro(id int64) error { _, err := d.sql.Exec(`DELETE FROM inbox_macros WHERE id=?`, id); return err }
func (d *DB) ListMacros() ([]InboxMacro, error) {
	rows, err := d.sql.Query(`SELECT id, name, actions, created_at FROM inbox_macros ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []InboxMacro
	for rows.Next() {
		var m InboxMacro
		rows.Scan(&m.ID, &m.Name, &m.Actions, &m.Created)
		out = append(out, m)
	}
	return out, nil
}

// ---- Webhook Retry ----
func (d *DB) QueueWebhookRetry(url, event, payload string) {
	d.sql.Exec(`INSERT INTO webhook_retries (url, event, payload) VALUES (?,?,?)`, url, event, payload)
}
func (d *DB) DueRetries() ([]WebhookRetry, error) {
	rows, err := d.sql.Query(`SELECT id, webhook_id, url, event, payload, attempts, max_attempts, next_retry, status, created_at FROM webhook_retries WHERE status='pending' AND next_retry <= NOW() AND attempts < max_attempts LIMIT 20`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []WebhookRetry
	for rows.Next() {
		var r WebhookRetry
		rows.Scan(&r.ID, &r.WebhookID, &r.URL, &r.Event, &r.Payload, &r.Attempts, &r.MaxAttempts, &r.NextRetry, &r.Status, &r.Created)
		out = append(out, r)
	}
	return out, nil
}
func (d *DB) UpdateRetry(id int64, success bool) {
	if success {
		d.sql.Exec(`UPDATE webhook_retries SET status='done' WHERE id=?`, id)
	} else {
		var att int
		d.sql.QueryRow(`SELECT attempts FROM webhook_retries WHERE id=?`, id).Scan(&att)
		delay := 30 + 15*att
		d.sql.Exec(`UPDATE webhook_retries SET attempts=attempts+1, next_retry=DATE_ADD(NOW(), INTERVAL ? SECOND) WHERE id=?`, delay, id)
	}
}

// ---- Contact Merge ----
func (d *DB) FindDuplicateContacts() []map[string]interface{} {
	rows, err := d.sql.Query(`SELECT phone, COUNT(id) as cnt, GROUP_CONCAT(id) as ids, GROUP_CONCAT(name) as names FROM contacts GROUP BY phone HAVING cnt > 1`)
	if err != nil { return nil }
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var phone string
		var cnt int
		var ids, names string
		rows.Scan(&phone, &cnt, &ids, &names)
		out = append(out, map[string]interface{}{"phone": phone, "cnt": cnt, "ids": ids, "names": names})
	}
	return out
}
func (d *DB) MergeContacts(keepID int64, mergeIDs []int64) error {
	for _, mid := range mergeIDs {
		d.sql.Exec(`DELETE FROM contacts WHERE id=?`, mid)
	}
	return nil
}

// ---- Priority ----
func (d *DB) SetContactPriority(phone string, priority int) error {
	_, err := d.sql.Exec(`UPDATE agent_assignments SET priority=? WHERE phone=?`, priority, phone)
	return err
}
func (d *DB) HighPriorityCount() int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM agent_assignments WHERE priority=1 AND status='open'`).Scan(&n)
	return n
}

// ---- Audit Log ----
func (d *DB) LogAudit(userID int64, action, detail, ip string) {
	d.sql.Exec(`INSERT INTO audit_logs (user_id, action, detail, ip) VALUES (?,?,?,?)`, userID, action, detail, ip)
}
func (d *DB) ListAuditLogs() ([]AuditLog, error) {
	rows, err := d.sql.Query(`SELECT id, user_id, action, detail, ip, created_at FROM audit_logs ORDER BY id DESC LIMIT 200`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []AuditLog
	for rows.Next() {
		var a AuditLog
		rows.Scan(&a.ID, &a.UserID, &a.Action, &a.Detail, &a.IP, &a.Created)
		out = append(out, a)
	}
	return out, nil
}

var _ = fmt.Sprintf
var _ = time.Now
