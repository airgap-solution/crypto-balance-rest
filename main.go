package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/airgap-solution/go-pkg/mux"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	hd "github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lamengao/go-electrum/electrum"
)

func AddressToScripthash(addr string) (string, error) {
	a, err := btcutil.DecodeAddress(addr, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	script, err := txscript.PayToAddrScript(a)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(script)
	for i, j := 0, len(h)-1; i < j; i, j = i+1, j-1 {
		h[i], h[j] = h[j], h[i]
	}
	return hex.EncodeToString(h[:]), nil
}

func DeriveTaprootAddresses(xpub string, externalCount, changeCount int) (external []btcutil.Address, change []btcutil.Address, err error) {
	key, err := hd.NewKeyFromString(xpub)
	if err != nil {
		return nil, nil, fmt.Errorf("bad xpub: %w", err)
	}

	var extRoot *hd.ExtendedKey
	var chRoot *hd.ExtendedKey

	switch key.Depth() {
	case 3:
		extRoot, err = key.Derive(0)
		if err != nil {
			return nil, nil, fmt.Errorf("derive external chain: %w", err)
		}
		chRoot, err = key.Derive(1)
		if err != nil {
			return nil, nil, fmt.Errorf("derive change chain: %w", err)
		}
	case 4:
		extRoot = key
	default:
		extRoot = key
	}

	makeTaprootAddress := func(pub *btcec.PublicKey) (btcutil.Address, error) {
		outputKey := txscript.ComputeTaprootKeyNoScript(pub)
		ser := schnorr.SerializePubKey(outputKey)
		addr, err := btcutil.NewAddressTaproot(ser, &chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}
		return addr, nil
	}

	if extRoot != nil {
		for i := range externalCount {
			child, err := extRoot.Derive(uint32(i))
			if err != nil {
				return nil, nil, fmt.Errorf("derive external child %d: %w", i, err)
			}
			pub, err := child.ECPubKey()
			if err != nil {
				return nil, nil, err
			}
			addr, err := makeTaprootAddress(pub)
			if err != nil {
				return nil, nil, err
			}
			external = append(external, addr)
		}
	}

	if chRoot != nil {
		for i := range changeCount {
			child, err := chRoot.Derive(uint32(i))
			if err != nil {
				return nil, nil, fmt.Errorf("derive change child %d: %w", i, err)
			}
			pub, err := child.ECPubKey()
			if err != nil {
				return nil, nil, err
			}
			addr, err := makeTaprootAddress(pub)
			if err != nil {
				return nil, nil, err
			}
			change = append(change, addr)
		}
	}

	return external, change, nil
}

func GetXpubBalance(node *electrum.Client, xpub string, externalCount, changeCount int) (btc float64, err error) {
	external, change, err := DeriveTaprootAddresses(xpub, externalCount, changeCount)
	if err != nil {
		return 0, err
	}

	totalSats := int64(0)

	checkAddr := func(addr btcutil.Address) error {
		sh, err := AddressToScripthash(addr.EncodeAddress())
		if err != nil {
			return err
		}

		balResp, err := node.GetBalance(context.Background(), sh)
		if err != nil {
			return err
		}
		confirmed := balResp.Confirmed
		unconfirmed := balResp.Unconfirmed
		totalSats += int64(confirmed) + int64(unconfirmed)
		return nil
	}

	for _, a := range external {
		if err := checkAddr(a); err != nil {
			return 0, err
		}
	}
	for _, a := range change {
		if err := checkAddr(a); err != nil {
			return 0, err
		}
	}

	return float64(totalSats) / 1e8, nil
}

func main() {
	node, err := electrum.NewClientTCP(context.Background(), "electrum.blockstream.info:50001")
	if err != nil {
		log.Fatal(err)
	}

	xpub := "xpub6DUrVe8dGXL1zJ1nFnFF6UTYe69sBboouZBFR1wLtCeLQ9FGoN34n2sFGKvimJ7VcpxVqZJgetm1tcdbQTjXjA57W3miQM8nbWY9Vju9b5a"
	btcBal, err := GetXpubBalance(node, xpub, 10, 10)
	if err != nil {
		log.Fatalf("GetXpubBalance error: %v", err)
	}
	fmt.Printf("Total BTC across derived addresses: %.8f\n", btcBal)

	r := mux.NewRouter(mux.Config{
		Address: "192.168.2.130",
		Port:    8082,
	})

	mux.HandleRoute(r, "/balance", func(q query) (string, mux.Error) {
		fmt.Println(q)
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/crypto/rate?currency=%s&fiat=%s", q.Currency, q.Fiat))
		if err != nil {
			return "", mux.NewError(500)
		}

		buf, _ := io.ReadAll(resp.Body)
		rate, _ := strconv.ParseFloat(string(buf), 64)

		bal, _ := GetXpubBalance(node, q.XPub, 10, 10)
		fmt.Println(bal)
		return fmt.Sprintf(`{"value":%.2f, "balance":%f}`, bal*rate, bal), nil
	})

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}

type query struct {
	XPub     string `query:"xpub"`
	Currency string `query:"currency"`
	Fiat     string `query:"fiat"`
}
