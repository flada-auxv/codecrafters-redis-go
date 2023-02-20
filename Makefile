.DEFAULT_GOAL := build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build-server: vet
	go mod tidy
	go build -o redis-server codecrafters-redis-go/cmd/server

build-cli: vet
	go mod tidy
	go build -o redis-cli codecrafters-redis-go/cmd/client
