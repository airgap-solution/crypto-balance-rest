package kaspa

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/kaspanet/go-secp256k1"
	"github.com/kaspanet/kaspad/cmd/kaspawallet/libkaspawallet/bip32"
)

func (a *Adapter) generateAllAddresses(masterKey *bip32.ExtendedKey) ([]string, error) {
	var allAddresses []string

	receiveAddresses, err := a.generateAddresses(masterKey, 0, 10)
	if err != nil {
		return nil, err
	}
	allAddresses = append(allAddresses, receiveAddresses...)

	changeAddresses, err := a.generateAddresses(masterKey, 1, 10)
	if err != nil {
		return nil, err
	}
	allAddresses = append(allAddresses, changeAddresses...)

	return allAddresses, nil
}

func (a *Adapter) generateAddresses(masterKey *bip32.ExtendedKey, purpose uint32, count int) ([]string, error) {
	addresses := make([]string, 0, count)

	purposeKey, err := masterKey.Child(purpose)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		childKey, err := purposeKey.Child(uint32(i))
		if err != nil {
			return nil, err
		}

		publicKey, err := childKey.PublicKey()
		if err != nil {
			return nil, err
		}

		address, err := a.publicKeyToAddress(publicKey)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (a *Adapter) publicKeyToAddress(publicKey *secp256k1.ECDSAPublicKey) (string, error) {
	serializedPubKey, err := publicKey.Serialize()
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(serializedPubKey[:])
	addressBytes := hash[:20]

	return fmt.Sprintf("kaspa:%s", hex.EncodeToString(addressBytes)), nil
}
