package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogsFilePath(t *testing.T) {
	const inspectNamespacesDir = "./testdata/mgs/validInspect/namespaces"
	cases := map[string]struct {
		namespacePath    string
		podName          string
		containerName    string
		expectedLogsPath string
	}{
		"stackrox-central": {
			namespacePath:    inspectNamespacesDir + "/stackrox",
			podName:          "central-66b5ffdfdf-9qtlq",
			containerName:    "central",
			expectedLogsPath: "./testdata/mgs/validInspect/namespaces/stackrox/pods/central-66b5ffdfdf-9qtlq/central/central/logs/current.log",
		},
		"stackrox-scanner-db": {
			namespacePath:    inspectNamespacesDir + "/stackrox",
			podName:          "scanner-db-d69986857-x4jrz",
			containerName:    "init-db",
			expectedLogsPath: "./testdata/mgs/validInspect/namespaces/stackrox/pods/scanner-db-d69986857-x4jrz/init-db/init-db/logs/current.log",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ns := newNamespace(tc.namespacePath)
			for _, pod := range ns.Pods {
				if pod.GetName() == tc.podName {
					assert.Equal(t, tc.expectedLogsPath, pod.GetLogsFilePath(tc.containerName))
				}
			}
		})
	}
}

func TestGetManifestFilePath(t *testing.T) {
	const inspectNamespacesDir = "./testdata/mgs/validInspect/namespaces"
	cases := map[string]struct {
		namespacePath        string
		podName              string
		expectedManifestPath string
	}{
		"stackrox-central": {
			namespacePath:        inspectNamespacesDir + "/stackrox",
			podName:              "central-66b5ffdfdf-9qtlq",
			expectedManifestPath: inspectNamespacesDir + "/stackrox/pods/central-66b5ffdfdf-9qtlq/central-66b5ffdfdf-9qtlq.yaml",
		},
		"postgres": {
			namespacePath:        inspectNamespacesDir + "/stackrox",
			podName:              "postgres-5d95c6fb56-7zwz9",
			expectedManifestPath: inspectNamespacesDir + "/stackrox/pods/postgres-5d95c6fb56-7zwz9/postgres-5d95c6fb56-7zwz9.yaml",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ns := newNamespace(tc.namespacePath)
			for _, pod := range ns.Pods {
				if pod.GetName() == tc.podName {
					assert.Equal(t, tc.expectedManifestPath, pod.GetManifestFilePath())
				}
			}
		})
	}
}
