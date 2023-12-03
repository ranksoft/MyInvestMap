package handlers

import (
	"database/sql"
	"encoding/json"
	"myinvestmap/models"
	"net/http"
)

func SaveAPIKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	var apiKeyRequest models.APIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&apiKeyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	upsertQuery := `INSERT INTO api_keys (user_id, api_key) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET api_key = excluded.api_key`
	if _, err := db.Exec(upsertQuery, userID, apiKeyRequest.APIKey); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.APIKeyResponse{APIKey: apiKeyRequest.APIKey})
}

func GetAPIKey(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	var apiKey string
	query := `SELECT api_key FROM api_keys WHERE user_id = ?`
	if err := db.QueryRow(query, userID).Scan(&apiKey); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models.APIKeyResponse{APIKey: apiKey})
}
