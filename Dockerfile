FROM --platform=$BUILDPLATFORM golang:1.22.5-alpine3.20 AS builder-main
ARG TARGETOS
ARG TARGETARCH
ARG ARG_COMPRESS=true
ARG GOPROXY=direct
COPY . /opt/src/
WORKDIR /opt/src
RUN apk add upx git \
    && GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o /usr/bin/iptv-toolkit /opt/src/main.go \
    && if [[ ${ARG_COMPRESS} == "true" ]]; then upx --best --lzma /usr/bin/iptv-toolkit; fi \
    && chmod +x /usr/bin/iptv-toolkit
FROM alpine:3.20.3 AS builder-udpxy
ARG ARG_UDPXY_VERSION=master
RUN apk add git \
        make \
        gcc \
        libc-dev \
    && git -C /opt clone --branch ${ARG_UDPXY_VERSION} https://github.com/pcherenkov/udpxy.git \
        && make -C /opt/udpxy/chipmunk
FROM alpine:3.20.3 AS app
ARG ARG_VERSION=latest
ENV IPTVTOOLKIT_VERSION=${ARG_VERSION} \
    IPTVTOOLKIT_EPG_DST=/www/iptv-toolkit/tvguide \
    IPTVTOOLKIT_PLAYLIST_DST=/www/iptv-toolkit/playlist \
    IPTVTOOLKIT_WEB_PATH=/www/iptv-toolkit
COPY --from=builder-main /usr/bin/iptv-toolkit /usr/bin/iptv-toolkit
COPY --from=builder-udpxy /opt/udpxy/chipmunk/udpxy /usr/bin/udpxy
RUN mkdir -p /www/iptv-toolkit/tvguide /www/iptv-toolkit/tvrecord /www/iptv-toolkit/playlist
WORKDIR /www/iptv-toolkit
ENTRYPOINT ["iptv-toolkit"]
CMD ["-S", "-W"]
