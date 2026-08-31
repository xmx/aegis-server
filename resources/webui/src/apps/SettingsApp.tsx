import { useState, useEffect } from 'react'
import { useThemeStore, type ThemeMode } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import wallpaperItems from 'virtual:wallpapers'
import lockscreenItems from 'virtual:lockscreens'
import Logo from '@/components/Logo'
import ColorIcon from '@/components/ColorIcon'
import type { BuildInfo } from '@/lib/types'
import {
  SettingsRegular,
  WeatherSunnyRegular,
  WeatherMoonRegular,
  DesktopMacRegular,
  ChevronRightRegular,
  SearchRegular,
} from '@fluentui/react-icons'

// ============================================================
// 导航定义
// ============================================================
interface NavItem {
  id: string
  label: string
  colorIcon: string
}

const navItems: NavItem[] = [
  { id: 'about', label: '关于', colorIcon: 'question_circle_48_color' },
  { id: 'home', label: '首页', colorIcon: 'home_48_color' },
  { id: 'system', label: '系统', colorIcon: 'settings_48_color' },
  { id: 'personalization', label: '个性化', colorIcon: 'paint_brush_32_color' },
  { id: 'account', label: '账户', colorIcon: 'person_48_color' },
]

// ============================================================
// SettingsApp
// ============================================================
export default function SettingsApp() {
  const [section, setSection] = useState('home')
  const { user } = useAuthStore()

  return (
    <div className="settings-layout">
      {/* 左侧导航 */}
      <nav className="settings-nav">
        {/* 账户 */}
        <div className="settings-nav-account" onClick={() => setSection('account')}>
          {user?.avatar_url ? (
            <img className="settings-nav-avatar" src={user.avatar_url} alt={user.login} />
          ) : (
            <div className="settings-nav-avatar settings-nav-avatar--placeholder">?</div>
          )}
          <div className="settings-nav-account-meta">
            <span className="settings-nav-account-name">{user?.name ?? user?.login ?? '未登录'}</span>
          </div>
        </div>

        {/* 搜索 */}
        <div className="settings-nav-search">
          <SearchRegular className="settings-nav-search-icon" />
          <input
            className="settings-nav-search-input"
            type="text"
            placeholder="查找设置"
          />
        </div>

        {/* 菜单项 */}
        <div className="settings-nav-items">
          {navItems.map((item) => {
            const active = section === item.id
            return (
              <button
                key={item.id}
                className={`settings-nav-item ${active ? 'settings-nav-item--active' : ''}`}
                onClick={() => setSection(item.id)}
              >
                <span className="settings-nav-icon">
                  <ColorIcon name={item.colorIcon} size={24} />
                </span>
                <span className="settings-nav-label">{item.label}</span>
              </button>
            )
          })}
        </div>
      </nav>

      {/* 右侧内容 */}
      <div className="settings-content">
        {section === 'home' && <SettingsHome onNavigate={setSection} />}
        {section === 'system' && <SettingsSystem />}
        {section === 'personalization' && <SettingsPersonalization />}
        {section === 'account' && <SettingsAccount />}
        {section === 'about' && <SettingsAbout />}
      </div>
    </div>
  )
}

// ============================================================
// 首页
// ============================================================
const homeCards = [
  {
    id: 'personalization',
    colorIcon: 'paint_brush_32_color',
    title: '个性化',
    sub: '背景、主题和颜色',
  },
  {
    id: 'account',
    colorIcon: 'person_48_color',
    title: '账户',
    sub: '用户信息与登录方式',
  },
  {
    id: 'system',
    colorIcon: 'settings_48_color',
    title: '系统',
    sub: '显示、通知等',
  },
  {
    id: 'about',
    colorIcon: 'question_circle_48_color',
    title: '关于',
    sub: '版本信息和依赖项',
  },
]

function SettingsHome({ onNavigate }: { onNavigate: (id: string) => void }) {
  const { user } = useAuthStore()
  const { mode } = useThemeStore()
  const [query, setQuery] = useState('')

  const modeLabel: Record<string, string> = {
    light: '浅色',
    dark: '深色',
    system: '跟随系统',
  }

  const q = query.trim().toLowerCase()
  const cards = homeCards.filter(
    (c) => !q || c.title.toLowerCase().includes(q) || c.sub.toLowerCase().includes(q),
  )

  return (
    <div className="settings-page">
      <h3 className="settings-section-title">设置首页</h3>
      <p className="settings-section-desc">管理中心的所有设置，包括个性化、账户和系统信息。</p>

      {/* 搜索框 */}
      <div className="settings-home-search">
        <SearchRegular className="settings-home-search-icon" />
        <input
          className="settings-home-search-input"
          type="text"
          placeholder="查找设置"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      {/* 推荐卡片 */}
      <div className="settings-home-grid">
        {cards.map((card) => (
          <button key={card.id} className="settings-home-card" onClick={() => onNavigate(card.id)}>
            <span className="settings-home-card-icon"><ColorIcon name={card.colorIcon} size={22} /></span>
            <span className="settings-home-card-text">
              <span className="settings-home-card-title">{card.title}</span>
              <span className="settings-home-card-sub">{card.sub}</span>
            </span>
            <ChevronRightRegular className="settings-home-card-arrow" />
          </button>
        ))}
        {cards.length === 0 && (
          <div className="app-empty">没有匹配的设置</div>
        )}
      </div>

      <div className="settings-home-summary">
        <span>当前主题：{modeLabel[mode]}</span>
        <span>用户：{user?.login ?? '—'}</span>
      </div>
    </div>
  )
}

// ============================================================
// 系统
// ============================================================
function SettingsSystem() {
  return (
    <div className="settings-page">
      <h3 className="settings-section-title">系统</h3>
      <p className="settings-section-desc">显示、通知和系统相关配置。</p>
      <div className="settings-placeholder">
        <SettingsRegular />
        <span>更多系统设置即将推出</span>
      </div>
    </div>
  )
}

// ============================================================
// 个性化
// ============================================================
const themes: { value: ThemeMode; label: string; icon: typeof WeatherSunnyRegular }[] = [
  { value: 'light', label: '浅色', icon: WeatherSunnyRegular },
  { value: 'dark', label: '深色', icon: WeatherMoonRegular },
  { value: 'system', label: '跟随系统', icon: DesktopMacRegular },
]

/** Win11 主题色色板 */
const accents: string[] = [
  '#0067c0', // 蓝色（默认）
  '#4cc2ff', // 浅蓝
  '#107c10', // 绿色
  '#ca5010', // 橙色
  '#d13438', // 红色
  '#b4009e', // 紫色
  '#8b2fd6', // 深紫
  '#038387', // 青色
  '#486c0a', // 橄榄绿
  '#6e4f00', // 琥珀
]

function SettingsPersonalization() {
  const {
    mode,
    wallpaper,
    lockWallpaper,
    windowTransparency,
    accent,
    singleClickOpen,
    setMode,
    setWallpaper,
    setLockWallpaper,
    setWindowTransparency,
    setAccent,
    setSingleClickOpen,
  } = useThemeStore()

  return (
    <div className="settings-page">
      <h3 className="settings-section-title">个性化</h3>
      <p className="settings-section-desc">自定义桌面的外观，包括主题、背景和颜色。</p>

      {/* 主题模式 */}
      <div className="settings-block">
        <h4 className="settings-block-title">主题模式</h4>
        <div className="settings-grid-4">
          {themes.map((t) => {
            const Icon = t.icon
            return (
              <button
                key={t.value}
                className={`settings-card ${mode === t.value ? 'settings-card--active' : ''}`}
                onClick={() => setMode(t.value)}
              >
                <Icon />
                <span>{t.label}</span>
              </button>
            )
          })}
        </div>
      </div>

      {/* 主题色 */}
      <div className="settings-block">
        <h4 className="settings-block-title">主题色</h4>
        <div className="settings-accent-row">
          {accents.map((color) => (
            <button
              key={color}
              className={`settings-accent-swatch ${accent === color ? 'settings-accent-swatch--active' : ''}`}
              style={{ backgroundColor: color }}
              onClick={() => setAccent(color)}
              title={color}
              aria-label={`主题色 ${color}`}
            >
              {accent === color && <span className="settings-accent-check">✓</span>}
            </button>
          ))}
        </div>
        <p className="settings-transparency-desc">
          主题色应用于按钮、链接和选中状态。
        </p>
      </div>

      {/* 桌面壁纸 */}
      <div className="settings-block">
        <h4 className="settings-block-title">桌面壁纸</h4>
        {wallpaperItems.length > 0 ? (
          <div className="settings-wallpaper-grid">
            {wallpaperItems.map((item) => (
              <button
                key={item.path}
                className={`settings-wallpaper-item ${wallpaper === item.path ? 'settings-wallpaper-item--active' : ''}`}
                onClick={() => setWallpaper(item.path)}
                title={item.name}
              >
                <img
                  className="settings-wallpaper-thumb"
                  src={`/wallpaper/${item.path}`}
                  alt={item.name}
                  loading="lazy"
                />
              </button>
            ))}
          </div>
        ) : (
          <p className="settings-transparency-desc">
            未发现桌面壁纸，请将图片放入 public/wallpaper/ 目录。
          </p>
        )}
      </div>

      {/* 锁屏壁纸 */}
      <div className="settings-block">
        <h4 className="settings-block-title">锁屏壁纸</h4>
        {lockscreenItems.length > 0 ? (
          <div className="settings-wallpaper-grid">
            {lockscreenItems.map((item) => (
              <button
                key={item.path}
                className={`settings-wallpaper-item ${lockWallpaper === item.path ? 'settings-wallpaper-item--active' : ''}`}
                onClick={() => setLockWallpaper(item.path)}
                title={item.name}
              >
                <img
                  className="settings-wallpaper-thumb"
                  src={`/lockscreen/${item.path}`}
                  alt={item.name}
                  loading="lazy"
                />
              </button>
            ))}
          </div>
        ) : (
          <p className="settings-transparency-desc">
            未发现锁屏壁纸，请将图片放入 public/lockscreen/ 目录。
          </p>
        )}
      </div>

      {/* 窗口透明度 */}
      <div className="settings-block">
        <h4 className="settings-block-title">窗口透明度</h4>
        <div className="settings-transparency">
          <input
            type="range"
            min="0.5"
            max="1"
            step="0.01"
            value={windowTransparency}
            onChange={(e) => setWindowTransparency(Number(e.target.value))}
            className="settings-transparency-slider"
            aria-label="窗口透明度"
            style={{
              background: `linear-gradient(to right, var(--accent) 0%, var(--accent) ${((windowTransparency - 0.5) / 0.5) * 100}%, var(--border) ${((windowTransparency - 0.5) / 0.5) * 100}%, var(--border) 100%)`,
            }}
          />
          <span className="settings-transparency-value">
            {Math.round(windowTransparency * 100)}%
          </span>
        </div>
        <p className="settings-transparency-desc">
          调整窗口的毛玻璃透明度，数值越低越透明。
        </p>
      </div>

      {/* 点击方式 */}
      <div className="settings-block">
        <h4 className="settings-block-title">点击方式</h4>
        <div className="settings-toggle-row">
          <div className="settings-toggle-text">
            <div className="settings-toggle-title">单击打开项目</div>
            <div className="settings-toggle-desc">指向时选中，单击打开（网页风格）</div>
          </div>
          <button
            type="button"
            role="switch"
            aria-checked={singleClickOpen}
            className={`settings-toggle ${singleClickOpen ? 'settings-toggle--on' : ''}`}
            onClick={() => setSingleClickOpen(!singleClickOpen)}
          >
            <span className="settings-toggle-thumb" />
          </button>
        </div>
        <p className="settings-transparency-desc">
          关闭时为 Win11 默认行为：单击选中、双击打开。
        </p>
      </div>
    </div>
  )
}

// ============================================================
// 账户
// ============================================================
function SettingsAccount() {
  const { user } = useAuthStore()

  if (!user) {
    return (
      <div className="settings-page">
        <h3 className="settings-section-title">账户</h3>
        <p className="settings-section-desc">未登录。请退出并重新登录以管理账户信息。</p>
      </div>
    )
  }

  return (
    <div className="settings-page">
      <h3 className="settings-section-title">账户</h3>
      <p className="settings-section-desc">管理你的账户和登录信息。</p>

      <div className="settings-block">
        <div className="settings-account-profile">
          <img
            className="settings-account-avatar"
            src={user.avatar_url}
            alt={user.login}
            onError={(e) => {
              (e.target as HTMLImageElement).style.display = 'none'
            }}
          />
          <div className="settings-account-info">
            <span className="settings-account-name">{user.name ?? user.login}</span>
            <span className="settings-account-login">{user.login}</span>
          </div>
        </div>
      </div>

      <div className="settings-block">
        <h4 className="settings-block-title">登录信息</h4>
        <div className="settings-info-grid">
          <SettingsInfoField label="用户名" value={user.login} />
          <SettingsInfoField label="显示名称" value={user.name ?? '—'} />
          <SettingsInfoField label="邮箱" value={user.email ?? '—'} />
          <SettingsInfoField label="登录方式" value={user.provider ?? '—'} />
        </div>
      </div>
    </div>
  )
}

// ============================================================
// 关于
// ============================================================
function SettingsAbout() {
  const [info, setInfo] = useState<BuildInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    api<BuildInfo>('/api/system/buildinfo')
      .then(setInfo)
      .catch(() => setError(true))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="settings-page">
      <h3 className="settings-section-title">关于</h3>
      <p className="settings-section-desc">系统的版本和构建信息。</p>

      {loading ? (
        <div className="app-loading">加载中...</div>
      ) : error ? (
        <div className="app-empty">无法获取构建信息</div>
      ) : info ? (
        <>
          <div className="settings-about-hero">
            <Logo size={48} />
            <div className="settings-about-version">v{info.version}</div>
          </div>

          <div className="settings-block">
            <h4 className="settings-block-title">基本信息</h4>
            <div className="settings-info-grid">
              <SettingsInfoField label="版本" value={info.version} />
              <SettingsInfoField label="修订" value={info.revision} />
              <SettingsInfoField label="用户名" value={info.username} />
              <SettingsInfoField label="工作目录" value={info.workdir} />
              <SettingsInfoField label="模块" value={info.module} />
              <SettingsInfoField label="提交时间" value={info.committed_at} />
              <SettingsInfoField label="编译时间" value={info.compiled_at} />
              <SettingsInfoField label="系统" value={`${info.goos} / ${info.goarch}`} />
            </div>
          </div>

          {info.build_info && (
            <>
              <div className="settings-block">
                <h4 className="settings-block-title">构建信息</h4>
                <div className="settings-info-grid">
                  <SettingsInfoField label="Go 版本" value={info.build_info.GoVersion} />
                  <SettingsInfoField
                    label="主模块"
                    value={`${info.build_info.Main.Path} ${info.build_info.Main.Version}`}
                  />
                </div>
              </div>

              <div className="settings-block">
                <h4 className="settings-block-title">依赖项 ({info.build_info.Deps.length})</h4>
                <div className="settings-deps">
                  {info.build_info.Deps.slice(0, 30).map((d) => (
                    <div key={d.Path} className="settings-dep">
                      <span className="settings-dep-path">{d.Path}</span>
                      <span className="settings-dep-ver">{d.Version}</span>
                    </div>
                  ))}
                </div>
              </div>
            </>
          )}
        </>
      ) : null}
    </div>
  )
}

// ============================================================
// 共享小组件
// ============================================================
function SettingsInfoField({ label, value }: { label: string; value?: string }) {
  return (
    <div className="settings-info-field">
      <span className="settings-info-label">{label}</span>
      <span className="settings-info-value">{value ?? '—'}</span>
    </div>
  )
}