$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$envFile = Join-Path $projectRoot ".env.local"
$binary = Join-Path $projectRoot "tmp-api.exe"
$stdout = Join-Path $projectRoot "tmp-api-live.log"
$stderr = Join-Path $projectRoot "tmp-api-live.err.log"

if (-not (Test-Path $binary)) {
  throw "tmp-api.exe not found. Build it first."
}

if (Test-Path $envFile) {
  Get-Content $envFile | ForEach-Object {
    $line = $_.Trim()
    if (-not $line -or $line.StartsWith("#")) {
      return
    }

    $parts = $line.Split("=", 2)
    if ($parts.Length -eq 2) {
      Set-Item -Path ("Env:" + $parts[0]) -Value $parts[1]
    }
  }
}

foreach ($logFile in @($stdout, $stderr)) {
  if (Test-Path $logFile) {
    try {
      Remove-Item $logFile -Force -ErrorAction Stop
    } catch {
      Clear-Content $logFile -ErrorAction SilentlyContinue
    }
  }
}

Start-Process -FilePath $binary -WorkingDirectory $projectRoot -RedirectStandardOutput $stdout -RedirectStandardError $stderr -WindowStyle Hidden
