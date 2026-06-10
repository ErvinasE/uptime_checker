import { NextRequest, NextResponse } from "next/server";

function backendUrl(): string {
  let url = process.env.API_URL || "http://localhost:8080";
  if (url.endsWith("/")) {
    url = url.slice(0, -1);
  }
  return url;
}

async function proxy(request: NextRequest, path: string[]) {
  const target = `${backendUrl()}/${path.join("/")}${request.nextUrl.search}`;

  const res = await fetch(target, {
    method: request.method,
    headers: { 
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    cache: "no-store",
    redirect: "manual", // Don't follow redirects to avoid losing POST method
  });

  // If the backend tries to redirect, we pass it through or handle it.
  // But with path.Clean and fixed proxy URL, it shouldn't redirect anymore.
  if (res.status >= 300 && res.status < 400) {
    const location = res.headers.get("location");
    return NextResponse.json({ error: "Backend redirected", location }, { status: res.status });
  }

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

export async function POST(
  request: NextRequest,
  context: { params: Promise<{ path: string[] }> }
) {
  const { path } = await context.params;
  return proxy(request, path);
}
