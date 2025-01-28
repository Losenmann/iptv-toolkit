FROM debian:12.9 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION=0.0.1 \
    PKG_LICENSE=None \
    PKG_HOME_URL=https://example.com \
    PKG_DESCRIPTION=None \
    PKG_MAINTAINER=example \
    PKG_MAINTAINER_EMAIL=example@example.com \
    PKG_CHANGELOG \
    PKG_REVISION
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_HOME_URL=${PKG_HOME_URL} \
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