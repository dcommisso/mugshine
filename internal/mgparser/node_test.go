package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNode(t *testing.T) {
	cases := map[string]struct {
		mgPath           string
		nodeName         string
		expectedNodePath string
		expectedStatus   string
		expectedRoles    string
	}{
		"master-0.clustername.domain.local": {
			mgPath:           "./testdata/mgs/validMg",
			nodeName:         "master-0.clustername.domain.local",
			expectedNodePath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/cluster-scoped-resources/core/nodes/master-0.clustername.domain.local.yaml",
			expectedStatus:   "Ready",
			expectedRoles:    "control-plane,master",
		},
		"worker-2.clustername.domain.local": {
			mgPath:           "./testdata/mgs/validMg",
			nodeName:         "worker-2.clustername.domain.local",
			expectedNodePath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/cluster-scoped-resources/core/nodes/worker-2.clustername.domain.local.yaml",
			expectedStatus:   "NotReady",
			expectedRoles:    "worker",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)

			node := mg01.Nodes[tc.nodeName]

			assert.Equal(t, tc.expectedNodePath, node.GetManifestFilePath())
			assert.Equal(t, tc.expectedStatus, node.GetStatus())
			assert.Equal(t, tc.expectedRoles, node.GetRoles())
		})
	}
}
