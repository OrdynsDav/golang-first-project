package database

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func Connect(databaseUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseUrl)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)

	return db, nil
}