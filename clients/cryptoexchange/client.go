package cryptoexchange

import (
	"fmt"
	"strings"

	"github.com/blesniewski/knm/helpers"
	"github.com/blesniewski/knm/models"
)

type Client struct {
	data map[string]cryptoData
}

type cryptoData struct {
	precision int
	rate      float64
}

func NewClient() *Client {
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
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	fromData, ok := c.data[from]
	if !ok {
		return models.CryptoPair{}, fmt.Errorf("unknown currency: %s", from)
	}

	toData, ok := c.data[to]
	if !ok {
		return models.CryptoPair{}, fmt.Errorf("unknown currency: %s", to)
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
