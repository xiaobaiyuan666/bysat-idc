# 无穷云 IDC 系统当前状态

更新时间：2026-03-23

## 当前系统状态

当前仓库已经完成一次大整理，只保留了后续继续开发真正需要的代码和文档。

当前运行入口：

- 后台：`http://localhost:5177`
- 用户端：`http://localhost:5178`
- API：`http://localhost:18080`

## 当前已经确认的事实

### 仓库层面

- 根目录已经压缩为“说明文件 + 有效文档 + active workspace”模式。
- 当前唯一业务工作区是 [`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)。
- 旧版后台、旧参考仓、legacy 大文档都已经从主线删除。

### 后台与 VNC

- 服务详情中的 VNC 能力已经迁移到新后台。
- 当前 VNC 页面入口在 [`idc-finance/apps/admin-web/src/views/vnc/index.vue`](/Users/a1/Documents/codex/bysat-idc/idc-finance/apps/admin-web/src/views/vnc/index.vue)。
- 当前 noVNC 使用本地静态 bundle，而不是直接依赖 npm 包。
- 这么做的原因是 npm 版 noVNC 在当前 Vite 产线构建中会触发打包问题。

### 魔方云链路

- 当前主线是“新后台 / 新用户端 / 新 API”对接魔方云，而不是再维护旧系统。
- 魔方云相关说明保留在根目录 [`docs/mofang-integration.md`](/Users/a1/Documents/codex/bysat-idc/docs/mofang-integration.md) 与 [`docs/mofang-portal-sync.md`](/Users/a1/Documents/codex/bysat-idc/docs/mofang-portal-sync.md)。
- 门户同步链当前重点是客户、门户账号、服务实例、资源明细、导入订单、通知、续费预览。

## 当前开发优先级

如果没有新的业务优先级覆盖，建议按下面顺序继续：

1. 继续补齐用户端云服务器相关真实业务链路，而不是只做静态页面。
2. 把魔方云实例动作、资源动作、同步结果在后台和用户端都做成可闭环的工作台。
3. 把门户服务详情页继续向真实运营场景推进，尤其是实例详情、资源明细、续费和状态展示。
4. 把关键业务动作产生的审计、同步日志和失败处理补完整。

## 当前必须遵守的技术决策

- 不要把旧后台或旧参考仓重新作为主线代码恢复进仓库。
- 不要把密钥、账号密码、机器路径写进 git 跟踪文件。
- 所有能影响后续接手的关键上下文，都要写进仓库文档，而不是只留在聊天记录里。
- 新增复杂业务逻辑时，注释必须说明“为什么这样做”，而不是只描述代码表面动作。

## 继续开发前建议先看的文件

1. [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
2. [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
3. [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
4. [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md)
5. 你本次任务对应的具体模块代码
