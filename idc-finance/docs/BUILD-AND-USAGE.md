# 编译与使用手册

这份文档专门讲清楚：

- 编译前要准备什么
- 后端怎么编译
- 后台怎么编译
- 用户中心怎么编译
- 编译后会生成什么
- 编译后的版本怎么用

## 一、编译前准备

必须先安装：

- Node.js 20+
- npm 10+
- Go 1.24+
- MySQL 8.x

建议先执行：

```bash
npm install
go mod download
```

然后复制环境文件：

```bash
cp .env.example .env.local
```

至少修改：

```env
STORAGE_DRIVER=mysql
STORAGE_STRICT=true
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

## 二、编译前先准备数据库

先创建数据库：

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

如果你只是开发联调：

```bash
npm run db:prepare:mysql
```

如果你只想导入结构：

```bash
npm run db:migrate:mysql
```

## 三、分别编译各个程序

### 1. 编译后端

```bash
go build ./...
```

作用：

- 检查 Go 代码是否能通过编译
- 生成后端程序

### 2. 编译后台

```bash
npm run build:admin
```

编译结果：

- 输出到 `apps/admin-web/dist`

### 3. 编译用户中心

```bash
npm run build:portal
```

编译结果：

- 输出到 `apps/portal-web/dist`

## 四、生成完整安装包

如果你要一个“别人可以直接部署”的版本，不要只编译单个前端或后端，要执行：

```bash
npm run build:release
```

它会自动完成：

- 编译后台
- 编译用户中心
- 编译后端主程序
- 复制文档
- 复制 SQL
- 生成安装目录
- 生成 zip 包

## 五、编译后会生成什么

默认生成位置：

- `output/release/windows-amd64/wuqiongyun-idc`
- `output/release/windows-amd64/wuqiongyun-idc.zip`

目录结构：

- `wuqiongyun-idc-api.exe`
- `web/admin`
- `web/portal`
- `database/install-clean.sql`
- `database/install-demo.sql`
- `.env.example`
- `start.ps1`
- `start.bat`
- `docs/*.md`

## 六、编译后怎么用

### 1. 导入数据库

纯净安装：

```bash
mysql -u root -p idc_finance < database/install-clean.sql
```

演示安装：

```bash
mysql -u root -p idc_finance < database/install-demo.sql
```

### 2. 配置环境文件

复制：

```bash
copy .env.example .env.local
```

至少改好：

```env
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

### 3. 启动安装版

PowerShell：

```powershell
.\start.ps1
```

或 Windows 双击：

```text
start.bat
```

### 4. 启动后地址

- 后台：`http://127.0.0.1:18080/admin/`
- 用户中心：`http://127.0.0.1:18080/portal/`
- 健康检查：`http://127.0.0.1:18080/api/v1/health`

## 七、源码模式和安装版模式的区别

### 源码模式

- 适合开发
- 前后端分别启动
- 后台走 `5177`
- 用户中心走 `5178`
- API 走 `18080`

### 安装版模式

- 适合部署和演示
- 只启动一个主程序
- 后台和用户中心都由 API 托管
- 默认统一走 `18080`

## 八、默认账号

只有导入演示数据时才有：

- 后台：`admin / Admin123!`
- 用户中心：`portal / Portal123!`

## 九、常见问题

### 1. GitHub 拉下来为什么不能直接双击运行

因为 GitHub 默认给的是源码仓库，不是编译后的安装版。你需要：

- 按源码模式启动
- 或执行 `npm run build:release` 生成安装版

### 2. 为什么有人说数据库缺少

因为源码仓库不会提交你的本地 MySQL 数据，只会提交：

- 迁移脚本
- 演示 seed
- 发布包 SQL 合并脚本结果

### 3. 编译后后台和用户中心为什么不用再单独启动

因为发布模式下，后端主程序已经会托管：

- `/admin/`
- `/portal/`

### 4. 如果要接魔方云或上游财务怎么办

在 `.env.local` 里继续填：

- `MOFANG_CLOUD_*`
- `FINANCE_UPSTREAM_*`

然后再启动后端即可。
