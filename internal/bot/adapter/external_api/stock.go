package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ivofreitas/chat/internal/bot/adapter/csv"
	"net/http"
	"strconv"
)

var (
	NotFoundError            = errors.New("email not found")
	UnprocessableEntityError = errors.New("email unprocessable")
)

type StockAPI interface {
	Lookup(ctx context.Context, stockCode string) ([]Stock, error)
}

type client struct {
	baseURL string
}

func NewClient(baseURL string) StockAPI {
	return &client{baseURL}
}

func (c *client) Lookup(ctx context.Context, stockCode string) ([]Stock, error) {
	url := fmt.Sprintf("%s?s=%s&f=sd2t2ohlcv&h&e=csv", c.baseURL, stockCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, NotFoundError
		case http.StatusUnprocessableEntity:
			return nil, UnprocessableEntityError
		default:
			return nil, fmt.Errorf("received non-200 response status: %s", resp.Status)
		}
	}

	stocks, err := csv.Decode(resp.Body, mapper)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	return stocks, nil
}

func mapper(record []string) (Stock, error) {
	open, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return Stock{}, err
	}
	high, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return Stock{}, err
	}
	low, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return Stock{}, err
	}
	c, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		return Stock{}, err
	}
	volume, err := strconv.ParseFloat(record[7], 64)
	if err != nil {
		return Stock{}, err
	}

	return Stock{
		Symbol: record[0],
		Date:   record[1],
		Time:   record[2],
		Open:   open,
		High:   high,
		Low:    low,
		Close:  c,
		Volume: volume,
	}, nil
}
