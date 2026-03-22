$ErrorActionPreference = "Stop"

$phpRoot = "C:\tools\php74\7.4.33-ts"
$phpIni = "C:\Users\Administrator\Desktop\IDC\scripts\legacy-mofang-php74.ini"
$sessionDir = "C:\tools\php74\session"
$tempDir = "C:\tools\php74\tmp"
$publicRoot = "C:\Users\Administrator\Desktop\IDC\idccw-src\3.7.4\public"
$router = "C:\Users\Administrator\Desktop\IDC\scripts\legacy-mofang-router.php"
$port = 8091

New-Item -ItemType Directory -Force $sessionDir | Out-Null
New-Item -ItemType Directory -Force $tempDir | Out-Null

Set-Location $publicRoot

& "$phpRoot\php.exe" -c $phpIni -S "127.0.0.1:$port" -t $publicRoot $router
