APP_NAME := go-svc-starter
APP_REPOSITORY := gitlab.com/t34-dev
%:
	@:
MAKEFLAGS += --no-print-directory
export GOPRIVATE=$(APP_REPOSITORY)/*
export GOPRIVATE=github.com/t34-dev/*
export COMPOSE_PROJECT_NAME
# ============================== Environments
PROJECT_NAME ?= $(APP_NAME)
ENV ?= local
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

info: set-env
	@echo "================ [ ENVIRONMENT: $(ENV) ] ================"
	@echo "PROJECT_NAME: $(PROJECT_NAME)"
	@echo ""

#GOPROXY=direct go list -m -versions PACKAGE
#go get -u=patch PACKAGE

################################# DEV
NAME_SERVER=server
NAME_CLIENT=client
NAME_OTHER_SVC=other-svc

build-server:
	@rm -f .bin/$(NAME_SERVER)$(APP_EXT)
	@go build -o .bin/$(NAME_SERVER)$(APP_EXT) cmd/server/*
server: build-server
	@.bin/$(NAME_SERVER)${APP_EXT}

build-other-service:
	@rm -f .bin/$(NAME_OTHER_SVC)$(APP_EXT)
	@go build -o .bin/$(NAME_OTHER_SVC)$(APP_EXT) cmd/other_service/*
other-service: build-other-service
	@.bin/$(NAME_OTHER_SVC)${APP_EXT}

build-client:
	@rm -f .bin/$(NAME_CLIENT)$(APP_EXT)
	go build -o .bin/$(NAME_CLIENT)$(APP_EXT) cmd/client/*
client: build-client
	@.bin/$(NAME_CLIENT)${APP_EXT}



cert-gen:
	mkdir -p cert
	openssl genrsa -out cert/ca.key 4096
	openssl req -new -x509 -key cert/ca.key -sha256 -subj "//C=US/ST=NJ/O=CA, Inc." -days 365 -out cert/ca.cert
	openssl genrsa -out cert/service.key 4096
	openssl req -new -key cert/service.key -out cert/service.csr -config cert/certificate.conf
	openssl x509 -req -in cert/service.csr -CA cert/ca.cert -CAkey cert/ca.key -CAcreateserial \
		-out cert/service.pem -days 365 -sha256 -extfile cert/certificate.conf -extensions req_ext
