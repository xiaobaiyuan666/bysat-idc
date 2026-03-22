param(
  [Parameter(ValueFromRemainingArguments = $true)]
  [string[]]$Args
)

$candidates = @(
  $env:GO_BIN,
  "C:\Program Files\Go\bin\go.exe",
  "C:\tools\go\bin\go.exe"
) | Where-Object { $_ -and (Test-Path $_) }

if (@($candidates).Count -gt 0) {
  & @($candidates)[0] @Args
  exit $LASTEXITCODE
}

$goCommand = Get-Command go -ErrorAction SilentlyContinue
if ($goCommand) {
  & $goCommand.Source @Args
  exit $LASTEXITCODE
}

Write-Error "Go executable not found. Set GO_BIN or install Go to C:\Program Files\Go\bin\go.exe."
exit 1
