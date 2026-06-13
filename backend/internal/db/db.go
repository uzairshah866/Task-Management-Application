package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/user/taskapp/internal/log"
)

var DB *sql.DB

func Init() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	maxOpenConns := 25
	if val, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS")); err == nil && val > 0 {
		maxOpenConns = val
	}

	maxIdleConns := 5
	if val, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS")); err == nil && val > 0 {
		maxIdleConns = val
	}

	DB.SetMaxOpenConns(maxOpenConns)
	DB.SetMaxIdleConns(maxIdleConns)

	log.Info("Database connected (maxOpenConns=%d, maxIdleConns=%d)", maxOpenConns, maxIdleConns)
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
