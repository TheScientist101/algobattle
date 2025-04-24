import { DocumentReference, Timestamp } from "firebase/firestore";

export type Bot = {
  name: string;
  apiKey: string;
  cash: number
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
}

export type HistoricalValueEntry = {
  date: Date;
  value: number;
};

export type HistoricalValueEntryFirebase = {
  date: Timestamp;
  value: number;
};

export type WithBot = {
  botId: string;
};