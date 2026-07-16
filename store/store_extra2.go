package store

type CustomerProfile struct {
	Phone       string
	Name        string
	TotalOrders int
	TotalSpent  float64
	LastActive  string
	CSATAvg     float64
	ChatCount   int
	Tags        string
}

func (d *DB) GetCustomerProfile(phone string) *CustomerProfile {
	p := &CustomerProfile{Phone: phone}
	d.sql.QueryRow(`SELECT IFNULL(MAX(name),'') FROM received WHERE phone=? ORDER BY id DESC LIMIT 1`, phone).Scan(&p.Name)
	d.sql.QueryRow(`SELECT COUNT(*), IFNULL(SUM(total),0) FROM store_orders WHERE phone=?`, phone).Scan(&p.TotalOrders, &p.TotalSpent)
	d.sql.QueryRow(`SELECT IFNULL(MAX(created_at),'') FROM received WHERE phone=?`, phone).Scan(&p.LastActive)
	d.sql.QueryRow(`SELECT IFNULL(AVG(rating),0) FROM csat_ratings WHERE phone=?`, phone).Scan(&p.CSATAvg)
	d.sql.QueryRow(`SELECT COUNT(*) FROM received WHERE phone=?`, phone).Scan(&p.ChatCount)
	rows, _ := d.sql.Query(`SELECT GROUP_CONCAT(t.name) FROM contact_tags ct JOIN tags t ON t.id=ct.tag_id WHERE ct.contact_id=(SELECT id FROM contacts WHERE phone=? LIMIT 1)`, phone)
	if rows != nil {
		defer rows.Close()
		if rows.Next() { rows.Scan(&p.Tags) }
	}
	return p
}

// Calendar events across campaigns, drips, recurring
type CalendarEvent struct {
	Date  string
	Title string
	Type  string
}

func (d *DB) GetCalendarEvents() []CalendarEvent {
	var out []CalendarEvent
	// campaigns
	rows, _ := d.sql.Query(`SELECT created_at, name FROM campaigns WHERE status IN ('running','pending') ORDER BY created_at`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var dt, name string
			rows.Scan(&dt, &name)
			out = append(out, CalendarEvent{Date: dt[:10], Title: name, Type: "Campaign"})
		}
	}
	// recurring
	rows2, _ := d.sql.Query(`SELECT created_at, name FROM recurring_campaigns WHERE status='active' ORDER BY created_at`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var dt, name string
			rows2.Scan(&dt, &name)
			out = append(out, CalendarEvent{Date: dt[:10], Title: name, Type: "Recurring"})
		}
	}
	// reminders
	rows3, _ := d.sql.Query(`SELECT due_date, CONCAT(name,' - Rp',amount) FROM payment_reminders WHERE status='pending' ORDER BY due_date`)
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var dt, name string
			rows3.Scan(&dt, &name)
			out = append(out, CalendarEvent{Date: dt, Title: name, Type: "Reminder"})
		}
	}
	return out
}

// Backup
func (d *DB) BackupDB(outPath string) error {
	return execBackup(outPath)
}

var execBackup = func(path string) error {
	// Placeholder - actual implementation uses mysqldump
	return nil
}
