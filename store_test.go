package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseTimeRange(t *testing.T) {
	tr, err := ParseTimeRange("Mon 16:00-18:00")
	require.NoError(t, err)
	assert.Equal(t, TimeRange{
		Weekday: time.Monday,
		Start: Clock{
			Hour:   16,
			Minute: 0,
		},
		End: Clock{
			Hour:   18,
			Minute: 0,
		},
	}, tr)
	assert.Equal(t, "Mon 16:00-18:00", tr.String())
}
