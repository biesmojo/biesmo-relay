package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Customer represents a CRM customer
type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt string   `json:"created_at"`
}

// Session represents a chat session
type Session struct {
	ID          int64      `json:"id"`
	CustomerID  int64      `json:"customer_id"`
	Status      string     `json:"status"` // active, ended
	Summary     string     `json:"summary,omitempty"`
	Sentiment   string     `json:"sentiment,omitempty"`
	CreatedAt   string     `json:"created_at"`
	EndedAt     *string    `json:"ended_at,omitempty"`
}

// Message represents a chat message
type Message struct {
	ID         int64     `json:"id"`
	SessionID  int64     `json:"session_id"`
	Sender     string    `json:"sender"`
	Content    string    `json:"content"`
	CreatedAt  string    `json:"created_at"`
}

// Ticket represents a support ticket
type Ticket struct {
	ID          int64     `json:"id"`
	CustomerID  int64     `json:"customer_id"`
	SessionID   *int64    `json:"session_id,omitempty"`
	Category    string    `json:"category"`
	Priority    string    `json:"priority"`
	Summary     string    `json:"summary"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// KBArticle represents knowledge base articles
type KBArticle struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Tags     string `json:"tags"`
}

// ChatRequest for web chat
type ChatRequest struct {
	SessionID  int64  `json:"session_id,omitempty"`
	CustomerID int64  `json:"customer_id,omitempty"`
	Message    string `json:"message"`
}

// ChatResponse from agent
type ChatResponse struct {
	SessionID  int64     `json:"session_id"`
	Reply      string    `json:"reply"`
	Sentiment  string    `json:"sentiment"`
}

var db *sql.DB

// EnableCORS middleware
func enableCORS(next http.Handler) http.Handler {
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

// RespondJSON sends JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// RespondError sends error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// StubHandler returns stub handlers
func stubHandler(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusNotImplemented, "Endpoint stub - coming soon")
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "relay")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Warning: Could not connect to DB: %v", err)
		log.Println("Running in demo mode without database")
		db = nil
	} else {
		if err = db.Ping(); err != nil {
			log.Printf("Warning: Could not ping DB: %v", err)
			db = nil
		} else {
			log.Println("Connected to database")
		}
	}

	// Router
	r := mux.NewRouter()
	r.Use(enableCORS)

	// Root endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"service": "Relay API",
			"status":  "running",
			"version": "1.0.0",
		})
	}).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "healthy",
			"db":     db != nil,
		})
	}).Methods("GET")

	// Customers endpoints
	r.HandleFunc("/api/customers", getCustomers).Methods("GET")
	r.HandleFunc("/api/customers", createCustomer).Methods("POST")
	r.HandleFunc("/api/customers/{id}", getCustomerByID).Methods("GET")

	// Sessions endpoints (stub)
	r.HandleFunc("/api/sessions", stubHandler).Methods("GET", "POST")
	r.HandleFunc("/api/sessions/{id}", stubHandler).Methods("GET")

	// Tickets endpoints (stub)
	r.HandleFunc("/api/tickets", stubHandler).Methods("GET", "POST")
	r.HandleFunc("/api/tickets/{id}", stubHandler).Methods("GET")

	// Messages endpoints (stub)
	r.HandleFunc("/api/messages", stubHandler).Methods("GET", "POST")

	// KB endpoints (stub)
	r.HandleFunc("/api/kb", stubHandler).Methods("GET", "POST")

	// Events endpoints (stub)
	r.HandleFunc("/api/events", stubHandler).Methods("GET", "POST")

	// Chat endpoint (stub)
	r.HandleFunc("/api/chat", stubHandler).Methods("POST")

	// Rules endpoint (stub)
	r.HandleFunc("/api/rules", stubHandler).Methods("GET")

	// Admin endpoints (stub)
	r.HandleFunc("/api/admin/sessions", stubHandler).Methods("GET")
	r.HandleFunc("/api/admin/events", stubHandler).Methods("GET")
	r.HandleFunc("/api/admin/rules", stubHandler).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("Relay API starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Customer handlers
func getCustomers(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		// Demo mode
		mock := []Customer{
			{ID: 1, Name: "Budi Santoso", Email: "budi@example.com", Phone: "+6281234567890"},
			{ID: 2, Name: "Ani Wijaya", Email: "ani@example.com", Phone: "+6282345678901"},
			{ID: 3, Name: "Dewi Lestari", Email: "dewi@example.com", Phone: "+6283456789012"},
		}
		respondJSON(w, http.StatusOK, mock)
		return
	}
	rows, err := db.Query("SELECT id, name, email, phone, created_at FROM customers ORDER BY id")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var c Customer
		var createdAt string
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &createdAt)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		c.CreatedAt = createdAt
		customers = append(customers, c)
	}
	respondJSON(w, http.StatusOK, customers)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	var id int64
	err := db.QueryRow("INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id", c.Name, c.Email, c.Phone).Scan(&id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.ID = id
	c.CreatedAt = time.Now().Format(time.RFC3339)
	respondJSON(w, http.StatusCreated, c)
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	if db == nil {
		// Demo mode
		if id == 1 {
			respondJSON(w, http.StatusOK, Customer{ID: 1, Name: "Budi Santoso", Email: "budi@example.com"})
			return
		}
		if id == 2 {
			respondJSON(w, http.StatusOK, Customer{ID: 2, Name: "Ani Wijaya", Email: "ani@example.com"})
			return
		}
		if id == 3 {
			respondJSON(w, http.StatusOK, Customer{ID: 3, Name: "Dewi Lestari", Email: "dewi@example.com"})
			return
		}
		respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	var c Customer
	var createdAt string
	err := db.QueryRow("SELECT id, name, email, phone, created_at FROM customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &createdAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Customer not found")
		return
	} else if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	c.CreatedAt = createdAt
	respondJSON(w, http.StatusOK, c)
}
