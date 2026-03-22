$ErrorActionPreference = "Stop"

$mariaRoot = "C:\tools\mariadb104\mariadb-10.4.34-winx64"
$dataRoot = "C:\mariadb104-legacy\data"
$configPath = "C:\Users\Administrator\Desktop\IDC\scripts\mariadb104-legacy.ini"
$serviceName = "MariaDB104Legacy"
$rootPassword = "MariaDB104Root!2026"

New-Item -ItemType Directory -Force $dataRoot | Out-Null

$existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if (!$existingService) {
  & "$mariaRoot\bin\mysql_install_db.exe" --datadir=$dataRoot --service=$serviceName --password=$rootPassword --port=3308 --allow-remote-root-access
}

Start-Service -Name $serviceName
Start-Sleep -Seconds 5

& "$mariaRoot\bin\mysql.exe" --protocol=tcp -h 127.0.0.1 -P 3308 -u root "-p$rootPassword" -e "CREATE DATABASE IF NOT EXISTS zjmf_legacy374 DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci; CREATE USER IF NOT EXISTS 'zjmf_legacy'@'127.0.0.1' IDENTIFIED BY 'ZjmfLegacy!2026'; GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'127.0.0.1'; CREATE USER IF NOT EXISTS 'zjmf_legacy'@'localhost' IDENTIFIED BY 'ZjmfLegacy!2026'; GRANT ALL PRIVILEGES ON zjmf_legacy374.* TO 'zjmf_legacy'@'localhost'; FLUSH PRIVILEGES;"
