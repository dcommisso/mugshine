package main

import (
	"github.com/dcommisso/mugshine/internal/mgparser"
)

type ActionableElement interface {
	MgDelegateItem
	// Init is called to initialize the ActionableElement with the provided
	// must-gather.
	Init(mg *mgparser.Mg)

	// Header return the header with the field names to print just one time
	Header() string

	// IsFailed returns the status of the selected element.
	IsFailed() bool

	// Selected returns the list of resources to show in the next column when
	// the element is is selected.
	Selected() []ActionableElement

	// Pressed returns the file to open when the enter key is pressed on the
	// selected element.
	Pressed() (fileToOpen string)
}
