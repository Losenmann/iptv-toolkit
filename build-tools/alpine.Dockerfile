FROM alpine:3.21.2 AS builder
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
    PKG_SIGN_ALPINE
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_HOMEGIT=${PKG_HOMEGIT} \
    PKG_HOMEPAGE=${PKG_HOMEPAGE} \
    PKG_DESCRIPTION="${PKG_DESCRIPTION}" \
    PKG_CHANGELOG="${PKG_CHANGELOG}" \
    PACKAGER="${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}>" \
    MAINTAINER="${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}>" \
    PKG_SIGN_ALPINE=${PKG_SIGN_ALPINE} \
    PKG_ARCH=${TARGETARCH}
RUN apk add alpine-sdk atools abuild-rootbld doas
COPY . /opt/pkg
WORKDIR /opt/pkg
RUN make build-apk

FROM scratch AS package
COPY --from=builder /home/alpine/*/*/apkbuild/packages/*/*/*.apk /
ENTRYPOINT ["/iptv-toolkit-*"]