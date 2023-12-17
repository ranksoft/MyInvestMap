// /backend/handlers/loginHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"myinvestmap/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	selectUserSQL = `SELECT id, password FROM users WHERE email = ?`
)

var jwtKey = []byte(os.Getenv("AUTH_SECRET_KEY"))

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword string
	if err := db.QueryRow(selectUserSQL, creds.Email).Scan(&userID, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.TokenResponse{Token: tokenString})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Logged out successfully"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
