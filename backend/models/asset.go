// /backend/models/asset.go

package models

import (
	"database/sql"
	"time"
)

type Asset struct {
	ID           int             `json:"id"`
	StockTag     string          `json:"stockTag"`
	Exchange     string          `json:"exchange"`
	Name         sql.NullString  `json:"name"`
	Price        float64         `json:"price"`
	Quantity     float64         `json:"quantity"`
	CurrentPrice sql.NullFloat64 `json:"currentPrice"`
	IsPurchase   bool            `json:"isPurchase"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type AssetResponce struct {
	ID           int     `json:"id"`
	StockTag     string  `json:"stockTag"`
	Exchange     string  `json:"exchange"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Quantity     float64 `json:"quantity"`
	CurrentPrice float64 `json:"currentPrice"`
	IsPurchase   bool    `json:"isPurchase"`
}
