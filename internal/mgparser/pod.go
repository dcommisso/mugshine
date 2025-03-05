package mgparser

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Pod struct {
	v1.Pod
	podDirectoryPath string
}

func newPod(pod *v1.Pod, namespaceDirectory string) *Pod {
	return &Pod{
		Pod:              *pod,
		podDirectoryPath: strings.TrimSuffix(namespaceDirectory, "/") + "/pods/" + pod.GetName(),
	}
}

func (p *Pod) GetLogsFilePath(containerName string) string {
	return strings.TrimSuffix(p.podDirectoryPath, "/") + "/" + containerName + "/" + containerName + "/logs/current.log"
}

func (p *Pod) GetManifestFilePath() string {
	return strings.TrimSuffix(p.podDirectoryPath, "/") + "/" + p.GetName() + ".yaml"
}
