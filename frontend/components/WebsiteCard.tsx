import Link from "next/link";
import type { Website } from "@/lib/api";
import { StatusBadge } from "./StatusBadge";

interface WebsiteCardProps {
  website: Website;
}

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString();
}

export function WebsiteCard({ website }: WebsiteCardProps) {
  const check = website.latest_check;
  const status = check?.status ?? "down";

  return (
    <Link href={`/websites/${website.id}`} className="card">
      <div className="card-header">
        <h3>{website.name}</h3>
        <StatusBadge status={status} />
      </div>
      <p className="muted">{website.url}</p>
      {check ? (
        <div className="card-meta">
          <span>
            Response:{" "}
            {check.response_time_ms != null ? `${check.response_time_ms} ms` : "—"}
          </span>
          <span>Last checked: {formatDateTime(check.checked_at)}</span>
        </div>
      ) : (
        <p className="muted">No checks yet</p>
      )}
    </Link>
  );
}
