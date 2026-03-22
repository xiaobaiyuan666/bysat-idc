param(
  [string]$BaseUrl = "http://127.0.0.1:8091",
  [string]$AdminPath = "admin01",
  [string]$Username = "admin",
  [string]$Password = "Admin123!",
  [string]$OutputPath = "C:\Users\Administrator\Desktop\IDC\docs\legacy-admin-login-response.json"
)

$ErrorActionPreference = "Stop"

$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession

$query = "request_time=$([DateTimeOffset]::Now.ToUnixTimeMilliseconds())&languagesys=CN"
$loginUri = "$BaseUrl/$AdminPath/login?$query"
$body = @{
  username = $Username
  password = $Password
  code = ""
  captcha = ""
}

$response = Invoke-WebRequest -Uri $loginUri -Method Post -Body $body -WebSession $session -UseBasicParsing -TimeoutSec 90
$json = $response.Content | ConvertFrom-Json

$dir = Split-Path -Parent $OutputPath
if (-not (Test-Path $dir)) {
  New-Item -ItemType Directory -Path $dir | Out-Null
}

$response.Content | Set-Content -Path $OutputPath -Encoding UTF8

[PSCustomObject]@{
  Status = $json.status
  Message = $json.msg
  OutputPath = $OutputPath
  TopMenus = @($json.data.rule | ForEach-Object { $_.title })
} | Format-List
