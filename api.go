package main

import (
	"chatgo/msgtemplate"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func checkAPIKey(r *http.Request) bool {
	apikey := r.Header.Get("X-API-Key")
	if apikey == "" { apikey = r.URL.Query().Get("apikey") }
	return apikey != "" && db.ValidAPIKey(apikey)
}

func handleAPISend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAPIError(w, "method not allowed", 405)
		return
	}
	if !checkAPIKey(r) {
		writeAPIError(w, "invalid api key", 401)
		return
	}
	var req struct {
		Phone        string `json:"phone"`
		Message      string `json:"message"`
		AccountPhone string `json:"account_phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, "invalid json body", 400)
		return
	}
	if req.Phone == "" || req.Message == "" {
		writeAPIError(w, "phone and message required", 400)
		return
	}
	if err := engine.SendFrom(0, strings.TrimPrefix(req.AccountPhone, "+"), req.Phone, msgtemplate.Render(req.Message, msgtemplate.Vars{Phone: req.Phone})); err != nil {
		writeAPIError(w, err.Error(), 500)
		return
	}
	writeAPIOK(w, map[string]string{"status": "ok", "message": "sent"})
}

func handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(r) { writeAPIError(w, "invalid api key", 401); return }
	status, phone := engine.Status()
	writeAPIOK(w, map[string]interface{}{
		"status":  status,
		"phone":   phone,
		"accounts": engine.Accounts(0),
	})
}

func handleAPIMessages(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(r) { writeAPIError(w, "invalid api key", 401); return }
	t := r.URL.Query().Get("type")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if page <= 0 { page = 1 }
	if perPage <= 0 { perPage = 20 }

	switch t {
	case "sent":
		list, _ := db.ListSentPaginated(page, perPage)
		writeAPIOK(w, map[string]interface{}{"type": "sent", "page": page, "data": list, "total": db.CountSent()})
	case "received":
		list, _ := db.ListReceivedPaginated(page, perPage)
		writeAPIOK(w, map[string]interface{}{"type": "received", "page": page, "data": list, "total": db.CountReceived()})
	default:
		sent, _ := db.ListSentPaginated(1, 10)
		recv, _ := db.ListReceivedPaginated(1, 10)
		writeAPIOK(w, map[string]interface{}{"sent": sent, "received": recv})
	}
}

func handleAPIContacts(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(r) { writeAPIError(w, "invalid api key", 401); return }
	contacts, _ := db.ListContacts(0)
	groups, _ := db.ListGroups(0)
	writeAPIOK(w, map[string]interface{}{"contacts": contacts, "groups": groups})
}

func handleAPIDevices(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(r) { writeAPIError(w, "invalid api key", 401); return }
	devices, _ := db.ListDevices()
	writeAPIOK(w, map[string]interface{}{"devices": devices})
}

func writeAPIOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
func writeAPIError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": msg})
}
