@echo off
setlocal
for /f "delims=" %%i in ('hg parent --template "{rev}({node|short})"') do set Revision=%%i
echo Building rev %Revision%...
go build -o GoSogouProxy.exe -ldflags "-s -w -X main.Revision %Revision%"
go build -o GoSogouProxy-no-console.exe -ldflags "-H windowsgui -s -w -X main.Revision %Revision%"
echo Done.
endlocal
