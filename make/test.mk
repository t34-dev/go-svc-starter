test:
	go clean -testcache
	go test internal/... -covermode=count  -count 1

test-coverage:
	go clean -testcache
	@mkdir -p $(TEMP_DIR)
	@CGO_ENABLED=0 go test \
	. \
	-coverprofile=$(TEMP_DIR)/coverage-report.out -covermode=count
	@go tool cover -html=$(TEMP_DIR)/coverage-report.out -o $(TEMP_DIR)/coverage-report.html
	@go tool cover -func=$(TEMP_DIR)/coverage-report.out

.PHONY: test
