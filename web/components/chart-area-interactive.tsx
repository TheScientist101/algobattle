"use client";

import * as React from "react";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";

import { useIsMobile } from "@/hooks/use-mobile";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { useEffect, useState } from "react";
import { getBotHistory } from "@/utils/botData";
import { WithBot } from "@/utils/types";

/**
 * A responsive area chart displaying historical bot performance over time.
 *
 * Features:
 * - Fetches historical account values from Firestore using the `botId`.
 * - Supports 3 time filters: last 7 days, 30 days, and 90 days.
 * - Mobile-adaptive layout with dropdown instead of toggle group.
 * - Uses Recharts for rendering, with custom styling and tooltip formatting.
 *
 * Props:
 * - `botId`: string (from `WithBot`) – used to fetch historical data for a specific bot.
 */

// Chart config defines how to label and color the data series
const chartConfig = {
  value: {
    label: "value",
    color: "var(--primary)",
  },
} satisfies ChartConfig;

export function ChartAreaInteractive({ botId }: WithBot) {
  const isMobile = useIsMobile();

  const [timeRange, setTimeRange] = React.useState("90d"); // Active time range filter
  const [historicalData, setHistoricalData] = useState<
    { date: string; value: number }[]
  >([]);

  /**
   * useEffect: Fetches historical bot data on mount or when `botId` changes.
   */
  useEffect(() => {
    if (!botId) return;
    getBotHistory(botId).then(setHistoricalData).catch(console.error);
  }, [botId]);

  /**
   * useEffect: On mobile, default to "7d" to reduce visual clutter.
   */
  useEffect(() => {
    if (isMobile) {
      setTimeRange("7d");
    }
  }, [isMobile]);

  /**
   * Filter the full dataset based on the selected time range.
   * Only includes entries within the last 7, 30, or 90 days.
   */
  const filteredData = historicalData.filter((item) => {
    const parsedDate = new Date(item.date);
    if (isNaN(parsedDate.getTime())) return false;

    const referenceDate = new Date();
    let daysToSubtract = 90;

    if (timeRange === "30d") {
      daysToSubtract = 30;
    } else if (timeRange === "7d") {
      daysToSubtract = 7;
    }

    const startDate = new Date(
      referenceDate.getTime() - daysToSubtract * 24 * 60 * 60 * 1000
    );

    return parsedDate >= startDate;
  });

  return (
    // Main container card for the chart
    <Card className="@container/card">
      {/* Header section: Title, optional description, and time range controls */}
      <CardHeader>
        <CardTitle>Historical Account Value</CardTitle>

        {/* Mobile-only description hinting at the default range */}
        <CardDescription>
          <span className="@[540px]/card:hidden">Last 3 months</span>
        </CardDescription>

        {/* Toggle controls for selecting the time range */}
        <CardAction>
          {/* Desktop: horizontal toggle buttons */}
          <ToggleGroup
            type="single"
            value={timeRange}
            onValueChange={setTimeRange}
            variant="outline"
            className="hidden *:data-[slot=toggle-group-item]:!px-4 @[767px]/card:flex"
          >
            <ToggleGroupItem value="90d">Last 3 months</ToggleGroupItem>
            <ToggleGroupItem value="30d">Last 30 days</ToggleGroupItem>
            <ToggleGroupItem value="7d">Last 7 days</ToggleGroupItem>
          </ToggleGroup>

          {/* Mobile: dropdown select for compact view */}
          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger
              className="flex w-40 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[767px]/card:hidden"
              size="sm"
              aria-label="Select a value"
            >
              <SelectValue placeholder="Last 3 months" />
            </SelectTrigger>
            <SelectContent className="rounded-xl">
              <SelectItem value="90d" className="rounded-lg">
                Last 3 months
              </SelectItem>
              <SelectItem value="30d" className="rounded-lg">
                Last 30 days
              </SelectItem>
              <SelectItem value="7d" className="rounded-lg">
                Last 7 days
              </SelectItem>
            </SelectContent>
          </Select>
        </CardAction>
      </CardHeader>

      {/* Chart rendering section */}
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer
          config={chartConfig}
          className="aspect-auto h-[250px] w-full"
        >
          {/* Area chart using Recharts */}
          <AreaChart data={filteredData}>
            {/* Gradient fill for area background */}
            <defs>
              <linearGradient id="fillValue" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--primary)"
                  stopOpacity={0.8}
                />
                <stop
                  offset="95%"
                  stopColor="var(--primary)"
                  stopOpacity={0.05}
                />
              </linearGradient>
            </defs>

            {/* Grid lines for better readability */}
            <CartesianGrid vertical={false} />

            {/* X-axis showing formatted dates */}
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => {
                try {
                  const d = new Date(value);
                  return d.toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                  });
                } catch {
                  return "Invalid";
                }
              }}
            />

            {/* Y-axis automatically adjusts to data range */}
            <YAxis
              domain={["auto", (dataMax: number) => dataMax * 1.05]}
              tickLine={false}
              axisLine={false}
              tickMargin={8}
            />

            {/* Custom tooltip on hover to show time and value */}
            <ChartTooltip
              cursor={false}
              defaultIndex={isMobile ? -1 : 10}
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => {
                    const date = new Date(value);
                    return isNaN(date.getTime())
                      ? "Invalid"
                      : date.toLocaleString("en-US", {
                          month: "short",
                          day: "numeric",
                          hour: "numeric",
                          minute: "2-digit",
                          hour12: true,
                        });
                  }}
                  indicator="dot"
                />
              }
            />

            {/* Actual area data line and fill */}
            <Area
              dataKey="value"
              type="natural"
              fill="url(#fillValue)"
              stroke="var(--primary)"
              strokeWidth={2}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
