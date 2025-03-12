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
	list       list.Model
	widthFunc  func(windowSize int) int
	heightFunc func(windowSize int) int
	active     bool
	status     PanelStatus
}

func (p *panel) SetStatus(status PanelStatus) {
	// set the new status
	p.status = status

	// The render of the elements should happens in delegate but from there it's
	// not possible/easy (without doing ugly things) to check the panel status
	d := NewMgDelegate()
	if status == PanelStatusNext {
		// avoid highlighting elements in next panel
		d.Styles.SelectedTitle = d.Styles.NormalTitle
		d.Styles.FailedSelectedTitle = d.Styles.FailedTitle
	}

	// the new delegate must be set in any status, otherwise when the panel
	// becomes focused again, it retains the previous style
	p.list.SetDelegate(d)
}

func (p *panel) getStyle() lipgloss.Style {
	if p.status == PanelStatusFocused {
		return getFocusedStyle()
	}
	return lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
}

func (p panel) Init() tea.Cmd {
	return nil
}

func (p panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return p.getStyle().Render(p.list.View())
}

func (p *panel) setSize(windowWidth, windowHeight int) {
	p.list.SetWidth(p.widthFunc(windowWidth))
	p.list.SetHeight(p.heightFunc(windowHeight))
}
