# Variables
PROJECT_NAME := college-appointment-system
BINARY_PATH := /home/pranith/Desktop/$(PROJECT_NAME)/cmd/server/main
GO_FILES := $(shell find . -name '*.go')

# Default target
all: deps build

# Install dependencies
deps:
	go mod download

# Build the application
build: $(GO_FILES)
	go build -o $(BINARY_PATH) ./cmd/server/main.go

# Run the application
run:
	go run ./cmd/server/main.go

# Clean up build artifacts
clean:
	rm -f $(BINARY_PATH)


# Format the code
fmt:
	go fmt ./...

.PHONY: all deps build run test clean lint fmt