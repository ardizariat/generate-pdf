.PHONY: build run test fmt lint race vuln

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

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
