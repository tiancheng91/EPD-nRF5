# Makefile for pusher - 多平台构建
# 输出文件名格式: pusher_<os>-<arch>

# 项目名称
PROJECT_NAME := pusher

# 输出目录
OUTPUT_DIR := dist

# Go 构建参数
LDFLAGS := -s -w
CGO_ENABLED := 1

# 支持的平台列表
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 默认目标：构建当前平台
.PHONY: all
all: build

# 构建当前平台
.PHONY: build
build:
	@echo "构建当前平台..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)"

# 构建所有平台
.PHONY: build-all
build-all: $(foreach platform,$(PLATFORMS),build-$(subst /,-,$(platform)))

# 清理构建产物
.PHONY: clean
clean:
	@echo "清理构建产物..."
	@rm -rf $(OUTPUT_DIR)
	@echo "清理完成"

# 构建特定平台
.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "构建 darwin/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(PROJECT_NAME)_darwin-amd64 .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)_darwin-amd64"

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "构建 darwin/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(PROJECT_NAME)_darwin-arm64 .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)_darwin-arm64"

.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "构建 linux/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(PROJECT_NAME)_linux-amd64 .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)_linux-amd64"

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "构建 linux/arm64..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(PROJECT_NAME)_linux-arm64 .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)_linux-arm64"

.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "构建 windows/amd64..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(PROJECT_NAME)_windows-amd64.exe .
	@echo "构建完成: $(OUTPUT_DIR)/$(PROJECT_NAME)_windows-amd64.exe"

# 显示帮助信息
.PHONY: help
help:
	@echo "Pusher 构建系统"
	@echo ""
	@echo "用法:"
	@echo "  make [目标]"
	@echo ""
	@echo "目标:"
	@echo "  build          - 构建当前平台（默认）"
	@echo "  build-all      - 构建所有支持的平台"
	@echo "  build-darwin-amd64   - 构建 macOS (Intel)"
	@echo "  build-darwin-arm64   - 构建 macOS (Apple Silicon)"
	@echo "  build-linux-amd64    - 构建 Linux (AMD64)"
	@echo "  build-linux-arm64    - 构建 Linux (ARM64)"
	@echo "  build-windows-amd64  - 构建 Windows (64-bit)"
	@echo "  clean          - 清理构建产物"
	@echo "  help           - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build                    # 构建当前平台"
	@echo "  make build-all                # 构建所有平台"
	@echo "  make build-darwin-amd64       # 仅构建 macOS Intel 版本"
	@echo ""
	@echo "输出目录: $(OUTPUT_DIR)"

