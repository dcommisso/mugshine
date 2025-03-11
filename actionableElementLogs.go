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
func (a aeLogs) Header() string      { return "" }
func (a aeLogs) Title() string       { return aeLogsSectionName }
func (a aeLogs) Description() string { return "" }
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
func (a aeNamespace) Header() string       { return "" }
func (a aeNamespace) Title() string        { return a.namespace.Name }
func (a aeNamespace) Description() string  { return "" }
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
	return false
}

func (a aeNamespace) Selected() []ActionableElement {
	podsToReturn := []ActionableElement{}

	// calculate the max width for name and status pod fields
	var nameLenghts, statusLenghts []int
	for _, podName := range a.namespace.GetPodsAlphabetical() {
		nameLenghts = append(nameLenghts, len(podName))
		statusLenghts = append(statusLenghts, len(a.namespace.Pods[podName].GetOcOutput().Status))
	}

	for _, podName := range a.namespace.GetPodsAlphabetical() {
		podsToReturn = append(podsToReturn, aePod{
			pod:            a.namespace.Pods[podName],
			nameLenght:     slices.Max(nameLenghts),
			readyLenght:    len("READY"),
			statusLenght:   slices.Max(statusLenghts),
			restartsLenght: len("RESTARTS"),
		})
	}
	return podsToReturn
}

func (a aeNamespace) Pressed() (fileToOpen string) {
	return ""
}

/* aePod */
type aePod struct {
	pod                                                   *mgparser.Pod
	nameLenght, readyLenght, statusLenght, restartsLenght int
}

func (a aePod) Init(mg *mgparser.Mg) {}
func (a aePod) Header() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.nameLenght).Render("NAME"),
		baseStyle.Width(a.readyLenght).Render("READY"),
		baseStyle.Width(a.statusLenght).Render("STATUS"),
		baseStyle.Width(a.restartsLenght).Render("RESTARTS"))
}

func (a aePod) Title() string {
	baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		baseStyle.Width(a.nameLenght).Render(a.pod.GetName()),
		baseStyle.Width(a.readyLenght).Render(a.pod.GetOcOutput().Ready),
		baseStyle.Width(a.statusLenght).Render(a.pod.GetOcOutput().Status),
		baseStyle.Width(a.restartsLenght).Render(strconv.Itoa(a.pod.GetOcOutput().Restarts)))
}
func (a aePod) Description() string { return a.pod.GetName() }
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
	return false
}

func (a aePod) Selected() []ActionableElement {
	containersToReturn := []ActionableElement{}

	// Add containers
	for _, containerName := range a.pod.GetContainersAlphabetical() {
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.Containers[containerName],
			isInit:    false,
		})
	}

	// Add initContainers
	for _, containerName := range a.pod.GetInitContainersAlphabetical() {
		containersToReturn = append(containersToReturn, aeContainer{
			container: a.pod.InitContainers[containerName],
			isInit:    true,
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
}

func (a aeContainer) Init(mg *mgparser.Mg) {}
func (a aeContainer) Header() string       { return "" }
func (a aeContainer) Title() string {
	initHeader := ""
	if a.isInit {
		initHeader = "[INIT] "
	}
	return initHeader + a.container.Name
}

func (a aeContainer) Description() string { return "" }
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
func (a aeContainer) IsFailed() bool                { return false }
func (a aeContainer) Selected() []ActionableElement { return nil }
func (a aeContainer) Pressed() (fileToOpen string) {
	return a.container.GetLogsFilename()
}
