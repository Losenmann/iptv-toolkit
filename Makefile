VERSION ?=
RUN_ARG ?=
ifneq ($(ARCH_ALL),)
ARCH_ALL = true 
endif
IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= ${VERSION}
ENV_IMAGE_COMPRESS_BIN ?= true
ENV_PATH_CONTEXT ?= .
ENV_PATH_DOCKERFILE ?= .
ENV_PATH_BUILD ?= ./build
ENV_PATH_BUILDKIT ?= ./buildkit

.PHONY: run docker

all: build docker

run:
	@go run ./main.go ${RUN_ARG}

build: compil compress

compil:
	@rm -rf ${ENV_PATH_BUILD}/iptv-toolkit*
ifeq ($(ARCH_ALL),)
	@go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit ./main.go
endif
ifneq ($(ARCH_ALL),)
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-amd64 ./main.go
	@GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-i386 ./main.go
	@GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-arm ./main.go
	@GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-arm64 ./main.go
endif

compress:
	@upx --best --lzma ${ENV_PATH_BUILD}/iptv-toolkit*

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

docker:
	@docker buildx build . \
		--platform=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64,linux/s390x\
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest --push