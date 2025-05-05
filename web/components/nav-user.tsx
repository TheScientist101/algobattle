"use client"

import {
  IconDotsVertical,
  IconLogout,
} from "@tabler/icons-react"

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { logOut } from "@/firebase/firebaseAuth"

/**
 * Component of the user section on the bottom of the sidebar
 * Displays user's avatar and name.
 * Provides dropdown actions including logout.
 *
 * Props:
 * - person: an object containing:
 *    - name: user's display name
 *    - avatar: URL of user's avatar image
 */
export function NavUser({
  person,
}: {
  person: {
    name: string
    avatar: string
  }
}) {
  const { isMobile } = useSidebar(); // Determine dropdown direction based on layout

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          {/* Menu trigger: full-width button showing avatar, name, and menu icon */}
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              {/* User avatar with fallback initials */}
              <Avatar className="h-8 w-8 rounded-lg grayscale">
                <AvatarImage src={person.avatar} alt={person.name} />
                <AvatarFallback className="rounded-lg">CN</AvatarFallback>
              </Avatar>

              {/* User name display */}
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{person.name}</span>
              </div>

              {/* Vertical dots icon indicating dropdown */}
              <IconDotsVertical className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>

          {/* Dropdown menu content, positioned based on device */}
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg dark"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            {/* Top label showing user info */}
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage src={person.avatar} alt={person.name} />
                  <AvatarFallback className="rounded-lg">CN</AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{person.name}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {/* Logout button, calls Firebase sign-out */}
            <DropdownMenuItem onClick={logOut}>
              <IconLogout />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}