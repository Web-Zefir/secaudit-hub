import axios from "axios";

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

let onUnauthorized: (() => void) | null = null;

export function setUnauthorizedHandler(handler: () => void): void {
  onUnauthorized = handler;
}

api.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      const url = error.config?.url ?? "";
      const authAttempt =
        url.includes("/auth/login") ||
        url.includes("/auth/register") ||
        url.includes("/auth/demo") ||
        url.includes("/auth/me");
      if (!authAttempt) {
        onUnauthorized?.();
      }
    }
    return Promise.reject(error);
  },
);

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

export function errorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data: unknown = error.response?.data;
    if (isRecord(data) && typeof data.error === "string") {
      return data.error;
    }
  }
  return "Не удалось выполнить запрос";
}
