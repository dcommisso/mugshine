package mgparser

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Pod struct {
	v1.Pod
	podDirectoryPath string
	Containers       map[string]*Container
	InitContainers   map[string]*Container
}

func newPod(pod *v1.Pod, namespaceDirectory string) *Pod {
	podDirectoryPath := strings.TrimSuffix(namespaceDirectory, "/") + "/pods/" + pod.GetName()

	containersToReturn := map[string]*Container{}
	initContainersToReturn := map[string]*Container{}

	for _, container := range pod.Status.ContainerStatuses {
		containersToReturn[container.Name] = newContainer(&container, podDirectoryPath)
	}

	for _, initContainer := range pod.Status.InitContainerStatuses {
		initContainersToReturn[initContainer.Name] = newContainer(&initContainer, podDirectoryPath)
	}

	return &Pod{
		Pod:              *pod,
		podDirectoryPath: podDirectoryPath,
		Containers:       containersToReturn,
		InitContainers:   initContainersToReturn,
	}
}

func (p *Pod) GetManifestFilePath() string {
	return strings.TrimSuffix(p.podDirectoryPath, "/") + "/" + p.GetName() + ".yaml"
}
