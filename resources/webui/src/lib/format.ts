export function formatTime(iso: string | undefined): string {
  if (!iso) return '-'
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '-'
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  const s = String(d.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${day} ${h}:${min}:${s}`
}

export function formatDuration(seconds: number | undefined): string {
  if (!seconds || seconds <= 0) return '-'
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m`
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (h < 24) return `${h}h ${m}m`
  const d = Math.floor(h / 24)
  const rh = h % 24
  return `${d}d ${rh}h`
}

export function formatBytes(bytes: number | undefined): string {
  if (!bytes || bytes <= 0) return '-'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

export function formatClock(d: Date): string {
  const h = String(d.getHours()).padStart(2, '0')
  const m = String(d.getMinutes()).padStart(2, '0')
  return `${h}:${m}`
}

export function formatDate(d: Date): string {
  const y = d.getFullYear()
  const m = d.getMonth() + 1
  const day = d.getDate()
  return `${y}/${m}/${day}`
}

/** Win11 锁屏风格的日期：如「2026年8月31日 星期一」 */
export function formatLockDate(d: Date): string {
  const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const y = d.getFullYear()
  const m = d.getMonth() + 1
  const day = d.getDate()
  return `${y}年${m}月${day}日 ${weekdays[d.getDay()]}`
}