package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	p.status = status
	p.list.SetDelegate(p.getDelegate())
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
