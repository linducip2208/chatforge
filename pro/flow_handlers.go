//go:build pro

package pro

import (
	"chatgo/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var flowDB *store.DB

func SetFlowDB(db *store.DB) { flowDB = db }

func requireUID(r *http.Request, w http.ResponseWriter) (int64, bool) {
	uid := getUID(r)
	if uid == 0 {
		http.Error(w, `{"error":"Login required"}`, 401)
		return 0, false
	}
	return uid, true
}

func loadFlows(uid int64) []*Flow {
	if flowDB == nil {
		return nil
	}
	raw, err := flowDB.LoadActiveFlows(uid)
	if err != nil {
		log.Printf("[FLOW] DB error: %v", err)
		return nil
	}
	var result []*Flow
	for _, r := range raw {
		f := &Flow{
			ID:      r.ID,
			UserID:  r.UserID,
			Name:    r.Name,
			Active:  r.Active == 1,
			AIKeyID: flowDB.GetFlowAIKey(r.ID),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
		json.Unmarshal([]byte(r.Trigger), &f.Trigger)
		json.Unmarshal([]byte(r.NodesJSON), &f.Nodes)
		json.Unmarshal([]byte(r.EdgesJSON), &f.Edges)
		result = append(result, f)
	}
	return result
}

func saveFlowDB(f *Flow) error {
	if flowDB == nil {
		return fmt.Errorf("no DB connection")
	}
	triggerJSON, _ := json.Marshal(f.Trigger)
	nodesJSON, _ := json.Marshal(f.Nodes)
	edgesJSON, _ := json.Marshal(f.Edges)

	if f.ID > 0 {
		err := flowDB.UpdateFlowRaw(f.ID, f.UserID, f.Name, string(triggerJSON), string(nodesJSON), string(edgesJSON))
		if err == nil {
			flowDB.SetFlowAIKey(f.ID, f.AIKeyID)
		}
		return err
	}
	id, err := flowDB.SaveFlowRaw(f.UserID, f.Name, string(triggerJSON), string(nodesJSON), string(edgesJSON))
	if err == nil {
		f.ID = id
		if f.AIKeyID > 0 {
			flowDB.SetFlowAIKey(f.ID, f.AIKeyID)
		}
	}
	return err
}

func deleteFlowDB(id, uid int64) bool {
	if flowDB == nil {
		return false
	}
	return flowDB.DeleteFlowRaw(id, uid) == nil
}

func getFlowDB(id, uid int64) *Flow {
	if flowDB == nil {
		return nil
	}
	r, err := flowDB.GetFlowRaw(id, uid)
	if err != nil {
		return nil
	}
	f := &Flow{
		ID:      r.ID,
		UserID:  r.UserID,
		Name:    r.Name,
		Active:  r.Active == 1,
		AIKeyID: flowDB.GetFlowAIKey(r.ID),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	json.Unmarshal([]byte(r.Trigger), &f.Trigger)
	json.Unmarshal([]byte(r.NodesJSON), &f.Nodes)
	json.Unmarshal([]byte(r.EdgesJSON), &f.Edges)
	return f
}

func toggleFlowDB(id, uid int64) (bool, error) {
	if flowDB == nil {
		return false, fmt.Errorf("no DB")
	}
	return flowDB.ToggleFlowRaw(id, uid)
}

func duplicateFlowDB(id, uid int64) (int64, error) {
	if flowDB == nil {
		return 0, fmt.Errorf("no DB")
	}
	return flowDB.DuplicateFlowRaw(id, uid)
}

// ═══════════════════════════════════════════
// HTTP Handlers
// ═══════════════════════════════════════════

func handleFlowList(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	all := loadFlows(uid)
	if all == nil {
		all = []*Flow{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func handleFlowSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid := getUID(r)
	if uid == 0 {
		http.Error(w, `{"error":"Login required"}`, 401)
		return
	}
	var f Flow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), 400)
		return
	}
	f.UserID = uid
	f.Active = true
	if f.ID > 0 {
		if existing := getFlowDB(f.ID, uid); existing == nil {
			http.Error(w, `{"error":"Not found or not authorized"}`, 403)
			return
		}
	}
	if err := saveFlowDB(&f); err != nil {
		http.Error(w, fmt.Sprintf("Save error: %v", err), 500)
		return
	}
	// Auto-save version history
	if flowDB != nil {
		nj, _ := json.Marshal(f.Nodes); ej, _ := json.Marshal(f.Edges)
		flowDB.SaveFlowVersion(f.ID, f.Name, string(nj), string(ej))
	}
	log.Printf("[FLOW] Saved: #%d %s (%d nodes, %d edges)", f.ID, f.Name, len(f.Nodes), len(f.Edges))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": f.ID, "status": "ok"})
}

func handleFlowLoad(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 { http.Error(w, "Missing id", 400); return }
	f := getFlowDB(id, uid)
	if f == nil {
		http.Error(w, "Not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func handleFlowDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		http.Error(w, "Missing id", 400)
		return
	}
	if deleteFlowDB(id, uid) {
		log.Printf("[FLOW] Deleted: #%d", id)
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	} else {
		http.Error(w, "Not found", 404)
	}
}

func handleFlowToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	active, err := toggleFlowDB(id, uid)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"active": active})
}

func handleFlowDuplicate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.FormValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	newID, err := duplicateFlowDB(id, uid)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	log.Printf("[FLOW] Duplicated: #%d → #%d", id, newID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": newID, "status": "ok"})
}

func handleFlowBuilderUI(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var name, email string
	flowDB.QueryRow(`SELECT name, email FROM users WHERE id=? LIMIT 1`, uid).Scan(&name, &email)
	html := flowBuilderHTML
	if name != "" {
		html = strings.Replace(html, `<span class="flow-name">Flow Builder Pro</span>`,
			`<span class="flow-name">`+name+`</span><span style="font-size:10px;opacity:.5;margin-left:6px">`+email+`</span>`, 1)
	}
	fmt.Fprint(w, html)
}

func handleFlowExport(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 { http.Error(w, "Missing id", 400); return }
	f := getFlowDB(id, uid)
	if f == nil {
		http.Error(w, "Not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name+".json")
	json.NewEncoder(w).Encode(f)
}

func handleFlowImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	var f Flow
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), 400)
		return
	}
	f.ID = 0
	f.UserID = uid
	f.Active = true
	if err := saveFlowDB(&f); err != nil {
		http.Error(w, "Save error: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": f.ID, "status": "imported"})
}

func handleFlowAnalytics(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	f := getFlowDB(id, uid)
	if f == nil {
		http.Error(w, "Not found", 404)
		return
	}
	a := GetFlowAnalytics(id)
	if a == nil {
		a = &FlowAnalytics{FlowID: id, NodeHits: make(map[string]int)}
	}
	a.NodeHits = GetNodeMetrics(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

// Version History
func handleFlowVersions(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if id == 0 { http.Error(w, "Missing id", 400); return }
	if getFlowDB(id, uid) == nil { http.Error(w, "Not found", 404); return }
	versions, _ := flowDB.GetFlowVersions(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

func handleFlowRollback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	fid, _ := strconv.ParseInt(r.FormValue("flow_id"), 10, 64)
	vid, _ := strconv.ParseInt(r.FormValue("version_id"), 10, 64)
	if getFlowDB(fid, uid) == nil { http.Error(w, "Not found", 404); return }
	if err := flowDB.RollbackFlow(fid, vid); err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(map[string]string{"status": "rolled_back"})
}

// Marketplace
func handleFlowPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err := flowDB.PublishFlow(id, uid); err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(map[string]string{"status": "published"})
}
func handleFlowUnpublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err := flowDB.UnpublishFlow(id, uid); err != nil { http.Error(w, err.Error(), 500); return }
	json.NewEncoder(w).Encode(map[string]string{"status": "unpublished"})
}
func handleFlowMarketplace(w http.ResponseWriter, r *http.Request) {
	raw, _ := flowDB.ListPublicFlows()
	var result []*Flow
	for _, r := range raw {
		f := &Flow{ID: r.ID, UserID: r.UserID, Name: r.Name, Active: r.Active == 1, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
		json.Unmarshal([]byte(r.Trigger), &f.Trigger)
		json.Unmarshal([]byte(r.NodesJSON), &f.Nodes)
		json.Unmarshal([]byte(r.EdgesJSON), &f.Edges)
		result = append(result, f)
	}
	if result == nil { result = []*Flow{} }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// AI Flow Generator
func handleFlowAIGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	var req struct{ Prompt string `json:"prompt"` }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Prompt == "" { req.Prompt = r.FormValue("prompt") }
	if req.Prompt == "" { http.Error(w, "Missing prompt", 400); return }

	// Generate flow using AI
	reply := ""
	if AICallback != nil {
		sysPrompt := `You are a flow builder generator. Generate a complete WhatsApp chatbot flow in JSON.
The flow must have: name, trigger (type: keyword, value: keywords), nodes array, edges array.
Node types available: message, question, condition, wait, ai_reply, transfer_agent, close_chat, poll, location, set_variable, tag_contact, api_call, counter, math, loop, random, split_merge.
Output ONLY valid JSON, no markdown, no explanation.`
		reply, _ = AICallback(sysPrompt, req.Prompt, nil, 0)
	}

	var flow Flow
	if err := json.Unmarshal([]byte(reply), &flow); err != nil {
		// Try to extract JSON from AI response
		start := strings.Index(reply, "{")
		end := strings.LastIndex(reply, "}")
		if start >= 0 && end > start {
			json.Unmarshal([]byte(reply[start:end+1]), &flow)
		}
	}
	if flow.Name == "" {
		http.Error(w, "AI could not generate flow: "+reply, 500)
		return
	}
	flow.UserID = uid
	flow.Active = true
	saveFlowDB(&flow)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": flow.ID, "status": "generated", "name": flow.Name})
}

func handleFlowLogs(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	id, _ := strconv.ParseInt(r.URL.Query().Get("flow_id"), 10, 64)
	if id == 0 { http.Error(w, "Missing flow_id", 400); return }
	if getFlowDB(id, uid) == nil { http.Error(w, "Not found", 404); return }
	logs, _ := flowDB.GetFlowExecutionLog(id, 100)
	if logs == nil { logs = []map[string]interface{}{} }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func getUID(r *http.Request) int64 {
	c, err := r.Cookie("chatgo_sess")
	if err != nil { return 0 }
	uid, _ := flowDB.GetSession(c.Value)
	return uid
}

// ═══════ WA Accounts for Flow Builder ═══════
type WAAccountInfo struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Type  string `json:"type"` // "wa" or "meta"
	Label string `json:"label"`
}

var AccountLister func(uid int64) []WAAccountInfo

func SetAccountLister(fn func(uid int64) []WAAccountInfo) { AccountLister = fn }

func handleFlowAccounts(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	var list []WAAccountInfo
	if AccountLister != nil {
		list = AccountLister(uid)
	}
	if list == nil { list = []WAAccountInfo{} }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// ═══════════════════════════════════════════
// Flow Simulator — execute flow server-side for testing
// POST /pro/flow-builder/simulate
// Body: { "flow": {...}, "phone": "62812...", "message": "halo", "name": "Test", "debug": true }
// ═══════════════════════════════════════════
func handleFlowSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid, ok := requireUID(r, w)
	if !ok { return }
	var req struct {
		Flow    *Flow  `json:"flow"`
		Phone   string `json:"phone"`
		Message string `json:"message"`
		Name    string `json:"name"`
		Debug   bool   `json:"debug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), 400)
		return
	}
	if req.Flow == nil {
		http.Error(w, "Missing flow", 400)
		return
	}
	if req.Phone == "" { req.Phone = "test" }
	if req.Name == "" { req.Name = "Tester" }
	req.Flow.UserID = uid
	result := SimulateFlow(req.Flow, req.Phone, req.Message, req.Name, req.Debug)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ═══════════════════════════════════════════
// Flow Debugger — step-by-step node execution
// POST /pro/flow-builder/debug
// ═══════════════════════════════════════════
var debugBreaks = make(map[string]map[string]bool) // flow_id -> set of node ids with breakpoints

func handleFlowDebug(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	action := r.URL.Query().Get("action")
	switch action {
	case "set_break":
		fid := r.FormValue("flow_id")
		nid := r.FormValue("node_id")
		if _, exists := debugBreaks[fid]; !exists {
			debugBreaks[fid] = make(map[string]bool)
		}
		debugBreaks[fid][nid] = true
		json.NewEncoder(w).Encode(map[string]string{"status": "breakpoint_set"})
	case "clear_break":
		fid := r.FormValue("flow_id")
		nid := r.FormValue("node_id")
		if breaks, exists := debugBreaks[fid]; exists {
			delete(breaks, nid)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "breakpoint_cleared"})
	case "clear_all":
		fid := r.FormValue("flow_id")
		delete(debugBreaks, fid)
		json.NewEncoder(w).Encode(map[string]string{"status": "all_cleared"})
	case "get_breaks":
		fid := r.URL.Query().Get("flow_id")
		breaks := debugBreaks[fid]
		if breaks == nil { breaks = make(map[string]bool) }
		ids := make([]string, 0, len(breaks))
		for k := range breaks { ids = append(ids, k) }
		json.NewEncoder(w).Encode(map[string]interface{}{"flow_id": fid, "breakpoints": ids})
	default:
		http.Error(w, "Unknown action. Use: set_break, clear_break, clear_all, get_breaks", 400)
	}
	_ = uid
}

// ═══════════════════════════════════════════
// Flow Reviews
// ═══════════════════════════════════════════
func handleFlowReviews(w http.ResponseWriter, r *http.Request) {
	fid, _ := strconv.ParseInt(r.URL.Query().Get("flow_id"), 10, 64)
	if fid == 0 { http.Error(w, "Missing flow_id", 400); return }
	if r.Method == http.MethodPost {
		uid, ok := requireUID(r, w)
		if !ok { return }
		rating, _ := strconv.Atoi(r.FormValue("rating"))
		review := r.FormValue("review")
		if rating < 1 { rating = 5 }
		if err := flowDB.AddFlowReview(fid, uid, rating, review); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "reviewed"})
		return
	}
	reviews, _ := flowDB.GetFlowReviews(fid)
	if reviews == nil { reviews = []map[string]interface{}{} }
	avg := flowDB.GetFlowAvgRating(fid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"reviews": reviews, "avg_rating": avg, "count": len(reviews)})
}

// ═══════════════════════════════════════════
// Flow Download tracking
// ═══════════════════════════════════════════
func handleFlowDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "POST only", 405); return }
	uid, ok := requireUID(r, w)
	if !ok { return }
	fid, _ := strconv.ParseInt(r.FormValue("flow_id"), 10, 64)
	if fid == 0 { http.Error(w, "Missing flow_id", 400); return }
	if !flowDB.HasDownloadedFlow(fid, uid) {
		flowDB.MarkFlowDownloaded(fid, uid)
		flowDB.IncFlowTrigger(fid) // reuse counter for download tracking
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "downloaded"})
}

// ═══════════════════════════════════════════
// AI Keys list for Flow Builder dropdown
// ═══════════════════════════════════════════
func handleFlowAIKeys(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireUID(r, w)
	if !ok { return }
	type AIKeyInfo struct {
		ID       int64  `json:"id"`
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}
	var keys []AIKeyInfo
	// Try to query ai_keys table
	rows, err := flowDB.Query(`SELECT id, provider, model FROM ai_keys WHERE user_id=? OR user_id=0 ORDER BY id`, uid)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var k AIKeyInfo
			rows.Scan(&k.ID, &k.Provider, &k.Model)
			keys = append(keys, k)
		}
	}
	if keys == nil { keys = []AIKeyInfo{} }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// ═══════════════════════════════════════════
// Omnichannel Inbox UI
// ═══════════════════════════════════════════
func handleOmniInboxUI(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, omniInboxHTML)
}

// ═══════════════════════════════════════════
// Agency Dashboard UI
// ═══════════════════════════════════════════
func handleAgencyUI(w http.ResponseWriter, r *http.Request) {
	uid := getUID(r)
	if uid == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, agencyHTML)
}
