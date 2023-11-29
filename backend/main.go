package main

import (
	"log"
	"myinvestmap/database"
	"myinvestmap/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db := database.InitDB("data/database.db")
	defer db.Close()

	r := mux.NewRouter()

	// r.HandleFunc("/api/stocks", handlers.GetBatchStockData).Methods("GET")

	r.HandleFunc("/api/assets/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddAsset(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/api/assets/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateAsset(db, w, r)
	}).Methods("PUT")

	r.HandleFunc("/api/assets/sell", func(w http.ResponseWriter, r *http.Request) {
		handlers.SellAsset(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/api/assets", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAssets(db, w, r)
	}).Methods("GET")

	r.HandleFunc("/api/assets/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteAsset(db, w, r)
	}).Methods("DELETE")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://myinvestmap.local:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := corsHandler.Handler(r)

	log.Println("Server is running on http://myinvestmap.local:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
