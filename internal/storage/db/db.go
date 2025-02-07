package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"github.com/kiper0808/api/internal/storage/config"
)

func New(cfg config.Database) (*sqlx.DB, error) {
	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("time load location failed: %w", err)
	}
	conf := mysql.NewConfig()
	conf.Net = cfg.Net
	conf.Addr = cfg.Server
	conf.User = cfg.User
	conf.Passwd = cfg.Password
	conf.DBName = cfg.DBName
	conf.Timeout = cfg.Timeout
	conf.Loc = location
	conf.ParseTime = true

	dbConn, err := otelsqlx.Connect("mysql", conf.FormatDSN(),
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithDBName(cfg.DBName),
	)
	if err != nil {
		return nil, fmt.Errorf("db connect failed: %w", err)
	}

	dbConn.SetMaxOpenConns(cfg.MaxOpenConnections)
	dbConn.SetMaxIdleConns(cfg.MaxIdleConnections)
	dbConn.SetConnMaxLifetime(15 * time.Minute)

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping failed: %w", err)
	}

	return dbConn, nil
}
