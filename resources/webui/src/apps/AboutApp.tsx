import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import Logo from '@/components/Logo'
import type { BuildInfo } from '@/lib/types'

export default function AboutApp() {
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
        <h2 className="app-page-title">关于</h2>
      </div>
      <div className="app-page-body">
        {loading ? (
          <div className="app-loading">加载中...</div>
        ) : error ? (
          <div className="app-empty">无法获取构建信息</div>
        ) : info ? (
          <div className="app-about">
            <div className="app-about-hero">
              <Logo size={48} />
              <div className="app-about-version">v{info.version}</div>
            </div>

            <div className="app-about-section">
              <h3 className="app-about-section-title">基本信息</h3>
              <div className="app-about-grid">
                <AboutField label="版本" value={info.version} />
                <AboutField label="修订" value={info.revision} />
                <AboutField label="用户名" value={info.username} />
                <AboutField label="工作目录" value={info.workdir} />
                <AboutField label="模块" value={info.module} />
                <AboutField label="提交时间" value={info.committed_at} />
                <AboutField label="编译时间" value={info.compiled_at} />
                <AboutField label="系统" value={`${info.goos}/${info.goarch}`} />
              </div>
            </div>

            {info.build_info && (
              <>
                <div className="app-about-section">
                  <h3 className="app-about-section-title">构建信息</h3>
                  <div className="app-about-grid">
                    <AboutField label="Go 版本" value={info.build_info.GoVersion} />
                    <AboutField label="主模块" value={`${info.build_info.Main.Path} ${info.build_info.Main.Version}`} />
                  </div>
                </div>

                <div className="app-about-section">
                  <h3 className="app-about-section-title">依赖项 ({info.build_info.Deps.length})</h3>
                  <div className="app-about-deps">
                    {info.build_info.Deps.slice(0, 30).map((d) => (
                      <div key={d.Path} className="app-about-dep">
                        <span className="app-about-dep-path">{d.Path}</span>
                        <span className="app-about-dep-ver">{d.Version}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            )}
          </div>
        ) : null}
      </div>
    </div>
  )
}

function AboutField({ label, value }: { label: string; value?: string }) {
  return (
    <div className="app-about-field">
      <span className="app-about-field-label">{label}</span>
      <span className="app-about-field-value">{value ?? '-'}</span>
    </div>
  )
}