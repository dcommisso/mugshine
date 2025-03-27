package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
)

var mugshineVersion = "v0.1.0"

var (
	fs *flag.FlagSet

	helpFlag    bool
	versionFlag bool
)

func init() {
	// needed because this bug: https://github.com/spf13/pflag/issues/352
	fs = flag.NewFlagSet("mugshine", flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, getLogo())
		fmt.Fprintln(os.Stderr, getUsage())
		fmt.Fprintln(os.Stderr) // put a blank line before parameters
		fs.PrintDefaults()
	}

	fs.BoolVarP(&helpFlag, "help", "h", false, "Show this help")
	fs.BoolVarP(&versionFlag, "version", "v", false, "Show version number")
}

func main() {
	fs.Parse(os.Args[1:])

	switch {
	case versionFlag:
		fmt.Fprint(os.Stderr, mugshineVersion)
		os.Exit(0)
	case helpFlag:
		fs.Usage()
		os.Exit(0)
	case fs.NArg() != 1:
		fs.Usage()
		os.Exit(1)
	}

	mgBoard, err := NewMgBoard(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(mgBoard, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
