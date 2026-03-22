# MySQL 本地初始化

当前这台机器已经完成本地 MySQL 安装，信息如下：

- 安装目录：`C:\tools\mysql\mysql\current`
- 数据目录：`C:\mysql-data\mysql\data`
- 服务名：`MySQLIDC`
- 端口：`3306`
- 项目数据库：`idc_finance`
- 项目账号：`idc_finance`

## 常用命令

查看服务状态：

```powershell
Get-Service -Name MySQLIDC
```

连接 MySQL：

```powershell
& "C:\tools\mysql\mysql\current\bin\mysql.exe" -u root -p
```

使用项目脚本初始化数据库：

```powershell
$env:MYSQL_DSN='idc_finance:IdcFinance!2026@tcp(127.0.0.1:3306)/idc_finance?parseTime=true&charset=utf8mb4'
npm run db:prepare:mysql
```

用 MySQL 模式启动 API：

```powershell
npm run dev:api:mysql
```

## 说明

- `.env.local` 已写入项目本地 MySQL DSN，且已被 `.gitignore` 忽略。
- `STORAGE_STRICT=true` 已启用；如果 MySQL 连接失败，API 会直接退出，不再静默回退到内存仓储。
