package handlers

import (
	"encoding/json"
	"net/http"
)

// EnableCORS middleware for CORS support
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// StubHandler returns a stub handler for unimplemented endpoints
func StubHandler(db interface{}, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondError(w, http.StatusNotImplemented, endpoint+" endpoint stub - coming soon")
	}
}
