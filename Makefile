.PHONY: version deps

version: ## Determine and output current version.
	@[[ -z "${VERSION}" ]] && git describe --tags --always | sed 's/v//' || echo $(VERSION)

deps: ## Sync dependencies.
	go mod tidy