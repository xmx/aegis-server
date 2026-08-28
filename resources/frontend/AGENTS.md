# AI Agent 行为约束

> 以下规则适用于所有 AI 编程助手会话，请严格遵守。

## 工作范围

- **可编辑目录**：仅 `aegis-server/resources/frontend/`（即本项目目录）
- **只读范围**：项目根目录下的其他模块（`aegis-agent/`、`aegis-common/`、Go 源码等）只允许读取，禁止修改
- 所有新建文件、编辑操作都限定在 `resources/frontend/` 目录内

## 角色与语言

- 扮演专业前端开发者
- 与用户对话使用中文
- 所有 UI 文字使用中文

## 工程上下文

- 本项目是 Aegis 后端（Go）的前端子项目，位于 `aegis-server/resources/frontend/`
- 后端 API 统一前缀 `/api`，Vite 开发代理到 `https://127.0.0.1:8060`
- 构建产出目录 `aegis-server/resources/static/dist/`

---

# Aegis Frontend 设计规范

## 项目概述

Aegis 前端，基于 React 19 + Vite 8 + TypeScript 构建，使用 Microsoft Fluent UI React v9 组件库。

## 技术栈

| 类别   | 技术                                      | 版本      |
|--------|-------------------------------------------|-----------|
| 框架   | React                                     | ^19.2     |
| 构建   | Vite                                      | ^8.2      |
| 类型   | TypeScript                                | ~6.0      |
| 样式   | Tailwind CSS (工具类) + Fluent UI Griffel | ^4.3 / v9 |
| 组件库 | @fluentui/react-components                | ^9.74     |
| 图标   | @fluentui/react-icons                     | ^2.0      |

## 目录结构

```
src/
├── components/
│   ├── ThemeProvider.tsx    # 暗黑/明亮模式状态管理
│   ├── AuthProvider.tsx     # 用户认证状态
│   ├── FontSettings.tsx     # 字体设置弹窗组件
│   └── Icons.tsx            # 共享图标组件
├── lib/
│   ├── api.ts               # fetch 封装 + 错误处理
│   └── toast.tsx            # Fluent UI Toaster 封装
├── pages/
│   ├── Dashboard.tsx        # 主布局（侧边栏 + 顶栏）
│   ├── Home.tsx             # 首页
│   ├── Login.tsx            # 登录页
│   ├── GitHubCallback.tsx   # GitHub OAuth 回调
│   └── User.tsx             # 用户管理
├── App.tsx                  # 路由配置
├── main.tsx                 # 应用入口
└── index.css                # 全局样式 + 字体
```

## 核心规则

### 语言

**所有 UI 文字必须使用中文**。

### 主题系统

- 使用 `ThemeProvider` 管理 theme 状态（system/light/dark）
- `FluentProvider` 根据 resolvedTheme 切换 `webLightTheme` / `webDarkTheme`
- 状态持久化到 `localStorage` 键 `theme`

### 样式

- **Fluent UI makeStyles**：组件级样式首选 `makeStyles` + `tokens`
- **Tailwind**：仅用于简单工具类（flex、gap、padding 等）
- 颜色使用 Fluent UI tokens（`tokens.colorNeutralForeground1` 等）
- 禁止手动 `dark:` 覆盖，暗黑模式由 FluentProvider 自动处理

### 组件

- 统一使用 `@fluentui/react-components` 组件
- Button、Input、Dialog、Select、Table、Tooltip、Spinner、Text 等
- 图标使用 `@fluentui/react-icons`

### 字体系统

- 默认英文：Space Mono（等宽）
- 字体切换通过 `getFontStyle()` + `style={{ fontFamily }}` 实现

### 时间格式

- `yyyy-MM-dd HH:mm:ss`，24 小时制

### 后端 API

- 统一前缀 `/api`，使用相对路径调用
- Vite proxy 代理到 `https://127.0.0.1:8060`（`secure: false`）

### 错误处理

- 使用 `api<T>()` 封装所有请求
- 401 → 跳转 `/login`
- 其他错误 → `showError()` Fluent UI Toast 通知