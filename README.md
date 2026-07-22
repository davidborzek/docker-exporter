<div align="center">

# docker-exporter

**Simple and lightweight Prometheus exporter for Docker container metrics.**

[![ci](https://github.com/davidborzek/docker-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/davidborzek/docker-exporter/actions/workflows/ci.yml)
[![license](https://img.shields.io/github/license/davidborzek/docker-exporter)](LICENSE)
[![release](https://img.shields.io/github/v/release/davidborzek/docker-exporter)](https://github.com/davidborzek/docker-exporter/releases)

</div>

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

| Flag             | Description                                                                                          | Default Value            | Environment Variable           |
| ---------------- | ---------------------------------------------------------------------------------------------------- | ------------------------ | ------------------------------ |
| `--port`         | The port of docker exporter server.                                                                  | `8080`                   | `DOCKER_EXPORTER_PORT`         |
| `--host`         | The host of docker exporter server.                                                                  |                          | `DOCKER_EXPORTER_HOST`         |
| `--auth-token`   | Optional auth token for the docker exporter server. If no token is set authentication is disabled.   |                          | `DOCKER_EXPORTER_AUTH_TOKEN`   |
| `--log-level`    | Log level for the exporter.                                                                          | `info`                   | `DOCKER_EXPORTER_LOG_LEVEL`    |
| `--ignore-label` | Set the label name for ignoring docker containers. (See [Ignoring Containers](#ignoring-containers)) | `docker-exporter.ignore` | `DOCKER_EXPORTER_IGNORE_LABEL` |
| `--container-label` | Docker label to expose as a `docker_container_labels` metric. Repeatable. (See [Exposing Container Labels](#exposing-container-labels)) | | `DOCKER_EXPORTER_CONTAINER_LABELS` |

### Exported Metrics

| Metric Name                                 | Description                        | Labels                  |
| ------------------------------------------- | ---------------------------------- | ----------------------- |
| docker_container_block_io_read_bytes        | Block I/O read bytes total         | name                    |
| docker_container_block_io_write_bytes       | Block I/O write bytes total        | name                    |
| docker_container_cpu_usage_percentage       | CPU usage in percentage            | name                    |
| docker_container_info                       | Infos about the container          | name, image_name, image |
| docker_container_labels                     | Configured container labels (value 1)  | name, container_label_* |
| docker_container_memory_total_bytes         | Total memory in bytes              | name                    |
| docker_container_memory_usage_bytes         | Memory usage in bytes              | name                    |
| docker_container_memory_usage_percentage    | Memory usage in percentage         | name                    |
| docker_container_network_rx_bytes           | Network received bytes total       | name, network           |
| docker_container_network_rx_dropped_packets | Network dropped packets total      | name, network           |
| docker_container_network_rx_errors          | Network received errors            | name, network           |
| docker_container_network_rx_packets         | Network received packets total     | name, network           |
| docker_container_network_tx_bytes           | Network sent bytes total           | name, network           |
| docker_container_network_tx_dropped_packets | Network dropped packets total      | name, network           |
| docker_container_network_tx_errors          | Network sent errors                | name, network           |
| docker_container_network_tx_packets         | Network sent packets total         | name, network           |
| docker_container_pids_current               | Current number of pids             | name                    |
| docker_container_state                      | State of the container             | name, state             |
| docker_container_uptime                     | Uptime of the container in seconds | name                    |
| docker_exporter_scrape_duration             | Duration of the scrape in seconds  |                         |
| docker_exporter_scrape_errors               | Number of scrape errors            |                         |

### Ignoring Containers

You can ignore containers by setting the label `docker-exporter.ignore` on the container. The label name can be configured with the `--ignore-label` flag.

```yaml
services:
  nginx:
    image: nginx
    labels:
      docker-exporter.ignore: "true"
```

### Exposing Container Labels

By default no container labels are exported. Selected labels are exposed on a
dedicated `docker_container_labels` metric (value `1`), following the
`kube_pod_labels` convention: the Docker label key is prefixed with
`container_label_` and any character outside `[a-zA-Z0-9_]` becomes `_`. A
series only carries the selected labels a container actually sets — absent
labels are omitted, not exported empty.

There are two ways to select labels, and they combine (union):

1. **Globally**, for every container, via the `--container-label` flag
   (repeatable) or a comma-separated `DOCKER_EXPORTER_CONTAINER_LABELS`
   environment variable:

   ```
   $ docker-exporter --container-label com.docker.compose.project --container-label maintainer
   ```

2. **Per container**, by setting the `docker-exporter.exposed-labels` label to a
   comma-separated list of that container's own label keys to expose:

   ```yaml
   services:
     web:
       image: nginx
       labels:
         docker-exporter.exposed-labels: "com.docker.compose.project,maintainer"
   ```

Either way the result is the same metric:

```
docker_container_labels{name="web",container_label_com_docker_compose_project="shop",container_label_maintainer="acme"} 1
```

Join labels onto other metrics in PromQL via the container `name`:

```promql
docker_container_cpu_usage_percentage
  * on(name) group_left(container_label_com_docker_compose_project) docker_container_labels
```

> Exporting *all* labels is intentionally not offered. Labels are selected by an
> explicit allowlist to keep metric cardinality bounded — note the per-container
> option delegates that choice to whoever can set container labels.

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) for the
development workflow and release process, [SECURITY.md](SECURITY.md) for
reporting vulnerabilities, and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). Notable
changes are tracked in [CHANGELOG.md](CHANGELOG.md).

## License

[MIT](LICENSE)
