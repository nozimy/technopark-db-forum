build:
	go build -v -o ./technopark-db-forum ./cmd/technopark-db-forum

run:
	go run ./cmd/technopark-db-forum/main.go

test:
	go test -v -cover -race -timeout 30s ./...

.DEFAULT_GOAL := build

.PHONY: build run test
