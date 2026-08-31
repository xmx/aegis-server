import { useEffect, useState, useMemo, useCallback } from 'react'
import { api } from '@/lib/api'
import type { AgentRecord, AgentConnRecord, PageResponse } from '@/lib/types'
import { formatTime, formatDuration, formatBytes } from '@/lib/format'
import { InfoRegular, LinkRegular } from '@fluentui/react-icons'

interface AgentDetailProps {
  windowId: string
  props?: Record<string, unknown>
}

export default function AgentDetailApp({ props }: AgentDetailProps) {
  const agentId = props?.agentId as string | undefined
  const [agent, setAgent] = useState<AgentRecord | null>(null)
  const [agentLoading, setAgentLoading] = useState(true)
  const [tab, setTab] = useState<'info' | 'connections'>('info')
  const [connData, setConnData] = useState<PageResponse<AgentConnRecord> | null>(null)
  const [connLoading, setConnLoading] = useState(false)
  const [connPage, setConnPage] = useState(1)

  useEffect(() => {
    if (!agentId) return
    setAgentLoading(true)
    api<AgentRecord>('/api/agent?id=' + agentId)
      .then(setAgent)
      .catch(() => {})
      .finally(() => setAgentLoading(false))
  }, [agentId])

  const fetchConnections = useCallback(async (pageNum: number) => {
    if (!agentId) return
    setConnLoading(true)
    try {
      const res = await api<PageResponse<AgentConnRecord>>(
        `/api/agent/records?id=${agentId}&page=${pageNum}&size=15`,
      )
      setConnData(res)
    } catch {
      // handled by api()
    }
    setConnLoading(false)
  }, [agentId])

  useEffect(() => {
    if (tab === 'connections') fetchConnections(connPage)
  }, [tab, connPage, fetchConnections])

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

  const connTotalPages = connData ? Math.ceil(connData.total / 15) : 0

  if (agentLoading) {
    return (
      <div className="app-page">
        <div className="app-page-body">
          <div className="app-loading">加载中...</div>
        </div>
      </div>
    )
  }

  if (!agent) {
    return (
      <div className="app-page">
        <div className="app-page-body">
          <div className="app-empty">未找到该终端节点</div>
        </div>
      </div>
    )
  }

  const p = agent.process_info
  const t = agent.tunnel_info

  return (
    <div className="app-page app-page--detail">
      <div className="app-detail-nav">
        <button
          className={`app-detail-nav-item ${tab === 'info' ? 'app-detail-nav-item--active' : ''}`}
          onClick={() => setTab('info')}
        >
          <InfoRegular className="app-detail-nav-icon" />
          <span>基本信息</span>
        </button>
        <button
          className={`app-detail-nav-item ${tab === 'connections' ? 'app-detail-nav-item--active' : ''}`}
          onClick={() => setTab('connections')}
        >
          <LinkRegular className="app-detail-nav-icon" />
          <span>连接记录</span>
        </button>
      </div>

      <div className="app-page-body">
        {tab === 'info' ? (
          <div className="app-detail-sections">
            <div className="app-detail-section">
              <h3 className="app-detail-section-title">基本信息</h3>
              <div className="app-detail-grid">
                <Field label="主机名" value={p.hostname} />
                <Field label="IP 地址" value={p.inet} />
                <Field label="机器码" value={agent.machine_id} />
                <Field label="状态" value={agent.status ? '在线' : '离线'} />
                <Field label="操作系统" value={`${p.goos ?? '-'} / ${p.goarch ?? '-'}`} />
                <Field label="Agent 版本" value={p.semver?.version} />
                <Field label="进程 PID" value={p.pid?.toString()} />
                <Field label="工作目录" value={p.workdir} />
                <Field label="可执行文件" value={p.executable} />
                <Field label="创建时间" value={formatTime(agent.created_at)} />
              </div>
            </div>

            <div className="app-detail-section">
              <h3 className="app-detail-section-title">隧道信息</h3>
              <div className="app-detail-grid">
                <Field label="库名称" value={t.library_name} />
                <Field label="库模块" value={t.library_module} />
                <Field label="服务端地址" value={t.server_addr} />
                <Field label="远程地址" value={t.remote_addr} />
                <Field label="接收字节" value={t.receive_bytes?.toString()} />
                <Field label="发送字节" value={t.transmit_bytes?.toString()} />
                <Field label="最后保活" value={formatTime(t.keepalive_at)} />
              </div>
            </div>
          </div>
        ) : (
          <div className="app-detail-section">
            <div className="app-detail-section-header">
              <h3 className="app-detail-section-title">连接记录</h3>
              <span className="app-detail-section-sub">共 {connData?.total ?? 0} 条</span>
            </div>

            {connStats && (
              <div className="app-detail-summary">
                <span className="app-detail-summary-item">本页 {connStats.count}</span>
                <span className="app-detail-summary-item">在线 {connStats.online}</span>
                <span className="app-detail-summary-item">累计 {formatDuration(connStats.totalSeconds)}</span>
                <span className="app-detail-summary-item">接收 {formatBytes(connStats.rx)}</span>
                <span className="app-detail-summary-item">发送 {formatBytes(connStats.tx)}</span>
              </div>
            )}

            {connLoading ? (
              <div className="app-loading">加载中...</div>
            ) : !connData || connData.records.length === 0 ? (
              <div className="app-empty">暂无记录</div>
            ) : (
              <div className="app-table-wrap">
                <table className="app-table app-table--dense">
                  <thead>
                    <tr>
                      <th>状态/主机名</th>
                      <th>IP</th>
                      <th>版本</th>
                      <th>连接时间</th>
                      <th>断开时间</th>
                      <th>持续</th>
                      <th>服务端</th>
                      <th>远程</th>
                      <th>库</th>
                      <th>接收</th>
                      <th>发送</th>
                    </tr>
                  </thead>
                  <tbody>
                    {connData.records.map((r) => {
                      const online = !r.disconnected_at
                      return (
                        <tr key={r.id}>
                          <td>
                            <span className={`app-status-dot ${online ? 'app-status-dot--online' : 'app-status-dot--offline'}`} />
                            {r.process_info.hostname ?? '-'}
                          </td>
                          <td className="app-td-mono">{r.process_info.inet ?? '-'}</td>
                          <td className="app-td-mono">{r.process_info.semver?.version ?? '-'}</td>
                          <td className="app-td-mono">{formatTime(r.tunnel_info.connected_at)}</td>
                          <td className="app-td-mono">{online ? '—' : formatTime(r.disconnected_at)}</td>
                          <td className="app-td-num">{formatDuration(r.active_seconds)}</td>
                          <td className="app-td-mono">{r.tunnel_info.server_addr ?? '-'}</td>
                          <td className="app-td-mono">{r.tunnel_info.remote_addr ?? '-'}</td>
                          <td>{r.tunnel_info.library_name ?? '-'}</td>
                          <td className="app-td-num">{formatBytes(r.tunnel_info.receive_bytes)}</td>
                          <td className="app-td-num">{formatBytes(r.tunnel_info.transmit_bytes)}</td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {connData && connTotalPages > 1 && (
              <div className="app-pagination">
                <span className="app-pagination-info">
                  第 {(connData.page - 1) * 15 + 1}-
                  {Math.min(connData.page * 15, connData.total)} 条，共 {connData.total} 条
                </span>
                <div className="app-pagination-btns">
                  <button
                    className="app-btn app-btn-sm"
                    disabled={connPage <= 1}
                    onClick={() => setConnPage(connPage - 1)}
                  >
                    上一页
                  </button>
                  <span className="app-pagination-page">{connPage} / {connTotalPages}</span>
                  <button
                    className="app-btn app-btn-sm"
                    disabled={connPage >= connTotalPages}
                    onClick={() => setConnPage(connPage + 1)}
                  >
                    下一页
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

function Field({ label, value }: { label: string; value?: string }) {
  return (
    <div className="app-field">
      <span className="app-field-label">{label}</span>
      <span className="app-field-value">{value ?? '-'}</span>
    </div>
  )
}