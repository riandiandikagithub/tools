export interface TopicMetrics {
  name: string;
  partitions: number;
  replicas: number;
  status: string;
  messageRate: number;
  size: string;
  retention: string;
  underReplicated: number;
}

export interface ClusterInfo {
  totalBrokers: number;
  totalTopics: number;
  totalPartitions: number;
  totalConsumerGroups: number;
  totalMessageRate: number; // messages per second
  totalDiskUsage: string; // contoh: "4.2TB"
  avgReplicationFactor: number;
}

export interface ConsumerGroup {
  name: string;
  topic: string;
  consumers: number;
  lag: number;
  status: "healthy" | "warning" | "offline";
}

export interface KafkaBroker {
  id: number | string;
  name: string;
  host: string;
  status: "online" | "offline" | "warning";

  diskUsed: string; // misal "120GB"
  diskTotal: string; // misal "500GB"
  diskUsage: number; // misal 24 (%)

  cpuUsage: number; // percent
  memoryUsage: number; // percent

  networkIn: string; // misal "12 MB/s"
  networkOut: string; // misal "10 MB/s"

  messageRate: number; // messages/s
}
