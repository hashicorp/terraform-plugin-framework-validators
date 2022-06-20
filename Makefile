default: build

build:
	go build -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

.PHONY: build lint fmt test
