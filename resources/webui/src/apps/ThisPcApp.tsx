import { useEffect, useState, useCallback, useMemo, useLayoutEffect, useRef } from 'react'
import { createPortal } from 'react-dom'
import { listDirectory, getWebdavClient, type DirEntry } from '@/lib/webdav'
import { formatBytes, formatTime } from '@/lib/format'
import { useWMStore } from '@/stores/wm'
import { useThemeStore } from '@/stores/theme'
import {
  SearchRegular,
  ChevronLeftRegular,
  ChevronRightRegular,
  ArrowUpRegular,
  ArrowClockwiseRegular,
  DocumentRegular,
  AddRegular,
  MoreHorizontalRegular,
  ListRegular,
  GridRegular,
  CopyRegular,
  DeleteRegular,
  RenameRegular,
  FolderZipFilled,
  ArrowDownloadRegular,
} from '@fluentui/react-icons'
import folderIcon from '@win11/icon/win/folder.png?url'
import thisPcIcon from '@win11/icon/win/thispc.png?url'
import photosIcon from '@win11/icon/photos.png?url'
import moviesIcon from '@win11/icon/movies.png?url'
import musicIcon from '@win11/icon/groove.png?url'
import codeIcon from '@win11/icon/code.png?url'
import terminalIcon from '@win11/icon/terminal.png?url'
import notepadIcon from '@win11/icon/notepad.png?url'
import excelIcon from '@win11/icon/excel.png?url'
import wordIcon from '@win11/icon/winWord.png?url'
import pptIcon from '@win11/icon/powerpoint.png?url'
import txtIcon from '@win11/icon/win/893.png?url'
import imageIcon from '@win11/icon/win/1085.png?url'
import isoIcon from '@win11/icon/win/1693.png?url'
// 具体文件类型图标（从 fluentui-system-icons/extract 提取）
import exeIcon from '@extract/taskmgr.png?url'
import docIcon from '@extract/doc.png?url'
import docxIcon from '@extract/docx.png?url'
import xlsIcon from '@extract/xls.png?url'
import xlsxIcon from '@extract/xlsx.png?url'
import pptLegacyIcon from '@extract/ppt.png?url'
import pptxIcon from '@extract/pptx.png?url'
import htmlIcon from '@extract/ie-html.png?url'
import cmdIcon from '@extract/cmd.png?url'
import rdpIcon from '@extract/mstsc.png?url'
import linuxIcon from '@extract/linux.png?url'

type ViewMode = 'list' | 'icons'
type SortKey = 'name' | 'date' | 'size'

interface Crumb {
  label: string
  path: string
}

function buildCrumbs(path: string): Crumb[] {
  const crumbs: Crumb[] = [{ label: '此电脑', path: '/' }]
  const segs = path.split('/').filter(Boolean)
  let acc = ''
  for (const seg of segs) {
    acc += `/${seg}`
    crumbs.push({ label: seg, path: acc })
  }
  return crumbs
}

function parentOf(path: string): string {
  if (path === '/') return '/'
  const idx = path.lastIndexOf('/')
  return idx <= 0 ? '/' : path.slice(0, idx)
}

function compareFn(key: SortKey, dir: 1 | -1) {
  return (a: DirEntry, b: DirEntry): number => {
    let r = 0
    if (key === 'name') r = a.name.localeCompare(b.name)
    else if (key === 'date') r = (a.lastmod ?? '').localeCompare(b.lastmod ?? '')
    else r = a.size - b.size
    return r * dir
  }
}

const IMAGE_EXT = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg', 'ico', 'tif', 'tiff', 'heic']
const VIDEO_EXT = ['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv', 'webm', 'm4v', 'mpg', 'mpeg', '3gp']
const AUDIO_EXT = ['mp3', 'wav', 'flac', 'aac', 'ogg', 'm4a', 'wma', 'opus', 'mid', 'midi']
const CODE_EXT = [
  'js', 'jsx', 'ts', 'tsx', 'go', 'py', 'java', 'c', 'h', 'cpp', 'hpp', 'cs', 'rs', 'rb',
  'php', 'swift', 'kt', 'kts', 'html', 'htm', 'css', 'scss', 'less', 'json', 'xml', 'yml',
  'yaml', 'toml', 'sql', 'vue', 'svelte', 'sh', 'bash', 'bat', 'cmd', 'ps1',
]
const DOC_EXT = ['doc', 'docx', 'docm', 'odt', 'rtf', 'pages']
const SHEET_EXT = ['xls', 'xlsx', 'xlsm', 'csv', 'ods', 'numbers']
const SLIDE_EXT = ['ppt', 'pptx', 'pptm', 'odp', 'key']
const TEXT_EXT = ['txt', 'md', 'markdown', 'log', 'text', 'conf', 'ini', 'cfg']
const BIN_EXT = ['exe', 'msi', 'bin', 'apk', 'deb', 'rpm', 'dmg', 'jar', 'appimage']
const ARCHIVE_EXT = ['zip', 'rar', '7z', 'tar', 'gz', 'tgz', 'bz2', 'xz', 'zst']

function extOf(name: string): string {
  const i = name.lastIndexOf('.')
  return i > 0 ? name.slice(i + 1).toLowerCase() : ''
}

function fileIconUrl(name: string): string | null {
  const ext = extOf(name)
  // 具体文件类型图标（优先级最高）
  if (ext === 'doc') return docIcon
  if (ext === 'docx' || ext === 'docm') return docxIcon
  if (ext === 'xls') return xlsIcon
  if (ext === 'xlsx' || ext === 'xlsm') return xlsxIcon
  if (ext === 'ppt') return pptLegacyIcon
  if (ext === 'pptx' || ext === 'pptm') return pptxIcon
  if (ext === 'html' || ext === 'htm') return htmlIcon
  if (ext === 'cmd' || ext === 'bat') return cmdIcon
  if (ext === 'sh' || ext === 'bash' || ext === 'zsh' || ext === 'fish') return linuxIcon
  if (ext === 'rdp') return rdpIcon
  if (ext === 'exe') return exeIcon
  // 明确规定的文件类型
  if (ext === 'txt' || ext === 'ini') return txtIcon
  if (ext === 'jpg' || ext === 'jpeg' || ext === 'png' || ext === 'gif') return imageIcon
  if (ext === 'iso') return isoIcon
  // 其余类型使用通用分类图标
  if (IMAGE_EXT.includes(ext)) return photosIcon
  if (VIDEO_EXT.includes(ext)) return moviesIcon
  if (AUDIO_EXT.includes(ext)) return musicIcon
  if (CODE_EXT.includes(ext)) return codeIcon
  if (BIN_EXT.includes(ext)) return terminalIcon
  if (DOC_EXT.includes(ext)) return wordIcon
  if (SHEET_EXT.includes(ext)) return excelIcon
  if (SLIDE_EXT.includes(ext)) return pptIcon
  if (TEXT_EXT.includes(ext)) return notepadIcon
  return null
}

function isTextFile(name: string): boolean {
  const ext = extOf(name)
  return TEXT_EXT.includes(ext) || CODE_EXT.includes(ext)
}

function isImageFile(name: string): boolean {
  return IMAGE_EXT.includes(extOf(name))
}

/** 是否能在应用内打开（目录 / 文本代码 / 图片） */
function canOpen(entry: DirEntry): boolean {
  return entry.isDir || isTextFile(entry.name) || isImageFile(entry.name)
}

function typeLabel(entry: DirEntry): string {
  if (entry.isDir) return '文件夹'
  const ext = extOf(entry.name)
  if (!ext) return '文件'
  const map: Record<string, string> = {
    png: 'PNG 图片', jpg: 'JPEG 图片', jpeg: 'JPEG 图片', gif: 'GIF 图片', webp: 'WebP 图片', svg: 'SVG 图片', bmp: 'BMP 图片', ico: '图标',
    mp4: 'MP4 视频', avi: 'AVI 视频', mkv: 'MKV 视频', mov: 'MOV 视频', webm: 'WebM 视频',
    mp3: 'MP3 音频', wav: 'WAV 音频', flac: 'FLAC 音频', aac: 'AAC 音频', ogg: 'OGG 音频', m4a: 'M4A 音频',
    doc: 'Word 文档', docx: 'Word 文档', xls: 'Excel 工作表', xlsx: 'Excel 工作表', csv: 'CSV 文件',
    ppt: 'PowerPoint 演示文稿', pptx: 'PowerPoint 演示文稿', pdf: 'PDF 文件',
    txt: '文本文档', md: 'Markdown 文件', markdown: 'Markdown 文件', json: 'JSON 文件',
    xml: 'XML 文件', html: 'HTML 文件', htm: 'HTML 文件', css: 'CSS 文件',
    js: 'JavaScript 文件', ts: 'TypeScript 文件', go: 'Go 源文件', py: 'Python 源文件',
    java: 'Java 源文件', c: 'C 源文件', cpp: 'C++ 源文件', cs: 'C# 源文件', rs: 'Rust 源文件',
    sh: 'Shell 脚本', bash: 'Shell 脚本', zsh: 'Shell 脚本', ps1: 'PowerShell 脚本',
    bat: 'Windows 批处理文件', cmd: 'Windows 批处理文件', rdp: '远程桌面连接',
    zip: '压缩文件', rar: '压缩文件', '7z': '压缩文件', tar: '压缩文件', gz: '压缩文件',
    tgz: '压缩文件', bz2: '压缩文件', xz: '压缩文件', zst: '压缩文件', iso: '光盘映像',
  }
  return map[ext] ?? `${ext.toUpperCase()} 文件`
}

function EntryIcon({ entry }: { entry: DirEntry }) {
  if (entry.isDir) return <img className="explorer-icon-img" src={folderIcon} alt="" draggable={false} />
  if (ARCHIVE_EXT.includes(extOf(entry.name))) {
    return <FolderZipFilled className="explorer-icon explorer-icon--zip" />
  }
  const src = fileIconUrl(entry.name)
  if (src) return <img className="explorer-icon-img" src={src} alt="" draggable={false} />
  return <DocumentRegular className="explorer-icon" />
}

export default function ThisPcApp() {
  const [path, setPath] = useState('/')
  const [entries, setEntries] = useState<DirEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [hist, setHist] = useState<string[]>(['/'])
  const [histIndex, setHistIndex] = useState(0)
  const [selected, setSelected] = useState<string | null>(null)
  const [query, setQuery] = useState('')
  const [sortKey, setSortKey] = useState<SortKey>('name')
  const [sortDir, setSortDir] = useState<1 | -1>(1)
  const [view, setView] = useState<ViewMode>('list')
  const { openApp } = useWMStore()
  const singleClickOpen = useThemeStore((s) => s.singleClickOpen)

  const [ctxMenu, setCtxMenu] = useState<{ x: number; y: number; entry: DirEntry | null } | null>(null)
  const [renameTarget, setRenameTarget] = useState<DirEntry | null>(null)
  const [renameValue, setRenameValue] = useState('')
  const [deleteTarget, setDeleteTarget] = useState<DirEntry | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)
  const [menuPos, setMenuPos] = useState<{ x: number; y: number } | null>(null)

  // 根据菜单实际尺寸把坐标收进视口内，保证边缘也能完整显示
  useLayoutEffect(() => {
    if (!ctxMenu) {
      setMenuPos(null)
      return
    }
    const el = menuRef.current
    if (!el) return
    const PAD = 8
    const w = el.offsetWidth
    const h = el.offsetHeight
    setMenuPos({
      x: Math.max(PAD, Math.min(ctxMenu.x, window.innerWidth - w - PAD)),
      y: Math.max(PAD, Math.min(ctxMenu.y, window.innerHeight - h - PAD)),
    })
  }, [ctxMenu])

  // 悬浮提示：跟随鼠标，边缘自动收进视口
  const tipRef = useRef<HTMLDivElement>(null)
  const tipTimer = useRef<number | null>(null)
  const [tip, setTip] = useState<{ x: number; y: number; entry: DirEntry } | null>(null)
  const [tipPos, setTipPos] = useState<{ x: number; y: number } | null>(null)

  useLayoutEffect(() => {
    if (!tip) {
      setTipPos(null)
      return
    }
    const el = tipRef.current
    if (!el) return
    const PAD = 8
    setTipPos({
      x: Math.max(PAD, Math.min(tip.x, window.innerWidth - el.offsetWidth - PAD)),
      y: Math.max(PAD, Math.min(tip.y, window.innerHeight - el.offsetHeight - PAD)),
    })
  }, [tip])

  const showTip = (e: React.MouseEvent, entry: DirEntry) => {
    const x = e.clientX + 16
    const y = e.clientY + 16
    if (tipTimer.current) window.clearTimeout(tipTimer.current)
    tipTimer.current = window.setTimeout(() => {
      setTip({ x, y, entry })
    }, 500)
  }
  const moveTip = (e: React.MouseEvent) => {
    setTip((t) => (t ? { ...t, x: e.clientX + 16, y: e.clientY + 16 } : t))
  }
  const hideTip = () => {
    if (tipTimer.current) window.clearTimeout(tipTimer.current)
    tipTimer.current = null
    setTip(null)
  }

  useEffect(() => {
    return () => {
      if (tipTimer.current) window.clearTimeout(tipTimer.current)
    }
  }, [])

  const load = useCallback(async (p: string) => {
    setLoading(true)
    setError('')
    try {
      setEntries(await listDirectory(p))
    } catch (e) {
      setError(e instanceof Error ? e.message : '无法读取目录')
      setEntries([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load(path)
  }, [path, load])

  const navigate = useCallback(
    (p: string) => {
      const next = [...hist.slice(0, histIndex + 1), p]
      setHist(next)
      setHistIndex(next.length - 1)
      setPath(p)
      setSelected(null)
      setQuery('')
    },
    [hist, histIndex],
  )

  const canBack = histIndex > 0
  const canForward = histIndex < hist.length - 1
  const canUp = path !== '/'

  const goBack = () => {
    if (!canBack) return
    setHistIndex((i) => {
      const idx = i - 1
      setPath(hist[idx])
      return idx
    })
    setSelected(null)
  }
  const goForward = () => {
    if (!canForward) return
    setHistIndex((i) => {
      const idx = i + 1
      setPath(hist[idx])
      return idx
    })
    setSelected(null)
  }
  const goUp = () => navigate(parentOf(path))
  const refresh = () => load(path)

  // 关闭右键菜单：点击外部或按 Esc
  useEffect(() => {
    if (!ctxMenu) return
    const onDown = (e: MouseEvent) => {
      const t = e.target as HTMLElement
      if (!t.closest('.explorer-context-menu')) setCtxMenu(null)
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setCtxMenu(null)
    }
    window.addEventListener('mousedown', onDown)
    window.addEventListener('keydown', onKey)
    return () => {
      window.removeEventListener('mousedown', onDown)
      window.removeEventListener('keydown', onKey)
    }
  }, [ctxMenu])

  const openItemMenu = (e: React.MouseEvent, entry: DirEntry) => {
    e.preventDefault()
    e.stopPropagation()
    setSelected(entry.path)
    setCtxMenu({ x: e.clientX, y: e.clientY, entry })
  }
  const openBlankMenu = (e: React.MouseEvent) => {
    e.preventDefault()
    setSelected(null)
    setCtxMenu({ x: e.clientX, y: e.clientY, entry: null })
  }

  const handleMenuAction = (action: string) => {
    const entry = ctxMenu?.entry
    setCtxMenu(null)
    if (!entry) {
      if (action === 'new') handleCreateFolder()
      else if (action === 'refresh') refresh()
      return
    }
    if (action === 'open') handleOpen(entry)
    else if (action === 'download') {
      window.open(getWebdavClient().getFileDownloadLink(entry.path), '_blank', 'noopener')
    } else if (action === 'copy') void navigator.clipboard?.writeText(entry.path)
    else if (action === 'rename') {
      setRenameValue(entry.name)
      setRenameTarget(entry)
    } else if (action === 'delete') {
      setDeleteTarget(entry)
    }
  }

  const confirmRename = async () => {
    const target = renameTarget
    if (!target) return
    const newName = renameValue.trim()
    if (!newName || newName === target.name || newName.includes('/')) {
      setRenameTarget(null)
      return
    }
    const dir = parentOf(target.path)
    const newPath = dir === '/' ? `/${newName}` : `${dir}/${newName}`
    try {
      await getWebdavClient().moveFile(target.path, newPath)
      await load(path)
    } catch (e) {
      setError(e instanceof Error ? e.message : '重命名失败')
    }
    setRenameTarget(null)
  }

  const confirmDelete = async () => {
    const target = deleteTarget
    if (!target) return
    try {
      await getWebdavClient().deleteFile(target.path)
      await load(path)
    } catch (e) {
      setError(e instanceof Error ? e.message : '删除失败')
    }
    setDeleteTarget(null)
  }

  const handleOpen = useCallback(
    (entry: DirEntry) => {
      if (entry.isDir) {
        navigate(entry.path)
      } else if (isTextFile(entry.name)) {
        openApp('editor', { props: { path: entry.path, name: entry.name }, title: entry.name })
      } else if (isImageFile(entry.name)) {
        openApp('imageviewer', { props: { path: entry.path, name: entry.name, size: entry.size }, title: entry.name })
      } else {
        const link = getWebdavClient().getFileDownloadLink(entry.path)
        window.open(link, '_blank', 'noopener')
      }
    },
    [navigate, openApp],
  )

  const handleCreateFolder = useCallback(async () => {
    const client = getWebdavClient()
    let name = '新建文件夹'
    let i = 2
    while (entries.some((e) => e.isDir && e.name === name)) {
      name = `新建文件夹 (${i})`
      i++
    }
    const target = path === '/' ? `/${name}` : `${path}/${name}`
    try {
      await client.createDirectory(target)
      await load(path)
    } catch (e) {
      setError(e instanceof Error ? e.message : '无法创建文件夹')
    }
  }, [entries, path, load])

  const setSort = (key: SortKey) => {
    if (key === sortKey) setSortDir((d) => (d === 1 ? -1 : 1))
    else {
      setSortKey(key)
      setSortDir(1)
    }
  }

  const display = useMemo(() => {
    const q = query.trim().toLowerCase()
    const filtered = q ? entries.filter((e) => e.name.toLowerCase().includes(q)) : entries
    const cmp = compareFn(sortKey, sortDir)
    const dirs = filtered.filter((e) => e.isDir).sort(cmp)
    const files = filtered.filter((e) => !e.isDir).sort(cmp)
    return [...dirs, ...files]
  }, [entries, query, sortKey, sortDir])

  const crumbs = useMemo(() => buildCrumbs(path), [path])
  const sortedIndicator = (key: SortKey) => (sortKey === key ? (sortDir === 1 ? '↑' : '↓') : '')

  return (
    <div className="app-page explorer">
      {/* 命令栏 */}
      <div className="explorer-commandbar">
        <button className="explorer-cmdbtn explorer-cmdbtn--primary" onClick={handleCreateFolder} title="新建文件夹">
          <AddRegular />
          <span>新建</span>
        </button>
        <span className="explorer-sep" />
        <button
          className={`explorer-cmdbtn ${view === 'list' ? 'explorer-cmdbtn--active' : ''}`}
          onClick={() => setView('list')}
          title="详细信息"
        >
          <ListRegular />
        </button>
        <button
          className={`explorer-cmdbtn ${view === 'icons' ? 'explorer-cmdbtn--active' : ''}`}
          onClick={() => setView('icons')}
          title="大图标"
        >
          <GridRegular />
        </button>
        <span className="explorer-sep" />
        <button className="explorer-cmdbtn" onClick={refresh} disabled={loading} title="刷新">
          <ArrowClockwiseRegular />
        </button>
        <button className="explorer-cmdbtn" title="更多">
          <MoreHorizontalRegular />
        </button>

        {/* 搜索框 */}
        <div className="explorer-search">
          <SearchRegular className="explorer-search-icon" />
          <input
            className="explorer-search-input"
            type="text"
            placeholder="在此电脑中搜索"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
      </div>

      {/* 地址栏 */}
      <div className="explorer-addressbar">
        <button className="explorer-navbtn" onClick={goBack} disabled={!canBack} title="后退">
          <ChevronLeftRegular />
        </button>
        <button className="explorer-navbtn" onClick={goForward} disabled={!canForward} title="前进">
          <ChevronRightRegular />
        </button>
        <button className="explorer-navbtn" onClick={goUp} disabled={!canUp} title="向上">
          <ArrowUpRegular />
        </button>

        <div className="explorer-breadcrumb">
          {crumbs.map((c, i) => (
            <span key={c.path} className="explorer-crumb">
              {i > 0 && <ChevronRightRegular className="explorer-crumb-sep" />}
              <button
                className={`explorer-crumb-btn ${i === crumbs.length - 1 ? 'explorer-crumb-btn--current' : ''}`}
                onClick={() => navigate(c.path)}
              >
                {i === 0 && <img className="explorer-crumb-icon" src={thisPcIcon} alt="" />}
                {c.label}
              </button>
            </span>
          ))}
        </div>
      </div>

      {/* 内容区 */}
      <div className="app-page-body" onContextMenu={openBlankMenu}>
        {loading ? (
          <div className="app-loading">正在加载...</div>
        ) : error ? (
          <div className="app-empty">{error}</div>
        ) : display.length === 0 ? (
          <div className="app-empty">{query ? '没有匹配的项目' : '此文件夹为空'}</div>
        ) : view === 'list' ? (
          <div className="explorer-list">
            <div className="explorer-list-header">
              <button className="explorer-th explorer-th--name" onClick={() => setSort('name')}>
                名称 <span className="explorer-sort-arrow">{sortedIndicator('name')}</span>
              </button>
              <button className="explorer-th" onClick={() => setSort('date')}>
                修改日期 <span className="explorer-sort-arrow">{sortedIndicator('date')}</span>
              </button>
              <span className="explorer-th explorer-th--static">类型</span>
              <button className="explorer-th explorer-th--right" onClick={() => setSort('size')}>
                大小 <span className="explorer-sort-arrow">{sortedIndicator('size')}</span>
              </button>
            </div>
            {display.map((entry) => (
              <div
                key={entry.path}
                className={`explorer-row ${selected === entry.path ? 'explorer-row--selected' : ''}`}
                onClick={() => (singleClickOpen ? handleOpen(entry) : setSelected(entry.path))}
                onDoubleClick={singleClickOpen ? undefined : () => handleOpen(entry)}
                onContextMenu={(e) => openItemMenu(e, entry)}
                onMouseEnter={(e) => {
                  if (singleClickOpen) setSelected(entry.path)
                  showTip(e, entry)
                }}
                onMouseMove={moveTip}
                onMouseLeave={hideTip}
              >
                <span className="explorer-cell explorer-cell--name">
                  <EntryIcon entry={entry} />
                  <span className="explorer-name">{entry.name}</span>
                </span>
                <span className="explorer-cell explorer-cell--muted">{formatTime(entry.lastmod ?? undefined)}</span>
                <span className="explorer-cell explorer-cell--muted">{entry.isDir ? '文件夹' : '文件'}</span>
                <span className="explorer-cell explorer-cell--size">{entry.isDir ? '' : formatBytes(entry.size)}</span>
              </div>
            ))}
          </div>
        ) : (
          <div className="explorer-icons">
            {display.map((entry) => (
              <div
                key={entry.path}
                className={`explorer-tile ${selected === entry.path ? 'explorer-tile--selected' : ''}`}
                onClick={() => (singleClickOpen ? handleOpen(entry) : setSelected(entry.path))}
                onDoubleClick={singleClickOpen ? undefined : () => handleOpen(entry)}
                onContextMenu={(e) => openItemMenu(e, entry)}
                onMouseEnter={(e) => {
                  if (singleClickOpen) setSelected(entry.path)
                  showTip(e, entry)
                }}
                onMouseMove={moveTip}
                onMouseLeave={hideTip}
              >
                <div className="explorer-tile-icon">
                  <EntryIcon entry={entry} />
                </div>
                <div className="explorer-tile-name">{entry.name}</div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* 状态栏 */}
      <div className="explorer-statusbar">
        <span>{display.length} 个项目</span>
        {selected && <span>已选择 1 个项目</span>}
      </div>

      {createPortal(
        <>
          {/* 悬浮提示 */}
          {tip && (
            <div ref={tipRef} className="explorer-tooltip" style={{ left: (tipPos ?? tip).x, top: (tipPos ?? tip).y }}>
              <div className="explorer-tooltip-head">
                <EntryIcon entry={tip.entry} />
                <span className="explorer-tooltip-name">{tip.entry.name}</span>
              </div>
              <div className="explorer-tooltip-row">
                <span className="explorer-tooltip-label">类型</span>
                <span>{typeLabel(tip.entry)}</span>
              </div>
              {!tip.entry.isDir && (
                <div className="explorer-tooltip-row">
                  <span className="explorer-tooltip-label">大小</span>
                  <span>{formatBytes(tip.entry.size)}</span>
                </div>
              )}
              <div className="explorer-tooltip-row">
                <span className="explorer-tooltip-label">修改日期</span>
                <span>{formatTime(tip.entry.lastmod ?? undefined)}</span>
              </div>
            </div>
          )}

          {/* 右键菜单 */}
          {ctxMenu && (
        <div ref={menuRef} className="context-menu explorer-context-menu" style={{ left: (menuPos ?? ctxMenu).x, top: (menuPos ?? ctxMenu).y }}>
          {ctxMenu.entry ? (
            <>
              {ctxMenu.entry.isDir ? (
              <button className="context-menu-item" onClick={() => handleMenuAction('open')}>
                <span className="context-menu-item-icon"><DocumentRegular /></span>
                <span>打开</span>
              </button>
            ) : canOpen(ctxMenu.entry) ? (
              <>
                <button className="context-menu-item" onClick={() => handleMenuAction('open')}>
                  <span className="context-menu-item-icon"><DocumentRegular /></span>
                  <span>打开</span>
                </button>
                <button className="context-menu-item" onClick={() => handleMenuAction('download')}>
                  <span className="context-menu-item-icon"><ArrowDownloadRegular /></span>
                  <span>下载</span>
                </button>
              </>
            ) : (
              <button className="context-menu-item" onClick={() => handleMenuAction('download')}>
                <span className="context-menu-item-icon"><ArrowDownloadRegular /></span>
                <span>下载</span>
              </button>
            )}
              <button className="context-menu-item" onClick={() => handleMenuAction('copy')}>
                <span className="context-menu-item-icon"><CopyRegular /></span>
                <span>复制路径</span>
              </button>
              <button className="context-menu-item" onClick={() => handleMenuAction('rename')}>
                <span className="context-menu-item-icon"><RenameRegular /></span>
                <span>重命名</span>
              </button>
              <button className="context-menu-item" onClick={() => handleMenuAction('delete')}>
                <span className="context-menu-item-icon"><DeleteRegular /></span>
                <span>删除</span>
              </button>
            </>
          ) : (
            <>
              <button className="context-menu-item" onClick={() => handleMenuAction('new')}>
                <span className="context-menu-item-icon"><AddRegular /></span>
                <span>新建文件夹</span>
              </button>
              <button className="context-menu-item" onClick={() => handleMenuAction('refresh')}>
                <span className="context-menu-item-icon"><ArrowClockwiseRegular /></span>
                <span>刷新</span>
              </button>
            </>
          )}
        </div>
      )}

      {/* 重命名对话框 */}
      {renameTarget && (
        <div className="win-dialog-overlay" onClick={(e) => { if (e.target === e.currentTarget) setRenameTarget(null) }}>
          <div className="win-dialog">
            <div className="win-dialog-title">重命名</div>
            <div className="win-dialog-body">
              <input
                className="explorer-rename-input"
                value={renameValue}
                onChange={(e) => setRenameValue(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') confirmRename()
                  else if (e.key === 'Escape') setRenameTarget(null)
                }}
                autoFocus
              />
            </div>
            <div className="win-dialog-footer">
              <button className="win-dialog-btn win-dialog-btn--subtle" onClick={() => setRenameTarget(null)}>取消</button>
              <button className="win-dialog-btn win-dialog-btn--accent" onClick={confirmRename}>确定</button>
            </div>
          </div>
        </div>
      )}

      {/* 删除确认对话框 */}
      {deleteTarget && (
        <div className="win-dialog-overlay" onClick={(e) => { if (e.target === e.currentTarget) setDeleteTarget(null) }}>
          <div className="win-dialog">
            <div className="win-dialog-title">删除</div>
            <div className="win-dialog-body">确定要删除「{deleteTarget.name}」吗？此操作不可撤销。</div>
            <div className="win-dialog-footer">
              <button className="win-dialog-btn win-dialog-btn--subtle" onClick={() => setDeleteTarget(null)}>取消</button>
              <button className="win-dialog-btn win-dialog-btn--accent" onClick={confirmDelete}>删除</button>
            </div>
          </div>
        </div>
      )}
        </>,
        document.body,
      )}
    </div>
  )
}