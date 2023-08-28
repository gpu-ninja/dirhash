VERSION 0.7
FROM golang:1.20-bookworm
WORKDIR /app

all:
  COPY (+dirhash/dirhash --GOARCH=amd64) ./dist/dirhash-linux-amd64
  COPY (+dirhash/dirhash --GOARCH=arm64) ./dist/dirhash-linux-arm64
  COPY (+dirhash/dirhash --GOOS=darwin --GOARCH=amd64) ./dist/dirhash-darwin-amd64
  COPY (+dirhash/dirhash --GOOS=darwin --GOARCH=arm64) ./dist/dirhash-darwin-arm64
  RUN cd dist && find . -type f -exec sha256sum {} \; >> checksums.txt
  SAVE ARTIFACT ./dist/dirhash-linux-amd64 AS LOCAL dist/dirhash-linux-amd64
  SAVE ARTIFACT ./dist/dirhash-linux-arm64 AS LOCAL dist/dirhash-linux-arm64
  SAVE ARTIFACT ./dist/dirhash-darwin-amd64 AS LOCAL dist/dirhash-darwin-amd64
  SAVE ARTIFACT ./dist/dirhash-darwin-arm64 AS LOCAL dist/dirhash-darwin-arm64
  SAVE ARTIFACT ./dist/checksums.txt AS LOCAL dist/checksums.txt

dirhash:
  ARG GOOS=linux
  ARG GOARCH=amd64
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=0 go build --ldflags '-s' -o dirhash cmd/dirhash/main.go
  SAVE ARTIFACT ./dirhash AS LOCAL dist/dirhash-${GOOS}-${GOARCH}

tidy:
  LOCALLY
  RUN go mod tidy
  RUN go fmt ./...

lint:
  FROM golangci/golangci-lint:v1.54.2
  WORKDIR /app
  COPY . ./
  RUN golangci-lint run --timeout 5m ./...

test:
  COPY . ./
  RUN go test -coverprofile=coverage.out -v ./...
  SAVE ARTIFACT ./coverage.out AS LOCAL coverage.out