# Makefile with Cross-Platform Build

## Complete Makefile

```makefile
.PHONY: help build build-web build-backend dev-web dev-backend clean \
        build-linux build-linux-arm64 build-macos build-macos-arm64 build-windows build-all

# Default target
.DEFAULT_GOAL := help

# Color output
GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
NC     := \033[0m

# Project configuration - MODIFY THESE FOR YOUR PROJECT
BINARY_NAME := app
CMD_PATH    := ./cmd/$(BINARY_NAME)/main.go
OUTPUT_DIR  := bin
LDFLAGS     := -s -w

help: ## Show help
	@printf "$(BLUE)$(BINARY_NAME) - Available Commands:$(NC)\n"
	@printf "\n$(YELLOW)Basic Commands:$(NC)\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -v "build-linux\|build-macos\|build-windows\|build-all" | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@printf "\n$(YELLOW)Cross-Platform Build:$(NC)\n"
	@grep -E '^build-(linux|macos|windows|all).*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

install-web: ## Install frontend dependencies
	@printf "$(BLUE)[Installing Frontend Dependencies]$(NC)\n"
	@cd site && npm install

build-web: ## Build frontend
	@printf "$(BLUE)[Building Frontend]$(NC)\n"
	@cd site && npm run build
	@printf "$(GREEN)[Success] Frontend built$(NC)\n"

build-backend: ## Build backend (current platform)
	@printf "$(BLUE)[Building Backend]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)$(NC)\n"

build: build-web build-backend ## Build frontend and backend

dev-web: ## Start frontend dev server (http://localhost:5173)
	@printf "$(BLUE)[Starting Frontend Dev Server]$(NC)\n"
	@printf "$(YELLOW)URL: http://localhost:5173$(NC)\n"
	@printf "$(YELLOW)API Proxy: http://localhost:8080$(NC)\n"
	@cd site && npm run dev

dev-backend: ## Start backend dev server (http://localhost:8080)
	@printf "$(BLUE)[Starting Backend Dev Server]$(NC)\n"
	@printf "$(YELLOW)URL: http://localhost:8080$(NC)\n"
	@go run $(CMD_PATH)

dev: ## Show instructions for running both servers
	@printf "$(YELLOW)Please run in two terminals:$(NC)\n"
	@printf "  Terminal 1: $(GREEN)make dev-backend$(NC)\n"
	@printf "  Terminal 2: $(GREEN)make dev-web$(NC)\n"

run: ## Run compiled binary
	@printf "$(BLUE)[Running Service]$(NC)\n"
	@$(OUTPUT_DIR)/$(BINARY_NAME)

clean: ## Clean build artifacts
	@printf "$(BLUE)[Cleaning Build Artifacts]$(NC)\n"
	@rm -rf site/dist $(OUTPUT_DIR)/$(BINARY_NAME) $(OUTPUT_DIR)/$(BINARY_NAME)-* $(OUTPUT_DIR)/*.exe
	@printf "$(GREEN)[Success] Cleaned$(NC)\n"

# ============================================================================
# Cross-Platform Build Targets
# ============================================================================

build-linux: build-web ## Build for Linux amd64
	@printf "$(BLUE)[Building Linux amd64]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64$(NC)\n"

build-linux-arm64: build-web ## Build for Linux arm64
	@printf "$(BLUE)[Building Linux arm64]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64$(NC)\n"

build-macos: build-web ## Build for macOS amd64 (Intel)
	@printf "$(BLUE)[Building macOS amd64]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64$(NC)\n"

build-macos-arm64: build-web ## Build for macOS arm64 (Apple Silicon)
	@printf "$(BLUE)[Building macOS arm64]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64$(NC)\n"

build-windows: build-web ## Build for Windows amd64
	@printf "$(BLUE)[Building Windows amd64]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@printf "$(GREEN)[Success] $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe$(NC)\n"

build-all: build-web ## Build for all platforms
	@printf "$(BLUE)[Building All Platforms]$(NC)\n"
	@mkdir -p $(OUTPUT_DIR)
	@printf "  $(YELLOW)→ Linux amd64$(NC)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@printf "  $(YELLOW)→ Linux arm64$(NC)\n"
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	@printf "  $(YELLOW)→ macOS amd64 (Intel)$(NC)\n"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@printf "  $(YELLOW)→ macOS arm64 (Apple Silicon)$(NC)\n"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	@printf "  $(YELLOW)→ Windows amd64$(NC)\n"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@printf "$(GREEN)[Success] All platforms built$(NC)\n"
	@printf "\n$(BLUE)Generated files:$(NC)\n"
	@ls -lh $(OUTPUT_DIR)/$(BINARY_NAME)-* 2>/dev/null | awk '{printf "  %s  %s\n", $$5, $$9}'
```

## Configuration Variables

Modify these at the top of Makefile for your project:

```makefile
BINARY_NAME := app              # Your binary name
CMD_PATH    := ./cmd/app/main.go # Path to main.go
OUTPUT_DIR  := bin              # Output directory
LDFLAGS     := -s -w            # Linker flags (-s -w for smaller binary)
```

## Build Targets

### Basic Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build` | Build frontend + backend |
| `make build-web` | Build frontend only |
| `make build-backend` | Build backend (current platform) |
| `make dev-web` | Start Vite dev server |
| `make dev-backend` | Start Go with hot reload |
| `make run` | Run compiled binary |
| `make clean` | Remove build artifacts |

### Cross-Platform Builds

| Command | Output |
|---------|--------|
| `make build-linux` | `bin/app-linux-amd64` |
| `make build-linux-arm64` | `bin/app-linux-arm64` |
| `make build-macos` | `bin/app-darwin-amd64` |
| `make build-macos-arm64` | `bin/app-darwin-arm64` |
| `make build-windows` | `bin/app-windows-amd64.exe` |
| `make build-all` | All of the above |

## LDFLAGS Explained

```makefile
LDFLAGS := -s -w
```

| Flag | Effect | Size Reduction |
|------|--------|----------------|
| `-s` | Strip symbol table | ~20% smaller |
| `-w` | Strip DWARF debug info | ~10% smaller |

## Build Output

```bash
$ make build-all
[Building All Platforms]
  → Linux amd64
  → Linux arm64
  → macOS amd64 (Intel)
  → macOS arm64 (Apple Silicon)
  → Windows amd64
[Success] All platforms built

Generated files:
  12M  bin/app-linux-amd64
  11M  bin/app-linux-arm64
  12M  bin/app-darwin-amd64
  11M  bin/app-darwin-arm64
  12M  bin/app-windows-amd64.exe
```
