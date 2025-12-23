// ==================== internal/infrastructure/mysql/client.go ====================
package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Danos/backend/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLManager struct {
	mu      sync.RWMutex
	clients map[string]*sql.DB
	config  *domain.MySQLConfig
}

func NewMySQLManager() *MySQLManager {
	return &MySQLManager{
		clients: make(map[string]*sql.DB),
	}
}

func (m *MySQLManager) Initialize(config *domain.MySQLConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config

	for _, dbConfig := range config.Databases {
		if err := m.connectDatabase(dbConfig); err != nil {
			log.Printf("Failed to connect to MySQL %s: %v", dbConfig.Name, err)
			continue
		}
		log.Printf("Connected to MySQL: %s", dbConfig.Name)
	}

	return nil
}

func (m *MySQLManager) connectDatabase(config domain.MySQLDatabase) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port,
		config.Database, config.Charset,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.Pool.MaxConnections)
	db.SetMaxIdleConns(config.Pool.MinConnections)
	db.SetConnMaxIdleTime(time.Duration(config.Pool.MaxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(config.Pool.ConnectionTimeout) * time.Second)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	m.clients[config.Name] = db
	return nil
}

func (m *MySQLManager) Reconnect(config *domain.MySQLConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Println("Reconnecting MySQL clients...")

	// Close existing connections
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing MySQL %s: %v", name, err)
		}
	}
	m.clients = make(map[string]*sql.DB)

	// Reconnect
	m.config = config
	for _, dbConfig := range config.Databases {
		if err := m.connectDatabase(dbConfig); err != nil {
			log.Printf("Failed to reconnect MySQL %s: %v", dbConfig.Name, err)
			continue
		}
		log.Printf("Reconnected to MySQL: %s", dbConfig.Name)
	}

	return nil
}

func (m *MySQLManager) GetClient(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("mysql client %s not found", name)
	}
	return client, nil
}

func (m *MySQLManager) GetAllClients() map[string]*sql.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*sql.DB)
	for k, v := range m.clients {
		result[k] = v
	}
	return result
}

func (m *MySQLManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing MySQL %s: %v", name, err)
		}
	}
	m.clients = make(map[string]*sql.DB)
	return nil
}

func (m *MySQLManager) GetConnectionStatus() map[string]interface{} {
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
