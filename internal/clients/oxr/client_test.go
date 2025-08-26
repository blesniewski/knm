package oxr

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/blesniewski/knm/internal/helpers"
	"github.com/blesniewski/knm/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	extraFn func()
}

func (rt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.extraFn()

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"base":"USD","rates":{"EUR":0.858611,"USD":1.00}}`)),
	}, nil
}

func setupClientForTesting(t *testing.T, extraRTFunc func()) *Client {
	t.Helper()

	client, err := NewClient(
		t.Context(),
		"https://openexchangerates.org/api",
		"test_app_id",
		WithHTTPClient(&http.Client{Transport: &mockRoundTripper{extraFn: extraRTFunc}}),
		WithUpdateInterval(1*time.Minute),
	)
	require.NoError(t, err)

	return client
}

func TestHappyPathWithFetchingRates(t *testing.T) {
	rtCalls := 0
	client := setupClientForTesting(t, func() {
		rtCalls += 1
	})

	response, err := client.GetRatesForCurrencies(t.Context(), []string{"USD", "EUR"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedPairs := []models.CurrencyPair{
		{From: "USD", To: "EUR", Rate: 0.86},
		{From: "EUR", To: "USD", Rate: 1.16},
	}

	if !reflect.DeepEqual(response, expectedPairs) {
		t.Errorf("expected %v, got %v", expectedPairs, response)
	}
	assert.Equal(t, rtCalls, 2, "expected RoundTrip to be called twice - one for init and one for fetching rates")
}

func TestHappyPathSyntheticData(t *testing.T) {
	tc := []struct {
		name          string
		currencies    []string
		expectedPairs []models.CurrencyPair
		shouldErr     bool
	}{
		{
			name:       "USD to EUR",
			currencies: []string{"USD", "EUR"},
			expectedPairs: []models.CurrencyPair{
				{From: "USD", To: "EUR", Rate: helpers.RoundToPrecision(0.8, 2)},
				{From: "EUR", To: "USD", Rate: 1.25},
			},
			shouldErr: false,
		},
		{
			name:       "PLN to GBP",
			currencies: []string{"PLN", "GBP"},
			expectedPairs: []models.CurrencyPair{
				{From: "PLN", To: "GBP", Rate: helpers.RoundToPrecision(0.75/3.5, 2)},
				{From: "GBP", To: "PLN", Rate: helpers.RoundToPrecision(3.5/0.75, 2)},
			},
			shouldErr: false,
		},
		{
			name:          "USD to JPY",
			currencies:    []string{"USD", "JPY"},
			expectedPairs: nil,
			shouldErr:     true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			rtCalls := 0
			client := setupClientForTesting(t, func() {
				rtCalls++
			})
			client.latestRates = map[string]float64{
				"USD": 1.00,
				"EUR": 0.8,
				"PLN": 3.5,
				"GBP": 0.75,
			}
			client.latestUpdate = time.Now()
			client.updateInterval = 60 * time.Minute

			response, err := client.GetRatesForCurrencies(t.Context(), tt.currencies)

			if tt.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPairs, response)
			}
			assert.Equal(t, 1, rtCalls, "expected RoundTrip to be called only once during init")
		})
	}
}
