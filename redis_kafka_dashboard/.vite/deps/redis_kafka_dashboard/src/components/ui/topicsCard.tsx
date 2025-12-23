import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { TopicMetrics } from "@/models/generalModel";

export function TopicsCard({ topics }: { topics: TopicMetrics[] }) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Kafka Topics</CardTitle>
      </CardHeader>

      <CardContent className="space-y-4">
        {topics.map((t) => (
          <div
            key={t.name}
            className="border p-4 rounded-lg flex justify-between items-center"
          >
            <div>
              <p className="font-semibold">{t.name}</p>
              <p className="text-sm mt-1">
                Partitions: {t.partitions} | Replicas: {t.replicas}
              </p>
            </div>

            <div className="text-right space-y-1">
              <p>Status: {t.status}</p>
              <p>Message Rate: {t.messageRate}/s</p>
              <p>Size: {t.size}</p>
              <p>Retention: {t.retention}</p>
              <p>Under-replicated: {t.underReplicated}</p>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
