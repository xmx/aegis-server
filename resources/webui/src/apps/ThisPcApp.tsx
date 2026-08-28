import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { BuildInfo } from '@/lib/types'

/** “此电脑”：展示本机/系统基本信息 */
export default function ThisPcApp() {
  const [info, setInfo] = useState<BuildInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    api<BuildInfo>('/api/system/buildinfo')
      .then(setInfo)
      .catch(() => setError(true))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="app-page">
      <div className="app-page-header">
        <h2 className="app-page-title">此电脑</h2>
      </div>
      <div className="app-page-body">
        {loading ? (
          <div className="app-loading">加载中...</div>
        ) : error || !info ? (
          <div className="app-empty">无法获取系统信息</div>
        ) : (
          <div className="app-about-grid">
            <Field label="操作系统" value={`${info.goos} / ${info.goarch}`} />
            <Field label="版本" value={info.version} />
            <Field label="修订" value={info.revision} />
            <Field label="用户名" value={info.username} />
            <Field label="工作目录" value={info.workdir} />
            <Field label="编译时间" value={info.compiled_at} />
          </div>
        )}
      </div>
    </div>
  )
}

function Field({ label, value }: { label: string; value?: string }) {
  return (
    <div className="app-about-field">
      <span className="app-about-field-label">{label}</span>
      <span className="app-about-field-value">{value ?? '—'}</span>
    </div>
  )
}