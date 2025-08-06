package intiator

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func InitDB(url string, log *zap.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	idleConnTimeoutString := os.Getenv("IDLE_CONN_TIMEOUT")
	idleConnTimeout, err := time.ParseDuration(idleConnTimeoutString)
	if err != nil || idleConnTimeout == 0 {
		log.Warn("invalid or missing IDLE_CONN_TIMEOUT", zap.String("value", idleConnTimeoutString), zap.Error(err))
		idleConnTimeout = 5 * time.Minute
	}
	config.MaxConnIdleTime = idleConnTimeout

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := conn.QueryRow(context.Background(), "SELECT 1").Scan(new(int)); err != nil {
		log.Fatal(fmt.Sprintf("Failed to ping database: %v", err))
	}

	return conn
}
