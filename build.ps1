# Build and Deploy Web Apps Script
# Replicates the Makefile build process for Windows

param(
    [switch]$SkipBuild = $false,
    [switch]$SkipServe = $false,
    [switch]$ServeOnly = $false,
    [switch]$InstallAll = $false,
    [switch]$CodeGen = $false,
    [switch]$Dev = $false,
    [switch]$Test = $false,
    [switch]$Lint = $false,
    [switch]$Fresh = $false,
    [switch]$Audit = $false,
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
if ($CodeGen) {
    Write-Host "Running Server Code Generation" -ForegroundColor $Green
} elseif ($Dev) {
    Write-Host "Running Development Build" -ForegroundColor $Green
} elseif ($Test) {
    Write-Host "Running Tests" -ForegroundColor $Green
} elseif ($Lint) {
    Write-Host "Running Linting" -ForegroundColor $Green
} elseif ($Fresh) {
    Write-Host "Running Fresh Build" -ForegroundColor $Green
} elseif ($Audit) {
    Write-Host "Running Security Audit" -ForegroundColor $Green
} elseif ($ServeOnly) {
    Write-Host "Serving Corteza API (no build)" -ForegroundColor $Green
} elseif ($InstallAll) {
    Write-Host "Installing and Building ALL Corteza Web Apps" -ForegroundColor $Green
} elseif ($App -eq "all") {
    Write-Host "Building and Deploying ALL Corteza Web Apps" -ForegroundColor $Green
} else {
    Write-Host "Building and Deploying: $App" -ForegroundColor $Green
}
Write-Host "==============================================" -ForegroundColor $Green

# Run server code generation if requested
if ($CodeGen) {
    Write-Host "`n---Running server code generation---" -ForegroundColor $Cyan
    Push-Location $ServerDir
    try {
        $gopath = & go env GOPATH
        $codegenPath = Join-Path $gopath "bin\corteza-codegen.exe"
        $jsontplPath = Join-Path $gopath "bin\corteza-json-tpl-exec.exe"
        $cuePath = Join-Path $gopath "bin\cue.exe"
        
        # Install CUE if not available
        if (-not (Test-Path $cuePath)) {
            Write-Host "---Installing CUE---" -ForegroundColor $Yellow
            & go install cuelang.org/go/cmd/cue@v0.4.2
        }
        
        # Build corteza-codegen if not available
        if (-not (Test-Path $codegenPath)) {
            Write-Host "---Building corteza-codegen---" -ForegroundColor $Yellow
            Push-Location "cmd\codegen"
            & go build -o $codegenPath main.go
            Pop-Location
        }
        
        # Rebuild corteza-json-tpl-exec to ensure we have the latest version with our changes
        Write-Host "---Rebuilding corteza-json-tpl-exec---" -ForegroundColor $Yellow
        Push-Location "codegen\tool"
        & go build -o $jsontplPath .
        Pop-Location
        
        # Run code generation for server files
        $serverCueFiles = Get-ChildItem "codegen\server.*.cue"
        foreach ($cueFile in $serverCueFiles) {
            Write-Host "`n---Generating server files from $($cueFile.Name)---" -ForegroundColor $Yellow
            & $cuePath eval $cueFile.FullName --out json | & $jsontplPath -v -p (Join-Path $PWD.Path "codegen\assets\templates") -b $PWD.Path
        }
        
        # Check if DOCS_DIR is set for docs generation
        if ($env:DOCS_DIR) {
            Write-Host "`n---DOCS_DIR is set, generating documentation files---" -ForegroundColor $Cyan
            $docsCueFiles = Get-ChildItem "codegen\docs.*.cue"
            foreach ($cueFile in $docsCueFiles) {
                Write-Host "`n---Generating doc files from $($cueFile.Name) (dst: $($env:DOCS_DIR))---" -ForegroundColor $Yellow
                & $cuePath eval $cueFile.FullName --out json | & $jsontplPath -v -p (Join-Path $PWD.Path "codegen\assets\templates") -b $env:DOCS_DIR
            }
        } else {
            Write-Host "`n---Skipping docs generation: DOCS_DIR is not set---" -ForegroundColor $Yellow
            Write-Host "  To generate docs, set the DOCS_DIR environment variable"
            Write-Host "  Example: `$env:DOCS_DIR='C:\path\to\corteza-docs'; .\build.ps1 -CodeGen"
        }
        
        Write-Host "`nCode generation completed successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to run code generation: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "Code generation done!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

# Run dev command (build libs and clients for development)
if ($Dev) {
    Write-Host "`n---Processing libs---" -ForegroundColor $Cyan
    Push-Location $LibDir
    try {
        Write-Host "---Building js---" -ForegroundColor $Yellow
        Push-Location "js"
        yarn
        yarn build
        yarn link
        Pop-Location
        
        Write-Host "---Building vue---" -ForegroundColor $Yellow
        Push-Location "vue"
        yarn
        yarn cdeps
        yarn build
        yarn link
        Pop-Location
        
        Write-Host "Libs processed successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to process libs: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n---Processing clients---" -ForegroundColor $Cyan
    $webSubdirs = Get-ChildItem -Path (Join-Path $ClientDir "web") -Directory
    foreach ($dir in $webSubdirs) {
        Write-Host "`n---Installing and linking clients $($dir.Name)---" -ForegroundColor $Yellow
        Push-Location $dir.FullName
        try {
            yarn
            yarn cdeps
            Write-Host "$($dir.Name) processed successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to process $($dir.Name): $_" -ForegroundColor $Red
            Pop-Location
            exit 1
        }
        Pop-Location
    }
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "Dev build completed successfully!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

# Run test command (run all tests)
if ($Test) {
    Write-Host "`n---Testing libs---" -ForegroundColor $Cyan
    Push-Location $LibDir
    try {
        Write-Host "---Testing js---" -ForegroundColor $Yellow
        Push-Location "js"
        yarn test
        Pop-Location
        
        Write-Host "---Testing vue---" -ForegroundColor $Yellow
        Push-Location "vue"
        yarn test
        Pop-Location
        
        Write-Host "Libs tested successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to test libs: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n---Testing clients---" -ForegroundColor $Cyan
    $webSubdirs = Get-ChildItem -Path (Join-Path $ClientDir "web") -Directory
    foreach ($dir in $webSubdirs) {
        Write-Host "`n---Testing $($dir.Name)---" -ForegroundColor $Yellow
        Push-Location $dir.FullName
        try {
            yarn test
            Write-Host "$($dir.Name) tested successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to test $($dir.Name): $_" -ForegroundColor $Red
            Pop-Location
            exit 1
        }
        Pop-Location
    }
    
    Write-Host "`n---Testing server---" -ForegroundColor $Cyan
    Push-Location $ServerDir
    try {
        make test
        Write-Host "Server tested successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to test server: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "All tests completed successfully!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

# Run lint command (run all linting)
if ($Lint) {
    Write-Host "`n---Linting libs---" -ForegroundColor $Cyan
    Push-Location $LibDir
    try {
        Write-Host "---Linting js---" -ForegroundColor $Yellow
        Push-Location "js"
        yarn lint
        Pop-Location
        
        Write-Host "---Linting vue---" -ForegroundColor $Yellow
        Push-Location "vue"
        yarn lint
        Pop-Location
        
        Write-Host "Libs linted successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to lint libs: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n---Linting clients---" -ForegroundColor $Cyan
    $webSubdirs = Get-ChildItem -Path (Join-Path $ClientDir "web") -Directory
    foreach ($dir in $webSubdirs) {
        Write-Host "`n---Linting $($dir.Name)---" -ForegroundColor $Yellow
        Push-Location $dir.FullName
        try {
            yarn lint
            Write-Host "$($dir.Name) linted successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to lint $($dir.Name): $_" -ForegroundColor $Red
            Pop-Location
            exit 1
        }
        Pop-Location
    }
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "All linting completed successfully!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

# Run fresh command (clean and rebuild everything)
if ($Fresh) {
    Write-Host "`n---Fresh libs---" -ForegroundColor $Cyan
    Push-Location $LibDir
    try {
        Write-Host "---Fresh js---" -ForegroundColor $Yellow
        Push-Location "js"
        if (Test-Path "node_modules") { Remove-Item -Recurse -Force "node_modules" }
        if (Test-Path "yarn.lock") { Remove-Item -Force "yarn.lock" }
        yarn
        yarn unlink
        yarn link
        yarn build
        Pop-Location
        
        Write-Host "---Fresh vue---" -ForegroundColor $Yellow
        Push-Location "vue"
        if (Test-Path "node_modules") { Remove-Item -Recurse -Force "node_modules" }
        if (Test-Path "yarn.lock") { Remove-Item -Force "yarn.lock" }
        yarn
        yarn cdeps
        yarn build
        yarn unlink
        yarn link
        Pop-Location
        
        Write-Host "Libs freshed successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Failed to fresh libs: $_" -ForegroundColor $Red
        Pop-Location
        exit 1
    }
    Pop-Location
    
    Write-Host "`n---Fresh clients---" -ForegroundColor $Cyan
    $webSubdirs = Get-ChildItem -Path (Join-Path $ClientDir "web") -Directory
    foreach ($dir in $webSubdirs) {
        Write-Host "`n---Fresh $($dir.Name)---" -ForegroundColor $Yellow
        Push-Location $dir.FullName
        try {
            if (Test-Path "node_modules") { Remove-Item -Recurse -Force "node_modules" }
            if (Test-Path "yarn.lock") { Remove-Item -Force "yarn.lock" }
            yarn
            yarn cdeps
            Write-Host "$($dir.Name) freshed successfully" -ForegroundColor $Green
        } catch {
            Write-Host "Failed to fresh $($dir.Name): $_" -ForegroundColor $Red
            Pop-Location
            exit 1
        }
        Pop-Location
    }
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "Fresh build completed successfully!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

# Run audit command (run security audit)
if ($Audit) {
    Write-Host "`n---Auditing libs---" -ForegroundColor $Cyan
    Push-Location $LibDir
    try {
        Write-Host "---Auditing js---" -ForegroundColor $Yellow
        Push-Location "js"
        yarn audit
        Pop-Location
        
        Write-Host "---Auditing vue---" -ForegroundColor $Yellow
        Push-Location "vue"
        yarn audit
        Pop-Location
        
        Write-Host "Libs audited successfully" -ForegroundColor $Green
    } catch {
        Write-Host "Lib audit failed (continuing): $_" -ForegroundColor $Yellow
    }
    Pop-Location
    
    Write-Host "`n---Auditing clients---" -ForegroundColor $Cyan
    $webSubdirs = Get-ChildItem -Path (Join-Path $ClientDir "web") -Directory
    foreach ($dir in $webSubdirs) {
        Write-Host "`n---Auditing $($dir.Name)---" -ForegroundColor $Yellow
        Push-Location $dir.FullName
        try {
            yarn audit
            Write-Host "$($dir.Name) audited successfully" -ForegroundColor $Green
        } catch {
            Write-Host "$($dir.Name) audit failed (continuing): $_" -ForegroundColor $Yellow
        }
        Pop-Location
    }
    
    Write-Host "`n==============================================" -ForegroundColor $Green
    Write-Host "Security audit completed!" -ForegroundColor $Green
    Write-Host "==============================================" -ForegroundColor $Green
    exit 0
}

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
        Fresh-Client
        Write-Host "Client fresh completed successfully" -ForegroundColor $Green
        Build-Client
        Write-Host "Client build completed successfully" -ForegroundColor $Green
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
