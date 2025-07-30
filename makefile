.PHONY: build dev clean generate serve import help

# Default target
all: build

# Generate templ files and build the binary
build: generate
	go build -o bin/fancykaraoke .

# Generate templ templates
generate:
	go run github.com/a-h/templ/cmd/templ generate

# Development server with auto-reload
dev: generate
	go run . serve

# Import songs
import: build 
	./bin/fancykaraoke import 1000

# Serve the application
serve: build
	./bin/fancykaraoke serve

# Clean generated files and binary
clean:
	rm -rf bin/
	find . -name "*_templ.go" -delete

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test: build 
	go test ./...

# Format code
fmt:
	go fmt ./...
	go run github.com/a-h/templ/cmd/templ fmt .

# Lint code
lint: generate
	go vet ./...

# Create bin directory
bin:
	mkdir -p bin

# Help target
help:
	@echo "Available targets:"
	@echo "  build     - Generate templates and build binary"
	@echo "  dev       - Run development server"
	@echo "  generate  - Generate templ templates"
	@echo "  serve     - Build and serve the application"
	@echo "  import    - Run song import"
	@echo "  clean     - Remove generated files and binary"
	@echo "  deps      - Install dependencies"
	@echo "  test      - Run tests"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  help      - Show this help"
