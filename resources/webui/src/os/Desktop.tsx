import { useState, useCallback, useRef, useEffect } from 'react'
import { useWMStore } from '@/stores/wm'
import { getPinnedDesktop } from '@/apps/registry'
import WindowFrame from '@/os/WindowFrame'
import DesktopContextMenu from '@/os/DesktopContextMenu'
import { useThemeStore } from '@/stores/theme'
import ColorIcon from '@/components/ColorIcon'

interface ContextMenuState {
  x: number
  y: number
}

export default function Desktop() {
  const { windows, openApp } = useWMStore()
  const singleClickOpen = useThemeStore((s) => s.singleClickOpen)
  const [cm, setCm] = useState<ContextMenuState | null>(null)
  const [selectedIcon, setSelectedIcon] = useState<string | null>(null)
  const desktopRef = useRef<HTMLDivElement>(null)

  const pinned = getPinnedDesktop()

  const handleContextMenu = useCallback((e: React.MouseEvent) => {
    const target = e.target as HTMLElement
    // 窗口内部右键交给应用自己处理，桌面不再弹出菜单
    if (target.closest('.window-layer')) return
    e.preventDefault()
    setCm({ x: e.clientX, y: e.clientY })
    setSelectedIcon(null)
  }, [])

  const handleClick = useCallback(() => {
    setCm(null)
    setSelectedIcon(null)
  }, [])

  const handleIconDoubleClick = useCallback(
    (appId: string) => {
      openApp(appId)
    },
    [openApp],
  )

  const handleIconClick = useCallback(
    (e: React.MouseEvent, appId: string) => {
      e.stopPropagation()
      if (singleClickOpen) {
        openApp(appId)
      } else {
        setSelectedIcon(appId)
      }
    },
    [singleClickOpen, openApp],
  )

  useEffect(() => {
    const handler = () => setCm(null)
    window.addEventListener('click', handler)
    return () => window.removeEventListener('click', handler)
  }, [])

  return (
    <div
      ref={desktopRef}
      className="desktop"
      onContextMenu={handleContextMenu}
      onClick={handleClick}
    >
      {/* 桌面图标 */}
      <div className="desktop-icons">
        {pinned.map((app) => (
          <div
            key={app.id}
            className={`desktop-icon ${selectedIcon === app.id ? 'desktop-icon--selected' : ''}`}
            onClick={(e) => handleIconClick(e, app.id)}
            onMouseEnter={singleClickOpen ? () => setSelectedIcon(app.id) : undefined}
            onDoubleClick={singleClickOpen ? undefined : () => handleIconDoubleClick(app.id)}
          >
            <div className="desktop-icon-img"><ColorIcon name={app.colorIcon} size={36} /></div>
            <span className="desktop-icon-label">{app.title}</span>
          </div>
        ))}
      </div>

      {/* 窗口层 */}
      <div className="window-layer">
        {windows.map((win) => (
          <WindowFrame key={win.id} win={win} />
        ))}
      </div>

      {/* 右键菜单 */}
      {cm && (
        <DesktopContextMenu
          x={cm.x}
          y={cm.y}
          onClose={() => setCm(null)}
        />
      )}
    </div>
  )
}