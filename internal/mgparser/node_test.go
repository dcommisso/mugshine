package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNodeManifestFilePath(t *testing.T) {
	cases := map[string]struct {
		mgPath           string
		node             string
		expectedNodePath string
	}{
		"master-0.clustername.domain.local": {
			mgPath:           "./testdata/mgs/validMg",
			node:             "master-0.clustername.domain.local",
			expectedNodePath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/cluster-scoped-resources/core/nodes/master-0.clustername.domain.local.yaml",
		},
		"worker-2.clustername.domain.local": {
			mgPath:           "./testdata/mgs/validMg",
			node:             "worker-2.clustername.domain.local",
			expectedNodePath: "./testdata/mgs/validMg/quay-io-openshift-release-dev-ocp-v4-0-art-dev-sha256/cluster-scoped-resources/core/nodes/worker-2.clustername.domain.local.yaml",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)

			assert.Equal(t, tc.expectedNodePath, mg01.Nodes[tc.node].GetManifestFilePath())
		})
	}
}
