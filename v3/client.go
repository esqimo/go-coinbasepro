package coinbasepro

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	WsURL      string
	BaseURL    string
	Secret     string
	Key        string
	Passphrase string
	HTTPClient *http.Client
	RetryCount int
}

type ClientConfig struct {
	WsURL      string
	BaseURL    string
	Key        string
	Passphrase string
	Secret     string
}

func NewClient() *Client {
	const defaultBaseUrl = "https://api.pro.coinbase.com"
	baseURL := os.Getenv("COINBASE_PRO_BASEURL")
	if baseURL == "" {
		baseURL = defaultBaseUrl
	}
	wsURL := os.Getenv("COINBASE_PRO_WSURL")
	// Only set the websocket url to a default i
	if wsURL == "" && baseURL == defaultBaseUrl {
		wsURL = "wss://ws-feed.pro.coinbase.com"
	}

	client := Client{
		BaseURL:    baseURL,
		WsURL:      wsURL,
		Key:        os.Getenv("COINBASE_PRO_KEY"),
		Passphrase: os.Getenv("COINBASE_PRO_PASSPHRASE"),
		Secret:     os.Getenv("COINBASE_PRO_SECRET"),
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		RetryCount: 0,
	}

	return &client
}

func (c *Client) Request(
	ctx context.Context, method string, url string, params, result interface{}) (res *http.Response, err error) {
	for i := 0; i < c.RetryCount+1; i++ {
		retryDuration := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond
		time.Sleep(retryDuration)

		res, err = c.request(ctx, method, url, params, result)
		if res != nil && res.StatusCode == 429 {
			continue
		} else {
			break
		}
	}

	return res, err
}

func (c *Client) request(ctx context.Context, method string, url string,
	params, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, url)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// XXX: Sandbox time is off right now
	if os.Getenv("TEST_COINBASE_OFFSET") != "" {
		inc, err := strconv.Atoi(os.Getenv("TEST_COINBASE_OFFSET"))
		if err != nil {
			return res, err
		}

		timestamp = strconv.FormatInt(time.Now().Unix()+int64(inc), 10)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Go Coinbase Pro Client 1.0")

	h, err := c.Headers(method, url, timestamp, string(data))
	if err != nil {
		return res, err
	}

	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		coinbaseError := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&coinbaseError); err != nil {
			return res, err
		}
		coinbaseError.StatusCode = res.StatusCode
		return res, error(coinbaseError)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(result); err != nil {
			return res, err
		}
	}

	return res, nil
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *Client) Headers(method, url, timestamp, data string) (map[string]string, error) {
	h := make(map[string]string)
	h["CB-ACCESS-KEY"] = c.Key
	h["CB-ACCESS-PASSPHRASE"] = c.Passphrase
	h["CB-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, c.Secret)
	if err != nil {
		return nil, err
	}
	h["CB-ACCESS-SIGN"] = sig
	return h, nil
}
