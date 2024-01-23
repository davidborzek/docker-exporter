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
		Parse(gomock.Any(), gomock.Any()).
		DoAndReturn(func(s1, s2 string) (time.Time, error) {
			return time.Parse(s1, s2)
		}).
		Times(1)

	mockClock.EXPECT().
		Since(gomock.Any()).
		Return(1 * time.Second).
		Times(1)

	dc := collector.NewWithClient(cli, mockClock, ignoreLabel)

	const expected = `
	# HELP docker_container_block_io_read_bytes Block I/O read bytes total
	# TYPE docker_container_block_io_read_bytes gauge
	docker_container_block_io_read_bytes{name="testName"} 9999
	# HELP docker_container_block_io_write_bytes Block I/O write bytes total
	# TYPE docker_container_block_io_write_bytes gauge
	docker_container_block_io_write_bytes{name="testName"} 7777
	# HELP docker_container_cpu_online_cpus Number of online CPUs
	# TYPE docker_container_cpu_online_cpus gauge
	docker_container_cpu_online_cpus{name="testName"} 4
	# HELP docker_container_cpu_usage_percentage CPU usage in percentage
	# TYPE docker_container_cpu_usage_percentage gauge
	docker_container_cpu_usage_percentage{name="testName"} 0
	# HELP docker_container_info Infos about the container
	# TYPE docker_container_info gauge
	docker_container_info{image="sha256:d3751d33f9cd5049c4af2b462735457e4d3baf130bcbb87f389e349fbaeb20b9",image_name="myImage",name="testName"} 1
	# HELP docker_container_memory_total_bytes Total memory in bytes
	# TYPE docker_container_memory_total_bytes gauge
	docker_container_memory_total_bytes{name="testName"} 8e+09
	# HELP docker_container_memory_usage_bytes Memory usage in bytes
	# TYPE docker_container_memory_usage_bytes gauge
	docker_container_memory_usage_bytes{name="testName"} 9999
	# HELP docker_container_memory_usage_percentage Memory usage in percentage
	# TYPE docker_container_memory_usage_percentage gauge
	docker_container_memory_usage_percentage{name="testName"} 0.0001249875
	# HELP docker_container_network_rx_bytes Network received bytes total
	# TYPE docker_container_network_rx_bytes gauge
	docker_container_network_rx_bytes{name="testName",network="eth0"} 135
	# HELP docker_container_network_rx_dropped_packets Network dropped packets total
	# TYPE docker_container_network_rx_dropped_packets gauge
	docker_container_network_rx_dropped_packets{name="testName",network="eth0"} 3
	# HELP docker_container_network_rx_errors Network received errors
	# TYPE docker_container_network_rx_errors gauge
	docker_container_network_rx_errors{name="testName",network="eth0"} 1
	# HELP docker_container_network_rx_packets Network received packets total
	# TYPE docker_container_network_rx_packets gauge
	docker_container_network_rx_packets{name="testName",network="eth0"} 246
	# HELP docker_container_network_tx_bytes Network sent bytes total
	# TYPE docker_container_network_tx_bytes gauge
	docker_container_network_tx_bytes{name="testName",network="eth0"} 975
	# HELP docker_container_network_tx_dropped_packets Network dropped packets total
	# TYPE docker_container_network_tx_dropped_packets gauge
	docker_container_network_tx_dropped_packets{name="testName",network="eth0"} 2
	# HELP docker_container_network_tx_errors Network sent errors
	# TYPE docker_container_network_tx_errors gauge
	docker_container_network_tx_errors{name="testName",network="eth0"} 4
	# HELP docker_container_network_tx_packets Network sent packets total
	# TYPE docker_container_network_tx_packets gauge
	docker_container_network_tx_packets{name="testName",network="eth0"} 864
	# HELP docker_container_pids_current Current number of pids
	# TYPE docker_container_pids_current gauge
	docker_container_pids_current{name="testName"} 12
	# HELP docker_container_state State of the container
	# TYPE docker_container_state gauge
	docker_container_state{name="testName",state="running"} 1
	# HELP docker_container_uptime Uptime of the container in seconds
	# TYPE docker_container_uptime gauge
	docker_container_uptime{name="testName"} 1.0
	`

	if err := testutil.CollectAndCompare(dc, strings.NewReader(expected)); err != nil {
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

func buildStatsResponse() types.StatsJSON {
	return types.StatsJSON{
		Stats: types.Stats{
			BlkioStats: types.BlkioStats{
				IoServiceBytesRecursive: []types.BlkioStatEntry{
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
			MemoryStats: types.MemoryStats{
				Usage:    9999,
				MaxUsage: 99999,
				Limit:    8000000000,
				Stats: map[string]uint64{
					"total_inactive_file": 121212,
				},
			},
			CPUStats: types.CPUStats{
				CPUUsage: types.CPUUsage{
					TotalUsage: 8888,
				},
				SystemUsage: 202,
				OnlineCPUs:  4,
			},
			PreCPUStats: types.CPUStats{
				CPUUsage: types.CPUUsage{
					TotalUsage: 1111,
				},
				SystemUsage: 1223,
				OnlineCPUs:  4,
			},
			PidsStats: types.PidsStats{
				Current: 12,
			},
		},
		Networks: map[string]types.NetworkStats{
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

	w.Write(raw)
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
