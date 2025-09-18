package main

import (
	"log"
	"net/http"
	"os"
"fmt"
	"github.com/gorilla/mux"
	"backend/db"
	"backend/handlers"
)

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func initializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/reports", handlers.CreateReport).Methods("POST")
	r.HandleFunc("/reports", handlers.GetReports).Methods("GET")
	r.HandleFunc("/reports/{id}", handlers.DeleteReport).Methods("DELETE")
	r.HandleFunc("/reports/{id}/status", handlers.UpdateReportStatus).Methods("PATCH")

	return r
}

func main() {
	db.InitDB() // подключаем БД

	r := initializeRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
fmt.Println("работай")
	log.Println("Server running on :" + port)
	// оборачиваем mux в CORS middleware
	http.ListenAndServe(":"+port, withCORS(r))
}
