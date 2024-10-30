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
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
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
	go build -ldflags " -s -w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME}" -o bin/yaml-merge ./cmd/yaml-merge

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
	@if [ -z "$$(git tag)" ]; then \
		echo "v0.1.0" > .version; \
	else \
		CURRENT_VERSION=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
		MAJOR=$$(echo $$CURRENT_VERSION | cut -d. -f1); \
		MINOR=$$(echo $$CURRENT_VERSION | cut -d. -f2); \
		PATCH=$$(echo $$CURRENT_VERSION | cut -d. -f3); \
		NEW_MINOR=$$((MINOR + 1)); \
		NEW_VERSION="$$MAJOR.$$NEW_MINOR.0"; \
		echo "$$NEW_VERSION" > .version; \
	fi
	@NEW_VERSION=$$(cat .version); \
	echo "New version: $$NEW_VERSION"

bump-patch: ## Bump patch version (0.0.x)
	@echo "Bumping patch version..."
	@if [ -z "$$(git tag)" ]; then \
		echo "v0.0.1" > .version; \
	else \
		CURRENT_VERSION=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
		MAJOR=$$(echo $$CURRENT_VERSION | cut -d. -f1); \
		MINOR=$$(echo $$CURRENT_VERSION | cut -d. -f2); \
		PATCH=$$(echo $$CURRENT_VERSION | cut -d. -f3); \
		NEW_PATCH=$$((PATCH + 1)); \
		NEW_VERSION="$$MAJOR.$$MINOR.$$NEW_PATCH"; \
		echo "$$NEW_VERSION" > .version; \
	fi
	@NEW_VERSION=$$(cat .version); \
	echo "New version: $$NEW_VERSION"

release-major: bump-major ## Create and push major version tag
	git tag -a $(shell cat .version) -m "Release $(shell cat .version)"
	git push origin $(shell cat .version)

release-minor: bump-minor ## Create and push minor version tag
	git tag -a $(shell cat .version) -m "Release $(shell cat .version)"
	git push origin $(shell cat .version)

release-patch: bump-patch ## Create and push patch version tag
	@NEW_VERSION=$$(cat .version); \
	git tag -a "$$NEW_VERSION" -m "Release $$NEW_VERSION"; \
	git push origin "$$NEW_VERSION"

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

# Add these new targets to your existing Makefile
create-branch: ## Create a new feature branch
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then \
		echo "Error: Please checkout main branch first"; \
		echo "Run: git checkout main"; \
		exit 1; \
	fi; \
	read -p "Enter branch type (feature/fix/docs/deps): " type; \
	read -p "Enter branch description: " desc; \
	BRANCH_NAME="$$type/$$desc"; \
	BRANCH_NAME=$$(echo "$$BRANCH_NAME" | tr ' ' '-'); \
	git fetch origin main; \
	git checkout -b "$$BRANCH_NAME" origin/main

create-pr: ## Create a pull request with appropriate labels
	@if [ -z "$(title)" ]; then \
		echo "Error: Please provide a PR title using 'make create-pr title=\"Your PR title\"'"; \
		exit 1; \
	fi; \
	BRANCH_NAME=$$(git rev-parse --abbrev-ref HEAD); \
	LABELS=""; \
	if echo "$$BRANCH_NAME" | grep -q "^feature\|feature/"; then \
		LABELS="--label enhancement"; \
	elif echo "$$BRANCH_NAME" | grep -q "^fix\|fix/"; then \
		LABELS="--label bug"; \
	elif echo "$$BRANCH_NAME" | grep -q "^docs\|docs/"; then \
		LABELS="--label documentation"; \
	elif echo "$$BRANCH_NAME" | grep -q "^deps\|deps/"; then \
		LABELS="--label dependencies"; \
	fi; \
	if [ -n "$$LABELS" ]; then \
		gh pr create --title "$(title)" --body-file .github/pull_request_template.md $$LABELS; \
	else \
		echo "Warning: Branch name '$$BRANCH_NAME' doesn't match expected patterns. Creating PR without labels."; \
		gh pr create --title "$(title)" --body-file .github/pull_request_template.md; \
	fi

update-pr: ## Update pull request with appropriate labels
	@if [ -z "$(title)" ]; then \
		echo "Error: Please provide a PR title using 'make update-pr title=\"Your PR title\"'"; \
		exit 1; \
	fi; \
	BRANCH_NAME=$$(git rev-parse --abbrev-ref HEAD); \
	LABELS=""; \
	if echo "$$BRANCH_NAME" | grep -q "^feature\|feature/"; then \
		LABELS="--label enhancement"; \
	elif echo "$$BRANCH_NAME" | grep -q "^fix\|fix/"; then \
		LABELS="--label bug"; \
	elif echo "$$BRANCH_NAME" | grep -q "^docs\|docs/"; then \
		LABELS="--label documentation"; \
	elif echo "$$BRANCH_NAME" | grep -q "^deps\|deps/"; then \
		LABELS="--label dependencies"; \
	fi; \
	if [ -n "$$LABELS" ]; then \
		gh pr edit --title "$(title)" --body-file .github/pull_request_template.md $$LABELS; \
	else \
		echo "Warning: Branch name '$$BRANCH_NAME' doesn't match expected patterns. Updating PR without labels."; \
		gh pr edit --title "$(title)" --body-file .github/pull_request_template.md; \
	fi

check-pr: ## Run pre-PR checks
	@make fmt
	@make vet
	@make lint
	@make test

# Also add a new target for switching branches
switch-branch: ## Switch to a new branch from main
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" = "main" ]; then \
		read -p "Enter branch type (feature/fix/docs/deps): " type; \
		read -p "Enter branch description: " desc; \
		BRANCH_NAME="$$type/$$desc"; \
		BRANCH_NAME=$$(echo "$$BRANCH_NAME" | tr ' ' '-'); \
		git checkout -b "$$BRANCH_NAME"; \
	else \
		echo "Current branch: $$(git rev-parse --abbrev-ref HEAD)"; \
		echo "First run: git checkout main"; \
		exit 1; \
	fi
