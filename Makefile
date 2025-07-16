SHELL:=/bin/bash
.PHONY: run
run:
# run with .env, source .env and run go main
	@echo "Running xBot..."
	@source .env && go run main.go

.PHONY: build
build:
	@echo "Building xBot..."
	@mkdir -p bin
	@go build -o ./bin/xbot  .
