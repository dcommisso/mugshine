package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func getBaseDelegate() list.DefaultDelegate {
	baseDelegate := list.NewDefaultDelegate()
	baseDelegate.ShowDescription = false
	return baseDelegate
}

func getNextDelegate() list.DefaultDelegate {
	nextDelegate := list.NewDefaultDelegate()

	// avoid highlighting elements in next panel
	nextDelegate.Styles.SelectedTitle = nextDelegate.Styles.NormalTitle
	nextDelegate.Styles.SelectedDesc = nextDelegate.Styles.NormalDesc

	nextDelegate.ShowDescription = false
	return nextDelegate
}

func (p *panel) getDelegate() list.DefaultDelegate {
	switch p.status {
	case PanelStatusNext:
		return getNextDelegate()
	default:
		return getBaseDelegate()
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
