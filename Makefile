#!/usr/bin/make
.DEFAULT_GOAL := help
.PHONY: help

DOCKER_COMPOSE ?= docker compose -f docker-compose.yml

export GOOS=linux
export GOARCH=amd64

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps: ## Install dependencies for protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps: ## Get dependencies for protobuf
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

fmt: ## Automatically format source code
	go fmt ./...
.PHONY:fmt

lint: fmt ## Check code (lint)
	golangci-lint run ./... --config .golangci.pipeline.yaml
.PHONY:lint

vet: fmt ## Check code (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

vet-shadow: fmt ## Check code with detect shadow (vet)
	go vet -vettool=$(which shadow) ./...
.PHONY:vet

build: ## Build service containers
	$(DOCKER_COMPOSE) build

up: vet ## Start services
	$(DOCKER_COMPOSE) up -d $(SERVICES)

down: ## Down services
	$(DOCKER_COMPOSE) down

clean: ## Delete all containers
	$(DOCKER_COMPOSE) down --remove-orphans

generate: ## Generate all API proto files
	make generate-auth-api
	make generate-access-api

generate-auth-api: ## Generate pb.go files for Auth API
	mkdir -p pkg/api/auth_v1
	protoc --proto_path=api/v1/pb \
    	--go_out=pkg/api/auth_v1 --go_opt=paths=source_relative \
    	--go-grpc_out=pkg/api/auth_v1 --go-grpc_opt=paths=source_relative \
    	--validate_out lang=go:pkg/api/auth_v1 --validate_opt=paths=source_relative \
    	auth.proto

generate-access-api: ## Generate pb.go files for Access API
	mkdir -p pkg/api/access_v1
	protoc --proto_path=api/v1/pb \
        	--go_out=pkg/api/access_v1 --go_opt=paths=source_relative \
        	--go-grpc_out=pkg/api/access_v1 --go-grpc_opt=paths=source_relative \
        	--validate_out lang=go:pkg/api/access_v1 --validate_opt=paths=source_relative \
        	access.proto

generate-secret-key:
	openssl rand -base64 30