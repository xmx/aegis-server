import { useEffect, useState, useCallback, useMemo } from 'react'
import Editor from '@monaco-editor/react'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { getWebdavClient } from '@/lib/webdav'
import '@/lib/monaco'

interface TextEditorProps {
  windowId: string
  props?: Record<string, unknown>
}

const LANG_MAP: Record<string, string> = {
  js: 'javascript', jsx: 'javascript', mjs: 'javascript', cjs: 'javascript',
  ts: 'typescript', tsx: 'typescript', mts: 'typescript',
  go: 'go', py: 'python', java: 'java',
  c: 'c', h: 'c', cpp: 'cpp', cc: 'cpp', hpp: 'cpp',
  cs: 'csharp', rs: 'rust', rb: 'ruby', php: 'php', swift: 'swift', kt: 'kotlin',
  html: 'html', htm: 'html', vue: 'html', svelte: 'html',
  css: 'css', scss: 'scss', less: 'less',
  json: 'json', xml: 'xml', yml: 'yaml', yaml: 'yaml', toml: 'ini',
  md: 'markdown', markdown: 'markdown',
  sql: 'sql', sh: 'shell', bash: 'shell', ps1: 'powershell',
  ini: 'ini', conf: 'ini', cfg: 'ini',
}

function langFor(name: string): string {
  const i = name.lastIndexOf('.')
  const ext = i > 0 ? name.slice(i + 1).toLowerCase() : ''
  return LANG_MAP[ext] ?? 'plaintext'
}

type ViewMode = 'code' | 'preview' | 'split'

export default function TextEditorApp({ props }: TextEditorProps) {
  const path = props?.path as string | undefined
  const name = (props?.name as string | undefined) ?? path?.split('/').filter(Boolean).pop() ?? '未命名'
  const isMarkdown = /\.(md|markdown)$/i.test(name)
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(Boolean(path))
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const [dirty, setDirty] = useState(false)
  const [mode, setMode] = useState<ViewMode>('code')

  useEffect(() => {
    if (!path) return
    let cancelled = false
    setLoading(true)
    getWebdavClient()
      .getFileContents(path, { format: 'text' })
      .then((v) => {
        if (cancelled) return
        const text = typeof v === 'string' ? v : new TextDecoder().decode(v as ArrayBuffer)
        setContent(text)
        setDirty(false)
      })
      .catch((e) => {
        if (!cancelled) setError(e instanceof Error ? e.message : '无法读取文件')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })
    return () => {
      cancelled = true
    }
  }, [path])

  const handleSave = useCallback(async () => {
    if (!path) return
    setSaving(true)
    try {
      const ok = await getWebdavClient().putFileContents(path, content)
      if (ok === false) throw new Error('保存失败')
      setDirty(false)
    } catch (e) {
      setError(e instanceof Error ? e.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }, [path, content])

  const previewHtml = useMemo(() => {
    if (!isMarkdown) return ''
    const raw = marked.parse(content, { async: false })
    return DOMPurify.sanitize(raw)
  }, [content, isMarkdown])

  return (
    <div className="editor-app">
      <div className="editor-toolbar">
        <span className="editor-filename">{name}</span>
        <span className="editor-meta" title={path}>{path}</span>
        <span className="editor-spacer" />

        {isMarkdown && (
          <>
            <button
              className="editor-toolbar-btn"
              onClick={() => setMode(mode === 'preview' ? 'code' : 'preview')}
              title="切换预览（Ctrl+Shift+V）"
            >
              {mode === 'preview' ? '源码' : '预览'}
            </button>
            <button
              className="editor-toolbar-btn"
              onClick={() => setMode(mode === 'split' ? 'code' : 'split')}
              title="并排预览（Ctrl+K V）"
            >
              {mode === 'split' ? '关闭分屏' : '分屏'}
            </button>
          </>
        )}

        {dirty && <span className="editor-dirty-dot" title="未保存的更改" />}
        <button className="editor-save-btn" onClick={handleSave} disabled={saving || !dirty || !path}>
          {saving ? '保存中…' : '保存'}
        </button>
      </div>

      <div className="editor-body">
        {loading ? (
          <div className="app-loading">正在加载文件...</div>
        ) : error ? (
          <div className="app-empty">{error}</div>
        ) : (
          <div className={`editor-stage ${mode === 'split' ? 'editor-stage--split' : ''}`}>
            {(mode === 'code' || mode === 'split') && (
              <div className="editor-pane">
                <Editor
                  height="100%"
                  language={langFor(name)}
                  value={content}
                  theme="vs-dark"
                  onChange={(v) => {
                    setContent(v ?? '')
                    setDirty(true)
                  }}
                  options={{
                    fontSize: 13,
                    minimap: { enabled: true },
                    automaticLayout: true,
                    scrollBeyondLastLine: false,
                    tabSize: 2,
                    wordWrap: 'off',
                  }}
                />
              </div>
            )}
            {(mode === 'preview' || mode === 'split') && (
              <div
                className="editor-preview markdown-body"
                dangerouslySetInnerHTML={{ __html: previewHtml }}
              />
            )}
          </div>
        )}
      </div>
    </div>
  )
}