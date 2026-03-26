# MySQL 初始化说明

`无穷云IDC业务管理系统` 目前要求使用 `MySQL 8.x`，并建议统一使用 `utf8mb4`。

## 1. 创建数据库

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 2. 配置连接

在 `.env.local` 中设置：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

## 3. 导入结构

仅导入表结构：

```bash
npm run db:migrate:mysql
```

导入结构和演示数据：

```bash
npm run db:prepare:mysql
```

如果你使用发布包，可以直接导入以下 SQL：

- `database/install-clean.sql`
- `database/install-demo.sql`

## 4. 启动 API

```bash
npm run dev:api:mysql
```

健康检查：

```bash
curl http://127.0.0.1:18080/api/v1/health
```

## 说明

- GitHub 仓库不提交你的本地 MySQL 数据目录、运行日志和 `.env.local`。
- `seed/mysql` 仅用于演示安装，不代表生产数据。
