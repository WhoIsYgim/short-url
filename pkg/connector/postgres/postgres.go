package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"runtime"
	"short-link/config"
	"strconv"
	"time"
)

func GetPostgresConnector(cfg *config.DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		cfg.User,
		cfg.Database,
		cfg.Password,
		cfg.Host,
		strconv.FormatUint(uint64(cfg.Port), 10))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	connectTries := 0
	for ; connectTries < cfg.ReconnectRetries; connectTries++ {
		err = db.PingContext(ctx)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if connectTries == cfg.ReconnectRetries {
		return nil, errors.New("can not establish connection with database")
	}

	const OpenConnsFactor = 3
	var consCount = runtime.NumCPU()
	db.SetMaxOpenConns(consCount * OpenConnsFactor)
	db.SetMaxIdleConns(consCount)
	db.SetConnMaxIdleTime(10 * time.Second)
	db.SetConnMaxLifetime(0)

	return db, nil
}

func GetSqlxConnector(db *sql.DB, driverName string) *sqlx.DB {
	return sqlx.NewDb(db, driverName)
}
