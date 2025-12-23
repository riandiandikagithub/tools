import { useState, useEffect } from "react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { monitoringService } from "@/services/monitoringService";

import { BrokersCard } from "@/components/ui/brokersCard";
import { TopicsCard } from "@/components/ui/topicsCard";
import { ConsumerGroupsCard } from "@/components/ui/consumerGroupsCard";
import { ClusterInfoCard } from "@/components/ui/clusterInfoCard";

import { RedisInfoCard } from "@/components/ui/redisInfoCard";
import { RedisStatsCard } from "@/components/ui/redisStatsCard";

import { MySQLInfoCard } from "@/components/ui/mysqlInfoCard";
import { MySQLOperationsCard } from "@/components/ui/mysqlinfoOperationsCard";

import { PostgresInfoCard } from "@/components/ui/postgreInfoCard";
import { PostgresStatsCard } from "@/components/ui/postgreStatsCard";

import { Button } from "@/components/ui/button";

export default function MonitoringDashboard() {
  const [activeTab, setActiveTab] = useState("kafka");
  const [data, setData] = useState<any>(null);
  const [running, setRunning] = useState(false);

  const serviceMap: unknown = {
    kafka: fetchKafkaMetrics,
    redis: fetchRedisMetrics,
    mysql: fetchMySQLMetrics,
    postgres: fetchPostgresMetrics,
  };

  const load = async () => {
    const fn = serviceMap[activeTab];
    const res = await fn();
    setData(res);
  };

  useEffect(() => {
    if (!running) return;
    load();
    const id = setInterval(load, 3000);
    return () => clearInterval(id);
  }, [running, activeTab]);

  return (
    <div className="p-6 space-y-6">
      <Tabs defaultValue="kafka" onValueChange={(v) => setActiveTab(v)}>
        <TabsList className="grid grid-cols-4 w-full">
          <TabsTrigger value="kafka">Kafka</TabsTrigger>
          <TabsTrigger value="redis">Redis</TabsTrigger>
          <TabsTrigger value="mysql">MySQL</TabsTrigger>
          <TabsTrigger value="postgres">PostgreSQL</TabsTrigger>
        </TabsList>

        <div className="mt-4">
          <Button onClick={() => setRunning(!running)}>
            {running ? "Stop Monitoring" : "Start Monitoring"}
          </Button>
        </div>

        {/* Kafka */}
        <TabsContent value="kafka">
          {data && (
            <div className="space-y-6">
              <BrokersCard brokers={data.brokers} />
              <TopicsCard topics={data.topics} />
              <ConsumerGroupsCard groups={data.consumerGroups} />
              <ClusterInfoCard info={data.clusterInfo} />
            </div>
          )}
        </TabsContent>

        {/* Redis */}
        <TabsContent value="redis">
          {data && (
            <div className="space-y-6">
              <RedisInfoCard info={data.info} />
              <RedisStatsCard stats={data.stats} />
            </div>
          )}
        </TabsContent>

        {/* MySQL */}
        <TabsContent value="mysql">
          {data && (
            <div className="space-y-6">
              <MySQLInfoCard info={data.info} />
              <MySQLOperationsCard ops={data.ops} />
            </div>
          )}
        </TabsContent>

        {/* PostgreSQL */}
        <TabsContent value="postgres">
          {data && (
            <div className="space-y-6">
              <PostgresInfoCard info={data.info} />
              <PostgresStatsCard stats={data.stats} />
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
