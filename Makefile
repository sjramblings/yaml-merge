# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=yaml-merge

# Directories
BIN_DIR=bin
DIST_DIR=dist

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS=-ldflags "\
	-s -w \
	-X main.version=${VERSION} \
	-X main.gitCommit=${GIT_COMMIT} \
	-X main.buildTime=${BUILD_TIME}"

# Cross compilation
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64 arm64

.PHONY: all build test clean help version release-major release-minor release-patch cross-build init-github-actions

all: test build

build: ensure-dirs ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/yaml-merge

ensure-dirs: ## Create necessary directories
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(DIST_DIR)

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)

coverage: ensure-dirs ## Run tests with coverage
	$(GOTEST) -coverprofile=$(BIN_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(BIN_DIR)/coverage.out -o $(BIN_DIR)/coverage.html

cross-build: clean ensure-dirs ## Build for multiple platforms
	$(foreach GOOS, $(PLATFORMS),\
		$(foreach GOARCH, $(ARCHITECTURES),\
			$(if $(filter $(GOOS)/$(GOARCH),windows/arm64), ,\
				GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,) ./cmd/yaml-merge; \
			) \
		) \
	)

# Version management
version: ## Display current version
	@echo "Current version: ${VERSION}"

bump-major: ## Bump major version (x.0.0)
	@echo "Bumping major version..."
	$(eval CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"))
	$(eval MAJOR_VERSION=$(shell echo $(CURRENT_VERSION) | cut -d. -f1 | tr -d v))
	$(eval NEW_VERSION="v$$(($(MAJOR_VERSION)+1)).0.0")
	@echo "New version: $(NEW_VERSION)"
	@echo "$(NEW_VERSION)" > .version

bump-minor: ## Bump minor version (0.x.0)
	@echo "Bumping minor version..."
	$(eval CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"))
	$(eval MAJOR_VERSION=$(shell echo $(CURRENT_VERSION) | cut -d. -f1))
	$(eval MINOR_VERSION=$(shell echo $(CURRENT_VERSION) | cut -d. -f2))
	$(eval NEW_VERSION="$(MAJOR_VERSION).$$(($(MINOR_VERSION)+1)).0")
	@echo "New version: $(NEW_VERSION)"
	@echo "$(NEW_VERSION)" > .version

bump-patch: ## Bump patch version (0.0.x)
	@echo "Bumping patch version..."
	$(eval CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"))
	$(eval MAJOR_MINOR_VERSION=$(shell echo $(CURRENT_VERSION) | cut -d. -f1,2))
	$(eval PATCH_VERSION=$(shell echo $(CURRENT_VERSION) | cut -d. -f3))
	$(eval NEW_VERSION="$(MAJOR_MINOR_VERSION).$$(($(PATCH_VERSION)+1))")
	@echo "New version: $(NEW_VERSION)"
	@echo "$(NEW_VERSION)" > .version

release-major: bump-major ## Create and push major version tag
	git tag -a $(shell cat .version) -m "Release $(shell cat .version)"
	git push origin $(shell cat .version)

release-minor: bump-minor ## Create and push minor version tag
	git tag -a $(shell cat .version) -m "Release $(shell cat .version)"
	git push origin $(shell cat .version)

release-patch: bump-patch ## Create and push patch version tag
	git tag -a $(shell cat .version) -m "Release $(shell cat .version)"
	git push origin $(shell cat .version)

# Development helpers
fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

lint: ## Run linter
	golangci-lint run

# Install development tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
