import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { ClusterInfo } from "@/models/kafkaModel";

export function ClusterInfoCard({ info }: { info: ClusterInfo }) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Cluster Info</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-3">
          <p>Total Brokers: {info.totalBrokers}</p>
          <p>Total Topics: {info.totalTopics}</p>
          <p>Total Partitions: {info.totalPartitions}</p>
          <p>Consumer Groups: {info.totalConsumerGroups}</p>
          <p>Message Rate: {info.totalMessageRate}/s</p>
          <p>Disk Usage: {info.totalDiskUsage}</p>
          <p>Avg Replication Factor: {info.avgReplicationFactor}</p>
        </div>
      </CardContent>
    </Card>
  );
}
