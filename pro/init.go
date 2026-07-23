//go:build pro

package pro

import (
	"log"
	"net/http"
)

// RegisterRoutes registers all Pro feature routes.
func RegisterRoutes(mux *http.ServeMux) {
	log.Println("[PRO] Registering Pro feature routes...")
	count := 0

	// ═══ Flow Builder ═══
	mux.HandleFunc("/pro/flow-builder", handleFlowBuilderUI)
	mux.HandleFunc("/pro/flow-builder/save", handleFlowSave)
	mux.HandleFunc("/pro/flow-builder/list", handleFlowList)
	mux.HandleFunc("/pro/flow-builder/load", handleFlowLoad)
	mux.HandleFunc("/pro/flow-builder/delete", handleFlowDelete)
	mux.HandleFunc("/pro/flow-builder/toggle", handleFlowToggle)
	mux.HandleFunc("/pro/flow-builder/duplicate", handleFlowDuplicate)
	mux.HandleFunc("/pro/flow-builder/templates", handleFlowTemplates)
	mux.HandleFunc("/pro/flow-builder/templates/import", handleFlowTemplateImport)
	mux.HandleFunc("/pro/flow/trigger/", handleFlowTriggerWebhook)
	mux.HandleFunc("/pro/flow-builder/export", handleFlowExport)
	mux.HandleFunc("/pro/flow-builder/import", handleFlowImport)
	mux.HandleFunc("/pro/flow-builder/analytics", handleFlowAnalytics)
	mux.HandleFunc("/pro/flow-builder/versions", handleFlowVersions)
	mux.HandleFunc("/pro/flow-builder/rollback", handleFlowRollback)
	mux.HandleFunc("/pro/flow-builder/publish", handleFlowPublish)
	mux.HandleFunc("/pro/flow-builder/unpublish", handleFlowUnpublish)
	mux.HandleFunc("/pro/flow-builder/marketplace", handleFlowMarketplace)
	mux.HandleFunc("/pro/flow-builder/ai-generate", handleFlowAIGenerate)
	mux.HandleFunc("/pro/flow-builder/accounts", handleFlowAccounts)
	mux.HandleFunc("/pro/flow-builder/logs", handleFlowLogs)
	mux.HandleFunc("/pro/flow-builder/simulate", handleFlowSimulate)
	mux.HandleFunc("/pro/flow-builder/debug", handleFlowDebug)
	mux.HandleFunc("/pro/flow-builder/reviews", handleFlowReviews)
	mux.HandleFunc("/pro/flow-builder/download", handleFlowDownload)
	mux.HandleFunc("/pro/flow-builder/ai-keys", handleFlowAIKeys)
	count += 26

	// ═══ Omnichannel ═══
	mux.HandleFunc("/pro/instagram/webhook", HandleInstagramWebhook)
	mux.HandleFunc("/pro/instagram/inbox", HandleInstagramInbox)
	mux.HandleFunc("/pro/facebook/webhook", HandleFacebookWebhook)
	mux.HandleFunc("/pro/facebook/inbox", HandleFacebookInbox)
	mux.HandleFunc("/pro/telegram/webhook", HandleTelegramWebhook)
	mux.HandleFunc("/pro/telegram/inbox", HandleTelegramInbox)
	count += 6

	// ═══ Advanced Messaging ═══
	mux.HandleFunc("/pro/message-status/", TrackMessageStatus)
	mux.HandleFunc("/pro/send/buttons", SendMessageWithButtons)
	mux.HandleFunc("/pro/check-number", HandleWANumberCheck)
	count += 3

	// ═══ Omnichannel Inbox ═══
	mux.HandleFunc("/pro/omni/inbox", handleOmniInboxUI)
	mux.HandleFunc("/pro/omni/events", HandleOmnichannelEvents)
	mux.HandleFunc("/pro/omni/send", HandleOmnichannelSend)
	mux.HandleFunc("/pro/omni/conversations", handleOmniConversations)
	mux.HandleFunc("/pro/omni/messages", handleOmniMessages)
	count += 5

	// ═══ Integrations ═══
	mux.HandleFunc("/pro/webhooks/n8n/", HandleN8nWebhook)
	mux.HandleFunc("/pro/webhooks/zapier/", HandleZapierWebhook)
	mux.HandleFunc("/pro/webhooks/make/", HandleZapierWebhook)
	mux.HandleFunc("/pro/sheets/sync", HandleGoogleSheetsSync)
	count += 4

	// ═══ Agency Dashboard ═══
	mux.HandleFunc("/pro/agency", handleAgencyUI)
	mux.HandleFunc("/pro/agency/clients", handleAgencyClients)
	mux.HandleFunc("/pro/agency/clients/add", HandleAgencyClientAdd)
	mux.HandleFunc("/pro/agency/clients/delete", HandleAgencyClientDelete)
	count += 4

	log.Printf("[PRO] %d routes registered", count)
}
