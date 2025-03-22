package mgparser

import (
	"errors"
	configv1 "github.com/openshift/api/config/v1"
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
	mgPath         string
	basePath       string
	timestamp      []string
	Namespaces     map[string]*Namespace
	infrastructure *configv1.Infrastructure
	clusterVersion *configv1.ClusterVersion
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
		namespaceDir            = "namespaces"
		infrastructuresFilePath = "cluster-scoped-resources/config.openshift.io/infrastructures.yaml"
		clusterVersionPath      = "cluster-scoped-resources/config.openshift.io/clusterversions/version.yaml"
		timestampPath           = "timestamp"
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

	// parse timestamp
	timestampFile, err := os.ReadFile(mgBasePath + "/" + timestampPath)
	if err != nil {
		return nil, err
	}
	timestampStartEnd := strings.Split(string(timestampFile), "\n")

	// parse Infrastructure, if present
	infrastructureFile := mgBasePath + "/" + infrastructuresFilePath
	var infrastructure *configv1.Infrastructure
	if _, err := os.Stat(infrastructureFile); err == nil {
		infralist, err := parseInfrastructureList(infrastructureFile)
		if err != nil {
			return nil, err
		}
		infrastructure = &infralist.Items[0]
	}

	// parse clusterVersion, if present
	clusterVersionFile := mgBasePath + "/" + clusterVersionPath
	var clusterVersion *configv1.ClusterVersion
	if _, err := os.Stat(clusterVersionFile); err == nil {
		cv, err := parseClusterVersion(clusterVersionFile)
		if err != nil {
			return nil, err
		}
		clusterVersion = &cv
	}

	return &Mg{
		mgPath:         strings.TrimSuffix(directory, "/"),
		basePath:       mgBasePath,
		timestamp:      timestampStartEnd,
		Namespaces:     namespacesToReturn,
		infrastructure: infrastructure,
		clusterVersion: clusterVersion,
	}, nil
}

func (m *Mg) GetMgPath() string {
	return m.mgPath
}

func (m *Mg) GetNamespacesAlphabetical() []string {
	return getAlphabeticalKeys(m.Namespaces)
}
