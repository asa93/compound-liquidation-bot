package priceapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewCoindeskPriceAPI creates new Coindesk price API
func NewCoindeskPriceAPI() PriceAPI {
	return &coindeskPriceAPI{
		client: &http.Client{},
	}
}

type coindeskPriceAPI struct {
	client *http.Client
}

func (pa *coindeskPriceAPI) GetName() string {
	return "CoinDesk"
}

func (pa *coindeskPriceAPI) GetPrice(ctx context.Context) (float64, error) {
	getReq, err := http.NewRequest("GET", coindeskGetPriceEndpoint, nil)
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

	coindeskResp := new(coindeskGetPriceResponse)
	err = json.Unmarshal(respData, &coindeskResp)
	if err != nil {
		return 0, err
	}

	if coindeskResp.Bpi.USD.RateFloat == 0.0 {
		return 0, errors.New("currency rate is 0")
	}

	return coindeskResp.Bpi.USD.RateFloat, nil
}

const (
	coindeskGetPriceEndpoint = "https://api.coindesk.com/v1/bpi/currentprice/USD.json"
)

type coindeskGetPriceResponse struct {
	Bpi struct {
		USD struct {
			Code        string  `json:"code"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
	} `json:"bpi"`
}
