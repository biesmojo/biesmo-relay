package handlers

import (
	"encoding/json"
	"net/http"

	"relay/api/db"
	"relay/api/engine"
	"relay/api/models"
)

type EventsResponse struct {
	ID int `json:"id"`
}

func EventsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.EventInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			RespondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if input.IdempotencyKey == "" {
			RespondError(w, http.StatusBadRequest, "idempotency_key required")
			return
		}

		// Check idempotency
		var existingID int
		err := pool.QueryRow(r.Context(), "SELECT id FROM events WHERE idempotency_key = $1", input.IdempotencyKey).Scan(&existingID)
		if err == nil {
			// Duplicate, return existing
			RespondJSON(w, http.StatusOK, EventsResponse{ID: existingID})
			return
		}

		// Insert new
		var newID int
		err = pool.QueryRow(r.Context(),
			"INSERT INTO events (type, source, payload, idempotency_key) VALUES ($1, $2, $3, $4) RETURNING id",
			input.Type, input.Source, input.Payload, input.IdempotencyKey).Scan(&newID)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to create event")
			return
		}

		// Async rule evaluation
		go func() {
			event := models.Event{
				ID:             newID,
				Type:           input.Type,
				Source:         input.Source,
				Payload:        input.Payload,
				IdempotencyKey: input.IdempotencyKey,
			}
			engine.EvaluateRules(pool, event)
		}()

		RespondJSON(w, http.StatusCreated, EventsResponse{ID: newID})
	}
}
