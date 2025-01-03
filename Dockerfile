FROM --platform=$BUILDPLATFORM golang:1.22.5-alpine3.20 AS builder-main
ARG TARGETOS
ARG TARGETARCH
ARG ARG_COMPRESS=true
ARG GOPROXY=direct
COPY . /opt/src/
WORKDIR /opt/src
RUN apk add upx git \
    && GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o /usr/bin/iptv-toolkit /opt/src/main.go \
    && if [[ ${ARG_COMPRESS} == "true" ]] \
        && [[ ! ${TARGETARCH} == "riscv64" ]] \
        && [[ ! ${TARGETARCH} == "s390x" ]]; then upx --best --lzma /usr/bin/iptv-toolkit; fi \
    && chmod +x /usr/bin/iptv-toolkit
FROM alpine:3.20.3 AS builder-udpxy
ARG ARG_UDPXY_VERSION=master \
    ARG_BUILD_BIN=false
RUN if [[ ${ARG_BUILD_BIN} == "false" ]]; \
    then \
        apk add git \
            make \
            gcc \
            libc-dev \
        && git -C /opt clone --branch ${ARG_UDPXY_VERSION} https://github.com/pcherenkov/udpxy.git \
        && make -C /opt/udpxy/chipmunk; \
    else \
        mkdir -p /opt/udpxy/chipmunk \
        && touch /opt/udpxy/chipmunk/udpxy; \
    fi
FROM alpine:3.20.3 AS app
ARG ARG_VERSION=latest \
    ARG_BUILD_BIN=false \
    TARGETOS \
    TARGETARCH
ENV IPTVTOOLKIT_VERSION=${ARG_VERSION} \
    IPTVTOOLKIT_EPG_DST="/www/iptv-toolkit/tvguide" \
    IPTVTOOLKIT_PLAYLIST_DST="/www/iptv-toolkit/playlist" \
    IPTVTOOLKIT_WEB_PORT="4023" \
    IPTVTOOLKIT_WEB_PATH="/www/iptv-toolkit" \
    IPTVTOOLKIT_CRONTAB="30 6 * * *"
COPY --from=builder-main /usr/bin/iptv-toolkit* /usr/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH}
COPY --from=builder-udpxy /opt/udpxy/chipmunk/udpxy /usr/bin/udpxy
RUN mkdir -p /www/iptv-toolkit/tvguide /www/iptv-toolkit/tvrecord /www/iptv-toolkit/playlist \
    ls -s /usr/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} /usr/bin/iptv-toolkit
WORKDIR /www/iptv-toolkit
ENTRYPOINT ["iptv-toolkit"]
CMD ["-S", "-U", "-W"]
