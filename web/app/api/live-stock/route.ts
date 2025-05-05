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
  // Extract API key from request headers
  const apiKey = req.headers.get("authorization");

  // Return 400 Bad Request if no API key is provided
  if (!apiKey) {
    return new NextResponse("Missing API key", { status: 400 });
  }

  // Forward request to the external live stock data API
  const res = await fetch("http://algobattle.thescientist101.hackclub.app/live_stock_data", {
    method: "GET",
    headers: {
      Authorization: apiKey,
    },
  });

  // If fetch fails (non-2xx response), forward the status
  if (!res.ok) {
    return new NextResponse("Failed to fetch stock data", { status: res.status });
  }

  // Return the successful response as JSON
  const data = await res.json();
  return NextResponse.json(data);
}
