package priceapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewCoingeckoPriceAPI creates new Coindesk price API
func NewCoingeckoPriceAPI() PriceAPI {
	return &coingeckoPriceAPI{
		client: &http.Client{},
	}
}

type coingeckoPriceAPI struct {
	client *http.Client
}

func (pa *coingeckoPriceAPI) GetName() string {
	return "CoinGecko"
}

func (pa *coingeckoPriceAPI) GetPrice(ctx context.Context) (float64, error) {
	getReq, err := http.NewRequest("GET", coingeckoGetPriceEndpoint, nil)
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

	coingeckoResps := new(coingeckoGetPriceResponse)
	err = json.Unmarshal(respData, &coingeckoResps)
	if err != nil {
		return 0, err
	}

	if len(*coingeckoResps) != 1 {
		return 0, fmt.Errorf("unexpected response array length: %d", len(*coingeckoResps))
	}

	coingeckoResp := (*coingeckoResps)[0]
	if coingeckoResp.CurrentPrice == 0.0 {
		return 0, errors.New("currency rate is 0")
	}

	return coingeckoResp.CurrentPrice, nil
}

type coingeckoGetPriceResponse []struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	CurrentPrice float64 `json:"current_price"`
}

const (
	coingeckoGetPriceEndpoint = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=bitcoin"
)
