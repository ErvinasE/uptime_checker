import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Server Uptime Checker",
  description: "Monitor popular website uptime",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
