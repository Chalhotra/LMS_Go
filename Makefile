# Define environment variables
ENV_FILE=.env

# Define directories
CMD_DIR=cmd
DATABASE_DIR=$(CMD_DIR)/database
MIGRATION_DIR=$(DATABASE_DIR)/migration
MAIN_DIR=$(CMD_DIR)/main

# Define executables
EXECUTABLE=BasicCrudApi

# Load environment variables
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

# Define the default target
.PHONY: all
all: build

# Target for building the application
.PHONY: build
build: clean
	go build -o $(EXECUTABLE) $(MAIN_DIR)/main.go

# Target for running the application
.PHONY: run
run: build
	./$(EXECUTABLE)

# Target for running migrations up
.PHONY: migrate-up
migrate-up:
	for f in $(MIGRATION_DIR)/*.up.sql; do \
		mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) -h$(DB_HOST) -P$(DB_PORT) $(DB_NAME) < $$f; \
	done

# Target for running migrations down
.PHONY: migrate-down
migrate-down:
	for f in $(MIGRATION_DIR)/*.down.sql; do \
		mysql -u$(DB_USERNAME) -p$(DB_PASSWORD) -h$(DB_HOST) -P$(DB_PORT) $(DB_NAME) < $$f; \
	done

# Target for cleaning the build
.PHONY: clean
clean:
	go clean
	rm -f $(EXECUTABLE)

# Target for running tests
.PHONY: test
test:
	go test ./...

# Target for linting the code
.PHONY: lint
lint:
	golangci-lint run

# Target for updating dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Target for formatting the code
.PHONY: fmt
fmt:
	go fmt ./...

# Target for vendor dependencies
.PHONY: vendor
vendor:
	go mod vendor

# Target for debugging environment variables
.PHONY: debug
debug:
	@echo DB_USERNAME=$(DB_USERNAME)
	@echo DB_PASSWORD=$(DB_PASSWORD)
	@echo DB_HOST=$(DB_HOST)
	@echo DB_PORT=$(DB_PORT)
	@echo DB_NAME=$(DB_NAME)
