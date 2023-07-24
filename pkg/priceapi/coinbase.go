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

// NewCoinbasePriceAPI creates new Coinbase price API
func NewCoinbasePriceAPI() PriceAPI {
	return &coinbasePriceAPI{
		client: &http.Client{},
	}
}

type coinbasePriceAPI struct {
	client *http.Client
}

func (pa *coinbasePriceAPI) GetName() string {
	return "Coinbase"
}

func (pa *coinbasePriceAPI) GetPrice(ctx context.Context) (float64, error) {
	getReq, err := http.NewRequest("GET", coinbaseGetPriceEndpoint, nil)
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

	coinbaseResp := new(coinbaseGetPriceResponse)
	err = json.Unmarshal(respData, &coinbaseResp)
	if err != nil {
		return 0, err
	}

	priceFloat, err := strconv.ParseFloat(coinbaseResp.Data.Amount, 64)
	if err != nil {
		return 0, err
	}

	if priceFloat == 0.0 {
		return 0, errors.New("currency rate is 0")
	}

	return priceFloat, nil
}

type coinbaseGetPriceResponse struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}

const (
	coinbaseGetPriceEndpoint = "https://api.coinbase.com/v2/prices/spot?currency=USD"
)
