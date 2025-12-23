import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { PostgresInfo } from "@/models/generalModel";

export function PostgresInfoCard({ info }: { info: PostgresInfo }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>PostgreSQL Info</CardTitle>
      </CardHeader>

      <CardContent className="space-y-2">
        <p>Version: {info.version}</p>
        <p>Max Connections: {info.maxConnections}</p>
        <p>Active Connections: {info.activeConnections}</p>
      </CardContent>
    </Card>
  );
}
