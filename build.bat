@echo off
setlocal

echo === Building Fin Compiler ===
go build -o fin.exe ./cmd/fin
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go build failed!
    exit /b 1
)
echo fin.exe built successfully.

echo.
echo === Creating Installer ===
cd scripts
makensis fin_installer.nsi
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: makensis failed!
    cd ..
    exit /b 1
)
cd ..

echo.
echo === Build Complete ===
echo Installer: scripts\Fin-v1.0.0-Setup.exe
endlocal