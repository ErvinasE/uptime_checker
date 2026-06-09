"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import type { Website } from "@/lib/api";
import { filterWebsites } from "@/lib/search";
import { StatusBadge } from "./StatusBadge";

const CLIENT_API_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export function WebsiteSearch() {
  const router = useRouter();
  const containerRef = useRef<HTMLDivElement>(null);
  const [websites, setWebsites] = useState<Website[]>([]);
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);

  useEffect(() => {
    fetch(`${CLIENT_API_URL}/websites`, { cache: "no-store" })
      .then((res) => {
        if (!res.ok) throw new Error("fetch failed");
        return res.json() as Promise<Website[]>;
      })
      .then(setWebsites)
      .catch(() => setWebsites([]));
  }, []);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    }

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const suggestions = filterWebsites(websites, query);
  const showSuggestions = open && query.trim().length > 0;

  function goToWebsite(id: number) {
    setQuery("");
    setOpen(false);
    router.push(`/websites/${id}`);
  }

  return (
    <div className="search" ref={containerRef}>
      <input
        type="search"
        className="search-input"
        placeholder="Search websites..."
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        aria-label="Search websites"
        aria-expanded={showSuggestions}
        aria-autocomplete="list"
        role="combobox"
      />

      {showSuggestions && (
        <ul className="search-suggestions" role="listbox">
          {suggestions.length > 0 ? (
            suggestions.map((site) => (
              <li key={site.id} role="option">
                <button
                  type="button"
                  className="search-suggestion"
                  onClick={() => goToWebsite(site.id)}
                >
                  <span className="search-suggestion-name">{site.name}</span>
                  <span className="search-suggestion-url">{site.url}</span>
                  {site.latest_check && (
                    <StatusBadge status={site.latest_check.status} />
                  )}
                </button>
              </li>
            ))
          ) : (
            <li className="search-empty" role="option">
              No matching websites
            </li>
          )}
        </ul>
      )}
    </div>
  );
}
