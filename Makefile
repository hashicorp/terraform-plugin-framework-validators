default: build

build:
	go build -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

# Generate copywrite headers
generate:
	cd tools; go generate ./...

.PHONY: build lint fmt test
