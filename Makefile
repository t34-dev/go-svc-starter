APP_NAME := go-svc-starter
APP_REPOSITORY := gitlab.com/zakon47
%:
	@:
MAKEFLAGS += --no-print-directory
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
TEMP_DIR := $(CURDIR)/.temp
ENV_FILE := $(CURDIR)/.env
SECRET_FILE := $(CURDIR)/.secrets
CONFIG_DIR := $(CURDIR)/configs
# ============================== Includes
include make/get-started.mk
include make/proto.mk
include make/lint.mk
include make/test.mk


#GOPROXY=direct go list -m -versions PACKAGE
#go get -u=patch PACKAGE

################################# DEV
NAME_SERVER=server
NAME_CLIENT=client

build-server:
	go build -o .bin/$(NAME_SERVER)$(APP_EXT) cmd/server/*
server: build-server
	@.bin/$(NAME_SERVER)${APP_EXT}

build-client:
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
