:: Windows version of the Makefile commands
@ECHO OFF

SETLOCAL
IF "%1" == "dev" (
    SET GIT_HEAD=CALL 'git rev-parse HEAD'
    ECHO %GIT_HEAD%
    SET GOARCH=amd64
    go build -ldflags="-X github.com/pteropackages/soar/cmd.Version=%GIT_HEAD%" soar.go
) ELSE (
    SET GOARCH=amd64
    SET GOOS=linux
    go build -o build/soar_linux soar.go

    SET GOOS=windows
    go build -o build/soar_win32.exe soar.go
)
ENDLOCAL
