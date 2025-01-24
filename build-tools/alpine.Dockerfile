FROM alpine:3.21.2 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION=0.0.1 \
    PKG_LICENSE=None \
    PKG_HOME_URL=https://example.com \
    PKG_DESCRIPTION=None \
    PKG_MAINTAINER=None \
    PKG_MAINTAINER_EMAIL=example@example.com \
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_LICENSE=${PKG_LICENSE} \
    PKG_HOME_URL=${PKG_HOME_URL} \
    PKG_DESCRIPTION=${PKG_DESCRIPTION} \
    PACKAGER=${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}> \
    MAINTAINER=${PKG_MAINTAINER} <${PKG_MAINTAINER_EMAIL}> \
    PKG_ARCH=${TARGETARCH}
RUN apk add alpine-sdk atools abuild-rootbld doas \
    && adduser -D alpine \
    && adduser alpine abuild \
    && adduser alpine wheel \
    && echo "permit nopass :wheel as root" > /etc/doas.d/doas.conf
COPY . /home/alpine/
WORKDIR /home/alpine
RUN chown -R alpine:alpine /home/alpine/
USER alpine
RUN make -e build-apk

FROM scratch AS package
COPY --from=builder ~/packages/*/*/*.apk /
ENTRYPOINT ["/iptv-toolkit-*"]