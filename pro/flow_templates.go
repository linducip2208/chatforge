//go:build pro

package pro

import (
	"encoding/json"
	"net/http"
)

type FlowTemplate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Flow        Flow   `json:"flow"`
}

var templates = []FlowTemplate{
	{
		ID: "restaurant", Name: "Restoran Ordering", Icon: "🍽️", Category: "F&B / Restoran",
		Description: "Menu → pilih kategori → konfirmasi pesanan → payment link",
		Flow: Flow{
			Name: "Restoran Ordering", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "menu,makan,pesan,order"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🍽️ Selamat datang! Mau pesan apa?\n\nKetik:\n1. Makanan\n2. Minuman\n3. Snack"}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[
					{"id":"b1","label":"Makanan","operator":"contains","value":"makan,1"},
					{"id":"b2","label":"Minuman","operator":"contains","value":"minum,2"},
					{"id":"b3","label":"Snack","operator":"contains","value":"snack,3"}]}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{"text":"📋 Menu Makanan:\n1. Nasi Goreng - 25k\n2. Mie Ayam - 20k\n3. Soto Ayam - 22k\n\nKetik nomor untuk pesan."}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 130}, Data: rawJSON(`{"text":"🥤 Menu Minuman:\n1. Es Teh - 5k\n2. Jeruk - 7k\n3. Kopi - 10k\n\nKetik nomor untuk pesan."}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 500, Y: 220}, Data: rawJSON(`{"text":"🍿 Menu Snack:\n1. Kentang Goreng - 15k\n2. Pisang Goreng - 12k\n3. Tahu Crispy - 10k"}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0", Label: "Makanan"},
				{ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1", Label: "Minuman"},
				{ID: "e4", Source: "n2", Target: "n5", SourceH: "out-2", Label: "Snack"},
			},
		},
	},
	{
		ID: "csat", Name: "CSAT Survey", Icon: "⭐",
		Description: "Close chat → minta rating 1-5 → follow-up / complaint",
		Flow: Flow{
			Name: "CSAT Survey", Active: true,
			Trigger: FlowTrigger{Type: TriggerCSAT},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"⭐ Terima kasih sudah chat dengan kami!\n\nBagaimana pengalaman Anda hari ini?\nBeri rating 1-5 (5 = sangat puas)"}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[
					{"id":"b1","label":"Happy","operator":"contains","value":"4,5"},
					{"id":"b2","label":"Neutral","operator":"contains","value":"3"},
					{"id":"b3","label":"Unhappy","operator":"contains","value":"1,2"}]}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{"text":"🌟 Senang mendengarnya! Kalau ada waktu, bantu review kami ya. Terima kasih!"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 130}, Data: rawJSON(`{"text":"Terima kasih feedback-nya! Kami akan terus berusaha lebih baik."}`)},
				{ID: "n5", Type: "transfer_agent", Position: FlowPosition{X: 500, Y: 220}, Data: rawJSON(`{}`)},
				{ID: "n6", Type: "close_chat", Position: FlowPosition{X: 700, Y: 150}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0", Label: "Happy"},
				{ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1", Label: "Neutral"},
				{ID: "e4", Source: "n2", Target: "n5", SourceH: "out-2", Label: "Unhappy"},
				{ID: "e5", Source: "n3", Target: "n6"},
				{ID: "e6", Source: "n4", Target: "n6"},
			},
		},
	},
	{
		ID: "lead", Name: "Lead Qualification", Icon: "🎯",
		Description: "Greet → tanya kebutuhan → qualify → transfer sales / nurturing",
		Flow: Flow{
			Name: "Lead Qualification", Active: true,
			Trigger: FlowTrigger{Type: TriggerWelcome},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"👋 Halo {{name}}! Terima kasih sudah menghubungi kami.\n\nBoleh tahu Anda tertarik dengan layanan apa?"}`)},
				{ID: "n2", Type: "ai_decide", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"system_prompt":"Classify the user into: Hot Lead (ready to buy), Warm Lead (needs info), Cold Lead (just browsing)","options":["Hot Lead","Warm Lead","Cold Lead"],"var_result":"lead_type"}`)},
				{ID: "n3", Type: "transfer_agent", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 130}, Data: rawJSON(`{"text":"📋 Terima kasih infonya! Ini beberapa info yang mungkin membantu:"}`)},
				{ID: "n5", Type: "tag_contact", Position: FlowPosition{X: 500, Y: 220}, Data: rawJSON(`{"tags":["cold-lead"]}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"},
				{ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1"},
				{ID: "e4", Source: "n2", Target: "n5", SourceH: "out-2"},
			},
		},
	},
	{
		ID: "appointment", Name: "Appointment Booking", Icon: "📅",
		Description: "Tanya tanggal → konfirmasi → reminder H-1",
		Flow: Flow{
			Name: "Appointment Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "booking,jadwal,janji,appointment"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"📅 Mau booking untuk kapan?\n\nFormat: DD/MM/YYYY\nContoh: 25/12/2026"}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Siapa nama dan nomor HP yang bisa dihubungi?","var_name":"contact_info","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Booking diterima!\n\nKami akan konfirmasi via WhatsApp H-1 sebelum jadwal.\n\nAda yang bisa dibantu lagi?"}`)},
				{ID: "n4", Type: "wait", Position: FlowPosition{X: 500, Y: 190}, Data: rawJSON(`{"seconds":86400}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 700, Y: 190}, Data: rawJSON(`{"text":"🔔 Reminder: Besok ada janji temu. Masih bisa kan?\n\nBalas YA untuk konfirmasi."}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
				{ID: "e4", Source: "n4", Target: "n5"},
			},
		},
	},
	{
		ID: "support", Name: "Customer Support Triage", Icon: "🛟",
		Description: "Klasifikasi → FAQ / transfer agent / escalate",
		Flow: Flow{
			Name: "Support Triage", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "bantu,help,masalah,error,komplain"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🛟 Ada yang bisa kami bantu?\n\nCeritakan masalah Anda secara singkat."}`)},
				{ID: "n2", Type: "ai_decide", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"system_prompt":"Classify support issue: Technical, Billing, General Question, Complaint, Urgent","options":["Technical","Billing","General","Complaint","Urgent"],"var_result":"issue_type"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 20}, Data: rawJSON(`{"text":"🔧 Tim teknis kami akan membantu. Mohon tunggu sebentar."}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 90}, Data: rawJSON(`{"text":"💳 Untuk pertanyaan billing, silakan cek invoice Anda."}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 500, Y: 160}, Data: rawJSON(`{"text":"📋 Ini FAQ kami, semoga membantu!"}`)},
				{ID: "n6", Type: "transfer_agent", Position: FlowPosition{X: 500, Y: 230}, Data: rawJSON(`{}`)},
				{ID: "n7", Type: "transfer_agent", Position: FlowPosition{X: 500, Y: 300}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"},
				{ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1"},
				{ID: "e4", Source: "n2", Target: "n5", SourceH: "out-2"},
				{ID: "e5", Source: "n2", Target: "n6", SourceH: "out-3"},
				{ID: "e6", Source: "n2", Target: "n7", SourceH: "out-4"},
			},
		},
	},
	{
		ID: "abandoned", Name: "Abandoned Cart Recovery", Icon: "🛒",
		Description: "Detect → reminder → discount → checkout",
		Flow: Flow{
			Name: "Cart Recovery", Active: true,
			Trigger: FlowTrigger{Type: TriggerDrip},
			Nodes: []FlowNode{
				{ID: "n1", Type: "wait", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"seconds":3600}`)},
				{ID: "n2", Type: "message", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"text":"🛒 Hai {{name}}! Masih ada barang di keranjang nih.\n\nJangan sampai kehabisan ya! Ada yang bisa dibantu?"}`)},
				{ID: "n3", Type: "wait", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"seconds":86400}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"🎁 Spesial buat kamu: diskon 10% dengan kode: COMEBACK10\n\nBerlaku 24 jam. Checkout sekarang!"}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
			},
		},
	},
	{
		ID: "event", Name: "Event Registration", Icon: "🎉",
		Description: "Greet → tanya nama/email → konfirmasi → reminder",
		Flow: Flow{
			Name: "Event Registration", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "daftar,register,event,webinar"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🎉 Selamat datang di pendaftaran event!\n\nSiapa nama Anda?"}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Alamat email Anda?","var_name":"email","type":"email"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Terdaftar!\n\nDetail event akan dikirim ke email {{email}}.\n\nSampai jumpa di event!"}`)},
				{ID: "n4", Type: "tag_contact", Position: FlowPosition{X: 300, Y: 200}, Data: rawJSON(`{"tags":["event-participant"]}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n2", Target: "n4"},
			},
		},
	},
	{
		ID: "product_reco", Name: "Product Recommendation", Icon: "💡",
		Description: "Tanya preferensi → rekomendasi AI → konfirmasi order",
		Flow: Flow{
			Name: "Product Recommendation", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "rekomendasi,rekomen,saran,cari"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💡 Mau cari apa? Ceritakan kebutuhan Anda.\n\nContoh: \"cari HP budget 3 juta untuk gaming\""}`)},
				{ID: "n2", Type: "ai_reply", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"system_prompt":"You are a helpful shopping assistant. Recommend 3 products based on user needs. Include price range. Reply in friendly Indonesian.","var_result":"recommendation"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"Semoga rekomendasi di atas membantu! Ada yang mau ditanyakan lagi?"}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
			},
		},
	},
	{
		ID: "order_status", Name: "Order Status Tracker", Icon: "📦",
		Description: "Cek status order → lookup DB → kirim status",
		Flow: Flow{
			Name: "Order Status", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "cek order,status order,order"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "question", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"question":"📦 Masukkan nomor order Anda.","var_name":"order_id","type":"text"}`)},
				{ID: "n2", Type: "db_query", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"query":"SELECT status FROM store_orders WHERE id={{order_id}}","var_result":"order_status"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"📋 Order #{{order_id}}: *{{order_status}}*\n\nAda yang bisa dibantu lagi?"}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
			},
		},
	},
	{
		ID: "payment_reminder", Name: "Payment Reminder", Icon: "💰",
		Description: "H-3 reminder → payment link → konfirmasi",
		Flow: Flow{
			Name: "Payment Reminder", Active: true,
			Trigger: FlowTrigger{Type: TriggerCron},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💰 Halo {{name}}! Tagihan sebesar Rp{{amount}} jatuh tempo {{due_date}}.\n\nJangan sampai terlambat ya!"}`)},
				{ID: "n2", Type: "wait", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"seconds":259200}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"⏰ Reminder: tagihan Rp{{amount}} jatuh tempo HARI INI.\n\nBayar sekarang: /subscribe"}`)},
				{ID: "n4", Type: "close_chat", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
			},
		},
	},
	{
		ID: "delivery", Name: "Delivery Confirmation", Icon: "🚚",
		Description: "Konfirmasi pengiriman → minta foto → rating",
		Flow: Flow{
			Name: "Delivery Confirmation", Active: true,
			Trigger: FlowTrigger{Type: TriggerDrip},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🚚 Pesanan Anda sedang dalam perjalanan!\n\nEstimasi sampai: {{eta}}"}`)},
				{ID: "n2", Type: "wait", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"seconds":86400}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"📸 Pesanan sudah sampai? Share foto barang yang diterima ya!"}`)},
				{ID: "n4", Type: "poll", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"question":"Bagaimana kondisi barang?","options":["Bagus!","Ada masalah","Belum sampai"],"var_result":"delivery_rating"}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
			},
		},
	},
	{
		ID: "loyalty", Name: "Loyalty Program", Icon: "🏆",
		Description: "Cek poin → redeem → voucher",
		Flow: Flow{
			Name: "Loyalty Program", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "poin,loyalty,reward,redeem"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🏆 Program Loyalty!\n\nCek poin Anda: ketik *POIN*\nTukar reward: ketik *TUKAR*"}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[{"id":"b1","label":"Cek Poin","operator":"contains","value":"poin,cek"},{"id":"b2","label":"Tukar","operator":"contains","value":"tukar,redeem"}]}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{"text":"⭐ Poin Anda: {{contact.total_points}}\n\nTerus belanja untuk kumpulkan poin!"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 130}, Data: rawJSON(`{"text":"🎁 Reward tersedia:\n1. Diskon 10% (100 poin)\n2. Gratis ongkir (200 poin)\n3. Cashback 50rb (500 poin)\n\nBalas nomor untuk tukar."}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"},
				{ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1"},
			},
		},
	},
	{
		ID: "price_quote", Name: "Price Quote Generator", Icon: "📊",
		Description: "Tanya spesifikasi → kalkulasi → kirim penawaran",
		Flow: Flow{
			Name: "Price Quote", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "harga,penawaran,quote,estimasi"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "question", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"question":"📊 Untuk bantu hitung penawaran, ceritakan kebutuhan Anda.","var_name":"requirement","type":"text"}`)},
				{ID: "n2", Type: "ai_reply", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"system_prompt":"You are a pricing specialist. Based on user requirements, provide a detailed price quote with:\n1. Package recommendation\n2. Price breakdown\n3. Estimated timeline\nReply in Indonesian.","var_result":"quote"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"📋 Butuh detail lebih lanjut? Ketik *LANJUT* untuk bicara dengan tim sales."}`)},
				{ID: "n4", Type: "condition", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[{"id":"b1","label":"Lanjut","operator":"contains","value":"lanjut,sales,manusia"}]}`)},
				{ID: "n5", Type: "transfer_agent", Position: FlowPosition{X: 900, Y: 100}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
				{ID: "e4", Source: "n4", Target: "n5", SourceH: "out-0"},
			},
		},
	},
	{
		ID: "complaint", Name: "Complaint Handler", Icon: "🚨",
		Description: "Terima komplain → escalate → follow-up",
		Flow: Flow{
			Name: "Complaint Handler", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "komplain,kecewa,marah,rusak,salah"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🚨 Kami minta maaf atas ketidaknyamanan ini.\n\nCeritakan masalahnya agar kami bisa bantu."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nomor order atau invoice yang bermasalah?","var_name":"order_ref","type":"text"}`)},
				{ID: "n3", Type: "set_variable", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"variables":{"priority":"HIGH"}}`)},
				{ID: "n4", Type: "tag_contact", Position: FlowPosition{X: 300, Y: 200}, Data: rawJSON(`{"tags":["complaint"]}`)},
				{ID: "n5", Type: "transfer_agent", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
				{ID: "e4", Source: "n2", Target: "n4"},
				{ID: "e5", Source: "n3", Target: "n5"},
			},
		},
	},
	{
		ID: "onboarding", Name: "Onboarding 7 Hari", Icon: "🚀",
		Description: "Welcome → day1 tips → day3 promo → day7 survey",
		Flow: Flow{
			Name: "Onboarding 7 Days", Active: true,
			Trigger: FlowTrigger{Type: TriggerWelcome},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 40}, Data: rawJSON(`{"text":"🚀 Selamat datang {{name}}!\n\nKami akan bantu maksimalkan pengalaman Anda. Cek tips harian dari kami."}`)},
				{ID: "n2", Type: "wait", Position: FlowPosition{X: 300, Y: 40}, Data: rawJSON(`{"seconds":86400}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{"text":"📌 Hari 1: Kenali fitur broadcast.\n\nCoba kirim pesan massal ke kontak Anda. Ada yang bisa dibantu?"}`)},
				{ID: "n4", Type: "wait", Position: FlowPosition{X: 300, Y: 130}, Data: rawJSON(`{"seconds":259200}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 500, Y: 130}, Data: rawJSON(`{"text":"🎁 Hari 3: Diskon 20% untuk upgrade Pro!\n\nKode: WELCOME20\nBerlaku 24 jam."}`)},
				{ID: "n6", Type: "wait", Position: FlowPosition{X: 300, Y: 220}, Data: rawJSON(`{"seconds":604800}`)},
				{ID: "n7", Type: "message", Position: FlowPosition{X: 500, Y: 220}, Data: rawJSON(`{"text":"⭐ Hari 7: Bagaimana pengalaman Anda?\n\nBeri rating 1-5 agar kami bisa lebih baik."}`)},
				{ID: "n8", Type: "close_chat", Position: FlowPosition{X: 700, Y: 220}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
				{ID: "e3", Source: "n3", Target: "n4"},
				{ID: "e4", Source: "n4", Target: "n5"},
				{ID: "e5", Source: "n5", Target: "n6"},
				{ID: "e6", Source: "n6", Target: "n7"},
				{ID: "e7", Source: "n7", Target: "n8"},
			},
		},
	},
	{
		ID: "referral", Name: "Referral Program", Icon: "👥",
		Description: "Ajak teman → kode referral → diskon",
		Flow: Flow{
			Name: "Referral Program", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "referral,ajak teman,referal,kode"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"👥 Program Referral!\n\nAjak teman & dapat diskon 20%.\n\nKode referral Anda: *{{contact.referral_code}}*\n\nBagikan kode ini ke teman."}`)},
				{ID: "n2", Type: "set_variable", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variables":{"referral_code":"REF{{phone}}"}}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"📋 Cara dapat diskon:\n1. Bagikan kode ke teman\n2. Teman daftar pakai kode Anda\n3. Anda berdua dapat diskon 20%!\n\nKetik *STATUS* untuk cek referral."}`)},
			},
			Edges: []FlowEdge{
				{ID: "e1", Source: "n1", Target: "n2"},
				{ID: "e2", Source: "n2", Target: "n3"},
			},
		},
	},
	{
		ID: "health_booking", Name: "Klinik Booking", Icon: "🏥", Category: "Healthcare / Klinik",
		Description: "Booking dokter → pilih jadwal → konfirmasi → reminder H-1",
		Flow: Flow{
			Name: "Klinik Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "booking,dokter,jadwal,klinik"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🏥 Selamat datang di Klinik Sehat!\n\nLayanan kami:\n1. Dokter Umum\n2. Dokter Gigi\n3. Dokter Anak\n\nKetik nomor untuk pilih."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Mau booking untuk tanggal berapa?\nFormat: DD/MM/YYYY","var_name":"tanggal","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Nama pasien dan nomor HP?","var_name":"pasien","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Booking berhasil!\n\n📅 Tanggal: {{tanggal}}\n👤 Pasien: {{pasien}}\n\nKami akan kirim reminder H-1."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "health_reminder", Name: "Obat Reminder", Icon: "💊", Category: "Healthcare / Klinik",
		Description: "Reminder minum obat → jadwal harian → konfirmasi",
		Flow: Flow{
			Name: "Obat Reminder", Active: true,
			Trigger: FlowTrigger{Type: TriggerCron},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💊 Reminder: Saatnya minum obat!\n\nJangan lupa minum sesuai resep dokter ya."}`)},
				{ID: "n2", Type: "wait", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"seconds":21600}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"🕐 Ingat lagi ya, minum obat tepat waktu itu penting untuk kesembuhan."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "edu_enrollment", Name: "School Enrollment", Icon: "🎓", Category: "Education / LMS",
		Description: "Info program → pendaftaran → upload dokumen → konfirmasi",
		Flow: Flow{
			Name: "School Enrollment", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "daftar,sekolah,pendaftaran,enroll"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🎓 Pendaftaran Siswa Baru!\n\nPilih program:\n1. TK/PAUD\n2. SD\n3. SMP\n4. SMA\n\nKetik nomor untuk lanjut."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nama lengkap calon siswa?","var_name":"nama","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Nama orang tua & nomor HP?","var_name":"ortu","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Pendaftaran diterima!\n\n👤 Siswa: {{nama}}\n👨‍👩‍👧 Ortu: {{ortu}}\n\nTim kami akan menghubungi untuk proses selanjutnya."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "edu_course", Name: "Online Course Selling", Icon: "📚", Category: "Education / LMS",
		Description: "Katalog kursus → detail → daftar → payment → akses",
		Flow: Flow{
			Name: "Online Course", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "kursus,belajar,kelas,pelatihan"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"📚 Kursus Online:\n\n1. Digital Marketing - Rp299k\n2. Graphic Design - Rp349k\n3. Web Development - Rp499k\n4. Data Science - Rp599k\n\nKetik nomor untuk info detail."}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[{"id":"b1","label":"DM","operator":"contains","value":"1"},{"id":"b2","label":"GD","operator":"contains","value":"2"},{"id":"b3","label":"WD","operator":"contains","value":"3"},{"id":"b4","label":"DS","operator":"contains","value":"4"}]}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 30}, Data: rawJSON(`{"text":"📊 Digital Marketing:\n• SEO, SEM, Social Media\n• 8 minggu · Sertifikat\n• Mentor praktisi\n\nKetik DAFTAR untuk enroll."}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 110}, Data: rawJSON(`{"text":"🎨 Graphic Design:\n• Canva, Photoshop, Illustrator\n• 8 minggu · Portfolio\n• Mentor profesional"}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 500, Y: 190}, Data: rawJSON(`{"text":"💻 Web Development:\n• HTML, CSS, PHP, Laravel\n• 12 minggu · Project\n• Full stack"}`)},
				{ID: "n6", Type: "message", Position: FlowPosition{X: 500, Y: 270}, Data: rawJSON(`{"text":"📈 Data Science:\n• Python, SQL, ML\n• 12 minggu · Case study\n• Industri ready"}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"}, {ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1"}, {ID: "e4", Source: "n2", Target: "n5", SourceH: "out-2"}, {ID: "e5", Source: "n2", Target: "n6", SourceH: "out-3"}},
		},
	},
	{
		ID: "property_listing", Name: "Property Inquiry", Icon: "🏠", Category: "Property / Real Estate",
		Description: "Cari properti → filter type/harga → info detail → jadwal survey",
		Flow: Flow{
			Name: "Property Inquiry", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "rumah,properti,cari rumah,apartemen,ruko"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🏠 Cari Properti Idaman?\n\nPilih tipe:\n1. Rumah\n2. Apartemen\n3. Ruko\n4. Tanah\n\nKetik nomor atau ceritakan kebutuhan Anda."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"💰 Budget berapa? (contoh: 500jt - 1M)","var_name":"budget","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"📍 Lokasi yang diinginkan? (contoh: Jakarta Selatan)","var_name":"lokasi","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"📋 Terima kasih! Tim kami akan kirim rekomendasi properti sesuai kriteria:\n\n💰 Budget: {{budget}}\n📍 Lokasi: {{lokasi}}\n\nAda jadwal survey juga? Ketik SURVEY."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "travel_booking", Name: "Travel Booking", Icon: "✈️", Category: "Travel / Hospitality",
		Description: "Cari destinasi → pilih paket → booking → payment → e-ticket",
		Flow: Flow{
			Name: "Travel Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "travel,tiket,liburan,wisata,tour"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"✈️ Mau liburan kemana?\n\nDestinasi populer:\n1. Bali\n2. Lombok\n3. Yogyakarta\n4. Bandung\n5. Labuan Bajo\n\nKetik nama kota atau nomor."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Tanggal berangkat? (DD/MM/YYYY)","var_name":"tanggal","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"👥 Berapa orang?","var_name":"pax","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"🎉 Paket tersedia!\n\n📅 {{tanggal}} · {{pax}} orang\n\nKetik LANJUT untuk lihat detail harga."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "hotel_booking", Name: "Hotel Reservation", Icon: "🏨", Category: "Travel / Hospitality",
		Description: "Cari kamar → check-in/out → pilih tipe → booking → konfirmasi",
		Flow: Flow{
			Name: "Hotel Reservation", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "hotel,kamar,menginap,reservasi"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🏨 Booking Hotel Mudah!\n\nPilih tipe kamar:\n1. Standard Single\n2. Standard Double\n3. Deluxe\n4. Suite\n\nKetik nomor untuk lanjut."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Check-in date? (DD/MM/YYYY)","var_name":"checkin","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"📅 Check-out date? (DD/MM/YYYY)","var_name":"checkout","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Booking diterima!\n\n📅 {{checkin}} → {{checkout}}\n\nKami akan kirim konfirmasi via WA."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "finance_loan", Name: "Loan / Kredit Inquiry", Icon: "💳", Category: "Finance / Banking",
		Description: "Cek eligibility → pilih produk → simulasi → apply",
		Flow: Flow{
			Name: "Loan Inquiry", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "pinjaman,kredit,dana,cicilan"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💳 Butuh dana cepat?\n\nProduk kami:\n1. KTA (tanpa agunan) - bunga 0.8%\n2. KKB (kendaraan) - bunga 0.6%\n3. KPR (rumah) - bunga 0.5%\n\nKetik nomor untuk simulasi."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"💰 Jumlah pinjaman yang dibutuhkan? (contoh: 50jt)","var_name":"jumlah","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"📅 Tenor yang diinginkan? (contoh: 12 bulan)","var_name":"tenor","type":"text"}`)},
				{ID: "n4", Type: "question", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"question":"Nama lengkap & nomor KTP?","var_name":"identitas","type":"text"}`)},
				{ID: "n5", Type: "transfer_agent", Position: FlowPosition{X: 900, Y: 100}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}, {ID: "e4", Source: "n4", Target: "n5"}},
		},
	},
	{
		ID: "auto_service", Name: "Bengkel Booking", Icon: "🔧", Category: "Automotive",
		Description: "Pilih service → booking jadwal → estimasi → reminder",
		Flow: Flow{
			Name: "Bengkel Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "bengkel,servis,service,montir"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🔧 Bengkel AutoPro\n\nService tersedia:\n1. Ganti Oli\n2. Tune Up\n3. AC Service\n4. Body Repair\n\nKetik nomor atau ceritakan masalah kendaraan Anda."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"🚗 Nama & plat nomor kendaraan?","var_name":"kendaraan","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"📅 Booking hari apa? (DD/MM/YYYY)","var_name":"jadwal","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Booking bengkel:\n🚗 {{kendaraan}}\n📅 {{jadwal}}\n\nEstimasi selesai akan diinfokan saat datang."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "logistics_track", Name: "Package Tracking", Icon: "📦", Category: "Logistics / Delivery",
		Description: "Input resi → cek status → estimasi tiba → komplain jika perlu",
		Flow: Flow{
			Name: "Package Tracking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "resi,lacak,tracking,paket"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "question", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"question":"📦 Masukkan nomor resi Anda.","var_name":"resi","type":"text"}`)},
				{ID: "n2", Type: "message", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"text":"📋 Resi: {{resi}}\n\nStatus: Dalam Pengiriman 🚚\nEstimasi tiba: 2-3 hari\n\nAda masalah? Ketik KOMPLAIN."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}},
		},
	},
	{
		ID: "membership_signup", Name: "Membership Signup", Icon: "🪪", Category: "SaaS / Subscription",
		Description: "Pilih plan → registrasi → payment → welcome email",
		Flow: Flow{
			Name: "Membership Signup", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "daftar,member,langganan,subs"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🪪 Pilih paket membership:\n\n1. Basic - Rp99k/bln (fitur dasar)\n2. Pro - Rp249k/bln (semua fitur)\n3. Enterprise - Rp999k/bln (custom)\n\nKetik nomor untuk daftar."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nama lengkap & email?","var_name":"identitas","type":"text"}`)},
				{ID: "n3", Type: "payment_link", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"amount":"249000"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"🎉 Selamat bergabung!\n\nCek inbox email untuk akses dashboard.\nAda yang bisa dibantu?"}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "job_apply", Name: "Job Application", Icon: "💼", Category: "HR / Recruitment",
		Description: "List lowongan → apply → upload CV → screening → interview",
		Flow: Flow{
			Name: "Job Application", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "lowongan,loker,karir,kerja,lamaran"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💼 Lowongan Tersedia:\n\n1. Frontend Developer\n2. Backend Developer\n3. UI/UX Designer\n4. Digital Marketing\n5. Customer Service\n\nKetik nomor untuk detail & apply."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nama lengkap & email Anda?","var_name":"identitas","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Lamaran diterima!\n\n👤 {{identitas}}\n\nHR akan review dan hubungi untuk interview jika cocok.\n\nTips: Siapkan portofolio terbaik Anda."}`)},
				{ID: "n4", Type: "tag_contact", Position: FlowPosition{X: 300, Y: 200}, Data: rawJSON(`{"tags":["job-applicant"]}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n2", Target: "n4"}},
		},
	},
	{
		ID: "social_contest", Name: "Giveaway / Contest", Icon: "🎁", Category: "Marketing / Growth",
		Description: "Join contest → share → invite friends → leaderboard → winner",
		Flow: Flow{
			Name: "Giveaway Contest", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "giveaway,kontes,hadiah,lomba"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🎁 GIVEAWAY!\n\nMenangkan iPhone 15!\n\nCara ikut:\n1. Share info ini ke 5 teman\n2. Minta mereka ketik REFERAL: {{phone}}\n3. Terbanyak referral menang!\n\nKetik JOIN untuk ikut."}`)},
				{ID: "n2", Type: "set_variable", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variables":{"contestant":"{{phone}}"}}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"🏆 Anda terdaftar!\n\nKode referral: REF{{phone}}\n\nShare ke teman, minta mereka ketik kode ini. Leaderboard diumumkan setiap minggu."}`)},
				{ID: "n4", Type: "tag_contact", Position: FlowPosition{X: 300, Y: 200}, Data: rawJSON(`{"tags":["contestant"]}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n2", Target: "n4"}},
		},
	},
	{
		ID: "fitness_plan", Name: "Fitness / Gym Plan", Icon: "💪", Category: "Health / Fitness",
		Description: "Cek program → pilih trainer → booking sesi → progress tracking",
		Flow: Flow{
			Name: "Fitness Plan", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "gym,fitness,olahraga,workout,diet"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💪 Mau transformasi badan?\n\nProgram kami:\n1. Weight Loss\n2. Muscle Building\n3. Yoga & Flexibility\n4. Personal Trainer\n\nKetik nomor untuk info."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Apa goal fitness Anda? Ceritakan singkat.","var_name":"goal","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Nama & nomor HP untuk free trial?","var_name":"kontak","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"🎯 Terima kasih!\n\nGoal: {{goal}}\nKontak: {{kontak}}\n\nKami akan hubungi untuk jadwal free trial session.\n\nSiap berubah? 💪"}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "insurance_claim", Name: "Insurance Claim", Icon: "🛡️", Category: "Finance / Banking",
		Description: "Lapor klaim → upload dokumen → verifikasi → pencairan",
		Flow: Flow{
			Name: "Insurance Claim", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "klaim,asuransi,polis,claim"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🛡️ Klaim Asuransi\n\nSilakan siapkan:\n• Nomor polis\n• Kronologi kejadian\n• Dokumen pendukung\n\nKetik LANJUT untuk mulai."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nomor polis Anda?","var_name":"polis","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Ceritakan kronologi kejadian secara singkat.","var_name":"kronologi","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"📋 Klaim terdaftar!\n\nPolis: {{polis}}\nStatus: Dalam proses verifikasi\n\nTim kami akan update dalam 1x24 jam."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "laundry_order", Name: "Laundry Order", Icon: "👕", Category: "Service / Jasa",
		Description: "Pilih layanan → estimasi harga → pickup → selesai → notifikasi",
		Flow: Flow{
			Name: "Laundry Order", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "laundry,cuci,cucian,bersih"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"👕 Laundry Cepat!\n\nLayanan:\n1. Cuci Setrika - Rp8k/kg\n2. Cuci Kering - Rp5k/kg\n3. Setrika - Rp5k/kg\n4. Sepatu - Rp20k/pasang\n\nKetik nomor + alamat pickup."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📍 Alamat lengkap untuk pickup?","var_name":"alamat","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Order laundry diterima!\n\n📍 Pickup: {{alamat}}\n⏱ Estimasi selesai: 2-3 hari\n\nKami akan kabari sebelum pickup."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "food_delivery", Name: "Food Delivery Order", Icon: "🛵", Category: "F&B / Restoran",
		Description: "Pilih resto → menu → order → payment → tracking delivery",
		Flow: Flow{
			Name: "Food Delivery", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "gofood,grabfood,pesan antar,delivery"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🛵 Pesan Makanan!\n\nMenu hari ini:\n1. Nasi Goreng Spesial - 25k\n2. Ayam Geprek - 22k\n3. Mie Ayam Komplit - 20k\n4. Sate Ayam - 30k\n\nKetik nomor + alamat antar."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📍 Alamat pengiriman lengkap?","var_name":"alamat","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Pesanan diterima!\n\n📍 Antar ke: {{alamat}}\n🚀 Estimasi: 30-45 menit\n\nPembayaran COD / Transfer."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "car_rental", Name: "Car Rental Booking", Icon: "🚗", Category: "Automotive",
		Description: "Pilih mobil → cek ketersediaan → booking → pembayaran → serah terima",
		Flow: Flow{
			Name: "Car Rental", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "rental,sewa mobil,car rental"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🚗 Rental Mobil\n\nPilihan:\n1. Avanza - 300k/hari\n2. Innova - 500k/hari\n3. Alphard - 1.5jt/hari\n4. Bus Medium - 2jt/hari\n\nKetik nomor untuk booking."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Tanggal sewa? (DD/MM - DD/MM)","var_name":"tanggal","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Nama & nomor HP penyewa?","var_name":"penyewa","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Booking rental:\n🚗 {{tanggal}}\n👤 {{penyewa}}\n\nSopir tersedia +50k/hari.\nKetik SOPIR jika butuh."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "salon_booking", Name: "Salon & Spa Booking", Icon: "💇", Category: "Service / Jasa",
		Description: "Pilih treatment → booking jadwal → konfirmasi → reminder",
		Flow: Flow{
			Name: "Salon Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "salon,spa,potong rambut,creambath"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💇 Salon & Spa\n\nTreatment:\n1. Haircut - 75k\n2. Creambath - 100k\n3. Facial - 150k\n4. Massage 60min - 200k\n\nKetik nomor untuk booking."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Hari & jam yang diinginkan?","var_name":"jadwal","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"✅ Booking terkonfirmasi!\n\n🕐 {{jadwal}}\n\nDatang 10 menit sebelum jadwal ya. Ada parkir luas."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "church_prayer", Name: "Doa / Prayer Request", Icon: "🙏", Category: "Religious / NGO",
		Description: "Submit doa → kategori → didoakan komunitas → follow-up",
		Flow: Flow{
			Name: "Prayer Request", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "doa,prayer,dukungan,rohani"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🙏 Tim Doa siap mendukung Anda.\n\nSilakan share pokok doa Anda. Kami akan doakan bersama komunitas."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Apa yang bisa kami doakan?","var_name":"doa","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"🙏 Terima kasih sudah berbagi. Doa Anda sudah dicatat dan akan didoakan.\n\nTuhan memberkati!"}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "ticket_concert", Name: "Ticket / Concert Booking", Icon: "🎫", Category: "Event / Entertainment",
		Description: "List event → pilih seat → payment → e-ticket → reminder",
		Flow: Flow{
			Name: "Ticket Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "tiket,konser,event,show,festival"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🎫 Tiket Event!\n\nUpcoming:\n1. Jazz Festival - 25 Des\n2. Rock Concert - 1 Jan\n3. Standup Comedy - 10 Jan\n\nKetik nomor untuk beli tiket."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"🎟 Jumlah tiket?","var_name":"qty","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Nama pemesan & email (untuk e-ticket)?","var_name":"pemesan","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"🎉 Tiket dikonfirmasi!\n\n🎟 {{qty}} tiket\n👤 {{pemesan}}\n\nE-ticket akan dikirim via email."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "voting_polling", Name: "Voting / Quick Poll", Icon: "🗳️", Category: "Marketing / Growth",
		Description: "Buat polling → broadcast → kumpulkan suara → hasil real-time",
		Flow: Flow{
			Name: "Quick Poll", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "vote,polling,pilih,pilihan"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "poll", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"question":"Apa yang paling penting buat kamu?","options":["Harga","Kualitas","Pelayanan","Kecepatan"],"var_result":"vote","max_select":1}`)},
				{ID: "n2", Type: "message", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"text":"✅ Suara tercatat: {{vote}}\n\nTerima kasih partisipasinya! Hasil polling akan diumumkan."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}},
		},
	},
	{
		ID: "donation_charity", Name: "Donation / Charity", Icon: "❤️", Category: "Religious / NGO",
		Description: "Pilih campaign → nominal → payment → receipt → terima kasih",
		Flow: Flow{
			Name: "Donation", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "donasi,sumbangan,sedekah,amal,bantu"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"❤️ Donasi Peduli\n\nCampaign aktif:\n1. Bantu Pendidikan\n2. Bencana Alam\n3. Panti Asuhan\n4. Masjid/Pesantren\n\nKetik nomor untuk donasi."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"💰 Nominal donasi? (min 10rb)","var_name":"nominal","type":"text"}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"text":"🤲 Terima kasih!\n\nDonasi Rp{{nominal}} akan disalurkan.\n\nTransfer ke:\nBCA 1234567890\n\nKonfirmasi: ketik KONFIRMASI."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}},
		},
	},
	{
		ID: "real_estate_survey", Name: "Survey / Open House", Icon: "🔑", Category: "Property / Real Estate",
		Description: "Daftar survey → pilih properti → jadwal → reminder → follow-up",
		Flow: Flow{
			Name: "Open House Survey", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "survey,open house,lihat rumah"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🔑 Open House Weekend!\n\nSabtu-Minggu, jam 10:00-16:00\n\nLokasi:\n1. Cluster Anggrek - Tipe 36/72\n2. Cluster Mawar - Tipe 45/90\n3. Cluster Melati - Tipe 60/120\n\nKetik nomor untuk daftar survey."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Nama & nomor HP untuk daftar survey?","var_name":"pendaftar","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Mau survey hari Sabtu atau Minggu?","var_name":"hari","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Terdaftar!\n\n👤 {{pendaftar}}\n📅 {{hari}}\n\nAlamat lengkap akan dikirim via WA pagi hari."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "consultation_booking", Name: "Konsultan Booking", Icon: "📋", Category: "Service / Jasa",
		Description: "Pilih konsultan → jadwal → brief → payment → session",
		Flow: Flow{
			Name: "Consulting Booking", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "konsultasi,konsultan,advisor,minta saran"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"📋 Layanan Konsultasi\n\n1. Bisnis & Startup - 300k/jam\n2. Keuangan - 250k/jam\n3. Hukum - 500k/jam\n4. Karir - 200k/jam\n\nKetik nomor atau ceritakan kebutuhan."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"Ceritakan secara singkat apa yang ingin dikonsultasikan.","var_name":"brief","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"Kapan jadwal yang diinginkan? (contoh: Senin jam 10:00)","var_name":"jadwal","type":"text"}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"text":"✅ Booking konsultasi:\n\n📋 {{brief}}\n🕐 {{jadwal}}\n\nLink Zoom akan dikirim 1 jam sebelum sesi."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}},
		},
	},
	{
		ID: "wedding_planner", Name: "Wedding / Event Organizer", Icon: "💍", Category: "Event / Entertainment",
		Description: "Tanya paket → pilih venue → vendor → quotation → booking",
		Flow: Flow{
			Name: "Wedding Planner", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "wedding,nikah,pernikahan,resepsi,event organizer,eo"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"💍 Wedding Planner\n\nPaket pernikahan:\n1. Silver - 50jt (200 tamu)\n2. Gold - 100jt (500 tamu)\n3. Platinum - 200jt+ (custom)\n\nTermasuk: venue, dekor, catering, foto, MC.\n\nKetik nomor untuk konsultasi gratis."}`)},
				{ID: "n2", Type: "question", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"question":"📅 Kapan rencana tanggal pernikahan?","var_name":"tanggal","type":"text"}`)},
				{ID: "n3", Type: "question", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{"question":"👥 Estimasi jumlah tamu?","var_name":"tamu","type":"text"}`)},
				{ID: "n4", Type: "question", Position: FlowPosition{X: 700, Y: 100}, Data: rawJSON(`{"question":"Nama pasangan & nomor HP?","var_name":"pasangan","type":"text"}`)},
				{ID: "n5", Type: "message", Position: FlowPosition{X: 900, Y: 100}, Data: rawJSON(`{"text":"💐 Terima kasih!\n\n📅 {{tanggal}} · {{tamu}} tamu\n💑 {{pasangan}}\n\nWedding consultant kami akan menghubungi untuk konsultasi GRATIS."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3"}, {ID: "e3", Source: "n3", Target: "n4"}, {ID: "e4", Source: "n4", Target: "n5"}},
		},
	},
	{
		ID: "customer_verification", Name: "OTP / Verifikasi", Icon: "🔐", Category: "Security / Verification",
		Description: "Kirim OTP → verifikasi → akses diberikan / ditolak",
		Flow: Flow{
			Name: "OTP Verification", Active: true,
			Trigger: FlowTrigger{Type: TriggerKeyword, Value: "otp,verifikasi,kode,verif"},
			Nodes: []FlowNode{
				{ID: "n1", Type: "message", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"text":"🔐 Verifikasi Akun\n\nKode OTP telah dikirim ke nomor Anda.\n\nMasukkan kode 6 digit untuk verifikasi."}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{input}}","branches":[{"id":"b1","label":"Benar","operator":"equals","value":"123456"},{"id":"b2","label":"Salah","operator":"contains","value":""}]}`)},
				{ID: "n3", Type: "message", Position: FlowPosition{X: 500, Y: 40}, Data: rawJSON(`{"text":"✅ Verifikasi berhasil! Akun Anda sudah aktif."}`)},
				{ID: "n4", Type: "message", Position: FlowPosition{X: 500, Y: 120}, Data: rawJSON(`{"text":"❌ Kode salah. Silakan coba lagi atau ketik ULANG untuk kirim OTP baru."}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"}, {ID: "e3", Source: "n2", Target: "n4", SourceH: "out-1"}},
		},
	},
	{
		ID: "faq_bot", Name: "FAQ Knowledge Base", Icon: "❓", Category: "Customer Service",
		Description: "Kategori FAQ → jawaban → eskalasi jika tidak ketemu",
		Flow: Flow{
			Name: "FAQ Bot", Active: true,
			Trigger: FlowTrigger{Type: TriggerAlways},
			Nodes: []FlowNode{
				{ID: "n1", Type: "ai_reply", Position: FlowPosition{X: 100, Y: 100}, Data: rawJSON(`{"system_prompt":"You are a helpful FAQ bot. Answer user questions based on knowledge base. If unsure, say 'Saya akan sambungkan ke tim support.' Reply in Indonesian.","var_result":"answer"}`)},
				{ID: "n2", Type: "condition", Position: FlowPosition{X: 300, Y: 100}, Data: rawJSON(`{"variable":"{{answer}}","branches":[{"id":"b1","label":"Eskalasi","operator":"contains","value":"sambungkan,tim support,tidak tahu"}]}`)},
				{ID: "n3", Type: "transfer_agent", Position: FlowPosition{X: 500, Y: 100}, Data: rawJSON(`{}`)},
			},
			Edges: []FlowEdge{{ID: "e1", Source: "n1", Target: "n2"}, {ID: "e2", Source: "n2", Target: "n3", SourceH: "out-0"}},
		},
	},
}

func rawJSON(s string) json.RawMessage { return json.RawMessage(s) }

func handleFlowTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func handleFlowTemplateImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	templateID := r.FormValue("id")
	if templateID == "" {
		http.Error(w, "Missing template id", 400)
		return
	}
	for _, t := range templates {
		if t.ID == templateID {
			t.Flow.UserID = getUID(r)
			saveFlowDB(&t.Flow)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"id": t.Flow.ID, "name": t.Flow.Name, "status": "imported"})
			return
		}
	}
	http.Error(w, "Template not found", 404)
}
