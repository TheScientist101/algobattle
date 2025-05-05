"use client";

import React, { createContext, useContext, useEffect, useState } from 'react';
import { User, onAuthStateChanged } from 'firebase/auth';
import { auth } from '@/firebase/firebase';

/**
 * authContext
 * --------------------------------------
 * This module sets up a global authentication context for use across the app.
 * It leverages Firebase Authentication to track and provide user login state.
 * 
 * Main components:
 * - `AuthProvider`: A React context provider that manages the user's auth state.
 *    - Uses Firebase's `onAuthStateChanged` to reactively update user info.
 *    - Sets a `loading` flag to prevent rendering UI before auth is ready.
 * 
 * - `useAuth`: A custom hook for accessing the current user and loading state from any component.
 *    - Ensures components are within the `AuthProvider`.
 * 
 * Usage:
 * - Wrap your application (typically `_app.tsx` or `layout.tsx`) with `<AuthProvider>`.
 * - Use `useAuth()` anywhere in the tree to get `{ user, loading }`.
 */

interface AuthContextType {
  user: User | null;  
  loading: boolean;    
}

// Create the context with an initial undefined value.
// Consumers must be within an AuthProvider to access it.
const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * AuthProvider
 * Wraps the app in a context provider to supply authentication state.
 * - Tracks the authenticated user with Firebase.
 * - Supplies `user` and `loading` state to children.
 * - Prevents rendering children until the auth state is resolved.
 */
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);     
  const [loading, setLoading] = useState<boolean>(true);  

  /**
   * Effect: Set up a Firebase auth state listener.
   * This runs once on component mount.
   * - Firebase will invoke this callback whenever the user logs in, logs out, or reloads the page.
   * - Sets the `user` state to the logged-in user (or null if not logged in).
   * - Once the initial check is complete, `loading` is set to false to allow UI rendering.
   */
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);      
      setLoading(false);        
    });

    return () => unsubscribe();
  }, []);

  const value = {
    user,
    loading,
  };

  /**
   * Provider component renders children **only after** auth state is known.
   * This avoids flashing of unauthenticated UI during initial load.
   */
  return <AuthContext.Provider value={value}>{!loading && children}</AuthContext.Provider>;
};

/**
 * useAuth
 * Custom hook to access the authentication context from any component.
 * - Returns the `user` and `loading` state.
 * - Throws an error if used outside of the `AuthProvider` to enforce proper usage.
 */
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
