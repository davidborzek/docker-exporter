# docker exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/davidborzek/docker-exporter)](https://goreportcard.com/report/github.com/davidborzek/docker-exporter)

Simple and lightweight Prometheus exporter for docker container metrics.

## Prerequisites

- [Go](https://golang.org/doc/)

## Installation

### Using Docker

The exporter is available as a Docker image.
You can run it using the following example:

```
$ docker run \
  -u root \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -p 8080:8080 \
  ghcr.io/davidborzek/docker-exporter:latest
```

> Note: To run Docker Exporter, you'll need to mount the Docker socket from your host system. This operation necessitates root privileges or the user running the command to be a member of the Docker group. It's important to note that mounting the Docker socket grants the container unrestricted access to Docker. For a more secure approach, consider utilizing the Docker Socket Proxy, which is further explained below for additional information.

### Running with [docker-socket-proxy](https://github.com/Tecnativa/docker-socket-proxy)

```
$ docker run \
  -e "DOCKER_HOST=tcp://localhost:2375" \
  -p 8080:8080 \
  ghcr.io/davidborzek/docker-exporter:latest
```

> Note: the [docker-socket-proxy](https://github.com/Tecnativa/docker-socket-proxy#not-always-needed) needs to have container access enabled. (`CONTAINERS=1`)

### Prometheus config

Once you have configured the exporter, update your `prometheus.yml` scrape config:

```yaml
scrape_configs:
  - job_name: "docker_exporter"
    static_configs:
      - targets: ["localhost:8080"]
```

### Config

| Flag             | Description                                                                                        | Default Value            | Environment Variable           |
| ---------------- | -------------------------------------------------------------------------------------------------- | ------------------------ | ------------------------------ |
| `--port`         | The port of docker exporter server.                                                                | `8080`                   | `DOCKER_EXPORTER_PORT`         |
| `--host`         | The host of docker exporter server.                                                                |                          | `DOCKER_EXPORTER_HOST`         |
| `--auth-token`   | Optional auth token for the docker exporter server. If no token is set authentication is disabled. |                          | `DOCKER_EXPORTER_AUTH_TOKEN`   |
| `--log-level`    | Log level for the exporter.                                                                        | `info`                   | `DOCKER_EXPORTER_LOG_LEVEL`    |
| `--ignore-label` | Set the label name for ignoring docker containers. (See [Ignoring Containers](#ignoring-containers))   | `docker-exporter.ignore` | `DOCKER_EXPORTER_IGNORE_LABEL` |

### Exported Metrics

| Metric Name                                 | Description                    | Labels                  |
| ------------------------------------------- | ------------------------------ | ----------------------- |
| docker_container_block_io_read_bytes        | Block I/O read bytes total     | name                    |
| docker_container_block_io_write_bytes       | Block I/O write bytes total    | name                    |
| docker_container_cpu_usage_percentage       | CPU usage in percentage        | name                    |
| docker_container_info                       | Infos about the container      | name, image_name, image |
| docker_container_memory_total_bytes         | Total memory in bytes          | name                    |
| docker_container_memory_usage_bytes         | Memory usage in bytes          | name                    |
| docker_container_memory_usage_percentage    | Memory usage in percentage     | name                    |
| docker_container_network_rx_bytes           | Network received bytes total   | name, network           |
| docker_container_network_rx_dropped_packets | Network dropped packets total  | name, network           |
| docker_container_network_rx_errors          | Network received errors        | name, network           |
| docker_container_network_rx_packets         | Network received packets total | name, network           |
| docker_container_network_tx_bytes           | Network sent bytes total       | name, network           |
| docker_container_network_tx_dropped_packets | Network dropped packets total  | name, network           |
| docker_container_network_tx_errors          | Network sent errors            | name, network           |
| docker_container_network_tx_packets         | Network sent packets total     | name, network           |
| docker_container_pids_current               | Current number of pids         | name                    |
| docker_container_state                      | State of the container         | name, state             |
| docker_container_uptime                     | Uptime of the container        | name                    |

### Ignoring Containers

You can ignore containers by setting the label `docker-exporter.ignore` on the container. The label name can be configured with the `--ignore-label` flag.

```yaml
services:
  nginx:
    image: nginx
    labels:
      docker-exporter.ignore: "true"
```
