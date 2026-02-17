package handlers

import (
	"encoding/json"
	"incident-report-app/database"
	"incident-report-app/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// helper: write error response
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /api/incidents — fetch all incidents, newest first
func GetIncidents(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, title, description, category, status, created_at, updated_at
		FROM incidents
		ORDER BY created_at DESC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch incidents")
		return
	}
	defer rows.Close()

	incidents := []models.Incident{}
	for rows.Next() {
		var inc models.Incident
		if err := rows.Scan(
			&inc.ID, &inc.Title, &inc.Description,
			&inc.Category, &inc.Status, &inc.CreatedAt, &inc.UpdatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan incident")
			return
		}
		incidents = append(incidents, inc)
	}

	writeJSON(w, http.StatusOK, incidents)
}

// POST /api/incidents — create a new incident
func CreateIncident(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	var req models.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body (max 1MB)")
		return
	}

	if msg := req.Validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO incidents (title, description, category, status)
		VALUES (?, ?, ?, ?)
	`, req.Title, req.Description, req.Category, req.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create incident")
		return
	}

	id, _ := result.LastInsertId()

	var inc models.Incident
	err = database.DB.QueryRow(`
		SELECT id, title, description, category, status, created_at, updated_at
		FROM incidents WHERE id = ?
	`, id).Scan(
		&inc.ID, &inc.Title, &inc.Description,
		&inc.Category, &inc.Status, &inc.CreatedAt, &inc.UpdatedAt,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve created incident")
		return
	}

	writeJSON(w, http.StatusCreated, inc)
}

// PUT /api/incidents/{id} — update an existing incident
func UpdateIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	var req models.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body (max 1MB)")
		return
	}

	if msg := req.Validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	result, err := database.DB.Exec(`
		UPDATE incidents
		SET title = ?, description = ?, category = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.Title, req.Description, req.Category, req.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update incident")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "Incident not found")
		return
	}

	var inc models.Incident
	if err := database.DB.QueryRow(`
		SELECT id, title, description, category, status, created_at, updated_at
		FROM incidents WHERE id = ?
	`, id).Scan(
		&inc.ID, &inc.Title, &inc.Description,
		&inc.Category, &inc.Status, &inc.CreatedAt, &inc.UpdatedAt,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve updated incident")
		return
	}

	writeJSON(w, http.StatusOK, inc)
}

// DELETE /api/incidents/{id} — delete an incident
func DeleteIncident(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	result, err := database.DB.Exec(`DELETE FROM incidents WHERE id = ?`, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete incident")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "Incident not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Incident deleted successfully"})
}