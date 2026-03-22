# 无穷云 IDC 系统当前交接面板

更新时间：2026-03-23

这份文件是跨端续开发的第一交接面板。

任何客户端开始工作前，都应该先看这里，再决定读哪些模块代码。

## 当前主线

- 仓库状态：私有仓库，当前用于跨设备续开发与存放主线资料
- active workspace：[`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)
- 当前主线方向：继续完善新后台、新用户端、新 API 的真实业务闭环，不再恢复旧系统

## 当前已完成的基础建设

- 根目录已经整理为单一主线仓库
- 续开发入口文档已经建立
- 可提交环境模板已经建立
- MySQL 初始化说明已清洗为通用版本
- MySQL 演示 seed 已扩成完整工作台链路
- 仓库级 AI 协作规则已经建立
- 跨端协作协议已经建立

## 当前默认恢复路径

1. 阅读 [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
2. 阅读 [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)
3. 阅读 [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
4. 阅读 [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)
5. 阅读 [`docs/cross-device-collaboration.md`](/Users/a1/Documents/codex/bysat-idc/docs/cross-device-collaboration.md)
6. 再进入 [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md)

## 当前推荐启动方式

### 快速界面开发

- 使用 `memory` 模式
- 适合前端、页面、普通接口联调

### 真实链路联调

- 使用 `mysql` 模式
- 执行 `npm run db:prepare:mysql`
- 适合报表、Provider 同步、资源明细、服务链路联调

## 当前演示账号

- 后台：`admin / Admin123!`
- 用户端：`portal / Portal123!`

## 当前优先级

1. 继续补齐用户端云服务器真实业务链路，而不是只做静态展示
2. 把魔方云实例动作、资源动作、同步结果做成后台和用户端都能闭环
3. 推进门户服务详情页，包括实例详情、资源明细、续费、状态展示
4. 继续补齐审计、自动化任务、失败处理和交付后的可追踪信息

## 本次交接时已知状态

- 续开发规范已经正式写入仓库
- 新设备启动入口已明确
- 真实敏感信息仍然只能保存在本地 `.env.local`
- MySQL seed 目前是演示链路用数据，不代表真实生产数据结构已经全部打通

## 当前风险与注意事项

- 如果远端客户端开发后没有 `git push`，其他端默认接不到这次工作
- 如果改了环境变量但没同步 `.env.example`，新设备会恢复失败
- 如果改了迁移或 seed 但没写验证状态，后续接手容易误判数据库可用性

## 结束当前开发前必须补的内容

- 如果本次改动改变了主线优先级，更新本文件
- 如果本次改动改变了环境要求，更新本文件和环境文档
- 如果本次改动改变了外部对接行为，更新相关说明文档
- 如果本次改动没有完成验证，在本文件里写清楚未验证项

## 下一步推荐

如果没有新的用户优先级覆盖，下一步建议继续做：

- 用户端云服务器详情页与资源明细真实联调
- 用户端服务动作和后台 Provider 动作结果联动
- 服务详情、同步日志、自动化任务之间的可追踪闭环
