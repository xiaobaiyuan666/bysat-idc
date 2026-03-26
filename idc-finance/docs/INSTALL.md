# 源码安装说明

适用于从 GitHub 拉取源码后，以“开发模式”运行无穷云IDC业务管理系统。

## 一、适用场景

你应该在以下场景使用这份文档：

- 你要继续开发这个系统
- 你要调试前后端
- 你要单独修改后台或用户中心
- 你要联调魔方云或上游财务接口

如果你只是想要一个能直接部署的版本，请改看 [发布包安装说明](./RELEASE.md)。

## 二、环境要求

- Node.js 20+
- npm 10+
- Go 1.24+
- MySQL 8.x

## 三、拉取仓库

```bash
git clone https://github.com/xiaobaiyuan666/bysat-idc.git
cd bysat-idc/idc-finance
```

## 四、安装依赖

```bash
npm install
go mod download
```

## 五、创建环境文件

复制环境模板：

```bash
cp .env.example .env.local
```

至少保证以下配置正确：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

如果需要对接魔方云：

```env
MOFANG_CLOUD_BASE_URL=
MOFANG_CLOUD_USERNAME=
MOFANG_CLOUD_PASSWORD=
MOFANG_CLOUD_LANG=zh-cn
```

如果需要对接上游财务：

```env
FINANCE_UPSTREAM_BASE_URL=
FINANCE_UPSTREAM_USERNAME=
FINANCE_UPSTREAM_PASSWORD=
FINANCE_UPSTREAM_SOURCE_NAME=上游财务
```

## 六、初始化数据库

先创建数据库：

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

纯净结构安装：

```bash
npm run db:migrate:mysql
```

导入演示数据：

```bash
npm run db:prepare:mysql
```

## 七、启动开发环境

### 1. 启动 API

```bash
npm run dev:api:mysql
```

### 2. 启动后台

```bash
npm run dev:admin
```

### 3. 启动用户中心

```bash
npm run dev:portal
```

## 八、开发环境地址

- 后台：`http://localhost:5177`
- 用户中心：`http://localhost:5178`
- API：`http://127.0.0.1:18080/api/v1/health`

## 九、默认演示账号

只有导入演示数据后才有：

- 后台：`admin / Admin123!`
- 用户中心：`portal / Portal123!`

## 十、源码模式的特点

- 后端、后台、用户中心分别启动
- 前端走 Vite 开发服务
- 修改代码后便于热更新
- 最适合继续开发和排错

## 十一、补充文档

- [编译与使用手册](./BUILD-AND-USAGE.md)
- [前后端组件说明](./COMPONENTS.md)
- [MySQL 初始化说明](./mysql-setup.md)
