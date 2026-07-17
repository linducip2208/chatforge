package store

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
	dsn := os.Getenv("CHATGO_MYSQL")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/chatgo?charset=utf8mb4"
	}
	// parse DSN: user:pass@tcp(host:port)/dbname?params
	user := "root"
	pass := ""
	host := "127.0.0.1"
	port := "3306"
	dbname := "chatgo"

	if idx := strings.Index(dsn, ":"); idx >= 0 {
		user = dsn[:idx]
		rest := dsn[idx+1:]
		if idx2 := strings.Index(rest, "@"); idx2 >= 0 {
			pass = rest[:idx2]
			rest = rest[idx2+1:]
			if strings.HasPrefix(rest, "tcp(") {
				rest = rest[4:]
				if idx3 := strings.Index(rest, ")"); idx3 >= 0 {
					hp := rest[:idx3]
					rest = rest[idx3+1:]
					if idx4 := strings.Index(hp, ":"); idx4 >= 0 {
						host = hp[:idx4]
						port = hp[idx4+1:]
					} else {
						host = hp
					}
				}
			}
			if strings.HasPrefix(rest, "/") {
				rest = rest[1:]
				if idx5 := strings.Index(rest, "?"); idx5 >= 0 {
					dbname = rest[:idx5]
				} else {
					dbname = rest
				}
			}
		}
	}

	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command("mysqldump",
		"-u"+user,
		"-p"+pass,
		"-h"+host,
		"-P"+port,
		"--single-transaction",
		"--routines",
		"--triggers",
		dbname,
	)
	cmd.Stdout = f
	cmd.Stderr = f
	return cmd.Run()
}
