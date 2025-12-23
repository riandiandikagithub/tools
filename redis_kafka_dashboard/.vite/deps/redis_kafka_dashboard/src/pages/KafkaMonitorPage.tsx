import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { 
  MessageSquare, 
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
  Database,
  TrendingUp
} from 'lucide-react';
import { Layout } from '@/components/Layout';

// Mock Kafka data
const kafkaBrokers = [
  {
    id: 'kafka-001',
    name: 'Kafka Broker 1',
    host: '10.0.2.10:9092',
    status: 'online',
    isController: true,
    diskUsage: 65,
    diskTotal: '500GB',
    diskUsed: '325GB',
    cpuUsage: 45,
    memoryUsage: 78,
    networkIn: 120,
    networkOut: 95,
    messageRate: 15000,
    uptime: '25d 8h 15m',
    lastUpdate: '10:04 PM'
  },
  {
    id: 'kafka-002',
    name: 'Kafka Broker 2',
    host: '10.0.2.11:9092',
    status: 'online',
    isController: false,
    diskUsage: 72,
    diskTotal: '500GB',
    diskUsed: '360GB',
    cpuUsage: 52,
    memoryUsage: 82,
    networkIn: 98,
    networkOut: 87,
    messageRate: 12500,
    uptime: '25d 8h 12m',
    lastUpdate: '10:04 PM'
  },
  {
    id: 'kafka-003',
    name: 'Kafka Broker 3',
    host: '10.0.2.12:9092',
    status: 'warning',
    isController: false,
    diskUsage: 88,
    diskTotal: '500GB',
    diskUsed: '440GB',
    cpuUsage: 75,
    memoryUsage: 92,
    networkIn: 145,
    networkOut: 132,
    messageRate: 18000,
    uptime: '20d 4h 30m',
    lastUpdate: '10:03 PM'
  }
];

const kafkaTopics = [
  {
    name: 'orders',
    partitions: 12,
    replicas: 3,
    underReplicated: 0,
    messageRate: 8500,
    size: '2.4GB',
    retention: '7d',
    status: 'healthy'
  },
  {
    name: 'payments',
    partitions: 8,
    replicas: 3,
    underReplicated: 0,
    messageRate: 5200,
    size: '1.8GB',
    retention: '30d',
    status: 'healthy'
  },
  {
    name: 'notifications',
    partitions: 6,
    replicas: 3,
    underReplicated: 2,
    messageRate: 12000,
    size: '950MB',
    retention: '3d',
    status: 'warning'
  },
  {
    name: 'logs',
    partitions: 24,
    replicas: 2,
    underReplicated: 0,
    messageRate: 25000,
    size: '8.2GB',
    retention: '1d',
    status: 'healthy'
  }
];

const consumerGroups = [
  {
    name: 'order-processor',
    topic: 'orders',
    consumers: 3,
    lag: 150,
    status: 'healthy'
  },
  {
    name: 'payment-handler',
    topic: 'payments',
    consumers: 2,
    lag: 45,
    status: 'healthy'
  },
  {
    name: 'notification-service',
    topic: 'notifications',
    consumers: 1,
    lag: 15000,
    status: 'warning'
  },
  {
    name: 'log-aggregator',
    topic: 'logs',
    consumers: 4,
    lag: 2500,
    status: 'healthy'
  }
];

const clusterInfo = {
  totalBrokers: 3,
  totalTopics: 15,
  totalPartitions: 180,
  totalConsumerGroups: 8,
  totalMessageRate: 89000,
  totalDiskUsage: '4.2TB',
  avgReplicationFactor: 2.8
};

export function KafkaMonitorPage() {
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
      case 'healthy':
        return <CheckCircle className="h-4 w-4 text-green-400" />;
      case 'warning':
        return <AlertTriangle className="h-4 w-4 text-yellow-400" />;
      case 'offline':
      case 'error':
        return <XCircle className="h-4 w-4 text-red-400" />;
      default:
        return <XCircle className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusBadge = (status: string) => {
    const variants = {
      online: 'default',
      healthy: 'default',
      warning: 'secondary', 
      offline: 'destructive',
      error: 'destructive'
    } as const;
    
    return (
      <Badge variant={variants[status as keyof typeof variants] || 'secondary'} className="capitalize">
        {status}
      </Badge>
    );
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex items-center gap-3">
          <MessageSquare className="h-8 w-8 text-blue-400" />
          <div>
            <h1 className="text-3xl font-bold">Kafka Monitor</h1>
            <p className="text-muted-foreground">Monitor Kafka brokers, topics, and consumer groups</p>
          </div>
        </div>

        <Tabs defaultValue="overview" className="space-y-6">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="brokers">Brokers</TabsTrigger>
            <TabsTrigger value="topics">Topics</TabsTrigger>
            <TabsTrigger value="consumers">Consumer Groups</TabsTrigger>
          </TabsList>

          <TabsContent value="overview" className="space-y-6">
            {/* Cluster Summary */}
            <div className="dashboard-grid">
              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Total Brokers</CardTitle>
                  <Server className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">{clusterInfo.totalBrokers}</div>
                  <p className="metric-label">Active brokers</p>
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Topics</CardTitle>
                  <Database className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">{clusterInfo.totalTopics}</div>
                  <p className="metric-label">{clusterInfo.totalPartitions} partitions</p>
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Message Rate</CardTitle>
                  <Zap className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">{clusterInfo.totalMessageRate.toLocaleString()}</div>
                  <p className="metric-label">Messages per second</p>
                </CardContent>
              </Card>

              <Card className="metric-card">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Consumer Groups</CardTitle>
                  <Users className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="metric-value">{clusterInfo.totalConsumerGroups}</div>
                  <p className="metric-label">Active groups</p>
                </CardContent>
              </Card>
            </div>

            {/* Quick Status */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {kafkaBrokers.map((broker) => (
                <Card key={broker.id} className="p-4">
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center gap-2">
                      <MessageSquare className="h-4 w-4 text-blue-400" />
                      <span className="font-medium text-sm">{broker.name}</span>
                      {broker.isController && (
                        <Badge variant="outline" className="text-xs">Controller</Badge>
                      )}
                    </div>
                    {getStatusIcon(broker.status)}
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">Disk</span>
                      <span>{broker.diskUsage}%</span>
                    </div>
                    <Progress value={broker.diskUsage} className="h-1" />
                    <div className="flex justify-between text-xs">
                      <span className="text-muted-foreground">Memory</span>
                      <span>{broker.memoryUsage}%</span>
                    </div>
                    <Progress value={broker.memoryUsage} className="h-1" />
                  </div>
                </Card>
              ))}
            </div>
          </TabsContent>

          <TabsContent value="brokers" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  Kafka Brokers
                </CardTitle>
                <div className="flex gap-4 mt-4">
                  <div className="flex-1">
                    <Input
                      placeholder="Search brokers..."
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
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {kafkaBrokers.map((broker) => (
                    <Card key={broker.id} className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <MessageSquare className="h-5 w-5 text-blue-400" />
                          <div>
                            <div className="font-medium flex items-center gap-2">
                              {broker.name}
                              {getStatusIcon(broker.status)}
                              {broker.isController && (
                                <Badge variant="outline">Controller</Badge>
                              )}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {broker.host}
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {getStatusBadge(broker.status)}
                          <Button variant="outline" size="sm">
                            View Details
                          </Button>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <HardDrive className="h-3 w-3" />
                            Disk Usage
                          </div>
                          <div className="text-sm font-medium">{broker.diskUsage}%</div>
                          <div className="text-xs text-muted-foreground">{broker.diskUsed}</div>
                          <Progress value={broker.diskUsage} className="h-1" />
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Cpu className="h-3 w-3" />
                            CPU
                          </div>
                          <div className="text-sm font-medium">{broker.cpuUsage}%</div>
                          <Progress value={broker.cpuUsage} className="h-1" />
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Database className="h-3 w-3" />
                            Memory
                          </div>
                          <div className="text-sm font-medium">{broker.memoryUsage}%</div>
                          <Progress value={broker.memoryUsage} className="h-1" />
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Network className="h-3 w-3" />
                            Network I/O
                          </div>
                          <div className="text-sm font-medium">{broker.networkIn}MB/s</div>
                          <div className="text-xs text-muted-foreground">↓{broker.networkIn} ↑{broker.networkOut}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Zap className="h-3 w-3" />
                            Messages/sec
                          </div>
                          <div className="text-sm font-medium">{broker.messageRate.toLocaleString()}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="flex items-center gap-1 text-xs text-muted-foreground">
                            <Clock className="h-3 w-3" />
                            Uptime
                          </div>
                          <div className="text-sm font-medium">{broker.uptime}</div>
                          <div className="text-xs text-muted-foreground">{broker.lastUpdate}</div>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="topics" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Kafka Topics
                </CardTitle>
                <div className="flex gap-4 mt-4">
                  <div className="flex-1">
                    <Input
                      placeholder="Search topics..."
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
                      <SelectItem value="healthy">Healthy</SelectItem>
                      <SelectItem value="warning">Warning</SelectItem>
                      <SelectItem value="error">Error</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {kafkaTopics.map((topic) => (
                    <Card key={topic.name} className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <Database className="h-5 w-5 text-blue-400" />
                          <div>
                            <div className="font-medium flex items-center gap-2">
                              {topic.name}
                              {getStatusIcon(topic.status)}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              {topic.partitions} partitions, {topic.replicas} replicas
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {getStatusBadge(topic.status)}
                          <Button variant="outline" size="sm">
                            View Details
                          </Button>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Partitions</div>
                          <div className="text-sm font-medium">{topic.partitions}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Replicas</div>
                          <div className="text-sm font-medium">{topic.replicas}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Under-replicated</div>
                          <div className={`text-sm font-medium ${topic.underReplicated > 0 ? 'text-yellow-400' : 'text-green-400'}`}>
                            {topic.underReplicated}
                          </div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Message Rate</div>
                          <div className="text-sm font-medium">{topic.messageRate.toLocaleString()}/s</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Size</div>
                          <div className="text-sm font-medium">{topic.size}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Retention</div>
                          <div className="text-sm font-medium">{topic.retention}</div>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="consumers" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Users className="h-5 w-5" />
                  Consumer Groups
                </CardTitle>
                <div className="flex gap-4 mt-4">
                  <div className="flex-1">
                    <Input
                      placeholder="Search consumer groups..."
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
                      <SelectItem value="healthy">Healthy</SelectItem>
                      <SelectItem value="warning">Warning</SelectItem>
                      <SelectItem value="error">Error</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {consumerGroups.map((group) => (
                    <Card key={group.name} className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                          <Users className="h-5 w-5 text-blue-400" />
                          <div>
                            <div className="font-medium flex items-center gap-2">
                              {group.name}
                              {getStatusIcon(group.status)}
                            </div>
                            <div className="text-sm text-muted-foreground">
                              Topic: {group.topic}
                            </div>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          {getStatusBadge(group.status)}
                          <Button variant="outline" size="sm">
                            View Details
                          </Button>
                        </div>
                      </div>

                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Active Consumers</div>
                          <div className="text-sm font-medium">{group.consumers}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Consumer Lag</div>
                          <div className={`text-sm font-medium ${group.lag > 10000 ? 'text-yellow-400' : 'text-green-400'}`}>
                            {group.lag.toLocaleString()} messages
                          </div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Topic</div>
                          <div className="text-sm font-medium">{group.topic}</div>
                        </div>

                        <div className="space-y-1">
                          <div className="text-xs text-muted-foreground">Status</div>
                          <div className="text-sm font-medium capitalize">{group.status}</div>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </Layout>
  );
}