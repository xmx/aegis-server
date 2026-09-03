import { useState, useEffect, useCallback, useLayoutEffect, useRef, useMemo, memo } from 'react'
import { useWMStore } from '@/stores/wm'
import { useTaskbarStore } from '@/stores/taskbar'
import { getApp, getPinnedTaskbar, type AppMeta } from '@/apps/registry'
import { formatClock, formatDate } from '@/lib/format'
import Logo from '@/components/Logo'
import ColorIcon from '@/components/ColorIcon'
import StartMenu from '@/os/StartMenu'
import Calendar from '@/os/Calendar'

function Clock({ onClick }: { onClick: () => void }) {
  const [now, setNow] = useState(new Date())
  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000)
    return () => clearInterval(id)
  }, [])
  return (
    <button className="tray-clock" onClick={onClick} title="时间和通知">
      <span className="tray-clock-time">{formatClock(now)}</span>
      <span className="tray-clock-date">{formatDate(now)}</span>
    </button>
  )
}

interface TaskbarMenuState {
  x: number
  y: number
  appId: string
}

function Taskbar() {
  const { windows, activeId, openApp, toggleMinimize, focusWindow, toggleMaximize, closeWindow } = useWMStore()
  const { unpinned, extraPinned, pin, unpin } = useTaskbarStore()
  const [startOpen, setStartOpen] = useState(false)
  const [calendarOpen, setCalendarOpen] = useState(false)
  const [menu, setMenu] = useState<TaskbarMenuState | null>(null)
  const [menuPos, setMenuPos] = useState<{ x: number; y: number } | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  const staticPinnedIds = useMemo(() => getPinnedTaskbar().map((a) => a.id), [])
  const pinnedIds = useMemo(() => {
    const set = new Set<string>()
    staticPinnedIds.forEach((id) => {
      if (!unpinned.includes(id)) set.add(id)
    })
    extraPinned.forEach((id) => set.add(id))
    return Array.from(set)
  }, [staticPinnedIds, unpinned, extraPinned])

  const pinned: AppMeta[] = pinnedIds.map((id) => getApp(id)).filter((a): a is AppMeta => Boolean(a))

  const handleAppClick = useCallback(
    (appId: string) => {
      const openWins = windows.filter((w) => w.appId === appId && !w.minimized)
      if (openWins.length > 0) {
        const top = openWins.reduce((a, b) => (a.z > b.z ? a : b))
        if (top.id === activeId) toggleMinimize(top.id)
        else focusWindow(top.id)
      } else {
        openApp(appId)
      }
    },
    [windows, activeId, openApp, toggleMinimize, focusWindow],
  )

  const handleStartClick = () => setStartOpen((v) => !v)
  const handleSearchClick = () => openApp('agents')

  const openMenu = (e: React.MouseEvent, appId: string) => {
    e.preventDefault()
    e.stopPropagation()
    setMenu({ x: e.clientX, y: e.clientY, appId })
  }

  const runMenuAction = useCallback(
    (action: string) => {
      if (!menu) return
      const appId = menu.appId
      const appWins = windows.filter((w) => w.appId === appId)
      const visible = appWins.filter((w) => !w.minimized)
      const pool = visible.length ? visible : appWins
      const top = pool.length ? pool.reduce((a, b) => (a.z > b.z ? a : b)) : null
      setMenu(null)
      if (action === 'open') openApp(appId)
      else if (action === 'toggle' && top) toggleMinimize(top.id)
      else if (action === 'maximize' && top) toggleMaximize(top.id)
      else if (action === 'close' && top) closeWindow(top.id)
      else if (action === 'pin') pin(appId)
      else if (action === 'unpin') unpin(appId)
    },
    [menu, windows, openApp, toggleMinimize, toggleMaximize, closeWindow, pin, unpin],
  )

  // 菜单位置钳制 + 点击外部/ Esc 关闭
  useLayoutEffect(() => {
    if (!menu) {
      setMenuPos(null)
      return
    }
    const el = menuRef.current
    if (!el) return
    const PAD = 8
    setMenuPos({
      x: Math.max(PAD, Math.min(menu.x, window.innerWidth - el.offsetWidth - PAD)),
      y: Math.max(PAD, Math.min(menu.y, window.innerHeight - el.offsetHeight - PAD)),
    })
  }, [menu])

  useEffect(() => {
    if (!menu) return
    const onDown = (e: MouseEvent) => {
      const t = e.target as HTMLElement
      if (!t.closest('.taskbar-context-menu')) setMenu(null)
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setMenu(null)
    }
    window.addEventListener('mousedown', onDown)
    window.addEventListener('keydown', onKey)
    return () => {
      window.removeEventListener('mousedown', onDown)
      window.removeEventListener('keydown', onKey)
    }
  }, [menu])

  // 点击开始菜单外部区域时关闭菜单
  useEffect(() => {
    if (!startOpen) return
    const handleMouseDown = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      if (!target.closest('.start-menu') && !target.closest('.taskbar-btn-start')) {
        setStartOpen(false)
      }
    }
    window.addEventListener('mousedown', handleMouseDown)
    return () => window.removeEventListener('mousedown', handleMouseDown)
  }, [startOpen])

  // 右键菜单的上下文数据
  const menuWins = menu ? windows.filter((w) => w.appId === menu.appId) : []
  const menuTop = (() => {
    const visible = menuWins.filter((w) => !w.minimized)
    const pool = visible.length ? visible : menuWins
    return pool.length ? pool.reduce((a, b) => (a.z > b.z ? a : b)) : null
  })()
  const menuIsPinned = menu ? pinnedIds.includes(menu.appId) : false

  return (
    <>
      {startOpen && <StartMenu onClose={() => setStartOpen(false)} />}

      {calendarOpen && <Calendar onClose={() => setCalendarOpen(false)} />}

      {menu && (
        <div
          ref={menuRef}
          className="context-menu taskbar-context-menu"
          style={{ left: (menuPos ?? menu).x, top: (menuPos ?? menu).y }}
        >
          {menuWins.length === 0 ? (
            <button className="context-menu-item" onClick={() => runMenuAction('open')}>
              <span>打开</span>
            </button>
          ) : (
            <>
              <button className="context-menu-item" onClick={() => runMenuAction('toggle')}>
                <span>{menuTop?.minimized ? '还原' : '最小化'}</span>
              </button>
              <button className="context-menu-item" onClick={() => runMenuAction('maximize')}>
                <span>{menuTop?.maximized ? '还原' : '最大化'}</span>
              </button>
              <button className="context-menu-item" onClick={() => runMenuAction('close')}>
                <span>关闭窗口</span>
              </button>
            </>
          )}
          <div className="context-menu-sep" />
          <button className="context-menu-item" onClick={() => runMenuAction(menuIsPinned ? 'unpin' : 'pin')}>
            <span>{menuIsPinned ? '从任务栏取消固定' : '固定到任务栏'}</span>
          </button>
        </div>
      )}

      <div className="taskbar">
        {/* 中间：应用图标（Win11 居中） */}
        <div className="taskbar-center">
          <button
            className={`taskbar-btn taskbar-btn-start ${startOpen ? 'taskbar-btn--active' : ''}`}
            onClick={handleStartClick}
            title="开始"
          >
            <Logo size={18} className="taskbar-start-icon" />
          </button>
          <button className="taskbar-btn" onClick={handleSearchClick} title="搜索">
            <ColorIcon name="search_sparkle_24_color" size={18} />
          </button>

          {pinned.map((app) => {
            const openWins = windows.filter((w) => w.appId === app.id)
            const hasOpen = openWins.length > 0
            const isActive = hasOpen && openWins.some((w) => w.id === activeId)
            return (
              <button
                key={app.id}
                className={`taskbar-btn ${isActive ? 'taskbar-btn--active' : ''}`}
                onClick={() => handleAppClick(app.id)}
                onContextMenu={(e) => openMenu(e, app.id)}
                title={app.title}
              >
                <span className="taskbar-app-icon"><ColorIcon name={app.colorIcon} size={22} /></span>
                {hasOpen && <span className={`taskbar-indicator ${isActive ? 'taskbar-indicator--active' : ''}`} />}
              </button>
            )
          })}

          {/* 运行中但未固定的窗口 */}
          {windows
            .filter((w) => !pinnedIds.includes(w.appId))
            .map((win) => {
              const app = getApp(win.appId)
              return (
                <button
                  key={win.id}
                  className={`taskbar-btn ${win.id === activeId ? 'taskbar-btn--active' : ''}`}
                  onClick={() => {
                    if (win.minimized) toggleMinimize(win.id)
                    else if (win.id === activeId) toggleMinimize(win.id)
                    else focusWindow(win.id)
                  }}
                  onContextMenu={(e) => openMenu(e, win.appId)}
                  title={win.title}
                >
                  <span className="taskbar-app-icon">{app && <ColorIcon name={app.colorIcon} size={22} />}</span>
                  <span className={`taskbar-indicator ${win.id === activeId ? 'taskbar-indicator--active' : ''}`} />
                </button>
              )
            })}
        </div>

        {/* 右侧：托盘 */}
        <div className="taskbar-right">
          <Clock onClick={() => setCalendarOpen((v) => !v)} />
          <button className="taskbar-tray-btn taskbar-notify" title="通知">
            <span className="taskbar-notify-badge">0</span>
          </button>
        </div>
      </div>
    </>
  )
}

export default memo(Taskbar)