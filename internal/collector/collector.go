package collector

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/davidborzek/docker-exporter/internal/clock"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type DockerCollector struct {
	client *client.Client
	clock  clock.Clock
}

func NewDockerCollector(clk clock.Clock) (*DockerCollector, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerCollector{
		client: client,
		clock:  clk,
	}, nil
}

func NewWithClient(client *client.Client, clk clock.Clock) *DockerCollector {
	return &DockerCollector{
		client: client,
		clock:  clk,
	}
}

func (c *DockerCollector) Describe(_ chan<- *prometheus.Desc) {}

func (c *DockerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	containers, err := c.client.ContainerList(
		ctx,
		types.ContainerListOptions{
			All: true,
		},
	)

	if err != nil {
		log.WithError(err).
			Error("failed to fetch container list")
		return
	}

	var wg sync.WaitGroup

	for _, container := range containers {
		wg.Add(1)
		go c.collectContainerMetrics(ctx, container, ch, &wg)
	}

	wg.Wait()
}

func (c *DockerCollector) collectContainerMetrics(ctx context.Context, container types.Container, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	name := containerName(container)

	ch <- prometheus.MustNewConstMetric(
		containerStateMetric, prometheus.GaugeValue, 1, name, container.State,
	)

	if container.State != "running" {
		return
	}

	inspect, err := c.client.ContainerInspect(ctx, container.ID)
	if err != nil {
		log.WithError(err).WithField("id", container.ID).
			Error("error inspecting container")
		return
	}

	ch <- prometheus.MustNewConstMetric(containerUptime,
		prometheus.GaugeValue,
		c.calculateUptime(inspect),
		name,
	)

	stats, err := c.containerStats(ctx, container.ID)
	if err != nil {
		log.WithError(err).WithField("id", container.ID).
			Error("error getting stats for container")
		return
	}

	c.cpuMetrics(ch, name, stats)
	c.memoryMetrics(ch, name, stats)
	c.networkMetrics(ch, name, stats)
	c.blockIOMetrics(ch, name, stats)
}

func (c *DockerCollector) cpuMetrics(ch chan<- prometheus.Metric, name string, stats *types.StatsJSON) {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	onlineCPUs := float64(stats.CPUStats.OnlineCPUs)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	cpuPercent := 0.0
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	ch <- prometheus.MustNewConstMetric(cpuUsagePercentage,
		prometheus.GaugeValue,
		cpuPercent,
		name,
	)
}

func (c *DockerCollector) memoryMetrics(ch chan<- prometheus.Metric, name string, stats *types.StatsJSON) {
	mem := calculateMemUsageUnixNoCache(stats.MemoryStats)
	memLimit := float64(stats.MemoryStats.Limit)

	memPercent := 0.0
	if memLimit > 0 {
		memPercent = mem / memLimit * 100.0
	}

	ch <- prometheus.MustNewConstMetric(memoryTotalBytes,
		prometheus.GaugeValue,
		memLimit,
		name,
	)

	ch <- prometheus.MustNewConstMetric(memoryUsageBytes,
		prometheus.GaugeValue,
		mem,
		name,
	)

	ch <- prometheus.MustNewConstMetric(memoryUsagePercentage,
		prometheus.GaugeValue,
		memPercent,
		name,
	)
}

func (c *DockerCollector) networkMetrics(ch chan<- prometheus.Metric, name string, stats *types.StatsJSON) {
	for networkName, network := range stats.Networks {
		ch <- prometheus.MustNewConstMetric(networkRxBytes,
			prometheus.GaugeValue,
			float64(network.RxBytes),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkRxPackets,
			prometheus.GaugeValue,
			float64(network.RxPackets),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkRxDroppedPackets,
			prometheus.GaugeValue,
			float64(network.RxDropped),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkRxErrors,
			prometheus.GaugeValue,
			float64(network.RxErrors),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkTxBytes,
			prometheus.GaugeValue,
			float64(network.TxBytes),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkTxPackets,
			prometheus.GaugeValue,
			float64(network.TxPackets),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkTxDroppedPackets,
			prometheus.GaugeValue,
			float64(network.TxDropped),
			name, networkName,
		)

		ch <- prometheus.MustNewConstMetric(networkTxErrors,
			prometheus.GaugeValue,
			float64(network.TxErrors),
			name, networkName,
		)
	}
}

func (c *DockerCollector) blockIOMetrics(ch chan<- prometheus.Metric, name string, stats *types.StatsJSON) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range stats.BlkioStats.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}

	ch <- prometheus.MustNewConstMetric(blockIOReadBytes,
		prometheus.GaugeValue,
		float64(blkRead),
		name,
	)

	ch <- prometheus.MustNewConstMetric(blockIOWriteBytes,
		prometheus.GaugeValue,
		float64(blkWrite),
		name,
	)
}

// containerStats gets the stats of a single containers.
func (c *DockerCollector) containerStats(ctx context.Context, containerID string) (*types.StatsJSON, error) {
	r, err := c.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}

	var stats types.StatsJSON
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, err
}

func (c *DockerCollector) calculateUptime(container types.ContainerJSON) float64 {
	startTime, err := c.clock.Parse(time.RFC3339Nano, container.State.StartedAt)
	if err != nil {
		return 0
	}

	return c.clock.Since(startTime).Seconds()
}

func calculateMemUsageUnixNoCache(mem types.MemoryStats) float64 {
	if v, isCgroup1 := mem.Stats["total_inactive_file"]; isCgroup1 && v < mem.Usage {
		return float64(mem.Usage - v)
	}
	if v := mem.Stats["inactive_file"]; v < mem.Usage {
		return float64(mem.Usage - v)
	}
	return float64(mem.Usage)
}

// containerName returns the first name of a container
// without the leading slash.
func containerName(c types.Container) string {
	if len(c.Names) == 0 {
		return ""
	}

	return strings.TrimLeft(c.Names[0], "/")
}
