package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID             int             `json:"id"`
	Type           string          `json:"type"`
	Source         string          `json:"source"`
	Payload        json.RawMessage `json:"payload"`
	IdempotencyKey string          `json:"idempotency_key"`
	ProcessedAt    *time.Time      `json:"processed_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

type EventInput struct {
	Type           string          `json:"type"`
	Source         string          `json:"source"`
	Payload        json.RawMessage `json:"payload"`
	IdempotencyKey string          `json:"idempotency_key"`
}

type Rule struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	EventType       string          `json:"event_type"`
	MatchConditions json.RawMessage `json:"match_conditions"`
	ActionType      string          `json:"action_type"`
	ActionPayload   json.RawMessage `json:"action_payload"`
	IsActive        bool            `json:"is_active"`
}

type Delivery struct {
	ID           int             `json:"id"`
	EventID      int             `json:"event_id"`
	RuleID       int             `json:"rule_id"`
	Channel      string          `json:"channel"`
	Recipient    string          `json:"recipient"`
	Payload      json.RawMessage `json:"payload"`
	Status       string          `json:"status"`
	AttemptCount int             `json:"attempt_count"`
	LastError    *string         `json:"last_error,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	SentAt       *time.Time      `json:"sent_at,omitempty"`
}

type RuleAudit struct {
	ID         int       `json:"id"`
	RuleID     int       `json:"rule_id"`
	EventID    int       `json:"event_id"`
	ActionType string    `json:"action_type"`
	Outcome    string    `json:"outcome"`
	ErrorMsg   *string   `json:"error_msg,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
