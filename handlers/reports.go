package handlers

import (
	"net/http"

	"github.com/Amirc0c/hike-reports/db"
	"github.com/Amirc0c/hike-reports/models"
	"github.com/gin-gonic/gin"
)

// Получить все отчёты
func GetReports(c *gin.Context) {
	var reports []models.Report
	if err := db.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отчетов"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// Получить один отчёт по ID
func GetReport(c *gin.Context) {
	id := c.Param("id")
	var report models.Report
	if err := db.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отчет не найден"})
		return
	}
	c.JSON(http.StatusOK, report)
}

// Создать новый отчёт
func CreateReport(c *gin.Context) {
	var report models.Report

	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if report.Status == "" {
		report.Status = "active"
	}

	if err := db.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании отчета"})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// Обновить статус отчёта
func UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var report models.Report

	if err := db.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отчет не найден"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report.Status = body.Status
	if err := db.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса"})
		return
	}

	c.JSON(http.StatusOK, report)
}
