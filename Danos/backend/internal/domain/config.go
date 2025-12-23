package domain

// ==================== Redis Configuration ====================
type RedisConfig struct {
	Nodes      []RedisNode     `yaml:"nodes" json:"nodes"`
	Single     *RedisSingle    `yaml:"single" json:"single"`
	Monitoring RedisMonitoring `yaml:"monitoring" json:"monitoring"`
}

type RedisNode struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
}

type RedisSingle struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	Database int    `yaml:"database" json:"database"`
}

type RedisMonitoring struct {
	Interval int      `yaml:"interval" json:"interval"` // seconds
	Metrics  []string `yaml:"metrics" json:"metrics"`
}

// ==================== Kafka Configuration ====================
type KafkaConfig struct {
	Brokers    []string        `yaml:"brokers" json:"brokers"`
	Security   KafkaSecurity   `yaml:"security" json:"security"`
	Monitoring KafkaMonitoring `yaml:"monitoring" json:"monitoring"`
}

type KafkaSecurity struct {
	Protocol      string `yaml:"protocol" json:"protocol"`             // PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL
	SASLMechanism string `yaml:"sasl_mechanism" json:"sasl_mechanism"` // PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
	Username      string `yaml:"username" json:"username"`
	Password      string `yaml:"password" json:"password"`
}

type KafkaMonitoring struct {
	Interval       int      `yaml:"interval" json:"interval"` // seconds
	Topics         []string `yaml:"topics" json:"topics"`
	ConsumerGroups []string `yaml:"consumer_groups" json:"consumer_groups"`
	Metrics        []string `yaml:"metrics" json:"metrics"`
}

// ==================== PostgreSQL Configuration ====================
type PostgreSQLConfig struct {
	Databases  []PostgreSQLDatabase `yaml:"databases" json:"databases"`
	Monitoring PostgreSQLMonitoring `yaml:"monitoring" json:"monitoring"`
}

type PostgreSQLDatabase struct {
	Name       string             `yaml:"name" json:"name"`
	Host       string             `yaml:"host" json:"host"`
	Port       int                `yaml:"port" json:"port"`
	Database   string             `yaml:"database" json:"database"`
	Username   string             `yaml:"username" json:"username"`
	Password   string             `yaml:"password" json:"password"`
	SSLMode    string             `yaml:"ssl_mode" json:"ssl_mode"` // disable, allow, prefer, require, verify-ca, verify-full
	Pool       ConnectionPool     `yaml:"pool" json:"pool"`
	Monitoring DatabaseMonitoring `yaml:"monitoring" json:"monitoring"`
}

type PostgreSQLMonitoring struct {
	Metrics     []string         `yaml:"metrics" json:"metrics"`
	HealthCheck HealthCheck      `yaml:"health_check" json:"health_check"`
	Alerts      MonitoringAlerts `yaml:"alerts" json:"alerts"`
}

// ==================== MySQL Configuration ====================
type MySQLConfig struct {
	Databases  []MySQLDatabase `yaml:"databases" json:"databases"`
	Monitoring MySQLMonitoring `yaml:"monitoring" json:"monitoring"`
}

type MySQLDatabase struct {
	Name       string             `yaml:"name" json:"name"`
	Host       string             `yaml:"host" json:"host"`
	Port       int                `yaml:"port" json:"port"`
	Database   string             `yaml:"database" json:"database"`
	Username   string             `yaml:"username" json:"username"`
	Password   string             `yaml:"password" json:"password"`
	Charset    string             `yaml:"charset" json:"charset"` // utf8mb4
	Pool       ConnectionPool     `yaml:"pool" json:"pool"`
	Monitoring DatabaseMonitoring `yaml:"monitoring" json:"monitoring"`
}

type MySQLMonitoring struct {
	Metrics     []string         `yaml:"metrics" json:"metrics"`
	HealthCheck HealthCheck      `yaml:"health_check" json:"health_check"`
	Alerts      MonitoringAlerts `yaml:"alerts" json:"alerts"`
}

// ==================== Shared Configuration Types ====================
type ConnectionPool struct {
	MinConnections    int `yaml:"min_connections" json:"min_connections"`
	MaxConnections    int `yaml:"max_connections" json:"max_connections"`
	MaxIdleTime       int `yaml:"max_idle_time" json:"max_idle_time"`           // seconds
	ConnectionTimeout int `yaml:"connection_timeout" json:"connection_timeout"` // seconds
}

type DatabaseMonitoring struct {
	Enabled            bool `yaml:"enabled" json:"enabled"`
	Interval           int  `yaml:"interval" json:"interval"` // seconds
	TrackActivity      bool `yaml:"track_activity" json:"track_activity"`
	LogSlowQueries     bool `yaml:"log_slow_queries" json:"log_slow_queries"`
	SlowQueryThreshold int  `yaml:"slow_query_threshold" json:"slow_query_threshold"` // milliseconds
}

type HealthCheck struct {
	Enabled  bool `yaml:"enabled" json:"enabled"`
	Interval int  `yaml:"interval" json:"interval"` // seconds
	Timeout  int  `yaml:"timeout" json:"timeout"`   // seconds
}

type MonitoringAlerts struct {
	MaxConnectionsPercent int `yaml:"max_connections_percent" json:"max_connections_percent"`
	SlowQueryThreshold    int `yaml:"slow_query_threshold" json:"slow_query_threshold"` // milliseconds
	ReplicationLagSeconds int `yaml:"replication_lag_seconds" json:"replication_lag_seconds"`
	CacheHitRatioMin      int `yaml:"cache_hit_ratio_min" json:"cache_hit_ratio_min"` // percentage
	DiskUsagePercent      int `yaml:"disk_usage_percent" json:"disk_usage_percent"`
}

// ==================== Configuration Wrapper ====================
type ConfigWrapper struct {
	Redis      RedisConfig      `yaml:"redis" json:"redis"`
	Kafka      KafkaConfig      `yaml:"kafka" json:"kafka"`
	PostgreSQL PostgreSQLConfig `yaml:"postgresql" json:"postgresql"`
	MySQL      MySQLConfig      `yaml:"mysql" json:"mysql"`
}
