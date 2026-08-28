import { useState, useEffect, useCallback, type CSSProperties } from 'react'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { api, ApiRequestError } from '@/lib/api'
import type { OAuthProviderInfo } from '@/lib/types'
import { formatClock, formatDate } from '@/lib/format'

export default function LockScreen({ onUnlock }: { onUnlock: () => void }) {
  const { user } = useAuthStore()
  const wallpaper = useThemeStore((s) => s.wallpaper)
  const customWallpaperUrl = useThemeStore((s) => s.customWallpaperUrl)
  const [showLogin, setShowLogin] = useState(false)
  const [now, setNow] = useState(new Date())
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const isCustom = wallpaper === 'custom' && customWallpaperUrl
  const wpStyle: CSSProperties = isCustom
    ? { backgroundImage: `url(/wallpapers/${customWallpaperUrl})`, backgroundSize: 'cover', backgroundPosition: 'center' }
    : {}

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000)
    return () => clearInterval(id)
  }, [])

  const handleUnlock = useCallback(() => {
    if (user) {
      onUnlock()
    } else {
      setShowLogin(true)
    }
  }, [user, onUnlock])

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        handleUnlock()
      }
    },
    [handleUnlock],
  )

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [handleKeyDown])

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
    <div className="lockscreen" onClick={!showLogin ? handleUnlock : undefined}>
      <div
        className={`lockscreen-wallpaper ${isCustom ? 'wallpaper-custom' : `wallpaper-${wallpaper}`}`}
        style={wpStyle}
      />

      {!showLogin ? (
        <div className="lockscreen-hero">
          <div className="lockscreen-time">{formatClock(now)}</div>
          <div className="lockscreen-date">{formatDate(now)}</div>
          <div className="lockscreen-hint">点击或按任意键登录</div>
        </div>
      ) : (
        <div className="lockscreen-panel" onClick={(e) => e.stopPropagation()}>
          <div className="lockscreen-panel-time">{formatClock(now)}</div>
          <div className="lockscreen-panel-date">{formatDate(now)}</div>

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
      )}
    </div>
  )
}