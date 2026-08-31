import { createClient, type WebDAVClient, type FileStat } from 'webdav'

let _client: WebDAVClient | null = null

/** 获取（并缓存）WebDAV 客户端，指向后端 /api/webdav */
export function getWebdavClient(): WebDAVClient {
  if (!_client) {
    _client = createClient(`${window.location.origin}/api/webdav`, {
      withCredentials: true,
    })
  }
  return _client
}

export interface DirEntry {
  name: string
  path: string
  isDir: boolean
  size: number
  lastmod: string | null
}

/** 列出目录内容，目录在前、按名称排序 */
export async function listDirectory(path: string): Promise<DirEntry[]> {
  const client = getWebdavClient()
  const items: FileStat[] = await client.getDirectoryContents(path)

  const result: DirEntry[] = []
  for (const it of items) {
    if (!it || !it.basename) continue
    result.push({
      name: it.basename,
      path: it.filename,
      isDir: it.type === 'directory',
      size: it.size ?? 0,
      lastmod: it.lastmod ?? null,
    })
  }

  result.sort(
    (a, b) =>
      Number(b.isDir) - Number(a.isDir) || a.name.localeCompare(b.name),
  )
  return result
}