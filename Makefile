IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= v0.0.3
ENV_PATH_CONTEXT ?= .
ENV_PATH_DOCKERFILE ?= Docker/Dockerfile
ENV_PATH_BUILD ?= ./build

.PHONY: all build

all: build compress

build:
	@go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit ./main.go

compress:
	@./upx --best --lzma ${ENV_PATH_BUILD}/iptv-toolkit
