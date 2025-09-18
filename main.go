package main

import (
	"log"
	"os"

	"backend/db"
	"backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к БД
	db.Connect()

	// Создаём роутер Gin
	r := gin.Default()

	// Включаем CORS
	r.Use(cors.Default())

	// Роуты
	r.GET("/reports", handlers.GetReports)
	r.GET("/reports/:id", handlers.GetReport)
	r.POST("/reports", handlers.CreateReport)
	r.PATCH("/reports/:id/status", handlers.UpdateStatus)

	// Порт (Render передаёт через PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
