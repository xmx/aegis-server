import { useEffect, useRef } from 'react'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'

interface TerminalAppProps {
  windowId: string
  props?: Record<string, unknown>
}

/** 把特殊键翻译为 VT 序列；普通可打印字符仍交由 xterm 的 onData 处理。 */
function encodeKey(e: KeyboardEvent): string | null {
  const key = e.key
  const ctrl = e.ctrlKey && !e.altKey && !e.metaKey
  const alt = e.altKey && !e.ctrlKey && !e.metaKey

  if (ctrl) {
    // Ctrl+@/A-Z/[\]^_ → 控制字符
    if (key.length === 1) {
      const code = key.toUpperCase().charCodeAt(0)
      if (code >= 0x40 && code <= 0x5f) return String.fromCharCode(code - 0x40)
      if (key === ' ') return '\x00'
    }
    switch (key) {
      case 'ArrowUp': return '\x1b[1;5A'
      case 'ArrowDown': return '\x1b[1;5B'
      case 'ArrowRight': return '\x1b[1;5C'
      case 'ArrowLeft': return '\x1b[1;5D'
      default: return null
    }
  }

  if (alt) {
    switch (key) {
      case 'ArrowUp': return '\x1b\x1b[A'
      case 'ArrowDown': return '\x1b\x1b[B'
      case 'ArrowRight': return '\x1b\x1b[C'
      case 'ArrowLeft': return '\x1b\x1b[D'
      default:
        return key.length === 1 ? '\x1b' + key.toLowerCase() : null
    }
  }

  switch (key) {
    case 'ArrowUp': return '\x1b[A'
    case 'ArrowDown': return '\x1b[B'
    case 'ArrowRight': return '\x1b[C'
    case 'ArrowLeft': return '\x1b[D'
    case 'Home': return '\x1b[H'
    case 'End': return '\x1b[F'
    case 'Insert': return '\x1b[2~'
    case 'Delete': return '\x1b[3~'
    case 'PageUp': return '\x1b[5~'
    case 'PageDown': return '\x1b[6~'
    case 'Tab': return e.shiftKey ? '\x1b[Z' : '\t'
    default: return null
  }
}

export default function TerminalApp(_props: TerminalAppProps) {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const term = new Terminal({
      cursorBlink: true,
      fontFamily: "'Cascadia Code', 'Consolas', 'Courier New', monospace",
      fontSize: 14,
      lineHeight: 1.2,
      scrollback: 5000,
      allowTransparency: true,
      theme: {
        background: 'rgba(0, 0, 0, 0)',
        foreground: '#e8e8e8',
        cursor: '#e8e8e8',
        cursorAccent: '#202020',
        selectionBackground: 'rgba(255, 255, 255, 0.28)',
      },
    })
    const fit = new FitAddon()
    term.loadAddon(fit)
    term.open(container)

    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(`${proto}://${window.location.host}/api/system/pty`)

    let disposed = false

    const send = (obj: Record<string, unknown>) => {
      if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(obj))
    }

    ws.onopen = () => {
      if (disposed) return
      fit.fit()
      send({ type: 'resize', data: { cols: term.cols, rows: term.rows } })
      term.focus()
    }

    ws.onmessage = (ev) => {
      if (disposed) return
      try {
        const msg = JSON.parse(ev.data)
        if (msg && msg.type === 'stdout' && typeof msg.data === 'string') {
          term.write(msg.data)
        }
      } catch {
        // 忽略无法解析的帧
      }
    }

    ws.onerror = () => {
      // 失败后 onclose 会接着触发，统一在那里提示
    }

    ws.onclose = () => {
      if (disposed) return
      term.write('\r\n\x1b[31m[终端连接已断开]\x1b[0m\r\n')
    }

    const dataDisposable = term.onData((data) => send({ type: 'stdin', data }))
    const resizeDisposable = term.onResize(({ cols, rows }) =>
      send({ type: 'resize', data: { cols, rows } }),
    )

    term.attachCustomKeyEventHandler((event) => {
      const e = event as KeyboardEvent
      // 只处理 keydown，避免 keyup 再次触发导致一次按键触发两次
      if (e.type !== 'keydown') return true
      // IME 组合输入交给 xterm 处理
      if (e.isComposing || e.key === 'Process') return true
      const seq = encodeKey(e)
      if (seq == null) return true
      send({ type: 'stdin', data: seq })
      return false
    })

    // 窗口/容器尺寸变化时重新适配 PTY 尺寸
    let hadSize = false
    const ro = new ResizeObserver((entries) => {
      if (disposed) return
      const rect = entries[0]?.contentRect
      const valid = !!rect && rect.width > 0 && rect.height > 0
      if (!valid) {
        // 最小化（display:none）时尺寸为 0，跳过 fit，避免终端被压成 2×1 破坏缓冲区
        hadSize = false
        return
      }
      try {
        fit.fit()
        // 从隐藏恢复后强制重绘一次，修复布局错乱
        if (!hadSize) term.refresh(0, term.rows - 1)
        hadSize = true
      } catch {
        // 布局过程中容器不可测量时忽略
      }
    })
    ro.observe(container)

    const raf = requestAnimationFrame(() => {
      if (!disposed) {
        try {
          fit.fit()
        } catch {
          // ignore
        }
      }
    })

    return () => {
      disposed = true
      cancelAnimationFrame(raf)
      ro.disconnect()
      dataDisposable.dispose()
      resizeDisposable.dispose()
      try {
        ws.close()
      } catch {
        // ignore
      }
      term.dispose()
    }
  }, [])

  return <div ref={containerRef} className="terminal-root" />
}