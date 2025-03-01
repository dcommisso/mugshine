package mgparser

import (
	"errors"
	"os"
	"slices"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// directoryContains looks into `directory` and returns true if a file or
// directory with `name` is present. It doesn't recurse into folders.
func directoryContains(directory, name string) bool {
	files, _ := os.ReadDir(directory)
	for _, file := range files {
		if file.Name() == name {
			return true
		}
	}
	return false
}

// getAlphabeticalKeys takes a map and return an alphabetical ordered slice of
// map's keys
func getAlphabeticalKeys[V any](mapWithKeysToOrder map[string]V) []string {
	alphabeticalKeys := make([]string, len(mapWithKeysToOrder))
	i := 0
	for k := range mapWithKeysToOrder {
		alphabeticalKeys[i] = k
		i++
	}
	slices.Sort(alphabeticalKeys)
	return alphabeticalKeys
}

// getPods parses a yaml manifest file containing a single Pod or a PodList and
// returns a slice of pods.
func getPods(filename string) ([]v1.Pod, error) {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	obj, groupVersionKind, err := decoder.Decode(
		manifest,
		nil,
		nil)
	if err != nil {
		return nil, err
	}

	if groupVersionKind.Group == "" && groupVersionKind.Version == "v1" {
		if groupVersionKind.Kind == "Pod" {
			pod, ok := obj.(*v1.Pod)
			if !ok {
				return nil, errors.New("Cannot convert obj to v1.Pod")
			}
			return []v1.Pod{*pod}, nil
		}

		if groupVersionKind.Kind == "PodList" {
			pods, ok := obj.(*v1.PodList)
			if !ok {
				return nil, errors.New("Cannot convert obj to v1.PodList")
			}
			return pods.Items, nil
		}
	}

	return nil, errors.New("Resource is neither Pod nor PodList")
}
