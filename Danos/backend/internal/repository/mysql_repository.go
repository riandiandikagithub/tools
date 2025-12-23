package repository

import (
	"database/sql"
	"log"

	"github.com/Danos/backend/internal/domain"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) GetMetrics(name string, schema string) (domain.MySQLMetrics, error) {
	metric := domain.MySQLMetrics{
		Name: name,
	}

	// Check connection
	if err := r.db.Ping(); err != nil {
		metric.Status = "DOWN"
		return metric, err
	}
	metric.Status = "UP"

	// Get database size
	var sizeBytes int64
	err := r.db.QueryRow(`
		SELECT SUM(data_length + index_length)
		FROM information_schema.tables
		WHERE table_schema = DATABASE();
	`).Scan(&sizeBytes)
	if err != nil {
		log.Printf("Error getting MySQL size: %v", err)
	}
	metric.SizeBytes = sizeBytes

	// Count tables
	var tableCount int
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = DATABASE();
	`).Scan(&tableCount)
	if err != nil {
		log.Printf("Error getting MySQL table count: %v", err)
	}
	metric.TableCount = tableCount

	// Version
	err = r.db.QueryRow(`SELECT VERSION()`).Scan(&metric.Version)
	if err != nil {
		log.Printf("Error getting MySQL version: %v", err)
	}

	// Uptime
	var uptime int64
	err = r.db.QueryRow(`SHOW GLOBAL STATUS LIKE 'Uptime'`).Scan(new(string), &uptime)
	if err != nil {
		log.Printf("Error getting MySQL uptime: %v", err)
	}
	metric.UptimeSeconds = uptime

	return metric, nil
}
