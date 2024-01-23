package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/davidborzek/docker-exporter/internal/clock"
	"github.com/davidborzek/docker-exporter/internal/collector"
	"github.com/davidborzek/docker-exporter/internal/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

const (
	version = "v0.2.0"
)

var (
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Value:   "8080",
			Usage:   "The port of docker exporter server",
			EnvVars: []string{"DOCKER_EXPORTER_PORT"},
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "The host of docker exporter server",
			EnvVars: []string{"DOCKER_EXPORTER_HOST"},
		},
		&cli.StringFlag{
			Name:    "auth-token",
			Usage:   "Optional auth token for the docker exporter server. If no token is set authentication is disabled.",
			EnvVars: []string{"DOCKER_EXPORTER_AUTH_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level",
			Value:   "info",
			EnvVars: []string{"DOCKER_EXPORTER_LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:    "ignore-label",
			Usage:   "Label to ignore containers",
			Value:   "docker-exporter.ignore",
			EnvVars: []string{"DOCKER_EXPORTER_IGNORE_LABEL"},
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

func start(ctx *cli.Context) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(
		parseLogLevel(ctx.String("log-level")),
	)

	log.WithField("pid", os.Getpid()).
		Info("docker prometheus exporter started")

	token := ctx.String("auth-token")
	if len(token) > 0 {
		log.Info("authentication is enabled")
	}

	dc, err := collector.NewDockerCollector(clock.NewClock(), ctx.String("ignore-label"))
	if err != nil {
		log.WithError(err).
			Fatal("failed to create docker collector")
	}

	prometheus.MustRegister(dc)

	h := handler.New(token)

	addr := net.JoinHostPort(
		ctx.String("host"), ctx.String("port"))
	log.WithField("addr", addr).
		Infof("starting the http server")

	return http.ListenAndServe(addr, h)
}

func Main(args []string) {
	app := cli.App{
		Name:    "Docker Prometheus exporter",
		Usage:   "Export Docker metrics to prometheus format",
		Action:  start,
		Flags:   flags,
		Version: version,
	}

	if err := app.Run(args); err != nil {
		fmt.Println(err.Error())
	}
}
