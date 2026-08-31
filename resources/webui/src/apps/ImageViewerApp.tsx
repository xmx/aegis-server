import { useState } from 'react'
import { getWebdavClient } from '@/lib/webdav'
import { formatBytes } from '@/lib/format'

interface ImageViewerProps {
  windowId: string
  props?: Record<string, unknown>
}

export default function ImageViewerApp({ props }: ImageViewerProps) {
  const path = props?.path as string | undefined
  const name = (props?.name as string | undefined) ?? path?.split('/').filter(Boolean).pop() ?? '图片'
  const size = typeof props?.size === 'number' ? (props.size as number) : undefined

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [dims, setDims] = useState<{ w: number; h: number } | null>(null)

  const src = path ? getWebdavClient().getFileDownloadLink(path) : ''

  return (
    <div className="image-viewer">
      <div className="image-viewer-canvas">
        {error || !src ? (
          <div className="app-empty">无法加载图片</div>
        ) : (
          <img
            className="image-viewer-img"
            src={src}
            alt={name}
            onLoad={(e) => {
              const img = e.currentTarget
              setDims({ w: img.naturalWidth, h: img.naturalHeight })
              setLoading(false)
            }}
            onError={() => {
              setError(true)
              setLoading(false)
            }}
          />
        )}
      </div>
      <div className="image-viewer-statusbar">
        <span className="image-viewer-name">{name}</span>
        {dims && <span className="image-viewer-meta">{dims.w} × {dims.h}</span>}
        {size !== undefined && <span className="image-viewer-meta">{formatBytes(size)}</span>}
        {loading && !error && <span className="image-viewer-meta">加载中…</span>}
      </div>
    </div>
  )
}