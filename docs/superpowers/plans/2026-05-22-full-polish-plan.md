# 全方面精修 实现计划

> **面向 AI 代理的工作者：** 使用 subagent-driven-development 或 spawn_agent 并发执行任务。步骤使用复选框进度跟踪。

**目标：** 对 campus-social 项目进行全方面 UI/UX 精修、动画增强、功能完善、反AI 设计

**架构：** 5 路并行 agent — 前端页面精修、Admin 大屏增强、动画系统、空状态完善、构建验证

**技术栈：** uni-app (Vue 3), uView, DataV, Element Plus, SCSS

---

### Task 1: 空状态三件套改造

**文件范围：** 36 个页面的空状态

- [ ] 全部 `暂无xxx` 纯文本 → `<u-empty>` 组件 + 描述 + action
- [ ] 优先级：square(3种) > found/list > circle/list > notification > course/search > chat/index

### Task 2: 动画系统增强

**文件：**
- 修改：所有页面引入 page-level 过渡
- 修改：各页面微交互

- [ ] Loading → Data 过渡动画 (fadeInUp 300ms)
- [ ] 列表条目 stagger 动画 (index.vue, 延时递增)
- [ ] 空状态渐入动画
- [ ] Button/like 点击反馈 scale bounce

### Task 3: Admin 大屏数据接入

**文件：**
- 创建：admin/src/api/dashboard.ts
- 修改：admin/src/views/DataScreen.vue

- [ ] 创建 API 层获取真实统计数据
- [ ] 替换 DataScreen.vue 的 mock 数据为 API 调用
- [ ] 添加自动刷新 (30s interval)
- [ ] 构建验证

### Task 4: 反AI 设计注入

**文件范围：** 全部页面

- [ ] 检查占位文本 → 替换为洛阳高校相关示例
- [ ] 添加 noise/grain 纹理 overlay
- [ ] 不对称布局元素插入
- [ ] 低饱和度配色一致性校验

### Task 5: 构建验证

- [ ] Frontend build
- [ ] Admin build
- [ ] 后端 test
- [ ] 修复发现的问题
