package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// TestGetContainersAlphabetical actually tests both
// TestGetContainersAlphabetical and TestGetInitContainersAlphabetical
func TestGetContainersAlphabetical(t *testing.T) {
	const (
		mgPath      = "./testdata/mgs/validMg"
		inspectPath = "./testdata/mgs/validInspect"
	)
	cases := map[string]struct {
		mgPath                             string
		namespaceName                      string
		podName                            string
		expectedContainersAlphabetical     []string
		expectedInitContainersAlphabetical []string
	}{
		"kube-multus-additional-cni-plugins": {
			mgPath:        mgPath,
			namespaceName: "openshift-multus",
			podName:       "multus-additional-cni-plugins-hpx5v",
			expectedContainersAlphabetical: []string{
				"kube-multus-additional-cni-plugins",
			},
			expectedInitContainersAlphabetical: []string{
				"bond-cni-plugin",
				"cni-plugins",
				"egress-router-binary-copy",
				"routeoverride-cni",
				"whereabouts-cni",
				"whereabouts-cni-bincopy",
			},
		},
		"multus-admission-controller-75968f7c47-554fb": {
			mgPath:        mgPath,
			namespaceName: "openshift-multus",
			podName:       "multus-admission-controller-75968f7c47-554fb",
			expectedContainersAlphabetical: []string{
				"kube-rbac-proxy",
				"multus-admission-controller",
			},
			expectedInitContainersAlphabetical: []string{},
		},
		"scanner-db-d69986857-x4jrz": {
			mgPath:        inspectPath,
			namespaceName: "stackrox",
			podName:       "scanner-db-d69986857-x4jrz",
			expectedContainersAlphabetical: []string{
				"db",
			},
			expectedInitContainersAlphabetical: []string{
				"init-db",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)
			pod01 := mg01.Namespaces[tc.namespaceName].Pods[tc.podName]

			assert.Equal(t, tc.expectedContainersAlphabetical, pod01.GetContainersAlphabetical())
			assert.Equal(t, tc.expectedInitContainersAlphabetical, pod01.GetInitContainersAlphabetical())
		})
	}
}

// TODO: test more statuses to test cases
func TestGetOcOutput(t *testing.T) {
	const (
		mgPath      = "./testdata/mgs/validMg"
		inspectPath = "./testdata/mgs/validInspect"
	)
	cases := map[string]struct {
		mgPath           string
		namespaceName    string
		podName          string
		expectedReady    string
		expectedRestarts int
		expectedStatus   string
	}{
		"multus-2fxcb": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-btg47",
			expectedReady:    "1/1",
			expectedRestarts: 1,
			expectedStatus:   "Running",
		},
		"multus-admission-controller-75968f7c47-554fb": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "multus-admission-controller-75968f7c47-554fb",
			expectedReady:    "2/2",
			expectedRestarts: 0,
			expectedStatus:   "Running",
		},
		"network-metrics-daemon-l5lr4": {
			mgPath:           mgPath,
			namespaceName:    "openshift-multus",
			podName:          "network-metrics-daemon-l5lr4",
			expectedReady:    "0/2",
			expectedRestarts: 0,
			expectedStatus:   "ContainerStatusUnknown",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)
			pod01 := mg01.Namespaces[tc.namespaceName].Pods[tc.podName]

			assert.Equal(t, tc.expectedReady, pod01.GetOcOutput().Ready)
			assert.Equal(t, tc.expectedRestarts, pod01.GetOcOutput().Restarts)
			assert.Equal(t, tc.expectedStatus, pod01.GetOcOutput().Status)
		})
	}
}
