@echo off
echo Stopping nginx...
taskkill /F /IM nginx.exe > nul
echo Stopping PHP FastCGI...
taskkill /F /IM php-cgi.exe > nul
echo Stopping Redis...
taskkill /F /IM redis-server.exe > nul

tasklist /nh^ | findstr /i /s /c:"nginx.exe"

tasklist /nh^ | findstr /i /s /c:"php-cgi.exe"

tasklist /nh^ | findstr /i /s /c:"redis-server.exe"

pause