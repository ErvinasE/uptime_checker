import type { CheckStatus } from "@/lib/api";

interface StatusBadgeProps {
  status: CheckStatus;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const isUp = status === "up";
  return (
    <span className={`badge ${isUp ? "badge-up" : "badge-down"}`}>
      {isUp ? "Up" : "Down"}
    </span>
  );
}
