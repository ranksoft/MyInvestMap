// /backend/handlers/assetHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"myinvestmap/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	insertAssetSQL      = `INSERT INTO assets (user_id, stockTag, exchange, price, quantity, IsPurchase) VALUES (?, ?, ?, ?, ?, ?)`
	selectAssetsSQL     = `SELECT id, stockTag, exchange, price, quantity, isPurchase, name, currentPrice, createdAt, updatedAt FROM assets WHERE user_id = ?`
	selectMaxUpdateSQL  = `SELECT MAX(updatedAt) FROM assets`
	distinctStockTagSQL = `SELECT DISTINCT stockTag FROM assets WHERE user_id = ? ORDER BY updatedAt ASC LIMIT 8`
	updateAssetSQL      = `UPDATE assets SET name = ?, currentPrice = ?, updatedAt = CURRENT_TIMESTAMP WHERE stockTag = ?`
	deleteAssetSQL      = `DELETE FROM assets WHERE id = ? AND user_id = ?`
	updateAssetByIDSQL  = `UPDATE assets SET stockTag = ?, exchange = ?, price = ?, quantity = ? WHERE id = ?`
	selectAPIKeySQL     = `SELECT api_key FROM api_keys WHERE user_id = ?`
)

func AddAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	var newAsset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&newAsset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	statement, err := db.Prepare(insertAssetSQL)
	if err != nil {
		http.Error(w, "failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	newAsset.IsPurchase = true
	if _, err = statement.Exec(userClaims.UserID, newAsset.StockTag, newAsset.Exchange, newAsset.Price, newAsset.Quantity, newAsset.IsPurchase); err != nil {
		http.Error(w, "failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	updateStockData(db, []string{newAsset.StockTag}, userClaims.UserID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAsset)
}

func UpdateSelectedAssets(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req models.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	updateStockData(db, req.Symbols, userClaims.UserID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func SellAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	var soldAsset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&soldAsset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	statement, err := db.Prepare(insertAssetSQL)
	if err != nil {
		http.Error(w, "failed to prepare SQL statement", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	soldAsset.IsPurchase = false
	if _, err = statement.Exec(userClaims.UserID, soldAsset.StockTag, soldAsset.Exchange, soldAsset.Price, soldAsset.Quantity, soldAsset.IsPurchase); err != nil {
		http.Error(w, "failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	updateStockData(db, []string{soldAsset.StockTag}, userClaims.UserID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(soldAsset)
}

func GetAssets(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	updateStockDataIfNeeded(db, userClaims.UserID)

	rows, err := db.Query(selectAssetsSQL, userClaims.UserID)
	if err != nil {
		http.Error(w, "failed to query assets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		if err := rows.Scan(&asset.ID, &asset.StockTag, &asset.Exchange, &asset.Price, &asset.Quantity, &asset.IsPurchase, &asset.Name, &asset.CurrentPrice, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			http.Error(w, "failed to scan asset row", http.StatusInternalServerError)
			return
		}
		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "error iterating over assets rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func UpdateAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid asset ID", http.StatusBadRequest)
		return
	}

	var updatedAsset models.Asset
	if err := json.NewDecoder(r.Body).Decode(&updatedAsset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if _, err = db.Exec(updateAssetByIDSQL, updatedAsset.StockTag, updatedAsset.Exchange, updatedAsset.Price, updatedAsset.Quantity, id); err != nil {
		http.Error(w, "failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	updateStockData(db, []string{updatedAsset.StockTag}, userClaims.UserID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAsset)
}

func DeleteAsset(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
	if !ok {
		http.Error(w, "invalid user claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := db.Exec(deleteAssetSQL, id, userClaims.UserID); err != nil {
		http.Error(w, "failed to execute SQL statement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted")
}

func updateStockDataIfNeeded(db *sql.DB, userID int) error {
	var lastUpdateStr string
	err := db.QueryRow(selectMaxUpdateSQL).Scan(&lastUpdateStr)
	if err != nil {
		return fmt.Errorf("error fetching last update time: %v", err)
	}

	lastUpdate, err := time.Parse("2006-01-02 15:04:05", lastUpdateStr)
	if err != nil {
		return fmt.Errorf("error parsing last update time: %v", err)
	}

	if time.Since(lastUpdate) >= time.Minute {
		symbols, err := getSymbolsToUpdate(db, userID)
		if err != nil {
			return err
		}
		if len(symbols) > 0 {
			if err := updateStockData(db, symbols, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func getSymbolsToUpdate(db *sql.DB, userID int) ([]string, error) {
	rows, err := db.Query(distinctStockTagSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching symbols for update: %v", err)
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var stockTag string
		if err := rows.Scan(&stockTag); err != nil {
			return nil, fmt.Errorf("error scanning stockTag: %v", err)
		}
		symbols = append(symbols, stockTag)
	}
	return symbols, nil
}

func updateStockData(db *sql.DB, symbols []string, userID int) error {
	stockData, err := fetchStockDataFromAPI(db, symbols, userID)
	if err != nil {
		return fmt.Errorf("error fetching stock data: %v", err)
	}

	for _, data := range stockData {
		if data.Price != "" {
			if _, err := db.Exec(updateAssetSQL, data.Name, data.Price, data.Symbol); err != nil {
				return fmt.Errorf("error updating asset in database: %v", err)
			}
		}
	}
	return nil
}

func fetchStockDataFromAPI(db *sql.DB, symbols []string, userID int) ([]models.StockQuote, error) {
	if len(symbols) == 0 {
		return nil, nil
	}
	var apiKey string
	err := db.QueryRow(selectAPIKeySQL, userID).Scan(&apiKey)
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

	var quotes []models.StockQuote
	if len(symbols) == 1 {
		var quote models.StockQuoteApiResponce
		if err := json.Unmarshal(bodyBytes, &quote); err != nil {
			return nil, fmt.Errorf("JSON Decode error: %v", err)
		}

		quotes = append(quotes, models.StockQuote{
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
			quotes = append(quotes, models.StockQuote{
				Symbol: symbol,
				Name:   quoteData["name"].(string),
				Price:  quoteData["close"].(string),
			})
		}
	}

	return quotes, nil
}
