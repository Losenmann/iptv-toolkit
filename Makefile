IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= v0.0.3
ENV_IMAGE_COMPRESS_BIN ?= true
ENV_PATH_CONTEXT ?= .
ENV_PATH_DOCKERFILE ?= .
ENV_PATH_BUILD ?= ./build
ENV_PATH_BUILDKIT ?= ./buildkit

.PHONY: all build

all: build compress

build:
	@go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit ./main.go

compress:
	@./upx --best --lzma ${ENV_PATH_BUILD}/iptv-toolkit

image:
	@mkdir -p ${ENV_PATH_BUILD}
	@pkill buildkitd || true
	@nohup ${ENV_PATH_BUILDKIT}/buildkitd >/dev/null 2>&1 &
	@sleep 1
	${ENV_PATH_BUILDKIT}/buildctl build \
		--frontend gateway.v0 \
		--opt source=docker/dockerfile \
		--local context=${ENV_PATH_CONTEXT} \
		--local dockerfile=${ENV_PATH_DOCKERFILE} \
		--opt build-arg:ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		--opt build-arg:ARG_VERSION=${IMAGE_VERSION} \
		--output type=docker,name=${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION} > ${ENV_PATH_BUILD}/image.tar
	@pkill buildkitd || true
