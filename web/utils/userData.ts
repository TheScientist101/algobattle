import { auth, db } from "@/firebase/firebase";
import {
  doc,
  getDoc,
  setDoc,
} from "firebase/firestore";

export async function ensureUserDocExists() {
  const user = auth.currentUser;

  if (!user) {
    console.warn("No authenticated user.");
    return;
  }

  const userRef = doc(db, "users", user.uid);
  const userSnap = await getDoc(userRef);

  if (!userSnap.exists()) {
    const userData = {
      displayName: user.displayName,
      createdAt: new Date().toISOString(),
      bots: []
    };
    await setDoc(userRef, userData);
    console.log("User document created.");
  } else {
    console.log("User document already exists.");
  }
}
