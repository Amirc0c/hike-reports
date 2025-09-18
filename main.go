package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/db"
	"backend/handlers"
)

func main() {
	// Подключение к БД
	db.InitDB()

	r := gin.Default()

	// Настройки CORS
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // можешь указать свой фронт, например "http://localhost:3000"
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	// Маршруты
	r.GET("/reports", handlers.GetReports)
	r.GET("/reports/:id", handlers.GetReport)
	r.POST("/reports", handlers.CreateReport)
	r.PATCH("/reports/:id", handlers.UpdateReportStatus)
	r.DELETE("/reports/:id", handlers.DeleteReport)

	// Запуск сервера
	r.Run(":8080")
}
