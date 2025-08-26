package oxr

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/blesniewski/knm/internal/helpers"
	"github.com/blesniewski/knm/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	response *http.Response
	extraFn  func()
}

func (rt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.extraFn()
	return rt.response, nil
}

func setupClientForTesting(extraRTFunc func()) *Client {
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"base":"USD","rates":{"EUR":0.858611,"USD":1.00}}`)),
	}
	return NewClient("https://openexchangerates.org/api", "test_app_id").
		WithOptions(
			WithHttpClient(&http.Client{Transport: &mockRoundTripper{response: mockResponse, extraFn: extraRTFunc}}),
			WithUpdateInterval(1*time.Minute),
		)
}

func TestHappyPathWithFetchingRates(t *testing.T) {
	rtCalls := 0
	client := setupClientForTesting(func() {
		rtCalls += 1
	})

	response, err := client.GetRatesForCurrencies([]string{"USD", "EUR"})
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
	assert.Equal(t, rtCalls, 1, "expected RoundTrip to be called once")
}

func TestHappyPathSyntheticData(t *testing.T) {
	rtCalls := 0
	client := setupClientForTesting(func() {
		rtCalls += 1
	})
	client.latestRates = map[string]float64{
		"USD": 1.00,
		"EUR": 0.8,
		"PLN": 3.5,
		"GBP": 0.75,
	}
	client.latestUpdate = time.Now()
	client.updateInterval = 60 * time.Minute

	tc := []struct {
		currencies    []string
		expectedPairs []models.CurrencyPair
		shouldErr     bool
	}{
		{
			currencies: []string{"USD", "EUR"},
			expectedPairs: []models.CurrencyPair{
				{From: "USD", To: "EUR", Rate: helpers.RoundToPrecision(0.8, 2)},
				{From: "EUR", To: "USD", Rate: 1.25},
			},
			shouldErr: false,
		},
		{
			currencies: []string{"PLN", "GBP"},
			expectedPairs: []models.CurrencyPair{
				{From: "PLN", To: "GBP", Rate: helpers.RoundToPrecision(0.75/3.5, 2)},
				{From: "GBP", To: "PLN", Rate: helpers.RoundToPrecision(3.5/0.75, 2)},
			},
			shouldErr: false,
		},
		{
			currencies:    []string{"USD", "JPY"},
			expectedPairs: nil,
			shouldErr:     true,
		},
	}

	for _, tc := range tc {
		t.Run(fmt.Sprintf("%v", tc.currencies), func(t *testing.T) {
			response, err := client.GetRatesForCurrencies(tc.currencies)

			if (err != nil) != tc.shouldErr {
				t.Fatalf("expected error: %v, got: %v", tc.shouldErr, err)
			}

			if !reflect.DeepEqual(response, tc.expectedPairs) {
				t.Errorf("expected %v, got %v", tc.expectedPairs, response)
			}
			assert.Equal(t, rtCalls, 0, "expected RoundTrip to be called zero times")
		})
	}
}
