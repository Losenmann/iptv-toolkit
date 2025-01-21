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
GREEN=
RED=
NC=
PKG_VERSION=0.0.1
PKG_MAINTAINER=`whoami`
PKG_MAINTAINER_EMAIL=root@unknown
PKG_LICENSE=Apache
PKG_DESCRIPTION=none
PKG_ARCH=any

.EXPORT_ALL_VARIABLES:

.PHONY: realesae run docker testing

all: build

test-go-run:
	@go run ./main.go ${RUN_ARG}

clear:
	@rm -rf ${ENV_PATH_BUILD}/*

build:
	@mkdir -p ${ENV_PATH_BUILD}/bin/ ${ENV_PATH_BUILD}/artifact/images ${ENV_PATH_BUILD}/artifact/cache
	@docker buildx build . \
		--platform=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm/v8,linux/arm64,linux/riscv64,linux/s390x,linux/ppc64le \
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		--output=type=local,dest=${ENV_PATH_BUILD}/artifact/images \
		--cache-to=type=local,dest=${ENV_PATH_BUILD}/artifact/cache
	@mv ${ENV_PATH_BUILD}/artifact/images/*/usr/bin/iptv-toolkit-* ${ENV_PATH_BUILD}/bin/
	@docker buildx build . \
		--build-arg ARG_VERSION=${IMAGE_VERSION} \
		--build-arg ARG_COMPRESS=${ENV_IMAGE_COMPRESS_BIN} \
		--cache-from=type=local,dest=${ENV_PATH_BUILD}/artifact/cache
		-t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest --load

testing:
	@netstat -tulpn 2>/dev/null |grep 4023 || exit 1
	@curl -sLo /dev/null -w "%{http_code}" http://localhost:4023 |grep "200" || exit 1
	@curl -sLo /dev/null -w "%{http_code}" http://localhost:4023/files |grep "200" || exit 1
	@wget --spider -qL http://localhost:4023/files/playlist || exit 1
	@wget --spider -qL http://localhost:4023/files/tvguide || exit 1
	@wget -qL http://localhost:4023/files/playlist/playlist.m3u -P ./testing_tmp || exit 1
	@wget -qL http://localhost:4023/files/playlist/playlist.xml -P ./testing_tmp || exit 1
	@wget -qL http://localhost:4023/files/tvguide/epg.xml -P ./testing_tmp || exit 1
	@wget -qL http://localhost:4023/files/tvguide/epg.xml.gz -P ./testing_tmp || exit 1
	@wget -qL http://localhost:4023/files/tvguide/epg.zip -P ./testing_tmp || exit 1
	@sha256sum -c ./testing/sha256sums || exit 1

realesae: build testing

docker-up:
	@docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	@docker compose -f ./deploy/docker-compose.yaml down

testing-pre-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env up -d

testing-post-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env down

build-rpm:
	@rpmdev-setuptree
	@mkdir ~/rpmbuild/iptv-toolkit-${PKG_VERSION}
	@cp -r ./pkg/rpmbuild/* ~/rpmbuild/
	@cp ./artifact/bin/*linux-${PKG_ARCH} ~/rpmbuild/iptv-toolkit-${PKG_VERSION}/iptv-toolkit
	@tar -czvf ~/rpmbuild/SOURCES/iptv-toolkit-${PKG_VERSION}.tar.gz -C ~/rpmbuild/ iptv-toolkit-${PKG_VERSION} --remove-files
	@sed -i -e '/^Version/s/$$/${PKG_VERSION}/g' \
		-e '/^License/s/$$/${PKG_LICENSE}/g' \
		-e '/^%description/s/$$/\n  ${PKG_DESCRIPTION}/g' ~/rpmbuild/SPECS/iptv-toolkit.spec
	@cat ~/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmlint ~/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmbuild -ba ~/rpmbuild/SPECS/iptv-toolkit.spec