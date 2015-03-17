@echo off
setlocal
for /f "delims=" %%i in ('git rev-parse --short HEAD') do set Revision=%%i
echo Building rev %Revision%...
%~d0
cd %~p0
go install -ldflags "-s -w -X main.Revision %Revision%"
rem go build -o GoSogouProxy-no-console.exe -ldflags "-H windowsgui -s -w -X main.Revision %Revision%"
echo Done.
endlocal
