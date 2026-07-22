FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w -X main.version=${VERSION}" -o /docker-exporter .

FROM gcr.io/distroless/static:nonroot
COPY --from=build /docker-exporter /docker-exporter
USER nonroot:nonroot
ENTRYPOINT ["/docker-exporter"]
