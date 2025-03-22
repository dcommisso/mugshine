package main

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func aeSliceToItem(ae []ActionableElement) []list.Item {
	itemToReturn := make([]list.Item, len(ae))
	for i, elem := range ae {
		itemToReturn[i] = elem
	}
	return itemToReturn
}

// formatFieldsInLines takes a map where the key is the name of the resource and
// the value is a slice with its fields in order. It returns a formatted string
// line for each resource, to use as title parameter for the resource to create.
// The header should be included as well, using special characters in key name
// so it cannot collide with pod names, e.g. #HEADER#.
func formatFieldsInLines(resourceFields map[string][]string) map[string]string {

	// valuesGroupedByField[i] contains a slice with all the values of the field
	// at position i
	valuesGroupedByField := map[int][]string{}
	maxWidthsByField := map[int]int{}
	formattedLines := map[string]string{}

	// group all the values for each kind of field and use them to populate
	// valuesGroupedByField
	for _, resource := range resourceFields {
		for i, field := range resource {
			valuesGroupedByField[i] = append(valuesGroupedByField[i], field)
		}
	}

	// get the max width for each field
	for i := 0; i < len(valuesGroupedByField); i++ {
		maxWidthsByField[i] = getMaxWidth(valuesGroupedByField[i])
	}

	for resource, _ := range resourceFields {
		baseStyle := lipgloss.NewStyle().Align(lipgloss.Left).MarginRight(2)
		formattedFields := []string{}

		// format each line using the width value of the related field. Don't add
		// margin to last field
		for i, field := range resourceFields[resource] {
			if i == len(resourceFields[resource])-1 {
				formattedFields = append(formattedFields,
					baseStyle.Width(maxWidthsByField[i]).UnsetMarginRight().Render(field))
			} else {
				formattedFields = append(formattedFields,
					baseStyle.Width(maxWidthsByField[i]).Render(field))
			}
		}
		formattedLines[resource] = lipgloss.JoinHorizontal(lipgloss.Top, formattedFields...)
	}
	return formattedLines
}

// getMaxWidth return the lenght of the longest string in the input slice. It
// manages multilines string too.
func getMaxWidth(values []string) int {
	var widths []int
	for _, value := range values {
		if strings.Contains(value, "\n") {
			widths = append(widths, getMaxWidth(strings.Split(value, "\n")))
			continue
		}
		widths = append(widths, len(value))
	}
	return slices.Max(widths)
}
