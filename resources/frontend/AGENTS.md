# Aegis Frontend 设计规范

## 项目概述

Aegis 前端，基于 React 19 + Vite 8 + TypeScript 构建，使用 shadcn/ui 组件库。

## 技术栈

| 类别   | 技术                                             | 版本  |
|--------|--------------------------------------------------|-------|
| 框架   | React                                            | ^19.2 |
| 构建   | Vite                                             | ^8.2  |
| 类型   | TypeScript                                       | ~6.0  |
| 样式   | Tailwind CSS                                     | ^4.3  |
| 组件库 | shadcn/ui (base-nova)                            | ^4.18 |
| 底层库 | @base-ui/react                                   | ^1.7  |
| 图标   | lucide-react                                     | ^1.33 |
| 动画   | tw-animate-css                                   | ^1.4  |
| 工具   | class-variance-authority / clsx / tailwind-merge | -     |

## 目录结构

```
src/
├── components/
│   ├── ThemeProvider.tsx    # 暗黑/明亮模式 Provider
│   ├── FontSettings.tsx     # 字体设置弹窗组件
│   └── ui/                  # shadcn 基础组件
│       ├── button.tsx
│       ├── dialog.tsx
│       └── select.tsx
├── lib/
│   └── utils.ts             # cn() 工具函数
├── hooks/                   # 自定义 hooks
├── App.tsx                  # 入口页面
├── main.tsx                 # 应用入口
└── index.css                # 全局样式 + Tailwind + 主题变量
```

## 路径别名

- `@/` → `src/`
- 配置于 `vite.config.ts` 和 `tsconfig.app.json`

## 核心规则

### 语言

**所有 UI 文字必须使用中文**，不允许出现英文文本。例外：专有名词（如 React、shadcn）、代码标识符、技术术语。

### 控件一致性（最高优先级）

- **优先复用已有组件**：开发任何新功能前，必须先检查 `src/components/ui/` 目录下是否已有可用的 shadcn 组件，直接复用，禁止重复造轮子
- **风格必须统一**：如果不存在对应组件，开发新组件时必须参考现有组件的代码风格、Tailwind class 写法、变体（variant）模式、尺寸（size）体系，确保外观和行为与已有组件保持一致
- **和谐统一**：所有控件的圆角、间距、颜色、字体、过渡动画、交互反馈必须遵循项目全局设计 token，禁止出现风格割裂的自定义控件

### 主题系统

- **样式**：base-nova（Base UI 底层）
- **主题色**：neutral（中性色）
- **圆角**：`--radius: 0.25rem`（硬朗风格，全局生效）
- **暗黑模式**：通过 `.dark` class 切换，由 `ThemeProvider` 管理
- **状态支持**：`system`（跟随系统）、`light`、`dark`
- **存储**：`localStorage` 键 `theme`

### 颜色语义（CSS 变量）

所有颜色通过 CSS 变量定义，禁止使用硬编码颜色值（如 `bg-blue-500`）。使用语义化 token：

| Token                                | 用途               |
|--------------------------------------|--------------------|
| `background` / `foreground`          | 页面背景和默认文字 |
| `card` / `card-foreground`           | 卡片表面           |
| `popover` / `popover-foreground`     | 浮层               |
| `primary` / `primary-foreground`     | 主按钮和强调       |
| `secondary` / `secondary-foreground` | 次要按钮           |
| `muted` / `muted-foreground`         | 弱化/禁用状态      |
| `accent` / `accent-foreground`       | 悬停和强调         |
| `destructive`                        | 危险操作           |
| `border`                             | 默认边框           |
| `input`                              | 输入框边框         |
| `ring`                               | 聚焦环             |

### 字体系统

- **默认英文**：Space Mono（等宽）
- **中文**：系统默认，可选切换（微软雅黑、宋体、黑体、仿宋、楷体、思源黑体、思源宋体、苹方）
- **英文可选**：Space Mono、JetBrains Mono、Fira Code、Cascadia Code、Source Code Pro、IBM Plex Mono、Consolas、Menlo、Courier
  New
- 所有 Google Fonts 已在 `index.css` 预加载，字体切换通过 `style={{ fontFamily }}` 实现

### 组件操作

- 按钮禁用时必须隐藏手型指针（已在 `index.css` 全局处理）
- 所有交互元素必须有 `cursor: pointer`（button 全局已设置）
- 表单控件使用 `data-invalid` + `aria-invalid` 表示错误状态

### 滚动条

- **禁止使用系统默认滚动条**
- 自定义滚动条：默认透明不可见，鼠标悬停到容器时显示滚动条
- 宽度 6px，圆角，使用 `var(--border)` 作为轨道颜色，`var(--muted-foreground)` 作为悬停高亮色
- 同时兼容 Firefox（`scrollbar-width: thin`）和 WebKit（`::-webkit-scrollbar`）
- 任何需要滚动条的地方，不应额外添加样式，全局规则已覆盖

### 时间格式

- 统一使用 `yyyy-MM-dd HH:mm:ss` 格式，24 小时制
- 仅日期：`yyyy-MM-dd`
- 仅时间：`HH:mm:ss`

### 组件导入

```tsx
import {Button} from "@/components/ui/button"
import {cn} from "@/lib/utils"
import {useTheme} from "@/components/ThemeProvider"
```

### 样式规范

- 使用 `cn()` 合并 class，不使用模板字符串拼接
- 使用 `flex` + `gap-*` 代替 `space-x-*` / `space-y-*`
- 使用 `size-*` 代替 `w-* h-*`（当宽高相等时）
- 禁止暗色模式手动 `dark:` 覆盖，使用语义 token
- 禁止手动设置 `z-index`（覆盖层组件自带层级管理）

### shadcn CLI

- 组件通过 `npx shadcn@latest add <name>` 添加
- 注意：CLI 生成的组件需手动从 `@/` 目录移至 `src/` 对应位置（路径别名解析 bug）
- 已安装组件：button、dialog、select、input、separator、skeleton、tooltip、sheet、sidebar
- 组件底层库为 Base UI（`@base-ui/react`），非 Radix
- 图标库为 lucide（`lucide-react`）

## 可用 UI 组件

| 组件                                                                                             | 导入路径                 |
|--------------------------------------------------------------------------------------------------|--------------------------|
| Button                                                                                           | `@/components/ui/button` |
| Dialog, DialogTrigger, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter | `@/components/ui/dialog` |
| Select, SelectTrigger, SelectContent, SelectItem, SelectGroup, SelectLabel, SelectValue          | `@/components/ui/select` |

## 业务组件

### ThemeProvider (`src/components/ThemeProvider.tsx`)

- 提供 `useTheme()` hook，返回 `{ theme, resolvedTheme, setTheme }`
- 包裹在 `main.tsx` 的 `<ThemeProvider defaultTheme="system">` 中

### FontSettings (`src/components/FontSettings.tsx`)

- 字体设置 Dialog，内含中英文下拉选择 + 实时预览
- 导出 `FontSettings`、`getFontStyle()`、`CN_FONTS`、`EN_FONTS`

## 命令

```bash
npm run dev      # 开发服务器
npm run build    # 生产构建 (tsc + vite)
npm run lint     # oxlint 检查
npm run preview  # 预览构建产物
```