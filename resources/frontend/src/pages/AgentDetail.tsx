import { useEffect, useMemo, useState } from "react"
import { useParams, useNavigate, useLocation } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"
import {
  Button,
  Text,
  makeStyles,
  tokens,
  Spinner,
} from "@fluentui/react-components"
import { ArrowLeftRegular } from "@fluentui/react-icons"

function formatTime(iso: string | undefined): string {
  if (!iso) return "-"
  const d = new Date(iso)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, "0")
  const day = String(d.getDate()).padStart(2, "0")
  const h = String(d.getHours()).padStart(2, "0")
  const min = String(d.getMinutes()).padStart(2, "0")
  const s = String(d.getSeconds()).padStart(2, "0")
  return y + "-" + m + "-" + day + " " + h + ":" + min + ":" + s
}

function formatDuration(seconds: number | undefined): string {
  if (!seconds || seconds <= 0) return "-"
  if (seconds < 60) return seconds + "s"
  if (seconds < 3600) return Math.floor(seconds / 60) + "m"
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (h < 24) return h + "h " + m + "m"
  const d = Math.floor(h / 24)
  const rh = h % 24
  return d + "d " + rh + "h"
}

function formatBytes(bytes: number | undefined): string {
  if (!bytes || bytes <= 0) return "-"
  if (bytes < 1024) return bytes + " B"
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB"
  if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(1) + " MB"
  return (bytes / (1024 * 1024 * 1024)).toFixed(1) + " GB"
}

interface AgentRecord {
  id: string
  machine_id: string
  status: boolean
  enabled: boolean
  process_info: {
    inet?: string
    hostname?: string
    goos?: string
    goarch?: string
    semver?: { version: string; number: string }
    pid?: number
    workdir?: string
    executable?: string
  }
  tunnel_info: {
    connected_at?: string
    keepalive_at?: string
    server_addr?: string
    remote_addr?: string
    library_name?: string
    library_module?: string
    receive_bytes?: number
    transmit_bytes?: number
  }
  created_at: string
}

interface ConnRecord {
  id: string
  agent_id: string
  machine_id: string
  active_seconds?: number
  disconnected_at?: string
  process_info: {
    inet?: string
    hostname?: string
    semver?: { version: string; number: string }
  }
  tunnel_info: {
    connected_at?: string
    keepalive_at?: string
    server_addr?: string
    remote_addr?: string
    library_name?: string
    library_module?: string
    receive_bytes?: number
    transmit_bytes?: number
  }
  created_at: string
}

type ConnPageResponse = {
  page: number
  size: number
  total: number
  records: ConnRecord[]
}

const useStyles = makeStyles({
  root: {
    flex: 1,
    display: "flex",
    gap: tokens.spacingHorizontalXS,
    padding: tokens.spacingHorizontalXXL,
  },
  nav: {
    display: "flex",
    flexDirection: "column",
    width: "180px",
    minWidth: "180px",
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 2px 8px rgba(0,0,0,0.06)",
    padding: tokens.spacingVerticalS,
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
  content: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
  },
  backRow: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalS,
  },
  section: {
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 2px 8px rgba(0,0,0,0.06)",
  },
  sectionTitle: {
    padding: tokens.spacingHorizontalL,
    borderBottom: "1px solid " + tokens.colorNeutralStroke2,
  },
  grid: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr",
    gap: tokens.spacingVerticalL,
    padding: tokens.spacingHorizontalL,
  },
  field: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalXXS,
  },
  label: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
  },
  value: {
    fontSize: tokens.fontSizeBase300,
  },
  pagination: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    padding: tokens.spacingHorizontalL,
  },
  paginationButtons: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalXS,
  },
  loading: {
    display: "flex",
    justifyContent: "center",
    padding: tokens.spacingVerticalXXL,
  },
  tableWrap: {
    overflowX: "auto" as const,
  },
  denseTable: {
    width: "100%",
    borderCollapse: "collapse" as const,
    fontSize: tokens.fontSizeBase200,
    whiteSpace: "nowrap" as const,
  },
  th: {
    padding: `${tokens.spacingVerticalXS} ${tokens.spacingHorizontalS}`,
    textAlign: "left" as const,
    fontSize: tokens.fontSizeBase200,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorNeutralForeground3,
    borderBottom: `1px solid ${tokens.colorNeutralStroke2}`,
    backgroundColor: tokens.colorNeutralBackground2,
    position: "sticky" as const,
    top: 0,
  },
  thNum: {
    textAlign: "right" as const,
  },
  td: {
    padding: `${tokens.spacingVerticalXXS} ${tokens.spacingHorizontalS}`,
    borderBottom: `1px solid ${tokens.colorNeutralStroke3}`,
    lineHeight: "20px",
  },
  tdNum: {
    textAlign: "right" as const,
    fontVariantNumeric: "tabular-nums",
    color: tokens.colorNeutralForeground2,
  },
  tdMono: {
    fontVariantNumeric: "tabular-nums",
    color: tokens.colorNeutralForeground2,
  },
  row: {
    ":hover": {
      backgroundColor: tokens.colorNeutralBackground1Hover,
    },
  },
  statusDot: {
    display: "inline-block",
    width: "6px",
    height: "6px",
    borderRadius: "50%",
    marginRight: tokens.spacingHorizontalXS,
    verticalAlign: "middle",
  },
  dotOnline: {
    backgroundColor: tokens.colorStatusSuccessForeground1,
  },
  dotOffline: {
    backgroundColor: tokens.colorNeutralForeground4,
  },
  summary: {
    display: "flex",
    flexWrap: "wrap" as const,
    gap: tokens.spacingHorizontalL,
    padding: `${tokens.spacingVerticalS} ${tokens.spacingHorizontalL}`,
    borderBottom: `1px solid ${tokens.colorNeutralStroke2}`,
    backgroundColor: tokens.colorNeutralBackground2,
  },
  summaryItem: {
    display: "flex",
    alignItems: "baseline",
    gap: tokens.spacingHorizontalXS,
  },
  summaryLabel: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
  },
  summaryValue: {
    fontSize: tokens.fontSizeBase300,
    fontWeight: tokens.fontWeightSemibold,
    fontVariantNumeric: "tabular-nums",
  },
  sectionHeadRow: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    padding: tokens.spacingHorizontalL,
    borderBottom: `1px solid ${tokens.colorNeutralStroke2}`,
  },
})

function AgentDetail() {
  const styles = useStyles()
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  useDocumentTitle("终端节点详情")
  const [activeTab, setActiveTab] = useState<"info" | "connections">("info")

  const [agent, setAgent] = useState<AgentRecord | null>(location.state as AgentRecord | null)
  const [agentLoading, setAgentLoading] = useState(!location.state)
  const [connData, setConnData] = useState<ConnPageResponse | null>(null)
  const [connLoading, setConnLoading] = useState(false)
  const [connPage, setConnPage] = useState(1)

  useEffect(() => {
    if (!id) {
      navigate("/agent", { replace: true })
      return
    }
  }, [id, navigate])

  useEffect(() => {
    if (agent || !id) return
    api<AgentRecord>("/api/agent?id=" + id)
      .then(setAgent)
      .catch(showError)
      .finally(() => setAgentLoading(false))
  }, [id, agent])

  useEffect(() => {
    if (!id) return
    setConnLoading(true)
    const params = "?id=" + id + "&page=" + connPage + "&size=10"
    api<ConnPageResponse>("/api/agent/records" + params)
      .then(setConnData)
      .catch(showError)
      .finally(() => setConnLoading(false))
  }, [id, connPage])

  const connStats = useMemo(() => {
    const records = connData?.records
    if (!records || records.length === 0) return null
    return {
      count: records.length,
      online: records.filter((r) => !r.disconnected_at).length,
      totalSeconds: records.reduce((sum, r) => sum + (r.active_seconds ?? 0), 0),
      rx: records.reduce((sum, r) => sum + (r.tunnel_info.receive_bytes ?? 0), 0),
      tx: records.reduce((sum, r) => sum + (r.tunnel_info.transmit_bytes ?? 0), 0),
    }
  }, [connData])

  if (agentLoading) {
    return (
      <div className={styles.root}>
        <div className={styles.loading}>
          <Spinner label="加载中..." />
        </div>
      </div>
    )
  }

  if (!agent) {
    return (
      <div className={styles.root}>
        <div className={styles.backRow}>
          <Button
            appearance="subtle"
            icon={<ArrowLeftRegular />}
            onClick={() => navigate("/agent")}
          >
            返回
          </Button>
        </div>
        <Text size={400} style={{ color: tokens.colorNeutralForeground3 }}>未找到该终端节点</Text>
      </div>
    )
  }

  const p = agent.process_info
  const t = agent.tunnel_info
  const connTotalPages = connData ? Math.ceil(connData.total / connData.size) : 0

  return (
    <div className={styles.root}>
      <nav className={styles.nav}>
        <button
          className={`${styles.navItem} ${activeTab === "info" ? styles.navItemActive : ""}`}
          onClick={() => setActiveTab("info")}
        >
          基本信息
        </button>
        <button
          className={`${styles.navItem} ${activeTab === "connections" ? styles.navItemActive : ""}`}
          onClick={() => setActiveTab("connections")}
        >
          连接记录
        </button>
      </nav>

      <div className={styles.content}>
        <div className={styles.backRow}>
          <Button
            appearance="subtle"
            icon={<ArrowLeftRegular />}
            onClick={() => navigate("/agent")}
          >
            返回
          </Button>
        </div>

        {activeTab === "info" ? (
          <>
      <div className={styles.section}>
        <div className={styles.sectionTitle}>
          <Text weight="semibold" size={400}>基本信息</Text>
        </div>
        <div className={styles.grid}>
          <div className={styles.field}>
            <span className={styles.label}>主机名</span>
            <span className={styles.value}>{p.hostname ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>IP 地址</span>
            <span className={styles.value}>{p.inet ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>机器码</span>
            <span className={styles.value}>{agent.machine_id}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>状态</span>
            <span className={styles.value}>{agent.status ? "在线" : "离线"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>操作系统</span>
            <span className={styles.value}>{p.goos ?? "-"} / {p.goarch ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>Agent 版本</span>
            <span className={styles.value}>{p.semver?.version ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>进程 PID</span>
            <span className={styles.value}>{p.pid ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>工作目录</span>
            <span className={styles.value}>{p.workdir ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>可执行文件</span>
            <span className={styles.value}>{p.executable ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>创建时间</span>
            <span className={styles.value}>{formatTime(agent.created_at)}</span>
          </div>
        </div>
      </div>

      <div className={styles.section}>
        <div className={styles.sectionTitle}>
          <Text weight="semibold" size={400}>隧道信息</Text>
        </div>
        <div className={styles.grid}>
          <div className={styles.field}>
            <span className={styles.label}>库名称</span>
            <span className={styles.value}>{t.library_name ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>库模块</span>
            <span className={styles.value}>{t.library_module ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>服务端地址</span>
            <span className={styles.value}>{t.server_addr ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>远程地址</span>
            <span className={styles.value}>{t.remote_addr ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>接收字节</span>
            <span className={styles.value}>{t.receive_bytes ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>发送字节</span>
            <span className={styles.value}>{t.transmit_bytes ?? "-"}</span>
          </div>
          <div className={styles.field}>
            <span className={styles.label}>最后保活</span>
            <span className={styles.value}>{formatTime(t.keepalive_at)}</span>
          </div>
        </div>
      </div>
          </>
        ) : (
          <>
            <div className={styles.section}>
              <div className={styles.sectionHeadRow}>
                <Text weight="semibold" size={400}>连接记录</Text>
                <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                  共 {connData?.total ?? 0} 条
                </Text>
              </div>

              {connStats && (
                <div className={styles.summary}>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>本页连接</span>
                    <span className={styles.summaryValue}>{connStats.count}</span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>在线</span>
                    <span className={styles.summaryValue}>{connStats.online}</span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>累计时长</span>
                    <span className={styles.summaryValue}>{formatDuration(connStats.totalSeconds)}</span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>接收</span>
                    <span className={styles.summaryValue}>{formatBytes(connStats.rx)}</span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>发送</span>
                    <span className={styles.summaryValue}>{formatBytes(connStats.tx)}</span>
                  </div>
                </div>
              )}

              {connLoading ? (
                <div className={styles.loading}>
                  <Spinner label="加载中..." />
                </div>
              ) : connData?.records?.length === 0 ? (
                <div className={styles.loading}>
                  <Text style={{ color: tokens.colorNeutralForeground3 }}>暂无记录</Text>
                </div>
              ) : (
                <div className={styles.tableWrap}>
                  <table className={styles.denseTable}>
                    <thead>
                      <tr>
                        <th className={styles.th}>状态 / 主机名</th>
                        <th className={styles.th}>IP 地址</th>
                        <th className={styles.th}>版本</th>
                        <th className={styles.th}>连接时间</th>
                        <th className={styles.th}>断开时间</th>
                        <th className={`${styles.th} ${styles.thNum}`}>持续</th>
                        <th className={styles.th}>服务端地址</th>
                        <th className={styles.th}>远程地址</th>
                        <th className={styles.th}>库</th>
                        <th className={`${styles.th} ${styles.thNum}`}>接收</th>
                        <th className={`${styles.th} ${styles.thNum}`}>发送</th>
                      </tr>
                    </thead>
                    <tbody>
                      {connData?.records?.map((r) => {
                        const online = !r.disconnected_at
                        return (
                          <tr key={r.id} className={styles.row}>
                            <td className={styles.td} title={r.machine_id}>
                              <span
                                className={`${styles.statusDot} ${online ? styles.dotOnline : styles.dotOffline}`}
                              />
                              {r.process_info.hostname ?? "-"}
                            </td>
                            <td className={`${styles.td} ${styles.tdMono}`}>{r.process_info.inet ?? "-"}</td>
                            <td className={`${styles.td} ${styles.tdMono}`}>{r.process_info.semver?.version ?? "-"}</td>
                            <td className={`${styles.td} ${styles.tdMono}`}>{formatTime(r.tunnel_info.connected_at)}</td>
                            <td className={`${styles.td} ${styles.tdMono}`}>
                              {online ? "—" : formatTime(r.disconnected_at)}
                            </td>
                            <td className={`${styles.td} ${styles.tdNum}`}>{formatDuration(r.active_seconds)}</td>
                            <td className={`${styles.td} ${styles.tdMono}`}>{r.tunnel_info.server_addr ?? "-"}</td>
                            <td className={`${styles.td} ${styles.tdMono}`}>{r.tunnel_info.remote_addr ?? "-"}</td>
                            <td className={styles.td} title={r.tunnel_info.library_module}>
                              {r.tunnel_info.library_name ?? "-"}
                            </td>
                            <td className={`${styles.td} ${styles.tdNum}`}>{formatBytes(r.tunnel_info.receive_bytes)}</td>
                            <td className={`${styles.td} ${styles.tdNum}`}>{formatBytes(r.tunnel_info.transmit_bytes)}</td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}

              {connData && connTotalPages > 1 && (
                <div className={styles.pagination}>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    第 {(connData.page - 1) * connData.size + 1}-
                    {Math.min(connData.page * connData.size, connData.total)} 条，共 {connData.total} 条
                  </Text>
                  <div className={styles.paginationButtons}>
                    <Button
                      appearance="outline"
                      size="small"
                      disabled={connPage <= 1}
                      onClick={() => setConnPage(connPage - 1)}
                    >
                      上一页
                    </Button>
                    <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                      {connPage} / {connTotalPages}
                    </Text>
                    <Button
                      appearance="outline"
                      size="small"
                      disabled={connPage >= connTotalPages}
                      onClick={() => setConnPage(connPage + 1)}
                    >
                      下一页
                    </Button>
                  </div>
                </div>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default AgentDetail