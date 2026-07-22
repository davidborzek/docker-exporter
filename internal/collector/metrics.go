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

	containerUptimeSeconds = prometheus.NewDesc(
		"docker_container_uptime_seconds",
		"Uptime of the container in seconds",
		[]string{"name"},
		nil,
	)

	scrapeDurationSeconds = prometheus.NewDesc(
		"docker_exporter_scrape_duration_seconds",
		"Duration of the scrape in seconds",
		nil,
		nil,
	)

	scrapeErrorsTotal = prometheus.NewDesc(
		"docker_exporter_scrape_errors_total",
		"Total number of scrape errors",
		nil,
		nil,
	)

	/*
		CPU Metrics
	*/

	cpuUsageSecondsTotal = prometheus.NewDesc(
		"docker_container_cpu_usage_seconds_total",
		"Total CPU time consumed in seconds",
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

	memoryLimitBytes = prometheus.NewDesc(
		"docker_container_memory_limit_bytes",
		"Memory limit in bytes",
		[]string{"name"},
		nil,
	)

	memoryUsageRatio = prometheus.NewDesc(
		"docker_container_memory_usage_ratio",
		"Memory usage as a ratio of the limit (0-1)",
		[]string{"name"},
		nil,
	)

	/*
		Network Metrics
	*/

	networkReceiveBytesTotal = prometheus.NewDesc(
		"docker_container_network_receive_bytes_total",
		"Total network bytes received",
		[]string{"name", "network"},
		nil,
	)

	networkReceivePacketsTotal = prometheus.NewDesc(
		"docker_container_network_receive_packets_total",
		"Total network packets received",
		[]string{"name", "network"},
		nil,
	)

	networkReceivePacketsDroppedTotal = prometheus.NewDesc(
		"docker_container_network_receive_packets_dropped_total",
		"Total network packets dropped while receiving",
		[]string{"name", "network"},
		nil,
	)

	networkReceiveErrorsTotal = prometheus.NewDesc(
		"docker_container_network_receive_errors_total",
		"Total network receive errors",
		[]string{"name", "network"},
		nil,
	)

	networkTransmitBytesTotal = prometheus.NewDesc(
		"docker_container_network_transmit_bytes_total",
		"Total network bytes transmitted",
		[]string{"name", "network"},
		nil,
	)

	networkTransmitPacketsTotal = prometheus.NewDesc(
		"docker_container_network_transmit_packets_total",
		"Total network packets transmitted",
		[]string{"name", "network"},
		nil,
	)

	networkTransmitPacketsDroppedTotal = prometheus.NewDesc(
		"docker_container_network_transmit_packets_dropped_total",
		"Total network packets dropped while transmitting",
		[]string{"name", "network"},
		nil,
	)

	networkTransmitErrorsTotal = prometheus.NewDesc(
		"docker_container_network_transmit_errors_total",
		"Total network transmit errors",
		[]string{"name", "network"},
		nil,
	)

	/*
		Filesystem (Block I/O) Metrics
	*/

	fsReadsBytesTotal = prometheus.NewDesc(
		"docker_container_fs_reads_bytes_total",
		"Total bytes read from block devices",
		[]string{"name"},
		nil,
	)

	fsWritesBytesTotal = prometheus.NewDesc(
		"docker_container_fs_writes_bytes_total",
		"Total bytes written to block devices",
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

// Deprecated descriptors are kept for backward compatibility and emitted
// alongside the standard metrics above. They will be removed in a future
// release; prefer the replacements named in each help string.
var (
	cpuUsagePercentage = prometheus.NewDesc(
		"docker_container_cpu_usage_percentage",
		"CPU usage in percentage (deprecated; use docker_container_cpu_usage_seconds_total)",
		[]string{"name"},
		nil,
	)

	memoryTotalBytes = prometheus.NewDesc(
		"docker_container_memory_total_bytes",
		"Total memory in bytes (deprecated; use docker_container_memory_limit_bytes)",
		[]string{"name"},
		nil,
	)

	memoryUsagePercentage = prometheus.NewDesc(
		"docker_container_memory_usage_percentage",
		"Memory usage in percentage (deprecated; use docker_container_memory_usage_ratio)",
		[]string{"name"},
		nil,
	)

	networkRxBytes = prometheus.NewDesc(
		"docker_container_network_rx_bytes",
		"Network received bytes total (deprecated; use docker_container_network_receive_bytes_total)",
		[]string{"name", "network"},
		nil,
	)

	networkRxPackets = prometheus.NewDesc(
		"docker_container_network_rx_packets",
		"Network received packets total (deprecated; use docker_container_network_receive_packets_total)",
		[]string{"name", "network"},
		nil,
	)

	networkRxDroppedPackets = prometheus.NewDesc(
		"docker_container_network_rx_dropped_packets",
		"Network dropped packets total (deprecated; use docker_container_network_receive_packets_dropped_total)",
		[]string{"name", "network"},
		nil,
	)

	networkRxErrors = prometheus.NewDesc(
		"docker_container_network_rx_errors",
		"Network received errors (deprecated; use docker_container_network_receive_errors_total)",
		[]string{"name", "network"},
		nil,
	)

	networkTxBytes = prometheus.NewDesc(
		"docker_container_network_tx_bytes",
		"Network sent bytes total (deprecated; use docker_container_network_transmit_bytes_total)",
		[]string{"name", "network"},
		nil,
	)

	networkTxPackets = prometheus.NewDesc(
		"docker_container_network_tx_packets",
		"Network sent packets total (deprecated; use docker_container_network_transmit_packets_total)",
		[]string{"name", "network"},
		nil,
	)

	networkTxDroppedPackets = prometheus.NewDesc(
		"docker_container_network_tx_dropped_packets",
		"Network dropped packets total (deprecated; use docker_container_network_transmit_packets_dropped_total)",
		[]string{"name", "network"},
		nil,
	)

	networkTxErrors = prometheus.NewDesc(
		"docker_container_network_tx_errors",
		"Network sent errors (deprecated; use docker_container_network_transmit_errors_total)",
		[]string{"name", "network"},
		nil,
	)

	blockIOReadBytes = prometheus.NewDesc(
		"docker_container_block_io_read_bytes",
		"Block I/O read bytes total (deprecated; use docker_container_fs_reads_bytes_total)",
		[]string{"name"},
		nil,
	)

	blockIOWriteBytes = prometheus.NewDesc(
		"docker_container_block_io_write_bytes",
		"Block I/O write bytes total (deprecated; use docker_container_fs_writes_bytes_total)",
		[]string{"name"},
		nil,
	)

	containerUptime = prometheus.NewDesc(
		"docker_container_uptime",
		"Uptime of the container in seconds (deprecated; use docker_container_uptime_seconds)",
		[]string{"name"},
		nil,
	)

	scrapeDuration = prometheus.NewDesc(
		"docker_exporter_scrape_duration",
		"Duration of the scrape in seconds (deprecated; use docker_exporter_scrape_duration_seconds)",
		nil,
		nil,
	)

	scrapeErrors = prometheus.NewDesc(
		"docker_exporter_scrape_errors",
		"Number of scrape errors (deprecated; use docker_exporter_scrape_errors_total)",
		nil,
		nil,
	)
)
