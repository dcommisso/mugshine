package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dcommisso/img/internal/mgparser"
)

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
	m := model{}
	m.addNewPanel(0, ocpResources)
	m.panels[0].active = true

	m.updateNextPanel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
