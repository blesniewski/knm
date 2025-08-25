package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoundToPrecision(t *testing.T) {
	tests := []struct {
		value     float64
		precision int
		expected  float64
	}{
		{123.456789, 2, 123.46},
		{123.454, 2, 123.45},
		{123.4, 0, 123},
		{123.5, 0, 124},
		{123.456789, -1, 123.456789}, // Negative precision should return the original value
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("value=%f_precision=%d", tt.value, tt.precision), func(t *testing.T) {
			result := RoundToPrecision(tt.value, tt.precision)
			require.Equal(t, tt.expected, result)
		})
	}
}
