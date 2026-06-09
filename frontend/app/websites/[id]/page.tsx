import Link from "next/link";
import { notFound } from "next/navigation";
import { calcUptimePercent, getWebsite, getWebsiteHistory } from "@/lib/api";
import { HistoryTable } from "@/components/HistoryTable";
import { StatusBadge } from "@/components/StatusBadge";

interface PageProps {
  params: Promise<{ id: string }>;
}

export const dynamic = "force-dynamic";

export default async function WebsiteDetailPage({ params }: PageProps) {
  const { id } = await params;

  let website;
  let checks;
  try {
    [website, checks] = await Promise.all([
      getWebsite(id),
      getWebsiteHistory(id, 24),
    ]);
  } catch {
    notFound();
  }

  const status = website.latest_check?.status ?? "down";
  const uptime = calcUptimePercent(checks);

  return (
    <div className="container">
      <Link href="/" className="back-link">
        ← Back to dashboard
      </Link>

      <div className="detail-header">
        <h1>{website.name}</h1>
        <p className="muted">{website.url}</p>
        <StatusBadge status={status} />
      </div>

      <div className="stats">
        <div className="stat">
          <div className="stat-label">24h Uptime</div>
          <div className="stat-value">{uptime}%</div>
        </div>
        <div className="stat">
          <div className="stat-label">Checks (24h)</div>
          <div className="stat-value">{checks.length}</div>
        </div>
        {website.latest_check?.response_time_ms != null && (
          <div className="stat">
            <div className="stat-label">Latest Response</div>
            <div className="stat-value">
              {website.latest_check.response_time_ms} ms
            </div>
          </div>
        )}
      </div>

      <section>
        <h2>Check History (last 24 hours)</h2>
        <HistoryTable checks={checks} />
      </section>
    </div>
  );
}
