:: Windows version of the Makefile commands
@ECHO OFF

SETLOCAL
SET GOARCH=amd64
SET GIT_HEAD=(CALL git rev-parse HEAD)
IF "%1" == "dev" (
    go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=%GIT_HEAD%" soar.go
) ELSE (
    SET GOOS=linux
    go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=%GIT_HEAD%" -o build/soar_linux soar.go

    SET GOOS=windows
    go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=%GIT_HEAD%" -o build/soar_win32.exe soar.go
)
ENDLOCAL
