import React, { useState, useEffect } from "react";
import {
  Database,
  Server,
  Activity,
  HardDrive,
  Clock,
  Zap,
  AlertTriangle,
  CheckCircle,
  XCircle,
} from "lucide-react";
import { Layout } from "@/components/Layout";

export function DatabaseMonitorPage() {
  const [selectedDb, setSelectedDb] = useState("all");
  const [searchTerm, setSearchTerm] = useState("");
  const [filterStatus, setFilterStatus] = useState("all");

  // Mock data untuk database instances
  const [databases] = useState([
    {
      id: "db-001",
      name: "MySQL Production",
      type: "MySQL",
      version: "8.0.35",
      host: "prod-mysql-01.us-east-1.rds.amazonaws.com",
      status: "online",
      connections: { current: 145, max: 500 },
      memory: 78,
      cpu: 45,
      storage: { used: 245, total: 500, unit: "GB" },
      queries: { total: 15420, slow: 12 },
      uptime: "45d 12h",
      latency: 1.8,
      lastBackup: "2 hours ago",
      databases: 24,
      tables: 456,
    },
    {
      id: "db-002",
      name: "PostgreSQL Analytics",
      type: "PostgreSQL",
      version: "15.4",
      host: "analytics-pg-01.us-east-1.rds.amazonaws.com",
      status: "online",
      connections: { current: 89, max: 300 },
      memory: 82,
      cpu: 62,
      storage: { used: 512, total: 1000, unit: "GB" },
      queries: { total: 8945, slow: 45 },
      uptime: "30d 8h",
      latency: 2.1,
      lastBackup: "4 hours ago",
      databases: 12,
      tables: 289,
    },
    {
      id: "db-003",
      name: "MySQL Development",
      type: "MySQL",
      version: "8.0.35",
      host: "dev-mysql-01.us-east-1.rds.amazonaws.com",
      status: "warning",
      connections: { current: 420, max: 500 },
      memory: 91,
      cpu: 88,
      storage: { used: 180, total: 250, unit: "GB" },
      queries: { total: 9821, slow: 156 },
      uptime: "12d 3h",
      latency: 4.5,
      lastBackup: "1 hour ago",
      databases: 8,
      tables: 123,
    },
    {
      id: "db-004",
      name: "PostgreSQL Archive",
      type: "PostgreSQL",
      version: "14.9",
      host: "archive-pg-01.us-west-2.rds.amazonaws.com",
      status: "offline",
      connections: { current: 0, max: 100 },
      memory: 0,
      cpu: 0,
      storage: { used: 890, total: 1000, unit: "GB" },
      queries: { total: 0, slow: 0 },
      uptime: "Offline",
      latency: 0,
      lastBackup: "8 hours ago",
      databases: 5,
      tables: 87,
    },
    {
      id: "db-005",
      name: "MySQL Staging",
      type: "MySQL",
      version: "8.0.34",
      host: "staging-mysql-01.us-east-1.rds.amazonaws.com",
      status: "online",
      connections: { current: 67, max: 200 },
      memory: 65,
      cpu: 38,
      storage: { used: 95, total: 200, unit: "GB" },
      queries: { total: 5234, slow: 8 },
      uptime: "22d 15h",
      latency: 2.3,
      lastBackup: "30 minutes ago",
      databases: 6,
      tables: 98,
    },
    {
      id: "db-006",
      name: "PostgreSQL Reports",
      type: "PostgreSQL",
      version: "15.4",
      host: "reports-pg-01.eu-west-1.rds.amazonaws.com",
      status: "online",
      connections: { current: 34, max: 150 },
      memory: 58,
      cpu: 42,
      storage: { used: 340, total: 500, unit: "GB" },
      queries: { total: 3421, slow: 23 },
      uptime: "60d 2h",
      latency: 3.2,
      lastBackup: "5 hours ago",
      databases: 8,
      tables: 145,
    },
  ]);

  // Calculate statistics
  const stats = {
    total: databases.length,
    online: databases.filter((db) => db.status === "online").length,
    warning: databases.filter((db) => db.status === "warning").length,
    offline: databases.filter((db) => db.status === "offline").length,
    avgMemory: Math.round(
      databases.reduce((acc, db) => acc + db.memory, 0) / databases.length,
    ),
    avgCpu: Math.round(
      databases.reduce((acc, db) => acc + db.cpu, 0) / databases.length,
    ),
    totalConnections: databases.reduce(
      (acc, db) => acc + db.connections.current,
      0,
    ),
    avgLatency: (
      databases.reduce((acc, db) => acc + db.latency, 0) / databases.length
    ).toFixed(1),
  };

  // Filter databases
  const filteredDatabases = databases.filter((db) => {
    const matchesSearch =
      db.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      db.host.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = selectedDb === "all" || db.type === selectedDb;
    const matchesStatus = filterStatus === "all" || db.status === filterStatus;
    return matchesSearch && matchesType && matchesStatus;
  });

  const getStatusColor = (status) => {
    switch (status) {
      case "online":
        return "text-green-400";
      case "warning":
        return "text-yellow-400";
      case "offline":
        return "text-red-400";
      default:
        return "text-gray-400";
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case "online":
        return <CheckCircle className="w-4 h-4" />;
      case "warning":
        return <AlertTriangle className="w-4 h-4" />;
      case "offline":
        return <XCircle className="w-4 h-4" />;
      default:
        return null;
    }
  };

  const getStatusBadge = (status) => {
    const colors = {
      online: "bg-green-500/20 text-green-400 border-green-500/30",
      warning: "bg-yellow-500/20 text-yellow-400 border-yellow-500/30",
      offline: "bg-red-500/20 text-red-400 border-red-500/30",
    };
    return `px-3 py-1 rounded-full text-xs font-medium border ${colors[status]}`;
  };

  const getMetricColor = (value) => {
    if (value >= 90) return "text-red-400";
    if (value >= 75) return "text-yellow-400";
    return "text-green-400";
  };

  return (
    <Layout>
      <div className="min-h-screen bg-gray-950 text-gray-100 p-6">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <Database className="w-8 h-8 text-blue-400" />
            <h1 className="text-3xl font-bold">Database Monitor</h1>
          </div>
          <p className="text-gray-400">
            Monitor and manage MySQL & PostgreSQL databases
          </p>
        </div>

        {/* Statistics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Total Databases</span>
              <Server className="w-5 h-5 text-gray-400" />
            </div>
            <div className="text-3xl font-bold">{stats.total}</div>
            <div className="text-xs text-gray-500 mt-1">MySQL + PostgreSQL</div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Online</span>
              <CheckCircle className="w-5 h-5 text-green-400" />
            </div>
            <div className="text-3xl font-bold text-green-400">
              {stats.online}
            </div>
            <div className="text-xs text-gray-500 mt-1">Healthy instances</div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Warning</span>
              <AlertTriangle className="w-5 h-5 text-yellow-400" />
            </div>
            <div className="text-3xl font-bold text-yellow-400">
              {stats.warning}
            </div>
            <div className="text-xs text-gray-500 mt-1">Need attention</div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Offline</span>
              <XCircle className="w-5 h-5 text-red-400" />
            </div>
            <div className="text-3xl font-bold text-red-400">
              {stats.offline}
            </div>
            <div className="text-xs text-gray-500 mt-1">Down instances</div>
          </div>
        </div>

        {/* Performance Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Avg Memory Usage</span>
              <HardDrive className="w-5 h-5 text-gray-400" />
            </div>
            <div
              className={`text-3xl font-bold ${getMetricColor(stats.avgMemory)}`}
            >
              {stats.avgMemory}%
            </div>
            <div className="text-xs text-gray-500 mt-1">
              Across all instances
            </div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Avg CPU Usage</span>
              <Activity className="w-5 h-5 text-gray-400" />
            </div>
            <div
              className={`text-3xl font-bold ${getMetricColor(stats.avgCpu)}`}
            >
              {stats.avgCpu}%
            </div>
            <div className="text-xs text-gray-500 mt-1">
              System load average
            </div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Total Connections</span>
              <Zap className="w-5 h-5 text-gray-400" />
            </div>
            <div className="text-3xl font-bold">
              {stats.totalConnections.toLocaleString()}
            </div>
            <div className="text-xs text-gray-500 mt-1">Active connections</div>
          </div>

          <div className="bg-gray-900 border border-gray-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-gray-400 text-sm">Avg Latency</span>
              <Clock className="w-5 h-5 text-gray-400" />
            </div>
            <div className="text-3xl font-bold">{stats.avgLatency}ms</div>
            <div className="text-xs text-gray-500 mt-1">Response time</div>
          </div>
        </div>

        {/* Filters and Search */}
        <div className="bg-gray-900 border border-gray-800 rounded-lg p-6 mb-6">
          <div className="flex flex-col lg:flex-row gap-4">
            <div className="flex-1">
              <input
                type="text"
                placeholder="Search databases..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500"
              />
            </div>
            <div className="flex gap-2">
              <select
                value={selectedDb}
                onChange={(e) => setSelectedDb(e.target.value)}
                className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500"
              >
                <option value="all">All Types</option>
                <option value="MySQL">MySQL</option>
                <option value="PostgreSQL">PostgreSQL</option>
              </select>
              <select
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500"
              >
                <option value="all">All Status</option>
                <option value="online">Online</option>
                <option value="warning">Warning</option>
                <option value="offline">Offline</option>
              </select>
            </div>
          </div>
        </div>

        {/* Database Instances */}
        <div className="space-y-4">
          {filteredDatabases.map((db) => (
            <div
              key={db.id}
              className="bg-gray-900 border border-gray-800 rounded-lg p-6 hover:border-gray-700 transition-colors"
            >
              <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4 mb-4">
                <div className="flex items-start gap-4">
                  <div
                    className={`p-3 rounded-lg ${db.type === "MySQL" ? "bg-blue-500/10" : "bg-purple-500/10"}`}
                  >
                    <Database
                      className={`w-6 h-6 ${db.type === "MySQL" ? "text-blue-400" : "text-purple-400"}`}
                    />
                  </div>
                  <div>
                    <div className="flex items-center gap-3 mb-1">
                      <h3 className="text-lg font-semibold">{db.name}</h3>
                      <span className={getStatusBadge(db.status)}>
                        {db.status.charAt(0).toUpperCase() + db.status.slice(1)}
                      </span>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-gray-400">
                      <span
                        className={
                          db.type === "MySQL"
                            ? "text-blue-400"
                            : "text-purple-400"
                        }
                      >
                        {db.type} {db.version}
                      </span>
                      <span>•</span>
                      <span>{db.host}</span>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button className="px-4 py-2 bg-blue-500 hover:bg-blue-600 rounded-lg text-sm font-medium transition-colors">
                    Connect
                  </button>
                  <button className="px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg text-sm font-medium transition-colors">
                    Details
                  </button>
                </div>
              </div>

              {/* Metrics Grid */}
              <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
                <div>
                  <div className="text-xs text-gray-400 mb-1">Memory</div>
                  <div
                    className={`text-lg font-semibold ${getMetricColor(db.memory)}`}
                  >
                    {db.memory}%
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">CPU</div>
                  <div
                    className={`text-lg font-semibold ${getMetricColor(db.cpu)}`}
                  >
                    {db.cpu}%
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Connections</div>
                  <div className="text-lg font-semibold">
                    {db.connections.current}/{db.connections.max}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Storage</div>
                  <div className="text-lg font-semibold">
                    {db.storage.used}/{db.storage.total}
                    {db.storage.unit}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Latency</div>
                  <div className="text-lg font-semibold">{db.latency}ms</div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Queries</div>
                  <div className="text-lg font-semibold">
                    {db.queries.total.toLocaleString()}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Slow Queries</div>
                  <div
                    className={`text-lg font-semibold ${db.queries.slow > 50 ? "text-red-400" : db.queries.slow > 20 ? "text-yellow-400" : "text-green-400"}`}
                  >
                    {db.queries.slow}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-gray-400 mb-1">Uptime</div>
                  <div className="text-lg font-semibold">{db.uptime}</div>
                </div>
              </div>

              {/* Additional Info */}
              <div className="flex flex-wrap gap-4 mt-4 pt-4 border-t border-gray-800 text-sm text-gray-400">
                <div className="flex items-center gap-2">
                  <Database className="w-4 h-4" />
                  <span>{db.databases} databases</span>
                </div>
                <div className="flex items-center gap-2">
                  <Server className="w-4 h-4" />
                  <span>{db.tables} tables</span>
                </div>
                <div className="flex items-center gap-2">
                  <Clock className="w-4 h-4" />
                  <span>Last backup: {db.lastBackup}</span>
                </div>
              </div>
            </div>
          ))}
        </div>

        {filteredDatabases.length === 0 && (
          <div className="bg-gray-900 border border-gray-800 rounded-lg p-12 text-center">
            <Database className="w-12 h-12 text-gray-600 mx-auto mb-4" />
            <p className="text-gray-400">
              No databases found matching your filters
            </p>
          </div>
        )}
      </div>
    </Layout>
  );
}
