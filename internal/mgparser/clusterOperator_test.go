package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClusterOperatorsAlphabetical(t *testing.T) {
	cases := map[string]struct {
		mgPath              string
		clusterOperatorList []string
	}{
		"validMg": {
			mgPath: "./testdata/mgs/validMg",
			clusterOperatorList: []string{
				"authentication",
				"baremetal",
				"cloud-controller-manager",
				"cloud-credential",
				"cluster-autoscaler",
				"config-operator",
				"console",
				"control-plane-machine-set",
				"csi-snapshot-controller",
				"dns",
				"etcd",
				"image-registry",
				"ingress",
				"insights",
				"kube-apiserver",
				"kube-controller-manager",
				"kube-scheduler",
				"kube-storage-version-migrator",
				"machine-api",
				"machine-approver",
				"machine-config",
				"marketplace",
				"monitoring",
				"network",
				"node-tuning",
				"openshift-apiserver",
				"openshift-controller-manager",
				"openshift-samples",
				"operator-lifecycle-manager",
				"operator-lifecycle-manager-catalog",
				"operator-lifecycle-manager-packageserver",
				"service-ca",
				"storage",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg01, _ := NewMg(tc.mgPath)

			assert.Equal(t, tc.clusterOperatorList, mg01.GetClusterOperatorsAlphabetical())
		})
	}
}
