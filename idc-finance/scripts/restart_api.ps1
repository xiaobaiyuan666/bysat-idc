$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$exe = Join-Path $runtime "api-live.exe"
$envFile = Join-Path $root ".env.local"
$stdout = Join-Path $runtime "api.out.log"
$stderr = Join-Path $runtime "api.err.log"

$listener = Get-NetTCPConnection -LocalPort 18080 -State Listen -ErrorAction SilentlyContinue | Select-Object -First 1
if ($listener) {
  Stop-Process -Id $listener.OwningProcess -Force
  Start-Sleep -Seconds 1
}

if (Test-Path $envFile) {
  Get-Content $envFile | ForEach-Object {
    if ($_ -match '^([^#=]+)=(.*)$') {
      [Environment]::SetEnvironmentVariable($matches[1], $matches[2], "Process")
    }
  }
}

if (-not (Test-Path $exe)) {
  throw "API executable not found: $exe"
}

Start-Process -FilePath $exe -WorkingDirectory $root -RedirectStandardOutput $stdout -RedirectStandardError $stderr -PassThru | Out-Null
Start-Sleep -Seconds 2
Invoke-WebRequest -UseBasicParsing http://127.0.0.1:18080/api/v1/health -TimeoutSec 5 | Select-Object -ExpandProperty Content
