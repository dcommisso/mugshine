package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PanelStatus int

const (
	PanelStatusPrevious PanelStatus = iota
	PanelStatusFocused
	PanelStatusNext
)

type panel struct {
	list   list.Model
	active bool
	status PanelStatus
}

func (p *panel) SetStatus(status PanelStatus) {
	// set the new status
	p.status = status
}

func (p *panel) getStyle() lipgloss.Style {
	if p.status == PanelStatusFocused {
		return getFocusedStyle()
	}
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
}

func (p panel) Update(msg tea.Msg) (panel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.setSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			filename := p.list.SelectedItem().(ActionableElement).Pressed()
			if filename != "" {
				return p, openEditor(filename)
			}
		}
	}
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p panel) View() string {
	// This is to avoid highlighting elements in next panel. Render of the
	// elements should controlled in delegate but from there it's not
	// possible/easy (without doing ugly things) to check the panel status
	d := NewMgDelegate()
	if p.status == PanelStatusNext {
		d.Styles.SelectedTitle = d.Styles.NormalTitle
		d.Styles.FailedSelectedTitle = d.Styles.FailedTitle
	}
	// the new delegate must be set in any status, otherwise when the panel
	// becomes focused again, it retains the previous style
	p.list.SetDelegate(d)

	return p.getStyle().Render(p.list.View())
}

func (p *panel) setSize(newWidth, newHeight int) {
	p.list.SetWidth(newWidth)
	p.list.SetHeight(newHeight)
}

// GetWantedWidth returns the width required by the panel for a full output
func (p *panel) GetWantedWidth() int {
	return lipgloss.Width(p.View())
}
