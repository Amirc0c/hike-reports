package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"backend/db"
)

// Report — модель для JSON
type Group struct {
    Number   string `json:"number"`
    Name     string `json:"name"`
    Telegram string `json:"telegram"`
}

type Checkpoint struct {
    Name string `json:"name"`
    Time string `json:"time"`
}

type Report struct {
    ID            int64        `json:"id"`
    RouteName     string       `json:"route_name"`
    GpxFile       string       `json:"gpx_file"`
    Checkpoints   []Checkpoint `json:"checkpoints"`
    MustContactBy string       `json:"must_contact_by"`
    Status        string       `json:"status"`
    Grp           Group        `json:"grp"`
}


// CREATE
func CreateReport(w http.ResponseWriter, r *http.Request) {
	var report Report

	// парсим JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// маршалим grp и checkpoints в JSON для сохранения в Postgres
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

	// ответ в JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// READ (all)
func GetReports(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, route_name, gpx_file, checkpoints, must_contact_by, status, grp FROM reports`)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []models.Report

	for rows.Next() {
		var report models.Report
		var checkpointsBytes []byte
		var grpBytes []byte

		err := rows.Scan(
			&report.ID,
			&report.RouteName,
			&report.GpxFile,
			&checkpointsBytes, // сначала []byte
			&report.MustContactBy,
			&report.Status,
			&grpBytes,
		)
		if err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// распарсим JSON в Go-структуру
		if err := json.Unmarshal(checkpointsBytes, &report.Checkpoints); err != nil {
			http.Error(w, "JSON unmarshal checkpoints error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(grpBytes, &report.Grp); err != nil {
			http.Error(w, "JSON unmarshal grp error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// UPDATE
func UpdateReportStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var input struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    _, err = db.DB.Exec("UPDATE reports SET status = $1 WHERE id = $2", input.Status, id)
    if err != nil {
        http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Status updated successfully"))
}


// DELETE
func DeleteReport(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`DELETE FROM reports WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

