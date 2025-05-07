"use client";

import * as React from "react";
import { IconDashboard, IconTrophy } from "@tabler/icons-react";

import { NavMain } from "@/components/nav-main";
import { NavUser } from "@/components/nav-user";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useAuth } from "@/hooks/authContext";

/**
 * Renders the app's persistent sidebar for layout and navigation.
 *
 * - Contains a logo/header section at the top.
 * - Provides main navigation items (e.g., Dashboard, Leaderboard).
 * - Displays the current user's avatar and name in the footer.
 *
 * Pulls user data via Firebase Auth (`useAuth`) and populates:
 * - Navigation via `NavMain` (with optional `onBotCreated` handler for bot creation)
 * - User identity via `NavUser` in the sidebar footer
 *
 * Props:
 * @param onBotCreated - Optional callback triggered when a bot is successfully created
 * @param ...props - Additional props passed to the `Sidebar` component (e.g., className)
 */
export function AppSidebar({
  onBotCreated,
  ...props
}: React.ComponentProps<typeof Sidebar> & { onBotCreated?: () => void }) {
  const { user } = useAuth(); // Access authenticated user info

  // Compose sidebar content using user data and static routes
  const data = {
    user: {
      name: user?.displayName || "Hey there!", // Fallback if display name is missing
      avatar: `https://robohash.org/${user?.uid}`, // Dynamically generate avatar using UID
    },
    navMain: [
      {
        title: "Dashboard",
        url: "/",
        icon: IconDashboard,
      },
      {
        title: "Leaderboard",
        url: "/leaderboard",
        icon: IconTrophy,
      },
    ],
  };

  return (
    <Sidebar collapsible="offcanvas" {...props}>
      {/* Top: App name/logo section */}
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <a href="">
                <span className="text-base font-semibold">AlgoBattle</span>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      {/* Middle: Navigation items like Dashboard, Leaderboard */}
      <SidebarContent>
        <NavMain items={data.navMain} onBotCreated={onBotCreated} />
      </SidebarContent>

      {/* Bottom: User avatar and name dropdown */}
      <SidebarFooter>
        <NavUser person={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}