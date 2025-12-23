// ==================== internal/infrastructure/redis/client.go ====================
package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Danos/backend/internal/domain"

	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	mu             sync.RWMutex
	AddressSingle  string
	AddressCluster []string
	clients        map[string]*redis.Client
	cluster        *redis.ClusterClient
	config         *domain.RedisConfig
	ctx            context.Context
}

func NewRedisManager() *RedisManager {
	return &RedisManager{
		clients: make(map[string]*redis.Client),
		ctx:     context.Background(),
	}
}

func (m *RedisManager) Initialize(config *domain.RedisConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config

	// if config.Mode == "cluster" {
	if config != nil && len(config.Nodes) > 0 {
		if len(config.Nodes) > 0 {
			return m.connectSingleNCluster()
		}
	}
	return m.connectSingle()
}


func (m *RedisManager) connectCluster() error {
	if len(m.config.Nodes) == 0 {
		return fmt.Errorf("single mode config is nil")
	}
	addrs := make([]string, 0, len(m.config.Nodes))
	var password string

	for _, node := range m.config.Nodes {
		addr := fmt.Sprintf("%s:%d", node.Host, node.Port)
		addrs = append(addrs, addr)
		if password == "" {
			password = node.Password
		}
	}

	m.cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	if err := m.cluster.Ping(m.ctx).Err(); err != nil {
		return fmt.Errorf("cluster connection failed: %w", err)
	}

	log.Println("Connected to Redis cluster")
	return nil
}

func (m *RedisManager) connectSingleNCluster() error {

	addrs := make([]string, 0, len(m.config.Nodes))

	//connect to single too
	m.connectSingle()

	var password string
	if len(m.config.Nodes) == 0 {
		fmt.Println("Just Connect Redis Single")
		return nil
	}
	for _, node := range m.config.Nodes {
		addr := fmt.Sprintf("%s:%d", node.Host, node.Port)
		addrs = append(addrs, addr)
		if password == "" {
			password = node.Password
		}
	}
	m.AddressCluster = addrs
	m.cluster = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	if err := m.cluster.Ping(m.ctx).Err(); err != nil {
		return fmt.Errorf("cluster connection failed: %w", err)
	}

	log.Println("Connected to Redis cluster")
	return nil
}
func (m *RedisManager) connectSingle() error {
	if m.config.Single == nil {
		return fmt.Errorf("single mode config is nil")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", m.config.Single.Host, m.config.Single.Port),
		Password: m.config.Single.Password,
		DB:       m.config.Single.Database,
	})
	m.AddressSingle = fmt.Sprintf("%s:%d", m.config.Single.Host, m.config.Single.Port)s
	if err := client.Ping(m.ctx).Err(); err != nil {
		return fmt.Errorf("single connection failed: %w", err)
	}

	m.clients["single"] = client

	log.Println("Connected to Redis single mode")
	return nil
}

func (m *RedisManager) Reconnect(config *domain.RedisConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Println("Reconnecting Redis clients...")

	// Close existing connections
	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis %s: %v", name, err)
		}
	}
	if m.cluster != nil {
		if err := m.cluster.Close(); err != nil {
			log.Printf("Error closing Redis cluster: %v", err)
		}
	}

	m.clients = make(map[string]*redis.Client)
	m.cluster = nil
	m.config = config

	if len(config.Nodes) > 0 {
		return m.connectCluster()
	}
	return m.connectSingle()
}

// modeRequest:
// - "auto"    → ikut config (default)
// - "single"  → paksa ambil single redis saja
// - "cluster" → paksa ambil cluster saja

func (m *RedisManager) GetMetrics(modeRequest string) ([]domain.RedisMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make([]domain.RedisMetrics, 0)

	// normalize request
	mode := strings.ToLower(modeRequest)
	if mode == "" || mode == "auto" {
		mode = "single" // follow config
	}

	switch mode {
	case "cluster":
		if m.cluster == nil {
			return nil, fmt.Errorf("redis is not configured as cluster")
		}

		for _, node := range m.config.Nodes {
			metric, err := m.getNodeMetrics(node)
			if err != nil {
				log.Printf("Error getting metrics for %s: %v", node.Host, err)
				continue
			}
			metrics = append(metrics, metric)
		}

	case "single":
		if m.config.Single == nil {
			return nil, fmt.Errorf("redis single is not configured")
		}

		client := m.clients["single"]
		if client == nil {
			return nil, fmt.Errorf("redis single client is nil")
		}

		metric, err := m.getSingleMetrics(client, m.config.Single)
		if err != nil {
			log.Printf("Error getting single metrics: %v", err)
		} else {
			metrics = append(metrics, metric)
		}

	default:
		return nil, fmt.Errorf("unknown redis mode '%s'", mode)
	}

	return metrics, nil
}

func (m *RedisManager) getNodeMetrics(node domain.RedisNode) (domain.RedisMetrics, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", node.Host, node.Port),
		Password: node.Password,
	})
	defer client.Close()

	info, err := client.Info(m.ctx).Result()
	if err != nil {
		return domain.RedisMetrics{}, err
	}

	return m.parseInfo(info, "cluster", node.Host, node.Port)
}

func (m *RedisManager) getSingleMetrics(client *redis.Client, config *domain.RedisSingle) (domain.RedisMetrics, error) {
	info, err := client.Info(m.ctx).Result()
	if err != nil {
		return domain.RedisMetrics{}, err
	}

	return m.parseInfo(info, "Single", config.Host, config.Port)
}
func (m *RedisManager) parseInfo(info, mode, host string, port int) (domain.RedisMetrics, error) {
	metrics := domain.RedisMetrics{
		Name:      fmt.Sprintf("%s:%d", host, port),
		Mode:      mode,
		Host:      host,
		Port:      port,
		Status:    "online",
		Timestamp: time.Now(),
	}

	lines := strings.Split(info, "\r\n")

	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {

			// =======================
			//   SERVER INFO
			// =======================
			case "role":
				metrics.Role = value

			// =======================
			//   CLIENTS
			// =======================
			case "connected_clients":
				metrics.ConnectedClients, _ = strconv.ParseInt(value, 10, 64)
			case "blocked_clients":
				metrics.BlockedClients, _ = strconv.ParseInt(value, 10, 64)

			// =======================
			//   MEMORY
			// =======================
			case "used_memory":
				metrics.UsedMemory, _ = strconv.ParseInt(value, 10, 64)
			case "used_memory_human":
				metrics.UsedMemoryHuman = value
			case "used_memory_rss":
				metrics.UsedMemoryRSS, _ = strconv.ParseInt(value, 10, 64)
			case "maxmemory":
				metrics.MaxMemory, _ = strconv.ParseInt(value, 10, 64)
			case "mem_fragmentation_ratio":
				metrics.FragmentationRatio, _ = strconv.ParseFloat(value, 64)

			// =======================
			//   CPU
			// =======================
			case "used_cpu_sys":
				metrics.CPUUsage, _ = strconv.ParseFloat(value, 64)

			// =======================
			//   COMMAND STATS
			// =======================
			case "total_commands_processed":
				metrics.TotalCommands, _ = strconv.ParseInt(value, 10, 64)
			case "instantaneous_ops_per_sec":
				metrics.CommandsPerSec, _ = strconv.ParseInt(value, 10, 64)

			// =======================
			//   KEYSPACE
			// =======================
			case "keyspace_hits":
				metrics.KeyspaceHits, _ = strconv.ParseInt(value, 10, 64)
			case "keyspace_misses":
				metrics.KeyspaceMisses, _ = strconv.ParseInt(value, 10, 64)

			case "evicted_keys":
				metrics.EvictedKeys, _ = strconv.ParseInt(value, 10, 64)
			case "expired_keys":
				metrics.ExpiredKeys, _ = strconv.ParseInt(value, 10, 64)

			// =======================
			//   NETWORK
			// =======================
			case "total_net_input_bytes":
				metrics.NetInputBytes, _ = strconv.ParseInt(value, 10, 64)
			case "total_net_output_bytes":
				metrics.NetOutputBytes, _ = strconv.ParseInt(value, 10, 64)

			// =======================
			//   UPTIME
			// =======================
			case "uptime_in_seconds":
				metrics.Uptime, _ = strconv.ParseInt(value, 10, 64)
				metrics.UptimeHuman = formatUptime(metrics.Uptime)

			// =======================
			//   REPLICATION
			// =======================
			case "connected_slaves":
				metrics.ConnectedSlaves, _ = strconv.ParseInt(value, 10, 64)
			case "master_repl_offset":
				metrics.MasterReplOffset, _ = strconv.ParseInt(value, 10, 64)
			case "slave_repl_offset": // untuk replica
				metrics.ReplicaOffset, _ = strconv.ParseInt(value, 10, 64)

			// =======================
			//   PERSISTENCE
			// =======================
			case "loading":
				metrics.Loading, _ = strconv.ParseInt(value, 10, 64)
			case "rdb_last_save_time":
				metrics.RDBLastSaveTime, _ = strconv.ParseInt(value, 10, 64)
			case "aof_enabled":
				metrics.AOFEnabled = (value == "1")
			}
		}

		// =======================
		//     PARSE KEYSPACE (db0, db1…)
		// =======================
		if strings.HasPrefix(line, "db") {
			// contoh: db0:keys=1000,expires=10,avg_ttl=5000
			parts := strings.SplitN(line, ":", 2)
			dbName := parts[0]
			stats := strings.Split(parts[1], ",")

			dbData := domain.RedisKeyspace{}

			for _, s := range stats {
				kv := strings.Split(s, "=")
				if len(kv) != 2 {
					continue
				}

				switch kv[0] {
				case "keys":
					dbData.Keys, _ = strconv.ParseInt(kv[1], 10, 64)
				case "expires":
					dbData.Expires, _ = strconv.ParseInt(kv[1], 10, 64)
				case "avg_ttl":
					dbData.AvgTTL, _ = strconv.ParseInt(kv[1], 10, 64)
				}
			}

			if metrics.Keyspace == nil {
				metrics.Keyspace = make(map[string]domain.RedisKeyspace)
			}
			metrics.Keyspace[dbName] = dbData
		}
	}

	// =======================
	//  DERIVED VALUES
	// =======================
	if metrics.MaxMemory > 0 {
		metrics.MemoryUsagePercent = float64(metrics.UsedMemory) / float64(metrics.MaxMemory) * 100
	}

	if metrics.KeyspaceHits+metrics.KeyspaceMisses > 0 {
		metrics.HitRate = float64(metrics.KeyspaceHits) /
			float64(metrics.KeyspaceHits+metrics.KeyspaceMisses) * 100
	}

	return metrics, nil
}

func (m *RedisManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis %s: %v", name, err)
		}
	}
	if m.cluster != nil {
		if err := m.cluster.Close(); err != nil {
			log.Printf("Error closing Redis cluster: %v", err)
		}
	}
	return nil
}

func formatUptime(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	return fmt.Sprintf("%dd %dh", days, hours)
}

// Add these methods to RedisManager

func (m *RedisManager) IsConnected(modeRequest string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if modeRequest == "cluster" {
		return m.cluster != nil
	}
	return len(m.clients) > 0
}

func (m *RedisManager) GetConnectionInfo(modeRequest string) []map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info := make([]map[string]interface{}, 0)

	if modeRequest == "cluster" {
		for _, node := range m.config.Nodes {
			info = append(info, map[string]interface{}{
				"host": node.Host,
				"port": node.Port,
				"mode": "cluster",
			})
		}
	} else if m.config.Single != nil {
		info = append(info, map[string]interface{}{
			"host": m.config.Single.Host,
			"port": m.config.Single.Port,
			"mode": "single",
		})
	}

	return info
}
