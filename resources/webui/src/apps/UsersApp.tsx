import { useEffect, useState, useCallback } from 'react'
import { api } from '@/lib/api'
import type { UserRecord, PageResponse } from '@/lib/types'
import { formatTime } from '@/lib/format'
import { SearchRegular, CheckmarkCircleRegular, DismissCircleRegular } from '@fluentui/react-icons'

export default function UsersApp() {
  const [data, setData] = useState<PageResponse<UserRecord> | null>(null)
  const [loading, setLoading] = useState(false)
  const [q, setQ] = useState('')
  const [page, setPage] = useState(1)
  const size = 20

  const fetchData = useCallback(async (pageNum: number, search: string) => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(pageNum), size: String(size) })
      if (search) params.set('q', search)
      const res = await api<PageResponse<UserRecord>>('/api/users?' + params)
      setData(res)
    } catch {
      // handled
    }
    setLoading(false)
  }, [])

  useEffect(() => {
    fetchData(page, q)
  }, [page, fetchData])

  const handleSearch = () => {
    setPage(1)
    fetchData(1, q)
  }

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  return (
    <div className="app-page">
      <div className="app-page-header">
        <h2 className="app-page-title">系统用户</h2>
        <div className="app-page-search">
          <input
            className="app-search-input"
            type="text"
            placeholder="搜索..."
            value={q}
            onChange={(e) => setQ(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button className="app-btn app-btn-icon" onClick={handleSearch}>
            <SearchRegular />
          </button>
        </div>
      </div>

      <div className="app-page-body">
        {loading ? (
          <div className="app-loading">加载中...</div>
        ) : !data || data.records.length === 0 ? (
          <div className="app-empty">暂无数据</div>
        ) : (
          <>
            <div className="app-table-wrap">
              <table className="app-table">
                <thead>
                  <tr>
                    <th>头像</th>
                    <th>用户名</th>
                    <th>昵称</th>
                    <th>认证渠道</th>
                    <th>PUID</th>
                    <th>邮箱</th>
                    <th>创建时间</th>
                  </tr>
                </thead>
                <tbody>
                  {data.records.map((user) => (
                    <tr key={user.id}>
                      <td>
                        {user.avatar_url ? (
                          <img src={user.avatar_url} alt={user.login} className="app-avatar" />
                        ) : (
                          <div className="app-avatar app-avatar--placeholder" />
                        )}
                      </td>
                      <td>
                        <span className="app-user-login">
                          <span className="app-user-name">{user.login}</span>
                          {user.enabled ? (
                            <CheckmarkCircleRegular className="app-status-online" />
                          ) : (
                            <DismissCircleRegular className="app-status-offline" />
                          )}
                        </span>
                      </td>
                      <td>{user.name ?? '-'}</td>
                      <td>
                        <span className="app-tag">{user.provider}</span>
                      </td>
                      <td>
                        <span className="app-td-mono">{user.puid}</span>
                      </td>
                      <td>{user.email ?? '-'}</td>
                      <td className="app-td-muted">{formatTime(user.created_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="app-pagination">
              <span className="app-pagination-info">共 {data.total} 条</span>
              <div className="app-pagination-btns">
                <button
                  className="app-btn app-btn-sm"
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                >
                  上一页
                </button>
                <span className="app-pagination-page">{page} / {totalPages}</span>
                <button
                  className="app-btn app-btn-sm"
                  disabled={page >= totalPages}
                  onClick={() => setPage(page + 1)}
                >
                  下一页
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}