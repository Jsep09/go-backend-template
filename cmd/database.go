package main

import (
	"context"
	"fmt"
	"log/slog"
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

	// Ping เพื่อเช็ค connection จริง
	// ถ้า fail แค่ warn — server ยังขึ้นได้ DB query จะ fail เองตอนใช้งาน
	if err := pool.Ping(ctx); err != nil {
		slog.Warn("database ping failed, continuing without DB", "error", err)
	}

	return pool, nil
}
