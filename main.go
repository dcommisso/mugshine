package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// send the msg to the focused panel
	var cmd tea.Cmd
	m.panels[m.focused], cmd = m.panels[m.focused].Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.panels[m.focused].View()
}

func main() {
	mgPathToLoad := "./internal/mgparser/testdata/mgs/validMg"
	mgToLoad, err := mgparser.NewMg(mgPathToLoad)
	if err != nil {
		fmt.Println("Error loading mg:", err)
		os.Exit(1)
	}

	ocpResources := []ActionableElement{
		aeLogs{},
	}

	firstPanelItems := make([]list.Item, len(ocpResources))

	// initialize the ActionableElement and add them to firstPanelItems
	for i, elem := range ocpResources {
		elem.Init(mgToLoad)
		firstPanelItems[i] = elem
	}

	m := model{
		focused: 0,
		panels: [maxNumberOfPanels]list.Model{
			list.New(firstPanelItems, list.NewDefaultDelegate(), 50, 10),
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
