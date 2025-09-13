package main

import (
	"a/internal/adapter/balance/bitcoin"
	httpy "a/internal/adapter/http"
	"a/internal/config"
	"a/internal/ports"
	"log"
	"net/http"
	"os"
	"strings"

	cryptobalancerest "github.com/airgap-solution/crypto-balance-rest/openapi/servergen/go"
	"github.com/restartfu/coinmarketcap/coinmarketcap"
	"github.com/restartfu/gophig"
	"github.com/samber/lo"

	cmcrest "github.com/airgap-solution/cmc-rest/openapi/clientgen/go"
)

func main() {
	configPath, _ := lo.Coalesce(os.Getenv("CONFIG_PATH"), "./config.toml")
	conf, err := loadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	cmcRestCfg := cmcrest.NewConfiguration()
	cmcRestCfg.Scheme = "http"
	cmcRestCfg.Host = conf.CMCRestAddr
	cmcRestClient := cmcrest.NewAPIClient(cmcRestCfg)

	balanceProviders := map[string]ports.BalanceProvider{
		strings.ToLower(coinmarketcap.CurrencyBTC.String()): bitcoin.NewAdapter(conf.ElectrumAddr),
	}

	servicer := httpy.NewAdapter(*cmcRestClient, balanceProviders)
	ctrl := cryptobalancerest.NewDefaultAPIController(servicer)

	router := cryptobalancerest.NewRouter(ctrl)
	err = http.ListenAndServe(conf.ListenAddr, router)
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
