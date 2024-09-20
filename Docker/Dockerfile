FROM alpine:3.20.3

COPY iptv.py \
    Docker/default.conf \
    Docker/xmltv.dtd \
    /opt/iptv-toolkit/

RUN apk add \
        --no-cache \
        --clean-protected \
        --upgrade \
        git \
        make \
        gcc \
        libc-dev \
        python3 \
        py3-magic \
        py3-schedule \
        py3-setuptools \
        libxml2-utils \
        nginx \
    && mkdir -p /www/iptv/tvguide /www/iptv/tvrecord /www/iptv/playlist /opt/iptv-toolkit/src \
    && rm -rf /etc/nginx/http.d/* \
    && git -C /opt/iptv-toolkit/src clone --branch ${ARG_UDPXY_VERSION:-master} https://github.com/pcherenkov/udpxy.git \
        && make -C /opt/iptv-toolkit/src/udpxy/chipmunk \
        && make -C /opt/iptv-toolkit/src/udpxy/chipmunk install \
    && git -C /opt/iptv-toolkit/src clone --branch ${ARG_JTV2XMLTV_VERSION:-master} https://github.com/tataranovich/jtv2xmltv.git \
        && cd /opt/iptv-toolkit/src/jtv2xmltv \
        && python3 setup.py install \
    && mv /opt/iptv-toolkit/default.conf /etc/nginx/http.d/ \
    && chmod +x /opt/iptv-toolkit/iptv.py \
    && apk del git make gcc libc-dev py3-setuptools \
    && rm -rf /opt/iptv-toolkit/src

EXPOSE 4022/TCP

WORKDIR /opt/iptv-toolkit

ENTRYPOINT ["./iptv.py"]