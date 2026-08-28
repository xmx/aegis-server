import type { SVGProps } from 'react'

interface LogoProps extends Omit<SVGProps<SVGSVGElement>, 'viewBox' | 'xmlns' | 'fill'> {
  /** 尺寸（同时设置 width + height），默认 24 */
  size?: number
}

/**
 * 系统 Logo 组件，渲染 favicon.svg 的盾牌图形。
 * 始终保持品牌色 #FF5543，不接受外部颜色覆盖。
 *
 * 用法：
 *   <Logo size={48} />
 *   <Logo size={20} />
 */
export default function Logo({ size = 24, className, ...rest }: LogoProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      width={size}
      height={size}
      fill="#FF5543"
      className={className}
      aria-label="Logo"
      {...rest}
    >
      <path d="m9.663 6.991 3.238.862L10.8 0 7.065 1.005l1.227 4.611c.178.67.704 1.197 1.371 1.375M10.894 15.852l.1-.099.037.135L13.2 24l3.735-1.005-.914-3.434-2.443-9.15L.993 7.058 0 10.81l8.103 2.165.134.037-.097.099-5.94 5.944 2.729 2.752zM16.131 11.086l.867 3.25c.18.67.705 1.197 1.374 1.375l4.637 1.233.992-3.753-7.868-2.105zM15.464 8.522l1.824.067h.07c.515 0 1-.201 1.363-.567l3.077-3.08-2.728-2.751-2.987 2.982a1.94 1.94 0 0 0-.57 1.329l-.05 2.02" />
    </svg>
  )
}