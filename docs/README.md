# IPTV Toolkit
[![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)
[![GitHub License](https://img.shields.io/github/license/losenmann/iptv-toolkit?style=for-the-badge)](https://github.com/Losenmann/iptv-toolkit/blob/master/LICENSE)
[![Docker Image Version](https://img.shields.io/docker/v/losenmann/iptv-toolkit?style=for-the-badge&label=Docker&color=%231D63ED)](https://hub.docker.com/r/losenmann/iptv-toolkit/tags)


## Overview
The service allows you to convert iptv and epg playlists. For playlists, conversion is carried out in the following formats: m3u, m3u8, xml. EPG conversion is carried out in the following formats: jtv, xml, xml.gz.<br>
The service supports automatic detection of the EPG type and playlists.<br>
Also, there is support for creating a playlist with a built-in udpxy link: from `udp://@238.0.0.0:1234` to `http://udpxy.local:4022/udp/238.0.0.0:1234` There is support for integrating links to EPG into the playlist file.<br>
Also in the background, the udpxy proxy server is launched, allowing you to convert Multicast traffic to HLS. Along with this, a simple file service is launched to download playlists and EPG.

> [!NOTE]
> Server file access endpoint: "/files/". For UDPXY: "/udp/" and "/status/".

### Supported tools
+ Playlist converter in 2 formats ([xspf](https://xspf.org), [m3u/m3u8](https://wikipedia.org/wiki/M3U))
+ EPG converter in 3 formats (jtv, [xmltv/xmltv.gz](https://xmltv.org))
+ Scheduler
+ Web server for converted files
+ [UDPXY](https://github.com/pcherenkov/udpxy) proxy

## Quick start
+ [Docker Compose](./deploy/docker-compose.yaml)
  ```bash
  docker-compose up -d -f ./docker-compose.yaml
  ```

+ [Kubernetes](./deploy/kubernetes.yaml)
  ```bash
  kubectl apply -n media -f https://raw.githubusercontent.com/Losenmann/iptv-toolkit/refs/heads/master/deploy/kubernetes.yaml
  ```
> [!IMPORTANT]
> For Docker, need to edit the [.env](./deploy/.env) file in place next to docker-compose.<br>
> For Kubernetes, need to edit and apply the [ConfigMap](./deploy/kubernetes-configmap.yaml) before deployment. Also edit the Ingress resource host, for example, via kubectl edit.

## Environment Variables and CLI Key
Environmental variables and key CLI applicable in all operating modes.

| Variables | Key  | Default | Description |
| :-------- | :--: | :-----: | :---------- |
| IPTVTOOLKIT_EPG | -e | none | Link to tv guide |
| IPTVTOOLKIT_EPG_DST | -E | ./files/tvguide | Path of export of EPGs |
| IPTVTOOLKIT_PLAYLIST | -p | none | Link to channels playlist |
| IPTVTOOLKIT_PLAYLIST_DST | -P | ./files/playlist | Path of export of Playlists |
| IPTVTOOLKIT_PLAYLIST_UDPXY | -u | none | Create a playlist with a formatted<br> udp link in udpxy format |
| IPTVTOOLKIT_PLAYLIST_EMBED_EPG | -i | none | Link to epg integrated into playlist |
| IPTVTOOLKIT_WEB_PATH | -f | ./files | Path to display files by web server |
| IPTVTOOLKIT_WEB_PORT | -P | 4023 | Web server port |
| IPTVTOOLKIT_CRONTAB | -c | 30 6 * * * | Сrontab style task schedule |
| none | -S | false | Start Schedule job |
| none | -U | false | Start UDPXY proxy |
| none | -W | false | Start Web Server |

> [!IMPORTANT]
> Environment variables repeat CLI.<br>
> CLI key take precedence over environment variables.

## Downloads
### Binary
+ Linux i386 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-386)
+ Linux amd64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-amd64)
+ Linux arm [downloads](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-arm)
+ Linux arm64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-arm64)
+ Linux ppc64le [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-ppc64le)
+ Linux riscv64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-riscv64)
+ Linux s390x [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-s390x)

### Package
+ RPM amd64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-amd64)
+ RPM arm64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-arm64)
+ RPM ppc64le [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-ppc64le)
+ RPM s390x [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-s390x)
+ DEB i386 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-386)
+ DEB amd64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-amd64)
+ DEB arm [downloads](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-arm)
+ DEB arm64 [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-arm64)
+ DEB ppc64le [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-ppc64le)
+ DEB s390x [download](https://github.com/losenmann/iptv-toolkit/releases/latest/download/iptv-toolkit-linux-s390x)
