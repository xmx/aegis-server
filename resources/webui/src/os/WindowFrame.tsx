import { type ComponentType, memo, useCallback } from 'react'
import { Rnd } from 'react-rnd'
import { DismissRegular, SubtractRegular, SquareMultipleRegular, SquareRegular } from '@fluentui/react-icons'
import { useWMStore } from '@/stores/wm'
import { getApp } from '@/apps/registry'
import ColorIcon from '@/components/ColorIcon'

interface WindowFrameProps {
  win: ReturnType<typeof useWMStore.getState>['windows'][0]
}

function WindowFrame({ win }: WindowFrameProps) {
  const { focusWindow, closeWindow, minimizeWindow, toggleMaximize, setBounds } = useWMStore()

  const meta = getApp(win.appId)
  const AppComponent: ComponentType<{ windowId: string; props?: Record<string, unknown> }> | undefined = meta?.component
  const minSize = meta?.minSize ?? { width: 360, height: 240 }

  // 最大化时 key 变化 → Rnd 销毁重建，尺寸回到最大化/还原后的值
  const rndKey = `${win.id}${win.maximized ? '-max' : ''}`

  // 非受控模式：初始值来自 win.bounds，Rnd 自行管理拖拽/缩放中的状态
  const defaultBounds = win.maximized
    ? { x: 0, y: 0, width: window.innerWidth, height: window.innerHeight - 48 }
    : { x: win.bounds.x, y: win.bounds.y, width: win.bounds.width, height: win.bounds.height }

  const handleDragStart = useCallback(() => {
    focusWindow(win.id)
  }, [focusWindow, win.id])

  const handleDragStop = useCallback(
    (_e: any, d: { x: number; y: number }) => {
      setBounds(win.id, {
        x: d.x,
        y: d.y,
        width: win.bounds.width,
        height: win.bounds.height,
      })
    },
    [setBounds, win.id, win.bounds.width, win.bounds.height],
  )

  const handleResizeStop = useCallback(
    (
      _e: MouseEvent | TouchEvent,
      _dir: any,
      ref: HTMLElement,
      _delta: any,
      position: { x: number; y: number },
    ) => {
      setBounds(win.id, {
        x: position.x,
        y: position.y,
        width: ref.offsetWidth,
        height: ref.offsetHeight,
      })
    },
    [setBounds, win.id],
  )

  return (
    <div
      data-win-id={win.id}
      className={`window-wrapper ${win.minimized ? 'window-minimized' : ''} ${win.maximized ? 'window-maximized' : ''}`}
      style={{ zIndex: win.z }}
      onPointerDown={() => focusWindow(win.id)}
    >
      <Rnd
        key={rndKey}
        default={{
          x: defaultBounds.x,
          y: defaultBounds.y,
          width: defaultBounds.width,
          height: defaultBounds.height,
        }}
        minWidth={minSize.width}
        minHeight={minSize.height}
        bounds="parent"
        dragHandleClassName="win-titlebar"
        cancel=".win-titlebar-buttons"
        disableDragging={win.maximized}
        enableResizing={!win.maximized}
        onDragStart={handleDragStart}
        onDragStop={handleDragStop}
        onResizeStop={handleResizeStop}
        style={{ position: 'absolute' }}
      >
        <div className={`window ${win.maximized ? 'window--maximized' : ''}`}>
          <div className="win-titlebar">
            <div className="win-titlebar-icon">
              {meta && <ColorIcon name={meta.colorIcon} size={16} />}
            </div>
            <span className="win-titlebar-text">{win.title}</span>
            <div className="win-titlebar-buttons">
              <button
                className="win-btn win-btn-min"
                onClick={() => minimizeWindow(win.id)}
                title="最小化"
              >
                <SubtractRegular />
              </button>
              <button
                className="win-btn win-btn-max"
                onClick={() => toggleMaximize(win.id)}
                title={win.maximized ? '还原' : '最大化'}
              >
                {win.maximized ? <SquareMultipleRegular /> : <SquareRegular />}
              </button>
              <button
                className="win-btn win-btn-close"
                onClick={() => closeWindow(win.id)}
                title="关闭"
              >
                <DismissRegular />
              </button>
            </div>
          </div>
          <div className="win-body">
            {AppComponent && <AppComponent windowId={win.id} props={win.props} />}
          </div>
        </div>
      </Rnd>
    </div>
  )
}

export default memo(WindowFrame)