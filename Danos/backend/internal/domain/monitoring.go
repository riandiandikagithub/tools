// ==================== internal/domain/monitoring.go ====================
package domain

import "time"

// ==================== Redis Metrics ====================
type RedisMetrics struct {
	Name               string  `json:"name"`
	Mode               string  `json:"mode"`
	Host               string  `json:"host"`
	Port               int     `json:"port"`
	Status             string  `json:"status"` // online, offline, warning
	Role               string  `json:"role"`
	FragmentationRatio float64 `json:"fragmentation_ratio"`
	EvictedKeys        int64   `json:"evicted_keys"`
	ExpiredKeys        int64   `json:"expired_keys"`
	NetInputBytes      int64   `json:"net_input_bytes"`
	NetOutputBytes     int64   `json:"net_output_bytes"`
	// Clients
	ConnectedClients int64 `json:"connected_clients"`
	BlockedClients   int64 `json:"blocked_clients"`

	// Memory
	UsedMemory          int64   `json:"used_memory"`       // bytes
	UsedMemoryHuman     string  `json:"used_memory_human"` // human readable
	UsedMemoryRSS       int64   `json:"used_memory_rss"`   // bytes
	UsedMemoryPeak      int64   `json:"used_memory_peak"`  // bytes
	UsedMemoryPeakHuman string  `json:"used_memory_peak_human"`
	MaxMemory           int64   `json:"max_memory"` // bytes
	MemoryUsagePercent  float64 `json:"memory_usage_percent"`
	MemoryFragmentation float64 `json:"memory_fragmentation_ratio"`

	// CPU
	CPUUsage    float64 `json:"cpu_usage"`     // used_cpu_user
	CPUUsageSys float64 `json:"cpu_usage_sys"` // used_cpu_sys

	// Commands & Ops
	TotalCommands    int64 `json:"total_commands"`
	CommandsPerSec   int64 `json:"commands_per_sec"`
	InstantaneousOps int64 `json:"instantaneous_ops_per_sec"`

	// Uptime
	Uptime      int64  `json:"uptime"` // seconds
	UptimeHuman string `json:"uptime_human"`

	// Keyspace Stats
	KeyspaceHits   int64   `json:"keyspace_hits"`
	KeyspaceMisses int64   `json:"keyspace_misses"`
	HitRate        float64 `json:"hit_rate"` // percentage
	TotalKeys      int64   `json:"total_keys"`
	DatabaseCount  int     `json:"database_count"`

	// Replication
	ReplicationRole  string `json:"replication_role"` // master, slave
	ConnectedSlaves  int64  `json:"connected_slaves"`
	MasterReplOffset int64  `json:"master_repl_offset"`
	ReplicaOffset    int64  `json:"replica_offset"`

	// Persistence
	Loading                 int64 `json:"loading"`
	RDBLastSaveTime         int64 `json:"rdb_last_save_time"`
	RDBChangesSinceLastSave int64 `json:"rdb_changes_since_last_save"`
	AOFEnabled              bool  `json:"aof_enabled"`

	// Network
	NetworkInputBytes   int64 `json:"network_input_bytes"`
	NetworkOutputBytes  int64 `json:"network_output_bytes"`
	RejectedConnections int64 `json:"rejected_connections"`

	// Keyspace detailed info (db0, db1, ...)
	Keyspace map[string]RedisKeyspace `json:"keyspace"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`
}

// ==================== Redis Keyspace ====================
type RedisKeyspace struct {
	Keys    int64 `json:"keys"`
	Expires int64 `json:"expires"`
	AvgTTL  int64 `json:"avg_ttl"`
}

// ==================== Kafka Metrics ====================
type KafkaMetrics struct {
	BrokerID          int                    `json:"broker_id"`
	Host              string                 `json:"host"`
	Port              int                    `json:"port"`
	Status            string                 `json:"status"` // online, offline, warning
	Version           string                 `json:"version"`
	ClusterID         string                 `json:"cluster_id"`
	ControllerID      int                    `json:"controller_id"`
	Topics            []KafkaTopicMetrics    `json:"topics"`
	ConsumerGroups    []KafkaConsumerMetrics `json:"consumer_groups"`
	TotalPartitions   int                    `json:"total_partitions"`
	TotalTopics       int                    `json:"total_topics"`
	UnderReplicated   int                    `json:"under_replicated_partitions"`
	OfflinePartitions int                    `json:"offline_partitions"`
	ActiveControllers int                    `json:"active_controllers"`
	BytesInPerSec     float64                `json:"bytes_in_per_sec"`
	BytesOutPerSec    float64                `json:"bytes_out_per_sec"`
	MessagesInPerSec  float64                `json:"messages_in_per_sec"`
	Timestamp         time.Time              `json:"timestamp"`
}

type KafkaTopicMetrics struct {
	Name              string  `json:"name"`
	Partitions        int     `json:"partitions"`
	ReplicationFactor int     `json:"replication_factor"`
	ISRCount          int     `json:"isr_count"` // In-Sync Replicas
	MessagesPerSec    float64 `json:"messages_per_sec"`
	BytesInPerSec     float64 `json:"bytes_in_per_sec"`
	BytesOutPerSec    float64 `json:"bytes_out_per_sec"`
	TotalSize         int64   `json:"total_size"` // bytes
	RetentionMs       int64   `json:"retention_ms"`
	SegmentBytes      int64   `json:"segment_bytes"`
}

type KafkaConsumerMetrics struct {
	GroupID     string           `json:"group_id"`
	State       string           `json:"state"` // Stable, Dead, Empty, PreparingRebalance, CompletingRebalance
	Members     int              `json:"members"`
	Lag         int64            `json:"lag"`
	TopicLags   []KafkaTopicLag  `json:"topic_lags"`
	Coordinator KafkaCoordinator `json:"coordinator"`
}

type KafkaTopicLag struct {
	Topic         string `json:"topic"`
	Partition     int32  `json:"partition"`
	CurrentOffset int64  `json:"current_offset"`
	LogEndOffset  int64  `json:"log_end_offset"`
	Lag           int64  `json:"lag"`
	ConsumerID    string `json:"consumer_id,omitempty"`
	Host          string `json:"host,omitempty"`
}

type KafkaCoordinator struct {
	ID   int    `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

// ==================== PostgreSQL Metrics ====================
type PostgreSQLMetrics struct {
	Name               string            `json:"name"`
	Host               string            `json:"host"`
	Port               int               `json:"port"`
	Status             string            `json:"status"` // online, offline, warning
	Version            string            `json:"version"`
	DatabaseSize       int64             `json:"database_size"` // bytes
	DatabaseSizeHuman  string            `json:"database_size_human"`
	Connections        int               `json:"connections"`
	MaxConnections     int               `json:"max_connections"`
	ConnectionPercent  float64           `json:"connection_percent"`
	ActiveConnections  int               `json:"active_connections"`
	IdleConnections    int               `json:"idle_connections"`
	IdleInTransaction  int               `json:"idle_in_transaction"`
	ActiveQueries      int               `json:"active_queries"`
	WaitingQueries     int               `json:"waiting_queries"`
	CacheHitRatio      float64           `json:"cache_hit_ratio"` // percentage
	TransactionsPerSec float64           `json:"transactions_per_sec"`
	QueriesPerSec      float64           `json:"queries_per_sec"`
	TuplesReturned     int64             `json:"tuples_returned"`
	TuplesFetched      int64             `json:"tuples_fetched"`
	TuplesInserted     int64             `json:"tuples_inserted"`
	TuplesUpdated      int64             `json:"tuples_updated"`
	TuplesDeleted      int64             `json:"tuples_deleted"`
	BlocksRead         int64             `json:"blocks_read"`
	BlocksHit          int64             `json:"blocks_hit"`
	ReplicationLag     int64             `json:"replication_lag"` // bytes
	ReplicationState   string            `json:"replication_state"`
	IsReplica          bool              `json:"is_replica"`
	SlowQueries        int64             `json:"slow_queries"`
	Deadlocks          int64             `json:"deadlocks"`
	ConflictCount      int64             `json:"conflict_count"`
	TempFiles          int64             `json:"temp_files"`
	TempBytes          int64             `json:"temp_bytes"`
	Uptime             int64             `json:"uptime"` // seconds
	Databases          []DatabaseInfo    `json:"databases"`
	TableStatistics    []TableStatistics `json:"table_statistics,omitempty"`
	IndexStatistics    []IndexStatistics `json:"index_statistics,omitempty"`
	Timestamp          time.Time         `json:"timestamp"`
	UptimeSeconds      int64             `json:"uptime_seconds"`
	TableCount         int               `json:"table_count"`
	SizeBytes          int64             `json:"size_bytes"`
}

type DatabaseInfo struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"` // bytes
	SizeHuman    string `json:"size_human"`
	Connections  int    `json:"connections"`
	Transactions int64  `json:"transactions"`
	Commits      int64  `json:"commits"`
	Rollbacks    int64  `json:"rollbacks"`
}

type TableStatistics struct {
	SchemaName  string     `json:"schema_name"`
	TableName   string     `json:"table_name"`
	RowCount    int64      `json:"row_count"`
	TableSize   int64      `json:"table_size"` // bytes
	IndexSize   int64      `json:"index_size"` // bytes
	TotalSize   int64      `json:"total_size"` // bytes
	SeqScan     int64      `json:"seq_scan"`
	SeqTupRead  int64      `json:"seq_tup_read"`
	IdxScan     int64      `json:"idx_scan"`
	IdxTupFetch int64      `json:"idx_tup_fetch"`
	InsertCount int64      `json:"insert_count"`
	UpdateCount int64      `json:"update_count"`
	DeleteCount int64      `json:"delete_count"`
	LastVacuum  *time.Time `json:"last_vacuum,omitempty"`
	LastAnalyze *time.Time `json:"last_analyze,omitempty"`
}

type IndexStatistics struct {
	SchemaName  string `json:"schema_name"`
	TableName   string `json:"table_name"`
	IndexName   string `json:"index_name"`
	IndexSize   int64  `json:"index_size"` // bytes
	IndexScans  int64  `json:"index_scans"`
	TuplesRead  int64  `json:"tuples_read"`
	TuplesFetch int64  `json:"tuples_fetch"`
}

// ==================== MySQL Metrics ====================
type MySQLMetrics struct {
	Name                 string            `json:"name"`
	Host                 string            `json:"host"`
	Port                 int               `json:"port"`
	Status               string            `json:"status"` // online, offline, warning
	Version              string            `json:"version"`
	DatabaseSize         int64             `json:"database_size"` // bytes
	DatabaseSizeHuman    string            `json:"database_size_human"`
	Connections          int               `json:"connections"`
	MaxConnections       int               `json:"max_connections"`
	ConnectionPercent    float64           `json:"connection_percent"`
	ThreadsRunning       int               `json:"threads_running"`
	ThreadsConnected     int               `json:"threads_connected"`
	ThreadsCached        int               `json:"threads_cached"`
	ThreadsCreated       int               `json:"threads_created"`
	SlowQueries          int64             `json:"slow_queries"`
	QueriesPerSec        float64           `json:"queries_per_sec"`
	Questions            int64             `json:"questions"`
	Queries              int64             `json:"queries"`
	ComSelect            int64             `json:"com_select"`
	ComInsert            int64             `json:"com_insert"`
	ComUpdate            int64             `json:"com_update"`
	ComDelete            int64             `json:"com_delete"`
	InnoDBBufferPoolSize int64             `json:"innodb_buffer_pool_size"` // bytes
	InnoDBBufferPoolUsed int64             `json:"innodb_buffer_pool_used"` // bytes
	InnoDBBufferPoolHit  float64           `json:"innodb_buffer_pool_hit"`  // percentage
	InnoDBRowsRead       int64             `json:"innodb_rows_read"`
	InnoDBRowsInserted   int64             `json:"innodb_rows_inserted"`
	InnoDBRowsUpdated    int64             `json:"innodb_rows_updated"`
	InnoDBRowsDeleted    int64             `json:"innodb_rows_deleted"`
	BytesSent            int64             `json:"bytes_sent"`
	BytesReceived        int64             `json:"bytes_received"`
	TableLocks           int64             `json:"table_locks"`
	TableLocksWaited     int64             `json:"table_locks_waited"`
	AbortedConnections   int64             `json:"aborted_connections"`
	AbortedClients       int64             `json:"aborted_clients"`
	OpenTables           int               `json:"open_tables"`
	OpenFiles            int               `json:"open_files"`
	OpenTableDefinitions int               `json:"open_table_definitions"`
	Uptime               int64             `json:"uptime"` // seconds
	ReplicationStatus    *MySQLReplication `json:"replication_status,omitempty"`
	Databases            []DatabaseInfo    `json:"databases"`
	TableStatistics      []MySQLTableStats `json:"table_statistics,omitempty"`
	Timestamp            time.Time         `json:"timestamp"`
	SizeBytes            int64             `json:"size_bytes"`
	TableCount           int               `json:"table_count"`
	UptimeSeconds        int64             `json:"uptime_seconds"`
}

type MySQLReplication struct {
	SlaveIORunning      bool   `json:"slave_io_running"`
	SlaveSQLRunning     bool   `json:"slave_sql_running"`
	MasterHost          string `json:"master_host"`
	MasterPort          int    `json:"master_port"`
	SecondsBehindMaster *int64 `json:"seconds_behind_master"`
	MasterLogFile       string `json:"master_log_file"`
	ReadMasterLogPos    int64  `json:"read_master_log_pos"`
	RelayLogFile        string `json:"relay_log_file"`
	RelayLogPos         int64  `json:"relay_log_pos"`
	LastIOError         string `json:"last_io_error,omitempty"`
	LastSQLError        string `json:"last_sql_error,omitempty"`
}

type MySQLTableStats struct {
	SchemaName    string     `json:"schema_name"`
	TableName     string     `json:"table_name"`
	Engine        string     `json:"engine"`
	RowCount      int64      `json:"row_count"`
	DataLength    int64      `json:"data_length"`  // bytes
	IndexLength   int64      `json:"index_length"` // bytes
	DataFree      int64      `json:"data_free"`    // bytes
	TotalSize     int64      `json:"total_size"`   // bytes
	AutoIncrement *int64     `json:"auto_increment,omitempty"`
	CreateTime    *time.Time `json:"create_time,omitempty"`
	UpdateTime    *time.Time `json:"update_time,omitempty"`
}

// ==================== Aggregated Response ====================
type MetricsResponse struct {
	Redis      []RedisMetrics      `json:"redis"`
	Kafka      []KafkaMetrics      `json:"kafka"`
	PostgreSQL []PostgreSQLMetrics `json:"postgresql"`
	MySQL      []MySQLMetrics      `json:"mysql"`
	Summary    MetricsSummary      `json:"summary"`
	Timestamp  time.Time           `json:"timestamp"`
}

type MetricsSummary struct {
	TotalServices    int     `json:"total_services"`
	OnlineServices   int     `json:"online_services"`
	OfflineServices  int     `json:"offline_services"`
	WarningServices  int     `json:"warning_services"`
	HealthPercentage float64 `json:"health_percentage"`
	TotalConnections int     `json:"total_connections"`
	TotalDatabases   int     `json:"total_databases"`
	TotalMemoryUsed  int64   `json:"total_memory_used"` // bytes
	TotalDiskUsed    int64   `json:"total_disk_used"`   // bytes
}

// ==================== Alert Types ====================
type Alert struct {
	ID           string      `json:"id"`
	Type         AlertType   `json:"type"`
	Severity     AlertLevel  `json:"severity"`
	Service      string      `json:"service"` // redis, kafka, postgresql, mysql
	Instance     string      `json:"instance"`
	Message      string      `json:"message"`
	Value        interface{} `json:"value"`
	Threshold    interface{} `json:"threshold"`
	Timestamp    time.Time   `json:"timestamp"`
	Acknowledged bool        `json:"acknowledged"`
	ResolvedAt   *time.Time  `json:"resolved_at,omitempty"`
}

type AlertType string

const (
	AlertTypeMemory        AlertType = "memory"
	AlertTypeCPU           AlertType = "cpu"
	AlertTypeConnections   AlertType = "connections"
	AlertTypeDisk          AlertType = "disk"
	AlertTypeReplication   AlertType = "replication"
	AlertTypeSlowQuery     AlertType = "slow_query"
	AlertTypeConsumerLag   AlertType = "consumer_lag"
	AlertTypePartition     AlertType = "partition"
	AlertTypeCacheHitRatio AlertType = "cache_hit_ratio"
)

type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// ==================== Health Check ====================
type HealthStatus struct {
	Service   string    `json:"service"`
	Instance  string    `json:"instance"`
	Status    string    `json:"status"` // healthy, degraded, unhealthy
	Message   string    `json:"message,omitempty"`
	Latency   int64     `json:"latency"` // milliseconds
	Timestamp time.Time `json:"timestamp"`
}

type SystemHealth struct {
	Overall   string         `json:"overall"` // healthy, degraded, unhealthy
	Services  []HealthStatus `json:"services"`
	Timestamp time.Time      `json:"timestamp"`
}
