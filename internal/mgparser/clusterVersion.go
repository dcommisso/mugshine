package mgparser

import (
	configv1 "github.com/openshift/api/config/v1"
	"os"
	"sigs.k8s.io/yaml"
)

func parseClusterVersion(filename string) (configv1.ClusterVersion, error) {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return configv1.ClusterVersion{}, err
	}

	cv := configv1.ClusterVersion{}
	err = yaml.Unmarshal(manifest, &cv)
	if err != nil {
		return configv1.ClusterVersion{}, err
	}

	return cv, nil
}

func (m *Mg) GetClusterID() string {
	if m.clusterVersion == nil {
		return ""
	}
	return string(m.clusterVersion.Spec.ClusterID)
}
