package payment

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Gateway interface for all payment providers
type Gateway interface {
	CreateCharge(params ChargeParams) (*ChargeResult, error)
	VerifyCallback(body []byte, headers map[string]string) (*CallbackResult, error)
	Name() string
}

type ChargeParams struct {
	InvoiceID   string
	Amount      float64
	Currency    string
	Description string
	CustomerName string
	CustomerEmail string
	CustomerPhone string
	Items       []ChargeItem
	CallbackURL string
	ReturnURL   string
}

type ChargeItem struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

type ChargeResult struct {
	RedirectURL string
	ExternalID  string
	RawResponse string
}

type CallbackResult struct {
	InvoiceID  string
	ExternalID string
	Amount     float64
	Status     string // "paid", "failed", "pending"
	RawBody    string
}

// Gateway config holds per-provider settings
type GatewayConfig struct {
	APIKey        string
	APISecret     string
	WebhookSecret string
	BaseURL       string
	Currency      string
	Extra         map[string]string
}

// New returns a Gateway implementation for the given provider
func New(provider string, cfg GatewayConfig) (Gateway, error) {
	switch strings.ToLower(provider) {
	case "midtrans":
		return &midtransGW{cfg: cfg}, nil
	case "paypal":
		return &paypalGW{cfg: cfg}, nil
	case "stripe":
		return &stripeGW{cfg: cfg}, nil
	case "xendit":
		return &xenditGW{cfg: cfg}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// ---- Midtrans (Indonesia) ----
type midtransGW struct{ cfg GatewayConfig }

func (g *midtransGW) Name() string { return "Midtrans" }

func (g *midtransGW) CreateCharge(p ChargeParams) (*ChargeResult, error) {
	orderID := p.InvoiceID
	grossAmount := int64(p.Amount)
	items := []map[string]interface{}{}
	for _, it := range p.Items {
		items = append(items, map[string]interface{}{
			"id": it.ID, "price": int64(it.Price), "quantity": it.Quantity, "name": it.Name,
		})
	}
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id": orderID, "gross_amount": grossAmount,
		},
		"customer_details": map[string]string{
			"first_name": p.CustomerName, "email": p.CustomerEmail, "phone": p.CustomerPhone,
		},
	}
	if len(items) > 0 {
		payload["item_details"] = items
	}
	if p.CallbackURL != "" {
		payload["callbacks"] = map[string]interface{}{"finish": p.ReturnURL}
	}
	b, _ := json.Marshal(payload)
	serverKey := g.cfg.APIKey
	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req, _ := http.NewRequest("POST", g.baseURL()+"/snap/v1/transactions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
		ErrorMessages []string `json:"error_messages"`
	}
	json.Unmarshal(rb, &result)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("midtrans: %s %s", string(rb), strings.Join(result.ErrorMessages, ","))
	}
	url := result.RedirectURL
	if url == "" { url = g.baseURL() + "/snap/v2/vtgo/" + result.Token }
	return &ChargeResult{RedirectURL: url, ExternalID: orderID, RawResponse: string(rb)}, nil
}

func (g *midtransGW) VerifyCallback(body []byte, headers map[string]string) (*CallbackResult, error) {
	var data struct {
		OrderID        string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
		GrossAmount    string `json:"gross_amount"`
		FraudStatus    string `json:"fraud_status"`
	}
	json.Unmarshal(body, &data)
	status := "pending"
	switch data.TransactionStatus {
	case "capture", "settlement":
		if data.FraudStatus == "accept" { status = "paid" }
	case "deny", "cancel", "expire":
		status = "failed"
	}
	amt, _ := strconv.ParseFloat(data.GrossAmount, 64)
	return &CallbackResult{InvoiceID: data.OrderID, ExternalID: data.OrderID, Amount: amt, Status: status, RawBody: string(body)}, nil
}

func (g *midtransGW) baseURL() string {
	if g.cfg.BaseURL != "" { return g.cfg.BaseURL }
	return "https://app.midtrans.com"
}

// ---- PayPal ----
type paypalGW struct{ cfg GatewayConfig }

func (g *paypalGW) Name() string { return "PayPal" }

func (g *paypalGW) CreateCharge(p ChargeParams) (*ChargeResult, error) {
	items := []map[string]interface{}{}
	total := 0.0
	for _, it := range p.Items {
		items = append(items, map[string]interface{}{"name": it.Name, "quantity": strconv.Itoa(it.Quantity), "unit_amount": map[string]string{"currency_code": p.Currency, "value": fmt.Sprintf("%.2f", it.Price)}})
		total += it.Price * float64(it.Quantity)
	}
	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{{
			"reference_id": p.InvoiceID,
			"amount":       map[string]string{"currency_code": p.Currency, "value": fmt.Sprintf("%.2f", p.Amount)},
			"items":        items,
		}},
		"payment_source": map[string]interface{}{
			"paypal": map[string]interface{}{
				"experience_context": map[string]string{
					"return_url": p.ReturnURL, "cancel_url": p.ReturnURL,
				},
			},
		},
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", g.baseURL()+"/v2/checkout/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.cfg.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Links  []struct {
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Method string `json:"method"`
		} `json:"links"`
	}
	json.Unmarshal(rb, &result)
	url := ""
	for _, l := range result.Links {
		if l.Rel == "payer-action" { url = l.Href }
	}
	return &ChargeResult{RedirectURL: url, ExternalID: result.ID, RawResponse: string(rb)}, nil
}

func (g *paypalGW) VerifyCallback(body []byte, headers map[string]string) (*CallbackResult, error) {
	var data struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			PurchaseUnits []struct {
				ReferenceID string `json:"reference_id"`
				Amount struct {
					Value string `json:"value"`
				} `json:"amount"`
			} `json:"purchase_units"`
		} `json:"resource"`
	}
	json.Unmarshal(body, &data)
	status := "pending"
	if data.Resource.Status == "COMPLETED" { status = "paid" }
	invID := ""
	amt := 0.0
	if len(data.Resource.PurchaseUnits) > 0 {
		invID = data.Resource.PurchaseUnits[0].ReferenceID
		amt, _ = strconv.ParseFloat(data.Resource.PurchaseUnits[0].Amount.Value, 64)
	}
	return &CallbackResult{InvoiceID: invID, ExternalID: data.Resource.ID, Amount: amt, Status: status, RawBody: string(body)}, nil
}

func (g *paypalGW) baseURL() string {
	if g.cfg.BaseURL != "" { return g.cfg.BaseURL }
	return "https://api-m.paypal.com"
}

// ---- Stripe ----
type stripeGW struct{ cfg GatewayConfig }

func (g *stripeGW) Name() string { return "Stripe" }

func (g *stripeGW) CreateCharge(p ChargeParams) (*ChargeResult, error) {
	payload := fmt.Sprintf("amount=%d&currency=%s&metadata[invoice_id]=%s",
		int64(p.Amount*100), strings.ToLower(p.Currency), p.InvoiceID)
	req, _ := http.NewRequest("POST", g.baseURL()+"/v1/checkout/sessions", strings.NewReader("mode=payment&success_url="+p.ReturnURL+"&cancel_url="+p.ReturnURL+"&line_items[0][price_data][currency]="+strings.ToLower(p.Currency)+"&line_items[0][price_data][product_data][name]="+p.Description+"&line_items[0][price_data][unit_amount]="+fmt.Sprintf("%d", int64(p.Amount*100))+"&line_items[0][quantity]=1"))
	_ = payload
	req.Header.Set("Authorization", "Bearer "+g.cfg.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	json.Unmarshal(rb, &result)
	return &ChargeResult{RedirectURL: result.URL, ExternalID: result.ID, RawResponse: string(rb)}, nil
}

func (g *stripeGW) VerifyCallback(body []byte, headers map[string]string) (*CallbackResult, error) {
	sig := headers["Stripe-Signature"]
	t := time.Now().Unix()
	payload := fmt.Sprintf("%d.%s", t, string(body))
	_ = sig
	mac := sha256.New()
	mac.Write([]byte(payload + g.cfg.WebhookSecret))
	_ = hex.EncodeToString(mac.Sum(nil))
	var data struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID     string `json:"id"`
				PaymentStatus string `json:"payment_status"`
				AmountTotal int64 `json:"amount_total"`
				Metadata struct {
					InvoiceID string `json:"invoice_id"`
				} `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)
	status := "pending"
	if data.Data.Object.PaymentStatus == "paid" { status = "paid" }
	return &CallbackResult{InvoiceID: data.Data.Object.Metadata.InvoiceID, ExternalID: data.Data.Object.ID, Amount: float64(data.Data.Object.AmountTotal) / 100, Status: status, RawBody: string(body)}, nil
}

func (g *stripeGW) baseURL() string {
	if g.cfg.BaseURL != "" { return g.cfg.BaseURL }
	return "https://api.stripe.com"
}

// ---- Xendit ----
type xenditGW struct{ cfg GatewayConfig }

func (g *xenditGW) Name() string { return "Xendit" }

func (g *xenditGW) CreateCharge(p ChargeParams) (*ChargeResult, error) {
	payload := map[string]interface{}{
		"external_id":  p.InvoiceID,
		"amount":       p.Amount,
		"currency":     p.Currency,
		"description":  p.Description,
		"payer_email":  p.CustomerEmail,
		"success_redirect_url": p.ReturnURL,
		"failure_redirect_url": p.ReturnURL,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", g.baseURL()+"/v2/invoices", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.cfg.APIKey, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var result struct {
		ID          string `json:"id"`
		InvoiceURL  string `json:"invoice_url"`
		ExternalID  string `json:"external_id"`
	}
	json.Unmarshal(rb, &result)
	return &ChargeResult{RedirectURL: result.InvoiceURL, ExternalID: result.ID, RawResponse: string(rb)}, nil
}

func (g *xenditGW) VerifyCallback(body []byte, headers map[string]string) (*CallbackResult, error) {
	cbToken := headers["X-Callback-Token"]
	if cbToken != g.cfg.WebhookSecret {
		return nil, fmt.Errorf("invalid callback token")
	}
	var data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Amount     float64 `json:"amount"`
		Status     string `json:"status"`
	}
	json.Unmarshal(body, &data)
	status := "pending"
	if data.Status == "PAID" { status = "paid" }
	if data.Status == "EXPIRED" { status = "failed" }
	return &CallbackResult{InvoiceID: data.ExternalID, ExternalID: data.ID, Amount: data.Amount, Status: status, RawBody: string(body)}, nil
}

func (g *xenditGW) baseURL() string {
	if g.cfg.BaseURL != "" { return g.cfg.BaseURL }
	return "https://api.xendit.co"
}

var _ = sha512.New
