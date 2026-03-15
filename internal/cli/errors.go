package cli

import "database/sql"

// Severity levels for error logging.
const (
	SeverityWarn  = "warn"
	SeverityError = "error"
	SeverityFatal = "fatal"
)

// LogError records a non-fatal system error to the castra_errors table.
// It is best-effort — failures to write the error log are silently ignored
// so that the original error handling path is never disrupted.
func LogError(db *sql.DB, role, command, message, severity string) {
	if db == nil {
		return
	}
	_, _ = db.Exec(
		`INSERT INTO castra_errors (role, command, message, severity) VALUES (?, ?, ?, ?)`,
		role, command, message, severity,
	)
}

// ListErrors returns all recorded errors, newest first.
func ListErrors(db *sql.DB) ([]ErrorEntry, error) {
	rows, err := db.Query(`SELECT id, role, command, message, severity, created_at FROM castra_errors ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ErrorEntry
	for rows.Next() {
		var e ErrorEntry
		if err := rows.Scan(&e.ID, &e.Role, &e.Command, &e.Message, &e.Severity, &e.Timestamp); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ErrorEntry represents a row in castra_errors.
type ErrorEntry struct {
	ID        int64
	Role      string
	Command   string
	Message   string
	Severity  string
	Timestamp string
}
