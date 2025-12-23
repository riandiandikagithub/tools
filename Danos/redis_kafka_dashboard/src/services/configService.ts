// src/services/configService.ts
import { get, postConfig } from "@/lib/apiClient";

const CONFIG_ENDPOINTS = {
  redis: "/api/v1/config/redis",
  kafka: "/api/v1/config/kafka",
  postgresql: "/api/v1/config/postgresql",
  mysql: "/api/v1/config/mysql",
} as const;

export type ConfigType = keyof typeof CONFIG_ENDPOINTS;

export type Configs = {
  [K in ConfigType]: string;
};

export const configService = {
  // Ambil semua konfigurasi
  async getAll(): Promise<Configs> {
    const [redis, kafka, postgresql, mysql] = await Promise.all([
      get(CONFIG_ENDPOINTS.redis),
      get(CONFIG_ENDPOINTS.kafka),
      get(CONFIG_ENDPOINTS.postgresql),
      get(CONFIG_ENDPOINTS.mysql),
    ]);

    return { redis, kafka, postgresql, mysql };
  },

  // Simpan semua konfigurasi
  async saveAll(configs: Configs): Promise<void> {
    await Promise.all([
      postConfig(CONFIG_ENDPOINTS.redis, configs.redis),
      postConfig(CONFIG_ENDPOINTS.kafka, configs.kafka),
      postConfig(CONFIG_ENDPOINTS.postgresql, configs.postgresql),
      postConfig(CONFIG_ENDPOINTS.mysql, configs.mysql),
    ]);
  },
};
