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
      const data = await api<OAuthProviderInfo>('/api/oauth/provider?provider=github')
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
              {loading ? '正在连接...' : '使用 GitHub 登录'}
            </button>
            {error && <div className="lockscreen-error">{error}</div>}
          </div>
        )}
      </div>
    </div>
  )
}