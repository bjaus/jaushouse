
ENV_FILE ?= .env
ifneq ("$(wildcard $(ENV_FILE))", "")
  include $(ENV_FILE)
  export $(shell sed -n 's/=.*//' $(ENV_FILE))
endif

.PHONY: help
help: # List available targets
	@grep -h '^[a-zA-Z_-]\+:.*#' Makefile make/*.mk \
		|sed 's/:.*#/:/' \
		|awk -F':' '{printf "  %-20s %s\n", $$1, $$2}'

.PHONY: run
run: # Run the service
	@echo "Running the service..."
	@go run ./cmd/service

.PHONY: build
build: # Build the service
	@echo "Building the service..."
	@go build -o ./bin/jacked ./cmd/service
