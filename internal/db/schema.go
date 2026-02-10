package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
	-- Projects: The top-level container
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'active', -- active, archived
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME -- Soft delete if NOT NULL
	);

	-- Sprints: Time-boxed iterations within a project
	CREATE TABLE IF NOT EXISTS sprints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		start_date DATE,
		end_date DATE,
		status TEXT NOT NULL DEFAULT 'planning', -- planning, active, completed
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME, -- Soft delete
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	-- Tasks: Work items
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		sprint_id INTEGER, -- NULL means Backlog
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'todo', -- todo, doing, review, pending, done, blocked
		priority TEXT DEFAULT 'medium',
		context_ref TEXT,
		qa_approved BOOLEAN DEFAULT FALSE,
		security_approved BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME, -- Soft delete
		FOREIGN KEY(project_id) REFERENCES projects(id),
		FOREIGN KEY(sprint_id) REFERENCES sprints(id)
	);

	-- Project Notes: Role-tagged context

	CREATE TABLE IF NOT EXISTS project_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		tags TEXT, -- Comma-separated tags e.g. "engineer,frontend"
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);

	-- Audit Log (Preserved for history)
	CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_type TEXT, -- project, sprint, task
		entity_id INTEGER,
		action TEXT NOT NULL,
		payload TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error creating schema: %v", err)
		return nil, err
	}

	return db, nil
}
