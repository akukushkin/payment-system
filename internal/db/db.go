package db

import (
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func New() *sqlx.DB {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		log.Fatalf("failed to init db: %s", err)
	}

	db := stdlib.OpenDB(config)
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err)
	}

	return sqlx.NewDb(db, "pgx")
}
