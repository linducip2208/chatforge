//go:build pro

package pro

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// ════════════════ Node Types ════════════════
type FlowNodeType string

const (
	NodeMessage       FlowNodeType = "message"
	NodeQuestion      FlowNodeType = "question"
	NodeCondition     FlowNodeType = "condition"
	NodeWait          FlowNodeType = "wait"
	NodeSetVariable   FlowNodeType = "set_variable"
	NodeTag           FlowNodeType = "tag_contact"
	NodeAPICall       FlowNodeType = "api_call"
	NodeTransfer      FlowNodeType = "transfer_agent"
	NodeClose         FlowNodeType = "close_chat"
	NodeAI            FlowNodeType = "ai_reply"
	NodeAIDecide      FlowNodeType = "ai_decide"
	NodeCarousel      FlowNodeType = "product_carousel"
	NodeOrderCreate   FlowNodeType = "order_create"
	NodePaymentLink   FlowNodeType = "payment_link"
	NodeMath          FlowNodeType = "math"
	NodeStringTemplate FlowNodeType = "string_template"
	NodeDateMath      FlowNodeType = "date_math"
	NodeLoop          FlowNodeType = "loop"
	NodeSubflow       FlowNodeType = "subflow"
	NodeRandom        FlowNodeType = "random"
	NodeContactLookup FlowNodeType = "contact_lookup"
	NodeCounter       FlowNodeType = "counter"
	NodeSequence      FlowNodeType = "message_sequence"
	NodePoll          FlowNodeType = "poll"
	NodeLocation      FlowNodeType = "location"
	NodeDocument      FlowNodeType = "document"
	NodeTyping        FlowNodeType = "typing"
	NodeDelayUntil    FlowNodeType = "delay_until"
	NodeGallery       FlowNodeType = "media_gallery"
	NodeVoice         FlowNodeType = "voice_note"
	NodeVCard         FlowNodeType = "contact_card"
	NodeDBQuery       FlowNodeType = "db_query"
	NodeSplit         FlowNodeType = "split_merge"
	NodeReadReceipt   FlowNodeType = "read_receipt"
	NodeSticker       FlowNodeType = "sticker"
	NodeGoogleSheets  FlowNodeType = "google_sheets"
	NodeWarmer        FlowNodeType = "warmer"
	NodeTGKeyboard    FlowNodeType = "tg_keyboard"
	NodeTGPayment     FlowNodeType = "tg_payment"
	NodeIGComment     FlowNodeType = "ig_comment_dm"
	NodeIGStoryReply  FlowNodeType = "ig_story_reply"
	NodeIGQuickReply  FlowNodeType = "ig_quick_reply"
	NodeFBCarousel    FlowNodeType = "fb_carousel"
	NodeFBMenu        FlowNodeType = "fb_persistent_menu"
	NodeFBWebview     FlowNodeType = "fb_webview"
	NodeButtons       FlowNodeType = "buttons"
	NodeReceiveMedia  FlowNodeType = "receive_media"
)

// ════════════════ Trigger Types ════════════════
type TriggerType string

const (
	TriggerKeyword     TriggerType = "keyword"
	TriggerWelcome     TriggerType = "welcome"
	TriggerFallback    TriggerType = "fallback"
	TriggerDrip        TriggerType = "drip"
	TriggerBroadcast   TriggerType = "broadcast"
	TriggerCSAT        TriggerType = "csat"
	TriggerClose       TriggerType = "close_chat"
	TriggerAlways      TriggerType = "always"
	TriggerCron        TriggerType = "cron"
	TriggerButton      TriggerType = "button"
	TriggerInactivity  TriggerType = "inactivity"
	TriggerExact       TriggerType = "exact"
)

// ════════════════ Data Structures ════════════════
type FlowPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type FlowNode struct {
	ID       string          `json:"id"`
	Type     FlowNodeType    `json:"type"`
	Position FlowPosition    `json:"position"`
	Data     json.RawMessage `json:"data"`
}
type FlowEdge struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	SourceH   string `json:"sourceHandle,omitempty"`
	Target    string `json:"target"`
	Label     string `json:"label,omitempty"`
	Fallback  bool   `json:"fallback,omitempty"`
}
type FlowTrigger struct {
	Type         TriggerType `json:"type"`
	Value        string      `json:"value"`
	AccountPhone string      `json:"account_phone,omitempty"`
	GroupIDs     string      `json:"group_ids,omitempty"`
	TagIDs       string      `json:"tag_ids,omitempty"`
	Priority     int         `json:"priority"`
}
type ConditionalEnable struct {
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
	DaysOfWeek string `json:"days_of_week"`
}
type Flow struct {
	ID          int64              `json:"id"`
	UserID      int64              `json:"user_id"`
	Name        string             `json:"name"`
	Trigger     FlowTrigger        `json:"trigger"`
	Nodes       []FlowNode         `json:"nodes"`
	Edges       []FlowEdge         `json:"edges"`
	Active      bool               `json:"active"`
	Priority    int                `json:"priority"`
	Conditional *ConditionalEnable `json:"conditional,omitempty"`
	Tags        string             `json:"tags,omitempty"`
	RateLimit   int                `json:"rate_limit,omitempty"`
	AIKeyID     int64              `json:"ai_key_id,omitempty"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}
type FlowContext struct {
	Flow         *Flow
	Phone        string
	AccountPhone string
	Message      string
	ContactName  string
	Variables    map[string]string
	CurrentNode  string
	VisitedNodes map[string]int
	Replies      []FlowReply
	Actions      []FlowAction
	Done         bool
	MaxSteps     int
}
type FlowReply struct {
	Text       string `json:"text,omitempty"`
	MediaURL   string `json:"media_url,omitempty"`
	MediaType  string `json:"media_type,omitempty"`
	Action     string `json:"action,omitempty"`
	ActionData map[string]interface{} `json:"action_data,omitempty"`
}
type FlowAction struct {
	Type   string
	Params map[string]string
}
type FlowAnalytics struct {
	FlowID          int64          `json:"flow_id"`
	TriggerCount    int            `json:"trigger_count"`
	CompletionCount int            `json:"completion_count"`
	NodeHits        map[string]int `json:"node_hits"`
}

// Node data structs
type MessageData struct {
	Text      string `json:"text"`
	MediaURL  string `json:"media_url,omitempty"`
	MediaType string `json:"media_type,omitempty"`
}
type QuestionData struct {
	Question string `json:"question"`
	VarName  string `json:"var_name"`
	Type     string `json:"type"`
}
type ConditionBranch struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
type ConditionData struct {
	Variable string            `json:"variable"`
	Branches []ConditionBranch `json:"branches"`
}
type WaitData struct{ Seconds int `json:"seconds"` }
type SetVarData struct{ Variables map[string]string `json:"variables"` }
type TagData struct{ Tags []string `json:"tags"` }
type APICallData struct {
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	VarName  string            `json:"var_name"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     string            `json:"body,omitempty"`
}
type AIData struct {
	SystemPrompt string `json:"system_prompt"`
	VarResult    string `json:"var_result"`
}
type AIDecideData struct {
	SystemPrompt string   `json:"system_prompt"`
	Options      []string `json:"options"`
	VarResult    string   `json:"var_result"`
}
type CarouselData struct {
	Category string `json:"category"`
	MaxItems int    `json:"max_items"`
}
type OrderCreateData struct {
	ProductVar string `json:"product_var"`
	QtyVar     string `json:"qty_var"`
	VarResult  string `json:"var_result"`
}
type PaymentLinkData struct {
	OrderVar  string `json:"order_var"`
	AmountVar string `json:"amount_var"`
	VarResult string `json:"var_result"`
}
type MathData struct {
	Formula   string `json:"formula"`
	VarResult string `json:"var_result"`
}
type StringTemplateData struct {
	Template  string `json:"template"`
	VarResult string `json:"var_result"`
}
type DateMathData struct {
	Operation string `json:"operation"`
	Value     int    `json:"value"`
	VarResult string `json:"var_result"`
}
type LoopData struct {
	MaxIterations int    `json:"max_iterations"`
	ConditionVar  string `json:"condition_var"`
	LoopStartNode string `json:"loop_start"`
}
type SubflowData struct {
	FlowID  int64             `json:"flow_id"`
	Inputs  map[string]string `json:"inputs"`
	Outputs map[string]string `json:"outputs"`
}
type RandomData struct{ Weights []int `json:"weights"` }
type ContactLookupData struct{ Fields []string `json:"fields"` }
type CounterData struct {
	Key       string `json:"key"`
	MaxCount  int    `json:"max_count"`
	VarResult string `json:"var_result"`
}
type SequenceMessage struct {
	Text  string `json:"text"`
	Delay int    `json:"delay"`
}
type SequenceData struct{ Messages []SequenceMessage `json:"messages"` }
type PollData struct {
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	MaxSelect int      `json:"max_select"`
	VarResult string   `json:"var_result"`
}
type LocationData struct {
	Action    string  `json:"action"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Name      string  `json:"name"`
	VarResult string  `json:"var_result"`
}
type DocumentData struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
}
type TypingData struct{ Seconds int `json:"seconds"` }
type DelayUntilData struct {
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
}
type GalleryImage struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
}
type GalleryData struct {
	Images  []GalleryImage `json:"images"`
	Caption string         `json:"caption"`
}
type VoiceData struct {
	URL string `json:"url"`
	TTS string `json:"tts"`
}
type VCardData struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Org   string `json:"org"`
}
type DBQueryData struct {
	Query     string `json:"query"`
	VarResult string `json:"var_result"`
}

type StickerData struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

type ReceiveMediaData struct {
	Prompt     string `json:"prompt"`
	VarResult  string `json:"var_result"`
	MediaType  string `json:"media_type"` // image, document, any
}

type GoogleSheetsData struct {
	Action   string `json:"action"` // append, read
	SheetID  string `json:"sheet_id"`
	Range    string `json:"range"`  // "Sheet1!A:C"
	Data     string `json:"data"`   // comma-separated for append
	VarResult string `json:"var_result"`
}

// Google Sheets callback
var GoogleSheetsFunc func(action, sheetID, sheetRange, data string) string
func SetGoogleSheetsFunc(fn func(string, string, string, string) string) { GoogleSheetsFunc = fn }

// ════════════════ Callbacks ════════════════
var (
	AICallback          func(string, string, []string, int64) (string, int)
	StoreListProducts   func(string, int) []map[string]string
	StoreCreateOrder    func(string, string, int) (int64, error)
	StorePaymentLink    func(int64, float64) string
	ContactLookupFunc   func(string) map[string]string
	DBQueryFunc         func(string) (string, error)
)

func SetAICallback(fn func(string, string, []string, int64) (string, int)) { AICallback = fn }
func SetStoreCallbacks(lp func(string, int) []map[string]string, co func(string, string, int) (int64, error), pl func(int64, float64) string) {
	StoreListProducts = lp; StoreCreateOrder = co; StorePaymentLink = pl
}
func SetContactLookup(fn func(string) map[string]string) { ContactLookupFunc = fn }
func SetDBQueryFunc(fn func(string) (string, error))        { DBQueryFunc = fn }

// ════════════════ Execution Engine ════════════════
func NewFlowContext(flow *Flow, phone, accountPhone, message, contactName string) *FlowContext {
	return &FlowContext{
		Flow: flow, Phone: phone, AccountPhone: accountPhone,
		Message: strings.TrimSpace(message), ContactName: contactName,
		Variables: make(map[string]string), VisitedNodes: make(map[string]int),
		Replies: make([]FlowReply, 0), Actions: make([]FlowAction, 0), MaxSteps: 50,
	}
}

func ExecuteFlow(ctx *FlowContext) ([]FlowReply, []FlowAction) {
	start := findStartNode(ctx.Flow)
	if start == nil { ctx.Done = true; return ctx.Replies, ctx.Actions }
	ctx.CurrentNode = start.ID
	for !ctx.Done && len(ctx.VisitedNodes) < ctx.MaxSteps {
		node := findNode(ctx.Flow, ctx.CurrentNode)
		if node == nil { break }
		ctx.VisitedNodes[ctx.CurrentNode]++
		if ctx.VisitedNodes[ctx.CurrentNode] > 3 {
			log.Printf("[FLOW] Loop detected at %s", ctx.CurrentNode); break
		}
		recordNodeHit(ctx.Flow.ID, node.ID)
		executeNode(node, ctx)
	}
	// Log execution
	if flowDB != nil {
		errMsg := ""
		if ctx.Done && len(ctx.Replies) == 0 { errMsg = "no_reply" }
		flowDB.LogFlowExecution(ctx.Flow.ID, ctx.Flow.Name, ctx.Phone, string(ctx.Flow.Trigger.Type), len(ctx.VisitedNodes), len(ctx.Replies), "completed", errMsg)
	}
	return ctx.Replies, ctx.Actions
}

func executeNode(node *FlowNode, ctx *FlowContext) {
	switch node.Type {
	case NodeMessage:      executeMessage(node, ctx)
	case NodeQuestion:     executeQuestion(node, ctx)
	case NodeCondition:    executeCondition(node, ctx)
	case NodeWait:         executeWait(node, ctx)
	case NodeSetVariable:  executeSetVar(node, ctx)
	case NodeTag:          executeTag(node, ctx)
	case NodeAPICall:      executeAPICall(node, ctx)
	case NodeTransfer:     executeTransfer(node, ctx)
	case NodeClose:        executeClose(node, ctx)
	case NodeAI:           executeAI(node, ctx)
	case NodeAIDecide:     executeAIDecide(node, ctx)
	case NodeCarousel:     executeCarousel(node, ctx)
	case NodeOrderCreate:  executeOrderCreate(node, ctx)
	case NodePaymentLink:  executePaymentLink(node, ctx)
	case NodeMath:         executeMath(node, ctx)
	case NodeStringTemplate: executeStringTemplate(node, ctx)
	case NodeDateMath:     executeDateMath(node, ctx)
	case NodeLoop:         executeLoop(node, ctx)
	case NodeSubflow:      executeSubflow(node, ctx)
	case NodeRandom:       executeRandom(node, ctx)
	case NodeContactLookup: executeContactLookup(node, ctx)
	case NodeCounter:      executeCounter(node, ctx)
	case NodeSequence:     executeSequence(node, ctx)
	case NodePoll:         executePoll(node, ctx)
	case NodeLocation:     executeLocation(node, ctx)
	case NodeDocument:     executeDocument(node, ctx)
	case NodeTyping:       executeTyping(node, ctx)
	case NodeDelayUntil:   executeDelayUntil(node, ctx)
	case NodeGallery:      executeGallery(node, ctx)
	case NodeVoice:        executeVoice(node, ctx)
	case NodeVCard:        executeVCard(node, ctx)
	case NodeDBQuery:      executeDBQuery(node, ctx)
	case NodeSplit:        ctx.CurrentNode = next(ctx, "") // no-op, just forward
	case NodeReadReceipt:  executeReadReceipt(node, ctx)
	case NodeSticker:      executeSticker(node, ctx)
	case NodeGoogleSheets: executeGoogleSheets(node, ctx)
	case NodeWarmer:       executeWarmer(node, ctx)
	case NodeTGKeyboard:   executeTGKeyboard(node, ctx)
	case NodeTGPayment:    executeTGPayment(node, ctx)
	case NodeIGComment:    executeIGComment(node, ctx)
	case NodeIGStoryReply: executeIGStoryReply(node, ctx)
	case NodeIGQuickReply: executeIGQuickReply(node, ctx)
	case NodeFBCarousel:   executeFBCarousel(node, ctx)
	case NodeFBMenu:       executeFBMenu(node, ctx)
	case NodeFBWebview:    executeFBWebview(node, ctx)
	case NodeButtons:      executeButtons(node, ctx)
	case NodeReceiveMedia: executeReceiveMedia(node, ctx)
	}
}

// ════════════════ Matching + Analytics ════════════════
func FindMatchingFlow(uid int64, triggerType TriggerType, accountPhone, phone, message string) *Flow {
	all := loadFlows(uid)
	var best *Flow; bestPrio := -1
	for _, f := range all {
		if !f.Active || f.Trigger.Type != triggerType { continue }
		if f.Trigger.AccountPhone != "" && f.Trigger.AccountPhone != accountPhone { continue }
		if !matchTriggerValue(f.Trigger, phone, message) { continue }
		if f.Conditional != nil && !checkConditional(f.Conditional) { continue }
		if f.RateLimit > 0 && !checkFlowRateLimit(f.ID, phone, f.RateLimit) { continue }
		if f.Trigger.Priority > bestPrio { best = f; bestPrio = f.Trigger.Priority }
	}
	return best
}
func FindFlowsByTrigger(uid int64, triggerType TriggerType, accountPhone string) []*Flow {
	var result []*Flow
	for _, f := range loadFlows(uid) {
		if !f.Active || f.Trigger.Type != triggerType { continue }
		if f.Trigger.AccountPhone != "" && f.Trigger.AccountPhone != accountPhone { continue }
		result = append(result, f)
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Trigger.Priority > result[i].Trigger.Priority {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
func matchTriggerValue(t FlowTrigger, phone, message string) bool {
	switch t.Type {
	case TriggerAlways: return true
	case TriggerKeyword:
		if t.Value == "" || t.Value == "*" { return true }
		for _, kw := range splitTrim(t.Value, ",") {
			if kw != "" && strings.Contains(strings.ToLower(message), strings.ToLower(strings.TrimSpace(kw))) { return true }
		}
		return false
	case TriggerExact:
		if t.Value == "" || t.Value == "*" { return true }
		for _, kw := range splitTrim(t.Value, ",") {
			if strings.EqualFold(strings.TrimSpace(message), strings.TrimSpace(kw)) { return true }
		}
		return false
	case TriggerWelcome, TriggerFallback:
		return true
	default: return true
	}
}
func checkConditional(c *ConditionalEnable) bool {
	if c.TimeStart == "" && c.TimeEnd == "" && c.DaysOfWeek == "" { return true }
	now := time.Now()
	if c.TimeStart != "" && c.TimeEnd != "" {
		t := now.Format("15:04")
		if t < c.TimeStart || t > c.TimeEnd { return false }
	}
	if c.DaysOfWeek != "" {
		dow := fmt.Sprintf("%d", now.Weekday())
		if now.Weekday() == 0 { dow = "7" }
		if !strings.Contains(c.DaysOfWeek, dow) { return false }
	}
	return true
}
func RecordFlowTrigger(flowID int64) { if flowDB != nil { flowDB.IncFlowTrigger(flowID) } }
func RecordFlowComplete(flowID int64) { if flowDB != nil { flowDB.IncFlowComplete(flowID) } }
func recordNodeHit(flowID int64, nodeID string) { if flowDB != nil { flowDB.IncNodeHit(flowID, nodeID) } }
func GetFlowAnalytics(flowID int64) *FlowAnalytics {
	a := &FlowAnalytics{FlowID: flowID, NodeHits: map[string]int{}}
	if flowDB != nil { a.TriggerCount, a.CompletionCount, a.NodeHits = flowDB.GetFlowStats(flowID) }
	return a
}
func incFlowCounter(key string) int {
	if flowDB != nil { return flowDB.IncFlowCounter(key) }
	return 0
}

// ════════════════ Flow Resume (paused flows - DB persisted) ════════════════
var pausedFlowsMem = make(map[string]*FlowContext)

func ResumeFlow(phone, message string) ([]FlowReply, []FlowAction, bool) {
	ctx, ok := pausedFlowsMem[phone]
	if !ok {
		if flowDB == nil { return nil, nil, false }
		fid, accPhone, varsJSON, currentNode, visitedJSON, contactName, pausedMsg, err := flowDB.LoadPausedFlow(phone)
		if err != nil || fid == 0 { return nil, nil, false }
		flow := getFlowDB(fid, 0)
		if flow == nil { flowDB.DeletePausedFlow(phone); return nil, nil, false }
		ctx = NewFlowContext(flow, phone, accPhone, message, contactName)
		if pausedMsg != "" { ctx.Message = pausedMsg }
		json.Unmarshal([]byte(varsJSON), &ctx.Variables)
		json.Unmarshal([]byte(visitedJSON), &ctx.VisitedNodes)
		ctx.Variables["__paused_at"] = currentNode
		flowDB.DeletePausedFlow(phone)
	} else {
		delete(pausedFlowsMem, phone)
	}

	varName := ctx.Variables["__paused_var"]
	// Handle media: [media:type:url] prefix
	isMedia := ctx.Variables["__paused_media"] == "true"
	if isMedia && strings.HasPrefix(message, "[media:") {
		if idx := strings.LastIndex(message, "]"); idx > 7 {
			mediaInfo := message[7:idx]
			parts := strings.SplitN(mediaInfo, ":", 2)
			if len(parts) >= 2 && varName != "" {
				ctx.Variables[varName] = parts[1]
				ctx.Variables[varName+"_type"] = parts[0]
			}
		}
	} else if varName != "" {
		ctx.Variables[varName] = message
	}
	nextID := ctx.Variables["__paused_at"]
	if nextID != "" { ctx.CurrentNode = getNextNode(ctx.Flow, nextID, "") }
	ctx.Done = false
	replies, actions := ExecuteFlow(ctx)
	return replies, actions, true
}

func pauseFlow(phone string, ctx *FlowContext) {
	pausedFlowsMem[phone] = ctx
	// Also persist to DB for server restart safety
	if flowDB != nil {
		varsJSON, _ := json.Marshal(ctx.Variables)
		visitedJSON, _ := json.Marshal(ctx.VisitedNodes)
		flowDB.SavePausedFlow(phone, ctx.Flow.ID, ctx.Flow.Trigger.AccountPhone,
			string(varsJSON), ctx.CurrentNode, string(visitedJSON),
			ctx.ContactName, ctx.Message)
	}
}

// Flow Tags
type FlowTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Flow rate limiter: per-flow per-phone per-hour, persisted in DB
func checkFlowRateLimit(flowID int64, phone string, maxPerHour int) bool {
	if maxPerHour <= 0 { return true }
	key := fmt.Sprintf("flow_rate_%d_%s_%s", flowID, phone, time.Now().Format("2006-01-02-15"))
	if flowDB != nil {
		count := flowDB.IncFlowCounterWithExpiry(key, 1)
		return count <= maxPerHour
	}
	return true
}
func GetNodeMetrics(flowID int64) map[string]int {
	if flowDB != nil { _, _, nh := flowDB.GetFlowStats(flowID); return nh }
	return map[string]int{}
}

// ════════════════ Node Executors ════════════════
func executeMessage(node *FlowNode, ctx *FlowContext) {
	var d MessageData; json.Unmarshal(node.Data, &d)
	t, m := renderVar(d.Text, ctx), renderVar(d.MediaURL, ctx)
	if t != "" || m != "" { ctx.Replies = append(ctx.Replies, FlowReply{Text: t, MediaURL: m, MediaType: d.MediaType}) }
	ctx.CurrentNode = next(ctx, "")
}
func executeQuestion(node *FlowNode, ctx *FlowContext) {
	var d QuestionData; json.Unmarshal(node.Data, &d)
	if d.Question != "" { ctx.Replies = append(ctx.Replies, FlowReply{Text: renderVar(d.Question, ctx)}) }
	if d.VarName != "" { ctx.Variables["__paused_at"] = node.ID; ctx.Variables["__paused_var"] = d.VarName; pauseFlow(ctx.Phone, ctx) }
	ctx.Done = true
}
func executeCondition(node *FlowNode, ctx *FlowContext) {
	var d ConditionData; json.Unmarshal(node.Data, &d)
	fv := resolveVar(d.Variable, ctx)
	for i, br := range d.Branches {
		if evalOp(fv, br.Operator, br.Value) { ctx.CurrentNode = next(ctx, fmt.Sprintf("out-%d", i)); return }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeWait(node *FlowNode, ctx *FlowContext) {
	var d WaitData; json.Unmarshal(node.Data, &d)
	ctx.Actions = append(ctx.Actions, FlowAction{Type: "wait", Params: map[string]string{"seconds": fmt.Sprintf("%d", d.Seconds)}})
	ctx.CurrentNode = next(ctx, "")
}
func executeSetVar(node *FlowNode, ctx *FlowContext) {
	var d SetVarData; json.Unmarshal(node.Data, &d)
	for k, v := range d.Variables { ctx.Variables[k] = renderVar(v, ctx) }
	ctx.CurrentNode = next(ctx, "")
}
func executeTag(node *FlowNode, ctx *FlowContext) {
	var d TagData; json.Unmarshal(node.Data, &d)
	for _, tag := range d.Tags { ctx.Actions = append(ctx.Actions, FlowAction{Type: "tag", Params: map[string]string{"tag": tag, "phone": ctx.Phone}}) }
	ctx.CurrentNode = next(ctx, "")
}
func executeAPICall(node *FlowNode, ctx *FlowContext) {
	var d APICallData; json.Unmarshal(node.Data, &d)
	ctx.Actions = append(ctx.Actions, FlowAction{Type: "api_call", Params: map[string]string{"url": d.URL, "method": d.Method, "var_name": d.VarName}})
	ctx.CurrentNode = next(ctx, "")
}
func executeTransfer(node *FlowNode, ctx *FlowContext) { ctx.Actions = append(ctx.Actions, FlowAction{Type: "transfer_agent"}); ctx.Done = true }
func executeClose(node *FlowNode, ctx *FlowContext)   { ctx.Actions = append(ctx.Actions, FlowAction{Type: "close_chat"}); ctx.Done = true }
func executeAI(node *FlowNode, ctx *FlowContext) {
	var d AIData; json.Unmarshal(node.Data, &d)
	prompt := renderVar(d.SystemPrompt, ctx)
	if prompt == "" { prompt = "You are a helpful WhatsApp assistant. Reply in Indonesian." }
	if AICallback != nil {
		reply, _ := AICallback(prompt, ctx.Message, nil, ctx.Flow.AIKeyID)
		if reply != "" { ctx.Replies = append(ctx.Replies, FlowReply{Text: reply}); if d.VarResult != "" { ctx.Variables[d.VarResult] = reply } }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeAIDecide(node *FlowNode, ctx *FlowContext) {
	var d AIDecideData; json.Unmarshal(node.Data, &d)
	prompt := renderVar(d.SystemPrompt, ctx)
	if prompt == "" { prompt = "Classify the user message." }
	if AICallback != nil && len(d.Options) > 0 {
		_, idx := AICallback(prompt, ctx.Message, d.Options, ctx.Flow.AIKeyID)
		if idx >= 0 && idx < len(d.Options) { ctx.CurrentNode = next(ctx, fmt.Sprintf("out-%d", idx)); if d.VarResult != "" { ctx.Variables[d.VarResult] = d.Options[idx] }; return }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeCarousel(node *FlowNode, ctx *FlowContext) {
	var d CarouselData; json.Unmarshal(node.Data, &d)
	if d.MaxItems <= 0 { d.MaxItems = 5 }
	if StoreListProducts != nil {
		products := StoreListProducts(d.Category, d.MaxItems)
		if len(products) > 0 {
			text := "🛍️ *" + d.Category + "*\n\n"
			for i, p := range products { text += fmt.Sprintf("%d. %s — %s\n", i+1, p["name"], p["price"]) }
			text += "\nKetik nomor untuk pesan."
			ctx.Replies = append(ctx.Replies, FlowReply{Text: text})
		}
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeOrderCreate(node *FlowNode, ctx *FlowContext) {
	var d OrderCreateData; json.Unmarshal(node.Data, &d)
	p, qs := renderVar(d.ProductVar, ctx), renderVar(d.QtyVar, ctx)
	q, _ := strconv.Atoi(qs); if q <= 0 { q = 1 }
	if StoreCreateOrder != nil {
		oid, _ := StoreCreateOrder(ctx.Phone, p, q)
		ctx.Replies = append(ctx.Replies, FlowReply{Text: fmt.Sprintf("✅ Pesanan #%d: %s x%d", oid, p, q)})
		if d.VarResult != "" { ctx.Variables[d.VarResult] = fmt.Sprintf("%d", oid) }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executePaymentLink(node *FlowNode, ctx *FlowContext) {
	var d PaymentLinkData; json.Unmarshal(node.Data, &d)
	oid, _ := strconv.ParseInt(renderVar(d.OrderVar, ctx), 10, 64)
	amt, _ := strconv.ParseFloat(renderVar(d.AmountVar, ctx), 64)
	if StorePaymentLink != nil {
		url := StorePaymentLink(oid, amt)
		ctx.Replies = append(ctx.Replies, FlowReply{Text: fmt.Sprintf("💳 Silakan bayar: %s", url)})
		if d.VarResult != "" { ctx.Variables[d.VarResult] = url }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeMath(node *FlowNode, ctx *FlowContext) {
	var d MathData; json.Unmarshal(node.Data, &d)
	r := evalMath(renderVar(d.Formula, ctx))
	if d.VarResult != "" { ctx.Variables[d.VarResult] = fmt.Sprintf("%.2f", r) }
	ctx.CurrentNode = next(ctx, "")
}
func executeStringTemplate(node *FlowNode, ctx *FlowContext) {
	var d StringTemplateData; json.Unmarshal(node.Data, &d)
	r := renderVar(d.Template, ctx)
	if d.VarResult != "" { ctx.Variables[d.VarResult] = r }
	ctx.CurrentNode = next(ctx, "")
}
func executeDateMath(node *FlowNode, ctx *FlowContext) {
	var d DateMathData; json.Unmarshal(node.Data, &d)
	var r string; now := time.Now()
	switch d.Operation {
	case "add_days": r = now.AddDate(0, 0, d.Value).Format("2006-01-02")
	case "add_hours": r = now.Add(time.Duration(d.Value) * time.Hour).Format("2006-01-02 15:04")
	case "now": r = now.Format("2006-01-02 15:04:05")
	default: r = now.Format("2006-01-02")
	}
	if d.VarResult != "" { ctx.Variables[d.VarResult] = r }
	ctx.CurrentNode = next(ctx, "")
}
func executeLoop(node *FlowNode, ctx *FlowContext) {
	var d LoopData; json.Unmarshal(node.Data, &d)
	if d.MaxIterations <= 0 { d.MaxIterations = 3 }
	key := fmt.Sprintf("__loop_%s", node.ID)
	iter, _ := strconv.Atoi(ctx.Variables[key]); iter++
	ctx.Variables[key] = fmt.Sprintf("%d", iter)
	if iter >= d.MaxIterations || (d.ConditionVar != "" && ctx.Variables[d.ConditionVar] == "") { ctx.CurrentNode = next(ctx, "done"); return }
	if d.LoopStartNode != "" { ctx.CurrentNode = d.LoopStartNode } else { ctx.CurrentNode = next(ctx, "loop") }
}
func executeSubflow(node *FlowNode, ctx *FlowContext) {
	var d SubflowData; json.Unmarshal(node.Data, &d)
	if d.FlowID == 0 { ctx.CurrentNode = next(ctx, ""); return }
	sf := getFlowDB(d.FlowID, ctx.Flow.UserID)
	if sf == nil { ctx.CurrentNode = next(ctx, ""); return }
	sc := NewFlowContext(sf, ctx.Phone, ctx.AccountPhone, ctx.Message, ctx.ContactName)
	for k, v := range d.Inputs { sc.Variables[k] = renderVar(v, ctx) }
	replies, _ := ExecuteFlow(sc)
	ctx.Replies = append(ctx.Replies, replies...)
	for ok, sv := range d.Outputs { if val, ok2 := sc.Variables[ok]; ok2 { ctx.Variables[sv] = val } }
	ctx.CurrentNode = next(ctx, "")
}
func executeRandom(node *FlowNode, ctx *FlowContext) {
	var d RandomData; json.Unmarshal(node.Data, &d)
	if len(d.Weights) == 0 { d.Weights = []int{1, 1} }
	total := 0; for _, w := range d.Weights { total += w }
	r := int(time.Now().UnixNano()) % total
	cumulative := 0
	for i, w := range d.Weights {
		cumulative += w
		if r < cumulative { ctx.CurrentNode = next(ctx, fmt.Sprintf("out-%d", i)); return }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeContactLookup(node *FlowNode, ctx *FlowContext) {
	var d ContactLookupData; json.Unmarshal(node.Data, &d)
	if ContactLookupFunc != nil {
		data := ContactLookupFunc(ctx.Phone)
		for k, v := range data { ctx.Variables["contact."+k] = v }
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeCounter(node *FlowNode, ctx *FlowContext) {
	var d CounterData; json.Unmarshal(node.Data, &d)
	if d.Key == "" { d.Key = fmt.Sprintf("flow_%s_%s", node.ID, ctx.Phone) }
	count := incFlowCounter(d.Key)
	if d.VarResult != "" { ctx.Variables[d.VarResult] = fmt.Sprintf("%d", count) }
	if d.MaxCount > 0 && count >= d.MaxCount { ctx.CurrentNode = next(ctx, "overflow"); return }
	ctx.CurrentNode = next(ctx, "")
}
func executeSequence(node *FlowNode, ctx *FlowContext) {
	var d SequenceData; json.Unmarshal(node.Data, &d)
	for _, m := range d.Messages { ctx.Replies = append(ctx.Replies, FlowReply{Text: renderVar(m.Text, ctx)}) }
	ctx.CurrentNode = next(ctx, "")
}
func executePoll(node *FlowNode, ctx *FlowContext) {
	var d PollData; json.Unmarshal(node.Data, &d)
	text := "📊 *" + renderVar(d.Question, ctx) + "*\n\n"
	for i, o := range d.Options { text += fmt.Sprintf("%d. %s\n", i+1, o) }
	text += "\nBalas nomor pilihan (1-" + fmt.Sprintf("%d", len(d.Options)) + ")."
	ctx.Replies = append(ctx.Replies, FlowReply{Text: text})
	if d.VarResult != "" {
		ctx.Variables["__paused_at"] = node.ID
		ctx.Variables["__paused_var"] = d.VarResult
		pauseFlow(ctx.Phone, ctx)
	}
	ctx.Done = true
}
func executeLocation(node *FlowNode, ctx *FlowContext) {
	var d LocationData; json.Unmarshal(node.Data, &d)
	if d.Action == "send" { ctx.Replies = append(ctx.Replies, FlowReply{Text: fmt.Sprintf("📍 %s\nhttps://maps.google.com/?q=%.6f,%.6f", d.Name, d.Lat, d.Lng)}) }
	if d.Action == "request" { ctx.Replies = append(ctx.Replies, FlowReply{Text: "📍 Silakan share lokasi Anda."}) }
	ctx.CurrentNode = next(ctx, "")
}
func executeDocument(node *FlowNode, ctx *FlowContext) {
	var d DocumentData; json.Unmarshal(node.Data, &d)
	ctx.Replies = append(ctx.Replies, FlowReply{Text: renderVar(d.Caption, ctx), MediaURL: d.URL, MediaType: "document"})
	ctx.CurrentNode = next(ctx, "")
}
func executeTyping(node *FlowNode, ctx *FlowContext) {
	var d TypingData; json.Unmarshal(node.Data, &d)
	if d.Seconds > 0 { time.Sleep(time.Duration(d.Seconds) * time.Second) }
	ctx.CurrentNode = next(ctx, "")
}
func executeDelayUntil(node *FlowNode, ctx *FlowContext) {
	var d DelayUntilData; json.Unmarshal(node.Data, &d)
	if d.Time != "" {
		parts := strings.Split(d.Time, ":")
		if len(parts) == 2 {
			h, _ := strconv.Atoi(parts[0]); m, _ := strconv.Atoi(parts[1])
			now := time.Now(); target := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())
			if target.Before(now) { target = target.Add(24 * time.Hour) }
			if delay := target.Sub(now); delay > 0 && delay < 24*time.Hour {
				ctx.Actions = append(ctx.Actions, FlowAction{Type: "wait", Params: map[string]string{"seconds": fmt.Sprintf("%d", int(delay.Seconds()))}})
			}
		}
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeGallery(node *FlowNode, ctx *FlowContext) {
	var d GalleryData; json.Unmarshal(node.Data, &d)
	caption := renderVar(d.Caption, ctx)
	for i, img := range d.Images {
		cap := renderVar(img.Caption, ctx)
		if i == 0 && caption != "" { cap = caption }
		ctx.Replies = append(ctx.Replies, FlowReply{Text: cap, MediaURL: img.URL, MediaType: "image"})
	}
	ctx.CurrentNode = next(ctx, "")
}
func executeVoice(node *FlowNode, ctx *FlowContext) {
	var d VoiceData; json.Unmarshal(node.Data, &d)
	if d.TTS != "" { ctx.Replies = append(ctx.Replies, FlowReply{Text: renderVar(d.TTS, ctx)}) }
	if d.URL != "" { ctx.Replies = append(ctx.Replies, FlowReply{MediaURL: d.URL, MediaType: "audio"}) }
	ctx.CurrentNode = next(ctx, "")
}
func executeVCard(node *FlowNode, ctx *FlowContext) {
	var d VCardData; json.Unmarshal(node.Data, &d)
	ctx.Replies = append(ctx.Replies, FlowReply{Text: fmt.Sprintf("📇 *%s*\n📞 %s\n🏢 %s", d.Name, d.Phone, d.Org)})
	ctx.Actions = append(ctx.Actions, FlowAction{Type: "send_vcard", Params: map[string]string{"name": d.Name, "phone": d.Phone, "org": d.Org}})
	ctx.CurrentNode = next(ctx, "")
}
func executeDBQuery(node *FlowNode, ctx *FlowContext) {
	var d DBQueryData; json.Unmarshal(node.Data, &d)
	q := renderVar(d.Query, ctx)
	if DBQueryFunc != nil && q != "" {
		if r, err := DBQueryFunc(q); err == nil && d.VarResult != "" { ctx.Variables[d.VarResult] = r }
	}
	ctx.CurrentNode = next(ctx, "")
}

func executeReadReceipt(node *FlowNode, ctx *FlowContext) {
	ctx.Actions = append(ctx.Actions, FlowAction{Type: "wait_for_read"})
	ctx.CurrentNode = next(ctx, "")
}

func executeSticker(node *FlowNode, ctx *FlowContext) {
	var d StickerData; json.Unmarshal(node.Data, &d)
	ctx.Replies = append(ctx.Replies, FlowReply{Text: d.Text, MediaURL: d.URL, MediaType: "sticker"})
	ctx.CurrentNode = next(ctx, "")
}

func executeGoogleSheets(node *FlowNode, ctx *FlowContext) {
	var d GoogleSheetsData; json.Unmarshal(node.Data, &d)
	data := renderVar(d.Data, ctx)
	if GoogleSheetsFunc != nil {
		result := GoogleSheetsFunc(d.Action, d.SheetID, d.Range, data)
		if d.VarResult != "" { ctx.Variables[d.VarResult] = result }
	}
	ctx.CurrentNode = next(ctx, "")
}

func executeWarmer(node *FlowNode, ctx *FlowContext) {
	// Number warmer: sends periodic messages to keep WA account active
	ctx.Actions = append(ctx.Actions, FlowAction{
		Type: "warmer",
		Params: map[string]string{"phone": ctx.Phone},
	})
	ctx.CurrentNode = next(ctx, "")
}

func executeTGKeyboard(node *FlowNode, ctx *FlowContext) {
	var d map[string]interface{}; json.Unmarshal(node.Data, &d)
	text := ""; if v, ok := d["text"].(string); ok { text = renderVar(v, ctx) }
	btns := []string{}; if arr, ok := d["buttons"].([]interface{}); ok { for _, a := range arr { if s, ok := a.(string); ok { btns = append(btns, s) } } }
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "tg_keyboard", Text: text, ActionData: map[string]interface{}{"buttons": btns}})
	ctx.CurrentNode = next(ctx, "")
}
func executeTGPayment(node *FlowNode, ctx *FlowContext) {
	var d map[string]interface{}; json.Unmarshal(node.Data, &d)
	title := ""; if v, ok := d["title"].(string); ok { title = renderVar(v, ctx) }
	amount := "0"; if v, ok := d["amount"].(string); ok { amount = v }
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "tg_payment", Text: title, ActionData: map[string]interface{}{"amount": amount, "currency": "IDR"}})
	ctx.CurrentNode = next(ctx, "")
}
func executeIGComment(node *FlowNode, ctx *FlowContext) {
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "ig_comment", Text: `Comment detected — sending DM`})
	ctx.CurrentNode = next(ctx, "")
}
func executeIGStoryReply(node *FlowNode, ctx *FlowContext) {
	var d map[string]interface{}; json.Unmarshal(node.Data, &d)
	text := ""; if v, ok := d["text"].(string); ok { text = renderVar(v, ctx) }
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "ig_story_reply", Text: text})
	ctx.CurrentNode = next(ctx, "")
}
func executeIGQuickReply(node *FlowNode, ctx *FlowContext) {
	var d map[string]interface{}; json.Unmarshal(node.Data, &d)
	text := ""; if v, ok := d["text"].(string); ok { text = renderVar(v, ctx) }
	replies := []string{}; if arr, ok := d["replies"].([]interface{}); ok { for _, a := range arr { if s, ok := a.(string); ok { replies = append(replies, s) } } }
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "ig_quick_reply", Text: text, ActionData: map[string]interface{}{"replies": replies}})
	ctx.CurrentNode = next(ctx, "")
}
func executeFBCarousel(node *FlowNode, ctx *FlowContext) {
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "fb_carousel", Text: "Carousel products"})
	ctx.CurrentNode = next(ctx, "")
}
func executeFBMenu(node *FlowNode, ctx *FlowContext) {
	ctx.Actions = append(ctx.Actions, FlowAction{Type: "fb_menu", Params: map[string]string{}})
	ctx.CurrentNode = next(ctx, "")
}
func executeFBWebview(node *FlowNode, ctx *FlowContext) {
	var d map[string]interface{}; json.Unmarshal(node.Data, &d)
	url := ""; if v, ok := d["url"].(string); ok { url = renderVar(v, ctx) }
	ctx.Replies = append(ctx.Replies, FlowReply{Action: "fb_webview", Text: url})
	ctx.CurrentNode = next(ctx, "")
}

func executeButtons(node *FlowNode, ctx *FlowContext) {
	// CTA Buttons: title + footer + buttons array
	var d map[string]interface{}
	json.Unmarshal(node.Data, &d)
	title := ""
	footer := ""
	btns := []string{}
	if v, ok := d["title"].(string); ok { title = renderVar(v, ctx) }
	if v, ok := d["footer"].(string); ok { footer = renderVar(v, ctx) }
	if arr, ok := d["buttons"].([]interface{}); ok {
		for _, a := range arr {
			if s, ok := a.(string); ok { btns = append(btns, renderVar(s, ctx)) }
		}
	}
	if len(btns) > 0 {
		ctx.Replies = append(ctx.Replies, FlowReply{
			Text: title + "\n" + footer,
			Action: "buttons",
			ActionData: map[string]interface{}{"title": title, "footer": footer, "buttons": btns},
		})
	}
	ctx.CurrentNode = next(ctx, "")
}

func executeReceiveMedia(node *FlowNode, ctx *FlowContext) {
	var d ReceiveMediaData; json.Unmarshal(node.Data, &d)
	prompt := renderVar(d.Prompt, ctx)
	if prompt == "" { prompt = "📎 Silakan kirim foto/dokumen." }
	ctx.Replies = append(ctx.Replies, FlowReply{Text: prompt})
	if d.VarResult != "" {
		ctx.Variables["__paused_at"] = node.ID
		ctx.Variables["__paused_var"] = d.VarResult
		ctx.Variables["__paused_media"] = "true"
		pauseFlow(ctx.Phone, ctx)
	}
	ctx.Done = true
}

// ════════════════ Graph Navigation ════════════════
func next(ctx *FlowContext, handle string) string   { return getNextNode(ctx.Flow, ctx.CurrentNode, handle) }
func findNode(flow *Flow, id string) *FlowNode      { for i := range flow.Nodes { if flow.Nodes[i].ID == id { return &flow.Nodes[i] } }; return nil }
func findStartNode(flow *Flow) *FlowNode {
	hasIncoming := map[string]bool{}
	for _, e := range flow.Edges { hasIncoming[e.Target] = true }
	for i := range flow.Nodes { if !hasIncoming[flow.Nodes[i].ID] { return &flow.Nodes[i] } }
	if len(flow.Nodes) > 0 { return &flow.Nodes[0] }; return nil
}
func getNextNode(flow *Flow, sourceID, handle string) string {
	for _, e := range flow.Edges {
		if e.Source == sourceID && (handle == "" || e.SourceH == handle) && !e.Fallback { return e.Target }
	}
	for _, e := range flow.Edges {
		if e.Source == sourceID && e.Fallback { return e.Target }
	}
	return ""
}

// ════════════════ Flow Simulate (for testing in Flow Builder) ════════════════
type SimulateResult struct {
	Replies     []FlowReply              `json:"replies"`
	Actions     []FlowAction             `json:"actions"`
	Variables   map[string]string        `json:"variables"`
	NodesVisited int                     `json:"nodes_visited"`
	Status      string                   `json:"status"`
	Error       string                   `json:"error,omitempty"`
	Trace       []SimulateTrace          `json:"trace,omitempty"`
}
type SimulateTrace struct {
	NodeID   string `json:"node_id"`
	NodeType string `json:"node_type"`
	Label    string `json:"label"`
}

func SimulateFlow(flow *Flow, phone, message, contactName string, debug bool) *SimulateResult {
	result := &SimulateResult{
		Variables: make(map[string]string),
		Trace:     make([]SimulateTrace, 0),
	}
	ctx := NewFlowContext(flow, phone, phone, message, contactName)
	ctx.MaxSteps = 100

	start := findStartNode(flow)
	if start == nil {
		result.Status = "error"
		result.Error = "No start node found"
		return result
	}
	ctx.CurrentNode = start.ID

	for !ctx.Done && len(ctx.VisitedNodes) < ctx.MaxSteps {
		node := findNode(flow, ctx.CurrentNode)
		if node == nil {
			break
		}
		ctx.VisitedNodes[ctx.CurrentNode]++
		if ctx.VisitedNodes[ctx.CurrentNode] > 5 {
			break
		}

		if debug {
			result.Trace = append(result.Trace, SimulateTrace{
				NodeID: node.ID, NodeType: string(node.Type),
				Label: getNodeLabelForTrace(node),
			})
		}

		executeNode(node, ctx)
	}

	result.Replies = ctx.Replies
	result.Actions = ctx.Actions
	result.Variables = ctx.Variables
	result.NodesVisited = len(ctx.VisitedNodes)
	if result.Status == "" {
		result.Status = "completed"
	}
	return result
}

func getNodeLabelForTrace(node *FlowNode) string {
	label := string(node.Type)
	switch node.Type {
	case NodeMessage:
		var d MessageData
		json.Unmarshal(node.Data, &d)
		if d.Text != "" {
			label = d.Text
		}
	case NodeCondition:
		var d ConditionData
		json.Unmarshal(node.Data, &d)
		label = "IF " + d.Variable
	case NodeQuestion:
		var d QuestionData
		json.Unmarshal(node.Data, &d)
		label = d.Question
	}
	if len(label) > 40 {
		label = label[:40] + "..."
	}
	return label
}

// ════════════════ Helpers ════════════════
func renderVar(text string, ctx *FlowContext) string {
	r := text
	r = strings.ReplaceAll(r, "{{input}}", ctx.Message)
	r = strings.ReplaceAll(r, "{{phone}}", ctx.Phone)
	r = strings.ReplaceAll(r, "{{name}}", ctx.ContactName)
	for k, v := range ctx.Variables { r = strings.ReplaceAll(r, "{{"+k+"}}", v) }
	return spintax(r)
}

// Spintax: {Halo|Hai|Hi} → random choice
func spintax(text string) string {
	for {
		start := strings.Index(text, "{")
		if start < 0 { break }
		end := strings.Index(text[start:], "}")
		if end < 0 { break }
		end += start
		inner := text[start+1 : end]
		options := strings.Split(inner, "|")
		if len(options) == 0 { break }
		pick := options[int(time.Now().UnixNano())%len(options)]
		text = text[:start] + pick + text[end+1:]
	}
	return text
}
func resolveVar(name string, ctx *FlowContext) string {
	n := strings.TrimSpace(name)
	switch {
	case n == "{{input}}" || n == "input": return ctx.Message
	case n == "{{phone}}" || n == "phone": return ctx.Phone
	case n == "{{name}}" || n == "name": return ctx.ContactName
	case strings.HasPrefix(n, "{{") && strings.HasSuffix(n, "}}"): vn := n[2 : len(n)-2]; if v, ok := ctx.Variables[vn]; ok { return v }; return ctx.Message
	default: return name
	}
}
func evalOp(fv, op, val string) bool {
	a, b := strings.ToLower(strings.TrimSpace(fv)), strings.ToLower(strings.TrimSpace(val))
	switch op {
	case "equals": return a == b
	case "contains": return strings.Contains(a, b)
	case "starts_with": return strings.HasPrefix(a, b)
	case "empty": return a == ""
	case "not_empty": return a != ""
	default: return a == b
	}
}
func evalMath(expr string) float64 {
	expr = strings.ReplaceAll(expr, " ", "")
	var nums []float64; var ops []byte; i := 0
	for i < len(expr) {
		if (expr[i] >= '0' && expr[i] <= '9') || expr[i] == '.' {
			j := i; for j < len(expr) && ((expr[j] >= '0' && expr[j] <= '9') || expr[j] == '.') { j++ }
			v, _ := strconv.ParseFloat(expr[i:j], 64); nums = append(nums, v); i = j
		} else if expr[i] == '+' || expr[i] == '-' || expr[i] == '*' || expr[i] == '/' { ops = append(ops, expr[i]); i++ } else { i++ }
	}
	if len(nums) == 1 { return nums[0] }
	if len(nums) < 2 || len(ops) == 0 { return 0 }
	r := nums[0]; for i, op := range ops { n := nums[i+1]; switch op { case '+': r += n; case '-': r -= n; case '*': r *= n; case '/': if n != 0 { r /= n } } }
	return r
}
func splitTrim(s, sep string) []string {
	var r []string; for _, p := range strings.Split(s, sep) { p = strings.TrimSpace(p); if p != "" { r = append(r, p) } }; return r
}
