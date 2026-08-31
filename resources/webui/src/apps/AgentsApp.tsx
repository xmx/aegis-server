import { useEffect, useState, useCallback } from 'react'
import { api } from '@/lib/api'
import type { AgentRecord, PageResponse } from '@/lib/types'
import { formatTime } from '@/lib/format'
import { useWMStore } from '@/stores/wm'
import {
  SearchRegular,
  ArrowClockwiseRegular,
  ChevronLeftRegular,
  ChevronRightRegular,
} from '@fluentui/react-icons'

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

  const handleRefresh = () => {
    fetchData(page, q)
  }

  const totalPages = data ? Math.ceil(data.total / data.size) : 0

  return (
    <div className="app-page">
      {/* 命令栏 */}
      <div className="agents-commandbar">
        <div className="agents-heading">
          <h2 className="app-page-title">终端节点</h2>
          {data && <span className="agents-count">共 {data.total} 台终端</span>}
        </div>
        <div className="agents-actions">
          <div className="agents-search">
            <SearchRegular className="agents-search-icon" />
            <input
              className="agents-search-input"
              type="text"
              placeholder="搜索终端节点"
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
          <div className="agents-empty">
            <SearchRegular />
            <span>未找到终端节点</span>
          </div>
        ) : (
          <>
            <div className="agents-list">
              <div className="agents-list-header">
                <span>主机名</span>
                <span>IP 地址</span>
                <span>状态</span>
                <span>版本</span>
                <span>系统</span>
                <span>创建时间</span>
              </div>
              {data.records.map((agent) => (
                <div
                  key={agent.id}
                  className="agents-row"
                  onClick={() =>
                    openApp('agent-detail', {
                      props: { agentId: agent.id },
                      title: agent.process_info.hostname || '终端详情',
                    })
                  }
                >
                  <span className="agents-cell" title={agent.process_info.hostname ?? '-'}>
                    {agent.process_info.hostname ?? '-'}
                  </span>
                  <span className="agents-cell agents-cell--mono">
                    {agent.process_info.inet ?? '-'}
                  </span>
                  <span className="agents-cell">
                    <span className={`agents-status ${agent.status ? 'agents-status--online' : 'agents-status--offline'}`}>
                      <span className="agents-status-dot" />
                      {agent.status ? '在线' : '离线'}
                    </span>
                  </span>
                  <span className="agents-cell agents-cell--mono">
                    {agent.process_info.semver?.version ?? '-'}
                  </span>
                  <span className="agents-cell agents-cell--mono">
                    {(agent.process_info.goos ?? '-')}/{agent.process_info.goarch ?? '-'}
                  </span>
                  <span className="agents-cell agents-cell--muted">
                    {formatTime(agent.created_at)}
                  </span>
                </div>
              ))}
            </div>

            <div className="agents-pagination">
              <span className="agents-pagination-info">第 {page} 页，共 {data.total} 条</span>
              <div className="agents-pagination-btns">
                <button
                  className="agents-pagination-btn"
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                  title="上一页"
                >
                  <ChevronLeftRegular />
                </button>
                <span className="agents-pagination-page">{page} / {totalPages}</span>
                <button
                  className="agents-pagination-btn"
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