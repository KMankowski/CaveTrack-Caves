package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/KMankowski/CaveTrack-Caves/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func parseAndImportCaves(alabamaXML io.Reader, dsn string) error {
	db, err := openDB(dsn)
	if err != nil {
		return fmt.Errorf("Unable to open database: %w", err)
	}
	defer db.Close()

	// Existing DB must first be backed up elsewhere.
	isCavesTableEmpty, err := checkIfCavesTableIsEmpty(db)
	if err != nil {
		return fmt.Errorf("Error checking if caves table is empty: %w", err)
	} else if !isCavesTableEmpty {
		return errors.New("Caves table must be empty.")
	}

	caves, err := parseCaves(alabamaXML)
	if err != nil {
		return fmt.Errorf("Error parsing XML: %w", err)
	}

	err = importCaves(caves, db)
	if err != nil {
		return fmt.Errorf("Error importing into database: %w", err)
	}

	return nil
}

func parseCaves(alabamaXML io.Reader) ([]models.Cave, error) {
	return parseAlabamaCaves(alabamaXML)
}

func importCaves(caves []models.Cave, db *sql.DB) error {
	cavesDB := models.CaveModel{DB: db}

	for _, cave := range caves {
		err := cavesDB.Insert(cave.Name, cave.State, cave.County, cave.Latitude, cave.Longitude)
		if err != nil {
			return fmt.Errorf("Error inserting cave with state %s and name %s: %w", cave.State, cave.Name, err)
		}
	}

	return nil
}

func checkIfCavesTableIsEmpty(db *sql.DB) (bool, error) {
	var table string
	err := db.QueryRow("SHOW TABLES LIKE 'caves';").Scan(&table)
	if err == sql.ErrNoRows {
		return false, errors.New("caves table does not exist")
	} else if err != nil {
		return false, err
	}

	var numRows int
	err = db.QueryRow("SELECT COUNT(*) FROM caves").Scan(&numRows)
	if err != nil {
		return false, err
	}
	return numRows == 0, nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
