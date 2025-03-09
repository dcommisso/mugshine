package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func getNextDelegate() list.DefaultDelegate {
	nextDelegate := list.NewDefaultDelegate()

	// avoid highlighting elements in next panel
	nextDelegate.Styles.SelectedTitle = nextDelegate.Styles.NormalTitle
	nextDelegate.Styles.SelectedDesc = nextDelegate.Styles.NormalDesc

	return nextDelegate
}

func (p *panel) getDelegate() list.DefaultDelegate {
	switch p.status {
	case PanelStatusNext:
		return getNextDelegate()
	default:
		return list.NewDefaultDelegate()
	}
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
