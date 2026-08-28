import { useEffect, useState, useCallback } from 'react'
import { api } from '@/lib/api'
import type { AgentRecord, PageResponse } from '@/lib/types'
import { formatTime } from '@/lib/format'
import { useWMStore } from '@/stores/wm'
import { CircleFilled, CircleRegular, SearchRegular } from '@fluentui/react-icons'

export default function AgentsApp() {
  const [data, setData] = useState<PageResponse<AgentRecord> | null>(null)
  const [loading, setLoading] = useState(false)
  const [q, setQ] = useState('')
  const [page, setPage] = useState(1)
  const size = 20
  const { openApp } = useWMStore()

  const fetchData = useCallback(async (pageNum: number, search: string) => {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(pageNum), size: String(size) })
      if (search) params.set('q', search)
      const res = await api<PageResponse<AgentRecord>>('/api/agents?' + params)
      setData(res)
    } catch {
      // error handled by api()
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

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  return (
    <div className="app-page">
      <div className="app-page-header">
        <h2 className="app-page-title">终端节点</h2>
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
                    <th>主机名</th>
                    <th>IP 地址</th>
                    <th>状态</th>
                    <th>版本</th>
                    <th>系统</th>
                    <th>创建时间</th>
                  </tr>
                </thead>
                <tbody>
                  {data.records.map((agent) => (
                    <tr
                      key={agent.id}
                      className="app-table-row--clickable"
                      onClick={() => openApp('agent-detail', {
                        props: { agentId: agent.id },
                        title: agent.process_info.hostname || '终端详情',
                      })}
                    >
                      <td>{agent.process_info.hostname ?? '-'}</td>
                      <td>{agent.process_info.inet ?? '-'}</td>
                      <td>
                        <span className="app-status">
                          {agent.status ? (
                            <CircleFilled className="app-status-online" />
                          ) : (
                            <CircleRegular className="app-status-offline" />
                          )}
                          {agent.status ? '在线' : '离线'}
                        </span>
                      </td>
                      <td>
                        <span className="app-tag">{agent.process_info.semver?.version ?? '-'}</span>
                      </td>
                      <td>
                        <span className="app-tag">
                          {(agent.process_info.goos ?? '-')}/{agent.process_info.goarch ?? '-'}
                        </span>
                      </td>
                      <td className="app-td-muted">{formatTime(agent.created_at)}</td>
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