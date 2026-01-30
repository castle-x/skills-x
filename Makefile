.PHONY: build clean install test build-all build-npm

# 项目信息
BINARY_NAME=skills-x
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go 命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# 目录
CMD_DIR=cmd/skills-x
OUT_DIR=bin
SKILLS_DATA=$(CMD_DIR)/skills/skills

# 默认目标
all: build

# 同步自研 skills 数据到 embed 目录
sync-skills:
	@rm -rf $(SKILLS_DATA)
	@cp -r skills $(SKILLS_DATA)
	@echo "Synced skills -> $(SKILLS_DATA)"

# 构建
build: sync-skills
	@mkdir -p $(OUT_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@rm -rf $(SKILLS_DATA)
	@echo "Built: $(OUT_DIR)/$(BINARY_NAME)"

# 安装到 GOPATH/bin
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Installed: $(GOPATH)/bin/$(BINARY_NAME)"

# 本地安装到 ~/.local/bin
install-local:
	@mkdir -p ~/.local/bin
	$(GOBUILD) $(LDFLAGS) -o ~/.local/bin/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Installed: ~/.local/bin/$(BINARY_NAME)"

# 清理
clean:
	@rm -rf $(OUT_DIR)
	@rm -rf /tmp/skills-*
	@echo "Cleaned"

# 测试
test:
	$(GOTEST) -v ./...

# 依赖
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 跨平台构建
build-all: sync-skills build-linux build-darwin build-windows
	@rm -rf $(SKILLS_DATA)
	@echo "All platforms built -> $(OUT_DIR)/"

# 跨平台构建并输出到 npm/bin (用于 npm 发布)
# 从 npm/package.json 读取版本号，确保二进制版本与 npm 版本一致
NPM_VERSION=$(shell grep '"version"' npm/package.json | sed 's/.*"version": "\(.*\)".*/\1/')
NPM_LDFLAGS=-ldflags "-X main.Version=$(NPM_VERSION) -X main.BuildTime=$(BUILD_TIME)"

build-npm: sync-skills
	@mkdir -p npm/bin
	@echo "Building version: $(NPM_VERSION)"
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(NPM_LDFLAGS) -o npm/bin/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(NPM_LDFLAGS) -o npm/bin/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(NPM_LDFLAGS) -o npm/bin/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(NPM_LDFLAGS) -o npm/bin/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(NPM_LDFLAGS) -o npm/bin/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@rm -rf $(SKILLS_DATA)
	@echo "All platforms built -> npm/bin/ (v$(NPM_VERSION))"
	@ls -lh npm/bin/

build-linux:
	@mkdir -p $(OUT_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

build-darwin:
	@mkdir -p $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

build-windows:
	@mkdir -p $(OUT_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

# 运行
run:
	$(GOBUILD) $(LDFLAGS) -o $(OUT_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	./$(OUT_DIR)/$(BINARY_NAME) $(ARGS)
