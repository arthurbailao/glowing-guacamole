#!make

GOPATH=$(shell go env GOPATH)

setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
	go get ./...

lint: setup
	$(GOPATH)/bin/golangci-lint run -c .golangci.yaml ./...

test: setup
	go test -v ./...

run:
	docker compose up

