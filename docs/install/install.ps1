# install.ps1 — Run with:
# irm https://diamondosas.github.io/errpipe/install/install.ps1 | iex

$ErrorActionPreference = "Stop"

$ToolName    = "errpipe"
$RepoOwner   = "diamondosas"
$RepoName    = "errpipe"
$InstallDir  = "$env:LOCALAPPDATA\$ToolName"   # per-user, no admin needed

# ── 1. Fetch latest release download URL from GitHub API ──────────────────────
Write-Host "Fetching latest release..."
try {
    $release = Invoke-RestMethod "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
} catch {
    throw "Failed to fetch latest release from GitHub API. Ensure the repository has a public release."
}

$asset = $release.assets | Where-Object { $_.name -match "windows" -and $_.name -like "*.exe" -or $_.name -eq "$ToolName.exe" } | Select-Object -First 1

if (-not $asset) { throw "No suitable .exe asset found in latest release." }

# ── 2. Download the exe ────────────────────────────────────────────────────────
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}
$dest = "$InstallDir\$ToolName.exe"

Write-Host "Downloading $($asset.name)..."
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $dest -UseBasicParsing

# ── 3. Add to PATH (user-level, no admin required) ────────────────────────────
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")

if ($userPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable(
        "PATH",
        "$userPath;$InstallDir",
        "User"
    )
    Write-Host "Added $InstallDir to PATH."
} else {
    Write-Host "$InstallDir already in PATH."
}

# ── 4. Refresh PATH in current session ────────────────────────────────────────
$env:PATH = [Environment]::GetEnvironmentVariable("PATH","Machine") + ";" +
            [Environment]::GetEnvironmentVariable("PATH","User")

Write-Host ""
Write-Host "✅ $ToolName installed! Run: $ToolName"
