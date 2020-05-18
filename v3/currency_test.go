package coinbasepro

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetCurrencies(t *testing.T) {
	client := NewTestClient()
	currencies, err := client.GetCurrencies(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, c := range currencies {
		spew.Dump(c)
	}
}
