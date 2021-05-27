::[Bat To Exe Converter]
::
::YAwzoRdxOk+EWAjk
::fBw5plQjdCyDJGyX8VAjFA1dQAWGAES0A5EO4f7+r6fHoUwQRPcsNprf3fqLJPUvw0jlYZEoxUZzm8QCHydUcRGgIAY3pg4=
::YAwzuBVtJxjWCl3EqQJgSA==
::ZR4luwNxJguZRRnk
::Yhs/ulQjdF+5
::cxAkpRVqdFKZSjk=
::cBs/ulQjdF+5
::ZR41oxFsdFKZSDk=
::eBoioBt6dFKZSDk=
::cRo6pxp7LAbNWATEpCI=
::egkzugNsPRvcWATEpCI=
::dAsiuh18IRvcCxnZtBJQ
::cRYluBh/LU+EWAnk
::YxY4rhs+aU+JeA==
::cxY6rQJ7JhzQF1fEqQJQ
::ZQ05rAF9IBncCkqN+0xwdVs0
::ZQ05rAF9IAHYFVzEqQJQ
::eg0/rx1wNQPfEVWB+kM9LVsJDGQ=
::fBEirQZwNQPfEVWB+kM9LVsJDGQ=
::cRolqwZ3JBvQF1fEqQJQ
::dhA7uBVwLU+EWDk=
::YQ03rBFzNR3SWATElA==
::dhAmsQZ3MwfNWATElA==
::ZQ0/vhVqMQ3MEVWAtB9wSA==
::Zg8zqx1/OA3MEVWAtB9wSA==
::dhA7pRFwIByZRRnk
::Zh4grVQjdCyDJGyX8VAjFA1dQAWGAE+/Fb4I5/jH6+WEqUgPGeY7dpyW3rGYJewc+nnXYZc/wklpsPQ4GRVWex7laxcxyQ==
::YB416Ek+ZG8=
::
::
::978f952a14a936cc963da21a135fa983
@echo off

@REM set secret=%1

set /p user="Enter username: "
@REM set /p pw="Enter password: "
@REM set user=%2
@REM set pw=%3%
shift
shift

SETLOCAL ENABLEDELAYEDEXPANSION
SET count=1
FOR /F "tokens=* USEBACKQ" %%F IN (`ResolvePublicIP_client.bat "lluisAprovama" %user%`) DO (
  SET var!count!=%%F
  SET /a count=!count!+1
)

set addr=%var2%

echo %addr%
(plink -ssh -P 22 -l %user% -D 8022 %addr% -no-antispoof -share -N -i "pi-private.ppk")


ENDLOCAL

pause > nul