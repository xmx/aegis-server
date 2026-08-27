import { useEffect, useState } from "react"
import { useSearchParams } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { ProviderIcon } from "@/components/Icons"
import { useDocumentTitle } from "@/hooks/useDocumentTitle"
import {
  Button,
  Input,
  Text,
  Tooltip,
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
import { SearchRegular, CheckmarkCircleRegular, DismissCircleRegular } from "@fluentui/react-icons"

function formatTime(iso: string): string {
  const d = new Date(iso)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, "0")
  const day = String(d.getDate()).padStart(2, "0")
  const h = String(d.getHours()).padStart(2, "0")
  const min = String(d.getMinutes()).padStart(2, "0")
  const s = String(d.getSeconds()).padStart(2, "0")
  return `${y}-${m}-${day} ${h}:${min}:${s}`
}

interface UserRecord {
  id: string
  enabled: boolean
  provider: string
  puid: string
  login: string
  name?: string
  avatar_url?: string
  email?: string
  created_at: string
}

interface UserPageResponse {
  page: number
  size: number
  total: number
  records: UserRecord[]
}

const useStyles = makeStyles({
  root: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
    padding: tokens.spacingHorizontalXXL,
  },
  card: {
    borderRadius: tokens.borderRadiusMedium,
    backgroundColor: tokens.colorNeutralBackground1,
    boxShadow: "0 2px 8px rgba(0,0,0,0.06)",
    overflow: "hidden",
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
  avatar: {
    width: "36px",
    height: "36px",
    borderRadius: "50%",
    objectFit: "cover" as const,
    border: `2px solid ${tokens.colorNeutralStroke1}`,
  },
  statusIcon: {
    display: "inline-flex",
    alignItems: "center",
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
  loading: {
    display: "flex",
    justifyContent: "center",
    padding: tokens.spacingVerticalXXL,
  },
  providerCell: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalXS,
  },
  loginCell: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalXS,
  },
})

function User() {
  const styles = useStyles()
  const [searchParams, setSearchParams] = useSearchParams()
  useDocumentTitle("系统用户")
  const [data, setData] = useState<UserPageResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [q, setQ] = useState(searchParams.get("q") ?? "")

  const page = Number(searchParams.get("page") ?? 1)
  const size = Number(searchParams.get("size") ?? 10)

  const fetchUsers = async (pageNum: number, pageSize: number, search: string) => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(pageNum), size: String(pageSize) })
      if (search) params.set("q", search)
      const res = await api<UserPageResponse>(`/api/users?${params}`)
      setData(res)
    } catch (err) {
      showError(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers(page, size, q)
  }, [page, size])

  const handleSearch = () => {
    setSearchParams({ page: "1", size: String(size), q })
    fetchUsers(1, size, q)
  }

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  const columns = [
    { columnKey: "avatar", label: "头像" },
    { columnKey: "login", label: "用户名" },
    { columnKey: "name", label: "昵称" },
    { columnKey: "provider", label: "认证渠道" },
    { columnKey: "puid", label: "PUID" },
    { columnKey: "email", label: "邮箱" },
    { columnKey: "created_at", label: "创建时间" },
  ]

  return (
    <div className={styles.root}>
      <div className={styles.header}>
        <Text size={500} weight="semibold">系统用户</Text>
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
              <TableHeaderCell key={col.columnKey}>{col.label}</TableHeaderCell>
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
                <Text style={{ color: tokens.colorNeutralForeground3, textAlign: "center", display: "block" }}>暂无数据</Text>
              </TableCell>
            </TableRow>
          )}
          {!loading &&
            data?.records.map((user) => (
              <TableRow key={user.id}>
                <TableCell>
                  <img src={user.avatar_url} alt={user.login} className={styles.avatar} />
                </TableCell>
                <TableCell>
                  <div className={styles.loginCell}>
                    {user.provider === "github" ? (
                      <a
                        href={`https://github.com/${user.login}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        style={{ textDecoration: "none", color: "inherit" }}
                      >
                        <Text weight="semibold">{user.login}</Text>
                      </a>
                    ) : (
                      <Text weight="semibold">{user.login}</Text>
                    )}
                    <Tooltip content={user.enabled ? "此用户可正常登录本系统" : "此用户已被禁止登录本系统"} relationship="label">
                      <span className={styles.statusIcon}>
                        {user.enabled ? (
                          <CheckmarkCircleRegular style={{ color: tokens.colorStatusSuccessForeground1 }} />
                        ) : (
                          <DismissCircleRegular style={{ color: tokens.colorNeutralForeground3 }} />
                        )}
                      </span>
                    </Tooltip>
                  </div>
                </TableCell>
                <TableCell>{user.name ?? "-"}</TableCell>
                <TableCell>
                  <div className={styles.providerCell}>
                    <ProviderIcon provider={user.provider} className="size-4" />
                    <Text>{user.provider}</Text>
                  </div>
                </TableCell>
                <TableCell>
                  <Text size={200} font="monospace" truncate style={{ maxWidth: "120px" }} title={user.puid}>
                    {user.puid}
                  </Text>
                </TableCell>
                <TableCell>{user.email ?? "-"}</TableCell>
                <TableCell>
                  <Text size={200} style={{ color: tokens.colorNeutralForeground3 }}>
                    {formatTime(user.created_at)}
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
            <Text size={200} style={{ padding: `0 ${tokens.spacingHorizontalS}`, color: tokens.colorNeutralForeground3 }}>
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

export default User