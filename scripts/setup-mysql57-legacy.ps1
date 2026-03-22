$ErrorActionPreference = "Stop"

$zipPath = "C:\tools\downloads\mysql-5.7.44-winx64.zip"
$extractRoot = "C:\tools\mysql57"
$mysqlRoot = "C:\tools\mysql57\mysql-5.7.44-winx64"
$dataRoot = "C:\mysql57-legacy\data"
$configPath = "C:\Users\Administrator\Desktop\IDC\scripts\mysql57-legacy.ini"
$serviceName = "MySQL57Legacy"

if (!(Test-Path $mysqlRoot)) {
  New-Item -ItemType Directory -Force $extractRoot | Out-Null
  Expand-Archive -Path $zipPath -DestinationPath $extractRoot -Force
}

New-Item -ItemType Directory -Force $dataRoot | Out-Null

$existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if (!$existingService) {
  & "$mysqlRoot\bin\mysqld.exe" --defaults-file=$configPath --initialize-insecure
  & "$mysqlRoot\bin\mysqld.exe" --install $serviceName --defaults-file=$configPath
}

Start-Service -Name $serviceName
Start-Sleep -Seconds 5

& "$mysqlRoot\bin\mysql.exe" -u root --protocol=tcp -P 3307 -e "CREATE DATABASE IF NOT EXISTS zjmf_legacy374 DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci; CREATE USER IF NOT EXISTS 'zjmf_legacy'@'127.0.0.1' IDENTIFIED BY 'ZjmfLegacy!2026'; GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'127.0.0.1'; CREATE USER IF NOT EXISTS 'zjmf_legacy'@'localhost' IDENTIFIED BY 'ZjmfLegacy!2026'; GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'localhost'; FLUSH PRIVILEGES;"
