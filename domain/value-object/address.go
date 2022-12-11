package valueobject

import (
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func NewAddress(password string) (string, string) {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, password)
	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	return publicKey.String(), strings.ReplaceAll(mnemonic, " ", "-")
}

func GetAddress(mnemonic, password string) string {
	mnemonic = strings.ReplaceAll(mnemonic, "-", " ")
	seed := bip39.NewSeed(mnemonic, password)
	masterKey, _ := bip32.NewMasterKey(seed)
	publicKey := masterKey.PublicKey()

	return publicKey.String()
}
