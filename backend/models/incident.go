package models

import (
	"strings"
	"time"
)

type Category string
type Status string

const (
	CategorySafety      Category = "Safety"
	CategoryMaintenance Category = "Maintenance"

	StatusOpen       Status = "Open"
	StatusInProgress Status = "In Progress"
	StatusSuccess    Status = "Success"
)

type Incident struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      Category  `json:"category"`
	Status        Status    `json:"status"`
	UserID        int64     `json:"user_id"`        // who created it
	OwnerUsername string    `json:"owner_username"` // joined from users table
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type IncidentRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
	Status      Status   `json:"status"`
}

func (r *IncidentRequest) Validate() string {
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)

	if r.Title == "" {
		return "title is required"
	}
	if len(r.Title) > 150 {
		return "title must be 150 characters or fewer"
	}
	if r.Description == "" {
		return "description is required"
	}
	if len(r.Description) > 2000 {
		return "description must be 2000 characters or fewer"
	}
	if r.Category != CategorySafety && r.Category != CategoryMaintenance {
		return "category must be 'Safety' or 'Maintenance'"
	}
	if r.Status != StatusOpen && r.Status != StatusInProgress && r.Status != StatusSuccess {
		return "status must be 'Open', 'In Progress', or 'Success'"
	}
	return ""
}