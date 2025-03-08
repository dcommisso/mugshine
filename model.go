package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxNumberOfPanels = 4
)

type model struct {
	panels       [maxNumberOfPanels]panel
	focused      int
	windowWidth  int
	windowHeight int
}

func (m *model) setSize(width, height int) {
	m.windowWidth = width
	m.windowHeight = height
}

func (m *model) getActivePanels() []panel {
	var activePanels []panel
	for _, panel := range m.panels {
		if panel.active {
			activePanels = append(activePanels, panel)
		}
	}
	return activePanels
}

func (m *model) increaseFocused() {
	// Do nothing if we're on the last panel or if the next panel is empty
	if m.focused == len(m.panels)-1 || !m.panels[m.focused+1].active {
		return
	}

	m.focused++
}

func (m *model) decreaseFocused() {
	if m.focused == 0 {
		return
	}

	m.focused--
	panelToDelete := m.focused + 2
	if panelToDelete < len(m.panels) {
		m.deletePanel(panelToDelete)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		for i, p := range m.getActivePanels() {
			res, cmd := p.Update(msg)
			m.panels[i] = res.(panel)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "right":
			m.increaseFocused()
		case "left":
			m.decreaseFocused()
		}
	}

	// send the msg to the focused panel. This must be done before getting the
	// selected item.
	var res tea.Model
	res, cmd = m.panels[m.focused].Update(msg)
	m.panels[m.focused] = res.(panel)

	m.updateNextPanel()

	return m, cmd
}

func (m model) View() string {
	activePanels := []string{}
	for _, panel := range m.getActivePanels() {
		activePanels = append(activePanels, panel.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, activePanels...)
}
