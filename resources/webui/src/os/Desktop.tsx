import { useState, useCallback, useRef, useEffect, type CSSProperties } from 'react'
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
  const wallpaper = useThemeStore((s) => s.wallpaper)
  const customWallpaperUrl = useThemeStore((s) => s.customWallpaperUrl)
  const [cm, setCm] = useState<ContextMenuState | null>(null)
  const [selectedIcon, setSelectedIcon] = useState<string | null>(null)
  const desktopRef = useRef<HTMLDivElement>(null)

  const pinned = getPinnedDesktop()

  const handleContextMenu = useCallback((e: React.MouseEvent) => {
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
      setSelectedIcon(appId)
    },
    [],
  )

  useEffect(() => {
    const handler = () => setCm(null)
    window.addEventListener('click', handler)
    return () => window.removeEventListener('click', handler)
  }, [])

  // 自定义图片壁纸用 inline style，内置壁纸用 CSS class
  const isCustom = wallpaper === 'custom' && customWallpaperUrl
  const wpClass = isCustom ? 'wallpaper-custom' : `wallpaper-${wallpaper}`
  const wpStyle: CSSProperties = isCustom
    ? { backgroundImage: `url(/wallpapers/${customWallpaperUrl})` }
    : {}

  return (
    <div
      ref={desktopRef}
      className={`desktop ${wpClass}`}
      style={wpStyle}
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
            onDoubleClick={() => handleIconDoubleClick(app.id)}
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