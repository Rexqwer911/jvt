.PHONY: build test clean install run help

# Build variables
BINARY_NAME=jvt.exe
BUILD_DIR=build
MAIN_PATH=cmd/jvt/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

run: ## Run the application
	@$(GOCMD) run $(MAIN_PATH)

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

release: ## Build release version
	@echo "Building release..."
	@mkdir -p dist
	@$(GOBUILD) -ldflags="-s -w" -o dist/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Release build complete: dist/$(BINARY_NAME)"

choco-pack: release ## Create Chocolatey package
	@echo "Creating Chocolatey package..."
	@cd chocolatey && choco pack
	@echo "Chocolatey package created"

installer: build ## Build Windows Installer (requires Inno Setup)
	@echo "Building Windows Installer..."
	@iscc installer/jvt.iss
	@echo "Installer created: installer/Output/jvt-setup.exe"
