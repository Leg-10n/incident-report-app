package models

import "strings"

// User represents a registered account stored in the DB.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // bcrypt hash
}

// AuthRequest is used for both /register and /login.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() string {
	r.Username = strings.TrimSpace(r.Username)
	if r.Username == "" {
		return "username is required"
	}
	if len(r.Username) < 3 {
		return "username must be at least 3 characters"
	}
	if len(r.Username) > 30 {
		return "username must be 30 characters or fewer"
	}
	if r.Password == "" {
		return "password is required"
	}
	if len(r.Password) < 6 {
		return "password must be at least 6 characters"
	}
	return ""
}