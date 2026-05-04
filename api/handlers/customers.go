package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Customer represents a CRM customer
type Customer struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}

// GetCustomers lists all customers
func GetCustomers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			// Demo mode - return mock data
			mock := []Customer{
				{ID: 1, Name: "Budi Santoso", Email: "budi@example.com", Phone: "+6281234567890"},
				{ID: 2, Name: "Ani Wijaya", Email: "ani@example.com", Phone: "+6282345678901"},
				{ID: 3, Name: "Dewi Lestari", Email: "dewi@example.com", Phone: "+6283456789012"},
			}
			RespondJSON(w, http.StatusOK, mock)
			return
		}
		rows, err := db.Query("SELECT id, name, email, phone, created_at FROM customers ORDER BY id")
		if err != nil {
			RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()
		var customers []Customer
		for rows.Next() {
			var c Customer
			var createdAt string
			err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &createdAt)
			if err != nil {
				RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}
			c.CreatedAt = createdAt
			customers = append(customers, c)
		}
		RespondJSON(w, http.StatusOK, customers)
	}
}

// CreateCustomer creates new customer
func CreateCustomer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			RespondError(w, http.StatusServiceUnavailable, "Database not available")
			return
		}
		var c Customer
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		var id int64
		err := db.QueryRow("INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING id", c.Name, c.Email, c.Phone).Scan(&id)
		if err != nil {
			RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		c.ID = id
		RespondJSON(w, http.StatusCreated, c)
	}
}

// GetCustomerByID gets single customer
func GetCustomerByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		if db == nil {
			if id == 1 {
				RespondJSON(w, http.StatusOK, Customer{ID: 1, Name: "Budi Santoso", Email: "budi@example.com"})
				return
			}
			if id == 2 {
				RespondJSON(w, http.StatusOK, Customer{ID: 2, Name: "Ani Wijaya", Email: "ani@example.com"})
				return
			}
			if id == 3 {
				RespondJSON(w, http.StatusOK, Customer{ID: 3, Name: "Dewi Lestari", Email: "dewi@example.com"})
				return
			}
			RespondError(w, http.StatusNotFound, "Customer not found")
			return
		}
		var c Customer
		var createdAt string
		err := db.QueryRow("SELECT id, name, email, phone, created_at FROM customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &createdAt)
		if err == sql.ErrNoRows {
			RespondError(w, http.StatusNotFound, "Customer not found")
			return
		} else if err != nil {
			RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		c.CreatedAt = createdAt
		RespondJSON(w, http.StatusOK, c)
	}
}
