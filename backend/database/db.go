package database

import (
	"database/sql"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite", "./incidents.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	createUsersTable()
	createIncidentsTable()
	migrateIncidentsAddUserID()
	log.Println("Database initialized successfully")
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		username   TEXT NOT NULL UNIQUE,
		password   TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(query); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
}

func createIncidentsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS incidents (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		description TEXT NOT NULL,
		category    TEXT NOT NULL CHECK(category IN ('Safety', 'Maintenance')),
		status      TEXT NOT NULL CHECK(status IN ('Open', 'In Progress', 'Success')),
		user_id     INTEGER REFERENCES users(id),
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(query); err != nil {
		log.Fatalf("Failed to create incidents table: %v", err)
	}
}

func migrateIncidentsAddUserID() {
	_, err := DB.Exec(`ALTER TABLE incidents ADD COLUMN user_id INTEGER REFERENCES users(id)`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column") {
		log.Printf("Migration note: %v", err)
	}
}