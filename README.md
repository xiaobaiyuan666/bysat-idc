# IDC 财务业务管理系统

一套面向中小型 IDC 运营场景的财务与业务管理系统，覆盖客户、产品、订单、服务、账单、收款、工单、资源与审计能力，并预留魔方云 `mf_cloud` 风格适配层。

## 技术栈

- `backend`：Next.js 16 + Prisma + SQLite
- `admin-ui`：Vue 3 + Element Plus
- `provider adapter`：可切换 mock / 真实接口的魔方云适配层

## 当前已具备能力

- 客户、产品、订单、服务、账单、收款、工单全链路管理
- 客户余额流水、手工余额调整、支付回调日志
- 退款记录、退款回冲与余额退款
- 计费引擎设置、续费账单生成、逾期标记、自动暂停、自动终止
- RBAC 菜单权限与 API 权限控制
- 审计日志中心
- 通知模板、通知消息、异步通知队列
- VPC、IP、云硬盘、快照、备份、安全组资源台账
- 服务生命周期操作：同步、开通、暂停、续费、终止
- 魔方云接口整理文档：
  [mofang-integration.md](/Users/Administrator/Desktop/IDC/docs/mofang-integration.md)

## 主要接口

- `/api/auth/*`
- `/api/dashboard`
- `/api/customers`
- `/api/customers/[id]/balance`
- `/api/customers/[id]/ledger`
- `/api/products`
- `/api/orders`
- `/api/services`
- `/api/services/[id]/action`
- `/api/invoices`
- `/api/payments`
- `/api/payments/[id]/refund`
- `/api/payments/callback`
- `/api/notifications`
- `/api/notifications/process`
- `/api/tickets`
- `/api/tickets/[id]/reply`
- `/api/billing`
- `/api/billing/settings`
- `/api/billing/run`
- `/api/resources`
- `/api/audits`

## 后台页面

- 运营概览
- 客户管理
- 产品管理
- 订单管理
- 服务管理
- 计费引擎
- 资源中心
- 审计日志
- 账单管理
- 收款管理
- 通知中心
- 工单中心

## 安装

```bash
npm install
npm run db:reset
npm --prefix admin-ui install
```

## 启动

后端：

```bash
npm run dev:backend
```

前端：

```bash
npm run dev:frontend
```

## 默认账号

- 管理后台：`http://localhost:5173`
- 健康检查：`http://127.0.0.1:3000/api/health`
- 超级管理员：`admin@idc.local`
- 默认密码：`Admin123!`

预置角色账号：

- 财务：`finance@idc.local`
- 客服：`support@idc.local`
- 运维：`ops@idc.local`

## 角色说明

- `SUPER_ADMIN`：全量权限
- `FINANCE`：客户、账单、收款、计费、审计
- `SUPPORT`：工单与只读查询
- `OPERATIONS`：服务、资源、工单与只读查询

## 验证命令

```bash
npm run db:reset
npm run build
npm --prefix admin-ui run build
npm run lint
```

## 关键环境变量

- `DATABASE_URL`
- `AUTH_SECRET`
- `SEED_ADMIN_EMAIL`
- `SEED_ADMIN_PASSWORD`
- `MOFANG_CLOUD_BASE_URL`
- `MOFANG_CLOUD_API_KEY`
- `MOFANG_CLOUD_API_SECRET`
- `MOFANG_CLOUD_USE_MOCK`
- `PAYMENT_CALLBACK_SECRET`

## 支付回调说明

- 回调接口：`/api/payments/callback`
- 请求头需带：`x-payment-secret`
- 未配置环境变量时，默认密钥为：`dev-callback-secret`

## 当前边界

- 默认仍以 `MOFANG_CLOUD_USE_MOCK=true` 运行，便于本地开发与演示。
- 已具备完整商用后台骨架，但真实生产接入仍建议继续补齐：
  - 真实支付渠道签名校验与渠道侧退款闭环
  - 发票、税率、对账导出
  - 更细粒度的商品映射与魔方云参数模板
  - 通知通道正式投递、失败重试与部署监控

## 官方公开参考

- 魔方文档入口：[my.idcsmart.com/doc](https://my.idcsmart.com/doc/)
- 文档真实数据接口：[my.idcsmart.com/v1/doc](https://my.idcsmart.com/v1/doc)
