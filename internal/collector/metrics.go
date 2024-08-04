package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	containerStateMetric = prometheus.NewDesc(
		"docker_container_state",
		"State of the container",
		[]string{"name", "state"},
		nil,
	)

	containerInfo = prometheus.NewDesc(
		"docker_container_info",
		"Infos about the container",
		[]string{"name", "image_name", "image"},
		nil,
	)

	containerUptime = prometheus.NewDesc(
		"docker_container_uptime",
		"Uptime of the container in seconds",
		[]string{"name"},
		nil,
	)

	scrapeDuration = prometheus.NewDesc(
		"docker_exporter_scrape_duration",
		"Duration of the scrape in seconds",
		nil,
		nil,
	)

	scrapeErrors = prometheus.NewDesc(
		"docker_exporter_scrape_errors",
		"Number of scrape errors",
		nil,
		nil,
	)

	/*
		CPU Metrics
	*/

	cpuUsagePercentage = prometheus.NewDesc(
		"docker_container_cpu_usage_percentage",
		"CPU usage in percentage",
		[]string{"name"},
		nil,
	)

	cpuOnlineCPUs = prometheus.NewDesc(
		"docker_container_cpu_online_cpus",
		"Number of online CPUs",
		[]string{"name"},
		nil,
	)

	/*
		Memory Metrics
	*/

	memoryUsageBytes = prometheus.NewDesc(
		"docker_container_memory_usage_bytes",
		"Memory usage in bytes",
		[]string{"name"},
		nil,
	)

	memoryTotalBytes = prometheus.NewDesc(
		"docker_container_memory_total_bytes",
		"Total memory in bytes",
		[]string{"name"},
		nil,
	)

	memoryUsagePercentage = prometheus.NewDesc(
		"docker_container_memory_usage_percentage",
		"Memory usage in percentage",
		[]string{"name"},
		nil,
	)

	/*
		Network Metrics
	*/

	networkRxBytes = prometheus.NewDesc(
		"docker_container_network_rx_bytes",
		"Network received bytes total",
		[]string{"name", "network"},
		nil,
	)

	networkRxPackets = prometheus.NewDesc(
		"docker_container_network_rx_packets",
		"Network received packets total",
		[]string{"name", "network"},
		nil,
	)

	networkRxDroppedPackets = prometheus.NewDesc(
		"docker_container_network_rx_dropped_packets",
		"Network dropped packets total",
		[]string{"name", "network"},
		nil,
	)

	networkRxErrors = prometheus.NewDesc(
		"docker_container_network_rx_errors",
		"Network received errors",
		[]string{"name", "network"},
		nil,
	)

	networkTxBytes = prometheus.NewDesc(
		"docker_container_network_tx_bytes",
		"Network sent bytes total",
		[]string{"name", "network"},
		nil,
	)

	networkTxPackets = prometheus.NewDesc(
		"docker_container_network_tx_packets",
		"Network sent packets total",
		[]string{"name", "network"},
		nil,
	)

	networkTxDroppedPackets = prometheus.NewDesc(
		"docker_container_network_tx_dropped_packets",
		"Network dropped packets total",
		[]string{"name", "network"},
		nil,
	)

	networkTxErrors = prometheus.NewDesc(
		"docker_container_network_tx_errors",
		"Network sent errors",
		[]string{"name", "network"},
		nil,
	)

	/*
		BlockIO Metrics
	*/

	blockIOReadBytes = prometheus.NewDesc(
		"docker_container_block_io_read_bytes",
		"Block I/O read bytes total",
		[]string{"name"},
		nil,
	)

	blockIOWriteBytes = prometheus.NewDesc(
		"docker_container_block_io_write_bytes",
		"Block I/O write bytes total",
		[]string{"name"},
		nil,
	)

	/*
		PIDs Metrics
	*/

	pidsCurrent = prometheus.NewDesc(
		"docker_container_pids_current",
		"Current number of pids",
		[]string{"name"},
		nil,
	)
)
