package repository

import (
	"database/sql"
	"log"

	"github.com/Danos/backend/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetMetrics(name string, schema string) (domain.PostgreSQLMetrics, error) {
	metric := domain.PostgreSQLMetrics{
		Name: name,
	}

	// Check connection
	if err := r.db.Ping(); err != nil {
		metric.Status = "DOWN"
		return metric, err
	}
	metric.Status = "UP"

	// Version
	err := r.db.QueryRow(`SELECT version()`).Scan(&metric.Version)
	if err != nil {
		log.Printf("Error getting PostgreSQL version: %v", err)
	}

	// Database size
	err = r.db.QueryRow(`SELECT pg_database_size(current_database())`).Scan(&metric.SizeBytes)
	if err != nil {
		log.Printf("Error getting PostgreSQL size: %v", err)
	}

	// Table count
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = 'public';
	`).Scan(&metric.TableCount)
	if err != nil {
		log.Printf("Error getting PostgreSQL table count: %v", err)
	}

	// Uptime
	err = r.db.QueryRow(`
		SELECT EXTRACT(EPOCH FROM (NOW() - pg_postmaster_start_time()))
	`).Scan(&metric.UptimeSeconds)
	if err != nil {
		log.Printf("Error getting PostgreSQL uptime: %v", err)
	}

	// Active connections
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM pg_stat_activity
		WHERE datname = current_database()
	`).Scan(&metric.ActiveConnections)
	if err != nil {
		log.Printf("Error getting PostgreSQL active connections: %v", err)
	}

	return metric, nil
}
