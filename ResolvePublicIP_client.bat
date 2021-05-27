echo off

set secret=%1
set user=%2
shift
shift

SETLOCAL ENABLEDELAYEDEXPANSION
SET count=1
FOR /F "tokens=* USEBACKQ" %%F IN (`curl --location --request GET "lmbfao.ddns.net:8053/resolve?secret=%secret%&domain=%user%"`) DO (
  SET var!count!=%%F
  SET /a count=!count!+1
)


(for /f "tokens=2,* delims=," %%a in ("%var1%") do set res=%%a)
(for /f "tokens=2,* delims=:" %%a in ("%res%") do set res=%%a)

set res=%res:"=%

ECHO %res%

ENDLOCAL