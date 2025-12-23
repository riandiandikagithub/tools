// ==================== internal/infrastructure/config/loader.go ====================
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Danos/backend/internal/domain"
	"gopkg.in/yaml.v3"
)

type ConfigLoader struct {
	mu              sync.RWMutex
	configPath      string
	redisConfig     *domain.RedisConfig
	kafkaConfig     *domain.KafkaConfig
	postgresConfig  *domain.PostgreSQLConfig
	mysqlConfig     *domain.MySQLConfig
	changeCallbacks []func()
}

func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath:      configPath,
		changeCallbacks: make([]func(), 0),
	}
}

// ensureConfigDir creates the config directory if it doesn't exist
func (c *ConfigLoader) ensureConfigDir() error {
	if err := os.MkdirAll(c.configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

// LoadAll loads all configuration files, creates them if they don't exist
func (c *ConfigLoader) LoadAll() error {
	// Ensure config directory exists
	if err := c.ensureConfigDir(); err != nil {
		return err
	}

	if err := c.LoadRedis(); err != nil {
		return fmt.Errorf("failed to load redis config: %w", err)
	}
	if err := c.LoadKafka(); err != nil {
		return fmt.Errorf("failed to load kafka config: %w", err)
	}
	if err := c.LoadPostgreSQL(); err != nil {
		return fmt.Errorf("failed to load postgresql config: %w", err)
	}
	if err := c.LoadMySQL(); err != nil {
		return fmt.Errorf("failed to load mysql config: %w", err)
	}
	return nil
}

// LoadRedis loads Redis configuration, creates default if not exists
func (c *ConfigLoader) LoadRedis() error {
	path := filepath.Join(c.configPath, "redis.yaml")

	// Check if file exists, if not create default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := c.createDefaultRedisConfig(path); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var wrapper struct {
		Redis domain.RedisConfig `yaml:"redis"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	c.mu.Lock()
	c.redisConfig = &wrapper.Redis
	c.mu.Unlock()

	return nil
}

// LoadKafka loads Kafka configuration, creates default if not exists
func (c *ConfigLoader) LoadKafka() error {
	path := filepath.Join(c.configPath, "kafka.yaml")

	// Check if file exists, if not create default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := c.createDefaultKafkaConfig(path); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var wrapper struct {
		Kafka domain.KafkaConfig `yaml:"kafka"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	c.mu.Lock()
	c.kafkaConfig = &wrapper.Kafka
	c.mu.Unlock()

	return nil
}

// LoadPostgreSQL loads PostgreSQL configuration, creates default if not exists
func (c *ConfigLoader) LoadPostgreSQL() error {
	path := filepath.Join(c.configPath, "postgresql.yaml")

	// Check if file exists, if not create default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := c.createDefaultPostgreSQLConfig(path); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var wrapper struct {
		PostgreSQL domain.PostgreSQLConfig `yaml:"postgresql"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	c.mu.Lock()
	c.postgresConfig = &wrapper.PostgreSQL
	c.mu.Unlock()

	return nil
}

// LoadMySQL loads MySQL configuration, creates default if not exists
func (c *ConfigLoader) LoadMySQL() error {
	path := filepath.Join(c.configPath, "mysql.yaml")

	// Check if file exists, if not create default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := c.createDefaultMySQLConfig(path); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var wrapper struct {
		MySQL domain.MySQLConfig `yaml:"mysql"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	c.mu.Lock()
	c.mysqlConfig = &wrapper.MySQL
	c.mu.Unlock()

	return nil
}

// ==================== Default Config Creators ====================

func (c *ConfigLoader) createDefaultRedisConfig(path string) error {
	defaultConfig := `# Redis Configuration
redis:
  # Single mode configuration
  single:
    host: "localhost"
    port: 6379
    password: ""
    database: 0

  # Cluster mode configuration (uncomment to use)
   nodes:
     - host: "redis-node-1"
       port: 6379
       password: ""
     - host: "redis-node-2"
       port: 6379
       password: ""
     - host: "redis-node-3"
       port: 6379
       password: ""

  monitoring:
    interval: 30  # seconds
    metrics:
      - memory_usage
      - cpu_usage
      - connected_clients
      - commands_per_sec
      - hit_rate
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

func (c *ConfigLoader) createDefaultKafkaConfig(path string) error {
	defaultConfig := `# Kafka Configuration
kafka:
  brokers:
    - "localhost:9092"

  security:
    protocol: "PLAINTEXT"  # PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL
    sasl_mechanism: "PLAIN"  # PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
    username: ""
    password: ""

  monitoring:
    interval: 30  # seconds
    topics: []  # Empty means monitor all topics
    consumer_groups: []  # Empty means monitor all consumer groups
    metrics:
      - broker_status
      - topic_partitions
      - consumer_lag
      - message_rate
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

func (c *ConfigLoader) createDefaultPostgreSQLConfig(path string) error {
	defaultConfig := `# PostgreSQL Configuration
postgresql:
  # Multi-database support
  databases:
    - name: "default"
      host: "localhost"
      port: 5432
      database: "postgres"
      username: "postgres"
      password: "postgres"
      ssl_mode: "disable"  # disable, allow, prefer, require, verify-ca, verify-full

      # Connection pool settings
      pool:
        min_connections: 5
        max_connections: 20
        max_idle_time: 300  # seconds
        connection_timeout: 10  # seconds

      # Monitoring settings
      monitoring:
        enabled: true
        interval: 30  # seconds
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1000  # milliseconds

  # Global monitoring settings
  monitoring:
    metrics:
      - connection_count
      - active_queries
      - database_size
      - cache_hit_ratio
      - transaction_rate
      - replication_lag
      - slow_queries
      - deadlocks

    # Health check settings
    health_check:
      enabled: true
      interval: 10  # seconds
      timeout: 5  # seconds

    # Alert thresholds
    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000  # milliseconds
      replication_lag_seconds: 10
      cache_hit_ratio_min: 90  # percentage
      disk_usage_percent: 85
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

func (c *ConfigLoader) createDefaultMySQLConfig(path string) error {
	defaultConfig := `# MySQL Configuration
mysql:
  # Multi-database support
  databases:
    - name: "default"
      host: "localhost"
      port: 3306
      database: "mysql"
      username: "root"
      password: "root"
      charset: "utf8mb4"

      # Connection pool settings
      pool:
        min_connections: 5
        max_connections: 20
        max_idle_time: 300  # seconds
        connection_timeout: 10  # seconds

      # Monitoring settings
      monitoring:
        enabled: true
        interval: 30  # seconds
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1000  # milliseconds

  # Global monitoring settings
  monitoring:
    metrics:
      - connection_count
      - query_cache_hit_ratio
      - slow_queries
      - table_locks
      - innodb_buffer_pool
      - threads_running
      - bytes_sent_received
      - aborted_connections
      - table_size
      - replication_status

    # Health check settings
    health_check:
      enabled: true
      interval: 10  # seconds
      timeout: 5  # seconds

    # Alert thresholds
    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000  # milliseconds
      query_cache_hit_ratio_min: 85  # percentage
      replication_lag_seconds: 10
      disk_usage_percent: 85
      threads_running_max: 100
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

// ==================== Getter Methods ====================

func (c *ConfigLoader) GetRedis() *domain.RedisConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redisConfig
}

func (c *ConfigLoader) GetKafka() *domain.KafkaConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.kafkaConfig
}

func (c *ConfigLoader) GetPostgreSQL() *domain.PostgreSQLConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.postgresConfig
}

func (c *ConfigLoader) GetMySQL() *domain.MySQLConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mysqlConfig
}

// ==================== Change Callback ====================

func (c *ConfigLoader) OnChange(callback func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changeCallbacks = append(c.changeCallbacks, callback)
}

func (c *ConfigLoader) notifyChanges() {
	c.mu.RLock()
	callbacks := c.changeCallbacks
	c.mu.RUnlock()

	for _, callback := range callbacks {
		go callback()
	}
}

// ==================== Helper Methods ====================

// ValidateConfig validates if a config string is valid YAML
func ValidateConfig(configStr string) error {
	var temp interface{}
	if err := yaml.Unmarshal([]byte(configStr), &temp); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	return nil
}

// BackupConfig creates a backup of existing config file
func (c *ConfigLoader) BackupConfig(filename string) error {
	path := filepath.Join(c.configPath, filename)
	backupPath := filepath.Join(c.configPath, filename+".backup")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to backup
		}
		return err
	}

	return os.WriteFile(backupPath, data, 0644)
}

// RestoreConfig restores config from backup
func (c *ConfigLoader) RestoreConfig(filename string) error {
	backupPath := filepath.Join(c.configPath, filename+".backup")
	path := filepath.Join(c.configPath, filename)

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// GetConfigPath returns the configuration directory path
func (c *ConfigLoader) GetConfigPath() string {
	return c.configPath
}

// ListConfigFiles returns list of all config files
func (c *ConfigLoader) ListConfigFiles() ([]string, error) {
	files, err := os.ReadDir(c.configPath)
	if err != nil {
		return nil, err
	}

	configFiles := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			configFiles = append(configFiles, file.Name())
		}
	}

	return configFiles, nil
}
