// Read at call time so Docker runtime env is used (not build-time inlined values).
function getApiUrl(): string {
  const url = process.env["API_URL"];
  if (url) {
    return url;
  }
  // Default for local dev without Docker.
  return "http://localhost:8080";
}

export type CheckStatus = "up" | "down";

export interface Check {
  id: number;
  website_id: number;
  status: CheckStatus;
  response_time_ms: number | null;
  error_message?: string | null;
  checked_at: string;
}

export interface Website {
  id: number;
  name: string;
  url: string;
  created_at: string;
  last_manual_check_at?: string | null;
  latest_check?: Check | null;
}

async function fetchJSON<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${getApiUrl()}${path}`, {
    cache: "no-store",
    ...options,
  });
  if (!res.ok) {
    const errorBody = await res.json().catch(() => ({}));
    throw new Error(errorBody.error || `API error: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export function getWebsites(): Promise<Website[]> {
  return fetchJSON<Website[]>("/websites");
}

export function getWebsite(id: string): Promise<Website> {
  return fetchJSON<Website>(`/websites/${id}`);
}

export function getWebsiteHistory(id: string, hours = 24): Promise<Check[]> {
  return fetchJSON<Check[]>(`/websites/${id}/history?hours=${hours}`);
}

export function triggerCheck(id: number): Promise<{ status: string; result: string }> {
  // Use relative URL when called from browser to use the Next.js proxy
  const baseUrl = typeof window !== "undefined" ? "/api" : getApiUrl();
  return fetchJSON<{ status: string; result: string }>(`${baseUrl}/websites/${id}/check`, {
    method: "POST",
  });
}

export function calcUptimePercent(checks: Check[]): number {
  if (checks.length === 0) return 0;
  const upCount = checks.filter((c) => c.status === "up").length;
  return Math.round((upCount / checks.length) * 100);
}
