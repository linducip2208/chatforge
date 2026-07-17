package store

import "strings"

type BlacklistEntry struct {
	ID      int64
	Phone   string
	Reason  string
	Created string
}

type CSATResponse struct {
	ID      int64
	Phone   string
	Rating  int
	Comment string
	Created string
}

type ValidatedNumber struct {
	Phone    string
	Valid    bool
	OnWA     bool
	Reason   string
}

func (d *DB) migrateSafety() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS blacklist (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL UNIQUE, reason VARCHAR(255) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS csat_ratings (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, rating TINYINT NOT NULL, comment VARCHAR(500) DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, INDEX idx_phone (phone), INDEX idx_created (created_at)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS spam_log (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, message_hash VARCHAR(64) NOT NULL, count INT NOT NULL DEFAULT 1, first_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, INDEX idx_phone_hash (phone, message_hash)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`ALTER TABLE contact_groups ADD COLUMN IF NOT EXISTS language VARCHAR(10) NOT NULL DEFAULT ''`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	d.safeAddColumn("blacklist", "user_id", "BIGINT NOT NULL DEFAULT 0")
	return nil
}

// ---- Blacklist ----
func (d *DB) AddBlacklist(phone, reason string) error {
	_, err := d.sql.Exec(`INSERT IGNORE INTO blacklist (phone, reason) VALUES (?,?)`, phone, reason)
	return err
}
func (d *DB) RemoveBlacklist(phone string) error {
	_, err := d.sql.Exec(`DELETE FROM blacklist WHERE phone=?`, phone)
	return err
}
func (d *DB) IsBlacklisted(phone string) bool {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM blacklist WHERE phone=?`, phone).Scan(&n)
	return n > 0
}
func (d *DB) ListBlacklist() ([]BlacklistEntry, error) {
	rows, err := d.sql.Query(`SELECT id, phone, reason, created_at FROM blacklist ORDER BY created_at DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []BlacklistEntry
	for rows.Next() {
		var b BlacklistEntry
		rows.Scan(&b.ID, &b.Phone, &b.Reason, &b.Created)
		out = append(out, b)
	}
	return out, nil
}

// ---- Spam Detection ----
func (d *DB) TrackSpam(phone, msgHash string) bool {
	d.sql.Exec(`INSERT INTO spam_log (phone, message_hash, first_seen, last_seen) VALUES (?,?,NOW(),NOW()) ON DUPLICATE KEY UPDATE count=count+1, last_seen=NOW()`, phone, msgHash)
	var count int
	d.sql.QueryRow(`SELECT SUM(count) FROM spam_log WHERE phone=? AND last_seen > DATE_SUB(NOW(), INTERVAL 10 MINUTE)`, phone).Scan(&count)
	return count > 8 // threshold: 8+ identical/rapid messages in 10 min = spam
}

// ---- CSAT ----
func (d *DB) SaveCSAT(phone string, rating int, comment string) {
	d.sql.Exec(`INSERT INTO csat_ratings (phone, rating, comment) VALUES (?,?,?)`, phone, rating, comment)
}
func (d *DB) CSATAverage(days int) float64 {
	var avg float64
	d.sql.QueryRow(`SELECT IFNULL(AVG(rating),0) FROM csat_ratings WHERE created_at > DATE_SUB(NOW(), INTERVAL ? DAY)`, days).Scan(&avg)
	return avg
}
func (d *DB) CSATCount() int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM csat_ratings`).Scan(&n)
	return n
}

// ---- Number Validator ----
var ValidFormat = func(phone string) bool {
	p := strings.TrimPrefix(strings.TrimPrefix(strings.TrimSpace(phone), "+"), "0")
	if len(p) < 7 || len(p) > 16 { return false }
	for _, c := range p {
		if c < '0' || c > '9' { return false }
	}
	return true
}

func (d *DB) ValidateNumbers(phones []string) []ValidatedNumber {
	var out []ValidatedNumber
	for _, p := range phones {
		v := ValidatedNumber{Phone: p}
		if !ValidFormat(p) {
			v.Reason = "invalid format"
			out = append(out, v)
			continue
		}
		if d.IsBlacklisted(p) {
			v.Reason = "blacklisted"
			out = append(out, v)
			continue
		}
		if d.IsUnsub(p) {
			v.Reason = "unsubscribed"
			out = append(out, v)
			continue
		}
		v.Valid = true
		out = append(out, v)
	}
	return out
}

// ---- Multi-Lang Group ----
func (d *DB) SetGroupLanguage(groupID int64, lang string) error {
	_, err := d.sql.Exec(`UPDATE contact_groups SET language=? WHERE id=?`, lang, groupID)
	return err
}
func (d *DB) GetGroupLanguage(groupID int64) string {
	var lang string
	d.sql.QueryRow(`SELECT language FROM contact_groups WHERE id=?`, groupID).Scan(&lang)
	return lang
}
func (d *DB) ListGroupsWithLang() ([]Group, error) {
	rows, err := d.sql.Query(`SELECT g.id, g.name, COUNT(c.id), IFNULL(g.language,''), g.created_at FROM contact_groups g LEFT JOIN contacts c ON FIND_IN_SET(g.id, REPLACE(c.groups, ' ', '')) GROUP BY g.id ORDER BY g.id DESC`)
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
