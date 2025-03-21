package mgparser

import (
	"errors"
	"os"
	"strings"
)

// The pod info as returned by kubectl/oc.
// It's also used for containers, even if
// kubectl/oc doesn't provide info about containers.
// TODO: add `Age`
type OcOutput struct {
	Ready, Status string
	Restarts      int
}

type Mg struct {
	Timestamp struct {
		Start string
		End   string
	}
	basePath   string
	Namespaces map[string]*Namespace
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

	const (
		namespaceDir = "namespaces"
	)

	dirNumber := 0
	timestampFileFound := false
	namespaceDirFound := false
	mgBasePath := ""

	for _, file := range files {
		if file.Name() == "timestamp" && !file.IsDir() {
			timestampFileFound = true
		}

		if file.IsDir() {
			dirNumber++
			if file.Name() == "namespaces" {
				// inspect case
				namespaceDirFound = true
				mgBasePath = strings.TrimSuffix(directory, "/")
			} else if directoryContains(directory+"/"+file.Name(), "namespaces") {
				// must-gather case
				namespaceDirFound = true
				mgBasePath = strings.TrimSuffix(directory, "/") + "/" + file.Name()
			}
		}
	}

	if !timestampFileFound || !namespaceDirFound || dirNumber != 1 {
		return nil, errors.New("bad must-gather format")
	}

	// create namespaces
	namespacesToReturn := map[string]*Namespace{}
	namespacePath := mgBasePath + "/" + namespaceDir
	files, _ = os.ReadDir(namespacePath)
	for _, file := range files {
		if file.IsDir() {
			namespacesToReturn[file.Name()] = newNamespace(namespacePath + "/" + file.Name())
		}
	}

	return &Mg{
		basePath:   mgBasePath,
		Namespaces: namespacesToReturn,
	}, nil
}

func (m *Mg) GetNamespacesAlphabetical() []string {
	return getAlphabeticalKeys(m.Namespaces)
}
