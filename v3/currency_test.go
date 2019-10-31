package coinbasepro

import (
	"context"
	"errors"
	"testing"
)

func TestGetCurrencies(t *testing.T) {
	client := NewTestClient()
	currencies, err := client.GetCurrencies(context.Background())
	if err != nil {
		t.Error(err)
	}

	for _, c := range currencies {
		if StructHasZeroValues(c) {
			t.Error(errors.New("Zero value"))
		}
	}
}
