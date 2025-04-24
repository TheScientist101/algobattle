"use client";

import React, { useEffect, useState } from "react";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { Bot } from "@/utils/types";
import { getBots } from "@/utils/botData";
import { TradeTable } from "@/components/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAuth } from "@/hooks/authContext";

export default function Page() {
  const [cards, setCards] = useState<Bot[]>([]);
  const [selectedKey, setSelectedKey] = useState<string | undefined>(undefined);
  const { user } = useAuth();

  useEffect(() => {
    const fetchBots = async () => {
      const data = await getBots(user?.uid as string);
      console.log(data);
      setCards(data);
      if (data.length > 0) {
        setSelectedKey(data[0].apiKey);
      }
    };
    if (user?.uid) fetchBots();
  }, [user?.uid]);

  return (
    <SidebarProvider className="dark">
      <AppSidebar variant="inset" />
      <SidebarInset>
        <div className="flex flex-1 flex-col px-8 relative">
          <div className="absolute top-20 right-12 z-50">
            {selectedKey && (
              <Select value={selectedKey} onValueChange={setSelectedKey}>
                <SelectTrigger className="w-72 rounded-lg border border-gray-700 bg-muted h-40">
                  <SelectValue placeholder="Select an account" />
                </SelectTrigger>
                <SelectContent className="w-72 z-50 bg-black text-white border border-gray-700">
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
            )}
          </div>

          <div className="pt-40 px-4">
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
