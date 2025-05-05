//all helper mehods that access the bots collection in Firebase

"use client";

import { db } from "@/firebase/firebase";
import {
  arrayUnion,
  collection,
  doc,
  DocumentReference,
  getDoc,
  getDocs,
  setDoc,
  Timestamp,
  updateDoc,
} from "firebase/firestore";
import { v4 } from "uuid";
import { Bot, Holdings, LeaderboardEntry, Trade } from "./types";

/**
 * Creates a new bot and links it to the specified user.
 *
 * - Generate a unique bot ID using `uuid`.
 * - Create a bot document with default values: empty transactions, initial cash, etc.
 * - Store the bot in the "bots" collection using `setDoc`.
 * - Update the user document by appending a reference to the new bot using `arrayUnion`.
 * - Ensures each user can manage multiple bots via stored references.
 *
 * @param {string} botName - The name of the bot to be created.
 * @param {string} user - The UID of the authenticated user creating the bot.
 */
export async function createBot(botName: string, user: string): Promise<void> {
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

/**
 * Retrieves all bots associated with a given user.
 *
 * - Fetch the user document using the UID.
 * - Extract the "bots" array field (document references).
 * - Resolve each reference using `getDoc` in parallel with `Promise.all`.
 * - Return an array of existing bot objects.
 *
 * @param {string} user - The UID of the user.
 * @returns {Promise<Bot[]>} - An array of Bot objects linked to the user.
 */
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
  return botSnaps
    .filter((snap) => snap.exists())
    .map((snap) => snap.data() as Bot);
}

/**
 * Fetches the historical account value of a bot, formatted for charting.
 *
 * - Get the bot document from Firestore.
 * - Extract the "historicalAccountValue" array, defaulting to an empty array if undefined.
 * - Convert Firestore Timestamps to ISO strings and return a list of date-value pairs.
 *
 * @param {string} botId - The ID of the bot.
 * @returns {Promise<{ date: string; value: number }[]>} - An array of date-value pairs.
 */
export async function getBotHistory(
  botId: string
): Promise<{ date: string; value: number }[]> {
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

/**
 * Retrieves all trade records associated with a bot.
 *
 * - Access the "transactions" array from the bot document (which holds trade doc references).
 * - Fetch each trade document in parallel using `Promise.all`.
 * - Convert each Firestore Timestamp to a native Date object for use in UI/charting.
 * - Return an array of complete Trade objects.
 *
 * @param {string} botId - The ID of the bot.
 * @returns {Promise<Trade[]>} - An array of Trade objects.
 */
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
  return tradeDocs
    .filter((doc) => doc.exists())
    .map((doc) => {
      const data = doc.data() as Omit<Trade, "time"> & { time: Timestamp };
      return {
        ...data,
        time: data.time.toDate(),
      } as Trade;
    });
}

/**
 * Retrieves the full bot document for a specific bot ID.
 *
 * - Fetch and return the entire bot object directly from Firestore.
 * - This is useful for detailed views or analytics involving a single bot.
 *
 * @param {string} botId - The ID of the bot.
 * @returns {Promise<Bot>} - The complete Bot object.
 */
export async function getBotData(botId: string): Promise<Bot> {
  const botRef = doc(db, "bots", botId);
  const botSnap = await getDoc(botRef);
  return botSnap.data() as Bot;
}

/**
 * Retrieves all leaderboard entries based on bots' performance data.
 *
 * - Fetch all bot documents from Firestore.
 * - For each bot, extract the "historicalAccountValue" array.
 * - Resolve the owner's user document (if available) to fetch display names.
 * - Compose and return an array of LeaderboardEntry objects with user, bot name, and value history.
 *
 * @returns {Promise<LeaderboardEntry[]>} - An array of leaderboard entries.
 */
export async function getLeaderboardEntries(): Promise<LeaderboardEntry[]> {
  const leaderboardSnaps = await getDocs(collection(db, "bots"));
  const entries: LeaderboardEntry[] = [];

  for (const snapshot of leaderboardSnaps.docs) {
    const data = snapshot.data();
    const values = Array.isArray(data.historicalAccountValue)
      ? data.historicalAccountValue.map(
          (entry: { date: Timestamp; value: number }) => ({
            value: entry.value,
            date: entry.date,
          })
        )
      : [];

    const userId = typeof data.owner === "string" ? data.owner : "";
    let displayName = userId;

    if (userId) {
      try {
        const userDocRef = doc(db, "users", userId);
        const userSnap = await getDoc(userDocRef);
        if (userSnap.exists()) {
          displayName = userSnap.data().displayName || userId;
        }
      } catch (error) {
        console.warn("Failed to fetch user doc:", error);
      }
    }

    entries.push({
      name: data.name,
      user: displayName,
      historicalAccountValue: values,
    });
  }

  return entries;
}

/**
 * Retrieves the holdings for a specific bot.
 *
 * - Fetch the bot document using the bot ID.
 * - If the bot exists, return the "holdings" object directly.
 * - If it doesn't exist, throw an error so the caller can handle it.
 *
 * @param {string} botId - The ID of the bot.
 * @returns {Promise<Holdings>} - The holdings data for the bot.
 * @throws Will throw an error if the bot document does not exist.
 */
export async function getHoldings(botId: string): Promise<Holdings> {
  const botRef = doc(db, "bots", botId);
  const botSnap = await getDoc(botRef);

  if (botSnap.exists()) {
    const data = botSnap.data();
    return data.holdings as Holdings;
  } else {
    throw new Error("Bot document not found");
  }
}
