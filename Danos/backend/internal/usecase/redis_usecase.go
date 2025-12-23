package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

// ClusterOverview is the result object you want.
type ClusterOverview struct {
	TotalNodes           int64  `json:"total_nodes"`
	MasterNodes          int64  `json:"master_nodes"`
	SlaveNodes           int64  `json:"slave_nodes"`
	TotalSlots           int64  `json:"total_slots"`
	AssignedSlots        int64  `json:"assigned_slots"`
	ClusterState         string `json:"cluster_state"`
	TotalMemoryBytes     int64  `json:"total_memory_bytes"`
	UsedMemoryBytes      int64  `json:"used_memory_bytes"`
	TotalMemoryHuman     string `json:"total_memory"`
	UsedMemoryHuman      string `json:"used_memory"`
	TotalConnections     int64  `json:"total_connections"`
	TotalCommandsPerSec  int64  `json:"total_commands_per_sec"`
	SampleNodeCountedFor int    `json:"sample_node_counted_for"`
}

// GetClusterOverview connects to provided node addresses (any node in cluster is fine)
// and returns an aggregated ClusterOverview.
// - addrs: list of "host:port" for cluster nodes. Use at least one reachable node.
// - ctx: context for timeouts/cancellation.
func (u *MonitoringUsecase) GetClusterOverview(ctx context.Context) (*ClusterOverview, error) {
	redisManager := u.redisManager

	addrs := redisManager.AddressCluster

	// Use first reachable node to fetch CLUSTER NODES/INFO (these commands are cluster-wide)
	var clusterNodesOutput string
	var clusterInfoOutput string
	var firstErr error
	// var firstAddr string

	// Try to get cluster nodes/cluster info from the first reachable node
	for _, addr := range addrs {
		// firstAddr = addr
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		defer client.Close()

		// CLUSTER NODES
		out, err := client.Do(ctx, "CLUSTER", "NODES").Text()
		if err != nil {
			firstErr = err
			continue
		}
		clusterNodesOutput = out

		// CLUSTER INFO
		outInfo, err := client.Do(ctx, "CLUSTER", "INFO").Text()
		if err != nil {
			// cluster info is nice-to-have - keep clusterNodesOutput but record err
			clusterInfoOutput = ""
		} else {
			clusterInfoOutput = outInfo
		}
		// success: break
		firstErr = nil
		break
	}

	if firstErr != nil {
		return nil, fmt.Errorf("unable to fetch CLUSTER NODES from any provided addr: %w", firstErr)
	}

	// parse cluster nodes
	nodes, err := parseClusterNodes(clusterNodesOutput)
	if err != nil {
		return nil, fmt.Errorf("parse cluster nodes failed: %w", err)
	}

	overview := &ClusterOverview{
		TotalSlots: 16384, // constant
	}

	// count nodes & assigned slots & roles based on parsed nodes
	for _, n := range nodes {
		overview.TotalNodes++
		if n.IsMaster {
			overview.MasterNodes++
			overview.AssignedSlots += n.AssignedSlotCount
		} else {
			overview.SlaveNodes++
		}
	}

	// parse cluster state from clusterInfoOutput
	overview.ClusterState = parseClusterState(clusterInfoOutput)

	// If clusterInfoOutput empty, fallback to "ok" for safety; else unknown
	if overview.ClusterState == "" {
		overview.ClusterState = "unknown"
	}

	// Now collect per-node metrics (INFO memory, clients, stats)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var sampleNodes int

	// We'll iterate nodes addresses parsed from CLUSTER NODES. If parse failed to
	// find address, fallback to `addrs` provided by caller.
	targetAddrs := make([]string, 0, len(nodes))
	for _, n := range nodes {
		if n.Addr != "" {
			targetAddrs = append(targetAddrs, n.Addr)
		}
	}
	if len(targetAddrs) == 0 {
		targetAddrs = addrs
	}

	// Limit concurrency (optional). We'll spawn one goroutine per node here.
	for _, addr := range targetAddrs {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			c := redis.NewClient(&redis.Options{Addr: a})
			defer c.Close()

			// INFO memory
			memRaw, err := c.Info(ctx, "memory").Result()
			if err != nil {
				// skip node on error but continue
				return
			}
			mem := parseInfoToMap(memRaw)

			// INFO clients
			clientsRaw, _ := c.Info(ctx, "clients").Result()
			clients := parseInfoToMap(clientsRaw)

			// INFO stats
			statsRaw, _ := c.Info(ctx, "stats").Result()
			stats := parseInfoToMap(statsRaw)

			var usedMemory int64
			var totalSystemMemory int64
			var maxMemory int64
			var connectedClients int64
			var instantaneousOps int64

			if v, ok := mem["used_memory"]; ok {
				usedMemory, _ = strconv.ParseInt(v, 10, 64)
			}
			// try total_system_memory if available
			if v, ok := mem["total_system_memory"]; ok {
				totalSystemMemory, _ = strconv.ParseInt(v, 10, 64)
			}
			// fall back to maxmemory if set
			if v, ok := mem["maxmemory"]; ok {
				maxMemory, _ = strconv.ParseInt(v, 10, 64)
			}

			if v, ok := clients["connected_clients"]; ok {
				connectedClients, _ = strconv.ParseInt(v, 10, 64)
			}

			if v, ok := stats["instantaneous_ops_per_sec"]; ok {
				instantaneousOps, _ = strconv.ParseInt(v, 10, 64)
			}

			mu.Lock()
			// If total_system_memory present use it, else if maxMemory > 0 use that, otherwise don't add to totalMemory
			if totalSystemMemory > 0 {
				overview.TotalMemoryBytes += totalSystemMemory
			} else if maxMemory > 0 {
				overview.TotalMemoryBytes += maxMemory
			}
			overview.UsedMemoryBytes += usedMemory
			overview.TotalConnections += connectedClients
			overview.TotalCommandsPerSec += instantaneousOps
			sampleNodes++
			mu.Unlock()
		}(addr)
	}

	wg.Wait()
	overview.SampleNodeCountedFor = sampleNodes

	overview.TotalMemoryHuman = bytesToHuman(overview.TotalMemoryBytes)
	overview.UsedMemoryHuman = bytesToHuman(overview.UsedMemoryBytes)

	return overview, nil
}

// parseClusterNodes parses CLUSTER NODES output and returns a slice of parsed nodes.
type parsedNode struct {
	ID                string
	Addr              string // host:port (without @bus-port)
	IsMaster          bool
	IsSlave           bool
	AssignedSlotCount int64
	Flags             []string
}

func parseClusterNodes(raw string) ([]parsedNode, error) {
	lines := strings.Split(raw, "\n")
	var res []parsedNode
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// fields: id addr flags master? ping pong epoch link-state slots...
		parts := strings.Fields(ln)
		if len(parts) < 3 {
			continue
		}
		id := parts[0]
		addr := parts[1]
		flags := strings.Split(parts[2], ",")
		p := parsedNode{
			ID:    id,
			Addr:  normalizeAddr(addr),
			Flags: flags,
		}
		// flags contain master or slave
		for _, f := range flags {
			if f == "master" {
				p.IsMaster = true
			}
			if f == "slave" {
				p.IsSlave = true
			}
		}

		// assigned slots appear after the 7th field usually. We'll parse any slot tokens present.
		assigned := int64(0)
		// slot tokens are typically at parts[8:], but safer to scan all parts for tokens like "0-5461" or "[10923->-]"
		for i := 8; i < len(parts); i++ {
			token := parts[i]
			if token == "" {
				continue
			}
			// ignore tokens starting with '[' (migrating/importing)
			if strings.HasPrefix(token, "[") {
				continue
			}
			// token can be a single number or range `0-5461`
			if strings.Contains(token, "-") {
				r := strings.Split(token, "-")
				if len(r) == 2 {
					start, err1 := strconv.ParseInt(r[0], 10, 64)
					end, err2 := strconv.ParseInt(r[1], 10, 64)
					if err1 == nil && err2 == nil && end >= start {
						assigned += (end - start + 1)
					}
				}
			} else {
				// single slot
				if _, err := strconv.ParseInt(token, 10, 64); err == nil {
					assigned++
				}
			}
		}
		p.AssignedSlotCount = assigned
		res = append(res, p)
	}
	if len(res) == 0 {
		return nil, errors.New("no cluster nodes parsed")
	}
	return res, nil
}

// normalizeAddr trims bus-port suffix "ip:port@bus-port" -> "ip:port"
func normalizeAddr(addr string) string {
	if strings.Contains(addr, "@") {
		parts := strings.SplitN(addr, "@", 2)
		return parts[0]
	}
	return addr
}

// parseClusterState reads cluster_state from CLUSTER INFO text
func parseClusterState(info string) string {
	lines := strings.Split(info, "\n")
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "cluster_state:") {
			parts := strings.SplitN(ln, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// parseInfoToMap parses a Redis INFO section into map[string]string
// Input is the full INFO section text (key:value lines and comments starting with #)
func parseInfoToMap(info string) map[string]string {
	out := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	if len(lines) == 1 {
		// sometimes server returns LF only
		lines = strings.Split(info, "\n")
	}
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		if !strings.Contains(ln, ":") {
			continue
		}
		parts := strings.SplitN(ln, ":", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		out[k] = v
	}
	return out
}

// bytesToHuman converts bytes to a human readable string with 1 decimal (GB/MB/KB)
func bytesToHuman(b int64) string {
	if b <= 0 {
		return "0B"
	}
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)
	f := float64(b)
	switch {
	case f >= TB:
		return fmt.Sprintf("%.1fT", f/TB)
	case f >= GB:
		return fmt.Sprintf("%.1fG", f/GB)
	case f >= MB:
		return fmt.Sprintf("%.1fM", f/MB)
	case f >= KB:
		return fmt.Sprintf("%.1fK", f/KB)
	default:
		return fmt.Sprintf("%dB", b)
	}
}
