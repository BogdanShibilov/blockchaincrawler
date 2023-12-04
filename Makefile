# Tools.
TOOLS = ./tools
TOOLS_BIN = $(TOOLS)/bin

generate-swagger-api:
	swag init -g router.go -d ./internal//apigateway/controller/http/v1/ --parseInternal --parseDependency
.PHONY: generate-swagger-api

fix-lint: $(TOOLS_BIN)/golangci-lint
	$(TOOLS_BIN)/golangci-lint run --fix
.PHONY: fix-lint

imports: $(TOOLS_BIN)/goimports
	$(TOOLS_BIN)/goimports -local "github.com/bogdanshibilov/blockchaincrawler" -w ./internal ./cmd

# INSTALL linter
$(TOOLS_BIN)/golangci-lint: export GOBIN = $(shell pwd)/$(TOOLS_BIN)
$(TOOLS_BIN)/golangci-lint:
	mkdir -p $(TOOLS_BIN)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2


# INSTALL goimports
$(TOOLS_BIN)/goimports: export GOBIN = $(shell pwd)/$(TOOLS_BIN)
$(TOOLS_BIN)/goimports:
	mkdir -p $(TOOLS_BIN)
	go install golang.org/x/tools/cmd/goimports@latest

generate-dockerimage-apigateway:
	docker build -t apigateway -f Dockerfile-apigateway .
.PHONY: generate-dockerimage-apigateway

generate-dockerimage-auth:
	docker build -t auth -f Dockerfile-auth .
.PHONY: generate-dockerimage-auth

generate-dockerimage-blockinfo:
	docker build -t blockinfo -f Dockerfile-blockinfo .
.PHONY: generate-dockerimage-blockinfo

generate-dockerimage-crawler:
	docker build -t crawler -f Dockerfile-crawler .
.PHONY: generate-dockerimage-crawler

generate-dockerimage-user:
	docker build -t user -f Dockerfile-user .
.PHONY: generate-dockerimage-user