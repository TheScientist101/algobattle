import { auth } from './firebase';
import {
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  signOut,
  User
} from 'firebase/auth';

const googleProvider = new GoogleAuthProvider();

export const signUp = async (email: string, password: string): Promise<User> => {
  try {
    const userCredential = await createUserWithEmailAndPassword(auth, email, password);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing up:', error.message);
      throw error;
    } else {
      console.error('An unknown error occurred during sign-up.');
      throw new Error('Unknown error during sign-up');
    }
  }
};

export const signIn = async (email: string, password: string): Promise<User> => {
  try {
    const userCredential = await signInWithEmailAndPassword(auth, email, password);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing in:', error.message);
      throw error;
    } else {
      console.error('An unknown error occurred during sign-in.');
      throw new Error('Unknown error during sign-in');
    }
  }
};

export const signInWithGoogle = async (): Promise<User> => {
  try {
    const userCredential = await signInWithPopup(auth, googleProvider);
    return userCredential.user;
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing in with Google:', error.message);
      throw error;
    } else {
      console.error('An unknown error occurred during Google sign-in.');
      throw new Error('Unknown error during Google sign-in');
    }
  }
};

export const logOut = async (): Promise<void> => {
  try {
    await signOut(auth);
    console.log('User signed out');
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error('Error signing out:', error.message);
      throw error;
    } else {
      console.error('An unknown error occurred during sign-out.');
      throw new Error('Unknown error during sign-out');
    }
  }
};
