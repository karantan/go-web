# Convenience makefile to build the dev env and run common commands

VERSION := 0.1.0

# Build the package
.PHONY: build
build: elm
	@go build -ldflags="-X 'goweb/version.Version=$(VERSION)'"
# Run the server
.PHONY: run
run:
	@go run main.go

# Run tests
.PHONY: test
test:
	@go test -v ./...

# Create & populate database
.PHONY: db
db:
	@docker compose up db
