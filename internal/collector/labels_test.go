package collector

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildContainerLabels(t *testing.T) {
	labels := buildContainerLabels([]string{
		"com.docker.compose.project",
		"my-label",
		"a.b",
		"a-b", // sanitizes to the same name as "a.b" -> conflict suffix
	})

	want := []containerLabel{
		{dockerKey: "com.docker.compose.project", promName: "container_label_com_docker_compose_project"},
		{dockerKey: "my-label", promName: "container_label_my_label"},
		{dockerKey: "a.b", promName: "container_label_a_b"},
		{dockerKey: "a-b", promName: "container_label_a_b_conflict2"},
	}

	if len(labels) != len(want) {
		t.Fatalf("got %d labels, want %d: %+v", len(labels), len(want), labels)
	}

	for i, w := range want {
		if labels[i] != w {
			t.Errorf("label %d = %+v, want %+v", i, labels[i], w)
		}
	}
}

func TestSelectedLabelKeys(t *testing.T) {
	c := &DockerCollector{containerLabelKeys: []string{"app", "team"}}
	container := types.Container{
		Labels: map[string]string{
			// "team" duplicates a global key; whitespace and empty entries are
			// trimmed/dropped.
			exposedLabelsLabel: "team, project ,, role",
		},
	}

	got := c.selectedLabelKeys(container)
	want := []string{"app", "team", "project", "role"}

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("key %d = %q, want %q", i, got[i], w)
		}
	}
}
