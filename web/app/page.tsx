"use client";

import React, { useEffect, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { Bot, CompleteHoldingData, Holdings } from "@/utils/types";
import { getBotData, getBots, getHoldings } from "@/utils/botData";
import { TradeTable } from "@/components/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuth } from "@/hooks/authContext";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import LoadingScreen from "@/components/loading";
import axios from "axios";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

export default function Page() {
  const [cards, setCards] = useState<Bot[]>([]);
  const [currBot, setCurrBot] = useState<Bot>();
  const [tickers, setTickers] = useState<Map<string, CompleteHoldingData>>();
  const [selectedKey, setSelectedKey] = useState<string | undefined>(undefined);
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    const fetchiInitial = async () => {
      let botId = "";
      setLoading(true);
      const data = await getBots(user?.uid as string);
      setCards(data);
      if (data.length > 0) {
        setSelectedKey(data[0].apiKey);
        botId = data[0].apiKey;
        await fetchHoldingsData(botId);
      }
      setLoading(true);
    };

    const fetchHoldingsData = async (botId: string) => {
      setLoading(true);

      try {
        const holdings: Holdings = await getHoldings(botId);
        const response = await axios.get("/api/live-stock", {
          headers: {
            Authorization: botId,
          },
        });
        const prices = response.data.payload;
        const dataMap = new Map<string, CompleteHoldingData>();
        const tickers = Object.keys(holdings);

        for (const ticker of tickers) {
          const holding = holdings[ticker];
          const currentPrice = prices[ticker];

          if (currentPrice === undefined) continue;

          const currentValue = holding.numShares * currentPrice;
          const gainLoss = currentValue - holding.purchaseValue;
          const percentChange = (gainLoss / holding.purchaseValue) * 100;

          dataMap.set(ticker, {
            ...holding,
            currentPrice,
            currentValue,
            gainLoss,
            percentChange,
          });
        }
        setTickers(dataMap);
        console.log(dataMap);
      } catch (error) {
        console.error("Error:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchiInitial();
  }, [user?.uid]);

  useEffect(() => {
    const fetchBotData = async () => {
      if (selectedKey) {
        setLoading(true);
        const data = await getBotData(selectedKey);
        setCurrBot(data);
        setLoading(false);
      }
    };

    fetchBotData();
  }, [selectedKey, user?.uid]);

  return (
    <SidebarProvider className="dark">
      <AppSidebar variant="inset" />
      <SidebarInset>
        <div className="flex flex-1 flex-col px-8 relative">
          {selectedKey && currBot && (
            <Card className="mx-4 my-5">
              <div className="flex justify-between items-start p-4">
                <CardHeader className="text-lg font-semibold">
                  {currBot.name}
                </CardHeader>
                <Select value={selectedKey} onValueChange={setSelectedKey}>
                  <SelectTrigger className="w-56 rounded-lg border border-gray-700 bg-muted h-10">
                    <SelectValue placeholder="Select an account" />
                  </SelectTrigger>
                  <SelectContent className="w-56 z-50 bg-black text-white border border-gray-700">
                    <div className="max-h-64 overflow-y-auto">
                      {cards.map((card) =>
                        card.apiKey ? (
                          <SelectItem
                            key={card.apiKey}
                            value={card.apiKey}
                            className="bg-black text-white hover:bg-gray-800 focus:bg-gray-800"
                          >
                            <span className="text-sm text-white">
                              {card.name}
                            </span>
                          </SelectItem>
                        ) : null
                      )}
                    </div>
                  </SelectContent>
                </Select>
              </div>
              <div className="ml-10 mb-5">
                <CardDescription>API KEY: {currBot.apiKey}</CardDescription>
                <CardDescription>
                  Cash:{" "}
                  {new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                  }).format(currBot.cash)}
                </CardDescription>
                <CardDescription>
                  Account value:{" "}
                  {new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "USD",
                  }).format(currBot.accountValue)}
                </CardDescription>
              </div>
            </Card>
          )}
          <div className="mx-4 mb-8">
            <h2 className="text-white text-xl font-semibold mb-4">Holdings</h2>
            {tickers && tickers.size > 0 && (
              <div className="relative">
                <Carousel opts={{ align: "start" }}>
                  <CarouselContent className="-ml-4">
                    {Array.from(tickers.entries()).map(([ticker, info]) => (
                      <CarouselItem
                        key={ticker}
                        className="pl-4 basis-full sm:basis-1/2 md:basis-1/3 lg:basis-1/4"
                      >
                        <Card className="bg-card text-white h-full">
                          <CardHeader>
                            <CardTitle className="text-md">{ticker}</CardTitle>
                          </CardHeader>
                          <CardContent className="text-sm space-y-1">
                            <p>
                              Current Value:{" "}
                              {info.currentValue.toLocaleString("en-US", {
                                style: "currency",
                                currency: "USD",
                              })}
                            </p>
                            <p>
                              Change in value:{" "}
                              <span
                                className={
                                  info.gainLoss >= 0
                                    ? "text-green-400"
                                    : "text-red-400"
                                }
                              >
                                {info.gainLoss.toFixed(2)}
                              </span>
                            </p>
                            <p>
                              Percent Change:{" "}
                              <span
                                className={
                                  info.percentChange >= 0
                                    ? "text-green-400"
                                    : "text-red-400"
                                }
                              >
                                {info.percentChange.toFixed(2)}%
                              </span>
                            </p>
                          </CardContent>
                        </Card>
                      </CarouselItem>
                    ))}
                  </CarouselContent>

                  <CarouselPrevious className="absolute left-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
                  <CarouselNext className="absolute right-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
                </Carousel>
              </div>
            )}
          </div>

          <div className="pt-8 px-4">
            {selectedKey ? (
              <>
                <ChartAreaInteractive botId={selectedKey} />
                <div className="mt-10 mb-10">
                  <h1 className="text-white text-xl my-8">
                    Transaction History
                  </h1>
                  <TradeTable botId={selectedKey} />
                </div>
              </>
            ) : (
              <div className="flex items-center justify-center h-96">
                {loading && <LoadingScreen />}
                {!loading && (
                  <p className="text-white text-lg">
                    Create a bot to get started
                  </p>
                )}
              </div>
            )}
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
