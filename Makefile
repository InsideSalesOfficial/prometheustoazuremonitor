# Default to linux.
GOOS?=linux
GOARCH?=amd64
CGO_ENABLED?=0
TAG?=$(shell git rev-parse --short HEAD)
# DOCKER_REGISTRY=test.azurecr.io/ # Include a trailing slash for something like ACR/ECR/Quay.io or empty for Docker Hub
# DOCKER_REPO=company/repo-name # Your docker repository name


all: linux docker

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o bin/cron cmd/cron/*.go

linux: CGO_ENABLED=0
linux: build


docker:
	docker build -t $(DOCKER_REGISTRY)$(DOCKER_REPO):$(TAG) .
