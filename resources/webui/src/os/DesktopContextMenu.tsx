import { useLayoutEffect, useRef, useState } from 'react'
import { useWMStore } from '@/stores/wm'
import { PaintBrushRegular, InfoRegular } from '@fluentui/react-icons'

interface DesktopContextMenuProps {
  x: number
  y: number
  onClose: () => void
}

export default function DesktopContextMenu({ x, y, onClose }: DesktopContextMenuProps) {
  const { openApp } = useWMStore()
  const ref = useRef<HTMLDivElement>(null)
  const [pos, setPos] = useState({ x, y })

  // 根据菜单实际尺寸把坐标收进视口内，边缘也能完整显示
  useLayoutEffect(() => {
    const el = ref.current
    if (!el) return
    const PAD = 8
    setPos({
      x: Math.max(PAD, Math.min(x, window.innerWidth - el.offsetWidth - PAD)),
      y: Math.max(PAD, Math.min(y, window.innerHeight - el.offsetHeight - PAD)),
    })
  }, [x, y])

  const items = [
    {
      label: '个性化',
      icon: PaintBrushRegular,
      action: () => openApp('settings'),
    },
    {
      label: '关于',
      icon: InfoRegular,
      action: () => openApp('about'),
    },
  ]

  return (
    <div
      ref={ref}
      className="context-menu"
      style={{ left: pos.x, top: pos.y }}
      onClick={(e) => e.stopPropagation()}
    >
      {items.map((item) => {
        const Icon = item.icon
        return (
          <button
            key={item.label}
            className="context-menu-item"
            onClick={() => {
              item.action()
              onClose()
            }}
          >
            <span className="context-menu-item-icon"><Icon /></span>
            <span>{item.label}</span>
          </button>
        )
      })}
    </div>
  )
}