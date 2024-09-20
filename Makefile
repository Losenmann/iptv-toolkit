IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= v0.0.1
ENV_PATH_CONTEXT ?= .
ENV_PATH_DOCKERFILE ?= Docker/Dockerfile
ENV_PATH_BUILD ?= ./build

.PHONY: all build

all: build

build:
	docker buildx build \
		-t "${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION}" \
		--platform linux/386,linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/riscv64 \
		-f ${ENV_PATH_DOCKERFILE} \
		.

