# 无穷云 IDC 系统跨端协作与续开发协议

这份协议适用于下面所有开发入口：

- 本地编辑器
- Codex desktop
- Codex web
- 远程开发机 / 云端工作区
- 未来新增的 AI 协作者或人工协作者

目标只有一个：

- 不依赖某个聊天窗口
- 不依赖某台机器
- 不依赖某个编辑器
- 只靠仓库内容就能继续开发

## 一、核心原则

### 1. 仓库才是权威上下文

以下内容才算“可续开发状态”：

- 已提交并推送的代码
- 已提交并推送的文档
- 已提交并推送的环境模板
- 已提交并推送的迁移、seed 和脚本

以下内容都不算可靠上下文：

- 只存在于聊天里的说明
- 只存在于本地未提交改动里的状态
- 只存在于编辑器记忆、插件记忆、私有笔记里的结论

### 2. 所有端都走同一恢复路径

无论从哪个端进入仓库，开始开发前都必须按这个顺序恢复：

1. `git pull origin main`
2. 阅读 [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
3. 阅读 [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)
4. 阅读 [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
5. 阅读 [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
6. 阅读 [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)
7. 阅读 [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
8. 阅读 [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)

## 二、开始开发前必须做的事

### 1. 检查仓库状态

至少确认：

- 当前分支是什么
- 当前工作区是否干净
- 远端是否已经拉到最新
- 当前唯一 active workspace 仍然是 [`idc-finance`](/Users/a1/Documents/codex/bysat-idc/idc-finance)

### 2. 检查当前交接面板

开始动手前，必须查看 [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)，确认：

- 当前主线目标
- 最近已完成内容
- 当前阻塞项
- 下一步推荐任务
- 本次改动需要同步哪些文档

### 3. 检查本地环境是否符合仓库要求

重点核对：

- [`idc-finance/.env.example`](/Users/a1/Documents/codex/bysat-idc/idc-finance/.env.example)
- [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)
- [`idc-finance/docs/mysql-setup.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/docs/mysql-setup.md)

## 三、开发过程中必须遵守的规则

### 1. 代码注释规则

复杂业务逻辑必须写“为什么这样做”，尤其是：

- 状态流转
- 第三方接口字段解释
- 兼容旧行为
- 降级逻辑
- 临时绕行逻辑
- 联调假设和风险边界

### 2. 文档同步规则

只要出现下面任一变化，就不能只改代码：

- 架构边界变化
- 接口行为变化
- 数据库结构变化
- seed 演示数据结构变化
- 环境变量变化
- 启动方式变化
- 当前主线优先级变化
- 废弃旧目录或旧方案

至少要同步更新相关文档，通常是这些文件之一或多个：

- [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
- [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
- [`docs/environment-setup.md`](/Users/a1/Documents/codex/bysat-idc/docs/environment-setup.md)
- [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)
- [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md)
- [`idc-finance/.env.example`](/Users/a1/Documents/codex/bysat-idc/idc-finance/.env.example)

### 3. 敏感信息规则

真实敏感信息只能放在本地未提交文件，例如：

- `.env.local`
- 本地数据库配置
- 私有测试账号
- 真实第三方密钥

禁止提交：

- 密码
- Token
- 机器路径
- 真实业务后台地址加凭据组合
- 编辑器私有状态文件

### 4. 数据库与演示数据规则

如果修改了以下内容：

- `migrations/mysql`
- `seed/mysql`
- MySQL 仓储加载逻辑
- 报表依赖的统计字段

则必须同时说明：

- 是否需要新迁移
- 旧 seed 是否还兼容
- 新设备初始化是否还能得到完整演示链
- 哪些内容已经验证，哪些还没验证

## 四、结束开发前必须做的事

### 1. 最低交接要求

每次结束一个可识别的开发阶段前，必须保证仓库能回答这些问题：

- 当前系统开发到哪一步了
- 当前主线优先级是什么
- 新设备怎么启动
- 当前环境变量有哪些
- 哪些能力已经可演示
- 哪些能力还没打通
- 下一步先做什么
- 当前是否存在未验证风险

如果仓库回答不了这些问题，就说明不能算“可续开发”。

### 2. 必须更新交接面板

只要本次工作改变了以下任一内容，就必须更新 [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)：

- 当前焦点任务
- 当前阻塞项
- 推荐下一步
- 环境要求
- 验证状态
- 需要注意的风险

### 3. 必须提交并推送

跨设备续开发的最低完成标准是：

1. `git add`
2. `git commit`
3. `git push`

如果没有 push：

- 这次工作不能视为已同步
- 其他设备或其他客户端默认看不到
- 下一个协作者不应该假定这些内容存在

## 五、当前推荐的交接载体

### 1. 固定入口

- [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
- [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)
- [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)

### 2. 固定规则

- [`AGENTS.md`](/Users/a1/Documents/codex/bysat-idc/AGENTS.md)
- [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
- [`docs/cross-device-collaboration.md`](/Users/a1/Documents/codex/bysat-idc/docs/cross-device-collaboration.md)

### 3. 固定交接面板

- [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md)

### 4. 固定模板

- [`docs/handoff-template.md`](/Users/a1/Documents/codex/bysat-idc/docs/handoff-template.md)

## 六、推荐工作流

### 1. 新端开始开发

```bash
git pull origin main
```

然后按恢复顺序阅读文档，再开始编码。

### 2. 本次开发完成

```bash
git status
git add <files>
git commit -m "<clear message>"
git push origin HEAD
```

然后确认 [`docs/active-handoff.md`](/Users/a1/Documents/codex/bysat-idc/docs/active-handoff.md) 已经和本次结果一致。

## 七、禁止事项

- 不要把续开发上下文只留在聊天里
- 不要把真实密钥写进仓库
- 不要让文档和代码长期脱节
- 不要在远端开发后不推送就默认别人能接上
- 不要恢复已经废弃的旧后台、旧参考仓作为当前主线
- 不要让“哪个 AI 记住了什么”替代“仓库里写了什么”
