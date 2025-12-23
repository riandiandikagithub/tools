import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { PostgresStats } from "@/models/generalModel";

export function PostgresStatsCard({ stats }: { stats: PostgresStats }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>PostgreSQL Stats</CardTitle>
      </CardHeader>

      <CardContent className="grid grid-cols-2 gap-4">
        <p>TPS: {stats.tps}</p>
        <p>Cache Hit Ratio: {stats.cacheHitRatio}%</p>
        <p>Deadlocks: {stats.deadlocks}</p>
        <p>Conflicts: {stats.conflicts}</p>
      </CardContent>
    </Card>
  );
}
