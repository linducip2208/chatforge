//go:build pro

package main

import (
	"chatgo/aiservice"
	"chatgo/pro"
	"chatgo/secret"
	"chatgo/wa"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func isProBuild() bool { return true }

func initProRoutes(mux *http.ServeMux) {
	pro.RegisterRoutes(mux)
}

func setupProEngine() {
	pro.SetFlowDB(db)

	// Meta flow callback
	wa.MetaFlowCallback = func(uid int64, accountPhone, phone, message, name string) ([]wa.FlowReply, bool) {
		replies, _, matched := pro.HandleIncomingMessage(uid, accountPhone, phone, message, name)
		var out []wa.FlowReply
		for _, r := range replies { out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData}) }
		return out, matched
	}

	// WA Number Checker
	pro.SetWAChecker(func(phone string) bool {
		exists, err := engine.CheckWANumber(phone)
		if err != nil { return false }
		return exists
	})

	// AI Engine for flow nodes — supports per-flow AI key selection
	pro.SetAICallback(func(systemPrompt, userMessage string, options []string, aiKeyID int64) (string, int) {
		// Guardrails: prompt injection protection
		lower := strings.ToLower(userMessage)
		blocked := []string{"ignore previous", "abaikan instruksi", "system prompt", "api key", "lupakan semua", "forget everything", "who made you", "siapa yang buat", "jailbreak", "dan now"}
		for _, b := range blocked {
			if strings.Contains(lower, b) {
				return "Maaf, saya tidak bisa menjawab pertanyaan itu.", -1
			}
		}
		// Get AI key — per-flow or first available
		keys, _ := db.ListAiKeys(0)
		if len(keys) == 0 {
			return "", -1
		}
		aik := keys[0] // default: first key
		if aiKeyID > 0 {
			for _, k := range keys {
				if k.ID == aiKeyID {
					aik = k
					break
				}
			}
		}
		decKey, _ := secret.Decrypt(aik.APIKey)
		if decKey == "" {
			decKey = aik.APIKey
		}
		reply, err := aiservice.Reply(decKey, aik.Provider, aik.Model, aik.BaseURL, systemPrompt, userMessage, nil, nil)
		if err != nil || reply == "" {
			return "", -1
		}
		if len(options) > 0 {
			replyLower := strings.ToLower(reply)
			for i, opt := range options {
				if strings.Contains(replyLower, strings.ToLower(opt)) {
					return reply, i
				}
			}
			return reply, 0
		}
		return reply, -1
	})

	// Contact Lookup — search by phone
	pro.SetContactLookup(func(phone string) map[string]string {
		result := map[string]string{"phone": phone}
		var name string
		if err := db.QueryRow(`SELECT name FROM contacts WHERE phone=? LIMIT 1`, phone).Scan(&name); err == nil {
			result["name"] = name
		}
		return result
	})

	// DB Query — read-only SQL (SELECT only)
	pro.SetDBQueryFunc(func(query string) (string, error) {
		q := strings.TrimSpace(query)
		if !strings.HasPrefix(strings.ToUpper(q), "SELECT") {
			return "", fmt.Errorf("only SELECT allowed")
		}
		var result string
		if err := db.QueryRow(q).Scan(&result); err != nil {
			return "", err
		}
		return result, nil
	})

	// Store callbacks for flow e-commerce nodes
	pro.SetStoreCallbacks(
		func(category string, maxItems int) []map[string]string {
			products, _ := db.ListProducts()
			var result []map[string]string
			for _, p := range products {
				if category == "" || strings.Contains(strings.ToLower(p.Category), strings.ToLower(category)) {
					result = append(result, map[string]string{"name": p.Name, "price": fmt.Sprintf("%.0f", p.Price)})
					if len(result) >= maxItems {
						break
					}
				}
			}
			return result
		},
		func(phone, product string, qty int) (int64, error) {
			return db.CreateOrder(phone, product, 0, qty, 0)
		},
		func(orderID int64, amount float64) string {
			return fmt.Sprintf("/subscribe?order=%d&amount=%.0f", orderID, amount)
		},
	)

// Telegram flow callback
	wa.TGCallback = func(uid int64, phone, message, name string) ([]wa.FlowReply, bool) {
		replies, _, matched := pro.HandleIncomingMessage(uid, "tg", phone, message, name)
		var out []wa.FlowReply
		for _, r := range replies { out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData}) }
		return out, matched
	}

	// Instagram flow callback
	wa.IGCallback = func(uid int64, phone, message, name string) ([]wa.FlowReply, bool) {
		replies, _, matched := pro.HandleIncomingMessage(uid, "ig", phone, message, name)
		var out []wa.FlowReply
		for _, r := range replies { out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData}) }
		return out, matched
	}

	// FB Messenger flow callback
	wa.FBCallback = func(uid int64, phone, message, name string) ([]wa.FlowReply, bool) {
		replies, _, matched := pro.HandleIncomingMessage(uid, "fb", phone, message, name)
		var out []wa.FlowReply
		for _, r := range replies { out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData}) }
		return out, matched
	}

	// WA Account lister for Flow Builder dropdown
	pro.SetAccountLister(func(uid int64) []pro.WAAccountInfo {
		var list []pro.WAAccountInfo
		for _, a := range engine.Accounts(uid) {
			list = append(list, pro.WAAccountInfo{ID: a.Phone, Phone: a.Phone, Type: "wa", Label: "WA: " + a.Phone + " (" + a.Status + ")"})
		}
		metaAccs, _ := db.ListMetaAccounts()
		for _, m := range metaAccs {
			if m.UserID == uid || m.ParentID == uid || uid == 0 {
				list = append(list, pro.WAAccountInfo{ID: m.PhoneNumberID, Phone: m.PhoneNumberID, Type: "meta", Label: "Meta: " + m.Name + " (" + m.PhoneNumberID + ")"})
			}
		}
		return list
	})

	engine.FlowCallbacks = &wa.FlowCallbacks{
		OnMessage: func(uid int64, accountPhone, phone, message, contactName string) ([]wa.FlowReply, bool) {
			replies, _, matched := pro.HandleIncomingMessage(uid, accountPhone, phone, message, contactName)
			var out []wa.FlowReply
			for _, r := range replies {
				out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData})
			}
			return out, matched
		},
		OnWelcome: func(uid int64, accountPhone, phone, contactName string) []wa.FlowReply {
			replies, _ := pro.HandleWelcomeFlow(uid, accountPhone, phone, contactName)
			var out []wa.FlowReply
			for _, r := range replies {
				out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData})
			}
			return out
		},
		OnFallback: func(uid int64, accountPhone, phone, message, contactName string) []wa.FlowReply {
			replies, _ := pro.HandleFallbackFlow(uid, accountPhone, phone, message, contactName)
			var out []wa.FlowReply
			for _, r := range replies {
				out = append(out, wa.FlowReply{Text: r.Text, MediaURL: r.MediaURL, MediaType: r.MediaType, Action: r.Action, ActionData: r.ActionData})
			}
			return out
		},
		OnCronTick: func() int {
			return pro.RunCronFlows(0, "")
		},
		OnInactivity: func(phones []string) int {
			return pro.RunInactivityFlows(0, "", phones)
		},
	}
	log.Println("[PRO] Flow engine wired to WhatsApp pipeline")

	// SSE real-time: push every incoming message to omnichannel inbox
	engine.OnMessageNotify = func(phone, name, message, channel string, ownerUID int64) {
		pro.PushOmniEvent(phone, name, message, channel, ownerUID)
	}
}
