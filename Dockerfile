FROM --platform=$BUILDPLATFORM golang:1.24.3-alpine3.20 AS builder-main
ARG TARGETOS \
    TARGETARCH \
    BIN_COMPRESS=true \
    GOPROXY=direct
WORKDIR /opt/src
RUN apk add make git upx
COPY . .
RUN go mod download
RUN make build-bin-main
FROM alpine:3.20.3 AS builder-udpxy
ARG UDPXY_VERSION=master
WORKDIR /opt/src
RUN apk add make git gcc libc-dev
COPY ./Makefile .
RUN make build-bin-udpxy
FROM alpine:3.20.3 AS app
ARG TARGETOS \
    TARGETARCH \
    VERSION=latest \
    ARG_WORKDIR=/www/iptv-toolkit
ENV IPTVTOOLKIT_VERSION=${VERSION} \
    IPTVTOOLKIT_ARCH=${TARGETARCH} \
    IPTVTOOLKIT_EPG_DST="{ARG_WORKDIR}/tvguide" \
    IPTVTOOLKIT_PLAYLIST_DST="${ARG_WORKDIR}/playlist" \
    IPTVTOOLKIT_WEB_PORT="4023" \
    IPTVTOOLKIT_WEB_PATH="{ARG_WORKDIR}" \
    IPTVTOOLKIT_CRONTAB="30 6 * * *"
WORKDIR /www/iptv-toolkit
COPY --from=builder-udpxy /opt/src/build/udpxy/chipmunk/udpxy /usr/bin/udpxy
COPY --from=builder-main /opt/src/artifact/bin/iptv-toolkit* /usr/bin/
RUN mkdir -p ${ARG_WORKDIR}/{playlist,tvguide,tvrecord} \
    && ln -s /usr/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} /usr/bin/iptv-toolkit
ENTRYPOINT ["iptv-toolkit"]
CMD ["-S", "-U", "-W"]
