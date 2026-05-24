# University Social — 新功能设计文档

> **设计系统**：Dark Terracotta Oklch（$brand-primary: oklch(0.62 0.12 16), $bg-app #1A1E2E, 1.25 type scale, 4px base）
 > **框架**：Vue 3 + Element Plus + Vue Router
 > **日期**：2026-05-24
 > **所有前端页面已精修完毕**（~25 页统一）, 现在进入新功能开发

---

## A — 评价系统 ★ 核心功能

### 1. 产品定位

**虎扑风格「万物皆可评」**——评价系统是大学社交平台的独立一等公民功能，
不是任何实体（课程/商品/活动）的附属模块。

用户可以自由创建评价对象（课程、教师、食堂菜品、宿舍楼、校园活动……），
对任何对象打分 + 标签 + 写评价 + 标记「有用」。

### 2. 路由设计

| 路由 | 页面 | 说明 |
|------|------|------|
| `/reviews` | **评价广场** | 搜索 + 分类筛选 + 卡片网格 + 最新评价流 |
| `/reviews/create` | **创建评价对象** | 名称/分类/描述表单, 查重 |
| `/reviews/:id` | **评价详情** | 虎扑风格评分概览 + 分布柱状图 + 评价列表 |
| 嵌入至 `course/:id` / `product/:id` | `<ReviewSection>` | 通用嵌入组件 |

### 3. 完整组件树

```
ReviewHub                          —— /reviews
├── ReviewSearchBar                —— 分类 tabs (课程/教师/商品/活动/生活/全部) + 搜索
├── ReviewTargetCardGrid           —— auto-fill minmax 260px, stagger-item
│   └── ReviewTargetCard           —— 名称 + 大号均分(1-10) + 评价数 + 分类标签
├── RatingsFeed                    —— 最新评价流 (compact cards)
│   └── RatingCard(self)           —— 用户头像 + 目标名称 + 10分制分数 + 简短摘要
└── CreateReviewTargetFAB          —— 浮动按钮 48px circle, brand gradient, 右下角固定

ReviewTargetDetail                 —— /reviews/:id
├── ScoreOverview                  —— 两栏
│   ├── 左: ScoreBadge             大号分数(3rem gold) + 评价总数
│   ├── 右: DistributionBarChart   5 段柱状分布 [1-2] [3-4] [5-6] [7-8] [9-10]
│   └── TargetMeta                创建者 + 分类 badge + 描述
├── RatingActions                  —— 评价/编辑浮层 trigger
├── RatingFilterBar                —— 热门 / 最新 / 高分 (el-radio-group)
└── ReviewList
    └── ReviewCard                 —— 8-color 匿名头像 + 10分徽标 + 正文 + 标签 pills
        └── HelpfulButton          —— 「有用」toggle, 点亮色 oklch(0.72 0.18 55)

CreateReviewTarget                 —— /reviews/create
├── NameField                      —— el-input, 带防抖查重(300ms)
├── CategorySelect                 —— el-select: 课程/教师/商品/活动/生活/其他
├── DescriptionField               —— el-input textarea, 可选, 200 字
└── SubmitButton                   —— loading 态 + 创建后跳转详情页

ReviewSection (通用嵌入组件)
├── ScoreOverview (compact)        —— 小号概览: 分数 + 分布条(迷你)
├── RatingCard (preview 3 条)      —— 预览
└── ViewAllLink                    —— 「查看全部 N 条评价」→ /reviews/:id
```

### 4. 评分交互流 (RatingDialog)

```
点击「写评价」按钮
  → el-dialog 弹窗 (50vw, max 600px)
    ├── ⭐ 10 分量表
    │  el-rate :max=10 show-score, 高亮色 gold oklch(0.72 0.18 55)
    │  hover 时实时显示
    ├── 📋 快捷标签
    │  max 3 个, el-checkbox-group, pill 样式
    │  按分类动态: 课程→{讲得清楚,给分高,水分大,收获多}
    │  生活→{实惠,拥挤,干净,方便}
    ├── ✍️ 详细评价
    │  el-input type=textarea, max 500 字, show-word-limit
    └── 提交按钮
        成功后: 按钮→「已评价 X 分」灰色, 可编辑
```

### 5. API 契约（后端需实现）

#### 评价对象 ReviewTarget

```
GET    /api/review-targets          ?page&limit&category&search&sort
→ { targets: ReviewTarget[], total: number }

POST   /api/review-targets
body: { name, category, description }
→ ReviewTarget

GET    /api/review-targets/:id
→ { target: ReviewTarget, stats: { avgScore, distribution[5], totalRatings } }

GET    /api/review-targets/search  ?q=xxx
→ { targets: ReviewTarget[] }   // 创建时查重
```

#### 评价 Rating

```
POST   /api/review-targets/:id/ratings
body: { score(1-10), tags(string[]), content(string) }
→ Rating

GET    /api/review-targets/:id/ratings  ?page&limit&sort(热门/最新/高分)
→ { ratings: Rating[], total: number }

DELETE /api/ratings/:id
// 仅作者可删

POST   /api/ratings/:id/helpful
// toggle 有用
→ { helpfulCount, isHelpfulByMe }
```

#### 数据模型

```typescript
interface ReviewTarget {
  id: number
  name: string
  category: 'course' | 'teacher' | 'product' | 'activity' | 'life' | 'other'
  description: string
  creatorId: number; creator: UserInfo
  avgScore: number; totalRatings: number
  distribution: number[]  // [1-2段, 3-4, 5-6, 7-8, 9-10]
  createdAt: string
}

interface Rating {
  id: number; targetId: number; userId: number; user: UserInfo
  score: number       // 1-10
  tags: string[]      // max 3
  content: string     // max 500 字
  helpfulCount: number; isHelpfulByMe: boolean
  createdAt: string
}
```

### 6. 实施阶段

| Phase | 内容 | 后端依赖 |
|-------|------|---------|
| **P1** | 课程评价升级 — 现有 course API + ReviewSection 组件 + RatingDialog | 无需新 API |
| **P2** | 评价广场 /reviews + 评价详情 /reviews/:id | 新 API review-targets CRUD |
| **P3** | 用户自由创建评价对象 /reviews/create | 新 API POST target |
| **P4** | 扩展到商品/活动详情页嵌入 | 新增各实体 ratings API |

### 7. 视觉 token

```scss
$review-star-color: oklch(0.72 0.18 55);   // 金色高亮
$review-score-large: 3rem;
$review-score-weight: 700;
$card-rating-bg: darken($bg-card, 2%);

.dist-bar {
  height: 8px; border-radius: 4px;
  background: linear-gradient(90deg, oklch(0.72 0.18 55), oklch(0.62 0.12 16));
}
.tag-pill { border-radius: 999px; padding: 2px 12px; font-size: 0.8rem; }
.helpful-active { color: oklch(0.72 0.18 55); }
.fab-review {
  position: fixed; bottom: 24px; right: 24px;
  width: 48px; height: 48px; border-radius: 50%;
  background: linear-gradient(135deg, $brand-primary, darken($brand-primary, 15%));
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}
```

### 8. 空状态 & 边界

| 场景 | 文案 |
|------|------|
| 评价广场无数据 | 🏟️「还没有评价对象，成为第一个发起评价的人」+ 创建按钮 |
| 对象无评价 | ✍️「还没有评价，来说两句？」+ 写评价弹窗 |
| 搜索无结果 | 🔍「没有找到匹配的评价对象」+ 创建按钮 |
| 网络错误 | 骨架屏 + retry 按钮 |
| 已评不可再评 | 按钮置灰「已评价 X 分」+ 可编辑入口 |

---

## B — 内容发现 Content Discovery

**目标**：Square 从简单时间线升级为内容发现系统。

### 四 Tab 结构

| Tab | API | 说明 |
|-----|-----|------|
| 最新 | `listPosts()` | 默认, 即时刷新 |
| 🔥 热门 | `trendingPosts()` | 24h 热度排序, fire 标签 |
| 关注 | `followingPosts()` | 已存在, 未登录→el-empty |
| ✨ 精选 | `listFeaturedPosts()` | 系统推荐 |

### 交互
- Tab 切换: `el-radio-group`, fade 0.3s transition
- 无限滚动: IntersectionObserver + usePagination
- 🔥 badge: PostCard 右上角, pulse opacity 1.5s infinite
- 每个 Tab 独立 skeleton/empty, 切换重置页码

---

## C — Banner & 公告

**目标**：利用现有 `listBanners()` / `listAnnouncements()` API。

### C1 首页轮播
- `el-carousel` 5s interval, arrow hover
- 渐变 overlay: `linear-gradient(to top, rgba(0,0,0,0.6), transparent)`
- 标题 + 副标题 + 行动按钮
- 数据空→隐藏组件

### C2 公告列表 `/announcements`
- PageHeader「校园公告」
- 置顶公告: el-tag type=warning
- 卡片 hover translateY+shadow（统一风格）

---

## D — 通知中心升级

**目标**：`notification/index.vue` 升级为多分类管理式通知中心。

### 四 Tab

| Tab | 过滤器 |
|-----|--------|
| 全部 | 无过滤 |
| 互动 | type ∈ {like, comment, follow} |
| 系统 | type === 'system' |
| 课程 | type === 'course' |

### 功能
- 「全部已读」固定在 tab 右侧
- 智能时间戳: 今天→HH:mm / 昨天→昨天 HH:mm / 本周→周X / 更早→MM-DD
- 每 Tab 独立空状态文案
- Tab 切换 badge 计数更新
