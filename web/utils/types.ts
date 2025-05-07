//all types used in the program

import { DocumentReference, Timestamp } from "firebase/firestore";

/**
 * Represents a trading bot created by a user.
 * Includes core financial state, transaction history, and reference to its owner.
 */
export type Bot = {
  name: string;
  apiKey: string;
  cash: number;
  accountValue: number;
  holdings: []; 
  owner: string;
  transactions: DocumentReference[];
  historicalAccountValue: []; 
};

/**
 * Represents a single buy or sell trade executed by a bot.
 */
export type Trade = {
  bot: DocumentReference;
  ticker: string;
  action: "buy" | "sell";
  numShares: number;
  unitCost: number;
  time: Date;
};

/**
 * Represents a mapping of stock tickers to the basic holding information for a bot.
 */
export type Holdings = {
  [ticker: string]: Holding;
};

/**
 * Represents enriched holding data with both static and real-time values.
 * Used for UI display and performance metrics.
 */
export type CompleteHoldingData = {
  numShares: number;
  purchaseValue: number;
  currentPrice: number;
  currentValue: number;
  gainLoss: number;
  percentChange: number;
};

/**
 * Represents a basic record of a stock holding: number of shares and total cost.
 */
type Holding = {
  numShares: number;
  purchaseValue: number;
};

/**
 * Represents a leaderboard row showing a bot's name, owner, and account history.
 */
export type LeaderboardEntry = {
  user: string;
  name: string;
  historicalAccountValue: { date: Timestamp; value: number }[];
};

/**
 * Utility type for passing a bot identifier between components or services.
 */
export type WithBot = {
  botId: string;
};

/**
 * Holds the information for the card of a bot at the top of the screen.
 */
export type BotInfoCardProps = {
  bots: Bot[]; 
  selectedKey: string;
  onSelect: (key: string) => void; 
  currentBot: Bot; 
}

/**
 * Holds the information for the NoBot component (the function to create a bot and the loadings state)
 */
export type NoBotPromptProps = {
  onCreate: (name: string) => void;
}

/**
 * Holds the information for each holding (all in a map)
 */
export type HoldingsCarouselProps = {
  holdings: Map<string, CompleteHoldingData>; 
}
