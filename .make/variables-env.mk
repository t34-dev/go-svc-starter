export GOPRIVATE=github.com/t34-dev/*

# for dynamic values and docker-compose
export COMPOSE_PROJECT_NAME

%:
	@:

# Define default environment
ENV ?= local
FOLDER ?= /root
VERSION ?=

# Paths
BIN_DIR := $(CURDIR)/.bin
DEV_DIR := $(CURDIR)/.devops
ENV_FILE := $(CURDIR)/.env
SECRET_FILE := $(CURDIR)/.secrets
CONFIG_DIR := $(CURDIR)/configs

# Types
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)


# Set environment
set-env: copy-files
	@cp .env.$(ENV) $(ENV_FILE)
	@$(eval include $(ENV_FILE))
	@$(eval include $(SECRET_FILE))
	@$(eval COMPOSE_PROJECT_NAME := $(PROJECT_NAME))
	@echo "================ [ ENVIRONMENT: $(ENV) ] ================"

copy-files:
	@cp .env.$(ENV) $(ENV_FILE)
	@$(eval include $(ENV_FILE))
	@$(eval include $(SECRET_FILE))


log:
	git config --local url."https://oauth2:glpat-Pdmkdx8ssHmJrvthTryk@gitlab.com".insteadOf "https://gitlab.com"

login-git: set-env
	@echo "Setting up local Git configuration for GitLab"
	git config --global url."https://oauth2:${GITLAB_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"
