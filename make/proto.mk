PROTO_IN := $(CURDIR)/api
PROTO_OUT := $(CURDIR)/pkg/api
VENDOR := $(CURDIR)/.vendor.proto

proto-bin:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN_DIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0
	GOBIN=$(BIN_DIR) go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.1

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
			mkdir -p $(VENDOR)/protoc-gen-openapiv2/options && \
			curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/annotations.proto > $(VENDOR)/protoc-gen-openapiv2/options/annotations.proto && \
			curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/openapiv2.proto > $(VENDOR)/protoc-gen-openapiv2/options/openapiv2.proto ;\
		fi

proto:
	@$(MAKE) proto-random

proto-random:
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
		--openapiv2_out=allow_merge=true,merge_file_name=api:$(PROTO_OUT) \
			--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2$(APP_EXT) \
		$(PROTO_IN)/random_v1/random.proto
	@echo "Done"

proto-test-random:
	$(BIN_DIR)/grpcurl$(APP_EXT) -plaintext \
		-proto $(PROTO_IN)/random_v1/random.proto \
		-import-path $(PROTO_IN) \
		-import-path $(VENDOR) \
		-d '{}' \
		127.0.0.1:50051 \
		random_v1.RandomService/GetPing


.PHONY: proto-bin proto-vendor proto proto-random proto-test-random
