GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=manager
VERSION?=0.0.0
SERVICE_PORT?=3000
DOCKER_REGISTRY?=manager/#if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true
BIN_FOLDER?=bin/
MAIN_PATH?=cmd/manager/main.go


GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build run vendor

all: help

lint: lint-go lint-dockerfile

lint-dockerfile:
# If dockerfile is present we lint it.
ifeq ($(shell test -e ./Dockerfile && echo -n yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm \
		-i $(CONFIG_OPTION) \
		hadolint/hadolint \
		hadolint \
		$(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

lint-go:
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:latest-alpine \
		golangci-lint run \
		--deadline=65s $(OUTPUT_OPTIONS)

clean:
	rm -fr ./bin

test:
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GOTEST) -v -race ./... $(OUTPUT_OPTIONS)

coverage:
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off go get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif

vendor:
	@echo 'Creating vendor folder'
	@$(GOCMD) mod vendor

build: vendor
	@echo 'Building ${BINARY_NAME}'
	@mkdir -p bin
	@$(GOCMD) build -mod vendor -o $(BIN_FOLDER)$(BINARY_NAME) $(MAIN_PATH)

docker-build: vendor
	docker build --rm --tag $(BINARY_NAME) .

proto-build:
	@if ! which protoc > /dev/null; then \
    		echo "error: protoc not installed" >&2; \
    		exit 1; \
	fi
	@go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	@go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pb/*.proto

docker-release:
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

run:
	@$(GOCMD) run $(MAIN_PATH)

docker-run: docker-build
	docker run --network host $(BINARY_NAME)

watch:
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) -p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air

help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}build           ${RESET} ${GREEN}Build your project and put the output binary in $(BIN_FOLDER)(BINARY_NAME)${RESET}"
	@echo "  ${YELLOW}clean           ${RESET} ${GREEN}Remove build related file${RESET}"
	@echo "  ${YELLOW}docker-build    ${RESET} ${GREEN}Use the dockerfile to build the container (name: $(BINARY_NAME))${RESET}"
	@echo "  ${YELLOW}proto-build     ${RESET} ${GREEN}Use protoc to compile gRPC service definitions ${RESET}"
	@echo "  ${YELLOW}docker-release  ${RESET} ${GREEN}Release the container \"$(DOCKER_REGISTRY)$(BINARY_NAME)\" with tag latest and $(VERSION) ${RESET}"
	@echo "  ${YELLOW}docker-run      ${RESET} ${GREEN}Build and run the container ${RESET}"
	@echo "  ${YELLOW}lint            ${RESET} ${GREEN}Run all available linters${RESET}"
	@echo "  ${YELLOW}lint-dockerfile ${RESET} ${GREEN}Lint your Dockerfile${RESET}"
	@echo "  ${YELLOW}lint-go         ${RESET} ${GREEN}Use golintci-lint on your project${RESET}"
	@echo "  ${YELLOW}test            ${RESET} ${GREEN}Run the tests of the project${RESET}"
	@echo "  ${YELLOW}vendor          ${RESET} ${GREEN}Copy of all packages needed to support builds and tests in the vendor directory${RESET}"
	@echo "  ${YELLOW}watch           ${RESET} ${GREEN}Run the code with cosmtrek/air to have automatic reload on changes${RESET}"
	@echo "  ${YELLOW}run-jaeger      ${RESET} ${GREEN}Run Jaeger to store traces${RESET}"
	@echo "  ${YELLOW}run-registry	  ${RESET} ${GREEN}Run a docker container registry on port 5000${RESET}"

