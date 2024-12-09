# Include variables
include .envrc

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage':
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# =============================================================================== #
# DEVELOPMENT
# =============================================================================== #

## run/api: Run the application
.PHONY: run/api
run/api:
	go run ./cmd/main