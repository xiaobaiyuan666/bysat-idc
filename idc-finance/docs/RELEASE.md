# 发布包安装说明

适用于把源码仓库构建成“编译后可直接部署”的 Windows 安装版。

## 一、适用场景

你应该在以下场景使用这份文档：

- 你要把系统发给别人部署
- 你要交给测试环境使用
- 你要做演示环境
- 你不想单独启动后台和用户中心开发服务

## 二、先在源码仓库生成安装包

进入项目目录：

```bash
cd idc-finance
```

安装依赖：

```bash
npm install
go mod download
```

生成安装包：

```bash
npm run build:release
```

## 三、会生成什么

默认生成：

- `output/release/windows-amd64/wuqiongyun-idc`
- `output/release/windows-amd64/wuqiongyun-idc.zip`

其中包含：

- `wuqiongyun-idc-api.exe`
- `web/admin`
- `web/portal`
- `database/install-clean.sql`
- `database/install-demo.sql`
- `.env.example`
- `start.ps1`
- `start.bat`
- `docs/*.md`

## 四、部署步骤

### 1. 创建数据库

```sql
CREATE DATABASE idc_finance CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 导入数据库

纯净安装：

```bash
mysql -u root -p idc_finance < database/install-clean.sql
```

演示安装：

```bash
mysql -u root -p idc_finance < database/install-demo.sql
```

### 3. 配置环境文件

复制：

```bash
copy .env.example .env.local
```

至少修改：

```env
MYSQL_DSN=idc_finance:你的密码@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4
```

如果需要对接魔方云或上游财务，再补充对应配置。

### 4. 启动安装版

PowerShell：

```powershell
.\start.ps1
```

或 Windows 双击：

```text
start.bat
```

## 五、启动后访问地址

- 后台：`http://127.0.0.1:18080/admin/`
- 用户中心：`http://127.0.0.1:18080/portal/`
- API：`http://127.0.0.1:18080/api/v1/health`

## 六、默认演示账号

只有导入 `install-demo.sql` 时才有：

- 后台：`admin / Admin123!`
- 用户中心：`portal / Portal123!`

## 七、发布版的特点

- 只启动一个主程序
- 后台和用户中心不需要再单独运行 Vite
- 更适合部署、演示、交付

## 八、补充文档

- [编译与使用手册](./BUILD-AND-USAGE.md)
- [前后端组件说明](./COMPONENTS.md)
- [源码安装说明](./INSTALL.md)
