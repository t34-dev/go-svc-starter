proto-plugin:
	GOBIN=$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN_DIR) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(BIN_DIR) go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	GOBIN=$(BIN_DIR) go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest

proto-vendor:
		@if [ ! -d api/google ]; then \
			git clone https://github.com/googleapis/googleapis api/googleapis &&\
			mkdir -p  api/google/ &&\
			mv api/googleapis/google/api api/google &&\
			rm -rf api/googleapis ;\
		fi
		@if [ ! -d api/protoc-gen-openapiv2 ]; then \
			mkdir -p api/protoc-gen-openapiv2/options && \
			curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/annotations.proto > api/protoc-gen-openapiv2/options/annotations.proto && \
			curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/protoc-gen-openapiv2/options/openapiv2.proto > api/protoc-gen-openapiv2/options/openapiv2.proto ;\
		fi

proto:
	@mkdir -p pkg/api
	@$(MAKE) proto-random

proto-random:
	@protoc --proto_path api \
		--go_out=pkg/api --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=$(BIN_DIR)/protoc-gen-go$(APP_EXT) \
		--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-go-grpc=$(BIN_DIR)/protoc-gen-go-grpc$(APP_EXT) \
		--grpc-gateway_out=pkg/api --grpc-gateway_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=$(BIN_DIR)/protoc-gen-grpc-gateway$(APP_EXT) \
		--openapiv2_out=pkg/api --openapiv2_opt logtostderr=true,allow_repeated_fields_in_body=true \
			--plugin=protoc-gen-openapiv2=$(BIN_DIR)/protoc-gen-openapiv2$(APP_EXT) \
	  	--openapi_out=pkg/api/random_v1 \
			--plugin=protoc-gen-openapi=$(BIN_DIR)/protoc-gen-openapi$(APP_EXT) \
		api/random_v1/random.proto
	@echo "Done"

proto-test-random:
	$(BIN_DIR)/grpcurl$(APP_EXT) -plaintext \
		-proto api/random_v1/random.proto \
		-import-path ./api \
		-d '{}' \
		127.0.0.1:50051 \
		random_v1.RandomService/GetPing


.PHONY: proto-plugin proto-vendor proto proto-random proto-test
