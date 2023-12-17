// /backend/handlers/assetHandler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"myinvestmap/models"
	"net/http"
)

const (
	upsertAPIKeySQL = `INSERT INTO api_keys (user_id, api_key) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET api_key = excluded.api_key`
	getAPIKeySQL    = `SELECT api_key FROM api_keys WHERE user_id = ?`
)

func SaveAPIKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	var apiKeyRequest models.APIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&apiKeyRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := db.Exec(upsertAPIKeySQL, userClaims.UserID, apiKeyRequest.APIKey); err != nil {
		http.Error(w, fmt.Sprintf("error saving API key: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.APIKeyResponse{APIKey: apiKeyRequest.APIKey})
}

func GetAPIKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	var apiKey string
	if err := db.QueryRow(getAPIKeySQL, userClaims.UserID).Scan(&apiKey); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("error retrieving API key: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.APIKeyResponse{APIKey: apiKey})
}
