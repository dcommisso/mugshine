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

func (p panel) Update(msg tea.Msg) (panel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.setSize(msg.Width, msg.Height-5)
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

func (p *panel) setSize(newWidth, newHeight int) {
	p.setWidth(newWidth)
	p.setHeight(newHeight)
}

func (p *panel) setWidth(newWidth int) {
	p.list.SetWidth(newWidth)
}

func (p *panel) setHeight(newHeight int) {
	p.list.SetHeight(newHeight)
}

// GetWantedWidth returns the width required by the panel for a full output
func (p *panel) GetWantedWidth() int {
	// Since the list.View() output is limited by the width of the list, here we
	// create a fake panel with arbitrarily large width, in order to run a
	// View() and find out the maximum width needed by the panel. Yes, it's ugly.
	fakePanel := panel{
		list:   list.New(p.list.Items(), NewMgDelegate(), 1000, 0),
		status: p.status,
	}

	// Set header. This is important to calculate width based on header.
	fakePanel.list.Styles.Title = getListTitleStyle()
	if header := fakePanel.list.Items()[0].(ActionableElement).Header(); header != "" {
		fakePanel.list.Title = header
	}

	return lipgloss.Width(fakePanel.View())
}
