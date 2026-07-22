package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/davidborzek/docker-exporter/internal/clock"
	"github.com/davidborzek/docker-exporter/internal/collector"
	"github.com/davidborzek/docker-exporter/internal/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli/v3"

	log "github.com/sirupsen/logrus"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Value:   "8080",
			Usage:   "The port of docker exporter server",
			Sources: cli.EnvVars("DOCKER_EXPORTER_PORT"),
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "The host of docker exporter server",
			Sources: cli.EnvVars("DOCKER_EXPORTER_HOST"),
		},
		&cli.StringFlag{
			Name:    "auth-token",
			Usage:   "Optional auth token for the docker exporter server. If no token is set authentication is disabled.",
			Sources: cli.EnvVars("DOCKER_EXPORTER_AUTH_TOKEN"),
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level",
			Value:   "info",
			Sources: cli.EnvVars("DOCKER_EXPORTER_LOG_LEVEL"),
		},
		&cli.StringFlag{
			Name:    "ignore-label",
			Usage:   "Label to ignore containers",
			Value:   "docker-exporter.ignore",
			Sources: cli.EnvVars("DOCKER_EXPORTER_IGNORE_LABEL"),
		},
		&cli.StringSliceFlag{
			Name:    "container-label",
			Usage:   "Docker label to expose as a `docker_container_labels` metric. Repeatable, or comma-separated via the environment variable.",
			Sources: cli.EnvVars("DOCKER_EXPORTER_CONTAINER_LABELS"),
		},
	}
)

func parseLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	}

	log.WithField("level", level).
		Warn("invalid log level provided - falling back to 'info'")

	return log.InfoLevel
}

func start(_ context.Context, cmd *cli.Command) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(
		parseLogLevel(cmd.String("log-level")),
	)

	log.WithField("pid", os.Getpid()).
		Info("docker prometheus exporter started")

	token := cmd.String("auth-token")
	if len(token) > 0 {
		log.Info("authentication is enabled")
	}

	dc, err := collector.NewDockerCollector(clock.NewClock(), cmd.String("ignore-label"), cmd.StringSlice("container-label"))
	if err != nil {
		log.WithError(err).
			Fatal("failed to create docker collector")
	}

	prometheus.MustRegister(dc)

	h := handler.New(token)

	addr := net.JoinHostPort(
		cmd.String("host"), cmd.String("port"))
	log.WithField("addr", addr).
		Infof("starting the http server")

	return http.ListenAndServe(addr, h)
}

func Main(version string, args []string) {
	cmd := &cli.Command{
		Name:    "Docker Prometheus exporter",
		Usage:   "Export Docker metrics to prometheus format",
		Action:  start,
		Flags:   flags,
		Version: version,
	}

	if err := cmd.Run(context.Background(), args); err != nil {
		fmt.Println(err.Error())
	}
}
