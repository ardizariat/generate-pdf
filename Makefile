.PHONY: build run test fmt lint race

build:
	go build ./...
	go build -o pdfgen ./cmd/pdfgen

run:
	go run ./cmd/pdfgen

test:
	go test ./... -count=1

fmt:
	gofmt -w .

lint:
	golangci-lint run

race:
	go test -race ./... -count=1
