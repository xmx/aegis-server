import type { Plugin } from 'vite'
import path from 'node:path'
import fs from 'node:fs'

const WALLPAPER_DIR = 'public/wallpapers'
const VIRTUAL_ID = 'virtual:wallpapers'
const RESOLVED_VIRTUAL_ID = '\0' + VIRTUAL_ID

/** 扫描 public/wallpapers 下的图片文件，通过虚拟模块暴露给前端 */
export function wallpaperListPlugin(): Plugin {
  return {
    name: 'aegis-wallpaper-list',
    resolveId(id) {
      if (id === VIRTUAL_ID) return RESOLVED_VIRTUAL_ID
      return undefined
    },
    load(id) {
      if (id !== RESOLVED_VIRTUAL_ID) return undefined
      const dir = path.resolve(import.meta.dirname, WALLPAPER_DIR)
      const items: { name: string; path: string }[] = []
      try {
        items.push(
          ...fs
            .readdirSync(dir)
            .filter((f) => /\.(jpe?g|png|webp)$/i.test(f))
            .map((f) => ({ name: path.parse(f).name, path: f }))
            .sort((a, b) => a.name.localeCompare(b.name)),
        )
      } catch {
        // 目录不存在时返回空列表
      }
      return `export default ${JSON.stringify(items)}`
    },
  }
}