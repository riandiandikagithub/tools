// ==================== internal/infrastructure/config/defaults.go ====================
package config

// This file contains default configuration templates
// Separated for better organization and maintainability

const (
	// DefaultRedisConfig is the default Redis configuration template
	DefaultRedisConfig = `# Redis Configuration
redis:
  single:
    host: "localhost"
    port: 6379
    password: ""
    database: 0
  monitoring:
    interval: 30
    metrics:
      - memory_usage
      - cpu_usage
      - connected_clients
      - commands_per_sec
`

	// DefaultKafkaConfig is the default Kafka configuration template
	DefaultKafkaConfig = `# Kafka Configuration
kafka:
  brokers:
    - "localhost:9092"
  security:
    protocol: "PLAINTEXT"
    sasl_mechanism: "PLAIN"
    username: ""
    password: ""
  monitoring:
    interval: 30
    topics: []
    consumer_groups: []
    metrics:
      - broker_status
      - topic_partitions
      - consumer_lag
`

	// DefaultPostgreSQLConfig is the default PostgreSQL configuration template
	DefaultPostgreSQLConfig = `# PostgreSQL Configuration
postgresql:
  databases:
    - name: "default"
      host: "localhost"
      port: 5432
      database: "postgres"
      username: "postgres"
      password: "postgres"
      ssl_mode: "disable"
      pool:
        min_connections: 5
        max_connections: 20
        max_idle_time: 300
        connection_timeout: 10
      monitoring:
        enabled: true
        interval: 30
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1000
  monitoring:
    metrics:
      - connection_count
      - active_queries
      - database_size
      - cache_hit_ratio
    health_check:
      enabled: true
      interval: 10
      timeout: 5
    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000
      cache_hit_ratio_min: 90
`

	// DefaultMySQLConfig is the default MySQL configuration template
	DefaultMySQLConfig = `# MySQL Configuration
mysql:
  databases:
    - name: "default"
      host: "localhost"
      port: 3306
      database: "mysql"
      username: "root"
      password: "root"
      charset: "utf8mb4"
      pool:
        min_connections: 5
        max_connections: 20
        max_idle_time: 300
        connection_timeout: 10
      monitoring:
        enabled: true
        interval: 30
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1000
  monitoring:
    metrics:
      - connection_count
      - slow_queries
      - innodb_buffer_pool
      - threads_running
    health_check:
      enabled: true
      interval: 10
      timeout: 5
    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000
      threads_running_max: 100
`
)
