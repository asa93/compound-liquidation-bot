package priceapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewBlockchainPriceAPI creates new Blockchain.com price API
func NewBlockchainPriceAPI() PriceAPI {
	return &blockchainPriceAPI{
		client: &http.Client{},
	}
}

type blockchainPriceAPI struct {
	client *http.Client
}

func (pa *blockchainPriceAPI) GetName() string {
	return "Blockchain.com"
}

func (pa *blockchainPriceAPI) GetPrice(ctx context.Context) (float64, error) {
	getReq, err := http.NewRequest("GET", blockchainInfoGetPriceEndpoint, nil)
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

	blockchainResp := new(blockchainGetPriceResponse)
	err = json.Unmarshal(respData, &blockchainResp)
	if err != nil {
		return 0, err
	}

	if blockchainResp.USD.Last == 0.0 {
		return 0, errors.New("currency rate is 0")
	}

	return blockchainResp.USD.Last, nil
}

type blockchainGetPriceResponse struct {
	USD struct {
		FifteenM float64 `json:"15m"`
		Last     float64 `json:"last"`
		Buy      float64 `json:"buy"`
		Sell     float64 `json:"sell"`
		Symbol   string  `json:"symbol"`
	} `json:"USD"`
}

const (
	blockchainInfoGetPriceEndpoint = "https://www.blockchain.com/ticker"
)
