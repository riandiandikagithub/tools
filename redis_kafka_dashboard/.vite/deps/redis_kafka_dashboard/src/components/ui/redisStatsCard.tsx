import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { RedisStats } from "@/models/generalModel";

export function RedisStatsCard({ stats }: { stats: RedisStats }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Redis Stats</CardTitle>
      </CardHeader>

      <CardContent className="grid grid-cols-2 gap-4">
        <p>Ops/sec: {stats.opsPerSec}</p>
        <p>Memory Used: {stats.memoryUsed}</p>
        <p>Memory Peak: {stats.memoryPeak}</p>
        <p>Keyspace Hits: {stats.keyHits}</p>
        <p>Keyspace Misses: {stats.keyMisses}</p>
        <p>Pub/Sub Channels: {stats.pubsubChannels}</p>
      </CardContent>
    </Card>
  );
}
