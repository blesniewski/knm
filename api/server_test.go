package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blesniewski/kryptonim/models"
	"github.com/stretchr/testify/assert"
)

type MockExchangeRateClient struct{}

func (m *MockExchangeRateClient) GetRatesForCurrencies(currencies []string) ([]models.CurrencyPair, error) {
	return []models.CurrencyPair{
		{From: "USD", To: "EUR", Rate: 0.85},
		{From: "EUR", To: "USD", Rate: 1.15},
	}, nil
}

type MockCryptoConversionClient struct{}

func (m *MockCryptoConversionClient) GetConversionRate(from, to string, amount float64) (models.CryptoPair, error) {
	return models.CryptoPair{From: "BTC", To: "USDT", Amount: 45000.00}, nil
}

func SetupForTesting() *Server {
	server := NewServer(&MockExchangeRateClient{}, &MockCryptoConversionClient{})
	return server
}

func TestRates(t *testing.T) {
	server := SetupForTesting()
	tc := []struct {
		currencies     string
		expectedStatus int
	}{
		{"USD,EUR", http.StatusOK},
		{"btc", http.StatusBadRequest},
	}

	for _, tt := range tc {
		t.Run(tt.currencies, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/rates?currencies="+tt.currencies, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestExchange(t *testing.T) {
	server := SetupForTesting()
	tc := []struct {
		from           string
		to             string
		amount         string
		expectedStatus int
	}{
		{"BTC", "USDT", "1", http.StatusOK},
		{"BTC", "USDT", "abc", http.StatusBadRequest},
		{"BTC", "", "1", http.StatusBadRequest},
	}

	for _, tt := range tc {
		t.Run(tt.from+"-"+tt.to, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/exchange?from="+tt.from+"&to="+tt.to+"&amount="+tt.amount, nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
