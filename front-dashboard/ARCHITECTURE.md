# 🏗️ CheYa 企业级 SaaS 架构文档

## 📁 新项目结构

```
front-dashboard/
├── src/
│   ├── layouts/
│   │   └── SaaSLayout.tsx          ✨ 新增：企业级布局（侧边栏+内容区）
│   ├── pages/
│   │   ├── LoginPage.tsx           ✅ 保留：登录页面
│   │   ├── RegisterPage.tsx        ✅ 保留：注册页面
│   │   ├── DashboardOverviewPage.tsx  ✨ 新增：概览页面（迁移自旧 Dashboard）
│   │   ├── DashboardMapPage.tsx    ✨ 新增：地图页面占位符
│   │   └── DashboardListPage.tsx   ✨ 新增：列表页面占位符
│   ├── App.tsx                     🔄 更新：嵌套路由配置
│   ├── main.tsx
│   └── index.css
└── ...
```

## 🛣️ 路由架构

### 嵌套路由结构

```
/                           → 重定向到 /dashboard/overview
/login                      → LoginPage（独立页面）
/register                   → RegisterPage（独立页面）
/dashboard                  → SaaSLayout（布局容器）
  ├── /dashboard/overview   → DashboardOverviewPage
  ├── /dashboard/map        → DashboardMapPage
  └── /dashboard/list       → DashboardListPage
```

### 路由代码

```tsx
<Routes>
  <Route path="/" element={<Navigate to="/dashboard/overview" replace />} />
  <Route path="/login" element={<LoginPage />} />
  <Route path="/register" element={<RegisterPage />} />
  
  {/* 嵌套路由 */}
  <Route path="/dashboard" element={<SaaSLayout />}>
    <Route path="overview" element={<DashboardOverviewPage />} />
    <Route path="map" element={<DashboardMapPage />} />
    <Route path="list" element={<DashboardListPage />} />
  </Route>
</Routes>
```

## 🎨 SaaSLayout 组件架构

### 布局结构

```
┌─────────────────────────────────────────────────────┐
│                    整体容器 (flex)                    │
├──────────┬──────────────────────────────────────────┤
│          │                                          │
│          │           Header (顶部导航栏)             │
│  Sidebar │  ┌────────────────────────────────────┐  │
│  (侧边栏) │  │ 标题 + 登出按钮                     │  │
│          │  └────────────────────────────────────┘  │
│  - Logo  │                                          │
│  - 导航  │           Main (主内容区)                 │
│    链接  │  ┌────────────────────────────────────┐  │
│          │  │                                    │  │
│  - 概览  │  │      <Outlet />                    │  │
│  - 地图  │  │  (子路由在此渲染)                   │  │
│  - 列表  │  │                                    │  │
│          │  │                                    │  │
│  - Footer│  └────────────────────────────────────┘  │
│          │                                          │
└──────────┴──────────────────────────────────────────┘
```

### Sidebar 特性

- **固定宽度**: `w-64` (256px)
- **深色背景**: `bg-gray-800`
- **全屏高度**: `h-screen`
- **导航链接**: 使用 `<NavLink>` 实现自动高亮
  - 激活状态: 蓝色背景 (`bg-blue-600`)
  - 未激活: 灰色文字，hover 时深灰背景
- **图标**: 使用 SVG Heroicons

### Content Area 特性

- **Header**: 白色背景，包含平台标题和登出按钮
- **Main**: 
  - 灰色背景 (`bg-gray-100`)
  - 可滚动 (`overflow-y-auto`)
  - 内边距 `p-6`
  - 渲染 `<Outlet />` 组件

## 📄 页面说明

### 1. DashboardOverviewPage（概览页面）

**内容**：完整迁移自旧的 DashboardPage
- 地图占位符（500px 高度）
- 车队列表表格（左侧 70%）
- 车辆仪表盘（右侧 30%）

**路由**：`/dashboard/overview`

### 2. DashboardMapPage（地图页面）

**内容**：全屏地图占位符
- 大图标 + 居中文字
- 灰色背景

**路由**：`/dashboard/map`

**未来扩展**：集成 Google Maps / Mapbox / Leaflet

### 3. DashboardListPage（列表页面）

**内容**：专用车辆列表占位符
- 大图标 + 居中文字
- 白色背景

**路由**：`/dashboard/list`

**未来扩展**：完整的车辆列表 CRUD 功能

## 🔄 迁移变更总结

| 项目 | 旧版本 | 新版本 | 状态 |
|------|-------|--------|------|
| 路由模式 | 平面路由 | 嵌套路由 | ✅ 已升级 |
| 布局方式 | 每页独立 | 统一 Layout | ✅ 已升级 |
| Dashboard | 单一页面 | 拆分为 3 个页面 | ✅ 已完成 |
| 导航方式 | 无侧边栏 | Sidebar 导航 | ✅ 已添加 |
| 代码复用 | 低 | 高 | ✅ 已优化 |

## 🚀 如何使用

### 启动开发服务器

```bash
cd front-dashboard
npm run dev
```

### 访问页面

- **默认首页**: http://localhost:5173 → 自动跳转到 `/dashboard/overview`
- **概览页面**: http://localhost:5173/dashboard/overview
- **地图页面**: http://localhost:5173/dashboard/map
- **列表页面**: http://localhost:5173/dashboard/list
- **登录页面**: http://localhost:5173/login
- **注册页面**: http://localhost:5173/register

### 导航测试

1. 访问任意 dashboard 子页面
2. 左侧 Sidebar 会自动高亮当前激活的导航项
3. 点击不同的导航链接，内容区域会切换，但 Layout 保持不变

## 🎯 架构优势

### 1. **可维护性提升**
- 统一的布局管理
- 清晰的文件组织结构
- 易于添加新页面

### 2. **用户体验改善**
- 持久化的侧边栏导航
- 页面切换无需重新渲染 Header/Sidebar
- 现代化的 SaaS 界面

### 3. **可扩展性**
- 嵌套路由支持无限层级
- Layout 可轻松添加权限控制
- 便于集成状态管理

### 4. **企业级特性**
- 专业的侧边栏导航
- 统一的顶部操作栏
- 响应式设计基础

## 🔮 后续优化建议

1. **权限管理**: 在 SaaSLayout 中添加路由守卫
2. **面包屑导航**: 在 Header 中显示当前路径
3. **响应式侧边栏**: 移动端支持收起/展开
4. **多语言支持**: i18n 国际化
5. **主题切换**: 支持亮色/暗色主题
6. **子菜单**: 支持侧边栏多级菜单
7. **用户信息**: 在 Sidebar 底部显示当前用户

## 📝 注意事项

- 所有 Dashboard 子页面默认继承 `bg-gray-100` 背景
- 子页面无需再包含 Header，统一由 Layout 管理
- 使用 `<Outlet />` 渲染子路由内容
- `<NavLink>` 会自动添加 `active` 类用于样式控制

---

**架构设计**: 企业级 SaaS 标准
**更新时间**: 2025-11-11
**版本**: v2.0 (嵌套路由架构)



















