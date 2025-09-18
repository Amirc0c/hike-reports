package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
"fmt"
	"backend/handlers"
)

func main() {
	r := gin.Default()

	// ✅ Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://твойфронт.vercel.app"}, // замени на свой фронт
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
fmt.Println("кал")
	// ✅ Роуты
	r.GET("/reports", handlers.GetReports)             // получить все отчёты
	r.GET("/reports/:id", handlers.GetReport)          // получить один отчёт
	r.POST("/reports", handlers.CreateReport)          // создать отчёт
	r.PATCH("/reports/:id/status", handlers.UpdateStatus) // обновить статус и вернуть отчёт

	// ✅ Запуск сервера
	r.Run(":8080")
}
