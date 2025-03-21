package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterVersion(t *testing.T) {
	const (
		inspectDir = "./testdata/mgs/validInspect"
		mgDir      = "./testdata/mgs/validMg"
	)

	cases := map[string]struct {
		mgpath    string
		clusterID string
	}{
		"validMG": {
			mgpath:    mgDir,
			clusterID: "asdfgh56-1234-asdf-12as-5655443fgrt6a",
		},
		"validInspect": {
			mgpath:    inspectDir,
			clusterID: "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg, _ := NewMg(tc.mgpath)
			assert.Equal(t, tc.clusterID, mg.GetClusterID())
		})
	}
}
