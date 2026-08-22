import { Command, Moon, Sun, Users } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useTheme } from "@/components/ThemeProvider"
import { useAuth } from "@/components/AuthProvider"
import { FontSettings, getFontStyle } from "@/components/FontSettings"
import { useState, useRef, useEffect } from "react"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { Outlet, useLocation, Link } from "react-router-dom"

function DashboardLayout() {
  const { resolvedTheme, setTheme } = useTheme()
  const { user } = useAuth()
  const [cnFont, setCnFont] = useState("system")
  const [enFont, setEnFont] = useState("space-mono")

  return (
    <SidebarProvider>
      <div
        className="flex min-h-screen w-full"
        style={{ fontFamily: getFontStyle(cnFont, enFont) }}
      >
        <AppSidebar />
        <div className="flex flex-1 flex-col">
          <header className="flex h-12 items-center gap-3 border-b px-4">
            <SidebarTrigger />
            <div className="flex-1" />
            {user && <UserMenu user={user} />}
            <Button
              variant="outline"
              size="icon"
              onClick={() =>
                setTheme(resolvedTheme === "dark" ? "light" : "dark")
              }
            >
              {resolvedTheme === "dark" ? (
                <Sun className="size-4" />
              ) : (
                <Moon className="size-4" />
              )}
            </Button>
            <FontSettings
              cnFont={cnFont}
              enFont={enFont}
              onCnFontChange={setCnFont}
              onEnFontChange={setEnFont}
            />
          </header>
          <main className="flex flex-1 flex-col">
            <Outlet />
          </main>
        </div>
      </div>
    </SidebarProvider>
  )
}

function AppSidebar() {
  const { pathname } = useLocation()

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              render={
                <Link to="/">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                    <Command className="size-4" />
                  </div>
                  <div className="flex flex-col gap-0.5 leading-none">
                    <span className="font-semibold">Aegis</span>
                    <span className="text-xs text-muted-foreground">
                      企业版
                    </span>
                  </div>
                </Link>
              }
            />
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>系统管理</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  isActive={pathname === "/user"}
                  render={<Link to="/user" />}
                >
                  <Users className="size-4" />
                  <span>系统用户</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton>
              <span>版本 1.0.0</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

function UserMenu({ user }: { user: NonNullable<ReturnType<typeof useAuth>["user"]> }) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    if (open) {
      document.addEventListener("click", handleClick)
    }
    return () => document.removeEventListener("click", handleClick)
  }, [open])

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        className="inline-flex size-8 cursor-pointer items-center justify-center rounded-lg border border-border bg-background outline-none hover:bg-muted"
        onClick={() => setOpen(!open)}
      >
        <img
          src={user.avatar_url}
          alt={user.login}
          className="size-5 rounded-full"
        />
      </button>
      {open && (
        <div className="absolute right-0 top-full z-50 mt-2 w-48 rounded-lg border bg-popover p-1 text-popover-foreground shadow-md">
          <div className="px-2 py-1.5">
            <div className="text-sm font-medium">{user.name ?? user.login}</div>
            <div className="text-xs text-muted-foreground">@{user.login}</div>
          </div>
          <div className="my-1 h-px bg-border" />
          <button
            type="button"
            className="w-full rounded-sm px-2 py-1.5 text-left text-sm outline-none hover:bg-accent hover:text-accent-foreground"
          >
            个人设置
          </button>
          <button
            type="button"
            className="w-full rounded-sm px-2 py-1.5 text-left text-sm outline-none hover:bg-accent hover:text-accent-foreground"
          >
            退出登录
          </button>
        </div>
      )}
    </div>
  )
}

export default DashboardLayout