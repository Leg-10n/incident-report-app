package handlers

import (
	"encoding/json"
	"incident-report-app/database"
	"incident-report-app/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if msg := req.Validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	// Check username availability before hashing (saves CPU if taken)
	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, req.Username).Scan(&count)
	if count > 0 {
		writeError(w, http.StatusConflict, "Username already taken")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO users (username, password) VALUES (?, ?)`,
		req.Username, string(hash),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	id, _ := result.LastInsertId()
	token, err := generateToken(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"token": token,
		"user":  map[string]any{"id": id, "username": req.Username},
	})
}

// POST /api/auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if msg := req.Validate(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, username, password FROM users WHERE username = ?`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		// Same message for "not found" and "wrong password" — prevents username enumeration
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  map[string]any{"id": user.ID, "username": user.Username},
	})
}

// generateToken creates a signed JWT that expires in 7 days.
func generateToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}