import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { RedisInfo } from "@/models/generalModel";

export function RedisInfoCard({ info }: { info: RedisInfo }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Redis Info</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        <p>Version: {info.version}</p>
        <p>Mode: {info.mode}</p>
        <p>Connected Clients: {info.connectedClients}</p>
        <p>Uptime: {info.uptime}</p>
      </CardContent>
    </Card>
  );
}
