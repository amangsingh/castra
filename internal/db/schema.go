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

	if err := RunMigrations(db); err != nil {
		log.Printf("Error running migrations: %v", err)
		return nil, err
	}

	return db, nil
}
