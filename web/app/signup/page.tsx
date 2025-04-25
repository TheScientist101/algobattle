"use client";

import Link from "next/link";
import { ArrowRightIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { signInWithGoogle, signUp } from "@/firebase/firebaseAuth";
import { useState } from "react";
import { useRouter } from "next/compat/router";
import { ensureUserDocExists } from "@/utils/userData";
import LoadingScreen from "@/components/loading";

export default function SignupPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmedPassword, setConfirmedPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [displayError, setDisplayError] = useState("");
  const router = useRouter();

  const createAccount = async () => {
    if (username && password && password === confirmedPassword) {
      try {
        setLoading(true);
        signUp(username, password);
        ensureUserDocExists();
        if (router) router.push("/");
      } catch (error) {
        console.error(error);
        if (error instanceof Error) {
          setDisplayError(error.message);
        } else {
          setDisplayError("Something went wrong.");
        }
      } finally {
        setLoading(false);
      }
    }
  };

  const googleSignIn = () => {
    try {
      setLoading(true);
      signInWithGoogle();
      ensureUserDocExists();
      if (router) router.push("/");
    } catch (error) {
      console.error(error);
      if (error instanceof Error) {
        setDisplayError(error.message);
      } else {
        setDisplayError("Something went wrong.");
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <LoadingScreen />;

  return (
    <div className="flex min-h-screen w-full items-center justify-center dark">
      <div className="grid w-150 gap-6">
        <div className="flex flex-col items-center gap-2 text-center">
          <div className="flex items-center gap-2 text-2xl font-semibold">
            <ArrowRightIcon className="h-6 w-6" />
            <span>AlgoBattle</span>
          </div>
          <h1 className="text-3xl font-bold text-white">Create an account</h1>
          <p className="text-muted-foreground">
            Enter your information to get started
          </p>
          {displayError && <p className="text-red-50">{displayError}</p>}
        </div>
        <Card>
          <CardHeader>
            <CardTitle>Sign Up</CardTitle>
            <CardDescription>
              Already have an account?{" "}
              <Link
                href="/signin"
                className="underline underline-offset-2 hover:text-primary"
              >
                Login
              </Link>
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                placeholder="johndoe"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirm-password">Confirm Password</Label>
              <Input
                id="confirm-password"
                type="password"
                required
                value={confirmedPassword}
                onChange={(e) => setConfirmedPassword(e.target.value)}
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col">
            <Button className="w-full" size="lg" onClick={createAccount}>
              Create Account
            </Button>

            <div className="relative my-4 w-full">
              <Separator className="my-4" />
              <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-card px-2 text-xs text-muted-foreground">
                OR CONTINUE WITH
              </div>
            </div>

            <Button
              variant="outline"
              className="w-full"
              size="lg"
              onClick={googleSignIn}
            >
              <svg viewBox="0 0 24 24" className="mr-2 h-5 w-5">
                <path
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  fill="#4285F4"
                />
                <path
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  fill="#34A853"
                />
                <path
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  fill="#FBBC05"
                />
                <path
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  fill="#EA4335"
                />
                <path d="M1 1h22v22H1z" fill="none" />
              </svg>
              Sign up with Google
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
