package coinbasepro

import (
	"context"
	"fmt"
)

// see https://docs.pro.coinbase.com/#the-status-channel for an enumeration of the available fields.
type Currency struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	MinSize       string                 `json:"min_size"`
	Status        string                 `json:"status"`
	StatusMessage string                 `json:"status_message"`
	MaxPrecision  string                 `json:"max_precision"`
	ConvertibleTo []string               `json:"convertible_to"`
	Details       map[string]interface{} `json:"details"`
}

func (c *Client) GetCurrencies(ctx context.Context) ([]Currency, error) {
	var currencies []Currency

	url := fmt.Sprintf("/currencies")
	_, err := c.Request(ctx, "GET", url, nil, &currencies)
	return currencies, err
}
