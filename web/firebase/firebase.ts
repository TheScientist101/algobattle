//Setup firebase to be used in the application
import { initializeApp } from "firebase/app"; 
import { getAuth } from "firebase/auth"; 
import { getFirestore } from "firebase/firestore"; 

// Firebase configuration object using environment variables for security
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_AUTH_DOMAIN, 
  projectId: process.env.NEXT_PUBLIC_PROJECT_ID, 
  storageBucket: process.env.NEXT_PUBLIC_STORAGE_BUCKET, 
  messagingSenderId: process.env.NEXT_PUBLIC_MESSAGING_SENDER_ID, 
  appId: process.env.NEXT_PUBLIC_APP_ID, 
  measurementId: process.env.NEXT_PUBLIC_MEASUREMENT_ID 
};

// Initialize the Firebase app using the provided config
const app = initializeApp(firebaseConfig);

// Set up Firebase Authentication using the initialized app
const auth = getAuth(app);

// Set up Firestore Database instance using the initialized app
const db = getFirestore(app);

// Export auth and db so they can be used throughout the application
export { auth, db };
