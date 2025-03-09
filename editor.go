package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func openEditor(filename string) tea.Cmd {
	cmd := tea.ExecProcess(exec.Command("less", filename), nil)
	return cmd
}
