# Variables
APP_NAME = simplegowebserver-linux-x86_64

# Default target
.PHONY: all
all: build

# Build the Go application
.PHONY: build
build:
	@echo "Building the Go application..."
	@go build -o $(APP_NAME) main.go

# Clean up the build and Docker artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(APP_NAME)

# Display help
.PHONY: help
help:
	@echo "Makefile targets:"
	@echo "  build          Build the Go application"
	@echo "  clean          Clean up the build and Docker artifacts"
	@echo "  help           Display this help message"

