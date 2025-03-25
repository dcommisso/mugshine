package mgparser

import (
	configv1 "github.com/openshift/api/config/v1"
	"os"
	"sigs.k8s.io/yaml"
)

type ClusterOperator struct {
	configv1.ClusterOperator
	coFilePath string
}

func parseClusterOperatorList(filename string) (configv1.ClusterOperatorList, error) {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return configv1.ClusterOperatorList{}, err
	}

	operatorList := configv1.ClusterOperatorList{}
	err = yaml.Unmarshal(manifest, &operatorList)
	if err != nil {
		return configv1.ClusterOperatorList{}, err
	}

	return operatorList, nil
}

func (m *Mg) GetClusterOperatorsAlphabetical() []string {
	return getAlphabeticalKeys(m.clusterOperators)
}
