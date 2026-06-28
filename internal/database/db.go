package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDBConnection(DBAddr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBAddr)
	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return db, err
	}

	log.Printf("Database Connection Established",)

	return db, nil
}