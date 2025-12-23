// src/lib/apiClient.ts
const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8085";

interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export async function get<T>(endpoint: string): Promise<string> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`);
  if (!response.ok) throw new Error(`GET ${endpoint}: ${response.status}`);

  const data: ApiResponse<string> = await response.json();
  if (!data.success) throw new Error(`GET ${endpoint} failed`);

  return data.data;
}

export async function postConfig(
  endpoint: string,
  yamlContent: string,
): Promise<void> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ config: yamlContent }), // ← sesuai struct Go
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`POST ${endpoint}: ${response.status} - ${errorText}`);
  }

  // Backend tidak mengembalikan body? Cukup pastikan status 2xx.
  // Jika perlu validasi success, sesuaikan.
}

export async function getApi<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`);
  if (!response.ok) throw new Error(`GET ${endpoint}: ${response.status}`);

  const data: ApiResponse<T> = await response.json();
  if (!data.success) throw new Error(`GET ${endpoint} failed`);

  return data.data;
}
