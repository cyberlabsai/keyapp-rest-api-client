PROJECT_NAME := "exemplo-integracao-api"
PKG := "github.com/cyberlabsai/${PROJECT_NAME}"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean race test coverage coverhtml run

all: build

test: ## Run unittests
	@go test -count 1 -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -count 1 -race -short ${PKG_LIST}

coverage: ## Run code coverage
	@go test -covermode=count -short ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./.../${PROJECT_NAME}

run: ## Run the project
	@go run ${PROJECT_NAME}

build: dep ## Build the binary file
	@go build -a -ldflags '-w -extldflags "-static"' -o ${PROJECT_NAME} cmd/*.go

clean: ## Remove previous build
	@rm -f ${PROJECT_NAME}

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'