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
// я срал посчему деплойка сломаласб
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

	// парсим JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// если статус пустой → ставим active
	if report.Status == "" {
		report.Status = "active"
	}

	// маршалим grp и checkpoints
	grpBytes, err := json.Marshal(report.Grp)
	if err != nil {
		http.Error(w, "Error encoding grp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	checkpointsBytes, err := json.Marshal(report.Checkpoints)
	if err != nil {
		http.Error(w, "Error encoding checkpoints: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// вставляем запись
	err = db.DB.QueryRow(
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
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}


// READ
func GetReport(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`
		SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp
		FROM reports
		WHERE id=$1
	`, id)

	var report Report
	var checkpointsJSON, grpJSON []byte

	err = row.Scan(
		&report.ID,
		&report.RouteName,
		&report.GpxFile,
		&checkpointsJSON,
		&report.MustContactBy,
		&report.Status,
		&grpJSON,
	)
	if err != nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	// декодируем JSON-поля
	if err := json.Unmarshal(checkpointsJSON, &report.Checkpoints); err != nil {
		http.Error(w, "Error decoding checkpoints: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(grpJSON, &report.Grp); err != nil {
		http.Error(w, "Error decoding grp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}


func GetReport(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`
		SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp
		FROM reports
		WHERE id=$1
	`, id)

	var report Report
	var checkpointsJSON, grpJSON []byte

	err = row.Scan(
		&report.ID,
		&report.RouteName,
		&report.GpxFile,
		&checkpointsJSON,
		&report.MustContactBy,
		&report.Status,
		&grpJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.Unmarshal(checkpointsJSON, &report.Checkpoints); err != nil {
		http.Error(w, "Error decoding checkpoints: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(grpJSON, &report.Grp); err != nil {
		http.Error(w, "Error decoding grp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
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
// UPDATE STATUS (PATCH)
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

	// обновляем статус
	_, err = db.DB.Exec("UPDATE reports SET status=$1 WHERE id=$2", payload.Status, id)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// получаем обновлённый отчёт
	var report Report
	var checkpointsBytes []byte
	var grpBytes []byte

	err = db.DB.QueryRow(
		`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp 
		 FROM reports WHERE id=$1`, id,
	).Scan(
		&report.ID,
		&report.RouteName,
		&report.GpxFile,
		&checkpointsBytes,
		&report.MustContactBy,
		&report.Status,
		&grpBytes,
	)
	if err != nil {
		http.Error(w, "DB fetch error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.Unmarshal(checkpointsBytes, &report.Checkpoints)
	json.Unmarshal(grpBytes, &report.Grp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
