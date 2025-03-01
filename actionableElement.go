package img

import "github.com/dcommisso/img/internal/mgparser"

type ActionableElement interface {
	// Init is called to initialize the ActionableElement with the provided
	// must-gather.
	Init(mg mgparser.Mg)

	// OutputLine returns the string to display in the list.
	OutputLine() string

	// IsFailed returns the status of the selected element.
	IsFailed() bool

	// Selected returns the list of resources to show in the next column when
	// the element is is selected.
	Selected(param string) []ActionableElement

	// Pressed returns the file to open when the enter key is pressed on the
	// selected element.
	Pressed() (fileToOpen string)
}
