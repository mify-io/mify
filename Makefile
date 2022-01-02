all: build lint

build:
	go build ./cmd/mify

lint:
	go vet ./...
	staticcheck ./...
