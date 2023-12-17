// /backend/handlers/registerHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"myinvestmap/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

const (
	insertUserSQL = `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
)

func RegisterHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(insertUserSQL, user.Username, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
