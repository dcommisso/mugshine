package mgparser

import (
	"errors"
	"os"
	"strings"
)

type Mg struct {
	Timestamp struct {
		Start string
		End   string
	}
	NamespacesPath string
	Namespaces     map[string]*Namespace
}

// NewMg return an instance of Mg.
//
// The rules for a good mg/inspect are:
// 1. directory must contain a file named "timestamp"
// 2. directory must contain exactly 1 directory
// 2b. the name of that directory must be "namespaces"
// (this is the case of inspect)
// 2c. if the name of that directory is different than
// "namespaces", then it must contain a directory named
// "namespaces" (this is the case of must-gather)
func NewMg(directory string) (*Mg, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	dirNumber := 0
	timestampFile := false
	namespaceDir := false
	namespacePath := ""

	for _, file := range files {
		if file.Name() == "timestamp" && !file.IsDir() {
			timestampFile = true
		}

		if file.IsDir() {
			dirNumber++
			if file.Name() == "namespaces" {
				// inspect case
				namespaceDir = true
				namespacePath = strings.TrimSuffix(directory, "/") + "/namespaces"
			} else if directoryContains(directory+"/"+file.Name(), "namespaces") {
				// must-gather case
				namespaceDir = true
				namespacePath = strings.TrimSuffix(directory, "/") + "/" + file.Name() + "/namespaces"
			}
		}
	}

	if !timestampFile || !namespaceDir || dirNumber != 1 {
		return nil, errors.New("bad must-gather format")
	}

	// create namespaces
	namespacesToReturn := map[string]*Namespace{}
	files, _ = os.ReadDir(namespacePath)
	for _, file := range files {
		if file.IsDir() {
			namespacesToReturn[file.Name()] = newNamespace(namespacePath + "/" + file.Name())
		}
	}

	return &Mg{
		NamespacesPath: namespacePath,
		Namespaces:     namespacesToReturn,
	}, nil
}

func (m *Mg) GetNamespacesAlphabetical() []string {
	return getAlphabeticalKeys(m.Namespaces)
}
