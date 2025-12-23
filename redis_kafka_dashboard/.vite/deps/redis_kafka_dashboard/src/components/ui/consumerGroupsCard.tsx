import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { ConsumerGroup } from "@/models/kafkaModel";

export function ConsumerGroupsCard({ groups }: { groups: ConsumerGroup[] }) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Consumer Groups</CardTitle>
      </CardHeader>

      <CardContent className="space-y-4">
        {groups.map((g) => (
          <div
            key={g.name}
            className="border p-4 rounded-lg flex justify-between items-center"
          >
            <div>
              <p className="font-semibold">{g.name}</p>
              <p className="text-sm text-muted-foreground">Topic: {g.topic}</p>
            </div>

            <div className="text-right space-y-1">
              <p>Consumers: {g.consumers}</p>
              <p>Lag: {g.lag}</p>
              <p>Status: {g.status}</p>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
