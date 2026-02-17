package database

import (
	"database/sql"
	"log"

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

	createTable()
	log.Println("Database initialized successfully")
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS incidents (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		description TEXT NOT NULL,
		category    TEXT NOT NULL CHECK(category IN ('Safety', 'Maintenance')),
		status      TEXT NOT NULL CHECK(status IN ('Open', 'In Progress', 'Success')),
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(query); err != nil {
		log.Fatalf("Failed to create incidents table: %v", err)
	}
}