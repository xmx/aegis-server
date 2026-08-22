import { useEffect, useState } from "react"
import { useSearchParams } from "react-router-dom"
import { api } from "@/lib/api"
import { showError } from "@/lib/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Card,
  CardContent,
} from "@/components/ui/card"
import { SearchIcon, CircleCheckIcon, CircleXIcon } from "lucide-react"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import { ProviderIcon } from "@/components/Icons"

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

function User() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [data, setData] = useState<UserPageResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [q, setQ] = useState(searchParams.get("q") ?? "")

  const page = Number(searchParams.get("page") ?? 1)
  const size = Number(searchParams.get("size") ?? 12)

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

  return (
    <div className="flex flex-1 flex-col gap-4 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">系统用户</h1>
        <div className="flex items-center gap-2">
          <Input
            placeholder="搜索..."
            value={q}
            onChange={(e) => setQ(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSearch()}
            className="w-48"
          />
          <Button variant="outline" size="icon" onClick={handleSearch}>
            <SearchIcon className="size-4" />
          </Button>
        </div>
      </div>

      {loading && (
        <p className="py-8 text-center text-muted-foreground">加载中...</p>
      )}
      {!loading && data && data.records.length === 0 && (
        <p className="py-8 text-center text-muted-foreground">暂无数据</p>
      )}
      {!loading && data && data.records.length > 0 && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {data.records.map((user) => (
            <Card key={user.id} className="transition-shadow hover:shadow-md">
              <CardContent className="flex flex-col items-center gap-3 p-5">
                <img
                  src={user.avatar_url}
                  alt={user.login}
                  className="size-16 rounded-full object-cover"
                />
                <div className="flex flex-col items-center gap-0.5">
                  <span className="font-medium">{user.name ?? user.login}</span>
                  <div className="flex items-center gap-1">
                    <span className="text-sm text-muted-foreground">{user.login}</span>
                    <Tooltip>
                      <TooltipTrigger>
                        {user.enabled ? (
                          <CircleCheckIcon className="size-4 text-emerald-500" />
                        ) : (
                          <CircleXIcon className="size-4 text-muted-foreground" />
                        )}
                      </TooltipTrigger>
                      <TooltipContent>
                        {user.enabled ? "此用户可正常登录本系统" : "此用户已被禁止登录本系统"}
                      </TooltipContent>
                    </Tooltip>
                  </div>
                </div>
                <div className="flex w-full items-center justify-between pt-2 text-xs text-muted-foreground">
                  <ProviderIcon provider={user.provider} className="size-4" />
                  <span className="truncate font-mono" title={user.puid}>{user.puid}</span>
                </div>
                <div className="flex w-full justify-center pt-1 text-xs text-muted-foreground">
                  <span>{formatTime(user.created_at)}</span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {data && totalPages > 1 && (
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">
            共 {data.total} 条
          </span>
          <div className="flex items-center gap-1">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() =>
                setSearchParams({ page: String(page - 1), size: String(size), q })
              }
            >
              上一页
            </Button>
            <span className="px-2 text-sm text-muted-foreground">
              {page} / {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= totalPages}
              onClick={() =>
                setSearchParams({ page: String(page + 1), size: String(size), q })
              }
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