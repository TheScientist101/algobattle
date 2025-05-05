//All helper methods for accessing/modifying data from the users collection in Firebase

import { auth, db } from "@/firebase/firebase";
import { updateProfile } from "firebase/auth";
import { doc, getDoc, setDoc } from "firebase/firestore";
import { generateFromEmail } from "unique-username-generator";

/**
 * Ensures a Firestore document exists for the currently authenticated user.
 *
 * - If the user has no existing document in the "users" collection:
 *   - A display name is generated from their email.
 *   - A new user document is created with display name, creation timestamp, and empty bots list.
 *
 * - If the document already exists:
 *   - The Firebase auth profile display name is updated to match the stored display name.
 *
 * @returns {Promise<void>}
 */
export async function ensureUserDocExists(): Promise<void> {
  const user = auth.currentUser;

  if (!user) {
    console.warn("No authenticated user.");
    return;
  }

  const userRef = doc(db, "users", user.uid);
  const userSnap = await getDoc(userRef);

  if (!userSnap.exists()) {
    await updateProfile(user, {
      displayName: generateFromEmail(user?.email || "anonymous@example.com", 5),
    });

    const userData = {
      displayName: user.displayName,
      createdAt: new Date().toISOString(),
      bots: [],
    };

    await setDoc(userRef, userData);
    console.log("User document created.");
  } else {
    await updateProfile(user, {
      displayName: userSnap.data().displayName,
    });

    console.log("User document already exists.");
  }
}
