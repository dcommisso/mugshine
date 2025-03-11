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

func (c *Container) GetOcOutput() OcOutput {
	var status string
	// check on "is not ready" must happens first, because a Completed container
	// has the Teminated state but it's also in Ready status
	if !c.Ready {
		if ok := c.State.Waiting; ok != nil {
			status = ok.Reason
		}
		if ok := c.State.Terminated; ok != nil {
			status = ok.Reason
		}
	} else if c.Ready {
		status = "Running"
	}

	return OcOutput{
		Status:   status,
		Restarts: int(c.RestartCount),
	}
}
