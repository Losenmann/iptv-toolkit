FROM debian:12.9 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION \
    PKG_LICENSE \
    PKG_HOMEGIT \
    PKG_HOMEPAGE \
    PKG_DESCRIPTION \
    PKG_MAINTAINER \
    PKG_MAINTAINER_EMAIL \
    PKG_REVISION \
    PKG_CHANGELOG \
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_HOMEGIT=${PKG_HOMEGIT} \
    PKG_HOMEPAGE=${PKG_HOMEPAGE} \
    PKG_DESCRIPTION="${PKG_DESCRIPTION}" \
    PKG_CHANGELOG="${PKG_CHANGELOG}" \
    PACKAGER="${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}>" \
    MAINTAINER="${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}>" \
    PKG_REVISION=${PKG_REVISION} \
    PKG_ARCH=${TARGETARCH}
RUN apt update && apt install -y --no-install-recommends make build-essential dpkg-dev dh-make
COPY . /opt/pkg
WORKDIR /opt/pkg
RUN make build-deb

FROM scratch AS package
COPY --from=builder /opt/pkg/pkg/debbuild/*.deb /
ENTRYPOINT ["/iptv-toolkit-*"]