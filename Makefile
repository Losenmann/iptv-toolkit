RUN_ARG ?=
IMAGE_REPO ?= losenmann
IMAGE_NAME ?= iptv-toolkit
BIN_COMPRESS?=true
MAKE_GOTMP=`pwd`/artifact/go
MAKE_USER!=whoami

ifeq ($(TARGETOS),)
	TARGETOS!=go env GOOS
endif
ifeq ($(TARGETARCH),)
	TARGETARCH!=go env GOARCH
endif
ifeq ($(PKG_VERSION),)
	PKG_VERSION!=git log -1 --tags --pretty="%S" |tr -d "v"
endif
ifeq ($(PKG_MAINTAINER),)
	PKG_MAINTAINER=${MAKE_USER}
endif
ifeq ($(PKG_MAINTAINER_EMAIL),)
	PKG_MAINTAINER_EMAIL=${PKG_MAINTAINER}@example.com
endif
ifeq ($(PKG_MAINTAINER_EMAIL),)
	PKG_MAINTAINER_EMAIL=${PKG_MAINTAINER}@example.com
endif
ifeq ($(PKG_USER),)
	PKG_USER!=id -u
endif
ifeq ($(PKG_GROUP),)
	PKG_GROUP!=id -g
endif
MAINTAINER=${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}>

.PHONY: realesae run docker testing pkg

all: build

test-go-run:
	@go run ./main.go ${RUN_ARG}

clear:
	@rm -rvf ./pkg/.config \
		./pkg/.docker \
		./pkg/rpmbuild/BUILD \
		./pkg/rpmbuild/BUILDROOT \
		./pkg/rpmbuild/RPMS \
		./pkg/rpmbuild/SOURCES \
		./pkg/rpmbuild/SRPMS

testing:
	@printf "Check listen port "
	@ss -tl src :4022 |grep -q 4022 && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Check web path '/files' "
	@curl -sfLo /dev/null http://localhost:4022/files && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Check web path '/files/playlist' "
	@curl -sfLo /dev/null http://localhost:4022/files/playlist && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Check web path '/files/tvguide' "
	@curl -sfLo /dev/null http://localhost:4022/files/tvguide && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Checking conversion 'playlist.m3u' "
	@wget -qL http://localhost:4022/files/playlist/playlist.m3u -P ./testing/_tmp && printf "✔\n" || ((printf "❌\n"; exit 1))
	@@printf "Checking conversion 'playlist.xml' "
	@wget -qL http://localhost:4022/files/playlist/playlist.xml -P ./testing/_tmp && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Checking conversion 'epg.xml' "
	@wget -qL http://localhost:4022/files/tvguide/epg.xml -P ./testing/_tmp && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Checking conversion 'epg.xml.gz' "
	@wget -qL http://localhost:4022/files/tvguide/epg.xml.gz -P ./testing/_tmp && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Checking conversion 'epg.zip' "
	@wget -qL http://localhost:4022/files/tvguide/epg.zip -P ./testing/_tmp && printf "✔\n" || ((printf "❌\n"; exit 1))
	@printf "Checksums verification...\n"
	@sha256sum -c ./testing/sha256sums || exit 1
	@echo "All checks are successful ✔"

docker-up:
	@docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	@docker compose -f ./deploy/docker-compose.yaml down

testing-pre-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env up -d

testing-post-stage:
	@docker compose -f ./deploy/docker-compose.yaml --env-file ./testing/testing.env down

build-deb:
	@git config --global --add safe.directory /opt/src
	@install -o ${PKG_USER} -g ${PKG_GROUP} -d ./artifact ./artifact/pkg
	@install -m755 -D ./artifact/bin/*linux-${TARGETARCH} ./pkg/debbuild/iptv-toolkit/iptv-toolkit
	@chmod +x ./pkg/debbuild/iptv-toolkit/debian/rules
	@sed -i '/^Maintainer/s/:.*/: ${MAINTAINER}/g' ./pkg/debbuild/iptv-toolkit/debian/control
	@sed -i -e "/; urgency=/s/([0-9.]*)/(${PKG_VERSION}-1)/" \
		-e '2,$$d' \
		-e "/; urgency=/s/$$/\n` \
			LANG=en_US git -P tag -l --sort=-v:refname --format='%(contents) -- %(*authorname) %(*authoremail)  %(*authordate:rfc)' $$(git describe --abbrev=0) \
			|sed -e "/^ -- /! s/^/  * /g" -e 's/$$/\\\n/' \
			|tr -d '\n' \
		`/" ./pkg/debbuild/iptv-toolkit/debian/changelog
	@cd ./pkg/debbuild/iptv-toolkit; dpkg-buildpackage -b -us -uc
	@install -o ${PKG_USER} -g ${PKG_GROUP} -m755 -D ./pkg/debbuild/*.deb -t ./artifact/pkg/

build-rpm:
	@git config --global --add safe.directory /opt/src
	@install -o ${PKG_USER} -g ${PKG_GROUP} -d ./artifact ./artifact/pkg
	@install -o ${PKG_USER} -g ${PKG_GROUP} -d ./pkg/rpmbuild/BUILD/../BUILDROOT/../RPMS/../SOURCES/../SPECS/../SRPMS/
	@install -m755 -D ./artifact/bin/*linux-${TARGETARCH} ./pkg/rpmbuild/iptv-toolkit-${PKG_VERSION}/iptv-toolkit
	@tar -czvf ./pkg/rpmbuild/SOURCES/v${PKG_VERSION}.tar.gz -C ./pkg/rpmbuild/ iptv-toolkit-${PKG_VERSION} --remove-files
	@sed -i -e '/^Version/s/:.*/: ${PKG_VERSION}/' \
		-e "0,/%changelog/!d" \
		-e "s|%changelog|%changelog\n` \
			LANG=en_US git -P tag -l --sort=-v:refname --format='* %(*authordate:format:%a %b %d %Y) %(*authorname) %(*authoremail) - %(tag)%0a%(contents)' \
			|sed -e 's/^[^*]/- /g' -e '/^*/s/ v/ /g' -e '/^*/s/$$/-1/g' -e 's/$$/\\\n/' |tr -d '\n'`|" \
		./pkg/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmlint ./pkg/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmbuild --define "_topdir `pwd`/pkg/rpmbuild" -ba ./pkg/rpmbuild/SPECS/iptv-toolkit.spec
	@rpmlint -r ./pkg/rpmbuild/.rpmlintrc ./pkg/rpmbuild/RPMS/*/*.rpm
	@install -o ${PKG_USER} -g ${PKG_GROUP} -m755 -D ./pkg/rpmbuild/RPMS/*/*.rpm -t ./artifact/pkg/

image:
	@docker buildx build . \
		-f Dockerfile \
		--build-arg BIN_COMPRESS=${BIN_COMPRESS} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest \
		--load

bin:
	@docker run \
		--rm \
		-w /opt/src \
		-v .:/opt/src \
		-e TARGETOS=${TARGETOS} \
		-e TARGETARCH=${TARGETARCH} \
		-e GOPATH=${MAKE_GOTMP}/path \
		-e GOCACHE=${MAKE_GOTMP}/cache \
		-e BIN_COMPRESS=${BIN_COMPRESS} \
		-e PKG_USER=${PKG_USER} \
		-e PKG_GROUP=${PKG_GROUP} \
		-t golang:1.24.3-alpine \
		sh -c "apk add make upx && make build-bin"

pkg:
	@docker buildx build \
		--build-arg PKG_VERSION=${PKG_VERSION} \
		--build-arg PKG_USER=${PKG_USER} \
		--build-arg PKG_GROUP=${PKG_GROUP} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest \
		--output=type=local,dest=./artifact/pkg \
		-f ./pkg/Dockerfile.debian .
	@docker buildx build \
		--build-arg PKG_VERSION=${PKG_VERSION} \
		--build-arg PKG_USER=${PKG_USER} \
		--build-arg PKG_GROUP=${PKG_GROUP} \
		-t ${IMAGE_REPO}/${IMAGE_NAME}:latest \
		--output=type=local,dest=./artifact/pkg \
		-f ./pkg/Dockerfile.rhel .

build-bin:
	@mkdir -p ./artifact/bin ./artifact/go/patch/../cache
	@go telemetry off
	@GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
		go build \
		-ldflags "-s -w" \
		-v \
		-o ./artifact/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} \
		.
ifeq (${BIN_COMPRESS},true)
ifneq (${TARGETARCH},mips64le)
ifneq (${TARGETARCH},riscv64)
ifneq (${TARGETARCH},s390x)
	@upx --best --lzma ./artifact/bin/iptv-toolkit-*
endif
endif
endif
endif
	@chown -R ${PKG_USER}:${PKG_GROUP} ./artifact
	@mkdir -p /tmp/app/files/playlist/../tvguide/../tvrecord
	@cp -p ./artifact/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} /tmp/app/
	@ln -s /tmp/app/iptv-toolkit* /tmp/app/iptv-toolkit

install:
	@mkdir -p /var/www/iptv-toolkit/files/playlist/../tvguide/../tvrecord
	@install -m755 -D ./artifact/bin/*linux-${TARGETARCH} /usr/bin/iptv-toolkit

uninstall:
	@rm -rvf /usr/bin/iptv-toolkit /var/www/iptv-toolkit