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

type model struct {
	panels  [maxNumberOfPanels]list.Model
	focused int
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
			if m.focused < maxNumberOfPanels-1 {
				m.focused++
			}
		case "left":
			if m.focused > 0 {
				m.focused--
			}
		}
	}

	// send the msg to the focused panel. This must be done before getting the
	// selected item.
	var cmd tea.Cmd
	m.panels[m.focused], cmd = m.panels[m.focused].Update(msg)

	// Get the selected items and build the next panel
	selected := m.panels[m.focused].SelectedItem()
	m.panels[m.focused+1] = list.New(aeSliceToItem(selected.(ActionableElement).Selected()), list.NewDefaultDelegate(), 50, 20)

	return m, cmd
}

func (m model) View() string {
	// return m.panels[m.focused].View()
	return lipgloss.JoinHorizontal(lipgloss.Top, m.panels[0].View(), m.panels[1].View(), m.panels[2].View(), m.panels[3].View())
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
		panels: [maxNumberOfPanels]list.Model{
			list.New(aeSliceToItem(ocpResources), list.NewDefaultDelegate(), 50, 10),
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
