APP_NAME := go-svc-starter
APP_REPOSITORY := gitlab.com/zakon47
%:
	@:
export GOPRIVATE=$(APP_REPOSITORY)/*
export GOPRIVATE=github.com/t34-dev/*
export COMPOSE_PROJECT_NAME
# ============================== Environments
ENV ?= local
FOLDER ?= /root
VERSION ?=
# ============================== Paths
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)
BIN_DIR := $(CURDIR)/.bin
DEVOPS_DIR := $(CURDIR)/.devops
ENV_FILE := $(CURDIR)/.env
SECRET_FILE := $(CURDIR)/.secrets
CONFIG_DIR := $(CURDIR)/configs
# ============================== Includes
include .make/get-started.mk
include .make/proto.mk
include .make/lint.mk



################################# DEV
build: env-development
	@go build -o .bin/$(APP_NAME)$(APP_EXT) cmd/*
run: build
	@.bin/$(APP_NAME)${APP_EXT}
