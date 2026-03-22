# 无穷云 IDC 系统环境与续开发落地说明

这份文档解决两个问题：

- 新设备第一次拉代码后，如何把项目跑起来
- 离开当前编辑器或当前 AI 会话后，如何继续开发

## 当前仓库策略

- 当前 Git 仓库保持私有。
- 目的不是“藏代码”，而是在系统尚未开发完成前，用仓库存放主线代码、关键文档和续开发上下文。
- 等系统达到可公开交付状态后，再决定是否重新开放。

## 新设备接手顺序

1. 克隆仓库并拉取主线：`git clone <repo-url>`，然后 `git pull origin main`
2. 阅读根目录 [`README.md`](/Users/a1/Documents/codex/bysat-idc/README.md)
3. 阅读 [`docs/project-continuity.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-continuity.md)
4. 阅读 [`docs/project-map.md`](/Users/a1/Documents/codex/bysat-idc/docs/project-map.md)
5. 阅读 [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
6. 阅读 [`docs/development-rules.md`](/Users/a1/Documents/codex/bysat-idc/docs/development-rules.md)
7. 阅读 [`idc-finance/README.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/README.md)

前 6 步解决“项目现在是什么状态”，第 7 步解决“这台机器怎么跑起来”。

## 环境基线

当前仓库明确依赖：

- Go `1.26.1`，定义在 [`idc-finance/go.mod`](/Users/a1/Documents/codex/bysat-idc/idc-finance/go.mod)
- Node.js：仓库没有锁死版本，建议使用当前 LTS 版本
- npm：跟随当前 Node.js LTS
- MySQL `8.x`，用于持久化演示和后续真实联调

## 本地环境文件

可提交模板文件：

- [`idc-finance/.env.example`](/Users/a1/Documents/codex/bysat-idc/idc-finance/.env.example)

本地真实环境文件：

- `idc-finance/.env.local`

标准做法：

```bash
cp idc-finance/.env.example idc-finance/.env.local
```

然后只修改 `.env.local`。真实账号、密码、外部地址只允许放在 `.env.local`，不要写进仓库文档和代码。

## 快速启动

### 纯演示模式

适合看界面、继续做页面和接口联调，不依赖 MySQL：

```bash
npm install --prefix idc-finance
npm run dev:api
npm run dev:admin
npm run dev:portal
```

此时建议 `.env.local` 使用：

```env
STORAGE_DRIVER=memory
STORAGE_STRICT=false
```

### MySQL 持久化模式

适合继续开发真实业务链路、报表、Provider 同步和资源明细：

1. 先准备 MySQL 数据库
2. 把 `STORAGE_DRIVER` 改成 `mysql`
3. 在 `.env.local` 写入真实 `MYSQL_DSN`
4. 执行迁移和演示数据

```bash
npm run db:prepare:mysql
npm run dev:api:mysql
```

MySQL 详细初始化见 [`idc-finance/docs/mysql-setup.md`](/Users/a1/Documents/codex/bysat-idc/idc-finance/docs/mysql-setup.md)。

## 默认演示账号

- 后台：`admin / Admin123!`
- 用户端：`portal / Portal123!`

这两个账号仅用于当前演示和开发联调，不是生产认证方案。

## MySQL 演示数据范围

当前演示 seed 已覆盖：

- 商品、价格矩阵、配置项、资源模板
- 客户、联系人、实名状态
- 工单
- 订单、账单、支付、服务
- Provider 账号
- Provider 同步日志
- 服务资源明细
- 自动化任务
- 改配单演示链路

这意味着新设备拉代码后，只要把 MySQL 初始化跑完，就能得到一套可直接演示的工作台数据。

## 继续开发的最小动作

每次开发结束前，至少完成这几件事：

1. 更新关键代码注释，说明“为什么这样做”
2. 如果架构边界或接入状态变化，更新 [`docs/current-state.md`](/Users/a1/Documents/codex/bysat-idc/docs/current-state.md)
3. 把变更提交到 git，并推送到当前私有仓库

这样你换设备、换编辑器、换 AI 对话后，才能直接靠仓库继续接手，而不是回忆聊天记录。
