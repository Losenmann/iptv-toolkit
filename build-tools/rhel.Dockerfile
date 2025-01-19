FROM redhat/ubi9:9.5 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION=0.0.1 \
    PKG_MAINTAINER=root \
    PKG_MAINTAINER_EMAIL=root@unknown \
    PKG_LICENSE=apache \
    PKG_DESCRIPTION=A set of tools for working with IPTV services: Playlists, EPG, UDPXY
ENV PKG_VERSION=${PKG_VERSION} \
    USER=${PKG_MAINTAINER} \
    PKG_MAINTAINER=${PKG_MAINTAINER} \
    PKG_MAINTAINER_EMAIL=${PKG_MAINTAINER_EMAIL} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_DESCRIPTION=${PKG_DESCRIPTION} \
    PKG_ARCH=${TARGETARCH}
WORKDIR /root
RUN dnf install -y make git rpmdevtools rpmlint
COPY ./Makefile ./artifact/bin/iptv-toolkit-${TARGETOS}-${TARGETARCH} /root/
RUN make -e build-rpm

FROM scratch AS package
COPY --from=builder /root/rpmbuild/RPMS/*/* /
ENTRYPOINT ["/iptv-toolkit-*"]
