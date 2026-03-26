# GitHub 拉取安装说明

适用于从 GitHub 直接拉取 `无穷云IDC业务管理系统` 后的首次安装。

## 环境要求

- Node.js 20+
- npm 10+
- Go 1.24+
- MySQL 8.x

## 安装步骤

### 1. 拉取仓库

```bash
git clone https://github.com/xiaobaiyuan666/bysat-idc.git
cd bysat-idc/idc-finance
```

### 2. 安装依赖

```bash
npm install
go mod download
```

### 3. 创建本地环境文件

```bash
cp .env.example .env.local
```

至少保证以下变量正确：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

### 4. 初始化数据库

纯净安装：

```bash
npm run db:migrate:mysql
```

如果你想直接看演示界面和流程：

```bash
npm run db:prepare:mysql
```

### 5. 启动服务

API：

```bash
npm run dev:api:mysql
```

后台：

```bash
npm run dev:admin
```

用户中心：

```bash
npm run dev:portal
```

## 默认访问地址

- API：`http://127.0.0.1:18080/api/v1/health`
- 后台：`http://localhost:5177`
- 用户中心：`http://localhost:5178`

## 默认演示账号

- 后台：`admin / Admin123!`
- 用户中心：`portal / Portal123!`

## 纯净项目与演示数据的区别

- GitHub 默认是纯项目源码，不包含你的本地运行状态和测试库快照
- `seed/mysql` 只是可选演示数据脚本
- 想保持纯净环境，只迁移，不执行 seed
