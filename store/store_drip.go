package store

import "time"

type Drip struct {
	ID        int64
	Name      string
	Status    string
	Created   string
	Steps     []DripStep
}

type DripStep struct {
	ID           int64
	DripID       int64
	DelayMinutes int
	Message      string
	SortOrder    int
}

type DripEnrollment struct {
	ID          int64
	DripID      int64
	Phone       string
	Name        string
	CurrentStep int
	NextSendAt  string
	Status      string
}

func (d *DB) migrateDrip() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS drips (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS drip_steps (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			drip_id BIGINT NOT NULL,
			delay_minutes INT NOT NULL DEFAULT 0,
			message TEXT NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			FOREIGN KEY (drip_id) REFERENCES drips(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS drip_enrollments (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			drip_id BIGINT NOT NULL,
			phone VARCHAR(64) NOT NULL,
			name VARCHAR(255) NOT NULL DEFAULT '',
			current_step INT NOT NULL DEFAULT 0,
			next_send_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_drip_phone (drip_id, phone),
			FOREIGN KEY (drip_id) REFERENCES drips(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Drips ----
func (d *DB) AddDrip(name string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO drips (name) VALUES (?)`, name)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateDripStatus(id int64, status string) error {
	_, err := d.sql.Exec(`UPDATE drips SET status=? WHERE id=?`, status, id)
	return err
}
func (d *DB) DeleteDrip(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM drips WHERE id=?`, id)
	return err
}
func (d *DB) ListDrips() ([]Drip, error) {
	rows, err := d.sql.Query(`SELECT d.id, d.name, d.status, d.created_at, IFNULL(s.id,0), IFNULL(s.delay_minutes,0), IFNULL(s.message,''), IFNULL(s.sort_order,0) FROM drips d LEFT JOIN drip_steps s ON s.drip_id=d.id ORDER BY d.id DESC, s.sort_order`)
	if err != nil { return nil, err }
	defer rows.Close()
	seen := map[int64]*Drip{}
	var order []int64
	for rows.Next() {
		var sid int64
		var delay, sortOrd int
		var msg string
		var dr Drip
		rows.Scan(&dr.ID, &dr.Name, &dr.Status, &dr.Created, &sid, &delay, &msg, &sortOrd)
		if existing, ok := seen[dr.ID]; ok {
			if sid > 0 {
				existing.Steps = append(existing.Steps, DripStep{ID: sid, DripID: dr.ID, DelayMinutes: delay, Message: msg, SortOrder: sortOrd})
			}
		} else {
			dr.Steps = []DripStep{}
			if sid > 0 {
				dr.Steps = append(dr.Steps, DripStep{ID: sid, DripID: dr.ID, DelayMinutes: delay, Message: msg, SortOrder: sortOrd})
			}
			seen[dr.ID] = &dr
			order = append(order, dr.ID)
		}
	}
	out := make([]Drip, 0, len(order))
	for _, id := range order { out = append(out, *seen[id]) }
	return out, nil
}
func (d *DB) GetDrip(id int64) (*Drip, error) {
	rows, err := d.sql.Query(`SELECT d.id, d.name, d.status, d.created_at, s.id, s.delay_minutes, s.message, s.sort_order FROM drips d LEFT JOIN drip_steps s ON s.drip_id=d.id WHERE d.id=? ORDER BY s.sort_order`, id)
	if err != nil { return nil, err }
	defer rows.Close()
	var dr *Drip
	for rows.Next() {
		var sid int64
		var delay int
		var msg string
		var sortOrd int
		if dr == nil {
			dr = &Drip{}
			rows.Scan(&dr.ID, &dr.Name, &dr.Status, &dr.Created, &sid, &delay, &msg, &sortOrd)
		} else {
			rows.Scan(new(int64), new(string), new(string), new(string), &sid, &delay, &msg, &sortOrd)
		}
		if sid > 0 {
			dr.Steps = append(dr.Steps, DripStep{ID: sid, DripID: dr.ID, DelayMinutes: delay, Message: msg, SortOrder: sortOrd})
		}
	}
	return dr, nil
}

// ---- Drip Steps ----
func (d *DB) AddDripStep(dripID int64, delayMinutes int, message string, sortOrder int) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO drip_steps (drip_id, delay_minutes, message, sort_order) VALUES (?, ?, ?, ?)`, dripID, delayMinutes, message, sortOrder)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteDripStep(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM drip_steps WHERE id=?`, id)
	return err
}

// ---- Enrollments ----
func (d *DB) EnrollInDrip(dripID int64, phone, name string) error {
	_, err := d.sql.Exec(`INSERT IGNORE INTO drip_enrollments (drip_id, phone, name) VALUES (?, ?, ?)`, dripID, phone, name)
	return err
}
func (d *DB) UnenrollFromDrip(phone string) error {
	_, err := d.sql.Exec(`UPDATE drip_enrollments SET status='stopped' WHERE phone=? AND status='active'`, phone)
	return err
}
func (d *DB) DueDripEnrollments() ([]DripEnrollment, error) {
	rows, err := d.sql.Query(`SELECT e.id, e.drip_id, e.phone, e.name, e.current_step, e.next_send_at, e.status FROM drip_enrollments e JOIN drips d ON d.id=e.drip_id WHERE e.status='active' AND d.status='active' AND e.next_send_at <= NOW() LIMIT 100`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []DripEnrollment
	for rows.Next() {
		var en DripEnrollment
		rows.Scan(&en.ID, &en.DripID, &en.Phone, &en.Name, &en.CurrentStep, &en.NextSendAt, &en.Status)
		out = append(out, en)
	}
	return out, nil
}
func (d *DB) AdvanceDripStep(enrollmentID int64, nextDelayMinutes int) error {
	if nextDelayMinutes > 0 {
		next := time.Now().Add(time.Duration(nextDelayMinutes) * time.Minute).Format("2006-01-02 15:04:05")
		_, err := d.sql.Exec(`UPDATE drip_enrollments SET current_step=current_step+1, next_send_at=? WHERE id=?`, next, enrollmentID)
		return err
	}
	_, err := d.sql.Exec(`UPDATE drip_enrollments SET status='completed', current_step=current_step+1 WHERE id=?`, enrollmentID)
	return err
}
func (d *DB) CountDripEnrollments(dripID int64) int {
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM drip_enrollments WHERE drip_id=? AND status='active'`, dripID).Scan(&n)
	return n
}

var _ = time.Now // ensure time import used
