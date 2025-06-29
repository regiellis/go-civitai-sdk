@echo off
REM Civitai API Tester - Windows build script
setlocal enabledelayedexpansion

set APP_NAME=civitai-tester
set VERSION=1.0.0
set BUILD_DIR=builds

echo ======================================================================
echo WARNING: DEVELOPMENT/TESTING TOOL ONLY
echo ======================================================================
echo Building Civitai API Tester v%VERSION%...
echo WARNING: Do not use these binaries in production environments!
echo ======================================================================

REM Clean previous builds
if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

REM Build for different platforms
set platforms=windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

for %%p in (%platforms%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set GOOS=%%a
        set GOARCH=%%b
        
        set output_name=%APP_NAME%
        if "!GOOS!"=="windows" set output_name=!output_name!.exe
        
        set output_path=%BUILD_DIR%\%APP_NAME%-!GOOS!-!GOARCH!
        if "!GOOS!"=="windows" set output_path=!output_path!.exe
        
        echo Building for !GOOS!/!GOARCH!...
        set GOOS=!GOOS!
        set GOARCH=!GOARCH!
        go build -ldflags="-s -w" -o !output_path! .
        
        if !errorlevel! neq 0 (
            echo Error building for !GOOS!/!GOARCH!
            exit /b 1
        )
    )
)

echo.
echo Build completed successfully!
echo Binaries available in .\%BUILD_DIR%\
dir %BUILD_DIR%