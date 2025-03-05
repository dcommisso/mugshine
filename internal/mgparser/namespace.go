package mgparser

import (
	"path/filepath"
)

type Namespace struct {
	Name string
	Pods map[string]*Pod
}

func newNamespace(namespaceDirectory string) *Namespace {
	nsPods := map[string]*Pod{}
	nsName := filepath.Base(namespaceDirectory)
	pods, _ := getPods(namespaceDirectory + "/core/pods.yaml")

	// Some namespace dirs don't contain pods file
	if pods != nil {
		for _, pod := range pods {
			nsPods[pod.GetName()] = newPod(&pod, namespaceDirectory)
		}
	}

	return &Namespace{
		Name: nsName,
		Pods: nsPods,
	}
}

func (n *Namespace) GetPodsAlphabetical() []string {
	return getAlphabeticalKeys(n.Pods)
}
