package oxr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/blesniewski/knm/internal/helpers"
	"github.com/blesniewski/knm/internal/models"
	"github.com/samber/lo"
)

type Client struct {
	baseURL        string
	appID          string
	httpClient     *http.Client
	updateInterval time.Duration

	// TODO: consider using sync.Map for the latestRates
	// TODO: consider using sync.Atomic for the latestUpdate
	mu           sync.Mutex // <- mutex hat
	latestRates  map[string]float64
	latestUpdate time.Time
}

// ^ Would consider RWmutex for a real world scenario where the rates would actually
// be updated before the program finishes
func NewClient(ctx context.Context, baseURL, appID string, opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:        baseURL,
		appID:          appID,
		httpClient:     http.DefaultClient,
		latestRates:    make(map[string]float64),
		updateInterval: 1 * time.Hour,
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.getLatestRates(ctx); err != nil {
		return nil, fmt.Errorf("failed to get latest rates: %w", err)
	}

	return c, nil
}

type Option func(*Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithUpdateInterval(interval time.Duration) Option {
	return func(c *Client) {
		c.updateInterval = interval
	}
}

func (c *Client) getLatestRates(ctx context.Context) error {
	u, err := url.ParseRequestURI(c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	// TODO: check which url manipulation is safer
	q := u.Query()
	q.Set("app_id", c.appID)
	u.RawQuery = q.Encode()
	u = u.JoinPath("latest.json")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get latest rates: %s", resp.Status)
	}

	var latestResponse LatestResponse
	err = json.NewDecoder(resp.Body).Decode(&latestResponse)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.latestRates = latestResponse.Rates
	c.latestUpdate = time.Unix(latestResponse.Timestamp, 0)

	return nil
}

func (c *Client) GetRatesForCurrencies(ctx context.Context, currencies []string) ([]models.CurrencyPair, error) {
	for i := range currencies {
		currencies[i] = strings.ToUpper(currencies[i])
	}

	if len(lo.Uniq(currencies)) != len(currencies) {
		return nil, fmt.Errorf("duplicate currencies")
	}

	if len(currencies) < 2 {
		return nil, fmt.Errorf("at least two currencies are required")
	}

	if c.shouldRefreshRates() {
		if err := c.getLatestRates(ctx); err != nil {
			return nil, fmt.Errorf("failed to get latest rates: %w", err)
		}
	}

	for _, currency := range currencies {
		if _, ok := c.getRate(currency); !ok {
			return nil, fmt.Errorf("unknown currency: %s", currency)
		}
	}

	var pairs []models.CurrencyPair
	for i := 0; i < len(currencies); i++ {
		for j := 0; j < len(currencies); j++ {
			if i == j {
				continue
			}

			iRate, ok := c.getRate(currencies[i])
			if !ok {
				return nil, fmt.Errorf("unknown currency: %s", currencies[i])
			}

			jRate, ok := c.getRate(currencies[j])
			if !ok {
				return nil, fmt.Errorf("unknown currency: %s", currencies[j])
			}

			rate := jRate / iRate

			rate = helpers.RoundToPrecision(rate, 2)

			pairs = append(pairs, models.CurrencyPair{
				From: currencies[i],
				To:   currencies[j],
				Rate: rate,
			})
		}
	}
	return pairs, nil
}

func (c *Client) shouldRefreshRates() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return time.Since(c.latestUpdate) > c.updateInterval
}

func (c *Client) getRate(currency string) (float64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rate, ok := c.latestRates[currency]
	return rate, ok
}
