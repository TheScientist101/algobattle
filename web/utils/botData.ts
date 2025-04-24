"use client";

import { db } from "@/firebase/firebase";
import {
  arrayUnion,
  doc,
  DocumentReference,
  getDoc,
  setDoc,
  Timestamp,
  updateDoc,
} from "firebase/firestore";
import { v4 } from "uuid";
import {
  Bot,
  Trade,
} from "./types";

export async function createBot(botName: string, user: string) {
  if (!user) {
    console.warn("No authenticated user.");
    return;
  }

  const botId = v4();
  const botRef = doc(db, "bots", botId);

  const botData = {
    name: botName,
    apiKey: botId,
    owner: user,
    transactions: [],
    cash: 2000,
    historicalAccountValue: [],
  };

  await setDoc(botRef, botData);
  const userDocRef = doc(db, "users", user);
  await updateDoc(userDocRef, {
    bots: arrayUnion(botRef),
  });

  console.log("Bot created and linked to user.");
}

export async function getBots(user: string): Promise<Bot[]> {
  if (!user) {
    console.warn("No authenticated user.");
    return [];
  }

  const userDocRef = doc(db, "users", user);
  const userDocSnap = await getDoc(userDocRef);

  if (!userDocSnap.exists()) {
    console.warn("User document not found.");
    return [];
  }

  const botRefs = userDocSnap.data().bots as DocumentReference[] | undefined;

  if (!botRefs || !Array.isArray(botRefs)) {
    console.warn("No bots found in user document.");
    return [];
  }
  const botSnaps = await Promise.all(botRefs.map((ref) => getDoc(ref)));
  const bots: Bot[] = botSnaps
    .filter((snap) => snap.exists())
    .map((snap) => snap.data() as Bot);

  return bots;
}

export async function getBotHistory(botId: string): Promise<{ date: string; value: number }[]> {
  const botRef = doc(db, "bots", botId);
  const botSnap = await getDoc(botRef);

  if (!botSnap.exists()) return [];

  const data = botSnap.data();
  const rawHistory = data.historicalAccountValue ?? [];

  return rawHistory.map((item: { date: Timestamp; value: number }) => ({
    date: item.date.toDate().toISOString(), 
    value: item.value,
  }));
}


export async function getTradesForBot(botId: string): Promise<Trade[]> {
  const botRef = doc(db, "bots", botId);
  const botSnap = await getDoc(botRef);

  if (!botSnap.exists()) {
    console.warn(`Bot with ID ${botId} not found.`);
    return [];
  }

  const data = botSnap.data() as Bot;
  const tradeRefs = data.transactions;

  const tradeDocs = await Promise.all(tradeRefs.map((ref) => getDoc(ref)));
  const trades: Trade[] = tradeDocs
    .filter((doc) => doc.exists())
    .map((doc) => {
      const data = doc.data() as Omit<Trade, "time"> & { time: Timestamp };
      return {
        ...data,
        time: data.time.toDate(),
      } as Trade;
    });

  return trades;
}
