# MySQL 本地初始化

这份文档只保留通用初始化步骤，不再记录某台机器的安装路径、账号密码或本地目录。

真实 MySQL 凭据请放在未提交的 `.env.local` 中。

## 建议版本

- MySQL `8.x`
- 数据库字符集：`utf8mb4`

## 1. 创建数据库和项目账号

下面是一个可直接改造的示例：

```sql
CREATE DATABASE IF NOT EXISTS idc_finance
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'idc_finance'@'%' IDENTIFIED BY 'change-me';
GRANT ALL PRIVILEGES ON idc_finance.* TO 'idc_finance'@'%';
FLUSH PRIVILEGES;
```

如果你只允许本机连接，也可以把 `'%'` 改成 `'localhost'` 或 `'127.0.0.1'`。

## 2. 填写本地环境变量

先复制模板：

```bash
cp .env.example .env.local
```

然后把 `.env.local` 改成 MySQL 模式，例如：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:change-me@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

## 3. 执行迁移和演示数据

在工作区目录执行：

```bash
npm run db:prepare:mysql
```

如果你只想分开执行，也可以：

```bash
npm run db:migrate:mysql
npm run db:seed:mysql
```

## 4. 以 MySQL 模式启动 API

```bash
npm run dev:api:mysql
```

## 5. 初始化后的效果

初始化完成后，MySQL 中会有一套可直接演示和联调的数据，覆盖：

- 商品、价格矩阵、配置项、资源模板
- 客户、联系人、实名状态
- 工单
- 订单、账单、支付、服务
- Provider 账号
- Provider 同步日志与资源明细
- 自动化任务与改配单示例

## 说明

- `.env.local` 已被 `.gitignore` 忽略，可以存放当前机器的真实 DSN。
- `STORAGE_STRICT=true` 时，如果 MySQL 不可用，API 会直接退出，不再静默回退到内存仓储。
- 如果只做页面开发或不需要持久化联调，可以改回 `STORAGE_DRIVER=memory`。
