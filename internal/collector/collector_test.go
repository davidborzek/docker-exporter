package collector_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/davidborzek/docker-exporter/internal/collector"
	"github.com/davidborzek/docker-exporter/internal/mock"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.uber.org/mock/gomock"
)

const ignoreLabel = "docker-exporter.ignore"

func TestCollectMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().
		Now().
		Return(time.Now()).
		Times(1)

	mockClock.EXPECT().
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)

	// The first call to Since is for the uptime of the container
	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(1 * time.Second).
		Times(1)

	// The second call to Since is for the scrape duration
	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(2 * time.Second).
		Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	const expected = `
	# HELP docker_container_cpu_online_cpus Number of online CPUs
	# TYPE docker_container_cpu_online_cpus gauge
	docker_container_cpu_online_cpus{name="testName"} 4
	# HELP docker_container_cpu_usage_seconds_total Total CPU time consumed in seconds
	# TYPE docker_container_cpu_usage_seconds_total counter
	docker_container_cpu_usage_seconds_total{name="testName"} 8.888e-06
	# HELP docker_container_fs_reads_bytes_total Total bytes read from block devices
	# TYPE docker_container_fs_reads_bytes_total counter
	docker_container_fs_reads_bytes_total{name="testName"} 9999
	# HELP docker_container_fs_writes_bytes_total Total bytes written to block devices
	# TYPE docker_container_fs_writes_bytes_total counter
	docker_container_fs_writes_bytes_total{name="testName"} 7777
	# HELP docker_container_info Infos about the container
	# TYPE docker_container_info gauge
	docker_container_info{image="sha256:d3751d33f9cd5049c4af2b462735457e4d3baf130bcbb87f389e349fbaeb20b9",image_name="myImage",name="testName"} 1
	# HELP docker_container_memory_limit_bytes Memory limit in bytes
	# TYPE docker_container_memory_limit_bytes gauge
	docker_container_memory_limit_bytes{name="testName"} 8e+09
	# HELP docker_container_memory_usage_bytes Memory usage in bytes
	# TYPE docker_container_memory_usage_bytes gauge
	docker_container_memory_usage_bytes{name="testName"} 9999
	# HELP docker_container_memory_usage_ratio Memory usage as a ratio of the limit (0-1)
	# TYPE docker_container_memory_usage_ratio gauge
	docker_container_memory_usage_ratio{name="testName"} 1.249875e-06
	# HELP docker_container_network_receive_bytes_total Total network bytes received
	# TYPE docker_container_network_receive_bytes_total counter
	docker_container_network_receive_bytes_total{name="testName",network="eth0"} 135
	# HELP docker_container_network_receive_errors_total Total network receive errors
	# TYPE docker_container_network_receive_errors_total counter
	docker_container_network_receive_errors_total{name="testName",network="eth0"} 1
	# HELP docker_container_network_receive_packets_dropped_total Total network packets dropped while receiving
	# TYPE docker_container_network_receive_packets_dropped_total counter
	docker_container_network_receive_packets_dropped_total{name="testName",network="eth0"} 3
	# HELP docker_container_network_receive_packets_total Total network packets received
	# TYPE docker_container_network_receive_packets_total counter
	docker_container_network_receive_packets_total{name="testName",network="eth0"} 246
	# HELP docker_container_network_transmit_bytes_total Total network bytes transmitted
	# TYPE docker_container_network_transmit_bytes_total counter
	docker_container_network_transmit_bytes_total{name="testName",network="eth0"} 975
	# HELP docker_container_network_transmit_errors_total Total network transmit errors
	# TYPE docker_container_network_transmit_errors_total counter
	docker_container_network_transmit_errors_total{name="testName",network="eth0"} 4
	# HELP docker_container_network_transmit_packets_dropped_total Total network packets dropped while transmitting
	# TYPE docker_container_network_transmit_packets_dropped_total counter
	docker_container_network_transmit_packets_dropped_total{name="testName",network="eth0"} 2
	# HELP docker_container_network_transmit_packets_total Total network packets transmitted
	# TYPE docker_container_network_transmit_packets_total counter
	docker_container_network_transmit_packets_total{name="testName",network="eth0"} 864
	# HELP docker_container_pids_current Current number of pids
	# TYPE docker_container_pids_current gauge
	docker_container_pids_current{name="testName"} 12
	# HELP docker_container_state State of the container
	# TYPE docker_container_state gauge
	docker_container_state{name="testName",state="running"} 1
	# HELP docker_container_uptime_seconds Uptime of the container in seconds
	# TYPE docker_container_uptime_seconds gauge
	docker_container_uptime_seconds{name="testName"} 1.0
	# HELP docker_exporter_scrape_duration_seconds Duration of the scrape in seconds
	# TYPE docker_exporter_scrape_duration_seconds gauge
	docker_exporter_scrape_duration_seconds 2
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected),
		"docker_container_cpu_online_cpus",
		"docker_container_cpu_usage_seconds_total",
		"docker_container_fs_reads_bytes_total",
		"docker_container_fs_writes_bytes_total",
		"docker_container_info",
		"docker_container_memory_limit_bytes",
		"docker_container_memory_usage_bytes",
		"docker_container_memory_usage_ratio",
		"docker_container_network_receive_bytes_total",
		"docker_container_network_receive_errors_total",
		"docker_container_network_receive_packets_dropped_total",
		"docker_container_network_receive_packets_total",
		"docker_container_network_transmit_bytes_total",
		"docker_container_network_transmit_errors_total",
		"docker_container_network_transmit_packets_dropped_total",
		"docker_container_network_transmit_packets_total",
		"docker_container_pids_current",
		"docker_container_state",
		"docker_container_uptime_seconds",
		"docker_exporter_scrape_duration_seconds",
	); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectDeprecatedMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Now()).Times(1)
	mockClock.EXPECT().
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(1 * time.Second).Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(2 * time.Second).Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	// The deprecated metrics are still emitted for backward compatibility.
	const expected = `
	# HELP docker_container_block_io_read_bytes Block I/O read bytes total (deprecated; use docker_container_fs_reads_bytes_total)
	# TYPE docker_container_block_io_read_bytes gauge
	docker_container_block_io_read_bytes{name="testName"} 9999
	# HELP docker_container_cpu_usage_percentage CPU usage in percentage (deprecated; use docker_container_cpu_usage_seconds_total)
	# TYPE docker_container_cpu_usage_percentage gauge
	docker_container_cpu_usage_percentage{name="testName"} 0
	# HELP docker_container_memory_total_bytes Total memory in bytes (deprecated; use docker_container_memory_limit_bytes)
	# TYPE docker_container_memory_total_bytes gauge
	docker_container_memory_total_bytes{name="testName"} 8e+09
	# HELP docker_container_memory_usage_percentage Memory usage in percentage (deprecated; use docker_container_memory_usage_ratio)
	# TYPE docker_container_memory_usage_percentage gauge
	docker_container_memory_usage_percentage{name="testName"} 0.0001249875
	# HELP docker_container_network_rx_bytes Network received bytes total (deprecated; use docker_container_network_receive_bytes_total)
	# TYPE docker_container_network_rx_bytes gauge
	docker_container_network_rx_bytes{name="testName",network="eth0"} 135
	# HELP docker_container_uptime Uptime of the container in seconds (deprecated; use docker_container_uptime_seconds)
	# TYPE docker_container_uptime gauge
	docker_container_uptime{name="testName"} 1.0
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected),
		"docker_container_block_io_read_bytes",
		"docker_container_cpu_usage_percentage",
		"docker_container_memory_total_bytes",
		"docker_container_memory_usage_percentage",
		"docker_container_network_rx_bytes",
		"docker_container_uptime",
	); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectMetricsShouldCollectErrorWhenContainerListFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockErrorDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().
		Now().
		Return(time.Now()).
		Times(1)

	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(2 * time.Second).
		Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	const expected = `
	# HELP docker_exporter_scrape_errors_total Total number of scrape errors
	# TYPE docker_exporter_scrape_errors_total counter
	docker_exporter_scrape_errors_total 1
	# HELP docker_exporter_scrape_duration_seconds Duration of the scrape in seconds
	# TYPE docker_exporter_scrape_duration_seconds gauge
	docker_exporter_scrape_duration_seconds 2
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected),
		"docker_exporter_scrape_errors_total",
		"docker_exporter_scrape_duration_seconds",
	); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectMetricsShouldCollectErrorWhenContainerInspectFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockContainerInspectErrorDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().
		Now().
		Return(time.Now()).
		Times(1)

	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(2 * time.Second).
		Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	const expected = `
	# HELP docker_exporter_scrape_errors_total Total number of scrape errors
	# TYPE docker_exporter_scrape_errors_total counter
	docker_exporter_scrape_errors_total 1
	# HELP docker_exporter_scrape_duration_seconds Duration of the scrape in seconds
	# TYPE docker_exporter_scrape_duration_seconds gauge
	docker_exporter_scrape_duration_seconds 2
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected),
		"docker_exporter_scrape_errors_total",
		"docker_exporter_scrape_duration_seconds",
	); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectMetricsShouldCollectErrorWhenContainerStatsFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockContainerStatsErrorDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().
		Now().
		Return(time.Now()).
		Times(1)

	mockClock.EXPECT().
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)

	// The first call to Since is for the uptime of the container
	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(1 * time.Second).
		Times(1)

	// The second call to Since is for the scrape duration
	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(2 * time.Second).
		Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	const expected = `
	# HELP docker_container_info Infos about the container
	# TYPE docker_container_info gauge
	docker_container_info{image="sha256:d3751d33f9cd5049c4af2b462735457e4d3baf130bcbb87f389e349fbaeb20b9",image_name="myImage",name="testName"} 1
	# HELP docker_container_state State of the container
	# TYPE docker_container_state gauge
	docker_container_state{name="testName",state="running"} 1
	# HELP docker_container_uptime_seconds Uptime of the container in seconds
	# TYPE docker_container_uptime_seconds gauge
	docker_container_uptime_seconds{name="testName"} 1
	# HELP docker_exporter_scrape_duration_seconds Duration of the scrape in seconds
	# TYPE docker_exporter_scrape_duration_seconds gauge
	docker_exporter_scrape_duration_seconds 2
	# HELP docker_exporter_scrape_errors_total Total number of scrape errors
	# TYPE docker_exporter_scrape_errors_total counter
	docker_exporter_scrape_errors_total 1
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected),
		"docker_container_info",
		"docker_container_state",
		"docker_container_uptime_seconds",
		"docker_exporter_scrape_duration_seconds",
		"docker_exporter_scrape_errors_total",
	); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func buildInspectResponse() types.ContainerJSON {
	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				StartedAt: "2023-09-17T12:00:00.00Z",
			},
			Image: "sha256:d3751d33f9cd5049c4af2b462735457e4d3baf130bcbb87f389e349fbaeb20b9",
		},
		Config: &container.Config{
			Image: "myImage",
		},
	}

}

func buildContainerListResponse() []types.Container {
	return []types.Container{
		{
			ID:    "testID",
			Names: []string{"/testName"},
			State: "running",
			Labels: map[string]string{
				"com.docker.compose.project": "web",
				"maintainer":                 "acme",
			},
		},
		{
			ID:    "testIDIgnored",
			Names: []string{"/testNameIgnored"},
			State: "running",
			Labels: map[string]string{
				ignoreLabel: "true",
			},
		},
	}
}

func buildStatsResponse() container.StatsResponse {
	return container.StatsResponse{
		Stats: container.Stats{
			BlkioStats: container.BlkioStats{
				IoServiceBytesRecursive: []container.BlkioStatEntry{
					{
						Op:    "read",
						Value: 9999,
					},
					{
						Op:    "write",
						Value: 7777,
					},
				},
			},
			MemoryStats: container.MemoryStats{
				Usage:    9999,
				MaxUsage: 99999,
				Limit:    8000000000,
				Stats: map[string]uint64{
					"total_inactive_file": 121212,
				},
			},
			CPUStats: container.CPUStats{
				CPUUsage: container.CPUUsage{
					TotalUsage: 8888,
				},
				SystemUsage: 202,
				OnlineCPUs:  4,
			},
			PreCPUStats: container.CPUStats{
				CPUUsage: container.CPUUsage{
					TotalUsage: 1111,
				},
				SystemUsage: 1223,
				OnlineCPUs:  4,
			},
			PidsStats: container.PidsStats{
				Current: 12,
			},
		},
		Networks: map[string]container.NetworkStats{
			"eth0": {
				RxBytes:   135,
				RxPackets: 246,
				RxDropped: 3,
				RxErrors:  1,
				TxBytes:   975,
				TxPackets: 864,
				TxDropped: 2,
				TxErrors:  4,
			},
		},
	}
}

func mockJsonResponse(w http.ResponseWriter, r *http.Request, body any) {
	raw, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	_, _ = w.Write(raw)
}

func mockDockerApi(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "stats") {
		mockJsonResponse(w, r, buildStatsResponse())
		return
	}

	if strings.Contains(r.URL.Path, "testID") {
		mockJsonResponse(w, r, buildInspectResponse())
		return
	}

	mockJsonResponse(w, r, buildContainerListResponse())
}

func mockErrorDockerApi(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func mockContainerInspectErrorDockerApi(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "stats") {
		mockJsonResponse(w, r, buildStatsResponse())
		return
	}

	if strings.Contains(r.URL.Path, "testID") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mockJsonResponse(w, r, buildContainerListResponse())
}

func mockContainerStatsErrorDockerApi(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "stats") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.URL.Path, "testID") {
		mockJsonResponse(w, r, buildInspectResponse())
		return
	}

	mockJsonResponse(w, r, buildContainerListResponse())
}

func TestCollectContainerLabels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := httptest.NewServer(http.HandlerFunc(mockDockerApi))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Now()).Times(1)
	mockClock.EXPECT().
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(1 * time.Second).Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(2 * time.Second).Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, []string{
		"com.docker.compose.project",
		"unset",
	})
	// "unset" is not present on the container, so (kube_pod_labels style) it is
	// simply omitted rather than exposed with an empty value. "maintainer" is
	// present but not selected, so it is omitted too.
	const expected = `
	# HELP docker_container_labels Container labels converted to Prometheus labels
	# TYPE docker_container_labels gauge
	docker_container_labels{container_label_com_docker_compose_project="web",name="testName"} 1
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected), "docker_container_labels"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}

func TestCollectContainerLabelsFromContainerLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// A container opts specific labels in via the docker-exporter/exposed-labels
	// label; no global keys are configured.
	list := []types.Container{
		{
			ID:    "testID",
			Names: []string{"/testName"},
			State: "running",
			Labels: map[string]string{
				"com.docker.compose.project":     "web",
				"maintainer":                     "acme",
				"docker-exporter.exposed-labels": "maintainer",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "stats"):
			mockJsonResponse(w, r, buildStatsResponse())
		case strings.Contains(r.URL.Path, "testID"):
			mockJsonResponse(w, r, buildInspectResponse())
		default:
			mockJsonResponse(w, r, list)
		}
	}))
	defer srv.Close()

	cli, err := client.NewClientWithOpts(
		client.WithHost(srv.URL),
		client.WithHTTPClient(&http.Client{}),
	)

	if err != nil {
		panic(err)
	}

	mockClock := mock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Now()).Times(1)
	mockClock.EXPECT().
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(1 * time.Second).Times(1)
	mockClock.EXPECT().Since(gomock.Any()).Return(2 * time.Second).Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel, nil)

	// Only "maintainer" was opted in; "com.docker.compose.project" is present
	// but not selected.
	const expected = `
	# HELP docker_container_labels Container labels converted to Prometheus labels
	# TYPE docker_container_labels gauge
	docker_container_labels{container_label_maintainer="acme",name="testName"} 1
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected), "docker_container_labels"); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
