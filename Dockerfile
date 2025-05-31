FROM --platform=$BUILDPLATFORM golang:1.24.3-alpine AS builder
ARG TARGETOS \
    TARGETARCH \
    BIN_COMPRESS=true
WORKDIR /opt/src
RUN apk add make upx
COPY . .
RUN make build-bin

FROM scratch AS app
COPY --from=builder /opt/src/artifact/bin/iptv-toolkit* /tmp/app /
ENTRYPOINT ["/iptv-toolkit"]
CMD ["-U", "-F", "-S"]
