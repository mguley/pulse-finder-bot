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

## run/auth-grpc-client: Run the Auth gRPC client
.PHONY: run/auth-grpc-client
run/auth-grpc-client:
	go run ./cmd/grpc/auth

# =============================================================================== #
# QUALITY CONTROL
# =============================================================================== #

## install/goimports: Install goimports for formatting
.PHONY: install/goimports
install/goimports:
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest

## install/linter: Install GolangCI-Lint
.PHONY: install/linter
install/linter:
	@echo "Installing GolangCI-Lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

## lint: Run linter on all Go files in each module directory
.PHONY: lint
lint: install/linter
	@echo "Running GolangCI-Lint on all Go files in each module directory..."
	@find ./application ./cmd ./domain ./infrastructure -name '*.go' -exec dirname {} \; | sort -u | xargs $(shell go env GOPATH)/bin/golangci-lint run

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Tidying root module dependencies...'
	(cd ./ && go mod tidy)
	@echo 'Verifying root module dependencies...'
	(cd ./ && go mod verify)

	@echo 'Tidying application module dependencies...'
	(cd ./application && go mod tidy)
	@echo 'Verifying application module dependencies...'
	(cd ./application && go mod verify)

	@echo 'Tidying cmd module dependencies...'
	(cd ./cmd && go mod tidy)
	@echo 'Verifying cmd module dependencies...'
	(cd ./cmd && go mod verify)

	@echo 'Tidying domain module dependencies...'
	(cd ./domain && go mod tidy)
	@echo 'Verifying domain module dependencies...'
	(cd ./domain && go mod verify)

	@echo 'Tidying infrastructure module dependencies...'
	(cd ./infrastructure && go mod tidy)
	@echo 'Verifying infrastructure module dependencies...'
	(cd ./infrastructure && go mod verify)

	@echo 'Tidying tests module dependencies...'
	(cd ./tests && go mod tidy)
	@echo 'Verifying cmd module dependencies...'
	(cd ./tests && go mod verify)

	@echo 'Vendoring workspace dependencies...'
	go work vendor

# =============================================================================== #
# TESTING
# =============================================================================== #

## test/integration: Run integration tests (uses Go's caching mechanism)
.PHONY: test/integration
test/integration:
	@echo 'Running integration tests (with caching, sequentially)...'
	go test -v -p=1 ./tests/integration/...

## test/integration/no-cache: Run integration tests (bypass cache)
.PHONY: test/integration/no-cache
test/integration/no-cache:
	@echo 'Running integration tests (no cache, sequentially)...'
	go test -v -count=1 -p=1 ./tests/integration/...

# =============================================================================== #
# BUILD
# =============================================================================== #

## build/api: Build application without optimizations
.PHONY: build/api
build/api:
	@echo 'Building application without optimizations...'
	@mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o=./bin/api ./cmd/main
	@echo 'Build for Linux (amd64) complete.'

## build/api/optimized: Build application with optimizations
.PHONY: build/api/optimized
build/api/optimized:
	@echo 'Building application with optimizations...'
	@mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-s -w" -o=./bin/linux_amd64/api-o ./cmd/main
	@echo 'Build for Linux (amd64) complete (optimized).'

# =============================================================================== #
# PRODUCTION DEPLOYMENT TASKS
# =============================================================================== #

## production/connect: Connect to the production server
.PHONY: production/connect
production/connect:
	ssh bot@${PRODUCTION_HOST_IP}

## production/deploy-bot-files: Deploy new binary
.PHONY: production/deploy-bot-files
production/deploy-bot-files:
	@echo 'Deploying new binary ...'
	rsync -P ./bin/linux_amd64/api-o bot@${PRODUCTION_HOST_IP}:/tmp/api-o
	ssh -t bot@${PRODUCTION_HOST_IP} 'set -e; \
	  sudo mkdir -p /opt/bot-client && \
	  sudo mv /tmp/api-o /opt/bot-client && \
	  sudo chown -R bot:bot /opt/bot-client && \
	  sudo chmod +x /opt/bot-client/api-o'

## production/deploy/bot: Deploy application to production
.PHONY: production/deploy/bot
production/deploy/bot:
	@$(MAKE) build/api/optimized
	@$(MAKE) production/deploy-bot-files
	@echo 'Deployment to production complete.'
