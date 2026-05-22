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

### 安装的自主性插件
- **owflow@0.1.5** — 规格驱动开发工作流（propose→spec→plan→implement→verify→archive）
- **@zcy2nn/agent-forge@1.2.6** — 多智能体编排（Orchestrator/Workflow/Researcher）
- **@spoons-and-mirrors/subtask2@0.3.5** — 条件驱动任务链（自动重复执行+自检）
- **@op1/workspace@0.5.9** — 跨会话持久计划管理（plan+notepad）
- **opencode-hive@1.4.8** — 蜂群编码（大规模并行任务）

### 更新配置
- AGENTS.md 重写为自主开发工作流
- opencode.json 新增5个插件声明
- package.json 新增5个依赖并 npm install
## 技术决策
- **前后端分离**: uni-app H5 + Gin REST API
- **认证**: JWT token（dev 环境 bypass WeChat）
- **存储**: PostgreSQL 本地，文件系统 uploads/
- **开发**: WSL2 环境，后端 :8081，Vite :5173
- **UI 框架**: uView Pro 0.6.2
- **树洞匿名策略**: pickCodename(userID*31 + whisperID*7) % poolSize 保证帖内一致、帖间变化

- **前后端分离**: uni-app H5 + Gin REST API
- **认证**: JWT token（dev 环境 bypass WeChat）
- **存储**: PostgreSQL 本地，文件系统 uploads/
- **开发**: WSL2 环境，后端 :8081，Vite :5173
- **UI 框架**: uView Pro 0.6.2
- **自主工作流**: owflow + agent-forge + subtask2 + workspace + hive（参考 08.txt）

## 已知问题

- 图片上传的绝对路径需在前端拼接 base URL
- Redis 已安装本地，但未用于业务缓存（仅 graphiti 插件热层）
