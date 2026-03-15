package cli

import (
	"database/sql"
	"fmt"
)

type Milestone struct {
	ID          int64
	ProjectID   int64
	ParentID    *int64 // self-referencing FK for nesting
	ArchetypeID *int64 // Pointer to handle NULL
	Name        string
	Description string
	Status      string
}

func AddMilestone(db *sql.DB, projectID int64, parentID *int64, archetypeID *int64, name string, description string) (int64, error) {
	query := `INSERT INTO milestones (project_id, parent_id, archetype_id, name, description, status) VALUES (?, ?, ?, ?, ?, 'open')`
	res, err := db.Exec(query, projectID, parentID, archetypeID, name, description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetMilestone(db *sql.DB, id int64) (*Milestone, error) {
	var m Milestone
	query := `SELECT id, project_id, parent_id, archetype_id, name, COALESCE(description, ''), status FROM milestones WHERE id = ? AND deleted_at IS NULL`
	err := db.QueryRow(query, id).Scan(&m.ID, &m.ProjectID, &m.ParentID, &m.ArchetypeID, &m.Name, &m.Description, &m.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("milestone not found")
		}
		return nil, err
	}
	return &m, nil
}

func ListMilestones(db *sql.DB, projectID int64, role string) ([]Milestone, error) {
	query := `SELECT id, project_id, parent_id, archetype_id, name, COALESCE(description, ''), status FROM milestones WHERE project_id = ? AND deleted_at IS NULL`
	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []Milestone
	for rows.Next() {
		var m Milestone
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.ParentID, &m.ArchetypeID, &m.Name, &m.Description, &m.Status); err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, nil
}

func UpdateMilestoneStatus(db *sql.DB, id int64, newStatus string, role string) error {
	if role != "architect" && role != "senior-engineer" {
		return fmt.Errorf("only architect or senior-engineer can update milestone status")
	}

	if newStatus != "open" && newStatus != "completed" {
		return fmt.Errorf("invalid milestone status: %s (must be 'open' or 'completed')", newStatus)
	}

	_, err := db.Exec(`UPDATE milestones SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, newStatus, id)
	if err == nil {
		LogTaskAction(db, id, "milestone_status_change", role, fmt.Sprintf("Milestone status changed to '%s'", newStatus))
	}
	return err
}

func SoftDeleteMilestone(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE milestones SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}
