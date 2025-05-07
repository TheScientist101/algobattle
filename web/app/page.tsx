"use client";

import React, { useEffect, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { Bot, CompleteHoldingData, Holdings } from "@/utils/types";
import { createBot, getBotData, getBots, getHoldings } from "@/utils/botData";
import { TradeTable } from "@/components/data-table";
import { useAuth } from "@/hooks/authContext";
import { HoldingsCarousel } from "@/components/bot-holdings";
import { BotInfoCard } from "@/components/bot-info";
import { NoBotPrompt } from "@/components/no-bot";
import LoadingScreen from "@/components/loading";
import axios from "axios";

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
    const fetchInitial = async () => {
      setLoading(true); 
      const data = await getBots(user?.uid as string);
      setCards(data);
      if (data.length > 0) setSelectedKey(data[0].apiKey);
      setLoading(false);
    };
    fetchInitial();
  }, [user?.uid]);

  //Runs whenever `selectedKey` (the selected bot) changes.
  useEffect(() => {
    if (!selectedKey) return; 
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
            const botValue = holding.purchaseValue * holding.numShares;
            const gainLoss = botValue - currentValue;
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
            <BotInfoCard
              bots={cards}
              selectedKey={selectedKey}
              onSelect={setSelectedKey}
              currentBot={currBot}
            />
          )}
          <div className="pt-8 px-4">
            {selectedKey ? (
              <>
                {/* Interactive chart component that shows if there is a value for the selectedKey useState variable */}
                <ChartAreaInteractive botId={selectedKey} />

                {tickers && tickers.size > 0 && (
                  <HoldingsCarousel holdings={tickers} />
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
                {/* Upload the bot to firebase and display it when the bot is created on the NoBotPrompt component */}
                {!loading && (
                  <NoBotPrompt
                    onCreate={async (botName) => {
                      await createBot(botName, user?.uid as string);
                      const updatedBots = await getBots(user?.uid as string);
                      setCards(updatedBots);
                      if (updatedBots.length > 0)
                        setSelectedKey(updatedBots[0].apiKey);
                    }}
                  />
                )}
              </div>
            )}
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
