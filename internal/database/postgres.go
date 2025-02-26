package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // pq adapter for sql

	"github.com/multitheftauto/community/internal/config"
)

// NewPostgres connects to the database and returns a query generator
func NewPostgres(cfg config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
