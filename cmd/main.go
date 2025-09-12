package main

import (
	"a/internal/adapter/balance/bitcoin"
	httpy "a/internal/adapter/http"
	"a/internal/ports"
	"net/http"
	"strings"

	cryptobalancerest "github.com/airgap-solution/crypto-balance-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
)

func main() {
	cfg := cmcrest.NewConfiguration()
	cfg.Host = "localhost:8083"
	cfg.Scheme = "http"
	cmcRestClient := cmcrest.NewAPIClient(cfg)

	balanceProviders := map[string]ports.BalanceProvider{
		strings.ToLower(coinmarketcap.CurrencyBTC.String()): bitcoin.NewAdapter("restartfu.com:5001"),
	}

	servicer := httpy.NewAdapter(*cmcRestClient, balanceProviders)
	ctrl := cryptobalancerest.NewDefaultAPIController(servicer)

	router := cryptobalancerest.NewRouter(ctrl)
	http.ListenAndServe("localhost:8082", router)
}
