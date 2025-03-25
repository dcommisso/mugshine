package mgparser

import (
	"errors"
	"os"
	"path"
	"strings"

	configv1 "github.com/openshift/api/config/v1"
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
	Nodes          map[string]*Node
	infrastructure *configv1.Infrastructure
	clusterVersion *configv1.ClusterVersion
	inspect        bool
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
		nodesDir                = "cluster-scoped-resources/core/nodes"
		infrastructuresFilePath = "cluster-scoped-resources/config.openshift.io/infrastructures.yaml"
		clusterVersionPath      = "cluster-scoped-resources/config.openshift.io/clusterversions/version.yaml"
		timestampPath           = "timestamp"
	)

	dirNumber := 0
	timestampFileFound := false
	namespaceDirFound := false
	mgBasePath := ""
	var inspect bool

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
				inspect = true
			} else if directoryContains(directory+"/"+file.Name(), "namespaces") {
				// must-gather case
				namespaceDirFound = true
				mgBasePath = strings.TrimSuffix(directory, "/") + "/" + file.Name()
				inspect = false
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

	// parse nodes, if present
	nodesPath := mgBasePath + "/" + nodesDir
	nodesToReturn := map[string]*Node{}
	nodeFiles, _ := os.ReadDir(nodesPath)
	for _, file := range nodeFiles {
		if path.Ext(file.Name()) == ".yaml" || path.Ext(file.Name()) == ".yml" {
			if node, err := parseNode(nodesPath + "/" + file.Name()); err == nil {
				nodesToReturn[node.GetName()] = &Node{
					Node:         node,
					nodeFilePath: nodesPath + "/" + file.Name(),
				}
			}
		}
	}

	return &Mg{
		mgPath:         strings.TrimSuffix(directory, "/"),
		basePath:       mgBasePath,
		timestamp:      timestampStartEnd,
		Namespaces:     namespacesToReturn,
		Nodes:          nodesToReturn,
		infrastructure: infrastructure,
		clusterVersion: clusterVersion,
		inspect:        inspect,
	}, nil
}

func (m *Mg) GetMgPath() string {
	return m.mgPath
}

func (m *Mg) IsInspect() bool {
	return m.inspect
}

func (m *Mg) GetNamespacesAlphabetical() []string {
	return getAlphabeticalKeys(m.Namespaces)
}

func (m *Mg) GetNodesAlphabetical() []string {
	return getAlphabeticalKeys(m.Nodes)
}
