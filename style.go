package main

import (
	"github.com/charmbracelet/lipgloss"
)

/* Change here the colors */
var (
	colorSelectedItem = lipgloss.Color("#35A7FF")
	colorFocusBorder  = lipgloss.Color("#E4B363")
	colorHeader       = lipgloss.Color("#048BA8")
	colorFailed       = lipgloss.Color("#E84855")
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
