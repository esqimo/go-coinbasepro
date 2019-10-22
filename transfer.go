package coinbasepro

import (
	"context"
	"fmt"
)

type Transfer struct {
	Type              string `json:"type"`
	Amount            string `json:"amount"`
	CoinbaseAccountID string `json:"coinbase_account_id,string"`
}

func (c *Client) CreateTransfer(ctx context.Context, transfer *Transfer) (Transfer, error) {
	var savedTransfer Transfer

	url := fmt.Sprintf("/transfers")
	_, err := c.Request(ctx, "POST", url, transfer, &savedTransfer)
	return savedTransfer, err
}
