FROM debian:12.9 AS builder
ARG TARGETOS \
    TARGETARCH \
    USER=root
ENV USER=${USER}
WORKDIR /root
RUN apt update && apt install -y --no-install-recommends make git build-essential dpkg-dev devscripts debhelper dh-make
COPY ./Makefile ./artifact/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} /root/
RUN make build-deb

FROM scratch AS package
COPY --from=builder /root/debbuild/RPMS/*/* /
ENTRYPOINT ["/iptv-toolkit-*"]
