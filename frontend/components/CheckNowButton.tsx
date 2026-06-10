"use client";

import { useState, useEffect } from "react";
import { triggerCheck } from "@/lib/api";
import { useRouter } from "next/navigation";

interface CheckNowButtonProps {
  websiteId: number;
  lastManualCheckAt?: string | null;
}

export function CheckNowButton({ websiteId, lastManualCheckAt }: CheckNowButtonProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [timeLeft, setTimeLeft] = useState<number>(0);
  const router = useRouter();

  useEffect(() => {
    if (!lastManualCheckAt) {
      setTimeLeft(0);
      return;
    }

    const updateTimer = () => {
      const lastCheck = new Date(lastManualCheckAt).getTime();
      const nextAvailable = lastCheck + 60 * 60 * 1000;
      const now = new Date().getTime();
      const diff = Math.max(0, Math.ceil((nextAvailable - now) / 1000));
      setTimeLeft(diff);
    };

    updateTimer();
    const interval = setInterval(updateTimer, 1000);
    return () => clearInterval(interval);
  }, [lastManualCheckAt]);

  const handleCheck = async () => {
    setLoading(true);
    setError(null);
    try {
      await triggerCheck(websiteId);
      router.refresh();
    } catch (err: any) {
      setError(err.message || "Failed to trigger check");
      // Refresh anyway in case the error was a rate limit we didn't know about
      router.refresh();
    } finally {
      setLoading(false);
    }
  };

  const formatTimeLeft = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: "0.25rem" }}>
      <button
        onClick={handleCheck}
        disabled={loading || timeLeft > 0}
        className="btn btn-primary"
      >
        {loading ? "Checking..." : timeLeft > 0 ? `Check available in ${formatTimeLeft(timeLeft)}` : "Check Now"}
      </button>
      {error && <span style={{ fontSize: "0.75rem", color: "var(--down)" }}>{error}</span>}
    </div>
  );
}
