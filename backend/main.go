package main

import (
	"log"
	"myinvestmap/database"
	"myinvestmap/handlers"
	"myinvestmap/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	db := database.InitDB("data/database.db")
	defer db.Close()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(db, w, r)
	}).Methods("POST")

	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(db, w, r)
	}).Methods("POST")

	secureApi := r.PathPrefix("/api").Subrouter()
	secureApi.Use(middleware.JWTMiddleware(db))

	secureApi.HandleFunc("/api-key", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAPIKey(db, w, r)
	}).Methods("GET")

	secureApi.HandleFunc("/api-key", func(w http.ResponseWriter, r *http.Request) {
		handlers.SaveAPIKey(db, w, r)
	}).Methods("POST")

	secureApi.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r)
	}).Methods("POST")

	secureApi.HandleFunc("/refresh-assets", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateSelectedAssets(db, w, r)
	}).Methods("POST")

	secureApi.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAssets(db, w, r)
	}).Methods("GET")

	secureApi.HandleFunc("/assets/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddAsset(db, w, r)
	}).Methods("POST")

	secureApi.HandleFunc("/assets/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateAsset(db, w, r)
	}).Methods("PUT")

	secureApi.HandleFunc("/assets/sell", func(w http.ResponseWriter, r *http.Request) {
		handlers.SellAsset(db, w, r)
	}).Methods("POST")

	secureApi.HandleFunc("/assets/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteAsset(db, w, r)
	}).Methods("DELETE")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://myinvestmap.local:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := corsHandler.Handler(r)

	log.Println("Server is running on http://myinvestmap.local:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
