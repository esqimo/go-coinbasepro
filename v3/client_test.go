package coinbasepro

import (
	"context"
	"errors"
	"testing"
)

func TestClientErrorsOnNotFound(t *testing.T) {
	client := NewTestClient()
	_, err := client.Request(context.Background(), "GET", "/fake", nil, nil)
	if err == nil {
		t.Error(errors.New("Should have thrown 404 error"))
	}
}
