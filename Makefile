# Confluent Kafka Go package tags:
GO_CKG_TAGS ?= dynamic

.PHONY: version deps test

version: ## Determine and output current version.
	@[[ -z "$(SOMEVAR)" ]] && git describe --tags --always | sed 's/v//' || echo "$(SOMEVAR)"

deps: ## Sync dependencies.
	go mod tidy

test: ## Run all unit tests.
	go test -tags unit,$(GO_CKG_TAGS) -count=1 -coverprofile coverage.out -covermode count -coverpkg=$(shell go list ./... | tr '\n' ',') -v ./...