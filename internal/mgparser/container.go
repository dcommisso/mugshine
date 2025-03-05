package mgparser

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

type Container struct {
	v1.ContainerStatus
	ContainerDirectoryPath string
}

func newContainer(container *v1.ContainerStatus, podDirectory string) *Container {
	return &Container{
		ContainerStatus:        *container,
		ContainerDirectoryPath: strings.TrimSuffix(podDirectory, "/") + "/" + container.Name + "/" + container.Name,
	}
}

func (c *Container) GetLogsFilename() string {
	return strings.TrimSuffix(c.ContainerDirectoryPath, "/") + "/logs/current.log"
}
