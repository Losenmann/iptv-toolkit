FROM --platform=$BUILDPLATFORM alpine:3.21.2 AS builder-main
ARG TARGETOS
ARG TARGETARCH
COPY ./Makefile /opt/pkg
RUN apk add make
WORKDIR /opt/pkg
ENTRYPOINT ["make"]
CMD ["build-apk"]