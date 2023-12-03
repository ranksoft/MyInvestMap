package middleware

import (
	"context"
	"database/sql"
	"log"
	"myinvestmap/models"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("YOUR_SECRET_KEY"))

func JWTMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Authorization start")
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !userExists(claims.UserID, db) {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			log.Println("Authorization success")
			ctx := context.WithValue(r.Context(), "userClaims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userExists(userID int, db *sql.DB) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user availability: %v", err)
		return false
	}
	return exists
}
