"use client";

import { TrophyIcon, MedalIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AppSidebar } from "../../components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { getLeaderboardEntries } from "@/utils/botData";
import LoadingScreen from "@/components/loading";

/**
 * Displays a sorted list of bots based on their performance (% gain).
 * - Fetches bot historical account value data from Firestore.
 * - Computes percent change for each bot.
 * - Ranks bots and displays them in a table.
 * - Highlights top 3 with icons (trophy and medals).
 */
export default function LeaderboardPage() {
  const [leaderboardData, setLeaderboardData] = useState<
    {
      userName: string;
      name: string;
      percentChange: number;
    }[]
  >([]);
  const [loading, setLoading] = useState(false);

  /**
   * useEffect: On initial mount, fetch and process leaderboard data.
   * 
   * - Calls `getLeaderboardEntries()` to retrieve bot documents.
   * - Each bot includes a `historicalAccountValue` array.
   * - Calculates percent change between first and last recorded values.
   * - Rounds to 2 decimal places.
   * - Sorts bots in descending order of performance.
   * - Updates `leaderboardData` with the processed list.
   */
  useEffect(() => {
    const fetchLeaderboard = async () => {
      setLoading(true); // Show loading screen during fetch

      try {
        // Step 1: Fetch raw leaderboard entries from Firestore
        const data = await getLeaderboardEntries();

        // Step 2: Compute percent change for each bot
        const computed = data.map((entry) => {
          const values = entry.historicalAccountValue;
          const first = values[0]?.value;
          const last = values[values.length - 1]?.value;

          let percentChange = 0;

          // Only calculate if both values exist and we have enough history
          if (first != null && last != null && values.length > 1) {
            percentChange = ((last - first) / first) * 100;
          }

          return {
            name: entry.name,
            percentChange: Math.round(percentChange * 100) / 100,
            userName: entry.user,
          };
        });

        // Step 3: Sort by highest percentChange first
        const sorted = computed.sort(
          (a, b) => b.percentChange - a.percentChange
        );

        // Step 4: Save sorted data to state
        setLeaderboardData(sorted);
      } catch (error) {
        console.error("Failed to fetch leaderboard data", error);
      } finally {
        setLoading(false); // Hide loading screen
      }
    };

    fetchLeaderboard();
  }, []);

  return (
    <SidebarProvider className="dark">
      <AppSidebar variant="inset" />
      <SidebarInset>
        <div className="flex flex-1 flex-col">
          <div className="flex flex-1 flex-col">
            <div className="flex flex-col my-8">
              <div className="mx-8">
                <Card>
                  <CardHeader className="relative">
                    <div className="flex flex-col gap-1">
                      <CardTitle className="text-2xl flex items-center gap-2">
                        Algobattle Leaderboard
                        {loading && <LoadingScreen />}
                      </CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="rounded-lg border">
                      <Table>
                        {/* Table header: column labels */}
                        <TableHeader>
                          <TableRow>
                            <TableHead className="w-30 text-center">
                              Rank
                            </TableHead>
                            <TableHead>Owner</TableHead>
                            <TableHead>Bot Name</TableHead>
                            <TableHead className="text-right pl-20 pr-10">
                              % Change
                            </TableHead>
                          </TableRow>
                        </TableHeader>

                        {/* Table body: list of leaderboard entries */}
                        <TableBody>
                          {leaderboardData.map((user, index) => (
                            <TableRow key={index}>
                              {/* Rank column with icons for top 3 */}
                              <TableCell className="text-center font-medium">
                                {index === 0 ? (
                                  <div className="flex justify-center">
                                    <TrophyIcon className="h-5 w-5 text-yellow-500" />
                                  </div>
                                ) : index === 1 ? (
                                  <div className="flex justify-center">
                                    <MedalIcon className="h-5 w-5 text-gray-400" />
                                  </div>
                                ) : index === 2 ? (
                                  <div className="flex justify-center">
                                    <MedalIcon className="h-5 w-5 text-amber-700" />
                                  </div>
                                ) : (
                                  index + 1
                                )}
                              </TableCell>

                              {/* Username */}
                              <TableCell>
                                <div>{user.userName}</div>
                              </TableCell>

                              {/* Bot name */}
                              <TableCell>
                                <div className="text-muted-foreground">
                                  {user.name}
                                </div>
                              </TableCell>

                              {/* Percent change with color formatting */}
                              <TableCell className="text-right pr-10">
                                <span
                                  className={`${
                                    user.percentChange < 0
                                      ? "text-red-500"
                                      : user.percentChange > 0
                                      ? "text-green-500"
                                      : "text-muted-foreground"
                                  }`}
                                >
                                  {`${
                                    user.percentChange > 0
                                      ? "+"
                                      : ""
                                  }${user.percentChange.toFixed(2)}%`}
                                </span>
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}