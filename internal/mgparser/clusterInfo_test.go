package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inspectDir = "./testdata/mgs/validInspect"
	mgDir      = "./testdata/mgs/validMg"
)

func TestInfrastructure(t *testing.T) {
	cases := map[string]struct {
		mgpath       string
		apiserverurl string
		platform     string
	}{
		"validMG": {
			mgpath:       mgDir,
			apiserverurl: "https://api.clustername.domain.local:6443",
			platform:     "VSphere",
		},
		"validInspect": {
			mgpath:       inspectDir,
			apiserverurl: "",
			platform:     "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg, _ := NewMg(tc.mgpath)
			assert.Equal(t, tc.apiserverurl, mg.GetApiServerURL())
			assert.Equal(t, tc.platform, mg.GetPlatform())
		})
	}
}
