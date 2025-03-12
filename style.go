package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

/* Change here the colors */
var (
	colorSelectedItem = lipgloss.Color("#35A7FF")
	colorFocusBorder  = lipgloss.Color("#E4B363")
	colorHeader       = lipgloss.Color("#048BA8")
	colorFailed       = lipgloss.Color("#E84855")
)

/* Change here the style of elements in panels */
func getItemStyles() list.DefaultItemStyles {
	elementStyle := list.NewDefaultItemStyles()

	/* SelectedTitle style */
	elementStyle.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(colorSelectedItem).
		Foreground(colorSelectedItem).
		Padding(0, 0, 0, 1)

	return elementStyle
}

/* Change here the style of focused border */
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
