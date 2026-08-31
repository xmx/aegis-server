import { useEffect, useState, useCallback } from 'react'
import { api } from '@/lib/api'
import type { UserRecord, PageResponse } from '@/lib/types'
import { formatTime } from '@/lib/format'
import {
  SearchRegular,
  ArrowClockwiseRegular,
  ChevronLeftRegular,
  ChevronRightRegular,
  CheckmarkCircleRegular,
  DismissCircleRegular,
} from '@fluentui/react-icons'

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
      // handled by api()
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchData(page, q)
  }, [page, fetchData])

  const handleSearch = () => {
    setPage(1)
    fetchData(1, q)
  }

  const handleRefresh = () => {
    fetchData(page, q)
  }

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  return (
    <div className="app-page">
      {/* 命令栏 */}
      <div className="users-commandbar">
        <div className="users-heading">
          <h2 className="app-page-title">系统用户</h2>
          {data && <span className="users-count">共 {data.total} 个用户</span>}
        </div>
        <div className="users-actions">
          <div className="users-search">
            <SearchRegular className="users-search-icon" />
            <input
              className="users-search-input"
              type="text"
              placeholder="搜索用户"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
            />
          </div>
          <button
            className="app-btn app-btn-icon"
            onClick={handleRefresh}
            title="刷新"
            disabled={loading}
          >
            <ArrowClockwiseRegular />
          </button>
        </div>
      </div>

      <div className="app-page-body">
        {loading ? (
          <div className="app-loading">正在加载...</div>
        ) : !data || data.records.length === 0 ? (
          <div className="users-empty">
            <SearchRegular />
            <span>未找到用户</span>
          </div>
        ) : (
          <>
            <div className="users-list">
              <div className="users-list-header">
                <span>头像</span>
                <span>用户名</span>
                <span>昵称</span>
                <span>认证渠道</span>
                <span>PUID</span>
                <span>邮箱</span>
                <span>创建时间</span>
              </div>
              {data.records.map((user) => (
                <div key={user.id} className="users-row">
                  <span className="users-cell users-cell--avatar">
                    {user.avatar_url ? (
                      <img className="users-avatar" src={user.avatar_url} alt={user.login} />
                    ) : (
                      <span className="users-avatar users-avatar--placeholder">
                        {user.login?.charAt(0).toUpperCase() ?? '?'}
                      </span>
                    )}
                  </span>
                  <span className="users-cell">
                    <span className="users-name">
                      {user.login}
                      {user.enabled ? (
                        <CheckmarkCircleRegular className="users-status users-status--on" title="已启用" />
                      ) : (
                        <DismissCircleRegular className="users-status users-status--off" title="已禁用" />
                      )}
                    </span>
                  </span>
                  <span className="users-cell">{user.name ?? '-'}</span>
                  <span className="users-cell">
                    <span className="users-provider">{user.provider}</span>
                  </span>
                  <span className="users-cell users-cell--mono">{user.puid}</span>
                  <span className="users-cell">{user.email ?? '-'}</span>
                  <span className="users-cell users-cell--muted">{formatTime(user.created_at)}</span>
                </div>
              ))}
            </div>

            <div className="users-pagination">
              <span className="users-pagination-info">第 {page} 页，共 {data.total} 条</span>
              <div className="users-pagination-btns">
                <button
                  className="users-pagination-btn"
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                  title="上一页"
                >
                  <ChevronLeftRegular />
                </button>
                <span className="users-pagination-page">{page} / {totalPages}</span>
                <button
                  className="users-pagination-btn"
                  disabled={page >= totalPages}
                  onClick={() => setPage(page + 1)}
                  title="下一页"
                >
                  <ChevronRightRegular />
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}