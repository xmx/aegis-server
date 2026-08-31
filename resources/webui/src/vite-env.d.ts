/// <reference types="vite/client" />

declare module 'virtual:wallpapers' {
  /** public/wallpaper 目录下扫描到的桌面壁纸 */
  interface WallpaperItem {
    /** 显示名（去掉扩展名），如「default」 */
    name: string
    /** 文件名（含扩展名），用于拼装 /wallpaper/<path>，如「default.jpg」 */
    path: string
  }
  const wallpaperItems: WallpaperItem[]
  export default wallpaperItems
}

declare module 'virtual:lockscreens' {
  /** public/lockscreen 目录下扫描到的锁屏壁纸 */
  interface LockscreenItem {
    /** 显示名（去掉扩展名），如「default」 */
    name: string
    /** 文件名（含扩展名），用于拼装 /lockscreen/<path>，如「default.jpg」 */
    path: string
  }
  const lockscreenItems: LockscreenItem[]
  export default lockscreenItems
}