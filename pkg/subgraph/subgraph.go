package subgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// NewBitstampPriceAPI creates new Bitstamp price API
func NewSubgraph() *subgraph {
	return &subgraph{
		client: &http.Client{},
	}
}

type subgraph struct {
	client *http.Client
}

func (s *subgraph) GetEndpoint() string {
	return subgraphEndpoint
}

func (s *subgraph) GetAccounts(ctx context.Context) ([]subgraphAccount, error) {

	respData, err := postQuery(ctx, s.client, accountQuery)
	if err != nil {
		return nil, err
	}

	response := new(subgraphResponse)
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}

	accounts := response.Data.Accounts
	return accounts, nil
}

func postQuery(ctx context.Context, client *http.Client, query string) ([]byte, error) {
	payload := []byte(query)
	req, err := http.NewRequest("POST", subgraphEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}

type subgraphResponse struct {
	Data subgraphData `json:"data"`
}
type subgraphData struct {
	Accounts []subgraphAccount `json:"accounts"`
}
type subgraphAccount struct {
	Id                    string `json:"id"`
	TotalBorrowValueInEth string `json:"totalBorrowValueInEth"`
	Health                string `json:"health"`
}

func (a *subgraphAccount) IsLiquidable() bool {

	totalBorrowValueInEth, err := strconv.ParseFloat(a.TotalBorrowValueInEth, 64)

	if err != nil {
		fmt.Println("Error parsing float:", err)

	}

	health, err := strconv.ParseFloat(a.Health, 64)

	if err != nil {
		fmt.Println("Error parsing float:", err)

	}

	return totalBorrowValueInEth > 0 && health > 0 //&& health < 1 // to uncomment

}

const subgraphEndpoint = "https://api.thegraph.com/subgraphs/name/graphprotocol/compound-v2"

const accountQuery = `
{"query":"{\n  accounts(first: 1, orderBy: id, where : {hasBorrowed: true}) {\n\tid \n  totalBorrowValueInEth \n  health \n tokens {id }\n  }\n}","variables":{}}
`
