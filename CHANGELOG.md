# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

From this release onward, entries are maintained automatically by
[release-please](https://github.com/googleapis/release-please) from the
Conventional Commit history.

## [0.3.0]

Baseline release documenting the current feature set:

- Prometheus exporter for Docker container metrics (CPU, memory, block I/O,
  network, pids, state, uptime, and container info).
- Optional bearer-token authentication for the metrics endpoint
  (`--auth-token`).
- Container exclusion via a configurable ignore label (`--ignore-label`,
  default `docker-exporter.ignore`).
- Configuration through flags and `DOCKER_EXPORTER_*` environment variables.
- Multi-arch container image (`linux/amd64`, `linux/arm64`).
