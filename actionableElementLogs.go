package main

import (
	"slices"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/dcommisso/img/internal/mgparser"
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
func (a aeLogs) FilterValue() string { return "" }
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

func (a aeLogs) Pressed() (fileToOpen string) {
	return ""
}

/* aeNamespace */
type aeNamespace struct {
	namespace *mgparser.Namespace
}

func (a aeNamespace) Init(mg *mgparser.Mg) {}
func (a aeNamespace) Header() string       { return "NAMESPACES" }
func (a aeNamespace) Title() string        { return a.namespace.Name }
func (a aeNamespace) FilterValue() string  { return a.namespace.Name }
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

func (a aeNamespace) Pressed() (fileToOpen string) {
	return ""
}

/* aePod */
type aePod struct {
	pod           *mgparser.Pod
	title, header string
}

func (a aePod) Init(mg *mgparser.Mg) {}
func (a aePod) Header() string {
	return a.header
}

func (a aePod) Title() string {
	return a.title
}
func (a aePod) FilterValue() string { return a.pod.GetName() }
func (a aePod) IsFailed() bool {
	if a.pod.GetOcOutput().Status != "Running" && a.pod.GetOcOutput().Status != "Completed" {
		return true
	}
	return false
}

func (a aePod) Selected() []ActionableElement {
	containersToReturn := []ActionableElement{}

	// Calculate the max width for name and status fields for containers and
	// initContainers
	var nameLengths, statusLengths []int
	for _, containerName := range a.pod.GetContainersAlphabetical() {
		nameLengths = append(nameLengths, len(containerName))
		statusLengths = append(statusLengths, len(a.pod.Containers[containerName].GetOcOutput().Status))
	}
	for _, containerName := range a.pod.GetInitContainersAlphabetical() {
		nameLengths = append(nameLengths, len(containerName))
		statusLengths = append(statusLengths, len(a.pod.InitContainers[containerName].GetOcOutput().Status))
	}

	// Add containers
	for _, containerName := range a.pod.GetContainersAlphabetical() {
		containerLengths := map[string]int{
			"name":     slices.Max(nameLengths),
			"status":   slices.Max(statusLengths),
			"restarts": len("RESTARTS"),
			"type":     len("TYPE"),
		}
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.Containers[containerName],
			isInit:    false,
			lengths:   containerLengths,
		})
	}

	// Add initContainers
	for _, containerName := range a.pod.GetInitContainersAlphabetical() {
		containerLengths := map[string]int{
			"name":     slices.Max(nameLengths),
			"status":   slices.Max(statusLengths),
			"restarts": len("RESTARTS"),
			"type":     len("TYPE"),
		}
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.InitContainers[containerName],
			isInit:    true,
			lengths:   containerLengths,
		})
	}

	return containersToReturn
}

func (a aePod) Pressed() (fileToOpen string) {
	return a.pod.GetManifestFilePath()
}

/* aeContainer */
type aeContainer struct {
	container *mgparser.Container
	isInit    bool
	lengths   map[string]int
}

func (a aeContainer) Init(mg *mgparser.Mg) {}
func (a aeContainer) Header() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.lengths["name"]).Render("NAME"),
		baseStyle.Width(a.lengths["status"]).Render("STATUS"),
		baseStyle.Width(a.lengths["restarts"]).Render("RESTARTS"),
		baseStyle.Width(a.lengths["type"]).Render("TYPE"),
	)
}
func (a aeContainer) Title() string {
	initHeader := ""
	if a.isInit {
		initHeader = "Init"
	}

	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.lengths["name"]).Render(a.container.Name),
		baseStyle.Width(a.lengths["status"]).Render(a.container.GetOcOutput().Status),
		baseStyle.Width(a.lengths["restarts"]).Render(strconv.Itoa(a.container.GetOcOutput().Restarts)),
		baseStyle.Width(a.lengths["init"]).Render(initHeader),
	)
}

func (a aeContainer) FilterValue() string { return a.container.Name }
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
