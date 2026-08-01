BINARY_NAME=kis-trading-mcp-go
BUILD_DIR=../bin
PREFIX?=$(HOME)/.local
INSTALL_DIR=$(PREFIX)/bin
CMD_DIR=cmd

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build install uninstall build-all clean test run deps fmt lint help

all: build

## 빌드: ../bin 디렉토리에 컴파일 (배포용)
build:
	@echo "==> Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "    Output: $(BUILD_DIR)/$(BINARY_NAME)"

## 설치: 시스템 PATH에 설치 (PREFIX=/path 로 경로 변경 가능)
## 예: make install PREFIX=~/.local
install:
	@echo "==> Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "    Binary not found. Building first..."; \
		$(MAKE) build; \
	fi
	@install -d $(INSTALL_DIR)
	@install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "    Installed: $(INSTALL_DIR)/$(BINARY_NAME)"

## 설치 제거
uninstall:
	@echo "==> Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "    Removed: $(INSTALL_DIR)/$(BINARY_NAME)"

## 모든 커맨드 빌드
build-all:
	@echo "Building all commands..."
	@mkdir -p $(BUILD_DIR)
	@for cmd in $(CMD_DIR)/*/; do \
		name=$$(basename $$cmd); \
		echo "  Building $$name..."; \
		$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$$name ./$$cmd; \
	done

## 클린
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

## 테스트
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## 의존성 설치
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## 포맷
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## 린트 (golangci-lint 필요)
lint:
	@echo "Running linter..."
	golangci-lint run ./...

## 서버 실행
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

## 도움말
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build targets:"
	@echo "  make build       - 빌드: ../bin/ 에 컴파일 (배포용 아티팩트)"
	@echo "  make install     - 설치: ~/.local/bin 에 설치 (기본)"
	@echo "                     변경: make install PREFIX=/usr/local"
	@echo "  make uninstall   - 설치 제거"
	@echo "  make build-all   - 모든 cmd/ 하위 명령어 빌드"
	@echo ""
	@echo "Development targets:"
	@echo "  make test        - 테스트 실행"
	@echo "  make run         - 빌드 후 실행"
	@echo "  make fmt         - 코드 포맷"
	@echo "  make lint        - 린트 검사 (golangci-lint 필요)"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean       - 빌드 결과물 삭제"
	@echo "  make deps        - 의존성 다운로드"
	@echo "  make help        - 이 도움말 표시"
