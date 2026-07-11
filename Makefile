.PHONY: run build test vet fmt tidy check

run:
	go run ./cmd/octyne

build:
	go build ./cmd/octyne

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

check:
	go test ./...
	go vet ./...