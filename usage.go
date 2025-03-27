package main

import (
	"bytes"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func getUsage() string {
	var (
		usage       = "Usage: mugshine MG-DIRECTORY"
		description = "Dive into must-gathers and inspects in a colorful and interactive way."
	)
	return usage + "\n" + description
}

func getLogo() string {
	colors := []lipgloss.Color{
		lipgloss.Color("#FFD300"),
		lipgloss.Color("#FFD300"),
		lipgloss.Color("#FF0000"),
		lipgloss.Color("#FF8700"),
		lipgloss.Color("#A1FF0A"),
		lipgloss.Color("#0AFF99"),
		lipgloss.Color("#147DF5"),
		lipgloss.Color("#BE0AFF"),
	}

	mugshineString := "MuGShine"
	letters := []string{}
	for i := 0; i < len(mugshineString); i++ {
		buf := new(bytes.Buffer)
		mugshineLetter := figure.NewFigure(string(mugshineString[i]), "big", false)
		figure.Write(buf, mugshineLetter)
		styledLetter := lipgloss.NewStyle().
			Foreground(colors[i]).
			Render(buf.String())
		letters = append(letters, styledLetter)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, letters...)
}
