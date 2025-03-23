package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Open   key.Binding
	Filter key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// not yet implemented
//func (k keyMap) ShortHelp() []key.Binding {
//	return []key.Binding{k.Help, k.Quit}
//}
//
//func (k keyMap) FullHelp() [][]key.Binding {
//	return [][]key.Binding{
//		{k.Up, k.Down, k.Left, k.Right}, // first column
//		{k.Help, k.Quit},                // second column
//	}
//}

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
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view manifest/logs"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
