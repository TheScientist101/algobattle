"use client";

import { type Icon } from "@tabler/icons-react";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { useRouter } from "next/navigation";


/**
 * Sidebar component for navigating the app and creating new bots.
 *
 * - Renders a list of navigation items with optional icons.
 * - Uses ShadCN and Tabler components for consistent UI.
 *
 * Props:
 * @param items - An array of sidebar links with a title, url, and optional icon.
 *
 *  */
export function NavMain({
  items,
}: {
  items: {
    title: string;
    url: string;
    icon?: Icon;
  }[];
}) {
  const router = useRouter();  

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
