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
func (a aeLogs) GetWidthFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return 20
	}
}
func (a aeLogs) GetHeightFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return windowSize
	}
}
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
func (a aeNamespace) GetWidthFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return 50
	}
}
func (a aeNamespace) GetHeightFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return windowSize
	}
}
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

	// calculate the max width for name and status pod fields
	var nameLengths, statusLengths []int
	for _, podName := range a.namespace.GetPodsAlphabetical() {
		nameLengths = append(nameLengths, len(podName))
		statusLengths = append(statusLengths, len(a.namespace.Pods[podName].GetOcOutput().Status))
	}

	for _, podName := range a.namespace.GetPodsAlphabetical() {
		podLengths := map[string]int{
			"name":     slices.Max(nameLengths),
			"ready":    len("READY"),
			"status":   slices.Max(statusLengths),
			"restarts": len("RESTARTS"),
		}
		podsToReturn = append(podsToReturn, aePod{
			pod:     a.namespace.Pods[podName],
			lengths: podLengths,
		})
	}
	return podsToReturn
}

func (a aeNamespace) Pressed() (fileToOpen string) {
	return ""
}

/* aePod */
type aePod struct {
	pod     *mgparser.Pod
	lengths map[string]int
}

func (a aePod) Init(mg *mgparser.Mg) {}
func (a aePod) Header() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.lengths["name"]).Render("NAME"),
		baseStyle.Width(a.lengths["ready"]).Render("READY"),
		baseStyle.Width(a.lengths["status"]).Render("STATUS"),
		baseStyle.Width(a.lengths["restarts"]).Render("RESTARTS"))
}

func (a aePod) Title() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.lengths["name"]).Render(a.pod.GetName()),
		baseStyle.Width(a.lengths["ready"]).Render(a.pod.GetOcOutput().Ready),
		baseStyle.Width(a.lengths["status"]).Render(a.pod.GetOcOutput().Status),
		baseStyle.Width(a.lengths["restarts"]).Render(strconv.Itoa(a.pod.GetOcOutput().Restarts)))
}
func (a aePod) FilterValue() string { return a.pod.GetName() }
func (a aePod) GetWidthFunc() func(windowSize int) int {
	return func(windowSize int) int {
		// return windowSize / (1 / 0.5)
		return 60
	}
}
func (a aePod) GetHeightFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return windowSize
	}
}
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
func (a aeContainer) GetWidthFunc() func(windowSize int) int {
	return func(windowSize int) int {
		//return windowSize / (1 / 0.2)
		return 60
	}
}
func (a aeContainer) GetHeightFunc() func(windowSize int) int {
	return func(windowSize int) int {
		return windowSize
	}
}
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
