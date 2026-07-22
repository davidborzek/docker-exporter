package collector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/prometheus/client_golang/prometheus"
)

// exposedLabelsLabel is the container label whose comma-separated value selects
// additional Docker label keys to expose for that specific container, on top of
// any globally configured keys.
const exposedLabelsLabel = "docker-exporter.exposed-labels"

const containerLabelPrefix = "container_label_"

var invalidLabelChar = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// containerLabel maps a Docker label key to the sanitized Prometheus label name
// it is exposed as.
type containerLabel struct {
	dockerKey string
	promName  string
}

// sanitizeLabelName turns a Docker label key into a valid Prometheus label name
// by prefixing it and replacing every invalid character with an underscore,
// e.g. "com.docker.compose.project" -> "container_label_com_docker_compose_project".
func sanitizeLabelName(key string) string {
	return containerLabelPrefix + invalidLabelChar.ReplaceAllString(key, "_")
}

// buildContainerLabels resolves Docker label keys to unique Prometheus label
// names, preserving order. Keys whose sanitized names collide get a
// "_conflictN" suffix so the resulting names stay unique.
func buildContainerLabels(keys []string) []containerLabel {
	labels := make([]containerLabel, 0, len(keys))
	seen := make(map[string]int, len(keys))

	for _, key := range keys {
		name := sanitizeLabelName(key)
		if n := seen[name]; n > 0 {
			seen[name] = n + 1
			name = fmt.Sprintf("%s_conflict%d", name, n+1)
		} else {
			seen[name] = 1
		}

		labels = append(labels, containerLabel{dockerKey: key, promName: name})
	}

	return labels
}

// selectedLabelKeys returns the deduplicated Docker label keys to expose for a
// container: the globally configured keys plus any listed in the container's
// exposedLabelsLabel value.
func (c *DockerCollector) selectedLabelKeys(container types.Container) []string {
	seen := make(map[string]struct{})
	keys := make([]string, 0, len(c.containerLabelKeys))

	add := func(k string) {
		k = strings.TrimSpace(k)
		if k == "" {
			return
		}
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = struct{}{}
		keys = append(keys, k)
	}

	for _, k := range c.containerLabelKeys {
		add(k)
	}
	if raw, ok := container.Labels[exposedLabelsLabel]; ok {
		for _, k := range strings.Split(raw, ",") {
			add(k)
		}
	}

	return keys
}

// collectContainerLabels emits the docker_container_labels metric for a
// container, exposing only the selected labels the container actually sets
// (following the kube_pod_labels convention). It is a no-op when nothing is
// selected or present, so every series carries a consistent, meaningful set.
func (c *DockerCollector) collectContainerLabels(ch chan<- prometheus.Metric, name string, container types.Container) {
	present := make([]string, 0)
	for _, key := range c.selectedLabelKeys(container) {
		if _, ok := container.Labels[key]; ok {
			present = append(present, key)
		}
	}

	if len(present) == 0 {
		return
	}

	labels := buildContainerLabels(present)

	names := make([]string, 0, len(labels)+1)
	values := make([]string, 0, len(labels)+1)
	names = append(names, "name")
	values = append(values, name)
	for _, l := range labels {
		names = append(names, l.promName)
		values = append(values, container.Labels[l.dockerKey])
	}

	desc := prometheus.NewDesc(
		"docker_container_labels",
		"Container labels converted to Prometheus labels",
		names,
		nil,
	)

	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1, values...)
}
