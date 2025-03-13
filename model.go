package main

import (
	"github.com/charmbracelet/bubbles/list"
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

func (m *model) AddNewPanel(index int, actElements []ActionableElement) {
	items := aeSliceToItem(actElements)
	model := list.New(items, NewMgDelegate(), 10000, 40)

	// disable help
	model.SetShowHelp(false)

	// remap left/right to pagedown/pageup to avoid overlapping to navigation
	model.KeyMap.NextPage.SetKeys("pgdown")
	model.KeyMap.PrevPage.SetKeys("pgup")

	// header section
	model.Styles.Title = getListTitleStyle()
	if header := actElements[0].Header(); header != "" {
		model.Title = header
	}

	m.panels[index] = panel{
		list:   model,
		active: true,
	}
}

// UpdatePanelsAfterMoving creates/updates/deletes the panels based on the item
// selected in focused panel
func (m *model) UpdatePanelsAfterMoving() {
	// update the status of the panels
	m.panels[m.focused].status = PanelStatusFocused
	if m.focused-1 >= 0 {
		m.panels[m.focused-1].status = PanelStatusPrevious
	}

	// do nothing else if we are at last panel
	lastPanelIndex := len(m.panels) - 1
	if m.focused == lastPanelIndex {
		return
	}

	// delete the next-next panel
	nextNextPanelIndex := m.focused + 2
	if nextNextPanelIndex < len(m.panels) {
		m.deletePanel(nextNextPanelIndex)
	}

	selectedItem := m.panels[m.focused].list.SelectedItem()

	// return if selectedItem is nil, otherwise it crashes when filtered with
	// unmatching strings
	if selectedItem == nil {
		return
	}

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

func (m *model) setSize(width, height int) {
	m.windowWidth = width
	m.windowHeight = height
}

func (m *model) GetActivePanels() []panel {
	var activePanels []panel
	for _, panel := range m.panels {
		if panel.active {
			activePanels = append(activePanels, panel)
		}
	}
	return activePanels
}

func (m *model) IncreaseFocused() {
	// Do nothing if we're on the last panel or if the next panel is empty
	if m.focused == len(m.panels)-1 || !m.panels[m.focused+1].active {
		return
	}
	m.focused++
}

func (m *model) DecreaseFocused() {
	if m.focused == 0 {
		return
	}

	m.focused--
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
		for i, p := range m.GetActivePanels() {
			m.panels[i], cmd = p.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "right", "tab":
			m.IncreaseFocused()
		case "left", "shift+tab":
			m.DecreaseFocused()
		}
	}

	m.panels[m.focused], cmd = m.panels[m.focused].Update(msg)
	m.UpdatePanelsAfterMoving()

	return m, cmd
}

func (m model) View() string {
	activePanels := []string{}
	for _, panel := range m.GetActivePanels() {
		activePanels = append(activePanels, panel.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, activePanels...)
}
