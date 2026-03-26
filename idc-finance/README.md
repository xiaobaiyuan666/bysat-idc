# 无穷云IDC业务管理系统

江苏白猿网络科技有限公司 - 猿创软件开发

基于 `Go + Vue 3 + Element Plus + vue-pure-admin` 的 IDC 业务管理系统，包含：

- 后端 API
- 后台运营端
- 用户中心
- 订单/账单/服务/工单
- 资源与上游对接
- 自动化任务与改配链路

## 这是什么仓库

这个仓库默认是“源码仓库”，适合：

- 本地开发
- 二次开发
- 测试环境编译
- 生成可分发安装包

仓库默认不提交以下本地内容：

- `.env.local`
- `.runtime/`
- `data/*.json`
- 本地日志
- 本地数据库目录
- 你的测试快照

## 四个核心程序

| 组件 | 目录 | 用途 | 开发启动 | 编译命令 |
| --- | --- | --- | --- | --- |
| 后端 API | `apps/api` | 提供所有业务接口、登录鉴权、静态托管 | `npm run dev:api:mysql` | `go build ./...` |
| 后台运营端 | `apps/admin-web` | 管理员、客服、财务、运维后台 | `npm run dev:admin` | `npm run build:admin` |
| 用户中心 | `apps/portal-web` | 客户登录、下单、账单、工单、服务管理 | `npm run dev:portal` | `npm run build:portal` |
| 发布安装包 | `output/release/...` | 编译后可直接部署的版本 | 无 | `npm run build:release` |

## 两种使用方式

### 1. 源码开发模式

适合你自己继续开发。

```bash
npm install
go mod download
cp .env.example .env.local
npm run db:prepare:mysql
npm run dev:api:mysql
npm run dev:admin
npm run dev:portal
```

访问地址：

- 后台：`http://localhost:5177`
- 用户中心：`http://localhost:5178`
- API：`http://127.0.0.1:18080/api/v1/health`

### 2. 编译安装模式

适合给别人部署、演示、测试。

```bash
npm install
go mod download
npm run build:release
```

生成目录：

- `output/release/windows-amd64/wuqiongyun-idc`
- `output/release/windows-amd64/wuqiongyun-idc.zip`

编译后的版本使用单端口：

- 后台：`http://127.0.0.1:18080/admin/`
- 用户中心：`http://127.0.0.1:18080/portal/`
- API：`http://127.0.0.1:18080/api/v1/health`

## 默认演示账号

只有导入演示数据后才有：

- 后台：`admin / Admin123!`
- 用户中心：`portal / Portal123!`

## 文档入口

- [源码安装说明](./docs/INSTALL.md)
- [发布包安装说明](./docs/RELEASE.md)
- [编译与使用手册](./docs/BUILD-AND-USAGE.md)
- [前后端组件说明](./docs/COMPONENTS.md)
- [MySQL 初始化说明](./docs/mysql-setup.md)

## 最常用命令

```bash
# 安装依赖
npm install
go mod download

# 初始化数据库
npm run db:migrate:mysql
npm run db:seed:mysql

# 开发启动
npm run dev:api:mysql
npm run dev:admin
npm run dev:portal

# 单独编译
npm run build:admin
npm run build:portal
go build ./...

# 生成发布安装包
npm run build:release
```
