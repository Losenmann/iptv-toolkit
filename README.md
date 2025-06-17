# IPTV Toolkit
[![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)
[![GitHub License](https://img.shields.io/github/license/losenmann/iptv-toolkit?style=for-the-badge)](https://github.com/Losenmann/iptv-toolkit/blob/master/LICENSE)
[![Workflow](https://img.shields.io/github/actions/workflow/status/losenmann/iptv-toolkit/ci.yml?style=for-the-badge&label=Workflow
)](https://github.com/Losenmann/iptv-toolkit/actions/workflows/ci.yml)
[![Docker Image Version](https://img.shields.io/docker/v/losenmann/iptv-toolkit?style=for-the-badge&label=Docker&color=%231D63ED)](https://hub.docker.com/r/losenmann/iptv-toolkit/tags)
[![Docker Pulls](https://img.shields.io/docker/pulls/losenmann/iptv-toolkit?style=for-the-badge&label=Docker%20Pull&color=%231D63ED)](https://hub.docker.com/r/losenmann/iptv-toolkit)

## Overview
The service allows you to convert iptv and epg playlists. For playlists, conversion is carried out in the following formats: m3u, m3u8, xml. EPG conversion is carried out in the following formats: jtv, xml, xml.gz.<br>
The service supports automatic detection of the EPG type and playlists.<br>
Also, there is support for creating a playlist with a built-in udpxy link: from `udp://@238.0.0.0:1234` to `http://udpxy.local:4022/udp/238.0.0.0:1234` There is support for integrating links to EPG into the playlist file.<br>
Also in the background, the udpxy proxy server is launched, allowing you to convert Multicast traffic to HLS. Along with this, a simple file service is launched to download playlists and EPG.

> [!NOTE]
> Server file access endpoint: "/files/". For UDPXY: "/udp/".

> [!WARNING]
> Please use tags or releases, the master branch may contain unstable changes.

### Supported tools
+ Playlist converter in 2 formats ([xspf](https://xspf.org), [m3u/m3u8](https://wikipedia.org/wiki/M3U))
+ EPG converter in 3 formats (jtv, [xmltv/xmltv.gz](https://xmltv.org))
+ Scheduler
+ Web server for converted files
+ UDP-to-HTTP

## Quick start
+ [Docker Compose](./deploy/docker-compose.yaml)
  ```bash
  docker-compose up -d -f ./docker-compose.yaml
  ```

+ [Docker Swarm](./deploy/docker-stack.yaml)
  ```bash
  docker stack deploy -c ./docker-stack.yaml iptv -d --prune
  ```

+ [Kubernetes](./deploy/kubernetes.yaml)
  ```bash
  kubectl apply -n iptv -f https://raw.githubusercontent.com/Losenmann/iptv-toolkit/master/deploy/kubernetes.yaml
  ```

+ [MikroTik RoutesOS](./deploy/routeros-v7.rsc)
  ```rsc
  :execute script=[([/tool/fetch url="https://raw.githubusercontent.com/Losenmann/iptv-toolkit/master/deploy/routeros-v7.rsc" output=user as-value]->"data")]
  ```

> [!WARNING]
> **MikroTik RoutesOS**
> If you are running a container on RouterOS devices, remember that the device resource is limited.
> The service requires at least 32 MB of RAM.
> You also need to manually set the address on the [veth](https://help.mikrotik.com/docs/spaces/ROS/pages/84901929/Container#Container-ContainerinLayer2network) `veth-iptv-toolkit` interface and add the interface to the bridge.
> You also need to enable [IGMP Proxy](https://help.mikrotik.com/docs/spaces/ROS/pages/128221386/IGMP+Proxy) and configure the firewall.

> [!IMPORTANT]
> For Docker, need to edit the [.env](./deploy/.env) file in place next to docker-compose.<br>
> For Kubernetes, need to edit and apply the [ConfigMap](./deploy/kubernetes-configmap.yaml) before deployment. Also edit the Ingress resource host, for example, via kubectl edit.

## Environment Variables and CLI Key
Environmental variables and key CLI applicable in all operating modes.
| Variables | Key  | Default | Description |
| :-------- | :--: | :-----: | :---------- |
| `IPTVTOOLKIT_ADDRESS` | `-a` | `none` | The address at which the service will be available |
| `IPTVTOOLKIT_EPG` | `-e` | `none` | Tvguide location |
| `IPTVTOOLKIT_EPG_DIR` | `-E` | `./files/tvguide` | The directory in which the TV programs will be placed |
| `IPTVTOOLKIT_PLAYLIST` | `-p` | `none` | Playlist location |
| `IPTVTOOLKIT_PLAYLIST_DIR` | `-P` | `./files/playlist` | The directory in which the playlists will be placed |
| `IPTVTOOLKIT_WEB_PATH_FILES` | `-f` | `/files` | URL path file server |
| `IPTVTOOLKIT_WEB_PATH_UDXPY` | `-u` | `/udp` | URL path udpxy |
| `IPTVTOOLKIT_WEB_DIR` | `-W` | `./files` | Location on host files for web server |
| `IPTVTOOLKIT_WEB_PORT` | `-w` | `4022` | Listening port |
| `IPTVTOOLKIT_CRONTAB` | `-c` | `30 6 * * *` | Set a schedule for updating files |
| `IPTVTOOLKIT_EMBED_UDPXY` | `-d` | `false` | Create an additional playlist |
| `IPTVTOOLKIT_EMBED_EPG` | `-t` | `false` | Add the program guide to the playlist |
| `none` | `-U` | `false` | Enable udpxy converter |
| `none` | `-F` | `false` | Enable File Server |
| `none` | `-S` | `false` | Enable Scheduler |

> [!IMPORTANT]
> Environment variables repeat CLI.<br>
> CLI key take precedence over environment variables.

## Build from source
### Build executable file from source
The project provides several assembly options. The possible options are listed below:
+ [Bin](#build-bin)
+ [Docker Image](#build-docker-image)
+ [Distro Package](#build-distro-package)

> [!IMPORTANT]
> Before any build you need to install docker.

### Build bin
1. Go to the project directory
2. Run `make bin` command
3. Check catalog `./build`
After that, the image can be used as you wish.

### Build Docker Image
The built image is placed in the local registry, similar to the `docker pull` command. After that, the image can be used as you wish.
1. Go to the project directory
2. Run `make image` command
2. Check the repository `docker image ls`

Once the image build is complete, you can deploy it by running the command: `make docker-up` or `make docker-down`