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

func TestNewPod(t *testing.T) {
	const (
		mgPath      = "./testdata/mgs/validMg"
		inspectPath = "./testdata/mgs/validInspect"
	)
	cases := map[string]struct {
		mgPath                              string
		namespaceName                       string
		podName                             string
		expectedContainerNumber             int
		expectedInitContainerNumber         int
		arbitraryContainerName              string
		arbitraryContainerDirectoryPath     string
		arbitraryInitContainerName          string
		arbitraryInitContainerDirectoryPath string
	}{
		"multus-additional-cni-plugins-2hz8k": {
			mgPath:                              mgPath,
			namespaceName:                       "openshift-multus",
			podName:                             "multus-additional-cni-plugins-2hz8k",
			expectedContainerNumber:             1,
			expectedInitContainerNumber:         6,
			arbitraryContainerName:              "kube-multus-additional-cni-plugins",
			arbitraryContainerDirectoryPath:     mgPath + "/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces/openshift-multus/pods/multus-additional-cni-plugins-2hz8k/kube-multus-additional-cni-plugins/kube-multus-additional-cni-plugins",
			arbitraryInitContainerName:          "egress-router-binary-copy",
			arbitraryInitContainerDirectoryPath: mgPath + "/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces/openshift-multus/pods/multus-additional-cni-plugins-2hz8k/egress-router-binary-copy/egress-router-binary-copy",
		},
		"scanner-db-d69986857-x4jrz": {
			mgPath:                              inspectPath,
			namespaceName:                       "stackrox",
			podName:                             "scanner-db-d69986857-x4jrz",
			expectedContainerNumber:             1,
			expectedInitContainerNumber:         1,
			arbitraryContainerName:              "db",
			arbitraryContainerDirectoryPath:     inspectPath + "/namespaces/stackrox/pods/scanner-db-d69986857-x4jrz/db/db",
			arbitraryInitContainerName:          "init-db",
			arbitraryInitContainerDirectoryPath: inspectPath + "/namespaces/stackrox/pods/scanner-db-d69986857-x4jrz/init-db/init-db",
		},
		"central-66b5ffdfdf-9qtlq": {
			mgPath:                          inspectPath,
			namespaceName:                   "stackrox",
			podName:                         "central-66b5ffdfdf-9qtlq",
			expectedContainerNumber:         1,
			expectedInitContainerNumber:     0,
			arbitraryContainerName:          "central",
			arbitraryContainerDirectoryPath: inspectPath + "/namespaces/stackrox/pods/central-66b5ffdfdf-9qtlq/central/central",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)
			pod01 := mg01.Namespaces[tc.namespaceName].Pods[tc.podName]

			containerElements := len(pod01.Containers)
			initContainerElements := len(pod01.InitContainers)

			assert.Equal(t, tc.expectedContainerNumber, containerElements)
			assert.Equal(t, tc.expectedInitContainerNumber, initContainerElements)

			if tc.expectedContainerNumber == 0 || tc.expectedInitContainerNumber == 0 {
				return
			}

			assert.Equal(t, tc.arbitraryContainerName, pod01.Containers[tc.arbitraryContainerName].Name)
			assert.Equal(t, tc.arbitraryContainerDirectoryPath, pod01.Containers[tc.arbitraryContainerName].ContainerDirectoryPath)

			if tc.arbitraryInitContainerName == "" {
				return
			}

			assert.Equal(t, tc.arbitraryInitContainerName, pod01.InitContainers[tc.arbitraryInitContainerName].Name)
			assert.Equal(t, tc.arbitraryInitContainerDirectoryPath, pod01.InitContainers[tc.arbitraryInitContainerName].ContainerDirectoryPath)
		})
	}
}
