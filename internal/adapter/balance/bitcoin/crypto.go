package bitcoin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	hd "github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lamengao/go-electrum/electrum"
)

func addressToScripthash(addr string) (string, error) {
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

func deriveTaprootAddresses(xpub string, externalCount, changeCount int) (external []btcutil.Address, change []btcutil.Address, err error) {
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

func getXpubBalance(node *electrum.Client, xpub string, externalCount, changeCount int) (btc float64, err error) {
	external, change, err := deriveTaprootAddresses(xpub, externalCount, changeCount)
	if err != nil {
		return 0, err
	}

	totalSats := int64(0)

	checkAddr := func(addr btcutil.Address) error {
		sh, err := addressToScripthash(addr.EncodeAddress())
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
