//Layout manager of the website
"use client";

import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import React, { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import { AuthProvider, useAuth } from "@/hooks/authContext";

// Load custom fonts from Google and assign them to CSS variables
const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

/**
 * AuthWrapper
 * Handles client-side authentication gating.
 * - Protects private routes from unauthenticated users.
 * - Redirects to `/signin` if no user is detected and the route is not public.
 * - Displays a simple loading screen until auth state is resolved.
 *
 * Only used inside the `AuthProvider` context.
 */
function AuthWrapper({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const publicRoutes = ["/signin", "/signup"];
    if (loading) return; 

    if (!user && !publicRoutes.includes(pathname)) {
      router.push("/signin");
    }
  }, [user, loading, pathname, router]);

  if (loading) {
    return (
      <div style={{ textAlign: "center", marginTop: "50px" }}>
        <h1>Loading...</h1>
      </div>
    );
  }

  return <>{children}</>;
}

/**
 * RootLayout
 * The base layout for the entire application.
 * - Wraps the app in global font styles and CSS.
 * - Sets up authentication context and gating via `AuthProvider` and `AuthWrapper`.
 * - Applies Google fonts using Tailwind-compatible CSS variables.
 */
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AuthProvider>
          <AuthWrapper>{children}</AuthWrapper>
        </AuthProvider>
      </body>
    </html>
  );
}
