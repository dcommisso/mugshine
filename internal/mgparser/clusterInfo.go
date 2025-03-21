package mgparser

import (
	configv1 "github.com/openshift/api/config/v1"
	"os"
	"sigs.k8s.io/yaml"
)

func (m *Mg) GetApiServerURL() string {
	return m.infrastructure.Status.APIServerURL
}

func parseInfrastructureList(filename string) (configv1.InfrastructureList, error) {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return configv1.InfrastructureList{}, err
	}

	infralist := configv1.InfrastructureList{}
	err = yaml.Unmarshal(manifest, &infralist)
	if err != nil {
		return configv1.InfrastructureList{}, err
	}

	return infralist, nil
}
