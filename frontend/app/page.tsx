import { getWebsites } from "@/lib/api";
import { WebsiteCard } from "@/components/WebsiteCard";

export const dynamic = "force-dynamic";

export default async function HomePage() {
  let websites;
  try {
    websites = await getWebsites();
  } catch {
    return (
      <div className="container">
        <p className="empty">Could not reach the API. Is the backend running?</p>
      </div>
    );
  }

  const up = websites.filter((w) => w.latest_check?.status === "up");
  const down = websites.filter(
    (w) => !w.latest_check || w.latest_check.status === "down"
  );

  return (
    <div className="container">
      <p className="page-subtitle">Status of {websites.length} popular websites</p>

      <section className="section">
        <h2>Currently Up ({up.length})</h2>
        {up.length > 0 ? (
          <div className="grid">
            {up.map((website) => (
              <WebsiteCard key={website.id} website={website} />
            ))}
          </div>
        ) : (
          <p className="muted">No websites are currently up.</p>
        )}
      </section>

      <section className="section">
        <h2>Currently Down ({down.length})</h2>
        {down.length > 0 ? (
          <div className="grid">
            {down.map((website) => (
              <WebsiteCard key={website.id} website={website} />
            ))}
          </div>
        ) : (
          <p className="muted">All websites are up.</p>
        )}
      </section>
    </div>
  );
}
