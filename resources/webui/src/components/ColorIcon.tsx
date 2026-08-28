import { memo } from 'react'

// 以 @icons 别名 + ?url 引入用到的彩色图标：Vite 只会打包这 9 个文件。
import homeIcon from '@icons/Home/SVG/ic_fluent_home_48_color.svg?url'
import settingsIcon from '@icons/Settings/SVG/ic_fluent_settings_48_color.svg?url'
import paintBrushIcon from '@icons/Paint Brush/SVG/ic_fluent_paint_brush_32_color.svg?url'
import personIcon from '@icons/Person/SVG/ic_fluent_person_48_color.svg?url'
import questionCircleIcon from '@icons/Question Circle/SVG/ic_fluent_question_circle_48_color.svg?url'
import chartMultipleIcon from '@icons/Chart Multiple/SVG/ic_fluent_chart_multiple_32_color.svg?url'
import agentsIcon from '@icons/Agents/SVG/ic_fluent_agents_48_color.svg?url'
import shieldIcon from '@icons/Shield/SVG/ic_fluent_shield_48_color.svg?url'
import peopleIcon from '@icons/People/SVG/ic_fluent_people_48_color.svg?url'
import laptopIcon from '@icons/Laptop/SVG/ic_fluent_laptop_48_color.svg?url'

interface ColorIconProps {
  /** 短键，见下方 ICON_PATHS */
  name: string
  size?: number
  className?: string
}

const ICON_PATHS: Record<string, string> = {
  home_48_color: homeIcon,
  settings_48_color: settingsIcon,
  paint_brush_32_color: paintBrushIcon,
  person_48_color: personIcon,
  question_circle_48_color: questionCircleIcon,
  chart_multiple_32_color: chartMultipleIcon,
  agents_48_color: agentsIcon,
  shield_48_color: shieldIcon,
  people_48_color: peopleIcon,
  laptop_48_color: laptopIcon,
}

/**
 * 渲染 Fluent System Icons 的彩色（color）图标。
 * 彩色图标为多色 SVG，无法用 currentColor 上色，因此以 <img> 方式加载。
 */
function ColorIcon({ name, size = 24, className }: ColorIconProps) {
  const src = ICON_PATHS[name]
  if (!src) return null
  return (
    <img
      src={src}
      width={size}
      height={size}
      alt=""
      className={className}
      style={{ display: 'block' }}
      draggable={false}
    />
  )
}

export default memo(ColorIcon)