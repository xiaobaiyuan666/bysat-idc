# 前后端组件说明

这份文档专门说明：后端、后台、用户中心分别是什么，开发时改哪里，编译后产物在哪。

## 1. 后端 API

目录：

- `apps/api`
- `internal`

作用：

- 提供后台接口 `/api/v1/admin/*`
- 提供用户中心接口 `/api/v1/portal/*`
- 负责登录鉴权、菜单、权限、业务逻辑
- 负责 MySQL 读写
- 编译安装版中负责托管后台和用户中心静态文件

开发启动：

```bash
npm run dev:api:mysql
```

源码编译：

```bash
go build ./...
```

发布包产物：

- `wuqiongyun-idc-api.exe`

## 2. 后台运营端

目录：

- `apps/admin-web`

作用：

- 管理员登录
- 客户、商品、订单、账单、服务、工单
- 资源与上游、自动化任务、报表、系统设置

开发启动：

```bash
npm run dev:admin
```

开发地址：

- `http://localhost:5177`

源码编译：

```bash
npm run build:admin
```

编译产物目录：

- `apps/admin-web/dist`

发布包内位置：

- `web/admin`

## 3. 用户中心

目录：

- `apps/portal-web`

作用：

- 客户登录
- 商城下单
- 服务管理
- 订单、账单、钱包、工单、账户中心

开发启动：

```bash
npm run dev:portal
```

开发地址：

- `http://localhost:5178`

源码编译：

```bash
npm run build:portal
```

编译产物目录：

- `apps/portal-web/dist`

发布包内位置：

- `web/portal`

## 4. 数据库

结构迁移目录：

- `migrations/mysql`

演示数据目录：

- `seed/mysql`

源码模式导入：

```bash
npm run db:migrate:mysql
npm run db:seed:mysql
```

发布包模式导入：

- `database/install-clean.sql`
- `database/install-demo.sql`

## 5. 编译前和编译后的区别

### 编译前

你看到的是源码：

- Go 源码
- Vue 源码
- SQL 迁移
- 演示数据
- 文档

这时候需要分别启动 API、后台、用户中心。

### 编译后

你得到的是安装版：

- 一个后端主程序
- 一套后台静态文件
- 一套用户中心静态文件
- 一套数据库安装 SQL
- 启动脚本

这时候只需要启动一个程序，后台和用户中心都由 API 托管。
