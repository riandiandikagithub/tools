import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import { X } from "lucide-react";
import { MonitoringRedisData } from "@/models/monitoringModel";
import { Badge } from "@/components/ui/badge";

import {
  AreaChart,
  Area,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  YAxis,
} from "recharts";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  instance: MonitoringRedisData | null;
  calcReplicationLag: (m: number, r: number) => number;
}

export default function RedisInstanceDrawer({
  isOpen,
  onClose,
  instance,
  calcReplicationLag,
}: Props) {
  const [cpuHistory, setCpuHistory] = useState<{ value: number }[]>([]);
  const [memHistory, setMemHistory] = useState<{ value: number }[]>([]);

  // ESC to close
  useEffect(() => {
    const listener = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", listener);
    return () => window.removeEventListener("keydown", listener);
  }, [onClose]);

  // Auto-generate history (simulate real-time)
  useEffect(() => {
    if (!instance) return;

    const interval = setInterval(() => {
      setCpuHistory((h) => [...h.slice(-30), { value: instance.cpu_usage }]);
      setMemHistory((h) => [
        ...h.slice(-30),
        { value: Number((instance.used_memory / 1024 / 1024).toFixed(2)) },
      ]);
    }, 1500);

    return () => clearInterval(interval);
  }, [instance]);

  if (!isOpen || !instance) return null;

  const getStatusColor = (status: string) => {
    if (status === "online") return "bg-green-500";
    if (status === "warning") return "bg-yellow-500";
    return "bg-red-500";
  };

  return (
    <>
      {/* overlay */}
      <div className="fixed inset-0 bg-black/40 z-50" onClick={onClose}></div>

      {/* drawer */}
      <div
        className={`fixed top-0 right-0 h-full w-[480px] bg-background shadow-xl z-50 transform transition-transform duration-300 ${
          isOpen ? "translate-x-0" : "translate-x-full"
        }`}
      >
        {/* header */}
        <div className="p-4 border-b flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div
              className={`w-3 h-3 rounded-full ${getStatusColor(
                instance.status,
              )}`}
            ></div>
            <h2 className="text-xl text-white font-semibold">
              {instance.name}
            </h2>
          </div>

          <Button variant="ghost" size="icon" onClick={onClose}>
            <X />
          </Button>
        </div>

        {/* TABS */}
        <Tabs defaultValue="overview" className="w-full h-full flex flex-col">
          <TabsList className="w-full text-white grid grid-cols-4 rounded-none border-b">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="memory">Memory</TabsTrigger>
            <TabsTrigger value="replication">Replication</TabsTrigger>
            <TabsTrigger value="charts">Charts</TabsTrigger>
          </TabsList>

          <ScrollArea className="flex-1 p-4">
            {/* OVERVIEW */}
            <TabsContent value="overview">
              <Section title="General Info">
                <Info
                  label="Host"
                  value={`${instance.host}:${instance.port}`}
                />
                <Info label="Status" value={getStatusBadge(instance.status)} />
                <Info label="Mode" value={instance.mode} />
                <Info label="Role" value={instance.role} />
                <Info label="Uptime" value={instance.uptime_human} />
              </Section>

              <Section title="Performance">
                <Info
                  label="Connected Clients"
                  value={instance.connected_clients}
                />
                <Info
                  label="Commands / Sec"
                  value={instance.commands_per_sec}
                />
                <Info
                  label="Ops / Sec"
                  value={instance.instantaneous_ops_per_sec}
                />
                <Info label="CPU User" value={`${instance.cpu_usage}%`} />
                <Info label="CPU Sys" value={`${instance.cpu_usage_sys}%`} />
              </Section>
            </TabsContent>

            {/* MEMORY */}
            <TabsContent value="memory">
              <Section title="Memory">
                <Info label="Used Memory" value={instance.used_memory_human} />
                <Info
                  label="Peak Memory"
                  value={instance.used_memory_peak_human}
                />
                <Info
                  label="Fragmentation Ratio"
                  value={instance.memory_fragmentation_ratio}
                />
                <Info
                  label="Memory Usage"
                  value={`${instance.memory_usage_percent.toFixed(2)}%`}
                />
              </Section>

              <Section title="Network">
                <Info label="Input" value={instance.network_input_bytes} />
                <Info label="Output" value={instance.network_output_bytes} />
              </Section>
            </TabsContent>

            {/* REPLICATION */}
            <TabsContent value="replication">
              <Section title="Replication">
                <Info label="Role" value={instance.replication_role} />
                <Info
                  label="Connected Slaves"
                  value={instance.connected_slaves}
                />

                <Info
                  label="Master Repl Offset"
                  value={instance.master_repl_offset}
                />

                <Info
                  label="Replica Offset"
                  value={instance.replica_offset ?? "-"}
                />

                {instance.role === "slave" ? (
                  <Info
                    label="Replication Lag"
                    value={`${calcReplicationLag(
                      instance.master_repl_offset,
                      instance.replica_offset,
                    )} bytes`}
                  />
                ) : (
                  <Info label="Replication Lag" value="-" />
                )}
              </Section>

              <Section title="Keyspace">
                <pre className="bg-muted p-3 rounded-md text-xs overflow-auto">
                  {JSON.stringify(instance.keyspace, null, 2)}
                </pre>
              </Section>
            </TabsContent>

            {/* CHARTS */}
            <TabsContent value="charts">
              <Section title="CPU Usage (last 30 samples)">
                <div className="h-40">
                  <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={cpuHistory}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <YAxis />
                      <Tooltip />
                      <Area
                        type="monotone"
                        dataKey="value"
                        stroke="currentColor"
                        fill="currentColor"
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              </Section>

              <Section title="Memory Usage (MB)">
                <div className="h-40">
                  <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={memHistory}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <YAxis />
                      <Tooltip />
                      <Area
                        type="monotone"
                        dataKey="value"
                        stroke="currentColor"
                        fill="currentColor"
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              </Section>
            </TabsContent>
          </ScrollArea>
        </Tabs>
      </div>
    </>
  );
}

function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="mb-6">
      <h3 className="mb-3 text-lg font-semibold tracking-wide text-white drop-shadow-sm">
        {title}
      </h3>
      <div className="space-y-2">{children}</div>
    </div>
  );
}

function Info({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex justify-between text-sm py-1">
      <span className="text-white/50 font-medium tracking-wide">{label}</span>
      <span className="text-white font-semibold tracking-wide">{value}</span>
    </div>
  );
}

const getStatusBadge = (status: string) => {
  const variants = {
    online: "default",
    warning: "secondary",
    offline: "destructive",
  } as const;

  return (
    <Badge
      variant={variants[status as keyof typeof variants] || "secondary"}
      className="capitalize"
    >
      {status}
    </Badge>
  );
};

const getRoleBadge = (role: string) => {
  return (
    <Badge
      variant={role === "master" ? "default" : "outline"}
      className="capitalize"
    >
      {role}
    </Badge>
  );
};
