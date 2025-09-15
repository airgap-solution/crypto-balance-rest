package main

import (
	"a/internal/adapter/balance/bitcoin"
	"a/internal/adapter/balance/kaspa"
	httpy "a/internal/adapter/http"
	"a/internal/config"
	"a/internal/ports"
	"log"
	"net/http"
	"os"
	"strings"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
	cryptobalancerest "github.com/airgap-solution/crypto-balance-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/restartfu/gophig"
	"github.com/samber/lo"
)

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	configPath, _ := lo.Coalesce(os.Getenv("CONFIG_PATH"), "./config.toml")
	conf, err := loadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmcRestCfg := cmcrest.NewConfiguration()
	cmcRestCfg.Scheme = "http"
	if conf.TLSEnabled {
		cmcRestCfg.Scheme += "s"
	}

	cmcRestCfg.Host = conf.CMCRestAddr
	cmcRestClient := cmcrest.NewAPIClient(cmcRestCfg)

	balanceProviders := map[string]ports.BalanceProvider{
		strings.ToLower(coinmarketcap.CurrencyBTC.String()): bitcoin.NewAdapter(conf.ElectrumAddr),
		strings.ToLower(coinmarketcap.CurrencyKAS.String()): kaspa.NewAdapter("https://api.kaspa.org/addresses/balances"),
	}

	servicer := httpy.NewAdapter(*cmcRestClient, balanceProviders)
	ctrl := cryptobalancerest.NewDefaultAPIController(servicer)
	router := cryptobalancerest.NewRouter(ctrl)

	// Wrap the router with CORS middleware
	handler := corsMiddleware(router)

	if conf.TLSEnabled {
		err = http.ListenAndServeTLS(conf.ListenAddr, conf.TLSConfig.CertificatePath, conf.TLSConfig.PrivateKeyPath, handler)
	} else {
		err = http.ListenAndServe(conf.ListenAddr, handler)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(configPath string) (config.Config, error) {
	defaultConfig := config.DefaultConfig()
	g := gophig.NewGophig[config.Config](configPath, gophig.TOMLMarshaler{}, 0777)
	conf, err := g.LoadConf()
	if err != nil {
		if os.IsNotExist(err) {
			err = g.SaveConf(defaultConfig)
			return defaultConfig, err
		}
		return config.Config{}, err
	}
	return conf, nil
}
