package coinbasepro

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) Subscribe(ctx context.Context, auth bool, msg *Message) (*websocket.Conn, error) {
	var dialer websocket.Dialer
	if c.WsURL == "" {
		return nil, errors.New("no websocket url set")
	}

	var m interface{} = msg

	if auth {
		var err error
		if m, err = c.signMessage(msg); err != nil {
			return nil, err
		}
	}

	if conn, _, err := dialer.DialContext(ctx, c.WsURL, nil); err != nil {
		return nil, err
	} else if err = conn.WriteJSON(m); err != nil {
		conn.Close()
		return nil, err
	} else {
		return conn, nil
	}
}

func generateSig(message, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}

type signedMessage struct {
	Message
	Key        string `json:"key"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Signature  string `json:"signature"`
}

func (c *Client) signMessage(msg *Message) (*signedMessage, error) {
	method := "GET"
	url := "/users/self/verify"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := generateSig(fmt.Sprintf("%s%s%s", timestamp, method, url), c.Secret)
	return &signedMessage{
		Message:    *msg,
		Key:        c.Key,
		Passphrase: c.Passphrase,
		Timestamp:  strconv.FormatInt(time.Now().Unix(), 10),
		Signature:  sig,
	}, err
}
