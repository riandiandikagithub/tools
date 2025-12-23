import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import {
  Database,
  MessageSquare,
  FileText,
  CheckCircle,
  AlertCircle,
  Settings,
  Server,
  Save,
} from "lucide-react";
import { motion } from "framer-motion";
import { configService, type Configs } from "@/services/configService";
import { useNavigate } from "react-router-dom";

// ==================== DEFAULT CONFIGS ====================
const defaultRedisConfig = `# Redis Configuration
redis:
  mode: "cluster"
  nodes:
    - host: "redis-node-1"
      port: 6379
      password: "your-password"
    - host: "redis-node-2"
      port: 6379
      password: "your-password"
    - host: "redis-node-3"
      port: 6379
      password: "your-password"

  monitoring:
    interval: 30
    metrics:
      - memory_usage
      - cpu_usage
      - connected_clients
      - commands_per_sec`;

const defaultKafkaConfig = `# Kafka Configuration
kafka:
  brokers:
    - "kafka-broker-1:9092"
    - "kafka-broker-2:9092"
    - "kafka-broker-3:9092"

  security:
    protocol: "SASL_SSL"
    sasl_mechanism: "PLAIN"
    username: "your-username"
    password: "your-password"

  monitoring:
    interval: 30
    topics:
      - "orders"
      - "payments"
      - "notifications"

    consumer_groups:
      - "order-processor"
      - "payment-handler"
      - "notification-service"

    metrics:
      - broker_status
      - topic_partitions
      - consumer_lag
      - message_rate`;

const defaultPostgresConfig = `# PostgreSQL Configuration
postgresql:
  databases:
    - name: "production"
      host: "prod-postgres-01.example.com"
      port: 5432
      database: "main_db"
      username: "postgres_user"
      password: "your-password"
      ssl_mode: "require"

      pool:
        min_connections: 5
        max_connections: 20
        max_idle_time: 300

      monitoring:
        enabled: true
        interval: 30
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1000

    - name: "analytics"
      host: "analytics-postgres-01.example.com"
      port: 5432
      database: "analytics_db"
      username: "analytics_user"
      password: "your-password"
      ssl_mode: "require"

      pool:
        min_connections: 3
        max_connections: 15
        max_idle_time: 300

      monitoring:
        enabled: true
        interval: 60
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 2000

    - name: "staging"
      host: "staging-postgres-01.example.com"
      port: 5432
      database: "staging_db"
      username: "staging_user"
      password: "your-password"
      ssl_mode: "prefer"

      pool:
        min_connections: 2
        max_connections: 10
        max_idle_time: 300

      monitoring:
        enabled: true
        interval: 60
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1500

  monitoring:
    metrics:
      - connection_count
      - active_queries
      - database_size
      - table_sizes
      - index_usage
      - cache_hit_ratio
      - transaction_rate
      - replication_lag
      - deadlocks
      - slow_queries

    health_check:
      enabled: true
      interval: 10
      timeout: 5

    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000
      replication_lag_seconds: 10
      cache_hit_ratio_min: 90
      disk_usage_percent: 85`;

const defaultMySQLConfig = `# MySQL Configuration
mysql:
  databases:
    - name: "production"
      host: "prod-mysql-01.example.com"
      port: 3306
      database: "main_db"
      username: "mysql_user"
      password: "your-password"
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

    - name: "reporting"
      host: "report-mysql-01.example.com"
      port: 3306
      database: "reports_db"
      username: "report_user"
      password: "your-password"
      charset: "utf8mb4"

      pool:
        min_connections: 3
        max_connections: 15
        max_idle_time: 300
        connection_timeout: 10

      monitoring:
        enabled: true
        interval: 60
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 2000

    - name: "development"
      host: "dev-mysql-01.example.com"
      port: 3306
      database: "dev_db"
      username: "dev_user"
      password: "your-password"
      charset: "utf8mb4"

      pool:
        min_connections: 2
        max_connections: 10
        max_idle_time: 300
        connection_timeout: 10

      monitoring:
        enabled: true
        interval: 60
        track_activity: true
        log_slow_queries: true
        slow_query_threshold: 1500

  monitoring:
    metrics:
      - connection_count
      - query_cache_hit_ratio
      - slow_queries
      - table_locks
      - innodb_buffer_pool
      - threads_running
      - bytes_sent_received
      - aborted_connections
      - table_size
      - replication_status

    health_check:
      enabled: true
      interval: 10
      timeout: 5

    alerts:
      max_connections_percent: 80
      slow_query_threshold: 1000
      query_cache_hit_ratio_min: 85
      replication_lag_seconds: 10
      disk_usage_percent: 85
      threads_running_max: 100`;

// ==================== COMPONENT ====================
export function ConfigurationPage() {
  const [redisConfig, setRedisConfig] = useState(defaultRedisConfig);
  const [kafkaConfig, setKafkaConfig] = useState(defaultKafkaConfig);
  const [postgresConfig, setPostgresConfig] = useState(defaultPostgresConfig);
  const [mysqlConfig, setMySQLConfig] = useState(defaultMySQLConfig);
  const [isSaving, setIsSaving] = useState(false);
  const [saveStatus, setSaveStatus] = useState("idle");
  const [activeTab, setActiveTab] = useState("redis");
  const navigate = useNavigate();

  useEffect(() => {
    const loadConfigs = async () => {
      try {
        const configs = await configService.getAll();
        setRedisConfig(configs.redis ?? defaultRedisConfig);
        setKafkaConfig(configs.kafka ?? defaultKafkaConfig);
        setPostgresConfig(configs.postgresql ?? defaultPostgresConfig);
        setMySQLConfig(configs.mysql ?? defaultMySQLConfig);
      } catch (error) {
        console.error("Failed to load configs", error);
      }
    };
    loadConfigs();
  }, []);

  const handleSave = async () => {
    setIsSaving(true);
    setSaveStatus("idle");

    try {
      const configs: Configs = {
        redis: redisConfig,
        kafka: kafkaConfig,
        postgresql: postgresConfig,
        mysql: mysqlConfig,
      };

      await configService.saveAll(configs);
      setSaveStatus("success");

      setTimeout(() => navigate("/dashboard"), 1500);
    } catch (error) {
      console.error("Save failed:", error);
      setSaveStatus("error");
    } finally {
      setIsSaving(false);
    }
  };
  const validateConfig = (config: string) => {
    try {
      const lines = config.split("\n");
      return lines.some((line) => line.trim().includes(":"));
    } catch {
      return false;
    }
  };

  const isRedisValid = validateConfig(redisConfig);
  const isKafkaValid = validateConfig(kafkaConfig);
  const isPostgresValid = validateConfig(postgresConfig);
  const isMySQLValid = validateConfig(mysqlConfig);

  const allConfigsValid =
    isRedisValid && isKafkaValid && isPostgresValid && isMySQLValid;

  const renderConfigCard = (
    title: string,
    Icon: string,
    value: string,
    setValue: (v: string) => void,
    desc: string,
  ) => (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.18 }}
    >
      <Card className="bg-[#121418] border border-white/5 shadow-xl/40 backdrop-blur-md border border-border/30 shadow">
        <CardHeader>
          <CardTitle className="flex items-center gap-3 text-lg">
            <Icon className="h-5 w-5 text-primary" />
            {title}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label
              htmlFor={`${title}-config`}
              className="flex items-center gap-2"
            >
              <FileText className="h-4 w-4 text-primary" /> YAML Configuration
            </Label>
            <Textarea
              id={`${title}-config`}
              value={value}
              onChange={(e) => setValue(e.target.value)}
              className="font-mono text-sm min-h-[360px] bg-background/60"
              placeholder="Enter configuration in YAML format..."
            />
          </div>
          <p className="text-sm text-muted-foreground leading-relaxed">
            {desc}
          </p>
        </CardContent>
      </Card>
    </motion.div>
  );

  return (
    <div className="min-h-screen bg-[#0d0f12] text-gray-200 p-6">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <div className="flex items-center gap-3 mb-2">
            <Settings className="h-8 w-8 text-primary" />
            <h1 className="text-3xl font-bold tracking-tight">
              Configuration Setup
            </h1>
          </div>
          <p className="text-muted-foreground text-sm">
            Manage your Redis, Kafka, PostgreSQL, and MySQL setup.
          </p>
        </motion.div>

        {saveStatus === "success" && (
          <Alert className="mb-6 border-green-500/20 bg-green-500/10">
            <CheckCircle className="h-4 w-4 text-green-500" />
            <AlertDescription className="text-green-500">
              Saved successfully. Redirecting...
            </AlertDescription>
          </Alert>
        )}

        {saveStatus === "error" && (
          <Alert className="mb-6 border-red-500/20 bg-red-500/10">
            <AlertCircle className="h-4 w-4 text-red-500" />
            <AlertDescription className="text-red-500">
              Save failed. Try again.
            </AlertDescription>
          </Alert>
        )}

        <div className="flex gap-6">
          <div className="w-64 space-y-3">
            {[
              {
                key: "redis",
                label: "Redis",
                icon: Database,
                valid: isRedisValid,
              },
              {
                key: "kafka",
                label: "Kafka",
                icon: MessageSquare,
                valid: isKafkaValid,
              },
              {
                key: "postgresql",
                label: "PostgreSQL",
                icon: Server,
                valid: isPostgresValid,
              },
              {
                key: "mysql",
                label: "MySQL",
                icon: Database,
                valid: isMySQLValid,
              },
            ].map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`w-full flex items-center gap-3 px-4 py-3 text-sm rounded-xl shadow-sm transition-all border backdrop-blur-md ${
                  activeTab === tab.key
                    ? "bg-green-600 text-white border-primary"
                    : "bg-white/5/30 hover:bg-white/5/60 border-border/20"
                }`}
              >
                <tab.icon className="h-4 w-4" />
                <span className="flex-1 text-left">{tab.label}</span>
                <Badge
                  variant={tab.valid ? "default" : "destructive"}
                  className="ml-auto"
                >
                  {tab.valid ? "✓" : "✗"}
                </Badge>
              </button>
            ))}
          </div>

          <div className="flex-1 space-y-4">
            {activeTab === "redis" &&
              renderConfigCard(
                "Redis Configuration",
                Database,
                redisConfig,
                setRedisConfig,
                "Configure Redis clusters, node hosts, passwords, and monitoring options.",
              )}

            {activeTab === "kafka" &&
              renderConfigCard(
                "Kafka Configuration",
                MessageSquare,
                kafkaConfig,
                setKafkaConfig,
                "Define brokers, topics, consumer groups, and security settings.",
              )}

            {activeTab === "postgresql" &&
              renderConfigCard(
                "PostgreSQL Configuration",
                Server,
                postgresConfig,
                setPostgresConfig,
                "Manage multiple PostgreSQL instances, pooling, SSL, and health checks.",
              )}

            {activeTab === "mysql" &&
              renderConfigCard(
                "MySQL Configuration",
                Database,
                mysqlConfig,
                setMySQLConfig,
                "Configure MySQL nodes, charset options, replication, and performance settings.",
              )}
          </div>
        </div>

        <div className="flex justify-end gap-4 mt-8">
          <Button
            variant="outline"
            onClick={() => navigate("/dashboard")}
            className="rounded-xl"
          >
            Cancel
          </Button>
          <Button
            onClick={handleSave}
            disabled={isSaving || !allConfigsValid}
            className="gap-2 rounded-xl shadow-md"
          >
            <Save className="h-4 w-4" />
            {isSaving ? "Saving..." : "Save All Configurations"}
          </Button>
        </div>
      </div>
    </div>
  );
}
