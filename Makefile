# ÁGORA — build, test and lint targets.
#
# The linter version is pinned here and in .github/workflows/go-ci.yml. Keep the
# two in step: a local run that passes against a different version than CI is
# worse than no local run.
BINDIR          := bin
PKG             := ./...
GOLANGCI_LINT   := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2

.PHONY: all build test test-race vet lint fmt tidy cover clean help

all: build test lint ## Build, test and lint.

build: ## Build every command into bin/.
	go build -trimpath -o $(BINDIR)/ ./cmd/...

test: ## Run the test suite.
	go test $(PKG)

test-race: ## Run the test suite with the race detector, as CI does.
	go test $(PKG) -race -coverprofile=coverage.out -covermode=atomic

vet: ## Run go vet.
	go vet $(PKG)

lint: ## Run golangci-lint at the version CI uses.
	go run $(GOLANGCI_LINT) run $(PKG)

fmt: ## Format all Go sources.
	gofmt -w .

tidy: ## Tidy the module graph.
	go mod tidy

cover: test-race ## Show per-function coverage.
	go tool cover -func=coverage.out

clean: ## Remove build and test artifacts.
	rm -rf $(BINDIR) coverage.out

help: ## List the available targets.
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
