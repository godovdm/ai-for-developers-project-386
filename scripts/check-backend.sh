#!/bin/sh
set -eu
cd backend
if [ -n "$(gofmt -l .)" ]; then
  echo 'Go files require gofmt' >&2
  gofmt -l .
  exit 1
fi
if [ -n "$(go run golang.org/x/tools/cmd/goimports@v0.35.0 -l .)" ]; then
  echo 'Go files require goimports' >&2
  go run golang.org/x/tools/cmd/goimports@v0.35.0 -l .
  exit 1
fi
golangci-lint run --enable revive --enable gocritic --enable errorlint --enable exhaustive --enable gosec --enable bodyclose --enable noctx --enable nilerr ./...
go build ./...
go test -race -count=1 ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...
