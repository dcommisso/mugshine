package mgparser

import (
	"strings"
	"time"
)

func (m *Mg) GetTimestampStart() (string, error) {
	time, err := getConciseTimestamp(m.timestamp[0])
	if err != nil {
		return "", err
	}
	return time, nil
}

func (m *Mg) GetTimestampEnd() (string, error) {
	time, err := getConciseTimestamp(m.timestamp[1])
	if err != nil {
		return "", err
	}
	return time, nil
}

func getConciseTimestamp(mgTimestamp string) (string, error) {
	// remove the m+0.123... part at the end
	lastSpaceIndex := strings.LastIndex(mgTimestamp, " ")
	mgTimestampCleaned := mgTimestamp[:lastSpaceIndex]

	const (
		mgLayout      = "2006-01-02 15:04:05.999999999 -0700 MST"
		conciseLayout = "2006-01-02 15:04:05 -0700 MST"
	)

	t, err := time.Parse(mgLayout, mgTimestampCleaned)
	if err != nil {
		return "", err
	}

	return t.Format(conciseLayout), nil
}
