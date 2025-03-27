package main

import (
	"strconv"

	"github.com/dcommisso/mugshine/internal/mgparser"
)

/*
 * This is the section with the pod manifests and logs
 *
 * This is the structure:
 *
 * |   aeLogs   |  aeNamespace |aePod | aeContainer  |
 * +++++++++++++++++++++++++++++++++++++++++++++++++++
 *  Namespaces --> ns01
 *                 ns02   ----> pod01
 *                 ns03         pod02
 *                              pod03 --> container01
 *                                        container02
 */

/* aeLogs */

const aeLogsSectionName = "Namespaces"

type aeLogs struct {
	mg *mgparser.Mg
}

func (a *aeLogs) Init(mg *mgparser.Mg) {
	a.mg = mg
}
func (a aeLogs) Header() string      { return "OCP RESOURCES" }
func (a aeLogs) Title() string       { return aeLogsSectionName }
func (a aeLogs) FilterValue() string { return a.Title() }
func (a aeLogs) IsFailed() bool {
	for _, ns := range a.Selected() {
		if ns.IsFailed() {
			return true
		}
	}
	return false
}

func (a aeLogs) Selected() []ActionableElement {
	nsToReturn := []ActionableElement{}
	for _, nsName := range a.mg.GetNamespacesAlphabetical() {
		nsToReturn = append(nsToReturn, aeNamespace{
			namespace: a.mg.Namespaces[nsName],
		})
	}
	return nsToReturn
}

func (a aeLogs) Pressed() (fileToOpen string) { return "" }

/* aeNamespace */
type aeNamespace struct {
	namespace *mgparser.Namespace
}

func (a aeNamespace) Init(mg *mgparser.Mg) {}
func (a aeNamespace) Header() string       { return "NAMESPACES" }
func (a aeNamespace) Title() string        { return a.namespace.Name }
func (a aeNamespace) FilterValue() string  { return a.Title() }
func (a aeNamespace) IsFailed() bool {
	for _, pod := range a.Selected() {
		if pod.IsFailed() {
			return true
		}
	}
	return false
}

func (a aeNamespace) Selected() []ActionableElement {
	podsToReturn := []ActionableElement{}

	resourceFields := map[string][]string{
		"#HEADER#": []string{"NAME", "READY", "STATUS", "RESTARTS"},
	}
	for _, podName := range a.namespace.GetPodsAlphabetical() {
		pod := a.namespace.Pods[podName]
		resourceFields[podName] = []string{
			pod.GetName(),
			pod.GetOcOutput().Ready,
			pod.GetOcOutput().Status,
			strconv.Itoa(pod.GetOcOutput().Restarts),
		}
	}
	formattedFields := formatFieldsInLines(resourceFields)

	for _, podName := range a.namespace.GetPodsAlphabetical() {
		podsToReturn = append(podsToReturn, aePod{
			pod:    a.namespace.Pods[podName],
			header: formattedFields["#HEADER#"],
			title:  formattedFields[podName],
		})
	}

	return podsToReturn
}

func (a aeNamespace) Pressed() (fileToOpen string) { return "" }

/* aePod */
type aePod struct {
	pod           *mgparser.Pod
	title, header string
}

func (a aePod) Init(mg *mgparser.Mg) {}
func (a aePod) Header() string       { return a.header }
func (a aePod) Title() string        { return a.title }
func (a aePod) FilterValue() string  { return a.pod.GetName() }
func (a aePod) IsFailed() bool {
	if a.pod.GetOcOutput().Status != "Running" && a.pod.GetOcOutput().Status != "Completed" {
		return true
	}
	return false
}

func (a aePod) Selected() []ActionableElement {
	containersToReturn := []ActionableElement{}

	// create resourceFields and add header to it
	resourceFields := map[string][]string{
		"#HEADER#": []string{"NAME", "STATUS", "RESTARTS", "TYPE"},
	}
	// populate resourceFields with containers
	for _, containerName := range a.pod.GetContainersAlphabetical() {
		container := a.pod.Containers[containerName]
		resourceFields[containerName] = []string{
			container.Name,
			container.GetOcOutput().Status,
			strconv.Itoa(container.GetOcOutput().Restarts),
			"", // type field is empty in regular container
		}
	}
	// populate resourceFields with init containers
	for _, containerName := range a.pod.GetInitContainersAlphabetical() {
		container := a.pod.InitContainers[containerName]
		resourceFields[containerName] = []string{
			container.Name,
			container.GetOcOutput().Status,
			strconv.Itoa(container.GetOcOutput().Restarts),
			"init",
		}
	}
	formattedFields := formatFieldsInLines(resourceFields)

	// create containers
	for _, containerName := range a.pod.GetContainersAlphabetical() {
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.Containers[containerName],
			header:    formattedFields["#HEADER#"],
			title:     formattedFields[containerName],
		})
	}
	// create init containers
	for _, containerName := range a.pod.GetInitContainersAlphabetical() {
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.InitContainers[containerName],
			header:    formattedFields["#HEADER#"],
			title:     formattedFields[containerName],
		})
	}

	return containersToReturn
}

func (a aePod) Pressed() (fileToOpen string) {
	return a.pod.GetManifestFilePath()
}

/* aeContainer */
type aeContainer struct {
	container     *mgparser.Container
	header, title string
}

func (a aeContainer) Init(mg *mgparser.Mg) {}
func (a aeContainer) Header() string       { return a.header }
func (a aeContainer) Title() string        { return a.title }
func (a aeContainer) FilterValue() string  { return a.container.Name }
func (a aeContainer) IsFailed() bool {
	if a.container.GetOcOutput().Status != "Running" && a.container.GetOcOutput().Status != "Completed" {
		return true
	}
	return false
}
func (a aeContainer) Selected() []ActionableElement { return nil }
func (a aeContainer) Pressed() (fileToOpen string) {
	return a.container.GetLogsFilename()
}
