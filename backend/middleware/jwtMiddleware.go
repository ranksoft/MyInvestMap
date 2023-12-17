// /backend/middleware/jwtMiddleware.go

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

const (
	checkUserExistsSQL = `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
)

var jwtKey = []byte(os.Getenv("AUTH_SECRET_KEY"))

func JWTMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				unauthorized(w, "Authorization header is missing")
				return
			}

			claims, err := validateToken(tokenString)
			if err != nil {
				unauthorized(w, "Invalid token: "+err.Error())
				return
			}

			if !userExists(claims.UserID, db) {
				unauthorized(w, "User not found")
				return
			}

			log.Println("Authorization successful for user:", claims.UserID)
			ctx := context.WithValue(r.Context(), "userClaims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	return strings.TrimPrefix(tokenString, "Bearer ")
}

func validateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func userExists(userID int, db *sql.DB) bool {
	var exists bool
	err := db.QueryRow(checkUserExistsSQL, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return false
	}
	return exists
}

func unauthorized(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusUnauthorized)
}
