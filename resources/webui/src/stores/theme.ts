import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type ThemeMode = 'light' | 'dark' | 'system'
export type Wallpaper = 'bloom' | 'sunrise' | 'forest' | 'midnight' | 'custom'

interface ThemeStore {
  mode: ThemeMode
  wallpaper: Wallpaper
  customWallpaperUrl: string
  windowTransparency: number
  accent: string
  setMode: (mode: ThemeMode) => void
  setWallpaper: (wp: Wallpaper) => void
  setCustomWallpaper: (url: string) => void
  setWindowTransparency: (value: number) => void
  setAccent: (color: string) => void
}

/** Win11 默认主题色 */
export const DEFAULT_ACCENT = '#0067c0'

/** 解析 hex 颜色为 RGB 分量 */
export function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  const m = hex.trim().replace('#', '')
  if (m.length !== 6) return null
  const num = parseInt(m, 16)
  if (Number.isNaN(num)) return null
  return { r: (num >> 16) & 255, g: (num >> 8) & 255, b: num & 255 }
}

/** 根据亮度判断 accent 上的文字颜色（浅色底用深字，深色底用白字） */
export function accentContrast(hex: string): string {
  const rgb = hexToRgb(hex)
  if (!rgb) return '#ffffff'
  const luminance = (0.299 * rgb.r + 0.587 * rgb.g + 0.114 * rgb.b) / 255
  return luminance > 0.6 ? '#1a1a1a' : '#ffffff'
}

/** 生成更深的 accent 用于 hover 描边 / shadow ring */
export function accentDim(hex: string, factor = 0.78): string {
  const rgb = hexToRgb(hex)
  if (!rgb) return hex
  const toHex = (v: number) => Math.max(0, Math.min(255, Math.round(v))).toString(16).padStart(2, '0')
  return `#${toHex(rgb.r * factor)}${toHex(rgb.g * factor)}${toHex(rgb.b * factor)}`
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set) => ({
      mode: 'system',
      wallpaper: 'bloom',
      customWallpaperUrl: '',
      windowTransparency: 0.96,
      accent: DEFAULT_ACCENT,
      setMode: (mode: ThemeMode) => set({ mode }),
      setWallpaper: (wallpaper: Wallpaper) => set({ wallpaper }),
      setCustomWallpaper: (url: string) => set({ customWallpaperUrl: url }),
      setWindowTransparency: (windowTransparency: number) => set({ windowTransparency }),
      setAccent: (accent: string) => set({ accent }),
    }),
    { name: 'aegis.theme' },
  ),
)