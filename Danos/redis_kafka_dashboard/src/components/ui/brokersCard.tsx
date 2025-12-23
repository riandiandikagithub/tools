import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { KafkaBroker } from "@/models/kafkaModel";

export function BrokersCard({ brokers }: { brokers: KafkaBroker[] }) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Kafka Brokers</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {brokers.map((b) => (
          <div
            key={b.id}
            className="border p-4 rounded-lg flex justify-between items-center"
          >
            <div>
              <p className="font-semibold">{b.name}</p>
              <p className="text-sm text-muted-foreground">{b.host}</p>
              <p className="text-sm mt-1">Status: {b.status}</p>
            </div>

            <div className="text-right space-y-1">
              <p>
                Disk: {b.diskUsed} / {b.diskTotal} ({b.diskUsage}%)
              </p>
              <p>CPU: {b.cpuUsage}%</p>
              <p>Memory: {b.memoryUsage}%</p>
              <p>
                Net In/Out: {b.networkIn} / {b.networkOut}
              </p>
              <p>Message Rate: {b.messageRate}/s</p>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
