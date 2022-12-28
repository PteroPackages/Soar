GIT_HEAD=$(shell git rev-parse HEAD | head -c8)

.PHONY: all
all:
	GOOS=linux GOARCH=amd64 go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=$(GIT_HEAD)" -o build/soar_linux soar.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=$(GIT_HEAD)" -o build/soar_win32 soar.go

.PHONY: dev
dev:
	go build -ldflags="-X github.com/pteropackages/soar/cmd.Build=$(GIT_HEAD)" soar.go
