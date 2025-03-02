package main

import "github.com/dcommisso/img/internal/mgparser"

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

func (a aeLogs) Init(mg *mgparser.Mg) {
	a.mg = mg
}

func (a aeLogs) Title() string       { return aeLogsSectionName }
func (a aeLogs) Description() string { return "" }
func (a aeLogs) FilterValue() string { return "" }

func (a aeLogs) IsFailed() bool {
	return false
}

func (a aeLogs) Selected(param string) []ActionableElement {
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

func (a aeNamespace) Title() string       { return a.namespace.Name }
func (a aeNamespace) Description() string { return "" }
func (a aeNamespace) FilterValue() string { return a.namespace.Name }

func (a aeNamespace) IsFailed() bool {
	return false
}

func (a aeNamespace) Selected(param string) []ActionableElement {
	podsToReturn := []ActionableElement{}
	for _, podName := range a.namespace.GetPodsAlphabetical() {
		podsToReturn = append(podsToReturn, aePod{
			pod: a.namespace.Pods[podName],
		})
	}
	return podsToReturn
}

func (a aeNamespace) Pressed() (fileToOpen string) {
	return ""
}

/* aePod */
type aePod struct {
	pod *mgparser.Pod
}

func (a aePod) Init(mg *mgparser.Mg) {}

func (a aePod) Title() string       { return a.pod.GetName() }
func (a aePod) Description() string { return a.pod.GetName() }
func (a aePod) FilterValue() string { return a.pod.GetName() }

func (a aePod) IsFailed() bool {
	return false
}

func (a aePod) Selected(param string) []ActionableElement {
	return []ActionableElement{}
}

func (a aePod) Pressed() (fileToOpen string) {
	return a.pod.GetManifestFilePath()
}

/* aeContainer */
type aeContainer struct {
	mg   mgparser.Mg
	name string
}

func (a aeContainer) Init(mg mgparser.Mg)                       {}
func (a aeContainer) Title() string                             { return "" }
func (a aeContainer) IsFailed() bool                            { return false }
func (a aeContainer) Selected(param string) []ActionableElement { return nil }
func (a aeContainer) Pressed() (fileToOpen string)              { return "" }
