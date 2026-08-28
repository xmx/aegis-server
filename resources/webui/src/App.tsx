import { useState, useEffect, useCallback, useRef } from 'react'
import { useThemeStore, accentContrast, accentDim } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import type { UserRecord } from '@/lib/types'
import Logo from '@/components/Logo'
import Desktop from '@/os/Desktop'
import Taskbar from '@/os/Taskbar'
import LockScreen from '@/os/LockScreen'

function GitHubCallback() {
  const [status, setStatus] = useState<'loading' | 'error'>('loading')
  const { setUser } = useAuthStore()
  const called = useRef(false)

  useEffect(() => {
    if (called.current) return
    called.current = true

    const params = new URLSearchParams(window.location.search)
    const code = params.get('code')
    if (!code) {
      setStatus('error')
      return
    }

    api<UserRecord>('/api/oauth/github', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code,
        redirect_uri: window.location.origin + window.location.pathname,
      }),
    })
      .then((user) => {
        setUser({
          login: user.login,
          name: user.name,
          avatar_url: user.avatar_url ?? '',
          provider: user.provider,
          email: user.email,
        })
        window.history.replaceState({}, '', '/')
      })
      .catch(() => {
        setStatus('error')
      })
  }, [setUser])

  return (
    <div className="lockscreen">
      <div className="lockscreen-wallpaper wallpaper-bloom" />
      <div className="lockscreen-panel">
        <div className="lockscreen-panel-time">
          {status === 'loading' ? '正在登录...' : '登录失败'}
        </div>
        {status === 'error' && (
          <button className="lockscreen-login-btn" onClick={() => (window.location.href = '/')}>
            返回
          </button>
        )}
      </div>
    </div>
  )
}

export default function App() {
  const { mode, wallpaper, windowTransparency, accent } = useThemeStore()
  const { user } = useAuthStore()
  const [phase, setPhase] = useState<'boot' | 'lock' | 'desktop'>('boot')
  const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>('dark')

  const isCallback = window.location.pathname === '/login/github'

  useEffect(() => {
    document.documentElement.style.setProperty('--window-opacity', String(windowTransparency))
  }, [windowTransparency])

  useEffect(() => {
    const root = document.documentElement
    root.style.setProperty('--accent', accent)
    root.style.setProperty('--accent-dim', accentDim(accent))
    root.style.setProperty('--accent-contrast', accentContrast(accent))
  }, [accent])

  useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const update = () => {
      if (mode === 'system') {
        setResolvedTheme(mq.matches ? 'dark' : 'light')
      }
    }
    update()
    mq.addEventListener('change', update)
    return () => mq.removeEventListener('change', update)
  }, [mode])

  useEffect(() => {
    if (mode !== 'system') {
      setResolvedTheme(mode)
    }
  }, [mode])

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', resolvedTheme)
  }, [resolvedTheme])

  useEffect(() => {
    const id = setTimeout(() => {
      setPhase(user ? 'desktop' : 'lock')
    }, 1800)
    return () => clearTimeout(id)
  }, [user])

  const handleUnlock = useCallback(() => {
    setPhase('desktop')
  }, [])

  useEffect(() => {
    if (!user && phase === 'desktop') {
      setPhase('lock')
    }
  }, [user, phase])

  if (isCallback) {
    return <GitHubCallback />
  }

  if (phase === 'boot') {
    return (
      <div className="boot-screen">
        <Logo size={64} className="boot-logo" />
        <div className="boot-loading">
          <div className="boot-spinner" />
        </div>
      </div>
    )
  }

  if (phase === 'lock') {
    return <LockScreen onUnlock={handleUnlock} />
  }

  return (
    <div className={`app-shell wallpaper-${wallpaper}`}>
      <Desktop />
      <Taskbar />
    </div>
  )
}