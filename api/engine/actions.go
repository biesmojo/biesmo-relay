package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"relay/api/db"
	"relay/api/models"
)

func ExecuteAction(pool *pgxpool.Pool, rule models.Rule, event models.Event) {
	// Create delivery
	var deliveryID int
	err := pool.QueryRow(context.Background(),
		"INSERT INTO deliveries (event_id, rule_id, channel, recipient, payload) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		event.ID, rule.ID, "web", "", rule.ActionPayload).Scan(&deliveryID)
	if err != nil {
		logRuleAudit(pool, rule.ID, event.ID, rule.ActionType, "delivery_failed", &err.Error())
		return
	}

	attempt := 0
	maxAttempts := 3
	backoff := 2 * time.Second

	for attempt < maxAttempts {
		attempt++
		success, errMsg := performAction(pool, rule.ActionType, rule.ActionPayload, event.Payload)
		if success {
			_, _ = pool.Exec(context.Background(), "UPDATE deliveries SET status = 'sent' WHERE id = $1", deliveryID)
			logRuleAudit(pool, rule.ID, event.ID, rule.ActionType, "success", nil)
			return
		}

		if attempt < maxAttempts {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	_, _ = pool.Exec(context.Background(),
		"UPDATE deliveries SET status = 'failed', attempt_count = $1, last_error = $2 WHERE id = $3",
		attempt, errMsg, deliveryID)
	logRuleAudit(pool, rule.ID, event.ID, rule.ActionType, "failed", &errMsg)
}

func performAction(pool *pgxpool.Pool, actionType string, payload, eventPayload json.RawMessage) (bool, string) {
	var p map[string]interface{}
	json.Unmarshal(payload, &p)
	var ep map[string]interface{}
	json.Unmarshal(eventPayload, &ep)

	switch actionType {
	case "create_ticket":
		_, err := pool.Exec(context.Background(),
			"INSERT INTO tickets (customer_id, session_id, category, priority, summary) VALUES ($1, $2, $3, $4, $5)",
			ep["customer_id"], ep["session_id"], p["category"], p["priority"], p["summary"])
		if err != nil {
			return false, err.Error()
		}
		return true, ""

	case "send_message":
		// Stub - real channel adapter later
		return true, ""

	case "update_customer":
		_, err := pool.Exec(context.Background(), "UPDATE customers SET phone = $1 WHERE id = $2", p["phone"], ep["customer_id"])
		if err != nil {
			return false, err.Error()
		}
		return true, ""

	case "webhook":
		url, ok := p["url"].(string)
		if !ok {
			return false, "no url in payload"
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		if err != nil {
			return false, err.Error()
		}
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return false, err.Error()
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return true, ""
		}
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return false, "unknown action_type"
}
