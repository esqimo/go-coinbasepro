package coinbasepro

import "context"

type Fee struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
	USDVolume    string `json:"usd_volume"`
}

func (c *Client) GetFees(ctx context.Context) (Fee, error) {
	var fees Fee
	_, err := c.Request(ctx, "GET", "/fees", nil, &fees)
	return fees, err
}
