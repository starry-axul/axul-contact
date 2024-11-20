.PHONY: build

build:
	@echo "=> Building service"
	dir=hello/get make build-file
	dir=hello/create make build-file
	sam build

local:
	sam local start-api

build-file:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(dir)/bootstrap		cmd/$(dir)/main.go