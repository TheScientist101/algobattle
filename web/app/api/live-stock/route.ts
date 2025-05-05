import { NextRequest, NextResponse } from "next/server";

/**
 * GET /api/live-stock
 * -------------------
 * A proxy route that fetches live stock data from an external service.
 * - Requires an API key passed via the `Authorization` header.
 * - Forwards the request to the external stock API.
 * - Returns the response JSON or an appropriate error.
 */
export async function GET(req: NextRequest) {
  const apiKey = req.headers.get("authorization");

  if (!apiKey) {
    return new NextResponse("Missing API key", { status: 400 });
  }

  const res = await fetch("http://algobattle.thescientist101.hackclub.app/live_stock_data", {
    method: "GET",
    headers: {
      Authorization: apiKey,
    },
  });

  if (!res.ok) {
    return new NextResponse("Failed to fetch stock data", { status: res.status });
  }

  const data = await res.json();
  return NextResponse.json(data);
}
