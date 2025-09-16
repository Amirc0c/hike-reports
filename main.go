package main

import (
	"log"
	"net/http"
"os"
	"github.com/gorilla/mux"
	"backend/db"
	"backend/handlers"
)

func initializeRouter() *mux.Router {
	r := mux.NewRouter()

	// CRUD
	r.HandleFunc("/reports", handlers.CreateReport).Methods("POST")
	r.HandleFunc("/reports", handlers.GetReports).Methods("GET")
	r.HandleFunc("/reports/{id}", handlers.DeleteReport).Methods("DELETE")

	// обновление только статуса
	r.HandleFunc("/reports/{id}/status", handlers.UpdateReportStatus).Methods("PUT")

	return r
}

func main() {
	// инициализация базы
	db.InitDB()

	// роутер
	r := initializeRouter()

	port := os.Getenv("PORT")
if port == "" {
    port = "8080" // локально
}

log.Println("Server running on :" + port)
http.ListenAndServe(":"+port, r)

}
