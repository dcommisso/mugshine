package main

import "github.com/dcommisso/mugshine/internal/mgparser"

/*
 * This is the section with the nodes
 *
 * This is the structure:
 *
 * |   aeNodeResource   |  aeNode    |
 * +++++++++++++++++++++++++++++++++++
 *          Nodes      --> master-01
 *                         master-02
 *                         worker-01
 */

/* aeNodeResource */
const aeNodesSectionName = "Nodes"

type aeNodesResource struct {
	mg *mgparser.Mg
}

func (a *aeNodesResource) Init(mg *mgparser.Mg) {
	a.mg = mg
}
func (a aeNodesResource) Header() string      { return "OCP RESOURCES" }
func (a aeNodesResource) Title() string       { return aeNodesSectionName }
func (a aeNodesResource) FilterValue() string { return a.Title() }
func (a aeNodesResource) IsFailed() bool {
	for _, node := range a.Selected() {
		if node.IsFailed() {
			return true
		}
	}
	return false
}

func (a aeNodesResource) Selected() []ActionableElement {
	nodesToReturn := []ActionableElement{}

	// format fields
	resourceFields := map[string][]string{
		"#HEADER#": []string{"NAME", "STATUS", "ROLES", "VERSION"},
	}
	for nodeName, node := range a.mg.Nodes {
		resourceFields[nodeName] = []string{
			nodeName,
			node.GetStatus(),
			node.GetRoles(),
			node.GetVersion(),
		}
	}
	formattedFields := formatFieldsInLines(resourceFields)

	for _, nodeName := range a.mg.GetNodesAlphabetical() {
		nodesToReturn = append(nodesToReturn, aeNode{
			node:   a.mg.Nodes[nodeName],
			title:  formattedFields[nodeName],
			header: formattedFields["#HEADER#"],
		})
	}
	return nodesToReturn
}

func (a aeNodesResource) Pressed() (fileToOpen string) { return "" }

/* aeNode */
type aeNode struct {
	node          *mgparser.Node
	title, header string
}

func (a aeNode) Init(mg *mgparser.Mg) {}
func (a aeNode) Header() string       { return a.header }
func (a aeNode) Title() string        { return a.title }
func (a aeNode) FilterValue() string  { return a.node.GetName() }
func (a aeNode) IsFailed() bool {
	return a.node.GetStatus() != "Ready"
}

func (a aeNode) Selected() []ActionableElement {
	return nil
}

func (a aeNode) Pressed() (fileToOpen string) {
	return a.node.GetManifestFilePath()
}
