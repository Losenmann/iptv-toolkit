FROM redhat/ubi9:9.5 AS builder
ARG TARGETOS \
    TARGETARCH \
    PKG_VERSION=0.0.1
ENV PKG_VERSION=${PKG_VERSION} \
    PKG_ARCH=${TARGETARCH}
RUN dnf install -y make rpmdevtools rpmlint
COPY . /opt/iptv-toolkit
WORKDIR /opt/iptv-toolkit
RUN make -e build-rpm

FROM scratch AS package
COPY --from=builder /root/rpmbuild/RPMS/*/* /
ENTRYPOINT ["/iptv-toolkit-*"]