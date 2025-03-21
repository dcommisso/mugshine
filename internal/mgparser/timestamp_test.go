package mgparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp(t *testing.T) {
	const (
		inspectDir = "./testdata/mgs/validInspect"
		mgDir      = "./testdata/mgs/validMg"
	)

	cases := map[string]struct {
		mgpath string
		start  string
		end    string
	}{
		"validMG": {
			mgpath: mgDir,
			start:  "2025-03-20 18:07:24 +0100 CET",
			end:    "2025-03-20 18:12:41 +0100 CET",
		},
		"validInspect": {
			mgpath: inspectDir,
			start:  "2025-02-24 16:25:27 +0100 CET",
			end:    "2025-02-24 16:25:36 +0100 CET",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mg, _ := NewMg(tc.mgpath)
			timestampStartGot, _ := mg.GetTimestampStart()
			assert.Equal(t, tc.start, timestampStartGot)

			timestampEndGot, _ := mg.GetTimestampEnd()
			assert.Equal(t, tc.end, timestampEndGot)
		})
	}
}
