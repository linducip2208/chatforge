package store

import (
	"database/sql"
	"fmt"
	"os"
)

type ChatNote struct {
	ID       int64
	Phone    string
	AgentID  int64
	Note     string
	Created  string
}

type Department struct {
	ID      int64
	Name    string
	Agents  string // comma-separated user IDs
	Created string
}

type RecurringCampaign struct {
	ID         int64
	Name       string
	Groups     string
	Message    string
	DayOfWeek  int // 0=sunday, 1=monday... or 0=daily
	Hour       int // 0-23
	Status     string
	LastRun    string
	Created    string
}

func (d *DB) migratePlus() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS conversation_notes (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, agent_id BIGINT NOT NULL, note TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, INDEX idx_phone (phone)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS departments (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, agents TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS recurring_campaigns (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, ` + "`groups`" + ` VARCHAR(512) NOT NULL DEFAULT '', message TEXT NOT NULL, day_of_week INT NOT NULL DEFAULT 0, hour INT NOT NULL DEFAULT 9, status VARCHAR(20) NOT NULL DEFAULT 'active', last_run DATETIME NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`ALTER TABLE agent_assignments ADD COLUMN IF NOT EXISTS dept VARCHAR(100) NOT NULL DEFAULT ''`,
		`ALTER TABLE agent_assignments ADD COLUMN IF NOT EXISTS last_activity DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`,
		`ALTER TABLE agent_assignments ADD COLUMN IF NOT EXISTS label VARCHAR(50) NOT NULL DEFAULT ''`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Labels ----
func (d *DB) SetConversationLabel(phone, label string) error {
	_, err := d.sql.Exec(`UPDATE agent_assignments SET label=? WHERE phone=?`, label, phone)
	return err
}
func (d *DB) GetConversationLabel(phone string) string {
	var l string
	d.sql.QueryRow(`SELECT label FROM agent_assignments WHERE phone=?`, phone).Scan(&l)
	return l
}
func (d *DB) InboxFiltered(userID int64, filter string) []map[string]string {
	where := ""
	switch filter {
	case "unread":
		where = " AND r.is_read=0"
	case "today":
		where = " AND DATE(r.created_at)=CURDATE()"
	case "week":
		where = " AND r.created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)"
	}
	var rows *sql.Rows
	var err error
	if userID == 0 {
		rows, err = d.sql.Query(`SELECT r.phone, COALESCE(g.name, MAX(r.name)) as name, MAX(r.is_group) as is_group, COUNT(CASE WHEN r.is_read=0 THEN 1 END) as unread, MAX(r.channel) as channel, IFNULL(a.label,'') as label FROM received r LEFT JOIN wa_groups g ON r.phone = g.jid LEFT JOIN agent_assignments a ON a.phone=r.phone WHERE 1=1` + where + ` GROUP BY r.phone ORDER BY MAX(r.id) DESC LIMIT 50`)
	} else {
		rows, err = d.sql.Query(`SELECT r.phone, COALESCE(g.name, MAX(r.name)) as name, MAX(r.is_group) as is_group, COUNT(CASE WHEN r.is_read=0 THEN 1 END) as unread, MAX(r.channel) as channel, IFNULL(a.label,'') as label FROM received r INNER JOIN wa_session_owners o ON r.wa_phone = o.phone AND o.user_id = ? LEFT JOIN wa_groups g ON r.phone = g.jid LEFT JOIN agent_assignments a ON a.phone=r.phone WHERE 1=1` + where + ` GROUP BY r.phone ORDER BY MAX(r.id) DESC LIMIT 50`, userID)
	}
	if err != nil { return nil }
	defer rows.Close()
	var out []map[string]string
	for rows.Next() {
		var phone, name, channel, label string
		var isGroup, unread int
		rows.Scan(&phone, &name, &isGroup, &unread, &channel, &label)
		out = append(out, map[string]string{"phone": phone, "name": name, "channel": channel, "unread": fmt.Sprintf("%d", unread), "label": label})
	}
	return out
}

// ---- Chat Notes ----
func (d *DB) AddNote(phone string, agentID int64, note string) error {
	_, err := d.sql.Exec(`INSERT INTO conversation_notes (phone, agent_id, note) VALUES (?,?,?)`, phone, agentID, note)
	return err
}
func (d *DB) GetNotes(phone string) ([]ChatNote, error) {
	rows, err := d.sql.Query(`SELECT id, phone, agent_id, note, created_at FROM conversation_notes WHERE phone=? ORDER BY id DESC LIMIT 50`, phone)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ChatNote
	for rows.Next() {
		var n ChatNote
		rows.Scan(&n.ID, &n.Phone, &n.AgentID, &n.Note, &n.Created)
		out = append(out, n)
	}
	return out, nil
}

// ---- Departments ----
func (d *DB) AddDept(name, agents string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO departments (name, agents) VALUES (?,?)`, name, agents)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteDept(id int64) error { _, err := d.sql.Exec(`DELETE FROM departments WHERE id=?`, id); return err }
func (d *DB) ListDepts() ([]Department, error) {
	rows, err := d.sql.Query(`SELECT id, name, agents, created_at FROM departments ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Department
	for rows.Next() {
		var de Department
		rows.Scan(&de.ID, &de.Name, &de.Agents, &de.Created)
		out = append(out, de)
	}
	return out, nil
}
func (d *DB) GetDept(name string) (*Department, error) {
	var de Department
	err := d.sql.QueryRow(`SELECT id, name, agents, created_at FROM departments WHERE name=?`, name).Scan(&de.ID, &de.Name, &de.Agents, &de.Created)
	if err != nil { return nil, err }
	return &de, nil
}
func (d *DB) AssignToDept(phone, dept string) error {
	_, err := d.sql.Exec(`UPDATE agent_assignments SET dept=? WHERE phone=?`, dept, phone)
	return err
}

// ---- Recurring Campaigns ----
func (d *DB) AddRecurring(name, groups, message string, dayOfWeek, hour int) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO recurring_campaigns (name, `+"`groups`"+`, message, day_of_week, hour) VALUES (?,?,?,?,?)`, name, groups, message, dayOfWeek, hour)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteRecurring(id int64) error { _, err := d.sql.Exec(`DELETE FROM recurring_campaigns WHERE id=?`, id); return err }
func (d *DB) ToggleRecurring(id int64) error {
	_, err := d.sql.Exec(`UPDATE recurring_campaigns SET status=IF(status='active','inactive','active') WHERE id=?`, id)
	return err
}
func (d *DB) ListRecurring() ([]RecurringCampaign, error) {
	rows, err := d.sql.Query(`SELECT id, name, `+"`groups`"+`, message, day_of_week, hour, status, IFNULL(last_run,''), created_at FROM recurring_campaigns ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []RecurringCampaign
	for rows.Next() {
		var r RecurringCampaign
		rows.Scan(&r.ID, &r.Name, &r.Groups, &r.Message, &r.DayOfWeek, &r.Hour, &r.Status, &r.LastRun, &r.Created)
		out = append(out, r)
	}
	return out, nil
}
func (d *DB) DueRecurring() ([]RecurringCampaign, error) {
	rows, err := d.sql.Query(`SELECT id, name, `+"`groups`"+`, message, day_of_week, hour, status, IFNULL(last_run,''), created_at FROM recurring_campaigns WHERE status='active' AND (last_run IS NULL OR last_run < DATE_SUB(NOW(), INTERVAL 23 HOUR))`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []RecurringCampaign
	for rows.Next() {
		var r RecurringCampaign
		rows.Scan(&r.ID, &r.Name, &r.Groups, &r.Message, &r.DayOfWeek, &r.Hour, &r.Status, &r.LastRun, &r.Created)
		out = append(out, r)
	}
	return out, nil
}
func (d *DB) MarkRecurringRun(id int64) error {
	_, err := d.sql.Exec(`UPDATE recurring_campaigns SET last_run=NOW() WHERE id=?`, id)
	return err
}

// ---- Auto-close ----
func (d *DB) IdleConversations(hours int) []string {
	rows, err := d.sql.Query(`SELECT phone FROM agent_assignments WHERE status='open' AND last_activity < DATE_SUB(NOW(), INTERVAL ? HOUR)`, hours)
	if err != nil { return nil }
	defer rows.Close()
	var phones []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		phones = append(phones, p)
	}
	return phones
}

// ---- File Browser ----
func (d *DB) ListUploads(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil { return nil }
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out
}
