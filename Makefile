VERSION ?= latest
RUN_ARG ?=
ARCH_ALL = false 
IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= ${VERSION}
ENV_IMAGE_COMPRESS_BIN ?= true
ENV_PATH_BUILD ?= ./build

.PHONY: run docker

all: build

run:
	@go run ./main.go ${RUN_ARG}

build: compil compress

compil:
ifeq ($(strip $(ARCH_ALL)),false)
	@go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit ./main.go
endif
ifeq ($(strip $(ARCH_ALL)),true)
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-amd64 ./main.go
	@GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-i386 ./main.go
	@GOOS=linux GOARCH=arm go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-arm ./main.go
	@GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-arm64 ./main.go
endif

compress:
	@upx --best --lzma ${ENV_PATH_BUILD}/iptv-toolkit*

docker:
ifneq ($(strip $(IMAGE_VERSION)),latest)
	@docker buildx build . \
		--platform=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64 \
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest --push
endif
ifeq ($(strip $(IMAGE_VERSION)),latest)
	@docker buildx build . \
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest --load
	@docker compose -f ./deploy/docker-compose.yaml up -d
endif

docker-up:
	@docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	@docker compose -f ./deploy/docker-compose.yaml down