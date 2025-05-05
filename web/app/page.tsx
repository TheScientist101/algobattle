//Dashboard
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
  const [cards, setCards] = useState<Bot[]>([]); // List of all bots belonging to the user
  const [currBot, setCurrBot] = useState<Bot>(); // Currently selected bot
  const [tickers, setTickers] = useState<Map<string, CompleteHoldingData>>(); // Map of stock tickers to their value, percent change, and change in value to the user
  const [selectedKey, setSelectedKey] = useState<string | undefined>(undefined); // Selected bot's API key
  const [loading, setLoading] = useState(false); // Loading state for async operations
  const { user } = useAuth(); // Get authenticated user from auth context

  /**
   * Effect: Runs once on initial mount when the user's UID becomes available.
   * 1. Set the loading state (`loading`) to true to indicate data fetching has started.
   * 2. Call `getBots(user?.uid)` to retrieve all bots associated with the current user from Firestore.
   *    Updates `cards` state with the list of `Bot` objects retrieved.
   * 3. If the user has at least one bot: Automatically select the first bot by updating `selectedKey` with its API key.
   * 4. Set the loading state (`loading`) back to false to signal data is ready.
   */
  useEffect(() => {
    const fetchiInitial = async () => {
      setLoading(true);
      const data = await getBots(user?.uid as string);
      setCards(data);
      if (data.length > 0) {
        setSelectedKey(data[0].apiKey);
      }
      setLoading(true);
    };
    fetchiInitial();
  }, [user?.uid]);

  //Runs whenever `selectedKey` (the selected bot) changes.
  useEffect(() => {
    /**
     * - Check if a bot is selected (`selectedKey` is defined).
     * - Call `getBotData(selectedKey)` to retrieve the full bot object from Firestore.
     * - Update `currBot` state with the retrieved bot data.
     */
    const fetchBotData = async () => {
      if (selectedKey) {
        const data = await getBotData(selectedKey);
        setCurrBot(data);
      }
    };

    /**
     * Fetches holdings and live prices for the selected bot and updates ticker data.
     * - Retrieve holdings from Firestore.
     * - Get live stock prices via API using the bot's API key.
     * - For each holding:
     *   - Calculate current value, gain/loss, and percent change.
     *   - Store the data in a Map with the ticker as key.
     * - Update the `tickers` state with the computed Map.
     * Advanced structure used: `Map<string, CompleteHoldingData>`
     */
    const fetchHoldingsData = async (botId: string) => {
      try {
        const holdings: Holdings = await getHoldings(botId);
        console.log(holdings);
        const response = await axios.get("/api/live-stock", {
          headers: {
            Authorization: botId,
          },
        });
        const prices = response.data.payload;

        const dataMap = new Map<string, CompleteHoldingData>();
        const tickers = Object.keys(holdings);

        if (holdings) {
          for (const ticker of tickers) {
            const holding = holdings[ticker];
            const currentPrice = prices[ticker];

            if (holding.numShares == 0) continue;
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
        } else {
          setTickers(Object.create(null));
        }
      } catch (error) {
        console.error("Error:", error);
      }
    };
    setLoading(true);
    fetchHoldingsData(selectedKey as string);
    fetchBotData();
    setLoading(false);
  }, [selectedKey, user?.uid]);

  return (
    <SidebarProvider className="dark">
      {/* Sidebar component */}
      <AppSidebar variant="inset" />

      {/* Main content area */}
      <SidebarInset>
        <div className="flex flex-1 flex-col px-8 relative">
          {/* Bot information card only shows if the current bot has data and there is a selected key to ensure the 
          bot displayed is the one currently selected */}
          {selectedKey && currBot && (
            <Card className="mx-4 my-5">
              <div className="flex justify-between items-start p-4">
                {/* Bot name */}
                <CardHeader className="text-lg font-semibold">
                  {currBot.name}
                </CardHeader>

                {/* Bot selection dropdown  - a list of the user's bot' names*/}
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

              {/* Bot financial information: Its api key, cash, and account value */}
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
          <div className="pt-8 px-4">
            {selectedKey ? (
              <>
                {/* Interactive chart component that shows if there is a value for the selectedKey useState variable */}
                <ChartAreaInteractive botId={selectedKey} />

                {/* Holdings Section that displays if the tickers map has items */}
                {tickers && tickers.size > 0 && (
                  <h1 className="text-white text-xl my-8">Current Holdings</h1>
                )}
                {tickers && tickers.size > 0 && (
                  <div className="relative">
                    <Carousel opts={{ align: "start" }}>
                      <CarouselContent className="-ml-4">
                        {/* Map through holdings to create cards of each holding by using the tickers map and using
                         the attributes each ticker is mapped to (using my first advanced data structire); if a ticker has a 
                         greater than 10% change, it is mark it as volitile with a yellow boarder - use of my second advanced data structure */}
                        {Array.from(tickers.entries()).map(([ticker, info]) => (
                          <CarouselItem
                            key={ticker}
                            className="pl-4 basis-full sm:basis-1/2 md:basis-1/3 lg:basis-1/4"
                          >
                            <Card className="bg-card text-white h-full">
                              <CardHeader>
                                {/*Ticker name */}
                                <CardTitle className="text-md">
                                  {ticker}
                                </CardTitle>
                              </CardHeader>
                              <CardContent className="text-sm space-y-1">
                                {/* Metrics data displayL Current value, Change in value, and the Percent Change in value the bot has
                                 achived with a particular stock */}
                                <p>
                                  Current Value:{" "}
                                  {info.currentValue.toLocaleString("en-US", {
                                    style: "currency",
                                    currency: "USD",
                                  })}
                                </p>
                                <p>
                                  Bot Value:{" "}
                                  {info.purchaseValue.toLocaleString("en-US", {
                                    style: "currency",
                                    currency: "USD",
                                  })}
                                </p>
                                <p>
                                  Change in Value:{" $"}
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

                      {/* Carousel navigation buttons to go forward or backward*/}
                      <CarouselPrevious className="absolute left-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
                      <CarouselNext className="absolute right-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
                    </Carousel>
                  </div>
                )}
                {/* Transaction history table that shows if there is a value for the selectedKey useState variable */}
                <div className="mt-10 mb-10">
                  <h1 className="text-white text-xl my-8">
                    Transaction History
                  </h1>
                  <TradeTable botId={selectedKey} />
                </div>
              </>
            ) : (
              /* Display when no bot is selected telling the user to create a bot or shows the loading screen if the loading state is true */
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
