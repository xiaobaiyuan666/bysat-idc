# MySQL 安装与初始化

本项目要求 `MySQL 8.x`，字符集必须使用 `utf8mb4`。

## 1. 创建数据库

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 2. 配置项目连接

在 `.env.local` 中填写：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

## 3. 执行迁移

只安装表结构：

```bash
npm run db:migrate:mysql
```

安装表结构并导入演示数据：

```bash
npm run db:prepare:mysql
```

## 4. 启动 API

```bash
npm run dev:api:mysql
```

健康检查：

```bash
curl http://127.0.0.1:18080/api/v1/health
```

## 说明

- GitHub 仓库不提交你的本地 MySQL 数据目录、运行日志和 `.env.local`
- `seed/mysql` 里的 SQL 是通用演示数据，不是本地真实生产数据
- 如果只想要纯净项目，请只执行迁移，不执行 seed
