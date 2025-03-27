package main

import (
	"github.com/charmbracelet/lipgloss"
)

/* Change here the colors */
var (
	colorSelectedItem      = lipgloss.Color("#4CE0B3")
	colorFocusBorder       = lipgloss.Color("#F9A03F")
	colorHeader            = lipgloss.Color("#084B83")
	colorFailed            = lipgloss.Color("#cc0000")
	colorClusterInfoBorder = lipgloss.Color("#FAA6FF")
)

/* Change here the style of focused panel */
func getFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorFocusBorder)
}

/* Change here the style of list header */
func getListTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(colorHeader)
}
