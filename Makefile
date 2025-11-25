APP_NAME=review-assign-service
MAIN=./cmd/app
BIN=./bin/$(APP_NAME)
ENV_FILE=.env.example

.PHONY: build run fmt lint

build:
	go build -o $(BIN) $(MAIN)

run:
	go run $(MAIN)

fmt:
	gofmt -s -w .
	go vet ./...

lint:
	golangci-lint run