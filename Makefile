# Output directory for binaries
BIN_DIR := ./bin

# Source directories and target binaries
SRC_SDLCLIENT := cmd/sdlclient/main.go
TARGET_SDLCLIENT := $(BIN_DIR)/midgarts

SRC_GRFEXPLORER := cmd/grfexplorer/main.go
TARGET_GRFEXPLORER := $(BIN_DIR)/grfexplorer

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## all: builds all binaries
.PHONY: all
all: build build-grfexplorer

## build: builds sdlclient
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(TARGET_SDLCLIENT) $(SRC_SDLCLIENT)

## build-grfexplorer: builds grfexplorer
build-grfexplorer:
	@mkdir -p $(BIN_DIR)
	go build -o $(TARGET_GRFEXPLORER) $(SRC_GRFEXPLORER)

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the cmd/web application
.PHONY: run
run:
	go run github.com/cosmtrek/air@v1.40.4 --c="./air.toml"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify

## test: run tests
.PHONY: test
test:
	go test -race -vet=off ./...

## clean: clean binaries
.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
