################################################################## CODE FORMATTING
gofmt:
	@echo "Formatting Go code..."
	@find . -name '*.go' -exec gofmt -s -w {} \;

padding:
	@echo "Checking gopad"
	go install github.com/t34-dev/go-pad-alignment/gopad/...@latest
	@echo "Checking /internal/models/ and /cmd"
	@gopad --files "./internal/models/, ./cmd/" --fix

################################################################## LINTING
GOLANG_LINT_CI := $(BIN_DIR)/golangci-lint${APP_EXT} run \
	./... \
	--config=./.golangci.yml \
	--timeout=10m

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

lint-init: $(BIN_DIR)
	@echo "Initializing linter..."
	@[ -f $(BIN_DIR)/golangci-lint${APP_EXT} ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v1.55.2
	@go install mvdan.cc/gofumpt@latest
	@go install golang.org/x/tools/cmd/goimports@latest

lint-check-install: $(BIN_DIR)
	@if [ ! -f $(BIN_DIR)/golangci-lint${APP_EXT} ]; then \
		echo "$(BIN_DIR)/golangci-lint${APP_EXT} not found, initializing..."; \
		$(MAKE) lint-init; \
	fi

lint-fix: gofmt lint-check-install
	@echo "Running linter with auto-fix..."
	@gofumpt -w .
	@goimports -w .

lint: lint-fix
	@echo "Running linter on new changes..."
	@$(GOLANG_LINT_CI) --new

lint-all: lint-fix
	@echo "Running linter on all files..."
	@$(GOLANG_LINT_CI)
