.PHONY: build-windows build-windows-installer build-linux build-macos clean deps

# Build variables
BINARY_NAME=jvt.exe
BUILD_DIR=build
MAIN_PATH=cmd/jvt/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

build-windows: ## Build the binary for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-windows-installer: build ## Build Windows Installer (requires Inno Setup)
	@echo "Building Windows Installer..."
	@iscc installer/jvt.iss
	@echo "Installer created: installer/Output/jvt-setup.exe"

build-linux: ## Build the binary for Linux
	@echo "Building jvt-linux..."
	@set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) -o $(BUILD_DIR)/jvt-linux $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/jvt-linux"

build-macos: ## Build the binary for macOS (arm64 and x86_64 separated)
	@echo "Building jvt-macos-amd64..."
	@set GOOS=darwin&& set GOARCH=amd64&& $(GOBUILD) -o $(BUILD_DIR)/jvt-macos-amd64 $(MAIN_PATH)
	@echo "Building jvt-macos-arm64..."
	@set GOOS=darwin&& set GOARCH=arm64&& $(GOBUILD) -o $(BUILD_DIR)/jvt-macos-arm64 $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/jvt-macos-*"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
