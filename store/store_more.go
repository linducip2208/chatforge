package store

// InstanceLog records session connect/disconnect events.
type InstanceLog struct {
	ID      int64
	Phone   string
	Event   string // "connected" | "disconnected" | "logged_out" | "heartbeat_fail"
	Created string
}

func (d *DB) migrateInstanceLog() error {
	_, err := d.sql.Exec(`CREATE TABLE IF NOT EXISTS instance_log (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		phone VARCHAR(64) NOT NULL DEFAULT '',
		event VARCHAR(32) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	return err
}

func (d *DB) LogInstance(phone, event string) {
	d.sql.Exec(`INSERT INTO instance_log (phone, event) VALUES (?, ?)`, phone, event)
}

func (d *DB) UptimeMinutes(phone string) int {
	var total *int
	rows, _ := d.sql.Query(`SELECT TIMESTAMPDIFF(MINUTE, created_at, LEAD(created_at) OVER (ORDER BY id)) FROM instance_log WHERE phone=? AND event='connected'`, phone)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var m int
			rows.Scan(&m)
			if total == nil { t := 0; total = &t }
			*total += m
		}
	}
	if total == nil { return 0 }
	return *total
}

type WAStatus struct {
	ID       int64
	Phone    string
	Name     string
	Caption  string
	MediaURL string
	Created  string
}

func (d *DB) migrateStatuses() error {
	_, err := d.sql.Exec(`CREATE TABLE IF NOT EXISTS wa_statuses (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		phone VARCHAR(64) NOT NULL,
		name VARCHAR(255) NOT NULL DEFAULT '',
		caption TEXT NOT NULL,
		media_url VARCHAR(512) NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	return err
}

func (d *DB) LogStatus(phone, name, caption, mediaURL string) {
	d.sql.Exec(`INSERT INTO wa_statuses (phone, name, caption, media_url) VALUES (?, ?, ?, ?)`, phone, name, caption, mediaURL)
}

func (d *DB) ListStatuses() ([]WAStatus, error) {
	rows, err := d.sql.Query(`SELECT phone, MAX(name) as name, MAX(caption) as caption, MAX(media_url) as media_url, MAX(created_at) as created_at FROM wa_statuses WHERE created_at > DATE_SUB(NOW(), INTERVAL 24 HOUR) GROUP BY phone ORDER BY MAX(created_at) DESC LIMIT 50`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []WAStatus
	for rows.Next() {
		var s WAStatus
		rows.Scan(&s.Phone, &s.Name, &s.Caption, &s.MediaURL, &s.Created)
		out = append(out, s)
	}
	return out, nil
}
