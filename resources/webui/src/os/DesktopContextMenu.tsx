import { useWMStore } from '@/stores/wm'
import ColorIcon from '@/components/ColorIcon'

interface DesktopContextMenuProps {
  x: number
  y: number
  onClose: () => void
}

export default function DesktopContextMenu({ x, y, onClose }: DesktopContextMenuProps) {
  const { openApp } = useWMStore()

  const items: { label: string; colorIcon: string; action: () => void }[] = [
    {
      label: '个性化',
      colorIcon: 'paint_brush_32_color',
      action: () => openApp('settings'),
    },
    {
      label: '关于',
      colorIcon: 'question_circle_48_color',
      action: () => openApp('about'),
    },
  ]

  return (
    <div
      className="context-menu"
      style={{ left: x, top: y }}
      onClick={(e) => e.stopPropagation()}
    >
      {items.map((item) => (
        <button
          key={item.label}
          className="context-menu-item"
          onClick={() => {
            item.action()
            onClose()
          }}
        >
          <span className="context-menu-item-icon">
            <ColorIcon name={item.colorIcon} size={16} />
          </span>
          <span>{item.label}</span>
        </button>
      ))}
    </div>
  )
}