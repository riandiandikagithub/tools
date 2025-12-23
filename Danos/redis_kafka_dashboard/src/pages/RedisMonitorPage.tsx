import { useState, useEffect } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import RedisInstanceDrawer from "@/components/ui/RedisInstanceDrawer";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Progress } from "@/components/ui/progress";
import {
  Database,
  Server,
  Activity,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Search,
  Filter,
  Users,
  Cpu,
  HardDrive,
  Network,
  Clock,
  Zap,
} from "lucide-react";
import { Layout } from "@/components/Layout";
import { monitoringService } from "@/services/monitoringService";
import { MonitoringRedisData } from "@/models/monitoringModel";
import { parseHumanSize, calculateReplicationLag } from "@/lib/utils";

// Mock Redis data
export const redisInstances: MonitoringRedisData[] = [
  {
    name: "Redis Primary",
    mode: "Single",
    host: "10.0.1.10",
    port: 6379,
    status: "online",

    connected_clients: 1200,
    blocked_clients: 0,

    used_memory: 6800000000,
    used_memory_human: "6.8GB",
    used_memory_peak: 7000000000,
    used_memory_peak_human: "7GB",
    used_memory_rss: 6900000000,

    max_memory: 8000000000,
    memory_usage_percent: 85,
    memory_fragmentation_ratio: 1.02,

    cpu_usage: 45,
    cpu_usage_sys: 12,

    total_commands: 15000,
    commands_per_sec: 15000,
    instantaneous_ops_per_sec: 15000,

    uptime: 15 * 24 * 3600 + 4 * 3600 + 32 * 60,
    uptime_human: "15d 4h 32m",
    role: "master",

    keyspace_hits: 12000000,
    keyspace_misses: 300000,
    hit_rate: 97.5,

    evicted_keys: 100,
    expired_keys: 200,
    total_keys: 950000,
    database_count: 1,

    replication_role: "master",
    connected_slaves: 2,

    master_repl_offset: 900000000,
    replica_offset: 0, // master = no offset
    loading: 0,

    rdb_last_save_time: Math.floor(Date.now() / 1000) - 120,
    rdb_changes_since_last_save: 10,

    aof_enabled: false,

    network_input_bytes: 500000000,
    network_output_bytes: 850000000,

    rejected_connections: 5,

    keyspace: {
      db0: { keys: 950000, expires: 10000, avg_ttl: 86400000 },
    },

    timestamp: new Date().toISOString(),
  },

  {
    name: "Redis Replica 1",
    mode: "Single",
    host: "10.0.1.11",
    port: 6379,
    status: "online",

    connected_clients: 800,
    blocked_clients: 0,

    used_memory: 6600000000,
    used_memory_human: "6.6GB",
    used_memory_peak: 6800000000,
    used_memory_peak_human: "6.8GB",
    used_memory_rss: 6700000000,

    max_memory: 8000000000,
    memory_usage_percent: 82,
    memory_fragmentation_ratio: 1.01,

    cpu_usage: 38,
    cpu_usage_sys: 8,

    total_commands: 12000,
    commands_per_sec: 12000,
    instantaneous_ops_per_sec: 12000,

    uptime: 15 * 24 * 3600 + 4 * 3600 + 30 * 60,
    uptime_human: "15d 4h 30m",
    role: "master",

    keyspace_hits: 10000000,
    keyspace_misses: 200000,
    hit_rate: 98,

    evicted_keys: 80,
    expired_keys: 150,
    total_keys: 930000,
    database_count: 1,

    replication_role: "slave",
    connected_slaves: 0,

    master_repl_offset: 900000000,
    replica_offset: 899999500,
    loading: 0,

    rdb_last_save_time: Math.floor(Date.now() / 1000) - 240,
    rdb_changes_since_last_save: 20,

    aof_enabled: false,

    network_input_bytes: 400000000,
    network_output_bytes: 700000000,

    rejected_connections: 3,

    keyspace: {
      db0: { keys: 930000, expires: 8000, avg_ttl: 72000000 },
    },

    timestamp: new Date().toISOString(),
  },

  {
    name: "Redis Replica 2",
    mode: "Single",
    host: "10.0.1.12",
    port: 6379,
    status: "warning",

    connected_clients: 600,
    blocked_clients: 0,

    used_memory: 7600000000,
    used_memory_human: "7.6GB",
    used_memory_peak: 7700000000,
    used_memory_peak_human: "7.7GB",
    used_memory_rss: 7650000000,

    max_memory: 8000000000,
    memory_usage_percent: 95,
    memory_fragmentation_ratio: 1.05,

    cpu_usage: 78,
    cpu_usage_sys: 25,

    total_commands: 8000,
    commands_per_sec: 8000,
    instantaneous_ops_per_sec: 8000,

    uptime: 12 * 24 * 3600 + 2 * 3600 + 15 * 60,
    uptime_human: "12d 2h 15m",
    role: "master",

    keyspace_hits: 8000000,
    keyspace_misses: 300000,
    hit_rate: 96.38,

    evicted_keys: 200,
    expired_keys: 300,
    total_keys: 900000,
    database_count: 1,

    replication_role: "slave",
    connected_slaves: 0,

    master_repl_offset: 900000000,
    replica_offset: 899998000,
    loading: 0,

    rdb_last_save_time: Math.floor(Date.now() / 1000) - 300,
    rdb_changes_since_last_save: 30,

    aof_enabled: false,

    network_input_bytes: 350000000,
    network_output_bytes: 650000000,

    rejected_connections: 10,

    keyspace: {
      db0: { keys: 900000, expires: 7500, avg_ttl: 70000000 },
    },

    timestamp: new Date().toISOString(),
  },

  {
    name: "Redis Cache",
    mode: "Single",
    host: "10.0.1.13",
    port: 6379,
    status: "offline",

    connected_clients: 0,
    blocked_clients: 0,

    used_memory: 0,
    used_memory_human: "0GB",
    used_memory_peak: 0,
    used_memory_peak_human: "0GB",
    used_memory_rss: 0,

    max_memory: 4000000000,
    memory_usage_percent: 0,
    memory_fragmentation_ratio: 0,

    cpu_usage: 0,
    cpu_usage_sys: 0,

    total_commands: 0,
    commands_per_sec: 0,
    instantaneous_ops_per_sec: 0,

    uptime: 0,
    uptime_human: "0m",
    role: "master",

    keyspace_hits: 0,
    keyspace_misses: 0,
    hit_rate: 0,

    evicted_keys: 0,
    expired_keys: 0,
    total_keys: 0,
    database_count: 0,

    replication_role: "master",
    connected_slaves: 0,

    master_repl_offset: 0,
    replica_offset: 0,
    loading: 1, // offline node biasanya "loading=1"

    rdb_last_save_time: 0,
    rdb_changes_since_last_save: 0,

    aof_enabled: false,

    network_input_bytes: 0,
    network_output_bytes: 0,

    rejected_connections: 0,

    keyspace: {},

    timestamp: new Date().toISOString(),
  },
];

const redisClusterInfo = {
  totalNodes: 6,
  masterNodes: 3,
  slaveNodes: 3,
  totalSlots: 16384,
  assignedSlots: 16384,
  clusterState: "ok",
  totalMemory: "48GB",
  usedMemory: "38.2GB",
  totalConnections: 4200,
  totalCommandsPerSec: 45000,
};

export function RedisMonitorPage() {
  const [redisMonitoring, setRedisMonitoring] = useState(redisInstances);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedInstance, setSelectedInstance] =
    useState<MonitoringRedisData>();

  useEffect(() => {
    const loadMonitoringRedis = async () => {
      try {
        const monitoringRedis = await monitoringService.getMonitoringRedis();
        console.log(monitoringRedis);
        setRedisMonitoring(monitoringRedis ?? redisInstances);
      } catch (error) {
        console.error("Failed to load configs", error);
      }
    };
    loadMonitoringRedis();
  }, []);
  const getStatusIcon = (status: string) => {
    switch (status) {
      case "online":
        return <CheckCircle className="h-4 w-4 text-green-400" />;
      case "warning":
        return <AlertTriangle className="h-4 w-4 text-yellow-400" />;
      case "offline":
        return <XCircle className="h-4 w-4 text-red-400" />;
      default:
        return <XCircle className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusBadge = (status: string) => {
    const variants = {
      online: "default",
      warning: "secondary",
      offline: "destructive",
    } as const;

    return (
      <Badge
        variant={variants[status as keyof typeof variants] || "secondary"}
        className="capitalize"
      >
        {status}
      </Badge>
    );
  };

  const getRoleBadge = (role: string) => {
    return (
      <Badge
        variant={role === "master" ? "default" : "outline"}
        className="capitalize"
      >
        {role}
      </Badge>
    );
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <Database className="h-8 w-8 text-red-400" />
          <div>
            <h1 className="text-3xl font-bold">Redis Monitor</h1>
            <p className="text-muted-foreground">
              Monitor Redis instances and cluster health
            </p>
          </div>
        </div>

        <Tabs defaultValue="overview" className="space-y-6">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="instances">Instances</TabsTrigger>
            <TabsTrigger value="cluster">Cluster Info</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Cluster Summary */}
            <div className="dashboard-grid">
              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Total Nodes
                  </CardTitle>
                  <Server className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">
                    {redisClusterInfo.totalNodes}
                  </div>
                  <p className="metric-label">
                    {redisClusterInfo.masterNodes} masters,{" "}
                    {redisClusterInfo.slaveNodes} slaves
                  </p>
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Memory Usage
                  </CardTitle>
                  <HardDrive className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">79.6%</div>
                  <p className="metric-label">
                    {redisClusterInfo.usedMemory} /{" "}
                    {redisClusterInfo.totalMemory}
                  </p>
                  <Progress value={79.6} className="mt-2" />
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Connections
                  </CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">
                    {(redisClusterInfo.totalConnections ?? 0).toLocaleString()}
                  </div>
                  <p className="metric-label">Active connections</p>
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    Commands/sec
                  </CardTitle>
                  <Zap className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">
                    {(
                      redisClusterInfo.totalCommandsPerSec ?? 0
                    ).toLocaleString()}
                  </div>
                  <p className="metric-label">Operations per second</p>
                </CardContent>
              </Card>
            </div>

            {/* Quick Status */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {redisMonitoring.map((instance) => (
                <Card key={instance.port} className="p-4">
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center gap-2">
                      <Database className="h-4 w-4 text-red-400" />
                      <span className="font-medium text-sm">
                        {instance.name}
                      </span>
                    </div>
                    {getStatusIcon(instance.status)}
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">Memory</span>
                      <span>{instance.used_memory_human}</span>
                    </div>
                    <Progress
                      value={parseHumanSize(instance.used_memory_human)}
                      className="h-1"
                    />
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">CPU</span>
                      <span>{instance.cpu_usage}%</span>
                    </div>
                    <Progress value={instance.cpu_usage} className="h-1" />
                  </div>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="instances" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  Redis Instances
                </CardTitle>
                <div className="flex gap-4 mt-4">
                  <div className="flex-1">
                    <Input
                      placeholder="Search instances..."
                      className="bg-background"
                      icon={<Search className="h-4 w-4" />}
                    />
                  </div>
                  <Select defaultValue="all-status">
                    <SelectTrigger className="w-[140px]">
                      <Filter className="h-4 w-4 mr-2" />
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all-status">All Status</SelectItem>
                      <SelectItem value="online">Online</SelectItem>
                      <SelectItem value="warning">Warning</SelectItem>
                      <SelectItem value="offline">Offline</SelectItem>
                    </SelectContent>
                  </Select>
                  <Select defaultValue="all-roles">
                    <SelectTrigger className="w-[140px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all-roles">All Roles</SelectItem>
                      <SelectItem value="master">Master</SelectItem>
                      <SelectItem value="slave">Slave</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {redisMonitoring.map((instance) => (
                    <Card key={instance.port} className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <Database className="h-5 w-5 text-red-400" />
                          <div>
                            <div className="font-medium flex items-center gap-2">
                              {instance.name}
                              {getStatusIcon(instance.status)}
                              {getRoleBadge(instance.role)}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {instance.role}
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {getStatusBadge(instance.status)}
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedInstance(instance);
                              setDrawerOpen(true);
                            }}
                          >
                            View Details
                          </Button>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <HardDrive className="h-3 w-3" />
                            Memory
                          </div>
                          <div className="text-sm font-medium">
                            {instance.used_memory_human}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {instance.used_memory_human}
                          </div>
                          <Progress
                            value={instance.used_memory_human}
                            className="h-1"
                          />
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Cpu className="h-3 w-3" />
                            CPU
                          </div>
                          <div className="text-sm font-medium">
                            {instance.cpu_usage}%
                          </div>
                          <Progress
                            value={instance.cpu_usage}
                            className="h-1"
                          />
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Users className="h-3 w-3" />
                            Connections
                          </div>
                          <div className="text-sm font-medium">
                            {instance.connected_clients}
                          </div>
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Zap className="h-3 w-3" />
                            Commands/sec
                          </div>
                          <div className="text-sm font-medium">
                            {(instance.commands_per_sec ?? 0).toLocaleString()}
                          </div>
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Network className="h-3 w-3" />
                            Repl Lag
                          </div>
                          <div className="text-sm font-medium">
                            {calculateReplicationLag(
                              instance.master_repl_offset,
                              instance.replica_offset,
                            )}{" "}
                            bytes
                          </div>
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Clock className="h-3 w-3" />
                            Uptime
                          </div>
                          <div className="text-sm font-medium">
                            {instance.uptime}
                          </div>
                          <div className="text-xs text-muted-foreground">
                            {instance.lastUpdate}
                          </div>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="cluster" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Cluster Information</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Cluster State</span>
                    <Badge
                      variant="default"
                      className="bg-green-500/10 text-green-500"
                    >
                      {redisClusterInfo.clusterState.toUpperCase()}
                    </Badge>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Nodes</span>
                    <span>{redisClusterInfo.totalNodes}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Master Nodes</span>
                    <span>{redisClusterInfo.masterNodes}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Slave Nodes</span>
                    <span>{redisClusterInfo.slaveNodes}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Total Slots</span>
                    <span>{redisClusterInfo.totalSlots}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Assigned Slots
                    </span>
                    <span>{redisClusterInfo.assignedSlots}</span>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Performance Metrics</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">
                        Memory Usage
                      </span>
                      <span>79.6%</span>
                    </div>
                    <Progress value={79.6} />
                    <div className="text-xs text-muted-foreground">
                      {redisClusterInfo.usedMemory} /{" "}
                      {redisClusterInfo.totalMemory}
                    </div>
                  </div>

                  <div className="flex justify-between">
                    <span className="text-muted-foreground">
                      Total Connections
                    </span>
                    <span>
                      {(
                        redisClusterInfo.totalConnections ?? 0
                      ).toLocaleString()}
                    </span>
                  </div>

                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Commands/sec</span>
                    <span>
                      {(
                        redisClusterInfo.totalCommandsPerSec ?? 0
                      ).toLocaleString()}
                    </span>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
        </Tabs>
        <RedisInstanceDrawer
          isOpen={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          instance={selectedInstance}
          calcReplicationLag={calculateReplicationLag}
        />
      </div>
    </Layout>
  );
}
