package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/users", createUser(db)).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("DB query error:", err)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
				log.Println("DB scan error:", err)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading database rows", http.StatusInternalServerError)
			log.Println("DB read error:", err)
			return
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("JSON encode error:", err)
		}
	}
}

// Get user by ID
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			log.Println("DB query error:", err)
			return
		}

		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("JSON encode error:", err)
		}
	}
}

// Create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("DB insert error:", err)
			return
		}

		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("JSON encode error:", err)
		}
	}
}

// Update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("DB update error:", err)
			return
		}

		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("JSON encode error:", err)
		}
	}
}

// Delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			log.Println("DB query error:", err)
			return
		}

		_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			log.Println("DB delete error:", err)
			return
		}

		if err := json.NewEncoder(w).Encode("User deleted"); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("JSON encode error:", err)
		}
	}
}
