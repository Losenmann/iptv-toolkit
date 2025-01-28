FROM alpine:3.21.2 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION=0.0.1 \
    PKG_LICENSE=None \
    PKG_HOME_URL=https://example.com \
    PKG_DESCRIPTION=None \
    PKG_MAINTAINER=None \
    PKG_MAINTAINER_EMAIL=example@example.com \
    PKG_SIGN_ALPINE="" \
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_HOME_URL=${PKG_HOME_URL} \
    PKG_DESCRIPTION=${PKG_DESCRIPTION} \
    PACKAGER=${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}> \
    MAINTAINER=${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}> \
    PKG_SIGN_ALPINE=${PKG_SIGN_ALPINE}
    PKG_ARCH=${TARGETARCH}
RUN apk add alpine-sdk atools abuild-rootbld doas
COPY . /opt/pkg
WORKDIR /opt/pkg
RUN make -e build-apk

FROM scratch AS package
COPY --from=builder /home/alpine/*/*/apkbuild/packages/*/*/*.apk /
ENTRYPOINT ["/iptv-toolkit-*"]