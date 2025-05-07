import { Card, CardDescription, CardHeader } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { BotInfoCardProps } from "@/utils/types";

/**
 * Renders a card showing data for a selected bot.
 *
 * - Displays the current bot's name, API key, cash, and account value.
 * - Allows switching between multiple bots using a dropdown.
 * - Styled using ShadCN UI components.
 */
export function BotInfoCard({
  bots,
  selectedKey,
  onSelect,
  currentBot,
}: BotInfoCardProps) {
  return (
    <Card className="mx-4 my-5">
      <div className="flex justify-between items-start p-4">
        {/* Bot Name Display */}
        <CardHeader className="text-lg font-semibold">
          {currentBot.name}
        </CardHeader>

        {/* Bot Selection Dropdown */}
        <Select value={selectedKey} onValueChange={onSelect}>
          <SelectTrigger className="w-56 rounded-lg border border-gray-700 bg-muted h-10">
            <SelectValue placeholder="Select an account" />
          </SelectTrigger>
          <SelectContent className="w-56 z-50 bg-black text-white border border-gray-700">
            <div className="max-h-64 overflow-y-auto">
              {bots.map((bot) =>
                bot.apiKey ? (
                  <SelectItem
                    key={bot.apiKey}
                    value={bot.apiKey}
                    className="bg-black text-white hover:bg-gray-800 focus:bg-gray-800"
                  >
                    <span className="text-sm text-white">{bot.name}</span>
                  </SelectItem>
                ) : null
              )}
            </div>
          </SelectContent>
        </Select>
      </div>

      {/* Financial Details Section (api key, cash, and account value) */}
      <div className="ml-10 mb-5">
        <CardDescription>API KEY: {currentBot.apiKey}</CardDescription>
        <CardDescription>
          Cash:{" "}
          {new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
          }).format(currentBot.cash)}
        </CardDescription>
        <CardDescription>
          Account value:{" "}
          {new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
          }).format(currentBot.accountValue)}
        </CardDescription>
      </div>
    </Card>
  );
}
