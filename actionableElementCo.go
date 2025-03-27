package main

import "github.com/dcommisso/mugshine/internal/mgparser"

/*
 * This is the section with the cluster operators
 *
 * This is the structure:
 *
 * |   aeCoResource     |    aeCo        |
 * +++++++++++++++++++++++++++++++++++++++
 *   Cluster operators --> authentication
 *                         etcd
 *                         ...
 */

/* aeCoResource */
const aeCoSectionName = "Cluster operators"

type aeCoResource struct {
	mg *mgparser.Mg
}

func (a *aeCoResource) Init(mg *mgparser.Mg) {
	a.mg = mg
}
func (a aeCoResource) Header() string      { return "OCP RESOURCES" }
func (a aeCoResource) Title() string       { return aeCoSectionName }
func (a aeCoResource) FilterValue() string { return a.Title() }
func (a aeCoResource) IsFailed() bool {
	for _, co := range a.Selected() {
		if co.IsFailed() {
			return true
		}
	}
	return false
}

func (a aeCoResource) Selected() []ActionableElement {
	coToReturn := []ActionableElement{}

	// format fields
	resourceFields := map[string][]string{
		"#HEADER#": []string{"NAME", "VERSION", "AVAILABLE", "PROGRESSING", "DEGRADED"},
	}
	for coName, co := range a.mg.ClusterOperators {
		resourceFields[coName] = []string{
			coName,
			co.GetVersion(),
			co.GetAvailableStatus(),
			co.GetProgressingStatus(),
			co.GetDegradedStatus(),
		}
	}
	formattedFields := formatFieldsInLines(resourceFields)

	for _, coName := range a.mg.GetClusterOperatorsAlphabetical() {
		coToReturn = append(coToReturn, aeCo{
			co:     a.mg.ClusterOperators[coName],
			title:  formattedFields[coName],
			header: formattedFields["#HEADER#"],
		})
	}
	return coToReturn
}

func (a aeCoResource) Pressed() (fileToOpen string) { return "" }

/* aeCo */
type aeCo struct {
	co            mgparser.ClusterOperator
	title, header string
}

func (a aeCo) Init(mg *mgparser.Mg) {}
func (a aeCo) Header() string       { return a.header }
func (a aeCo) Title() string        { return a.title }
func (a aeCo) FilterValue() string  { return a.co.GetName() }
func (a aeCo) IsFailed() bool {
	if a.co.GetAvailableStatus() != "True" || a.co.GetDegradedStatus() == "True" {
		return true
	}
	return false
}

func (a aeCo) Selected() []ActionableElement {
	return nil
}

func (a aeCo) Pressed() (fileToOpen string) {
	return a.co.GetManifestFilePath()
}
