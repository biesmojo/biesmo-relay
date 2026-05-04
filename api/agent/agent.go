package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/jackc/pgx/v5/pgxpool"

	"relay/api/models"
)

const SYSTEM_PROMPT = `Kamu adalah agen layanan pelanggan AI untuk perusahaan kami.
Tugasmu adalah membantu pelanggan dengan pertanyaan, keluhan, dan permintaan mereka.
Aturan yang harus kamu ikuti:

Selalu balas dalam Bahasa Indonesia yang sopan dan profesional
Gunakan tools yang tersedia untuk mendapatkan data — jangan mengarang informasi pelanggan atau produk
Jika pertanyaan tidak jelas setelah 1 kali bertanya, eskalasi ke manusia
Jika pelanggan meminta berbicara dengan manusia, langsung panggil escalate_to_human
Jangan membuat janji yang tidak bisa kamu penuhi

Saat sesi berakhir, kamu akan diminta menghasilkan: ringkasan, sentimen, dan prediksi CSAT.`

const END_SESSION_PROMPT = `Berdasarkan transkrip sesi ini, hasilkan:
1. ringkasan: 2-4 kalimat ringkasan dalam Bahasa Indonesia
2. sentimen: Positif | Netral | Negatif | Sangat Negatif
3. csat_score: integer 1-5
4. csat_rationale: satu kalimat penjelasan

Format JSON exact:
{
  "ringkasan": "...",
  "sentimen": "Positif",
  "csat_score": 4,
  "csat_rationale": "Customer puas dengan solusi."
}`

type ChatInput struct {
	SessionID   *int64  `json:"session_id"`
	Message     string  `json:"message"`
	CustomerPhone *string `json:"customer_phone"`
}

type ChatOutput struct {
	Reply       string `json:"reply"`
	SessionID   int64  `json:"session_id"`
	Summary     *string `json:"summary,omitempty"`
	Sentiment   *string `json:"sentiment,omitempty"`
	CSATScore   *int    `json:"csat_score,omitempty"`
	CSATRationale *string `json:"csat_rationale,omitempty"`
}

func ChatHandler(pool *pgxpool.Pool) http.HandlerFunc {
	client := anthropic.NewClient(
		an throp ic.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var input ChatInput
		json.NewDecoder(r.Body).Decode(&input)

		sessionID, customerID := getOrCreateSession(pool, input.CustomerPhone, input.SessionID)

		// Save user message
		_, err := pool.Exec(r.Context(), "INSERT INTO messages (session_id, role, content) VALUES ($1, 'user', $2)", sessionID, input.Message)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Save message failed")
			return
		}

		// Get history
		rows, err := pool.Query(r.Context(), "SELECT role, content FROM messages WHERE session_id = $1 ORDER BY created_at", sessionID)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "History failed")
			return
		}
		history := make([]anthropic.Message, 0)
		for rows.Next() {
			var role, content string
			rows.Scan(&role, &content)
			history = append(history, anthropic.Message{
				Role: anthropic.Role(role),
				Content: []anthropic.Content{{Type: "text", Text: content}},
			})
		}
		rows.Close()

		// Agent loop
		maxIterations := 10
		for i := 0; i < maxIterations; i++ {
			resp, err := client.Messages.Create(r.Context(), &anthropic.MessagesCreateRequest{
				Model: "claude-sonnet-4-20250514",
				MaxTokens: 1024,
				Messages: history,
				System: SYSTEM_PROMPT,
				Tools: tools,
			})
			if err != nil {
				log.Printf("Claude error: %v", err)
				break
			}

			for _, content := range resp.Content {
				if content.Type == "text" {
					reply := content.Text
					// Save assistant message
					_, _ = pool.Exec(r.Context(), "INSERT INTO messages (session_id, role, content) VALUES ($1, 'assistant', $2)", sessionID, reply)
					history = append(history, anthropic.Message{
						Role: anthropic.RoleAssistant,
						Content: []anthropic.Content{{Type: "text", Text: reply}},
					})
					RespondJSON(w, http.StatusOK, ChatOutput{SessionID: sessionID, Reply: reply})
					return
				}
				if content.Type == "tool_use" {
					toolResult := executeTool(pool, content.Name, content.Input)
					history = append(history, anthropic.Message{
						Role: anthropic.RoleUser,
						Content: []anthropic.Content{
							{
								Type: "tool_result",
								ToolUseID: content.ID,
								Content: []anthropic.Content{{Type: "text", Text: toolResult}},
							},
						},
					})
				}
				if content.Name == "escalate_to_human" {
					// End session
					updateSessionEnd(pool, sessionID)
					reply := "Terima kasih. Saya akan menghubungkan Anda dengan agen manusia."
					fireEvent(pool, "session.escalated", map[string]interface{}{"session_id": sessionID, "reason": content.Input["reason"]})
					RespondJSON(w, http.StatusOK, ChatOutput{SessionID: sessionID, Reply: reply})
					return
				}
			}
		}

		// End session
		updateSessionEnd(pool, sessionID)
		RespondJSON(w, http.StatusOK, ChatOutput{SessionID: sessionID, Reply: "Terima kasih atas pertanyaannya!"})
	}
}

var tools = []anthropic.Tool{
	{
		Type: "lookup_customer",
		Name: "lookup_customer",
		Description: "Lookup customer by phone or email",
		InputSchema: objectSchema(map[string]schema{
			"phone": {"type": "string"},
			"email": {"type": "string"},
		}),
	},
	{
		Type: "search_kb",
		Name: "search_kb",
		Description: "Search knowledge base articles",
		InputSchema: objectSchema(map[string]schema{
			"query": {"type": "string"},
			"limit": {"type": "integer", "default": 3},
		}),
	},
	{
		Type: "create_ticket",
		Name: "create_ticket",
		Description: "Create or append to customer ticket",
		InputSchema: objectSchema(map[string]schema{
			"customer_id": {"type": "integer"},
			"category": {"type": "string"},
			"priority": {"type": "string"},
			"summary": {"type": "string"},
		}),
	},
	{
		Type: "escalate_to_human",
		Name: "escalate_to_human",
		Description: "Escalate to human agent",
		InputSchema: objectSchema(map[string]schema{
			"reason": {"type": "string"},
			"urgency": {"type": "string", "enum": ["low", "high"]},
		}),
	},
}

func executeTool(pool *pgxpool.Pool, name string, input map[string]interface{}) string {
	switch name {
	case "lookup_customer":
		phone, _ := input["phone"].(string)
		email, _ := input["email"].(string)
		var customer models.Customer
		query := "SELECT id, name, phone, email FROM customers WHERE phone = $1 OR email = $2"
		err := pool.QueryRow(context.Background(), query, phone, email).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Email)
		if err != nil {
			return "Customer not found"
		}
		data, _ := json.Marshal(customer)
		return string(data)
	case "search_kb":
		query, _ := input["query"].(string)
		limit := 3
		if l, ok := input["limit"].(float64); ok {
			limit = int(l)
		}
		rows, err := pool.Query(context.Background(), "SELECT id, title, LEFT(content, 300) as excerpt FROM kb_articles WHERE content ILIKE $1 OR title ILIKE $1 LIMIT $2", "%"+query+"%", limit)
		if err != nil {
			return "Search failed"
		}
		defer rows.Close()
		var results []map[string]interface{}
		for rows.Next() {
			var id int
			var title, excerpt string
			rows.Scan(&id, &title, &excerpt)
			results = append(results, map[string]interface{}{"id": id, "title": title, "excerpt": excerpt})
		}
		data, _ := json.Marshal(results)
		return string(data)
	case "create_ticket":
		cid := int(input["customer_id"].(float64))
		category := input["category"].(string)
		priority := input["priority"].(string)
		summary := input["summary"].(string)
		// Check existing open
		var existingID int
		err := pool.QueryRow(context.Background(), 
			"SELECT id FROM tickets WHERE customer_id = $1 AND category = $2 AND status != 'resolved' AND created_at > NOW() - INTERVAL '24 hours' ORDER BY created_at DESC LIMIT 1", 
			cid, category).Scan(&existingID)
		if err == nil {
			_, _ = pool.Exec(context.Background(), "UPDATE tickets SET summary = summary || '\n\n' || $1 WHERE id = $2", summary, existingID)
			return fmt.Sprintf(`{"ticket_id": %d, "action": "appended"}`, existingID)
		}
		var newID int
		err = pool.QueryRow(context.Background(), 
			"INSERT INTO tickets (customer_id, category, priority, summary) VALUES ($1, $2, $3, $4) RETURNING id", 
			cid, category, priority, summary).Scan(&newID)
		if err != nil {
			return "Create ticket failed"
		}
		return fmt.Sprintf(`{"ticket_id": %d, "action": "created"}`, newID)
	case "escalate_to_human":
		_, _ = pool.Exec(context.Background(), "UPDATE sessions SET status = 'escalated' WHERE id = $1", ??? ) // from context sessionID
		fireEvent(...) // implement fire event
		return "Escalation confirmed, human agent notified"
	}
	return "Unknown tool"
}

func updateSessionEnd(pool *pgxpool.Pool, sessionID int) {
	// Get transcript
	rows, err := pool.Query(context.Background(), "SELECT content FROM messages WHERE session_id = $1 ORDER BY created_at", sessionID)
	if err != nil {
		return
	}
	var transcript strings.Builder
	for rows.Next() {
		var content string
		rows.Scan(&content)
		transcript.WriteString(content + "\n")
	}
	rows.Close()

	// Claude session end
	client.Messages.Create(...) with END_SESSION_PROMPT input transcript
	parse JSON update session summary sentiment csat_score csat_rationale
	fireEvent "session.closed" payload {session_id, sentiment}
}

func getOrCreateSession(pool *pgxpool.Pool, phone *string, sessionID *int64) (int64, int64) {
	if sessionID != nil {
		var cid int64
		pool.QueryRow(context.Background(), "SELECT customer_id FROM sessions WHERE id = $1", *sessionID).Scan(&cid)
		return *sessionID, cid
	}
	// Lookup customer
	var cid int64
	if phone != nil {
		pool.QueryRow(context.Background(), "SELECT id FROM customers WHERE phone = $1", *phone).Scan(&cid)
	}
	if cid == 0 {
		// Create anonymous
		pool.QueryRow(context.Background(), "INSERT INTO customers (name, channel) VALUES ('Anonymous', 'web') RETURNING id").Scan(&cid)
	}
	var newSession int64
	pool.QueryRow(context.Background(), "INSERT INTO sessions (customer_id) VALUES ($1) RETURNING id", cid).Scan(&newSession)
	return newSession, cid
}

func fireEvent(pool *pgxpool.Pool, typ string, payload interface{}) {
	idempKey := fmt.Sprintf("%s-%d", typ, time.Now().UnixNano())
	pBytes, _ := json.Marshal(payload)
	_, _ = pool.Exec(context.Background(), "INSERT INTO events (type, payload, idempotency_key) VALUES ($1, $2, $3)", typ, pBytes, idempKey)
}

