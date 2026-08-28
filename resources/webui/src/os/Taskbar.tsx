import { useState, useEffect, useCallback, memo } from 'react'
import { useWMStore } from '@/stores/wm'
import { getApp, getPinnedTaskbar } from '@/apps/registry'
import { formatClock, formatDate } from '@/lib/format'
import Logo from '@/components/Logo'
import ColorIcon from '@/components/ColorIcon'
import StartMenu from '@/os/StartMenu'
import { SearchRegular } from '@fluentui/react-icons'

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

function Taskbar() {
  const { windows, activeId, openApp, toggleMinimize, focusWindow } = useWMStore()
  const [startOpen, setStartOpen] = useState(false)
  const pinned = getPinnedTaskbar()

  const handleAppClick = useCallback(
    (appId: string) => {
      const openWins = windows.filter((w) => w.appId === appId && !w.minimized)
      if (openWins.length > 0) {
        const top = openWins.reduce((a, b) => (a.z > b.z ? a : b))
        if (top.id === activeId) {
          // 已激活 → 最小化
          toggleMinimize(top.id)
        } else {
          focusWindow(top.id)
        }
      } else {
        openApp(appId)
      }
    },
    [windows, activeId, openApp, toggleMinimize, focusWindow],
  )

  const handleStartClick = () => setStartOpen((v) => !v)
  const handleSearchClick = () => openApp('agents')

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

  return (
    <>
      {startOpen && <StartMenu onClose={() => setStartOpen(false)} />}

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
            <SearchRegular />
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
                title={app.title}
              >
                <span className="taskbar-app-icon"><ColorIcon name={app.colorIcon} size={22} /></span>
                {hasOpen && <span className={`taskbar-indicator ${isActive ? 'taskbar-indicator--active' : ''}`} />}
              </button>
            )
          })}

          {/* 运行中但未固定的窗口 */}
          {windows
            .filter((w) => !pinned.some((p) => p.id === w.appId))
            .map((win) => {
              const app = getApp(win.appId)
              return (
                <button
                  key={win.id}
                  className={`taskbar-btn ${win.id === activeId ? 'taskbar-btn--active' : ''}`}
                  onClick={() => {
                    if (win.id === activeId) toggleMinimize(win.id)
                    else focusWindow(win.id)
                  }}
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
          <Clock onClick={() => openApp('settings')} />
          <button className="taskbar-tray-btn taskbar-notify" title="通知">
            <span className="taskbar-notify-badge">0</span>
          </button>
        </div>
      </div>
    </>
  )
}

export default memo(Taskbar)