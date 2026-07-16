package store

import "database/sql"

// Admin / AI / Android entities — the remaining Zender menus.
// Auto-reply core & WhatsApp engine are untouched.

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Role     string
	Country  string
	Created  string
}
type Role struct {
	ID          int64
	Name        string
	Permissions string
	Created     string
}
type Package struct {
	ID             int64
	Name           string
	Price          string
	SendLimit      int
	ReceiveLimit   int
	DeviceLimit    int
	UssdLimit      int
	WaSendLimit    int
	WaReceiveLimit int
	WaAccountLimit int
	ContactLimit   int
	ScheduledLimit int
	KeyLimit       int
	WebhookLimit   int
	ActionLimit    int
	Services       string
	Hidden         int
	Footermark     int
	Created        string
}
type Voucher struct {
	ID       int64
	Code     string
	Pkg      string
	Duration int
	Created  string
}
type Subscription struct {
	ID      int64
	User    string
	Pkg     string
	Expire  string
	Created string
}
type Transaction struct {
	ID       int64
	User     string
	Amount   string
	Provider string
	Created  string
}
type Payout struct {
	ID       int64
	User     string
	Amount   string
	Address  string
	Status   string
	Created  string
}
type Page struct {
	ID      int64
	Title   string
	Slug    string
	Content string
	Created string
}
type Marketing struct {
	ID      int64
	Title   string
	Content string
	Created string
}
type Language struct {
	ID      int64
	Name    string
	ISO     string
	Created string
}
type WaServer struct {
	ID       int64
	Name     string
	URL      string
	Port     string
	Secret   string
	Accounts string
	Packages string
	Created  string
}
type Gateway struct {
	ID         int64
	Name       string
	Callback   string
	CallbackID string
	Created    string
}
type Shortener struct {
	ID      int64
	Name    string
	Created string
}
type Plugin struct {
	ID      int64
	Name    string
	Dir     string
	Created string
}
type AiKey struct {
	ID           int64
	Name         string
	Provider     string
	Model        string
	APIKey       string
	BaseURL      string
	SystemPrompt string
	Created      string
}
type AiPlugin struct {
	ID       int64
	Name     string
	Endpoint string
	Created  string
}
type Device struct {
	ID           int64
	Name         string
	DID          string
	Manufacturer string
	Created      string
}
type Ussd struct {
	ID       int64
	Code     string
	Response string
	Status   string
	Created  string
}

func (d *DB) migrateAdmin() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, email VARCHAR(255) NOT NULL, role VARCHAR(64) NOT NULL DEFAULT 'user', country VARCHAR(8) NOT NULL DEFAULT 'ID', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS roles (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, permissions TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS packages (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, price VARCHAR(32) NOT NULL DEFAULT '0', send_limit INT NOT NULL DEFAULT 0, device_limit INT NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS vouchers (id BIGINT AUTO_INCREMENT PRIMARY KEY, code VARCHAR(64) NOT NULL, pkg VARCHAR(255) NOT NULL DEFAULT '', duration INT NOT NULL DEFAULT 30, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS subscriptions (id BIGINT AUTO_INCREMENT PRIMARY KEY, user VARCHAR(255) NOT NULL DEFAULT '', pkg VARCHAR(255) NOT NULL DEFAULT '', expire DATE NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS transactions (id BIGINT AUTO_INCREMENT PRIMARY KEY, user VARCHAR(255) NOT NULL DEFAULT '', amount VARCHAR(32) NOT NULL DEFAULT '0', provider VARCHAR(64) NOT NULL DEFAULT 'manual', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS payouts (id BIGINT AUTO_INCREMENT PRIMARY KEY, user VARCHAR(255) NOT NULL DEFAULT '', amount VARCHAR(32) NOT NULL DEFAULT '0', address VARCHAR(255) NOT NULL DEFAULT '', status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS pages (id BIGINT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255) NOT NULL, slug VARCHAR(255) NOT NULL DEFAULT '', content TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS marketing (id BIGINT AUTO_INCREMENT PRIMARY KEY, title VARCHAR(255) NOT NULL, content TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS languages_admin (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, iso VARCHAR(8) NOT NULL DEFAULT 'us', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS waservers (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, url VARCHAR(255) NOT NULL DEFAULT '', port VARCHAR(16) NOT NULL DEFAULT '', secret VARCHAR(255) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS gateways (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS shorteners (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS plugins (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, dir VARCHAR(255) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS ai_keys (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, provider VARCHAR(64) NOT NULL DEFAULT 'openai', model VARCHAR(64) NOT NULL DEFAULT '', apikey VARCHAR(512) NOT NULL DEFAULT '', base_url VARCHAR(255) NOT NULL DEFAULT '', system_prompt TEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS ai_plugins (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, endpoint VARCHAR(255) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS devices (id BIGINT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, did VARCHAR(128) NOT NULL DEFAULT '', manufacturer VARCHAR(128) NOT NULL DEFAULT '', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS ussd (id BIGINT AUTO_INCREMENT PRIMARY KEY, code VARCHAR(64) NOT NULL, response TEXT NOT NULL, status VARCHAR(20) NOT NULL DEFAULT 'pending', created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) MustListAiTrainings() []AiTraining {
	list, _ := d.ListAiTrainings()
	return list
}

// AI Token Usage



func (d *DB) GetUserPackageName(userID int64) string {
	u, err := d.GetUserByID(userID)
	if err != nil { return "" }
	var pkgName string
	err = d.sql.QueryRow(`SELECT pkg FROM subscriptions WHERE user=? ORDER BY id DESC LIMIT 1`, u.Email).Scan(&pkgName)
	if err != nil { return "" }
	return pkgName
}

func (d *DB) RecordAiUsage(userID int64, tokens int, provider, model string) {
	d.sql.Exec(`INSERT INTO ai_usage (user_id, tokens, provider, model) VALUES (?,?,?,?)`, userID, tokens, provider, model)
}
func (d *DB) GetAiTokenUsage(userID int64) int64 {
	var total sql.NullInt64
	d.sql.QueryRow(`SELECT SUM(tokens) FROM ai_usage WHERE user_id=?`, userID).Scan(&total)
	return total.Int64
}
func (d *DB) GetUserAiQuota(userID int64) int64 {
	u, err := d.GetUserByID(userID)
	if err != nil { return 100000 }
	var pkgName string
	err = d.sql.QueryRow(`SELECT pkg FROM subscriptions WHERE user=? ORDER BY id DESC LIMIT 1`, u.Email).Scan(&pkgName)
	if err != nil { return 100000 }
	var quota sql.NullInt64
	d.sql.QueryRow(`SELECT ai_token_quota FROM packages WHERE name=? LIMIT 1`, pkgName).Scan(&quota)
	if quota.Int64 <= 0 { return 100000 }
	return quota.Int64
}
type AiTraining struct {
	ID           int64
	Name         string
	SystemPrompt string
	AiKeyID      int64
	Created      string
}

func (d *DB) AddAiTraining(name, systemPrompt string, aiKeyID int64) (int64, error) {
	return d.exec(`INSERT INTO ai_trainings (name,system_prompt,ai_key_id) VALUES (?,?,?)`, name, systemPrompt, aiKeyID)
}
func (d *DB) DeleteAiTraining(id int64) error { return d.del("ai_trainings", id) }
func (d *DB) ListAiTrainings() ([]AiTraining, error) {
	rows, err := d.sql.Query(`SELECT id,name,system_prompt,ai_key_id,created_at FROM ai_trainings ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []AiTraining
	for rows.Next() {
		var t AiTraining
		rows.Scan(&t.ID, &t.Name, &t.SystemPrompt, &t.AiKeyID, &t.Created)
		out = append(out, t)
	}
	return out, nil
}

// generic exec helpers
func (d *DB) exec(q string, args ...interface{}) (int64, error) {
	res, err := d.sql.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (d *DB) del(table string, id int64) error {
	_, err := d.sql.Exec("DELETE FROM "+table+" WHERE id=?", id)
	return err
}

// Users
func (d *DB) AddUser(name, email, role, country string) (int64, error) {
	return d.exec(`INSERT INTO users (name,email,role,country) VALUES (?,?,?,?)`, name, email, role, country)
}
func (d *DB) DeleteUser(id int64) error { return d.del("users", id) }
func (d *DB) GetUserByEmail(email string) (*User, error) {
	var u User
	err := d.sql.QueryRow(`SELECT id,name,email,role,country,password,created_at FROM users WHERE email=?`, email).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Country, &u.Password, &u.Created)
	if err != nil { return nil, err }
	return &u, nil
}
func (d *DB) GetUserByID(id int64) (*User, error) {
	var u User
	err := d.sql.QueryRow(`SELECT id,name,email,role,country,password,created_at FROM users WHERE id=?`, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Country, &u.Password, &u.Created)
	if err != nil { return nil, err }
	return &u, nil
}
func (d *DB) GetUserPackageLimit(userID int64) int {
	u, err := d.GetUserByID(userID)
	if err != nil { return 0 }
	var pkgName string
	err = d.sql.QueryRow(`SELECT pkg FROM subscriptions WHERE user=? ORDER BY id DESC LIMIT 1`, u.Email).Scan(&pkgName)
	if err != nil { return 0 }
	var limit int
	err = d.sql.QueryRow(`SELECT device_limit FROM packages WHERE name=? LIMIT 1`, pkgName).Scan(&limit)
	if err != nil { return 0 }
	return limit
}

func (d *DB) GetUserSendLimit(userID int64) int {
	u, err := d.GetUserByID(userID)
	if err != nil { return 0 }
	var pkgName string
	err = d.sql.QueryRow(`SELECT pkg FROM subscriptions WHERE user=? ORDER BY id DESC LIMIT 1`, u.Email).Scan(&pkgName)
	if err != nil { return 0 }
	var limit int
	err = d.sql.QueryRow(`SELECT wa_send_limit FROM packages WHERE name=? LIMIT 1`, pkgName).Scan(&limit)
	if err != nil { return 0 }
	return limit
}

func (d *DB) CountSentByUser(userID int64) int {
	u, err := d.GetUserByID(userID)
	if err != nil { return 0 }
	var n int
	d.sql.QueryRow(`SELECT COUNT(*) FROM sent WHERE phone IN (SELECT phone FROM wa_accounts WHERE user_id=?) OR status='sent'`, userID).Scan(&n)
	return n
}
func (d *DB) SetUserPassword(id int64, hash string) error {
	_, err := d.sql.Exec(`UPDATE users SET password=? WHERE id=?`, hash, id)
	return err
}
func (d *DB) UpdateUser(id int64, name, email, role string) error {
	_, err := d.sql.Exec(`UPDATE users SET name=?, email=?, role=? WHERE id=?`, name, email, role, id)
	return err
}
func (d *DB) ListUsers() ([]User, error) {
	rows, err := d.sql.Query(`SELECT id,name,email,role,country,created_at FROM users ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Country, &u.Created)
		out = append(out, u)
	}
	return out, nil
}

// Roles
func (d *DB) AddRole(name, perms string) (int64, error) {
	return d.exec(`INSERT INTO roles (name,permissions) VALUES (?,?)`, name, perms)
}
func (d *DB) DeleteRole(id int64) error { return d.del("roles", id) }
func (d *DB) ListRoles() ([]Role, error) {
	rows, err := d.sql.Query(`SELECT id,name,permissions,created_at FROM roles ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Role
	for rows.Next() {
		var x Role
		rows.Scan(&x.ID, &x.Name, &x.Permissions, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Packages
func (d *DB) AddPackage(name, price string, send, receive, dev, ussd, waSend, waReceive, waAcc, contact, scheduled, keyL, webhookL, actionL int, services string, hidden, footermark int) (int64, error) {
	return d.exec(`INSERT INTO packages (name,price,send_limit,receive_limit,device_limit,ussd_limit,wa_send_limit,wa_receive_limit,wa_account_limit,contact_limit,scheduled_limit,key_limit,webhook_limit,action_limit,services,hidden,footermark) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		name, price, send, receive, dev, ussd, waSend, waReceive, waAcc, contact, scheduled, keyL, webhookL, actionL, services, hidden, footermark)
}
func (d *DB) DeletePackage(id int64) error { return d.del("packages", id) }
func (d *DB) ListPackages() ([]Package, error) {
	rows, err := d.sql.Query(`SELECT id,name,price,send_limit,receive_limit,device_limit,ussd_limit,wa_send_limit,wa_receive_limit,wa_account_limit,contact_limit,scheduled_limit,key_limit,webhook_limit,action_limit,IFNULL(services,''),hidden,footermark,created_at FROM packages ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Package
	for rows.Next() {
		var x Package
		rows.Scan(&x.ID, &x.Name, &x.Price, &x.SendLimit, &x.ReceiveLimit, &x.DeviceLimit, &x.UssdLimit, &x.WaSendLimit, &x.WaReceiveLimit, &x.WaAccountLimit, &x.ContactLimit, &x.ScheduledLimit, &x.KeyLimit, &x.WebhookLimit, &x.ActionLimit, &x.Services, &x.Hidden, &x.Footermark, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Vouchers
func (d *DB) AddVoucher(code, pkg string, dur int) (int64, error) {
	return d.exec(`INSERT INTO vouchers (code,pkg,duration) VALUES (?,?,?)`, code, pkg, dur)
}
func (d *DB) DeleteVoucher(id int64) error { return d.del("vouchers", id) }
func (d *DB) ListVouchers() ([]Voucher, error) {
	rows, err := d.sql.Query(`SELECT id,code,pkg,duration,created_at FROM vouchers ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Voucher
	for rows.Next() {
		var x Voucher
		rows.Scan(&x.ID, &x.Code, &x.Pkg, &x.Duration, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Subscriptions
func (d *DB) AddSubscription(user, pkg, expire string) (int64, error) {
	return d.exec(`INSERT INTO subscriptions (user,pkg,expire) VALUES (?,?,?)`, user, pkg, expire)
}
func (d *DB) DeleteSubscription(id int64) error { return d.del("subscriptions", id) }
func (d *DB) ListSubscriptions() ([]Subscription, error) {
	rows, err := d.sql.Query(`SELECT id,user,pkg,IFNULL(expire,''),created_at FROM subscriptions ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subscription
	for rows.Next() {
		var x Subscription
		rows.Scan(&x.ID, &x.User, &x.Pkg, &x.Expire, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Transactions
func (d *DB) ListTransactions() ([]Transaction, error) {
	rows, err := d.sql.Query(`SELECT id,user,amount,provider,created_at FROM transactions ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Transaction
	for rows.Next() {
		var x Transaction
		rows.Scan(&x.ID, &x.User, &x.Amount, &x.Provider, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Payouts
func (d *DB) DeletePayout(id int64) error { return d.del("payouts", id) }
func (d *DB) ListPayouts() ([]Payout, error) {
	rows, err := d.sql.Query(`SELECT id,user,amount,address,status,created_at FROM payouts ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Payout
	for rows.Next() {
		var x Payout
		rows.Scan(&x.ID, &x.User, &x.Amount, &x.Address, &x.Status, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Pages
func (d *DB) AddPage(title, slug, content string) (int64, error) {
	return d.exec(`INSERT INTO pages (title,slug,content) VALUES (?,?,?)`, title, slug, content)
}
func (d *DB) DeletePage(id int64) error { return d.del("pages", id) }
func (d *DB) ListPages() ([]Page, error) {
	rows, err := d.sql.Query(`SELECT id,title,slug,content,created_at FROM pages ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Page
	for rows.Next() {
		var x Page
		rows.Scan(&x.ID, &x.Title, &x.Slug, &x.Content, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Marketing
func (d *DB) AddMarketing(title, content string) (int64, error) {
	return d.exec(`INSERT INTO marketing (title,content) VALUES (?,?)`, title, content)
}
func (d *DB) DeleteMarketing(id int64) error { return d.del("marketing", id) }
func (d *DB) ListMarketing() ([]Marketing, error) {
	rows, err := d.sql.Query(`SELECT id,title,content,created_at FROM marketing ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Marketing
	for rows.Next() {
		var x Marketing
		rows.Scan(&x.ID, &x.Title, &x.Content, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Languages (admin)
func (d *DB) AddLanguageAdmin(name, iso string) (int64, error) {
	return d.exec(`INSERT INTO languages_admin (name,iso) VALUES (?,?)`, name, iso)
}
func (d *DB) DeleteLanguageAdmin(id int64) error { return d.del("languages_admin", id) }
func (d *DB) ListLanguagesAdmin() ([]Language, error) {
	rows, err := d.sql.Query(`SELECT id,name,iso,created_at FROM languages_admin ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Language
	for rows.Next() {
		var x Language
		rows.Scan(&x.ID, &x.Name, &x.ISO, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// WaServers
func (d *DB) AddWaServer(name, url, port, secret string, accounts int, packages string) (int64, error) {
	return d.exec(`INSERT INTO waservers (name,url,port,secret,accounts,packages) VALUES (?,?,?,?,?,?)`, name, url, port, secret, accounts, packages)
}
func (d *DB) DeleteWaServer(id int64) error { return d.del("waservers", id) }
func (d *DB) GetWaServer(id int64) (*WaServer, error) {
	var w WaServer
	err := d.sql.QueryRow(`SELECT id,name,url,port,secret,IFNULL(accounts,0),IFNULL(packages,''),created_at FROM waservers WHERE id=?`, id).Scan(&w.ID, &w.Name, &w.URL, &w.Port, &w.Secret, &w.Accounts, &w.Packages, &w.Created)
	if err != nil { return nil, err }
	return &w, nil
}
func (d *DB) UpdateWaServer(id int64, name, url, port, secret string, accounts int, packages string) error {
	_, err := d.sql.Exec(`UPDATE waservers SET name=?, url=?, port=?, secret=?, accounts=?, packages=? WHERE id=?`, name, url, port, secret, accounts, packages, id)
	return err
}
func (d *DB) ListWaServers() ([]WaServer, error) {
	rows, err := d.sql.Query(`SELECT id,name,url,port,secret,IFNULL(accounts,0),IFNULL(packages,''),created_at FROM waservers ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []WaServer
	for rows.Next() {
		var x WaServer
		rows.Scan(&x.ID, &x.Name, &x.URL, &x.Port, &x.Secret, &x.Accounts, &x.Packages, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Gateways
func (d *DB) AddGateway(name string) (int64, error) {
	return d.exec(`INSERT INTO gateways (name) VALUES (?)`, name)
}
func (d *DB) DeleteGateway(id int64) error { return d.del("gateways", id) }
func (d *DB) ListGateways() ([]Gateway, error) {
	rows, err := d.sql.Query(`SELECT id,name,IFNULL(callback,''),IFNULL(callback_id,''),created_at FROM gateways ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Gateway
	for rows.Next() {
		var x Gateway
		rows.Scan(&x.ID, &x.Name, &x.Callback, &x.CallbackID, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Shorteners
func (d *DB) AddShortener(name string) (int64, error) {
	return d.exec(`INSERT INTO shorteners (name) VALUES (?)`, name)
}
func (d *DB) DeleteShortener(id int64) error { return d.del("shorteners", id) }
func (d *DB) ListShorteners() ([]Shortener, error) {
	rows, err := d.sql.Query(`SELECT id,name,created_at FROM shorteners ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Shortener
	for rows.Next() {
		var x Shortener
		rows.Scan(&x.ID, &x.Name, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Plugins
func (d *DB) AddPlugin(name, dir string) (int64, error) {
	return d.exec(`INSERT INTO plugins (name,dir) VALUES (?,?)`, name, dir)
}
func (d *DB) DeletePlugin(id int64) error { return d.del("plugins", id) }
func (d *DB) ListPlugins() ([]Plugin, error) {
	rows, err := d.sql.Query(`SELECT id,name,dir,created_at FROM plugins ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Plugin
	for rows.Next() {
		var x Plugin
		rows.Scan(&x.ID, &x.Name, &x.Dir, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// AI Keys
func (d *DB) AddAiKey(name, provider, model, apikey, baseURL, systemPrompt string) (int64, error) {
	return d.exec(`INSERT INTO ai_keys (name,provider,model,apikey,base_url,system_prompt) VALUES (?,?,?,?,?,?)`, name, provider, model, apikey, baseURL, systemPrompt)
}
func (d *DB) DeleteAiKey(id int64) error { return d.del("ai_keys", id) }
func (d *DB) GetAiKey(id int64) (*AiKey, error) {
	var a AiKey
	err := d.sql.QueryRow(`SELECT id,name,provider,model,apikey,base_url,system_prompt,created_at FROM ai_keys WHERE id=?`, id).Scan(&a.ID, &a.Name, &a.Provider, &a.Model, &a.APIKey, &a.BaseURL, &a.SystemPrompt, &a.Created)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (d *DB) ListAiKeys() ([]AiKey, error) {
	rows, err := d.sql.Query(`SELECT id,name,provider,model,apikey,base_url,system_prompt,created_at FROM ai_keys ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AiKey
	for rows.Next() {
		var x AiKey
		rows.Scan(&x.ID, &x.Name, &x.Provider, &x.Model, &x.APIKey, &x.BaseURL, &x.SystemPrompt, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// AI Plugins
func (d *DB) AddAiPlugin(name, endpoint string) (int64, error) {
	return d.exec(`INSERT INTO ai_plugins (name,endpoint) VALUES (?,?)`, name, endpoint)
}
func (d *DB) DeleteAiPlugin(id int64) error { return d.del("ai_plugins", id) }
func (d *DB) ListAiPlugins() ([]AiPlugin, error) {
	rows, err := d.sql.Query(`SELECT id,name,endpoint,created_at FROM ai_plugins ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AiPlugin
	for rows.Next() {
		var x AiPlugin
		rows.Scan(&x.ID, &x.Name, &x.Endpoint, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// Devices
func (d *DB) AddDevice(name, did, manuf string) (int64, error) {
	return d.exec(`INSERT INTO devices (name,did,manufacturer) VALUES (?,?,?)`, name, did, manuf)
}
func (d *DB) DeleteDevice(id int64) error { return d.del("devices", id) }
func (d *DB) ListDevices() ([]Device, error) {
	rows, err := d.sql.Query(`SELECT id,name,did,manufacturer,created_at FROM devices ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Device
	for rows.Next() {
		var x Device
		rows.Scan(&x.ID, &x.Name, &x.DID, &x.Manufacturer, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

// USSD
func (d *DB) AddUssd(code string) (int64, error) {
	return d.exec(`INSERT INTO ussd (code,response) VALUES (?,'')`, code)
}
func (d *DB) DeleteUssd(id int64) error { return d.del("ussd", id) }
func (d *DB) ListUssd() ([]Ussd, error) {
	rows, err := d.sql.Query(`SELECT id,code,response,status,created_at FROM ussd ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Ussd
	for rows.Next() {
		var x Ussd
		rows.Scan(&x.ID, &x.Code, &x.Response, &x.Status, &x.Created)
		out = append(out, x)
	}
	return out, nil
}

