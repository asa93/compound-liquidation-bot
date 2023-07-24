package config

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Config represents price feed oracle config
type Config interface {
	RPCURL() *url.URL
	AccountAddress() common.Address
	AccountKey() *ecdsa.PrivateKey
	UpdateInterval() time.Duration
	ContractComptrollerAddress() common.Address
	ContractCusdcAddress() common.Address
}

// FromEnv creates config from environment variables
func FromEnv() (Config, error) {
	rpcURLStr, ok := os.LookupEnv("RPC_URL")
	if !ok {
		return nil, errors.New("RPC_URL: not set")
	}

	rpcURL, err := url.Parse(rpcURLStr)
	if err != nil {
		return nil, fmt.Errorf("RPC_URL: %v", err)
	}

	contractAddressStr, ok := os.LookupEnv("CONTRACT_ADDRESS")
	if !ok {
		return nil, errors.New("CONTRACT_ADDRESS: not set")
	}

	if !common.IsHexAddress(contractAddressStr) {
		return nil, errors.New("CONTRACT_ADDRESS: invalid")
	}

	privateKeyStr, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		return nil, errors.New("PRIVATE_KEY: not set")
	}

	updateIntervalStr, ok := os.LookupEnv("UPDATE_INTERVAL_SECONDS")
	if !ok {
		return nil, errors.New("UPDATE_INTERVAL_SECONDS: not set")
	}

	updateIntervalInt64, err := strconv.ParseInt(updateIntervalStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("UPDATE_INTERVAL_SECONDS: %v", err)
	}

	privateKeyBytes := common.FromHex(privateKeyStr)
	accountKey := crypto.ToECDSAUnsafe(privateKeyBytes)
	accountAddress := crypto.PubkeyToAddress(accountKey.PublicKey)
	contractAddress := common.HexToAddress(contractAddressStr)
	updateInterval := time.Second * time.Duration(updateIntervalInt64)

	contractComptrollerAddressStr, ok := os.LookupEnv("CONTRACT_COMPTROLLER_ADDRESS")
	contractComptrollerAddress := common.HexToAddress(contractComptrollerAddressStr)

	if !ok {
		return nil, errors.New("CONTRACT_COMPTROLLER_ADDRESS: not set")
	}

	contractCusdcAddressStr, ok := os.LookupEnv("CONTRACT_CUSDC_ADDRESS")
	contractCusdcAddress := common.HexToAddress(contractCusdcAddressStr)

	if !ok {
		return nil, errors.New("CONTRACT_CUSDC_ADDRESS: not set")
	}

	return &config{
		rpcURL:                     rpcURL,
		accountAddress:             accountAddress,
		accountKey:                 accountKey,
		contractAddress:            contractAddress,
		updateInverval:             updateInterval,
		contractComptrollerAddress: contractComptrollerAddress,
		contractCusdcAddress:       contractCusdcAddress,
	}, nil
}

type config struct {
	rpcURL                     *url.URL
	contractAddress            common.Address
	accountAddress             common.Address
	accountKey                 *ecdsa.PrivateKey
	updateInverval             time.Duration
	contractComptrollerAddress common.Address
	contractCusdcAddress       common.Address
}

func (c *config) RPCURL() *url.URL {
	return c.rpcURL
}

func (c *config) ContractAddress() common.Address {
	return c.contractAddress
}

func (c *config) AccountAddress() common.Address {
	return c.accountAddress
}

func (c *config) AccountKey() *ecdsa.PrivateKey {
	return c.accountKey
}

func (c *config) UpdateInterval() time.Duration {
	return c.updateInverval
}

func (c *config) setUpdateInterval(duration time.Duration) {
	c.updateInverval = duration
}

func (c *config) ContractComptrollerAddress() common.Address {
	return c.contractComptrollerAddress
}

func (c *config) ContractCusdcAddress() common.Address {
	return c.contractCusdcAddress
}
