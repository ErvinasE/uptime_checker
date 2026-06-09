import { NextRequest, NextResponse } from "next/server";

function backendUrl(): string {
  return process.env.API_URL || "http://localhost:8080";
}

async function proxy(request: NextRequest, path: string[]) {
  const target = `${backendUrl()}/${path.join("/")}${request.nextUrl.search}`;

  const res = await fetch(target, {
    method: request.method,
    headers: { Accept: "application/json" },
    cache: "no-store",
  });

  const body = await res.text();
  return new NextResponse(body, {
    status: res.status,
    headers: { "Content-Type": "application/json" },
  });
}

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  const { path } = await context.params;
  return proxy(request, path);
}
