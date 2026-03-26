# 无穷云IDC业务管理系统

江苏白猿网络科技有限公司-猿创软件开发

基于 `Go + Vue 3 + Element Plus + vue-pure-admin` 的 IDC 财务、服务、工单、订单、账单、资源与自动化管理系统。

## 仓库说明

- GitHub 仓库默认只保留项目源码、迁移脚本、可选演示 seed、开发文档。
- 不提交本地开发环境数据，例如 `.env.local`、`.runtime/`、`data/*.json`、本地日志和临时二进制。
- 演示数据是通用示例数据，不是你本地真实测试环境快照。

## 目录结构

- `apps/api`: Go API 服务
- `apps/admin-web`: 管理后台
- `apps/portal-web`: 用户中心
- `internal`: 后端业务模块
- `migrations/mysql`: MySQL 迁移脚本
- `seed/mysql`: 可选演示数据
- `docs`: 安装、环境与分析文档

## 快速开始

### 1. 安装依赖

```bash
npm install
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为你自己的本地环境文件：

```bash
cp .env.example .env.local
```

核心变量示例：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:IdcFinance!2026@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

如果需要对接魔方云，再补充：

```env
MOFANG_CLOUD_BASE_URL=
MOFANG_CLOUD_USERNAME=
MOFANG_CLOUD_PASSWORD=
MOFANG_CLOUD_LANG=zh-cn
MOFANG_CLOUD_INSECURE_SKIP_VERIFY=true
```

### 3. 初始化数据库

先手动创建数据库：

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

只安装纯净项目：

```bash
npm run db:migrate:mysql
```

如果你需要演示数据：

```bash
npm run db:seed:mysql
```

或者一键完成迁移加演示数据：

```bash
npm run db:prepare:mysql
```

## 启动项目

### API

```bash
npm run dev:api:mysql
```

健康检查：

- `http://127.0.0.1:18080/api/v1/health`

### 管理后台

```bash
npm run dev:admin
```

- 地址：`http://localhost:5177`
- 默认账号：`admin / Admin123!`

### 用户中心

```bash
npm run dev:portal
```

- 地址：`http://localhost:5178`
- 默认账号：`portal / Portal123!`

## 是否包含演示数据

包含，但属于可选内容。

- `seed/mysql/0001_demo_data.sql`: 基础演示业务数据
- `seed/mysql/0002_demo_finance_accounts.sql`: 财务账户演示数据
- `seed/mysql/0003_demo_finance_cleanup.sql`: 财务演示数据修补
- `seed/mysql/0004_demo_text_cleanup.sql`: 文本脏数据修补

如果你想保持 GitHub 项目为纯净安装状态，只执行迁移，不执行 seed 即可。

## 安装文档

- 通用安装说明：[docs/INSTALL.md](./docs/INSTALL.md)
- MySQL 初始化说明：[docs/mysql-setup.md](./docs/mysql-setup.md)

## 当前主线能力

- 后台：客户、商品、订单、账单、支付、退款、服务、工单、Provider、报表、系统管理
- 用户中心：控制台、商城、服务、订单、账单、工单、钱包、账户
- 自动化：服务动作、改配单、自动关单、上游商品同步记录
- 对接能力：魔方云、上游财务商品同步、多渠道 Provider 工作台

## 构建验证

```bash
go build ./...
npm run build:admin
npm run build:portal
```
