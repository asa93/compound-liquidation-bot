package priceapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// NewBitstampPriceAPI creates new Bitstamp price API
func NewBitstampPriceAPI() PriceAPI {
	return &bitstampPriceAPI{
		client: &http.Client{},
	}
}

type bitstampPriceAPI struct {
	client *http.Client
}

func (pa *bitstampPriceAPI) GetName() string {
	return "Bitstamp"
}

func (pa *bitstampPriceAPI) GetPrice(ctx context.Context) (float64, error) {
	getReq, err := http.NewRequest("GET", bitstampGetPriceEndpoint, nil)
	if err != nil {
		return 0, err
	}
	getReq.WithContext(ctx)

	resp, err := pa.client.Do(getReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	bitstampResp := new(bitstampGetPriceResponse)
	err = json.Unmarshal(respData, &bitstampResp)
	if err != nil {
		return 0, err
	}

	priceFloat, err := strconv.ParseFloat(bitstampResp.Last, 64)
	if err != nil {
		return 0, err
	}

	if priceFloat == 0.0 {
		return 0, errors.New("currency rate is 0")
	}

	return priceFloat, nil
}

type bitstampGetPriceResponse struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

const (
	bitstampGetPriceEndpoint = "https://www.bitstamp.net/api/v2/ticker/btcusd/"
)
