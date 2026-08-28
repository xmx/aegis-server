import { useTheme } from "@/components/ThemeProvider"
import { useAuth } from "@/components/AuthProvider"
import { FontSettings, getFontStyle } from "@/components/FontSettings"
import React, { useState, useRef, useEffect } from "react"
import { api } from "@/lib/api"
import {
  Button,
  Dialog,
  DialogTrigger,
  DialogSurface,
  DialogTitle,
  DialogBody,
  DialogActions,
  Table,
  TableHeader,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
  makeStyles,
  tokens,
  Text,
  Spinner,
} from "@fluentui/react-components"
import {
  WeatherSunnyRegular,
  WeatherMoonRegular,
  PeopleRegular,
  DesktopRegular,
} from "@fluentui/react-icons"
import { Outlet, useNavigate, useLocation } from "react-router-dom"

const useStyles = makeStyles({
  root: {
    display: "flex",
    minHeight: "100vh",
    width: "100%",
    backgroundColor: tokens.colorNeutralBackground2,
    gap: tokens.spacingHorizontalXS,
    padding: tokens.spacingHorizontalXS,
  },
  sidebar: {
    display: "flex",
    flexDirection: "column",
    width: "240px",
    minWidth: "240px",
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 2px 8px rgba(0,0,0,0.06)",
    overflow: "hidden",
  },
  sidebarBrand: {
    display: "flex",
    alignItems: "center",
    padding: tokens.spacingHorizontalL,
    height: "52px",
    background: "none",
    border: "none",
    cursor: "pointer",
    width: "100%",
  },
  sidebarBody: {
    flex: 1,
    padding: `${tokens.spacingVerticalS} ${tokens.spacingHorizontalS}`,
    overflowY: "auto" as const,
  },
  sidebarFooter: {
    padding: tokens.spacingHorizontalM,
    borderTop: `1px solid ${tokens.colorNeutralStroke2}`,
  },
  navItem: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalS,
    padding: `${tokens.spacingVerticalSNudge} ${tokens.spacingHorizontalM}`,
    borderRadius: tokens.borderRadiusMedium,
    fontSize: tokens.fontSizeBase300,
    color: tokens.colorNeutralForeground2,
    cursor: "pointer",
    width: "100%",
    background: "none",
    border: "none",
    textAlign: "left" as const,
    ":hover": {
      backgroundColor: tokens.colorNeutralBackground1Hover,
      color: tokens.colorNeutralForeground1,
    },
  },
  navItemActive: {
    backgroundColor: tokens.colorNeutralBackground1Selected,
    color: tokens.colorNeutralForeground1,
    fontWeight: tokens.fontWeightSemibold,
  },
  navSection: {
    padding: `${tokens.spacingVerticalS} ${tokens.spacingHorizontalM} ${tokens.spacingVerticalXS}`,
    fontSize: tokens.fontSizeBase200,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorNeutralForeground3,
  },
  main: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
    minWidth: 0,
  },
  avatarButton: {
    width: "32px",
    height: "32px",
    minWidth: "32px",
    padding: 0,
    borderRadius: tokens.borderRadiusMedium,
    overflow: "hidden",
    border: `1px solid ${tokens.colorNeutralStroke1}`,
    cursor: "pointer",
  },
  avatarImg: {
    width: "100%",
    height: "100%",
    objectFit: "cover" as const,
  },
  popover: {
    position: "absolute" as const,
    right: 0,
    top: "100%",
    marginTop: tokens.spacingVerticalS,
    zIndex: 50,
    width: "192px",
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 4px 16px rgba(0,0,0,0.12)",
    padding: tokens.spacingVerticalXS,
    border: "none",
  },
  popoverUser: {
    padding: `${tokens.spacingVerticalSNudge} ${tokens.spacingHorizontalS}`,
  },
  popoverName: {
    fontSize: tokens.fontSizeBase300,
    fontWeight: tokens.fontWeightSemibold,
  },
  popoverLogin: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
  },
  popoverDivider: {
    margin: `${tokens.spacingVerticalXS} 0`,
    height: "1px",
    backgroundColor: tokens.colorNeutralStroke2,
  },
  popoverItem: {
    width: "100%",
    padding: `${tokens.spacingVerticalSNudge} ${tokens.spacingHorizontalS}`,
    borderRadius: tokens.borderRadiusMedium,
    fontSize: tokens.fontSizeBase300,
    textAlign: "left" as const,
    cursor: "pointer",
    background: "none",
    border: "none",
    color: tokens.colorNeutralForeground2,
    ":hover": {
      backgroundColor: tokens.colorNeutralBackground1Hover,
      color: tokens.colorNeutralForeground1,
    },
  },
})

function DashboardLayout() {
  const styles = useStyles()
  const { resolvedTheme, setTheme } = useTheme()
  const { user } = useAuth()
  const [cnFont, setCnFont] = useState("system")
  const [enFont, setEnFont] = useState("space-mono")
  const navigate = useNavigate()
  const location = useLocation()
  interface BuildInfo {
  goos: string
  goarch: string
  version: string
  revision: string
  username: string
  workdir: string
  module: string
  committed_at: string
  build_info: {
    GoVersion: string
    Path: string
    Main: { Path: string; Version: string }
    Deps: { Path: string; Version: string }[]
    Settings: { Key: string; Value: string }[]
  }
}
  const [buildinfo, setBuildinfo] = useState<BuildInfo | null>(null)
  const [buildinfoLoading, setBuildinfoLoading] = useState(false)
  const [buildinfoOpen, setBuildinfoOpen] = useState(false)
  const [buildinfoError, setBuildinfoError] = useState(false)

  const handleLogoClick = () => {
    setBuildinfoOpen(true)
    setBuildinfoError(false)
    if (!buildinfoLoading && !buildinfo) {
      setBuildinfoLoading(true)
      api<BuildInfo>("/api/system/buildinfo")
        .then(setBuildinfo)
        .catch(() => setBuildinfoError(true))
        .finally(() => setBuildinfoLoading(false))
    }
  }

  return (
    <div className={styles.root} style={{ "--app-font": getFontStyle(cnFont, enFont) } as React.CSSProperties}>
      <nav className={styles.sidebar}>
        <button className={styles.sidebarBrand} onClick={handleLogoClick}>
          <div>
            <Text weight="semibold" size={400}>Aegis</Text>
            <br />
            <Text size={100} style={{ color: tokens.colorNeutralForeground3 }}>版本 1.0.0</Text>
          </div>
        </button>
        <div className={styles.sidebarBody}>
          <div className={styles.navSection}>系统管理</div>
          <button
            className={`${styles.navItem} ${location.pathname === "/user" ? styles.navItemActive : ""}`}
            onClick={() => navigate("/user")}
          >
            <PeopleRegular />
            <span>系统用户</span>
          </button>
          <button
            className={`${styles.navItem} ${location.pathname === "/agent" ? styles.navItemActive : ""}`}
            onClick={() => navigate("/agent")}
          >
            <DesktopRegular />
            <span>终端节点</span>
          </button>
        </div>
        <div className={styles.sidebarFooter}>
          <div style={{ display: "flex", alignItems: "center", gap: tokens.spacingHorizontalXS, marginBottom: tokens.spacingVerticalXS }}>
            {user && <UserMenu user={user} styles={styles} />}
            <Button
              appearance="subtle"
              size="small"
              icon={resolvedTheme === "dark" ? <WeatherSunnyRegular /> : <WeatherMoonRegular />}
              onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}
            />
            <FontSettings cnFont={cnFont} enFont={enFont} onCnFontChange={setCnFont} onEnFontChange={setEnFont} />
          </div>
        </div>
      </nav>

      <Dialog open={buildinfoOpen} onOpenChange={(_, d) => setBuildinfoOpen(d.open)}>
        <DialogSurface>
          <DialogBody>
            <DialogTitle>Aegis 构建信息</DialogTitle>
            {buildinfoLoading ? (
              <div style={{ display: "flex", justifyContent: "center", padding: tokens.spacingVerticalXL }}>
                <Spinner label="加载中..." />
              </div>
            ) : buildinfoError ? (
              <Text style={{ color: tokens.colorNeutralForeground3 }}>后端接口尚未就绪</Text>
            ) : buildinfo ? (
              <div style={{ display: "flex", flexDirection: "column", gap: tokens.spacingVerticalL, maxHeight: "60vh", overflowY: "auto" }}>
                <div>
                  <Text weight="semibold" size={300}>基本信息</Text>
                  <div style={{ display: "grid", gridTemplateColumns: "auto 1fr", gap: "4px 16px", marginTop: tokens.spacingVerticalS }}>
                    {[
                      ["版本", buildinfo.version],
                      ["修订", buildinfo.revision],
                      ["用户", buildinfo.username],
                      ["工作目录", buildinfo.workdir],
                      ["模块", buildinfo.module],
                      ["提交时间", buildinfo.committed_at],
                      ["OS", buildinfo.goos + "/" + buildinfo.goarch],
                    ].map(([label, value]) => (
                      <React.Fragment key={label}>
                        <Text size={200} style={{ color: tokens.colorNeutralForeground3, whiteSpace: "nowrap" }}>{label}</Text>
                        <Text size={200} style={{ wordBreak: "break-all" }}>{value}</Text>
                      </React.Fragment>
                    ))}
                  </div>
                </div>

                <div>
                  <Text weight="semibold" size={300}>构建信息</Text>
                  <div style={{ display: "grid", gridTemplateColumns: "auto 1fr", gap: "4px 16px", marginTop: tokens.spacingVerticalS }}>
                    {[
                      ["Go 版本", buildinfo.build_info.GoVersion],
                      ["主模块", buildinfo.build_info.Main.Path + " " + buildinfo.build_info.Main.Version],
                    ].map(([label, value]) => (
                      <React.Fragment key={label}>
                        <Text size={200} style={{ color: tokens.colorNeutralForeground3, whiteSpace: "nowrap" }}>{label}</Text>
                        <Text size={200} style={{ wordBreak: "break-all" }}>{value}</Text>
                      </React.Fragment>
                    ))}
                  </div>
                </div>

                <div>
                  <Text weight="semibold" size={300}>依赖项 ({buildinfo.build_info.Deps.length})</Text>
                  <Table size="small" style={{ marginTop: tokens.spacingVerticalS }}>
                    <TableHeader>
                      <TableRow>
                        <TableHeaderCell>模块</TableHeaderCell>
                        <TableHeaderCell>版本</TableHeaderCell>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {buildinfo.build_info.Deps.map((d) => (
                        <TableRow key={d.Path}>
                          <TableCell style={{ wordBreak: "break-all" }}>
                            <Text size={200}>{d.Path}</Text>
                          </TableCell>
                          <TableCell>
                            <Text size={200} font="monospace">{d.Version}</Text>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>

                <div>
                  <Text weight="semibold" size={300}>构建参数</Text>
                  <Table size="small" style={{ marginTop: tokens.spacingVerticalS }}>
                    <TableHeader>
                      <TableRow>
                        <TableHeaderCell>Key</TableHeaderCell>
                        <TableHeaderCell>Value</TableHeaderCell>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {buildinfo.build_info.Settings.map((s) => (
                        <TableRow key={s.Key}>
                          <TableCell>
                            <Text size={200} font="monospace">{s.Key}</Text>
                          </TableCell>
                          <TableCell>
                            <Text size={200}>{s.Value}</Text>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              </div>
            ) : null}
          </DialogBody>
          <DialogActions>
            <DialogTrigger disableButtonEnhancement>
              <Button appearance="secondary">关闭</Button>
            </DialogTrigger>
          </DialogActions>
        </DialogSurface>
      </Dialog>

      <div className={styles.main}>
        <main style={{ flex: 1, display: "flex", flexDirection: "column" }}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}

function UserMenu({ user, styles }: { user: NonNullable<ReturnType<typeof useAuth>["user"]>; styles: ReturnType<typeof useStyles> }) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    if (open) document.addEventListener("click", handleClick)
    return () => document.removeEventListener("click", handleClick)
  }, [open])

  return (
    <div ref={ref} style={{ position: "relative" }}>
      <button className={styles.avatarButton} onClick={() => setOpen(!open)}>
        <img src={user.avatar_url} alt={user.login} className={styles.avatarImg} />
      </button>
      {open && (
        <div className={styles.popover}>
          <div className={styles.popoverUser}>
            <div className={styles.popoverName}>{user.name ?? user.login}</div>
            <div className={styles.popoverLogin}>@{user.login}</div>
          </div>
          <div className={styles.popoverDivider} />
          <button className={styles.popoverItem}>个人设置</button>
          <button className={styles.popoverItem}>退出登录</button>
        </div>
      )}
    </div>
  )
}

export default DashboardLayout