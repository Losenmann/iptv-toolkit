FROM golang:1.22.5-alpine3.20 AS build
ARG ARG_COMPRESS
COPY . /src/
RUN apk add upx \
    && go build -ldflags "-s -w" -o /usr/bin/iptv-toolkit /src/main.go \
    && if [[ ${ARG_COMPRESS} == "true" ]]; then upx --best --lzma /usr/bin/iptv-toolkit; fi \
    && chmod +x /usr/bin/iptv-toolkit

FROM alpine:3.20.3 AS udpxy
ARG ARG_UDPXY_VERSION
RUN apk add git \
        make \
        gcc \
        libc-dev \
    && git -C /opt clone --branch ${ARG_UDPXY_VERSION:-master} https://github.com/pcherenkov/udpxy.git \
        && make -C /opt/udpxy/chipmunk

FROM alpine:3.20.3 AS app
ENV VERSION=${ARG_VERSION}
RUN mkdir -p /www/iptv-toolkit/tvguide /www/iptv-toolkit/tvrecord /www/iptv-toolkit/playlist /opt/iptv-toolkit/src
COPY --from=build /usr/bin/iptv-toolkit /usr/bin/iptv-toolkit
COPY --from=udpxy /opt/udpxy/chipmunk/udpxy /usr/bin/udpxy
WORKDIR /opt/iptv-toolkit
ENTRYPOINT ["iptv-toolkit"]
