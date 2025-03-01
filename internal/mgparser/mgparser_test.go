package mgparser

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMg(t *testing.T) {
	cases := map[string]struct {
		directory      string
		namespacesPath string
		namespaceNames []string
		podNames       map[string][]string
		expectedError  string
	}{
		"valid mg": {
			directory:      "./testdata/mgs/validMg",
			namespacesPath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces",
			namespaceNames: []string{
				"openshift-multus",
				"openshift-nfs-storage",
				"openshift-operator-lifecycle-manager",
			},
			podNames: map[string][]string{
				"openshift-multus": []string{
					"multus-2fxcb",
					"multus-4qflw",
					"multus-8c65f",
					"multus-additional-cni-plugins-2hz8k",
					"multus-additional-cni-plugins-7nlzg",
					"multus-additional-cni-plugins-fdq79",
					"multus-additional-cni-plugins-hpx5v",
					"multus-additional-cni-plugins-nr42g",
					"multus-additional-cni-plugins-t94zk",
					"multus-admission-controller-75968f7c47-554fb",
					"multus-admission-controller-75968f7c47-fzrvb",
					"multus-btg47",
					"multus-j9fx2",
					"multus-l6rj4",
					"network-metrics-daemon-489hj",
					"network-metrics-daemon-7n4k7",
					"network-metrics-daemon-k9jqg",
					"network-metrics-daemon-l5lr4",
					"network-metrics-daemon-mrgrr",
					"network-metrics-daemon-prss9",
				},
				"openshift-nfs-storage": []string{},
				"openshift-operator-lifecycle-manager": []string{
					"catalog-operator-869689d5d5-fkn5r",
					"collect-profiles-29006445-kvmhz",
					"collect-profiles-29006460-jnr5j",
					"collect-profiles-29006475-hf5hr",
					"olm-operator-676d4fb97-7v8pb",
					"package-server-manager-6984bc476c-w6jf7",
					"packageserver-dffff7b4-8lw4h",
					"packageserver-dffff7b4-pblk7",
				},
			},
		},
		"valid inspect": {
			directory:      "./testdata/mgs/validInspect",
			namespacesPath: "./testdata/mgs/validInspect/namespaces",
			namespaceNames: []string{
				"stackrox",
			},
			podNames: map[string][]string{
				"stackrox": {
					"central-66b5ffdfdf-9qtlq",
					"config-controller-6b79f4c4bb-zjtqm",
					"postgres-5d95c6fb56-7zwz9",
					"scanner-74dbd48b9d-pw8dx",
					"scanner-db-d69986857-x4jrz",
				},
			},
		},
		"valid mg with trailing slash": {
			directory:      "./testdata/mgs/validMg/",
			namespacesPath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces",
			namespaceNames: []string{
				"openshift-multus",
				"openshift-nfs-storage",
				"openshift-operator-lifecycle-manager",
			},
			podNames: map[string][]string{
				"openshift-multus": []string{
					"multus-2fxcb",
					"multus-4qflw",
					"multus-8c65f",
					"multus-additional-cni-plugins-2hz8k",
					"multus-additional-cni-plugins-7nlzg",
					"multus-additional-cni-plugins-fdq79",
					"multus-additional-cni-plugins-hpx5v",
					"multus-additional-cni-plugins-nr42g",
					"multus-additional-cni-plugins-t94zk",
					"multus-admission-controller-75968f7c47-554fb",
					"multus-admission-controller-75968f7c47-fzrvb",
					"multus-btg47",
					"multus-j9fx2",
					"multus-l6rj4",
					"network-metrics-daemon-489hj",
					"network-metrics-daemon-7n4k7",
					"network-metrics-daemon-k9jqg",
					"network-metrics-daemon-l5lr4",
					"network-metrics-daemon-mrgrr",
					"network-metrics-daemon-prss9",
				},
				"openshift-nfs-storage": []string{},
				"openshift-operator-lifecycle-manager": []string{
					"catalog-operator-869689d5d5-fkn5r",
					"collect-profiles-29006445-kvmhz",
					"collect-profiles-29006460-jnr5j",
					"collect-profiles-29006475-hf5hr",
					"olm-operator-676d4fb97-7v8pb",
					"package-server-manager-6984bc476c-w6jf7",
					"packageserver-dffff7b4-8lw4h",
					"packageserver-dffff7b4-pblk7",
				},
			},
		},
		"valid inspect with trailing slash": {
			directory:      "./testdata/mgs/validInspect/",
			namespacesPath: "./testdata/mgs/validInspect/namespaces",
			namespaceNames: []string{
				"stackrox",
			},
			podNames: map[string][]string{
				"stackrox": {
					"central-66b5ffdfdf-9qtlq",
					"config-controller-6b79f4c4bb-zjtqm",
					"postgres-5d95c6fb56-7zwz9",
					"scanner-74dbd48b9d-pw8dx",
					"scanner-db-d69986857-x4jrz",
				},
			},
		},
		"invalid mg: wrong number of dirs (0), timestamp ok": {
			directory:     "./testdata/mgs/mgDirBadTimestampOK",
			expectedError: "bad must-gather format",
		},
		"invalid mg: namespaces dir absent, timestamp ok": {
			directory:     "./testdata/mgs/mgNamespacesBadTimestampOK",
			expectedError: "bad must-gather format",
		},
		"invalid mg: wrong number of dirs (2), timestamp absent": {
			directory:     "./testdata/mgs/mgDirBadTimestampBad",
			expectedError: "bad must-gather format",
		},
		"invalid mg: dirs OK, timestamp absent": {
			directory:     "./testdata/mgs/mgDirOKTimestampBad",
			expectedError: "bad must-gather format",
		},
		"invalid inspect: wrong number of dirs (2), timestamp ok": {
			directory:     "./testdata/mgs/inspectDirBadTimestampOK",
			expectedError: "bad must-gather format",
		},
		"invalid inspect: namespaces dir absent, timestamp ok": {
			directory:     "./testdata/mgs/inspectNamespacesBadTimestampOK",
			expectedError: "bad must-gather format",
		},
		"invalid inspect: wrong number of dirs (0), timestamp absent": {
			directory:     "./testdata/mgs/inspectDirBadTimestampBad",
			expectedError: "bad must-gather format",
		},
		"invalid inspect: dirs OK, timestamp absent": {
			directory:     "./testdata/mgs/inspectDirOKTimestampBad",
			expectedError: "bad must-gather format",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg, err := NewMg(tc.directory)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.namespacesPath, mg.NamespacesPath)

			// get the namespace names and sort them for comparison
			namespaceNames := make([]string, len(mg.Namespaces))
			i := 0
			for k := range mg.Namespaces {
				namespaceNames[i] = k
				i++
			}
			slices.Sort(namespaceNames)
			assert.Equal(t, tc.namespaceNames, namespaceNames)

			// check if the Pods slice of each namespace is properly populated
			podNames := map[string][]string{}
			for _, ns := range namespaceNames {
				podNames[ns] = []string{}
				for _, pod := range mg.Namespaces[ns].Pods {
					podNames[ns] = append(podNames[ns], pod.GetName())
				}
				// sort each pod list for comparison
				slices.Sort(podNames[ns])
			}
			assert.Equal(t, tc.podNames, podNames)
		})
	}
}

func TestGetPods(t *testing.T) {
	cases := map[string]struct {
		manifestFile string
		firstPodName string
	}{
		"cluster-cloud-controller-manager-operator": {
			manifestFile: "./testdata/cluster-cloud-controller-manager-operator.yaml",
			firstPodName: "cluster-cloud-controller-manager-operator-64bbd8597f-g2gnt",
		},
		"openshift-console PodList": {
			manifestFile: "./testdata/openshift-console-pods.yaml",
			firstPodName: "console-f7996d89c-bczrn",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			pods, _ := getPods(tc.manifestFile)
			assert.Equal(t, tc.firstPodName, pods[0].GetName())
		})
	}
}

func TestNewNamespace(t *testing.T) {
	const mgNamespacesDir = "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces"
	const inspectNamespacesDir = "./testdata/mgs/validInspect/namespaces"
	cases := map[string]struct {
		namespaceDirectory string
		namespaceName      string
		firstPodName       string
		lastPodName        string
		podNumber          int
	}{
		"mg-openshift-multus": {
			namespaceDirectory: mgNamespacesDir + "/openshift-multus",
			namespaceName:      "openshift-multus",
			firstPodName:       "multus-2fxcb",
			lastPodName:        "network-metrics-daemon-prss9",
			podNumber:          20,
		},
		"mg-openshift-nfs-storage": {
			namespaceDirectory: mgNamespacesDir + "/openshift-nfs-storage",
			namespaceName:      "openshift-nfs-storage",
			firstPodName:       "",
			lastPodName:        "",
			podNumber:          0,
		},
		"mg-openshift-operator-lifecycle-manager": {
			namespaceDirectory: mgNamespacesDir + "/openshift-operator-lifecycle-manager",
			namespaceName:      "openshift-operator-lifecycle-manager",
			firstPodName:       "catalog-operator-869689d5d5-fkn5r",
			lastPodName:        "packageserver-dffff7b4-pblk7",
			podNumber:          8,
		},
		"inspect-stackrox": {
			namespaceDirectory: inspectNamespacesDir + "/stackrox",
			namespaceName:      "stackrox",
			firstPodName:       "central-66b5ffdfdf-9qtlq",
			lastPodName:        "scanner-db-d69986857-x4jrz",
			podNumber:          5,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ns := newNamespace(tc.namespaceDirectory)
			podsLen := len(ns.Pods)
			assert.Equal(t, tc.namespaceName, ns.Name)
			assert.Equal(t, tc.podNumber, podsLen)

			// This is necessary to avoid failure on checking an empty slice
			if tc.podNumber == 0 {
				return
			}

			assert.Equal(t, tc.firstPodName, ns.Pods[tc.firstPodName].GetName())
			assert.Equal(t, tc.lastPodName, ns.Pods[tc.lastPodName].GetName())
		})
	}
}

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

func TestGetNamespacesAlphabetical(t *testing.T) {
	cases := map[string]struct {
		mgPath                string
		namespaceAlphabetical []string
	}{
		"validMg": {
			mgPath: "./testdata/mgs/validMg",
			namespaceAlphabetical: []string{
				"openshift-multus",
				"openshift-nfs-storage",
				"openshift-operator-lifecycle-manager",
			},
		},
		"validInspect": {
			mgPath: "./testdata/mgs/validInspect",
			namespaceAlphabetical: []string{
				"stackrox",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)

			namespaces := mg01.GetNamespacesAlphabetical()
			assert.Equal(t, tc.namespaceAlphabetical, namespaces)
		})
	}
}

func TestGetPodsAlphabetical(t *testing.T) {
	const mgNamespacesDir = "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/namespaces"
	const inspectNamespacesDir = "./testdata/mgs/validInspect/namespaces"

	cases := map[string]struct {
		namespacesDirectory string
		namespaceName       string
		podsAlphabetical    []string
	}{
		"validMg": {
			namespacesDirectory: mgNamespacesDir,
			namespaceName:       "openshift-multus",
			podsAlphabetical: []string{
				"multus-2fxcb",
				"multus-4qflw",
				"multus-8c65f",
				"multus-additional-cni-plugins-2hz8k",
				"multus-additional-cni-plugins-7nlzg",
				"multus-additional-cni-plugins-fdq79",
				"multus-additional-cni-plugins-hpx5v",
				"multus-additional-cni-plugins-nr42g",
				"multus-additional-cni-plugins-t94zk",
				"multus-admission-controller-75968f7c47-554fb",
				"multus-admission-controller-75968f7c47-fzrvb",
				"multus-btg47",
				"multus-j9fx2",
				"multus-l6rj4",
				"network-metrics-daemon-489hj",
				"network-metrics-daemon-7n4k7",
				"network-metrics-daemon-k9jqg",
				"network-metrics-daemon-l5lr4",
				"network-metrics-daemon-mrgrr",
				"network-metrics-daemon-prss9",
			},
		},
		"validInspect": {
			namespacesDirectory: inspectNamespacesDir,
			namespaceName:       "stackrox",
			podsAlphabetical: []string{
				"central-66b5ffdfdf-9qtlq",
				"config-controller-6b79f4c4bb-zjtqm",
				"postgres-5d95c6fb56-7zwz9",
				"scanner-74dbd48b9d-pw8dx",
				"scanner-db-d69986857-x4jrz",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ns01 := newNamespace(tc.namespacesDirectory + "/" + tc.namespaceName)
			ordered := ns01.GetPodsAlphabetical()
			assert.Equal(t, tc.podsAlphabetical, ordered)
		})
	}
}
