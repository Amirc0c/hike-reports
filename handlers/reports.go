package handlers

import (
	"backend/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Структуры
type Checkpoint struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type GroupMember struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Telegram string `json:"telegram"`
}

type Report struct {
	ID            int64         `json:"id"`
	RouteName     string        `json:"route_name"`
	GpxFile       string        `json:"gpx_file"`
	Checkpoints   []Checkpoint  `json:"checkpoints"`
	MustContactBy string        `json:"must_contact_by"`
	Status        string        `json:"status"`
	Grp           []GroupMember `json:"grp"`
}

// CREATE
func CreateReport(c *gin.Context) {
	var report Report
	if err := c.BindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if report.Status == "" {
		report.Status = "active"
	}

	grpBytes, _ := json.Marshal(report.Grp)
	checkpointsBytes, _ := json.Marshal(report.Checkpoints)

	err := db.DB.QueryRow(
		`INSERT INTO reports (route_name, gpx_file, checkpoints, must_contact_by, status, grp) 
         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		report.RouteName,
		report.GpxFile,
		checkpointsBytes,
		report.MustContactBy,
		report.Status,
		grpBytes,
	).Scan(&report.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// READ (все)
func GetReports(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp FROM reports ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var report Report
		var checkpointsJSON, grpJSON []byte

		rows.Scan(&report.ID, &report.RouteName, &report.GpxFile, &checkpointsJSON, &report.MustContactBy, &report.Status, &grpJSON)

		json.Unmarshal(checkpointsJSON, &report.Checkpoints)
		json.Unmarshal(grpJSON, &report.Grp)

		reports = append(reports, report)
	}

	c.JSON(http.StatusOK, reports)
}

// READ (один)
func GetReport(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var report Report
	var checkpointsJSON, grpJSON []byte

	err := db.DB.QueryRow(
		`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp FROM reports WHERE id=$1`, id,
	).Scan(&report.ID, &report.RouteName, &report.GpxFile, &checkpointsJSON, &report.MustContactBy, &report.Status, &grpJSON)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	json.Unmarshal(checkpointsJSON, &report.Checkpoints)
	json.Unmarshal(grpJSON, &report.Grp)

	c.JSON(http.StatusOK, report)
}

// DELETE
func DeleteReport(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	_, err := db.DB.Exec("DELETE FROM reports WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB delete error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UPDATE STATUS
func UpdateReportStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var payload struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	_, err := db.DB.Exec("UPDATE reports SET status=$1 WHERE id=$2", payload.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
