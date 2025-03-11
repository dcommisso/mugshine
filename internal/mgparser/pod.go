package mgparser

import (
	"fmt"
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

func (p *Pod) GetContainersAlphabetical() []string {
	return getAlphabeticalKeys(p.Containers)
}

func (p *Pod) GetInitContainersAlphabetical() []string {
	return getAlphabeticalKeys(p.InitContainers)
}

func (p *Pod) GetOcOutput() OcOutput {
	var containersReady, totalContainers, containerRestarts int
	var containerNotReadyReason string

	for _, containerStatus := range p.Status.ContainerStatuses {
		// check on "is not ready" must happens first, because a Completed container
		// has the Teminated state but it's also in Ready status
		if !containerStatus.Ready {
			if ok := containerStatus.State.Waiting; ok != nil {
				containerNotReadyReason = ok.Reason
			}
			if ok := containerStatus.State.Terminated; ok != nil {
				containerNotReadyReason = ok.Reason
			}
		} else if containerStatus.Ready {
			containersReady++
		}

		containerRestarts += int(containerStatus.RestartCount)
		totalContainers++
	}

	var actualStatus string
	if containerNotReadyReason == "" {
		actualStatus = string(p.Status.Phase)
	} else {
		actualStatus = containerNotReadyReason
	}

	return OcOutput{
		Ready:    fmt.Sprintf("%v/%v", containersReady, totalContainers),
		Restarts: containerRestarts,
		Status:   actualStatus,
	}
}
