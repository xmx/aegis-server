import { type ComponentType } from 'react'
import OverviewApp from '@/apps/OverviewApp'
import ThisPcApp from '@/apps/ThisPcApp'
import AgentsApp from '@/apps/AgentsApp'
import AgentDetailApp from '@/apps/AgentDetailApp'
import UsersApp from '@/apps/UsersApp'
import SettingsApp from '@/apps/SettingsApp'
import AboutApp from '@/apps/AboutApp'

export interface AppMeta {
  id: string
  title: string
  /** ColorIcon 里的短键，映射到 fluentui-system-icons/assets/ 下的彩色 SVG */
  colorIcon: string
  component: ComponentType<{ windowId: string; props?: Record<string, unknown> }>
  defaultSize: { width: number; height: number }
  minSize?: { width: number; height: number }
  single?: boolean
  pinnedDesktop?: boolean
  pinnedStart?: boolean
  pinnedTaskbar?: boolean
}

export interface AppProps {
  windowId: string
  props?: Record<string, unknown>
}

const registry: AppMeta[] = [
  {
    id: 'computer',
    title: '此电脑',
    colorIcon: 'laptop_48_color',
    component: ThisPcApp,
    defaultSize: { width: 720, height: 480 },
    minSize: { width: 480, height: 360 },
    single: true,
    pinnedDesktop: true,
  },
  {
    id: 'overview',
    title: '概览',
    colorIcon: 'chart_multiple_32_color',
    component: OverviewApp,
    defaultSize: { width: 800, height: 520 },
    minSize: { width: 480, height: 320 },
    pinnedStart: true,
    pinnedTaskbar: true,
  },
  {
    id: 'agents',
    title: '终端节点',
    colorIcon: 'agents_48_color',
    component: AgentsApp,
    defaultSize: { width: 960, height: 600 },
    minSize: { width: 640, height: 400 },
    pinnedStart: true,
    pinnedTaskbar: true,
  },
  {
    id: 'agent-detail',
    title: '终端详情',
    colorIcon: 'shield_48_color',
    component: AgentDetailApp,
    defaultSize: { width: 800, height: 600 },
    minSize: { width: 560, height: 420 },
    single: false,
  },
  {
    id: 'users',
    title: '系统用户',
    colorIcon: 'people_48_color',
    component: UsersApp,
    defaultSize: { width: 960, height: 600 },
    minSize: { width: 640, height: 400 },
    pinnedStart: true,
    pinnedTaskbar: true,
  },
  {
    id: 'settings',
    title: '设置',
    colorIcon: 'settings_48_color',
    component: SettingsApp,
    defaultSize: { width: 960, height: 600 },
    minSize: { width: 640, height: 420 },
    pinnedDesktop: true,
    pinnedStart: true,
  },
  {
    id: 'about',
    title: '关于',
    colorIcon: 'question_circle_48_color',
    component: AboutApp,
    defaultSize: { width: 640, height: 500 },
    minSize: { width: 400, height: 340 },
    pinnedStart: true,
  },
]

export function getApp(id: string): AppMeta | undefined {
  return registry.find((a) => a.id === id)
}

export function getPinnedDesktop(): AppMeta[] {
  return registry.filter((a) => a.pinnedDesktop)
}

export function getPinnedStart(): AppMeta[] {
  return registry.filter((a) => a.pinnedStart)
}

export function getPinnedTaskbar(): AppMeta[] {
  return registry.filter((a) => a.pinnedTaskbar)
}

export default registry