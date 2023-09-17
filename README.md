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
  -p 8080:8080 \
  ghcr.io/davidborzek/docker-exporter:latest
```

### Prometheus config

Once you have configured the exporter, update your `prometheus.yml` scrape config:

```yaml
scrape_configs:
  - job_name: "docker_exporter"
    static_configs:
      - targets: ["localhost:8080"]
```

### Config

| Flag           | Description                                                                                        | Default Value | Environment Variable         |
| -------------- | -------------------------------------------------------------------------------------------------- | ------------- | ---------------------------- |
| `--port`       | The port of docker exporter server.                                                                | `8080`        | `DOCKER_EXPORTER_PORT`       |
| `--host`       | The host of docker exporter server.                                                                |               | `DOCKER_EXPORTER_HOST`       |
| `--auth-token` | Optional auth token for the docker exporter server. If no token is set authentication is disabled. |               | `DOCKER_EXPORTER_AUTH_TOKEN` |
| `--log-level`  | Log level for the exporter.                                                                        | `info`        | `DOCKER_EXPORTER_LOG_LEVEL`  |

### Exported Metrics

Currently the exporter exports all numeric and boolean states of a sensor into its own gauge:

| Metric Name                                 | Description                    | Labels        |
| ------------------------------------------- | ------------------------------ | ------------- |
| docker_container_block_io_read_bytes        | Block I/O read bytes total     | name          |
| docker_container_block_io_write_bytes       | Block I/O write bytes total    | name          |
| docker_container_cpu_usage_percentage       | CPU usage in percentage        | name          |
| docker_container_memory_total_bytes         | Total memory in bytes          | name          |
| docker_container_memory_usage_bytes         | Memory usage in bytes          | name          |
| docker_container_memory_usage_percentage    | Memory usage in percentage     | name          |
| docker_container_network_rx_bytes           | Network received bytes total   | name, network |
| docker_container_network_rx_dropped_packets | Network dropped packets total  | name, network |
| docker_container_network_rx_errors          | Network received errors        | name, network |
| docker_container_network_rx_packets         | Network received packets total | name, network |
| docker_container_network_tx_bytes           | Network sent bytes total       | name, network |
| docker_container_network_tx_dropped_packets | Network dropped packets total  | name, network |
| docker_container_network_tx_errors          | Network sent errors            | name, network |
| docker_container_network_tx_packets         | Network sent packets total     | name, network |
| docker_container_state                      | State of the container         | name, state   |
