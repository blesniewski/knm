package cryptoexchange

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHappyPath(t *testing.T) {
	client := NewClient()

	pair, err := client.GetConversionRate("WBTC", "USDT", 1.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedAmount := 57094.314314
	require.Equal(t, expectedAmount, pair.Amount)
}

func TestParis(t *testing.T) {
	tc := []struct {
		from      string
		to        string
		amount    float64
		expected  float64
		shouldErr bool
	}{
		{"WBTC", "USDT", 1.0, 57094.314314, false},
		{"USDT", "WBTC", 1.0, 0.00001751, false},
		{"USDT", "BEER", 1.0, 40593.25477448192, false},
		{"MATIC", "GATE", 0.999, 0.0, true},
		{"USDT", "GATE", 0.0, 0.0, true},
	}
	client := NewClient()
	for _, tt := range tc {
		t.Run(fmt.Sprintf("from=%s_to=%s_amount=%f", tt.from, tt.to, tt.amount), func(t *testing.T) {
			pair, err := client.GetConversionRate(tt.from, tt.to, tt.amount)
			if (err != nil) != tt.shouldErr {
				t.Fatalf("expected error: %v, got: %v", tt.shouldErr, err)
			}
			require.Equal(t, tt.expected, pair.Amount)
		})
	}
}
