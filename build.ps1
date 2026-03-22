# Build and Deploy Web Apps Script
# Replicates the Makefile build process for Windows

param(
    [switch]$SkipBuild = $false,
    [switch]$SkipServe = $false,
    [switch]$ServeOnly = $false,
    [switch]$InstallAll = $false,
    [ValidateSet("all", "one", "admin", "compose", "workflow", "reporter", "discovery", "privacy")]
    [string]$App = "all"
)

$ErrorActionPreference = "Stop"
$ProjectRoot = $PSScriptRoot
$LibDir = Join-Path $ProjectRoot "lib"
$ClientDir = Join-Path $ProjectRoot "client"
$WebAppsDir = Join-Path $ClientDir "web"
$ServerPublicDir = Join-Path $ProjectRoot "server\webapp\public"
$ServerDir = Join-Path $ProjectRoot "server"

# Colors for output
$Cyan = [System.ConsoleColor]::Cyan
$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$Red = [System.ConsoleColor]::Red

# Function to build libraries with proper symlink handling
function Build-Libraries {
    Write-Host "`n---Building js---" -ForegroundColor $Cyan
    Push-Location (Join-Path $LibDir "js")
    try {
        yarn
        Write-Host "---Cleaning old symlinks for js---" -ForegroundColor $Yellow
        yarn unlink 2>$null | Out-Null
        Write-Host "---Creating sym link for js---" -ForegroundColor $Yellow
        yarn link
        yarn build
        Write-Host "JS library built successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to build js: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location

    Write-Host "`n---Building vue---" -ForegroundColor $Cyan
    Push-Location (Join-Path $LibDir "vue")
    try {
        yarn
        Write-Host "---Cleaning old symlinks for vue---" -ForegroundColor $Yellow
        yarn unlink 2>$null | Out-Null
        Write-Host "---Linking vue to js library---" -ForegroundColor $Yellow
        yarn link @cortezaproject/corteza-js
        Write-Host "---Creating sym link for vue---" -ForegroundColor $Yellow
        yarn link
        yarn build
        Write-Host "Vue library built successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to build vue: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
}

Write-Host "==============================================" -ForegroundColor $Green
if ($ServeOnly) {
    Write-Host "Serving Corteza API (no build)" -ForegroundColor $Green
} elseif ($InstallAll) {
    Write-Host "Installing and Building ALL Corteza Web Apps" -ForegroundColor $Green
} elseif ($App -eq "all") {
    Write-Host "Building and Deploying ALL Corteza Web Apps" -ForegroundColor $Green
} else {
    Write-Host "Building and Deploying: $App" -ForegroundColor $Green
}
Write-Host "==============================================" -ForegroundColor $Green

# Check if gin is available (and not aliased to Get-ComputerInfo)
if (-not $SkipServe) {
    Remove-Item Alias:gin -ErrorAction SilentlyContinue
    if (-not (Get-Command "gin" -ErrorAction SilentlyContinue)) {
        Write-Host "gin not found. Install it with: go install github.com/codegangsta/gin@latest" -ForegroundColor $Red
        Write-Host "Or re-run with -SkipServe to skip the server step." -ForegroundColor $Yellow
        exit 1
    }
}

# If ServeOnly, skip straight to serve
if ($ServeOnly) {
    Write-Host "`nStarting gin server..." -ForegroundColor $Cyan
    Write-Host "Press Ctrl+C to stop`n" -ForegroundColor $Yellow
    Push-Location $ServerDir
    try {
        gin --laddr localhost --build cmd/corteza --bin build/gin-bin -- --env-file .env serve-api
    } finally {
        Pop-Location
    }
    exit 0
}

# Check if Make is available
if (-not (Get-Command "make" -ErrorAction SilentlyContinue)) {
    Write-Host "Make is not available. Please install Make for Windows (e.g., via Chocolatey or WSL)" -ForegroundColor $Red
    Write-Host "Press any key to exit..."
    $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}

# All app definitions
$allApps = @(
    @{ Source = "one"; Destination = "" },
    @{ Source = "one"; Destination = "one" },
    @{ Source = "admin"; Destination = "admin" },
    @{ Source = "compose"; Destination = "compose" },
    @{ Source = "workflow"; Destination = "workflow" },
    @{ Source = "reporter"; Destination = "reporter" },
    @{ Source = "discovery"; Destination = "discovery" },
    @{ Source = "privacy"; Destination = "privacy" }
)

# Filter to just the requested app
$appsToCopy = if ($App -eq "all") {
    $allApps
} else {
    $allApps | Where-Object { $_.Source -eq $App }
}

# Build process
if ($InstallAll) {
    Write-Host "`nBuilding libraries..." -ForegroundColor $Cyan
    Build-Libraries
    Write-Host "Libraries built successfully" -ForegroundColor $Green

    Write-Host "`nInstalling and building all client apps..." -ForegroundColor $Cyan
    Push-Location $ClientDir
    try {
        make fresh
        Write-Host "make fresh completed successfully" -ForegroundColor $Green
        make build
        Write-Host "make build completed successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to install/build client apps: $_" -ForegroundColor $Red
        exit 1
    }
    Pop-Location
} elseif (-not $SkipBuild) {
    if ($App -eq "all") {
        Write-Host "`nBuilding libraries..." -ForegroundColor $Cyan
        Build-Libraries
        Write-Host "Libraries built successfully" -ForegroundColor $Green

        Write-Host "`nBuilding all client apps..." -ForegroundColor $Cyan
        Push-Location $ClientDir
        try {
            make build
            Write-Host "Client apps built successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to build client apps: $_" -ForegroundColor $Red
            exit 1
        }
        Pop-Location
    } else {
        Write-Host "`nBuilding libraries..." -ForegroundColor $Cyan
        Build-Libraries
        Write-Host "Libraries built successfully" -ForegroundColor $Green
        Write-Host "`nBuilding $App only..." -ForegroundColor $Cyan
        Push-Location (Join-Path $WebAppsDir $App)
        try {
            yarn build
            Write-Host "$App built successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to build ${App}: $_" -ForegroundColor $Red
            exit 1
        }
        Pop-Location
    }
} else {
    Write-Host "Skipping build - copying existing dist files..." -ForegroundColor $Yellow
}

# Clean destination dirs for selected apps
Write-Host "`nCleaning destinations..." -ForegroundColor $Cyan
foreach ($entry in $appsToCopy) {
    $destPath = Join-Path $ServerPublicDir $entry.Destination
    if (Test-Path $destPath) {
        Write-Host "Removing $destPath..." -ForegroundColor $Yellow
        Remove-Item -Path $destPath -Recurse -Force
    }
    New-Item -Path $destPath -ItemType Directory -Force | Out-Null
}

# Copy
Write-Host "`nCopying files..." -ForegroundColor $Cyan
foreach ($entry in $appsToCopy) {
    $srcPath = Join-Path $WebAppsDir "$($entry.Source)\dist"
    $destPath = Join-Path $ServerPublicDir $entry.Destination

    if (-not (Test-Path $srcPath)) {
        Write-Host "Dist directory not found for $($entry.Source): $srcPath" -ForegroundColor $Red
        exit 1
    }

    Write-Host "Copying $($entry.Source) -> $destPath..." -ForegroundColor $Yellow
    Copy-Item -Path "$srcPath\*" -Destination $destPath -Recurse -Force
}

Write-Host "`n==============================================" -ForegroundColor $Green
Write-Host "Done! App: $App" -ForegroundColor $Green
Write-Host "Output: $ServerPublicDir" -ForegroundColor $Cyan
Write-Host "==============================================" -ForegroundColor $Green

# Serve
if (-not $SkipServe) {
    Write-Host "`nStarting gin server..." -ForegroundColor $Cyan
    Write-Host "Press Ctrl+C to stop`n" -ForegroundColor $Yellow
    Push-Location $ServerDir
    try {
        gin --laddr localhost --build cmd/corteza --bin build/gin-bin -- --env-file .env serve-api
    } finally {
        Pop-Location
    }
} else {
    Write-Host "`nPress any key to exit..."
    $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}
