import type { Plugin } from 'vite'
import path from 'node:path'
import fs from 'node:fs'

/** 桌面壁纸与锁屏壁纸所在的 public 子目录 */
const DIRS = {
  wallpaper: 'public/wallpaper',
  lockscreen: 'public/lockscreen',
} as const

const VIRTUAL_IDS = ['virtual:wallpapers', 'virtual:lockscreens'] as const

function scan(rel: string): { name: string; path: string }[] {
  const dir = path.resolve(import.meta.dirname, rel)
  try {
    return fs
      .readdirSync(dir)
      .filter((f) => /\.(jpe?g|png|webp)$/i.test(f))
      .map((f) => ({ name: path.parse(f).name, path: f }))
      .sort((a, b) => a.name.localeCompare(b.name))
  } catch {
    // 目录不存在时返回空列表
    return []
  }
}

/** 扫描 public 下的壁纸目录，通过虚拟模块暴露给前端 */
export function wallpaperListPlugin(): Plugin {
  return {
    name: 'aegis-wallpaper-list',
    resolveId(id) {
      if ((VIRTUAL_IDS as readonly string[]).includes(id)) return '\0' + id
      return undefined
    },
    load(id) {
      if (id === '\0virtual:wallpapers') return `export default ${JSON.stringify(scan(DIRS.wallpaper))}`
      if (id === '\0virtual:lockscreens') return `export default ${JSON.stringify(scan(DIRS.lockscreen))}`
      return undefined
    },
  }
}