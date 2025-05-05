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
 * Sidebar component for navigating the app and creating new bots.
 *
 * - Includes a "Create a bot" dialog form connected to Firebase.
 * - Renders a list of navigation items with optional icons.
 * - Uses ShadCN and Tabler components for consistent UI.
 *
 * Props:
 * @param items - An array of sidebar links with a title, url, and optional icon.
 *
 * - Integrates with Firebase to create a bot for the logged-in user.
 */
export function NavMain({
  items,
}: {
  items: {
    title: string;
    url: string;
    icon?: Icon;
  }[];
}) {
  const [open, setOpen] = useState(false); // Dialog open/closed state
  const [name, setName] = useState(""); // Controlled input for bot name
  const router = useRouter();
  const { user } = useAuth(); // Firebase auth context

  /**
   * Handles the bot creation process:
   * - Calls Firestore `createBot` with the bot name and current user ID.
   * - Closes the dialog and resets the input on success.
   */
  const createNewBot = async () => {
    if (name) {
      await createBot(name, user?.uid as string);
      setOpen(false);
      setName("");
      setOpen(false);
    }
  };

  /**
   * Navigates to the specified sidebar route.
   * @param url - The destination path.
   */
  const redirect = (url: string) => {
    router.push(url);
  };

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          {/* Bot creation section with modal dialog */}
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

        {/* Navigation items rendered from props that redirect to the url of each item */}
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
