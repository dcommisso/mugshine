package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	mgPathToLoad := "./internal/mgparser/testdata/mgs/validMg"
	//mgPathToLoad := "./internal/mgparser/testdata/mgs/validInspect"

	mgBoard, err := NewMgBoard(mgPathToLoad)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := tea.NewProgram(mgBoard, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
