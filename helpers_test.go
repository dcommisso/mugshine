package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFieldsInLines(t *testing.T) {
	cases := map[string]struct {
		input          map[string][]string
		expectedOutput map[string]string
	}{
		"first case": {
			input: map[string][]string{
				"#HEADER#":                                []string{"NAME", "READY", "STATUS", "RESTARTS"},
				"catalog-operator-869689d5d5-fkn5r":       []string{"catalog-operator-869689d5d5-fkn5r", "1/1", "Running", "0"},
				"collect-profiles-29006445-kvmhz":         []string{"collect-profiles-29006445-kvmhz", "0/1", "Completed", "0"},
				"collect-profiles-29006460-jnr5j":         []string{"collect-profiles-29006460-jnr5j", "0/1", "Completed", "0"},
				"collect-profiles-29006475-hf5hr":         []string{"collect-profiles-29006475-hf5hr", "0/1", "Completed", "0"},
				"olm-operator-676d4fb97-7v8pb":            []string{"olm-operator-676d4fb97-7v8pb", "1/1", "Running", "0"},
				"package-server-manager-6984bc476c-w6jf7": []string{"package-server-manager-6984bc476c-w6jf7", "2/2", "Running", "1"},
				"packageserver-dffff7b4-8lw4h":            []string{"packageserver-dffff7b4-8lw4h", "1/1", "Running", "0"},
				"packageserver-dffff7b4-pblk7":            []string{"packageserver-dffff7b4-pblk7", "1/1", "Running", "0"},
			},
			expectedOutput: map[string]string{
				"#HEADER#":                                "NAME                                     READY  STATUS     RESTARTS",
				"catalog-operator-869689d5d5-fkn5r":       "catalog-operator-869689d5d5-fkn5r        1/1    Running    0       ",
				"collect-profiles-29006445-kvmhz":         "collect-profiles-29006445-kvmhz          0/1    Completed  0       ",
				"collect-profiles-29006460-jnr5j":         "collect-profiles-29006460-jnr5j          0/1    Completed  0       ",
				"collect-profiles-29006475-hf5hr":         "collect-profiles-29006475-hf5hr          0/1    Completed  0       ",
				"olm-operator-676d4fb97-7v8pb":            "olm-operator-676d4fb97-7v8pb             1/1    Running    0       ",
				"package-server-manager-6984bc476c-w6jf7": "package-server-manager-6984bc476c-w6jf7  2/2    Running    1       ",
				"packageserver-dffff7b4-8lw4h":            "packageserver-dffff7b4-8lw4h             1/1    Running    0       ",
				"packageserver-dffff7b4-pblk7":            "packageserver-dffff7b4-pblk7             1/1    Running    0       ",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, formatFieldsInLines(tc.input))
		})
	}
}
