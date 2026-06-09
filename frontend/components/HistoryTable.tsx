import type { Check } from "@/lib/api";
import { StatusBadge } from "./StatusBadge";

interface HistoryTableProps {
  checks: Check[];
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString();
}

function formatTime(iso: string): string {
  return new Date(iso).toLocaleTimeString();
}

export function HistoryTable({ checks }: HistoryTableProps) {
  if (checks.length === 0) {
    return <p className="muted">No checks in the last 24 hours.</p>;
  }

  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Time</th>
            <th>Response (ms)</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {checks.map((check) => (
            <tr key={check.id}>
              <td>{formatDate(check.checked_at)}</td>
              <td>{formatTime(check.checked_at)}</td>
              <td>
                {check.status === "up" && check.response_time_ms != null
                  ? check.response_time_ms
                  : "—"}
              </td>
              <td>
                <StatusBadge status={check.status} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
