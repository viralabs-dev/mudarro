# Mudarro installer for Windows (Windows PowerShell 5.1 and PowerShell 7+).
#
#   irm https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.ps1 | iex
#   powershell -NoProfile -ExecutionPolicy Bypass -File install.ps1 [-Uninstall]
#   & ([scriptblock]::Create((irm https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.ps1))) -Uninstall
#
# Environment:
#   MUDARRO_VERSION      latest (default) or a tag vX.Y.Z
#   MUDARRO_INSTALL_DIR  absolute directory (default %LOCALAPPDATA%\Programs\mudarro)
#   MUDARRO_REPOSITORY   owner/name on GitHub (default viralabs-dev/mudarro)
#
# Test-only (scripts/test-installer.ps1), ignored unless MUDARRO_INSTALLER_TEST=1:
#   MUDARRO_TEST_ASSETS_DIR  absolute local directory that replaces the HTTPS download base
#   MUDARRO_TEST_PATH_FILE   file that stands in for the user PATH when not running on Windows
#
# The script never changes the execution policy, never needs administrator rights,
# never reads from stdin and only touches the user PATH (HKCU\Environment).
param([switch]$Uninstall)

try {
& {
  # Child scope: strict mode and preferences do not leak into the caller's session (irm | iex).
  # Failures throw; the handler at the end exits 1 only in -File mode, so iex never
  # closes the user's terminal.
  param([bool]$Uninstall)
  Set-StrictMode -Version 2.0
  $ErrorActionPreference = 'Stop'
  $ProgressPreference = 'SilentlyContinue'

  # The file stays pure ASCII (Windows PowerShell 5.1 reads BOM-less -File scripts as ANSI,
  # and a BOM breaks irm | iex), so accented text is written as \uXXXX and expanded here.
  # Arguments are formatted after expansion so paths are never unescaped.
  function Text([string]$Template, [object[]]$Arguments) {
    $t = [regex]::Replace($Template, '\\u([0-9a-fA-F]{4})', { param($m) [string][char][Convert]::ToInt32($m.Groups[1].Value, 16) })
    if ($Arguments) { return $t -f $Arguments }
    return $t
  }
  function Fail([string]$Template, [object[]]$Arguments) { throw [System.InvalidOperationException]::new((Text $Template $Arguments)) }
  function Say([string]$Template, [object[]]$Arguments) { Write-Host (Text $Template $Arguments) }

  $onWindows = [System.Environment]::OSVersion.Platform -eq [System.PlatformID]::Win32NT
  $testMode = $env:MUDARRO_INSTALLER_TEST -eq '1'

  function Test-Reparse([string]$Path) {
    $item = Get-Item -LiteralPath $Path -Force -ErrorAction SilentlyContinue
    if ($null -eq $item) { return $false }
    if ($item.Attributes -band [System.IO.FileAttributes]::ReparsePoint) { return $true }
    $link = $item.PSObject.Properties['LinkType']
    return ($null -ne $link -and $null -ne $link.Value -and "$($link.Value)" -ne '')
  }

  function Get-PathEntries([string]$Value) {
    if ([string]::IsNullOrEmpty($Value)) { return @() }
    return @($Value.Split(';') | Where-Object { $_.Trim() -ne '' })
  }

  function Test-SamePath([string]$A, [string]$B) {
    $x = [System.Environment]::ExpandEnvironmentVariables($A.Trim().Trim('"')).TrimEnd('\', '/')
    $y = $B.TrimEnd('\', '/')
    return [string]::Equals($x, $y, [System.StringComparison]::OrdinalIgnoreCase)
  }

  function Read-UserPath {
    if ($onWindows) {
      $key = [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment', $false)
      if ($null -eq $key) { return '' }
      try {
        return [string]$key.GetValue('Path', '', [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames)
      } finally { $key.Close() }
    }
    $file = $env:MUDARRO_TEST_PATH_FILE
    if ($testMode -and $file -and (Test-Path -LiteralPath $file -PathType Leaf)) {
      return [System.IO.File]::ReadAllText($file).Trim()
    }
    return ''
  }

  function Write-UserPath([string]$Value) {
    if ($onWindows) {
      $key = [Microsoft.Win32.Registry]::CurrentUser.CreateSubKey('Environment')
      try {
        # Keep REG_EXPAND_SZ so entries such as %USERPROFILE%\bin stay unexpanded.
        $key.SetValue('Path', $Value, [Microsoft.Win32.RegistryValueKind]::ExpandString)
      } finally { $key.Close() }
      # Setting and clearing a user variable makes Windows broadcast WM_SETTINGCHANGE,
      # so new terminals see the PATH change without compiling P/Invoke code.
      [System.Environment]::SetEnvironmentVariable('MUDARRO_INSTALLER_REFRESH', '1', 'User')
      [System.Environment]::SetEnvironmentVariable('MUDARRO_INSTALLER_REFRESH', $null, 'User')
      return
    }
    if ($testMode -and $env:MUDARRO_TEST_PATH_FILE) {
      [System.IO.File]::WriteAllText($env:MUDARRO_TEST_PATH_FILE, $Value)
    }
  }

  function Write-Atomic([string]$Source, [string]$Destination) {
    $dir = Split-Path -Parent $Destination
    $temp = Join-Path $dir ('.mudarro-install-' + [guid]::NewGuid().ToString('N') + '.tmp')
    try {
      [System.IO.File]::Copy($Source, $temp, $true)
      if (Test-Path -LiteralPath $Destination) {
        [System.IO.File]::Replace($temp, $Destination, [NullString]::Value)
      } else {
        [System.IO.File]::Move($temp, $Destination)
      }
    } finally {
      if (Test-Path -LiteralPath $temp) { Remove-Item -LiteralPath $temp -Force }
    }
  }

  # ---- configuration and validation -------------------------------------------------
  $repo = if ($env:MUDARRO_REPOSITORY) { $env:MUDARRO_REPOSITORY } else { 'viralabs-dev/mudarro' }
  $version = if ($env:MUDARRO_VERSION) { $env:MUDARRO_VERSION } else { 'latest' }
  if ($env:MUDARRO_INSTALL_DIR) {
    $installDir = $env:MUDARRO_INSTALL_DIR
  } else {
    $localAppData = $env:LOCALAPPDATA
    if (-not $localAppData) { $localAppData = [System.Environment]::GetFolderPath('LocalApplicationData') }
    if (-not $localAppData) { Fail 'LOCALAPPDATA indefinido; defina MUDARRO_INSTALL_DIR.' }
    $installDir = Join-Path (Join-Path $localAppData 'Programs') 'mudarro'
  }

  if (-not $onWindows -and -not $testMode) { Fail 'Sistema n\u00e3o suportado; use install.sh no Linux e no macOS.' }
  if ($version -ne 'latest' -and $version -cnotmatch '^v[0-9]+\.[0-9]+\.[0-9]+([.-][A-Za-z0-9.-]+)?$') {
    Fail 'MUDARRO_VERSION deve ser latest ou uma tag vX.Y.Z.'
  }
  if ($repo -cnotmatch '^[A-Za-z0-9][A-Za-z0-9-]{0,38}/[A-Za-z0-9._-]{1,100}$' -or $repo -match '/\.{1,2}$') {
    Fail 'MUDARRO_REPOSITORY deve ter o formato dono/reposit\u00f3rio.'
  }
  $absolute = if ($onWindows) { '^([A-Za-z]:[\\/]|\\\\[^\\/]+[\\/][^\\/]+)' } else { '^/' }
  if ($installDir -notmatch $absolute -or $installDir.IndexOfAny([char[]]';"*?<>|') -ge 0 -or
      $installDir -match '(^|[\\/])\.\.?([\\/]|$)' -or ($onWindows -and $installDir.Substring(2).Contains(':'))) {
    Fail 'MUDARRO_INSTALL_DIR deve ser um caminho absoluto, sem "..", ";" ou caracteres especiais.'
  }
  $installDir = [System.IO.Path]::GetFullPath($installDir).TrimEnd('\', '/')
  if ($installDir -eq '' -or $installDir -match '^[A-Za-z]:$') { Fail 'MUDARRO_INSTALL_DIR n\u00e3o pode ser a raiz do disco.' }

  $exeName = 'mudarro.exe'
  $exePath = Join-Path $installDir $exeName
  $noticeDir = Join-Path $installDir 'mudarro-licenses'
  $noticeLicenses = Join-Path $noticeDir 'LICENSES'

  # ---- uninstall ---------------------------------------------------------------------
  if ($Uninstall) {
    foreach ($p in @($installDir, $exePath, $noticeDir, $noticeLicenses)) {
      if (Test-Reparse $p) { Fail 'Destino \u00e9 reparse point (link/jun\u00e7\u00e3o): {0}; remo\u00e7\u00e3o interrompida.' $p }
    }
    if (Test-Path -LiteralPath $exePath) { Remove-Item -LiteralPath $exePath -Force }
    if (Test-Path -LiteralPath $noticeDir) { Remove-Item -LiteralPath $noticeDir -Recurse -Force }
    # Remove only our entry; keep every other byte of the user PATH (empty entries, %VARS%).
    $userPath = Read-UserPath
    $parts = @($userPath.Split(';'))
    $kept = @($parts | Where-Object { $_.Trim() -eq '' -or -not (Test-SamePath $_ $installDir) })
    if ($kept.Count -ne $parts.Count) { Write-UserPath ($kept -join ';') }
    if ((Test-Path -LiteralPath $installDir) -and -not (Get-ChildItem -LiteralPath $installDir -Force)) {
      Remove-Item -LiteralPath $installDir -Force
    }
    Say 'Mudarro removido de {0}; os arquivos dos projetos foram preservados.' $installDir
    return
  }

  # ---- download ----------------------------------------------------------------------
  $arch = $env:PROCESSOR_ARCHITEW6432
  if (-not $arch) { $arch = $env:PROCESSOR_ARCHITECTURE }
  if (-not $arch) {
    try { $arch = [string][System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture } catch { $arch = '' }
  }
  switch -regex ($arch) {
    '^(AMD64|X64)$' { $arch = 'amd64'; break }
    '^ARM64$' { $arch = 'arm64'; break }
    default { Fail 'Arquitetura n\u00e3o suportada; use Windows amd64 ou arm64.' }
  }
  $asset = "mudarro_windows_$arch.zip"

  $testAssets = $null
  if ($testMode -and $env:MUDARRO_TEST_ASSETS_DIR) {
    $testAssets = $env:MUDARRO_TEST_ASSETS_DIR
    if (-not [System.IO.Path]::IsPathRooted($testAssets) -or -not (Test-Path -LiteralPath $testAssets -PathType Container)) {
      Fail 'MUDARRO_TEST_ASSETS_DIR deve ser um diret\u00f3rio local absoluto existente.'
    }
  }
  $base = if ($version -eq 'latest') { "https://github.com/$repo/releases/latest/download" } else { "https://github.com/$repo/releases/download/$version" }

  $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ('mudarro-install-' + [guid]::NewGuid().ToString('N'))
  New-Item -ItemType Directory -Path $tmp | Out-Null
  try {
    function Get-ReleaseFile([string]$Name) {
      $out = Join-Path $tmp $Name
      if ($testAssets) {
        $src = Join-Path $testAssets $Name
        if (-not (Test-Path -LiteralPath $src -PathType Leaf)) { Fail 'Falha ao baixar {0}.' $Name }
        Copy-Item -LiteralPath $src -Destination $out
        return $out
      }
      $uri = [System.Uri]"$base/$Name"
      if ($uri.Scheme -ne 'https') { Fail 'Somente HTTPS \u00e9 aceito.' }
      try {
        $response = Invoke-WebRequest -Uri $uri -OutFile $out -UseBasicParsing -MaximumRedirection 10 -PassThru
      } catch {
        Fail 'Falha ao baixar {0} de {1}. {2}' $Name, $uri, $_.Exception.Message
      }
      # Refuse a redirect that left HTTPS (Windows PowerShell 5.1 does not block it itself).
      $final = $null
      $raw = $response.BaseResponse
      if ($null -ne $raw) {
        if ($raw.PSObject.Properties['ResponseUri']) { $final = $raw.ResponseUri }
        elseif ($raw.PSObject.Properties['RequestMessage'] -and $raw.RequestMessage) { $final = $raw.RequestMessage.RequestUri }
      }
      if ($null -ne $final -and $final.Scheme -ne 'https') { Fail 'Download redirecionado para fora de HTTPS; instala\u00e7\u00e3o interrompida.' }
      return $out
    }

    if (-not $testAssets) {
      # TLS 1.2+ only. PowerShell 7 already negotiates TLS 1.2/1.3 through the OS.
      $protocols = [System.Net.SecurityProtocolType]::Tls12
      if ([enum]::GetNames([System.Net.SecurityProtocolType]) -contains 'Tls13') {
        $protocols = $protocols -bor [System.Net.SecurityProtocolType]::Tls13
      }
      [System.Net.ServicePointManager]::SecurityProtocol = $protocols
    }

    $zipPath = Get-ReleaseFile $asset
    $sumsPath = Get-ReleaseFile 'checksums.txt'

    # ---- checksum ----------------------------------------------------------------------
    $expected = @(Get-Content -LiteralPath $sumsPath | ForEach-Object {
        $parts = $_.Trim() -split '\s+'
        if ($parts.Count -eq 2 -and $parts[1].TrimStart('*') -ceq $asset) { $parts[0] }
      })
    if ($expected.Count -ne 1 -or $expected[0] -cnotmatch '^[0-9a-f]{64}$') { Fail 'Checksum ausente ou inv\u00e1lido.' }
    $actual = (Get-FileHash -LiteralPath $zipPath -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -cne $expected[0]) { Fail 'Checksum divergente; instala\u00e7\u00e3o interrompida.' }

    # ---- extraction of expected members only ------------------------------------------
    Add-Type -AssemblyName System.IO.Compression
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $stage = Join-Path $tmp 'stage'
    New-Item -ItemType Directory -Path (Join-Path $stage 'LICENSES') | Out-Null
    $zip = [System.IO.Compression.ZipFile]::OpenRead($zipPath)
    try {
      $entries = @{}
      foreach ($entry in $zip.Entries) {
        $name = $entry.FullName.Replace('\', '/')
        if ($name.StartsWith('/') -or $name.Contains(':') -or $name -match '(^|/)\.\.(/|$)') {
          Fail 'Pacote com caminho inseguro: {0}' $entry.FullName
        }
        if ($name.EndsWith('/')) { continue }
        if ($entries.ContainsKey($name)) { Fail 'Pacote com membro duplicado: {0}' $name }
        $entries[$name] = $entry
      }
      if (-not $entries.ContainsKey('THIRD_PARTY_NOTICES.md')) { Fail 'Pacote incompleto: THIRD_PARTY_NOTICES.md ausente.' }
      $wanted = New-Object System.Collections.Generic.List[string]
      foreach ($n in @($exeName, 'LICENSE', 'THIRD_PARTY_NOTICES.md')) { $wanted.Add($n) }

      function Expand-Member([string]$Name, [long]$Limit) {
        if (-not $entries.ContainsKey($Name)) { Fail 'Pacote incompleto: {0} ausente.' $Name }
        $entry = $entries[$Name]
        if ($entry.Length -gt $Limit) { Fail 'Membro grande demais no pacote: {0}' $Name }
        $dest = Join-Path $stage ($Name.Replace('/', [System.IO.Path]::DirectorySeparatorChar))
        $in = $entry.Open()
        try {
          $out = [System.IO.File]::Open($dest, [System.IO.FileMode]::CreateNew)
          try { $in.CopyTo($out) } finally { $out.Dispose() }
        } finally { $in.Dispose() }
        if ((Get-Item -LiteralPath $dest).Length -ne $entry.Length) { Fail 'Membro corrompido no pacote: {0}' $Name }
      }

      Expand-Member 'THIRD_PARTY_NOTICES.md' 1MB
      # Every license referenced by the notices must ship; a partial payload is refused.
      $notices = [System.IO.File]::ReadAllText((Join-Path $stage 'THIRD_PARTY_NOTICES.md'))
      $licenses = @([regex]::Matches($notices, 'LICENSES/([A-Za-z0-9][A-Za-z0-9._-]*\.txt)') | ForEach-Object { $_.Groups[1].Value } | Sort-Object -Unique)
      if ($licenses.Count -eq 0) { Fail 'Pacote incompleto: avisos sem licen\u00e7as.' }
      foreach ($l in $licenses) { $wanted.Add("LICENSES/$l") }
      foreach ($n in $wanted) {
        if ($n -eq 'THIRD_PARTY_NOTICES.md') { continue }
        $limit = if ($n -eq $exeName) { 512MB } else { 1MB }
        Expand-Member $n $limit
      }
    } finally { $zip.Dispose() }

    # ---- smoke test before replacing anything ----------------------------------------
    $stagedExe = Join-Path $stage $exeName
    if (-not $onWindows) { & chmod 755 $stagedExe }
    try {
      # Windows PowerShell 5.1 turns native stderr into terminating errors under 'Stop'.
      $ErrorActionPreference = 'Continue'
      $versionOutput = & $stagedExe version 2>&1 | Out-String
      $code = $LASTEXITCODE
      $ErrorActionPreference = 'Stop'
    } catch {
      Fail 'O bin\u00e1rio baixado n\u00e3o executou: {0}' $_.Exception.Message
    }
    if ($code -ne 0) { Fail "O bin\u00e1rio baixado falhou em 'mudarro version'; instala\u00e7\u00e3o interrompida." }

    # ---- destination checks -------------------------------------------------------------
    if (Test-Reparse $installDir) { Fail 'Diret\u00f3rio de instala\u00e7\u00e3o \u00e9 reparse point (link/jun\u00e7\u00e3o); escolha outro diret\u00f3rio.' }
    if (-not (Test-Path -LiteralPath $installDir)) { New-Item -ItemType Directory -Path $installDir | Out-Null }
    if (-not (Test-Path -LiteralPath $installDir -PathType Container)) { Fail 'MUDARRO_INSTALL_DIR n\u00e3o \u00e9 um diret\u00f3rio.' }
    if (Test-Reparse $exePath) { Fail 'Destino \u00e9 reparse point (link/jun\u00e7\u00e3o); escolha outro diret\u00f3rio.' }
    if ((Test-Reparse $noticeDir) -or (Test-Reparse $noticeLicenses)) {
      Fail 'Destino dos avisos \u00e9 reparse point (link/jun\u00e7\u00e3o); escolha outro diret\u00f3rio.'
    }
    foreach ($n in $wanted) {
      if ($n -eq $exeName) { continue }
      if (Test-Reparse (Join-Path $noticeDir $n)) { Fail 'Destino de licen\u00e7a \u00e9 reparse point.' }
    }

    # ---- install ----------------------------------------------------------------------------
    try {
      Write-Atomic $stagedExe $exePath
    } catch {
      Fail 'N\u00e3o foi poss\u00edvel substituir {0} (feche o mudarro em execu\u00e7\u00e3o ou verifique o antiv\u00edrus). {1}' $exePath, $_.Exception.Message
    }
    if (-not (Test-Path -LiteralPath $noticeLicenses)) { New-Item -ItemType Directory -Path $noticeLicenses -Force | Out-Null }
    foreach ($n in $wanted) {
      if ($n -eq $exeName) { continue }
      $rel = $n.Replace('/', [System.IO.Path]::DirectorySeparatorChar)
      Write-Atomic (Join-Path $stage $rel) (Join-Path $noticeDir $rel)
    }
  } finally {
    Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
  }

  Say 'Mudarro instalado: {0}' $exePath

  # ---- user PATH ----------------------------------------------------------------------------
  $userPath = Read-UserPath
  if (-not (@(Get-PathEntries $userPath) | Where-Object { Test-SamePath $_ $installDir })) {
    $newPath = if ($userPath -eq '' -or $userPath.EndsWith(';')) { $userPath + $installDir } else { $userPath + ';' + $installDir }
    Write-UserPath $newPath
    Say 'Adicionado ao PATH do usu\u00e1rio: {0} (abra um novo terminal).' $installDir
  }
  # Also refresh this session (useful with irm | iex) on Windows.
  if ($onWindows -and -not (@(Get-PathEntries $env:PATH) | Where-Object { Test-SamePath $_ $installDir })) {
    $env:PATH = (@(Get-PathEntries $env:PATH) + $installDir) -join ';'
  }
  Write-Host -NoNewline $versionOutput
} $Uninstall.IsPresent
} catch {
  $message = "Erro: $($_.Exception.Message)"
  # $PSCommandPath is set only when this file runs with -File (not under irm | iex).
  if ($PSCommandPath) { [Console]::Error.WriteLine($message); exit 1 }
  throw $message
}
