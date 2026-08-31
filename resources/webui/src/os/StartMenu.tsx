import { memo, useState } from 'react'
import { useAuthStore } from '@/stores/auth'
import { useWMStore } from '@/stores/wm'
import { getAllApps } from '@/apps/registry'
import ColorIcon from '@/components/ColorIcon'
import { SearchRegular, PowerRegular } from '@fluentui/react-icons'

interface StartMenuProps {
  onClose: () => void
}

function StartMenu({ onClose }: StartMenuProps) {
  const { user, logout } = useAuthStore()
  const { openApp } = useWMStore()
  const allApps = getAllApps()
  const [query, setQuery] = useState('')
  const [confirming, setConfirming] = useState(false)

  const q = query.trim().toLowerCase()
  const matches = q
    ? allApps.filter((a) => a.title.toLowerCase().includes(q) || a.id.toLowerCase().includes(q))
    : allApps

  const handleAppClick = (appId: string) => {
    openApp(appId)
    onClose()
  }

  const handleSearchEnter = () => {
    if (matches.length === 1) handleAppClick(matches[0].id)
  }

  const handlePowerClick = () => {
    setConfirming(true)
  }

  const handleConfirmLogout = () => {
    logout()
    onClose()
  }

  const handleCancelLogout = () => {
    setConfirming(false)
  }

  return (
    <div className="start-menu" onClick={(e) => e.stopPropagation()}>
      {/* 搜索框 */}
      <div className="start-search">
        <SearchRegular className="start-search-icon" />
        <input
          className="start-search-input"
          type="text"
          placeholder="键入此处搜索"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') handleSearchEnter()
          }}
        />
      </div>

      {/* 所有应用 */}
      <div className="start-section">
        <div className="start-section-title">所有应用</div>
        {matches.length === 0 ? (
          <div className="start-empty">未找到匹配的应用</div>
        ) : (
          <div className="start-grid">
            {matches.map((app) => (
              <button
                key={app.id}
                className="start-grid-item"
                onClick={() => handleAppClick(app.id)}
              >
                <div className="start-grid-icon"><ColorIcon name={app.colorIcon} size={30} /></div>
                <span className="start-grid-label">{app.title}</span>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* 底部：用户 + 电源 */}
      <div className="start-footer">
        <div className="start-footer-user">
          {user?.avatar_url ? (
            <img src={user.avatar_url} alt={user.login} className="start-footer-avatar" />
          ) : (
            <div className="start-footer-avatar start-footer-avatar--placeholder" />
          )}
          <span className="start-footer-name">{user?.name ?? user?.login ?? '未登录'}</span>
        </div>
        <button className="start-footer-power" onClick={handlePowerClick} title="退出登录">
          <PowerRegular />
        </button>
      </div>

      {/* 退出确认 */}
      {confirming && (
        <div className="win-dialog-overlay" onClick={handleCancelLogout}>
          <div className="win-dialog" onClick={(e) => e.stopPropagation()}>
            <div className="win-dialog-title">退出登录</div>
            <div className="win-dialog-body">确定要退出当前账户吗？</div>
            <div className="win-dialog-footer">
              <button className="win-dialog-btn win-dialog-btn--subtle" onClick={handleCancelLogout}>
                取消
              </button>
              <button className="win-dialog-btn win-dialog-btn--accent" onClick={handleConfirmLogout}>
                退出登录
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default memo(StartMenu)