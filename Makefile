# Set shell
SHELL=/bin/bash
# Include .env file, if any
-include ./config/.env.example
export
-include Makefile.kind
-include Makefile.loader

# Variables
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
PROTOC_OS := $(OS)
ifeq ($(PROTOC_OS), darwin)
	PROTOC_OS := osx
endif
ARCH := $(shell uname -m)
PROTOC_ARCH := $(ARCH)
ifeq ($(PROTOC_ARCH), arm64)
	PROTOC_ARCH := aarch_64
endif
PROTOC_VERSION := 23.3
PROTOC_GEN_GO_VERSION := 1.30.0
PROTOC_GEN_GO_GRPC_VERSION := 1.3.0
PROTOC_GEN_GO_GRPC_GATEWAY_VERSION := 2.16.0
PROTOC_GEN_GO_GRPC_OPENAPI_VERSION := 0.6.8
GOLANGCI_LINT_VERSION := 1.61.0
READ_REPO_TOKEN ?= $(shell cat ~/.netrc | cut -d ' ' -f 6)

PROTO_EVENTS := $(wildcard api/events/*.proto)
PROTO_V1 := $(wildcard api/v1/*.proto)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.DEFAULT_GOAL := help

.PHONY: help

help: ## Показать подсказки по командам make
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

clean: require ## Удалить все бинарные файлы, созданные make-ом
	@rm -rf ./bin

run-local: require ## Собрать и запустить проект локально
	@go build -o ./bin/draft-project gitlab.bronevik.space/bronevik/backend/draft-project/cmd/draft-project
	@./bin/draft-project

test: require ## Запустить тесты
	@CGO_ENABLED=1 go test ./... -vet=all -race -coverprofile=cover.out
	@go tool cover -html=cover.out

lint: lint-install update-lint-config ## Запустить линтер
	@${GOPATH}/bin/golangci-lint version | tail -n 2;
	@${GOPATH}/bin/golangci-lint run ./... --timeout "5m"

lint-install: require ## Установить линтер
ifeq ($(shell ${GOPATH}/bin/golangci-lint --version 2>/dev/null), )
	@echo "Installing golangci-lint for ${OS} ${ARCH}"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)
endif

lint-uninstall: require ## Удалить линтер
ifneq ($(shell ${GOPATH}/bin/golangci-lint --version 2>/dev/null), )
	@rm -f ${GOPATH}/bin/golangci-lint
	@echo "golangci-lint uninstalled"
endif

proto: protobuf-install ## Генерация кода из protobuf
	@if [ ! -d ./pkg ] ; then mkdir pkg; fi
ifneq ($(PROTO_EVENTS),)
	@protoc -I . \
		-I ./api \
		-I ./api/google/api \
    	--go_out ./pkg/ --go_opt paths=source_relative \
    	--go-grpc_out ./pkg/ --go-grpc_opt paths=source_relative \
    	--grpc-gateway_out ./pkg/ \
        --grpc-gateway_opt paths=source_relative \
        --openapi_out ./api/events/ \
    	api/events/*.proto
	@sed -i.bak 's/version: [0-9.]*/version: $${API_VERSION}/' api/events/openapi.yaml && rm api/events/openapi.yaml.bak
else
	@echo "No proto files found in api/events/, skipping code generation"
endif
ifneq ($(PROTO_V1),)
	@protoc -I . \
		-I ./api \
		-I ./api/google/api \
    	--go_out ./pkg/ --go_opt paths=source_relative \
    	--go-grpc_out ./pkg/ --go-grpc_opt paths=source_relative \
    	--grpc-gateway_out ./pkg/ \
        --grpc-gateway_opt paths=source_relative \
        --openapi_out ./api/v1/ \
        --openapi_opt version=0.0.0 \
        --docgen_out=. \
    	api/v1/*.proto
	@sed -i.bak 's/version: [0-9.]*/version: $${API_VERSION}/' api/v1/openapi.yaml && rm api/v1/openapi.yaml.bak
else
	@echo "No proto files found in api/v1/, skipping code generation"
endif

protobuf-install: require ## Установить protoc и плагины для генерации кода из protobuf
ifeq ($(shell ${GOPATH}/bin/protocc --version 2>/dev/null), )
	@echo "Installing protoc v$(PROTOC_VERSION) for $(PROTOC_OS) $(PROTOC_ARCH)"
	curl -sLO "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip"
	@unzip -o protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip -d $(GOPATH)/
	@rm -f ${GOPATH}/readme.txt && rm -f protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-$(PROTOC_ARCH).zip
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v$(PROTOC_GEN_GO_GRPC_GATEWAY_VERSION)
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v$(PROTOC_GEN_GO_GRPC_OPENAPI_VERSION)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$(PROTOC_GEN_GO_GRPC_VERSION)
endif

protobuf-uninstall: require ## Удалить protoc и его плагины
ifneq ($(shell ${GOPATH}/bin/protoc --version 2>/dev/null), )
	@rm -f ${GOPATH}/bin/protoc && rm -f ${GOPATH}/bin/protoc-gen* && rm -rf ${GOPATH}/include/google/protobuf
	@echo "protoc uninstalled"
endif

require: ## Проверить наличие необходимых программ
	@mkdir -p ./bin
	@which go &> /dev/null || (echo "error: go is required."; exit 1)
	@which curl &> /dev/null || (echo "error: curl is required."; exit 1)
	@which git &> /dev/null || (echo "error: git is required."; exit 1)

docker: ## Собрать образ сервиса для отправки в репозиторий
	@go mod vendor; \
	docker build \
	--build-arg CGO_ENABLED="${CGO_ENABLED}" \
	--build-arg GOOS="${GOOS}" \
	--build-arg GOARCH="${GOARCH}" \
	--build-arg SERVICE_NAME="$(SERVICE_NAME)" \
	--build-arg API_VERSION="$(OPENAPI_VERSION)" \
	--build-arg SKIP_OPENAPI="$(SKIP_OPENAPI)" \
	-t $(DOCKER_IMAGE) -f $(SERVICE_DOCKERFILE) . ; \
	rm -rf vendor
