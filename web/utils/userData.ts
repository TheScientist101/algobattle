import { auth, db } from "@/firebase/firebase";
import { updateProfile } from "firebase/auth";
import { doc, getDoc, setDoc } from "firebase/firestore";
import { generateFromEmail } from "unique-username-generator";

export async function ensureUserDocExists() {
  const user = auth.currentUser;

  if (!user) {
    console.warn("No authenticated user.");
    return;
  }

  const userRef = doc(db, "users", user.uid);
  const userSnap = await getDoc(userRef);

  if (!userSnap.exists()) {
    await updateProfile(user, {
      displayName: generateFromEmail(user?.email || "anonymous@example.com ", 5)
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
      displayName: userSnap.data().displayName
    });
    console.log("User document already exists.");
  }
}
