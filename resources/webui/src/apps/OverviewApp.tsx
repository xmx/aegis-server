import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { AgentRecord, PageResponse } from '@/lib/types'
import { useWMStore } from '@/stores/wm'
import ColorIcon from '@/components/ColorIcon'

interface Stats {
  total: number
  online: number
  offline: number
}

export default function OverviewApp() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)
  const { openApp } = useWMStore()

  useEffect(() => {
    api<PageResponse<AgentRecord>>('/api/agents?page=1&size=100')
      .then((data) => {
        setStats({
          total: data.total,
          online: data.records.filter((a) => a.status).length,
          offline: data.records.filter((a) => !a.status).length,
        })
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="app-page">
      <div className="app-page-header">
        <h2 className="app-page-title">概览</h2>
      </div>
      <div className="app-page-body">
        {loading ? (
          <div className="app-loading">加载中...</div>
        ) : stats ? (
          <div className="overview-grid">
            <div className="overview-card" onClick={() => openApp('agents')}>
              <div className="overview-card-value">{stats.total}</div>
              <div className="overview-card-label">终端总数</div>
            </div>
            <div className="overview-card overview-card--online" onClick={() => openApp('agents')}>
              <div className="overview-card-value">{stats.online}</div>
              <div className="overview-card-label">在线</div>
            </div>
            <div className="overview-card overview-card--offline">
              <div className="overview-card-value">{stats.offline}</div>
              <div className="overview-card-label">离线</div>
            </div>
          </div>
        ) : (
          <div className="app-empty">无法获取数据</div>
        )}

        <div className="overview-links">
          <h3 className="overview-links-title">快捷操作</h3>
          <div className="overview-links-grid">
            <button className="overview-link-card" onClick={() => openApp('agents')}>
              <span className="overview-link-icon"><ColorIcon name="agents_48_color" size={28} /></span>
              <span>终端节点管理</span>
            </button>
            <button className="overview-link-card" onClick={() => openApp('users')}>
              <span className="overview-link-icon"><ColorIcon name="people_48_color" size={28} /></span>
              <span>系统用户管理</span>
            </button>
            <button className="overview-link-card" onClick={() => openApp('settings')}>
              <span className="overview-link-icon"><ColorIcon name="settings_48_color" size={28} /></span>
              <span>系统设置</span>
            </button>
            <button className="overview-link-card" onClick={() => openApp('about')}>
              <span className="overview-link-icon"><ColorIcon name="question_circle_48_color" size={28} /></span>
              <span>关于</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}