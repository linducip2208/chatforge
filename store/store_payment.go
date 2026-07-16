package store

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

type PaymentGateway struct {
	ID        int64
	Name      string
	Provider  string
	APIKey    string
	APISecret string
	WebhookSecret string
	BaseURL   string
	Currency  string
	Config    string
	Status    string
	Created   string
}

type PaymentTransaction struct {
	ID         int64
	UserID     int64
	PackageID  int64
	GatewayID  int64
	Amount     float64
	Currency   string
	InvoiceID  string
	ExternalID string
	Status     string
	PaymentURL string
	PaidAt     string
	Created    string
}

func (d *DB) migratePayment() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS payment_gateways (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			provider VARCHAR(50) NOT NULL,
			api_key TEXT NOT NULL,
			api_secret TEXT NOT NULL DEFAULT '',
			webhook_secret TEXT NOT NULL DEFAULT '',
			base_url VARCHAR(512) NOT NULL DEFAULT '',
			currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
			config TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS payment_transactions (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			package_id BIGINT NOT NULL,
			gateway_id BIGINT NOT NULL,
			amount DECIMAL(12,2) NOT NULL DEFAULT 0,
			currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
			invoice_id VARCHAR(64) NOT NULL UNIQUE,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			payment_url TEXT NOT NULL,
			external_id VARCHAR(128) NOT NULL DEFAULT '',
			paid_at DATETIME NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}
	for _, s := range stmts {
		if _, err := d.sql.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// ---- Gateways ----
func (d *DB) AddPaymentGateway(name, provider, apiKey, apiSecret, webhookSecret, baseURL, currency, config string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO payment_gateways (name, provider, api_key, api_secret, webhook_secret, base_url, currency, config) VALUES (?,?,?,?,?,?,?,?)`,
		name, provider, apiKey, apiSecret, webhookSecret, baseURL, currency, config)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) DeletePaymentGateway(id int64) error {
	_, err := d.sql.Exec(`DELETE FROM payment_gateways WHERE id=?`, id)
	return err
}
func (d *DB) ListPaymentGateways() ([]PaymentGateway, error) {
	rows, err := d.sql.Query(`SELECT id, name, provider, api_key, api_secret, webhook_secret, base_url, currency, config, status, created_at FROM payment_gateways ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []PaymentGateway
	for rows.Next() {
		var g PaymentGateway
		rows.Scan(&g.ID, &g.Name, &g.Provider, &g.APIKey, &g.APISecret, &g.WebhookSecret, &g.BaseURL, &g.Currency, &g.Config, &g.Status, &g.Created)
		out = append(out, g)
	}
	return out, nil
}
func (d *DB) GetPaymentGateway(id int64) (*PaymentGateway, error) {
	var g PaymentGateway
	err := d.sql.QueryRow(`SELECT id, name, provider, api_key, api_secret, webhook_secret, base_url, currency, config, status, created_at FROM payment_gateways WHERE id=?`, id).Scan(&g.ID, &g.Name, &g.Provider, &g.APIKey, &g.APISecret, &g.WebhookSecret, &g.BaseURL, &g.Currency, &g.Config, &g.Status, &g.Created)
	if err != nil { return nil, err }
	return &g, nil
}
func (d *DB) TogglePaymentGateway(id int64) error {
	_, err := d.sql.Exec(`UPDATE payment_gateways SET status=IF(status='active','inactive','active') WHERE id=?`, id)
	return err
}

// ---- Transactions ----
func (d *DB) CreateTransaction(userID, packageID, gatewayID int64, amount float64, currency, invoiceID, paymentURL, externalID string) (int64, error) {
	res, err := d.sql.Exec(`INSERT INTO payment_transactions (user_id, package_id, gateway_id, amount, currency, invoice_id, status, payment_url, external_id) VALUES (?,?,?,?,?,?,'pending',?,?)`,
		userID, packageID, gatewayID, amount, currency, invoiceID, paymentURL, externalID)
	if err != nil { return 0, err }
	return res.LastInsertId()
}
func (d *DB) UpdateTransactionStatus(invoiceID, status string) error {
	_, err := d.sql.Exec(`UPDATE payment_transactions SET status=?, paid_at=IF(?='paid',NOW(),NULL) WHERE invoice_id=?`, status, status, invoiceID)
	return err
}
func (d *DB) GetTransactionByInvoice(invoiceID string) (*PaymentTransaction, error) {
	var t PaymentTransaction
	err := d.sql.QueryRow(`SELECT id, user_id, package_id, gateway_id, amount, currency, invoice_id, status, payment_url, external_id, IFNULL(paid_at,''), created_at FROM payment_transactions WHERE invoice_id=?`, invoiceID).Scan(&t.ID, &t.UserID, &t.PackageID, &t.GatewayID, &t.Amount, &t.Currency, &t.InvoiceID, &t.Status, &t.PaymentURL, &t.ExternalID, &t.PaidAt, &t.Created)
	if err != nil { return nil, err }
	return &t, nil
}
func (d *DB) ListPayTransactions() ([]PaymentTransaction, error) {
	rows, err := d.sql.Query(`SELECT id, user_id, package_id, gateway_id, amount, currency, invoice_id, status, payment_url, external_id, IFNULL(paid_at,''), created_at FROM payment_transactions ORDER BY id DESC LIMIT 100`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []PaymentTransaction
	for rows.Next() {
		var t PaymentTransaction
		rows.Scan(&t.ID, &t.UserID, &t.PackageID, &t.GatewayID, &t.Amount, &t.Currency, &t.InvoiceID, &t.Status, &t.PaymentURL, &t.ExternalID, &t.PaidAt, &t.Created)
		out = append(out, t)
	}
	return out, nil
}
func (d *DB) ActivateSubscription(userID, packageID int64) error {
	var pkgName string
	var price string
	d.sql.QueryRow(`SELECT name, price FROM packages WHERE id=?`, packageID).Scan(&pkgName, &price)
	days := 30
	expire := time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05")
	_, err := d.sql.Exec(`INSERT INTO subscriptions (user_id, package_id, expire, status) VALUES (?,?,?,'active') ON DUPLICATE KEY UPDATE package_id=VALUES(package_id), expire=VALUES(expire), status='active'`, userID, packageID, expire)
	_ = pkgName
	_ = price
	return err
}

func GenInvoiceID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "INV-" + strings.ToUpper(hex.EncodeToString(b))
}
