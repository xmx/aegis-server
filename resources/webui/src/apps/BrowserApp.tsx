import { useState } from 'react'
import {
  ArrowLeftRegular,
  ArrowRightRegular,
  ArrowClockwiseRegular,
  HomeRegular,
  SearchRegular,
  GlobeRegular,
} from '@fluentui/react-icons'

interface BrowserAppProps {
  windowId: string
  props?: Record<string, unknown>
}

interface Bookmark {
  name: string
  url: string
  color: string
  letter: string
}

const BOOKMARKS: Bookmark[] = [
  { name: '必应', url: 'https://www.bing.com', color: '#008373', letter: 'b' },
  { name: '百度', url: 'https://www.baidu.com', color: '#2932e1', letter: '百' },
  { name: '维基百科', url: 'https://www.wikipedia.org', color: '#54595d', letter: 'W' },
  { name: 'MDN', url: 'https://developer.mozilla.org', color: '#151718', letter: 'M' },
  { name: 'GitHub', url: 'https://github.com', color: '#24292f', letter: 'G' },
  { name: 'Gitee', url: 'https://gitee.com', color: '#c71d23', letter: 'G' },
]

/** 把地址栏输入规范化为完整 URL：无协议补 https，非网址则走搜索引擎 */
function normalizeInput(raw: string): string {
  const t = raw.trim()
  if (!t) return ''
  if (/^https?:\/\//i.test(t)) return t
  if (/^(localhost|[\w-]+(\.[\w-]+)+)(:\d+)?(\/\S*)?$/i.test(t)) return 'https://' + t
  return 'https://www.bing.com/search?q=' + encodeURIComponent(t)
}

export default function BrowserApp(_props: BrowserAppProps) {
  const [input, setInput] = useState('')
  const [current, setCurrent] = useState('') // '' 表示首页
  const [backStack, setBackStack] = useState<string[]>([])
  const [forwardStack, setForwardStack] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [frameKey, setFrameKey] = useState(0)

  const applyPage = (url: string) => {
    setCurrent(url)
    setInput(url)
    setLoading(url !== '')
  }

  const navigate = (raw: string) => {
    const target = normalizeInput(raw)
    if (!target) return
    if (current === target) {
      // 重复进入同一地址视为刷新
      setFrameKey((k) => k + 1)
      setLoading(true)
      return
    }
    setBackStack((s) => [...s, current])
    setForwardStack([])
    applyPage(target)
  }

  const goBack = () => {
    if (backStack.length === 0) return
    const prev = backStack[backStack.length - 1]
    setForwardStack((f) => [current, ...f])
    setBackStack(backStack.slice(0, -1))
    applyPage(prev)
  }

  const goForward = () => {
    if (forwardStack.length === 0) return
    const next = forwardStack[0]
    setBackStack((s) => [...s, current])
    setForwardStack(forwardStack.slice(1))
    applyPage(next)
  }

  const reload = () => {
    if (!current) return
    setFrameKey((k) => k + 1)
    setLoading(true)
  }

  const goHome = () => {
    if (current) setBackStack((s) => [...s, current])
    setForwardStack([])
    setCurrent('')
    setInput('')
    setLoading(false)
  }

  const handleLoad = () => setLoading(false)

  return (
    <div className="browser">
      {/* 工具栏 */}
      <div className="browser-toolbar">
        <button className="browser-btn" onClick={goBack} disabled={backStack.length === 0} title="后退">
          <ArrowLeftRegular />
        </button>
        <button className="browser-btn" onClick={goForward} disabled={forwardStack.length === 0} title="前进">
          <ArrowRightRegular />
        </button>
        <button className="browser-btn" onClick={reload} disabled={!current} title="刷新">
          <ArrowClockwiseRegular />
        </button>
        <button className="browser-btn" onClick={goHome} title="主页">
          <HomeRegular />
        </button>
        <form
          className="browser-address"
          onSubmit={(e) => {
            e.preventDefault()
            navigate(input)
          }}
        >
          <SearchRegular className="browser-address-icon" />
          <input
            className="browser-address-input"
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="搜索或输入网址"
            spellCheck={false}
          />
        </form>
      </div>

      {/* 加载进度条 */}
      {loading && <div className="browser-loading" />}

      {/* 内容区 */}
      {!current ? (
        <div className="browser-home">
          <div className="browser-home-inner">
            <div className="browser-home-logo">
              <GlobeRegular />
            </div>
            <div className="browser-home-title">浏览器</div>
            <form
              className="browser-home-search"
              onSubmit={(e) => {
                e.preventDefault()
                navigate(input)
              }}
            >
              <SearchRegular />
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="搜索或输入网址"
                spellCheck={false}
                autoFocus
              />
            </form>
            <div className="browser-bookmarks">
              {BOOKMARKS.map((b) => (
                <button key={b.url} className="browser-bookmark" onClick={() => navigate(b.url)}>
                  <span className="browser-bookmark-letter" style={{ background: b.color }}>
                    {b.letter}
                  </span>
                  <span className="browser-bookmark-name">{b.name}</span>
                </button>
              ))}
            </div>
            <p className="browser-home-hint">部分网站因安全策略（X-Frame-Options / CSP）可能无法显示。</p>
          </div>
        </div>
      ) : (
        <div className="browser-viewport">
          <iframe
            key={frameKey}
            src={current}
            onLoad={handleLoad}
            title="浏览器内容"
            referrerPolicy="no-referrer"
          />
        </div>
      )}
    </div>
  )
}