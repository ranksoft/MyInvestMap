// /backend/handlers/assetHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myinvestmap/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func AddAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	var newAsset models.Asset
	err := json.NewDecoder(r.Body).Decode(&newAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertAssetSQL := `INSERT INTO assets (user_id, stockTag, exchange, price, quantity, IsPurchase) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertAssetSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	newAsset.IsPurchase = true
	_, err = statement.Exec(userID, newAsset.StockTag, newAsset.Exchange, newAsset.Price, newAsset.Quantity, newAsset.IsPurchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var symbols []string
	symbols = append(symbols, newAsset.StockTag)
	updateStockData(db, symbols, userID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAsset)
}

type updateStockRequest struct {
	Symbols []string `json:"symbols"`
}

func UpdateSelectedAssets(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateSelectedAssets")
	var updateStockRequest updateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStockRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	updateStockData(db, updateStockRequest.Symbols, userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func SellAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	var soldAsset models.Asset
	err := json.NewDecoder(r.Body).Decode(&soldAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertAssetSQL := `INSERT INTO assets (stockTag, exchange, price, quantity, IsPurchase) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertAssetSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	soldAsset.IsPurchase = false
	_, err = statement.Exec(soldAsset.StockTag, soldAsset.Exchange, soldAsset.Price, soldAsset.Quantity, soldAsset.IsPurchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var symbols []string
	symbols = append(symbols, soldAsset.StockTag)
	updateStockData(db, symbols, userID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(soldAsset)
}

func GetAssets(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	updateStockDataIfNeeded(db, userID)

	rows, err := db.Query("SELECT id, stockTag, exchange, price, quantity, isPurchase, name, currentPrice, createdAt, updatedAt FROM assets WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(&asset.ID, &asset.StockTag, &asset.Exchange, &asset.Price, &asset.Quantity, &asset.IsPurchase, &asset.Name, &asset.CurrentPrice, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func updateStockDataIfNeeded(db *sql.DB, userID int) {
	var lastUpdateStr string
	err := db.QueryRow("SELECT MAX(updatedAt) FROM assets").Scan(&lastUpdateStr)
	if err != nil {
		log.Println("Error fetching last update time:", err)
		return
	}

	lastUpdate, err := time.Parse("2006-01-02 15:04:05", lastUpdateStr)
	if err != nil {
		log.Println("Error parsing last update time:", err)
		return
	}

	if time.Since(lastUpdate) >= time.Minute {
		symbols := getSymbolsToUpdate(db, userID)
		if len(symbols) > 0 {
			updateStockData(db, symbols, userID)
		}
	}
}

func getSymbolsToUpdate(db *sql.DB, userID int) []string {
	rows, err := db.Query("SELECT DISTINCT stockTag FROM assets WHERE user_id = ? ORDER BY updatedAt ASC LIMIT 8", userID)
	if err != nil {
		log.Println("Error fetching symbols for update:", err)
		return nil
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var stockTag string
		if err := rows.Scan(&stockTag); err != nil {
			log.Println("Error scanning stockTag:", err)
			continue
		}
		symbols = append(symbols, stockTag)
	}
	return symbols
}

func updateStockData(db *sql.DB, symbols []string, userID int) {
	stockData, err := fetchStockDataFromAPI(db, symbols, userID)
	log.Println("fetchStockDataFromAPI:")
	log.Println(stockData)
	if err != nil {
		log.Println("Error fetching stock data:", err)
		return
	}

	for _, data := range stockData {
		if data.Price != "" {
			updateSQL := `UPDATE assets SET name = ?, currentPrice = ?, updatedAt = CURRENT_TIMESTAMP WHERE stockTag = ?`
			_, err := db.Exec(updateSQL, data.Name, data.Price, data.Symbol)
			if err != nil {
				log.Println("Error updating asset in database:", err)
			}
		}
	}
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

func fetchStockDataFromAPI(db *sql.DB, symbols []string, userID int) ([]StockQuote, error) {
	if len(symbols) == 0 {
		return nil, nil
	}
	var apiKey string
	err := db.QueryRow("SELECT api_key FROM api_keys WHERE user_id = ?", userID).Scan(&apiKey)
	if err != nil {
		return nil, fmt.Errorf("Error fetching API key: %v", err)
	}

	symbolsQuery := strings.Join(symbols, ",")
	apiURL := fmt.Sprintf("https://api.twelvedata.com/quote?symbol=%s&apikey=%s", symbolsQuery, apiKey)
	fmt.Println("API Request:", string(apiURL))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}
	fmt.Println("API Response:", string(bodyBytes))

	var quotes []StockQuote
	if len(symbols) == 1 {
		var quote StockQuoteApiResponce
		if err := json.Unmarshal(bodyBytes, &quote); err != nil {
			return nil, fmt.Errorf("JSON Decode error: %v", err)
		}

		quotes = append(quotes, StockQuote{
			Symbol: quote.Symbol,
			Name:   quote.Name,
			Price:  quote.Price,
		})
	} else {
		var response map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			return nil, fmt.Errorf("JSON Decode error: %v", err)
		}

		for symbol, data := range response {
			quoteData, ok := data.(map[string]interface{})
			if !ok {
				continue
			}
			quotes = append(quotes, StockQuote{
				Symbol: symbol,
				Name:   quoteData["name"].(string),
				Price:  quoteData["close"].(string),
			})
		}
	}

	return quotes, nil
}

func DeleteAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID
	vars := mux.Vars(r)
	id := vars["id"]

	deleteSQL := `DELETE FROM assets WHERE id = ? AND user_id = ?`
	_, err := db.Exec(deleteSQL, id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted")
}

func UpdateAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("userClaims").(*models.Claims)
	userID := userClaims.UserID

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	var updatedAsset models.Asset
	err = json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateSQL := `UPDATE assets SET stockTag = ?, exchange = ?, price = ?, quantity = ? WHERE id = ?`
	_, err = db.Exec(updateSQL, updatedAsset.StockTag, updatedAsset.Exchange, updatedAsset.Price, updatedAsset.Quantity, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var symbols []string
	symbols = append(symbols, updatedAsset.StockTag)
	updateStockData(db, symbols, userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAsset)
}
