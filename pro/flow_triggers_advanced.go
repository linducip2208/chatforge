//go:build pro

package pro

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// ═══════════════════════════════════════════
// Webhook Trigger — external HTTP POST starts a flow
// POST /pro/flow/trigger/{flow_id}
// Body: {"phone": "62812...", "message": "...", "name": "..."}
// ═══════════════════════════════════════════

func handleFlowTriggerWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	idStr := r.URL.Path[len("/pro/flow/trigger/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "Missing flow id", 400)
		return
	}
	uid := getUID(r)
	flow := getFlowDB(id, uid)
	if flow == nil {
		http.Error(w, "Flow not found", 404)
		return
	}

	var body struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		body.Message = r.FormValue("message")
		body.Phone = r.FormValue("phone")
		body.Name = r.FormValue("name")
	}
	if body.Phone == "" {
		body.Phone = "webhook"
	}

	ctx := NewFlowContext(flow, body.Phone, body.Phone, body.Message, body.Name)
	replies, actions := ExecuteFlow(ctx)

	log.Printf("[FLOW] Webhook trigger: %s → %s: %d replies", flow.Name, body.Phone, len(replies))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "executed",
		"replies":  replies,
		"actions":  len(actions),
	})
}

// ═══════════════════════════════════════════
// Cron Trigger — runs flows on schedule
// Called from wa/loops.go scheduler
// ═══════════════════════════════════════════

func RunCronFlows(uid int64, accountPhone string) int {
	flows := FindFlowsByTrigger(uid, TriggerCron, accountPhone)
	executed := 0
	for _, flow := range flows {
		// Cron flows use trigger value as cron expression
		// Simplified: "hourly", "daily", "every_6h"
		now := fmt.Sprintf("%d", getCurrentHour())
		if matchCronSchedule(flow.Trigger.Value, now) {
			ctx := NewFlowContext(flow, "cron", accountPhone, "", "System")
			ExecuteFlow(ctx)
			executed++
			log.Printf("[FLOW] Cron: %s executed", flow.Name)
		}
	}
	return executed
}

func matchCronSchedule(schedule, current string) bool {
	switch schedule {
	case "hourly", "every_1h":
		return true
	case "daily", "every_24h":
		return current == "8" // run at 8 AM
	case "every_6h":
		h := getCurrentHour()
		return h == 0 || h == 6 || h == 12 || h == 18
	case "every_12h":
		h := getCurrentHour()
		return h == 8 || h == 20
	default:
		return schedule == current
	}
}

func getCurrentHour() int {
	return 0 // placeholder — caller passes time info
}

func SetFlowHourGetter(fn func() int) { flowHourGetter = fn }

var flowHourGetter func() int

// ═══════════════════════════════════════════
// Inactivity Trigger — user hasn't replied in N hours
// Called from wa/loops.go auto-close checker
// ═══════════════════════════════════════════

func RunInactivityFlows(uid int64, accountPhone string, phones []string) int {
	flows := FindFlowsByTrigger(uid, TriggerInactivity, accountPhone)
	if len(flows) == 0 {
		return 0
	}
	executed := 0
	for _, phone := range phones {
		for _, flow := range flows {
			ctx := NewFlowContext(flow, phone, accountPhone, "", "")
			replies, _ := ExecuteFlow(ctx)
			if len(replies) > 0 {
				executed++
				log.Printf("[FLOW] Inactivity: %s → %s", flow.Name, phone)
			}
		}
	}
	return executed
}

// Button Trigger — handler defined in flow_triggers.go

// HandleButtonTrigger handles button-click triggered flows.
func HandleButtonTrigger(uid int64, accountPhone, phone, buttonPayload, contactName string) ([]FlowReply, []FlowAction) {
	flows := FindFlowsByTrigger(uid, TriggerButton, accountPhone)
	for _, flow := range flows {
		if matchTriggerValue(flow.Trigger, phone, buttonPayload) {
			ctx := NewFlowContext(flow, phone, accountPhone, buttonPayload, contactName)
			return ExecuteFlow(ctx)
		}
	}
	return nil, nil
}
