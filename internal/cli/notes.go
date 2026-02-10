package cli

import (
	"database/sql"
	"strings"
)

type Note struct {
	ID        int64
	ProjectID int64
	Content   string
	Tags      string
}

func AddNote(db *sql.DB, projectID int64, content, tags string) (int64, error) {
	query := `INSERT INTO project_notes (project_id, content, tags) VALUES (?, ?, ?)`
	res, err := db.Exec(query, projectID, content, tags)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListNotes(db *sql.DB, projectID int64, role string) ([]Note, error) {
	// Architect and Doc Writer see all
	if role == "architect" || role == "doc-writer" {
		query := `SELECT id, project_id, content, tags FROM project_notes WHERE project_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`
		return queryNotes(db, query, projectID)
	}

	// Others see if tags contain their role OR "all" (convention)
	// SQLite LIKE is case-insensitive by default for ASCII
	// We handle filtering via query or code. Code is safer for comma-separated checks.
	query := `SELECT id, project_id, content, tags FROM project_notes WHERE project_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`
	allNotes, err := queryNotes(db, query, projectID)
	if err != nil {
		return nil, err
	}

	var filtered []Note
	for _, n := range allNotes {
		if containsRole(n.Tags, role) {
			filtered = append(filtered, n)
		}
	}
	return filtered, nil
}

func queryNotes(db *sql.DB, query string, args ...interface{}) ([]Note, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.ProjectID, &n.Content, &n.Tags); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func containsRole(tags, role string) bool {
	// Simple tag check: "engineer,frontend" contains "engineer"
	parts := strings.Split(tags, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == role {
			return true
		}
	}
	return false
}
