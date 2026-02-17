package handlers

import (
	"encoding/json"
	"incident-report-app/database"
	"incident-report-app/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /api/incidents — all incidents, newest first, with owner username joined
func GetIncidents(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT i.id, i.title, i.description, i.category, i.status,
		       COALESCE(i.user_id, 0),
		       COALESCE(u.username, 'Unknown'),
		       i.created_at, i.updated_at
		FROM incidents i
		LEFT JOIN users u ON i.user_id = u.id
		ORDER BY i.created_at DESC
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
			&inc.Category, &inc.Status,
			&inc.UserID, &inc.OwnerUsername,
			&inc.CreatedAt, &inc.UpdatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to read incidents")
			return
		}
		incidents = append(incidents, inc)
	}
	writeJSON(w, http.StatusOK, incidents)
}

// POST /api/incidents — create incident, owned by the authenticated user
func CreateIncident(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey).(int64)

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
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
		INSERT INTO incidents (title, description, category, status, user_id)
		VALUES (?, ?, ?, ?, ?)
	`, req.Title, req.Description, req.Category, req.Status, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create incident")
		return
	}

	id, _ := result.LastInsertId()

	var inc models.Incident
	if err := database.DB.QueryRow(`
		SELECT i.id, i.title, i.description, i.category, i.status,
		       COALESCE(i.user_id, 0), COALESCE(u.username, 'Unknown'),
		       i.created_at, i.updated_at
		FROM incidents i
		LEFT JOIN users u ON i.user_id = u.id
		WHERE i.id = ?
	`, id).Scan(
		&inc.ID, &inc.Title, &inc.Description,
		&inc.Category, &inc.Status,
		&inc.UserID, &inc.OwnerUsername,
		&inc.CreatedAt, &inc.UpdatedAt,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve created incident")
		return
	}
	writeJSON(w, http.StatusCreated, inc)
}

// PUT /api/incidents/{id} — only the owner can update
func UpdateIncident(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	// Separate 404 from 403 so the client knows the difference
	var ownerID int64
	if err := database.DB.QueryRow(
		`SELECT COALESCE(user_id, 0) FROM incidents WHERE id = ?`, id,
	).Scan(&ownerID); err != nil {
		writeError(w, http.StatusNotFound, "Incident not found")
		return
	}
	if ownerID != userID {
		writeError(w, http.StatusForbidden, "You can only edit your own incidents")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req models.IncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body (max 1MB)")
		return
	}
	if msg := req.Validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	if _, err := database.DB.Exec(`
		UPDATE incidents
		SET title = ?, description = ?, category = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, req.Title, req.Description, req.Category, req.Status, id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update incident")
		return
	}

	var inc models.Incident
	if err := database.DB.QueryRow(`
		SELECT i.id, i.title, i.description, i.category, i.status,
		       COALESCE(i.user_id, 0), COALESCE(u.username, 'Unknown'),
		       i.created_at, i.updated_at
		FROM incidents i
		LEFT JOIN users u ON i.user_id = u.id
		WHERE i.id = ?
	`, id).Scan(
		&inc.ID, &inc.Title, &inc.Description,
		&inc.Category, &inc.Status,
		&inc.UserID, &inc.OwnerUsername,
		&inc.CreatedAt, &inc.UpdatedAt,
	); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve updated incident")
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

// DELETE /api/incidents/{id} — only the owner can delete
func DeleteIncident(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid incident ID")
		return
	}

	var ownerID int64
	if err := database.DB.QueryRow(
		`SELECT COALESCE(user_id, 0) FROM incidents WHERE id = ?`, id,
	).Scan(&ownerID); err != nil {
		writeError(w, http.StatusNotFound, "Incident not found")
		return
	}
	if ownerID != userID {
		writeError(w, http.StatusForbidden, "You can only delete your own incidents")
		return
	}

	if _, err := database.DB.Exec(`DELETE FROM incidents WHERE id = ?`, id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete incident")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Incident deleted successfully"})
}