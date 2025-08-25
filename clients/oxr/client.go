package oxr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/blesniewski/knm/helpers"
	"github.com/blesniewski/knm/models"
	"github.com/samber/lo"
)

type Client struct {
	baseURL        string
	appID          string
	httpClient     *http.Client
	latestRates    map[string]float64
	latestUpdate   time.Time
	updateInterval time.Duration
	mutex          sync.Mutex
}

// ^ Would consider RWmutex for a real world scenario where the rates would actually
// be updated before the program finishes

func NewClient(baseURL, appID string) *Client {
	c := &Client{
		baseURL:        baseURL,
		appID:          appID,
		httpClient:     &http.Client{},
		latestRates:    make(map[string]float64),
		updateInterval: 1 * time.Hour,
		mutex:          sync.Mutex{},
	}
	c.getLatestRates()
	return c
}

func (c *Client) WithOptions(options ...Option) *Client {
	for _, opt := range options {
		opt(c)
	}
	return c
}

type Option func(*Client)

func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithUpdateInterval(interval time.Duration) Option {
	return func(c *Client) {
		c.updateInterval = interval
	}
}

func (c *Client) getLatestRates() error {
	requestUrl, err := url.JoinPath(c.baseURL, "/latest.json")
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Set("app_id", c.appID)
	requestUrl = requestUrl + "?" + params.Encode()

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get latest rates: %s", resp.Status)
	}

	var latestResponse LatestResponse
	if err := json.NewDecoder(resp.Body).Decode(&latestResponse); err != nil {
		return err
	}

	c.mutex.Lock()
	c.latestRates = latestResponse.Rates
	c.latestUpdate = time.Unix(latestResponse.Timestamp, 0)
	c.mutex.Unlock()

	return nil
}

func (c *Client) GetRatesForCurrencies(currencies []string) ([]models.CurrencyPair, error) {
	currencies = lo.Uniq(currencies)
	if len(currencies) < 2 {
		return nil, fmt.Errorf("at least two currencies are required")
	}
	for i := range currencies {
		currencies[i] = strings.ToUpper(currencies[i])
	}

	if c.shouldRefreshRates() {
		if err := c.getLatestRates(); err != nil {
			return nil, err
		}
	}

	for _, currency := range currencies {
		if _, ok := c.latestRates[currency]; !ok {
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
	return time.Since(c.latestUpdate) > c.updateInterval
}

func (c *Client) getRate(currency string) (float64, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	rate, ok := c.latestRates[currency]
	return rate, ok
}
