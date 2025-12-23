export interface MySQLOperations {
  qps: number; // Queries per second
  tps: number; // Transactions per second
  innodbReads: number; // InnoDB read ops/sec
  innodbWrites: number; // InnoDB write ops/sec
}

export interface MySQLInfo {
  version: string;
  threads: number;
  uptime: string; // contoh: "12d 4h 33m" atau "36000s"
  activeConnections: number;
}

export interface PostgresInfo {
  version: string;
  maxConnections: number;
  activeConnections: number;
}

export interface PostgresStats {
  tps: number;
  cacheHitRatio: number; // dalam persen
  deadlocks: number;
  conflicts: number;
}

export interface RedisInfo {
  version: string;
  mode: "standalone" | "cluster" | "replica"; // umum di Redis
  connectedClients: number;
  uptime: string; // contoh: "5d 3h 12m"
}

export interface RedisStats {
  opsPerSec: number;
  memoryUsed: string; // contoh "120MB"
  memoryPeak: string; // contoh "256MB"
  keyHits: number;
  keyMisses: number;
  pubsubChannels: number;
}
