# Tools.
TOOLS = ./tools
TOOLS_BIN = $(TOOLS)/bin

.PHONY: generate-swagger-api
generate-swagger-api:
	swag init -g router.go -d ./internal//apigateway/controller/http/v1/ --parseInternal --parseDependency

.PHONY: fix-lint
fix-lint: $(TOOLS_BIN)/golangci-lint
	$(TOOLS_BIN)/golangci-lint run --fix

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

