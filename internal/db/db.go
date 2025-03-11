package db

import (
	"context"
	"database/sql"
	"time"
)

func New(conn_string string, maxOpenConns int, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn_string)
	
	if err != nil {
		return nil, err

	}
	db.SetConnMaxIdleTime(maxIdleTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
