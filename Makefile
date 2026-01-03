GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
SERVICE_NAME=medical-service

ifeq ($(GOHOSTOS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
endif

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: wire
wire:
	wire ./cmd

.PHONY: build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$(SERVICE_NAME) ./cmd

.PHONY: run
run:
	go run ./cmd -conf ./configs

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...
	gofmt -s -w .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	find . -name "wire_gen.go" -delete

.PHONY: dev
dev: wire run

.PHONY: check
check: fmt test

.PHONY: help
help:
	@echo 'Medical Service - Makefile'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@echo '  init     Install development tools'
	@echo '  wire     Generate Wire DI files'
	@echo '  build    Build binary'
	@echo '  run      Run service'
	@echo '  test     Run tests'
	@echo '  fmt      Format code'
	@echo '  tidy     Tidy Go modules'
	@echo '  clean    Clean build artifacts'
	@echo '  dev      Wire + Run (recommended for development)'
	@echo '  check    Format + Test (run before commit)'
	@echo ''

.DEFAULT_GOAL := help
