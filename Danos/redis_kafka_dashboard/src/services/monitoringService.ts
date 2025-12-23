import { get, getApi } from "@/lib/apiClient";
import { MonitoringRedisData } from "@/models/monitoringModel";

const MONITORING_ENDPOINTS = {
  redis: "/api/v1/monitoring/redis",
  kafka: "/api/v1/monitoring/kafka",
  postgresql: "/api/v1/monitoring/postgresql",
  mysql: "/api/v1/monitoring/mysql",
} as const;

export type MonitoringType = keyof typeof MONITORING_ENDPOINTS;

export type Monitoring = {
  [K in MonitoringType]: string;
};

// Generic API response
export interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export interface RedisResponse {
  redis: MonitoringRedisData[];
}

export const monitoringService = {
  async getAll(): Promise<Monitoring> {
    const [redis, kafka, postgresql, mysql] = await Promise.all([
      get(MONITORING_ENDPOINTS.redis),
      get(MONITORING_ENDPOINTS.kafka),
      get(MONITORING_ENDPOINTS.postgresql),
      get(MONITORING_ENDPOINTS.mysql),
    ]);

    return { redis, kafka, postgresql, mysql };
  },

  async getMonitoringRedis(): Promise<MonitoringRedisData[]> {
    const res = await getApi<RedisResponse>(MONITORING_ENDPOINTS.redis);
    return res.redis;
  },
};
