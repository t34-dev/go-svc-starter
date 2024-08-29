# Default target
.DEFAULT_GOAL := help

# Help
help:
	@echo "Get Started:"
	@echo "  download       - Download Go module dependencies"
	@echo "  tidy           - Tidy Go module dependencies"
	@echo "  swag           - Install Swag for API documentation"
	@echo "  upgrade        - Upgrade all dependencies"
	@echo "  kill-port      - Kill process on port 3000"
	@echo "  login          - Login to GitLab Docker registry"
	@echo "  workspace      - Set up Go workspace"
	@echo ""
	@echo "Tag:"
	@echo "  tag            - Show current git tag"
	@echo "  tag-up         - Update git tag"
	@echo ""
	@echo "Shared:"
	@echo "  shared-upgrade      - Update shared module to latest version"
	@echo "  shared-upgrade-tag  - Update shared module to specific tag"
	@echo ""
	@echo "Lint:"
	@echo "  gofmt          - Format Go code"
	@echo "  padding        - Fix struct field alignment"
	@echo "  lint-init      - Initialize linter"
	@echo "  lint-fix       - Run linter with auto-fix"
	@echo "  lint           - Run linter on new changes"
	@echo "  lint-all 		- Run linter on all files"
	@echo ""
	@echo "Helps:"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "APP:"
	@echo "  build          - Build the project"
	@echo "  run            - Run the project"
	@echo ""
	@echo "Example:"
	@echo "  make run"
	@echo "  make ENV=prod run  (to use a specific environment)"



# Phony targets
.PHONY: help
