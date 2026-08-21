import {
  BookOpen,
  Bot,
  ChevronRight,
  Command,
  Frame,
  LifeBuoy,
  Map,
  Moon,
  PieChart,
  Send,
  Settings2,
  SquareTerminal,
  Sun,
} from "lucide-react"
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
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { Outlet } from "react-router-dom"

type NavSubItem = { title: string; url: string }
type NavItem = {
  title: string
  url: string
  isActive?: boolean
  subItems?: NavSubItem[]
}

const navMain: {
  title: string
  icon: React.ComponentType<{ className?: string }>
  isActive?: boolean
  items: NavItem[]
}[] = [
  {
    title: "仪表盘",
    icon: SquareTerminal,
    isActive: true,
    items: [
      {
        title: "概览",
        url: "#",
        isActive: true,
        subItems: [
          { title: "今日概要", url: "#" },
          { title: "实时流量", url: "#" },
        ],
      },
      { title: "实时监控", url: "#" },
      { title: "统计分析", url: "#" },
    ],
  },
  {
    title: "项目管理",
    icon: Bot,
    items: [
      {
        title: "项目列表",
        url: "#",
        subItems: [
          { title: "进行中", url: "#" },
          { title: "已归档", url: "#" },
        ],
      },
      { title: "任务看板", url: "#" },
      { title: "甘特图", url: "#" },
    ],
  },
  {
    title: "文档中心",
    icon: BookOpen,
    items: [
      { title: "API 文档", url: "#" },
      {
        title: "使用指南",
        url: "#",
        subItems: [
          { title: "快速入门", url: "#" },
          { title: "进阶教程", url: "#" },
          { title: "常见问题", url: "#" },
        ],
      },
      { title: "更新日志", url: "#" },
    ],
  },
  {
    title: "系统设置",
    icon: Settings2,
    items: [
      {
        title: "用户管理",
        url: "#",
        subItems: [
          { title: "用户列表", url: "#" },
          { title: "邀请成员", url: "#" },
        ],
      },
      {
        title: "角色权限",
        url: "#",
        subItems: [
          { title: "角色列表", url: "#" },
          { title: "权限模板", url: "#" },
        ],
      },
      { title: "系统配置", url: "#" },
    ],
  },
]

const navSecondary = [
  { title: "帮助中心", url: "#", icon: LifeBuoy },
  { title: "反馈建议", url: "#", icon: Send },
]

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
  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              render={
                <a href="/">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                    <Command className="size-4" />
                  </div>
                  <div className="flex flex-col gap-0.5 leading-none">
                    <span className="font-semibold">Aegis</span>
                    <span className="text-xs text-muted-foreground">
                      企业版
                    </span>
                  </div>
                </a>
              }
            />
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>平台导航</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton>
                  <Map className="size-4" />
                  <span>全局搜索</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton>
                  <Frame className="size-4" />
                  <span>快捷入口</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton>
                  <PieChart className="size-4" />
                  <span>数据报表</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {navMain.map((group) => (
          <SidebarGroup key={group.title}>
            <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      isActive={item.isActive}
                      data-collapsed-icon={item.title.charAt(0)}
                    >
                      {item.title}
                    </SidebarMenuButton>
                    {item.subItems && item.subItems.length > 0 && (
                      <SidebarMenuSub>
                        {item.subItems.map((sub) => (
                          <SidebarMenuSubItem key={sub.title}>
                            <SidebarMenuSubButton
                              data-collapsed-icon={sub.title.charAt(0)}
                            >
                              {sub.title}
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    )}
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
        <SidebarGroup>
          <SidebarGroupLabel>其他</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navSecondary.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton>
                    <item.icon className="size-4" />
                    <span>{item.title}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton>
              <ChevronRight className="size-4" />
              <span>版本 1.0.0</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

export default DashboardLayout

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