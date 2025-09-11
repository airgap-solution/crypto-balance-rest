package bitcoin

import (
	"a/internal/ports"

	"github.com/lamengao/go-electrum/electrum"
)

var _ ports.BalanceProvider = (*Adapter)(nil)

type Adapter struct {
	electrumClient *electrum.Client
}

func NewAdapter(electrumClient *electrum.Client) *Adapter {
	return &Adapter{electrumClient: electrumClient}
}

func (a *Adapter) Balance(address string) (float64, error) {
	return getXpubBalance(a.electrumClient, address, 10, 10)
}
