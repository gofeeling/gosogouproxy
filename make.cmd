@echo off
setlocal
for /f "delims=" %%i in ('git log -1 --pretty^=format:"%%h %%ai"') do set Revision=%%i
set Revision=%Revision:~0,18%
echo Building rev %Revision%...
%~d0
cd %~p0
go generate
go install -ldflags "-s -w -X main.Revision """%Revision%""""
rem go build -o GoSogouProxy-no-console.exe -ldflags "-H windowsgui -s -w -X main.Revision """%Revision%""""
echo Done.
endlocal
