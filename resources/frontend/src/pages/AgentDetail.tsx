import { useEffect, useState } from "react"
import { useParams, useNavigate, useLocation } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"
import {
  Button,
  Text,
  Table,
  TableHeader,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
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
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
    padding: tokens.spacingHorizontalXXL,
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
})

function AgentDetail() {
  const styles = useStyles()
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  useDocumentTitle("终端节点详情")

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

  if (agentLoading || !agent) {
    return (
      <div className={styles.root}>
        <div className={styles.loading}>
          <Spinner label="加载中..." />
        </div>
      </div>
    )
  }

  const p = agent.process_info
  const t = agent.tunnel_info
  const connTotalPages = connData ? Math.ceil(connData.total / connData.size) : 0

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

      <div className={styles.section}>
        <div className={styles.sectionTitle}>
          <Text weight="semibold" size={400}>连接记录</Text>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHeaderCell>主机名</TableHeaderCell>
              <TableHeaderCell>IP 地址</TableHeaderCell>
              <TableHeaderCell>连接时间</TableHeaderCell>
              <TableHeaderCell>断开时间</TableHeaderCell>
              <TableHeaderCell>服务端地址</TableHeaderCell>
              <TableHeaderCell>远程地址</TableHeaderCell>
              <TableHeaderCell>持续时间</TableHeaderCell>
              <TableHeaderCell>接收字节</TableHeaderCell>
              <TableHeaderCell>发送字节</TableHeaderCell>
            </TableRow>
          </TableHeader>
          <TableBody>
            {connLoading && (
              <TableRow>
                <TableCell colSpan={9}>
                  <div className={styles.loading}>
                    <Spinner label="加载中..." />
                  </div>
                </TableCell>
              </TableRow>
            )}
            {!connLoading && connData?.records?.length === 0 && (
              <TableRow>
                <TableCell colSpan={9}>
                  <Text style={{ color: tokens.colorNeutralForeground3, textAlign: "center" as const, display: "block" }}>
                    暂无记录
                  </Text>
                </TableCell>
              </TableRow>
            )}
            {!connLoading && connData?.records?.map((r) => (
              <TableRow key={r.id}>
                <TableCell>{r.process_info.hostname ?? "-"}</TableCell>
                <TableCell>{r.process_info.inet ?? "-"}</TableCell>
                <TableCell>{formatTime(r.tunnel_info.connected_at)}</TableCell>
                <TableCell>{formatTime(r.disconnected_at)}</TableCell>
                <TableCell>{r.tunnel_info.server_addr ?? "-"}</TableCell>
                <TableCell>{r.tunnel_info.remote_addr ?? "-"}</TableCell>
                <TableCell>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    {formatDuration(r.active_seconds)}
                  </Text>
                </TableCell>
                <TableCell>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    {formatBytes(r.tunnel_info.receive_bytes)}
                  </Text>
                </TableCell>
                <TableCell>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    {formatBytes(r.tunnel_info.transmit_bytes)}
                  </Text>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {connData && connTotalPages > 1 && (
          <div className={styles.pagination}>
            <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
              共 {connData.total} 条
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
    </div>
  )
}

export default AgentDetail