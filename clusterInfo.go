package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dcommisso/img/internal/mgparser"
	"path/filepath"
)

type ClusterInfoPanel struct {
	params map[ClusterInfoParam]struct {
		label, value string
	}
	isInspect     bool
	labelStyle    lipgloss.Style
	valueStyle    lipgloss.Style
	labelValueSep string
}

type ClusterInfoParam string

func NewClusterInfoPanel(mg *mgparser.Mg) ClusterInfoPanel {
	timestampStart, _ := mg.GetTimestampStart()
	timestampEnd, _ := mg.GetTimestampEnd()
	infos := map[ClusterInfoParam]struct {
		label, value string
	}{
		"ApiServerURL": {
			label: "Api server URL",
			value: mg.GetApiServerURL(),
		},
		"Platform": {
			label: "Platform",
			value: mg.GetPlatform(),
		},
		"ClusterVersion": {
			label: "Cluster version",
			value: mg.GetClusterVersion(),
		},
		"ClusterID": {
			label: "Cluster ID",
			value: mg.GetClusterID(),
		},
		"Timestamp": {
			label: "Timestamp",
			value: timestampStart + "\n" + timestampEnd,
		},
		"MgFilename": {
			label: "Must-gather",
			value: filepath.Base(mg.GetMgPath()),
		},
	}

	return ClusterInfoPanel{
		params:        infos,
		isInspect:     mg.IsInspect(),
		labelStyle:    lipgloss.NewStyle(),
		valueStyle:    lipgloss.NewStyle(),
		labelValueSep: ":",
	}
}

func (c ClusterInfoPanel) GetFormattedColumn(params ...ClusterInfoParam) string {
	var labels, values, formatted []string

	for _, param := range params {
		labels = append(labels, c.params[param].label)
		values = append(values, c.params[param].value)
	}

	labelsMaxWidth := getMaxWidth(labels) + len(c.labelValueSep)
	valuesMaxWidth := getMaxWidth(values)

	formattedLabelStyle := c.labelStyle.Width(labelsMaxWidth).Margin(0, 1, 0, 1).Bold(true)
	formattedValueStyle := c.valueStyle.Width(valuesMaxWidth).Margin(0, 1, 0, 1)

	for _, param := range params {
		labelFormatted := formattedLabelStyle.Render(c.params[param].label + c.labelValueSep)
		valueFormatted := formattedValueStyle.Render(c.params[param].value)
		labelAndValueFormatted := lipgloss.JoinHorizontal(lipgloss.Top, labelFormatted, valueFormatted)
		formatted = append(formatted, labelAndValueFormatted)
	}
	return lipgloss.JoinVertical(lipgloss.Left, formatted...)
}

func (c ClusterInfoPanel) Render(availableWidth int) string {
	var clusterType string
	if c.isInspect == true {
		clusterType = "inspect"
	} else {
		clusterType = "mg"
	}

	switch clusterType {
	case "inspect":
		minimumHeader := c.GetFormattedColumn("Timestamp", "MgFilename")
		fullHeader := clusterInfoBorder.Render(minimumHeader)

		fullHeaderWidth := lipgloss.Width(fullHeader)

		switch {
		case availableWidth < fullHeaderWidth:
			return minimumHeader
		default:
			return fullHeader
		}
	case "mg":
		firstColumn := c.GetFormattedColumn("ApiServerURL", "Platform", "ClusterVersion", "MgFilename")
		secondColumn := c.GetFormattedColumn("ClusterID", "Timestamp")

		minimumHeader := lipgloss.JoinVertical(lipgloss.Left, firstColumn, secondColumn)
		reducedHeader := clusterInfoBorder.Render(minimumHeader)
		fullHeader := clusterInfoBorder.Render(
			lipgloss.JoinHorizontal(lipgloss.Top, firstColumn, secondColumn))

		reducedHeaderWidth := lipgloss.Width(reducedHeader)
		fullHeaderWidth := lipgloss.Width(fullHeader)

		switch {
		case availableWidth < reducedHeaderWidth:
			return minimumHeader
		case availableWidth < fullHeaderWidth:
			return reducedHeader
		default:
			return fullHeader
		}
	}
	return ""
}

/* STYLE SECTION  */
var (
	clusterInfoBorder = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true, false, true, false).
		BorderForeground(colorClusterInfoBorder)
)
