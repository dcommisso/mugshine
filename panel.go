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

func (m *model) AddNewPanel(index int, actElements []ActionableElement) {
	items := aeSliceToItem(actElements)
	widthFunc := actElements[0].GetWidthFunc()
	heightFunc := actElements[0].GetHeightFunc()
	model := list.New(items, list.NewDefaultDelegate(), widthFunc(m.windowWidth), heightFunc(m.windowHeight))
	m.panels[index] = panel{
		list:       model,
		widthFunc:  widthFunc,
		heightFunc: heightFunc,
		active:     true,
	}
}

// updateNextPanel creates/updates/deletes the next panel based on the item
// selected in focused panel
func (m *model) UpdateNextPanel() {
	// do nothing if we are at last panel
	lastPanelIndex := len(m.panels) - 1
	if m.focused == lastPanelIndex {
		return
	}

	selectedItem := m.panels[m.focused].list.SelectedItem()
	itemsInNextPanel := selectedItem.(ActionableElement).Selected()

	// Add a new panel only if there are elements to show, otherwise delete it.
	if len(itemsInNextPanel) > 0 {
		m.AddNewPanel(m.focused+1, itemsInNextPanel)
		m.panels[m.focused+1].SetStatus(PanelStatusNext)
	} else {
		m.deletePanel(m.focused + 1)
	}
}

func (m *model) deletePanel(index int) {
	m.panels[index].active = false
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
