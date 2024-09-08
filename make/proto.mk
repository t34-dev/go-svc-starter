PROTO_IN := $(CURDIR)/api
PROTO_OUT := $(CURDIR)/pkg/api
VENDOR := $(CURDIR)/.proto_vendor

proto-bin:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN_DIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0
	GOBIN=$(BIN_DIR) go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.1
	GOBIN=$(BIN_DIR) go install github.com/go-swagger/go-swagger/cmd/swagger@latest

proto-vendor:
	@mkdir -p $(VENDOR)
		@if [ ! -d $(VENDOR)/google ]; then \
			mkdir -p  $(VENDOR)/google/ &&\
			git clone https://github.com/googleapis/googleapis $(VENDOR)/googleapis &&\
			mv $(VENDOR)/googleapis/google/api $(VENDOR)/google &&\
			rm -rf $(VENDOR)/googleapis ;\
		fi
		@if [ ! -d $(VENDOR)/validate ]; then \
			mkdir -p  $(VENDOR)/validate/ &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate $(VENDOR)/protoc-gen-validate &&\
			mv $(VENDOR)/protoc-gen-validate/validate/*.proto $(VENDOR)/validate &&\
			rm -rf $(VENDOR)/protoc-gen-validate ;\
		fi
		@if [ ! -d $(VENDOR)/protoc-gen-openapiv2 ]; then \
			mkdir -p $(VENDOR)/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR)/openapiv2 &&\
			mv $(VENDOR)/openapiv2/protoc-gen-openapiv2/options/*.proto $(VENDOR)/protoc-gen-openapiv2/options &&\
			rm -rf $(VENDOR)/openapiv2 ;\
		fi

proto:
	@$(MAKE) proto-common
	@$(MAKE) proto-auth
	@$(MAKE) proto-access
	@$(MAKE) proto-merge

proto-common:
	@mkdir -p $(PROTO_OUT)
	@protoc -I $(PROTO_IN) -I=$(VENDOR) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go$(APP_EXT) \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc$(APP_EXT) \
	  	--validate_out lang=go:$(PROTO_OUT) --validate_opt=paths=source_relative \
			--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate$(APP_EXT) \
		--grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway$(APP_EXT) \
		--openapiv2_out=$(PROTO_OUT) \
			--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2$(APP_EXT) \
		$(PROTO_IN)/common_v1/common.proto
	@echo "common - Done"

proto-auth:
	@mkdir -p $(PROTO_OUT)
	@protoc -I $(PROTO_IN) -I=$(VENDOR) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go$(APP_EXT) \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc$(APP_EXT) \
	  	--validate_out lang=go:$(PROTO_OUT) --validate_opt=paths=source_relative \
			--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate$(APP_EXT) \
		--grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway$(APP_EXT) \
		--openapiv2_out=$(PROTO_OUT) \
			--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2$(APP_EXT) \
		$(PROTO_IN)/auth_v1/auth.proto
	@echo "auth - Done"

proto-access:
	@mkdir -p $(PROTO_OUT)
	@protoc -I $(PROTO_IN) -I=$(VENDOR) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go$(APP_EXT) \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc$(APP_EXT) \
	  	--validate_out lang=go:$(PROTO_OUT) --validate_opt=paths=source_relative \
			--plugin=protoc-gen-validate=$(BIN_DIR)/protoc-gen-validate$(APP_EXT) \
		--grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway$(APP_EXT) \
		--openapiv2_out=$(PROTO_OUT) \
			--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2$(APP_EXT) \
		$(PROTO_IN)/access_v1/access.proto
	@echo "access - Done"


proto-merge:
	@echo "Merging Swagger files in specified order..."
	@#Common should be first and contain a general description
	@$(eval SWAGGER_FILES := \
		$(PROTO_OUT)/common_v1/common.swagger.json \
		$(PROTO_OUT)/auth_v1/auth.swagger.json \
		$(PROTO_OUT)/access_v1/access.swagger.json \
	)
	@$(BIN_DIR)/swagger$(APP_EXT) mixin $(shell echo $(SWAGGER_FILES)) -o $(PROTO_OUT)/api.swagger.json || true
	@if [ ! -f $(PROTO_OUT)/api.swagger.json ]; then echo "Error: Swagger merge failed"; exit 1; fi
	@echo "Swagger files merged"

proto-ping:
	$(BIN_DIR)/grpcurl$(APP_EXT) -plaintext \
		-proto $(PROTO_IN)/common_v1/common.proto \
		-import-path $(PROTO_IN) \
		-import-path $(VENDOR) \
		-d '{}' \
		127.0.0.1:50051 \
		common_v1.CommonV1/GetTime


.PHONY: proto-bin proto-vendor proto proto-common proto-ping
