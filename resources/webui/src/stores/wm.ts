import { create } from 'zustand'
import { getApp } from '@/apps/registry'

export interface WindowBounds {
  x: number
  y: number
  width: number
  height: number
}

export interface WinInstance {
  id: string
  appId: string
  title: string
  props?: Record<string, unknown>
  bounds: WindowBounds
  prevBounds: WindowBounds | null
  z: number
  minimized: boolean
  maximized: boolean
}

interface WMStore {
  windows: WinInstance[]
  activeId: string | null
  zTop: number

  openApp: (appId: string, opts?: { props?: Record<string, unknown>; title?: string }) => void
  closeWindow: (id: string) => void
  focusWindow: (id: string) => void
  minimizeWindow: (id: string) => void
  toggleMaximize: (id: string) => void
  toggleMinimize: (id: string) => void
  setBounds: (id: string, bounds: WindowBounds) => void
}

let _nextId = 1
function nextId(): string {
  return `w_${_nextId++}_${Date.now() % 100000}`
}

/** 为 appId(+props) 生成唯一键，用于单例去重 */
function singletonKey(appId: string, props?: Record<string, unknown>): string {
  if (!props) return appId
  const extra = Object.entries(props)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([k, v]) => `${k}=${v}`)
    .join('&')
  return extra ? `${appId}?${extra}` : appId
}

/** 生成新窗口的默认位置（居中 + 层叠偏移） */
function defaultBounds(appId: string, existing: WinInstance[]): WindowBounds {
  const meta = getApp(appId)
  const w = meta?.defaultSize?.width ?? 800
  const h = meta?.defaultSize?.height ?? 560
  const vw = window.innerWidth
  const vh = window.innerHeight - 48 // 任务栏

  const cx = Math.max(0, (vw - w) / 2)
  const cy = Math.max(0, (vh - h) / 2)
  const offset = (existing.length % 8) * 28

  return { x: cx + offset, y: cy + offset, width: w, height: h }
}

export const useWMStore = create<WMStore>((set, get) => ({
  windows: [],
  activeId: null,
  zTop: 100,

  openApp: (appId: string, opts?: { props?: Record<string, unknown>; title?: string }) => {
    const meta = getApp(appId)
    const title = opts?.title ?? meta?.title ?? appId
    const props = opts?.props

    // 单例模式下，检查是否已有同 appId 的窗口
    if (meta?.single !== false) {
      const key = singletonKey(appId, props)
      const existing = get().windows.find((w) => {
        const ek = singletonKey(w.appId, w.props)
        return ek === key
      })
      if (existing) {
        // 还原并聚焦
        get().focusWindow(existing.id)
        if (existing.minimized) get().toggleMinimize(existing.id)
        return
      }
    }

    const bounds = defaultBounds(appId, get().windows)
    const id = nextId()
    const z = get().zTop + 1

    set((s) => ({
      windows: [...s.windows, { id, appId, title, props, bounds, prevBounds: null, z, minimized: false, maximized: false }],
      activeId: id,
      zTop: z,
    }))
  },

  closeWindow: (id: string) => {
    set((s) => {
      const remaining = s.windows.filter((w) => w.id !== id)
      let activeId = s.activeId
      if (activeId === id) {
        const top = remaining.length > 0
          ? (remaining.reduce((a, b) => (a.z > b.z ? a : b)) as WinInstance | null)
          : null
        activeId = top?.id ?? null
      }
      return { windows: remaining, activeId }
    })
  },

  focusWindow: (id: string) => {
    set((s) => {
      const z = s.zTop + 1
      return {
        windows: s.windows.map((w) => (w.id === id ? { ...w, z } : w)),
        activeId: id,
        zTop: z,
      }
    })
  },

  minimizeWindow: (id: string) => {
    set((s) => {
      const remaining = s.windows.filter((w) => w.id !== id)
      const updated = s.windows.map((w) => (w.id === id ? { ...w, minimized: true } : w))
      // 最小化后，聚焦剩余最上层窗口
      let activeId = s.activeId
      if (activeId === id) {
        const top = remaining.length > 0
          ? (remaining.reduce((a, b) => (a.z > b.z ? a : b)) as WinInstance | null)
          : null
        activeId = top?.id ?? null
      }
      return { windows: updated, activeId }
    })
  },

  toggleMaximize: (id: string) => {
    set((s) => ({
      windows: s.windows.map((w) => {
        if (w.id !== id) return w
        if (w.maximized) {
          return {
            ...w,
            maximized: false,
            bounds: w.prevBounds ?? w.bounds,
            prevBounds: null,
          }
        }
        return {
          ...w,
          maximized: true,
          prevBounds: { ...w.bounds },
          bounds: {
            x: 0,
            y: 0,
            width: window.innerWidth,
            height: window.innerHeight - 48,
          },
        }
      }),
    }))
  },

  toggleMinimize: (id: string) => {
    set((s) => ({
      windows: s.windows.map((w) =>
        w.id === id ? { ...w, minimized: !w.minimized } : w,
      ),
    }))
    // 恢复时聚焦
    const win = get().windows.find((w) => w.id === id)
    if (win?.minimized) {
      get().focusWindow(id)
    }
  },

  setBounds: (id: string, bounds: WindowBounds) => {
    set((s) => ({
      windows: s.windows.map((w) =>
        w.id === id ? { ...w, bounds } : w,
      ),
    }))
  },
}))