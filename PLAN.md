# 大学社交平台 — 项目计划

## 项目定位
面向大学生的轻量社交平台，微信小程序为主端，Web 后台管理为辅。

## 技术栈

| 层 | 技术选型 | 参考项目 |
|---|---------|---------|
| 前端框架 | uni-app 3.x + Vue 3 + TypeScript | socialuni, beaver-mobile, uni-vue3 |
| 状态管理 | Pinia | uni-vue3 template |
| UI 组件 | Wot Design Uni / uView Plus | uniapp-vite-template |
| 样式方案 | UnoCSS (原子化CSS) | uni-vue3 |
| 后端框架 | Go + Gin/Go-Zero | beaver (go-zero) |
| 数据库 | PostgreSQL + Redis | — |
| IM | 腾讯云IM SDK / 自建WebSocket | 参考腾讯云IM文档 |
| 部署 | Docker + Docker Compose | — |

## 核心功能

### 小程序端（Phase A）
1. 用户系统 — 微信登录 + 手机号绑定
2. 校园认证 — 学籍邮箱/学生证认证
3. 动态广场 — 发帖、评论、点赞、话题
4. 好友系统 — 搜索、添加、私信
5. 圈子/社团 — 创建加入、圈内动态
6. 个人主页 — 资料、动态、好友列表

### 管理后台（Phase B）
1. 用户管理 — 审核、封禁、角色
2. 内容管理 — 帖子审核、违规处理
3. 圈子管理 — 创建审核、管理员任命
4. 数据统计 — DAU、帖量、注册趋势

## 架构设计

```
wechat-mini-program (uni-app)
        ↓ HTTPS/WebSocket
   Go API Gateway
        ↓
  Service Layer (user/post/circle/im)
        ↓
   PostgreSQL + Redis
```

## 里程碑

| 里程碑 | 内容 | 预计步数 |
|-------|------|---------|
| M1 | 项目脚手架 + 用户认证系统 | 10 |
| M2 | 校园认证 + 动态广场 | 8 |
| M3 | 好友系统 + IM 通讯 | 8 |
| M4 | 圈子/社团功能 | 6 |
| M5 | 管理后台 | 6 |
| M6 | 部署上线 + 测试 | 4 |

## 完成标准
- 微信小程序可正常扫码打开、登录、发帖、加好友、聊天
- 管理后台可审核用户和帖子
- Docker 一键部署
