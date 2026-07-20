@echo off
echo =====================================
echo   ChatGo Build Script
echo =====================================
echo.

REM Build Free/Standard binary
echo [1/2] Building Standard binary...
go build -ldflags="-s -w" -o dist\chatgo-standard.exe .
if %ERRORLEVEL% NEQ 0 (
    echo FAILED!
    exit /b 1
)
echo       chatgo-standard.exe OK

REM Build Pro binary (requires pro/*.go files)
echo [2/2] Building Pro binary...
go build -tags pro -ldflags="-s -w" -o dist\chatgo-pro.exe .
if %ERRORLEVEL% NEQ 0 (
    echo Pro features not available (stub only).
    echo Run: go build -tags pro -ldflags="-s -w" -o dist\chatgo-pro.exe .
) else (
    echo       chatgo-pro.exe OK
)

echo.
echo =====================================
echo   Build Complete
echo =====================================
dir dist\*.exe 2>nul
