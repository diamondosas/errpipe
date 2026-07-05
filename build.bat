@echo off
setlocal

echo ============================================
echo          ERRPIPE BUILD SCRIPT
echo ============================================
echo.

:: Check Go is installed
where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH.
    exit /b 1
)

:: Install rsrc if not already available
where rsrc >nul 2>&1
if errorlevel 1 (
    echo [*] Installing rsrc tool for manifest embedding...
    go install github.com/akavel/rsrc@latest
    if errorlevel 1 (
        echo [ERROR] Failed to install rsrc.
        exit /b 1
    )
    echo [OK] rsrc installed.
) else (
    echo [OK] rsrc already available.
)

:: Generate the .syso resource file from the manifest
echo.
echo [*] Embedding Windows manifest into resource file...
rsrc -manifest errpipe.manifest -o errpipe.syso
if errorlevel 1 (
    echo [ERROR] Failed to generate errpipe.syso from manifest.
    exit /b 1
)
echo [OK] errpipe.syso generated.

:: Build the binary
echo.
echo [*] Building errpipe.exe...
go build -ldflags="-s -w" -o errpipe.exe .
if errorlevel 1 (
    echo [ERROR] Build failed.
    exit /b 1
)

echo.
echo ============================================
echo   Build successful! -> errpipe.exe
echo ============================================
endlocal
