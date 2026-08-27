import { useEffect, useState } from "react"
import { useSearchParams, useNavigate } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"
import {
  Button,
  Input,
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
import { SearchRegular, CircleFilled, CircleRegular } from "@fluentui/react-icons"

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

interface AgentPageResponse {
  page: number
  size: number
  total: number
  records: AgentRecord[]
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
  }
  tunnel_info: {
    connected_at?: string
    keepalive_at?: string
    library_name?: string
    server_addr?: string
    remote_addr?: string
    receive_bytes?: number
    transmit_bytes?: number
  }
  created_at: string
}

const useStyles = makeStyles({
  root: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
    padding: tokens.spacingHorizontalXXL,
  },
  header: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  searchRow: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalS,
  },
  card: {
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 2px 8px rgba(0,0,0,0.06)",
    overflow: "hidden",
  },
  pagination: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  paginationButtons: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalXS,
  },
  statusCell: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalXS,
  },
  tag: {
    display: "inline-flex",
    alignItems: "center",
    borderRadius: tokens.borderRadiusSmall,
    backgroundColor: tokens.colorNeutralBackground2,
    padding: "0 4px",
    fontSize: tokens.fontSizeBase200,
    lineHeight: "20px",
  },
  loading: {
    display: "flex",
    justifyContent: "center",
    padding: tokens.spacingVerticalXXL,
  },
})

function Agent() {
  const styles = useStyles()
  const [searchParams, setSearchParams] = useSearchParams()
  const navigate = useNavigate()
  useDocumentTitle("终端节点")
  const [data, setData] = useState<AgentPageResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [q, setQ] = useState(searchParams.get("q") ?? "")

  const page = Number(searchParams.get("page") ?? 1)
  const size = Number(searchParams.get("size") ?? 10)

  const fetchAgents = async (pageNum: number, pageSize: number, search: string) => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(pageNum), size: String(pageSize) })
      if (search) params.set("q", search)
      const res = await api<AgentPageResponse>("/api/agents?" + params)
      setData(res)
    } catch (err) {
      showError(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchAgents(page, size, q)
  }, [page, size])

  const handleSearch = () => {
    setSearchParams({ page: "1", size: String(size), q })
    fetchAgents(1, size, q)
  }

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  const columns = [
    { key: "hostname", label: "主机名" },
    { key: "inet", label: "IP 地址" },
    { key: "status", label: "状态" },
    { key: "enabled", label: "启用" },
    { key: "version", label: "版本" },
    { key: "os", label: "系统" },
    { key: "created_at", label: "创建时间" },
  ]

  return (
    <div className={styles.root}>
      <div className={styles.header}>
        <Text size={500} weight="semibold">终端节点</Text>
        <div className={styles.searchRow}>
          <Input
            placeholder="搜索..."
            value={q}
            onChange={(e) => setQ(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSearch()}
            style={{ width: "192px" }}
          />
          <Button appearance="outline" icon={<SearchRegular />} onClick={handleSearch} />
        </div>
      </div>

      <div className={styles.card}>
        <Table>
          <TableHeader>
            <TableRow>
              {columns.map((col) => (
                <TableHeaderCell key={col.key}>{col.label}</TableHeaderCell>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading && (
              <TableRow>
                <TableCell colSpan={columns.length}>
                  <div className={styles.loading}>
                    <Spinner label="加载中..." />
                  </div>
                </TableCell>
              </TableRow>
            )}
            {!loading && data && data.records.length === 0 && (
              <TableRow>
                <TableCell colSpan={columns.length}>
                  <Text style={{ color: tokens.colorNeutralForeground3, textAlign: "center" as const, display: "block" }}>
                    暂无数据
                  </Text>
                </TableCell>
              </TableRow>
            )}
            {!loading && data?.records.map((agent) => (
              <TableRow
                key={agent.id}
                onClick={() => navigate("/agent/" + agent.id, { state: agent })}
                style={{ cursor: "pointer" }}
              >
                <TableCell>{agent.process_info.hostname ?? "-"}</TableCell>
                <TableCell>{agent.process_info.inet ?? "-"}</TableCell>
                <TableCell>
                  <div className={styles.statusCell}>
                    {agent.status ? (
                      <CircleFilled style={{ color: tokens.colorStatusSuccessForeground1, fontSize: "8px" }} />
                    ) : (
                      <CircleRegular style={{ color: tokens.colorNeutralForeground3, fontSize: "8px" }} />
                    )}
                    <Text size={200}>{agent.status ? "在线" : "离线"}</Text>
                  </div>
                </TableCell>
                <TableCell>
                  <Text size={200}>{agent.enabled ? "是" : "否"}</Text>
                </TableCell>
                <TableCell>
                  <span className={styles.tag}>{agent.process_info.semver?.version ?? "-"}</span>
                </TableCell>
                <TableCell>
                  <span className={styles.tag}>
                    {(agent.process_info.goos ?? "-") + "/" + (agent.process_info.goarch ?? "-")}
                  </span>
                </TableCell>
                <TableCell>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    {formatTime(agent.created_at)}
                  </Text>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {data && (
        <div className={styles.pagination}>
          <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
            共 {data.total} 条
          </Text>
          <div className={styles.paginationButtons}>
            <Button
              appearance="outline"
              size="small"
              disabled={page <= 1}
              onClick={() => setSearchParams({ page: String(page - 1), size: String(size), q })}
            >
              上一页
            </Button>
            <Text size={200} style={{ padding: "0 " + tokens.spacingHorizontalS, color: tokens.colorNeutralForeground3 }}>
              {page} / {totalPages}
            </Text>
            <Button
              appearance="outline"
              size="small"
              disabled={page >= totalPages}
              onClick={() => setSearchParams({ page: String(page + 1), size: String(size), q })}
            >
              下一页
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}

export default Agent