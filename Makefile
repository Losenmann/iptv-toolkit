VERSION ?= latest
RUN_ARG ?=

IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
IMAGE_VERSION ?= ${VERSION}
IMAGE_BIN_COMPRESS ?= true

MAKE_PATH_BUILD ?= ./build
MAKE_DATE!=`date`
MAKE_DATE_U!=`date +%s -d '${MAKE_DATE}'`
MAKE_DATE_Y!=`date +%Y -d '${MAKE_DATE}'`
MAKE_DATE_R!=`date -R -d '${MAKE_DATE}'`
MAKE_DATE_C!=`date '+%a %b %d %Y' -d '${MAKE_DATE}'`
MAKE_USER!=`whoami`
MAKE_PWD!=`pwd`
HOME=${MAKE_PWD}/pkg

ifeq ($(PKG_VERSION),)
	PKG_VERSION=0.0.1
endif
ifeq ($(PKG_REVISION),)
	PKG_REVISION=1
endif
ifeq ($(PKG_MAINTAINER),)
	PKG_MAINTAINER=${MAKE_USER}
endif
ifeq ($(PKG_MAINTAINER_EMAIL),)
	PKG_MAINTAINER_EMAIL=${PKG_MAINTAINER}@example.com
endif
ifeq ($(PKG_LICENSE),)
	PKG_LICENSE=Apache-2.0
endif
ifeq ($(PKG_DESCRIPTION),)
	PKG_DESCRIPTION=No description
endif
ifeq ($(PKG_CHANGELOG),)
	PKG_CHANGELOG=Unknown
endif
ifeq ($(PKG_ARCH),)
	PKG_ARCH=any
endif

.PHONY: realesae run docker testing

all: build

test-go-run:
	@go run ./main.go ${RUN_ARG}

clear:
	@rm -rf ${MAKE_PATH_BUILD}/*

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

docker-up:
	@docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	@docker compose -f ./deploy/docker-compose.yaml down

testing-pre-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env up -d

testing-post-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env down

build-apk:
	@install -m755 -D ./artifact/bin/*linux-${PKG_ARCH} ./pkg/apkbuild/iptv-toolkit/iptv-toolkit
	@printf "${PKG_SIGN_ALPINE}" > /root/.abuild/root-${MAKE_TIME}.rsa
	@openssl rsa -in /root/.abuild/root-${MAKE_TIME}.rsa -pubout > /root/.abuild/root-${MAKE_TIME}.rsa.pub
	@sed -i '/^[^#]/s/^/#/g' /root/.abuild/abuild.conf
	@sed -i -e '/^pkgver/s/$$/${PKG_VERSION}/g' \
		-e '/^pkgdesc/s/$$/"${PKG_DESCRIPTION}"/g' \
		-e '/^url/s|$$|"${PKG_HOME_URL}"|g' \
		-e '/^license/s/$$/"${PKG_LICENSE}"/g' ./pkg/apkbuild/iptv-toolkit/APKBUILD
	@apkbuild-lint ./pkg/apkbuild/iptv-toolkit/APKBUILD
	@abuild -F checksum ./pkg/apkbuild/iptv-toolkit
	@cd ./pkg/apkbuild/iptv-toolkit; abuild -Fr
	@cat ./pkg/apkbuild/iptv-toolkit/APKBUILD
	@efwegqa

build-rpm:
	@mkdir -p ./pkg/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
	@install -m755 -D ./artifact/bin/*linux-${PKG_ARCH} ./pkg/rpmbuild/iptv-toolkit-${PKG_VERSION}/iptv-toolkit
	@tar -czvf ./pkg/rpmbuild/SOURCES/iptv-toolkit-${PKG_VERSION}.tar.gz -C ./pkg/rpmbuild/ iptv-toolkit-${PKG_VERSION} --remove-files
	@sed -i -e '/^Version/s/$$/${PKG_VERSION}/g' \
		-e '/^License/s/$$/${PKG_LICENSE}/g' \
		-e '/^URL/s|$$|${PKG_HOME_URL}|g' \
		-e '/^%description/s/$$/\n  ${PKG_DESCRIPTION}/g' \
		-e '/^%changelog/s/$$/\n  * ${MAKE_DATE_C} ${MAINTAINER}/g' ./pkg/rpmbuild/SPECS/iptv-toolkit.spec
ifdef PKG_CHANGELOG
	@echo '${PKG_CHANGELOG}' |sed 's/^/  - /g' |tee -a ./pkg/rpmbuild/SPECS/iptv-toolkit.spec 1> /dev/null
endif
	@rpmlint ./pkg/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmbuild -ba ./pkg/rpmbuild/SPECS/iptv-toolkit.spec

build-deb:
	@chmod +x ./pkg/debbuild/iptv-toolkit/debian/rules
	@install -m755 -D ./artifact/bin/*linux-${PKG_ARCH} ./pkg/debbuild/iptv-toolkit/usr/bin/iptv-toolkit
	@sed -i -e '/^Description/s/$$/\n ${PKG_DESCRIPTION}/g' \
		-e '/^Homepage/s| .*| ${PKG_HOME_URL}|g' \
		-e '/^Architecture/s/ .*/ ${PKG_ARCH}/g' \
		-e '/^Maintainer/s/ .*/ ${MAINTAINER}/g' \
		-e '/^Vcs-Browser/s| .*| ${PKG_HOME_URL}|g' \
		-e '/^Vcs-Git/s| .*| ${PKG_HOME_URL}.git|g' \
		-e '16,$$d' ./pkg/debbuild/iptv-toolkit/debian/control
	@sed -i -e '/^Copyright/s/$$/\n ${MAKE_DATE_Y}/g' \
		-e '/^Upstream-Contact/s/$$/${PKG_MAINTAINER_EMAIL}/g' ./pkg/debbuild/iptv-toolkit/debian/copyright
	@sed -i -e '/urgency=/s/(.*)/(${PKG_VERSION}-${PKG_REVISION})/g' \
		-e '2,$$d' ./pkg/debbuild/iptv-toolkit/debian/changelog
ifdef PKG_CHANGELOG
	@echo '${PKG_CHANGELOG}' |sed 's/^/  * /g' |tee -a ./pkg/debbuild/iptv-toolkit/debian/changelog 1> /dev/null
endif
	@echo '\n -- ${PACKAGER}  ${MAKE_DATE_R}\n' >> ./pkg/debbuild/iptv-toolkit/debian/changelog
	@cd ./pkg/debbuild/iptv-toolkit; dpkg-buildpackage -b -us -uc

install:
	@install -m755 -D ./artifact/bin/*linux-${PKG_ARCH} /usr/bin/iptv-toolkit

build-image:
	@docker buildx build . \
		-f Dockerfile \
		--build-arg VERSION=${IMAGE_VERSION} \
		--build-arg BIN_COMPRESS=${IMAGE_BIN_COMPRESS} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:${IMAGE_VERSION} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest --load

build-bin:
	@docker buildx build \
		-f Dockerfile \
		--build-arg VERSION=${IMAGE_VERSION} \
		--build-arg BIN_COMPRESS=${IMAGE_BIN_COMPRESS} \
		--output=type=local,dest=${MAKE_PATH_BUILD}/artifact/images \
		.
	@cp -p ./artifact/images/usr/bin/iptv-toolkit-* ./artifact/bin/

build-bin-main:
	@GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
		go build \
		-ldflags "-s -w" \
		-o ./artifact/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} \
		./main.go
ifeq (${BIN_COMPRESS},true)
ifneq (${TARGETARCH},riscv64)
ifneq (${TARGETARCH},s390x)
	@upx --best --lzma ./artifact/bin/iptv-toolkit-*
endif
endif
endif

build-bin-udpxy:
	@mkdir -p ./build
	@git -C ./build clone --branch ${UDPXY_VERSION} https://github.com/pcherenkov/udpxy.git
	@make -C ./build/udpxy/chipmunk