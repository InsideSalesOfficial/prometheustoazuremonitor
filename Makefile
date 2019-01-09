# Default to linux.
GOOS?=linux
GOARCH?=amd64
CGO_ENABLED?=0
TAG?=$(shell git rev-parse --short HEAD)

all: linux

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o bin/cron cmd/cron/*.go

linux: CGO_ENABLED=0
linux: build

