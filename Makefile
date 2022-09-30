GOLANG_VERSION := 1.15.6
ALPINE_VERSION := 3.13

DOCKER_REPO := michalsw
APPNAME := whatismyip

SERVER_PORT ?= 8080
PPROF_PORT ?= 5050

.DEFAULT_GOAL := help
.PHONY: build docker-build docker-run-host docker-run-bridge docker-stop

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build bin
	GOOS=linux \
	CGO_ENABLED=0 \
	go build \
	-v \
	-o $(APPNAME) .

docker-build: ## Build docker image
	docker build \
	--platform linux/amd64 \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg APPNAME="$(APPNAME)" \
	--tag="$(DOCKER_REPO)/$(APPNAME):latest" \
	.

docker-run-host: ## Run docker - host network
	docker run --rm -d \
	--network host \
	--name $(APPNAME) \
	--env SERVER_PORT=$(SERVER_PORT) \
	--env PPROF_PORT=$(PPROF_PORT) \
	-p $(SERVER_PORT):$(SERVER_PORT) \
	-p $(PPROF_PORT):$(PPROF_PORT) \
	$(DOCKER_REPO)/$(APPNAME):latest &&\
	docker ps

docker-run-bridge: ## Run docker - bridge network
	docker run --rm -d \
	--name $(APPNAME) \
	--env SERVER_PORT=$(SERVER_PORT) \
	--env PPROF_PORT=$(PPROF_PORT) \
	-p $(SERVER_PORT):$(SERVER_PORT) \
	-p $(PPROF_PORT):$(PPROF_PORT) \
	$(DOCKER_REPO)/$(APPNAME):latest &&\
	docker ps

docker-stop: ## Stop docker
	docker stop $(APPNAME)
