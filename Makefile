APP_NAME ?= go-svc-starter
APP_REPOSITORY := gitlab.com/t34-dev
%:
	@:
MAKEFLAGS += --no-print-directory
export GOPRIVATE=$(APP_REPOSITORY)/*
export GOPRIVATE=github.com/t34-dev/*
export COMPOSE_PROJECT_NAME
# ============================== Environments
SERVICE_NAME = $(APP_NAME)
ENV = local
FOLDER ?= /root
VERSION ?=
# ============================== Paths
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)
BIN_DIR := $(CURDIR)/.bin
DEVOPS_DIR := $(CURDIR)/.devops
TEMP_DIR := $(CURDIR)/.temp
ENV_FILE := $(CURDIR)/.env
CONFIG_DIR := $(CURDIR)/configs
# ============================== Includes
include make/set-env.mk
include make/get-started.mk
include make/proto.mk
include make/lint.mk
include make/test.mk
include make/cert.mk
include make/compose.mk

info: set-env
	@echo "================ [ ENVIRONMENT: $(ENV) ] ================"
	@echo "SERVICE_NAME: $(SERVICE_NAME)"
	@echo ""

#GOPROXY=direct go list -m -versions PACKAGE
#go get -u=patch PACKAGE

################################# DEV
NAME_SERVER=server
NAME_CLIENT=client
NAME_OTHER_SVC=other-service

build-server: set-env info
	@rm -f $(BIN_DIR)/$(NAME_SERVER)$(APP_EXT)
	@go build -o $(BIN_DIR)/$(NAME_SERVER)$(APP_EXT) cmd/$(NAME_SERVER)/*
server: build-server
	@export ENV=$(ENV) APP_NAME=$(APP_NAME) && $(BIN_DIR)/$(NAME_SERVER)${APP_EXT}

build-other-service: set-env
	@rm -f $(BIN_DIR)/$(NAME_OTHER_SVC)$(APP_EXT)
	@go build -o $(BIN_DIR)/$(NAME_OTHER_SVC)$(APP_EXT) cmd/$(NAME_OTHER_SVC)/*
other-service: build-other-service
	@$(BIN_DIR)/$(NAME_OTHER_SVC)${APP_EXT}

build-client: set-env
	@rm -f $(BIN_DIR)/$(NAME_CLIENT)$(APP_EXT)
	go build -o $(BIN_DIR)/$(NAME_CLIENT)$(APP_EXT) cmd/$(NAME_CLIENT)/*
client: build-client
	@$(BIN_DIR)/$(NAME_CLIENT)${APP_EXT}



