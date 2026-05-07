package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func connectDB(databaseURL string) (*pgxpool.Pool, error) {
	// pgxpool = connection pool (หลาย goroutine ใช้ connection ร่วมกันได้)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// ทดสอบว่าเชื่อมได้จริงด้วย Ping
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
