BUILD_DIR=bin
BINARY_NAME=previewer
MAIN_FILE=cmd/previewer/main.go

.PHONY: install
install:
	go mod tidy

.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

.PHONY: dev
dev:
	PORT=8080 HOST=0.0.0.0 CACHE_SIZE=10 CACHE_DIR=tmp go run $(MAIN_FILE)

.PHONY: run
run:
	docker-compose up --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: test-integration
test-integration:
	docker-compose -f test/integration/docker-compose.yml up -d --build
	go test -tags=integration ./...
	docker-compose -f test/integration/docker-compose.yml down

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -tags '!integration' ./...