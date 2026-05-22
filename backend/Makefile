.PHONY: build test test-integration test-all

build:
	go build -o server ./cmd/server/

test:
	go test -v -count=1 ./...

test-integration:
	go test -tags=integration -p=1 -count=1 ./...
	go test -tags=integration -v -count=1 ./...

test-all: test test-integration
