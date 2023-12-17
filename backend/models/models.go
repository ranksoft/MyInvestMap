// /backend/models/models.go

package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type APIKeyRequest struct {
	APIKey string `json:"api_key"`
}

type APIKeyResponse struct {
	APIKey string `json:"api_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateStockRequest struct {
	Symbols []string `json:"symbols"`
}

type StockQuote struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Price  string `json:"price"`
}

type StockQuoteApiResponce struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Price  string `json:"close"`
}
