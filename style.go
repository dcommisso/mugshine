package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)


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
