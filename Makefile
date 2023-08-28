LOCALBIN ?= $(shell pwd)/bin

SRCS := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

$(LOCALBIN)/dirhash: $(LOCALBIN) $(SRCS)
	CGO_ENABLED=0 go build --ldflags '-s' -o $@ cmd/dirhash/main.go

tidy:
	go mod tidy
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test -coverprofile=coverage.out -v ./...

clean:
	-rm -rf bin
	go clean -testcache

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: tidy lint test clean