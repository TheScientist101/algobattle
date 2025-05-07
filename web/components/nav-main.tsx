"use client";

import { type Icon } from "@tabler/icons-react";
import { Button } from "@/components/ui/button";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Input } from "./ui/input";
import { createBot } from "@/utils/botData";
import { useAuth } from "@/hooks/authContext";
import { useRouter } from "next/navigation";

/**
 * Sidebar component that:
 * - Provides navigational access to different parts of the app via provided `items`.
 * - Includes a dialog interface for creating new bots linked to the current user via Firebase.
 * - Uses ShadCN UI components and Tabler icons for a polished and consistent look.
 * - Accepts an optional `onBotCreated` callback to notify parent components when a bot is created.
 *
 * Props:
 * @param items - Array of navigation entries with `title`, `url`, and optional `icon`.
 * @param onBotCreated - Optional callback invoked after a bot is successfully created.
 */
export function NavMain({
  items,
  onBotCreated,
}: {
  items: {
    title: string;
    url: string;
    icon?: Icon;
  }[];
  onBotCreated?: () => void;
}) {
  const [open, setOpen] = useState(false); // Controls whether the "Create a bot" dialog is open
  const [name, setName] = useState(""); // Stores input for the bot name
  const router = useRouter();
  const { user } = useAuth(); // Gets current authenticated user from custom auth context

  /**
   * Creates a new bot by:
   * - Sending the name and user ID to Firestore via `createBot()`.
   * - Resetting dialog state and invoking optional callback.
   */
  const createNewBot = async () => {
    if (name) {
      await createBot(name, user?.uid as string); 
      setOpen(false); 
      setName(""); 
      onBotCreated?.();
    }
  };

  /**
   * Navigates to a new route based on the provided `url`.
   * Used by sidebar menu buttons.
   *
   * @param url - The route to navigate to
   */
  const redirect = (url: string) => {
    router.push(url);
  };

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">

        {/* Section: Create Bot Dialog Trigger */}
        <SidebarMenu>
          <SidebarMenuItem className="flex items-center gap-2">
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">Create a bot</Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-[425px] rounded-2xl bg-zinc-900 shadow-xl p-6">
                <DialogHeader>
                  <DialogTitle className="text-white text-lg font-semibold">
                    Create your bot
                  </DialogTitle>
                  <DialogDescription className="text-zinc-300">
                    Enter the name of the bot
                  </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="flex flex-col gap-2">
                    <label
                      htmlFor="name"
                      className="text-white text-sm font-medium"
                    >
                      Bot Name
                    </label>
                    <Input
                      id="name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="w-full bg-zinc-800 border border-zinc-700 text-white placeholder-zinc-500 focus:ring-primary focus:border-primary"
                      placeholder="Enter bot name"
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button
                    type="button"
                    onClick={createNewBot}
                    className="bg-primary text-white hover:bg-primary/90 transition-all"
                  >
                    Save
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </SidebarMenuItem>
        </SidebarMenu>

        {/* Section: Sidebar Navigation Items */}
        <SidebarMenu>
          {items.map((item) => (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton
                tooltip={item.title}
                onClick={() => redirect(item.url)}
              >
                {item.icon && <item.icon />}
                <span>{item.title}</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}