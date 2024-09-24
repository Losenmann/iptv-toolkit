IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= v0.0.1
ENV_PATH_CONTEXT ?= .
ENV_PATH_DOCKERFILE ?= Docker/Dockerfile
ENV_PATH_BUILD ?= ./build

.PHONY: all build

all: build

build:
	@go build -ldflags "-s -w" -o ./build/iptv ./main.go

