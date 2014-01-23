@echo off
setlocal
for /f "delims=" %%i in ('hg parent --template "{rev}"') do set Revision=%%i
echo Building rev %Revision%...
go build -ldflags "-s -w -X main.Revision %Revision%"
echo Done.
endlocal