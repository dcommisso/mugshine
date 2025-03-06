package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dcommisso/img/internal/mgparser"
)

const (
	maxNumberOfPanels = 4
)

type panel struct {
	list.Model
	active bool
}

type model struct {
	panels  [maxNumberOfPanels]panel
	focused int
}

func (m *model) newPanel(index int, model list.Model) {
	m.panels[index].Model = model
	m.panels[index].active = true
}

func (m *model) deletePanel(index int) {
	m.panels[index].active = false
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
	switch msg := msg.(type) {
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
	var cmd tea.Cmd
	m.panels[m.focused].Model, cmd = m.panels[m.focused].Update(msg)

	// Get the selected items and build the next panel
	selectedItem := m.panels[m.focused].SelectedItem()
	itemsInNextPanel := aeSliceToItem(selectedItem.(ActionableElement).Selected())
	lastPanelIndex := len(m.panels) - 1
	// Add a new panel only if there are elements to show and we're not on the
	// last panel. Delete next panel if there's nothing to show (and we're not
	// on the last one)
	if len(itemsInNextPanel) > 0 && m.focused < lastPanelIndex {
		m.newPanel(m.focused+1, list.New(itemsInNextPanel, list.NewDefaultDelegate(), 50, 20))
	} else if len(itemsInNextPanel) == 0 && m.focused < lastPanelIndex {
		m.deletePanel(m.focused + 1)
	}

	return m, cmd
}

func (m model) View() string {
	activePanels := []string{}
	for _, panel := range m.panels {
		if panel.active {
			activePanels = append(activePanels, panel.View())
		} else {
			break
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, activePanels...)
}

func main() {
	mgPathToLoad := "./internal/mgparser/testdata/mgs/validMg"
	mgToLoad, err := mgparser.NewMg(mgPathToLoad)
	if err != nil {
		fmt.Println("Error loading mg:", err)
		os.Exit(1)
	}

	ocpResources := []ActionableElement{
		// a pointer is needed since aeLogs Init method has a pointer receiver
		new(aeLogs),
	}

	// initialize the ActionableElement and add them to firstPanelItems
	for _, elem := range ocpResources {
		elem.Init(mgToLoad)
	}

	m := model{
		focused: 0,
		panels: [maxNumberOfPanels]panel{
			panel{
				Model:  list.New(aeSliceToItem(ocpResources), list.NewDefaultDelegate(), 50, 10),
				active: true,
			},
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
