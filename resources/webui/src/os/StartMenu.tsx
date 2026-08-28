import { memo, useState } from 'react'
import { useAuthStore } from '@/stores/auth'
import { useWMStore } from '@/stores/wm'
import { getPinnedStart } from '@/apps/registry'
import ColorIcon from '@/components/ColorIcon'
import { SearchRegular, PowerRegular } from '@fluentui/react-icons'

interface StartMenuProps {
  onClose: () => void
}

function StartMenu({ onClose }: StartMenuProps) {
  const { user, logout } = useAuthStore()
  const { openApp } = useWMStore()
  const pinned = getPinnedStart()
  const [confirming, setConfirming] = useState(false)

  const handleAppClick = (appId: string) => {
    openApp(appId)
    onClose()
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
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              handleAppClick('agents')
            }
          }}
        />
      </div>

      {/* 已固定应用 */}
      <div className="start-section">
        <div className="start-section-title">已固定</div>
        <div className="start-grid">
          {pinned.map((app) => (
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
      </div>

      {/* 推荐 */}
      <div className="start-section">
        <div className="start-section-title">推荐</div>
        <div className="start-recommend">
          <button className="start-recommend-item" onClick={() => handleAppClick('overview')}>
            安全运营中心概览
          </button>
          <button className="start-recommend-item" onClick={() => handleAppClick('agents')}>
            终端节点管理
          </button>
          <button className="start-recommend-item" onClick={() => handleAppClick('users')}>
            系统用户管理
          </button>
          <button className="start-recommend-item" onClick={() => handleAppClick('about')}>
            关于
          </button>
        </div>
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