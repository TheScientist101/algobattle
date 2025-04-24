"use client";

import React, { useEffect, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { Bot } from "@/utils/types";
import { getBotData, getBots } from "@/utils/botData";
import { TradeTable } from "@/components/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuth } from "@/hooks/authContext";
import { Card, CardDescription, CardHeader } from "@/components/ui/card";
import LoadingScreen from "@/components/loading";

export default function Page() {
  const [cards, setCards] = useState<Bot[]>([]);
  const [currBot, setCurrBot] = useState<Bot>();
  const [selectedKey, setSelectedKey] = useState<string | undefined>(undefined);
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  useEffect(() => {
    const fetchBots = async () => {
      const data = await getBots(user?.uid as string);
      setCards(data);
      if (data.length > 0) {
        setSelectedKey(data[0].apiKey);
      }
    };
    setLoading(true);
    if (user?.uid) fetchBots();
    setLoading(false);
  }, [user?.uid]);

  useEffect(() => {
    const fetchBots = async () => {
      if (selectedKey) setCurrBot(await getBotData(selectedKey));
    };
    setLoading(true);
    fetchBots();
    setLoading(false);
  }, [selectedKey]);
  
  if (loading) return <LoadingScreen />;

  return (
    <SidebarProvider className="dark">
      <AppSidebar variant="inset" />
      <SidebarInset>
        <div className="flex flex-1 flex-col px-8 relative">
          {selectedKey && currBot && (
            <Card className="ml-4 w-270 my-5">
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
                <CardDescription>Cash: {currBot.cash}</CardDescription>
                <CardDescription>
                  Account value: {currBot.accountValue}
                </CardDescription>
              </div>
            </Card>
          )}
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
                <p className="text-white text-lg">
                  Create a bot to get started
                </p>
              </div>
            )}
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
