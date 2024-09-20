# IPTV Toolkit
[![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)](https://www.python.org)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)

## Overview
The service allows you to convert iptv and epg playlists. For playlists, conversion is carried out in the following formats: m3u, m3u8, xml. EPG conversion is carried out in the following formats: jtv, xml, xml.gz.<br>
The service supports automatic detection of the EPG type and playlists.<br>
Also, there is support for creating a playlist with a built-in udpxy link: from `udp://@238.0.0.0:1234` to `http://udpxy.local:4022/udp/238.0.0.0:1234` There is support for integrating links to EPG into the playlist file.<br>
Also in the background, the udpxy proxy server is launched, allowing you to convert Multicast traffic to HLS. Along with this, a simple file service is launched to download playlists and EPG.

> [!NOTE]
> Server file access endpoint: "/iptv/". For UDPXY: "/udp/" and "/status/"

## Quick start
<details><summary>Docker Compose</summary>

```yaml
version: "3.9"
services:
  iptv-toolkit:
    container_name: "IPTV-Toolkit"
    image: "losenmann/iptv-toolkit:latest"
    networks:
      - network-multicast
    expose:
      - 4022
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
    environment:
      - EPG_URL="http://localhost/epg.xml"
      - PLAYLIST_URL="http://localhost/playlist.m3u"
      - PLAYLIST_TVG_URL="http://localhost/epg.xml"
      - PLAYLIST_UDPXY_URL="http://udpxy.local:4022"
    tty: true
    restart: on-failure:3

networks:
  network-multicast:
    name: "network-multicast"
    driver: "macvlan"
    driver_opts:
      parent: "eth0"
    ipam:
      config:
        - subnet: "192.168.8.0/24"
          gateway: "192.168.8.1"
          ip_range: "192.168.8.24/29"
```

</details>

<details><summary>kubernetes</summary>

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iptv-toolkit
  namespace: default
  labels:
    app: iptv-toolkit
    name: iptv-toolkit
spec:
  replicas: 1
  selector:
    matchLabels:
      app: iptv-toolkit
      task: iptv-toolkit
  template:
    metadata:
      labels:
        app: iptv-toolkit
        task: iptv-toolkit
    spec:
      hostNetwork: true
      terminationGracePeriodSeconds: 0
      containers:
        - name: iptv-toolkit
          image: losenmann/iptv-toolkit:latest
          ports:
            - containerPort: 4022
          env:
            - name: EPG_URL
              value: "http://localhost/epg.xml"
            - name: PLAYLIST_URL
              value: "http://localhost/playlist.m3u"
            - name: PLAYLIST_TVG_URL
              value: "http://localhost/epg.xml"
            - name: PLAYLIST_UDPXY_URL
              value: "http://udpxy.local:4022"
```

</details>

## Environment Variables
* `EPG_URL` - Link to tv guide
* `PLAYLIST_URL` - Link to channels playlist
* `PLAYLIST_TVG_URL` - Link to epg integrated into playlist (usually it is iptv-toolkit address)
* `PLAYLIST_UDPXY_URL` - Create a playlist with a formatted udp link in udpxy format

> [!IMPORTANT]
> Environment variables repeat CLI<br>
> CLI key take precedence over environment variables
