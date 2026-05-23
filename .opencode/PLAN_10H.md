# 10-Hour Autonomous Optimization Plan

> **目标**: 持续10小时自主执行，对 university-social 项目进行全方位打磨
> **范围**: PC前端/Admin后端/Uni-app移动端/CI-CD
> **策略**: 每完成一个Phase记录状态，构建验证不中断

---

## Phase 1: PC Frontend Deep Polish (~3h)

### 1.1 Mobile Responsive (30min)
- **文件**: `pc/src/layouts/MainLayout.vue` + `pc/src/styles/global.scss`
- **任务**:
  - `<768px` 侧边栏折叠为 hamburger menu
  - Sidebar `position: fixed` → `transform: translateX(-100%)` + toggle
  - Content area padding/margin 响应式适配
  - 卡片网格 `repeat(auto-fill, minmax(300px, 1fr))`
- **验证**: Playwright 截图 375px viewport

### 1.2 Skeleton Loading Screens (20min)
- **文件**: `pc/src/components/LoadingWrapper.vue` + 各列表页
- **任务**:
  - LoadingWrapper 添加 skeleton variant（`mode="skeleton"`）
  - Skeleton 卡片: `animate-pulse` + 占位块
  - 替换所有 `u-loading` → skeleton
- **验证**: 观察页面加载过渡

### 1.3 Post Detail with Threaded Comments (30min)
- **文件**: `pc/src/pages/post/detail.vue` + `pc/src/api/post.ts`
- **任务**:
  - 评论区改为 threaded（缩进回复）
  - 新增 `replyTo` 状态 + `@username`
  - 添加删除评论确认
  - 评论分页 (load-more)
- **验证**: 手动测试回复链

### 1.4 User Profile Page (20min)
- **文件**: `pc/src/pages/user/detail.vue`
- **任务**:
  - Tabs: 帖子/动态/粉丝/关注
  - 使用 PostCard/CircleCard 组件
  - Follow/unfollow 按钮
  - 统计信息卡片
- **验证**: 数据展示完整

### 1.5 Notification System (20min)
- **文件**: `pc/src/pages/notification/` + MainLayout dropdown
- **任务**:
  - Header 通知图标 + 未读徽标
  - 下拉通知列表 (最近5条)
  - 通知页面 (全列表 + 分页)
  - 标记已读 API
- **验证**: 通知标记已读

### 1.6 Search Functionality (20min)
- **文件**: `pc/src/pages/search/` + `pc/src/api/search.ts`
- **任务**:
  - 搜索框输入 → 实时搜索 (debounce 300ms)
  - 结果分类: 帖子/用户/圈子
  - 搜索结果卡片
  - 空状态 + 搜索历史
- **验证**: 搜索关键词返回结果

### 1.7 Image Upload Polish (20min)
- **文件**: `pc/src/components/ImageUploader.vue`
- **任务**:
  - 拖拽上传区域
  - 多图预览 + 删除
  - 上传进度条
  - 图片裁剪 (可选)
- **验证**: 上传图片成功

### 1.8 Follow/Unfollow Feed Integration (20min)
- **文件**: `pc/src/pages/user/detail.vue` + `pc/src/api/user.ts`
- **任务**:
  - Follow/unfollow 按钮 (optimistic update)
  - 关注者列表/正在关注列表
  - 主页"关注"tab 显示关注者动态
- **验证**: 关注/取关 状态同步

---

## Phase 2: Admin Console Enhancement (~2h)

### 2.1 Admin Login Page (20min)
- **文件**: `admin/src/pages/Login.vue`
- **任务**:
  - 独立登录页（非 dev-login）
  - admin 凭据验证
  - JWT 存 localStorage + 路由守卫
  - 开发模式快捷登录
- **验证**: 登录成功 → 跳转 Dashboard

### 2.2 Dashboard with Real Charts (30min)
- **文件**: `admin/src/pages/Dashboard.vue` + `DataScreen.vue`
- **任务**:
  - ECharts 折线图: 每日 DAU/新增帖子
  - ECharts 饼图: 圈子分布/学校分布
  - DataV digital-flop 动态数字
  - 数据刷新按钮
- **验证**: 图表渲染真实数据

### 2.3 User Management (20min)
- **文件**: `admin/src/pages/Users.vue`
- **任务**:
  - 用户搜索 (nickname/openID)
  - Ban/unban 确认对话框
  - 用户详情弹窗
  - 批量操作 (多选 + ban)
- **验证**: Ban 后用户无法登录

### 2.4 Content Moderation (20min)
- **文件**: `admin/src/pages/Posts.vue` + `Reports.vue`
- **任务**:
  - 帖子列表 + 内容预览
  - 删除帖子确认
  - 举报处理 (resolve/dismiss)
  - 敏感词管理 (CRUD)
- **验证**: 删除帖子 → 前端不可见

### 2.5 Banner/Announcement CRUD (15min)
- **文件**: `admin/src/pages/Banners.vue` + `Announcements.vue`
- **任务**:
  - 轮播图上传 + 排序
  - 公告创建/编辑/删除
  - 启用/禁用切换
- **验证**: 前端展示最新公告

### 2.6 DataV Polish (15min)
- **文件**: `admin/src/pages/Dashboard.vue` + `admin/src/App.vue`
- **任务**:
  - BorderBox 统一风格
  - ScrollBoard 实时滚动
  - 全屏按钮
- **验证**: 大屏效果截图

---

## Phase 3: Backend Hardening (~2h)

### 3.1 Fix Test Deadlocks (30min)
- **文件**: `backend/internal/handler/circle_test.go` + `search_test.go`
- **任务**:
  - 分析死锁原因 (goroutine + db tx)
  - 修复 `context.Background()` → `context.WithTimeout`
  - 添加 `defer cancel()`
  - `go test -race` 验证
- **验证**: `go test ./...` 通过

### 3.2 Rate Limiting Middleware (20min)
- **文件**: `backend/internal/middleware/ratelimit.go`
- **任务**:
  - Token bucket / sliding window
  - 每 IP 每 endpoint 限制
  - 429 响应
  - 可配置 limit/duration
- **验证**: 超限返回 429

### 3.3 Image Upload Validation (20min)
- **文件**: `backend/internal/handler/upload.go`
- **任务**:
  - 文件类型白名单 (jpg/png/webp/gif)
  - 文件大小限制 (10MB)
  - 图片压缩 (webp 转换)
  - 防重名 (uuid)
- **验证**: 上传非法类型 → 400

### 3.4 WebSocket Connection Management (20min)
- **文件**: `backend/internal/ws/hub.go`
- **任务**:
  - 连接池最大连接数
  - 心跳检测 (ping/pong)
  - 自动清理过期连接
  - 优雅关闭
- **验证**: 断线重连

### 3.5 API Pagination Optimization (15min)
- **文件**: `backend/internal/handler/post.go` + `circle.go`
- **任务**:
  - 统一分页参数 (offset/limit)
  - 默认 limit 20, max 100
  - cursor-based 分页 (可选)
- **验证**: offset+limit 正确

### 3.6 Error Handling Standardization (15min)
- **文件**: `backend/internal/handler/*.go`
- **任务**:
  - 统一错误响应格式 `{error: string, code: int}`
  - 日志记录 (请求ID + 错误堆栈)
  - 敏感错误不暴露给客户端
- **验证**: 错误格式一致

---

## Phase 4: CI/CD + DevOps (~1.5h)

### 4.1 Fix GitHub Actions Test Workflow (30min)
- **文件**: `.github/workflows/test.yml`
- **任务**:
  - 生成 `package-lock.json` (frontend + admin)
  - PC 前端构建步骤
  - 后端 `go test -race`
  - 缓存 Go modules + node_modules
- **验证**: Workflow 绿色

### 4.2 PC Frontend CI Job (15min)
- **文件**: `.github/workflows/deploy.yml`
- **任务**:
  - 添加 `pc-frontend` 构建 job
  - 上传 dist 到 Pages artifact
  - 添加 `pc/` 子目录
- **验证**: 部署后 PC 前端可访问

### 4.3 Dockerfile Multi-stage (20min)
- **文件**: `backend/Dockerfile`
- **任务**:
  - Stage 1: `golang:1.23-alpine` → build
  - Stage 2: `alpine:3.19` → binary only (15MB)
  - Health check 端点
  - 非 root 用户
- **验证**: `docker build` + `docker run`

### 4.4 Environment Configuration (15min)
- **文件**: `backend/config/config.go`
- **任务**:
  - `.env` → `os.Getenv` 映射
  - 开发/生产/测试 配置分离
  - 敏感字段 (DB password/JWT secret) 环境变量
- **验证**: 配置覆盖正确

### 4.5 Production Deploy Scripts (10min)
- **文件**: `scripts/deploy.sh` + `scripts/init-db.sh`
- **任务**:
  - 数据库 migration 自动执行
  - 二进制替换 + 优雅重启
  - 日志轮转
- **验证**: 部署流程可执行

---

## Phase 5: Uni-app Mobile Refinement (~1.5h)

### 5.1 Loading State Unification (20min)
- **文件**: `frontend/src/pages/**/*.vue`
- **任务**:
  - 所有页面 `u-loading` 一致性检查
  - 添加 missing loading ref
  - LoadingWrapper pattern 导入
  - 统一 loading 样式 (flower 模式)
- **验证**: 所有页面有 loading 状态

### 5.2 Pull-to-Refresh Patterns (15min)
- **文件**: `frontend/src/pages/square/index.vue` + 列表页
- **任务**:
  - `u-refresh` 组件
  - 下拉刷新 → 重新加载
  - 加载更多 (触底)
- **验证**: 下拉触发刷新

### 5.3 Empty State CTAs (15min)
- **文件**: `frontend/src/pages/**/*.vue`
- **任务**:
  - 空状态 → `u-empty` + CTA 按钮
  - "发布第一条动态" / "创建圈子" 等
  - 引导文案
- **验证**: 空状态有操作引导

### 5.4 Page Transitions Polish (10min)
- **文件**: `frontend/src/pages.json`
- **任务**:
  - 统一页面切换动画
  - fade-in 200ms
  - slide-left/slide-right (导航方向)
- **验证**: 切换流畅

### 5.5 Error Handling Consistency (15min)
- **文件**: `frontend/src/pages/**/*.vue`
- **任务**:
  - 所有 try/catch 有用户提示
  - 网络错误 → "网络异常，请重试"
  - 401 → 跳转登录
- **验证**: 断网显示提示

### 5.6 Image Lazy Load (15min)
- **文件**: `frontend/src/pages/**/*.vue`
- **任务**:
  - `u-lazy-load` 组件
  - 图片占位符
  - 加载失败 fallback
- **验证**: 图片懒加载生效

---

## Execution Rules

1. **顺序执行** Phase1→2→3→4→5, 每 Phase 内任务可并行
2. **构建验证** 每完成 3 个任务运行一次 build
3. **状态记录** 每个 Phase 完成后写入 `.opencode/PROJECT_STATUS.md`
4. **遇到卡住** 最多尝试 3 次修复, 否则跳过 + 记录
5. **截图验证** 每个 Phase 完成后 Playwright 截图
6. **ELF 学习** 每次故障修复 → 保存 ELF learning
7. **10小时预算** Phase1-3h / Phase2-2h / Phase3-2h / Phase4-1.5h / Phase5-1.5h
