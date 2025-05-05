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
 * Renders the persistent sidebar used throughout the app.
 * 
 * - Displays a top logo/header
 * - Provides main navigation links ( Dashboard, Leaderboard)
 * - Shows the currently authenticated user's avatar and name in the footer
 * 
 * Pulls user data from Firebase Auth via `useAuth` and passes props to `Sidebar`.
 */
export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { user } = useAuth();

  // Prepare sidebar data for navigation and user section
  const data = {
    user: {
      name: user?.displayName || "Hey there!", 
      avatar: `https://robohash.org/${user?.uid}`, 
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
      {/* Sidebar top section with app name */}
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

      {/* Sidebar middle section with navigation links */}
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>

      {/* Sidebar bottom section with user dropdown menu */}
      <SidebarFooter>
        <NavUser person={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}