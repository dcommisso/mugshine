package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogsFilename(t *testing.T) {
	const (
		mgPath      = "./testdata/mgs/validMg"
		inspectPath = "./testdata/mgs/validInspect"
	)

	cases := map[string]struct {
		mgPath            string
		namespaceName     string
		podName           string
		containerName     string
		initContainerName string
		expectedLogsPath  string
	}{
		"central": {
			mgPath:           inspectPath,
			namespaceName:    "stackrox",
			podName:          "central-66b5ffdfdf-9qtlq",
			containerName:    "central",
			expectedLogsPath: inspectPath + "/namespaces/stackrox/pods/central-66b5ffdfdf-9qtlq/central/central/logs/current.log",
		},
		"init-db": {
			mgPath:            inspectPath,
			namespaceName:     "stackrox",
			podName:           "scanner-db-d69986857-x4jrz",
			initContainerName: "init-db",
			expectedLogsPath:  inspectPath + "/namespaces/stackrox/pods/scanner-db-d69986857-x4jrz/init-db/init-db/logs/current.log",
		},
		"kube-multus-additional-cni-plugins": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-additional-cni-plugins-2hz8k",
			containerName:    "kube-multus-additional-cni-plugins",
			expectedLogsPath: mgPath + "/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces/openshift-multus/pods/multus-additional-cni-plugins-2hz8k/kube-multus-additional-cni-plugins/kube-multus-additional-cni-plugins/logs/current.log",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)

			if tc.containerName != "" {
				assert.Equal(t, tc.expectedLogsPath, mg01.Namespaces[tc.namespaceName].Pods[tc.podName].Containers[tc.containerName].GetLogsFilename())
			} else {
				assert.Equal(t, tc.expectedLogsPath, mg01.Namespaces[tc.namespaceName].Pods[tc.podName].InitContainers[tc.initContainerName].GetLogsFilename())
			}
		})
	}
}

func TestContainerGetOcOutput(t *testing.T) {
	const (
		mgPath      = "./testdata/mgs/validMg"
		inspectPath = "./testdata/mgs/validInspect"
	)
	cases := map[string]struct {
		mgPath           string
		namespaceName    string
		podName          string
		containerName    string
		expectedRestarts int
		expectedStatus   string
	}{
		"multus-2fxcb": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-btg47",
			containerName:    "kube-multus",
			expectedRestarts: 1,
			expectedStatus:   "Running",
		},
		"multus-admission-controller-75968f7c47-554fb": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-admission-controller-75968f7c47-554fb",
			containerName:    "kube-rbac-proxy",
			expectedRestarts: 7,
			expectedStatus:   "Running",
		},
		"multus-admission-controller-75968f7c47-": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-admission-controller-75968f7c47-554fb",
			containerName:    "multus-admission-controller",
			expectedRestarts: 5,
			expectedStatus:   "Running",
		},
		"network-metrics-daemon-l5lr4": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "network-metrics-daemon-l5lr4",
			containerName:    "kube-rbac-proxy",
			expectedRestarts: 0,
			expectedStatus:   "ContainerStatusUnknown",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)
			container01 := mg01.Namespaces[tc.namespaceName].Pods[tc.podName].Containers[tc.containerName]

			assert.Equal(t, tc.expectedRestarts, container01.GetOcOutput().Restarts)
			assert.Equal(t, tc.expectedStatus, container01.GetOcOutput().Status)
		})
	}
}
