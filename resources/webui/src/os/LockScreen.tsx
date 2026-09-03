import { useState, useEffect } from 'react'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { api, ApiRequestError } from '@/lib/api'
import type { OAuthProviderInfo } from '@/lib/types'
import { formatClock, formatLockDate } from '@/lib/format'

export default function LockScreen({ onUnlock }: { onUnlock: () => void }) {
  const { user } = useAuthStore()
  const lockWallpaper = useThemeStore((s) => s.lockWallpaper)
  const [now, setNow] = useState(new Date())
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000)
    return () => clearInterval(id)
  }, [])

  const handleGitHubLogin = async () => {
    setLoading(true)
    setError('')
    try {
      const data = await api<OAuthProviderInfo>(
        `/api/oauth/provider?provider=github&origin=${encodeURIComponent(window.origin)}`,
      )
      const authParams = new URLSearchParams({
        client_id: data.client_id,
        redirect_uri: data.redirect_uri,
        scope: (data.scopes ?? []).join(' '),
      })
      window.location.href = `${data.auth_url}?${authParams.toString()}`
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.detail || err.message)
      } else {
        setError('无法连接认证服务')
      }
      setLoading(false)
    }
  }

  return (
    <div className="lockscreen">
      <div
        className="lockscreen-wallpaper"
        style={{ backgroundImage: `url(/lockscreen/${lockWallpaper})` }}
      />

      <div className="lockscreen-top">
        <div className="lockscreen-top-time">{formatClock(now)}</div>
        <div className="lockscreen-top-date">{formatLockDate(now)}</div>
      </div>

      <div className="lockscreen-panel">
        {user ? (
          <div className="lockscreen-user">
            <img src={user.avatar_url} alt={user.login} className="lockscreen-avatar" />
            <div className="lockscreen-user-name">{user.name ?? user.login}</div>
            <button className="lockscreen-login-btn" onClick={onUnlock}>
              登录
            </button>
          </div>
        ) : (
          <div className="lockscreen-oauth">
            <button
              className="lockscreen-login-btn"
              onClick={handleGitHubLogin}
              disabled={loading}
            >
              {loading ? (
                '正在连接...'
              ) : (
                <>
                  <svg className="lockscreen-btn-icon" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
                    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
                  </svg>
                  GitHub 登录
                </>
              )}
            </button>
            {error && <div className="lockscreen-error">{error}</div>}
          </div>
        )}
      </div>
    </div>
  )
}