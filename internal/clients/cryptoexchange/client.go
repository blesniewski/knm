package cryptoexchange

import (
	"fmt"
	"strings"

	"github.com/blesniewski/knm/internal/helpers"
	"github.com/blesniewski/knm/internal/models"
)

type Client struct {
	data map[string]cryptoData
}

type cryptoData struct {
	precision int
	rate      float64
}

func NewClient() *Client {
	// TODO: add fetching data from file ?
	return &Client{
		data: map[string]cryptoData{
			"BEER":  {precision: 18, rate: 0.00002461},
			"FLOKI": {precision: 18, rate: 0.0001428},
			"GATE":  {precision: 18, rate: 6.87},
			"USDT":  {precision: 6, rate: 0.999},
			"WBTC":  {precision: 8, rate: 57037.22},
		},
	}
}

func (c *Client) GetConversionRate(from, to string, amount float64) (models.CryptoPair, error) {
	if amount <= 0 {
		return models.CryptoPair{}, fmt.Errorf("invalid amount: %f", amount)
	}

	fromData, err := c.currency(from)
	if err != nil {
		return models.CryptoPair{}, fmt.Errorf("from currency: %w", err)
	}

	toData, err := c.currency(to)
	if err != nil {
		return models.CryptoPair{}, fmt.Errorf("to currency: %w", err)
	}

	fromAmountUSD := fromData.rate * amount
	resultAmount := fromAmountUSD / toData.rate

	resultAmount = helpers.RoundToPrecision(resultAmount, toData.precision)

	return models.CryptoPair{
		From:   from,
		To:     to,
		Amount: resultAmount,
	}, nil
}

func (c *Client) currency(cur string) (cryptoData, error) {
	curData, ok := c.data[strings.ToUpper(cur)]
	if !ok {
		return cryptoData{}, fmt.Errorf("unknown currency: %s", cur)
	}

	return curData, nil
}
