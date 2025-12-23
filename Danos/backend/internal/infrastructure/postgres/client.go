// ==================== internal/infrastructure/postgres/client.go ====================
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Danos/backend/internal/domain"

	_ "github.com/lib/pq"
)

type PostgresManager struct {
	mu      sync.RWMutex
	clients map[string]*sql.DB
	config  *domain.PostgreSQLConfig
}

func NewPostgresManager() *PostgresManager {
	return &PostgresManager{
		clients: make(map[string]*sql.DB),
	}
}

func (m *PostgresManager) Initialize(config *domain.PostgreSQLConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config

	for _, dbConfig := range config.Databases {
		if err := m.connectDatabase(dbConfig); err != nil {
			log.Printf("Failed to connect to PostgreSQL %s: %v", dbConfig.Name, err)
			continue
		}
		log.Printf("Connected to PostgreSQL: %s", dbConfig.Name)
	}

	return nil
}

func (m *PostgresManager) connectDatabase(config domain.PostgreSQLDatabase) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password,
		config.Database, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.Pool.MaxConnections)
	db.SetMaxIdleConns(config.Pool.MinConnections)
	db.SetConnMaxIdleTime(time.Duration(config.Pool.MaxIdleTime) * time.Second)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	m.clients[config.Name] = db
	return nil
}

func (m *PostgresManager) Reconnect(config *domain.PostgreSQLConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Println("Reconnecting PostgreSQL clients...")

	// Close existing connections
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing PostgreSQL %s: %v", name, err)
		}
	}
	m.clients = make(map[string]*sql.DB)

	// Reconnect
	m.config = config
	for _, dbConfig := range config.Databases {
		if err := m.connectDatabase(dbConfig); err != nil {
			log.Printf("Failed to reconnect PostgreSQL %s: %v", dbConfig.Name, err)
			continue
		}
		log.Printf("Reconnected to PostgreSQL: %s", dbConfig.Name)
	}

	return nil
}

func (m *PostgresManager) GetClient(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("postgres client %s not found", name)
	}
	return client, nil
}

func (m *PostgresManager) GetAllClients() map[string]*sql.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*sql.DB)
	for k, v := range m.clients {
		result[k] = v
	}
	return result
}

func (m *PostgresManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing PostgreSQL %s: %v", name, err)
		}
	}
	m.clients = make(map[string]*sql.DB)
	return nil
}
func (m *PostgresManager) GetConnectionStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	databases := make([]map[string]interface{}, 0)

	for name, db := range m.clients {
		status := "disconnected"
		if err := db.Ping(); err == nil {
			status = "connected"
		}

		databases = append(databases, map[string]interface{}{
			"name":   name,
			"status": status,
		})
	}

	return map[string]interface{}{
		"status":    getOverallStatus(databases),
		"databases": databases,
	}
}

func getOverallStatus(databases []map[string]interface{}) string {
	if len(databases) == 0 {
		return "disconnected"
	}

	allConnected := true
	for _, db := range databases {
		if db["status"] != "connected" {
			allConnected = false
			break
		}
	}

	if allConnected {
		return "connected"
	}
	return "partial"
}
