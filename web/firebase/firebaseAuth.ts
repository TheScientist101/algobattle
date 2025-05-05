// Firebase authentication utility functions
import { auth } from './firebase';
import {
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  signOut,
  User
} from 'firebase/auth';

// Initialize Google OAuth provider
const googleProvider = new GoogleAuthProvider();

/**
 * Registers a new user using Firebase Authentication with email and password.
 *
 * @param email - Email address for the new account
 * @param password - Password for the new account
 * @returns A Promise resolving to the created Firebase User
 * @throws If Firebase fails to create the account (e.g., email already in use)
 */
export const signUp = async (email: string, password: string): Promise<User> => {
  try {
    const userCredential = await createUserWithEmailAndPassword(auth, email, password);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing up:', error.message);
      throw error;
    } else {
      console.error('Unknown error during sign-up');
      throw new Error('Unknown error during sign-up');
    }
  }
};

/**
 * Authenticates an existing user via Firebase Authentication using email and password.
 *
 * @param email - Registered email address
 * @param password - Corresponding password
 * @returns A Promise resolving to the authenticated Firebase User
 * @throws If Firebase fails to authenticate the user (e.g., wrong credentials)
 */
export const signIn = async (email: string, password: string): Promise<User> => {
  try {
    const userCredential = await signInWithEmailAndPassword(auth, email, password);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing in:', error.message);
      throw error;
    } else {
      console.error('Unknown error during sign-in');
      throw new Error('Unknown error during sign-in');
    }
  }
};

/**
 * Signs in a user using Firebase Authentication with Google OAuth via a popup.
 *
 * @returns A Promise resolving to the authenticated Firebase User
 * @throws If Firebase fails to complete the Google sign-in process
 */
export const signInWithGoogle = async (): Promise<User> => {
  try {
    const userCredential = await signInWithPopup(auth, googleProvider);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing in with Google:', error.message);
      throw error;
    } else {
      console.error('Unknown error during Google sign-in');
      throw new Error('Unknown error during Google sign-in');
    }
  }
};

/**
 * Signs out the currently authenticated user via Firebase Authentication.
 *
 * @returns A Promise that resolves once the sign-out is complete
 * @throws If Firebase fails to sign the user out
 */
export const logOut = async (): Promise<void> => {
  try {
    await signOut(auth);
    console.log('User signed out');
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing out:', error.message);
      throw error;
    } else {
      console.error('Unknown error during sign-out');
      throw new Error('Unknown error during sign-out');
    }
  }
};