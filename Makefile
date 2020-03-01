.PHONY: build
build: ## Build nmcli Ansible module
	go build -mod=vendor -o nmcli-go-ansible cmd/nmcli-go/main.go

.PHONY: generate
generate: ## Generate test mocks
	go generate -mod=vendor

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: lint
lint: ## Lint source code
	gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`
	revive -config revive.toml -formatter stylish -exclude ./vendor/... ./...

.PHONY: test
test: ## Run unit tests
	go test -mod=vendor ./...

.PHONY: HELP
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

