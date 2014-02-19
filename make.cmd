@echo off
setlocal
for /f "delims=" %%i in ('hg parent --template "{rev}"') do set Revision=%%i
echo Building rev %Revision%...
go build -o gosogouproxy-console.exe -ldflags "-s -w -X main.Revision %Revision% -X main.Log gosogouproxy.log"
go build -o gosogouproxy.exe -ldflags "-H windowsgui -s -w -X main.Revision %Revision% -X main.Log gosogouproxy.log"
echo Done.
endlocal