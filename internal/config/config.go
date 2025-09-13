package config

type Config struct {
	ListenAddr   string `toml:"listen_addr"`
	CMCRestAddr  string `toml:"cmc_rest_addr"`
	ElectrumAddr string `toml:"electrum_addr"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddr:   "restartfu.com:8987",
		CMCRestAddr:  "restartfu.com:8765",
		ElectrumAddr: "electrum.blockstream.info:50001",
	}
}
