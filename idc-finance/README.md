# 无穷云 IDC 系统 Workspace

基于 `Go + Vue 3 + Element Plus + vue-pure-admin` 的无穷云 IDC 财务与云业务系统工作区。

如果你是中断后重新接手这个项目，先回到根目录阅读：

- `README.md`
- `docs/project-continuity.md`
- `docs/project-map.md`
- `docs/current-state.md`
- `docs/development-rules.md`

## 目录

- `apps/api`：Go API 服务
- `apps/admin-web`：后台管理端
- `apps/portal-web`：用户中心
- `internal`：后端模块化业务代码
- `migrations/mysql`：MySQL 迁移脚本
- `seed/mysql`：MySQL 演示数据
- `docs`：拆解、重建和执行计划文档

## 当前阶段

当前已推进到 `Phase 8`，已完成：

- Phase 1：认证、RBAC、菜单、审计日志、客户域、联系人、实名认证工作台
- Phase 2：商品、订单、账单、支付、服务最小闭环
- Phase 3：商品详情、订单详情、账单详情、后台登记收款
- Phase 4：服务工作台、服务动作、账单退款与状态回写
- Phase 5：价格矩阵、配置项、资源模板、支付记录、服务资源快照、实例动作
- Phase 5.5：门户下单支持配置项选择和价格预览
- Phase 6：MySQL 仓储入口、商品/订单/账单/服务/客户 MySQL 仓储、演示种子
- Phase 7：审计日志 MySQL 持久化、MySQL 迁移与种子执行器
- Phase 8：魔方云真实 Provider、实例拉取同步、资源落库、同步日志

## 启动

如果从仓库根目录启动：

```bash
npm run dev:api
npm run dev:admin
npm run dev:portal
```

### API

macOS / Linux：

```bash
npm run dev:api:unix
```

或直接：

```bash
go run ./apps/api/cmd/server
```

Windows：

```bash
& "C:\Program Files\Go\bin\go.exe" run ./apps/api/cmd/server
```

或直接使用本地 MySQL 启动脚本：

```bash
npm run dev:api:mysql
```

macOS / Linux 的 MySQL 模式：

```bash
npm run dev:api:mysql:unix
```

健康检查：

- `http://127.0.0.1:18080/api/v1/health`

### 后台

```bash
npm install
npm run dev:admin
```

访问地址：

- `http://localhost:5177`

默认账号：

- `admin / Admin123!`

### 用户中心

```bash
npm install
npm run dev:portal
```

访问地址：

- `http://localhost:5178`

默认账号：

- `portal / Portal123!`

## 环境变量

核心配置：

- `STORAGE_DRIVER=memory`
- `STORAGE_DRIVER=mysql`
- `STORAGE_STRICT=true`
- `MYSQL_DSN=idc_finance:password@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4`

魔方云配置：

- `MOFANG_CLOUD_BASE_URL`
- `MOFANG_CLOUD_USERNAME`
- `MOFANG_CLOUD_PASSWORD`
- `MOFANG_CLOUD_LANG=zh-cn`
- `MOFANG_CLOUD_INSECURE_SKIP_VERIFY=true`
- `MOFANG_CLOUD_LIST_PATH=/v1/clouds`
- `MOFANG_CLOUD_INSTANCE_DETAIL_PATH=/v1/clouds/:id`

当 `STORAGE_DRIVER=mysql` 且 `MYSQL_DSN` 可用时：

- 客户、商品、订单、账单、服务走 MySQL
- 审计日志走 MySQL
- 魔方云同步会把实例、资源明细和同步日志写入 MySQL

## MySQL 初始化

先创建数据库：

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

执行迁移与种子：

```bash
npm run db:migrate:mysql
npm run db:seed:mysql
```

macOS / Linux：

```bash
npm run db:migrate:mysql:unix
npm run db:seed:mysql:unix
```

一次完成初始化：

```bash
npm run db:prepare:mysql
```

macOS / Linux：

```bash
npm run db:prepare:mysql:unix
```

本机 MySQL 安装说明见：

- `docs/mysql-setup.md`

## 当前可用范围

### 后台

- 登录
- 动态菜单
- 工作台
- 客户列表与客户详情工作台
- 联系人新增、编辑、删除
- 实名审核
- 审计日志
- 商品列表与商品详情工作台
- 价格矩阵、配置项、资源模板编辑
- 订单列表与订单详情工作台
- 账单列表与账单详情工作台
- 线下收款登记
- 支付记录与退款记录查看
- 服务列表与服务详情工作台
- 服务恢复、暂停、终止、重启、重装、重置密码
- 服务详情直接打开 VNC 控制台
- 魔方云健康检查
- 魔方云实例列表、实例详情、远程动作
- 魔方云实例拉取同步
- 服务资源明细查看
- Provider 同步日志查看

### 用户中心

- 控制台概览
- 商品商城
- 配置项选择与价格预览
- 下单生成订单和账单
- 账单支付
- 服务列表
- 订单列表
- 工单概览
- 钱包概览
- 账户资料概览

### API

后台：

- `GET /api/v1/health`
- `POST /api/v1/admin/auth/login`
- `GET /api/v1/admin/menus`
- `GET /api/v1/admin/permissions`
- `GET /api/v1/admin/audit-logs`
- `GET /api/v1/admin/customer-groups`
- `GET /api/v1/admin/customer-levels`
- `GET /api/v1/admin/customer-identities`
- `GET /api/v1/admin/customers`
- `POST /api/v1/admin/customers`
- `GET /api/v1/admin/customers/:id`
- `PATCH /api/v1/admin/customers/:id`
- `GET /api/v1/admin/customers/:id/contacts`
- `POST /api/v1/admin/customers/:id/contacts`
- `PUT /api/v1/admin/customers/:id/contacts/:contactId`
- `DELETE /api/v1/admin/customers/:id/contacts/:contactId`
- `GET /api/v1/admin/customers/:id/identities`
- `POST /api/v1/admin/customers/:id/identities/:identityId/review`
- `GET /api/v1/admin/customers/:id/services`
- `GET /api/v1/admin/customers/:id/invoices`
- `GET /api/v1/admin/customers/:id/tickets`
- `GET /api/v1/admin/customers/:id/audit-logs`
- `GET /api/v1/admin/products`
- `GET /api/v1/admin/products/:id`
- `POST /api/v1/admin/products`
- `PATCH /api/v1/admin/products/:id`
- `GET /api/v1/admin/orders`
- `GET /api/v1/admin/orders/:id`
- `PATCH /api/v1/admin/orders/:id`
- `GET /api/v1/admin/invoices`
- `GET /api/v1/admin/invoices/:id`
- `PATCH /api/v1/admin/invoices/:id`
- `POST /api/v1/admin/invoices/:id/receive-payment`
- `POST /api/v1/admin/invoices/:id/refund`
- `GET /api/v1/admin/services`
- `GET /api/v1/admin/services/:id`
- `POST /api/v1/admin/services/:id/actions/:action`
- `GET /api/v1/admin/providers/mofang/health`
- `GET /api/v1/admin/providers/mofang/instances`
- `GET /api/v1/admin/providers/mofang/instances/:id`
- `POST /api/v1/admin/providers/mofang/instances/:id/actions/:action`
- `POST /api/v1/admin/providers/mofang/sync`
- `POST /api/v1/admin/providers/mofang/services/:id/sync`
- `GET /api/v1/admin/providers/mofang/services/:id/resources`
- `GET /api/v1/admin/providers/mofang/sync-logs`

门户：

- `GET /api/v1/portal/dashboard`
- `GET /api/v1/portal/products`
- `GET /api/v1/portal/orders`
- `GET /api/v1/portal/invoices`
- `GET /api/v1/portal/services`
- `GET /api/v1/portal/tickets`
- `GET /api/v1/portal/account`
- `POST /api/v1/portal/orders/checkout`
- `POST /api/v1/portal/invoices/:id/pay`

## 当前验证

已通过：

- `& "C:\Program Files\Go\bin\go.exe" build ./...`
- `npm run build:admin`
- `npm run build:portal`
- `npm run db:migrate:mysql`

已做过的主链验证：

- 门户下单生成 `ORD-00000003 / INV-00000003`
- 门户按配置项下单后金额按增价计算到 `629`
- 门户支付账单后激活 `SRV-00000002`
- 后台登记线下收款生成 `PAY-00000002`
- 后台服务详情可执行重启、重置密码、重装
- 魔方云真实环境实例 `585 / codexsync032001` 已成功同步到本地
- 同步后本地生成 `SRV-MF-585`
- 资源明细已落库：网卡 `1`、IP `1`、磁盘 `1`
- Provider 同步日志已成功写入 MySQL
