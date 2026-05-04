package handlers

import (
	"net/http"

	"relay/api/db"
	"relay/api/models"
)

func AdminEventsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := pool.Query(r.Context(), `
			SELECT 
				e.id, e.type, e.source, e.payload, e.created_at,
				d.id as delivery_id, d.status as delivery_status,
				ral.id as audit_id, ral.outcome, ral.error_msg
			FROM events e
			LEFT JOIN deliveries d ON e.id = d.event_id
			LEFT JOIN rule_audit_log ral ON e.id = ral.event_id
			ORDER BY e.created_at DESC LIMIT 50
		`)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to fetch events")
			return
		}
		defer rows.Close()

		events := []models.EventWithJoins{}
		for rows.Next() {
			var ev models.EventWithJoins
			rows.Scan(&ev.ID, &ev.Type, &ev.Source, &ev.Payload, &ev.CreatedAt, &ev.DeliveryID, &ev.DeliveryStatus, &ev.AuditID, &ev.Outcome, &ev.ErrorMsg)
			events = append(events, ev)
		}

		RespondJSON(w, http.StatusOK, events)
	}
}

type EventWithJoins struct {
	models.Event
	DeliveryID     *int    `json:"delivery_id"`
	DeliveryStatus *string `json:"delivery_status"`
	AuditID        *int    `json:"audit_id"`
	Outcome        *string `json:"outcome"`
	ErrorMsg       *string `json:"error_msg"`
}
