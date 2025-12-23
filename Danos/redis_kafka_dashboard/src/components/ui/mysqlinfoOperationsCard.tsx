import { Card, CardHeader, CardContent, CardTitle } from "@/components/ui/card";
import type { MySQLOperations } from "@/models/generalModel";

export function MySQLOperationsCard({ ops }: { ops: MySQLOperations }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Operations</CardTitle>
      </CardHeader>

      <CardContent className="grid grid-cols-2 gap-4">
        <p>Queries/sec: {ops.qps}</p>
        <p>Transactions/sec: {ops.tps}</p>
        <p>InnoDB Reads: {ops.innodbReads}</p>
        <p>InnoDB Writes: {ops.innodbWrites}</p>
      </CardContent>
    </Card>
  );
}
