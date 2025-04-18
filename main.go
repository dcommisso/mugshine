/*
   Copyright 2025 Domenico Commisso

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
)

var mugshineVersion = "v0.1.2"

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
		fmt.Fprintln(os.Stderr, mugshineVersion)
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
