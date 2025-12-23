import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { MySQLInfo } from "@/models/generalModel";

export function MySQLInfoCard({ info }: { info: MySQLInfo }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>MySQL Info</CardTitle>
      </CardHeader>

      <CardContent className="space-y-2">
        <p>Version: {info.version}</p>
        <p>Threads: {info.threads}</p>
        <p>Uptime: {info.uptime}</p>
        <p>Active Connections: {info.activeConnections}</p>
      </CardContent>
    </Card>
  );
}
