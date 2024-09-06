################################################################## CODE FORMATTING
gofmt:
	@echo "Formatting Go code..."
	@find . -name '*.go' -exec gofmt -s -w {} \;

padding:
	GOBIN=$(BIN_DIR) go install github.com/t34-dev/go-field-alignment/v2/cmd/gofield@v2.0.4
	@$(BIN_DIR)/gofield$(APP_EXT) --files "." --fix

################################################################## LINTING
LINT_CI := $(BIN_DIR)/golangci-lint${APP_EXT} run \
	./... \
	--config=./.golangci.yaml \
	--timeout=10m

lint-install:
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@latest
	GOBIN=$(BIN_DIR) go install golang.org/x/tools/cmd/goimports@latest

lint-prepare:
	@gofumpt -w .
	@goimports -w .

lint-fix-all: lint-prepare
	@$(LINT_CI) --fix

lint-fix: lint-prepare
	@$(LINT_CI) --fix --new

lint:
	@$(LINT_CI)


.PHONY: gofmt padding lint-install lint-prepare lint-fix-all lint-fix lint
