@echo off
title FIFA Wonderkid Scout v2

echo.
echo  ============================================
echo   FIFA WONDERKID SCOUT v2.0
echo   Stack: Go + Fiber / MariaDB / Svelte
echo  ============================================
echo.

:: Matar procesos previos en los puertos
for /f "tokens=5" %%a in ('netstat -aon 2^>nul ^| findstr ":8080 "') do taskkill /f /PID %%a >nul 2>&1
for /f "tokens=5" %%a in ('netstat -aon 2^>nul ^| findstr ":5173 "') do taskkill /f /PID %%a >nul 2>&1

:: Backend Go (binario compilado)
echo [1/2] Iniciando backend Go + Fiber en :8080 ...
start "Backend Go :8080" cmd /k ""%~dp0backend\fifa-scout.exe""

timeout /t 3 /nobreak >nul

:: Frontend Svelte
echo [2/2] Iniciando frontend Svelte en :5173 ...
start "Frontend Svelte :5173" cmd /k "cd /d "%~dp0frontend" && npm.cmd run dev"

:: Esperar hasta que el frontend este listo (max 30 segundos)
echo [3/3] Esperando que el frontend este listo...
set /a tries=0
:WAIT_LOOP
timeout /t 2 /nobreak >nul
curl -s http://localhost:5173 >nul 2>&1
if %errorlevel%==0 goto OPEN_BROWSER
set /a tries+=1
if %tries% lss 15 goto WAIT_LOOP

:OPEN_BROWSER
:: Abrir navegador automaticamente
echo  Abriendo navegador...
start http://localhost:5173

echo.
echo  ============================================
echo   APP INICIADA
echo   Backend Go:   http://localhost:8080/api/players
echo   Frontend:     http://localhost:5173
echo  ============================================
echo.
pause
