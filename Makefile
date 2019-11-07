.PHONY: build
build:
	go build -v ./cmd/technopark-db-forum

.DEFAULT_GOAL := build