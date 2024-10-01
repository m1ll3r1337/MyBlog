package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "balls",
		Database: "postgres",
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host %s port %s user %s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.SSLMode)
}
