package mgparser

import (
	"os"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type Node struct {
	v1.Node
	nodeFilePath string
}

func parseNode(filename string) (v1.Node, error) {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return v1.Node{}, err
	}

	node := v1.Node{}
	err = yaml.Unmarshal(manifest, &node)
	if err != nil {
		return v1.Node{}, err
	}

	return node, nil
}

func (n *Node) GetManifestFilePath() string {
	return n.nodeFilePath
}

func (n *Node) GetStatus() string {
	var status string
	conditions := n.Status.Conditions
	for _, cond := range conditions {
		if cond.Type == "Ready" {
			if cond.Status != "True" {
				status = "NotReady"
			} else {
				status = "Ready"
			}
		}
	}
	return status
}
