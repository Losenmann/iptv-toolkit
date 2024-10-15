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
+ Playlist converter in 2 formats (xml, m3u)
+ EPG converter in 3 formats (jtv, xml, xml.gz)
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
| `IPTVTOOLKIT_EPG` | `-e` | `none` | Link to tv guide |
| `IPTVTOOLKIT_EPG_DST` | `-E` | `./files/tvguide` | Path of export of EPGs |
| `IPTVTOOLKIT_PLAYLIST` | `-p` | `none` | Link to channels playlist |
| `IPTVTOOLKIT_PLAYLIST_DST` | `-P` | `./files/playlist` | Path of export of Playlists |
| `IPTVTOOLKIT_PLAYLIST_UDPXY` | `-u` | `none` | Create a playlist with a formatted<br> udp link in udpxy format |
| `IPTVTOOLKIT_PLAYLIST_EMBED_EPG` | `-i` | `none` | Link to epg integrated into playlist |
| `IPTVTOOLKIT_WEB_PATH` | `-f` | `./files` | Path to display files by web server |
| `IPTVTOOLKIT_WEB_PORT` | `-P` | `4023` | Web server port |
| `IPTVTOOLKIT_CRONTAB` | `-c` | `30 6 * * *` | Сrontab style task schedule |
| `none` | `-S` | `false` | Start Schedule job |
| `none` | `-U` | `false` | Start UDPXY proxy |
| `none` | `-W` | `false` | Start Web Server |

> [!IMPORTANT]
> Environment variables repeat CLI.<br>
> CLI key take precedence over environment variables.

## Build from source
> [!IMPORTANT]
> Before starting the build, you need to install the executable packager [upx](https://github.com/upx/upx/releases). This step can be skipped if you are only building a Docker image.

### Build executable file
1. Install the latest version of [golang](https://go.dev/dl).
2. Go to the project directory.
3. Run `make` command.

### Build Docker image
1. Go to the project directory.
2. Run `make docker` command.
> [!NOTE]
> You can pass the `ARCH_ALL=true` argument to the `make` command for cross-compilation executable files.<br>
> Example: `make build ARCH_ALL=true` or `make ARCH_ALL=true`
