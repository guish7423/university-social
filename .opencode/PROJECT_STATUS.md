---
project: university-social
path: ~/.opencode-workspace/university-social
milestone: "本地迭代完成 ✅"
overall: 88
last_updated: 2026-05-21
---

## 项目概述

校园社交小程序（校园社），uni-app + Vue 3 + uView Pro 前端，Go/Gin + PostgreSQL 后端。

## 完成功能

- [x] 用户注册/登录（dev-login + JWT）
- [x] 首页（hero + 功能卡片导航）
- [x] 广场（帖子流、点赞、评论数）
- [x] 搜索性能优化（ILIKE → pg_trgm GIN 索引 + similarity 排序）
- [x] 帖子详情（评论区渲染、发表评论）
- [x] 圈子（列表、创建、详情、帖子、成员）
- [x] 好友（搜索、发送请求、收件箱/接收/拒绝、长按删除）
- [x] 个人中心（头像、昵称、学校）
- [x] 个人资料编辑（头像上传、昵称、学校）
- [x] 消息通知列表
- [x] uView Pro UI 框架（6+页面改造）
- [x] 图片上传后端（/upload + /uploads/ static）
- [x] opencode-graphiti 图谱记忆插件（配置就绪，需重启生效）

## 待办功能

- [x] 搜索性能优化（PostgreSQL pg_trgm 全文检索）
- [ ] WebSocket/IM 实时消息
- [ ] DB schema 审查与索引优化
- [ ] 后台管理面板完善
- [ ] 自动化测试
- [ ] 部署到服务器

## 技术决策

- **前后端分离**: uni-app H5 + Gin REST API
- **认证**: JWT token（dev 环境 bypass WeChat）
- **存储**: PostgreSQL 本地，文件系统 uploads/
- **开发**: WSL2 环境，后端 :8081，Vite :5173
- **UI 框架**: uView Pro 0.6.2
- **跨会话记忆**: opencode-graphiti（Redis 本地热层 ✅, Graphiti MCP 待网络就绪）

## 已知问题

- 图片上传的绝对路径需在前端拼接 base URL
- Redis 已安装本地，但未用于业务缓存（仅 graphiti 插件热层）
- 搜索仅 LIMIT 10，无分页
- 无 WebSocket，消息靠轮询
