package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	PreviousPage key.Binding
	NextPage     key.Binding
	Left         key.Binding
	Right        key.Binding
	View         key.Binding
	Filter       key.Binding
	ClearFilter  key.Binding
	Help         key.Binding
	Quit         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PreviousPage, k.NextPage, k.Left, k.Right}, // first column
		{k.View, k.Filter, k.ClearFilter},                           // second column
		{k.Help, k.Quit},                                            // third column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "shift+tab"),
		key.WithHelp("←/shift+tab", "move to the previous panel"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "tab"),
		key.WithHelp("→/tab", "move to the next panel"),
	),
	PreviousPage: key.NewBinding(
		key.WithKeys("pgup", "pgup"),
		key.WithHelp("pgup", "move to the previous page"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("pgdown", "pgdown"),
		key.WithHelp("pgdown", "move to the next page"),
	),
	View: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view manifest/logs"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
