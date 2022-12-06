package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	DBPath string
	db     *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", dbPath, err)
	}
	CreateTable := `CREATE TABLE IF NOT EXISTS deeds (
		id TEXT PRIMARY KEY,
		name TEXT,
		period DATETIME,
		last DATETIME);`
	if _, err := db.Exec(CreateTable); err != nil {
		return nil, fmt.Errorf("%s: %w", dbPath, err)
	}
	return &DB{
		dbPath, db,
	}, nil
}

func (db *DB) AddDeed(d *Deed) error {
	stmt := "INSERT OR REPLACE INTO deeds (id, name, period, last) " +
		"VALUES ($1, $2, $3, $4)"
	_, err := db.db.Exec(stmt, d.ID, d.Name, d.Period, d.Last)
	return db.Error(err)
}

func (db *DB) Update(id string) error {
	stmt := "UPDATE deeds  SET last=$1 WHERE id=$2"
	_, err := db.db.Exec(stmt, time.Now(), id)
	return db.Error(err)
}

func (db *DB) Delete(id string) error {
	stmt := "DELETE FROM deeds WHERE id=$1"
	_, err := db.db.Exec(stmt, id)
	return db.Error(err)
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Iterate(callback func(*Deed)) error {
	stmt := "SELECT id, name, period, last FROM deeds ORDER BY name"
	rows, err := db.db.Query(stmt)
	if err != nil {
		return db.Error(err)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return db.Error(rows.Err())
	}
	for rows.Next() {
		var deed Deed
		err = rows.Scan(&deed.ID, &deed.Name, &deed.Period, &deed.Last)
		if err != nil {
			return db.Error(err)
		}
		callback(&deed)
	}
	return nil
}

func (db *DB) Error(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", db.DBPath, err)
}
