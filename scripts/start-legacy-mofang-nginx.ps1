$ErrorActionPreference = "Stop"

$phpCgi = "C:\tools\php74\7.4.33-ts\php-cgi.exe"
$phpIni = "C:\Users\Administrator\Desktop\IDC\scripts\legacy-mofang-php74.ini"
$nginx = "C:\tools\nginx\nginx-1.28.2\nginx.exe"
$nginxConf = "C:\Users\Administrator\Desktop\IDC\scripts\legacy-mofang-nginx.conf"
$publicRoot = "C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public"
$phpPorts = @(9074, 9075, 9076, 9077)

Get-NetTCPConnection -LocalPort 8091 -State Listen -ErrorAction SilentlyContinue |
  Select-Object -ExpandProperty OwningProcess -Unique |
  ForEach-Object {
    Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue
  }

foreach ($port in $phpPorts) {
  Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue |
    Select-Object -ExpandProperty OwningProcess -Unique |
    ForEach-Object {
      Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue
    }
}

Get-Process php-cgi -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Get-Process nginx -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue

foreach ($port in $phpPorts) {
  Start-Process -FilePath $phpCgi -ArgumentList @("-c", $phpIni, "-b", "127.0.0.1:$port") -WorkingDirectory $publicRoot | Out-Null
}
Start-Sleep -Seconds 1

Start-Process -FilePath $nginx -ArgumentList @("-p", "C:\tools\nginx\nginx-1.28.2\", "-c", $nginxConf) -WorkingDirectory "C:\tools\nginx\nginx-1.28.2" | Out-Null
Start-Sleep -Seconds 1

$listenPorts = @(8091) + $phpPorts

Get-NetTCPConnection -LocalPort $listenPorts -State Listen |
  Select-Object LocalAddress,LocalPort,OwningProcess |
  Sort-Object LocalPort
