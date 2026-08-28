/// <reference types="vite/client" />

declare module 'virtual:wallpapers' {
  /** public/wallpapers 目录下扫描到的壁纸 */
  interface WallpaperItem {
    /** 显示名（去掉扩展名），如「壁纸」 */
    name: string
    /** 文件名（含扩展名），用于拼装 /wallpapers/<path>，如「壁纸.jpg」 */
    path: string
  }
  const wallpaperItems: WallpaperItem[]
  export default wallpaperItems
}