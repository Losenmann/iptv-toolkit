VERSION ?= latest
RUN_ARG ?=
ARCH_ALL = false 
IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= ${VERSION}
ENV_IMAGE_COMPRESS_BIN ?= true
ENV_IMAGE_BUILD_BIN ?= true
ENV_PATH_BUILD ?= ./build
ENV_BUILD_ARCH = 386 amd64 arm arm64 riscv64 s390x ppc64le
GREEN=$(shell tput setaf 2 bold)
RED=$(shell tput setaf 1 bold)
NC=$(shell tput sgr0)

.PHONY: run docker testing

all: build

run:
	@go run ./main.go ${RUN_ARG}

build: docker-bin docker-extract

build-local: bin compress
	@echo "Build completed: ./build"

bin: 
ifeq ($(strip $(ARCH_ALL)),false)
	@go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit ./main.go
endif
ifeq ($(strip $(ARCH_ALL)),true)
	@for i in $(ENV_BUILD_ARCH); do \
		GOOS=linux GOARCH=$$i go build -ldflags "-s -w" -o ${ENV_PATH_BUILD}/iptv-toolkit-linux-$$i ./main.go; \
	done
endif

compress:
	@-upx --best --lzma ${ENV_PATH_BUILD}/iptv-toolkit* 2> /dev/null

docker-bin:
	@docker buildx build . \
		--platform=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64,linux/riscv64,linux/s390x,linux/ppc64le \
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		--build-arg ARG_BUILD_BIN=${ENV_IMAGE_BUILD_BIN} \
		--output=type=local,dest=${ENV_PATH_BUILD}

docker-extract:
	@for i in $(patsubst %/,%,$(dir $(wildcard ./build/*/))); do \
		mv $$i/iptv-toolkit* ${ENV_PATH_BUILD}/ \
		&& rm -rf $$i; \
	done

docker:
ifneq ($(strip $(IMAGE_VERSION)),latest)
	@docker buildx build . \
		--platform=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64,linux/riscv64,linux/s390x,linux/ppc64le \
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
endif

docker-up:
	@docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	@docker compose -f ./deploy/docker-compose.yaml down

testing:
	@netstat -tulpn 2>/dev/null |grep 4023 1>/dev/null \
		&& $(info Check webserver port open - ${GREEN}success${NC}) \
		|| $(error Check webserver port open - ${RED}failure${NC})
	@curl -sLo /dev/null -w "%{http_code}" http://localhost:4023 |grep "200" 1>/dev/null \
		&& $(info Check webserver path root redirect - ${GREEN}success${NC}) \
		|| $(error Check webserver path root redirect - ${RED}failure${NC})
	@curl -sLo /dev/null -w "%{http_code}" http://localhost:4023/files |grep "200" 1>/dev/null \
		&& $(info Check webserver path files open - ${GREEN}success${NC}) \
		|| $(error Check webserver path files open - ${RED}failure${NC})
	@wget --spider -qL http://localhost:4023/files/playlist 1>/dev/null \
		&& $(info Check webserver path playlist exist - ${GREEN}success${NC}) \
		|| $(error Check webserver path playlist exist - ${RED}failure${NC})
	@wget --spider -qL http://localhost:4023/files/tvguide 1>/dev/null \
		&& $(info Check webserver path tvguide exist - ${GREEN}success${NC}) \
		|| $(error Check webserver path tvguide exist - ${RED}failure${NC})
	@wget -qL http://localhost:4023/files/tvguide/playlist.m3u -P ./playlist 2>/dev/null || $(error Download playlist.m3u ${RED}failure${NC})
	@wget -qL http://localhost:4023/files/tvguide/playlist.xml -P ./playlist 2>/dev/null || $(error Download playlist.xml ${RED}failure${NC})
	@wget -qL http://localhost:4023/files/tvguide/epg.xml -P ./tvguide 2>/dev/null || $(error Download epg.xml ${RED}failure${NC})
	@wget -qL http://localhost:4023/files/tvguide/epg.xml.gz -P ./tvguide 2>/dev/null || $(error Download epg.xml.gz ${RED}failure${NC})
	@wget -qL http://localhost:4023/files/tvguide/epg.zip -P ./tvguide 2>/dev/null || $(error Download epg.zip ${RED}failure${NC})
	@sha256sum -c ./testing/sha256sums \
		&& printf "Check in work of converters - ${GREEN}success${NC}\n" \
		|| printf "Check in work of converters - ${RED}failure${NC}\n"