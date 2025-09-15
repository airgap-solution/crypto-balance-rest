package http

import (
	"a/internal/ports"
	"context"
	"strings"

	"github.com/samber/lo"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	cryptobalancerest "github.com/airgap-solution/crypto-balance-rest/openapi/servergen/go"
)

type Adapter struct {
	cmcRestClient    cmcrest.APIClient
	balanceProviders map[string]ports.BalanceProvider
}

func NewAdapter(cmcRestClient cmcrest.APIClient, balanceProviders map[string]ports.BalanceProvider) *Adapter {
	return &Adapter{cmcRestClient: cmcRestClient, balanceProviders: balanceProviders}
}

func (a *Adapter) BalanceGet(ctx context.Context, xpub string, currency string, fiat string) (cryptobalancerest.ImplResponse, error) {
	rate, _, err := a.cmcRestClient.DefaultAPI.V1RateCurrencyFiatGet(ctx, currency, fiat).Execute()
	if err != nil {
		return cryptobalancerest.Response(500, nil), err
	}

	balance, err := a.balanceProviders[strings.ToLower(currency)].Balance(xpub)
	if err != nil {
		return cryptobalancerest.Response(500, nil), err
	}
	return cryptobalancerest.Response(200, cryptobalancerest.BalanceGet200Response{
		Value:     balance * lo.FromPtr(rate.Rate),
		Balance:   balance,
		Change24h: balance * lo.FromPtr(rate.Change24h),
	}), err

}
