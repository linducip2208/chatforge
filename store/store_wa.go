package store

import "time"

// ---- Products ----
type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	ImageURL    string
	Category    string
	Stock       int
	Status      string
	Created     string
}

type Order struct {
	ID        int64
	Phone     string
	Name      string
	ProductID int64
	Quantity  int
	Total     float64
	Status    string
	Created   string
}

type ProductCategory struct {
	ID      int64
	Name    string
	Created string
}

type BehaviorEvent struct {
	ID        int64
	Phone     string
	Event     string
	Meta      string
	Created   string
}

type PaymentReminder struct {
	ID      int64
	Phone   string
	Name    string
	Amount  float64
	DueDate string
	Status  string
	Message string
	Created string
}

type ChatForm struct {
	ID      int64
	Name    string
	Fields  string
	Status  string
	Created string
}

type FormSubmission struct {
	ID       int64
	FormID   int64
	Phone    string
	Data     string
	Created  string
}

type AgentMetric struct {
	AgentName string
	Chats     int
	Replied   int
	AvgTime   float64
}

func (d *DB) migrateStore() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS product_categories (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS products (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, description TEXT, price DECIMAL(12,2) NOT NULL DEFAULT 0, image_url VARCHAR(1024) DEFAULT '', category VARCHAR(255) NOT NULL DEFAULT '', stock INT NOT NULL DEFAULT 0, status VARCHAR(20) NOT NULL DEFAULT 'active', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS store_orders (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, name VARCHAR(255) NOT NULL DEFAULT '', product_id BIGINT NOT NULL, quantity INT NOT NULL DEFAULT 1, total DECIMAL(12,2) NOT NULL DEFAULT 0, status VARCHAR(20) NOT NULL DEFAULT 'new', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS behavior_events (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, event VARCHAR(32) NOT NULL, meta VARCHAR(512) DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, INDEX idx_phone (phone), INDEX idx_event (event)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS payment_reminders (id BIGINT AUTO_INCREMENT PRIMARY KEY, phone VARCHAR(64) NOT NULL, name VARCHAR(255) DEFAULT '', amount DECIMAL(12,2) NOT NULL DEFAULT 0, due_date DATE NOT NULL, status VARCHAR(20) NOT NULL DEFAULT 'pending', message TEXT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS chat_forms (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, fields TEXT NOT NULL, status VARCHAR(20) NOT NULL DEFAULT 'active', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS form_submissions (id BIGINT AUTO_INCREMENT PRIMARY KEY, form_id BIGINT NOT NULL, phone VARCHAR(64) NOT NULL, data TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Products ----
func (d *DB) AddProduct(name, desc string, price float64, imageURL, category string, stock int) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO products (name, description, price, image_url, category, stock) VALUES (?,?,?,?,?,?)`, name, desc, price, imageURL, category, stock)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteProduct(id int64) error { _, err := d.sql.Exec(`DELETE FROM products WHERE id=?`, id); return err }
func (d *DB) ListProducts() ([]Product, error) {
	rows, err := d.sql.Query(`SELECT id, name, IFNULL(description,''), price, IFNULL(image_url,''), category, stock, status, created_at FROM products ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Category, &p.Stock, &p.Status, &p.Created)
		out = append(out, p)
	}
	return out, nil
}
func (d *DB) GetProduct(id int64) (*Product, error) {
	var p Product
	err := d.sql.QueryRow(`SELECT id, name, IFNULL(description,''), price, IFNULL(image_url,''), category, stock, status, created_at FROM products WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Category, &p.Stock, &p.Status, &p.Created)
	if err != nil { return nil, err }
	return &p, nil
}
func (d *DB) ProductsByCategory(cat string) ([]Product, error) {
	rows, err := d.sql.Query(`SELECT id, name, IFNULL(description,''), price, IFNULL(image_url,''), category, stock, status, created_at FROM products WHERE category=? AND status='active' ORDER BY name`, cat)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Category, &p.Stock, &p.Status, &p.Created)
		out = append(out, p)
	}
	return out, nil
}

// ---- Categories ----
func (d *DB) AddCategory(name string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO product_categories (name) VALUES (?)`, name)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteCategory(id int64) error { _, err := d.sql.Exec(`DELETE FROM product_categories WHERE id=?`, id); return err }
func (d *DB) ListCategories() ([]ProductCategory, error) {
	rows, err := d.sql.Query(`SELECT id, name, created_at FROM product_categories ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ProductCategory
	for rows.Next() {
		var c ProductCategory
		rows.Scan(&c.ID, &c.Name, &c.Created)
		out = append(out, c)
	}
	return out, nil
}

// ---- Orders ----
func (d *DB) CreateOrder(phone, name string, productID int64, quantity int, total float64) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO store_orders (phone, name, product_id, quantity, total) VALUES (?,?,?,?,?)`, phone, name, productID, quantity, total)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateOrderStatus(id int64, status string) error {
	_, err := d.sql.Exec(`UPDATE store_orders SET status=? WHERE id=?`, status, id)
	return err
}
func (d *DB) ListOrders() ([]Order, error) {
	rows, err := d.sql.Query(`SELECT id, phone, name, product_id, quantity, total, status, created_at FROM store_orders ORDER BY id DESC LIMIT 200`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.Phone, &o.Name, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.Created)
		out = append(out, o)
	}
	return out, nil
}

// ---- Behavior Tracking ----
func (d *DB) TrackBehavior(phone, event, meta string) {
	d.sql.Exec(`INSERT INTO behavior_events (phone, event, meta) VALUES (?,?,?)`, phone, event, meta)
}
func (d *DB) GetBehaviorSegments() map[string][]string {
	rows, err := d.sql.Query(`SELECT phone, event, COUNT(*) as cnt FROM behavior_events GROUP BY phone, event HAVING cnt > 2`)
	if err != nil { return nil }
	defer rows.Close()
	seg := map[string][]string{}
	for rows.Next() {
		var phone, event string
		var cnt int
		rows.Scan(&phone, &event, &cnt)
		seg[phone] = append(seg[phone], event)
	}
	return seg
}

// ---- Payment Reminders ----
func (d *DB) AddReminder(phone, name string, amount float64, dueDate, message string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO payment_reminders (phone, name, amount, due_date, message) VALUES (?,?,?,?,?)`, phone, name, amount, dueDate, message)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) ListReminders() ([]PaymentReminder, error) {
	rows, err := d.sql.Query(`SELECT id, phone, name, amount, due_date, status, IFNULL(message,''), created_at FROM payment_reminders ORDER BY due_date`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []PaymentReminder
	for rows.Next() {
		var r PaymentReminder
		rows.Scan(&r.ID, &r.Phone, &r.Name, &r.Amount, &r.DueDate, &r.Status, &r.Message, &r.Created)
		out = append(out, r)
	}
	return out, nil
}
func (d *DB) DueReminders() ([]PaymentReminder, error) {
	rows, err := d.sql.Query(`SELECT id, phone, name, amount, due_date, status, IFNULL(message,''), created_at FROM payment_reminders WHERE status='pending' AND due_date <= CURDATE() + INTERVAL 3 DAY`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []PaymentReminder
	for rows.Next() {
		var r PaymentReminder
		rows.Scan(&r.ID, &r.Phone, &r.Name, &r.Amount, &r.DueDate, &r.Status, &r.Message, &r.Created)
		out = append(out, r)
	}
	return out, nil
}
func (d *DB) MarkReminderSent(id int64) error {
	_, err := d.sql.Exec(`UPDATE payment_reminders SET status='sent' WHERE id=?`, id)
	return err
}

// ---- Forms ----
func (d *DB) AddForm(name, fieldsJSON string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO chat_forms (name, fields) VALUES (?,?)`, name, fieldsJSON)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeleteForm(id int64) error { _, err := d.sql.Exec(`DELETE FROM chat_forms WHERE id=?`, id); return err }
func (d *DB) ListForms() ([]ChatForm, error) {
	rows, err := d.sql.Query(`SELECT id, name, fields, status, created_at FROM chat_forms ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []ChatForm
	for rows.Next() {
		var f ChatForm
		rows.Scan(&f.ID, &f.Name, &f.Fields, &f.Status, &f.Created)
		out = append(out, f)
	}
	return out, nil
}
func (d *DB) SubmitForm(formID int64, phone, data string) error {
	_, err := d.sql.Exec(`INSERT INTO form_submissions (form_id, phone, data) VALUES (?,?,?)`, formID, phone, data)
	return err
}
func (d *DB) ListSubmissions(formID int64) ([]FormSubmission, error) {
	rows, err := d.sql.Query(`SELECT id, form_id, phone, data, created_at FROM form_submissions WHERE form_id=? ORDER BY id DESC LIMIT 200`, formID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []FormSubmission
	for rows.Next() {
		var s FormSubmission
		rows.Scan(&s.ID, &s.FormID, &s.Phone, &s.Data, &s.Created)
		out = append(out, s)
	}
	return out, nil
}

// ---- Agent Analytics ----
func (d *DB) AgentMetrics() []AgentMetric {
	rows, err := d.sql.Query(`SELECT IFNULL(u.name,'Unassigned') as agent, COUNT(DISTINCT s.phone) as chats, COUNT(*) as msgs, IFNULL(AVG(TIMESTAMPDIFF(SECOND, r.created_at, s.created_at)),0) as avg_time FROM sent s LEFT JOIN received r ON r.phone=s.phone AND r.created_at < s.created_at LEFT JOIN users u ON u.id = (SELECT agent_id FROM agent_assignments WHERE phone=s.phone LIMIT 1) WHERE s.created_at > DATE_SUB(NOW(), INTERVAL 30 DAY) GROUP BY agent`)
	if err != nil { return nil }
	defer rows.Close()
	var out []AgentMetric
	for rows.Next() {
		var m AgentMetric
		rows.Scan(&m.AgentName, &m.Chats, &m.Replied, &m.AvgTime)
		out = append(out, m)
	}
	return out
}

var _ = time.Now
