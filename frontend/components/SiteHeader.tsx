import Link from "next/link";
import { WebsiteSearch } from "./WebsiteSearch";

export function SiteHeader() {
  return (
    <header className="site-header">
      <div className="container site-header-inner">
        <Link href="/" className="site-title">
          <h1>Server Uptime Checker</h1>
        </Link>
        <WebsiteSearch />
      </div>
    </header>
  );
}
