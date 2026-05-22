# Project: university-social (校园社)

## 自主开发工作流 (Autonomous Development)

本项目已配置全栈自主性插件链。智能体应按以下优先级自动选择工作流：

### 指令优先级
1. 用户明确指令 (最高)
2. 本自主工作流规则
3. 全局 AGENTS.md / 系统提示

### 核心自主性能力

#### 🏭 owflow — 规格驱动开发
- 新功能必须从规范开始，而非从编码开始
- 流程: `propose → spec → plan → implement → verify → archive`
- 使用 `/proposal` 提变更方案
- 使用 `openspec_ff` 生成脚手架
- 使用 `openspec_apply + Ralph Loop` 批量执行任务

#### 🔄 @zcy2nn/agent-forge — 多智能体编排
- 复杂任务（4步以上）自动派发给专业智能体
- 内置 Agent 角色：Orchestrator（项目经理）、Workflow（流程专家）、Researcher（研究员）
- 使用 `@agent-name <task>` 指派任务

#### 🧵 @spoons-and-mirrors/subtask2 — 条件驱动任务链
- 为每个任务设定明确的"完成条件"
- 智能体自动重复执行并自检，直到满足条件
- 用于质量把关和自动化回归

#### 📋 @op1/workspace — 跨会话持久计划
- 使用 `plan` 和 `notepad` 维护跨会话的持久上下文
- 新会话启动时自动恢复上一个会话的计划状态
- 避免"失忆"导致的重复工作

#### 🐝 opencode-hive — 蜂群编码
- 大规模并行任务使用 Hive 模式
- 智能体集群分片处理，主智能体编排
- 适合独立模块的批量开发

### 自动工作流选择

| 任务类型 | 工作流 | 触发条件 |
|---------|--------|---------|
| 新功能开发 | owflow: propose → spec → plan → implement | 功能请求 |
| Bug 修复 | subtask2 条件驱动循环 | 错误报告 |
| 批量任务 (>3项) | agent-forge 多智能体派发 | 复杂需求 |
| 跨会话工作 | @op1/workspace plan 恢复 | 新会话 |
| 大规模并行 | opencode-hive 蜂群 | 大量独立模块 |
| 任何任务 | 先检视是否适用上述工作流 | 任何请求 |

### 完成标准 (Definition of Done)
1. 所有计划任务执行完成
2. 验证证据齐全（bash 输出 / lint 通过 / 构建通过）
3. 回归测试通过
4. 更新 .opencode/PROJECT_STATUS.md
5. 记录 ELF 学习经验 (elf add-rule)

## Session 启动
1. 读取 `.opencode/PROJECT_STATUS.md`
2. `@op1/workspace` 恢复计划状态
3. `elf search query="university-social"` 加载历史经验
4. `opencode-mem search semantic query="university social app"` 加载持久记忆

## Session 结束
1. 更新 `.opencode/PROJECT_STATUS.md`
2. `@op1/workspace` 保存计划状态
3. `elf add-rule type=learning` 保存学习经验
4. 重复 >2 次的模式 → `type=golden_rule` 固化
