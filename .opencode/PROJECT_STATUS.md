---
project: university-social
path: ~/.opencode-workspace/university-social
milestone: "生产就绪"
overall: 99
last_updated: 2026-05-22
---

## 项目概述

校园社交小程序（校园社），uni-app + Vue 3 + uView Pro 前端，Go/Gin + PostgreSQL 后端。

## 完成功能

- [x] 用户注册/登录（dev-login + JWT）
- [x] 首页（hero + 功能卡片导航）
- [x] 广场（帖子流、点赞、评论数）
- [x] 搜索性能优化（ILIKE → pg_trgm GIN 索引 + similarity 排序）
- [x] 帖子详情（评论区渲染、发表评论）
- [x] 圈子（列表、创建、详情、帖子、成员、搜索）
- [x] 好友（搜索、发送请求、收件箱/接收/拒绝、长按删除）
- [x] 个人中心 & 资料编辑（头像、昵称、学校、简介）
- [x] 消息通知列表
- [x] uView Pro UI 框架改造
- [x] 图片上传后端（/upload + /uploads/ static）
- [x] 后台管理面板（总览统计、用户/帖子/圈子/举报/敏感词管理）
- [x] 树洞系统（匿名发帖、随机代号池 50+、点赞、评论、敏感词过滤）
- [x] 举报系统（ReportPopup组件 + admin处理/驳回）
- [x] 敏感词系统（admin管理 + 内容过滤）
- [x] 校园认证（学籍邮箱验证码 + 认证状态 + 前端页面）
- [x] 关注系统（follow/unfollow + 关注动态流 + 双信息流）
- [x] 内容发现（热帖榜 trending + 搜索 + 关注流）
- [x] 邀请码系统（生成/领取 + 每人限5个）
- [x] 积分系统（签到 + 余额 + 排行榜）
- [x] 二手交易市场（发布/列表/分类/搜索/评价）
- [x] 校园活动系统（发布/报名/离开/留言/参与者）
- [x] 课程评价（搜索课程/评分/列表/有用标记）
- [x] 失物招领（发布/列表/搜索/认领状态）
- [x] WebSocket/IM 实时私信（gorilla/websocket + 对话/消息/未读）
- [x] Docker Compose 一键部署（PostgreSQL + Redis + Backend）
- [x] 20个数据库迁移文件
- [x] Banner/公告管理后台：Banner CRUD + 公告发布 + 精华帖标记
- [x] 圈子审核流：加入类型选择 + 待审批申请 + 圈主通过/拒绝
- [x] 成员头像墙 API + 前端展示
- [x] Banner轮播 + 公告栏 + 精华帖（后端 CRUD + 前端组件）
- [x] Redis 缓存：go-redis 初始化 + Banner/公告 5min TTL 缓存 + CRUD 自动失效
- [x] 前端 Banner/公告 API + 首页轮播组件 + 公告列表页
- [x] 广场精华帖 Tab + 圈子审核申请 UI
- [x] Banner/公告管理后台
- [x] ECharts 3D 看板（DAU/增量/占比图表）
- [x] 迁移 023_daily_stats.sql
- [x] Admin 构建验证通过
- [x] 圈子弹幕动态：migration 024 + API + 前端 Tab
- [x] 校园工具：校历/黄页/空教室（后端 API + 前端 3 页面）
- [x] ECharts 3D 管理后台看板（DAU/增量/占比饼图）
- [x] Redis 缓存集成（Banner/公告 5min TTL）
- [x] Handler 集成测试扩展（26 tests：auth/post/circle/friend）
- [x] 测试修复：GetUser handler 修复 nil→404
- [x] Cache 单元测试（JSON roundtrip）
- [x] 生产部署：Nginx 配置 + Dockerfile.web + deploy.sh
- [x] CI 扩展：3 job 并行（backend + frontend + admin）
- [x] Admin 敏感词管理页面 + 举报处理页面
- [x] GH Pages 部署配置（pages.yml + base 路径）

## 待办功能

- [ ] 部署到生产服务器（push main → GH Actions + render.com）
- [ ] e2e 测试（可选）

## 自主性优化 (2026-05-22)

参考 ~/sucai_0/08.txt，完成了以下自主性提升：

### 安全安装的插件 ✅
- **owflow@0.1.5** — 规格驱动开发工作流（需 Node ^25，已手动初始化 .owflow/）
- **@spoons-and-mirrors/subtask2@0.3.5** — 条件驱动任务链
- **@op1/workspace@0.5.9** — 跨会话持久计划管理
- **opencode-hive@1.4.8** — 蜂群编码（已初始化 .hive/）

### 冲突修复 🛠️
- **@zcy2nn/agent-forge@1.2.6 移除** — 经诊断发现其为 oh-my-opencode 的 fork（同一同源），
  两者同时运行会导致钩子冲突和功能异常。保留已深度配置的 oh-my-opencode。

### 最终插件架构
```
oh-my-opencode  (主编排) + owflow (规格驱动) + subtask2 (条件链)
        + workspace (跨会话) + hive (蜂群) + Ralph Loop (迭代)
```
| 组件 | 技术栈 |
|------|--------|
| **前端** | uni-app + Vue 3 + uView Pro + Vite |
| **后端** | Go/Gin + PostgreSQL + Redis |
| **认证** | JWT token (dev bypass WeChat) |
| **Admin** | Vue 3 + Element Plus + DataV (@kjgl77/datav-vue3) |
| **开发环境** | WSL2, Backend :8081, Vite :5173 |
| **部署** | Docker Compose + Nginx + CI (3-job parallel) |
| **自主工作流** | oh-my-opencode + owflow + subtask2 + workspace + hive |

## 已知问题

- 图片上传的绝对路径需在前端拼接 base URL
- Redis 已安装本地，但未用于业务缓存（仅 graphiti 插件热层）

## Ralph Loop 优化记录 (2026-05-22)

### 已完成的优化 ✅
|- DataV 装饰边框集成：Dashboard (dv-border-box-1, dv-digital-flop), App.vue (sidebar/header 边框), 6 个表格页 (dv-border-box-13)
|- Bug 修复：square/index.vue 重复 refreshing 行, admin/index.vue 作用域错误 (u→p) 和重复 onMounted 代码
|- Loading 状态增强：admin/index.vue Promise.all 并行加载, friend/index.vue 重复行移除+loading, chat/index.vue loading 状态, product/list.vue u-loading 组件化
|- 构建验证：frontend + admin 均通过，backend test 通过
|- 自主工作流：Node v25 升级, 4 插件安装, agent-forge 冲突移除
|- Admin SPA 登录页缺失 (低优先级, 开发模式 JWT 正常)

