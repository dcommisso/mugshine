package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

/* Change here the style of elements in panels */
func getItemStyles() list.DefaultItemStyles {
	elementStyle := list.NewDefaultItemStyles()

	/* SelectedTitle style */
	elementStyle.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#35A7FF")).
		Foreground(lipgloss.Color("#35A7FF")).
		Padding(0, 0, 0, 1)

	return elementStyle
}

/* Change here the style of focused border */
func getFocusedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6bf178"))
}

/* Change here the style of list header */
func getListTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#048ba8"))
}
