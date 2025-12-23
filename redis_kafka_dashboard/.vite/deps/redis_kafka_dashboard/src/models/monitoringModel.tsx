export interface RedisKeyspace {
  keys: number;
  expires: number;
  avg_ttl: number;
}

export interface MonitoringRedisData {
  name: string;
  mode: string;
  host: string;
  port: number;
  status: string;

  connected_clients: number;
  blocked_clients: number;

  used_memory: number;
  used_memory_human: string;
  used_memory_rss: number;
  used_memory_peak: number;
  used_memory_peak_human: string;

  max_memory: number;
  memory_usage_percent: number;
  memory_fragmentation_ratio: number;

  cpu_usage: number; // used_cpu_user
  cpu_usage_sys: number; // used_cpu_sys

  total_commands: number;
  commands_per_sec: number;
  instantaneous_ops_per_sec: number;

  uptime: number;
  uptime_human: string;
  role: string;

  keyspace_hits: number;
  keyspace_misses: number;
  hit_rate: number;

  evicted_keys: number;
  expired_keys: number;
  total_keys: number;
  database_count: number;

  replication_role: string;
  connected_slaves: number;

  master_repl_offset: number;
  replica_offset: number;

  loading: number;

  rdb_last_save_time: number;
  rdb_changes_since_last_save: number;

  aof_enabled: boolean;

  network_input_bytes: number;
  network_output_bytes: number;

  rejected_connections: number;

  keyspace: Record<string, RedisKeyspace>;

  timestamp: string;
}
