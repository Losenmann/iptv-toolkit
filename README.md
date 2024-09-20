# IPTV Toolkit
[![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)](https://www.python.org)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)
[![GitHub License](https://img.shields.io/github/license/losenmann/iptv-toolkit?style=for-the-badge)](https://github.com/Losenmann/iptv-toolkit/blob/master/LICENSE)
[![Docker Image Version](https://img.shields.io/docker/v/losenmann/iptv-toolkit?style=for-the-badge&label=Docker&color=%231D63ED)](https://hub.docker.com/r/losenmann/iptv-toolkit/tags)

## Overview
The service allows you to convert iptv and epg playlists. For playlists, conversion is carried out in the following formats: m3u, m3u8, xml. EPG conversion is carried out in the following formats: jtv, xml, xml.gz.<br>
The service supports automatic detection of the EPG type and playlists.<br>
Also, there is support for creating a playlist with a built-in udpxy link: from `udp://@238.0.0.0:1234` to `http://udpxy.local:4022/udp/238.0.0.0:1234` There is support for integrating links to EPG into the playlist file.<br>
Also in the background, the udpxy proxy server is launched, allowing you to convert Multicast traffic to HLS. Along with this, a simple file service is launched to download playlists and EPG.

> [!NOTE]
> Server file access endpoint: `/iptv/`. For UDPXY: `/udp/` and `/status/`

## Quick start
+ [Docker Compose](./docker-compose.yaml)
  ```bash
  docker-compose up -d -f ./docker-compose.yaml
  ```

+ [Kubernetes](./iptv-toolkit.yaml)
  ```bash
  kubectl apply -n media -f https://raw.githubusercontent.com/Losenmann/iptv-toolkit/refs/heads/master/iptv-toolkit.yaml
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
