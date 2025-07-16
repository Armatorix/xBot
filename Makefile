SHELL:=/bin/bash
.PHONY: run
run:
# run with .env, source .env and run go main
	@echo "Running xBot..."
	@source .env && go run main.go