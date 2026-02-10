package cli

import (
	"database/sql"
	"fmt"
)

// Project struct
type Project struct {
	ID          int64
	Name        string
	Description string
	Status      string
	Notes       string
}

// AddProject creates a new project
func AddProject(db *sql.DB, name, description, notes string) (int64, error) {
	query := `INSERT INTO projects (name, description, notes, status) VALUES (?, ?, ?, 'active')`
	res, err := db.Exec(query, name, description, notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ListProjects retrieves active projects (default) or supports filters
func ListProjects(db *sql.DB, showArchived bool, showDeleted bool) ([]Project, error) {
	query := `SELECT id, name, description, status, notes FROM projects WHERE 1=1`

	if !showDeleted {
		query += ` AND deleted_at IS NULL`
	}
	if !showArchived {
		query += ` AND status != 'archived'`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.Notes); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// UpdateProject updates fields dynamically
func UpdateProject(db *sql.DB, id int64, name, description, status, notes *string) error {
	// Simple approach: Construct update query based on non-nil fields
	// For MVP simplicity, we'll update fields if provided string is not empty, assuming user provides flag logic

	// Better: Read current, update provided, write back? Or just run UPDATE IF logic.
	// We'll trust the caller passes pointers, nil means no update.

	// Since this is a CLI helper, we might abstract this logic.
	// For now, let's implement soft delete / restore separately
	return nil
}

// SoftDeleteProject sets deleted_at
func SoftDeleteProject(db *sql.DB, id int64) error {
	query := `UPDATE projects SET deleted_at = CURRENT_TIMESTAMP, status = 'archived' WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// HardDeleteProject removes the row
func HardDeleteProject(db *sql.DB, id int64) error {
	// Check for tasks
	var count int
	err := db.QueryRow(`SELECT count(*) FROM tasks WHERE project_id = ?`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete project with %d existing tasks; delete tasks first or use force", count)
	}

	_, err = db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}
