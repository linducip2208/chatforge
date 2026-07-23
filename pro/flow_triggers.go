//go:build pro

package pro

import "log"

func HandleIncomingMessage(uid int64, accountPhone, phone, message, contactName string) ([]FlowReply, []FlowAction, bool) {
	// Check paused flows first — resume if user was in middle of a flow
	if replies, actions, ok := ResumeFlow(phone, message); ok {
		log.Printf("[FLOW] Resumed: %s → %s", phone, accountPhone)
		return replies, actions, true
	}

	flow := FindMatchingFlow(uid, TriggerKeyword, accountPhone, phone, message)
	if flow == nil {
		flow = FindMatchingFlow(uid, TriggerAlways, accountPhone, phone, message)
	}
	if flow == nil {
		return nil, nil, false
	}
	ctx := NewFlowContext(flow, phone, accountPhone, message, contactName)
	replies, actions := ExecuteFlow(ctx)
	log.Printf("[FLOW] %s → %s (via %s): %d replies", flow.Name, phone, accountPhone, len(replies))
	return replies, actions, true
}

func HandleDripFlow(uid int64, accountPhone, phone, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerDrip, accountPhone) {
		ctx := NewFlowContext(flow, phone, accountPhone, "", contactName)
		replies, actions := ExecuteFlow(ctx)
		if len(replies) > 0 || len(actions) > 0 {
			return replies, actions
		}
	}
	return nil, nil
}

func HandleWelcomeFlow(uid int64, accountPhone, phone, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerWelcome, accountPhone) {
		ctx := NewFlowContext(flow, phone, accountPhone, "", contactName)
		replies, actions := ExecuteFlow(ctx)
		if len(replies) > 0 || len(actions) > 0 {
			return replies, actions
		}
	}
	return nil, nil
}

func HandleBroadcastFlow(uid int64, accountPhone, phone, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerBroadcast, accountPhone) {
		ctx := NewFlowContext(flow, phone, accountPhone, "", contactName)
		replies, actions := ExecuteFlow(ctx)
		if len(replies) > 0 || len(actions) > 0 {
			return replies, actions
		}
	}
	return nil, nil
}

func HandleCSATFlow(uid int64, accountPhone, phone, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerCSAT, accountPhone) {
		ctx := NewFlowContext(flow, phone, accountPhone, "", contactName)
		replies, actions := ExecuteFlow(ctx)
		if len(replies) > 0 || len(actions) > 0 {
			return replies, actions
		}
	}
	return nil, nil
}

func HandleCloseFlow(uid int64, accountPhone, phone, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerClose, accountPhone) {
		ctx := NewFlowContext(flow, phone, accountPhone, "", contactName)
		replies, actions := ExecuteFlow(ctx)
		if len(replies) > 0 || len(actions) > 0 {
			return replies, actions
		}
	}
	return nil, nil
}

func HandleFallbackFlow(uid int64, accountPhone, phone, message, contactName string) ([]FlowReply, []FlowAction) {
	for _, flow := range FindFlowsByTrigger(uid, TriggerFallback, accountPhone) {
		if matchTriggerValue(flow.Trigger, phone, message) {
			ctx := NewFlowContext(flow, phone, accountPhone, message, contactName)
			return ExecuteFlow(ctx)
		}
	}
	return nil, nil
}
