.PHONY: all
all:
	GOOS=linux GOARCH=amd64 go build -o build/soar_linux soar.go
	GOOS=windows GOARCH=amd64 go build -o build/soar_win32 soar.go

.PHONY: dev
dev:
	GIT_HEAD=$(shell git rev-parse HEAD | head -c8)
	go build -ldflags="-X github.com/pteropackages/soar/cmd.Version=$(GIT_HEAD)" -o build/ soar.go
