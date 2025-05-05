import { DocumentReference, Timestamp } from "firebase/firestore";

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

export type Trade = {
  bot: DocumentReference;
  ticker: string;
  action: "buy" | "sell";
  numShares: number;
  unitCost: number;
  time: Date;
};

export type Holdings = {
  [ticker: string]: Holding;
};

export type CompleteHoldingData = {
  numShares: number;
  purchaseValue: number;
  currentPrice: number;
  currentValue: number;
  gainLoss: number;
  percentChange: number;
};


type Holding = {
  numShares: number;
  purchaseValue: number;
};

export type LeaderboardEntry = {
  user: string;
  name: string;
  historicalAccountValue: { date: Timestamp; value: number }[];
};

export type WithBot = {
  botId: string;
};
