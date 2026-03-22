# 无穷云 IDC 系统续开发说明

这份文档的目标只有一个：让任何人在离开当前编辑器、离开当前 AI 对话、甚至换一台机器之后，仍然可以基于仓库内容继续开发。

## 结论

这件事不需要依赖某个编辑器的专有功能，也不需要 `git submodule`。

正确做法是把“可续开发上下文”直接作为普通文本文件提交到仓库：

- 项目结构说明
- 当前进度与最近决策
- 开发规则与注释标准
- 外部接口边界与对接现状
- 下一阶段待办

这些文件跟代码一起走 git 版本管理，才具备跨编辑器、跨机器、跨会话的连续性。

## 中断后恢复开发的标准流程

1. 拉取最新代码：`git pull origin main`
2. 阅读根目录 [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
3. 阅读 [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
4. 阅读 [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
5. 阅读 [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)
6. 阅读 [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
7. 再进入 [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md) 看运行和模块说明
8. 最后根据本次任务，再去读对应模块代码和相关文档

## 当前约定

- 根目录只保留当前还会继续使用的主线工程和文档。
- 当前唯一 active workspace 是 [`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)。
- 与魔方云相关的补充资料保留在 [`docs`](/Users/a1/Documents/codex/bysat-idc/docs)。
- 当前仓库保持私有，作为未完工阶段的主线开发和续开发存储。
- 运行相关的本地敏感配置放在未提交的 `.env.local`，不要写进仓库文档。
- `.env.local` 的来源是 [`idc-finance/.env.example`](/Users/a1/Documents/codex/bysat-idc/idc-finance/.env.example)，模板可以提交，真实配置不能提交。

## 每次开发结束前必须补哪些内容

如果本次改动涉及下列任一项，就必须同时更新文档：

- 架构边界变化
- 外部接口接入变化
- 重要业务流程变化
- 当前阶段目标变化
- 下一步开发顺序变化
- 兼容性或构建限制变化

最少要更新：

- [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
- 必要时更新 [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
- 必要时更新 [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)

## 这套机制解决的不是“备份代码”，而是“备份上下文”

代码本身只能说明“现在长什么样”，但不能完整说明：

- 为什么这么做
- 还有哪些没做完
- 哪些旧方案已经废弃
- 哪些目录不能再作为参考
- 哪些接口现在是真实可用的

所以续开发能力必须依赖“代码 + 文档 + git 提交记录”三者同时存在。
