# IPTV Toolkit
[![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com)
[![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![Maintainer](https://img.shields.io/badge/MAINTAINER-%40Losenmann-red?style=for-the-badge)](https://github.com/Losenmann)
[![GitHub License](https://img.shields.io/github/license/losenmann/iptv-toolkit?style=for-the-badge)](https://github.com/Losenmann/iptv-toolkit/blob/master/LICENSE)
[![Docker Image Version](https://img.shields.io/docker/v/losenmann/iptv-toolkit?style=for-the-badge&label=Docker&color=%231D63ED)](https://hub.docker.com/r/losenmann/iptv-toolkit/tags)


## Overview
A set of tools for working with IPTV. The set includes: playlist converter, program guide converter, small web server for publishing files, scheduler and udpxy.


## Quick start
+ [Docker Compose](./deploy/docker-compose.yaml)
  ```bash
  docker-compose -f ./docker-compose.yaml up -d
  ```

+ [Docker Swarm](./deploy/docker-stack.yaml)
  ```bash
  docker stack deploy -c ./docker-stack.yaml iptv -d --prune
  ```

+ [Kubernetes](./deploy/kubernetes.yaml)
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/Losenmann/iptv-toolkit/refs/heads/master/deploy/kubernetes.yaml
  ```

## Downloads
+ [Latest](https://github.com/Losenmann/iptv-toolkit/releases/latest)
+ [Old](https://github.com/Losenmann/iptv-toolkit/releases)