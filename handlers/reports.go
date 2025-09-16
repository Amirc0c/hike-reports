package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/db"
	"github.com/gorilla/mux"
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
func CreateReport(w http.ResponseWriter, r *http.Request) {
	var report Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkpointsJSON, _ := json.Marshal(report.Checkpoints)
	grpJSON, _ := json.Marshal(report.Grp)

	err := db.DB.QueryRow(
		`INSERT INTO reports (route_name, gpx_file, checkpoints, must_contact_by, status, grp)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		report.RouteName, report.GpxFile, checkpointsJSON, report.MustContactBy, report.Status, grpJSON,
	).Scan(&report.ID)

	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// READ
func GetReports(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp FROM reports`)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []Report

	for rows.Next() {
		var report Report
		var checkpointsBytes []byte
		var grpBytes []byte

		if err := rows.Scan(
			&report.ID,
			&report.RouteName,
			&report.GpxFile,
			&checkpointsBytes,
			&report.MustContactBy,
			&report.Status,
			&grpBytes,
		); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.Unmarshal(checkpointsBytes, &report.Checkpoints)
		json.Unmarshal(grpBytes, &report.Grp)

		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// DELETE
func DeleteReport(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Удаляем отчёт
	_, err = db.DB.Exec("DELETE FROM reports WHERE id=$1", id)
	if err != nil {
		http.Error(w, "DB delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем оставшиеся отчёты
	rows, err := db.DB.Query(`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp FROM reports`)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var report Report
		var checkpointsBytes []byte
		var grpBytes []byte

		if err := rows.Scan(
			&report.ID,
			&report.RouteName,
			&report.GpxFile,
			&checkpointsBytes,
			&report.MustContactBy,
			&report.Status,
			&grpBytes,
		); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.Unmarshal(checkpointsBytes, &report.Checkpoints)
		json.Unmarshal(grpBytes, &report.Grp)

		reports = append(reports, report)
	}

	// Возвращаем оставшиеся отчёты
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}


// UPDATE STATUS
func UpdateReportStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE reports SET status=$1 WHERE id=$2", payload.Status, id)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"status": payload.Status,
	})
}
