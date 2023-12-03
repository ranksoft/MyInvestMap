package models

type APIKeyRequest struct {
	APIKey string `json:"api_key"`
}

type APIKeyResponse struct {
	APIKey string `json:"api_key"`
}
