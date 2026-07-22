# Security policy

## Reporting a vulnerability

Please **do not** open a public issue for security vulnerabilities.

Instead, report them privately via GitHub's
[security advisories](https://github.com/davidborzek/docker-exporter/security/advisories/new)
("Report a vulnerability"). You will receive a response as soon as possible, and
disclosure will be coordinated with you.

## Scope

docker-exporter reads container metrics from the Docker daemon, typically via
the mounted Docker socket. Mounting the socket grants the container broad access
to Docker — prefer a read-only [docker-socket-proxy](https://github.com/Tecnativa/docker-socket-proxy)
where possible, and keep any configured `--auth-token` out of version control.
