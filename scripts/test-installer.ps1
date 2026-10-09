# Simulated release tests for install.ps1. Mirrors scripts/test-installer.sh.
#
#   pwsh -NoProfile -File scripts/test-installer.ps1                         # installer under pwsh
#   pwsh -NoProfile -File scripts/test-installer.ps1 -InstallerShell powershell  # Windows PowerShell 5.1
#
# Runs on Windows and on pwsh for Linux (e.g. mcr.microsoft.com/powershell). Off Windows the
# fake mudarro.exe is a shell script and the user PATH is a file (MUDARRO_TEST_PATH_FILE).
# Only Windows proves: the HKCU\Environment PATH (REG_EXPAND_SZ), Windows PowerShell 5.1,
# NTFS junctions as reparse points and a real PE mudarro.exe (built here with Go).
# On Windows this test edits the real user PATH and restores it at the end.
param([string]$InstallerShell = 'pwsh')
Set-StrictMode -Version 2.0
$ErrorActionPreference = 'Stop'

$repoDir = Split-Path -Parent $PSScriptRoot
$installer = Join-Path $repoDir 'install.ps1'
$onWindows = [System.Environment]::OSVersion.Platform -eq [System.PlatformID]::Win32NT
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ('mudarro-test-' + [guid]::NewGuid().ToString('N'))
$payload = Join-Path $tmp 'payload'
$assets = Join-Path $tmp 'assets'
$install = Join-Path $tmp 'install'
$foreign = Join-Path $tmp 'foreign'
New-Item -ItemType Directory -Path $payload, $assets, $foreign, (Join-Path $tmp 'fake') | Out-Null
Add-Type -AssemblyName System.IO.Compression
Add-Type -AssemblyName System.IO.Compression.FileSystem

function Assert([bool]$Condition, [string]$Message) { if (-not $Condition) { throw "FALHOU: $Message" } }
function Assert-SameFile([string]$A, [string]$B) {
  Assert ((Get-FileHash -LiteralPath $A).Hash -eq (Get-FileHash -LiteralPath $B).Hash) "$A difere de $B"
}

$arch = $env:PROCESSOR_ARCHITEW6432
if (-not $arch) { $arch = $env:PROCESSOR_ARCHITECTURE }
if (-not $arch) { $arch = [string][System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture }
$arch = if ($arch -match '^(ARM64|Arm64)$') { 'arm64' } else { 'amd64' }
$asset = "mudarro_windows_$arch.zip"
$exe = Join-Path $install 'mudarro.exe'

# Fake binaries: print a fixed line for "version", or fail.
function New-FakeBinary([string]$Path, [string]$Output, [int]$Code) {
  if ($onWindows) {
    $src = Join-Path $tmp 'fake\main.go'
    $go = "package main`nimport (`"fmt`"; `"os`")`nvar out, code string`nfunc main() { fmt.Println(out); if code != `"0`" { os.Exit(1) } }`n"
    [System.IO.File]::WriteAllText($src, $go)
    $env:GOOS = 'windows'; $env:GOARCH = $arch; $env:CGO_ENABLED = '0'
    & go build -o $Path -ldflags "-X main.out=$Output -X main.code=$Code" $src
    if ($LASTEXITCODE -ne 0) { throw 'go build do binário falso falhou' }
  } else {
    [System.IO.File]::WriteAllText($Path, "#!/bin/sh`necho $Output`nexit $Code`n")
  }
}

function New-Release([string[]]$Members, [string[]]$ExtraUnsafe = @()) {
  $zipPath = Join-Path $assets $asset
  if (Test-Path -LiteralPath $zipPath) { Remove-Item -LiteralPath $zipPath -Force }
  $zip = [System.IO.Compression.ZipFile]::Open($zipPath, [System.IO.Compression.ZipArchiveMode]::Create)
  try {
    foreach ($m in $Members) {
      [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile($zip, (Join-Path $payload $m), $m) | Out-Null
    }
    foreach ($u in $ExtraUnsafe) {
      $w = New-Object System.IO.StreamWriter($zip.CreateEntry($u).Open())
      try { $w.Write('evil') } finally { $w.Dispose() }
    }
  } finally { $zip.Dispose() }
  $hash = (Get-FileHash -LiteralPath $zipPath -Algorithm SHA256).Hash.ToLowerInvariant()
  [System.IO.File]::WriteAllText((Join-Path $assets 'checksums.txt'), "$hash  $asset`n")
}

function Invoke-Installer([string[]]$Arguments = @(), [switch]$Iex) {
  if ($Iex) {
    $cmd = "Invoke-Expression ([System.IO.File]::ReadAllText('$($installer.Replace("'", "''"))'))"
    & $InstallerShell -NoProfile -NonInteractive -Command $cmd *>&1 | Out-Host
  } else {
    & $InstallerShell -NoProfile -NonInteractive -ExecutionPolicy Bypass -File $installer @Arguments *>&1 | Out-Host
  }
  return ($LASTEXITCODE -eq 0)
}
function Assert-Refused([string]$Message, [string[]]$Arguments = @()) {
  Assert (-not (Invoke-Installer $Arguments)) "aceitou: $Message"
  Write-Host "ok (recusado): $Message"
}

function Get-TestUserPath {
  if ($onWindows) {
    $key = [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment')
    try { return [string]$key.GetValue('Path', '', 'DoNotExpandEnvironmentNames') } finally { $key.Close() }
  }
  return [System.IO.File]::ReadAllText($env:MUDARRO_TEST_PATH_FILE)
}
function Get-PathCount {
  @((Get-TestUserPath).Split(';') | Where-Object { $_.TrimEnd('\', '/') -ieq $install }).Count
}

$savedUserPath = $null
$savedVars = @{}
foreach ($v in 'MUDARRO_INSTALLER_TEST', 'MUDARRO_TEST_ASSETS_DIR', 'MUDARRO_TEST_PATH_FILE', 'MUDARRO_INSTALL_DIR', 'MUDARRO_VERSION', 'MUDARRO_REPOSITORY') {
  $savedVars[$v] = [System.Environment]::GetEnvironmentVariable($v)
}
try {
  if ($onWindows) {
    $key = [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment')
    $savedUserPath = $key.GetValue('Path', $null, 'DoNotExpandEnvironmentNames'); $key.Close()
    $seed = if ($savedUserPath) { $savedUserPath } else { '%USERPROFILE%\mudarro-test-unexpanded' }
    $k = [Microsoft.Win32.Registry]::CurrentUser.CreateSubKey('Environment')
    $k.SetValue('Path', $seed, 'ExpandString'); $k.Close()
  } else {
    $env:MUDARRO_TEST_PATH_FILE = Join-Path $tmp 'user-path.txt'
    [System.IO.File]::WriteAllText($env:MUDARRO_TEST_PATH_FILE, '/opt/outra;;%USERPROFILE%\bin;')
  }
  $pathBefore = Get-TestUserPath
  $env:MUDARRO_INSTALLER_TEST = '1'
  $env:MUDARRO_TEST_ASSETS_DIR = $assets
  $env:MUDARRO_INSTALL_DIR = $install
  $env:MUDARRO_VERSION = 'v0.1.0'
  Remove-Item Env:MUDARRO_REPOSITORY -ErrorAction SilentlyContinue

  Copy-Item -LiteralPath (Join-Path $repoDir 'LICENSE'), (Join-Path $repoDir 'THIRD_PARTY_NOTICES.md') -Destination $payload
  Copy-Item -LiteralPath (Join-Path $repoDir 'LICENSES') -Destination $payload -Recurse
  $licenseFiles = @(Get-ChildItem -LiteralPath (Join-Path $payload 'LICENSES') -File | ForEach-Object { "LICENSES/$($_.Name)" })
  $full = @('mudarro.exe', 'LICENSE', 'THIRD_PARTY_NOTICES.md') + $licenseFiles
  New-FakeBinary (Join-Path $payload 'mudarro.exe') 'test-version' 0

  # Simulated release, installed twice (update), then once through irm | iex.
  New-Release $full
  Assert (Invoke-Installer) 'primeira instalação'
  Assert (Invoke-Installer) 'repetição da instalação'
  Assert (Invoke-Installer -Iex) 'instalação via Invoke-Expression'
  $out = (& $exe version | Out-String).Trim()
  Assert ($out -eq 'test-version') "versão instalada: $out"
  Assert-SameFile (Join-Path $payload 'LICENSE') (Join-Path $install 'mudarro-licenses\LICENSE')
  Assert-SameFile (Join-Path $payload 'THIRD_PARTY_NOTICES.md') (Join-Path $install 'mudarro-licenses\THIRD_PARTY_NOTICES.md')
  foreach ($l in $licenseFiles) { Assert-SameFile (Join-Path $payload $l) (Join-Path $install "mudarro-licenses\$l") }
  Assert ((Get-PathCount) -eq 1) 'PATH do usuário deve conter o diretório exatamente uma vez'
  Assert ((Get-TestUserPath).StartsWith($pathBefore)) 'PATH anterior preservado (inclusive %VAR% não expandida)'
  Assert (@(Get-ChildItem -LiteralPath $install -Force -Filter '.mudarro-install-*').Count -eq 0) 'sobrou temporário'
  Write-Host 'ok: instalação, repetição, iex, LICENSE, avisos e PATH sem duplicata'
  Copy-Item -LiteralPath $exe -Destination (Join-Path $tmp 'installed-before')

  # Invalid configuration.
  foreach ($bad in '1.0.0', 'v1.0', 'v1.0.0;calc', 'latest ', '../v1.0.0') {
    $env:MUDARRO_VERSION = $bad; Assert-Refused "versão inválida '$bad'"
  }
  $env:MUDARRO_VERSION = 'v0.1.0'
  foreach ($bad in 'evil/../x', 'a b/c', 'owner', 'owner/..', 'https://x/y', 'a/b/c') {
    $env:MUDARRO_REPOSITORY = $bad; Assert-Refused "repositório inválido '$bad'"
  }
  Remove-Item Env:MUDARRO_REPOSITORY
  foreach ($bad in 'relativo\mudarro', "$install;C:\x", (Join-Path $install '..\x')) {
    $env:MUDARRO_INSTALL_DIR = $bad; Assert-Refused "diretório inválido '$bad'"
  }
  $env:MUDARRO_INSTALL_DIR = $install
  $env:MUDARRO_TEST_ASSETS_DIR = 'relativo'; Assert-Refused 'base de teste relativa'
  $env:MUDARRO_TEST_ASSETS_DIR = $assets

  # Incomplete payload (a referenced license is missing) preserves binary and notices.
  New-FakeBinary (Join-Path $payload 'mudarro.exe') 'forbidden-partial-update' 0
  New-Release (@('mudarro.exe', 'LICENSE', 'THIRD_PARTY_NOTICES.md') + $licenseFiles[0])
  Assert-Refused 'payload de licenças incompleto'
  New-Release (@('mudarro.exe', 'THIRD_PARTY_NOTICES.md') + $licenseFiles)
  Assert-Refused 'payload sem LICENSE'
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe
  Assert-SameFile (Join-Path $payload 'THIRD_PARTY_NOTICES.md') (Join-Path $install 'mudarro-licenses\THIRD_PARTY_NOTICES.md')

  # A binary that fails "version" never replaces the installed one.
  New-FakeBinary (Join-Path $payload 'mudarro.exe') 'broken' 1
  New-Release $full
  Assert-Refused 'binário que falha em version'
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe

  # Unsafe archive members.
  New-FakeBinary (Join-Path $payload 'mudarro.exe') 'test-version' 0
  foreach ($unsafe in '../evil.txt', '/abs/evil.txt', 'C:/evil.txt', 'LICENSES/..\..\evil.txt', 'ads.txt:stream') {
    New-Release $full @($unsafe)
    Assert-Refused "membro inseguro '$unsafe'"
  }
  Assert (-not (Test-Path -LiteralPath (Join-Path $tmp 'evil.txt'))) 'escreveu fora do destino'
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe

  # Checksum: divergent, absent, malformed.
  New-Release $full
  $sums = Join-Path $assets 'checksums.txt'
  $good = [System.IO.File]::ReadAllText($sums)
  [System.IO.File]::WriteAllText($sums, ('0' * 64) + "  $asset`n"); Assert-Refused 'checksum divergente'
  [System.IO.File]::WriteAllText($sums, "$($good.Split(' ')[0])  outro.zip`n"); Assert-Refused 'checksum ausente'
  [System.IO.File]::WriteAllText($sums, "XYZ  $asset`n"); Assert-Refused 'checksum malformado'
  Remove-Item -LiteralPath $sums; Assert-Refused 'checksums.txt ausente'
  [System.IO.File]::WriteAllText($sums, $good)
  [System.IO.File]::AppendAllText((Join-Path $assets $asset), 'corrupt'); Assert-Refused 'pacote corrompido'
  New-Release $full
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe

  # Reparse-point destinations (symlink here; NTFS junction on Windows).
  [System.IO.File]::WriteAllText((Join-Path $foreign 'sentinel'), 'owned sentinel')
  function New-Link([string]$Path, [string]$Target) {
    if ($onWindows) { New-Item -ItemType Junction -Path $Path -Target $Target | Out-Null }
    else { New-Item -ItemType SymbolicLink -Path $Path -Target $Target | Out-Null }
  }
  function Remove-Link([string]$Path) {
    if ($onWindows) { [System.IO.Directory]::Delete($Path, $false) } else { & rm -- $Path }
  }
  function Assert-ForeignUntouched {
    Assert ([System.IO.File]::ReadAllText((Join-Path $foreign 'sentinel')) -eq 'owned sentinel') 'sentinela alterada'
    Assert (@(Get-ChildItem -LiteralPath $foreign -Recurse -Force).Count -eq 1) 'arquivos novos no destino alheio'
  }
  $notices = Join-Path $install 'mudarro-licenses'
  Move-Item -LiteralPath $notices -Destination (Join-Path $tmp 'saved-licenses')
  New-Link $notices $foreign
  Assert-Refused 'mudarro-licenses como reparse point'
  Assert-ForeignUntouched
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe
  Remove-Link $notices
  Move-Item -LiteralPath (Join-Path $tmp 'saved-licenses') -Destination $notices
  Move-Item -LiteralPath (Join-Path $notices 'LICENSES') -Destination (Join-Path $tmp 'saved-inner')
  New-Link (Join-Path $notices 'LICENSES') $foreign
  Assert-Refused 'LICENSES como reparse point'
  Assert-ForeignUntouched
  Remove-Link (Join-Path $notices 'LICENSES')
  Move-Item -LiteralPath (Join-Path $tmp 'saved-inner') -Destination (Join-Path $notices 'LICENSES')
  $linkedInstall = Join-Path $tmp 'linked-install'
  New-Link $linkedInstall $foreign
  $env:MUDARRO_INSTALL_DIR = $linkedInstall
  Assert-Refused 'diretório de instalação como reparse point'
  Assert-ForeignUntouched
  Remove-Link $linkedInstall
  $env:MUDARRO_INSTALL_DIR = $install
  if (-not $onWindows) {
    # File symlink as mudarro.exe (on Windows creating one needs privilege; covered by the junctions).
    Move-Item -LiteralPath $exe -Destination (Join-Path $tmp 'saved-exe')
    New-Item -ItemType SymbolicLink -Path $exe -Target (Join-Path $foreign 'sentinel') | Out-Null
    Assert-Refused 'mudarro.exe como symlink'
    Assert-ForeignUntouched
    Remove-Item -LiteralPath $exe -Force
    Move-Item -LiteralPath (Join-Path $tmp 'saved-exe') -Destination $exe
  }

  # Test variables are ignored outside test mode.
  Remove-Item Env:MUDARRO_INSTALLER_TEST
  if ($onWindows) {
    # Without test mode the HTTPS base is used; this tag does not exist, so the install fails.
    $env:MUDARRO_VERSION = 'v0.0.0-installer-test-does-not-exist'
    Assert-Refused 'variáveis de teste fora do modo de teste'
    $env:MUDARRO_VERSION = 'v0.1.0'
  } else {
    Assert-Refused 'execução fora do Windows sem modo de teste'
  }
  $env:MUDARRO_INSTALLER_TEST = '1'
  Assert-SameFile (Join-Path $tmp 'installed-before') $exe
  Assert ((Get-PathCount) -eq 1) 'PATH duplicado após recusas'

  # Uninstall keeps unrelated files and the rest of the user PATH.
  [System.IO.File]::WriteAllText((Join-Path $install 'project-file.txt'), 'keep')
  Assert (Invoke-Installer @('-Uninstall')) 'desinstalação'
  Assert (-not (Test-Path -LiteralPath $exe)) 'mudarro.exe não removido'
  Assert (-not (Test-Path -LiteralPath $notices)) 'avisos não removidos'
  Assert (Test-Path -LiteralPath (Join-Path $install 'project-file.txt')) 'arquivo alheio removido'
  Assert ((Get-PathCount) -eq 0) 'PATH ainda contém o diretório'
  # A trailing ';' before our entry is not restored; everything else must be identical.
  Assert ((Get-TestUserPath).TrimEnd(';') -eq $pathBefore.TrimEnd(';')) 'PATH anterior alterado pela desinstalação'
  Assert (Invoke-Installer @('-Uninstall')) 'desinstalação repetida'
  Write-Host "Instalação, atualização e remoção simuladas ($InstallerShell): OK"
} finally {
  if ($onWindows) {
    $k = [Microsoft.Win32.Registry]::CurrentUser.CreateSubKey('Environment')
    if ($null -eq $savedUserPath) { $k.DeleteValue('Path', $false) } else { $k.SetValue('Path', $savedUserPath, 'ExpandString') }
    $k.Close()
  }
  foreach ($v in $savedVars.Keys) { [System.Environment]::SetEnvironmentVariable($v, $savedVars[$v]) }
  Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
}
