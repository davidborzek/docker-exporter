FROM golang:1.21.5-alpine3.17 AS base

RUN adduser -D -H docker-exporter

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go mod download

FROM base as build

RUN go build -o docker-exporter -tags prod main.go

FROM scratch as prod

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /build/docker-exporter /

USER docker-exporter:docker-exporter

CMD ["./docker-exporter"]