@echo off
title Compilando Go...
echo Compilando backend Go...
cd /d "%~dp0backend"
"C:\Program Files\Go\bin\go.exe" build -o fifa-scout.exe .
if %errorlevel% == 0 (
    echo [OK] Compilado exitosamente: backend\fifa-scout.exe
) else (
    echo [ERROR] Fallo la compilacion. Revisa los errores arriba.
)
pause
