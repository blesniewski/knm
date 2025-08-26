package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundToPrecision(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		precision int
		expected  float64
	}{
		{"round up", 123.456789, 2, 123.46},
		{"round down", 123.454, 2, 123.45},
		{"round to integer", 123.4, 0, 123},
		{"round up to next integer", 123.5, 0, 124},
		{"round to original value", 123.456789, -1, 123.456789}, // Negative precision should return the original value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundToPrecision(tt.value, tt.precision)
			require.Equal(t, tt.expected, result)
		})
	}
}
