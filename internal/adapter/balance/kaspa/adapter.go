package kaspa

import (
	"a/internal/ports"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kaspanet/kaspad/cmd/kaspawallet/libkaspawallet/bip32"
)

var _ ports.BalanceProvider = (*Adapter)(nil)

type Adapter struct {
	explorerBaseURL string
}

func NewAdapter(url string) *Adapter {
	return &Adapter{
		explorerBaseURL: url,
	}
}

func (a *Adapter) Balance(kpub string) (float64, error) {
	xpub, err := bip32.DeserializeExtendedKey(kpub)
	if err != nil {
		return 0, err
	}

	allAddresses, err := a.generateAllAddresses(xpub)
	if err != nil {
		return 0, err
	}

	return a.getMultipleAddressBalances(allAddresses)
}

func (a *Adapter) getMultipleAddressBalances(addresses []string) (float64, error) {
	requestBody := map[string]interface{}{
		"addresses": addresses,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(a.explorerBaseURL+"addresses/balances", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var response struct {
		Address string  `json:"address"`
		Balance float64 `json:"balance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.Balance, nil
}
