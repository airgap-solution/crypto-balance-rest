package bitcoin

import (
	"a/internal/ports"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lamengao/go-electrum/electrum"
)

var _ ports.BalanceProvider = (*Adapter)(nil)

type Adapter struct {
	mu             sync.RWMutex
	electrumClient *electrum.Client
	addresses      map[string][]btcutil.Address
	addr           string
}

// NewAdapter creates an adapter with auto-reconnect
func NewAdapter(addr string) *Adapter {
	a := &Adapter{
		addresses: make(map[string][]btcutil.Address),
		addr:      addr,
	}
	a.connectWithRetry()
	return a
}

// connectWithRetry keeps trying until success
func (a *Adapter) connectWithRetry() {
	for {
		client, err := electrum.NewClientTCP(context.Background(), a.addr)
		if err == nil {
			log.Printf("[bitcoin] connected to Electrum %s", a.addr)
			a.mu.Lock()
			a.electrumClient = client
			a.mu.Unlock()
			return
		}
		log.Printf("[bitcoin] electrum connection failed: %v, retrying in 5s...", err)
		time.Sleep(5 * time.Second)
	}
}

func (a *Adapter) getClient() *electrum.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.electrumClient
}

func (a *Adapter) Balance(xpub string) (float64, error) {
	addresses, ok := a.addresses[xpub]
	if !ok {
		external, change, err := deriveTaprootAddresses(xpub, 1000, 1000)
		if err != nil {
			return 0, err
		}
		addresses = append(external, change...)
		a.addresses[xpub] = addresses
	}

	var lastErr error
	for i := range 3 {
		client := a.getClient()
		if client.IsShutdown() {
			a.connectWithRetry()
			continue
		}

		bal, err := getXpubBalance(client, addresses)
		if err == nil {
			return bal, nil
		}

		lastErr = err
		log.Printf("[bitcoin] balance fetch failed (attempt %d): %v", i+1, err)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			a.connectWithRetry()
		}

		time.Sleep(time.Second * 2)
	}

	return 0, lastErr
}
