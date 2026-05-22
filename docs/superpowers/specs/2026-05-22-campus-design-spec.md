# 校园社 设计规格 (v1.0)

> 洛阳高校 · 校园社交平台 — UI/UX 终极精修设计规范

## 一、设计系统

### 1.1 色彩体系

**品牌色 (Primary):** `#C67A6A` — 暖褐红，唤起校园温暖感
- Light: `#E8A08A` (hover, disabled)
- Dark: `#A85A4A` (active, pressed)
- Subtle: `rgba(198,122,106,0.08)` (bg tint)

**深色 (Dark):** `#1E2A3A` — 深夜蓝灰
- Surface: `#2A3A4E`
- Border: `rgba(255,255,255,0.1)`

**浅色 (Light):**
- Canvas: `#F7F4F0` (页面背景)
- Surface: `#ffffff` (卡片背景)
- Hairline: `#EAE6E0` (分割线)

**语义色:**
- Ink: `#1E2A3A` (主文本)
- Ink-muted: `#8E9BAB` (次要文本)
- Ink-subtle: `#B8C2CE` (辅助文本)
- Ink-tertiary: `#D0D5DD` (占位符)
- Success: `#67C23A`
- Warning: `#E6A23C`
- Danger: `#F56C6C`

### 1.2 字体规范
- 中文: `PingFang SC`, `Noto Sans SC`, `-apple-system`, `sans-serif`
- 英文/数字: `SF Pro Text`, `Helvetica Neue`
- 字号体系: 20/22/24/26/28/30/32/36/40/44 (rpx)
- 字重体系: 400(regular) / 500(medium) / 600(semibold) / 700(bold)

### 1.3 间距体系
- xs: 8rpx
- sm: 12rpx
- md: 16rpx
- lg: 24rpx
- xl: 32rpx
- 2xl: 40rpx

### 1.4 圆角体系
- sm: 8rpx
- md: 12rpx
- lg: 16rpx
- xl: 24rpx
- full: 9999rpx

### 1.5 阴影体系
- sm: `0 2rpx 8rpx rgba(30,42,58,0.03)`
- md: `0 2rpx 12rpx rgba(30,42,58,0.04)`
- lg: `0 4rpx 20rpx rgba(30,42,58,0.06)`
- glow: `0 0 20rpx rgba(198,122,106,0.15)`

### 1.6 动画体系
- easeOutQuart: `cubic-bezier(0.25, 0.46, 0.45, 0.94)` — 入场
- easeOutBack: `cubic-bezier(0.34, 1.56, 0.64, 1)` — 弹性效果
- duration-fast: 200ms
- duration-normal: 300ms
- duration-slow: 400ms
- stagger-base: 80ms (notifications, lists)

## 二、前端页面规范

### 2.1 通用状态
每个列表页面必须实现三态：
1. **Loading**: `<u-loading mode="flower" />`
2. **Empty**: `<u-empty mode="..." text="..." />` + 描述 + action button
3. **Data**: 列表内容

### 2.2 页面过渡
- SwitchTab: 默认 (uni-app handles)
- NavigateTo: slide-in-right (200ms easeOutQuart)
- Modal/Popup: fade + scale (300ms easeOutBack)

### 2.3 微交互
- Feature card hover: scale(0.96) + shadow 过渡 200ms
- Button press: scale(0.95) + 背景色加深 150ms
- Like/Star: 图标 + 数字动画 (bounce 300ms)
- Pull refresh: uView 默认

## 三、Admin 控制台规范

### 3.1 布局
- FullScreenContainer 包裹
- dv-border-box-12 侧栏 | dv-border-box-7 头部 | dv-border-box-13 内容区
- 3D 大屏: 独立路由 `/datascreen`，不侵入普通管理操作

### 3.2 DataV 组件使用
- 数字: dv-digital-flop (关键指标)
- 地图: dv-flyline-chart-enhanced (洛阳/河南)
- 排行: dv-scroll-ranking-board (高校排行)
- 滚动: dv-scroll-board (实时动态)
- 环图: dv-active-ring-chart (内容分布)
- 胶囊: dv-capsule-chart (板块分布)
- 锥柱: dv-conical-column-chart (今日新增)
- 水球: dv-water-level-pond (负载)
- 百分比: dv-percent-pond (使用率)

## 四、反AI 设计规则

1. 避免无意义的占位文本 — 所有示例数据应有意义（洛阳高校相关）
2. 避免过度对称 — 不对称布局更自然
3. 避免纯色块 — 使用渐变、纹理、装饰元素
4. 颜色饱和度不宜过高 — 优先使用低饱和度、高级感配色
5. 每个空状态必须有 icon + 描述 + 操作引导（三件套）
6. 动画必须有用 — 不做纯粹的装饰动画，所有动画服务于理解
