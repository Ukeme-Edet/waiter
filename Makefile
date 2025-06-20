.DEFAULT_GOAL := build

.PHONY: build
fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet: fmt
	@echo "Running go vet..."
	@go vet ./...

build: vet
	@echo "Building the application..."
	@go build -o bin/app
	@echo "Build complete."

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@rm -f coverage.out
	@echo "Clean complete."
