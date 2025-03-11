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

func (p *panel) getStyle() lipgloss.Style {
	var (
		focused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6bf178"))
	)

	switch p.status {
	case PanelStatusFocused:
		return focused
	default:
		return lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder())
	}
}
