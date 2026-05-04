package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/lib/pq" // for unmarshal
	"os"
)

type ChatRequest struct {
	SessionID  *int64  `json:"session_id"`
	CustomerID int64   `json:"customer_id"`
	Message    string  `json:"message"`
}

type ChatResponse struct {
	SessionID int64  `json:"session_id"`
	Reply     string `json:"reply"`
	Sentiment string `json:"sentiment"`
}

// ChatHandler handles AI chat using Claude
func ChatHandler(db *sql.DB) http.HandlerFunc {
	client := anthropic.NewClient(
		anthropic.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request")

			return
		}

		// Get or create session
		var sessionID int64
		if req.SessionID != nil {
			sessionID = *req.SessionID
		} else {
			// Create new session
			err := db.QueryRow("INSERT INTO sessions (customer_id) VALUES ($1) RETURNING id", req.CustomerID).Scan(&sessionID)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "Failed to create session")
				return
			}
		}

		// Save customer message
		_, err := db.Exec("INSERT INTO messages (session_id, sender, content) VALUES ($1, 'customer', $2)", sessionID, req.Message)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to save message")
			return
		}

		// Get conversation history
		rows, err := db.Query(`
			SELECT sender, content FROM messages 
			WHERE session_id = $1 
			ORDER BY created_at ASC
			LIMIT 20
		`, sessionID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to get history")
			return
		}
		defer rows.Close()

		var history []string
		for rows.Next() {
			var sender, content string
			rows.Scan(&sender, &content)
			if sender == "customer" {
				history = append(history, "Human: " + content)
			} else {
				history = append(history, "Assistant: " + content)
			}
		}

		// Claude conversation
		ctx := context.Background()
		messages := []anthropic.Message{{Role: anthropic.RoleUser, Content: anthropic.ContentText(req.Message)}}
		
		resp, err := client.Messages.Create(ctx, &anthropic.MessagesCreateRequest{
			Model: anthropic.ModelClaude35Sonnet20240620,

			MaxTokens: 1024,
			Messages:    messages,
			System: "You are Relay AI Agent for Indonesian CRM. ALWAYS reply in Bahasa Indonesia only. Analyze sentiment. Use tools if needed:

Tools:
1. lookup_customer: id → GET /api/customers/{id}
2. search_kb: query → SELECT from kb_articles LIKE title/content
3. create_ticket: customer_id, session_id, summary → POST /api/tickets

For negative sentiment ('Negatif', 'Sangat Negatif'), create ticket high priority 'keluhan'.

Be helpful, empathetic.",
			Tools: []anthropic.Tool{
				{
					Name: "lookup_customer",
					Description: "Lookup customer by ID",
					InputSchema: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"id": map[string]interface{}{"type": "integer"},
						},
						"required": []string{"id"},
					},
				},
				// Add more tools...
			},
		})
		if err != nil {
			respondError(w, http.StatusInternalServerError, "AI error")
			return
		}

		reply := resp.Content[0].Text
		sentiment := "Netral" // TODO analyze or from tool

		// Save agent message
		_, err = db.Exec("INSERT INTO messages (session_id, sender, content) VALUES ($1, 'agent', $2)", sessionID, reply)
		if err != nil {
			// non fatal
		}

		respondJSON(w, http.StatusOK, ChatResponse{
			SessionID: sessionID,
			Reply:     reply,
			Sentiment: sentiment,
		})
	}
}

