package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Cave struct {
	ID        int
	Name      string
	State     string
	County    string
	Latitude  float64
	Longitude float64
	CreatedAt time.Time
}

type CaveModel struct {
	DB *sql.DB
}

func (c CaveModel) Insert(name, state, county string, latitude, longitude float64) error {
	stmt := "INSERT INTO caves (name, state, county, latitude, longitude) VALUES (?, ?, ?, ?, ?)"

	result, err := c.DB.Exec(stmt, name, state, county, latitude, longitude)
	if err != nil {
		return fmt.Errorf("Error inserting into database: %w", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error getting last insert id: %w", err)
	}

	return nil
}
