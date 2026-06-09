import type { Website } from "./api";

export function filterWebsites(websites: Website[], query: string): Website[] {
  const trimmed = query.trim().toLowerCase();
  if (!trimmed) {
    return [];
  }

  return websites
    .filter(
      (site) =>
        site.name.toLowerCase().includes(trimmed) ||
        site.url.toLowerCase().includes(trimmed)
    )
    .slice(0, 8);
}
