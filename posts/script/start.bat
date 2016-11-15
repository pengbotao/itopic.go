@echo off

echo Stopping nginx...
taskkill /F /IM nginx.exe > nul

echo Stopping PHP FastCGI...
taskkill /F /IM php-cgi.exe > nul

echo Stopping Redis...
taskkill /F /IM redis-server.exe > nul


set path_php=C:/soft/php-7.0.4-nts-Win32-VC14-x64/
set path_nginx=C:/soft/openresty-1.9.7.3-win32/
set path_redis=C:/soft/Redis-x64-2.8.2400/


echo =============================================================================
echo Start PHP FastCGI...
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9000 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9001 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9002 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9003 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9004 -c %path_php%php.ini
RunHiddenConsole.exe %path_php%php-cgi.exe -b 127.0.0.1:9005 -c %path_php%php.ini

echo Start nginx...
cd /d %path_nginx%
%path_nginx%nginx.exe -t
start  %path_nginx%nginx.exe -c %path_nginx%conf/nginx.conf

echo Start Redis
cd /d %path_redis%
RunHiddenConsole.exe %path_redis%redis-server.exe %path_redis%redis.windows.conf

echo =============================================================================

tasklist /nh^ | findstr /i /s /c:"nginx.exe"

tasklist /nh^ | findstr /i /s /c:"php-cgi.exe"

tasklist /nh^ | findstr /i /s /c:"redis-server.exe"

pause