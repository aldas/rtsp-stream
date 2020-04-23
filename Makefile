DOCKER_VERSION?=2
PKG := "github.com/Roverr/rtsp-stream"
PKG_LIST := $(shell go list ${PKG}/...)

init:
	git config core.hooksPath ./scripts/.githooks
	@go get -u golang.org/x/lint/golint

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Vet the files
	@go vet ${PKG_LIST}

test: ## Runs tests
	go test ./...

check: lint vet test ## check project

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/coverage.sh html;

run:  ## Builds & Runs the application
	go build . && ./rtsp-stream

docker-build:  ## Builds normal docker container
	docker build -t roverr/rtsp-stream:${DOCKER_VERSION} .

docker-debug: ## Builds the image and starts it in debug mode
	rm -rf ./log && mkdir log && \
	$(MAKE) docker-build && \
	docker run -d \
	-v `pwd`/log:/var/log \
	-e RTSP_STREAM_DEBUG=true \
	-e RTSP_STREAM_AUTH_JWT_ENABLED=true \
	-e RTSP_STREAM_AUTH_JWT_SECRET=your-256-bit-secret \
	-e RTSP_STREAM_BLACKLIST_COUNT=2 \
	-p 8080:8080 \
	roverr/rtsp-stream:${DOCKER_VERSION}

docker-build-mg:  ## Builds docker container with management UI
	docker build -t roverr/rtsp-stream:${DOCKER_VERSION}-management -f Dockerfile.management .

docker-debug-mg: ## Builds management image and starts it in debug mode
	rm -rf ./log && mkdir log && \
	$(MAKE) docker-build-mg && \
	docker run -d \
	-v `pwd`/log:/var/log \
	-e RTSP_STREAM_DEBUG=true \
	-e RTSP_STREAM_BLACKLIST_COUNT=2 \
	-p 3000:80 -p 8080:8080 \
	roverr/rtsp-stream:${DOCKER_VERSION}-management

docker-all: ## Runs tests then builds all versions of docker images
	$(MAKE) test && $(MAKE) docker-build && $(MAKE) docker-build-mg

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
