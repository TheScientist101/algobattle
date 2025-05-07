import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { HoldingsCarouselProps } from "@/utils/types";

/**
 * Displays a horizontal scrollable carousel of the current stock holdings for a bot.
 *
 * - Uses a `Map` (first advanced data structure) to hold ticker info
 * - Shows detailed metrics per ticker: current value, original purchase value, change, and percent gain/loss.
 */
export function HoldingsCarousel({ holdings }: HoldingsCarouselProps) {
  const entries = Array.from(holdings.entries()); // Convert Map to array for iteration

  return (
    <div className="relative">
      {/* Section title */}
      <h1 className="text-white text-xl my-8">Current Holdings</h1>

      <Carousel opts={{ align: "start" }}>
        {/* 
          Iterate over each ticker in holdings.
          Render a card for each holding with relevant financial data.
        */}
        <CarouselContent className="-ml-4">
          {entries.map(([ticker, info]) => (
            <CarouselItem
              key={ticker}
              className="pl-4 basis-full sm:basis-1/2 md:basis-1/3 lg:basis-1/4"
            >
              <Card
                className={`bg-card text-white h-full ${
                  Math.abs(info.percentChange) > 10
                    ? "border-2 border-yellow-400"
                    : ""
                }`}
              >
                <CardHeader>
                  {/* Ticker symbol */}
                  <CardTitle className="text-md">{ticker}</CardTitle>
                </CardHeader>

                <CardContent className="text-sm space-y-1">
                  {/* Current market value of the holding */}
                  <p>
                    Current Value:{" "}
                    {info.currentValue.toLocaleString("en-US", {
                      style: "currency",
                      currency: "USD",
                    })}
                  </p>

                  {/* Purchase value (bot's investment) */}
                  <p>
                    Bot Value:{" "}
                    {info.purchaseValue.toLocaleString("en-US", {
                      style: "currency",
                      currency: "USD",
                    })}
                  </p>

                  {/* Gain or loss in dollar amount */}
                  <p>
                    Change in Value:{" "}
                    <span
                      className={
                        info.gainLoss >= 0 ? "text-green-400" : "text-red-400"
                      }
                    >
                      ${info.gainLoss.toFixed(2)}
                    </span>
                  </p>

                  {/* Percentage change with color indicator */}
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

        {/* Carousel navigation arrows that appear if there is a need to scroll (more than 4 cards) */}
        {holdings.size > 4 && (
          <>
            <CarouselPrevious className="absolute left-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
            <CarouselNext className="absolute right-[-1.5rem] top-1/2 -translate-y-1/2 z-10" />
          </>
        )}
      </Carousel>
    </div>
  );
}
