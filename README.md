# IPTV Toolkit
[![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)

## Overview
The service allows you to convert iptv and epg playlists. For playlists, conversion is carried out in the following formats: m3u, m3u8, xml. EPG conversion is carried out in the following formats: jtv, xml, xml.gz.<br>
The service supports automatic detection of the EPG type and playlists.<br>
Also, there is support for creating a playlist with a built-in udpxy link: from `udp://@238.0.0.0:1234` to `http://udpxy.local:4022/udp/238.0.0.0:1234` There is support for integrating links to EPG into the playlist file.<br>
Also in the background, the udpxy proxy server is launched, allowing you to convert Multicast traffic to HLS. Along with this, a simple file service is launched to download playlists and EPG.

> [!NOTE]
> Server file access endpoint: "/files/". For UDPXY: "/udp/" and "/status/"

### Supported tools
+ Playlist converter in 2 formats (xml, m3u)
+ EPG converter in 3 formats (jtv, xml, xml.gz)
+ Scheduler
+ Web server for converted files
+ UDPXY proxy

## Quick start
+ [Docker Compose](./deploy/docker-compose.yaml)
  ```bash
  docker-compose up -d -f ./docker-compose.yaml
  ```

+ [Kubernetes](./deploy/kubernetes.yaml)
  ```bash
  kubectl apply -n media -f https://raw.githubusercontent.com/Losenmann/iptv-toolkit/refs/heads/master/deploy/kubernetes.yaml
  ```

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
| `IPTVTOOLKIT_WEB_PORT` | `-P` | `4022` | Web server port |
| `IPTVTOOLKIT_CRONTAB` | `-c` | `30 6 * * *` | Сrontab style task schedule |

> [!IMPORTANT]
> Environment variables repeat CLI<br>
> CLI key take precedence over environment variables
