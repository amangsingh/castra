package cli

import (
	"database/sql"
)

type Sprint struct {
	ID        int64
	ProjectID int64
	Name      string
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
	Status    string
}

func AddSprint(db *sql.DB, projectID int64, name, start, end string) (int64, error) {
	query := `INSERT INTO sprints (project_id, name, start_date, end_date, status) VALUES (?, ?, ?, ?, 'planning')`
	res, err := db.Exec(query, projectID, name, start, end)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListSprints(db *sql.DB, projectID int64) ([]Sprint, error) {
	query := `SELECT id, project_id, name, start_date, end_date, status FROM sprints WHERE project_id = ? AND deleted_at IS NULL ORDER BY start_date DESC`
	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sprints []Sprint
	for rows.Next() {
		var s Sprint
		// Handle nullable dates if needed, scanning into sql.NullString
		var sd, ed sql.NullString
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &sd, &ed, &s.Status); err != nil {
			return nil, err
		}
		if sd.Valid {
			s.StartDate = sd.String
		}
		if ed.Valid {
			s.EndDate = ed.String
		}
		sprints = append(sprints, s)
	}
	return sprints, nil
}
