package cli

import (
	"database/sql"
	"fmt"
)

type AuditEntry struct {
	ID         int64
	EntityType string
	EntityID   int64
	Action     string
	Role       string
	Payload    string
	Timestamp  string
}

func AddAuditEntry(db *sql.DB, entityType string, entityID int64, action, role, payload string) error {
	query := `INSERT INTO audit_log (entity_type, entity_id, action, role, payload) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, entityType, entityID, action, role, payload)
	return err
}

// AddAuditEntryReturnID is like AddAuditEntry but returns the inserted row ID.
func AddAuditEntryReturnID(db *sql.DB, entityType string, entityID int64, action, role, payload string) (int64, error) {
	query := `INSERT INTO audit_log (entity_type, entity_id, action, role, payload) VALUES (?, ?, ?, ?, ?)`
	res, err := db.Exec(query, entityType, entityID, action, role, payload)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListAuditEntries(db *sql.DB, entityType string, entityID *int64) ([]AuditEntry, error) {
	query := `SELECT id, COALESCE(entity_type, ''), COALESCE(entity_id, 0), COALESCE(action, ''), COALESCE(role, ''), COALESCE(payload, ''), timestamp FROM audit_log WHERE 1=1`
	args := []interface{}{}

	if entityType != "" {
		query += ` AND entity_type = ?`
		args = append(args, entityType)
	}
	if entityID != nil {
		query += ` AND entity_id = ?`
		args = append(args, *entityID)
	}

	query += ` ORDER BY timestamp DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.Role, &e.Payload, &e.Timestamp); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// LogTaskAction is a convenience wrapper for task-related audit entries
func LogTaskAction(db *sql.DB, taskID int64, action, role, details string) {
	payload := fmt.Sprintf("[%s] %s", role, details)
	// Best-effort logging — don't fail the operation if audit logging fails
	_ = AddAuditEntry(db, "task", taskID, action, role, payload)
}
