# Makefile for ham-radio-assistant

# Variables
APP_NAME = ham-radio-assistant
MAIN_PATH = ./cmd/ham-radio-assistant
BUILD_DIR = ./bin
DOCKER_TAG = ham-radio-assistant:latest

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	@go test -v ./...

# Run go mod tidy to clean up dependencies
.PHONY: tidy
tidy:
	@go mod tidy

# Build Docker image
.PHONY: docker
docker:
	@docker build -t $(DOCKER_TAG) .

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)

# Run the application
.PHONY: run
run: build
	@$(BUILD_DIR)/$(APP_NAME)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make         : Build the application"
	@echo "  make build   : Same as above"
	@echo "  make test    : Run tests"
	@echo "  make tidy    : Run go mod tidy to clean up dependencies"
	@echo "  make docker  : Build Docker image"
	@echo "  make clean   : Remove build artifacts"
	@echo "  make run     : Build and run the application"
	@echo "  make help    : Show this help message"
