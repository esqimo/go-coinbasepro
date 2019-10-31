package coinbasepro

import (
	"context"
	"fmt"
)

type WithdrawalCrypto struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	CryptoAddress string `json:"crypto_address"`
}

type WithdrawalCoinbase struct {
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	CoinbaseAccountID string `json:"coinbase_account_id"`
}

func (c *Client) CreateWithdrawalCrypto(ctx context.Context, withdrawal *WithdrawalCrypto) (WithdrawalCrypto, error) {
	var savedWithdrawal WithdrawalCrypto
	url := fmt.Sprintf("/withdrawals/crypto")
	_, err := c.Request(ctx, "POST", url, withdrawal, &savedWithdrawal)
	return savedWithdrawal, err
}

func (c *Client) CreateWithdrawalCoinbase(ctx context.Context, withdrawal *WithdrawalCoinbase) (WithdrawalCoinbase, error) {
	var savedWithdrawal WithdrawalCoinbase
	url := fmt.Sprintf("/withdrawals/coinbase-account")
	_, err := c.Request(ctx, "POST", url, withdrawal, &savedWithdrawal)
	return savedWithdrawal, err
}
